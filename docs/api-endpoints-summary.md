# OneGo服务器API端点完整汇总

## 📊 总览统计
- **总计API端点**: 28个
- **Task任务模块**: 14个端点
- **Device设备模块**: 14个端点
- **健康检查**: 2个端点
- **API版本**: v1 (/api/v1)
- **文档风格**: RESTful + Swagger

---

## 📋 Task任务模块 (14个端点)

### 🔧 基础CRUD操作 (5个)
| 方法 | 路径 | 功能描述 | 参数说明 |
|------|------|----------|----------|
| `POST` | `/api/v1/tasks` | 创建任务 | Body: 任务信息 |
| `GET` | `/api/v1/tasks` | 获取任务列表 | Query: page, size, status |
| `GET` | `/api/v1/tasks/:id` | 获取任务详情 | Path: 任务ID |
| `PUT` | `/api/v1/tasks/:id` | 更新任务 | Path: 任务ID, Body: 更新信息 |
| `DELETE` | `/api/v1/tasks/:id` | 删除任务 | Path: 任务ID |

### 📊 任务状态管理 (2个)
| 方法 | 路径 | 功能描述 | 参数说明 |
|------|------|----------|----------|
| `GET` | `/api/v1/tasks/:id/status` | 获取任务状态 | Path: 任务ID |
| `PUT` | `/api/v1/tasks/:id/status` | 修改任务状态 | Path: 任务ID, Body: 状态信息 |

### ⚡ 任务执行控制 (3个)
| 方法 | 路径 | 功能描述 | 参数说明 |
|------|------|----------|----------|
| `POST` | `/api/v1/tasks/:id/execute` | 执行任务 | Path: 任务ID |
| `POST` | `/api/v1/tasks/:id/cancel` | 取消任务 | Path: 任务ID |
| `POST` | `/api/v1/tasks/:id/dispatch` | 分发任务 | Path: 任务ID |

### 🔍 任务信息查询 (4个)
| 方法 | 路径 | 功能描述 | 参数说明 |
|------|------|----------|----------|
| `GET` | `/api/v1/tasks/query` | 查询任务 | Query: 查询条件 |
| `GET` | `/api/v1/tasks/:id/detail` | 获取任务详情 | Path: 任务ID |
| `GET` | `/api/v1/tasks/:id/result` | 获取任务结果 | Path: 任务ID |
| `GET` | `/api/v1/tasks/:id/execution` | 获取执行记录 | Path: 任务ID |
| `GET` | `/api/v1/tasks/:id/stats` | 获取任务统计 | Path: 任务ID |

---

## 📱 Device设备模块 (14个端点)

### 🔧 基础CRUD操作 (5个)
| 方法 | 路径 | 功能描述 | 参数说明 |
|------|------|----------|----------|
| `POST` | `/api/v1/devices` | 注册设备 | Body: 设备信息 |
| `GET` | `/api/v1/devices` | 获取设备列表 | Query: page, size, status, type |
| `GET` | `/api/v1/devices/:id` | 获取设备详情 | Path: 设备ID |
| `PUT` | `/api/v1/devices/:id` | 更新设备 | Path: 设备ID, Body: 更新信息 |
| `DELETE` | `/api/v1/devices/:id` | 删除设备 | Path: 设备ID |

### 📊 设备状态管理 (3个)
| 方法 | 路径 | 功能描述 | 参数说明 |
|------|------|----------|----------|
| `GET` | `/api/v1/devices/:id/status` | 获取设备状态 | Path: 设备ID |
| `POST` | `/api/v1/devices/online` | 设备上线 | Body: 设备信息 |
| `POST` | `/api/v1/devices/offline` | 设备下线 | Body: 设备信息 |

### 🎮 设备控制操作 (3个)
| 方法 | 路径 | 功能描述 | 参数说明 |
|------|------|----------|----------|
| `POST` | `/api/v1/devices/:id/command` | 发送设备命令 | Path: 设备ID, Body: 命令信息 |
| `GET` | `/api/v1/devices/:id/heartbeat` | 获取设备心跳 | Path: 设备ID |
| `POST` | `/api/v1/devices/:id/heartbeat` | 更新设备心跳 | Path: 设备ID |

### 🔍 设备信息查询 (3个)
| 方法 | 路径 | 功能描述 | 参数说明 |
|------|------|----------|----------|
| `GET` | `/api/v1/devices/query` | 查询设备 | Query: 查询条件 |
| `GET` | `/api/v1/devices/:id/logs` | 获取设备日志 | Path: 设备ID |
| `GET` | `/api/v1/devices/:id/stats` | 获取设备统计 | Path: 设备ID |

---

## 🏥 系统健康检查 (2个端点)

### 🔍 健康检查接口
| 方法 | 路径 | 功能描述 | 响应内容 |
|------|------|----------|----------|
| `GET` | `/health` | 完整健康检查 | 系统状态、数据库、Redis、Kafka连接状态 |
| `GET` | `/ping` | 简单存活检查 | pong响应，服务可用性确认 |

---

## 📋 API设计规范

### 🎯 RESTful设计原则
- **资源命名**: 使用名词复数形式 (`tasks`, `devices`)
- **HTTP方法**: 标准CRUD操作映射
  - `GET`: 查询操作
  - `POST`: 创建操作  
  - `PUT`: 更新操作
  - `DELETE`: 删除操作
- **状态码**: 标准HTTP状态码
  - `200`: 成功
  - `201`: 创建成功
  - `400`: 请求参数错误
  - `404`: 资源不存在
  - `500`: 服务器内部错误

### 📝 统一响应格式
```json
{
  "message": "操作描述信息",
  "data": {
    // 具体数据内容
  }
}
```

### 📖 Swagger文档支持
- **完整注解**: 每个API都包含完整的Swagger文档注解
- **参数说明**: 详细的请求参数和响应格式说明
- **标签分组**: 按功能模块进行API分组
  - `tasks`: 任务管理相关API
  - `devices`: 设备管理相关API

### 🔍 查询参数支持
- **分页参数**: `page`, `size`
- **筛选参数**: `status`, `type` 等
- **默认值**: 合理的默认值设置

---

## 🚀 开发状态

### ✅ 已完成
- [x] **API架构设计**: 28个端点完整规划
- [x] **控制器实现**: Task和Device控制器
- [x] **路由配置**: 完整的路由映射
- [x] **中间件系统**: CORS、日志、恢复中间件
- [x] **文档注解**: Swagger API文档
- [x] **统一响应**: 标准JSON响应格式

### 🔄 开发中
- [ ] **业务逻辑**: 控制器方法具体实现
- [ ] **数据层**: PostgreSQL数据库集成
- [ ] **缓存层**: Redis缓存优化
- [ ] **消息队列**: Kafka异步处理

### 📅 待开发
- [ ] **认证系统**: JWT用户认证
- [ ] **参数验证**: 请求数据校验
- [ ] **API限流**: 访问频率控制
- [ ] **监控告警**: 性能指标收集
- [ ] **单元测试**: API接口测试
- [ ] **集成测试**: 端到端测试

---

## 🔗 相关文档
- [路由配置实现报告](./router-configuration-implementation.md)
- [Kafka功能实现总结](./kafka-producer-consumer-implementation-summary.md)
- [数据库集成报告](./postgresql-integration-test-report.md)

---

## 📞 使用示例

### 创建任务
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"name": "数据备份任务", "type": "backup"}'
```

### 获取设备列表
```bash
curl -X GET "http://localhost:8080/api/v1/devices?page=1&size=10&status=online"
```

### 健康检查
```bash
curl -X GET http://localhost:8080/health
```

通过以上API端点，OneGo服务器提供了完整的任务管理和设备管理功能，为IoT和任务调度场景提供了强大的后端支持。 