# 设备注册调试指南

> 更新时间：{{date}}

## 前置条件
1. PostgreSQL 已运行，数据库 `onego` 初始化并执行全部迁移脚本 `internal/data/migrations/*.sql`
2. Redis 已运行（默认 `localhost:6379`，无密码或按 `config.yaml` 配置）。
3. Kafka 已运行并确保 `broker`、`zookeeper` 正常；Topic `device-task-dispatch` 已创建（或开启自动创建）。
4. `.env` / `config/config.yaml` 指向正确的连接信息。
5. 执行 `go mod tidy` 保证依赖完整。

## 启动服务
```bash
# 1. 启动依赖（示例 docker-compose）
docker compose -f docker-compose.kafka-ui-final.yml up -d postgres redis kafka zookeeper

# 2. 运行应用
go run ./cmd/root.go   # 或 go run ./main.go
```
你应看到日志：`Server started on port xxxx`。

## 调试步骤
1. **注册设备**
   ```bash
   curl -X POST http://localhost:8080/api/v1/devices \
        -H "Content-Type: application/json" \
        -d '{"device_id":"dev001","name":"测试网关","type":"gateway"}'
   ```
   预期返回 `201 Created`，JSON `data` 包含 `id`。
2. **获取设备状态**
   ```bash
   curl http://localhost:8080/api/v1/devices/{id}/status
   ```
3. **心跳上报**
   ```bash
   curl -X POST http://localhost:8080/api/v1/devices/{id}/heartbeat \
        -H "Content-Type: application/json" \
        -d '{"status":"online","response_time":123}'
   ```
4. **发送设备命令**
   ```bash
   curl -X POST http://localhost:8080/api/v1/devices/{id}/command \
        -H "Content-Type: application/json" \
        -d '{"command_type":"reboot","command_data":{"delay":5}}'
   ```
5. **任务创建 & 分发（示例）**
   ```bash
   curl -X POST http://localhost:8080/api/v1/tasks \
        -H "Content-Type: application/json" \
        -d '{"name":"备份","type":"backup"}'
   # 分发
   curl -X POST http://localhost:8080/api/v1/tasks/{taskId}/dispatch
   ```
6. **查看队列**
   ```bash
   curl http://localhost:8080/api/v1/device-queues/{deviceID}/metrics
   ```

## 调试清单
- [ ] 依赖容器全部 healthy
- [ ] `go run` 启动无 panic
- [ ] 设备注册成功，DB 表 `devices` 记录正确
- [ ] 心跳写入 `device_heartbeats`
- [ ] 命令写入 Kafka topic，Redis list 观察
- [ ] 任务创建 & 状态流转
- [ ] 中间件：无 token 请求返回 401，带 `Bearer demo` 通过
- [ ] 速率限制：连续压测出现 429

完成以上检查即表明核心链路可用，可继续扩展 Worker 与监控。 