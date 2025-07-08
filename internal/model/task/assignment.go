package task

// ===============================
// 任务分配相关Input/Output结构体
// ===============================

// AssignTaskToDeviceInput 分配任务到设备输入
type AssignTaskToDeviceInput struct {
	TaskId    string   `json:"task_id"`
	DeviceIds []string `json:"device_ids"`
	Strategy  string   `json:"strategy"`
	Force     bool     `json:"force"`
}

// AssignTaskToDeviceOutput 分配任务到设备输出
type AssignTaskToDeviceOutput struct {
	Result map[string]interface{} `json:"result"`
}

// SelectBestDeviceInput 选择最佳设备输入
type SelectBestDeviceInput struct {
	TaskInfo map[string]interface{}   `json:"task_info"`
	Devices  []map[string]interface{} `json:"devices"`
	Strategy string                   `json:"strategy"`
	Force    bool                     `json:"force"`
}

// SelectBestDeviceOutput 选择最佳设备输出
type SelectBestDeviceOutput struct {
	Device map[string]interface{} `json:"device"`
	Score  float64                `json:"score"`
	Reason string                 `json:"reason"`
}

// FilterCompatibleDevicesInput 过滤兼容设备输入
type FilterCompatibleDevicesInput struct {
	TaskInfo map[string]interface{}   `json:"task_info"`
	Devices  []map[string]interface{} `json:"devices"`
}

// FilterCompatibleDevicesOutput 过滤兼容设备输出
type FilterCompatibleDevicesOutput struct {
	CompatibleDevices []map[string]interface{} `json:"compatible_devices"`
}

// CalculateAutoAssignmentInput 计算自动分配输入
type CalculateAutoAssignmentInput struct {
	Devices  []map[string]interface{} `json:"devices"`
	TaskInfo map[string]interface{}   `json:"task_info"`
}

// CalculateAutoAssignmentOutput 计算自动分配输出
type CalculateAutoAssignmentOutput struct {
	Device map[string]interface{} `json:"device"`
	Score  float64                `json:"score"`
	Reason string                 `json:"reason"`
}

// CalculateLoadBalancedAssignmentInput 计算负载均衡分配输入
type CalculateLoadBalancedAssignmentInput struct {
	Devices  []map[string]interface{} `json:"devices"`
	TaskInfo map[string]interface{}   `json:"task_info"`
}

// CalculateLoadBalancedAssignmentOutput 计算负载均衡分配输出
type CalculateLoadBalancedAssignmentOutput struct {
	Device map[string]interface{} `json:"device"`
	Score  float64                `json:"score"`
	Reason string                 `json:"reason"`
}

// CalculatePriorityAssignmentInput 计算优先级分配输入
type CalculatePriorityAssignmentInput struct {
	Devices  []map[string]interface{} `json:"devices"`
	TaskInfo map[string]interface{}   `json:"task_info"`
}

// CalculatePriorityAssignmentOutput 计算优先级分配输出
type CalculatePriorityAssignmentOutput struct {
	Device map[string]interface{} `json:"device"`
	Score  float64                `json:"score"`
	Reason string                 `json:"reason"`
}

// CalculateManualAssignmentInput 计算手动分配输入
type CalculateManualAssignmentInput struct {
	Devices  []map[string]interface{} `json:"devices"`
	TaskInfo map[string]interface{}   `json:"task_info"`
}

// CalculateManualAssignmentOutput 计算手动分配输出
type CalculateManualAssignmentOutput struct {
	Device map[string]interface{} `json:"device"`
	Score  float64                `json:"score"`
	Reason string                 `json:"reason"`
}

// CalculateDeviceAssignmentScoreInput 计算设备分配评分输入
type CalculateDeviceAssignmentScoreInput struct {
	Device   map[string]interface{} `json:"device"`
	TaskInfo map[string]interface{} `json:"task_info"`
}

// CalculateDeviceAssignmentScoreOutput 计算设备分配评分输出
type CalculateDeviceAssignmentScoreOutput struct {
	Score float64 `json:"score"`
}

// CalculateDeviceLoadScoreInput 计算设备负载评分输入
type CalculateDeviceLoadScoreInput struct {
	Device map[string]interface{} `json:"device"`
}

// CalculateDeviceLoadScoreOutput 计算设备负载评分输出
type CalculateDeviceLoadScoreOutput struct {
	Score float64 `json:"score"`
}

// CalculateCompatibilityScoreInput 计算兼容性评分输入
type CalculateCompatibilityScoreInput struct {
	Device   map[string]interface{} `json:"device"`
	TaskInfo map[string]interface{} `json:"task_info"`
}

// CalculateCompatibilityScoreOutput 计算兼容性评分输出
type CalculateCompatibilityScoreOutput struct {
	Score float64 `json:"score"`
}

// CalculateAvailabilityScoreInput 计算可用性评分输入
type CalculateAvailabilityScoreInput struct {
	Device map[string]interface{} `json:"device"`
}

// CalculateAvailabilityScoreOutput 计算可用性评分输出
type CalculateAvailabilityScoreOutput struct {
	Score float64 `json:"score"`
}

// CalculatePerformanceScoreInput 计算性能评分输入
type CalculatePerformanceScoreInput struct {
	Device map[string]interface{} `json:"device"`
}

// CalculatePerformanceScoreOutput 计算性能评分输出
type CalculatePerformanceScoreOutput struct {
	Score float64 `json:"score"`
}

// CheckResourceAvailabilityInput 检查资源可用性输入
type CheckResourceAvailabilityInput struct {
	Device   map[string]interface{} `json:"device"`
	TaskInfo map[string]interface{} `json:"task_info"`
}

// CheckResourceAvailabilityOutput 检查资源可用性输出
type CheckResourceAvailabilityOutput struct {
	IsAvailable bool   `json:"is_available"`
	Reason      string `json:"reason,omitempty"`
}

// CheckResourceCompatibilityInput 检查资源兼容性输入
type CheckResourceCompatibilityInput struct {
	Device   map[string]interface{} `json:"device"`
	TaskInfo map[string]interface{} `json:"task_info"`
}

// CheckResourceCompatibilityOutput 检查资源兼容性输出
type CheckResourceCompatibilityOutput struct {
	IsCompatible bool   `json:"is_compatible"`
	Reason       string `json:"reason,omitempty"`
}

// GetTaskInfoInput 获取任务信息输入
type GetTaskInfoInput struct {
	TaskId string `json:"task_id"`
}

// GetTaskInfoOutput 获取任务信息输出
type GetTaskInfoOutput struct {
	TaskInfo map[string]interface{} `json:"task_info"`
}

// GetAvailableDevicesInput 获取可用设备输入
type GetAvailableDevicesInput struct {
	DeviceIds []string `json:"device_ids"`
}

// GetAvailableDevicesOutput 获取可用设备输出
type GetAvailableDevicesOutput struct {
	Devices []map[string]interface{} `json:"devices"`
}

// GetAllAvailableDevicesInput 获取所有可用设备输入
type GetAllAvailableDevicesInput struct{}

// GetAllAvailableDevicesOutput 获取所有可用设备输出
type GetAllAvailableDevicesOutput struct {
	Devices []map[string]interface{} `json:"devices"`
}

// ExecuteTaskAssignmentInput 执行任务分配输入
type ExecuteTaskAssignmentInput struct {
	TaskId   string `json:"task_id"`
	DeviceId string `json:"device_id"`
}

// ExecuteTaskAssignmentOutput 执行任务分配输出
type ExecuteTaskAssignmentOutput struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}
