package queue

import (
	"context"
)

// ===============================
// 工具方法
// ===============================

// GetQueueInstance 获取队列实例（保留原有方法）
func (s *sQueue) GetQueueInstance(ctx context.Context) (interface{}, error) {
	// 这里可以返回具体的队列客户端实例
	// 目前返回nil，实际使用时需要根据具体需求实现
	return nil, nil
}

// TODO: Removed duplicate generateQueueSelectionReason method - already exists in balance.go
