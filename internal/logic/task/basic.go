package task

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"

	task "OneGfServer/internal/model/task"
)

// ===============================
// 基础任务管理业务逻辑
// ===============================

// sTask 任务服务实现
type sTask struct{}

// New 创建任务服务实例
func New() *sTask {
	return &sTask{}
}

// CreateTask 创建任务
func (s *sTask) CreateTask(ctx context.Context, input *task.CreateTaskInput) (*task.CreateTaskOutput, error) {
	// 验证任务数据
	validateInput := &task.ValidateTaskInput{
		TaskData: map[string]interface{}{
			"name":        input.Name,
			"type":        input.Type,
			"priority":    input.Priority,
			"description": input.Description,
			"timeout":     input.Timeout,
			"retry_count": input.RetryCount,
		},
	}
	validateOutput := s.validateTask(validateInput)
	if !validateOutput.IsValid {
		return nil, gerror.NewCode(gcode.CodeValidationFailed, validateOutput.Message)
	}

	// 计算任务优先级
	priorityInput := &task.CalculateTaskPriorityInput{
		TaskType:     input.Type,
		Parameters:   input.Parameters,
		UserPriority: input.Priority,
	}
	priorityOutput := s.calculateTaskPriority(priorityInput)

	// 生成任务ID
	taskID := s.generateTaskID()

	// 创建任务记录
	// 这里应该保存到数据库
	// 目前返回模拟结果

	return &task.CreateTaskOutput{
		TaskID:    taskID,
		Status:    "created",
		CreatedAt: gtime.Now().Format("2006-01-02 15:04:05"),
		Message:   "任务创建成功",
	}, nil
}

// GetTask 获取任务
func (s *sTask) GetTask(ctx context.Context, input *task.GetTaskInput) (*task.GetTaskOutput, error) {
	// 这里应该从数据库获取任务
	// 目前返回模拟数据
	taskData := map[string]interface{}{
		"id":          input.TaskID,
		"name":        "示例任务",
		"type":        "data_processing",
		"priority":    5,
		"status":      "pending",
		"description": "这是一个示例任务",
		"parameters":  map[string]interface{}{},
		"timeout":     3600,
		"retry_count": 3,
		"progress":    0.0,
		"created_at":  gtime.Now().Add(-time.Hour).Format("2006-01-02 15:04:05"),
		"updated_at":  gtime.Now().Format("2006-01-02 15:04:05"),
	}

	return &task.GetTaskOutput{
		Task: taskData,
	}, nil
}

// UpdateTask 更新任务
func (s *sTask) UpdateTask(ctx context.Context, input *task.UpdateTaskInput) (*task.UpdateTaskOutput, error) {
	// 验证任务数据
	validateInput := &task.ValidateTaskInput{
		TaskData: map[string]interface{}{
			"name":        input.Name,
			"type":        input.Type,
			"priority":    input.Priority,
			"description": input.Description,
			"timeout":     input.Timeout,
			"retry_count": input.RetryCount,
		},
	}
	validateOutput := s.validateTask(validateInput)
	if !validateOutput.IsValid {
		return nil, gerror.NewCode(gcode.CodeValidationFailed, validateOutput.Message)
	}

	// 这里应该更新数据库中的任务
	// 目前返回模拟结果

	return &task.UpdateTaskOutput{
		TaskID:    input.TaskID,
		Status:    "updated",
		UpdatedAt: gtime.Now().Format("2006-01-02 15:04:05"),
		Message:   "任务更新成功",
	}, nil
}

// DeleteTask 删除任务
func (s *sTask) DeleteTask(ctx context.Context, input *task.DeleteTaskInput) (*task.DeleteTaskOutput, error) {
	// 这里应该从数据库删除任务
	// 目前返回模拟结果

	return &task.DeleteTaskOutput{
		TaskID:  input.TaskID,
		Status:  "deleted",
		Message: "任务删除成功",
	}, nil
}

// ListTasks 任务列表
func (s *sTask) ListTasks(ctx context.Context, input *task.ListTasksInput) (*task.ListTasksOutput, error) {
	// 这里应该从数据库获取任务列表
	// 目前返回模拟数据
	var tasks []map[string]interface{}
	for i := 0; i < input.Size; i++ {
		task := map[string]interface{}{
			"id":          "task-" + string(rune(i+1)),
			"name":        "任务 " + string(rune(i+1)),
			"type":        "data_processing",
			"priority":    5,
			"status":      "pending",
			"description": "任务描述",
			"created_at":  gtime.Now().Add(-time.Duration(i) * time.Hour).Format("2006-01-02 15:04:05"),
		}
		tasks = append(tasks, task)
	}

	return &task.ListTasksOutput{
		List:  tasks,
		Total: 100,
		Page:  input.Page,
		Size:  input.Size,
	}, nil
}

// validateTask 验证任务数据
func (s *sTask) validateTask(input *task.ValidateTaskInput) *task.ValidateTaskOutput {
	var errors []string

	// 检查任务名称
	if name, ok := input.TaskData["name"].(string); ok {
		if name == "" {
			errors = append(errors, "任务名称不能为空")
		}
	} else {
		errors = append(errors, "任务名称不能为空")
	}

	// 检查任务类型
	if taskType, ok := input.TaskData["type"].(string); ok {
		if taskType == "" {
			errors = append(errors, "任务类型不能为空")
		}
	} else {
		errors = append(errors, "任务类型不能为空")
	}

	// 检查优先级
	if priority, ok := input.TaskData["priority"].(int); ok {
		if priority < 1 || priority > 10 {
			errors = append(errors, "优先级必须在1-10之间")
		}
	}

	// 检查超时时间
	if timeout, ok := input.TaskData["timeout"].(int); ok {
		if timeout <= 0 || timeout > 86400 {
			errors = append(errors, "超时时间必须在1-86400秒之间")
		}
	}

	isValid := len(errors) == 0
	message := "验证通过"
	if !isValid {
		message = "验证失败"
	}

	return &task.ValidateTaskOutput{
		IsValid: isValid,
		Message: message,
		Errors:  errors,
	}
}

// calculateTaskPriority 计算任务优先级
func (s *sTask) calculateTaskPriority(input *task.CalculateTaskPriorityInput) *task.CalculateTaskPriorityOutput {
	priority := input.UserPriority

	// 根据任务类型调整优先级
	switch input.TaskType {
	case "urgent":
		priority += 3
	case "high":
		priority += 2
	case "normal":
		priority += 1
	case "low":
		priority -= 1
	}

	// 根据参数调整优先级
	if input.Parameters != nil {
		if _, ok := input.Parameters["critical"]; ok {
			priority += 2
		}
	}

	// 确保优先级在有效范围内
	if priority < 1 {
		priority = 1
	} else if priority > 10 {
		priority = 10
	}

	return &task.CalculateTaskPriorityOutput{
		Priority:     priority,
		Reason:       "基于任务类型和参数计算",
		CalculatedAt: gtime.Now().Format("2006-01-02 15:04:05"),
	}
}

// generateTaskID 生成任务ID
func (s *sTask) generateTaskID() string {
	return "task-" + gtime.Now().Format("20060102150405") + "-" + string(rune(gtime.Now().UnixNano()%1000))
}
