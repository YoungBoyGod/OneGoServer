-- 创建设备主表
CREATE TABLE IF NOT EXISTS devices (
    -- 设备信息
    id              BIGSERIAL PRIMARY KEY,   -- 设备ID
    device_id      VARCHAR(100) NOT NULL UNIQUE, -- 设备编号
    name            VARCHAR(255) NOT NULL, -- 设备名称 
    type            VARCHAR(50) NOT NULL, -- 设备类型 （如：普通pc 服务器）
    model           VARCHAR(100), -- 设备型号 （如：host，soc）
    board_id        VARCHAR(100), -- 板卡ESN （如：1234567890）

   -- 状态信息
    status          VARCHAR(20) NOT NULL DEFAULT 'unkown', -- 设备状态 （如：online offline maintenance 从未上线过为unkown）
    health_score    INTEGER DEFAULT 100, -- 健康度(0-100) 

    --登陆信息
    login_username        VARCHAR(100), -- 用户名
    login_port      INTEGER, -- 登录端口
    login_public_key      VARCHAR(100), -- 公钥

    -- 连接信息
    ip_address      INET,   -- 设备IP地址
    port            INTEGER, -- 设备端口
    protocol        VARCHAR(20), -- 设备协议 ssh http
    endpoint        VARCHAR(500), -- 设备端点 （如：/api/v1/devices）
    

    -- 认证信息
    reg_time        TIMESTAMP, -- 注册时间

    -- 配置信息
    metadata        JSONB, -- 设备元数据 （如：{"cpu": "4", "memory": "8", "disk": "100"}）
    tags            JSONB, -- 设备标签 （如：{"tag1": "value1", "tag2": "value2"}）
   
    -- 统计信息
    uptime_hours    NUMERIC(10,2) DEFAULT 0, -- 在线时长(h) 计算方式：每次下线时间减去上线时间
    first_online_time TIMESTAMP, -- 首次上线时间
    last_online_time TIMESTAMP, -- 最后上线时间
    last_offline_time TIMESTAMP, -- 最后下线时间

    total_online_duration BIGINT DEFAULT 0, -- 在线时长(h) 总在线时长 计算方式：每次下线时间减去上线时间
    total_offline_duration BIGINT DEFAULT 0, -- 离线时长(h) 总离线时长 计算方式：每次上线时间减去下线时间

    total_heartbeats BIGINT DEFAULT 0, -- 总心跳数
    total_alerts BIGINT DEFAULT 0, -- 总告警数

    -- 任务信息
    total_tasks BIGINT DEFAULT 0, -- 总任务数
    total_success_tasks BIGINT DEFAULT 0, -- 总成功任务数
    total_failed_tasks BIGINT DEFAULT 0, -- 总失败任务数
    total_canceled_tasks BIGINT DEFAULT 0, -- 总取消任务数
    total_pending_tasks BIGINT DEFAULT 0, -- 总待执行任务数
    total_running_tasks BIGINT DEFAULT 0, -- 总执行中任务数
    total_completed_tasks BIGINT DEFAULT 0, -- 总完成任务数


    -- 审计字段
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by      VARCHAR(100) NOT NULL DEFAULT 'system', -- 创建者
    updated_by      VARCHAR(100) NOT NULL DEFAULT 'system' -- 更新者
);


-- 创建设备心跳表
CREATE TABLE IF NOT EXISTS device_heartbeats (
    
    id              BIGSERIAL PRIMARY KEY, -- 心跳ID
    device_id      VARCHAR(100) NOT NULL, -- 设备编号
    heartbeat_time  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 心跳时间
    
    -- 心跳数据
    status          VARCHAR(20) NOT NULL, -- 心跳状态 （如：online offline maintenance 从未上线过为unkown）
    metadata         JSONB,  -- 心跳数据 （如：{"cpu": "4", "memory": "8", "disk": "100"} ） 心跳数据为设备上报的指标数据，如cpu使用率，内存使用率，磁盘使用率等

    -- 网络信息
    ip_address      INET, -- 设备IP地址
    response_time   INTEGER -- 响应时间(毫秒) 心跳响应时间
    
    -- 注意：移除外键约束，改为应用层维护数据一致性
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
    correlation_id  VARCHAR(100)
    
    -- 注意：移除外键约束，改为应用层维护数据一致性
);

-- 创建设备命令表
CREATE TABLE IF NOT EXISTS device_commands (
    id              BIGSERIAL PRIMARY KEY, -- 命令ID
    device_id       BIGINT NOT NULL, -- 设备编号
    command_id      VARCHAR(100) NOT NULL UNIQUE, -- 命令ID
    
    -- 命令信息
    command_type    VARCHAR(50) NOT NULL, -- 命令类型 （如：shell 文件传输 数据库操作）
    command_data    JSONB NOT NULL, -- 命令数据 （如：{"command": "ls -l", "args": ["-l"]}）
    
    -- 执行状态
    status          VARCHAR(20) NOT NULL DEFAULT 'pending', -- 命令状态 （如：pending running completed failed canceled）
    sent_time       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 发送时间
    executed_time   TIMESTAMP, -- 执行时间
    completed_time  TIMESTAMP, -- 完成时间
    
    -- 结果信息
    response_data   JSONB, -- 命令执行结果 （如：{"output": "ls -l", "exit_code": 0}）
    error_message   TEXT, -- 错误信息 （如："command not found"）
    
    -- 审计信息
    created_by      BIGINT -- 创建者
    
    -- 注意：移除外键约束，改为应用层维护数据一致性
);

-- 创建设备任务表
CREATE TABLE IF NOT EXISTS device_tasks (
    id              BIGSERIAL PRIMARY KEY, -- 任务ID
    device_id       BIGINT NOT NULL, -- 设备编号
    task_id      VARCHAR(100) NOT NULL UNIQUE, -- 任务ID

    
    
);

-- 创建设备基础查询索引
CREATE UNIQUE INDEX IF NOT EXISTS idx_devices_device_id ON devices(device_id);
CREATE INDEX IF NOT EXISTS idx_devices_type ON devices(type);
CREATE INDEX IF NOT EXISTS idx_devices_status ON devices(status);
CREATE INDEX IF NOT EXISTS idx_devices_ip_address ON devices(ip_address);
CREATE INDEX IF NOT EXISTS idx_devices_last_online_time ON devices(last_online_time);

-- 创建设备复合索引
CREATE INDEX IF NOT EXISTS idx_devices_type_status ON devices(type, status);
CREATE INDEX IF NOT EXISTS idx_devices_status_last_online ON devices(status, last_online_time);

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

-- 注意：移除外键约束，任务表与设备表的关联改为应用层维护数据一致性
-- ALTER TABLE tasks 
-- ADD CONSTRAINT fk_tasks_device_id 
-- FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE SET NULL; 