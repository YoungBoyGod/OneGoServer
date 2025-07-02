# OneGo系统实施路线图

## 🗺️ 实施路线概览

```
SQL设计完成 → Go代码实现 → 测试验证 → 部署上线
     ↓
[当前阶段] → 下一步行动
```

## 🎯 当前状态

✅ **已完成**：
- 数据库设计 (4个migration文件)
- 任务生命周期流程设计
- 系统架构规划
- 文档和流程图

🔄 **下一步**：Go代码实现

## 🚀 具体实施步骤

### 第一步：环境准备和验证 (预计时间：30分钟)

#### 1.1 验证数据库环境
```bash
# 检查数据库连接
docker-compose ps postgres

# 如果需要启动数据库
docker-compose up -d postgres

# 验证连接
psql -h localhost -U onego -d onego_db -c "SELECT version();"
```

#### 1.2 执行Migration
```bash
# 检查现有migration文件
ls -la internal/data/migrations/

# 执行migration (需要先实现migration机制)
go run main.go migrate

# 验证表结构
psql -h localhost -U onego -d onego_db -c "\dt"
```

#### 1.3 创建Migration执行器
```go
// 在 main.go 中添加migration命令
func main() {
    if len(os.Args) > 1 && os.Args[1] == "migrate" {
        if err := runMigrations(); err != nil {
            log.Fatal("Migration failed:", err)
        }
        fmt.Println("Migration completed successfully")
        return
    }
    
    // 正常启动服务
    runServer()
}

func runMigrations() error {
    // 实现migration逻辑
    return nil
}
```

### 第二步：核心模型定义 (预计时间：1小时)

#### 2.1 创建任务模型
```go
// internal/biz/task/models.go
package task

import (
    "time"
    "encoding/json"
)

type Task struct {
    ID              int64           `json:"id" db:"id"`
    TaskID          string          `json:"task_id" db:"task_id"`
    Name            string          `json:"name" db:"name"`
    Description     string          `json:"description" db:"description"`
    Type            string          `json:"type" db:"type"`
    Status          string          `json:"status" db:"status"`
    Priority        int             `json:"priority" db:"priority"`
    Parameters      json.RawMessage `json:"parameters" db:"parameters"`
    ExecutorType    string          `json:"executor_type" db:"executor_type"`
    Timeout         int             `json:"timeout" db:"timeout"`
    MaxRetries      int             `json:"max_retries" db:"max_retries"`
    EstimatedDuration int           `json:"estimated_duration" db:"estimated_duration"`
    CreatedAt       time.Time       `json:"created_at" db:"created_at"`
    UpdatedAt       time.Time       `json:"updated_at" db:"updated_at"`
}

type TaskAssignmentQueue struct {
    ID                  int64           `json:"id" db:"id"`
    TaskID              string          `json:"task_id" db:"task_id"`
    Priority            int             `json:"priority" db:"priority"`
    QueueStatus         string          `json:"queue_status" db:"queue_status"`
    RequiredDeviceType  string          `json:"required_device_type" db:"required_device_type"`
    RequiredCapabilities json.RawMessage `json:"required_capabilities" db:"required_capabilities"`
    PreferredDeviceIDs  []int64         `json:"preferred_device_ids" db:"preferred_device_ids"`
    ExcludedDeviceIDs   []int64         `json:"excluded_device_ids" db:"excluded_device_ids"`
    AssignedDeviceID    *int64          `json:"assigned_device_id" db:"assigned_device_id"`
    AssignedAt          *time.Time      `json:"assigned_at" db:"assigned_at"`
    AssignmentScore     *float64        `json:"assignment_score" db:"assignment_score"`
    QueuePosition       *int            `json:"queue_position" db:"queue_position"`
    EstimatedWaitTime   *int            `json:"estimated_wait_time" db:"estimated_wait_time"`
    RetryCount          int             `json:"retry_count" db:"retry_count"`
    MaxRetries          int             `json:"max_retries" db:"max_retries"`
    QueuedAt            time.Time       `json:"queued_at" db:"queued_at"`
    UpdatedAt           time.Time       `json:"updated_at" db:"updated_at"`
}

type TaskAssignmentHistory struct {
    ID              int64           `json:"id" db:"id"`
    TaskID          string          `json:"task_id" db:"task_id"`
    DeviceID        *int64          `json:"device_id" db:"device_id"`
    Action          string          `json:"action" db:"action"`
    PreviousStatus  *string         `json:"previous_status" db:"previous_status"`
    NewStatus       *string         `json:"new_status" db:"new_status"`
    Reason          *string         `json:"reason" db:"reason"`
    Details         json.RawMessage `json:"details" db:"details"`
    CreatedAt       time.Time       `json:"created_at" db:"created_at"`
}
```

#### 2.2 创建设备模型
```go
// internal/biz/device/models.go
package device

import (
    "time"
    "encoding/json"
)

type Device struct {
    ID              int64           `json:"id" db:"id"`
    DeviceID        string          `json:"device_id" db:"device_id"`
    Name            string          `json:"name" db:"name"`
    Type            string          `json:"type" db:"type"`
    Model           *string         `json:"model" db:"model"`
    BoardID         *string         `json:"board_id" db:"board_id"`
    Status          string          `json:"status" db:"status"`
    HealthScore     int             `json:"health_score" db:"health_score"`
    LoginUsername   *string         `json:"login_username" db:"login_username"`
    LoginPort       *int            `json:"login_port" db:"login_port"`
    LoginPublicKey  *string         `json:"login_public_key" db:"login_public_key"`
    IPAddress       *string         `json:"ip_address" db:"ip_address"`
    Port            *int            `json:"port" db:"port"`
    Protocol        *string         `json:"protocol" db:"protocol"`
    Endpoint        *string         `json:"endpoint" db:"endpoint"`
    RegTime         *time.Time      `json:"reg_time" db:"reg_time"`
    Metadata        json.RawMessage `json:"metadata" db:"metadata"`
    Tags            json.RawMessage `json:"tags" db:"tags"`
    UptimeHours     float64         `json:"uptime_hours" db:"uptime_hours"`
    FirstOnlineTime *time.Time      `json:"first_online_time" db:"first_online_time"`
    LastOnlineTime  *time.Time      `json:"last_online_time" db:"last_online_time"`
    LastOfflineTime *time.Time      `json:"last_offline_time" db:"last_offline_time"`
    TotalOnlineDuration int64       `json:"total_online_duration" db:"total_online_duration"`
    TotalOfflineDuration int64      `json:"total_offline_duration" db:"total_offline_duration"`
    TotalHeartbeats int64           `json:"total_heartbeats" db:"total_heartbeats"`
    TotalAlerts     int64           `json:"total_alerts" db:"total_alerts"`
    TotalTasks      int64           `json:"total_tasks" db:"total_tasks"`
    TotalSuccessTasks int64         `json:"total_success_tasks" db:"total_success_tasks"`
    TotalFailedTasks int64          `json:"total_failed_tasks" db:"total_failed_tasks"`
    TotalCanceledTasks int64        `json:"total_canceled_tasks" db:"total_canceled_tasks"`
    TotalPendingTasks int64         `json:"total_pending_tasks" db:"total_pending_tasks"`
    TotalRunningTasks int64         `json:"total_running_tasks" db:"total_running_tasks"`
    TotalCompletedTasks int64       `json:"total_completed_tasks" db:"total_completed_tasks"`
    CreatedAt       time.Time       `json:"created_at" db:"created_at"`
    UpdatedAt       time.Time       `json:"updated_at" db:"updated_at"`
    CreatedBy       string          `json:"created_by" db:"created_by"`
    UpdatedBy       string          `json:"updated_by" db:"updated_by"`
}

type DeviceTaskQueue struct {
    ID                  int64           `json:"id" db:"id"`
    DeviceID            int64           `json:"device_id" db:"device_id"`
    TaskID              string          `json:"task_id" db:"task_id"`
    QueuePriority       int             `json:"queue_priority" db:"queue_priority"`
    QueuePosition       int             `json:"queue_position" db:"queue_position"`
    Status              string          `json:"status" db:"status"`
    EstimatedStartTime  *time.Time      `json:"estimated_start_time" db:"estimated_start_time"`
    ActualStartTime     *time.Time      `json:"actual_start_time" db:"actual_start_time"`
    ActualEndTime       *time.Time      `json:"actual_end_time" db:"actual_end_time"`
    EstimatedDuration   int             `json:"estimated_duration" db:"estimated_duration"`
    MaxRetryCount       int             `json:"max_retry_count" db:"max_retry_count"`
    CurrentRetry        int             `json:"current_retry" db:"current_retry"`
    TimeoutSeconds      int             `json:"timeout_seconds" db:"timeout_seconds"`
    RequeueCount        int             `json:"requeue_count" db:"requeue_count"`
    IsRequeued          bool            `json:"is_requeued" db:"is_requeued"`
    LastRequeueAt       *time.Time      `json:"last_requeue_at" db:"last_requeue_at"`
    CancelReason        *string         `json:"cancel_reason" db:"cancel_reason"`
    OriginalPriority    int             `json:"original_priority" db:"original_priority"`
    IsManualPriority    bool            `json:"is_manual_priority" db:"is_manual_priority"`
    IsManualPosition    bool            `json:"is_manual_position" db:"is_manual_position"`
    LastPriorityChangeAt *time.Time     `json:"last_priority_change_at" db:"last_priority_change_at"`
    LastPositionChangeAt *time.Time     `json:"last_position_change_at" db:"last_position_change_at"`
    LastAction          *string         `json:"last_action" db:"last_action"`
    CreatedAt           time.Time       `json:"created_at" db:"created_at"`
    UpdatedAt           time.Time       `json:"updated_at" db:"updated_at"`
}

type DeviceLoadMonitor struct {
    ID                  int64           `json:"id" db:"id"`
    DeviceID            int64           `json:"device_id" db:"device_id"`
    CurrentTasks        int             `json:"current_tasks" db:"current_tasks"`
    MaxConcurrentTasks  int             `json:"max_concurrent_tasks" db:"max_concurrent_tasks"`
    CPULoad             *float64        `json:"cpu_load" db:"cpu_load"`
    MemoryUsage         *float64        `json:"memory_usage" db:"memory_usage"`
    DiskUsage           *float64        `json:"disk_usage" db:"disk_usage"`
    NetworkLatency      *int            `json:"network_latency" db:"network_latency"`
    Status              string          `json:"status" db:"status"`
    LastHeartbeat       *time.Time      `json:"last_heartbeat" db:"last_heartbeat"`
    LoadScore           float64         `json:"load_score" db:"load_score"`
    TotalAssigned       int64           `json:"total_assigned" db:"total_assigned"`
    TotalCompleted      int64           `json:"total_completed" db:"total_completed"`
    TotalFailed         int64           `json:"total_failed" db:"total_failed"`
    SuccessRate         float64         `json:"success_rate" db:"success_rate"`
    UpdatedAt           time.Time       `json:"updated_at" db:"updated_at"`
}
```

### 第三步：Repository层实现 (预计时间：2小时)

#### 3.1 创建Repository接口
```go
// internal/data/task_repository.go
package data

import (
    "context"
    "github.com/your-org/OneGoServer002/internal/biz/task"
)

type TaskRepository interface {
    Create(ctx context.Context, task *task.Task) error
    GetByID(ctx context.Context, taskID string) (*task.Task, error)
    Update(ctx context.Context, task *task.Task) error
    Delete(ctx context.Context, taskID string) error
    List(ctx context.Context, offset, limit int) ([]*task.Task, error)
    Count(ctx context.Context) (int64, error)
}

type TaskAssignmentRepository interface {
    Enqueue(ctx context.Context, assignment *task.TaskAssignmentQueue) error
    Dequeue(ctx context.Context, limit int) ([]*task.TaskAssignmentQueue, error)
    AssignDevice(ctx context.Context, taskID string, deviceID int64) error
    UpdateStatus(ctx context.Context, taskID string, status string) error
    GetByTaskID(ctx context.Context, taskID string) (*task.TaskAssignmentQueue, error)
    UpdatePriority(ctx context.Context, taskID string, priority int) error
    GetQueueStatus(ctx context.Context) (*QueueStatus, error)
}

type DeviceTaskRepository interface {
    AddToQueue(ctx context.Context, deviceTask *task.DeviceTaskQueue) error
    GetNextTask(ctx context.Context, deviceID int64) (*task.DeviceTaskQueue, error)
    UpdateStatus(ctx context.Context, id int64, status string) error
    Requeue(ctx context.Context, id int64) error
    GetByID(ctx context.Context, id int64) (*task.DeviceTaskQueue, error)
    GetQueueStatus(ctx context.Context, deviceID int64) (*DeviceQueueStatus, error)
}

type DeviceLoadMonitorRepository interface {
    GetByDeviceID(ctx context.Context, deviceID int64) (*task.DeviceLoadMonitor, error)
    Update(ctx context.Context, monitor *task.DeviceLoadMonitor) error
    GetAvailableDevices(ctx context.Context, deviceType string, maxLoad float64) ([]*task.DeviceLoadMonitor, error)
    UpdateLoadScore(ctx context.Context, deviceID int64, loadScore float64) error
}
```

#### 3.2 实现PostgreSQL版本
```go
// internal/data/postgres_task_repository.go
package data

import (
    "context"
    "database/sql"
    "encoding/json"
    "fmt"
    "time"
    
    "github.com/your-org/OneGoServer002/internal/biz/task"
)

type postgresTaskRepository struct {
    db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
    return &postgresTaskRepository{db: db}
}

func (r *postgresTaskRepository) Create(ctx context.Context, t *task.Task) error {
    query := `
        INSERT INTO tasks (
            task_id, name, description, type, status, priority,
            parameters, executor_type, timeout, max_retries, estimated_duration
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        RETURNING id, created_at, updated_at
    `
    
    return r.db.QueryRowContext(ctx, query,
        t.TaskID, t.Name, t.Description, t.Type, t.Status, t.Priority,
        t.Parameters, t.ExecutorType, t.Timeout, t.MaxRetries, t.EstimatedDuration,
    ).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *postgresTaskRepository) GetByID(ctx context.Context, taskID string) (*task.Task, error) {
    query := `SELECT * FROM tasks WHERE task_id = $1`
    
    var t task.Task
    err := r.db.QueryRowContext(ctx, query, taskID).Scan(
        &t.ID, &t.TaskID, &t.Name, &t.Description, &t.Type, &t.Status, &t.Priority,
        &t.Parameters, &t.ExecutorType, &t.Timeout, &t.MaxRetries, &t.EstimatedDuration,
        &t.CreatedAt, &t.UpdatedAt,
    )
    
    if err != nil {
        return nil, err
    }
    
    return &t, nil
}

// 实现其他方法...
```

### 第四步：Service层实现 (预计时间：3小时)

#### 4.1 创建Service接口
```go
// internal/service/task.go
package service

import (
    "context"
    "github.com/your-org/OneGoServer002/internal/biz/task"
)

type TaskService interface {
    CreateTask(ctx context.Context, req *CreateTaskRequest) (*task.Task, error)
    GetTask(ctx context.Context, taskID string) (*task.Task, error)
    UpdateTask(ctx context.Context, taskID string, req *UpdateTaskRequest) error
    DeleteTask(ctx context.Context, taskID string) error
    ListTasks(ctx context.Context, offset, limit int) ([]*task.Task, error)
}

type CreateTaskRequest struct {
    Name            string          `json:"name" validate:"required"`
    Description     string          `json:"description"`
    Type            string          `json:"type" validate:"required"`
    Priority        int             `json:"priority" validate:"min=1,max=10"`
    Parameters      json.RawMessage `json:"parameters"`
    ExecutorType    string          `json:"executor_type" validate:"required"`
    Timeout         int             `json:"timeout"`
    MaxRetries      int             `json:"max_retries"`
    EstimatedDuration int           `json:"estimated_duration"`
}

type UpdateTaskRequest struct {
    Name            *string         `json:"name"`
    Description     *string         `json:"description"`
    Priority        *int            `json:"priority" validate:"omitempty,min=1,max=10"`
    Parameters      json.RawMessage `json:"parameters"`
    Timeout         *int            `json:"timeout"`
    MaxRetries      *int            `json:"max_retries"`
    EstimatedDuration *int          `json:"estimated_duration"`
}
```

#### 4.2 实现业务逻辑
```go
// internal/service/task_service.go
package service

import (
    "context"
    "fmt"
    "time"
    
    "github.com/your-org/OneGoServer002/internal/biz/task"
    "github.com/your-org/OneGoServer002/internal/data"
)

type taskService struct {
    taskRepo            data.TaskRepository
    assignmentRepo      data.TaskAssignmentRepository
    deviceTaskRepo      data.DeviceTaskRepository
    loadMonitorRepo     data.DeviceLoadMonitorRepository
}

func NewTaskService(
    taskRepo data.TaskRepository,
    assignmentRepo data.TaskAssignmentRepository,
    deviceTaskRepo data.DeviceTaskRepository,
    loadMonitorRepo data.DeviceLoadMonitorRepository,
) TaskService {
    return &taskService{
        taskRepo:        taskRepo,
        assignmentRepo:  assignmentRepo,
        deviceTaskRepo:  deviceTaskRepo,
        loadMonitorRepo: loadMonitorRepo,
    }
}

func (s *taskService) CreateTask(ctx context.Context, req *CreateTaskRequest) (*task.Task, error) {
    // 1. 验证请求
    if err := validateCreateTaskRequest(req); err != nil {
        return nil, fmt.Errorf("invalid request: %w", err)
    }
    
    // 2. 生成任务ID
    taskID := generateTaskID()
    
    // 3. 创建任务
    t := &task.Task{
        TaskID:          taskID,
        Name:            req.Name,
        Description:     req.Description,
        Type:            req.Type,
        Status:          "pending",
        Priority:        req.Priority,
        Parameters:      req.Parameters,
        ExecutorType:    req.ExecutorType,
        Timeout:         req.Timeout,
        MaxRetries:      req.MaxRetries,
        EstimatedDuration: req.EstimatedDuration,
    }
    
    if err := s.taskRepo.Create(ctx, t); err != nil {
        return nil, fmt.Errorf("failed to create task: %w", err)
    }
    
    // 4. 自动加入分配队列
    assignment := &task.TaskAssignmentQueue{
        TaskID:             taskID,
        Priority:           req.Priority,
        QueueStatus:        "queued",
        RequiredDeviceType: req.Type, // 可以根据需要设置
        MaxRetries:         req.MaxRetries,
    }
    
    if err := s.assignmentRepo.Enqueue(ctx, assignment); err != nil {
        return nil, fmt.Errorf("failed to enqueue task: %w", err)
    }
    
    return t, nil
}

func (s *taskService) GetTask(ctx context.Context, taskID string) (*task.Task, error) {
    return s.taskRepo.GetByID(ctx, taskID)
}

// 实现其他方法...
```

### 第五步：Controller层实现 (预计时间：2小时)

#### 5.1 创建Controller
```go
// internal/controller/task.go
package controller

import (
    "net/http"
    "strconv"
    
    "github.com/gin-gonic/gin"
    "github.com/your-org/OneGoServer002/internal/service"
)

type TaskController struct {
    taskService service.TaskService
}

func NewTaskController(taskService service.TaskService) *TaskController {
    return &TaskController{
        taskService: taskService,
    }
}

func (c *TaskController) CreateTask(ctx *gin.Context) {
    var req service.CreateTaskRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request body",
            "details": err.Error(),
        })
        return
    }
    
    task, err := c.taskService.CreateTask(ctx, &req)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to create task",
            "details": err.Error(),
        })
        return
    }
    
    ctx.JSON(http.StatusCreated, gin.H{
        "message": "Task created successfully",
        "data": task,
    })
}

func (c *TaskController) GetTask(ctx *gin.Context) {
    taskID := ctx.Param("task_id")
    if taskID == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "Task ID is required",
        })
        return
    }
    
    task, err := c.taskService.GetTask(ctx, taskID)
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{
            "error": "Task not found",
            "details": err.Error(),
        })
        return
    }
    
    ctx.JSON(http.StatusOK, gin.H{
        "data": task,
    })
}

func (c *TaskController) ListTasks(ctx *gin.Context) {
    offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
    limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
    
    if limit > 100 {
        limit = 100
    }
    
    tasks, err := c.taskService.ListTasks(ctx, offset, limit)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to list tasks",
            "details": err.Error(),
        })
        return
    }
    
    ctx.JSON(http.StatusOK, gin.H{
        "data": tasks,
        "pagination": gin.H{
            "offset": offset,
            "limit":  limit,
        },
    })
}

// 实现其他方法...
```

### 第六步：路由配置 (预计时间：30分钟)

#### 6.1 更新路由配置
```go
// internal/router/router.go
func (r *Router) setupTaskRoutes() {
    // 任务管理路由
    taskGroup := r.engine.Group("/api/v1/tasks")
    {
        taskGroup.POST("/", r.taskController.CreateTask)
        taskGroup.GET("/", r.taskController.ListTasks)
        taskGroup.GET("/:task_id", r.taskController.GetTask)
        taskGroup.PUT("/:task_id", r.taskController.UpdateTask)
        taskGroup.DELETE("/:task_id", r.taskController.DeleteTask)
    }
    
    // 任务分配路由
    assignmentGroup := r.engine.Group("/api/v1/assignments")
    {
        assignmentGroup.POST("/enqueue", r.assignmentController.EnqueueTask)
        assignmentGroup.GET("/status", r.assignmentController.GetQueueStatus)
        assignmentGroup.PUT("/:task_id/priority", r.assignmentController.UpdatePriority)
    }
    
    // 设备任务路由
    deviceGroup := r.engine.Group("/api/v1/devices")
    {
        deviceGroup.GET("/:device_id/next-task", r.deviceTaskController.GetNextTask)
        deviceGroup.POST("/:device_id/tasks/:task_id/start", r.deviceTaskController.StartTask)
        deviceGroup.POST("/:device_id/tasks/:task_id/complete", r.deviceTaskController.CompleteTask)
        deviceGroup.POST("/:device_id/tasks/:task_id/fail", r.deviceTaskController.FailTask)
    }
}
```

### 第七步：测试验证 (预计时间：1小时)

#### 7.1 创建测试文件
```go
// internal/service/task_service_test.go
package service

import (
    "context"
    "testing"
    "encoding/json"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestTaskService_CreateTask(t *testing.T) {
    // 创建mock repository
    mockTaskRepo := &MockTaskRepository{}
    mockAssignmentRepo := &MockTaskAssignmentRepository{}
    
    // 设置期望
    mockTaskRepo.On("Create", mock.Anything, mock.AnythingOfType("*task.Task")).Return(nil)
    mockAssignmentRepo.On("Enqueue", mock.Anything, mock.AnythingOfType("*task.TaskAssignmentQueue")).Return(nil)
    
    // 创建service
    service := NewTaskService(mockTaskRepo, mockAssignmentRepo, nil, nil)
    
    // 测试请求
    req := &CreateTaskRequest{
        Name:         "Test Task",
        Description:  "Test Description",
        Type:         "test",
        Priority:     5,
        ExecutorType: "shell",
    }
    
    // 执行测试
    result, err := service.CreateTask(context.Background(), req)
    
    // 验证结果
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, req.Name, result.Name)
    assert.Equal(t, "pending", result.Status)
    
    // 验证mock调用
    mockTaskRepo.AssertExpectations(t)
    mockAssignmentRepo.AssertExpectations(t)
}
```

#### 7.2 集成测试
```go
// test/integration/task_test.go
package integration

import (
    "testing"
    "net/http"
    "net/http/httptest"
    "bytes"
    "encoding/json"
    
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func TestTaskAPI_CreateTask(t *testing.T) {
    // 设置测试模式
    gin.SetMode(gin.TestMode)
    
    // 创建路由
    router := setupTestRouter()
    
    // 创建请求
    reqBody := map[string]interface{}{
        "name":         "Integration Test Task",
        "description":  "Test Description",
        "type":         "test",
        "priority":     5,
        "executor_type": "shell",
    }
    
    bodyBytes, _ := json.Marshal(reqBody)
    
    // 创建HTTP请求
    req := httptest.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(bodyBytes))
    req.Header.Set("Content-Type", "application/json")
    
    // 创建响应记录器
    w := httptest.NewRecorder()
    
    // 执行请求
    router.ServeHTTP(w, req)
    
    // 验证响应
    assert.Equal(t, http.StatusCreated, w.Code)
    
    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    
    assert.Equal(t, "Task created successfully", response["message"])
    assert.NotNil(t, response["data"])
}
```

## 🎯 实施时间估算

| 阶段 | 预计时间 | 关键交付物 |
|------|----------|------------|
| 环境准备 | 30分钟 | 数据库连接、Migration执行器 |
| 模型定义 | 1小时 | 核心数据结构 |
| Repository层 | 2小时 | 数据访问接口和实现 |
| Service层 | 3小时 | 业务逻辑实现 |
| Controller层 | 2小时 | API接口实现 |
| 路由配置 | 30分钟 | 路由定义 |
| 测试验证 | 1小时 | 单元测试和集成测试 |
| **总计** | **10小时** | **完整的基础功能** |

## 🚀 开始实施

### 立即行动项

1. **验证数据库环境**
   ```bash
   docker-compose up -d postgres
   psql -h localhost -U onego -d onego_db -c "SELECT version();"
   ```

2. **创建第一个模型文件**
   ```bash
   touch internal/biz/task/models.go
   # 开始编写Task结构体
   ```

3. **实现Migration执行器**
   ```bash
   # 在main.go中添加migrate命令
   ```

### 下一步建议

建议按照以下顺序实施：

1. **第一阶段**：环境准备 + 模型定义 (1.5小时)
2. **第二阶段**：Repository层 (2小时)
3. **第三阶段**：Service层 (3小时)
4. **第四阶段**：Controller层 (2小时)
5. **第五阶段**：测试验证 (1小时)

每个阶段完成后，建议进行简单的测试验证，确保功能正常后再进入下一阶段。

如果您需要任何具体阶段的详细实现指导，我可以为您提供更详细的代码示例和最佳实践。 