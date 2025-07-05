# 设备Model层流程图

## 概述

本文档描述了OneGoServer项目中设备Model层的数据流转和结构体关系图。

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
    
    B --> C1[业务逻辑处理]
    C --> D1[数据访问接口]
    D --> E1[具体数据操作]
```

## 2. Input/Output结构体关系图

```mermaid
flowchart LR
    subgraph "基础CRUD操作"
        A1[CreateDeviceInput] --> A2[CreateDeviceOutput]
        B1[GetDeviceByIDInput] --> B2[GetDeviceByIDOutput]
        C1[UpdateDeviceInput] --> C2[UpdateDeviceOutput]
        D1[DeleteDeviceInput] --> D2[DeleteDeviceOutput]
    end
    
    subgraph "查询操作"
        E1[GetDeviceListInput] --> E2[GetDeviceListOutput]
        F1[GetDevicesByTypeInput] --> F2[GetDevicesByTypeOutput]
        G1[GetDevicesByStatusInput] --> G2[GetDevicesByStatusOutput]
    end
    
    subgraph "状态管理"
        H1[UpdateDeviceStatusInput] --> H2[UpdateDeviceStatusOutput]
        I1[UpdateDeviceHealthScoreInput] --> I2[UpdateDeviceHealthScoreOutput]
        J1[BatchUpdateDeviceStatusInput] --> J2[BatchUpdateDeviceStatusOutput]
    end
    
    subgraph "统计查询"
        K1[GetDeviceStatisticsInput] --> K2[GetDeviceStatisticsOutput]
        L1[CountDevicesByStatusInput] --> L2[CountDevicesByStatusOutput]
        M1[CountDevicesByTypeInput] --> M2[CountDevicesByTypeOutput]
    end
```

## 3. 数据模型关系图

```mermaid
erDiagram
    Device ||--o{ DeviceHeartbeat : "has"
    Device ||--o{ DeviceLog : "generates"
    Device ||--o{ DeviceCommand : "receives"
    
    Device {
        int64 id PK
        string device_id
        string name
        string type
        string status
        int health_score
        string ip_address
        int port
        string protocol
        timestamp created_at
        timestamp updated_at
    }
    
    DeviceHeartbeat {
        int64 id PK
        int64 device_id FK
        timestamp heartbeat_time
        float cpu_usage
        float memory_usage
        float disk_usage
        string network_status
        int network_latency
        int running_tasks
        timestamp created_at
    }
    
    DeviceLog {
        int64 id PK
        int64 device_id FK
        string level
        string category
        string message
        timestamp log_time
        timestamp created_at
    }
    
    DeviceCommand {
        int64 id PK
        int64 device_id FK
        string command_id
        string command_type
        string command_data
        string status
        timestamp sent_time
        timestamp executed_time
        timestamp completed_time
        jsonb response_data
        string error_message
        timestamp created_at
    }
```

## 4. 过滤和排序选项流程图

```mermaid
flowchart TD
    A[GetDeviceListInput] --> B[DeviceFilter]
    A --> C[DeviceSortOption]
    A --> D[PaginationOption]
    
    B --> B1[Status过滤]
    B --> B2[Type过滤]
    B --> B3[健康度范围]
    B --> B4[时间范围]
    B --> B5[关键词搜索]
    
    C --> C1[排序字段]
    C --> C2[排序方向]
    
    D --> D1[页码]
    D --> D2[每页数量]
    
    B1 --> E[过滤结果]
    B2 --> E
    B3 --> E
    B4 --> E
    B5 --> E
    
    E --> F[排序处理]
    C1 --> F
    C2 --> F
    
    F --> G[分页处理]
    D1 --> G
    D2 --> G
    
    G --> H[GetDeviceListOutput]
```

## 5. 心跳管理流程图

```mermaid
flowchart TD
    A[设备发送心跳] --> B[CreateDeviceHeartbeatInput]
    B --> C[心跳数据验证]
    C --> D{验证通过?}
    
    D -->|是| E[保存心跳数据]
    D -->|否| F[返回错误]
    
    E --> G[更新设备状态]
    G --> H[CreateDeviceHeartbeatOutput]
    
    I[查询最新心跳] --> J[GetLatestDeviceHeartbeatInput]
    J --> K[获取设备ID]
    K --> L[查询数据库]
    L --> M[GetLatestDeviceHeartbeatOutput]
    
    N[查询心跳历史] --> O[GetDeviceHeartbeatHistoryInput]
    O --> P[设置查询参数]
    P --> Q[查询历史数据]
    Q --> R[GetDeviceHeartbeatHistoryOutput]
```

## 6. 日志管理流程图

```mermaid
flowchart TD
    A[设备生成日志] --> B[CreateDeviceLogInput]
    B --> C[日志数据验证]
    C --> D{验证通过?}
    
    D -->|是| E[保存日志数据]
    D -->|否| F[返回错误]
    
    E --> G[CreateDeviceLogOutput]
    
    H[查询设备日志] --> I[GetDeviceLogsInput]
    I --> J[设置过滤条件]
    J --> K[查询日志数据]
    K --> L[GetDeviceLogsOutput]
    
    M[按级别查询] --> N[GetDeviceLogsByLevelInput]
    N --> O[设置级别和限制]
    O --> P[查询指定级别日志]
    P --> Q[GetDeviceLogsByLevelOutput]
```

## 7. 命令管理流程图

```mermaid
flowchart TD
    A[发送设备命令] --> B[CreateDeviceCommandInput]
    B --> C[命令数据验证]
    C --> D{验证通过?}
    
    D -->|是| E[保存命令数据]
    D -->|否| F[返回错误]
    
    E --> G[CreateDeviceCommandOutput]
    
    H[查询命令状态] --> I[GetDeviceCommandInput]
    I --> J[获取命令ID]
    J --> K[查询命令详情]
    K --> L[GetDeviceCommandOutput]
    
    M[更新命令状态] --> N[UpdateDeviceCommandStatusInput]
    N --> O[设置新状态]
    O --> P[更新数据库]
    P --> Q[UpdateDeviceCommandStatusOutput]
    
    R[更新命令响应] --> S[UpdateDeviceCommandResponseInput]
    S --> T[设置响应数据]
    T --> U[更新响应信息]
    U --> V[UpdateDeviceCommandResponseOutput]
```

## 8. 统计查询流程图

```mermaid
flowchart TD
    A[获取设备统计] --> B[GetDeviceStatisticsInput]
    B --> C[设置过滤条件]
    C --> D[查询设备总数]
    D --> E[查询在线设备数]
    E --> F[查询离线设备数]
    F --> G[按状态分组统计]
    G --> H[按类型分组统计]
    H --> I[计算平均健康度]
    I --> J[GetDeviceStatisticsOutput]
    
    K[按状态统计] --> L[CountDevicesByStatusInput]
    L --> M[查询所有状态]
    M --> N[分组统计数量]
    N --> O[CountDevicesByStatusOutput]
    
    P[按类型统计] --> Q[CountDevicesByTypeInput]
    Q --> R[查询所有类型]
    R --> S[分组统计数量]
    S --> T[CountDevicesByTypeOutput]
```

## 9. 数据流转时序图

```mermaid
sequenceDiagram
    participant API as API层
    participant Model as Model层
    participant Logic as Logic层
    participant Repo as Repository层
    participant DAO as DAO层
    participant DB as 数据库
    
    API->>Model: 创建Input结构体
    Model->>Logic: 传递Input参数
    Logic->>Logic: 业务逻辑处理
    Logic->>Repo: 调用Repository方法
    Repo->>DAO: 调用DAO方法
    DAO->>DB: 执行SQL查询
    DB-->>DAO: 返回查询结果
    DAO-->>Repo: 返回数据
    Repo-->>Logic: 返回处理结果
    Logic->>Model: 创建Output结构体
    Model-->>API: 返回Output结果
```

## 10. 错误处理流程图

```mermaid
flowchart TD
    A[输入数据] --> B[数据验证]
    B --> C{验证通过?}
    
    C -->|是| D[业务处理]
    C -->|否| E[返回验证错误]
    
    D --> F{业务处理成功?}
    F -->|是| G[返回成功结果]
    F -->|否| H[返回业务错误]
    
    E --> I[错误响应]
    H --> I
    G --> J[成功响应]
    
    I --> K[统一错误格式]
    J --> L[统一成功格式]
    
    K --> M[Output结构体]
    L --> M
```

## 总结

这些流程图展示了设备Model层的完整数据流转过程：

1. **清晰的层次结构**：从API层到数据库的完整调用链
2. **标准化的Input/Output**：每个操作都有对应的输入输出结构体
3. **灵活的数据过滤**：支持多种过滤条件和排序选项
4. **完整的数据模型**：涵盖设备、心跳、日志、命令等所有相关数据
5. **规范的错误处理**：统一的错误处理和响应格式

这种设计确保了数据流转的清晰性和一致性，为整个设备管理系统提供了可靠的数据基础。 