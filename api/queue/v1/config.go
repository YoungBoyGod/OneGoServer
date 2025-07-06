package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ===============================
// 队列配置管理 API (3个)
// ===============================

// GetQueueConfigReq 获取队列配置请求
type GetQueueConfigReq struct {
	g.Meta  `path:"/queue/{queueId}/config" method:"get" tags:"队列配置" summary:"获取队列配置"`
	QueueId string `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
}

type GetQueueConfigRes struct {
	Config QueueConfig `json:"config"`
}

// UpdateQueueConfigReq 更新队列配置请求
type UpdateQueueConfigReq struct {
	g.Meta               `path:"/queue/{queueId}/config" method:"put" tags:"队列配置" summary:"更新队列配置"`
	QueueId              string                 `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
	ProcessingStrategy   string                 `json:"processing_strategy,omitempty" v:"in:parallel,sequential,batch#处理策略无效"`
	RetryStrategy        string                 `json:"retry_strategy,omitempty" v:"in:immediate,exponential,linear,fixed#重试策略无效"`
	RetryDelayMs         int                    `json:"retry_delay_ms,omitempty" v:"min:100#重试延迟最小100毫秒"`
	MaxRetryInterval     int                    `json:"max_retry_interval,omitempty" v:"min:1000#最大重试间隔最小1000毫秒"`
	DeadLetterEnabled    bool                   `json:"dead_letter_enabled,omitempty"`
	DeadLetterMaxRetries int                    `json:"dead_letter_max_retries,omitempty" v:"between:0,10#死信最大重试次数为0-10"`
	BatchSize            int                    `json:"batch_size,omitempty" v:"between:1,100#批量大小为1-100"`
	TimeoutMs            int                    `json:"timeout_ms,omitempty" v:"min:1000#超时时间最小1000毫秒"`
	CustomSettings       map[string]interface{} `json:"custom_settings,omitempty"`
}

type UpdateQueueConfigRes struct {
	Message string `json:"message"`
}

// GetQueueConfigHistoryReq 获取队列配置历史请求
type GetQueueConfigHistoryReq struct {
	g.Meta  `path:"/queue/{queueId}/config/history" method:"get" tags:"队列配置" summary:"获取队列配置历史"`
	QueueId string `json:"queue_id" v:"required|max-length:50#队列ID不能为空"`
	Page    int    `json:"page" d:"1" v:"min:1#页码最小为1"`
	Size    int    `json:"size" d:"10" v:"between:1,50#每页数量为1-50"`
}

type GetQueueConfigHistoryRes struct {
	List  []QueueConfigHistory `json:"list"`
	Total int64                `json:"total"`
	Page  int                  `json:"page"`
	Size  int                  `json:"size"`
}
