package task

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"OneGfServer/api/task/v1"
)

func (c *ControllerV1) GetTaskPerformanceReport(ctx context.Context, req *v1.GetTaskPerformanceReportReq) (res *v1.GetTaskPerformanceReportRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
