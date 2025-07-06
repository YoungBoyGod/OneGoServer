package common

import (
	"time"

	"github.com/gogf/gf/v2/os/gtime"
)

// ===============================
// 通用常量定义
// ===============================

// 通用状态常量
const (
	StatusActive    = "active"
	StatusInactive  = "inactive"
	StatusPending   = "pending"
	StatusRunning   = "running"
	StatusStopped   = "stopped"
	StatusPaused    = "paused"
	StatusError     = "error"
	StatusCompleted = "completed"
	StatusFailed    = "failed"
	StatusCanceled  = "canceled"
	StatusLocked    = "locked"
)

// 通用排序常量
const (
	SortOrderAsc  = "asc"
	SortOrderDesc = "desc"
)

// 通用操作类型常量
const (
	ActionTypeCreate   = "create"
	ActionTypeUpdate   = "update"
	ActionTypeDelete   = "delete"
	ActionTypeView     = "view"
	ActionTypeStart    = "start"
	ActionTypeStop     = "stop"
	ActionTypePause    = "pause"
	ActionTypeResume   = "resume"
	ActionTypeCancel   = "cancel"
	ActionTypeRestart  = "restart"
	ActionTypeRetry    = "retry"
	ActionTypeAssign   = "assign"
	ActionTypeUnassign = "unassign"
)

// 通用日志级别常量
const (
	LogLevelDebug    = "debug"
	LogLevelInfo     = "info"
	LogLevelWarn     = "warn"
	LogLevelError    = "error"
	LogLevelCritical = "critical"
)

// 通用优先级常量
const (
	PriorityLowest  = 1
	PriorityLow     = 3
	PriorityNormal  = 5
	PriorityHigh    = 7
	PriorityUrgent  = 9
	PriorityHighest = 10
)

// ===============================
// 通用分页结构 (参考API层)
// ===============================

// PaginationResponse 通用分页响应结构
type PaginationResponse[T any] struct {
	List  []T   `json:"list"`  // 数据列表
	Total int64 `json:"total"` // 总记录数
	Page  int   `json:"page"`  // 当前页码
	Size  int   `json:"size"`  // 每页大小
}

// PaginationRequest 通用分页请求结构
type PaginationRequest struct {
	Page      int    `json:"page" d:"1" v:"min:1#页码最小为1"`                // 页码
	Size      int    `json:"size" d:"10" v:"between:1,100#每页数量为1-100"`   // 每页大小
	SortBy    string `json:"sort_by,omitempty"`                          // 排序字段
	SortOrder string `json:"sort_order" d:"desc" v:"in:asc,desc#排序方向无效"` // 排序方向
}

// PaginationInfo 分页信息
type PaginationInfo struct {
	CurrentPage  int   `json:"current_page"`  // 当前页码
	PageSize     int   `json:"page_size"`     // 每页大小
	TotalPages   int   `json:"total_pages"`   // 总页数
	TotalRecords int64 `json:"total_records"` // 总记录数
	HasNext      bool  `json:"has_next"`      // 是否有下一页
	HasPrev      bool  `json:"has_prev"`      // 是否有上一页
}

// NewPaginationResponse 创建分页响应
func NewPaginationResponse[T any](list []T, total int64, page, size int) PaginationResponse[T] {
	return PaginationResponse[T]{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	}
}

// GetPaginationInfo 获取分页信息
func GetPaginationInfo(total int64, page, size int) PaginationInfo {
	totalPages := int((total + int64(size) - 1) / int64(size))

	return PaginationInfo{
		CurrentPage:  page,
		PageSize:     size,
		TotalPages:   totalPages,
		TotalRecords: total,
		HasNext:      page < totalPages,
		HasPrev:      page > 1,
	}
}

// ===============================
// 通用基础结构体
// ===============================

// BaseEntity 基础实体结构体
type BaseEntity struct {
	ID        int64       `json:"id"`
	CreatedAt *gtime.Time `json:"created_at"`
	UpdatedAt *gtime.Time `json:"updated_at"`
	CreatedBy int64       `json:"created_by"`
	UpdatedBy int64       `json:"updated_by"`
}

// BaseEntityWithStatus 带状态的基础实体结构体
type BaseEntityWithStatus struct {
	BaseEntity
	Status string `json:"status"`
}

// BaseEntityWithName 带名称的基础实体结构体
type BaseEntityWithName struct {
	BaseEntity
	Name        string `json:"name"`
	Description string `json:"description"`
}

// BaseEntityWithPriority 带优先级的基础实体结构体
type BaseEntityWithPriority struct {
	BaseEntity
	Priority int `json:"priority"`
}

// ===============================
// 通用Input/Output结构体 (保持向后兼容)
// ===============================

// BaseCreateInput 基础创建输入
type BaseCreateInput struct {
	Entity interface{} `json:"entity"`
}

// BaseCreateOutput 基础创建输出
type BaseCreateOutput struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

// BaseGetByIDInput 基础根据ID获取输入
type BaseGetByIDInput struct {
	ID int64 `json:"id"`
}

// BaseGetByIDOutput 基础根据ID获取输出
type BaseGetByIDOutput struct {
	Entity interface{} `json:"entity"`
}

// BaseGetByCodeInput 基础根据编码获取输入
type BaseGetByCodeInput struct {
	Code string `json:"code"`
}

// BaseGetByCodeOutput 基础根据编码获取输出
type BaseGetByCodeOutput struct {
	Entity interface{} `json:"entity"`
}

// BaseUpdateInput 基础更新输入
type BaseUpdateInput struct {
	Entity interface{} `json:"entity"`
}

// BaseUpdateOutput 基础更新输出
type BaseUpdateOutput struct {
	Message string `json:"message"`
}

// BaseDeleteInput 基础删除输入
type BaseDeleteInput struct {
	ID    string `json:"id"`
	Force bool   `json:"force"`
}

// BaseDeleteOutput 基础删除输出
type BaseDeleteOutput struct {
	Message string `json:"message"`
}

// BaseListInput 基础列表输入
type BaseListInput struct {
	Filter     interface{}       `json:"filter"`
	Sort       interface{}       `json:"sort"`
	Pagination *PaginationOption `json:"pagination"`
}

// BaseListOutput 基础列表输出
type BaseListOutput struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

// BaseBatchInput 基础批量操作输入
type BaseBatchInput struct {
	IDs    []string    `json:"ids"`
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
}

// BaseBatchOutput 基础批量操作输出
type BaseBatchOutput struct {
	Total     int      `json:"total"`
	Success   int      `json:"success"`
	Failed    int      `json:"failed"`
	FailedIDs []string `json:"failed_ids"`
	Message   string   `json:"message"`
}

// BaseActionInput 基础操作输入
type BaseActionInput struct {
	ID     string `json:"id"`
	Action string `json:"action"`
	Reason string `json:"reason"`
	Force  bool   `json:"force"`
}

// BaseActionOutput 基础操作输出
type BaseActionOutput struct {
	Message string `json:"message"`
}

// ===============================
// 通用过滤和排序选项
// ===============================

// BaseFilter 基础过滤条件
type BaseFilter struct {
	Status    []string   `json:"status"`
	Type      []string   `json:"type"`
	Keyword   *string    `json:"keyword"`
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
	Enabled   *bool      `json:"enabled"`
}

// BaseSortOption 基础排序选项
type BaseSortOption struct {
	Field string `json:"field"` // id, name, status, priority, created_at, updated_at
	Order string `json:"order"` // asc, desc
}

// PaginationOption 分页选项 (保持向后兼容)
type PaginationOption struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

// TimeRangeFilter 时间范围过滤
type TimeRangeFilter struct {
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
}

// KeywordFilter 关键词过滤
type KeywordFilter struct {
	Keyword *string `json:"keyword"`
}

// StatusFilter 状态过滤
type StatusFilter struct {
	Status []string `json:"status"`
}

// TypeFilter 类型过滤
type TypeFilter struct {
	Type []string `json:"type"`
}

// PriorityFilter 优先级过滤
type PriorityFilter struct {
	Priority *int `json:"priority"`
}

// ===============================
// 通用统计和报告结构体
// ===============================

// BaseStatistics 基础统计信息
type BaseStatistics struct {
	Total        int64            `json:"total"`
	ByStatus     map[string]int64 `json:"by_status"`
	ByType       map[string]int64 `json:"by_type"`
	RecentActive int64            `json:"recent_active"`
	TrendData    []TrendPoint     `json:"trend_data"`
}

// TrendPoint 趋势数据点
type TrendPoint struct {
	Timestamp *gtime.Time `json:"timestamp"`
	Count     int         `json:"count"`
	Value     float64     `json:"value"`
}

// BasePerformanceReport 基础性能报告
type BasePerformanceReport struct {
	ID               string             `json:"id"`
	Name             string             `json:"name"`
	ReportPeriod     string             `json:"report_period"`
	OverallScore     float64            `json:"overall_score"`
	PerformanceTrend []TrendPoint       `json:"performance_trend"`
	Issues           []PerformanceIssue `json:"issues"`
	Recommendations  []string           `json:"recommendations"`
}

// PerformanceIssue 性能问题
type PerformanceIssue struct {
	IssueType      string  `json:"issue_type"`
	Severity       string  `json:"severity"`
	Description    string  `json:"description"`
	Impact         float64 `json:"impact"`
	Recommendation string  `json:"recommendation"`
}

// ===============================
// 通用日志和活动结构体
// ===============================

// BaseLog 基础日志结构体
type BaseLog struct {
	LogID     string      `json:"log_id"`
	EntityID  string      `json:"entity_id"`
	Level     string      `json:"level"`
	Message   string      `json:"message"`
	Timestamp *gtime.Time `json:"timestamp"`
	Source    string      `json:"source"`
	IPAddress string      `json:"ip_address"`
	UserAgent string      `json:"user_agent"`
}

// BaseActivity 基础活动结构体
type BaseActivity struct {
	ActivityID  string                 `json:"activity_id"`
	EntityID    string                 `json:"entity_id"`
	ActionType  string                 `json:"action_type"`
	Resource    string                 `json:"resource"`
	ResourceID  string                 `json:"resource_id"`
	Description string                 `json:"description"`
	IPAddress   string                 `json:"ip_address"`
	UserAgent   string                 `json:"user_agent"`
	Timestamp   *gtime.Time            `json:"timestamp"`
	Result      string                 `json:"result"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// ===============================
// 通用会话和认证结构体
// ===============================

// BaseSession 基础会话结构体
type BaseSession struct {
	SessionID    string      `json:"session_id"`
	EntityID     string      `json:"entity_id"`
	DeviceInfo   string      `json:"device_info"`
	IPAddress    string      `json:"ip_address"`
	UserAgent    string      `json:"user_agent"`
	LoginTime    *gtime.Time `json:"login_time"`
	LastActivity *gtime.Time `json:"last_activity"`
	ExpiresAt    *gtime.Time `json:"expires_at"`
	Status       string      `json:"status"`
	Location     string      `json:"location"`
}

// BaseAuth 基础认证结构体
type BaseAuth struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int64    `json:"expires_in"`
	TokenType    string   `json:"token_type"`
	Permissions  []string `json:"permissions"`
	Roles        []string `json:"roles"`
}

// ===============================
// 通用配置和设置结构体
// ===============================

// BaseConfig 基础配置结构体
type BaseConfig struct {
	ConfigID    string                 `json:"config_id"`
	EntityID    string                 `json:"entity_id"`
	ConfigType  string                 `json:"config_type"`
	ConfigData  map[string]interface{} `json:"config_data"`
	IsDefault   bool                   `json:"is_default"`
	Description string                 `json:"description"`
	CreatedAt   *gtime.Time            `json:"created_at"`
	UpdatedAt   *gtime.Time            `json:"updated_at"`
}

// BaseSettings 基础设置结构体
type BaseSettings struct {
	SettingID   string                 `json:"setting_id"`
	EntityID    string                 `json:"entity_id"`
	SettingType string                 `json:"setting_type"`
	SettingData map[string]interface{} `json:"setting_data"`
	IsEnabled   bool                   `json:"is_enabled"`
	CreatedAt   *gtime.Time            `json:"created_at"`
	UpdatedAt   *gtime.Time            `json:"updated_at"`
}

// ===============================
// 通用验证结构体
// ===============================

// BaseValidationInput 基础验证输入
type BaseValidationInput struct {
	Data interface{} `json:"data"`
}

// BaseValidationOutput 基础验证输出
type BaseValidationOutput struct {
	Valid bool   `json:"valid"`
	Error string `json:"error"`
}

// ===============================
// 通用接口定义
// ===============================

// Entity 实体接口
type Entity interface {
	GetID() int64
	GetCreatedAt() *gtime.Time
	GetUpdatedAt() *gtime.Time
}

// Filter 过滤接口
type Filter interface {
	GetStatus() []string
	GetType() []string
	GetKeyword() *string
	GetStartTime() *time.Time
	GetEndTime() *time.Time
}

// Sort 排序接口
type Sort interface {
	GetField() string
	GetOrder() string
}

// Statistics 统计接口
type Statistics interface {
	GetTotal() int64
	GetByStatus() map[string]int64
	GetByType() map[string]int64
}

// ===============================
// 通用工具函数
// ===============================

// IsValidStatus 检查状态是否有效
func IsValidStatus(status string) bool {
	validStatuses := []string{
		StatusActive, StatusInactive, StatusPending, StatusRunning,
		StatusStopped, StatusPaused, StatusError, StatusCompleted,
		StatusFailed, StatusCanceled, StatusLocked,
	}
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

// IsValidSortOrder 检查排序方向是否有效
func IsValidSortOrder(order string) bool {
	return order == SortOrderAsc || order == SortOrderDesc
}

// IsValidPriority 检查优先级是否有效
func IsValidPriority(priority int) bool {
	return priority >= PriorityLowest && priority <= PriorityHighest
}

// IsValidLogLevel 检查日志级别是否有效
func IsValidLogLevel(level string) bool {
	validLevels := []string{
		LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError, LogLevelCritical,
	}
	for _, validLevel := range validLevels {
		if level == validLevel {
			return true
		}
	}
	return false
}

// GetDefaultPagination 获取默认分页选项
func GetDefaultPagination() *PaginationOption {
	return &PaginationOption{
		Page: 1,
		Size: 20,
	}
}

// GetDefaultSort 获取默认排序选项
func GetDefaultSort(field string) *BaseSortOption {
	return &BaseSortOption{
		Field: field,
		Order: SortOrderDesc,
	}
}

// ===============================
// 通用错误定义
// ===============================

// CommonError 通用错误
type CommonError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
}

// Error 实现error接口
func (e *CommonError) Error() string {
	return e.Message
}

// NewCommonError 创建通用错误
func NewCommonError(code, message, details string) *CommonError {
	return &CommonError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// ===============================
// 通用响应结构体
// ===============================

// BaseResponse 基础响应结构体
type BaseResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	TraceID string      `json:"trace_id"`
}

// NewSuccessResponse 创建成功响应
func NewSuccessResponse(data interface{}) *BaseResponse {
	return &BaseResponse{
		Code:    200,
		Message: "success",
		Data:    data,
	}
}

// NewErrorResponse 创建错误响应
func NewErrorResponse(code int, message string) *BaseResponse {
	return &BaseResponse{
		Code:    code,
		Message: message,
		Data:    nil,
	}
}
