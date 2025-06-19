package models

import (
	"fmt"
	"time"
)

// ClientStatus 客户端状态枚举
type ClientStatus string

const (
	ClientStatusOnline  ClientStatus = "online"  // 在线
	ClientStatusOffline ClientStatus = "offline" // 离线
	ClientStatusIdle    ClientStatus = "idle"    // 空闲
	ClientStatusBusy    ClientStatus = "busy"    // 繁忙
)

// ClientType 客户端类型枚举
type ClientType string

const (
	ClientTypeWeb     ClientType = "web"     // Web浏览器
	ClientTypeMobile  ClientType = "mobile"  // 移动应用
	ClientTypeDesktop ClientType = "desktop" // 桌面应用
	ClientTypeAPI     ClientType = "api"     // API客户端
	ClientTypeBot     ClientType = "bot"     // 机器人/爬虫
	ClientTypeOther   ClientType = "other"   // 其他
)

// RegisteredClient 已注册的客户端信息
type RegisteredClient struct {
	// 基本信息
	ClientID    string     `json:"client_id"`   // 客户端唯一标识
	Name        string     `json:"name"`        // 客户端名称
	Type        ClientType `json:"type"`        // 客户端类型
	Version     string     `json:"version"`     // 客户端版本
	Description string     `json:"description"` // 客户端描述

	// 网络信息
	IPAddress   string `json:"ip_address"`  // IP地址
	UserAgent   string `json:"user_agent"`  // User-Agent
	Fingerprint string `json:"fingerprint"` // 客户端指纹

	// 状态信息
	Status         ClientStatus `json:"status"`           // 当前状态
	LastHeartbeat  time.Time    `json:"last_heartbeat"`   // 最后心跳时间
	RegisterTime   time.Time    `json:"register_time"`    // 注册时间
	LastActiveTime time.Time    `json:"last_active_time"` // 最后活跃时间

	// 统计信息
	RequestCount   int64         `json:"request_count"`   // 请求总数
	ErrorCount     int64         `json:"error_count"`     // 错误总数
	OnlineDuration time.Duration `json:"online_duration"` // 在线时长

	// 扩展信息
	Metadata map[string]interface{} `json:"metadata,omitempty"` // 元数据
	Tags     []string               `json:"tags,omitempty"`     // 标签
}

// ClientRegisterRequest 客户端注册请求
type ClientRegisterRequest struct {
	Name        string                 `json:"name" binding:"required"`    // 客户端名称
	Type        ClientType             `json:"type" binding:"required"`    // 客户端类型
	Version     string                 `json:"version" binding:"required"` // 客户端版本
	Description string                 `json:"description,omitempty"`      // 客户端描述
	Metadata    map[string]interface{} `json:"metadata,omitempty"`         // 元数据
	Tags        []string               `json:"tags,omitempty"`             // 标签
}

// ClientHeartbeatRequest 客户端心跳请求
type ClientHeartbeatRequest struct {
	ClientID string                 `json:"client_id" binding:"required"` // 客户端ID
	Status   ClientStatus           `json:"status,omitempty"`             // 当前状态
	Metadata map[string]interface{} `json:"metadata,omitempty"`           // 更新的元数据
}

// ClientRegisterResponse 客户端注册响应
type ClientRegisterResponse struct {
	ClientID          string    `json:"client_id"`          // 分配的客户端ID
	Success           bool      `json:"success"`            // 注册是否成功
	Message           string    `json:"message"`            // 响应消息
	RegisterTime      time.Time `json:"register_time"`      // 注册时间
	HeartbeatUrl      string    `json:"heartbeat_url"`      // 心跳URL
	HeartbeatInterval int       `json:"heartbeat_interval"` // 心跳间隔(秒)
}

// ClientHeartbeatResponse 客户端心跳响应
type ClientHeartbeatResponse struct {
	Success       bool      `json:"success"`        // 心跳是否成功
	Message       string    `json:"message"`        // 响应消息
	ServerTime    time.Time `json:"server_time"`    // 服务器时间
	NextHeartbeat int       `json:"next_heartbeat"` // 下次心跳间隔(秒)
}

// ClientListResponse 客户端列表响应
type ClientListResponse struct {
	Total        int                 `json:"total"`         // 总数
	OnlineCount  int                 `json:"online_count"`  // 在线数量
	OfflineCount int                 `json:"offline_count"` // 离线数量
	Clients      []*RegisteredClient `json:"clients"`       // 客户端列表
	Timestamp    time.Time           `json:"timestamp"`     // 响应时间
}

// IsOnline 检查客户端是否在线
func (c *RegisteredClient) IsOnline() bool {
	// 如果超过5分钟没有心跳，认为离线
	return c.Status == ClientStatusOnline && time.Since(c.LastHeartbeat) < 5*time.Minute
}

// UpdateHeartbeat 更新心跳信息
func (c *RegisteredClient) UpdateHeartbeat(status ClientStatus, metadata map[string]interface{}) {
	c.LastHeartbeat = time.Now()
	c.LastActiveTime = time.Now()

	if status != "" {
		c.Status = status
	}

	// 更新元数据
	if metadata != nil {
		if c.Metadata == nil {
			c.Metadata = make(map[string]interface{})
		}
		for k, v := range metadata {
			c.Metadata[k] = v
		}
	}
}

// IncrementRequestCount 增加请求计数
func (c *RegisteredClient) IncrementRequestCount() {
	c.RequestCount++
	c.LastActiveTime = time.Now()
}

// IncrementErrorCount 增加错误计数
func (c *RegisteredClient) IncrementErrorCount() {
	c.ErrorCount++
}

// UpdateOnlineDuration 更新在线时长
func (c *RegisteredClient) UpdateOnlineDuration() {
	if c.Status == ClientStatusOnline {
		c.OnlineDuration = time.Since(c.RegisterTime)
	}
}

// GetSummary 获取客户端摘要信息
func (c *RegisteredClient) GetSummary() map[string]interface{} {
	return map[string]interface{}{
		"client_id":       c.ClientID,
		"name":            c.Name,
		"type":            c.Type,
		"status":          c.Status,
		"ip_address":      c.IPAddress,
		"is_online":       c.IsOnline(),
		"last_heartbeat":  c.LastHeartbeat,
		"register_time":   c.RegisterTime,
		"request_count":   c.RequestCount,
		"error_count":     c.ErrorCount,
		"online_duration": c.OnlineDuration.String(),
	}
}

// String 返回客户端的字符串表示
func (c *RegisteredClient) String() string {
	return fmt.Sprintf("Client[%s] %s (%s) - %s @ %s",
		c.ClientID, c.Name, c.Type, c.Status, c.IPAddress)
}

// Validate 验证客户端注册请求
func (req *ClientRegisterRequest) Validate() error {
	if req.Name == "" {
		return fmt.Errorf("客户端名称不能为空")
	}
	if req.Type == "" {
		return fmt.Errorf("客户端类型不能为空")
	}
	if req.Version == "" {
		return fmt.Errorf("客户端版本不能为空")
	}

	// 验证客户端类型
	validTypes := []ClientType{
		ClientTypeWeb, ClientTypeMobile, ClientTypeDesktop,
		ClientTypeAPI, ClientTypeBot, ClientTypeOther,
	}

	valid := false
	for _, validType := range validTypes {
		if req.Type == validType {
			valid = true
			break
		}
	}

	if !valid {
		return fmt.Errorf("无效的客户端类型: %s", req.Type)
	}

	return nil
}
