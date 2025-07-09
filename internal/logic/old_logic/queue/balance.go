// package queue

// import (
// 	"context"
// 	"fmt"
// 	"math"
// 	"strings"

// 	"github.com/gogf/gf/v2/errors/gerror"

// 	queue "OneGfServer/internal/model/queue"
// )

// // ===============================
// // 负载均衡业务逻辑
// // ===============================

// // CalculateLoadBalance 计算负载均衡
// func (s *sQueue) CalculateLoadBalance(ctx context.Context, input *queue.CalculateLoadBalanceInput) (*queue.CalculateLoadBalanceOutput, error) {
// 	loadBalance := make(map[string]interface{})

// 	// 计算总负载
// 	totalLoad := 0.0
// 	queueLoads := make(map[string]float64)

// 	for _, queue := range input.Queues {
// 		if queueID, ok := queue["queue_id"].(string); ok {
// 			// 计算单个队列的负载
// 			load := s.calculateQueueLoad(queue)
// 			queueLoads[queueID] = load
// 			totalLoad += load
// 		}
// 	}

// 	// 计算负载分布
// 	loadDistribution := make(map[string]float64)
// 	if totalLoad > 0 {
// 		for queueID, load := range queueLoads {
// 			loadDistribution[queueID] = load / totalLoad
// 		}
// 	}

// 	// 计算负载均衡度
// 	balanceScore := s.calculateBalanceScore(queueLoads)

// 	// 生成负载均衡建议
// 	recommendations := s.generateLoadBalanceRecommendations(queueLoads, balanceScore)

// 	loadBalance["total_load"] = totalLoad
// 	loadBalance["queue_loads"] = queueLoads
// 	loadBalance["load_distribution"] = loadDistribution
// 	loadBalance["balance_score"] = balanceScore
// 	loadBalance["recommendations"] = recommendations

// 	return &queue.CalculateLoadBalanceOutput{
// 		LoadBalance: loadBalance,
// 	}, nil
// }

// // SelectOptimalQueue 选择最优队列
// func (s *sQueue) SelectOptimalQueue(ctx context.Context, input *queue.SelectOptimalQueueInput) (*queue.SelectOptimalQueueOutput, error) {
// 	var optimalQueue map[string]interface{}
// 	var bestScore float64
// 	var reason string

// 	for _, queue := range input.AvailableQueues {
// 		score := s.calculateQueueScore(queue, input.TaskData)
// 		if score > bestScore {
// 			bestScore = score
// 			optimalQueue = queue
// 		}
// 	}

// 	if optimalQueue == nil {
// 		return nil, gerror.New("没有可用的队列")
// 	}

// 	// 生成选择原因
// 	reason = s.generateQueueSelectionReason(optimalQueue, bestScore)

// 	return &queue.SelectOptimalQueueOutput{
// 		OptimalQueue: optimalQueue,
// 		Score:        bestScore,
// 		Reason:       reason,
// 	}, nil
// }

// // calculateQueueScore 计算队列评分
// func (s *sQueue) calculateQueueScore(queue map[string]interface{}, taskData map[string]interface{}) float64 {
// 	score := 0.0

// 	// 1. 队列健康度 (权重: 30%)
// 	if healthScore, ok := queue["health_score"].(float64); ok {
// 		score += healthScore * 0.30
// 	}

// 	// 2. 负载情况 (权重: 25%)
// 	if currentLoad, ok := queue["current_load"].(float64); ok {
// 		// 负载越低分数越高
// 		loadScore := 100 - currentLoad
// 		score += loadScore * 0.25
// 	}

// 	// 3. 响应时间 (权重: 20%)
// 	if avgResponseTime, ok := queue["avg_response_time"].(float64); ok {
// 		// 响应时间越短分数越高
// 		responseScore := math.Max(0, 100-avgResponseTime/10)
// 		score += responseScore * 0.20
// 	}

// 	// 4. 错误率 (权重: 15%)
// 	if errorRate, ok := queue["error_rate"].(float64); ok {
// 		// 错误率越低分数越高
// 		errorScore := 100 - errorRate*100
// 		score += errorScore * 0.15
// 	}

// 	// 5. 队列类型匹配度 (权重: 10%)
// 	if queueType, ok := queue["type"].(string); ok {
// 		if taskType, ok := taskData["type"].(string); ok {
// 			if s.isTypeCompatible(queueType, taskType) {
// 				score += 100 * 0.10
// 			}
// 		}
// 	}

// 	return score
// }

// // isTypeCompatible 检查类型兼容性
// func (s *sQueue) isTypeCompatible(queueType, taskType string) bool {
// 	// 这里可以根据实际业务逻辑定义类型兼容性规则
// 	compatibilityMap := map[string][]string{
// 		"priority": {"high_priority", "normal", "low_priority"},
// 		"fifo":     {"normal", "batch"},
// 		"lifo":     {"normal", "batch"},
// 		"weighted": {"high_priority", "normal", "low_priority"},
// 	}

// 	if compatibleTypes, exists := compatibilityMap[queueType]; exists {
// 		for _, compatibleType := range compatibleTypes {
// 			if taskType == compatibleType {
// 				return true
// 			}
// 		}
// 	}

// 	return false
// }

// // calculateQueueLoad 计算队列负载
// func (s *sQueue) calculateQueueLoad(queue map[string]interface{}) float64 {
// 	load := 0.0

// 	// 基于当前长度和容量计算负载
// 	if currentLength, ok := queue["current_length"].(int); ok {
// 		if capacity, ok := queue["capacity"].(int); ok && capacity > 0 {
// 			load += float64(currentLength) / float64(capacity)
// 		}
// 	}

// 	// 基于处理速率计算负载
// 	if processingRate, ok := queue["processing_rate"].(float64); ok {
// 		load += processingRate / 100.0 // 假设100是最大处理速率
// 	}

// 	// 基于错误率计算负载
// 	if errorRate, ok := queue["error_rate"].(float64); ok {
// 		load += errorRate
// 	}

// 	return load
// }

// // calculateBalanceScore 计算负载均衡度评分
// func (s *sQueue) calculateBalanceScore(queueLoads map[string]float64) float64 {
// 	if len(queueLoads) == 0 {
// 		return 0.0
// 	}

// 	// 计算平均负载
// 	totalLoad := 0.0
// 	for _, load := range queueLoads {
// 		totalLoad += load
// 	}
// 	avgLoad := totalLoad / float64(len(queueLoads))

// 	// 计算负载方差
// 	variance := 0.0
// 	for _, load := range queueLoads {
// 		variance += math.Pow(load-avgLoad, 2)
// 	}
// 	variance /= float64(len(queueLoads))

// 	// 计算均衡度评分 (方差越小，均衡度越高)
// 	balanceScore := 100.0 - math.Sqrt(variance)*10.0
// 	return math.Max(0, math.Min(balanceScore, 100))
// }

// // generateLoadBalanceRecommendations 生成负载均衡建议
// func (s *sQueue) generateLoadBalanceRecommendations(queueLoads map[string]float64, balanceScore float64) []string {
// 	var recommendations []string

// 	if balanceScore < 50 {
// 		recommendations = append(recommendations, "负载分布不均衡，建议重新分配任务")
// 	}

// 	// 找出负载最高的队列
// 	var maxLoad float64
// 	var maxLoadQueue string
// 	for queueID, load := range queueLoads {
// 		if load > maxLoad {
// 			maxLoad = load
// 			maxLoadQueue = queueID
// 		}
// 	}

// 	if maxLoad > 0.8 {
// 		recommendations = append(recommendations, fmt.Sprintf("队列 %s 负载过高，建议扩容或分流", maxLoadQueue))
// 	}

// 	// 找出负载最低的队列
// 	var minLoad float64 = 1.0
// 	var minLoadQueue string
// 	for queueID, load := range queueLoads {
// 		if load < minLoad {
// 			minLoad = load
// 			minLoadQueue = queueID
// 		}
// 	}

// 	if minLoad < 0.2 && maxLoad > 0.6 {
// 		recommendations = append(recommendations, fmt.Sprintf("建议将部分任务从 %s 迁移到 %s", maxLoadQueue, minLoadQueue))
// 	}

// 	return recommendations
// }

// // generateQueueSelectionReason 生成队列选择原因
// func (s *sQueue) generateQueueSelectionReason(queue map[string]interface{}, score float64) string {
// 	reasons := []string{}

// 	// 基于健康度
// 	if healthScore, ok := queue["health_score"].(float64); ok {
// 		if healthScore > 80 {
// 			reasons = append(reasons, "健康度高")
// 		}
// 	}

// 	// 基于负载
// 	if currentLoad, ok := queue["current_load"].(float64); ok {
// 		if currentLoad < 0.5 {
// 			reasons = append(reasons, "负载较低")
// 		}
// 	}

// 	// 基于响应时间
// 	if avgResponseTime, ok := queue["avg_response_time"].(float64); ok {
// 		if avgResponseTime < 1000 {
// 			reasons = append(reasons, "响应时间短")
// 		}
// 	}

// 	// 基于错误率
// 	if errorRate, ok := queue["error_rate"].(float64); ok {
// 		if errorRate < 0.01 {
// 			reasons = append(reasons, "错误率低")
// 		}
// 	}

// 	if len(reasons) == 0 {
// 		return "综合评分最优"
// 	}

// 	return strings.Join(reasons, "，")
// }
