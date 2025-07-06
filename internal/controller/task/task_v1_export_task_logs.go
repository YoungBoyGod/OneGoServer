package task

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"OneGfServer/api/task/v1"
)

func (c *ControllerV1) ExportTaskLogs(ctx context.Context, req *v1.ExportTaskLogsReq) (res *v1.ExportTaskLogsRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
