## 下一步工作计划

1. **完成 Service 层实现**
   - 将 biz 层的业务逻辑封装到 `internal/service` 中，对外提供统一的接口。
   - 补全 `internal/service/queue.go`、`device.go`、`task.go`，实现对 biz 的调用和事务控制。

2. **完善 Repository 与持久化层**
   - 确认 `internal/data/repository` 中的 `device.go`、`task.go` 新文件实现是否完整。
   - 根据需要为队列、日志等新增 Repository。

3. **实现 Controller 层（API 层）**
   - 在 `internal/controller` 中实现 `queue.go`，并完善 `device.go`、`task.go` 的接口。
   - 定义请求/响应 DTO，调用 Service 层完成业务。

4. **配置 Router 与路由注册**
   - 在 `internal/router/router.go` 中注册新的 API 路由。
   - 确保中间件、鉴权、日志等在路由链路中生效。

5. **补充数据库迁移脚本**
   - 如果新增表或字段，更新 `internal/data/migrations`。
   - 保证迁移脚本向前兼容且可回滚。

6. **单元测试与集成测试**
   - 使用 `test/` 目录编写 Service、Controller 层测试。
   - 引入 mock 或 testcontainer 保证测试隔离。

7. **文档与示例**
   - 更新 `docs/` 下 API 接口文档、数据库设计文档、运行指南。
   - 如有需要，补充 Postman collection 或 Swagger/OpenAPI 描述。

8. **CI/CD 与部署**
   - 在 `docker-compose` 或 Kubernetes 部署脚本中加入新组件配置。
   - 确保流水线中包含 lint、测试、镜像构建、迁移步骤。

---

> 以上步骤建议按顺序逐步推进，每完成一个阶段便进行自测和代码评审，确保质量可控。 