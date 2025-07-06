package queue

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"OneGfServer/api/queue/v1"
)

func (c *ControllerV1) GetQueuePerformanceReport(ctx context.Context, req *v1.GetQueuePerformanceReportReq) (res *v1.GetQueuePerformanceReportRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
