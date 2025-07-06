// ===============================
// 任务状态相关Input/Output结构体
// ===============================

// ValidateTaskStatusTransitionInput 验证任务状态转换输入
type ValidateTaskStatusTransitionInput struct {
	CurrentStatus string `json:"current_status"`
	TargetStatus  string `json:"target_status"`
}

// ValidateTaskStatusTransitionOutput 验证任务状态转换输出
type ValidateTaskStatusTransitionOutput struct {
	IsValid bool   `json:"is_valid"`
	Message string `json:"message,omitempty"`
}

// DetermineTaskStatusInput 确定任务状态输入
type DetermineTaskStatusInput struct {
	TaskData map[string]interface{} `json:"task_data"`
}

// DetermineTaskStatusOutput 确定任务状态输出
type DetermineTaskStatusOutput struct {
	Status string `json:"status"`
}

// CanTransitionToStatusInput 检查是否可以转换到指定状态输入
type CanTransitionToStatusInput struct {
	TaskData     map[string]interface{} `json:"task_data"`
	TargetStatus string                 `json:"target_status"`
}

// CanTransitionToStatusOutput 检查是否可以转换到指定状态输出
type CanTransitionToStatusOutput struct {
	CanTransition bool   `json:"can_transition"`
	Reason        string `json:"reason,omitempty"`
} 