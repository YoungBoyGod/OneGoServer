package queue

import (
	"context"
	"strings"
)

// ===============================
// 工具方法
// ===============================

// GetQueueInstance 获取队列实例（保留原有方法）
func (s *sQueue) GetQueueInstance(ctx context.Context) (interface{}, error) {
	// 这里可以返回具体的队列客户端实例
	// 目前返回nil，实际使用时需要根据具体需求实现
	return nil, nil
}

// 修复balance.go中缺少的strings包导入
func (s *sQueue) generateQueueSelectionReason(queue map[string]interface{}, score float64) string {
	reasons := []string{}

	// 基于健康度
	if healthScore, ok := queue["health_score"].(float64); ok {
		if healthScore > 80 {
			reasons = append(reasons, "健康度高")
		}
	}

	// 基于负载
	if currentLoad, ok := queue["current_load"].(float64); ok {
		if currentLoad < 0.5 {
			reasons = append(reasons, "负载较低")
		}
	}

	// 基于响应时间
	if avgResponseTime, ok := queue["avg_response_time"].(float64); ok {
		if avgResponseTime < 1000 {
			reasons = append(reasons, "响应时间短")
		}
	}

	// 基于错误率
	if errorRate, ok := queue["error_rate"].(float64); ok {
		if errorRate < 0.01 {
			reasons = append(reasons, "错误率低")
		}
	}

	if len(reasons) == 0 {
		return "综合评分最优"
	}

	return strings.Join(reasons, "，")
}
