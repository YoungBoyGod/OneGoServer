// package queue

// import (
// 	"context"
// 	"sort"
// 	"time"

// 	"github.com/gogf/gf/v2/os/gtime"

// 	queue "OneGfServer/internal/model/queue"
// )

// // ===============================
// // 监控统计业务逻辑
// // ===============================

// // CalculateQueueStatistics 计算队列统计信息
// func (s *sQueue) CalculateQueueStatistics(ctx context.Context, input *queue.CalculateQueueStatisticsInput) (*queue.CalculateQueueStatisticsOutput, error) {
// 	statistics := make(map[string]interface{})

// 	// 计算利用率
// 	utilizationInput := &queue.CalculateUtilizationRateInput{
// 		QueueData: input.QueueData,
// 	}
// 	utilizationOutput := s.calculateUtilizationRate(utilizationInput)
// 	statistics["utilization_rate"] = utilizationOutput.UtilizationRate

// 	// 计算平均处理时间
// 	avgProcessingTimeInput := &queue.CalculateAverageProcessingTimeInput{
// 		HistoricalData: input.HistoricalData,
// 	}
// 	avgProcessingTimeOutput := s.calculateAverageProcessingTime(avgProcessingTimeInput)
// 	statistics["average_processing_time"] = avgProcessingTimeOutput.AverageProcessingTime

// 	// 计算吞吐量
// 	throughputInput := &queue.CalculateThroughputInput{
// 		HistoricalData: input.HistoricalData,
// 	}
// 	throughputOutput := s.calculateThroughput(throughputInput)
// 	statistics["throughput"] = throughputOutput.Throughput

// 	// 计算错误率
// 	errorRateInput := &queue.CalculateErrorRateInput{
// 		HistoricalData: input.HistoricalData,
// 	}
// 	errorRateOutput := s.calculateErrorRate(errorRateInput)
// 	statistics["error_rate"] = errorRateOutput.ErrorRate

// 	// 分析趋势
// 	trendInput := &queue.AnalyzeQueueTrendInput{
// 		HistoricalData: input.HistoricalData,
// 	}
// 	trendOutput := s.analyzeQueueTrend(trendInput)
// 	statistics["trend"] = trendOutput.Trend

// 	// 预测行为
// 	predictionInput := &queue.PredictQueueBehaviorInput{
// 		HistoricalData: input.HistoricalData,
// 	}
// 	predictionOutput := s.predictQueueBehavior(predictionInput)
// 	statistics["prediction"] = predictionOutput.Prediction

// 	return &queue.CalculateQueueStatisticsOutput{
// 		Statistics: statistics,
// 	}, nil
// }

// // calculateUtilizationRate 计算利用率
// func (s *sQueue) calculateUtilizationRate(input *queue.CalculateUtilizationRateInput) *queue.CalculateUtilizationRateOutput {
// 	currentLength, _ := input.QueueData["current_length"].(int)
// 	capacity, _ := input.QueueData["capacity"].(int)

// 	if capacity <= 0 {
// 		return &queue.CalculateUtilizationRateOutput{UtilizationRate: 0}
// 	}

// 	utilizationRate := float64(currentLength) / float64(capacity)
// 	return &queue.CalculateUtilizationRateOutput{UtilizationRate: utilizationRate}
// }

// // calculateAverageProcessingTime 计算平均处理时间
// func (s *sQueue) calculateAverageProcessingTime(input *queue.CalculateAverageProcessingTimeInput) *queue.CalculateAverageProcessingTimeOutput {
// 	if len(input.HistoricalData) == 0 {
// 		return &queue.CalculateAverageProcessingTimeOutput{AverageProcessingTime: 0}
// 	}

// 	totalTime := 0.0
// 	count := 0

// 	for _, record := range input.HistoricalData {
// 		if processTime, ok := record["process_time"].(float64); ok {
// 			totalTime += processTime
// 			count++
// 		}
// 	}

// 	if count == 0 {
// 		return &queue.CalculateAverageProcessingTimeOutput{AverageProcessingTime: 0}
// 	}

// 	avgTime := totalTime / float64(count)
// 	return &queue.CalculateAverageProcessingTimeOutput{AverageProcessingTime: avgTime}
// }

// // calculateThroughput 计算吞吐量
// func (s *sQueue) calculateThroughput(input *queue.CalculateThroughputInput) *queue.CalculateThroughputOutput {
// 	if len(input.HistoricalData) == 0 {
// 		return &queue.CalculateThroughputOutput{Throughput: 0}
// 	}

// 	// 计算时间窗口内的任务数量
// 	oneHourAgo := gtime.Now().Add(-time.Hour)
// 	completedTasks := 0

// 	for _, record := range input.HistoricalData {
// 		if completedAt, ok := record["completed_at"].(*gtime.Time); ok {
// 			if completedAt.After(oneHourAgo) {
// 				completedTasks++
// 			}
// 		}
// 	}

// 	throughput := float64(completedTasks) / 1.0 // 每小时的任务数
// 	return &queue.CalculateThroughputOutput{Throughput: throughput}
// }

// // calculateErrorRate 计算错误率
// func (s *sQueue) calculateErrorRate(input *queue.CalculateErrorRateInput) *queue.CalculateErrorRateOutput {
// 	if len(input.HistoricalData) == 0 {
// 		return &queue.CalculateErrorRateOutput{ErrorRate: 0}
// 	}

// 	totalTasks := len(input.HistoricalData)
// 	failedTasks := 0

// 	for _, record := range input.HistoricalData {
// 		if status, ok := record["status"].(string); ok {
// 			if status == "failed" || status == "error" {
// 				failedTasks++
// 			}
// 		}
// 	}

// 	errorRate := float64(failedTasks) / float64(totalTasks)
// 	return &queue.CalculateErrorRateOutput{ErrorRate: errorRate}
// }

// // analyzeQueueTrend 分析队列趋势
// func (s *sQueue) analyzeQueueTrend(input *queue.AnalyzeQueueTrendInput) *queue.AnalyzeQueueTrendOutput {
// 	trend := make(map[string]interface{})

// 	if len(input.HistoricalData) < 2 {
// 		trend["trend"] = "insufficient_data"
// 		trend["slope"] = 0.0
// 		return &queue.AnalyzeQueueTrendOutput{Trend: trend}
// 	}

// 	// 按时间排序
// 	sort.Slice(input.HistoricalData, func(i, j int) bool {
// 		timeI, _ := input.HistoricalData[i]["timestamp"].(*gtime.Time)
// 		timeJ, _ := input.HistoricalData[j]["timestamp"].(*gtime.Time)
// 		if timeI == nil || timeJ == nil {
// 			return false
// 		}
// 		return timeI.Before(timeJ)
// 	})

// 	// 计算趋势斜率
// 	var values []float64
// 	for _, record := range input.HistoricalData {
// 		if value, ok := record["value"].(float64); ok {
// 			values = append(values, value)
// 		}
// 	}

// 	if len(values) < 2 {
// 		trend["trend"] = "insufficient_data"
// 		trend["slope"] = 0.0
// 		return &queue.AnalyzeQueueTrendOutput{Trend: trend}
// 	}

// 	// 简单线性回归
// 	n := float64(len(values))
// 	sumX := 0.0
// 	sumY := 0.0
// 	sumXY := 0.0
// 	sumX2 := 0.0

// 	for i, y := range values {
// 		x := float64(i)
// 		sumX += x
// 		sumY += y
// 		sumXY += x * y
// 		sumX2 += x * x
// 	}

// 	slope := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)

// 	// 判断趋势
// 	var trendType string
// 	if slope > 0.1 {
// 		trendType = "increasing"
// 	} else if slope < -0.1 {
// 		trendType = "decreasing"
// 	} else {
// 		trendType = "stable"
// 	}

// 	trend["trend"] = trendType
// 	trend["slope"] = slope

// 	return &queue.AnalyzeQueueTrendOutput{Trend: trend}
// }

// // predictQueueBehavior 预测队列行为
// func (s *sQueue) predictQueueBehavior(input *queue.PredictQueueBehaviorInput) *queue.PredictQueueBehaviorOutput {
// 	prediction := make(map[string]interface{})

// 	if len(input.HistoricalData) < 10 {
// 		prediction["confidence"] = "low"
// 		prediction["prediction"] = "insufficient_data"
// 		return &queue.PredictQueueBehaviorOutput{Prediction: prediction}
// 	}

// 	// 基于历史数据预测未来行为
// 	// 这里使用简单的移动平均预测
// 	var recentValues []float64
// 	for i := len(input.HistoricalData) - 10; i < len(input.HistoricalData); i++ {
// 		if value, ok := input.HistoricalData[i]["value"].(float64); ok {
// 			recentValues = append(recentValues, value)
// 		}
// 	}

// 	if len(recentValues) == 0 {
// 		prediction["confidence"] = "low"
// 		prediction["prediction"] = "no_data"
// 		return &queue.PredictQueueBehaviorOutput{Prediction: prediction}
// 	}

// 	// 计算平均值作为预测值
// 	sum := 0.0
// 	for _, value := range recentValues {
// 		sum += value
// 	}
// 	avgValue := sum / float64(len(recentValues))

// 	prediction["predicted_value"] = avgValue
// 	prediction["confidence"] = "medium"
// 	prediction["method"] = "moving_average"

// 	return &queue.PredictQueueBehaviorOutput{Prediction: prediction}
// }
