package task

import (
	"context"
	"errors"
	"time"
)

type TaskBiz struct {
}

func (b *TaskBiz) ValidateTask(ctx context.Context, task *Task) error {
	// 复杂业务规则验证
	if task.Priority > 10 {
		return errors.New("任务优先级不能超过10")
	}

	// 业务状态检查
	if task.ExecuteTime.Before(time.Now()) {
		return errors.New("执行时间不能早于当前时间")
	}

	return nil
}

func (b *TaskBiz) CalculateTaskScore(task *Task) int {
	// 业务算法：根据优先级和紧急程度计算分数
	score := task.Priority * 10
	if task.IsUrgent {
		score += 50
	}
	return score
}
