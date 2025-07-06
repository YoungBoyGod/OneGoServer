package device

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"

	device "OneGfServer/internal/model/device"
)

// ===============================
// 设备监控相关业务逻辑
// ===============================

// GetDeviceLoadMetrics 获取设备负载指标
func (s *sDevice) GetDeviceLoadMetrics(ctx context.Context, input *device.GetDeviceLoadMetricsInput) (*device.GetDeviceLoadMetricsOutput, error) {
	// 验证设备是否存在
	if err := s.validateDeviceExists(ctx, input.DeviceID); err != nil {
		return nil, err
	}

	// 解析时间周期
	duration, err := s.parsePeriod(input.Period)
	if err != nil {
		return nil, err
	}

	// 获取历史数据
	history, err := s.getDeviceLoadHistory(ctx, input.DeviceID, duration)
	if err != nil {
		return nil, err
	}

	// 生成时间序列数据
	timeSeriesData := s.generateTimeSeriesData(
		gtime.Now().Add(-duration).Time,
		gtime.Now().Time,
		300, // 5分钟间隔
		0,   // 最小值
		100, // 最大值
	)

	return &device.GetDeviceLoadMetricsOutput{
		DeviceID:   input.DeviceID,
		Period:     input.Period,
		Duration:   duration.String(),
		TimeSeries: timeSeriesData,
		History:    history,
		Summary:    s.calculateMetricsSummary(history),
		Timestamp:  gtime.Now().Format("2006-01-02 15:04:05"),
	}, nil
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
	case "30d":
		return 30 * 24 * time.Hour, nil
	default:
		return time.Hour, gerror.New("不支持的时间周期")
	}
}

// getDeviceLoadHistory 获取设备负载历史
func (s *sDevice) getDeviceLoadHistory(ctx context.Context, deviceId string, duration time.Duration) (map[string]interface{}, error) {
	// 这里应该从数据库获取历史数据
	// 目前返回模拟数据
	return map[string]interface{}{
		"cpu_history":     []float64{45, 52, 48, 61, 55, 49, 58},
		"memory_history":  []float64{62, 65, 68, 71, 69, 66, 70},
		"disk_history":    []float64{78, 79, 80, 81, 82, 83, 84},
		"network_history": []float64{25, 28, 22, 30, 26, 24, 29},
	}, nil
}

// generateTimeSeriesData 生成时间序列数据
func (s *sDevice) generateTimeSeriesData(startTime, endTime time.Time, intervalSeconds int, minValue, maxValue float64) []map[string]interface{} {
	var data []map[string]interface{}
	currentTime := startTime

	for currentTime.Before(endTime) {
		data = append(data, map[string]interface{}{
			"timestamp":       currentTime.Format("2006-01-02 15:04:05"),
			"cpu_usage":       s.randomFloat()*(maxValue-minValue) + minValue,
			"memory_usage":    s.randomFloat()*(maxValue-minValue) + minValue,
			"disk_usage":      s.randomFloat()*(maxValue-minValue) + minValue,
			"network_latency": s.randomFloat()*(maxValue-minValue) + minValue,
		})
		currentTime = currentTime.Add(time.Duration(intervalSeconds) * time.Second)
	}

	return data
}

// calculateMetricsSummary 计算指标摘要
func (s *sDevice) calculateMetricsSummary(history map[string]interface{}) map[string]interface{} {
	summary := make(map[string]interface{})

	// CPU摘要
	if cpuHistory, ok := history["cpu_history"].([]float64); ok {
		summary["cpu"] = map[string]interface{}{
			"avg": s.calculateAverage(cpuHistory),
			"max": s.calculateMax(cpuHistory),
			"min": s.calculateMin(cpuHistory),
		}
	}

	// 内存摘要
	if memHistory, ok := history["memory_history"].([]float64); ok {
		summary["memory"] = map[string]interface{}{
			"avg": s.calculateAverage(memHistory),
			"max": s.calculateMax(memHistory),
			"min": s.calculateMin(memHistory),
		}
	}

	// 磁盘摘要
	if diskHistory, ok := history["disk_history"].([]float64); ok {
		summary["disk"] = map[string]interface{}{
			"avg": s.calculateAverage(diskHistory),
			"max": s.calculateMax(diskHistory),
			"min": s.calculateMin(diskHistory),
		}
	}

	// 网络摘要
	if netHistory, ok := history["network_history"].([]float64); ok {
		summary["network"] = map[string]interface{}{
			"avg": s.calculateAverage(netHistory),
			"max": s.calculateMax(netHistory),
			"min": s.calculateMin(netHistory),
		}
	}

	return summary
}

// calculateAverage 计算平均值
func (s *sDevice) calculateAverage(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

// calculateMax 计算最大值
func (s *sDevice) calculateMax(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	max := values[0]
	for _, v := range values {
		if v > max {
			max = v
		}
	}
	return max
}

// calculateMin 计算最小值
func (s *sDevice) calculateMin(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	min := values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
	}
	return min
}

// randomFloat 生成随机浮点数
func (s *sDevice) randomFloat() float64 {
	// 这里应该使用更好的随机数生成器
	// 目前返回固定值用于演示
	return 0.5
}
