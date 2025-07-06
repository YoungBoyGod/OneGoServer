package queue

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"OneGfServer/api/queue/v1"
)

func (c *ControllerV1) GetQueueList(ctx context.Context, req *v1.GetQueueListReq) (res *v1.GetQueueListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
