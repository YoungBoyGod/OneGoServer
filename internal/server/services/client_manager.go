package services

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"learngo0619/internal/logger"
	"learngo0619/internal/server/models"

	"go.uber.org/zap"
)

const (
	TokenFilePath    = "register_token.json"
	ClientFilePath   = "clients.json"
	ClientStateFile  = "client_states.json" // 客户端状态持久化文件
	TokenValidDays   = 30
	OfflineThreshold = 15 * time.Minute // 修改：15分钟无心跳认为离线（之前是10分钟）
	CleanupThreshold = 2 * time.Hour    // 修改：2小时后才清理离线客户端（之前是10分钟）
)

var (
	registerToken *models.RegisterToken
	clientList    *models.ClientList
	once          sync.Once
)

// ClientState 客户端状态持久化结构
type ClientState struct {
	ClientID          string                 `json:"client_id"`
	Name              string                 `json:"name"`
	Type              models.ClientType      `json:"type"`
	IPAddress         string                 `json:"ip_address"`
	LastHeartbeat     time.Time              `json:"last_heartbeat"`
	RegisterTime      time.Time              `json:"register_time"`
	Metadata          map[string]interface{} `json:"metadata"`
	HeartbeatInterval time.Duration          `json:"heartbeat_interval"`
}

// 生成随机token
func generateToken(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// 初始化注册token（首次生成或过期自动更新）
func InitRegisterToken() error {
	once.Do(func() {
		f := models.TokenFile(TokenFilePath)
		t, err := f.Load()
		if err != nil || t == nil || t.ExpiresAt.Before(time.Now()) {
			tokenStr, _ := generateToken(32)
			now := time.Now()
			t = &models.RegisterToken{
				Token:     tokenStr,
				CreatedAt: now,
				ExpiresAt: now.Add(TokenValidDays * 24 * time.Hour),
			}
			_ = f.Save(t)
		}
		registerToken = t
	})
	return nil
}

func GetRegisterToken() *models.RegisterToken {
	return registerToken
}

// 校验token合法性
func ValidateToken(token string) bool {
	return registerToken != nil && registerToken.Token == token && registerToken.ExpiresAt.After(time.Now())
}

// 初始化client列表
func InitClientList() error {
	f := models.ClientFile(ClientFilePath)
	list, err := f.Load()
	if err != nil {
		return err
	}
	clientList = list
	return nil
}

func SaveClientList() error {
	f := models.ClientFile(ClientFilePath)
	return f.Save(clientList)
}

// 注册client，绑定ip，返回client_id
func RegisterClient(ip string) (string, bool) {
	clientList.Lock()
	defer clientList.Unlock()
	for _, c := range clientList.Clients {
		if c.IP == ip {
			return c.ClientID, false // 已注册
		}
	}
	clientID, _ := generateToken(16)
	ci := models.ClientInfo{
		ClientID:     clientID,
		IP:           ip,
		RegisteredAt: time.Now(),
	}
	clientList.Clients = append(clientList.Clients, ci)
	_ = SaveClientList()
	return clientID, true
}

// 获取所有已注册client
func GetAllClients() []models.ClientInfo {
	clientList.Lock()
	defer clientList.Unlock()
	return append([]models.ClientInfo{}, clientList.Clients...)
}

// ClientManager 客户端管理器（增强版）
type ClientManager struct {
	clients           map[string]*models.RegisteredClient // 客户端映射
	mutex             sync.RWMutex                        // 读写锁
	heartbeatInterval time.Duration                       // 心跳间隔
	cleanupInterval   time.Duration                       // 清理间隔
	stopCleanup       chan struct{}                       // 停止清理信号
	stateFile         string                              // 状态文件路径
}

// NewClientManager 创建新的客户端管理器（增强版）
func NewClientManager() *ClientManager {
	cm := &ClientManager{
		clients:           make(map[string]*models.RegisteredClient),
		heartbeatInterval: 30 * time.Second, // 默认30秒心跳
		cleanupInterval:   5 * time.Minute,  // 每5分钟清理一次（增加频率）
		stopCleanup:       make(chan struct{}),
		stateFile:         ClientStateFile,
	}

	// 尝试恢复客户端状态
	cm.recoverClientStates()

	// 启动清理协程
	go cm.startCleanupRoutine()

	logger.Info("Enhanced client manager initialized",
		zap.String("state_file", cm.stateFile),
		zap.Duration("offline_threshold", OfflineThreshold),
		zap.Duration("cleanup_threshold", CleanupThreshold),
	)

	return cm
}

// recoverClientStates 恢复客户端状态
func (cm *ClientManager) recoverClientStates() {
	if _, err := os.Stat(cm.stateFile); os.IsNotExist(err) {
		logger.Info("No client state file found, starting fresh")
		return
	}

	data, err := os.ReadFile(cm.stateFile)
	if err != nil {
		logger.Error("Failed to read client state file", zap.Error(err))
		return
	}

	var states []ClientState
	if err := json.Unmarshal(data, &states); err != nil {
		logger.Error("Failed to unmarshal client states", zap.Error(err))
		return
	}

	recoveredCount := 0
	now := time.Now()

	for _, state := range states {
		// 只恢复最近活跃的客户端（24小时内有心跳的）
		if now.Sub(state.LastHeartbeat) < 24*time.Hour {
			client := &models.RegisteredClient{
				ClientID:       state.ClientID,
				Name:           state.Name,
				Type:           state.Type,
				IPAddress:      state.IPAddress,
				Status:         models.ClientStatusOffline, // 初始为离线状态
				LastHeartbeat:  state.LastHeartbeat,
				RegisterTime:   state.RegisterTime,
				LastActiveTime: state.LastHeartbeat,
				Metadata:       state.Metadata,
				RequestCount:   0,
				ErrorCount:     0,
				OnlineDuration: 0,
			}

			cm.clients[state.ClientID] = client
			recoveredCount++
		}
	}

	logger.Info("Client states recovered",
		zap.Int("total_states", len(states)),
		zap.Int("recovered_count", recoveredCount),
	)
}

// saveClientStates 保存客户端状态
func (cm *ClientManager) saveClientStates() {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	var states []ClientState
	for _, client := range cm.clients {
		state := ClientState{
			ClientID:          client.ClientID,
			Name:              client.Name,
			Type:              client.Type,
			IPAddress:         client.IPAddress,
			LastHeartbeat:     client.LastHeartbeat,
			RegisterTime:      client.RegisterTime,
			Metadata:          client.Metadata,
			HeartbeatInterval: cm.heartbeatInterval,
		}
		states = append(states, state)
	}

	data, err := json.MarshalIndent(states, "", "  ")
	if err != nil {
		logger.Error("Failed to marshal client states", zap.Error(err))
		return
	}

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(cm.stateFile), 0755); err != nil {
		logger.Error("Failed to create state directory", zap.Error(err))
		return
	}

	if err := os.WriteFile(cm.stateFile, data, 0644); err != nil {
		logger.Error("Failed to write client state file", zap.Error(err))
		return
	}

	logger.Debug("Client states saved", zap.Int("count", len(states)))
}

// RegisterClient 注册新客户端（增强版）
func (cm *ClientManager) RegisterClient(req *models.ClientRegisterRequest, ip, userAgent string) (*models.ClientRegisterResponse, error) {
	// 验证请求
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("注册请求验证失败: %w", err)
	}

	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// 生成客户端ID
	clientID := cm.generateClientID(ip, userAgent, req.Name)

	// 生成客户端指纹
	fingerprint := cm.generateFingerprint(ip, userAgent)

	now := time.Now()

	// 检查是否已经注册
	if existingClient, exists := cm.clients[clientID]; exists {
		// 更新现有客户端信息
		existingClient.Name = req.Name
		existingClient.Type = req.Type
		existingClient.Version = req.Version
		existingClient.Description = req.Description
		existingClient.Status = models.ClientStatusOnline
		existingClient.LastHeartbeat = now
		existingClient.LastActiveTime = now
		existingClient.Metadata = req.Metadata
		existingClient.Tags = req.Tags
		existingClient.IPAddress = ip // 更新IP（可能发生变化）
		existingClient.UserAgent = userAgent

		logger.Info("Client re-registered",
			zap.String("client_id", clientID),
			zap.String("name", req.Name),
			zap.String("ip", ip),
			zap.Duration("offline_duration", now.Sub(existingClient.LastHeartbeat)),
		)

		// 保存状态
		go cm.saveClientStates()

		return &models.ClientRegisterResponse{
			ClientID:          clientID,
			Success:           true,
			Message:           "客户端重新注册成功",
			RegisterTime:      existingClient.RegisterTime,
			HeartbeatUrl:      "/api/v1/clients/heartbeat",
			HeartbeatInterval: int(cm.heartbeatInterval.Seconds()),
		}, nil
	}

	// 创建新客户端
	client := &models.RegisteredClient{
		ClientID:       clientID,
		Name:           req.Name,
		Type:           req.Type,
		Version:        req.Version,
		Description:    req.Description,
		IPAddress:      ip,
		UserAgent:      userAgent,
		Fingerprint:    fingerprint,
		Status:         models.ClientStatusOnline,
		LastHeartbeat:  now,
		RegisterTime:   now,
		LastActiveTime: now,
		RequestCount:   0,
		ErrorCount:     0,
		OnlineDuration: 0,
		Metadata:       req.Metadata,
		Tags:           req.Tags,
	}

	// 保存客户端
	cm.clients[clientID] = client

	// 保存状态
	go cm.saveClientStates()

	// 记录注册日志
	logger.Info("New client registered",
		zap.String("client_id", clientID),
		zap.String("name", req.Name),
		zap.String("type", string(req.Type)),
		zap.String("version", req.Version),
		zap.String("ip", ip),
		zap.Int("total_clients", len(cm.clients)),
	)

	return &models.ClientRegisterResponse{
		ClientID:          clientID,
		Success:           true,
		Message:           "客户端注册成功",
		RegisterTime:      now,
		HeartbeatUrl:      "/api/v1/clients/heartbeat",
		HeartbeatInterval: int(cm.heartbeatInterval.Seconds()),
	}, nil
}

// ProcessHeartbeat 处理客户端心跳（增强版）
func (cm *ClientManager) ProcessHeartbeat(req *models.ClientHeartbeatRequest) (*models.ClientHeartbeatResponse, error) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	client, exists := cm.clients[req.ClientID]
	if !exists {
		return &models.ClientHeartbeatResponse{
			Success:       false,
			Message:       "客户端未注册，请重新注册",
			ServerTime:    time.Now(),
			NextHeartbeat: int(cm.heartbeatInterval.Seconds()),
		}, fmt.Errorf("client not found: %s", req.ClientID)
	}

	// 更新心跳信息
	client.UpdateHeartbeat(req.Status, req.Metadata)

	// 如果客户端之前离线，现在重新上线，记录日志
	if client.Status != models.ClientStatusOnline {
		client.Status = models.ClientStatusOnline
		logger.Info("Client back online",
			zap.String("client_id", req.ClientID),
			zap.String("client_name", client.Name),
		)
	}

	// 定期保存状态（每10次心跳保存一次，减少IO）
	client.RequestCount++
	if client.RequestCount%10 == 0 {
		go cm.saveClientStates()
	}

	// 记录心跳日志
	logger.Debug("Client heartbeat received",
		zap.String("client_id", req.ClientID),
		zap.String("status", string(req.Status)),
		zap.String("client_name", client.Name),
	)

	return &models.ClientHeartbeatResponse{
		Success:       true,
		Message:       "心跳接收成功",
		ServerTime:    time.Now(),
		NextHeartbeat: int(cm.heartbeatInterval.Seconds()),
	}, nil
}

// GetClient 获取指定客户端
func (cm *ClientManager) GetClient(clientID string) (*models.RegisteredClient, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	client, exists := cm.clients[clientID]
	return client, exists
}

// GetAllClients 获取所有客户端
func (cm *ClientManager) GetAllClients() []*models.RegisteredClient {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	clients := make([]*models.RegisteredClient, 0, len(cm.clients))
	for _, client := range cm.clients {
		// 更新在线时长
		client.UpdateOnlineDuration()
		clients = append(clients, client)
	}

	return clients
}

// GetClientsList 获取客户端列表响应
func (cm *ClientManager) GetClientsList() *models.ClientListResponse {
	clients := cm.GetAllClients()

	onlineCount := 0
	offlineCount := 0

	for _, client := range clients {
		if client.IsOnline() {
			onlineCount++
		} else {
			offlineCount++
		}
	}

	return &models.ClientListResponse{
		Total:        len(clients),
		OnlineCount:  onlineCount,
		OfflineCount: offlineCount,
		Clients:      clients,
		Timestamp:    time.Now(),
	}
}

// UnregisterClient 注销客户端
func (cm *ClientManager) UnregisterClient(clientID string) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	client, exists := cm.clients[clientID]
	if !exists {
		return fmt.Errorf("client not found: %s", clientID)
	}

	// 记录注销日志
	logger.Info("Client unregistered",
		zap.String("client_id", clientID),
		zap.String("name", client.Name),
		zap.String("ip", client.IPAddress),
		zap.Duration("online_duration", client.OnlineDuration),
	)

	// 删除客户端
	delete(cm.clients, clientID)

	return nil
}

// IncrementRequestCount 增加客户端请求计数
func (cm *ClientManager) IncrementRequestCount(ip, userAgent string) {
	clientID := cm.generateClientID(ip, userAgent, "")

	cm.mutex.RLock()
	client, exists := cm.clients[clientID]
	cm.mutex.RUnlock()

	if exists {
		client.IncrementRequestCount()
	}
}

// IncrementErrorCount 增加客户端错误计数
func (cm *ClientManager) IncrementErrorCount(ip, userAgent string) {
	clientID := cm.generateClientID(ip, userAgent, "")

	cm.mutex.RLock()
	client, exists := cm.clients[clientID]
	cm.mutex.RUnlock()

	if exists {
		client.IncrementErrorCount()
	}
}

// GetStats 获取统计信息
func (cm *ClientManager) GetStats() map[string]interface{} {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	onlineCount := 0
	offlineCount := 0
	totalRequests := int64(0)
	totalErrors := int64(0)

	typeStats := make(map[models.ClientType]int)

	for _, client := range cm.clients {
		if client.IsOnline() {
			onlineCount++
		} else {
			offlineCount++
		}

		totalRequests += client.RequestCount
		totalErrors += client.ErrorCount
		typeStats[client.Type]++
	}

	return map[string]interface{}{
		"total_clients":      len(cm.clients),
		"online_clients":     onlineCount,
		"offline_clients":    offlineCount,
		"total_requests":     totalRequests,
		"total_errors":       totalErrors,
		"error_rate":         calculateErrorRate(totalRequests, totalErrors),
		"types":              typeStats,
		"heartbeat_interval": cm.heartbeatInterval.Seconds(),
	}
}

// startCleanupRoutine 启动清理协程
func (cm *ClientManager) startCleanupRoutine() {
	ticker := time.NewTicker(cm.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cm.cleanupOfflineClients()
		case <-cm.stopCleanup:
			return
		}
	}
}

// cleanupOfflineClients 清理离线客户端（增强版）
func (cm *ClientManager) cleanupOfflineClients() {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	now := time.Now()
	var toDelete []string
	var toMarkOffline []string

	for clientID, client := range cm.clients {
		offlineDuration := now.Sub(client.LastHeartbeat)

		// 超过清理阈值的客户端删除
		if offlineDuration > CleanupThreshold {
			toDelete = append(toDelete, clientID)
		} else if offlineDuration > OfflineThreshold && client.Status == models.ClientStatusOnline {
			// 超过离线阈值但未到清理阈值的客户端标记为离线
			toMarkOffline = append(toMarkOffline, clientID)
		}
	}

	// 标记离线客户端
	for _, clientID := range toMarkOffline {
		client := cm.clients[clientID]
		client.Status = models.ClientStatusOffline
		logger.Info("Client marked as offline",
			zap.String("client_id", clientID),
			zap.String("name", client.Name),
			zap.Duration("offline_duration", now.Sub(client.LastHeartbeat)),
		)
	}

	// 删除长期离线的客户端
	for _, clientID := range toDelete {
		client := cm.clients[clientID]
		logger.Info("Cleaned up long-term offline client",
			zap.String("client_id", clientID),
			zap.String("name", client.Name),
			zap.Duration("offline_duration", now.Sub(client.LastHeartbeat)),
		)
		delete(cm.clients, clientID)
	}

	// 如果有变化，保存状态
	if len(toMarkOffline) > 0 || len(toDelete) > 0 {
		go cm.saveClientStates()

		logger.Info("Client cleanup completed",
			zap.Int("marked_offline", len(toMarkOffline)),
			zap.Int("cleaned_count", len(toDelete)),
			zap.Int("remaining_clients", len(cm.clients)),
		)
	}
}

// generateClientID 生成客户端ID
func (cm *ClientManager) generateClientID(ip, userAgent, name string) string {
	// 使用IP + UserAgent + Name生成ID，保证唯一性
	data := fmt.Sprintf("%s|%s|%s", ip, userAgent, name)
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", hash)[:16]
}

// generateFingerprint 生成客户端指纹
func (cm *ClientManager) generateFingerprint(ip, userAgent string) string {
	data := fmt.Sprintf("%s|%s|%d", ip, userAgent, time.Now().Unix())
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", hash)[:12]
}

// Stop 停止客户端管理器
func (cm *ClientManager) Stop() {
	close(cm.stopCleanup)
	logger.Info("Client manager stopped")
}

// calculateErrorRate 计算错误率
func calculateErrorRate(totalRequests, totalErrors int64) float64 {
	if totalRequests == 0 {
		return 0.0
	}
	return float64(totalErrors) / float64(totalRequests) * 100
}

// SetHeartbeatInterval 设置心跳间隔
func (cm *ClientManager) SetHeartbeatInterval(interval time.Duration) {
	cm.heartbeatInterval = interval
}

// GetOnlineClients 获取在线客户端
func (cm *ClientManager) GetOnlineClients() []*models.RegisteredClient {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	var onlineClients []*models.RegisteredClient
	for _, client := range cm.clients {
		if client.IsOnline() {
			onlineClients = append(onlineClients, client)
		}
	}

	return onlineClients
}

// K8s风格认证验证函数

// ValidatePSKCredential 验证PSK凭证
func ValidatePSKCredential(namespace, clientName, secretKey string) bool {
	// 这里应该从配置文件或数据库中验证PSK凭证
	// 为演示目的，使用简单的硬编码验证

	// 示例：验证默认的PSK凭证
	if namespace == "default" && clientName == "OneGoClient" {
		// 这里应该是从安全存储中获取的密钥
		expectedKey := "your-psk-secret-key-here"
		return secretKey == expectedKey
	}

	// 可以添加更多的验证逻辑
	return false
}

// ValidateJWTCredential 验证JWT凭证
func ValidateJWTCredential(namespace, token string) bool {
	// 这里应该使用真正的JWT库验证token
	// 为演示目的，使用简单的验证

	if namespace == "default" {
		// 简化的JWT验证
		return len(token) > 10 // 基本的长度检查
	}

	return false
}

// ValidateBasicAuthCredential 验证基础认证凭证
func ValidateBasicAuthCredential(namespace, username, password string) bool {
	// 这里应该从用户数据库验证用户名和密码
	// 为演示目的，使用简单的硬编码验证

	if namespace == "default" {
		// 示例用户
		validUsers := map[string]string{
			"admin":       "admin123",
			"client":      "client123",
			"onegoclient": "onegoclient123",
		}

		if expectedPassword, exists := validUsers[username]; exists {
			return password == expectedPassword
		}
	}

	return false
}

// ValidateCertificateCredential 验证证书凭证
func ValidateCertificateCredential(namespace, certificate string) bool {
	// 这里应该验证客户端证书的有效性
	// 为演示目的，使用简单的验证

	if namespace == "default" {
		// 基本的证书格式检查
		return len(certificate) > 50 &&
			(strings.Contains(certificate, "BEGIN CERTIFICATE") ||
				strings.Contains(certificate, "CERTIFICATE"))
	}

	return false
}
