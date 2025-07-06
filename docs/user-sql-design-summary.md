# 用户模块SQL表结构设计总结

## 概述

基于用户模块的API接口、Logic业务逻辑和Entity实体定义，设计了完整的用户管理数据库表结构，支持用户认证、授权、会话管理、活动跟踪、安全日志、设备管理、任务管理等所有功能。

## 表结构设计

### 1. 核心表结构

#### 1.1 用户主表 (users)
```sql
CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    phone VARCHAR(20),
    password VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    -- 基本信息
    avatar VARCHAR(255),
    nickname VARCHAR(50),
    real_name VARCHAR(50),
    gender VARCHAR(10),
    birthday DATE,
    department VARCHAR(100),
    position VARCHAR(100),
    employee_id VARCHAR(50),
    -- 登录信息
    last_login_at TIMESTAMP NULL,
    last_login_ip VARCHAR(45),
    last_user_agent TEXT,
    failed_login_attempts INT DEFAULT 0,
    account_locked_until TIMESTAMP NULL,
    -- 安全信息
    force_password_change BOOLEAN DEFAULT FALSE,
    password_changed_at TIMESTAMP NULL,
    recent_login_count INT DEFAULT 0,
    recent_security_events INT DEFAULT 0,
    -- 状态信息
    is_verified BOOLEAN DEFAULT FALSE,
    login_count INT DEFAULT 0,
    -- 时间戳
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);
```

#### 1.2 角色表 (roles)
```sql
CREATE TABLE roles (
    id VARCHAR(36) PRIMARY KEY,
    role_name VARCHAR(50) NOT NULL UNIQUE,
    role_desc TEXT,
    role_type VARCHAR(20) NOT NULL DEFAULT 'custom',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

#### 1.3 权限表 (permissions)
```sql
CREATE TABLE permissions (
    id VARCHAR(36) PRIMARY KEY,
    permission_name VARCHAR(100) NOT NULL UNIQUE,
    resource VARCHAR(100) NOT NULL,
    actions JSON NOT NULL,
    scope VARCHAR(20) DEFAULT 'global',
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### 2. 关联表结构

#### 2.1 用户角色关联表 (user_roles)
```sql
CREATE TABLE user_roles (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    role_id VARCHAR(36) NOT NULL,
    assigned_by VARCHAR(36),
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    UNIQUE KEY uk_user_role (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (assigned_by) REFERENCES users(id) ON DELETE SET NULL
);
```

#### 2.2 角色权限关联表 (role_permissions)
```sql
CREATE TABLE role_permissions (
    id VARCHAR(36) PRIMARY KEY,
    role_id VARCHAR(36) NOT NULL,
    permission_id VARCHAR(36) NOT NULL,
    granted_by VARCHAR(36),
    granted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    UNIQUE KEY uk_role_permission (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE,
    FOREIGN KEY (granted_by) REFERENCES users(id) ON DELETE SET NULL
);
```

### 3. 会话和活动表

#### 3.1 用户会话表 (user_sessions)
```sql
CREATE TABLE user_sessions (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    session_id VARCHAR(64) NOT NULL UNIQUE,
    ip_address VARCHAR(45),
    user_agent TEXT,
    device_info TEXT,
    location VARCHAR(100),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_activity TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

#### 3.2 用户活动表 (user_activities)
```sql
CREATE TABLE user_activities (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    session_id VARCHAR(64),
    action_type VARCHAR(50) NOT NULL,
    action_detail TEXT,
    resource_type VARCHAR(50),
    resource_id VARCHAR(36),
    ip_address VARCHAR(45),
    user_agent TEXT,
    session_duration FLOAT DEFAULT 0,
    result VARCHAR(20) DEFAULT 'success',
    metadata JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

### 4. 安全和统计表

#### 4.1 用户安全日志表 (user_security_logs)
```sql
CREATE TABLE user_security_logs (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36),
    event_type VARCHAR(50) NOT NULL,
    event_detail TEXT,
    level VARCHAR(20) NOT NULL DEFAULT 'info',
    ip_address VARCHAR(45),
    user_agent TEXT,
    risk_score FLOAT DEFAULT 0,
    risk_level VARCHAR(20) DEFAULT 'low',
    handled BOOLEAN DEFAULT FALSE,
    handled_by VARCHAR(36),
    handled_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (handled_by) REFERENCES users(id) ON DELETE SET NULL
);
```

#### 4.2 用户安全设置表 (user_security_settings)
```sql
CREATE TABLE user_security_settings (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL UNIQUE,
    enable_two_factor BOOLEAN DEFAULT FALSE,
    enable_email_notification BOOLEAN DEFAULT TRUE,
    enable_sms_notification BOOLEAN DEFAULT FALSE,
    session_timeout INT DEFAULT 7200,
    max_concurrent_sessions INT DEFAULT 5,
    password_complexity BOOLEAN DEFAULT TRUE,
    force_password_change BOOLEAN DEFAULT FALSE,
    password_expiry_days INT DEFAULT 90,
    login_notification BOOLEAN DEFAULT TRUE,
    suspicious_activity_alert BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

#### 4.3 用户统计表 (user_statistics)
```sql
CREATE TABLE user_statistics (
    id VARCHAR(36) PRIMARY KEY,
    stat_date DATE NOT NULL,
    total_users INT DEFAULT 0,
    active_users INT DEFAULT 0,
    inactive_users INT DEFAULT 0,
    new_users INT DEFAULT 0,
    login_count INT DEFAULT 0,
    failed_login_count INT DEFAULT 0,
    security_events INT DEFAULT 0,
    status_distribution JSON,
    role_distribution JSON,
    department_distribution JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    UNIQUE KEY uk_stat_date (stat_date)
);
```

### 5. 新增：用户设备管理表

#### 5.1 用户设备关联表 (user_devices)
```sql
CREATE TABLE user_devices (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    device_id VARCHAR(36) NOT NULL,
    device_name VARCHAR(100),
    device_type VARCHAR(50),
    access_level VARCHAR(20) NOT NULL DEFAULT 'read',
    is_primary BOOLEAN DEFAULT FALSE,
    is_trusted BOOLEAN DEFAULT FALSE,
    last_access_at TIMESTAMP NULL,
    access_count INT DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    UNIQUE KEY uk_user_device (user_id, device_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

#### 5.2 用户设备操作日志表 (user_device_logs)
```sql
CREATE TABLE user_device_logs (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    device_id VARCHAR(36) NOT NULL,
    operation_type VARCHAR(50) NOT NULL,
    operation_detail TEXT,
    result VARCHAR(20) DEFAULT 'success',
    ip_address VARCHAR(45),
    user_agent TEXT,
    session_id VARCHAR(64),
    duration FLOAT DEFAULT 0,
    error_message TEXT,
    metadata JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

### 6. 新增：用户任务管理表

#### 6.1 用户任务关联表 (user_tasks)
```sql
CREATE TABLE user_tasks (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    task_id VARCHAR(36) NOT NULL,
    assignment_type VARCHAR(20) NOT NULL DEFAULT 'assigned',
    role VARCHAR(50),
    priority VARCHAR(20) DEFAULT 'normal',
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMP NULL,
    completed_at TIMESTAMP NULL,
    due_date TIMESTAMP NULL,
    estimated_hours FLOAT DEFAULT 0,
    actual_hours FLOAT DEFAULT 0,
    progress_percentage INT DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    UNIQUE KEY uk_user_task (user_id, task_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

#### 6.2 用户任务操作日志表 (user_task_logs)
```sql
CREATE TABLE user_task_logs (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    task_id VARCHAR(36) NOT NULL,
    operation_type VARCHAR(50) NOT NULL,
    operation_detail TEXT,
    old_status VARCHAR(20),
    new_status VARCHAR(20),
    old_priority VARCHAR(20),
    new_priority VARCHAR(20),
    progress_change INT DEFAULT 0,
    time_spent FLOAT DEFAULT 0,
    ip_address VARCHAR(45),
    user_agent TEXT,
    session_id VARCHAR(64),
    metadata JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

## 索引设计

### 1. 主键索引
- 所有表都使用UUID作为主键，提供全局唯一性

### 2. 唯一索引
- `users.username` - 用户名唯一
- `users.email` - 邮箱唯一
- `users.phone` - 手机号唯一（可选）
- `user_roles(user_id, role_id)` - 用户角色关联唯一
- `role_permissions(role_id, permission_id)` - 角色权限关联唯一
- `user_sessions.session_id` - 会话ID唯一
- `user_security_settings.user_id` - 用户安全设置唯一
- `user_statistics.stat_date` - 统计日期唯一
- `user_devices(user_id, device_id)` - 用户设备关联唯一
- `user_tasks(user_id, task_id)` - 用户任务关联唯一

### 3. 普通索引
- `users.status` - 用户状态查询
- `users.department` - 部门查询
- `users.created_at` - 创建时间查询
- `users.last_login_at` - 最后登录时间查询
- `user_sessions.user_id` - 用户会话查询
- `user_sessions.status` - 会话状态查询
- `user_sessions.expires_at` - 会话过期查询
- `user_activities.user_id` - 用户活动查询
- `user_activities.action_type` - 活动类型查询
- `user_activities.created_at` - 活动时间查询
- `user_devices.user_id` - 用户设备查询
- `user_devices.device_id` - 设备用户查询
- `user_devices.device_type` - 设备类型查询
- `user_devices.access_level` - 访问级别查询
- `user_devices.status` - 设备状态查询
- `user_tasks.user_id` - 用户任务查询
- `user_tasks.task_id` - 任务用户查询
- `user_tasks.status` - 任务状态查询
- `user_tasks.priority` - 任务优先级查询
- `user_tasks.due_date` - 任务截止时间查询

## 外键约束

### 1. 级联删除
- `user_roles.user_id` -> `users.id` (CASCADE)
- `user_roles.role_id` -> `roles.id` (CASCADE)
- `role_permissions.role_id` -> `roles.id` (CASCADE)
- `role_permissions.permission_id` -> `permissions.id` (CASCADE)
- `user_sessions.user_id` -> `users.id` (CASCADE)
- `user_activities.user_id` -> `users.id` (CASCADE)
- `user_security_settings.user_id` -> `users.id` (CASCADE)
- `user_devices.user_id` -> `users.id` (CASCADE)
- `user_tasks.user_id` -> `users.id` (CASCADE)
- `user_device_logs.user_id` -> `users(id) (CASCADE)
- `user_task_logs.user_id` -> `users(id) (CASCADE)

### 2. 设置空值
- `user_roles.assigned_by` -> `users.id` (SET NULL)
- `role_permissions.granted_by` -> `users.id` (SET NULL)
- `user_security_logs.user_id` -> `users.id` (SET NULL)
- `user_security_logs.handled_by` -> `users.id` (SET NULL)

## 初始化数据

### 1. 默认角色
- `admin` - 系统管理员
- `user` - 普通用户
- `guest` - 访客用户

### 2. 默认权限
- `user:read` - 读取用户信息
- `user:write` - 创建和更新用户
- `user:delete` - 删除用户
- `role:manage` - 角色管理
- `permission:manage` - 权限管理
- `system:admin` - 系统管理

### 3. 默认管理员
- 用户名: `admin`
- 邮箱: `admin@example.com`
- 角色: `admin`
- 状态: `active`

## 功能支持

### 1. 用户认证
- 用户名/邮箱/手机号登录
- 密码加密存储
- 登录失败次数限制
- 账户锁定机制
- 强制密码修改

### 2. 用户授权
- 基于角色的访问控制(RBAC)
- 细粒度权限控制
- 权限继承和组合
- 动态权限分配

### 3. 会话管理
- 多会话支持
- 会话超时控制
- 会话状态跟踪
- 设备信息记录

### 4. 活动跟踪
- 用户行为记录
- 操作审计
- 资源访问日志
- 会话时长统计

### 5. 安全监控
- 安全事件记录
- 风险评分
- 异常行为检测
- 安全设置管理

### 6. 统计分析
- 用户数量统计
- 活跃度分析
- 登录趋势
- 安全事件统计

### 7. 新增：设备管理
- 用户设备关联
- 设备访问权限控制
- 设备操作日志
- 设备状态跟踪
- 可信设备管理

### 8. 新增：任务管理
- 用户任务分配
- 任务角色管理
- 任务进度跟踪
- 工时统计
- 任务操作日志

## 扩展性设计

### 1. 字段扩展
- 使用JSON字段存储灵活数据
- 预留扩展字段
- 支持自定义属性

### 2. 功能扩展
- 支持多租户
- 支持组织架构
- 支持工作流
- 支持通知系统
- 支持设备集群管理
- 支持任务工作流

### 3. 性能优化
- 合理的索引设计
- 分区表支持
- 读写分离支持
- 缓存友好设计

## 总结

该SQL设计完全支持用户模块的所有API功能和业务逻辑，提供了：

1. **完整的功能支持** - 覆盖认证、授权、会话、活动、安全、设备、任务等所有功能
2. **良好的性能** - 合理的索引设计和查询优化
3. **数据完整性** - 外键约束和数据验证
4. **扩展性** - 支持未来功能扩展
5. **安全性** - 密码加密、权限控制、审计日志
6. **设备管理** - 用户设备关联、权限控制、操作日志
7. **任务管理** - 任务分配、进度跟踪、工时统计

该设计为OneGoServer项目的用户管理功能提供了坚实的数据基础，支持完整的用户生命周期管理。 