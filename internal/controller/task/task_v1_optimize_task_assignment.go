package task

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"OneGfServer/api/task/v1"
)

func (c *ControllerV1) OptimizeTaskAssignment(ctx context.Context, req *v1.OptimizeTaskAssignmentReq) (res *v1.OptimizeTaskAssignmentRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
