package device

import (
	"context"
	"math"
	"time"

	"github.com/gogf/gf/v2/os/gtime"

	"OneGfServer/internal/model/device"
)

// ===============================
// 设备性能分析相关业务逻辑
// ===============================

// AnalyzeDevicePerformance 分析设备性能
func (s *sDevice) AnalyzeDevicePerformance(ctx context.Context, input *device.AnalyzeDevicePerformanceInput) (*device.AnalyzeDevicePerformanceOutput, error) {
	/*
		if len(input.PerformanceData) == 0 {
			return &device.AnalyzeDevicePerformanceOutput{
				BasicStats:       nil,
				Trends:           nil,
				Bottlenecks:      nil,
				PerformanceScore: 0,
				Recommendations:  []string{"没有性能数据可供分析"},
				TimeRange:        nil,
			}, nil
		}
	*/

	// 1. 基础统计信息
	basicStatsInput := &device.CalculateBasicStatsInput{
		PerformanceData: input.PerformanceData,
	}
	basicStatsOutput := s.calculateBasicStats(basicStatsInput)

	// 2. 性能趋势分析
	trendsInput := &device.AnalyzePerformanceTrendsInput{
		PerformanceData: input.PerformanceData,
	}
	trendsOutput := s.analyzePerformanceTrends(trendsInput)

	// 3. 瓶颈识别
	bottlenecksInput := &device.IdentifyBottlenecksInput{
		PerformanceData: input.PerformanceData,
	}
	bottlenecksOutput := s.identifyBottlenecks(bottlenecksInput)

	// 4. 性能评分
	performanceScoreInput := &device.CalculatePerformanceScoreFromDataInput{
		PerformanceData: input.PerformanceData,
	}
	performanceScoreOutput := s.calculatePerformanceScoreFromData(performanceScoreInput)

	// 5. 优化建议
	recommendationsInput := &device.GeneratePerformanceRecommendationsInput{
		Analysis: map[string]interface{}{
			"basic_stats":       basicStatsOutput.Stats,
			"trends":            trendsOutput.Trends,
			"bottlenecks":       bottlenecksOutput.Bottlenecks,
			"performance_score": performanceScoreOutput.Score,
		},
	}
	recommendationsOutput := s.generatePerformanceRecommendations(recommendationsInput)

	// 6. 时间范围
	earliestInput := &device.GetEarliestTimestampInput{
		Data: input.PerformanceData,
	}
	earliestOutput := s.getEarliestTimestamp(earliestInput)

	latestInput := &device.GetLatestTimestampInput{
		Data: input.PerformanceData,
	}
	latestOutput := s.getLatestTimestamp(latestInput)

	timeRange := map[string]interface{}{
		"start": earliestOutput.Timestamp,
		"end":   latestOutput.Timestamp,
		"count": len(input.PerformanceData),
	}

	return &device.AnalyzeDevicePerformanceOutput{
		BasicStats:       basicStatsOutput.Stats,
		Trends:           trendsOutput.Trends,
		Bottlenecks:      bottlenecksOutput.Bottlenecks,
		PerformanceScore: performanceScoreOutput.Score,
		Recommendations:  recommendationsOutput.Recommendations,
		TimeRange:        timeRange,
	}, nil
}

// calculateBasicStats 计算基础统计信息
func (s *sDevice) calculateBasicStats(input *device.CalculateBasicStatsInput) *device.CalculateBasicStatsOutput {
	stats := make(map[string]interface{})
	/*
		// CPU统计
		cpuValuesOutput := s.extractValues(&device.ExtractValuesInput{
			Data: input.PerformanceData,
			Key:  "cpu_usage",
		})
		stats["cpu"] = map[string]interface{}{
			"avg": s.calculateAverage(cpuValuesOutput.Values),
			"max": s.calculateMax(cpuValuesOutput.Values),
			"min": s.calculateMin(cpuValuesOutput.Values),
			"std": s.calculateStandardDeviation(&device.CalculateStandardDeviationInput{Values: cpuValuesOutput.Values}).StandardDeviation,
			"p95": s.calculatePercentile(&device.CalculatePercentileInput{Values: cpuValuesOutput.Values, Percentile: 95}).PercentileValue,
			"p99": s.calculatePercentile(&device.CalculatePercentileInput{Values: cpuValuesOutput.Values, Percentile: 99}).PercentileValue,
		}

		// 内存统计
		memValuesOutput := s.extractValues(&device.ExtractValuesInput{
			Data: input.PerformanceData,
			Key:  "memory_usage",
		})
		stats["memory"] = map[string]interface{}{
			"avg": s.calculateAverage(memValuesOutput.Values),
			"max": s.calculateMax(memValuesOutput.Values),
			"min": s.calculateMin(memValuesOutput.Values),
			"std": s.calculateStandardDeviation(&device.CalculateStandardDeviationInput{Values: memValuesOutput.Values}).StandardDeviation,
			"p95": s.calculatePercentile(&device.CalculatePercentileInput{Values: memValuesOutput.Values, Percentile: 95}).PercentileValue,
			"p99": s.calculatePercentile(&device.CalculatePercentileInput{Values: memValuesOutput.Values, Percentile: 99}).PercentileValue,
		}

		// 磁盘统计
		diskValuesOutput := s.extractValues(&device.ExtractValuesInput{
			Data: input.PerformanceData,
			Key:  "disk_usage",
		})
		stats["disk"] = map[string]interface{}{
			"avg": s.calculateAverage(diskValuesOutput.Values),
			"max": s.calculateMax(diskValuesOutput.Values),
			"min": s.calculateMin(diskValuesOutput.Values),
			"std": s.calculateStandardDeviation(&device.CalculateStandardDeviationInput{Values: diskValuesOutput.Values}).StandardDeviation,
			"p95": s.calculatePercentile(&device.CalculatePercentileInput{Values: diskValuesOutput.Values, Percentile: 95}).PercentileValue,
			"p99": s.calculatePercentile(&device.CalculatePercentileInput{Values: diskValuesOutput.Values, Percentile: 99}).PercentileValue,
		}

		// 网络统计
		netValuesOutput := s.extractValues(&device.ExtractValuesInput{
			Data: input.PerformanceData,
			Key:  "network_latency",
		})
		stats["network"] = map[string]interface{}{
			"avg": s.calculateAverage(netValuesOutput.Values),
			"max": s.calculateMax(netValuesOutput.Values),
			"min": s.calculateMin(netValuesOutput.Values),
			"std": s.calculateStandardDeviation(&device.CalculateStandardDeviationInput{Values: netValuesOutput.Values}).StandardDeviation,
			"p95": s.calculatePercentile(&device.CalculatePercentileInput{Values: netValuesOutput.Values, Percentile: 95}).PercentileValue,
			"p99": s.calculatePercentile(&device.CalculatePercentileInput{Values: netValuesOutput.Values, Percentile: 99}).PercentileValue,
		}
	*/
	return &device.CalculateBasicStatsOutput{
		Stats: stats,
	}
}

// analyzePerformanceTrends 分析性能趋势
func (s *sDevice) analyzePerformanceTrends(input *device.AnalyzePerformanceTrendsInput) *device.AnalyzePerformanceTrendsOutput {
	trends := make(map[string]interface{})
	/*
		// 按时间排序数据
		sortInput := &device.SortByTimestampInput{
			Data: input.PerformanceData,
		}
		sortOutput := s.sortByTimestamp(sortInput)

		// 分析CPU趋势
		cpuTrendInput := &device.AnalyzeMetricTrendInput{
			SortedData: sortOutput.SortedData,
			MetricKey:  "cpu_usage",
		}
		cpuTrendOutput := s.analyzeMetricTrend(cpuTrendInput)
		trends["cpu"] = map[string]interface{}{
			"trend":      cpuTrendOutput.Trend,
			"slope":      cpuTrendOutput.Slope,
			"volatility": cpuTrendOutput.Volatility,
		}

		// 分析内存趋势
		memTrendInput := &device.AnalyzeMetricTrendInput{
			SortedData: sortOutput.SortedData,
			MetricKey:  "memory_usage",
		}
		memTrendOutput := s.analyzeMetricTrend(memTrendInput)
		trends["memory"] = map[string]interface{}{
			"trend":      memTrendOutput.Trend,
			"slope":      memTrendOutput.Slope,
			"volatility": memTrendOutput.Volatility,
		}

		// 分析磁盘趋势
		diskTrendInput := &device.AnalyzeMetricTrendInput{
			SortedData: sortOutput.SortedData,
			MetricKey:  "disk_usage",
		}
		diskTrendOutput := s.analyzeMetricTrend(diskTrendInput)
		trends["disk"] = map[string]interface{}{
			"trend":      diskTrendOutput.Trend,
			"slope":      diskTrendOutput.Slope,
			"volatility": diskTrendOutput.Volatility,
		}

		// 分析网络趋势
		netTrendInput := &device.AnalyzeMetricTrendInput{
			SortedData: sortOutput.SortedData,
			MetricKey:  "network_latency",
		}
		netTrendOutput := s.analyzeMetricTrend(netTrendInput)
		trends["network"] = map[string]interface{}{
			"trend":      netTrendOutput.Trend,
			"slope":      netTrendOutput.Slope,
			"volatility": netTrendOutput.Volatility,
		}
	*/
	return &device.AnalyzePerformanceTrendsOutput{
		Trends: trends,
	}
}

// analyzeMetricTrend 分析单个指标趋势
func (s *sDevice) analyzeMetricTrend(input *device.AnalyzeMetricTrendInput) *device.AnalyzeMetricTrendOutput {
	/*
		extractInput := &device.ExtractValuesInput{
			Data: input.SortedData,
			Key:  input.MetricKey,
		}
		valuesOutput := s.extractValues(extractInput)
		values := valuesOutput.Values

		if len(values) < 2 {
			return &device.AnalyzeMetricTrendOutput{
				Trend:      "insufficient_data",
				Slope:      0,
				Volatility: 0,
			}
		}

		// 计算线性回归斜率
		slope := s.calculateLinearRegressionSlope(values)

		// 判断趋势
		var trend string
		if slope > 0.1 {
			trend = "increasing"
		} else if slope < -0.1 {
			trend = "decreasing"
		} else {
			trend = "stable"
		}

		volatilityInput := &device.CalculateVolatilityInput{
			Values: values,
		}
		volatilityOutput := s.calculateVolatility(volatilityInput)
	*/
	return &device.AnalyzeMetricTrendOutput{}
}

// identifyBottlenecks 识别瓶颈
func (s *sDevice) identifyBottlenecks(input *device.IdentifyBottlenecksInput) *device.IdentifyBottlenecksOutput {
	var bottlenecks []map[string]interface{}

	// 检查CPU瓶颈
	cpuInput := &device.CheckCPUBottleneckInput{
		PerformanceData: input.PerformanceData,
	}
	cpuOutput := s.checkCPUBottleneck(cpuInput)
	if cpuOutput.HasIssue {
		bottlenecks = append(bottlenecks, cpuOutput.Bottleneck)
	}

	// 检查内存瓶颈
	memInput := &device.CheckMemoryBottleneckInput{
		PerformanceData: input.PerformanceData,
	}
	memOutput := s.checkMemoryBottleneck(memInput)
	if memOutput.HasIssue {
		bottlenecks = append(bottlenecks, memOutput.Bottleneck)
	}

	// 检查磁盘瓶颈
	diskInput := &device.CheckDiskBottleneckInput{
		PerformanceData: input.PerformanceData,
	}
	diskOutput := s.checkDiskBottleneck(diskInput)
	if diskOutput.HasIssue {
		bottlenecks = append(bottlenecks, diskOutput.Bottleneck)
	}

	// 检查网络瓶颈
	netInput := &device.CheckNetworkBottleneckInput{
		PerformanceData: input.PerformanceData,
	}
	netOutput := s.checkNetworkBottleneck(netInput)
	if netOutput.HasIssue {
		bottlenecks = append(bottlenecks, netOutput.Bottleneck)
	}

	return &device.IdentifyBottlenecksOutput{
		Bottlenecks: bottlenecks,
	}
}

// checkCPUBottleneck 检查CPU瓶颈
func (s *sDevice) checkCPUBottleneck(input *device.CheckCPUBottleneckInput) *device.CheckCPUBottleneckOutput {
	/*
		extractInput := &device.ExtractValuesInput{
			Data: input.PerformanceData,
			Key:  "cpu_usage",
		}
		cpuValuesOutput := s.extractValues(extractInput)
		cpuValues := cpuValuesOutput.Values

		avgCPU := s.calculateAverage(cpuValues)
		p95Input := &device.CalculatePercentileInput{
			Values:     cpuValues,
			Percentile: 95,
		}
		p95Output := s.calculatePercentile(p95Input)
		p95CPU := p95Output.PercentileValue

		if avgCPU > 80 || p95CPU > 95 {
			severityInput := &device.GetBottleneckSeverityInput{
				Usage: avgCPU,
			}
			severityOutput := s.getBottleneckSeverity(severityInput)

			return &device.CheckCPUBottleneckOutput{
				HasIssue: true,
				Bottleneck: map[string]interface{}{
					"type":           "cpu",
					"severity":       severityOutput.Severity,
					"avg_usage":      avgCPU,
					"p95_usage":      p95CPU,
					"description":    "CPU使用率过高，可能影响系统性能",
					"recommendation": "考虑优化CPU密集型任务或增加CPU资源",
				},
			}
		}
	*/
	return &device.CheckCPUBottleneckOutput{}
}

// checkMemoryBottleneck 检查内存瓶颈
func (s *sDevice) checkMemoryBottleneck(input *device.CheckMemoryBottleneckInput) *device.CheckMemoryBottleneckOutput {
	/*
		memValuesOutput := s.extractValues(&device.ExtractValuesInput{
			Data: input.PerformanceData,
			Key:  "memory_usage",
		})
		memValues := memValuesOutput.Values
		avgMem := s.calculateAverage(memValues)
		p95Mem := s.calculatePercentile(&device.CalculatePercentileInput{Values: memValues, Percentile: 95}).PercentileValue

		if avgMem > 85 || p95Mem > 95 {
			severity := s.getBottleneckSeverity(&device.GetBottleneckSeverityInput{Usage: avgMem}).Severity
			return &device.CheckMemoryBottleneckOutput{
				HasIssue: true,
				Bottleneck: map[string]interface{}{
					"type":           "memory",
					"severity":       severity,
					"avg_usage":      avgMem,
					"p95_usage":      p95Mem,
					"description":    "内存使用率过高，可能导致系统不稳定",
					"recommendation": "考虑增加内存或优化内存使用",
				},
			}
		}
	*/
	return &device.CheckMemoryBottleneckOutput{
		HasIssue:   false,
		Bottleneck: nil,
	}
}

// checkDiskBottleneck 检查磁盘瓶颈
func (s *sDevice) checkDiskBottleneck(input *device.CheckDiskBottleneckInput) *device.CheckDiskBottleneckOutput {
	/*
		diskValuesOutput := s.extractValues(&device.ExtractValuesInput{
			Data: input.PerformanceData,
			Key:  "disk_usage",
		})
		diskValues := diskValuesOutput.Values
		avgDisk := s.calculateAverage(diskValues)
		p95Disk := s.calculatePercentile(&device.CalculatePercentileInput{Values: diskValues, Percentile: 95}).PercentileValue

		if avgDisk > 90 || p95Disk > 98 {
			severity := s.getBottleneckSeverity(&device.GetBottleneckSeverityInput{Usage: avgDisk}).Severity
			return &device.CheckDiskBottleneckOutput{
				HasIssue: true,
				Bottleneck: map[string]interface{}{
					"type":           "disk",
					"severity":       severity,
					"avg_usage":      avgDisk,
					"p95_usage":      p95Disk,
					"description":    "磁盘使用率过高，可能影响I/O性能",
					"recommendation": "考虑清理磁盘空间或扩容",
				},
			}
		}
	*/
	return &device.CheckDiskBottleneckOutput{
		HasIssue:   false,
		Bottleneck: nil,
	}
}

// checkNetworkBottleneck 检查网络瓶颈
func (s *sDevice) checkNetworkBottleneck(input *device.CheckNetworkBottleneckInput) *device.CheckNetworkBottleneckOutput {
	/*
		netValuesOutput := s.extractValues(&device.ExtractValuesInput{
			Data: input.PerformanceData,
			Key:  "network_latency",
		})
		netValues := netValuesOutput.Values
		avgNet := s.calculateAverage(netValues)
		p95Net := s.calculatePercentile(&device.CalculatePercentileInput{Values: netValues, Percentile: 95}).PercentileValue

		if avgNet > 100 || p95Net > 200 {
			severity := s.getBottleneckSeverity(&device.GetBottleneckSeverityInput{Usage: avgNet / 10}).Severity
			return &device.CheckNetworkBottleneckOutput{
				HasIssue: true,
				Bottleneck: map[string]interface{}{
					"type":           "network",
					"severity":       severity,
					"avg_latency":    avgNet,
					"p95_latency":    p95Net,
					"description":    "网络延迟过高，可能影响响应时间",
					"recommendation": "检查网络连接或优化网络配置",
				},
			}
		}
	*/
	return &device.CheckNetworkBottleneckOutput{
		HasIssue:   false,
		Bottleneck: nil,
	}
}

// calculatePerformanceScoreFromData 计算性能评分
func (s *sDevice) calculatePerformanceScoreFromData(input *device.CalculatePerformanceScoreFromDataInput) *device.CalculatePerformanceScoreFromDataOutput {
	/*
		score := 100.0

		// CPU评分 (权重: 30%)
		extractInput := &device.ExtractValuesInput{
			Data: input.PerformanceData,
			Key:  "cpu_usage",
		}
		cpuValuesOutput := s.extractValues(extractInput)
		avgCPU := s.calculateAverage(cpuValuesOutput.Values)
		if avgCPU > 90 {
			score -= 30
		} else if avgCPU > 70 {
			score -= 15
		} else if avgCPU > 50 {
			score -= 5
		}

		// 内存评分 (权重: 25%)
		extractInput.Key = "memory_usage"
		memValuesOutput := s.extractValues(extractInput)
		avgMem := s.calculateAverage(memValuesOutput.Values)
		if avgMem > 90 {
			score -= 25
		} else if avgMem > 80 {
			score -= 12
		} else if avgMem > 60 {
			score -= 5
		}

		// 磁盘评分 (权重: 20%)
		extractInput.Key = "disk_usage"
		diskValuesOutput := s.extractValues(extractInput)
		avgDisk := s.calculateAverage(diskValuesOutput.Values)
		if avgDisk > 95 {
			score -= 20
		} else if avgDisk > 85 {
			score -= 10
		} else if avgDisk > 70 {
			score -= 3
		}

		// 网络评分 (权重: 15%)
		extractInput.Key = "network_latency"
		netValuesOutput := s.extractValues(extractInput)
		avgNet := s.calculateAverage(netValuesOutput.Values)
		if avgNet > 200 {
			score -= 15
		} else if avgNet > 100 {
			score -= 8
		} else if avgNet > 50 {
			score -= 3
		}

		// 稳定性评分 (权重: 10%)
		stabilityInput := &device.CalculateStabilityScoreInput{
			PerformanceData: input.PerformanceData,
		}
		stabilityOutput := s.calculateStabilityScore(stabilityInput)
		score += stabilityOutput.Score * 0.1
	*/
	return &device.CalculatePerformanceScoreFromDataOutput{}
}

// generatePerformanceRecommendations 生成性能优化建议
func (s *sDevice) generatePerformanceRecommendations(input *device.GeneratePerformanceRecommendationsInput) *device.GeneratePerformanceRecommendationsOutput {
	var recommendations []string

	// 基于性能评分生成建议
	if performanceScore, ok := input.Analysis["performance_score"].(float64); ok {
		if performanceScore < 50 {
			recommendations = append(recommendations, "性能严重不足，建议立即优化")
		} else if performanceScore < 70 {
			recommendations = append(recommendations, "性能偏低，建议进行优化")
		}
	}

	// 基于瓶颈分析生成建议
	if bottlenecks, ok := input.Analysis["bottlenecks"].([]map[string]interface{}); ok {
		for _, bottleneck := range bottlenecks {
			if recommendation, ok := bottleneck["recommendation"].(string); ok {
				recommendations = append(recommendations, recommendation)
			}
		}
	}

	// 基于趋势分析生成建议
	if trends, ok := input.Analysis["trends"].(map[string]interface{}); ok {
		if cpuTrend, ok := trends["cpu"].(map[string]interface{}); ok {
			if trend, ok := cpuTrend["trend"].(string); ok && trend == "increasing" {
				recommendations = append(recommendations, "CPU使用率呈上升趋势，建议监控并优化")
			}
		}
	}

	return &device.GeneratePerformanceRecommendationsOutput{
		Recommendations: recommendations,
	}
}

// 辅助方法
func (s *sDevice) extractValues(input *device.ExtractValuesInput) *device.ExtractValuesOutput {
	var values []float64
	for _, item := range input.Data {
		if value, ok := item[input.Key].(float64); ok {
			values = append(values, value)
		}
	}
	return &device.ExtractValuesOutput{
		Values: values,
	}
}

func (s *sDevice) sortByTimestamp(input *device.SortByTimestampInput) *device.SortByTimestampOutput {
	// 这里应该按时间戳排序
	// 目前返回原数据
	return &device.SortByTimestampOutput{
		SortedData: input.Data,
	}
}

func (s *sDevice) calculateLinearRegressionSlope(values []float64) float64 {
	if len(values) < 2 {
		return 0
	}

	n := float64(len(values))
	sumX := 0.0
	sumY := 0.0
	sumXY := 0.0
	sumX2 := 0.0

	for i, y := range values {
		x := float64(i)
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	slope := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
	return slope
}

func (s *sDevice) calculateVolatility(input *device.CalculateVolatilityInput) *device.CalculateVolatilityOutput {
	if len(input.Values) < 2 {
		return &device.CalculateVolatilityOutput{
			Volatility: 0,
		}
	}

	mean := s.calculateAverage(input.Values)
	variance := 0.0

	for _, value := range input.Values {
		variance += math.Pow(value-mean, 2)
	}

	variance /= float64(len(input.Values))
	return &device.CalculateVolatilityOutput{
		Volatility: math.Sqrt(variance),
	}
}

// calculateStandardDeviation 计算标准差
func (s *sDevice) calculateStandardDeviation(input *device.CalculateStandardDeviationInput) *device.CalculateStandardDeviationOutput {
	if len(input.Values) < 2 {
		return &device.CalculateStandardDeviationOutput{StandardDeviation: 0}
	}
	mean := s.calculateAverage(input.Values)
	variance := 0.0
	for _, value := range input.Values {
		variance += math.Pow(value-mean, 2)
	}
	variance /= float64(len(input.Values))
	return &device.CalculateStandardDeviationOutput{
		StandardDeviation: math.Sqrt(variance),
	}
}

func (s *sDevice) calculatePercentile(input *device.CalculatePercentileInput) *device.CalculatePercentileOutput {
	values := input.Values
	percentile := input.Percentile
	if len(values) == 0 {
		return &device.CalculatePercentileOutput{PercentileValue: 0}
	}
	index := int(float64(len(values)-1) * float64(percentile) / 100.0)
	if index >= len(values) {
		index = len(values) - 1
	}
	return &device.CalculatePercentileOutput{PercentileValue: values[index]}
}

func (s *sDevice) calculateStabilityScore(input *device.CalculateStabilityScoreInput) *device.CalculateStabilityScoreOutput {
	// 计算性能稳定性评分
	// 基于各指标的变异系数
	score := 100.0

	extractInput := &device.ExtractValuesInput{
		Data: input.PerformanceData,
		Key:  "cpu_usage",
	}
	cpuValuesOutput := s.extractValues(extractInput)
	cpuValues := cpuValuesOutput.Values

	if len(cpuValues) > 0 {
		volatilityInput := &device.CalculateVolatilityInput{
			Values: cpuValues,
		}
		volatilityOutput := s.calculateVolatility(volatilityInput)
		cv := volatilityOutput.Volatility / s.calculateAverage(cpuValues)
		if cv > 0.5 {
			score -= 20
		} else if cv > 0.3 {
			score -= 10
		}
	}

	return &device.CalculateStabilityScoreOutput{
		Score: math.Max(0, score),
	}
}

func (s *sDevice) getBottleneckSeverity(input *device.GetBottleneckSeverityInput) *device.GetBottleneckSeverityOutput {
	usage := input.Usage
	severity := "low"
	if usage > 95 {
		severity = "critical"
	} else if usage > 85 {
		severity = "high"
	} else if usage > 70 {
		severity = "medium"
	}
	return &device.GetBottleneckSeverityOutput{
		Severity: severity,
	}
}

func (s *sDevice) getEarliestTimestamp(input *device.GetEarliestTimestampInput) *device.GetEarliestTimestampOutput {
	// 这里应该返回最早的时间戳
	return &device.GetEarliestTimestampOutput{
		Timestamp: gtime.Now().Add(-24 * time.Hour).Format("2006-01-02 15:04:05"),
	}
}

func (s *sDevice) getLatestTimestamp(input *device.GetLatestTimestampInput) *device.GetLatestTimestampOutput {
	// 这里应该返回最晚的时间戳
	return &device.GetLatestTimestampOutput{
		Timestamp: gtime.Now().Format("2006-01-02 15:04:05"),
	}
}
