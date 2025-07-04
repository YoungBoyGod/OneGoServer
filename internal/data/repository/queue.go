package repository

import (
	"context"

	"github.com/YoungBoyGod/OneGoServer/internal/biz/queue"
	"gorm.io/gorm"
)

// QueueRepository 队列数据访问接口
type QueueRepository interface {
	Create(ctx context.Context, item *queue.DeviceTaskQueue) error
	ListByDevice(ctx context.Context, deviceID int64) ([]queue.DeviceTaskQueue, error)
	Update(ctx context.Context, item *queue.DeviceTaskQueue) error
}

type queueRepoImpl struct{ db *gorm.DB }

// NewQueueRepository 构造
func NewQueueRepository(db *gorm.DB) QueueRepository {
	return &queueRepoImpl{db: db}
}

func (r *queueRepoImpl) Create(ctx context.Context, item *queue.DeviceTaskQueue) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *queueRepoImpl) ListByDevice(ctx context.Context, deviceID int64) ([]queue.DeviceTaskQueue, error) {
	var list []queue.DeviceTaskQueue
	err := r.db.WithContext(ctx).Where("device_id = ?", deviceID).Order("queue_position").Find(&list).Error
	return list, err
}

func (r *queueRepoImpl) Update(ctx context.Context, item *queue.DeviceTaskQueue) error {
	return r.db.WithContext(ctx).Save(item).Error
}
