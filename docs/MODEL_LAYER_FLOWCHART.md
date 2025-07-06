# Model层流程图

## 概述

本文档描述了OneGoServer项目中Queue、Task、User三个模块Model层的数据流转和结构体关系图。

## 1. 整体架构流程图

```mermaid
flowchart TD
    A[API层] --> B[Model层]
    B --> C[Logic层]
    C --> D[Repository层]
    D --> E[DAO层]
    E --> F[数据库]
    
    B1[Input结构体] --> B
    B2[Output结构体] --> B
    B3[数据模型] --> B
    B4[过滤选项] --> B
    B5[排序选项] --> B
    
    B --> C1[业务逻辑处理]
    C --> D1[数据访问接口]
    D --> E1[具体数据操作]
```

## 2. Queue模块数据流转图

```mermaid
flowchart LR
    subgraph "队列基础操作"
        A1[CreateQueueInput] --> A2[CreateQueueOutput]
        B1[GetQueueListInput] --> B2[GetQueueListOutput]
        C1[UpdateQueueInput] --> C2[UpdateQueueOutput]
        D1[DeleteQueueInput] --> D2[DeleteQueueOutput]
    end
    
    subgraph "队列控制操作"
        E1[StartQueueInput] --> E2[StartQueueOutput]
        F1[PauseQueueInput] --> F2[PauseQueueOutput]
        G1[ResumeQueueInput] --> G2[ResumeQueueOutput]
        H1[StopQueueInput] --> H2[StopQueueOutput]
    end
    
    subgraph "队列任务管理"
        I1[GetQueueTasksInput] --> I2[GetQueueTasksOutput]
        J1[ReorderQueueTasksInput] --> J2[ReorderQueueTasksOutput]
        K1[RemoveQueueTaskInput] --> K2[RemoveQueueTaskOutput]
    end
    
    subgraph "队列监控统计"
        L1[GetQueueStatisticsInput] --> L2[GetQueueStatisticsOutput]
        M1[GetQueuePerformanceReportInput] --> M2[GetQueuePerformanceReportOutput]
    end
```

## 3. Task模块数据流转图

```mermaid
flowchart LR
    subgraph "任务基础操作"
        A1[CreateTaskInput] --> A2[CreateTaskOutput]
        B1[GetTaskListInput] --> B2[GetTaskListOutput]
        C1[UpdateTaskInput] --> C2[UpdateTaskOutput]
        D1[DeleteTaskInput] --> D2[DeleteTaskOutput]
    end
    
    subgraph "任务执行控制"
        E1[StartTaskInput] --> E2[StartTaskOutput]
        F1[StopTaskInput] --> F2[StopTaskOutput]
        G1[RestartTaskInput] --> G2[RestartTaskOutput]
        H1[CancelTaskInput] --> H2[CancelTaskOutput]
    end
    
    subgraph "任务优先级管理"
        I1[GetTaskPrioritiesInput] --> I2[GetTaskPrioritiesOutput]
        J1[UpdateTaskPriorityInput] --> J2[UpdateTaskPriorityOutput]
        K1[BatchUpdateTaskPriorityInput] --> K2[BatchUpdateTaskPriorityOutput]
    end
    
    subgraph "任务日志管理"
        L1[GetTaskLogsInput] --> L2[GetTaskLogsOutput]
        M1[GetTaskLogDetailInput] --> M2[GetTaskLogDetailOutput]
        N1[ClearTaskLogsInput] --> N2[ClearTaskLogsOutput]
    end
    
    subgraph "任务执行管理"
        O1[CreateTaskExecutionInput] --> O2[CreateTaskExecutionOutput]
        P1[GetTaskExecutionInput] --> P2[GetTaskExecutionOutput]
        Q1[UpdateTaskExecutionInput] --> Q2[UpdateTaskExecutionOutput]
    end
```

## 4. User模块数据流转图

```mermaid
flowchart LR
    subgraph "用户基础操作"
        A1[CreateUserInput] --> A2[CreateUserOutput]
        B1[GetUserListInput] --> B2[GetUserListOutput]
        C1[UpdateUserInput] --> C2[UpdateUserOutput]
        D1[DeleteUserInput] --> D2[DeleteUserOutput]
    end
    
    subgraph "用户认证授权"
        E1[LoginUserInput] --> E2[LoginUserOutput]
        F1[LogoutUserInput] --> F2[LogoutUserOutput]
        G1[ValidateUserLoginInput] --> G2[ValidateUserLoginOutput]
    end
    
    subgraph "用户权限管理"
        H1[GetUserRolesInput] --> H2[GetUserRolesOutput]
        I1[AssignUserRoleInput] --> I2[AssignUserRoleOutput]
        J1[RemoveUserRoleInput] --> J2[RemoveUserRoleOutput]
        K1[CheckUserPermissionInput] --> K2[CheckUserPermissionOutput]
    end
    
    subgraph "用户会话管理"
        L1[GetUserSessionsInput] --> L2[GetUserSessionsOutput]
        M1[RefreshTokenInput] --> M2[RefreshTokenOutput]
        N1[RevokeUserSessionInput] --> N2[RevokeUserSessionOutput]
    end
    
    subgraph "用户安全管理"
        O1[ChangePasswordInput] --> O2[ChangePasswordOutput]
        P1[ResetPasswordInput] --> P2[ResetPasswordOutput]
        Q1[GetUserSecurityLogInput] --> Q2[GetUserSecurityLogOutput]
    end
```

## 5. 数据模型关系图

```mermaid
erDiagram
    Queue ||--o{ QueueTask : contains
    Queue ||--o{ QueueConfigHistory : has
    Queue ||--|| QueueConfig : configures
    Queue ||--o{ QueueStatistics : generates
    
    Task ||--o{ TaskExecution : executes
    Task ||--o{ TaskAssignment : assigns
    Task ||--o{ TaskLog : logs
    Task ||--o{ TaskQueue : queues
    Task ||--o{ TaskSchedule : schedules
    
    User ||--o{ SessionInfo : has
    User ||--o{ ActivityInfo : performs
    User ||--o{ SecurityLogInfo : generates
    User ||--o{ RoleInfo : has
    User ||--o{ PermissionInfo : has
    
    QueueTask {
        string task_id
        string queue_id
        string status
        int priority
        int position
    }
    
    TaskExecution {
        string task_id
        string execution_id
        string device_id
        string status
        datetime start_time
        datetime end_time
    }
    
    TaskAssignment {
        string task_id
        int priority
        string queue_status
        string required_device_type
        int64 assigned_device_id
    }
    
    SessionInfo {
        string session_id
        string user_id
        string device_info
        string ip_address
        datetime login_time
        datetime expires_at
    }
```

## 6. 过滤和排序选项关系图

```mermaid
flowchart TD
    subgraph "Queue过滤选项"
        QF1[QueueFilter]
        QF2[TaskFilter]
        QF3[QueueSortOption]
        QF4[PaginationOption]
    end
    
    subgraph "Task过滤选项"
        TF1[TaskFilter]
        TF2[LogFilter]
        TF3[QueueFilter]
        TF4[ScheduleFilter]
        TF5[ExecutionFilter]
        TF6[AssignmentFilter]
        TF7[TaskSortOption]
    end
    
    subgraph "User过滤选项"
        UF1[UserFilter]
        UF2[UserSortOption]
        UF3[PaginationOption]
    end
    
    QF1 --> QF3
    QF3 --> QF4
    TF1 --> TF7
    TF7 --> QF4
    UF1 --> UF2
    UF2 --> UF3
```

## 7. 常量定义关系图

```mermaid
flowchart LR
    subgraph "Queue常量"
        QC1[QueueStatusActive]
        QC2[QueueStatusPaused]
        QC3[QueueStatusStopped]
        QC4[QueueTypeFIFO]
        QC5[QueueTypePriority]
        QC6[QueueOperationAdd]
        QC7[QueueOperationRemove]
    end
    
    subgraph "Task常量"
        TC1[TaskStatusPending]
        TC2[TaskStatusRunning]
        TC3[TaskStatusCompleted]
        TC4[TaskTypeBackup]
        TC5[TaskTypeSync]
        TC6[ExecutorTypeLocal]
        TC7[ExecutorTypeRemote]
    end
    
    subgraph "User常量"
        UC1[UserStatusActive]
        UC2[UserStatusInactive]
        UC3[UserRoleAdmin]
        UC4[UserRoleManager]
        UC5[LoginTypePassword]
        UC6[LoginTypeEmail]
        UC7[LogTypeLogin]
    end
```

## 8. Input/Output结构体分类图

```mermaid
flowchart TD
    subgraph "基础CRUD操作"
        CRUD1[CreateInput/Output]
        CRUD2[GetInput/Output]
        CRUD3[UpdateInput/Output]
        CRUD4[DeleteInput/Output]
    end
    
    subgraph "查询操作"
        QUERY1[ListInput/Output]
        QUERY2[FilterInput/Output]
        QUERY3[SortInput/Output]
        QUERY4[PaginationInput/Output]
    end
    
    subgraph "业务操作"
        BUSINESS1[ControlInput/Output]
        BUSINESS2[ManageInput/Output]
        BUSINESS3[ValidateInput/Output]
        BUSINESS4[StatisticsInput/Output]
    end
    
    subgraph "安全操作"
        SECURITY1[AuthInput/Output]
        SECURITY2[PermissionInput/Output]
        SECURITY3[SessionInput/Output]
        SECURITY4[SecurityInput/Output]
    end
    
    CRUD1 --> QUERY1
    QUERY1 --> BUSINESS1
    BUSINESS1 --> SECURITY1
```

## 9. 数据验证流程

```mermaid
flowchart TD
    A[接收Input数据] --> B{数据完整性检查}
    B -->|通过| C{业务规则验证}
    B -->|失败| D[返回验证错误]
    C -->|通过| E[数据转换处理]
    C -->|失败| F[返回业务错误]
    E --> G[调用Logic层]
    G --> H[Logic层处理]
    H --> I[返回Output数据]
    I --> J[数据格式化]
    J --> K[返回响应]
    
    D --> K
    F --> K
```

## 10. 错误处理流程

```mermaid
flowchart TD
    A[输入数据] --> B{Input验证}
    B -->|失败| C[InputValidationError]
    B -->|成功| D[业务处理]
    D --> E{业务逻辑验证}
    E -->|失败| F[BusinessLogicError]
    E -->|成功| G[数据处理]
    G --> H{数据处理验证}
    H -->|失败| I[DataProcessingError]
    H -->|成功| J[生成Output]
    J --> K[Output验证]
    K -->|失败| L[OutputValidationError]
    K -->|成功| M[返回成功响应]
    
    C --> N[错误响应]
    F --> N
    I --> N
    L --> N
```

## 总结

这些流程图展示了Queue、Task、User三个模块Model层的完整数据流转过程，包括：

1. **整体架构**：清晰的分层结构和数据流向
2. **模块内部流转**：每个模块的Input/Output结构体关系
3. **数据模型关系**：实体间的关联关系
4. **过滤排序选项**：查询条件的组织方式
5. **常量定义**：状态和类型的统一管理
6. **结构体分类**：按功能分类的组织方式
7. **数据验证**：完整的验证流程
8. **错误处理**：统一的错误处理机制

这些图表为开发团队提供了清晰的架构视图，有助于理解和使用Model层的各种结构体。 