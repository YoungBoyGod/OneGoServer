-- =================================================================================
-- 用户模块数据库表结构
-- 创建时间: 2024-12-30
-- 描述: 用户管理相关的所有表结构
-- =================================================================================

-- 1. 用户主表
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY COMMENT '用户ID，UUID格式',
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名，唯一',
    email VARCHAR(100) NOT NULL UNIQUE COMMENT '邮箱，唯一',
    password VARCHAR(255) NOT NULL COMMENT '密码哈希',
    status VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '用户状态：active, inactive, locked, deleted',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL COMMENT '删除时间',  
    INDEX idx_username (username),
    INDEX idx_email (email),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
) COMMENT '用户主表';

-- 2. 角色表
CREATE TABLE IF NOT EXISTS roles (
    id VARCHAR(36) PRIMARY KEY COMMENT '角色ID，UUID格式',
    role_name VARCHAR(50) NOT NULL UNIQUE COMMENT '角色名称',
    role_desc TEXT COMMENT '角色描述',
    role_type VARCHAR(20) NOT NULL DEFAULT 'custom' COMMENT '角色类型：system, custom',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    INDEX idx_role_name (role_name),
    INDEX idx_role_type (role_type),
) COMMENT '角色表';

-- 3. 权限表
CREATE TABLE IF NOT EXISTS permissions (
    id VARCHAR(36) PRIMARY KEY COMMENT '权限ID，UUID格式',
    permission_name VARCHAR(100) NOT NULL UNIQUE COMMENT '权限名称',
    resource VARCHAR(100) NOT NULL COMMENT '资源',
    actions JSON NOT NULL COMMENT '操作列表，JSON格式',
    scope VARCHAR(20) DEFAULT 'global' COMMENT '权限范围：global, department, personal',
    description TEXT COMMENT '权限描述',
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否激活',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    INDEX idx_permission_name (permission_name),
    INDEX idx_resource (resource),
    INDEX idx_scope (scope),
    INDEX idx_is_active (is_active)
) COMMENT '权限表';

-- 4. 用户角色关联表
CREATE TABLE IF NOT EXISTS user_roles (
    id VARCHAR(36) PRIMARY KEY COMMENT '关联ID，UUID格式',
    user_id VARCHAR(36) NOT NULL COMMENT '用户ID',
    role_id VARCHAR(36) NOT NULL COMMENT '角色ID',
    assigned_by VARCHAR(36) COMMENT '分配人ID',
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '分配时间',
    expires_at TIMESTAMP NULL COMMENT '过期时间',
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否激活',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    UNIQUE KEY uk_user_role (user_id, role_id),
    INDEX idx_user_id (user_id),
    INDEX idx_role_id (role_id),
    INDEX idx_assigned_by (assigned_by),
    INDEX idx_is_active (is_active),
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (assigned_by) REFERENCES users(id) ON DELETE SET NULL
) COMMENT '用户角色关联表';

-- 5. 角色权限关联表
CREATE TABLE IF NOT EXISTS role_permissions (
    id VARCHAR(36) PRIMARY KEY COMMENT '关联ID，UUID格式',
    role_id VARCHAR(36) NOT NULL COMMENT '角色ID',
    permission_id VARCHAR(36) NOT NULL COMMENT '权限ID',
    granted_by VARCHAR(36) COMMENT '授权人ID',
    granted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '授权时间',
    expires_at TIMESTAMP NULL COMMENT '过期时间',
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否激活',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    UNIQUE KEY uk_role_permission (role_id, permission_id),
    INDEX idx_role_id (role_id),
    INDEX idx_permission_id (permission_id),
    INDEX idx_granted_by (granted_by),
    INDEX idx_is_active (is_active),
    
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE,
    FOREIGN KEY (granted_by) REFERENCES users(id) ON DELETE SET NULL
) COMMENT '角色权限关联表';

-- 6. 用户会话表
CREATE TABLE IF NOT EXISTS user_sessions (
    id VARCHAR(36) PRIMARY KEY COMMENT '会话ID，UUID格式',
    user_id VARCHAR(36) NOT NULL COMMENT '用户ID',
    session_id VARCHAR(64) NOT NULL UNIQUE COMMENT '会话标识',
    ip_address VARCHAR(45) COMMENT 'IP地址',
    user_agent TEXT COMMENT '用户代理',
    status VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '会话状态：active, expired, revoked',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    last_activity TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后活动时间',
    expires_at TIMESTAMP NOT NULL COMMENT '过期时间',
    
    INDEX idx_user_id (user_id),
    INDEX idx_session_id (session_id),
    INDEX idx_status (status),
    INDEX idx_expires_at (expires_at),
    INDEX idx_last_activity (last_activity),
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) COMMENT '用户会话表';

-- 7. 用户活动表
CREATE TABLE IF NOT EXISTS user_activities (
    id VARCHAR(36) PRIMARY KEY COMMENT '活动ID，UUID格式',
    user_id VARCHAR(36) NOT NULL COMMENT '用户ID',
    session_id VARCHAR(64) COMMENT '会话ID',
    action_type VARCHAR(50) NOT NULL COMMENT '操作类型',
    action_detail TEXT COMMENT '操作详情',
    resource_type VARCHAR(50) COMMENT '资源类型',
    resource_id VARCHAR(36) COMMENT '资源ID',
    ip_address VARCHAR(45) COMMENT 'IP地址',
    result VARCHAR(20) DEFAULT 'success' COMMENT '操作结果：success, failed',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    
    INDEX idx_user_id (user_id),
    INDEX idx_session_id (session_id),
    INDEX idx_action_type (action_type),
    INDEX idx_resource_type (resource_type),
    INDEX idx_created_at (created_at),
    INDEX idx_result (result),
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) COMMENT '用户活动表';

-- 8. 用户安全日志表
CREATE TABLE IF NOT EXISTS user_security_logs (
    id VARCHAR(36) PRIMARY KEY COMMENT '日志ID，UUID格式',
    user_id VARCHAR(36) COMMENT '用户ID，可为空（系统级日志）',
    event_type VARCHAR(50) NOT NULL COMMENT '事件类型',
    event_detail TEXT COMMENT '事件详情',
    level VARCHAR(20) NOT NULL DEFAULT 'info' COMMENT '日志级别：info, warning, error, critical',
    ip_address VARCHAR(45) COMMENT 'IP地址',
    risk_level VARCHAR(20) DEFAULT 'low' COMMENT '风险等级：low, medium, high, critical',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    
    INDEX idx_user_id (user_id),
    INDEX idx_event_type (event_type),
    INDEX idx_level (level),
    INDEX idx_risk_level (risk_level),
    INDEX idx_created_at (created_at),
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
) COMMENT '用户安全日志表';

-- 9. 用户安全设置表
CREATE TABLE IF NOT EXISTS user_security_settings (
    id VARCHAR(36) PRIMARY KEY COMMENT '设置ID，UUID格式',
    user_id VARCHAR(36) NOT NULL UNIQUE COMMENT '用户ID',
    enable_two_factor BOOLEAN DEFAULT FALSE COMMENT '是否启用双因素认证',
    session_timeout INT DEFAULT 7200 COMMENT '会话超时时间（秒）',
    max_concurrent_sessions INT DEFAULT 5 COMMENT '最大并发会话数',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    INDEX idx_user_id (user_id),
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) COMMENT '用户安全设置表';

-- 10. 用户统计表（用于缓存统计数据）
CREATE TABLE IF NOT EXISTS user_statistics (
    id VARCHAR(36) PRIMARY KEY COMMENT '统计ID，UUID格式',
    stat_date DATE NOT NULL COMMENT '统计日期',
    total_users INT DEFAULT 0 COMMENT '总用户数',
    active_users INT DEFAULT 0 COMMENT '活跃用户数',
    new_users INT DEFAULT 0 COMMENT '新增用户数',
    login_count INT DEFAULT 0 COMMENT '登录次数',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    UNIQUE KEY uk_stat_date (stat_date),
    INDEX idx_stat_date (stat_date)
) COMMENT '用户统计表';

-- 11. 用户设备关联表
CREATE TABLE IF NOT EXISTS user_devices (
    id VARCHAR(36) PRIMARY KEY COMMENT '关联ID，UUID格式',
    user_id VARCHAR(36) NOT NULL COMMENT '用户ID',
    device_id VARCHAR(36) NOT NULL COMMENT '设备ID',
    device_name VARCHAR(100) COMMENT '设备别名',
    device_type VARCHAR(50) COMMENT '设备类型：mobile, desktop, tablet, other',
    access_level VARCHAR(20) NOT NULL DEFAULT 'read' COMMENT '访问级别：read, write, admin',
    last_access_at TIMESTAMP NULL COMMENT '最后访问时间',
    status VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '状态：active, inactive, blocked',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    UNIQUE KEY uk_user_device (user_id, device_id),
    INDEX idx_user_id (user_id),
    INDEX idx_device_id (device_id),
    INDEX idx_device_type (device_type),
    INDEX idx_access_level (access_level),
    INDEX idx_status (status),
    INDEX idx_last_access_at (last_access_at),
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) COMMENT '用户设备关联表';

-- 12. 用户任务关联表
CREATE TABLE IF NOT EXISTS user_tasks (
    id VARCHAR(36) PRIMARY KEY COMMENT '关联ID，UUID格式',
    user_id VARCHAR(36) NOT NULL COMMENT '用户ID',
    task_id VARCHAR(36) NOT NULL COMMENT '任务ID',
    role VARCHAR(50) COMMENT '在任务中的角色：owner, executor, reviewer, watcher',
    priority VARCHAR(20) DEFAULT 'normal' COMMENT '优先级：low, normal, high, urgent',
    status VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '状态：pending, in_progress, completed, cancelled, failed',
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '分配时间',
    started_at TIMESTAMP NULL COMMENT '开始时间',
    completed_at TIMESTAMP NULL COMMENT '完成时间',
    due_date TIMESTAMP NULL COMMENT '截止时间',
    progress_percentage INT DEFAULT 0 COMMENT '进度百分比',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    UNIQUE KEY uk_user_task (user_id, task_id),
    INDEX idx_user_id (user_id),
    INDEX idx_task_id (task_id),
    INDEX idx_role (role),
    INDEX idx_priority (priority),
    INDEX idx_status (status),
    INDEX idx_due_date (due_date),
    INDEX idx_assigned_at (assigned_at),
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) COMMENT '用户任务关联表';

-- 13. 用户设备操作日志表
CREATE TABLE IF NOT EXISTS user_device_logs (
    id VARCHAR(36) PRIMARY KEY COMMENT '日志ID，UUID格式',
    user_id VARCHAR(36) NOT NULL COMMENT '用户ID',
    device_id VARCHAR(36) NOT NULL COMMENT '设备ID',
    operation_type VARCHAR(50) NOT NULL COMMENT '操作类型：connect, disconnect, command, monitor, configure',
    operation_detail TEXT COMMENT '操作详情',
    result VARCHAR(20) DEFAULT 'success' COMMENT '操作结果：success, failed, timeout',
    ip_address VARCHAR(45) COMMENT 'IP地址',
    session_id VARCHAR(64) COMMENT '会话ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    
    INDEX idx_user_id (user_id),
    INDEX idx_device_id (device_id),
    INDEX idx_operation_type (operation_type),
    INDEX idx_result (result),
    INDEX idx_created_at (created_at),
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) COMMENT '用户设备操作日志表';

-- 14. 用户任务操作日志表
CREATE TABLE IF NOT EXISTS user_task_logs (
    id VARCHAR(36) PRIMARY KEY COMMENT '日志ID，UUID格式',
    user_id VARCHAR(36) NOT NULL COMMENT '用户ID',
    task_id VARCHAR(36) NOT NULL COMMENT '任务ID',
    operation_type VARCHAR(50) NOT NULL COMMENT '操作类型：assign, start, update, complete, cancel, comment',
    operation_detail TEXT COMMENT '操作详情',
    old_status VARCHAR(20) COMMENT '原状态',
    new_status VARCHAR(20) COMMENT '新状态',
    progress_change INT DEFAULT 0 COMMENT '进度变化',
    ip_address VARCHAR(45) COMMENT 'IP地址',
    session_id VARCHAR(64) COMMENT '会话ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    
    INDEX idx_user_id (user_id),
    INDEX idx_task_id (task_id),
    INDEX idx_operation_type (operation_type),
    INDEX idx_created_at (created_at),
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) COMMENT '用户任务操作日志表';

-- =================================================================================
-- 初始化数据
-- =================================================================================

-- 插入默认角色
INSERT INTO roles (id, role_name, role_desc, role_type) VALUES
('role-admin', 'admin', '系统管理员', 'system'),
('role-user', 'user', '普通用户', 'system'),
('role-guest', 'guest', '访客用户', 'system');

-- 插入默认权限
INSERT INTO permissions (id, permission_name, resource, actions, scope, description) VALUES
('perm-user-read', 'user:read', 'user', '["read"]', 'global', '读取用户信息'),
('perm-user-write', 'user:write', 'user', '["create", "update"]', 'global', '创建和更新用户'),
('perm-user-delete', 'user:delete', 'user', '["delete"]', 'global', '删除用户'),
('perm-role-manage', 'role:manage', 'role', '["read", "create", "update", "delete"]', 'global', '角色管理'),
('perm-permission-manage', 'permission:manage', 'permission', '["read", "create", "update", "delete"]', 'global', '权限管理'),
('perm-system-admin', 'system:admin', 'system', '["*"]', 'global', '系统管理');

-- 为管理员角色分配所有权限
INSERT INTO role_permissions (id, role_id, permission_id) VALUES
('rp-admin-1', 'role-admin', 'perm-user-read'),
('rp-admin-2', 'role-admin', 'perm-user-write'),
('rp-admin-3', 'role-admin', 'perm-user-delete'),
('rp-admin-4', 'role-admin', 'perm-role-manage'),
('rp-admin-5', 'role-admin', 'perm-permission-manage'),
('rp-admin-6', 'role-admin', 'perm-system-admin');

-- 为普通用户角色分配基本权限
INSERT INTO role_permissions (id, role_id, permission_id) VALUES
('rp-user-1', 'role-user', 'perm-user-read');

-- 为访客角色分配只读权限
INSERT INTO role_permissions (id, role_id, permission_id) VALUES
('rp-guest-1', 'role-guest', 'perm-user-read');

-- 创建默认管理员用户（密码需要在实际使用时修改）
INSERT INTO users (id, username, email, password, status, nickname, department, position, is_verified) VALUES
('user-admin', 'admin', 'admin@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'active', '系统管理员', 'IT部门', '系统管理员', TRUE);

-- 为默认管理员分配管理员角色
INSERT INTO user_roles (id, user_id, role_id, assigned_by) VALUES
('ur-admin-1', 'user-admin', 'role-admin', 'user-admin');

-- 为默认管理员创建安全设置
INSERT INTO user_security_settings (id, user_id, enable_two_factor, session_timeout, max_concurrent_sessions) VALUES
('us-admin-1', 'user-admin', FALSE, 7200, 10); 