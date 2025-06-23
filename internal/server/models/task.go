package models

import (
	"time"
)

// TaskStatus 任务状态枚举
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"   // 待执行
	TaskStatusRunning   TaskStatus = "running"   // 执行中
	TaskStatusCompleted TaskStatus = "completed" // 已完成
	TaskStatusFailed    TaskStatus = "failed"    // 执行失败
	TaskStatusTimeout   TaskStatus = "timeout"   // 执行超时
)

// ScriptType 脚本类型枚举
type ScriptType string

const (
	ScriptTypeShell  ScriptType = "shell"  // Shell脚本
	ScriptTypePython ScriptType = "python" // Python脚本
	ScriptTypeNode   ScriptType = "node"   // Node.js脚本
	ScriptTypeBatch  ScriptType = "batch"  // Windows批处理
)

// Task 任务数据模型
type Task struct {
	ID            string     `json:"id"`             // 任务唯一标识
	Name          string     `json:"name"`           // 任务名称
	Description   string     `json:"description"`    // 任务描述
	ScriptType    ScriptType `json:"script_type"`    // 脚本类型
	ScriptContent string     `json:"script_content"` // 脚本内容
	TargetClients []string   `json:"target_clients"` // 目标客户端ID列表，空表示所有客户端
	Priority      int        `json:"priority"`       // 优先级（1-10，数字越大优先级越高）
	Timeout       int        `json:"timeout"`        // 超时时间（秒）
	Status        TaskStatus `json:"status"`         // 任务状态
	CreatedBy     string     `json:"created_by"`     // 创建者
	CreatedAt     time.Time  `json:"created_at"`     // 创建时间
	UpdatedAt     time.Time  `json:"updated_at"`     // 更新时间
	ScheduledAt   *time.Time `json:"scheduled_at"`   // 计划执行时间（可选）
}

// TaskResult 任务执行结果
type TaskResult struct {
	ID        string     `json:"id"`         // 结果唯一标识
	TaskID    string     `json:"task_id"`    // 任务ID
	ClientID  string     `json:"client_id"`  // 客户端ID
	Status    TaskStatus `json:"status"`     // 执行状态
	Output    string     `json:"output"`     // 执行输出
	Error     string     `json:"error"`      // 错误信息
	StartTime time.Time  `json:"start_time"` // 开始执行时间
	EndTime   time.Time  `json:"end_time"`   // 结束执行时间
	Duration  int64      `json:"duration"`   // 执行耗时（毫秒）
}

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	Name          string     `json:"name" binding:"required"`
	Description   string     `json:"description"`
	ScriptType    ScriptType `json:"script_type" binding:"required"`
	ScriptContent string     `json:"script_content" binding:"required"`
	TargetClients []string   `json:"target_clients"`
	Priority      int        `json:"priority"`
	Timeout       int        `json:"timeout"`
	ScheduledAt   *time.Time `json:"scheduled_at"`
}

// UpdateTaskStatusRequest 更新任务状态请求
type UpdateTaskStatusRequest struct {
	Status    TaskStatus `json:"status" binding:"required"`
	Output    string     `json:"output"`
	Error     string     `json:"error"`
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
}

// TaskListResponse 任务列表响应
type TaskListResponse struct {
	Tasks []Task `json:"tasks"`
	Total int    `json:"total"`
}

// ValidateScriptType 验证脚本类型是否有效
func (st ScriptType) IsValid() bool {
	switch st {
	case ScriptTypeShell, ScriptTypePython, ScriptTypeNode, ScriptTypeBatch:
		return true
	default:
		return false
	}
}

// ValidateTaskStatus 验证任务状态是否有效
func (ts TaskStatus) IsValid() bool {
	switch ts {
	case TaskStatusPending, TaskStatusRunning, TaskStatusCompleted, TaskStatusFailed, TaskStatusTimeout:
		return true
	default:
		return false
	}
}

// GetFileExtension 根据脚本类型获取文件扩展名
func (st ScriptType) GetFileExtension() string {
	switch st {
	case ScriptTypeShell:
		return ".sh"
	case ScriptTypePython:
		return ".py"
	case ScriptTypeNode:
		return ".js"
	case ScriptTypeBatch:
		return ".bat"
	default:
		return ".txt"
	}
}

// GetExecutor 根据脚本类型获取执行器命令
func (st ScriptType) GetExecutor() string {
	switch st {
	case ScriptTypeShell:
		return "bash"
	case ScriptTypePython:
		return "python3"
	case ScriptTypeNode:
		return "node"
	case ScriptTypeBatch:
		return "cmd"
	default:
		return ""
	}
}
