package device

// ===============================
// Logic层专用 Input/Output 结构体
// ===============================

// ValidateDeviceRegistrationInput 验证设备注册输入
type ValidateDeviceRegistrationInput struct {
	Name       string `json:"name"`        // 设备名称
	MacAddress string `json:"mac_address"` // MAC地址
	IPAddress  string `json:"ip_address"`  // IP地址
	DeviceType string `json:"device_type"` // 设备类型
	Model      string `json:"model"`       // 设备型号
	Protocol   string `json:"protocol"`    // 通信协议
	Port       int    `json:"port"`        // 端口
}

// ValidateDeviceRegistrationOutput 验证设备注册输出
type ValidateDeviceRegistrationOutput struct {
	IsValid bool     `json:"is_valid"` // 是否有效
	Message string   `json:"message"`  // 验证消息
	Errors  []string `json:"errors"`   // 错误列表
}

// ValidateDeviceStatusInput 验证设备状态转换输入
type ValidateDeviceStatusInput struct {
	CurrentStatus string `json:"current_status"` // 当前状态
	TargetStatus  string `json:"target_status"`  // 目标状态
}

// ValidateDeviceStatusOutput 验证设备状态转换输出
type ValidateDeviceStatusOutput struct {
	IsValid bool   `json:"is_valid"` // 是否有效
	Message string `json:"message"`  // 验证消息
}

// ValidateDeviceConfigurationInput 验证设备配置输入
type ValidateDeviceConfigurationInput struct {
	BasicConfig       *BasicConfig       `json:"basic_config"`       // 基础配置
	PerformanceConfig *PerformanceConfig `json:"performance_config"` // 性能配置
	SecurityConfig    *SecurityConfig    `json:"security_config"`    // 安全配置
}

// ValidateDeviceConfigurationOutput 验证设备配置输出
type ValidateDeviceConfigurationOutput struct {
	IsValid bool     `json:"is_valid"` // 是否有效
	Message string   `json:"message"`  // 验证消息
	Errors  []string `json:"errors"`   // 错误列表
}

// HandleDeviceRegistrationInput 处理设备注册输入
type HandleDeviceRegistrationInput struct {
	DeviceData map[string]interface{} `json:"device_data"` // 设备数据
}

// HandleDeviceRegistrationOutput 处理设备注册输出
type HandleDeviceRegistrationOutput struct {
	DeviceID   string                 `json:"device_id"`   // 设备ID
	DeviceData map[string]interface{} `json:"device_data"` // 处理后的设备数据
	Message    string                 `json:"message"`     // 处理消息
	IsSuccess  bool                   `json:"is_success"`  // 是否成功
}

// HandleDeviceHeartbeatInput 处理设备心跳输入
type HandleDeviceHeartbeatInput struct {
	DeviceID      string                 `json:"device_id"`      // 设备ID
	HeartbeatData map[string]interface{} `json:"heartbeat_data"` // 心跳数据
}

// HandleDeviceHeartbeatOutput 处理设备心跳输出
type HandleDeviceHeartbeatOutput struct {
	DeviceID      string                 `json:"device_id"`      // 设备ID
	HeartbeatData map[string]interface{} `json:"heartbeat_data"` // 处理后的心跳数据
	HealthScore   float64                `json:"health_score"`   // 健康度评分
	Status        string                 `json:"status"`         // 设备状态
	Message       string                 `json:"message"`        // 处理消息
}

// HandleDeviceDeactivationInput 处理设备停用输入
type HandleDeviceDeactivationInput struct {
	DeviceID string `json:"device_id"` // 设备ID
	Reason   string `json:"reason"`    // 停用原因
}

// HandleDeviceDeactivationOutput 处理设备停用输出
type HandleDeviceDeactivationOutput struct {
	DeviceID  string `json:"device_id"`  // 设备ID
	Message   string `json:"message"`    // 处理消息
	IsSuccess bool   `json:"is_success"` // 是否成功
}

// CalculateDeviceHealthScoreInput 计算设备健康度评分输入
type CalculateDeviceHealthScoreInput struct {
	DeviceData map[string]interface{} `json:"device_data"` // 设备数据
}

// CalculateDeviceHealthScoreOutput 计算设备健康度评分输出
type CalculateDeviceHealthScoreOutput struct {
	HealthScore float64                `json:"health_score"` // 健康度评分
	Components  map[string]interface{} `json:"components"`   // 评分组件
}

// DetermineDeviceStatusInput 确定设备状态输入
type DetermineDeviceStatusInput struct {
	DeviceData map[string]interface{} `json:"device_data"` // 设备数据
}

// DetermineDeviceStatusOutput 确定设备状态输出
type DetermineDeviceStatusOutput struct {
	Status      string  `json:"status"`       // 设备状态
	Reason      string  `json:"reason"`       // 状态确定原因
	HealthScore float64 `json:"health_score"` // 健康度评分
}

// CanAcceptNewTaskInput 判断设备是否可以接受新任务输入
type CanAcceptNewTaskInput struct {
	DeviceData map[string]interface{} `json:"device_data"` // 设备数据
}

// CanAcceptNewTaskOutput 判断设备是否可以接受新任务输出
type CanAcceptNewTaskOutput struct {
	CanAccept bool    `json:"can_accept"` // 是否可以接受
	Reason    string  `json:"reason"`     // 原因
	Score     float64 `json:"score"`      // 接受评分
}

// CalculateTaskAssignmentScoreInput 计算任务分配评分输入
type CalculateTaskAssignmentScoreInput struct {
	DeviceData map[string]interface{} `json:"device_data"` // 设备数据
	TaskData   map[string]interface{} `json:"task_data"`   // 任务数据
}

// CalculateTaskAssignmentScoreOutput 计算任务分配评分输出
type CalculateTaskAssignmentScoreOutput struct {
	Score          float64                `json:"score"`          // 分配评分
	Components     map[string]interface{} `json:"components"`     // 评分组件
	Recommendation string                 `json:"recommendation"` // 分配建议
}

// CheckDeviceHealthInput 检查设备健康状态输入
type CheckDeviceHealthInput struct {
	DeviceData map[string]interface{} `json:"device_data"` // 设备数据
}

// CheckDeviceHealthOutput 检查设备健康状态输出
type CheckDeviceHealthOutput struct {
	HealthScore     float64                  `json:"health_score"`    // 健康度评分
	Status          string                   `json:"status"`          // 健康状态
	Checks          map[string]interface{}   `json:"checks"`          // 健康检查结果
	Alerts          []map[string]interface{} `json:"alerts"`          // 健康告警
	Recommendations []string                 `json:"recommendations"` // 优化建议
	Timestamp       string                   `json:"timestamp"`       // 检查时间
}

// CalculateDeviceLoadScoreInput 计算设备负载评分输入
type CalculateDeviceLoadScoreInput struct {
	DeviceID string `json:"device_id"` // 设备ID
}

// CalculateDeviceLoadScoreOutput 计算设备负载评分输出
type CalculateDeviceLoadScoreOutput struct {
	DeviceID        string                 `json:"device_id"`       // 设备ID
	LoadScore       float64                `json:"load_score"`      // 负载评分
	Components      map[string]interface{} `json:"components"`      // 负载组件
	Recommendations []string               `json:"recommendations"` // 优化建议
	Timestamp       string                 `json:"timestamp"`       // 计算时间
	Metrics         map[string]interface{} `json:"metrics"`         // 负载指标
}

// OptimizeDeviceLoadInput 优化设备负载输入
type OptimizeDeviceLoadInput struct {
	DeviceID string `json:"device_id"` // 设备ID
	Strategy string `json:"strategy"`  // 优化策略
}

// OptimizeDeviceLoadOutput 优化设备负载输出
type OptimizeDeviceLoadOutput struct {
	DeviceID       string   `json:"device_id"`       // 设备ID
	Strategy       string   `json:"strategy"`        // 优化策略
	CurrentScore   float64  `json:"current_score"`   // 当前评分
	OptimizedScore float64  `json:"optimized_score"` // 优化后评分
	Improvement    float64  `json:"improvement"`     // 改善程度
	Actions        []string `json:"actions"`         // 优化动作
	Timestamp      string   `json:"timestamp"`       // 优化时间
}

// SetDeviceLoadThresholdInput 设置设备负载阈值输入
type SetDeviceLoadThresholdInput struct {
	DeviceID string  `json:"device_id"` // 设备ID
	Warning  float64 `json:"warning"`   // 警告阈值
	Critical float64 `json:"critical"`  // 严重阈值
	MaxTasks int     `json:"max_tasks"` // 最大任务数
}

// SetDeviceLoadThresholdOutput 设置设备负载阈值输出
type SetDeviceLoadThresholdOutput struct {
	DeviceID  string `json:"device_id"`  // 设备ID
	Message   string `json:"message"`    // 设置消息
	IsSuccess bool   `json:"is_success"` // 是否成功
}

// GetDeviceLoadThresholdInput 获取设备负载阈值输入
type GetDeviceLoadThresholdInput struct {
	DeviceID string `json:"device_id"` // 设备ID
}

// GetDeviceLoadThresholdOutput 获取设备负载阈值输出
type GetDeviceLoadThresholdOutput struct {
	DeviceID string  `json:"device_id"` // 设备ID
	Warning  float64 `json:"warning"`   // 警告阈值
	Critical float64 `json:"critical"`  // 严重阈值
	MaxTasks int     `json:"max_tasks"` // 最大任务数
}

// GetDeviceLoadMetricsInput 获取设备负载指标输入
type GetDeviceLoadMetricsInput struct {
	DeviceID string `json:"device_id"` // 设备ID
	Period   string `json:"period"`    // 时间周期
}

// GetDeviceLoadMetricsOutput 获取设备负载指标输出
type GetDeviceLoadMetricsOutput struct {
	DeviceID   string                   `json:"device_id"`   // 设备ID
	Period     string                   `json:"period"`      // 时间周期
	Duration   string                   `json:"duration"`    // 持续时间
	TimeSeries []map[string]interface{} `json:"time_series"` // 时间序列数据
	History    map[string]interface{}   `json:"history"`     // 历史数据
	Summary    map[string]interface{}   `json:"summary"`     // 数据摘要
	Timestamp  string                   `json:"timestamp"`   // 获取时间
}

// AnalyzeDevicePerformanceInput 分析设备性能输入
type AnalyzeDevicePerformanceInput struct {
	PerformanceData []map[string]interface{} `json:"performance_data"` // 性能数据
}

// AnalyzeDevicePerformanceOutput 分析设备性能输出
type AnalyzeDevicePerformanceOutput struct {
	BasicStats       map[string]interface{}   `json:"basic_stats"`       // 基础统计
	Trends           map[string]interface{}   `json:"trends"`            // 性能趋势
	Bottlenecks      []map[string]interface{} `json:"bottlenecks"`       // 性能瓶颈
	PerformanceScore float64                  `json:"performance_score"` // 性能评分
	Recommendations  []string                 `json:"recommendations"`   // 优化建议
	TimeRange        map[string]interface{}   `json:"time_range"`        // 时间范围
}

// ProcessDeviceDataInput 处理设备数据输入
type ProcessDeviceDataInput struct {
	RawData map[string]interface{} `json:"raw_data"` // 原始数据
}

// ProcessDeviceDataOutput 处理设备数据输出
type ProcessDeviceDataOutput struct {
	ProcessedData map[string]interface{} `json:"processed_data"` // 处理后的数据
	IsValid       bool                   `json:"is_valid"`       // 是否有效
	Message       string                 `json:"message"`        // 处理消息
}

// ValidateDeviceDataInput 验证设备数据输入
type ValidateDeviceDataInput struct {
	DeviceData map[string]interface{} `json:"device_data"` // 设备数据
}

// ValidateDeviceDataOutput 验证设备数据输出
type ValidateDeviceDataOutput struct {
	IsValid bool     `json:"is_valid"` // 是否有效
	Message string   `json:"message"`  // 验证消息
	Errors  []string `json:"errors"`   // 错误列表
}

// ===============================
// 配置结构体
// ===============================

// BasicConfig 基础配置
type BasicConfig struct {
	Name        string   `json:"name"`        // 设备名称
	Description string   `json:"description"` // 设备描述
	Tags        []string `json:"tags"`        // 设备标签
}

// PerformanceConfig 性能配置
type PerformanceConfig struct {
	MaxConcurrentTasks int     `json:"max_concurrent_tasks"` // 最大并发任务数
	CPUThreshold       float64 `json:"cpu_threshold"`        // CPU阈值
	MemoryThreshold    float64 `json:"memory_threshold"`     // 内存阈值
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	AccessControl *AccessControl `json:"access_control"` // 访问控制
	Encryption    *Encryption    `json:"encryption"`     // 加密设置
}

// AccessControl 访问控制
type AccessControl struct {
	AllowedIPs      []string `json:"allowed_ips"`      // 允许的IP地址
	UserPermissions []string `json:"user_permissions"` // 用户权限
}

// Encryption 加密设置
type Encryption struct {
	Algorithm string `json:"algorithm"`  // 加密算法
	KeyLength int    `json:"key_length"` // 密钥长度
}

// ===============================
// 性能分析相关 Input/Output
// ===============================

// CalculateBasicStatsInput 计算基础统计信息输入
type CalculateBasicStatsInput struct {
	PerformanceData []map[string]interface{} `json:"performance_data"` // 性能数据
}

// CalculateBasicStatsOutput 计算基础统计信息输出
type CalculateBasicStatsOutput struct {
	Stats map[string]interface{} `json:"stats"` // 统计信息
}

// AnalyzePerformanceTrendsInput 分析性能趋势输入
type AnalyzePerformanceTrendsInput struct {
	PerformanceData []map[string]interface{} `json:"performance_data"` // 性能数据
}

// AnalyzePerformanceTrendsOutput 分析性能趋势输出
type AnalyzePerformanceTrendsOutput struct {
	Trends map[string]interface{} `json:"trends"` // 趋势分析结果
}

// AnalyzeMetricTrendInput 分析单个指标趋势输入
type AnalyzeMetricTrendInput struct {
	SortedData []map[string]interface{} `json:"sorted_data"` // 排序后的数据
	MetricKey  string                   `json:"metric_key"`  // 指标键名
}

// AnalyzeMetricTrendOutput 分析单个指标趋势输出
type AnalyzeMetricTrendOutput struct {
	Trend      string  `json:"trend"`      // 趋势方向
	Slope      float64 `json:"slope"`      // 斜率
	Volatility float64 `json:"volatility"` // 波动性
}

// IdentifyBottlenecksInput 识别性能瓶颈输入
type IdentifyBottlenecksInput struct {
	PerformanceData []map[string]interface{} `json:"performance_data"` // 性能数据
}

// IdentifyBottlenecksOutput 识别性能瓶颈输出
type IdentifyBottlenecksOutput struct {
	Bottlenecks []map[string]interface{} `json:"bottlenecks"` // 瓶颈列表
}

// CheckCPUBottleneckInput 检查CPU瓶颈输入
type CheckCPUBottleneckInput struct {
	PerformanceData []map[string]interface{} `json:"performance_data"` // 性能数据
}

// CheckCPUBottleneckOutput 检查CPU瓶颈输出
type CheckCPUBottleneckOutput struct {
	Bottleneck map[string]interface{} `json:"bottleneck"` // 瓶颈信息
	HasIssue   bool                   `json:"has_issue"`  // 是否存在问题
}

// CheckMemoryBottleneckInput 检查内存瓶颈输入
type CheckMemoryBottleneckInput struct {
	PerformanceData []map[string]interface{} `json:"performance_data"` // 性能数据
}

// CheckMemoryBottleneckOutput 检查内存瓶颈输出
type CheckMemoryBottleneckOutput struct {
	Bottleneck map[string]interface{} `json:"bottleneck"` // 瓶颈信息
	HasIssue   bool                   `json:"has_issue"`  // 是否存在问题
}

// CheckDiskBottleneckInput 检查磁盘瓶颈输入
type CheckDiskBottleneckInput struct {
	PerformanceData []map[string]interface{} `json:"performance_data"` // 性能数据
}

// CheckDiskBottleneckOutput 检查磁盘瓶颈输出
type CheckDiskBottleneckOutput struct {
	Bottleneck map[string]interface{} `json:"bottleneck"` // 瓶颈信息
	HasIssue   bool                   `json:"has_issue"`  // 是否存在问题
}

// CheckNetworkBottleneckInput 检查网络瓶颈输入
type CheckNetworkBottleneckInput struct {
	PerformanceData []map[string]interface{} `json:"performance_data"` // 性能数据
}

// CheckNetworkBottleneckOutput 检查网络瓶颈输出
type CheckNetworkBottleneckOutput struct {
	Bottleneck map[string]interface{} `json:"bottleneck"` // 瓶颈信息
	HasIssue   bool                   `json:"has_issue"`  // 是否存在问题
}

// CalculatePerformanceScoreFromDataInput 计算性能评分输入
type CalculatePerformanceScoreFromDataInput struct {
	PerformanceData []map[string]interface{} `json:"performance_data"` // 性能数据
}

// CalculatePerformanceScoreFromDataOutput 计算性能评分输出
type CalculatePerformanceScoreFromDataOutput struct {
	Score float64 `json:"score"` // 性能评分
}

// GeneratePerformanceRecommendationsInput 生成性能优化建议输入
type GeneratePerformanceRecommendationsInput struct {
	Analysis map[string]interface{} `json:"analysis"` // 分析结果
}

// GeneratePerformanceRecommendationsOutput 生成性能优化建议输出
type GeneratePerformanceRecommendationsOutput struct {
	Recommendations []string `json:"recommendations"` // 优化建议列表
}

// ExtractValuesInput 提取数值输入
type ExtractValuesInput struct {
	Data []map[string]interface{} `json:"data"` // 数据
	Key  string                   `json:"key"`  // 键名
}

// ExtractValuesOutput 提取数值输出
type ExtractValuesOutput struct {
	Values []float64 `json:"values"` // 数值列表
}

// SortByTimestampInput 按时间戳排序输入
type SortByTimestampInput struct {
	Data []map[string]interface{} `json:"data"` // 数据
}

// SortByTimestampOutput 按时间戳排序输出
type SortByTimestampOutput struct {
	SortedData []map[string]interface{} `json:"sorted_data"` // 排序后的数据
}

// CalculateLinearRegressionSlopeInput 计算线性回归斜率输入
type CalculateLinearRegressionSlopeInput struct {
	Values []float64 `json:"values"` // 数值列表
}

// CalculateLinearRegressionSlopeOutput 计算线性回归斜率输出
type CalculateLinearRegressionSlopeOutput struct {
	Slope float64 `json:"slope"` // 斜率
}

// CalculateVolatilityInput 计算波动性输入
type CalculateVolatilityInput struct {
	Values []float64 `json:"values"` // 数值列表
}

// CalculateVolatilityOutput 计算波动性输出
type CalculateVolatilityOutput struct {
	Volatility float64 `json:"volatility"` // 波动性
}

// CalculateStandardDeviationInput 计算标准差输入
type CalculateStandardDeviationInput struct {
	Values []float64 `json:"values"` // 数值列表
}

// CalculateStandardDeviationOutput 计算标准差输出
type CalculateStandardDeviationOutput struct {
	StandardDeviation float64 `json:"standard_deviation"` // 标准差
}

// CalculatePercentileInput 计算百分位数输入
type CalculatePercentileInput struct {
	Values     []float64 `json:"values"`     // 数值列表
	Percentile int       `json:"percentile"` // 百分位数
}

// CalculatePercentileOutput 计算百分位数输出
type CalculatePercentileOutput struct {
	PercentileValue float64 `json:"percentile_value"` // 百分位数值
}

// CalculateStabilityScoreInput 计算稳定性评分输入
type CalculateStabilityScoreInput struct {
	PerformanceData []map[string]interface{} `json:"performance_data"` // 性能数据
}

// CalculateStabilityScoreOutput 计算稳定性评分输出
type CalculateStabilityScoreOutput struct {
	Score float64 `json:"score"` // 稳定性评分
}

// GetBottleneckSeverityInput 获取瓶颈严重程度输入
type GetBottleneckSeverityInput struct {
	Usage float64 `json:"usage"` // 使用率
}

// GetBottleneckSeverityOutput 获取瓶颈严重程度输出
type GetBottleneckSeverityOutput struct {
	Severity string `json:"severity"` // 严重程度
}

// GetEarliestTimestampInput 获取最早时间戳输入
type GetEarliestTimestampInput struct {
	Data []map[string]interface{} `json:"data"` // 数据
}

// GetEarliestTimestampOutput 获取最早时间戳输出
type GetEarliestTimestampOutput struct {
	Timestamp string `json:"timestamp"` // 时间戳
}

// GetLatestTimestampInput 获取最晚时间戳输入
type GetLatestTimestampInput struct {
	Data []map[string]interface{} `json:"data"` // 数据
}

// GetLatestTimestampOutput 获取最晚时间戳输出
type GetLatestTimestampOutput struct {
	Timestamp string `json:"timestamp"` // 时间戳
}
