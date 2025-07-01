-- 创建设备主表
CREATE TABLE IF NOT EXISTS devices (
    id              BIGSERIAL PRIMARY KEY,
    device_id       VARCHAR(100) NOT NULL UNIQUE,
    name            VARCHAR(255) NOT NULL,
    type            VARCHAR(50) NOT NULL,
    model           VARCHAR(100),
    manufacturer    VARCHAR(100),
    
    -- 连接信息
    ip_address      INET,
    port            INTEGER,
    protocol        VARCHAR(20),
    endpoint        VARCHAR(500),
    
    -- 认证信息
    auth_type       VARCHAR(20),
    credentials     JSONB, -- 加密存储
    
    -- 状态信息
    status          VARCHAR(20) NOT NULL DEFAULT 'offline',
    health_score    INTEGER DEFAULT 100, -- 健康度(0-100)
    last_seen       TIMESTAMP,
    
    -- 配置信息
    config          JSONB,
    capabilities    JSONB,
    metadata        JSONB,
    
    -- 统计信息
    total_commands  BIGINT DEFAULT 0,
    success_commands BIGINT DEFAULT 0,
    failed_commands BIGINT DEFAULT 0,
    uptime_hours    NUMERIC(10,2) DEFAULT 0,
    
    -- 审计字段
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by      BIGINT,
    updated_by      BIGINT
);

-- 创建设备心跳表
CREATE TABLE IF NOT EXISTS device_heartbeats (
    id              BIGSERIAL PRIMARY KEY,
    device_id       BIGINT NOT NULL,
    heartbeat_time  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- 心跳数据
    status          VARCHAR(20) NOT NULL,
    metrics         JSONB,
    system_info     JSONB,
    
    -- 网络信息
    ip_address      INET,
    response_time   INTEGER, -- 响应时间(毫秒)
    
    FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE
);

-- 创建设备日志表
CREATE TABLE IF NOT EXISTS device_logs (
    id              BIGSERIAL PRIMARY KEY,
    device_id       BIGINT NOT NULL,
    log_time        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- 日志分类
    level           VARCHAR(10) NOT NULL, -- DEBUG, INFO, WARN, ERROR
    category        VARCHAR(50), -- system, communication, sensor, etc.
    
    -- 日志内容
    message         TEXT NOT NULL,
    details         JSONB,
    
    -- 上下文信息
    source          VARCHAR(100),
    correlation_id  VARCHAR(100),
    
    FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE
);

-- 创建设备命令表
CREATE TABLE IF NOT EXISTS device_commands (
    id              BIGSERIAL PRIMARY KEY,
    device_id       BIGINT NOT NULL,
    command_id      VARCHAR(100) NOT NULL UNIQUE,
    
    -- 命令信息
    command_type    VARCHAR(50) NOT NULL,
    command_data    JSONB NOT NULL,
    
    -- 执行状态
    status          VARCHAR(20) NOT NULL DEFAULT 'pending',
    sent_time       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    executed_time   TIMESTAMP,
    completed_time  TIMESTAMP,
    
    -- 结果信息
    response_data   JSONB,
    error_message   TEXT,
    
    -- 审计信息
    created_by      BIGINT,
    
    FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE
);

-- 创建设备基础查询索引
CREATE UNIQUE INDEX IF NOT EXISTS idx_devices_device_id ON devices(device_id);
CREATE INDEX IF NOT EXISTS idx_devices_type ON devices(type);
CREATE INDEX IF NOT EXISTS idx_devices_status ON devices(status);
CREATE INDEX IF NOT EXISTS idx_devices_ip_address ON devices(ip_address);
CREATE INDEX IF NOT EXISTS idx_devices_last_seen ON devices(last_seen);

-- 创建设备复合索引
CREATE INDEX IF NOT EXISTS idx_devices_type_status ON devices(type, status);
CREATE INDEX IF NOT EXISTS idx_devices_status_last_seen ON devices(status, last_seen);

-- 创建设备心跳索引
CREATE INDEX IF NOT EXISTS idx_device_heartbeats_device_id ON device_heartbeats(device_id);
CREATE INDEX IF NOT EXISTS idx_device_heartbeats_heartbeat_time ON device_heartbeats(heartbeat_time);

-- 创建设备日志索引
CREATE INDEX IF NOT EXISTS idx_device_logs_device_id ON device_logs(device_id);
CREATE INDEX IF NOT EXISTS idx_device_logs_level ON device_logs(level);
CREATE INDEX IF NOT EXISTS idx_device_logs_log_time ON device_logs(log_time);
CREATE INDEX IF NOT EXISTS idx_device_logs_category ON device_logs(category);

-- 创建设备命令索引
CREATE INDEX IF NOT EXISTS idx_device_commands_device_id ON device_commands(device_id);
CREATE INDEX IF NOT EXISTS idx_device_commands_status ON device_commands(status);
CREATE UNIQUE INDEX IF NOT EXISTS idx_device_commands_command_id ON device_commands(command_id);

-- 为devices表创建自动更新触发器
CREATE TRIGGER update_devices_updated_at 
    BEFORE UPDATE ON devices 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- 添加任务表与设备表的外键关联
ALTER TABLE tasks 
ADD CONSTRAINT fk_tasks_device_id 
FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE SET NULL; 