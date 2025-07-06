# API分页响应重构流程图

## 重构前结构
```mermaid
graph TD
    A[用户API] --> A1[GetUserListRes<br/>List []UserInfo<br/>Total int64<br/>Page int<br/>Size int]
    A --> A2[GetUserSessionsRes<br/>List []SessionInfo<br/>Total int64<br/>Page int<br/>Size int]
    A --> A3[GetUserSecurityLogRes<br/>List []SecurityLogInfo<br/>Total int64<br/>Page int<br/>Size int]
    
    B[任务API] --> B1[GetTaskListRes<br/>List []TaskInfo<br/>Total int64<br/>Page int<br/>Size int]
    B --> B2[GetTaskLogsRes<br/>List []TaskLogInfo<br/>Total int64<br/>Page int<br/>Size int]
    B --> B3[GetTaskAssignmentsRes<br/>List []TaskAssignmentInfo<br/>Total int64<br/>Page int<br/>Size int]
    
    C[系统API] --> C1[GetSystemLogsRes<br/>List []SystemLog<br/>Total int64<br/>Page int<br/>Size int]
    C --> C2[GetSystemAlertsRes<br/>List []SystemAlert<br/>Total int64<br/>Page int<br/>Size int]
    
    D[设备API] --> D1[GetDeviceListRes<br/>List []DeviceInfo<br/>Total int64<br/>Page int<br/>Size int]
    
    E[队列API] --> E1[GetQueueListRes<br/>List []QueueInfo<br/>Total int64<br/>Page int<br/>Size int]
```

## 重构后结构
```mermaid
graph TD
    A[通用分页结构] --> A1[PaginationResponse[T]<br/>List []T<br/>Total int64<br/>Page int<br/>Size int]
    A --> A2[PaginationRequest<br/>Page int<br/>Size int<br/>SortBy string<br/>SortOrder string]
    A --> A3[PaginationInfo<br/>CurrentPage int<br/>PageSize int<br/>TotalPages int<br/>TotalRecords int64<br/>HasNext bool<br/>HasPrev bool]
    
    B[用户API] --> B1[GetUserListRes<br/>common.PaginationResponse[UserInfo]]
    B --> B2[GetUserSessionsRes<br/>common.PaginationResponse[SessionInfo]]
    B --> B3[GetUserSecurityLogRes<br/>common.PaginationResponse[SecurityLogInfo]]
    
    C[任务API] --> C1[GetTaskListRes<br/>common.PaginationResponse[TaskInfo]]
    C --> C2[GetTaskLogsRes<br/>common.PaginationResponse[TaskLogInfo]]
    C --> C3[GetTaskAssignmentsRes<br/>common.PaginationResponse[TaskAssignmentInfo]]
    
    D[系统API] --> D1[GetSystemLogsRes<br/>common.PaginationResponse[SystemLog]]
    D --> D2[GetSystemAlertsRes<br/>common.PaginationResponse[SystemAlert]]
    
    E[设备API] --> E1[GetDeviceListRes<br/>common.PaginationResponse[DeviceInfo]]
    
    F[队列API] --> F1[GetQueueListRes<br/>common.PaginationResponse[QueueInfo]]
    
    A1 -.-> B1
    A1 -.-> B2
    A1 -.-> B3
    A1 -.-> C1
    A1 -.-> C2
    A1 -.-> C3
    A1 -.-> D1
    A1 -.-> D2
    A1 -.-> E1
    A1 -.-> F1
```

## 重构流程
```mermaid
flowchart TD
    A[分析现有分页结构] --> B[识别重复代码]
    B --> C[设计通用分页结构]
    C --> D[创建api/common/pagination.go]
    D --> E[更新用户API]
    E --> F[更新任务API]
    F --> G[更新系统API]
    G --> H[更新设备API]
    H --> I[更新队列API]
    I --> J[验证重构结果]
    J --> K[生成文档]
    K --> L[提交代码]
```

## 代码变化对比
```mermaid
graph LR
    A[重构前] --> B[重复代码<br/>每个API文件都定义<br/>相同的分页结构]
    C[重构后] --> D[统一结构<br/>使用泛型<br/>代码复用]
    
    B --> E[维护困难<br/>修改需要<br/>更新多个文件]
    D --> F[维护便利<br/>修改只需<br/>更新一个文件]
    
    E --> G[代码冗余<br/>类型不一致<br/>扩展困难]
    F --> H[代码简洁<br/>类型安全<br/>易于扩展]
```

## 模块依赖关系
```mermaid
graph TD
    A[api/common/pagination.go] --> B[用户API模块]
    A --> C[任务API模块]
    A --> D[系统API模块]
    A --> E[设备API模块]
    A --> F[队列API模块]
    
    B --> G[GetUserListRes]
    B --> H[GetUserSessionsRes]
    B --> I[GetUserSecurityLogRes]
    
    C --> J[GetTaskListRes]
    C --> K[GetTaskLogsRes]
    C --> L[GetTaskAssignmentsRes]
    
    D --> M[GetSystemLogsRes]
    D --> N[GetSystemAlertsRes]
    
    E --> O[GetDeviceListRes]
    
    F --> P[GetQueueListRes]
```

## 类型安全提升
```mermaid
graph LR
    A[泛型设计] --> B[编译时类型检查]
    B --> C[运行时类型安全]
    C --> D[减少类型错误]
    
    E[统一结构] --> F[一致的API响应]
    F --> G[便于前端处理]
    G --> H[提升开发效率]
```

## 维护性对比
```mermaid
graph TD
    A[重构前维护] --> B[需要修改多个文件]
    B --> C[容易遗漏某些文件]
    C --> D[维护成本高]
    
    E[重构后维护] --> F[只需修改一个文件]
    F --> G[自动应用到所有使用处]
    G --> H[维护成本低]
    
    I[功能扩展] --> J[重构前: 每个文件单独扩展]
    I --> K[重构后: 统一扩展]
    
    J --> L[扩展困难]
    K --> M[扩展便利]
```

## 代码质量提升
```mermaid
graph LR
    A[代码复用] --> B[减少重复代码]
    B --> C[提高代码质量]
    
    D[类型安全] --> E[减少运行时错误]
    E --> F[提高系统稳定性]
    
    G[统一规范] --> H[便于团队协作]
    H --> I[提高开发效率]
    
    C --> J[整体质量提升]
    F --> J
    I --> J
```

## 后续扩展计划
```mermaid
graph TD
    A[当前状态] --> B[基础分页结构]
    B --> C[计划扩展]
    
    C --> D[分页元数据]
    C --> E[分页缓存]
    C --> F[分页性能优化]
    
    D --> G[添加更多分页信息]
    E --> H[缓存分页结果]
    F --> I[优化分页查询]
    
    G --> J[完整的分页解决方案]
    H --> J
    I --> J
``` 