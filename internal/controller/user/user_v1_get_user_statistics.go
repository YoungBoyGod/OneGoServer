package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"OneGfServer/api/user/v1"
)

func (c *ControllerV1) GetUserStatistics(ctx context.Context, req *v1.GetUserStatisticsReq) (res *v1.GetUserStatisticsRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
