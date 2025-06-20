package handlers

import (
	"net/http"
	"time"

	"learngo0619/internal/logger"
	"learngo0619/internal/server/models"
	"learngo0619/internal/server/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	// 全局客户端管理器实例
	clientManager *services.ClientManager
)

// InitClientManager 初始化客户端管理器
func InitClientManager() {
	clientManager = services.NewClientManager()
	logger.Info("Client manager initialized")
}

// StopClientManager 停止客户端管理器
func StopClientManager() {
	if clientManager != nil {
		clientManager.Stop()
	}
}

type RegisterRequest struct {
	Token string `json:"token" binding:"required"`
	IP    string `json:"ip"`
}

type RegisterResponse struct {
	ClientID   string `json:"client_id"`
	Registered bool   `json:"registered"`
	ExpiresAt  string `json:"expires_at"`
	Msg        string `json:"msg"`
}

// ClientRegisterHandler 客户端注册接口
func ClientRegisterHandler(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if !services.ValidateToken(req.Token) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效或过期的token"})
		return
	}
	ip := req.IP
	if ip == "" {
		ip = c.ClientIP()
	}
	clientID, isNew := services.RegisterClient(ip)
	resp := RegisterResponse{
		ClientID:   clientID,
		Registered: isNew,
		ExpiresAt:  services.GetRegisterToken().ExpiresAt.Format("2006-01-02 15:04:05"),
		Msg:        "注册成功",
	}
	c.JSON(http.StatusOK, resp)
}

// ClientHeartbeatHandler 客户端心跳处理器
func ClientHeartbeatHandler(c *gin.Context) {
	var req models.ClientHeartbeatRequest

	// 绑定JSON请求
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求格式错误: " + err.Error(),
			"code":    400,
		})
		return
	}

	// 处理心跳
	response, err := clientManager.ProcessHeartbeat(&req)
	if err != nil {
		logger.WarnForClient(c.ClientIP(), c.Request.UserAgent(),
			"Client heartbeat failed",
			zap.Error(err),
			zap.String("client_id", req.ClientID),
		)

		c.JSON(http.StatusNotFound, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

// ClientListHandler 获取客户端列表处理器
func ClientListHandler(c *gin.Context) {
	// 获取查询参数
	status := c.Query("status")   // online, offline, all
	clientType := c.Query("type") // web, mobile, desktop, api, bot, other

	// 获取所有客户端
	response := clientManager.GetClientsList()

	// 根据查询参数过滤
	if status != "" || clientType != "" {
		filteredClients := make([]*models.RegisteredClient, 0)

		for _, client := range response.Clients {
			// 状态过滤
			if status != "" {
				switch status {
				case "online":
					if !client.IsOnline() {
						continue
					}
				case "offline":
					if client.IsOnline() {
						continue
					}
				case "all":
					// 不过滤
				default:
					// 具体状态过滤
					if string(client.Status) != status {
						continue
					}
				}
			}

			// 类型过滤
			if clientType != "" && string(client.Type) != clientType {
				continue
			}

			filteredClients = append(filteredClients, client)
		}

		// 更新响应
		response.Clients = filteredClients
		response.Total = len(filteredClients)

		// 重新计算在线/离线数量
		onlineCount := 0
		for _, client := range filteredClients {
			if client.IsOnline() {
				onlineCount++
			}
		}
		response.OnlineCount = onlineCount
		response.OfflineCount = response.Total - onlineCount
	}

	// 记录查询日志
	logger.InfoForClient(c.ClientIP(), c.Request.UserAgent(),
		"Client list requested",
		zap.String("status_filter", status),
		zap.String("type_filter", clientType),
		zap.Int("total_returned", response.Total),
	)

	c.JSON(http.StatusOK, response)
}

// ClientDetailHandler 获取客户端详情处理器
func ClientDetailHandler(c *gin.Context) {
	clientID := c.Param("id")

	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "客户端ID不能为空",
			"code":    400,
		})
		return
	}

	// 获取客户端信息
	client, exists := clientManager.GetClient(clientID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "客户端不存在",
			"code":    404,
		})
		return
	}

	// 更新在线时长
	client.UpdateOnlineDuration()

	// 记录查询日志
	logger.InfoForClient(c.ClientIP(), c.Request.UserAgent(),
		"Client detail requested",
		zap.String("requested_client_id", clientID),
		zap.String("client_name", client.Name),
	)

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"client":    client,
		"timestamp": time.Now(),
	})
}

// ClientUnregisterHandler 客户端注销处理器
func ClientUnregisterHandler(c *gin.Context) {
	clientID := c.Param("id")

	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "客户端ID不能为空",
			"code":    400,
		})
		return
	}

	// 注销客户端
	err := clientManager.UnregisterClient(clientID)
	if err != nil {
		logger.WarnForClient(c.ClientIP(), c.Request.UserAgent(),
			"Client unregister failed",
			zap.Error(err),
			zap.String("client_id", clientID),
		)

		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "客户端不存在或注销失败",
			"code":    404,
		})
		return
	}

	// 记录注销日志
	logger.InfoForClient(c.ClientIP(), c.Request.UserAgent(),
		"Client unregistered",
		zap.String("unregistered_client_id", clientID),
	)

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"message":   "客户端注销成功",
		"timestamp": time.Now(),
	})
}

// ClientStatsHandler 获取客户端统计信息处理器
func ClientStatsHandler(c *gin.Context) {
	stats := clientManager.GetStats()

	// 记录统计查询日志
	logger.InfoForClient(c.ClientIP(), c.Request.UserAgent(),
		"Client stats requested",
		zap.Any("stats", stats),
	)

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"stats":     stats,
		"timestamp": time.Now(),
	})
}

// ClientTypesHandler 获取支持的客户端类型处理器
func ClientTypesHandler(c *gin.Context) {
	types := []map[string]string{
		{"value": "web", "label": "Web浏览器", "description": "网页应用客户端"},
		{"value": "mobile", "label": "移动应用", "description": "手机/平板应用"},
		{"value": "desktop", "label": "桌面应用", "description": "桌面软件客户端"},
		{"value": "api", "label": "API客户端", "description": "API调用客户端"},
		{"value": "bot", "label": "机器人", "description": "爬虫或自动化程序"},
		{"value": "other", "label": "其他", "description": "其他类型客户端"},
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"types":     types,
		"timestamp": time.Now(),
	})
}

// ClientOnlineHandler 获取在线客户端处理器
func ClientOnlineHandler(c *gin.Context) {
	onlineClients := clientManager.GetOnlineClients()

	// 生成摘要信息
	summaries := make([]map[string]interface{}, len(onlineClients))
	for i, client := range onlineClients {
		summaries[i] = client.GetSummary()
	}

	// 记录查询日志
	logger.InfoForClient(c.ClientIP(), c.Request.UserAgent(),
		"Online clients requested",
		zap.Int("online_count", len(onlineClients)),
	)

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"online_count": len(onlineClients),
		"clients":      summaries,
		"timestamp":    time.Now(),
	})
}

// GetClientManager 获取客户端管理器实例（用于其他模块调用）
func GetClientManager() *services.ClientManager {
	return clientManager
}
