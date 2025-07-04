package service

import (
	"context"
	"fmt"
	"time"

	"github.com/YoungBoyGod/OneGoServer/internal/biz/queue"
	"github.com/YoungBoyGod/OneGoServer/internal/data/repository"
	pkglog "github.com/YoungBoyGod/OneGoServer/pkg/log"
	"github.com/YoungBoyGod/OneGoServer/pkg/sql"
	"go.uber.org/zap"
)

// QueueService 队列服务接口
// TODO: 后续根据业务场景补充更完整的方法
type QueueService interface {
	EnqueueTask(ctx context.Context, deviceID int64, taskID int64, priority int) error
	DequeueTask(ctx context.Context, deviceID int64) (*queue.DeviceTaskQueue, error)
	ReorderQueue(ctx context.Context, deviceID int64, strategy string) error
	GetQueueMetrics(ctx context.Context, deviceID int64) (map[string]interface{}, error)
}

type queueServiceImpl struct {
	queueBiz   *queue.QueueBusiness
	taskRepo   repository.TaskRepository
	deviceRepo repository.DeviceRepository
	queueRepo  repository.QueueRepository
	logger     *zap.Logger
}

// NewQueueService 创建队列服务实例
func NewQueueService(taskRepo repository.TaskRepository, deviceRepo repository.DeviceRepository) QueueService {
	return &queueServiceImpl{
		queueBiz:   queue.NewQueueBusiness(),
		taskRepo:   taskRepo,
		deviceRepo: deviceRepo,
		queueRepo:  repository.NewQueueRepository(sql.GetDB()),
		logger:     pkglog.GetAppLogger(nil),
	}
}

// EnqueueTask 将任务加入设备队列
func (s *queueServiceImpl) EnqueueTask(ctx context.Context, deviceID int64, taskID int64, priority int) error {
	// 校验设备与任务存在
	if _, err := s.deviceRepo.GetByID(ctx, deviceID); err != nil {
		return err
	}
	if _, err := s.taskRepo.GetByID(ctx, taskID); err != nil {
		return err
	}

	item := queue.DeviceTaskQueue{
		DeviceID:      deviceID,
		TaskID:        taskID,
		QueuePriority: priority,
		CreatedAt:     time.Now(),
	}

	if err := s.queueRepo.Create(ctx, &item); err != nil {
		return err
	}
	s.logger.Info("任务加入队列", zap.Int64("device_id", deviceID), zap.Int64("task_id", taskID), zap.Int("priority", priority))
	return nil
}

// DequeueTask 从队列中取出下一个任务
func (s *queueServiceImpl) DequeueTask(ctx context.Context, deviceID int64) (*queue.DeviceTaskQueue, error) {
	return nil, fmt.Errorf("not implemented")
}

// ReorderQueue 重新排序队列
func (s *queueServiceImpl) ReorderQueue(ctx context.Context, deviceID int64, strategy string) error {
	return fmt.Errorf("not implemented")
}

// GetQueueMetrics 获取队列指标
func (s *queueServiceImpl) GetQueueMetrics(ctx context.Context, deviceID int64) (map[string]interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}
