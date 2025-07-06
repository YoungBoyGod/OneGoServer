# 通用抽象层设计流程图

## 1. 通用抽象层架构图

```mermaid
graph TB
    subgraph "通用抽象层 (internal/model/common)"
        A[通用常量] --> A1[状态常量]
        A --> A2[排序常量]
        A --> A3[操作类型常量]
        A --> A4[优先级常量]
        
        B[通用基础结构体] --> B1[BaseEntity]
        B --> B2[BaseEntityWithStatus]
        B --> B3[BaseEntityWithName]
        B --> B4[BaseEntityWithPriority]
        
        C[通用Input/Output] --> C1[BaseCreateInput/Output]
        C --> C2[BaseGetByIDInput/Output]
        C --> C3[BaseListInput/Output]
        C --> C4[BaseActionInput/Output]
        
        D[通用过滤排序] --> D1[PaginationOption]
        D --> D2[BaseFilter]
        D --> D3[BaseSortOption]
        
        E[通用接口] --> E1[Entity接口]
        E --> E2[Filter接口]
        E --> E3[Sort接口]
        
        F[通用工具函数] --> F1[验证函数]
        F --> F2[默认值函数]
        F --> F3[错误处理]
    end
    
    subgraph "业务模块层"
        G[Device模块] --> G1[Device结构体]
        G --> G2[DeviceFilter]
        G --> G3[DeviceInput/Output]
        
        H[Task模块] --> H1[Task结构体]
        H --> H2[TaskFilter]
        H --> H3[TaskInput/Output]
        
        I[Queue模块] --> I1[Queue结构体]
        I --> I2[QueueFilter]
        I --> I3[QueueInput/Output]
        
        J[User模块] --> J1[User结构体]
        J --> J2[UserFilter]
        J --> J3[UserInput/Output]
    end
    
    subgraph "继承关系"
        G1 -.-> B2
        G2 -.-> D2
        G3 -.-> C1
        
        H1 -.-> B4
        H2 -.-> D2
        H3 -.-> C1
        
        I1 -.-> B3
        I2 -.-> D2
        I3 -.-> C1
        
        J1 -.-> B2
        J2 -.-> D2
        J3 -.-> C1
    end
    
    subgraph "接口实现"
        G1 -.-> E1
        G2 -.-> E2
        H1 -.-> E1
        H2 -.-> E2
        I1 -.-> E1
        I2 -.-> E2
        J1 -.-> E1
        J2 -.-> E2
    end
```

## 2. 通用抽象层设计流程

```mermaid
flowchart TD
    A[分析现有代码] --> B[识别重复模式]
    B --> C[提取通用部分]
    C --> D[设计通用结构体]
    D --> E[定义通用接口]
    E --> F[创建通用工具函数]
    F --> G[编写使用指南]
    G --> H[重构现有模块]
    H --> I[测试验证]
    I --> J[文档更新]
    J --> K[部署应用]
    
    subgraph "重复模式识别"
        B1[PaginationOption重复]
        B2[BaseFilter类似]
        B3[状态常量重复]
        B4[验证逻辑重复]
    end
    
    subgraph "通用部分提取"
        C1[基础字段提取]
        C2[常量统一]
        C3[接口标准化]
        C4[工具函数复用]
    end
    
    subgraph "结构体设计"
        D1[BaseEntity设计]
        D2[BaseFilter设计]
        D3[BaseSortOption设计]
        D4[Input/Output设计]
    end
    
    B --> B1
    B --> B2
    B --> B3
    B --> B4
    
    C --> C1
    C --> C2
    C --> C3
    C --> C4
    
    D --> D1
    D --> D2
    D --> D3
    D --> D4
```

## 3. 模块继承关系图

```mermaid
classDiagram
    class BaseEntity {
        +ID int64
        +CreatedAt *gtime.Time
        +UpdatedAt *gtime.Time
        +CreatedBy int64
        +UpdatedBy int64
    }
    
    class BaseEntityWithStatus {
        +Status string
    }
    
    class BaseEntityWithName {
        +Name string
        +Description string
    }
    
    class BaseEntityWithPriority {
        +Priority int
    }
    
    class BaseFilter {
        +Status []string
        +Type []string
        +Keyword *string
        +StartTime *time.Time
        +EndTime *time.Time
        +Enabled *bool
    }
    
    class BaseSortOption {
        +Field string
        +Order string
    }
    
    class Device {
        +DeviceID string
        +DeviceType string
        +Model string
        +Manufacturer string
        +IPAddress string
        +Port int
        +Protocol string
        +HealthScore int
        +LastSeen *gtime.Time
        +IsOnline bool
        +IsEnabled bool
    }
    
    class Task {
        +TaskID string
        +TaskType string
        +Status string
        +ExecutorType string
        +DeviceID string
        +Parameters string
        +Result string
        +StartTime *gtime.Time
        +EndTime *gtime.Time
    }
    
    class Queue {
        +QueueID string
        +QueueType string
        +Status string
        +Capacity int
        +CurrentSize int
        +MaxConcurrency int
        +ProcessingRate float64
        +IsEnabled bool
    }
    
    class User {
        +UserID string
        +Username string
        +Email string
        +Status string
        +Role string
        +Department string
        +LastLogin *gtime.Time
        +IsVerified bool
    }
    
    class DeviceFilter {
        +Manufacturer *string
        +Model *string
        +Protocol *string
        +MinHealth *int
        +MaxHealth *int
        +LastSeenFrom *time.Time
        +LastSeenTo *time.Time
    }
    
    class TaskFilter {
        +Priority *int
        +DeviceID *string
        +ExecutorType *string
        +IsUrgent *bool
        +StartTime *time.Time
        +EndTime *time.Time
    }
    
    class QueueFilter {
        +Priority *int
        +Keyword *string
        +StartTime *time.Time
        +EndTime *time.Time
        +Enabled *bool
    }
    
    class UserFilter {
        +Role []string
        +Department []string
        +IsVerified *bool
        +StartTime *time.Time
        +EndTime *time.Time
    }
    
    BaseEntity <|-- BaseEntityWithStatus
    BaseEntity <|-- BaseEntityWithName
    BaseEntity <|-- BaseEntityWithPriority
    
    BaseEntityWithName <|-- Device
    BaseEntityWithPriority <|-- Task
    BaseEntityWithName <|-- Queue
    BaseEntityWithStatus <|-- User
    
    BaseFilter <|-- DeviceFilter
    BaseFilter <|-- TaskFilter
    BaseFilter <|-- QueueFilter
    BaseFilter <|-- UserFilter
```

## 4. Input/Output结构体继承关系

```mermaid
classDiagram
    class BaseCreateInput {
        +Entity interface{}
    }
    
    class BaseCreateOutput {
        +ID string
        +Message string
    }
    
    class BaseGetByIDInput {
        +ID int64
    }
    
    class BaseGetByIDOutput {
        +Entity interface{}
    }
    
    class BaseListInput {
        +Filter interface{}
        +Sort interface{}
        +Pagination *PaginationOption
    }
    
    class BaseListOutput {
        +List interface{}
        +Total int64
    }
    
    class BaseActionInput {
        +ID string
        +Action string
        +Reason string
        +Force bool
    }
    
    class BaseActionOutput {
        +Message string
    }
    
    class CreateDeviceInput {
        +Device *Device
    }
    
    class CreateDeviceOutput {
        +DeviceID string
    }
    
    class GetDeviceByIDInput {
        +ID int64
    }
    
    class GetDeviceByIDOutput {
        +Device *Device
    }
    
    class GetDeviceListInput {
        +Filter *DeviceFilter
        +Sort *DeviceSortOption
        +Pagination *PaginationOption
    }
    
    class GetDeviceListOutput {
        +List []Device
        +Total int64
    }
    
    class UpdateDeviceStatusInput {
        +DeviceID string
        +Status string
        +Reason string
    }
    
    BaseCreateInput <|-- CreateDeviceInput
    BaseCreateOutput <|-- CreateDeviceOutput
    BaseGetByIDInput <|-- GetDeviceByIDInput
    BaseGetByIDOutput <|-- GetDeviceByIDOutput
    BaseListInput <|-- GetDeviceListInput
    BaseListOutput <|-- GetDeviceListOutput
    BaseActionInput <|-- UpdateDeviceStatusInput
    BaseActionOutput <|-- UpdateDeviceStatusOutput
```

## 5. 接口实现关系图

```mermaid
classDiagram
    class Entity {
        <<interface>>
        +GetID() int64
        +GetCreatedAt() *gtime.Time
        +GetUpdatedAt() *gtime.Time
    }
    
    class Filter {
        <<interface>>
        +GetStatus() []string
        +GetType() []string
        +GetKeyword() *string
        +GetStartTime() *time.Time
        +GetEndTime() *time.Time
    }
    
    class Sort {
        <<interface>>
        +GetField() string
        +GetOrder() string
    }
    
    class Statistics {
        <<interface>>
        +GetTotal() int64
        +GetByStatus() map[string]int64
        +GetByType() map[string]int64
    }
    
    class Device {
        +GetID() int64
        +GetCreatedAt() *gtime.Time
        +GetUpdatedAt() *gtime.Time
    }
    
    class Task {
        +GetID() int64
        +GetCreatedAt() *gtime.Time
        +GetUpdatedAt() *gtime.Time
    }
    
    class Queue {
        +GetID() int64
        +GetCreatedAt() *gtime.Time
        +GetUpdatedAt() *gtime.Time
    }
    
    class User {
        +GetID() int64
        +GetCreatedAt() *gtime.Time
        +GetUpdatedAt() *gtime.Time
    }
    
    class DeviceFilter {
        +GetStatus() []string
        +GetType() []string
        +GetKeyword() *string
        +GetStartTime() *time.Time
        +GetEndTime() *time.Time
    }
    
    class TaskFilter {
        +GetStatus() []string
        +GetType() []string
        +GetKeyword() *string
        +GetStartTime() *time.Time
        +GetEndTime() *time.Time
    }
    
    class DeviceSortOption {
        +GetField() string
        +GetOrder() string
    }
    
    class TaskSortOption {
        +GetField() string
        +GetOrder() string
    }
    
    class DeviceStatistics {
        +GetTotal() int64
        +GetByStatus() map[string]int64
        +GetByType() map[string]int64
    }
    
    class TaskStatistics {
        +GetTotal() int64
        +GetByStatus() map[string]int64
        +GetByType() map[string]int64
    }
    
    Entity <|.. Device
    Entity <|.. Task
    Entity <|.. Queue
    Entity <|.. User
    
    Filter <|.. DeviceFilter
    Filter <|.. TaskFilter
    
    Sort <|.. DeviceSortOption
    Sort <|.. TaskSortOption
    
    Statistics <|.. DeviceStatistics
    Statistics <|.. TaskStatistics
```

## 6. 工具函数使用流程

```mermaid
flowchart TD
    A[数据输入] --> B{数据类型}
    
    B -->|状态| C[IsValidStatus]
    B -->|排序| D[IsValidSortOrder]
    B -->|优先级| E[IsValidPriority]
    B -->|日志级别| F[IsValidLogLevel]
    
    C --> G{状态有效?}
    D --> H{排序有效?}
    E --> I{优先级有效?}
    F --> J{日志级别有效?}
    
    G -->|是| K[继续处理]
    G -->|否| L[返回错误]
    
    H -->|是| K
    H -->|否| L
    
    I -->|是| K
    I -->|否| L
    
    J -->|是| K
    J -->|否| L
    
    K --> M[设置默认值]
    M --> N[GetDefaultPagination]
    M --> O[GetDefaultSort]
    
    N --> P[返回结果]
    O --> P
    
    L --> Q[NewCommonError]
    Q --> R[返回错误响应]
    
    subgraph "验证函数"
        C
        D
        E
        F
    end
    
    subgraph "默认值函数"
        N
        O
    end
    
    subgraph "错误处理"
        Q
    end
```

## 7. 重构前后对比图

```mermaid
graph LR
    subgraph "重构前"
        A1[Device模块]
        A2[Task模块]
        A3[Queue模块]
        A4[User模块]
        
        B1[重复的PaginationOption]
        B2[重复的状态常量]
        B3[重复的验证逻辑]
        B4[重复的基础字段]
        
        A1 --> B1
        A1 --> B2
        A1 --> B3
        A1 --> B4
        
        A2 --> B1
        A2 --> B2
        A2 --> B3
        A2 --> B4
        
        A3 --> B1
        A3 --> B2
        A3 --> B3
        A3 --> B4
        
        A4 --> B1
        A4 --> B2
        A4 --> B3
        A4 --> B4
    end
    
    subgraph "重构后"
        C1[Device模块]
        C2[Task模块]
        C3[Queue模块]
        C4[User模块]
        
        D1[通用抽象层]
        D2[BaseEntity]
        D3[BaseFilter]
        D4[通用工具函数]
        
        C1 --> D1
        C2 --> D1
        C3 --> D1
        C4 --> D1
        
        D1 --> D2
        D1 --> D3
        D1 --> D4
    end
    
    A1 -.-> C1
    A2 -.-> C2
    A3 -.-> C3
    A4 -.-> C4
```

## 8. 代码复用效果图

```mermaid
pie title 代码复用效果
    "减少重复代码" : 40
    "统一常量管理" : 25
    "标准化接口" : 20
    "工具函数复用" : 15
```

## 9. 开发效率提升图

```mermaid
graph TD
    A[新模块开发] --> B{使用通用抽象层?}
    
    B -->|是| C[继承通用结构体]
    B -->|否| D[从头编写所有代码]
    
    C --> E[实现特定接口]
    E --> F[添加业务逻辑]
    F --> G[测试验证]
    G --> H[完成开发]
    
    D --> I[定义基础字段]
    I --> J[编写验证逻辑]
    J --> K[实现接口方法]
    K --> L[添加业务逻辑]
    L --> M[测试验证]
    M --> N[完成开发]
    
    H --> O[开发时间: 2天]
    N --> P[开发时间: 5天]
    
    subgraph "使用通用抽象层"
        C
        E
        F
        G
        H
        O
    end
    
    subgraph "传统开发方式"
        D
        I
        J
        K
        L
        M
        N
        P
    end
```

## 10. 维护成本对比图

```mermaid
graph LR
    subgraph "维护成本"
        A1[通用逻辑修改]
        A2[模块特定修改]
        A3[接口变更]
        A4[常量更新]
    end
    
    subgraph "传统方式"
        B1[修改4个模块]
        B2[修改4个模块]
        B3[修改4个模块]
        B4[修改4个模块]
    end
    
    subgraph "通用抽象层"
        C1[修改1个通用文件]
        C2[修改4个模块]
        C3[修改1个通用接口]
        C4[修改1个常量文件]
    end
    
    A1 --> B1
    A1 --> C1
    
    A2 --> B2
    A2 --> C2
    
    A3 --> B3
    A3 --> C3
    
    A4 --> B4
    A4 --> C4
    
    B1 --> D1[成本: 4倍]
    C1 --> D2[成本: 1倍]
    
    B2 --> D3[成本: 4倍]
    C2 --> D4[成本: 4倍]
    
    B3 --> D5[成本: 4倍]
    C3 --> D6[成本: 1倍]
    
    B4 --> D7[成本: 4倍]
    C4 --> D8[成本: 1倍]
```

## 总结

通过以上流程图，我们可以清楚地看到通用抽象层的设计优势：

1. **代码复用**：通过继承通用结构体，大幅减少重复代码
2. **一致性维护**：统一的接口和常量定义，保证代码一致性
3. **开发效率**：新模块可以快速继承通用功能，提高开发速度
4. **维护成本**：通用逻辑的修改只需要在一个地方进行，降低维护成本
5. **扩展性**：新增模块可以轻松继承通用功能，增强系统扩展性

通用抽象层的设计为OneGoServer项目提供了强大的代码复用和一致性维护能力，是提高代码质量和开发效率的重要工具。 