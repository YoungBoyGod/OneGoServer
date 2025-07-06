// ===============================
// 任务优先级相关Input/Output结构体
// ===============================

// CalculateTaskPriorityInput 计算任务优先级输入
type CalculateTaskPriorityInput struct {
	TaskData map[string]interface{} `json:"task_data"`
}

// CalculateTaskPriorityOutput 计算任务优先级输出
type CalculateTaskPriorityOutput struct {
	Priority int `json:"priority"`
}

// CalculateUrgencyScoreInput 计算紧急度评分输入
type CalculateUrgencyScoreInput struct {
	TaskData map[string]interface{} `json:"task_data"`
}

// CalculateUrgencyScoreOutput 计算紧急度评分输出
type CalculateUrgencyScoreOutput struct {
	Score float64 `json:"score"`
}

// CalculateBusinessImportanceScoreInput 计算业务重要性评分输入
type CalculateBusinessImportanceScoreInput struct {
	TaskData map[string]interface{} `json:"task_data"`
}

// CalculateBusinessImportanceScoreOutput 计算业务重要性评分输出
type CalculateBusinessImportanceScoreOutput struct {
	Score float64 `json:"score"`
}

// CalculateDeadlineScoreInput 计算截止时间评分输入
type CalculateDeadlineScoreInput struct {
	TaskData map[string]interface{} `json:"task_data"`
}

// CalculateDeadlineScoreOutput 计算截止时间评分输出
type CalculateDeadlineScoreOutput struct {
	Score float64 `json:"score"`
}

// CalculateResourceScoreInput 计算资源评分输入
type CalculateResourceScoreInput struct {
	TaskData map[string]interface{} `json:"task_data"`
}

// CalculateResourceScoreOutput 计算资源评分输出
type CalculateResourceScoreOutput struct {
	Score float64 `json:"score"`
} 