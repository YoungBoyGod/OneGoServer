package system

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	system "OneGfServer/internal/model/system"
)

// ===============================
// 系统性能业务逻辑
// ===============================

// GetSystemPerformance 获取系统性能
func (s *sSystem) GetSystemPerformance(ctx context.Context, input *system.GetSystemPerformanceInput) (*system.GetSystemPerformanceOutput, error) {
	// 解析时间周期
	parseInput := &system.ParsePeriodInput{Period: input.Period}
	parseOutput := s.parsePeriod(parseInput)
	if parseOutput.Duration == "" {
		return nil, gerror.NewCode(gcode.CodeValidationFailed, "不支持的时间周期")
	}

	// 获取性能数据
	calculateInput := &system.CalculateSystemPerformanceInput{Duration: parseOutput.Duration}
	performance := s.calculateSystemPerformance(calculateInput)

	return &system.GetSystemPerformanceOutput{
		Period:      input.Period,
		Performance: performance.Performance,
	}, nil
}

// calculateSystemPerformance 计算系统性能
func (s *sSystem) calculateSystemPerformance(input *system.CalculateSystemPerformanceInput) *system.CalculateSystemPerformanceOutput {
	// 这里应该基于实际数据计算性能指标
	// 目前返回模拟数据
	return &system.CalculateSystemPerformanceOutput{
		Performance: map[string]interface{}{
			"apiLatency": map[string]interface{}{
				"avg": 150.5,
				"p95": 300.0,
				"p99": 500.0,
				"max": 1000.0,
			},
			"throughput": map[string]interface{}{
				"requestsPerSecond": 500.0,
				"totalRequests":     1000000,
			},
			"errorRate": map[string]interface{}{
				"rate":        0.02,
				"totalErrors": 20000,
			},
			"resourceUsage": map[string]interface{}{
				"cpuUsage":    45.5,
				"memoryUsage": 60.2,
				"diskUsage":   75.8,
			},
		},
	}
}
