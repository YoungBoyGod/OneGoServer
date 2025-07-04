package cmd

import (
	"context"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			// 初始化数据库
			initDB(ctx)
			s := g.Server()
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(
					device.NewV1(),
				)
			})
			s.Run()
			return nil
		},
	}
)

func initDB(ctx context.Context) {
	if err := g.DB().PingMaster(); err != nil {
		g.Log().Fatalf(ctx, "database connection failed: %v", err)
	}
	devices, err := g.DB().GetAll(ctx, "select * from devices")
	if err != nil {
		g.Log().Fatalf(ctx, "database query failed: %v", err)
	}
	g.Log().Info(ctx, devices)
}
