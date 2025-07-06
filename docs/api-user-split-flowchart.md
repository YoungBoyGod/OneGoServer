# 用户API拆分流程图

## 拆分前结构
```mermaid
graph TD
    A[user.go - 500行] --> B[用户基础管理API]
    A --> C[用户权限管理API]
    A --> D[用户会话管理API]
    A --> E[用户安全管理API]
    A --> F[用户统计分析API]
    A --> G[通用数据模型]
    
    B --> B1[RegisterUserReq/Res]
    B --> B2[LoginUserReq/Res]
    B --> B3[LogoutUserReq/Res]
    B --> B4[GetUserListReq/Res]
    B --> B5[GetUserDetailReq/Res]
    B --> B6[UpdateUserReq/Res]
    
    C --> C1[GetUserRolesReq/Res]
    C --> C2[AssignUserRoleReq/Res]
    C --> C3[RemoveUserRoleReq/Res]
    C --> C4[GetUserPermissionsReq/Res]
    C --> C5[CheckUserPermissionReq/Res]
    
    D --> D1[GetUserSessionsReq/Res]
    D --> D2[RefreshTokenReq/Res]
    D --> D3[RevokeUserSessionReq/Res]
    D --> D4[GetUserActivityReq/Res]
    
    E --> E1[ChangePasswordReq/Res]
    E --> E2[ResetPasswordReq/Res]
    E --> E3[GetUserSecurityLogReq/Res]
    E --> E4[UpdateUserSecuritySettingsReq/Res]
    
    F --> F1[GetUserStatisticsReq/Res]
    F --> F2[GetUserBehaviorAnalysisReq/Res]
    
    G --> G1[UserInfo]
    G --> G2[UserDetailInfo]
    G --> G3[RoleInfo]
    G --> G4[PermissionInfo]
    G --> G5[SessionInfo]
    G --> G6[ActivityInfo]
    G --> G7[SecurityLogInfo]
    G --> G8[SecuritySettings]
    G --> G9[ActivitySummary]
    G --> G10[UserStatistics]
    G --> G11[UserTrendPoint]
    G --> G12[UserActivityRank]
    G --> G13[UserBehaviorAnalysis]
    G --> G14[TimeSlotActivity]
```

## 拆分后结构
```mermaid
graph TD
    A[api/user/v1/] --> B[basic.go]
    A --> C[permission.go]
    A --> D[session.go]
    A --> E[security.go]
    A --> F[statistics.go]
    A --> G[models.go]
    
    B --> B1[RegisterUserReq/Res]
    B --> B2[LoginUserReq/Res]
    B --> B3[LogoutUserReq/Res]
    B --> B4[GetUserListReq/Res]
    B --> B5[GetUserDetailReq/Res]
    B --> B6[UpdateUserReq/Res]
    
    C --> C1[GetUserRolesReq/Res]
    C --> C2[AssignUserRoleReq/Res]
    C --> C3[RemoveUserRoleReq/Res]
    C --> C4[GetUserPermissionsReq/Res]
    C --> C5[CheckUserPermissionReq/Res]
    
    D --> D1[GetUserSessionsReq/Res]
    D --> D2[RefreshTokenReq/Res]
    D --> D3[RevokeUserSessionReq/Res]
    D --> D4[GetUserActivityReq/Res]
    
    E --> E1[ChangePasswordReq/Res]
    E --> E2[ResetPasswordReq/Res]
    E --> E3[GetUserSecurityLogReq/Res]
    E --> E4[UpdateUserSecuritySettingsReq/Res]
    
    F --> F1[GetUserStatisticsReq/Res]
    F --> F2[GetUserBehaviorAnalysisReq/Res]
    
    G --> G1[UserInfo]
    G --> G2[UserDetailInfo]
    G --> G3[RoleInfo]
    G --> G4[PermissionInfo]
    G --> G5[SessionInfo]
    G --> G6[ActivityInfo]
    G --> G7[SecurityLogInfo]
    G --> G8[SecuritySettings]
    G --> G9[ActivitySummary]
    G --> G10[UserStatistics]
    G --> G11[UserTrendPoint]
    G --> G12[UserActivityRank]
    G --> G13[UserBehaviorAnalysis]
    G --> G14[TimeSlotActivity]
    
    G -.-> B
    G -.-> C
    G -.-> D
    G -.-> E
    G -.-> F
```

## 模块依赖关系
```mermaid
graph LR
    A[models.go] --> B[basic.go]
    A --> C[permission.go]
    A --> D[session.go]
    A --> E[security.go]
    A --> F[statistics.go]
    
    B --> G[用户基础管理]
    C --> H[用户权限管理]
    D --> I[用户会话管理]
    E --> J[用户安全管理]
    F --> K[用户统计分析]
```

## 功能分类图
```mermaid
graph TD
    A[用户API] --> B[管理类]
    A --> C[安全类]
    A --> D[分析类]
    
    B --> B1[basic.go - 基础管理]
    B --> B2[permission.go - 权限管理]
    B --> B3[session.go - 会话管理]
    
    C --> C1[security.go - 安全管理]
    
    D --> D1[statistics.go - 统计分析]
    
    E[models.go - 通用模型] --> F[支持所有模块]
```

## 拆分流程
```mermaid
flowchart TD
    A[分析原始user.go文件] --> B[识别功能模块]
    B --> C[确定拆分策略]
    C --> D[创建basic.go]
    C --> E[创建permission.go]
    C --> F[创建session.go]
    C --> G[创建security.go]
    C --> H[创建statistics.go]
    C --> I[创建models.go]
    
    D --> J[迁移基础管理API]
    E --> K[迁移权限管理API]
    F --> L[迁移会话管理API]
    G --> M[迁移安全管理API]
    H --> N[迁移统计分析API]
    I --> O[迁移通用数据模型]
    
    J --> P[删除原始user.go]
    K --> P
    L --> P
    M --> P
    N --> P
    O --> P
    
    P --> Q[验证拆分结果]
    Q --> R[生成文档]
    R --> S[提交代码]
```

## 模块职责分工
```mermaid
graph LR
    subgraph "管理模块"
        A[basic.go<br/>基础管理<br/>6个API]
        B[permission.go<br/>权限管理<br/>5个API]
        C[session.go<br/>会话管理<br/>4个API]
    end
    
    subgraph "安全模块"
        D[security.go<br/>安全管理<br/>4个API]
    end
    
    subgraph "分析模块"
        E[statistics.go<br/>统计分析<br/>2个API]
    end
    
    subgraph "通用模块"
        F[models.go<br/>数据模型<br/>15个模型]
    end
    
    F -.-> A
    F -.-> B
    F -.-> C
    F -.-> D
    F -.-> E
```

## 代码组织优化
```mermaid
graph TD
    A[原始状态] --> B[功能混杂<br/>500行单一文件]
    B --> C[拆分后状态]
    C --> D[6个模块文件<br/>功能清晰]
    
    D --> E[basic.go<br/>基础管理]
    D --> F[permission.go<br/>权限管理]
    D --> G[session.go<br/>会话管理]
    D --> H[security.go<br/>安全管理]
    D --> I[statistics.go<br/>统计分析]
    D --> J[models.go<br/>数据模型]
    
    E --> K[6个API]
    F --> L[5个API]
    G --> M[4个API]
    H --> N[4个API]
    I --> O[2个API]
    J --> P[15个模型]
```

## 维护性提升
```mermaid
graph LR
    A[拆分前] --> B[问题定位困难]
    A --> C[修改影响范围大]
    A --> D[团队协作冲突]
    
    E[拆分后] --> F[问题定位精确]
    E --> G[修改影响可控]
    E --> H[团队协作顺畅]
    
    B --> I[维护成本高]
    C --> I
    D --> I
    
    F --> J[维护成本低]
    G --> J
    H --> J
```

## API分布统计
```mermaid
pie title 用户API分布
    "基础管理" : 6
    "权限管理" : 5
    "会话管理" : 4
    "安全管理" : 4
    "统计分析" : 2
``` 