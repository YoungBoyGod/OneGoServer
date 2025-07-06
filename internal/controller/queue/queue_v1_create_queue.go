package queue

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"OneGfServer/api/queue/v1"
)

func (c *ControllerV1) CreateQueue(ctx context.Context, req *v1.CreateQueueReq) (res *v1.CreateQueueRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
