package services

import (
	"crypto/md5"
	"fmt"
	"sync"
	"time"

	"learngo0619/internal/logger"
	"learngo0619/internal/server/models"

	"go.uber.org/zap"
)

// ClientManager 客户端管理器
type ClientManager struct {
	clients           map[string]*models.RegisteredClient // 客户端映射
	mutex             sync.RWMutex                        // 读写锁
	heartbeatInterval time.Duration                       // 心跳间隔
	cleanupInterval   time.Duration                       // 清理间隔
	stopCleanup       chan struct{}                       // 停止清理信号
}

// NewClientManager 创建新的客户端管理器
func NewClientManager() *ClientManager {
	cm := &ClientManager{
		clients:           make(map[string]*models.RegisteredClient),
		heartbeatInterval: 30 * time.Second, // 默认30秒心跳
		cleanupInterval:   1 * time.Minute,  // 每分钟清理一次
		stopCleanup:       make(chan struct{}),
	}

	// 启动清理协程
	go cm.startCleanupRoutine()

	return cm
}

// RegisterClient 注册新客户端
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

	// 检查是否已经注册
	if existingClient, exists := cm.clients[clientID]; exists {
		// 更新现有客户端信息
		existingClient.Name = req.Name
		existingClient.Type = req.Type
		existingClient.Version = req.Version
		existingClient.Description = req.Description
		existingClient.Status = models.ClientStatusOnline
		existingClient.LastHeartbeat = time.Now()
		existingClient.LastActiveTime = time.Now()
		existingClient.Metadata = req.Metadata
		existingClient.Tags = req.Tags

		logger.Info("Client re-registered",
			zap.String("client_id", clientID),
			zap.String("name", req.Name),
			zap.String("ip", ip),
		)

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
	now := time.Now()
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

// ProcessHeartbeat 处理客户端心跳
func (cm *ClientManager) ProcessHeartbeat(req *models.ClientHeartbeatRequest) (*models.ClientHeartbeatResponse, error) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	client, exists := cm.clients[req.ClientID]
	if !exists {
		return &models.ClientHeartbeatResponse{
			Success:       false,
			Message:       "客户端未注册",
			ServerTime:    time.Now(),
			NextHeartbeat: int(cm.heartbeatInterval.Seconds()),
		}, fmt.Errorf("client not found: %s", req.ClientID)
	}

	// 更新心跳信息
	client.UpdateHeartbeat(req.Status, req.Metadata)

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

// cleanupOfflineClients 清理离线客户端
func (cm *ClientManager) cleanupOfflineClients() {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	now := time.Now()
	offlineThreshold := 10 * time.Minute // 10分钟无心跳认为离线

	var toDelete []string

	for clientID, client := range cm.clients {
		if now.Sub(client.LastHeartbeat) > offlineThreshold {
			toDelete = append(toDelete, clientID)
		}
	}

	// 删除离线客户端
	for _, clientID := range toDelete {
		client := cm.clients[clientID]
		logger.Info("Cleaned up offline client",
			zap.String("client_id", clientID),
			zap.String("name", client.Name),
			zap.Duration("offline_duration", now.Sub(client.LastHeartbeat)),
		)
		delete(cm.clients, clientID)
	}

	if len(toDelete) > 0 {
		logger.Info("Client cleanup completed",
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
