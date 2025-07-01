# OneGo服务器自底向上实施指南

## 🎯 实施策略总览

您的思路非常正确！**从底层往上层实现**是构建稳健应用的最佳实践。以下是完整的实施路线图：

```
📊 数据库层 (Database Layer)
    ↓
💾 数据访问层 (Data/Repository Layer)  
    ↓
🧠 业务逻辑层 (Business/Biz Layer)
    ↓
⚙️ 服务编排层 (Service Layer)
    ↓
🌐 控制器层 (Controller/API Layer)
```

---

## 📅 4阶段实施计划

### 🏗️ 阶段一：数据基础建设 (1-2天)

#### ✅ 已完成
- [x] **数据库设计分析** (`docs/database-design-analysis.md`)
- [x] **GORM模型创建**
  - `internal/biz/task/models.go` - 任务相关模型
  - `internal/biz/device/models.go` - 设备相关模型
- [x] **数据库迁移脚本**
  - `internal/data/migrations/001_create_tasks_table.sql`
  - `internal/data/migrations/002_create_devices_table.sql`
- [x] **Repository接口设计** (`internal/data/task_repository.go`)

#### 🔄 待完成
```bash
# 1. 执行数据库迁移
psql -U postgres -d onegodb -f internal/data/migrations/001_create_tasks_table.sql
psql -U postgres -d onegodb -f internal/data/migrations/002_create_devices_table.sql

# 2. 创建DeviceRepository接口
# 3. 编写Repository单元测试
```

---

### 🏛️ 阶段二：数据访问层完善 (2-3天)

#### 📋 核心任务

**A. 完善TaskRepository实现**
```go
// 需要实现的关键方法
- BatchCreate()          // 批量创建任务
- UpdateBatch()          // 批量更新
- GetPendingTasks()      // 获取待执行任务
- GetExpiredTasks()      // 获取超时任务
- ArchiveCompletedTasks() // 归档已完成任务
```

**B. 创建DeviceRepository实现**
```go
// internal/data/device_repository.go
type DeviceRepository interface {
    // 基础CRUD
    Create/GetByID/Update/Delete
    
    // 设备特定操作
    UpdateHeartbeat()     // 更新心跳
    UpdateStatus()        // 更新状态
    GetOnlineDevices()    // 获取在线设备
    GetDeviceCommands()   // 获取设备命令
    CreateCommand()       // 创建命令
    LogDeviceEvent()      // 记录设备事件
}
```

**C. 事务管理设计**
```go
// internal/data/transaction.go
type TransactionManager interface {
    WithTx(fn func(*gorm.DB) error) error
    BeginTx() (*gorm.DB, error)
    CommitTx(*gorm.DB) error
    RollbackTx(*gorm.DB) error
}
```

---

### 🧠 阶段三：业务逻辑层 (3-4天)

#### 📋 TaskBiz层增强
```go
// internal/biz/task/task_biz.go
type TaskBiz struct {
    taskRepo TaskRepository
    deviceRepo DeviceRepository
    logger   Logger
}

// 核心业务方法
func (b *TaskBiz) CreateTask(req *CreateTaskRequest) (*Task, error)
func (b *TaskBiz) ExecuteTask(taskID int64) (*TaskExecution, error)
func (b *TaskBiz) CancelTask(taskID int64) error
func (b *TaskBiz) RetryTask(taskID int64) error
func (b *TaskBiz) ScheduleTask(taskID int64, executeTime time.Time) error
```

#### 🏗️ 业务规则实现
- **任务优先级算法**: 动态计算任务执行优先级
- **设备负载均衡**: 智能分配任务到最优设备
- **故障恢复机制**: 自动重试和故障转移
- **资源限制检查**: 验证设备资源可用性
- **并发控制**: 防止任务重复执行

#### 📱 DeviceBiz层设计
```go
// internal/biz/device/device_biz.go  
type DeviceBiz struct {
    deviceRepo DeviceRepository
    logger     Logger
}

// 核心业务方法
func (b *DeviceBiz) RegisterDevice(req *RegisterDeviceRequest) (*Device, error)
func (b *DeviceBiz) UpdateDeviceStatus(deviceID string, status string) error
func (b *DeviceBiz) SendCommand(deviceID string, cmd *Command) error
func (b *DeviceBiz) ProcessHeartbeat(deviceID string, heartbeat *Heartbeat) error
func (b *DeviceBiz) GetDeviceHealth(deviceID string) (*HealthReport, error)
```

---

### ⚙️ 阶段四：服务编排层 (2-3天)

#### 🎼 TaskService实现
```go
// internal/service/task_service.go
type TaskService struct {
    taskBiz   *TaskBiz
    deviceBiz *DeviceBiz
    logger    Logger
    cache     Cache
    queue     MessageQueue
}

// 流程编排方法
func (s *TaskService) SubmitTask(req *SubmitTaskRequest) (*TaskResponse, error)
func (s *TaskService) GetTaskStatus(taskID int64) (*TaskStatusResponse, error)
func (s *TaskService) ListTasks(filter *TaskFilter) (*TaskListResponse, error)
func (s *TaskService) ExecuteTaskWorkflow(taskID int64) error
```

#### 🔄 异步处理设计
```go
// 任务执行工作流
1. 验证任务参数
2. 检查设备可用性  
3. 分配执行资源
4. 创建执行记录
5. 发送执行命令
6. 监控执行进度
7. 处理执行结果
8. 更新任务状态
```

#### 📡 DeviceService实现
```go
// internal/service/device_service.go
type DeviceService struct {
    deviceBiz *DeviceBiz
    logger    Logger
    cache     Cache
}

// 综合管理方法
func (s *DeviceService) ManageDevice(deviceID string) error
func (s *DeviceService) MonitorDeviceHealth() error
func (s *DeviceService) HandleDeviceEvents(events []DeviceEvent) error
func (s *DeviceService) OptimizeDevicePerformance() error
```

---

### 🌐 阶段五：API层集成 (1-2天)

#### 🔌 Controller层更新
```go
// 更新现有Controller，集成Service层
func (c *TaskController) CreateTask(ctx *gin.Context) {
    // 1. 参数验证
    // 2. 调用TaskService
    // 3. 返回响应
    result, err := c.taskService.SubmitTask(req)
}
```

#### 🧪 集成测试
- API端点测试
- 数据库事务测试  
- 错误处理测试
- 性能压力测试

---

## 🛠️ 技术实施要点

### 1. 依赖注入设计
```go
// internal/container/container.go
type Container struct {
    DB           *gorm.DB
    Redis        *redis.Client
    Kafka        *kafka.Producer
    TaskRepo     TaskRepository
    DeviceRepo   DeviceRepository
    TaskBiz      *TaskBiz
    DeviceBiz    *DeviceBiz
    TaskService  *TaskService
    DeviceService *DeviceService
}
```

### 2. 配置管理
```yaml
# config/config.yaml
database:
  host: localhost
  port: 5432
  dbname: onegodb
  
redis:
  host: localhost
  port: 6379
  
kafka:
  brokers: ["localhost:9092"]
```

### 3. 错误处理策略
```go
// pkg/errors/business_errors.go
var (
    ErrTaskNotFound    = errors.New("任务不存在")
    ErrDeviceOffline   = errors.New("设备离线")
    ErrInsufficientResource = errors.New("资源不足")
)
```

### 4. 日志记录规范
```go
// 结构化日志
logger.Info("任务创建成功",
    zap.Int64("task_id", task.ID),
    zap.String("task_type", task.Type),
    zap.String("user_id", userID),
)
```

---

## 📊 开发进度追踪

### ✅ 已完成 (约30%)
- [x] 数据库表设计
- [x] GORM模型定义
- [x] 数据库迁移脚本
- [x] TaskRepository接口定义
- [x] API端点占位实现

### 🔄 进行中 (下一步)
- [ ] 完善TaskRepository实现
- [ ] 创建DeviceRepository
- [ ] 编写Repository单元测试
- [ ] 实现基础的Biz层方法

### ⏳ 待启动
- [ ] Service层流程编排
- [ ] Controller层集成
- [ ] 集成测试编写
- [ ] 性能优化
- [ ] 生产环境配置

---

## 🎯 立即行动计划

### 今天就可以开始
1. **执行数据库迁移**
   ```bash
   cd /Users/haitang/code/OneGo/OneGoServer002
   psql -U postgres -d onegodb -f internal/data/migrations/001_create_tasks_table.sql
   psql -U postgres -d onegodb -f internal/data/migrations/002_create_devices_table.sql
   ```

2. **完善TaskRepository实现**
   - 修复导入路径问题
   - 实现剩余的Repository方法
   - 添加错误处理和日志记录

3. **创建DeviceRepository**
   - 基于TaskRepository模式
   - 实现设备特有的业务方法

4. **编写单元测试**
   - Repository层测试
   - 模型验证测试

这个自底向上的实施策略确保了：
- **数据一致性**: 从数据库层开始确保数据完整性
- **逐层验证**: 每层都可以独立测试和验证
- **稳固基础**: 为复杂业务逻辑提供可靠的数据支撑
- **渐进式开发**: 可以随时停止并投入使用

您现在拥有了完整的技术架构和清晰的实施路径！🚀 