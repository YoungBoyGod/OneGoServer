package cmd

import (
	"context"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	deviceController "OneGfServer/internal/controller/device"
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
			// s.Group("/", func(group *ghttp.RouterGroup) {
			// 	group.Middleware(ghttp.MiddlewareHandlerResponse)
			// 	group.Bind(
			// 	// user.New(),
			// 	)
			// })

			// 注册设备相关路由
			s.Group("/api/v1", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				// 设备管理路由
				group.Group("/device", func(group *ghttp.RouterGroup) {
					group.Bind(deviceController.NewV1())
				})
			})
			s.Run()
			return nil
		},
	}
)

func initDB(ctx context.Context) {
	// 等待配置加载完成
	g.Log().Info(ctx, "正在初始化数据库连接...")

	if err := g.DB().PingMaster(); err != nil {
		g.Log().Errorf(ctx, "数据库连接失败: %v", err)
		return
	}

	g.Log().Info(ctx, "数据库连接成功")

	// 测试查询
	devices, err := g.DB().GetAll(ctx, "select count(*) as count from devices")
	if err != nil {
		g.Log().Warningf(ctx, "数据库查询测试失败: %v", err)
		return
	}

	g.Log().Infof(ctx, "数据库查询测试成功，设备表记录数: %v", devices)
}
