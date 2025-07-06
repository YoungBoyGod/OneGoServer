package device

import (
	"context"
	"math"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
)

// ===============================
// 设备性能分析相关业务逻辑
// ===============================

// AnalyzeDevicePerformance 分析设备性能
func (s *sDevice) AnalyzeDevicePerformance(ctx context.Context, performanceData []map[string]interface{}) map[string]interface{} {
	if len(performanceData) == 0 {
		return map[string]interface{}{
			"error": "没有性能数据可供分析",
		}
	}

	analysis := make(map[string]interface{})

	// 1. 基础统计信息
	analysis["basic_stats"] = s.calculateBasicStats(performanceData)

	// 2. 性能趋势分析
	analysis["trends"] = s.analyzePerformanceTrends(performanceData)

	// 3. 瓶颈识别
	analysis["bottlenecks"] = s.identifyBottlenecks(performanceData)

	// 4. 性能评分
	analysis["performance_score"] = s.calculatePerformanceScore(performanceData)

	// 5. 优化建议
	analysis["recommendations"] = s.generatePerformanceRecommendations(analysis)

	// 6. 时间范围
	analysis["time_range"] = map[string]interface{}{
		"start": s.getEarliestTimestamp(performanceData),
		"end":   s.getLatestTimestamp(performanceData),
		"count": len(performanceData),
	}

	return analysis
}

// calculateBasicStats 计算基础统计信息
func (s *sDevice) calculateBasicStats(performanceData []map[string]interface{}) map[string]interface{} {
	stats := make(map[string]interface{})

	// CPU统计
	cpuValues := s.extractValues(performanceData, "cpu_usage")
	stats["cpu"] = map[string]interface{}{
		"avg": s.calculateAverage(cpuValues),
		"max": s.calculateMax(cpuValues),
		"min": s.calculateMin(cpuValues),
		"std": s.calculateStandardDeviation(cpuValues),
		"p95": s.calculatePercentile(cpuValues, 95),
		"p99": s.calculatePercentile(cpuValues, 99),
	}

	// 内存统计
	memValues := s.extractValues(performanceData, "memory_usage")
	stats["memory"] = map[string]interface{}{
		"avg": s.calculateAverage(memValues),
		"max": s.calculateMax(memValues),
		"min": s.calculateMin(memValues),
		"std": s.calculateStandardDeviation(memValues),
		"p95": s.calculatePercentile(memValues, 95),
		"p99": s.calculatePercentile(memValues, 99),
	}

	// 磁盘统计
	diskValues := s.extractValues(performanceData, "disk_usage")
	stats["disk"] = map[string]interface{}{
		"avg": s.calculateAverage(diskValues),
		"max": s.calculateMax(diskValues),
		"min": s.calculateMin(diskValues),
		"std": s.calculateStandardDeviation(diskValues),
		"p95": s.calculatePercentile(diskValues, 95),
		"p99": s.calculatePercentile(diskValues, 99),
	}

	// 网络统计
	netValues := s.extractValues(performanceData, "network_latency")
	stats["network"] = map[string]interface{}{
		"avg": s.calculateAverage(netValues),
		"max": s.calculateMax(netValues),
		"min": s.calculateMin(netValues),
		"std": s.calculateStandardDeviation(netValues),
		"p95": s.calculatePercentile(netValues, 95),
		"p99": s.calculatePercentile(netValues, 99),
	}

	return stats
}

// analyzePerformanceTrends 分析性能趋势
func (s *sDevice) analyzePerformanceTrends(performanceData []map[string]interface{}) map[string]interface{} {
	trends := make(map[string]interface{})

	// 按时间排序数据
	sortedData := s.sortByTimestamp(performanceData)

	// 分析CPU趋势
	cpuTrend := s.analyzeMetricTrend(sortedData, "cpu_usage")
	trends["cpu"] = cpuTrend

	// 分析内存趋势
	memTrend := s.analyzeMetricTrend(sortedData, "memory_usage")
	trends["memory"] = memTrend

	// 分析磁盘趋势
	diskTrend := s.analyzeMetricTrend(sortedData, "disk_usage")
	trends["disk"] = diskTrend

	// 分析网络趋势
	netTrend := s.analyzeMetricTrend(sortedData, "network_latency")
	trends["network"] = netTrend

	return trends
}

// analyzeMetricTrend 分析单个指标趋势
func (s *sDevice) analyzeMetricTrend(sortedData []map[string]interface{}, metricKey string) map[string]interface{} {
	values := s.extractValues(sortedData, metricKey)
	if len(values) < 2 {
		return map[string]interface{}{
			"trend": "insufficient_data",
			"slope": 0,
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

	return map[string]interface{}{
		"trend":      trend,
		"slope":      slope,
		"volatility": s.calculateVolatility(values),
	}
}

// identifyBottlenecks 识别性能瓶颈
func (s *sDevice) identifyBottlenecks(performanceData []map[string]interface{}) []map[string]interface{} {
	var bottlenecks []map[string]interface{}

	// 检查CPU瓶颈
	if cpuBottleneck := s.checkCPUBottleneck(performanceData); cpuBottleneck != nil {
		bottlenecks = append(bottlenecks, cpuBottleneck)
	}

	// 检查内存瓶颈
	if memBottleneck := s.checkMemoryBottleneck(performanceData); memBottleneck != nil {
		bottlenecks = append(bottlenecks, memBottleneck)
	}

	// 检查磁盘瓶颈
	if diskBottleneck := s.checkDiskBottleneck(performanceData); diskBottleneck != nil {
		bottlenecks = append(bottlenecks, diskBottleneck)
	}

	// 检查网络瓶颈
	if netBottleneck := s.checkNetworkBottleneck(performanceData); netBottleneck != nil {
		bottlenecks = append(bottlenecks, netBottleneck)
	}

	return bottlenecks
}

// checkCPUBottleneck 检查CPU瓶颈
func (s *sDevice) checkCPUBottleneck(performanceData []map[string]interface{}) map[string]interface{} {
	cpuValues := s.extractValues(performanceData, "cpu_usage")
	avgCPU := s.calculateAverage(cpuValues)
	p95CPU := s.calculatePercentile(cpuValues, 95)

	if avgCPU > 80 || p95CPU > 95 {
		return map[string]interface{}{
			"type":           "cpu",
			"severity":       s.getBottleneckSeverity(avgCPU),
			"avg_usage":      avgCPU,
			"p95_usage":      p95CPU,
			"description":    "CPU使用率过高，可能影响系统性能",
			"recommendation": "考虑优化CPU密集型任务或增加CPU资源",
		}
	}

	return nil
}

// checkMemoryBottleneck 检查内存瓶颈
func (s *sDevice) checkMemoryBottleneck(performanceData []map[string]interface{}) map[string]interface{} {
	memValues := s.extractValues(performanceData, "memory_usage")
	avgMem := s.calculateAverage(memValues)
	p95Mem := s.calculatePercentile(memValues, 95)

	if avgMem > 85 || p95Mem > 95 {
		return map[string]interface{}{
			"type":           "memory",
			"severity":       s.getBottleneckSeverity(avgMem),
			"avg_usage":      avgMem,
			"p95_usage":      p95Mem,
			"description":    "内存使用率过高，可能导致系统不稳定",
			"recommendation": "考虑增加内存或优化内存使用",
		}
	}

	return nil
}

// checkDiskBottleneck 检查磁盘瓶颈
func (s *sDevice) checkDiskBottleneck(performanceData []map[string]interface{}) map[string]interface{} {
	diskValues := s.extractValues(performanceData, "disk_usage")
	avgDisk := s.calculateAverage(diskValues)
	p95Disk := s.calculatePercentile(diskValues, 95)

	if avgDisk > 90 || p95Disk > 98 {
		return map[string]interface{}{
			"type":           "disk",
			"severity":       s.getBottleneckSeverity(avgDisk),
			"avg_usage":      avgDisk,
			"p95_usage":      p95Disk,
			"description":    "磁盘使用率过高，可能影响I/O性能",
			"recommendation": "考虑清理磁盘空间或扩容",
		}
	}

	return nil
}

// checkNetworkBottleneck 检查网络瓶颈
func (s *sDevice) checkNetworkBottleneck(performanceData []map[string]interface{}) map[string]interface{} {
	netValues := s.extractValues(performanceData, "network_latency")
	avgNet := s.calculateAverage(netValues)
	p95Net := s.calculatePercentile(netValues, 95)

	if avgNet > 100 || p95Net > 200 {
		return map[string]interface{}{
			"type":           "network",
			"severity":       s.getBottleneckSeverity(avgNet / 10), // 标准化到0-100
			"avg_latency":    avgNet,
			"p95_latency":    p95Net,
			"description":    "网络延迟过高，可能影响响应时间",
			"recommendation": "检查网络连接或优化网络配置",
		}
	}

	return nil
}

// calculatePerformanceScore 计算性能评分
func (s *sDevice) calculatePerformanceScore(performanceData []map[string]interface{}) float64 {
	score := 100.0

	// CPU评分 (权重: 30%)
	cpuValues := s.extractValues(performanceData, "cpu_usage")
	avgCPU := s.calculateAverage(cpuValues)
	if avgCPU > 90 {
		score -= 30
	} else if avgCPU > 70 {
		score -= 15
	} else if avgCPU > 50 {
		score -= 5
	}

	// 内存评分 (权重: 25%)
	memValues := s.extractValues(performanceData, "memory_usage")
	avgMem := s.calculateAverage(memValues)
	if avgMem > 90 {
		score -= 25
	} else if avgMem > 80 {
		score -= 12
	} else if avgMem > 60 {
		score -= 5
	}

	// 磁盘评分 (权重: 20%)
	diskValues := s.extractValues(performanceData, "disk_usage")
	avgDisk := s.calculateAverage(diskValues)
	if avgDisk > 95 {
		score -= 20
	} else if avgDisk > 85 {
		score -= 10
	} else if avgDisk > 70 {
		score -= 3
	}

	// 网络评分 (权重: 15%)
	netValues := s.extractValues(performanceData, "network_latency")
	avgNet := s.calculateAverage(netValues)
	if avgNet > 200 {
		score -= 15
	} else if avgNet > 100 {
		score -= 8
	} else if avgNet > 50 {
		score -= 3
	}

	// 稳定性评分 (权重: 10%)
	stabilityScore := s.calculateStabilityScore(performanceData)
	score += stabilityScore * 0.1

	return math.Max(0, math.Min(score, 100))
}

// generatePerformanceRecommendations 生成性能优化建议
func (s *sDevice) generatePerformanceRecommendations(analysis map[string]interface{}) []string {
	var recommendations []string

	// 基于性能评分生成建议
	if performanceScore, ok := analysis["performance_score"].(float64); ok {
		if performanceScore < 50 {
			recommendations = append(recommendations, "性能严重不足，建议立即优化")
		} else if performanceScore < 70 {
			recommendations = append(recommendations, "性能偏低，建议进行优化")
		}
	}

	// 基于瓶颈分析生成建议
	if bottlenecks, ok := analysis["bottlenecks"].([]map[string]interface{}); ok {
		for _, bottleneck := range bottlenecks {
			if recommendation, ok := bottleneck["recommendation"].(string); ok {
				recommendations = append(recommendations, recommendation)
			}
		}
	}

	// 基于趋势分析生成建议
	if trends, ok := analysis["trends"].(map[string]interface{}); ok {
		if cpuTrend, ok := trends["cpu"].(map[string]interface{}); ok {
			if trend, ok := cpuTrend["trend"].(string); ok && trend == "increasing" {
				recommendations = append(recommendations, "CPU使用率呈上升趋势，建议监控并优化")
			}
		}
	}

	return recommendations
}

// 辅助方法
func (s *sDevice) extractValues(data []map[string]interface{}, key string) []float64 {
	var values []float64
	for _, item := range data {
		if value, ok := item[key].(float64); ok {
			values = append(values, value)
		}
	}
	return values
}

func (s *sDevice) sortByTimestamp(data []map[string]interface{}) []map[string]interface{} {
	// 这里应该按时间戳排序
	// 目前返回原数据
	return data
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

func (s *sDevice) calculateVolatility(values []float64) float64 {
	if len(values) < 2 {
		return 0
	}

	mean := s.calculateAverage(values)
	variance := 0.0

	for _, value := range values {
		variance += math.Pow(value-mean, 2)
	}

	variance /= float64(len(values))
	return math.Sqrt(variance)
}

func (s *sDevice) calculateStandardDeviation(values []float64) float64 {
	return s.calculateVolatility(values)
}

func (s *sDevice) calculatePercentile(values []float64, percentile int) float64 {
	if len(values) == 0 {
		return 0
	}

	// 简单实现，实际应该使用更精确的算法
	index := int(float64(len(values)-1) * float64(percentile) / 100.0)
	if index >= len(values) {
		index = len(values) - 1
	}

	return values[index]
}

func (s *sDevice) calculateStabilityScore(performanceData []map[string]interface{}) float64 {
	// 计算性能稳定性评分
	// 基于各指标的变异系数
	score := 100.0

	cpuValues := s.extractValues(performanceData, "cpu_usage")
	if len(cpuValues) > 0 {
		cv := s.calculateVolatility(cpuValues) / s.calculateAverage(cpuValues)
		if cv > 0.5 {
			score -= 20
		} else if cv > 0.3 {
			score -= 10
		}
	}

	return math.Max(0, score)
}

func (s *sDevice) getBottleneckSeverity(usage float64) string {
	if usage > 95 {
		return "critical"
	} else if usage > 85 {
		return "high"
	} else if usage > 70 {
		return "medium"
	} else {
		return "low"
	}
}

func (s *sDevice) getEarliestTimestamp(data []map[string]interface{}) string {
	// 这里应该返回最早的时间戳
	return gtime.Now().Add(-24 * time.Hour).Format("2006-01-02 15:04:05")
}

func (s *sDevice) getLatestTimestamp(data []map[string]interface{}) string {
	// 这里应该返回最晚的时间戳
	return gtime.Now().Format("2006-01-02 15:04:05")
}
