# 用户模块重构流程图

## 整体重构流程

```mermaid
graph TD
    A[开始重构] --> B[分析API结构]
    B --> C[创建模型层结构体]
    C --> D[定义实体和服务接口]
    D --> E[拆分Logic层文件]
    E --> F[结构化改造方法]
    F --> G[删除原始文件]
    G --> H[生成文档]
    H --> I[完成重构]

    style A fill:#e1f5fe
    style I fill:#c8e6c9
```

## 文件拆分流程

```mermaid
graph LR
    A[原始user.go] --> B[auth.go]
    A --> C[validation.go]
    A --> D[session.go]
    A --> E[permission.go]
    A --> F[behavior.go]
    A --> G[risk.go]
    A --> H[utils.go]

    style A fill:#ffcdd2
    style B fill:#c8e6c9
    style C fill:#c8e6c9
    style D fill:#c8e6c9
    style E fill:#c8e6c9
    style F fill:#c8e6c9
    style G fill:#c8e6c9
    style H fill:#c8e6c9
```

## 结构化改造流程

```mermaid
graph TD
    A[原始方法] --> B[分析参数和返回值]
    B --> C[设计Input结构体]
    C --> D[设计Output结构体]
    D --> E[修改方法签名]
    E --> F[更新方法实现]
    F --> G[更新调用方]
    G --> H[验证类型安全]

    style A fill:#ffcdd2
    style H fill:#c8e6c9
```

## 模块职责划分

```mermaid
graph TB
    subgraph "认证授权模块 (auth.go)"
        A1[ValidateUserLogin]
        A2[HashPassword]
        A3[VerifyPassword]
        A4[CalculateUserRiskScore]
    end

    subgraph "用户验证模块 (validation.go)"
        B1[ValidateUserRegistration]
        B2[ValidateUsername]
        B3[ValidatePasswordStrength]
        B4[ValidateEmail]
        B5[ValidatePhone]
    end

    subgraph "会话管理模块 (session.go)"
        C1[CreateUserSession]
        C2[ValidateUserSession]
        C3[RefreshUserSession]
        C4[GenerateSessionId]
    end

    subgraph "权限管理模块 (permission.go)"
        D1[CalculateUserPermissions]
        D2[CheckUserPermission]
    end

    subgraph "行为分析模块 (behavior.go)"
        E1[AnalyzeUserBehavior]
        E2[FindMostActiveHour]
        E3[IdentifyBehaviorPattern]
        E4[CalculateActivityScore]
    end

    subgraph "风险评分模块 (risk.go)"
        F1[CalculateLoginRiskScore]
        F2[CalculateLocationRiskScore]
        F3[CalculateDeviceRiskScore]
        F4[CalculateTimeRiskScore]
        F5[CalculateBehaviorRiskScore]
        F6[CalculateHistoryRiskScore]
    end

    subgraph "辅助功能模块 (utils.go)"
        G1[IsSameNetwork]
        G2[IsKnownMaliciousIP]
        G3[IsMobileDevice]
        G4[IsAbnormalLoginPattern]
    end

    style A1 fill:#e3f2fd
    style B1 fill:#e3f2fd
    style C1 fill:#e3f2fd
    style D1 fill:#e3f2fd
    style E1 fill:#e3f2fd
    style F1 fill:#e3f2fd
    style G1 fill:#e3f2fd
```

## 数据流向

```mermaid
graph LR
    A[API层] --> B[Logic层]
    B --> C[Model层]
    C --> D[Entity层]
    
    subgraph "Logic层"
        B1[auth.go]
        B2[validation.go]
        B3[session.go]
        B4[permission.go]
        B5[behavior.go]
        B6[risk.go]
        B7[utils.go]
    end
    
    subgraph "Model层"
        C1[Input/Output结构体]
        C2[服务接口]
    end
    
    subgraph "Entity层"
        D1[User实体]
        D2[UserSession实体]
        D3[UserRole实体]
        D4[UserActivity实体]
    end

    style A fill:#fff3e0
    style B fill:#e8f5e8
    style C fill:#e3f2fd
    style D fill:#f3e5f5
```

## 重构前后对比

### 重构前
```mermaid
graph TD
    A[user.go] --> B[711行代码]
    B --> C[所有功能混在一起]
    C --> D[使用map[string]interface{}]
    D --> E[难以维护和扩展]
```

### 重构后
```mermaid
graph TD
    A[7个功能模块] --> B[清晰的功能边界]
    B --> C[结构化输入输出]
    C --> D[类型安全]
    D --> E[易于维护和扩展]
    E --> F[符合DDD架构]
```

## 质量提升指标

```mermaid
graph LR
    A[代码行数] --> B[平均每个文件100行]
    B --> C[可读性提升]
    
    D[类型安全] --> E[编译时检查]
    E --> F[运行时错误减少]
    
    G[模块化] --> H[单一职责]
    H --> I[可测试性提升]
    
    J[结构化] --> K[统一接口]
    K --> L[可维护性提升]

    style C fill:#c8e6c9
    style F fill:#c8e6c9
    style I fill:#c8e6c9
    style L fill:#c8e6c9
```

## 总结

通过本次重构，用户模块实现了：

1. **模块化拆分**: 从1个文件拆分为7个功能模块
2. **结构化改造**: 所有方法使用Input/Output结构体
3. **类型安全**: 编译时类型检查，减少运行时错误
4. **可维护性**: 清晰的功能边界，便于维护和扩展
5. **符合规范**: 与API层保持一致的设计风格
6. **DDD架构**: 遵循领域驱动设计原则

重构后的代码更加健壮、可维护、可扩展，为项目的长期发展奠定了良好的基础。 