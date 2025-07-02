# OneGo系统下一步实施指南

## 概述

现在您已经完成了数据库设计和任务生命周期流程规划，下一步需要进行**Go代码实现**。本文档提供详细的实施步骤和优先级建议。

## 🎯 实施优先级

### 第一阶段：核心基础设施 (优先级：🔥 最高)

#### 1.1 数据库连接和迁移
```bash
# 1. 确保数据库连接正常
# 2. 执行所有migration文件
# 3. 验证表结构创建成功
```

**需要实现**：
- 数据库连接池配置
- Migration执行机制
- 连接健康检查

#### 1.2 基础模型定义
```go
// 需要创建的核心模型
type Task struct { ... }
type TaskAssignmentQueue struct { ... }
type DeviceTaskQueue struct { ... }
type DeviceLoadMonitor struct { ... }
type TaskAssignmentHistory struct { ... }
```

**文件位置**：
- `internal/biz/task/models.go`
- `internal/biz/device/models.go`

### 第二阶段：数据访问层 (优先级：🔥 高)

#### 2.1 Repository层实现
```go
// 核心Repository接口
type TaskRepository interface {
    Create(ctx context.Context, task *Task) error
    GetByID(ctx context.Context, taskID string) (*Task, error)
    Update(ctx context.Context, task *Task) error
    Delete(ctx context.Context, taskID string) error
}

type TaskAssignmentRepository interface {
    Enqueue(ctx context.Context, assignment *TaskAssignmentQueue) error
    Dequeue(ctx context.Context, limit int) ([]*TaskAssignmentQueue, error)
    AssignDevice(ctx context.Context, taskID string, deviceID int64) error
    UpdateStatus(ctx context.Context, taskID string, status string) error
}

type DeviceTaskRepository interface {
    AddToQueue(ctx context.Context, deviceTask *DeviceTaskQueue) error
    GetNextTask(ctx context.Context, deviceID int64) (*DeviceTaskQueue, error)
    UpdateStatus(ctx context.Context, id int64, status string) error
    Requeue(ctx context.Context, id int4) error
}
```

**文件位置**：
- `internal/data/task_repository.go`
- `internal/data/task_assignment_repository.go`
- `internal/data/device_task_repository.go`

#### 2.2 数据库操作实现
```go
// PostgreSQL实现
type postgresTaskRepository struct {
    db *sql.DB
}

// 实现所有Repository接口方法
func (r *postgresTaskRepository) Create(ctx context.Context, task *Task) error {
    // SQL实现
}

func (r *postgresTaskRepository) Enqueue(ctx context.Context, assignment *TaskAssignmentQueue) error {
    // 插入分配队列
}
```

### 第三阶段：业务逻辑层 (优先级：🔥 高)

#### 3.1 Service层实现
```go
// 任务管理服务
type TaskService interface {
    CreateTask(ctx context.Context, req *CreateTaskRequest) (*Task, error)
    GetTask(ctx context.Context, taskID string) (*Task, error)
    UpdateTask(ctx context.Context, taskID string, req *UpdateTaskRequest) error
    DeleteTask(ctx context.Context, taskID string) error
}

// 任务分配服务
type TaskAssignmentService interface {
    EnqueueTask(ctx context.Context, taskID string, req *EnqueueRequest) error
    AssignTasks(ctx context.Context) error
    GetQueueStatus(ctx context.Context) (*QueueStatus, error)
    UpdatePriority(ctx context.Context, taskID string, priority int) error
}

// 设备任务服务
type DeviceTaskService interface {
    GetNextTask(ctx context.Context, deviceID int64) (*DeviceTask, error)
    StartTask(ctx context.Context, taskID string) error
    CompleteTask(ctx context.Context, taskID string, result *TaskResult) error
    FailTask(ctx context.Context, taskID string, reason string) error
    RequeueTask(ctx context.Context, taskID string) error
}
```

**文件位置**：
- `internal/service/task.go`
- `internal/service/task_assignment.go`
- `internal/service/device_task.go`

#### 3.2 核心业务逻辑
```go
// 任务分配算法
func (s *taskAssignmentService) assignTasks(ctx context.Context) error {
    // 1. 获取待分配任务
    // 2. 获取可用设备
    // 3. 执行分配算法
    // 4. 更新分配结果
    // 5. 记录分配历史
}

// 设备选择算法
func (s *taskAssignmentService) selectBestDevice(ctx context.Context, task *TaskAssignmentQueue) (*Device, error) {
    // 1. 过滤设备类型
    // 2. 应用白/黑名单
    // 3. 检查负载情况
    // 4. 计算分配评分
    // 5. 选择最优设备
}
```

### 第四阶段：API接口层 (优先级：🔥 中)

#### 4.1 Controller层实现
```go
// 任务管理API
type TaskController struct {
    taskService TaskService
}

func (c *TaskController) CreateTask(ctx *gin.Context) {
    // 处理创建任务请求
}

func (c *TaskController) GetTask(ctx *gin.Context) {
    // 处理获取任务请求
}

// 任务分配API
type TaskAssignmentController struct {
    assignmentService TaskAssignmentService
}

func (c *TaskAssignmentController) EnqueueTask(ctx *gin.Context) {
    // 处理任务入队请求
}

func (c *TaskAssignmentController) GetQueueStatus(ctx *gin.Context) {
    // 处理获取队列状态请求
}

// 设备任务API
type DeviceTaskController struct {
    deviceTaskService DeviceTaskService
}

func (c *DeviceTaskController) GetNextTask(ctx *gin.Context) {
    // 处理获取下一个任务请求
}

func (c *DeviceTaskController) CompleteTask(ctx *gin.Context) {
    // 处理任务完成请求
}
```

**文件位置**：
- `internal/controller/task.go`
- `internal/controller/task_assignment.go`
- `internal/controller/device_task.go`

#### 4.2 路由配置
```go
// 在 router.go 中添加路由
func (r *Router) setupTaskRoutes() {
    taskGroup := r.engine.Group("/api/v1/tasks")
    {
        taskGroup.POST("/", r.taskController.CreateTask)
        taskGroup.GET("/:task_id", r.taskController.GetTask)
        taskGroup.PUT("/:task_id", r.taskController.UpdateTask)
        taskGroup.DELETE("/:task_id", r.taskController.DeleteTask)
    }
    
    assignmentGroup := r.engine.Group("/api/v1/assignments")
    {
        assignmentGroup.POST("/enqueue", r.assignmentController.EnqueueTask)
        assignmentGroup.GET("/status", r.assignmentController.GetQueueStatus)
        assignmentGroup.PUT("/:task_id/priority", r.assignmentController.UpdatePriority)
    }
    
    deviceGroup := r.engine.Group("/api/v1/devices")
    {
        deviceGroup.GET("/:device_id/next-task", r.deviceTaskController.GetNextTask)
        deviceGroup.POST("/:device_id/tasks/:task_id/start", r.deviceTaskController.StartTask)
        deviceGroup.POST("/:device_id/tasks/:task_id/complete", r.deviceTaskController.CompleteTask)
        deviceGroup.POST("/:device_id/tasks/:task_id/fail", r.deviceTaskController.FailTask)
    }
}
```

### 第五阶段：调度器和执行器 (优先级：🔥 中)

#### 5.1 任务调度器
```go
// 任务调度器
type TaskScheduler struct {
    assignmentService TaskAssignmentService
    ticker           *time.Ticker
    stopChan         chan struct{}
}

func (s *TaskScheduler) Start() {
    // 启动调度器
    s.ticker = time.NewTicker(30 * time.Second) // 每30秒执行一次
    go s.run()
}

func (s *TaskScheduler) run() {
    for {
        select {
        case <-s.ticker.C:
            s.assignTasks()
        case <-s.stopChan:
            return
        }
    }
}

func (s *TaskScheduler) assignTasks() {
    // 执行任务分配逻辑
    ctx := context.Background()
    if err := s.assignmentService.AssignTasks(ctx); err != nil {
        log.Error("Failed to assign tasks", "error", err)
    }
}
```

#### 5.2 设备执行器
```go
// 设备执行器
type DeviceExecutor struct {
    deviceID         int64
    deviceTaskService DeviceTaskService
    taskProcessor    TaskProcessor
    stopChan         chan struct{}
}

func (e *DeviceExecutor) Start() {
    go e.run()
}

func (e *DeviceExecutor) run() {
    for {
        select {
        case <-e.stopChan:
            return
        default:
            e.processNextTask()
            time.Sleep(5 * time.Second) // 等待5秒
        }
    }
}

func (e *DeviceExecutor) processNextTask() {
    ctx := context.Background()
    
    // 获取下一个任务
    task, err := e.deviceTaskService.GetNextTask(ctx, e.deviceID)
    if err != nil {
        if err != ErrNoTaskAvailable {
            log.Error("Failed to get next task", "device_id", e.deviceID, "error", err)
        }
        return
    }
    
    // 开始执行任务
    if err := e.deviceTaskService.StartTask(ctx, task.TaskID); err != nil {
        log.Error("Failed to start task", "task_id", task.TaskID, "error", err)
        return
    }
    
    // 执行任务
    result, err := e.taskProcessor.Execute(ctx, task)
    if err != nil {
        // 任务失败
        e.deviceTaskService.FailTask(ctx, task.TaskID, err.Error())
    } else {
        // 任务成功
        e.deviceTaskService.CompleteTask(ctx, task.TaskID, result)
    }
}
```

### 第六阶段：监控和日志 (优先级：🔥 中)

#### 6.1 监控指标
```go
// 监控指标
type Metrics struct {
    TasksCreated    prometheus.Counter
    TasksAssigned   prometheus.Counter
    TasksCompleted  prometheus.Counter
    TasksFailed     prometheus.Counter
    QueueLength     prometheus.Gauge
    AssignmentTime  prometheus.Histogram
    ExecutionTime   prometheus.Histogram
}

// 在关键操作点记录指标
func (s *taskService) CreateTask(ctx context.Context, req *CreateTaskRequest) (*Task, error) {
    defer func(start time.Time) {
        s.metrics.TasksCreated.Inc()
    }(time.Now())
    
    // 创建任务逻辑
}
```

#### 6.2 日志记录
```go
// 结构化日志
func (s *assignmentService) assignTasks(ctx context.Context) error {
    logger := log.With("operation", "assign_tasks")
    
    logger.Info("Starting task assignment")
    
    // 分配逻辑
    
    logger.Info("Task assignment completed", 
        "assigned_count", assignedCount,
        "duration", time.Since(start))
    
    return nil
}
```

## 🚀 实施步骤建议

### 第一步：环境准备
```bash
# 1. 确保数据库运行正常
docker-compose up -d postgres

# 2. 执行migration
go run main.go migrate

# 3. 验证表结构
psql -h localhost -U onego -d onego_db -c "\dt"
```

### 第二步：核心模型定义
```bash
# 1. 创建模型文件
touch internal/biz/task/models.go
touch internal/biz/device/models.go

# 2. 定义核心结构体
# 3. 添加JSON标签和验证规则
```

### 第三步：Repository层
```bash
# 1. 创建Repository接口
# 2. 实现PostgreSQL版本
# 3. 添加单元测试
```

### 第四步：Service层
```bash
# 1. 实现业务逻辑
# 2. 添加事务处理
# 3. 实现错误处理
```

### 第五步：API层
```bash
# 1. 创建Controller
# 2. 配置路由
# 3. 添加中间件
```

### 第六步：调度器
```bash
# 1. 实现任务调度器
# 2. 实现设备执行器
# 3. 添加健康检查
```

## 📋 开发检查清单

### 代码质量
- [ ] 添加单元测试 (覆盖率 > 80%)
- [ ] 添加集成测试
- [ ] 代码格式化 (gofmt)
- [ ] 静态代码分析 (golangci-lint)
- [ ] 添加API文档 (Swagger)

### 性能优化
- [ ] 数据库连接池配置
- [ ] 查询优化和索引使用
- [ ] 缓存策略 (Redis)
- [ ] 并发控制
- [ ] 内存使用优化

### 监控和运维
- [ ] 健康检查接口
- [ ] 监控指标收集
- [ ] 日志记录和轮转
- [ ] 错误处理和恢复
- [ ] 配置管理

### 安全考虑
- [ ] 输入验证和清理
- [ ] SQL注入防护
- [ ] 认证和授权
- [ ] 敏感数据加密
- [ ] 访问控制

## 🎯 预期成果

完成实施后，您将拥有：

1. **完整的任务管理系统**
   - 任务创建、分配、执行、监控
   - 智能设备选择和负载均衡
   - 优先级调整和重试机制

2. **高性能的API服务**
   - RESTful API接口
   - 实时状态查询
   - 批量操作支持

3. **可靠的调度系统**
   - 自动任务分配
   - 设备执行器
   - 故障恢复机制

4. **完善的监控体系**
   - 性能指标收集
   - 实时监控面板
   - 告警机制

## 🚀 开始实施

建议从**第一步：环境准备**开始，逐步实施每个阶段。每个阶段完成后，建议进行测试验证，确保功能正常后再进入下一阶段。

如果您需要任何具体阶段的详细实现指导，我可以为您提供更详细的代码示例和最佳实践。 