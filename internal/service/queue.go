// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	queue "OneGfServer/internal/model/queue"
	"context"
)

type (
	IQueue interface {
		// CalculateLoadBalance 计算负载均衡
		CalculateLoadBalance(ctx context.Context, input *queue.CalculateLoadBalanceInput) (*queue.CalculateLoadBalanceOutput, error)
		// SelectOptimalQueue 选择最优队列
		SelectOptimalQueue(ctx context.Context, input *queue.SelectOptimalQueueInput) (*queue.SelectOptimalQueueOutput, error)
		// ValidateQueueCreation 验证队列创建
		ValidateQueueCreation(ctx context.Context, input *queue.ValidateQueueCreationInput) (*queue.ValidateQueueCreationOutput, error)
		// ValidateQueueConfiguration 验证队列配置
		ValidateQueueConfiguration(ctx context.Context, input *queue.ValidateQueueConfigurationInput) (*queue.ValidateQueueConfigurationOutput, error)
		// ValidateQueueOperation 验证队列操作
		ValidateQueueOperation(ctx context.Context, input *queue.ValidateQueueOperationInput) (*queue.ValidateQueueOperationOutput, error)
		// CalculateQueueHealthScore 计算队列健康度评分
		CalculateQueueHealthScore(ctx context.Context, input *queue.CalculateQueueHealthScoreInput) (*queue.CalculateQueueHealthScoreOutput, error)
		// SortQueueTasks 队列任务排序
		SortQueueTasks(ctx context.Context, input *queue.SortQueueTasksInput) (*queue.SortQueueTasksOutput, error)
		// CalculateQueueStatistics 计算队列统计信息
		CalculateQueueStatistics(ctx context.Context, input *queue.CalculateQueueStatisticsInput) (*queue.CalculateQueueStatisticsOutput, error)
		// ValidateTaskEnqueue 验证任务入队
		ValidateTaskEnqueue(ctx context.Context, input *queue.ValidateTaskEnqueueInput) (*queue.ValidateTaskEnqueueOutput, error)
		// GetQueueInstance 获取队列实例（保留原有方法）
		GetQueueInstance(ctx context.Context) (interface{}, error)
	}
)

var (
	localQueue IQueue
)

func Queue() IQueue {
	if localQueue == nil {
		panic("implement not found for interface IQueue, forgot register?")
	}
	return localQueue
}

func RegisterQueue(i IQueue) {
	localQueue = i
}
