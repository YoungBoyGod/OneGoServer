package task

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"OneGfServer/api/task/v1"
)

func (c *ControllerV1) GetTaskQueueStats(ctx context.Context, req *v1.GetTaskQueueStatsReq) (res *v1.GetTaskQueueStatsRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
