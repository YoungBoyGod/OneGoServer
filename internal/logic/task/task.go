package task

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

type sTask struct{}

func New() *sTask {
	return &sTask{}
}

func init() {
	// service.RegisterTask(New())
}

// ===============================
// 任务优先级计算相关业务逻辑
// ===============================

// CalculateTaskPriority 计算任务智能优先级
func (s *sTask) CalculateTaskPriority(ctx context.Context, taskData map[string]interface{}) int {
	baseScore := 0.0

	// 1. 基础优先级权重 (权重: 30%)
	if priority, ok := taskData["priority"].(int); ok && priority >= 1 && priority <= 10 {
		baseScore += float64(priority) * 10 * 0.3
	} else {
		baseScore += 50 * 0.3 // 默认中等优先级
	}

	// 2. 紧急程度权重 (权重: 25%)
	urgencyScore := s.calculateUrgencyScore(taskData)
	baseScore += urgencyScore * 0.25

	// 3. 业务重要性权重 (权重: 20%)
	businessScore := s.calculateBusinessImportanceScore(taskData)
	baseScore += businessScore * 0.20

	// 4. 截止时间影响 (权重: 15%)
	deadlineScore := s.calculateDeadlineScore(taskData)
	baseScore += deadlineScore * 0.15

	// 5. 资源需求权重 (权重: 10%)
	resourceScore := s.calculateResourceScore(taskData)
	baseScore += resourceScore * 0.10

	// 将评分转换为1-10的优先级
	priority := int(math.Round(baseScore / 10))
	if priority < 1 {
		priority = 1
	}
	if priority > 10 {
		priority = 10
	}

	return priority
}

// calculateUrgencyScore 计算紧急程度评分
func (s *sTask) calculateUrgencyScore(taskData map[string]interface{}) float64 {
	score := 50.0 // 默认分数

	// 检查任务类型的紧急性
	if taskType, ok := taskData["task_type"].(string); ok {
		switch taskType {
		case "emergency":
			score = 100
		case "urgent":
			score = 80
		case "normal":
			score = 50
		case "low":
			score = 20
		}
	}

	// 检查是否有紧急标签
	if tags, ok := taskData["tags"].([]string); ok {
		for _, tag := range tags {
			if strings.Contains(strings.ToLower(tag), "urgent") ||
				strings.Contains(strings.ToLower(tag), "emergency") {
				score = math.Max(score, 90)
			}
		}
	}

	return score
}

// calculateBusinessImportanceScore 计算业务重要性评分
func (s *sTask) calculateBusinessImportanceScore(taskData map[string]interface{}) float64 {
	score := 50.0 // 默认分数

	// 检查任务类别
	if category, ok := taskData["category"].(string); ok {
		switch category {
		case "critical":
			score = 100
		case "important":
			score = 80
		case "normal":
			score = 50
		case "minor":
			score = 30
		}
	}

	// 检查影响范围
	if impact, ok := taskData["impact_scope"].(string); ok {
		switch impact {
		case "global":
			score += 20
		case "regional":
			score += 10
		case "local":
			score += 5
		}
	}

	return math.Min(score, 100)
}

// calculateDeadlineScore 计算截止时间评分
func (s *sTask) calculateDeadlineScore(taskData map[string]interface{}) float64 {
	score := 50.0 // 默认分数

	if deadlineAt, ok := taskData["deadline_at"].(*gtime.Time); ok && deadlineAt != nil {
		now := time.Now()
		timeToDeadline := deadlineAt.Time.Sub(now)

		if timeToDeadline < 0 {
			// 已过期
			score = 100
		} else if timeToDeadline < 1*time.Hour {
			// 1小时内截止
			score = 95
		} else if timeToDeadline < 6*time.Hour {
			// 6小时内截止
			score = 85
		} else if timeToDeadline < 24*time.Hour {
			// 24小时内截止
			score = 75
		} else if timeToDeadline < 7*24*time.Hour {
			// 一周内截止
			score = 60
		} else {
			// 超过一周
			score = 30
		}
	}

	return score
}

// calculateResourceScore 计算资源需求评分
func (s *sTask) calculateResourceScore(taskData map[string]interface{}) float64 {
	score := 50.0 // 默认分数

	// 检查预估执行时间
	if estimatedTime, ok := taskData["estimated_duration"].(int); ok {
		if estimatedTime < 300 { // 5分钟以内
			score = 80 // 快速任务优先级高
		} else if estimatedTime < 1800 { // 30分钟以内
			score = 60
		} else if estimatedTime < 3600 { // 1小时以内
			score = 40
		} else {
			score = 20 // 长时间任务优先级低
		}
	}

	// 检查资源复杂度
	if complexity, ok := taskData["complexity"].(string); ok {
		switch complexity {
		case "low":
			score += 10
		case "medium":
			score += 0
		case "high":
			score -= 10
		}
	}

	return math.Max(score, 0)
}

// ===============================
// 任务状态流转相关业务逻辑
// ===============================

// ValidateTaskStatusTransition 验证任务状态转换是否合法
func (s *sTask) ValidateTaskStatusTransition(ctx context.Context, currentStatus, targetStatus string) error {
	// 定义合法的状态转换规则
	validTransitions := map[string][]string{
		"pending":   {"running", "cancelled", "failed"},
		"running":   {"completed", "paused", "failed", "cancelled"},
		"paused":    {"running", "cancelled", "failed"},
		"completed": {"pending"},              // 可以重新执行
		"failed":    {"pending", "cancelled"}, // 可以重试或取消
		"cancelled": {"pending"},              // 可以重新激活
	}

	allowedStatuses, exists := validTransitions[currentStatus]
	if !exists {
		return gerror.NewCode(gcode.CodeInvalidParameter, fmt.Sprintf("无效的当前状态: %s", currentStatus))
	}

	// 检查目标状态是否在允许的转换列表中
	for _, status := range allowedStatuses {
		if status == targetStatus {
			return nil
		}
	}

	return gerror.NewCode(gcode.CodeInvalidParameter,
		fmt.Sprintf("不允许从状态 %s 转换到 %s", currentStatus, targetStatus))
}

// DetermineTaskStatus 根据任务数据自动确定任务状态
func (s *sTask) DetermineTaskStatus(ctx context.Context, taskData map[string]interface{}) string {
	// 检查是否有明确的状态设置
	if status, ok := taskData["status"].(string); ok && status != "" {
		return status
	}

	// 检查是否已取消
	if cancelled, ok := taskData["is_cancelled"].(bool); ok && cancelled {
		return "cancelled"
	}

	// 检查是否已完成
	if completed, ok := taskData["is_completed"].(bool); ok && completed {
		return "completed"
	}

	// 检查是否已开始执行
	if startedAt, ok := taskData["started_at"].(*gtime.Time); ok && startedAt != nil {
		// 检查是否暂停
		if paused, ok := taskData["is_paused"].(bool); ok && paused {
			return "paused"
		}

		// 检查是否执行失败
		if failed, ok := taskData["is_failed"].(bool); ok && failed {
			return "failed"
		}

		return "running"
	}

	// 默认为待执行状态
	return "pending"
}

// CanTransitionToStatus 判断任务是否可以转换到目标状态
func (s *sTask) CanTransitionToStatus(ctx context.Context, taskData map[string]interface{}, targetStatus string) bool {
	currentStatus := s.DetermineTaskStatus(ctx, taskData)
	err := s.ValidateTaskStatusTransition(ctx, currentStatus, targetStatus)
	return err == nil
}

// ===============================
// 任务验证相关业务逻辑
// ===============================

// ValidateTaskCreation 验证任务创建数据
func (s *sTask) ValidateTaskCreation(ctx context.Context, taskData map[string]interface{}) error {
	// 验证必需字段
	requiredFields := []string{"task_name", "task_type"}
	for _, field := range requiredFields {
		if value, ok := taskData[field]; !ok || value == nil || value == "" {
			return gerror.NewCode(gcode.CodeInvalidParameter, fmt.Sprintf("缺少必需字段: %s", field))
		}
	}

	// 验证任务名称
	if taskName, ok := taskData["task_name"].(string); ok {
		if len(taskName) < 1 || len(taskName) > 100 {
			return gerror.NewCode(gcode.CodeInvalidParameter, "任务名称长度必须在1-100字符之间")
		}
	}

	// 验证任务类型
	if taskType, ok := taskData["task_type"].(string); ok {
		validTypes := []string{"backup", "sync", "monitor", "custom", "system", "user"}
		if !contains(validTypes, taskType) {
			return gerror.NewCode(gcode.CodeInvalidParameter, "无效的任务类型")
		}
	}

	// 验证优先级
	if priority, ok := taskData["priority"]; ok {
		if p := gconv.Int(priority); p < 1 || p > 10 {
			return gerror.NewCode(gcode.CodeInvalidParameter, "任务优先级必须在1-10之间")
		}
	}

	// 验证重试次数
	if retryCount, ok := taskData["retry_count"]; ok {
		if count := gconv.Int(retryCount); count < 0 || count > 10 {
			return gerror.NewCode(gcode.CodeInvalidParameter, "重试次数必须在0-10之间")
		}
	}

	// 验证时间设置
	if err := s.validateTaskTiming(taskData); err != nil {
		return err
	}

	return nil
}

// validateTaskTiming 验证任务时间设置
func (s *sTask) validateTaskTiming(taskData map[string]interface{}) error {
	var scheduleAt, deadlineAt *gtime.Time

	if sa, ok := taskData["schedule_at"].(*gtime.Time); ok {
		scheduleAt = sa
	}
	if da, ok := taskData["deadline_at"].(*gtime.Time); ok {
		deadlineAt = da
	}

	// 如果都设置了，验证逻辑关系
	if scheduleAt != nil && deadlineAt != nil {
		if scheduleAt.After(deadlineAt.Time) {
			return gerror.NewCode(gcode.CodeInvalidParameter, "调度时间不能晚于截止时间")
		}
	}

	// 验证调度时间不能是过去时间（允许一定的容差）
	if scheduleAt != nil && scheduleAt.Before(time.Now().Add(-1*time.Minute)) {
		return gerror.NewCode(gcode.CodeInvalidParameter, "调度时间不能是过去时间")
	}

	return nil
}

// ValidateTaskConfiguration 验证任务配置数据
func (s *sTask) ValidateTaskConfiguration(ctx context.Context, config map[string]interface{}) error {
	// 验证超时设置
	if timeout, ok := config["timeout"]; ok {
		if t := gconv.Int(timeout); t < 30 || t > 86400 {
			return gerror.NewCode(gcode.CodeInvalidParameter, "任务超时时间必须在30-86400秒之间")
		}
	}

	// 验证并发设置
	if concurrency, ok := config["max_concurrency"]; ok {
		if c := gconv.Int(concurrency); c < 1 || c > 100 {
			return gerror.NewCode(gcode.CodeInvalidParameter, "最大并发数必须在1-100之间")
		}
	}

	// 验证资源限制
	if resourceLimits, ok := config["resource_limits"].(map[string]interface{}); ok {
		if cpu, exists := resourceLimits["cpu"]; exists {
			if limit := gconv.Float64(cpu); limit < 0 || limit > 100 {
				return gerror.NewCode(gcode.CodeInvalidParameter, "CPU限制必须在0-100之间")
			}
		}
		if memory, exists := resourceLimits["memory"]; exists {
			if limit := gconv.Float64(memory); limit < 0 || limit > 100 {
				return gerror.NewCode(gcode.CodeInvalidParameter, "内存限制必须在0-100之间")
			}
		}
	}

	return nil
}

// ===============================
// 任务调度策略相关业务逻辑
// ===============================

// CalculateTaskSchedulingScore 计算任务调度评分
func (s *sTask) CalculateTaskSchedulingScore(ctx context.Context, taskData map[string]interface{}, deviceData map[string]interface{}) float64 {
	score := 0.0

	// 1. 任务优先级 (权重: 30%)
	if priority, ok := taskData["priority"].(int); ok {
		score += float64(priority) * 10 * 0.3
	}

	// 2. 设备适配度 (权重: 25%)
	deviceScore := s.calculateDeviceCompatibilityScore(taskData, deviceData)
	score += deviceScore * 0.25

	// 3. 紧急程度 (权重: 20%)
	urgencyScore := s.calculateUrgencyScore(taskData)
	score += urgencyScore * 0.20

	// 4. 资源效率 (权重: 15%)
	resourceScore := s.calculateResourceEfficiencyScore(taskData, deviceData)
	score += resourceScore * 0.15

	// 5. 时间因素 (权重: 10%)
	timeScore := s.calculateTimeFactorScore(taskData)
	score += timeScore * 0.10

	return math.Round(score*100) / 100
}

// calculateDeviceCompatibilityScore 计算设备兼容性评分
func (s *sTask) calculateDeviceCompatibilityScore(taskData, deviceData map[string]interface{}) float64 {
	score := 50.0 // 基础分数

	// 检查设备类型匹配
	if taskDeviceType, ok := taskData["preferred_device_type"].(string); ok {
		if deviceType, ok := deviceData["device_type"].(string); ok {
			if taskDeviceType == deviceType {
				score += 30
			} else {
				score += 10 // 部分兼容
			}
		}
	}

	// 检查设备能力
	if requiredCapabilities, ok := taskData["required_capabilities"].([]string); ok {
		if deviceCapabilities, ok := deviceData["capabilities"].([]string); ok {
			matchCount := 0
			for _, required := range requiredCapabilities {
				if contains(deviceCapabilities, required) {
					matchCount++
				}
			}
			if len(requiredCapabilities) > 0 {
				matchRatio := float64(matchCount) / float64(len(requiredCapabilities))
				score += matchRatio * 20
			}
		}
	}

	return math.Min(score, 100)
}

// calculateResourceEfficiencyScore 计算资源效率评分
func (s *sTask) calculateResourceEfficiencyScore(taskData, deviceData map[string]interface{}) float64 {
	score := 50.0

	// 检查设备负载
	if currentLoad, ok := deviceData["current_load"].(float64); ok {
		if currentLoad < 0.3 {
			score += 20 // 低负载，效率高
		} else if currentLoad < 0.7 {
			score += 10 // 中等负载
		} else {
			score -= 10 // 高负载，效率低
		}
	}

	// 检查任务资源需求
	if resourceReq, ok := taskData["resource_requirements"].(map[string]interface{}); ok {
		if cpu, ok := resourceReq["cpu"].(float64); ok && cpu > 80 {
			score -= 10 // 高CPU需求降低评分
		}
		if memory, ok := resourceReq["memory"].(float64); ok && memory > 80 {
			score -= 10 // 高内存需求降低评分
		}
	}

	return math.Max(score, 0)
}

// calculateTimeFactorScore 计算时间因素评分
func (s *sTask) calculateTimeFactorScore(taskData map[string]interface{}) float64 {
	score := 50.0

	now := time.Now()

	// 检查调度时间
	if scheduleAt, ok := taskData["schedule_at"].(*gtime.Time); ok && scheduleAt != nil {
		timeDiff := scheduleAt.Time.Sub(now)
		if timeDiff <= 0 {
			score += 30 // 应该立即执行
		} else if timeDiff < 5*time.Minute {
			score += 20 // 即将执行
		} else if timeDiff < 30*time.Minute {
			score += 10 // 近期执行
		}
	}

	// 检查截止时间
	if deadlineAt, ok := taskData["deadline_at"].(*gtime.Time); ok && deadlineAt != nil {
		timeToDeadline := deadlineAt.Time.Sub(now)
		if timeToDeadline < 1*time.Hour {
			score += 25 // 紧急
		} else if timeToDeadline < 6*time.Hour {
			score += 15 // 较紧急
		} else if timeToDeadline < 24*time.Hour {
			score += 5 // 普通
		}
	}

	return math.Min(score, 100)
}

// ===============================
// 任务依赖管理相关业务逻辑
// ===============================

// ValidateTaskDependencies 验证任务依赖关系
func (s *sTask) ValidateTaskDependencies(ctx context.Context, taskId string, dependencies []string, allTasks map[string]map[string]interface{}) error {
	// 检查是否存在循环依赖
	if s.hasCircularDependency(taskId, dependencies, allTasks, make(map[string]bool)) {
		return gerror.NewCode(gcode.CodeInvalidParameter, "检测到循环依赖")
	}

	// 检查依赖的任务是否存在
	for _, depId := range dependencies {
		if _, exists := allTasks[depId]; !exists {
			return gerror.NewCode(gcode.CodeInvalidParameter, fmt.Sprintf("依赖的任务不存在: %s", depId))
		}
	}

	return nil
}

// hasCircularDependency 检查是否存在循环依赖
func (s *sTask) hasCircularDependency(taskId string, dependencies []string, allTasks map[string]map[string]interface{}, visited map[string]bool) bool {
	if visited[taskId] {
		return true
	}

	visited[taskId] = true

	for _, depId := range dependencies {
		if depTask, exists := allTasks[depId]; exists {
			if depDependencies, ok := depTask["dependencies"].([]string); ok {
				if s.hasCircularDependency(depId, depDependencies, allTasks, visited) {
					return true
				}
			}
		}
	}

	delete(visited, taskId)
	return false
}

// SortTasksByDependencies 根据依赖关系排序任务
func (s *sTask) SortTasksByDependencies(ctx context.Context, tasks []map[string]interface{}) []map[string]interface{} {
	// 构建依赖图
	taskMap := make(map[string]map[string]interface{})
	for _, task := range tasks {
		if taskId, ok := task["task_id"].(string); ok {
			taskMap[taskId] = task
		}
	}

	// 拓扑排序
	sorted := make([]map[string]interface{}, 0, len(tasks))
	visited := make(map[string]bool)
	visiting := make(map[string]bool)

	var visit func(taskId string) bool
	visit = func(taskId string) bool {
		if visiting[taskId] {
			// 检测到循环依赖
			return false
		}
		if visited[taskId] {
			return true
		}

		visiting[taskId] = true

		task := taskMap[taskId]
		if dependencies, ok := task["dependencies"].([]string); ok {
			for _, depId := range dependencies {
				if !visit(depId) {
					return false
				}
			}
		}

		visiting[taskId] = false
		visited[taskId] = true
		sorted = append(sorted, task)
		return true
	}

	// 访问所有任务
	for _, task := range tasks {
		if taskId, ok := task["task_id"].(string); ok {
			if !visited[taskId] {
				visit(taskId)
			}
		}
	}

	return sorted
}

// ===============================
// 任务执行监控相关业务逻辑
// ===============================

// CalculateTaskExecutionScore 计算任务执行评分
func (s *sTask) CalculateTaskExecutionScore(ctx context.Context, executionData map[string]interface{}) float64 {
	score := 100.0

	// 检查执行时间效率
	if estimatedTime, ok := executionData["estimated_duration"].(int); ok {
		if actualTime, ok := executionData["actual_duration"].(int); ok && estimatedTime > 0 {
			ratio := float64(actualTime) / float64(estimatedTime)
			if ratio <= 1.0 {
				score += 10 // 提前或按时完成
			} else if ratio <= 1.2 {
				score -= 5 // 稍微超时
			} else if ratio <= 1.5 {
				score -= 15 // 明显超时
			} else {
				score -= 30 // 严重超时
			}
		}
	}

	// 检查错误率
	if errorCount, ok := executionData["error_count"].(int); ok {
		score -= float64(errorCount) * 5
	}

	// 检查重试次数
	if retryCount, ok := executionData["retry_count"].(int); ok {
		score -= float64(retryCount) * 3
	}

	// 检查资源使用效率
	if resourceUsage, ok := executionData["resource_usage"].(map[string]interface{}); ok {
		if cpu, ok := resourceUsage["cpu"].(float64); ok {
			if cpu > 90 {
				score -= 10
			} else if cpu < 20 {
				score -= 5 // 资源利用率过低也不好
			}
		}
	}

	return math.Max(score, 0)
}

// AnalyzeTaskPerformance 分析任务性能表现
func (s *sTask) AnalyzeTaskPerformance(ctx context.Context, performanceData []map[string]interface{}) map[string]interface{} {
	if len(performanceData) == 0 {
		return map[string]interface{}{
			"status":  "insufficient_data",
			"message": "性能数据不足",
		}
	}

	analysis := make(map[string]interface{})

	var totalDuration, totalScore float64
	var successCount, failureCount int
	var executionScores []float64

	for _, data := range performanceData {
		// 统计执行时间
		if duration, ok := data["duration"].(float64); ok {
			totalDuration += duration
		}

		// 统计成功失败率
		if status, ok := data["status"].(string); ok {
			if status == "completed" {
				successCount++
			} else if status == "failed" {
				failureCount++
			}
		}

		// 计算执行评分
		score := s.CalculateTaskExecutionScore(ctx, data)
		executionScores = append(executionScores, score)
		totalScore += score
	}

	count := float64(len(performanceData))
	analysis["avg_duration"] = math.Round((totalDuration/count)*100) / 100
	analysis["avg_execution_score"] = math.Round((totalScore/count)*100) / 100
	analysis["success_rate"] = float64(successCount) / count * 100
	analysis["failure_rate"] = float64(failureCount) / count * 100

	// 生成性能建议
	recommendations := s.generateTaskPerformanceRecommendations(analysis)
	analysis["recommendations"] = recommendations

	return analysis
}

// generateTaskPerformanceRecommendations 生成任务性能优化建议
func (s *sTask) generateTaskPerformanceRecommendations(analysis map[string]interface{}) []string {
	var recommendations []string

	if failureRate, ok := analysis["failure_rate"].(float64); ok && failureRate > 10 {
		recommendations = append(recommendations, "失败率较高，建议检查任务逻辑和环境配置")
	}

	if avgScore, ok := analysis["avg_execution_score"].(float64); ok && avgScore < 70 {
		recommendations = append(recommendations, "执行评分较低，建议优化任务性能")
	}

	if avgDuration, ok := analysis["avg_duration"].(float64); ok && avgDuration > 3600 {
		recommendations = append(recommendations, "执行时间较长，建议优化算法或增加资源")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "任务执行表现良好，保持当前配置")
	}

	return recommendations
}

// ===============================
// 任务分配相关业务逻辑
// ===============================

// AssignTaskToDevice 分配任务给设备
func (s *sTask) AssignTaskToDevice(ctx context.Context, taskId string, deviceIds []string, strategy string, force bool) (map[string]interface{}, error) {
	// 验证任务存在
	if err := s.validateTaskExists(ctx, taskId); err != nil {
		return nil, err
	}

	// 获取任务信息
	taskInfo, err := s.getTaskInfo(ctx, taskId)
	if err != nil {
		return nil, err
	}

	// 获取可用设备列表
	availableDevices, err := s.getAvailableDevices(ctx, deviceIds)
	if err != nil {
		return nil, err
	}

	// 根据策略选择最佳设备
	selectedDevice, assignmentScore, reason, err := s.selectBestDevice(ctx, taskInfo, availableDevices, strategy, force)
	if err != nil {
		return nil, err
	}

	// 执行任务分配
	deviceId := s.getStringValue(selectedDevice, "device_id", "")
	if err := s.executeTaskAssignment(ctx, taskId, deviceId); err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"taskId":           taskId,
		"assignedDeviceId": deviceId,
		"assignmentScore":  assignmentScore,
		"reason":           reason,
		"status":           "assigned",
	}

	return result, nil
}

// selectBestDevice 选择最佳设备
func (s *sTask) selectBestDevice(ctx context.Context, taskInfo map[string]interface{}, devices []map[string]interface{}, strategy string, force bool) (map[string]interface{}, float64, string, error) {
	// 过滤兼容设备
	compatibleDevices := s.filterCompatibleDevices(taskInfo, devices)
	if len(compatibleDevices) == 0 {
		return nil, 0, "", gerror.NewCode(gcode.CodeValidationFailed, "没有兼容的设备可用")
	}

	// 根据策略计算分配评分
	switch strategy {
	case "auto":
		return s.calculateAutoAssignment(compatibleDevices, taskInfo)
	case "manual":
		return s.calculateManualAssignment(compatibleDevices, taskInfo)
	case "load_balanced":
		return s.calculateLoadBalancedAssignment(compatibleDevices, taskInfo)
	case "priority":
		return s.calculatePriorityAssignment(compatibleDevices, taskInfo)
	default:
		return s.calculateAutoAssignment(compatibleDevices, taskInfo)
	}
}

// filterCompatibleDevices 过滤兼容设备
func (s *sTask) filterCompatibleDevices(taskInfo map[string]interface{}, devices []map[string]interface{}) []map[string]interface{} {
	var compatibleDevices []map[string]interface{}
	taskType := s.getStringValue(taskInfo, "type", "")
	requiredCapabilities := s.getStringValue(taskInfo, "required_capabilities", "")

	for _, device := range devices {
		deviceType := s.getStringValue(device, "type", "")
		deviceCapabilities := s.getStringValue(device, "capabilities", "")
		deviceStatus := s.getStringValue(device, "status", "")

		// 检查设备状态
		if deviceStatus != "online" && deviceStatus != "idle" {
			continue
		}

		// 检查设备类型兼容性
		if taskType != "" && deviceType != "" && taskType != deviceType {
			continue
		}

		// 检查设备能力兼容性
		if requiredCapabilities != "" && deviceCapabilities != "" {
			if !s.checkCapabilityCompatibility(requiredCapabilities, deviceCapabilities) {
				continue
			}
		}

		compatibleDevices = append(compatibleDevices, device)
	}

	return compatibleDevices
}

// calculateAutoAssignment 自动分配策略
func (s *sTask) calculateAutoAssignment(devices []map[string]interface{}, taskInfo map[string]interface{}) (map[string]interface{}, float64, string, error) {
	var bestDevice map[string]interface{}
	var bestScore float64 = -1
	var reason string

	for _, device := range devices {
		// 计算综合评分
		loadScore := s.calculateDeviceLoadScore(device)
		compatibilityScore := s.calculateCompatibilityScore(device, taskInfo)
		availabilityScore := s.calculateAvailabilityScore(device)
		performanceScore := s.calculatePerformanceScore(device)

		// 加权综合评分
		totalScore := loadScore*0.3 + compatibilityScore*0.3 + availabilityScore*0.2 + performanceScore*0.2

		if totalScore > bestScore {
			bestScore = totalScore
			bestDevice = device
			reason = "自动选择综合评分最高的设备"
		}
	}

	if bestDevice == nil {
		return nil, 0, "", gerror.NewCode(gcode.CodeValidationFailed, "没有合适的设备")
	}

	return bestDevice, bestScore, reason, nil
}

// calculateLoadBalancedAssignment 负载均衡分配策略
func (s *sTask) calculateLoadBalancedAssignment(devices []map[string]interface{}, taskInfo map[string]interface{}) (map[string]interface{}, float64, string, error) {
	var bestDevice map[string]interface{}
	var bestScore float64 = -1
	var reason string

	for _, device := range devices {
		// 主要考虑负载评分
		loadScore := s.calculateDeviceLoadScore(device)
		compatibilityScore := s.calculateCompatibilityScore(device, taskInfo)

		// 负载越低分数越高
		balancedScore := (1.0-loadScore)*0.7 + compatibilityScore*0.3

		if balancedScore > bestScore {
			bestScore = balancedScore
			bestDevice = device
			reason = "负载均衡策略选择负载最低的设备"
		}
	}

	if bestDevice == nil {
		return nil, 0, "", gerror.NewCode(gcode.CodeValidationFailed, "没有合适的设备")
	}

	return bestDevice, bestScore, reason, nil
}

// calculatePriorityAssignment 优先级分配策略
func (s *sTask) calculatePriorityAssignment(devices []map[string]interface{}, taskInfo map[string]interface{}) (map[string]interface{}, float64, string, error) {
	var bestDevice map[string]interface{}
	var bestScore float64 = -1
	var reason string

	taskPriority := s.getIntValue(taskInfo, "priority", 0)

	for _, device := range devices {
		// 主要考虑设备优先级和任务优先级匹配
		devicePriority := s.getIntValue(device, "priority", 0)
		priorityMatch := 1.0 - math.Abs(float64(taskPriority-devicePriority))/100.0

		compatibilityScore := s.calculateCompatibilityScore(device, taskInfo)
		loadScore := s.calculateDeviceLoadScore(device)

		// 优先级匹配度权重最高
		totalScore := priorityMatch*0.5 + compatibilityScore*0.3 + (1.0-loadScore)*0.2

		if totalScore > bestScore {
			bestScore = totalScore
			bestDevice = device
			reason = "优先级策略选择优先级最匹配的设备"
		}
	}

	if bestDevice == nil {
		return nil, 0, "", gerror.NewCode(gcode.CodeValidationFailed, "没有合适的设备")
	}

	return bestDevice, bestScore, reason, nil
}

// calculateManualAssignment 手动分配策略
func (s *sTask) calculateManualAssignment(devices []map[string]interface{}, taskInfo map[string]interface{}) (map[string]interface{}, float64, string, error) {
	// 手动分配通常选择第一个可用设备
	if len(devices) > 0 {
		device := devices[0]
		score := s.calculateDeviceLoadScore(device)
		return device, score, "手动分配选择第一个可用设备", nil
	}

	return nil, 0, "", gerror.NewCode(gcode.CodeValidationFailed, "没有可用的设备")
}

// calculateDeviceLoadScore 计算设备负载评分
func (s *sTask) calculateDeviceLoadScore(device map[string]interface{}) float64 {
	cpuUsage := s.getFloatValue(device, "cpu_usage", 0.0)
	memoryUsage := s.getFloatValue(device, "memory_usage", 0.0)
	currentTasks := s.getIntValue(device, "current_tasks", 0)
	maxTasks := s.getIntValue(device, "max_tasks", 1)

	// 计算综合负载评分 (0-1, 越高表示负载越重)
	cpuScore := cpuUsage / 100.0
	memoryScore := memoryUsage / 100.0
	taskScore := float64(currentTasks) / float64(maxTasks)

	// 加权平均
	loadScore := cpuScore*0.4 + memoryScore*0.3 + taskScore*0.3
	return math.Min(loadScore, 1.0)
}

// calculateCompatibilityScore 计算兼容性评分
func (s *sTask) calculateCompatibilityScore(device map[string]interface{}, taskInfo map[string]interface{}) float64 {
	deviceType := s.getStringValue(device, "type", "")
	taskType := s.getStringValue(taskInfo, "type", "")
	deviceCapabilities := s.getStringValue(device, "capabilities", "")
	taskCapabilities := s.getStringValue(taskInfo, "required_capabilities", "")

	score := 0.0

	// 类型匹配
	if deviceType == taskType {
		score += 0.5
	}

	// 能力匹配
	if deviceCapabilities != "" && taskCapabilities != "" {
		if s.checkCapabilityCompatibility(taskCapabilities, deviceCapabilities) {
			score += 0.5
		}
	}

	return score
}

// calculateAvailabilityScore 计算可用性评分
func (s *sTask) calculateAvailabilityScore(device map[string]interface{}) float64 {
	status := s.getStringValue(device, "status", "")
	uptime := s.getFloatValue(device, "uptime_hours", 0.0)
	healthScore := s.getFloatValue(device, "health_score", 0.0)

	score := 0.0

	// 状态评分
	switch status {
	case "online":
		score += 0.4
	case "idle":
		score += 0.3
	case "busy":
		score += 0.2
	default:
		score += 0.0
	}

	// 运行时间评分
	if uptime > 24 {
		score += 0.3
	} else if uptime > 1 {
		score += 0.2
	} else {
		score += 0.1
	}

	// 健康评分
	score += healthScore / 100.0 * 0.3

	return math.Min(score, 1.0)
}

// calculatePerformanceScore 计算性能评分
func (s *sTask) calculatePerformanceScore(device map[string]interface{}) float64 {
	// 基于设备性能指标计算评分
	performanceScore := s.getFloatValue(device, "performance_score", 50.0)
	return performanceScore / 100.0
}

// checkCapabilityCompatibility 检查能力兼容性
func (s *sTask) checkCapabilityCompatibility(required, available string) bool {
	// 简单的字符串包含检查，实际应该更复杂的匹配逻辑
	return available != "" && (required == "" || available == required)
}

// GetTaskAssignmentOptions 获取任务分配选项
func (s *sTask) GetTaskAssignmentOptions(ctx context.Context, taskId string) (map[string]interface{}, error) {
	// 验证任务存在
	if err := s.validateTaskExists(ctx, taskId); err != nil {
		return nil, err
	}

	// 获取任务信息
	taskInfo, err := s.getTaskInfo(ctx, taskId)
	if err != nil {
		return nil, err
	}

	// 获取所有可用设备
	allDevices, err := s.getAllAvailableDevices(ctx)
	if err != nil {
		return nil, err
	}

	// 计算每个设备的分配选项
	var options []map[string]interface{}
	for _, device := range allDevices {
		option := s.calculateDeviceAssignmentOption(device, taskInfo)
		if option != nil {
			options = append(options, option)
		}
	}

	// 按分配评分排序
	sort.Slice(options, func(i, j int) bool {
		scoreI := s.getFloatValue(options[i], "assignmentScore", 0.0)
		scoreJ := s.getFloatValue(options[j], "assignmentScore", 0.0)
		return scoreI > scoreJ
	})

	result := map[string]interface{}{
		"taskId":  taskId,
		"options": options,
		"total":   len(options),
	}

	return result, nil
}

// calculateDeviceAssignmentOption 计算设备分配选项
func (s *sTask) calculateDeviceAssignmentOption(device map[string]interface{}, taskInfo map[string]interface{}) map[string]interface{} {
	deviceId := s.getStringValue(device, "device_id", "")
	if deviceId == "" {
		return nil
	}

	// 计算各项评分
	loadScore := s.calculateDeviceLoadScore(device)
	compatibilityScore := s.calculateCompatibilityScore(device, taskInfo)
	availabilityScore := s.calculateAvailabilityScore(device)
	performanceScore := s.calculatePerformanceScore(device)

	// 计算综合分配评分
	assignmentScore := loadScore*0.3 + compatibilityScore*0.3 + availabilityScore*0.2 + performanceScore*0.2

	// 生成推荐
	recommendation := s.generateAssignmentRecommendation(assignmentScore, loadScore, compatibilityScore)

	option := map[string]interface{}{
		"deviceId":           deviceId,
		"deviceName":         s.getStringValue(device, "name", ""),
		"deviceType":         s.getStringValue(device, "type", ""),
		"loadScore":          loadScore * 100,
		"compatibilityScore": compatibilityScore * 100,
		"assignmentScore":    assignmentScore * 100,
		"currentTasks":       s.getIntValue(device, "current_tasks", 0),
		"maxTasks":           s.getIntValue(device, "max_tasks", 1),
		"estimatedWaitTime":  s.calculateEstimatedWaitTime(device),
		"status":             s.getStringValue(device, "status", ""),
		"recommendation":     recommendation,
	}

	return option
}

// generateAssignmentRecommendation 生成分配推荐
func (s *sTask) generateAssignmentRecommendation(assignmentScore, loadScore, compatibilityScore float64) string {
	if assignmentScore > 0.8 {
		return "强烈推荐"
	} else if assignmentScore > 0.6 {
		return "推荐"
	} else if assignmentScore > 0.4 {
		return "一般"
	} else {
		return "不推荐"
	}
}

// calculateEstimatedWaitTime 计算预估等待时间
func (s *sTask) calculateEstimatedWaitTime(device map[string]interface{}) int {
	currentTasks := s.getIntValue(device, "current_tasks", 0)
	avgTaskDuration := s.getFloatValue(device, "avg_task_duration", 300.0) // 默认5分钟

	// 简单计算：当前任务数 * 平均任务时长
	estimatedTime := int(float64(currentTasks) * avgTaskDuration)
	return estimatedTime
}

// BatchAssignTasks 批量分配任务
func (s *sTask) BatchAssignTasks(ctx context.Context, assignments []map[string]interface{}, strategy string) (map[string]interface{}, error) {
	var results []map[string]interface{}
	success := 0
	failed := 0

	for _, assignment := range assignments {
		taskId := s.getStringValue(assignment, "taskId", "")
		deviceIds := s.getStringSliceValue(assignment, "deviceIds", []string{})

		if taskId == "" || len(deviceIds) == 0 {
			failed++
			results = append(results, map[string]interface{}{
				"taskId": taskId,
				"status": "failed",
				"error":  "任务ID或设备ID列表为空",
			})
			continue
		}

		// 执行单个任务分配
		result, err := s.AssignTaskToDevice(ctx, taskId, deviceIds, strategy, false)
		if err != nil {
			failed++
			results = append(results, map[string]interface{}{
				"taskId": taskId,
				"status": "failed",
				"error":  err.Error(),
			})
		} else {
			success++
			results = append(results, result)
		}
	}

	return map[string]interface{}{
		"total":       len(assignments),
		"success":     success,
		"failed":      failed,
		"assignments": results,
	}, nil
}

// executeTaskAssignment 执行任务分配
func (s *sTask) executeTaskAssignment(ctx context.Context, taskId, deviceId string) error {
	// 这里应该执行实际的任务分配逻辑
	// 包括更新数据库、发送通知等

	g.Log().Info(ctx, "执行任务分配", g.Map{
		"task_id":   taskId,
		"device_id": deviceId,
		"time":      gtime.Now(),
	})

	return nil
}

// 工具方法
func (s *sTask) getStringValue(data map[string]interface{}, key string, defaultValue string) string {
	if value, ok := data[key]; ok {
		if strValue, ok := value.(string); ok {
			return strValue
		}
	}
	return defaultValue
}

func (s *sTask) getFloatValue(data map[string]interface{}, key string, defaultValue float64) float64 {
	if value, ok := data[key]; ok {
		if floatValue, ok := value.(float64); ok {
			return floatValue
		}
	}
	return defaultValue
}

func (s *sTask) getIntValue(data map[string]interface{}, key string, defaultValue int) int {
	if value, ok := data[key]; ok {
		if intValue, ok := value.(int); ok {
			return intValue
		}
	}
	return defaultValue
}

func (s *sTask) getStringSliceValue(data map[string]interface{}, key string, defaultValue []string) []string {
	if value, ok := data[key]; ok {
		if sliceValue, ok := value.([]string); ok {
			return sliceValue
		}
	}
	return defaultValue
}

// 模拟方法（实际应该从数据库获取）
func (s *sTask) validateTaskExists(ctx context.Context, taskId string) error {
	// 这里应该验证任务是否存在
	return nil
}

func (s *sTask) getTaskInfo(ctx context.Context, taskId string) (map[string]interface{}, error) {
	// 这里应该从数据库获取任务信息
	return map[string]interface{}{
		"task_id":               taskId,
		"type":                  "data_processing",
		"priority":              5,
		"required_capabilities": "high_performance",
	}, nil
}

func (s *sTask) getAvailableDevices(ctx context.Context, deviceIds []string) ([]map[string]interface{}, error) {
	// 这里应该从数据库获取设备信息
	var devices []map[string]interface{}
	for _, deviceId := range deviceIds {
		device := map[string]interface{}{
			"device_id":         deviceId,
			"name":              "device-" + deviceId,
			"type":              "data_processing",
			"status":            "online",
			"cpu_usage":         45.5,
			"memory_usage":      60.2,
			"current_tasks":     2,
			"max_tasks":         10,
			"capabilities":      "high_performance",
			"priority":          5,
			"uptime_hours":      168.5,
			"health_score":      85.0,
			"performance_score": 78.0,
			"avg_task_duration": 300.0,
		}
		devices = append(devices, device)
	}
	return devices, nil
}

func (s *sTask) getAllAvailableDevices(ctx context.Context) ([]map[string]interface{}, error) {
	// 这里应该从数据库获取所有可用设备
	return []map[string]interface{}{
		{
			"device_id":         "dev001",
			"name":              "device-001",
			"type":              "data_processing",
			"status":            "online",
			"cpu_usage":         45.5,
			"memory_usage":      60.2,
			"current_tasks":     2,
			"max_tasks":         10,
			"capabilities":      "high_performance",
			"priority":          5,
			"uptime_hours":      168.5,
			"health_score":      85.0,
			"performance_score": 78.0,
			"avg_task_duration": 300.0,
		},
		{
			"device_id":         "dev002",
			"name":              "device-002",
			"type":              "data_processing",
			"status":            "idle",
			"cpu_usage":         25.0,
			"memory_usage":      40.0,
			"current_tasks":     0,
			"max_tasks":         8,
			"capabilities":      "high_performance",
			"priority":          3,
			"uptime_hours":      120.0,
			"health_score":      90.0,
			"performance_score": 85.0,
			"avg_task_duration": 250.0,
		},
	}, nil
}

// ===============================
// 辅助函数
// ===============================

// contains 检查切片是否包含指定元素
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
