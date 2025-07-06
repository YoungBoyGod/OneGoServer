package device

import (
	"context"
	"math"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
)

// ===============================
// 设备健康检查相关业务逻辑
// ===============================

// CalculateDeviceHealthScore 计算设备健康度评分
func (s *sDevice) CalculateDeviceHealthScore(ctx context.Context, deviceData map[string]interface{}) float64 {
	var totalScore float64 = 100.0

	// 1. CPU使用率影响 (权重: 25%)
	if cpuUsage, ok := deviceData["cpu_usage"].(float64); ok {
		if cpuUsage > 90 {
			totalScore -= 25
		} else if cpuUsage > 70 {
			totalScore -= 15
		} else if cpuUsage > 50 {
			totalScore -= 5
		}
	}

	// 2. 内存使用率影响 (权重: 25%)
	if memUsage, ok := deviceData["memory_usage"].(float64); ok {
		if memUsage > 90 {
			totalScore -= 25
		} else if memUsage > 80 {
			totalScore -= 15
		} else if memUsage > 60 {
			totalScore -= 5
		}
	}

	// 3. 磁盘使用率影响 (权重: 20%)
	if diskUsage, ok := deviceData["disk_usage"].(float64); ok {
		if diskUsage > 95 {
			totalScore -= 20
		} else if diskUsage > 85 {
			totalScore -= 10
		} else if diskUsage > 70 {
			totalScore -= 3
		}
	}

	// 4. 网络连接状态影响 (权重: 15%)
	if networkStatus, ok := deviceData["network_status"].(string); ok {
		switch networkStatus {
		case "disconnected":
			totalScore -= 15
		case "unstable":
			totalScore -= 8
		case "slow":
			totalScore -= 3
		}
	}

	// 5. 最后心跳时间影响 (权重: 10%)
	if lastHeartbeat, ok := deviceData["last_heartbeat"].(*gtime.Time); ok && lastHeartbeat != nil {
		timeDiff := time.Since(lastHeartbeat.Time)
		if timeDiff > 10*time.Minute {
			totalScore -= 10
		} else if timeDiff > 5*time.Minute {
			totalScore -= 5
		} else if timeDiff > 2*time.Minute {
			totalScore -= 2
		}
	}

	// 6. 错误次数影响 (权重: 5%)
	if errorCount, ok := deviceData["error_count"].(int); ok {
		if errorCount > 10 {
			totalScore -= 5
		} else if errorCount > 5 {
			totalScore -= 3
		} else if errorCount > 2 {
			totalScore -= 1
		}
	}

	// 确保评分在0-100范围内
	if totalScore < 0 {
		totalScore = 0
	}
	if totalScore > 100 {
		totalScore = 100
	}

	return math.Round(totalScore*100) / 100
}

// CheckDeviceHealth 检查设备健康状态
func (s *sDevice) CheckDeviceHealth(ctx context.Context, deviceData map[string]interface{}) map[string]interface{} {
	healthScore := s.CalculateDeviceHealthScore(ctx, deviceData)
	healthChecks := s.performHealthChecks(ctx, deviceData)
	alerts := s.generateHealthAlerts(ctx, deviceData, healthScore)

	return map[string]interface{}{
		"health_score":    healthScore,
		"status":          s.getHealthStatus(healthScore),
		"checks":          healthChecks,
		"alerts":          alerts,
		"recommendations": s.generateHealthRecommendations(healthScore, healthChecks),
		"timestamp":       gtime.Now().Format("2006-01-02 15:04:05"),
	}
}

// performHealthChecks 执行健康检查
func (s *sDevice) performHealthChecks(ctx context.Context, deviceData map[string]interface{}) map[string]interface{} {
	checks := make(map[string]interface{})

	// CPU检查
	if cpuUsage, ok := deviceData["cpu_usage"].(float64); ok {
		checks["cpu"] = map[string]interface{}{
			"usage":   cpuUsage,
			"status":  s.getCheckStatus(cpuUsage, 70, 90),
			"message": s.getCPUMessage(cpuUsage),
		}
	}

	// 内存检查
	if memUsage, ok := deviceData["memory_usage"].(float64); ok {
		checks["memory"] = map[string]interface{}{
			"usage":   memUsage,
			"status":  s.getCheckStatus(memUsage, 80, 90),
			"message": s.getMemoryMessage(memUsage),
		}
	}

	// 磁盘检查
	if diskUsage, ok := deviceData["disk_usage"].(float64); ok {
		checks["disk"] = map[string]interface{}{
			"usage":   diskUsage,
			"status":  s.getCheckStatus(diskUsage, 85, 95),
			"message": s.getDiskMessage(diskUsage),
		}
	}

	// 网络检查
	if networkStatus, ok := deviceData["network_status"].(string); ok {
		checks["network"] = map[string]interface{}{
			"status":  s.getNetworkStatus(networkStatus),
			"message": s.getNetworkMessage(networkStatus),
		}
	}

	return checks
}

// generateHealthAlerts 生成健康告警
func (s *sDevice) generateHealthAlerts(ctx context.Context, deviceData map[string]interface{}, healthScore float64) []map[string]interface{} {
	var alerts []map[string]interface{}

	// 健康度告警
	if healthScore < 30 {
		alerts = append(alerts, map[string]interface{}{
			"level":   "critical",
			"message": "设备健康度严重不足，需要立即处理",
			"code":    "HEALTH_CRITICAL",
		})
	} else if healthScore < 70 {
		alerts = append(alerts, map[string]interface{}{
			"level":   "warning",
			"message": "设备健康度偏低，建议检查",
			"code":    "HEALTH_WARNING",
		})
	}

	// CPU告警
	if cpuUsage, ok := deviceData["cpu_usage"].(float64); ok {
		if cpuUsage > 90 {
			alerts = append(alerts, map[string]interface{}{
				"level":   "critical",
				"message": "CPU使用率过高",
				"code":    "CPU_CRITICAL",
			})
		} else if cpuUsage > 70 {
			alerts = append(alerts, map[string]interface{}{
				"level":   "warning",
				"message": "CPU使用率偏高",
				"code":    "CPU_WARNING",
			})
		}
	}

	// 内存告警
	if memUsage, ok := deviceData["memory_usage"].(float64); ok {
		if memUsage > 90 {
			alerts = append(alerts, map[string]interface{}{
				"level":   "critical",
				"message": "内存使用率过高",
				"code":    "MEMORY_CRITICAL",
			})
		} else if memUsage > 80 {
			alerts = append(alerts, map[string]interface{}{
				"level":   "warning",
				"message": "内存使用率偏高",
				"code":    "MEMORY_WARNING",
			})
		}
	}

	return alerts
}

// generateHealthRecommendations 生成健康建议
func (s *sDevice) generateHealthRecommendations(healthScore float64, checks map[string]interface{}) []string {
	var recommendations []string

	if healthScore < 50 {
		recommendations = append(recommendations, "建议立即检查设备状态")
	}

	if cpuCheck, ok := checks["cpu"].(map[string]interface{}); ok {
		if status, ok := cpuCheck["status"].(string); ok && status == "critical" {
			recommendations = append(recommendations, "建议优化CPU密集型任务")
		}
	}

	if memCheck, ok := checks["memory"].(map[string]interface{}); ok {
		if status, ok := memCheck["status"].(string); ok && status == "critical" {
			recommendations = append(recommendations, "建议增加内存或优化内存使用")
		}
	}

	return recommendations
}

// getHealthStatus 获取健康状态
func (s *sDevice) getHealthStatus(healthScore float64) string {
	if healthScore >= 80 {
		return "excellent"
	} else if healthScore >= 60 {
		return "good"
	} else if healthScore >= 40 {
		return "fair"
	} else if healthScore >= 20 {
		return "poor"
	} else {
		return "critical"
	}
}

// getCheckStatus 获取检查状态
func (s *sDevice) getCheckStatus(value, warningThreshold, criticalThreshold float64) string {
	if value >= criticalThreshold {
		return "critical"
	} else if value >= warningThreshold {
		return "warning"
	} else {
		return "normal"
	}
}

// getCPUMessage 获取CPU消息
func (s *sDevice) getCPUMessage(cpuUsage float64) string {
	if cpuUsage > 90 {
		return "CPU使用率过高，需要立即处理"
	} else if cpuUsage > 70 {
		return "CPU使用率偏高，建议优化"
	} else {
		return "CPU使用率正常"
	}
}

// getMemoryMessage 获取内存消息
func (s *sDevice) getMemoryMessage(memUsage float64) string {
	if memUsage > 90 {
		return "内存使用率过高，需要立即处理"
	} else if memUsage > 80 {
		return "内存使用率偏高，建议优化"
	} else {
		return "内存使用率正常"
	}
}

// getDiskMessage 获取磁盘消息
func (s *sDevice) getDiskMessage(diskUsage float64) string {
	if diskUsage > 95 {
		return "磁盘使用率过高，需要立即处理"
	} else if diskUsage > 85 {
		return "磁盘使用率偏高，建议清理"
	} else {
		return "磁盘使用率正常"
	}
}

// getNetworkStatus 获取网络状态
func (s *sDevice) getNetworkStatus(networkStatus string) string {
	switch networkStatus {
	case "connected":
		return "normal"
	case "slow":
		return "warning"
	case "unstable":
		return "warning"
	case "disconnected":
		return "critical"
	default:
		return "unknown"
	}
}

// getNetworkMessage 获取网络消息
func (s *sDevice) getNetworkMessage(networkStatus string) string {
	switch networkStatus {
	case "connected":
		return "网络连接正常"
	case "slow":
		return "网络连接缓慢"
	case "unstable":
		return "网络连接不稳定"
	case "disconnected":
		return "网络连接断开"
	default:
		return "网络状态未知"
	}
}
