package device

import (
	"context"
	"math"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ===============================
// 设备负载评分相关业务逻辑
// ===============================

// CalculateDeviceLoadScore 计算设备负载评分
func (s *sDevice) CalculateDeviceLoadScore(ctx context.Context, deviceId string) (map[string]interface{}, error) {
	// 获取设备信息
	_, err := s.getDeviceData(ctx, deviceId)
	if err != nil {
		return nil, err
	}

	// 获取设备指标
	metrics, err := s.getDeviceMetrics(ctx, deviceId)
	if err != nil {
		return nil, err
	}

	// 计算负载评分
	loadScore, components := s.calculateLoadScore(metrics)

	// 生成优化建议
	recommendations := s.generateLoadOptimizationRecommendations(loadScore, components, metrics)

	result := map[string]interface{}{
		"deviceId":        deviceId,
		"loadScore":       loadScore,
		"components":      components,
		"recommendations": recommendations,
	}

	return result, nil
}

// calculateLoadScore 计算负载评分
func (s *sDevice) calculateLoadScore(metrics map[string]interface{}) (float64, map[string]interface{}) {
	// 权重配置
	weights := map[string]float64{
		"cpu":     0.3,
		"memory":  0.25,
		"disk":    0.2,
		"network": 0.15,
		"tasks":   0.1,
	}

	// 获取各项指标
	cpuUsage := s.getFloatValue(metrics, "cpu_usage", 0.0)
	memoryUsage := s.getFloatValue(metrics, "memory_usage", 0.0)
	diskUsage := s.getFloatValue(metrics, "disk_usage", 0.0)
	networkLatency := s.getFloatValue(metrics, "network_latency", 0.0)
	currentTasks := s.getIntValue(metrics, "current_tasks", 0)
	maxTasks := s.getIntValue(metrics, "max_tasks", 1)

	// 标准化指标 (0-1)
	cpuScore := s.normalizeMetric(cpuUsage, 100)
	memoryScore := s.normalizeMetric(memoryUsage, 100)
	diskScore := s.normalizeMetric(diskUsage, 100)
	networkScore := s.normalizeLatency(networkLatency)
	taskScore := s.normalizeTaskLoad(currentTasks, maxTasks)

	// 加权求和
	score := weights["cpu"]*cpuScore +
		weights["memory"]*memoryScore +
		weights["disk"]*diskScore +
		weights["network"]*networkScore +
		weights["tasks"]*taskScore

	// 转换为0-100分制
	finalScore := score * 100

	// 构建组件评分
	components := map[string]interface{}{
		"cpuScore":     cpuScore * 100,
		"memoryScore":  memoryScore * 100,
		"diskScore":    diskScore * 100,
		"networkScore": networkScore * 100,
		"taskScore":    taskScore * 100,
	}

	return finalScore, components
}

// normalizeMetric 标准化指标
func (s *sDevice) normalizeMetric(value, maxValue float64) float64 {
	if maxValue <= 0 {
		return 0.0
	}
	normalized := value / maxValue
	return math.Max(0, math.Min(normalized, 1))
}

// normalizeLatency 标准化网络延迟
func (s *sDevice) normalizeLatency(latency float64) float64 {
	// 延迟越低分数越高
	if latency <= 10 {
		return 1.0
	} else if latency <= 50 {
		return 0.8
	} else if latency <= 100 {
		return 0.6
	} else if latency <= 200 {
		return 0.4
	} else if latency <= 500 {
		return 0.2
	} else {
		return 0.0
	}
}

// normalizeTaskLoad 标准化任务负载
func (s *sDevice) normalizeTaskLoad(currentTasks, maxTasks int) float64 {
	if maxTasks <= 0 {
		return 0.0
	}
	utilization := float64(currentTasks) / float64(maxTasks)
	// 任务负载越低分数越高
	return 1.0 - utilization
}

// generateLoadOptimizationRecommendations 生成负载优化建议
func (s *sDevice) generateLoadOptimizationRecommendations(loadScore float64, components map[string]interface{}, metrics map[string]interface{}) []string {
	var recommendations []string

	// 基于总体负载评分
	if loadScore > 80 {
		recommendations = append(recommendations, "设备负载过高，建议减少任务分配或升级设备配置")
	} else if loadScore > 60 {
		recommendations = append(recommendations, "设备负载较高，建议监控任务分配")
	}

	// 基于CPU使用率
	if cpuScore, ok := components["cpuScore"].(float64); ok && cpuScore > 80 {
		recommendations = append(recommendations, "CPU使用率过高，建议优化任务或增加CPU资源")
	}

	// 基于内存使用率
	if memoryScore, ok := components["memoryScore"].(float64); ok && memoryScore > 80 {
		recommendations = append(recommendations, "内存使用率过高，建议清理内存或增加内存容量")
	}

	// 基于磁盘使用率
	if diskScore, ok := components["diskScore"].(float64); ok && diskScore > 80 {
		recommendations = append(recommendations, "磁盘使用率过高，建议清理磁盘空间或扩容")
	}

	// 基于网络延迟
	if networkScore, ok := components["networkScore"].(float64); ok && networkScore < 40 {
		recommendations = append(recommendations, "网络延迟较高，建议检查网络连接或优化网络配置")
	}

	// 基于任务负载
	currentTasks := s.getIntValue(metrics, "current_tasks", 0)
	maxTasks := s.getIntValue(metrics, "max_tasks", 1)
	if maxTasks > 0 && float64(currentTasks)/float64(maxTasks) > 0.8 {
		recommendations = append(recommendations, "任务队列接近满载，建议调整任务分配策略")
	}

	return recommendations
}

// GetDeviceLoadMetrics 获取设备负载指标
func (s *sDevice) GetDeviceLoadMetrics(ctx context.Context, deviceId, period string) (map[string]interface{}, error) {
	// 验证设备存在
	if err := s.validateDeviceExists(ctx, deviceId); err != nil {
		return nil, err
	}

	// 解析时间周期
	duration, err := s.parsePeriod(period)
	if err != nil {
		return nil, err
	}

	// 获取历史指标数据
	metrics, err := s.getDeviceLoadHistory(ctx, deviceId, duration)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"deviceId": deviceId,
		"period":   period,
		"metrics":  metrics,
	}

	return result, nil
}

// parsePeriod 解析时间周期
func (s *sDevice) parsePeriod(period string) (time.Duration, error) {
	switch period {
	case "1h":
		return time.Hour, nil
	case "6h":
		return 6 * time.Hour, nil
	case "24h":
		return 24 * time.Hour, nil
	case "7d":
		return 7 * 24 * time.Hour, nil
	default:
		return 0, gerror.NewCode(gcode.CodeValidationFailed, "不支持的时间周期")
	}
}

// getDeviceLoadHistory 获取设备负载历史
func (s *sDevice) getDeviceLoadHistory(ctx context.Context, deviceId string, duration time.Duration) (map[string]interface{}, error) {
	// 这里应该从数据库或缓存中获取历史数据
	// 目前返回模拟数据
	endTime := time.Now()
	startTime := endTime.Add(-duration)

	metrics := map[string]interface{}{
		"cpuUsage":       s.generateTimeSeriesData(startTime, endTime, 60, 30, 80),
		"memoryUsage":    s.generateTimeSeriesData(startTime, endTime, 60, 40, 70),
		"diskUsage":      s.generateTimeSeriesData(startTime, endTime, 60, 50, 90),
		"networkLatency": s.generateTimeSeriesData(startTime, endTime, 60, 10, 100),
		"currentTasks":   s.generateTimeSeriesData(startTime, endTime, 60, 1, 10),
		"loadScore":      s.generateTimeSeriesData(startTime, endTime, 60, 20, 80),
	}

	return metrics, nil
}

// generateTimeSeriesData 生成时间序列数据
func (s *sDevice) generateTimeSeriesData(startTime, endTime time.Time, intervalSeconds int, minValue, maxValue float64) []map[string]interface{} {
	var data []map[string]interface{}
	currentTime := startTime

	for currentTime.Before(endTime) {
		// 生成随机值
		value := minValue + (maxValue-minValue)*s.randomFloat()

		data = append(data, map[string]interface{}{
			"timestamp": currentTime.Unix(),
			"value":     value,
		})

		currentTime = currentTime.Add(time.Duration(intervalSeconds) * time.Second)
	}

	return data
}

// OptimizeDeviceLoad 优化设备负载
func (s *sDevice) OptimizeDeviceLoad(ctx context.Context, deviceId, strategy string) (map[string]interface{}, error) {
	// 获取当前负载评分
	currentScore, _, err := s.getCurrentLoadScore(ctx, deviceId)
	if err != nil {
		return nil, err
	}

	// 根据策略生成优化建议
	actions := s.generateOptimizationActions(strategy, currentScore)

	// 计算优化后的评分
	optimizedScore := s.calculateOptimizedScore(currentScore, actions)

	// 计算改进程度
	improvement := optimizedScore - currentScore

	result := map[string]interface{}{
		"deviceId":       deviceId,
		"originalScore":  currentScore,
		"optimizedScore": optimizedScore,
		"improvement":    improvement,
		"actions":        actions,
		"status":         "optimized",
	}

	return result, nil
}

// generateOptimizationActions 生成优化动作
func (s *sDevice) generateOptimizationActions(strategy string, currentScore float64) []string {
	var actions []string

	switch strategy {
	case "auto":
		if currentScore > 80 {
			actions = append(actions, "减少任务分配", "优化任务调度", "增加资源监控")
		} else if currentScore > 60 {
			actions = append(actions, "监控负载变化", "优化任务优先级")
		}
	case "manual":
		actions = append(actions, "手动调整任务分配", "手动优化配置")
	case "conservative":
		actions = append(actions, "保守优化", "逐步调整", "监控效果")
	case "aggressive":
		actions = append(actions, "激进优化", "快速调整", "强制负载均衡")
	}

	return actions
}

// calculateOptimizedScore 计算优化后的评分
func (s *sDevice) calculateOptimizedScore(currentScore float64, actions []string) float64 {
	// 根据优化动作计算改进
	improvement := 0.0
	for _, action := range actions {
		switch action {
		case "减少任务分配":
			improvement += 15
		case "优化任务调度":
			improvement += 10
		case "增加资源监控":
			improvement += 5
		case "手动调整任务分配":
			improvement += 20
		case "保守优化":
			improvement += 8
		case "激进优化":
			improvement += 25
		}
	}

	optimizedScore := currentScore - improvement
	return math.Max(0, math.Min(optimizedScore, 100))
}

// SetDeviceLoadThreshold 设置设备负载阈值
func (s *sDevice) SetDeviceLoadThreshold(ctx context.Context, deviceId string, warning, critical float64, maxTasks int) error {
	// 验证参数
	if warning >= critical {
		return gerror.NewCode(gcode.CodeValidationFailed, "警告阈值必须小于严重阈值")
	}

	if warning < 0 || warning > 100 || critical < 0 || critical > 100 {
		return gerror.NewCode(gcode.CodeValidationFailed, "阈值必须在0-100之间")
	}

	if maxTasks <= 0 {
		return gerror.NewCode(gcode.CodeValidationFailed, "最大任务数必须大于0")
	}

	// 保存阈值配置
	threshold := map[string]interface{}{
		"device_id":  deviceId,
		"warning":    warning,
		"critical":   critical,
		"max_tasks":  maxTasks,
		"updated_at": gtime.Now(),
	}

	// 这里应该保存到数据库
	g.Log().Info(ctx, "设置设备负载阈值", g.Map{
		"device_id": deviceId,
		"threshold": threshold,
	})

	return nil
}

// GetDeviceLoadThreshold 获取设备负载阈值
func (s *sDevice) GetDeviceLoadThreshold(ctx context.Context, deviceId string) (map[string]interface{}, error) {
	// 这里应该从数据库获取阈值配置
	// 目前返回默认值
	threshold := map[string]interface{}{
		"device_id":    deviceId,
		"warning":      70.0,
		"critical":     90.0,
		"max_tasks":    10,
		"current_load": 45.0,
		"status":       "normal",
	}

	return threshold, nil
}

// 工具方法
func (s *sDevice) getFloatValue(data map[string]interface{}, key string, defaultValue float64) float64 {
	if value, ok := data[key]; ok {
		if floatValue, ok := value.(float64); ok {
			return floatValue
		}
	}
	return defaultValue
}

func (s *sDevice) getIntValue(data map[string]interface{}, key string, defaultValue int) int {
	if value, ok := data[key]; ok {
		if intValue, ok := value.(int); ok {
			return intValue
		}
	}
	return defaultValue
}

func (s *sDevice) randomFloat() float64 {
	// 简单的随机数生成，实际应该使用更好的随机数生成器
	return float64(time.Now().UnixNano()%100) / 100.0
}

func (s *sDevice) getCurrentLoadScore(ctx context.Context, deviceId string) (float64, map[string]interface{}, error) {
	// 获取设备指标
	metrics, err := s.getDeviceMetrics(ctx, deviceId)
	if err != nil {
		return 0, nil, err
	}

	// 计算负载评分
	loadScore, components := s.calculateLoadScore(metrics)
	return loadScore, components, nil
}

func (s *sDevice) getDeviceData(ctx context.Context, deviceId string) (map[string]interface{}, error) {
	// 这里应该从数据库获取设备数据
	// 目前返回模拟数据
	return map[string]interface{}{
		"device_id": deviceId,
		"name":      "test-device",
		"type":      "sensor",
		"status":    "online",
	}, nil
}

func (s *sDevice) getDeviceMetrics(ctx context.Context, deviceId string) (map[string]interface{}, error) {
	// 这里应该从数据库或监控系统获取设备指标
	// 目前返回模拟数据
	return map[string]interface{}{
		"cpu_usage":       65.5,
		"memory_usage":    45.2,
		"disk_usage":      78.9,
		"network_latency": 25.0,
		"current_tasks":   3,
		"max_tasks":       10,
	}, nil
}

func (s *sDevice) validateDeviceExists(ctx context.Context, deviceId string) error {
	// 这里应该验证设备是否存在
	// 目前简单返回nil
	return nil
}
