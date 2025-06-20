package handlers

import (
	"net/http"
	"time"

	"learngo0619/internal/config"
	"learngo0619/internal/logger"
	"learngo0619/internal/server/models"
	"learngo0619/internal/server/services"

	"encoding/base64"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	// 全局客户端管理器实例
	clientManager *services.ClientManager
	// 全局服务器配置实例
	globalServerConfig *config.Config
)

// InitClientManager 初始化客户端管理器
func InitClientManager() {
	clientManager = services.NewClientManager()
	logger.Info("Client manager initialized")
}

// SetGlobalConfig 设置全局配置
func SetGlobalConfig(cfg *config.Config) {
	globalServerConfig = cfg
	logger.Info("Global server config set for handlers")
}

// StopClientManager 停止客户端管理器
func StopClientManager() {
	if clientManager != nil {
		clientManager.Stop()
	}
}

// ClientRegisterHandler 客户端注册接口（简化版 - 仅预定义Token认证）
func ClientRegisterHandler(c *gin.Context) {
	var req models.ClientRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求格式错误: " + err.Error(),
			"code":    400,
		})
		return
	}

	// 获取客户端IP和User-Agent
	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()

	// 验证预定义Token（唯一认证方式）
	if err := validatePredefinedToken(c, ip, userAgent, &req); err != nil {
		logger.WarnForClient(ip, userAgent,
			"Predefined token authentication failed",
			zap.String("client_name", req.Name),
			zap.Error(err),
		)

		c.JSON(http.StatusUnauthorized, gin.H{
			"success":    false,
			"message":    err.Error(),
			"code":       401,
			"error_type": "token_invalid",
		})
		return
	}

	logger.InfoForClient(ip, userAgent,
		"Predefined token authentication successful",
		zap.String("client_name", req.Name),
	)

	// 认证通过，继续注册流程
	response, err := clientManager.RegisterClient(&req, ip, userAgent)
	if err != nil {
		logger.WarnForClient(ip, userAgent,
			"Client registration failed",
			zap.Error(err),
			zap.String("client_name", req.Name),
			zap.String("client_type", string(req.Type)),
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"code":    400,
		})
		return
	}

	// 记录注册成功日志
	logger.InfoForClient(ip, userAgent,
		"Client registered successfully",
		zap.String("client_id", response.ClientID),
		zap.String("client_name", req.Name),
		zap.String("client_type", string(req.Type)),
		zap.String("client_version", req.Version),
		zap.String("auth_method", "predefined-token"),
	)

	c.JSON(http.StatusOK, response)
}

// getServerConfig 获取服务器配置
func getServerConfig() *config.Config {
	return globalServerConfig
}

// validatePredefinedToken 验证预定义Token（简化版认证）
func validatePredefinedToken(c *gin.Context, ip, userAgent string, req *models.ClientRegisterRequest) error {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return fmt.Errorf("缺少授权令牌，请在请求头中添加 'Authorization: Bearer <token>'")
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		return fmt.Errorf("授权令牌格式错误，应为 'Bearer <token>'")
	}

	token := authHeader[len(bearerPrefix):]

	// 获取配置中的预定义Token
	serverConfig := getServerConfig()
	if serverConfig == nil {
		return fmt.Errorf("服务器配置未找到")
	}

	logger.InfoForClient(ip, userAgent,
		"Debug: checking predefined token",
		zap.String("client_name", req.Name),
		zap.String("token_value", serverConfig.Security.PredefinedToken),
		zap.Bool("token_empty", serverConfig.Security.PredefinedToken == ""),
	)

	if serverConfig.Security.PredefinedToken == "" {
		return fmt.Errorf("服务器未配置预定义Token")
	}

	// 验证Token是否匹配
	if token != serverConfig.Security.PredefinedToken {
		return fmt.Errorf("预定义Token验证失败，请检查Token是否正确")
	}

	logger.InfoForClient(ip, userAgent,
		"Predefined token validated successfully",
		zap.String("client_name", req.Name),
		zap.String("token_prefix", token[:8]+"..."),
	)

	return nil
}

// validateK8sStyleCredential 验证K8s风格凭证
func validateK8sStyleCredential(c *gin.Context, credentialType, ip, userAgent string, req *models.ClientRegisterRequest) error {
	switch credentialType {
	case "psk":
		return validatePSKCredential(c, ip, userAgent, req)
	case "jwt":
		return validateJWTCredential(c, ip, userAgent, req)
	case "basic":
		return validateBasicAuthCredential(c, ip, userAgent, req)
	case "certificate":
		return validateCertificateCredential(c, ip, userAgent, req)
	default:
		return fmt.Errorf("unsupported credential type: %s", credentialType)
	}
}

// validatePSKCredential 验证PSK凭证
func validatePSKCredential(c *gin.Context, ip, userAgent string, req *models.ClientRegisterRequest) error {
	namespace := c.GetHeader("X-Client-Namespace")
	clientName := c.GetHeader("X-Client-Name")
	secretKey := c.GetHeader("X-Client-Secret")

	if namespace == "" || clientName == "" || secretKey == "" {
		return fmt.Errorf("PSK凭证信息不完整")
	}

	// 这里应该从数据库或配置文件验证PSK
	// 为演示目的，使用简单的验证逻辑
	if !services.ValidatePSKCredential(namespace, clientName, secretKey) {
		return fmt.Errorf("PSK凭证验证失败")
	}

	return nil
}

// validateJWTCredential 验证JWT凭证
func validateJWTCredential(c *gin.Context, ip, userAgent string, req *models.ClientRegisterRequest) error {
	authHeader := c.GetHeader("Authorization")
	namespace := c.GetHeader("X-Client-Namespace")

	if authHeader == "" {
		return fmt.Errorf("缺少JWT令牌")
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		return fmt.Errorf("JWT令牌格式错误")
	}

	token := authHeader[len(bearerPrefix):]

	// 验证JWT令牌
	if !services.ValidateJWTCredential(namespace, token) {
		return fmt.Errorf("JWT令牌验证失败")
	}

	return nil
}

// validateBasicAuthCredential 验证基础认证凭证
func validateBasicAuthCredential(c *gin.Context, ip, userAgent string, req *models.ClientRegisterRequest) error {
	authHeader := c.GetHeader("Authorization")
	namespace := c.GetHeader("X-Client-Namespace")

	if authHeader == "" {
		return fmt.Errorf("缺少基础认证信息")
	}

	const basicPrefix = "Basic "
	if len(authHeader) < len(basicPrefix) || authHeader[:len(basicPrefix)] != basicPrefix {
		return fmt.Errorf("基础认证格式错误")
	}

	// 解析用户名密码
	auth := authHeader[len(basicPrefix):]
	decoded, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		return fmt.Errorf("基础认证解码失败")
	}

	credentials := strings.SplitN(string(decoded), ":", 2)
	if len(credentials) != 2 {
		return fmt.Errorf("基础认证格式错误")
	}

	username, password := credentials[0], credentials[1]

	// 验证用户名密码
	if !services.ValidateBasicAuthCredential(namespace, username, password) {
		return fmt.Errorf("用户名或密码错误")
	}

	return nil
}

// validateCertificateCredential 验证证书凭证
func validateCertificateCredential(c *gin.Context, ip, userAgent string, req *models.ClientRegisterRequest) error {
	certificate := c.GetHeader("X-Client-Certificate")
	namespace := c.GetHeader("X-Client-Namespace")

	if certificate == "" {
		return fmt.Errorf("缺少客户端证书")
	}

	// 验证证书
	if !services.ValidateCertificateCredential(namespace, certificate) {
		return fmt.Errorf("客户端证书验证失败")
	}

	return nil
}

// validateBearerToken 验证Bearer Token（支持预定义Token和动态Token）
func validateBearerToken(c *gin.Context, ip, userAgent string, req *models.ClientRegisterRequest) error {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return fmt.Errorf("缺少授权令牌，请在请求头中添加 'Authorization: Bearer <token>'")
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		return fmt.Errorf("授权令牌格式错误，应为 'Bearer <token>'")
	}

	token := authHeader[len(bearerPrefix):]

	// 首先检查是否为预定义Token
	serverConfig := getServerConfig()
	if serverConfig != nil && serverConfig.Security.PredefinedToken != "" {
		if token == serverConfig.Security.PredefinedToken {
			logger.InfoForClient(ip, userAgent,
				"Predefined token authentication successful",
				zap.String("client_name", req.Name),
				zap.String("token_prefix", token[:8]+"..."),
			)
			return nil
		}
	}

	// 如果不是预定义Token，则检查动态Token（向后兼容）
	if !services.ValidateToken(token) {
		return fmt.Errorf("注册令牌无效或已过期。请使用有效的注册令牌，或联系管理员获取")
	}

	logger.InfoForClient(ip, userAgent,
		"Dynamic token authentication successful",
		zap.String("client_name", req.Name),
		zap.String("token_prefix", token[:8]+"..."),
	)

	return nil
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
