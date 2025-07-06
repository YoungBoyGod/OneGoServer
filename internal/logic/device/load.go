package device

import (
	"OneGfServer/internal/model/device"
	"context"
	"math"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
)

// ===============================
// 设备负载管理相关业务逻辑
// ===============================

// CalculateDeviceLoadScore 计算设备负载评分
func (s *sDevice) CalculateDeviceLoadScore(ctx context.Context, input *device.CalculateDeviceLoadScoreInput) (*device.CalculateDeviceLoadScoreOutput, error) {
	// 获取设备数据
	_, err := s.getDeviceData(ctx, input.DeviceID)
	if err != nil {
		return nil, err
	}

	// 获取设备指标
	metrics, err := s.getDeviceMetrics(ctx, input.DeviceID)
	if err != nil {
		return nil, err
	}

	// 计算负载评分
	loadScore, components := s.calculateLoadScore(metrics)

	// 生成优化建议
	recommendations := s.generateLoadOptimizationRecommendations(loadScore, components, metrics)

	return &device.CalculateDeviceLoadScoreOutput{
		DeviceID:        input.DeviceID,
		LoadScore:       loadScore,
		Components:      components,
		Recommendations: recommendations,
		Timestamp:       gtime.Now().Format("2006-01-02 15:04:05"),
		Metrics:         metrics,
	}, nil
}

// calculateLoadScore 计算负载评分
func (s *sDevice) calculateLoadScore(metrics map[string]interface{}) (float64, map[string]interface{}) {
	var totalScore float64
	components := make(map[string]interface{})

	// 1. CPU负载评分 (权重: 30%)
	cpuUsage := s.getFloatValue(metrics, "cpu_usage", 0)
	cpuScore := s.normalizeMetric(cpuUsage, 100)
	components["cpu"] = map[string]interface{}{
		"usage":  cpuUsage,
		"score":  cpuScore,
		"weight": 0.30,
	}
	totalScore += cpuScore * 0.30

	// 2. 内存负载评分 (权重: 25%)
	memUsage := s.getFloatValue(metrics, "memory_usage", 0)
	memScore := s.normalizeMetric(memUsage, 100)
	components["memory"] = map[string]interface{}{
		"usage":  memUsage,
		"score":  memScore,
		"weight": 0.25,
	}
	totalScore += memScore * 0.25

	// 3. 磁盘负载评分 (权重: 20%)
	diskUsage := s.getFloatValue(metrics, "disk_usage", 0)
	diskScore := s.normalizeMetric(diskUsage, 100)
	components["disk"] = map[string]interface{}{
		"usage":  diskUsage,
		"score":  diskScore,
		"weight": 0.20,
	}
	totalScore += diskScore * 0.20

	// 4. 网络负载评分 (权重: 15%)
	networkLatency := s.getFloatValue(metrics, "network_latency", 0)
	networkScore := s.normalizeLatency(networkLatency)
	components["network"] = map[string]interface{}{
		"latency": networkLatency,
		"score":   networkScore,
		"weight":  0.15,
	}
	totalScore += networkScore * 0.15

	// 5. 任务负载评分 (权重: 10%)
	currentTasks := s.getIntValue(metrics, "running_task_count", 0)
	maxTasks := s.getIntValue(metrics, "max_concurrent_tasks", 5)
	taskScore := s.normalizeTaskLoad(currentTasks, maxTasks)
	components["tasks"] = map[string]interface{}{
		"current": currentTasks,
		"max":     maxTasks,
		"score":   taskScore,
		"weight":  0.10,
	}
	totalScore += taskScore * 0.10

	return math.Round(totalScore*100) / 100, components
}

// normalizeMetric 标准化指标值
func (s *sDevice) normalizeMetric(value, maxValue float64) float64 {
	if maxValue <= 0 {
		return 0
	}
	normalized := (value / maxValue) * 100
	return math.Max(0, math.Min(normalized, 100))
}

// normalizeLatency 标准化网络延迟
func (s *sDevice) normalizeLatency(latency float64) float64 {
	// 延迟越低越好，所以需要反转评分
	if latency <= 10 {
		return 100 // 延迟小于10ms，评分100
	} else if latency <= 50 {
		return 80 // 延迟10-50ms，评分80
	} else if latency <= 100 {
		return 60 // 延迟50-100ms，评分60
	} else if latency <= 200 {
		return 40 // 延迟100-200ms，评分40
	} else if latency <= 500 {
		return 20 // 延迟200-500ms，评分20
	} else {
		return 0 // 延迟超过500ms，评分0
	}
}

// normalizeTaskLoad 标准化任务负载
func (s *sDevice) normalizeTaskLoad(currentTasks, maxTasks int) float64 {
	if maxTasks <= 0 {
		return 0
	}
	utilization := float64(currentTasks) / float64(maxTasks)
	if utilization <= 0.5 {
		return 100 // 负载低于50%，评分100
	} else if utilization <= 0.7 {
		return 80 // 负载50-70%，评分80
	} else if utilization <= 0.8 {
		return 60 // 负载70-80%，评分60
	} else if utilization <= 0.9 {
		return 40 // 负载80-90%，评分40
	} else if utilization <= 1.0 {
		return 20 // 负载90-100%，评分20
	} else {
		return 0 // 负载超过100%，评分0
	}
}

// generateLoadOptimizationRecommendations 生成负载优化建议
func (s *sDevice) generateLoadOptimizationRecommendations(loadScore float64, components map[string]interface{}, metrics map[string]interface{}) []string {
	var recommendations []string

	if loadScore < 50 {
		recommendations = append(recommendations, "设备负载过高，建议立即优化")
	}

	// CPU优化建议
	if cpuComp, ok := components["cpu"].(map[string]interface{}); ok {
		if score, ok := cpuComp["score"].(float64); ok && score > 80 {
			recommendations = append(recommendations, "CPU使用率过高，建议优化计算密集型任务")
		}
	}

	// 内存优化建议
	if memComp, ok := components["memory"].(map[string]interface{}); ok {
		if score, ok := memComp["score"].(float64); ok && score > 80 {
			recommendations = append(recommendations, "内存使用率过高，建议增加内存或优化内存使用")
		}
	}

	// 磁盘优化建议
	if diskComp, ok := components["disk"].(map[string]interface{}); ok {
		if score, ok := diskComp["score"].(float64); ok && score > 80 {
			recommendations = append(recommendations, "磁盘使用率过高，建议清理临时文件或扩容")
		}
	}

	// 网络优化建议
	if netComp, ok := components["network"].(map[string]interface{}); ok {
		if score, ok := netComp["score"].(float64); ok && score < 60 {
			recommendations = append(recommendations, "网络延迟较高，建议检查网络连接")
		}
	}

	// 任务优化建议
	if taskComp, ok := components["tasks"].(map[string]interface{}); ok {
		if score, ok := taskComp["score"].(float64); ok && score < 60 {
			recommendations = append(recommendations, "任务负载过高，建议调整任务分配策略")
		}
	}

	return recommendations
}

// OptimizeDeviceLoad 优化设备负载
func (s *sDevice) OptimizeDeviceLoad(ctx context.Context, input *device.OptimizeDeviceLoadInput) (*device.OptimizeDeviceLoadOutput, error) {
	// 获取当前负载评分
	currentScore, _, err := s.getCurrentLoadScore(ctx, input.DeviceID)
	if err != nil {
		return nil, err
	}

	// 生成优化动作
	actions := s.generateOptimizationActions(input.Strategy, currentScore)

	// 计算优化后的评分
	optimizedScore := s.calculateOptimizedScore(currentScore, actions)

	return &device.OptimizeDeviceLoadOutput{
		DeviceID:       input.DeviceID,
		Strategy:       input.Strategy,
		CurrentScore:   currentScore,
		OptimizedScore: optimizedScore,
		Improvement:    optimizedScore - currentScore,
		Actions:        actions,
		Timestamp:      gtime.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

// generateOptimizationActions 生成优化动作
func (s *sDevice) generateOptimizationActions(strategy string, currentScore float64) []string {
	var actions []string

	switch strategy {
	case "aggressive":
		if currentScore < 50 {
			actions = append(actions, "立即停止非关键任务")
			actions = append(actions, "重启设备服务")
			actions = append(actions, "清理系统缓存")
		}
	case "moderate":
		if currentScore < 70 {
			actions = append(actions, "调整任务优先级")
			actions = append(actions, "优化资源分配")
		}
	case "conservative":
		if currentScore < 80 {
			actions = append(actions, "监控负载变化")
			actions = append(actions, "准备扩容方案")
		}
	default:
		actions = append(actions, "分析负载原因")
		actions = append(actions, "制定优化计划")
	}

	return actions
}

// calculateOptimizedScore 计算优化后的评分
func (s *sDevice) calculateOptimizedScore(currentScore float64, actions []string) float64 {
	improvement := float64(len(actions)) * 5 // 每个动作改善5分
	optimizedScore := currentScore + improvement
	return math.Min(optimizedScore, 100)
}

// SetDeviceLoadThreshold 设置设备负载阈值
func (s *sDevice) SetDeviceLoadThreshold(ctx context.Context, input *device.SetDeviceLoadThresholdInput) (*device.SetDeviceLoadThresholdOutput, error) {
	// 验证设备是否存在
	if err := s.validateDeviceExists(ctx, input.DeviceID); err != nil {
		return &device.SetDeviceLoadThresholdOutput{
			DeviceID:  input.DeviceID,
			Message:   "设备不存在",
			IsSuccess: false,
		}, err
	}

	// 验证阈值参数
	if input.Warning <= 0 || input.Critical <= 0 || input.MaxTasks <= 0 {
		return &device.SetDeviceLoadThresholdOutput{
			DeviceID:  input.DeviceID,
			Message:   "阈值参数无效",
			IsSuccess: false,
		}, gerror.New("阈值参数无效")
	}

	if input.Warning >= input.Critical {
		return &device.SetDeviceLoadThresholdOutput{
			DeviceID:  input.DeviceID,
			Message:   "警告阈值必须小于严重阈值",
			IsSuccess: false,
		}, gerror.New("警告阈值必须小于严重阈值")
	}

	// 这里应该保存到数据库
	// 目前只是返回成功消息

	return &device.SetDeviceLoadThresholdOutput{
		DeviceID:  input.DeviceID,
		Message:   "负载阈值设置成功",
		IsSuccess: true,
	}, nil
}

// GetDeviceLoadThreshold 获取设备负载阈值
func (s *sDevice) GetDeviceLoadThreshold(ctx context.Context, input *device.GetDeviceLoadThresholdInput) (*device.GetDeviceLoadThresholdOutput, error) {
	// 验证设备是否存在
	if err := s.validateDeviceExists(ctx, input.DeviceID); err != nil {
		return nil, err
	}

	// 这里应该从数据库获取阈值
	// 目前返回默认值
	return &device.GetDeviceLoadThresholdOutput{
		DeviceID: input.DeviceID,
		Warning:  70.0, // 默认警告阈值70%
		Critical: 90.0, // 默认严重阈值90%
		MaxTasks: 5,    // 默认最大任务数5
	}, nil
}

func (s *sDevice) getCurrentLoadScore(ctx context.Context, deviceId string) (float64, map[string]interface{}, error) {
	result, err := s.CalculateDeviceLoadScore(ctx, &device.CalculateDeviceLoadScoreInput{DeviceID: deviceId})
	if err != nil {
		return 0, nil, err
	}

	return result.LoadScore, result.Components, nil
}

func (s *sDevice) getDeviceData(ctx context.Context, deviceId string) (map[string]interface{}, error) {
	// 这里应该从数据库获取设备数据
	// 目前返回模拟数据
	return map[string]interface{}{
		"device_id": deviceId,
		"status":    "online",
	}, nil
}

func (s *sDevice) getDeviceMetrics(ctx context.Context, deviceId string) (map[string]interface{}, error) {
	// 这里应该从监控系统获取设备指标
	// 目前返回模拟数据
	return map[string]interface{}{
		"cpu_usage":            45.5,
		"memory_usage":         62.3,
		"disk_usage":           78.9,
		"network_latency":      25.0,
		"running_task_count":   2,
		"max_concurrent_tasks": 5,
	}, nil
}

func (s *sDevice) validateDeviceExists(ctx context.Context, deviceId string) error {
	// 这里应该检查设备是否存在
	// 目前只是简单验证
	if deviceId == "" {
		return gerror.New("设备ID不能为空")
	}
	return nil
}
