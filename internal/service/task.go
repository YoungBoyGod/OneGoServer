package service

import (
	"context"
	"fmt"

	"github.com/YoungBoyGod/OneGoServer/pkg/cache"
	"github.com/YoungBoyGod/OneGoServer/pkg/queue"
)

type TaskService struct {
	taskBiz  biz.TaskBiz
	taskData data.TaskRepository
	kafka    queue.KafkaManager
	redis    cache.RedisManager
}

func (s *TaskService) CreateTask(ctx context.Context, req *CreateTaskRequest) (*TaskResponse, error) {
	// 1. 参数转换和预处理
	task := convertRequestToTask(req)

	// 2. 调用biz层业务验证
	if err := s.taskBiz.ValidateTask(ctx, task); err != nil {
		return nil, err
	}

	// 3. 数据持久化
	savedTask, err := s.taskData.Create(ctx, task)
	if err != nil {
		return nil, err
	}

	// 4. 发送Kafka消息通知
	go s.kafka.SendMessage(ctx, &KafkaMessage{
		Topic: "task-created",
		Value: savedTask,
	})

	// 5. 缓存更新
	s.redis.Set(ctx, fmt.Sprintf("task:%s", savedTask.ID), savedTask)

	return convertTaskToResponse(savedTask), nil
}
