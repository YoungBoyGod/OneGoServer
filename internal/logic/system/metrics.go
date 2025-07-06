package system

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	system "OneGfServer/internal/model/system"
)

// ===============================
// 系统指标业务逻辑
// ===============================

// GetSystemMetrics 获取系统指标
func (s *sSystem) GetSystemMetrics(ctx context.Context, input *system.GetSystemMetricsInput) (*system.GetSystemMetricsOutput, error) {
	// 解析时间周期
	parseInput := &system.ParsePeriodInput{Period: input.Period}
	parseOutput := s.parsePeriod(parseInput)
	if parseOutput.Duration == "" {
		return nil, gerror.NewCode(gcode.CodeValidationFailed, "不支持的时间周期")
	}

	// 获取历史指标数据
	metrics, err := s.getSystemMetricsHistory(ctx, input.Period)
	if err != nil {
		return nil, err
	}

	return &system.GetSystemMetricsOutput{
		Period:  input.Period,
		Metrics: metrics,
	}, nil
}

// getSystemMetricsHistory 获取系统指标历史
func (s *sSystem) getSystemMetricsHistory(ctx context.Context, period string) (map[string]interface{}, error) {
	// 这里应该从数据库或监控系统获取历史数据
	// 目前返回模拟数据
	endTime := time.Now()
	var startTime time.Time

	switch period {
	case "1h":
		startTime = endTime.Add(-time.Hour)
	case "6h":
		startTime = endTime.Add(-6 * time.Hour)
	case "24h":
		startTime = endTime.Add(-24 * time.Hour)
	case "7d":
		startTime = endTime.Add(-7 * 24 * time.Hour)
	default:
		startTime = endTime.Add(-time.Hour)
	}

	generateInput := &system.GenerateTimeSeriesDataInput{
		StartTime:       startTime.Format("2006-01-02 15:04:05"),
		EndTime:         endTime.Format("2006-01-02 15:04:05"),
		IntervalSeconds: 60,
		MinValue:        0,
		MaxValue:        100,
	}

	metrics := map[string]interface{}{
		"cpuUsage":          s.generateTimeSeriesData(generateInput).Data,
		"memoryUsage":       s.generateTimeSeriesData(generateInput).Data,
		"diskUsage":         s.generateTimeSeriesData(generateInput).Data,
		"networkIO":         s.generateTimeSeriesData(generateInput).Data,
		"activeConnections": s.generateTimeSeriesData(generateInput).Data,
		"requestRate":       s.generateTimeSeriesData(generateInput).Data,
		"errorRate":         s.generateTimeSeriesData(generateInput).Data,
		"responseTime":      s.generateTimeSeriesData(generateInput).Data,
	}

	return metrics, nil
}

// parsePeriod 解析时间周期
func (s *sSystem) parsePeriod(input *system.ParsePeriodInput) *system.ParsePeriodOutput {
	switch input.Period {
	case "1h":
		return &system.ParsePeriodOutput{Duration: "1h"}
	case "6h":
		return &system.ParsePeriodOutput{Duration: "6h"}
	case "24h":
		return &system.ParsePeriodOutput{Duration: "24h"}
	case "7d":
		return &system.ParsePeriodOutput{Duration: "7d"}
	default:
		return &system.ParsePeriodOutput{Duration: ""}
	}
}

// generateTimeSeriesData 生成时间序列数据
func (s *sSystem) generateTimeSeriesData(input *system.GenerateTimeSeriesDataInput) *system.GenerateTimeSeriesDataOutput {
	var data []map[string]interface{}
	startTime, _ := time.Parse("2006-01-02 15:04:05", input.StartTime)
	endTime, _ := time.Parse("2006-01-02 15:04:05", input.EndTime)
	currentTime := startTime

	for currentTime.Before(endTime) {
		// 生成随机值
		value := input.MinValue + (input.MaxValue-input.MinValue)*s.randomFloat()

		data = append(data, map[string]interface{}{
			"timestamp": currentTime.Unix(),
			"value":     value,
		})

		currentTime = currentTime.Add(time.Duration(input.IntervalSeconds) * time.Second)
	}

	return &system.GenerateTimeSeriesDataOutput{
		Data: data,
	}
}

// randomFloat 生成随机浮点数
func (s *sSystem) randomFloat() float64 {
	// 简单的随机数生成
	return float64(time.Now().UnixNano()%100) / 100.0
}
