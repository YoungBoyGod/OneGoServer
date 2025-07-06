-- 创建触发器函数：用于自动更新时间戳
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

-- 1. 设备主表 (无依赖)
CREATE TABLE IF NOT EXISTS devices (
    id              BIGSERIAL PRIMARY KEY,                    -- 设备内部唯一标识符，自增主键
    device_id       VARCHAR(100) NOT NULL UNIQUE,             -- 设备业务ID，外部系统使用的设备标识
    name            VARCHAR(255) NOT NULL,                    -- 设备名称，用于显示和识别
    type            VARCHAR(50) NOT NULL,                     -- 设备类型，如：服务器、路由器、交换机等
    model           VARCHAR(100),                             -- 设备型号，制造商提供的具体型号
    board_id        VARCHAR(100),                             -- 主板ID，硬件层面的唯一标识
    status          VARCHAR(20) NOT NULL DEFAULT 'unknown',   -- 设备状态：online/offline/unknown/maintenance
    health_score    INTEGER DEFAULT 100,                      -- 设备健康评分，0-100，100为最佳状态
    login_username  VARCHAR(100),                             -- 登录用户名，用于设备访问认证
    login_port      INTEGER,                                  -- 登录端口，SSH/Telnet等服务的端口号
    login_public_key VARCHAR(100),                            -- 登录公钥，用于SSH密钥认证
    ip_address      INET,                                     -- 设备IP地址，支持IPv4和IPv6
    port            INTEGER,                                  -- 设备服务端口，主要服务监听的端口
    protocol        VARCHAR(20),                              -- 通信协议，如：SSH/Telnet/HTTP/HTTPS
    endpoint        VARCHAR(500),                             -- 设备访问端点，完整的访问URL
    reg_time        TIMESTAMP,                                -- 设备注册时间，首次添加到系统的时间
    metadata        JSONB,                                    -- 设备元数据，存储额外的设备信息
    tags            JSONB,                                    -- 设备标签，用于分类和筛选
    uptime_hours    NUMERIC(10,2) DEFAULT 0,                  -- 设备运行时长（小时），累计在线时间
    first_online_time TIMESTAMP,                              -- 首次上线时间，设备第一次连接的时间
    last_online_time TIMESTAMP,                               -- 最后上线时间，最近一次连接的时间
    last_offline_time TIMESTAMP,                              -- 最后下线时间，最近一次断开的时间
    total_online_duration BIGINT DEFAULT 0,                   -- 总在线时长（秒），累计在线时间
    total_offline_duration BIGINT DEFAULT 0,                  -- 总离线时长（秒），累计离线时间
    total_heartbeats BIGINT DEFAULT 0,                        -- 总心跳次数，设备发送心跳的总数
    total_alerts    BIGINT DEFAULT 0,                         -- 总告警次数，设备产生的告警总数
    total_tasks     BIGINT DEFAULT 0,                         -- 总任务数，分配给该设备的所有任务
    total_success_tasks BIGINT DEFAULT 0,                     -- 成功任务数，设备成功完成的任务数
    total_failed_tasks BIGINT DEFAULT 0,                      -- 失败任务数，设备执行失败的任务数
    total_canceled_tasks BIGINT DEFAULT 0,                    -- 取消任务数，被取消的任务数
    total_pending_tasks BIGINT DEFAULT 0,                     -- 待执行任务数，等待执行的任务数
    total_running_tasks BIGINT DEFAULT 0,                     -- 运行中任务数，正在执行的任务数
    total_completed_tasks BIGINT DEFAULT 0,                   -- 已完成任务数，已完成的任务总数
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 记录创建时间
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 记录更新时间
    created_by      VARCHAR(100) NOT NULL DEFAULT 'system',   -- 创建者，记录创建人
    updated_by      VARCHAR(100) NOT NULL DEFAULT 'system'    -- 更新者，记录最后修改人
);

-- 2. 任务主表 (依赖 devices)
CREATE TABLE IF NOT EXISTS tasks (
    id              BIGSERIAL PRIMARY KEY,                    -- 任务内部唯一标识符，自增主键
    device_id       varchar(100),                             -- 设备业务ID，关联到devices表的device_id
    task_id         VARCHAR(100) NOT NULL UNIQUE,             -- 任务业务ID，外部系统使用的任务标识
    name            VARCHAR(255) NOT NULL,                    -- 任务名称，用于显示和识别
    description     TEXT,                                     -- 任务描述，详细的任务说明
    type            VARCHAR(50) NOT NULL,                     -- 任务类型，如：backup/sync/monitor/custom
    status          VARCHAR(20) NOT NULL DEFAULT 'pending',   -- 任务状态：pending/running/completed/failed/canceled
    priority        INTEGER NOT NULL DEFAULT 5,               -- 任务优先级，1-10，数字越大优先级越高
    execute_time    TIMESTAMP,                                -- 计划执行时间，任务计划开始执行的时间
    timeout         INTEGER DEFAULT 86400,                    -- 超时时间（秒），任务执行的最大允许时间
    retry_count     INTEGER DEFAULT 0,                        -- 重试次数，当前已重试的次数
    max_retries     INTEGER DEFAULT 3,                        -- 最大重试次数，允许的最大重试次数
    is_urgent       BOOLEAN DEFAULT FALSE,                    -- 是否紧急任务，紧急任务优先执行
    parameters      JSONB,                                    -- 任务参数，执行任务所需的参数
    result          JSONB,                                    -- 任务结果，任务执行的结果数据
    error_message   TEXT,                                     -- 错误信息，任务失败时的错误描述
    executor_type   VARCHAR(50),                              -- 执行器类型，如：local/remote/container
    executor_id     VARCHAR(100),                             -- 执行器ID，具体执行器的标识
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 记录创建时间
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 记录更新时间
    created_by      VARCHAR(100),                             -- 创建者，任务创建人
    updated_by      VARCHAR(100),                             -- 更新者，任务最后修改人
    CONSTRAINT fk_tasks_device
      FOREIGN KEY (device_id)
      REFERENCES devices(device_id)
      ON DELETE SET NULL
);

-- 3. 任务执行记录表 (依赖 tasks)
CREATE TABLE IF NOT EXISTS task_executions (
    id              BIGSERIAL PRIMARY KEY,                    -- 执行记录内部唯一标识符，自增主键
    task_id         VARCHAR(100) NOT NULL,                    -- 任务业务ID，关联到tasks表的task_id
    execution_id    VARCHAR(100) NOT NULL UNIQUE,             -- 执行记录业务ID，本次执行的唯一标识
    device_id      VARCHAR(100),                              -- 设备业务ID，执行任务的设备
    status          VARCHAR(20) NOT NULL DEFAULT 'started',   -- 执行状态：started/running/completed/failed/canceled
    start_time      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 开始时间，任务开始执行的时间
    end_time        TIMESTAMP,                                -- 结束时间，任务执行完成的时间
    duration        INTEGER,                                  -- 执行时长（秒），任务实际执行的时间
    executor_info   JSONB,                                    -- 执行器信息，执行器的详细信息
    logs            TEXT,                                     -- 执行日志，任务执行过程中的日志
    metrics         JSONB,                                    -- 执行指标，性能指标和监控数据
    output          JSONB,                                    -- 执行输出，任务执行的结果输出
    error_details   JSONB,                                    -- 错误详情，详细的错误信息和堆栈
    cpu_usage_avg   NUMERIC(5,2),                             -- CPU使用率平均值，执行期间的平均CPU使用率
    cpu_usage_peak  NUMERIC(5,2),                             -- CPU使用率峰值，执行期间的最高CPU使用率
    memory_usage_avg NUMERIC(10,2),                           -- 内存使用率平均值，执行期间的平均内存使用率
    memory_usage_peak NUMERIC(10,2),                          -- 内存使用率峰值，执行期间的最高内存使用率
    io_operations_total BIGINT,                               -- IO操作总数，执行期间的IO操作次数
    io_bytes_total  BIGINT,                                   -- IO字节总数，执行期间的IO数据传输量
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 记录创建时间
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 记录更新时间
    CONSTRAINT fk_task_executions_task FOREIGN KEY (task_id) REFERENCES tasks(task_id) ON DELETE CASCADE
);

-- 4. 设备心跳表 (依赖 devices)
CREATE TABLE IF NOT EXISTS device_heartbeats (
    id              BIGSERIAL PRIMARY KEY,                    -- 心跳记录内部唯一标识符，自增主键
    device_id       VARCHAR(100) NOT NULL,                    -- 设备业务ID，关联到devices表的device_id
    heartbeat_time  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 心跳时间，设备发送心跳的时间
    status          VARCHAR(20) NOT NULL,                     -- 心跳状态：online/offline/error
    metadata        JSONB,                                    -- 心跳元数据，设备状态和性能信息
    ip_address      INET,                                     -- 心跳IP地址，设备当前使用的IP地址
    response_time   INTEGER,                                  -- 响应时间（毫秒），心跳响应的延迟时间
    CONSTRAINT fk_device_heartbeats_device FOREIGN KEY (device_id) REFERENCES devices(device_id) ON DELETE CASCADE
);

-- 5. 设备日志表 (依赖 devices)
CREATE TABLE IF NOT EXISTS device_logs (
    id              BIGSERIAL PRIMARY KEY,                    -- 日志记录内部唯一标识符，自增主键
    device_id       VARCHAR(100) NOT NULL,                    -- 设备业务ID，关联到devices表的device_id
    log_time        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 日志时间，日志产生的时间
    level           VARCHAR(10) NOT NULL,                     -- 日志级别：DEBUG/INFO/WARN/ERROR/FATAL
    category        VARCHAR(50),                              -- 日志分类，如：system/application/security
    message         TEXT NOT NULL,                            -- 日志消息，具体的日志内容
    details         JSONB,                                    -- 日志详情，额外的日志信息
    source          VARCHAR(100),                             -- 日志来源，产生日志的组件或模块
    correlation_id  VARCHAR(100),                             -- 关联ID，用于关联相关日志记录
    CONSTRAINT fk_device_logs_device FOREIGN KEY (device_id) REFERENCES devices(device_id) ON DELETE CASCADE
);

-- 6. 设备命令表 (依赖 devices)
CREATE TABLE IF NOT EXISTS device_commands (
    id              BIGSERIAL PRIMARY KEY,                    -- 命令记录内部唯一标识符，自增主键
    device_id       VARCHAR(100) NOT NULL,                    -- 设备业务ID，关联到devices表的device_id
    command_type    VARCHAR(50) NOT NULL,                     -- 命令类型，如：restart/shutdown/update/config
    command_data    JSONB NOT NULL,                           -- 命令数据，命令的具体参数和配置
    status          VARCHAR(20) NOT NULL DEFAULT 'pending',   -- 命令状态：pending/sent/executed/completed/failed
    sent_time       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 发送时间，命令发送的时间
    executed_time   TIMESTAMP,                                -- 执行时间，命令开始执行的时间
    completed_time  TIMESTAMP,                                -- 完成时间，命令执行完成的时间
    response_data   JSONB,                                    -- 响应数据，设备返回的响应信息
    error_message   TEXT,                                     -- 错误信息，命令执行失败的错误描述
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 记录创建时间
    created_by      VARCHAR(100),                             -- 创建者，命令发送人
    CONSTRAINT fk_device_commands_device FOREIGN KEY (device_id) REFERENCES devices(device_id) ON DELETE CASCADE
);

-- 7. 设备负载监控表 (依赖 devices)
CREATE TABLE IF NOT EXISTS device_load_monitor (
    id              BIGSERIAL PRIMARY KEY,                    -- 监控记录内部唯一标识符，自增主键
    device_id       VARCHAR(100) NOT NULL,                    -- 设备业务ID，关联到devices表的device_id
    current_tasks   INTEGER DEFAULT 0,                        -- 当前任务数，设备正在执行的任务数量
    max_concurrent_tasks INTEGER DEFAULT 1,                   -- 最大并发任务数，设备能同时处理的最大任务数
    cpu_load        NUMERIC(5,2),                             -- CPU负载，当前CPU使用率百分比
    memory_usage    NUMERIC(5,2),                             -- 内存使用率，当前内存使用率百分比
    disk_usage      NUMERIC(5,2),                             -- 磁盘使用率，当前磁盘使用率百分比
    network_latency INTEGER,                                  -- 网络延迟（毫秒），设备网络响应时间
    status          VARCHAR(20) NOT NULL DEFAULT 'online',    -- 监控状态：online/offline/overloaded
    last_heartbeat  TIMESTAMP,                                -- 最后心跳时间，最近一次心跳的时间
    load_score      NUMERIC(5,2),                             -- 负载评分，综合负载评估分数
    total_assigned  INTEGER DEFAULT 0,                        -- 总分配任务数，历史分配的任务总数
    total_completed INTEGER DEFAULT 0,                        -- 总完成任务数，历史完成的任务总数
    total_failed    INTEGER DEFAULT 0,                        -- 总失败任务数，历史失败的任务总数
    success_rate    NUMERIC(5,2),                             -- 成功率，任务执行成功率百分比
    avg_task_duration NUMERIC(10,2),                          -- 平均任务时长（秒），任务执行的平均时间
    last_task_completion TIMESTAMP,                           -- 最后任务完成时间，最近一次任务完成的时间
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 记录更新时间
    CONSTRAINT fk_device_load_monitor_device FOREIGN KEY (device_id) REFERENCES devices(device_id) ON DELETE CASCADE
);

-- 8. 设备任务表 (依赖 devices 和 tasks)
CREATE TABLE IF NOT EXISTS device_tasks (
    id              BIGSERIAL PRIMARY KEY,                    -- 设备任务关联记录内部唯一标识符，自增主键
    device_id       VARCHAR(100) NOT NULL,                    -- 设备业务ID，关联到devices表的device_id
    task_id         VARCHAR(100) NOT NULL UNIQUE,             -- 任务业务ID，关联到tasks表的task_id
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 记录创建时间
    created_by      VARCHAR(100),                             -- 创建者，关联关系创建人
    CONSTRAINT fk_device_tasks_device FOREIGN KEY (device_id) REFERENCES devices(device_id) ON DELETE CASCADE,
    CONSTRAINT fk_device_tasks_task FOREIGN KEY (task_id) REFERENCES tasks(task_id) ON DELETE CASCADE
);

-- 9. 任务分配队列表 (依赖 tasks 和 devices)
CREATE TABLE IF NOT EXISTS task_assignment_queue (
    id              BIGSERIAL PRIMARY KEY,                    -- 分配记录内部唯一标识符，自增主键
    task_id         VARCHAR(100) NOT NULL UNIQUE,             -- 任务业务ID，关联到tasks表的task_id
    priority        INTEGER NOT NULL DEFAULT 5,               -- 队列优先级，1-10，数字越大优先级越高
    queue_status    VARCHAR(20) NOT NULL DEFAULT 'queued',    -- 队列状态：queued/assigning/assigned/failed/canceled
    required_device_type VARCHAR(50),                         -- 要求设备类型，任务需要的设备类型
    required_capabilities JSONB,                              -- 要求能力，任务需要的设备能力
    preferred_device_ids BIGINT[],                            -- 偏好设备ID列表，优先选择的设备ID
    excluded_device_ids BIGINT[],                             -- 排除设备ID列表，不能选择的设备ID
    assignment_strategy VARCHAR(50) DEFAULT 'load_balance',   -- 分配策略：load_balance/round_robin/priority
    affinity_rules  JSONB,                                    -- 亲和性规则，设备选择的亲和性配置
    assigned_device_id varchar(100),                          -- 分配的设备ID，实际分配的设备业务ID
    assigned_at     TIMESTAMP,                                -- 分配时间，任务分配给设备的时间
    assignment_score NUMERIC(5,2),                            -- 分配评分，设备匹配度的评分
    queue_position  INTEGER,                                  -- 队列位置，任务在队列中的位置
    estimated_wait_time INTEGER,                              -- 预估等待时间（秒），预计等待执行的时间
    retry_count     INTEGER DEFAULT 0,                        -- 重试次数，分配失败后的重试次数
    max_retries     INTEGER DEFAULT 3,                        -- 最大重试次数，允许的最大重试次数
    original_priority INTEGER,                                -- 原始优先级，任务创建时的初始优先级
    last_priority_change_at TIMESTAMP,                        -- 最后优先级变更时间，最近一次优先级调整的时间
    priority_change_reason VARCHAR(255),                      -- 优先级变更原因，调整优先级的理由
    priority_boost_reason VARCHAR(100),                       -- 优先级提升原因，提升优先级的具体原因
    operation_source VARCHAR(50) DEFAULT 'system',            -- 操作来源：system/manual/api
    queued_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 入队时间，任务进入队列的时间
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 记录更新时间
    CONSTRAINT fk_task_assignment_queue_task FOREIGN KEY (task_id) REFERENCES tasks(task_id) ON DELETE CASCADE,
    CONSTRAINT fk_task_assignment_queue_device FOREIGN KEY (assigned_device_id) REFERENCES devices(device_id) ON DELETE SET NULL
);

-- 10. 任务分配历史表 (依赖 tasks 和 devices)
CREATE TABLE IF NOT EXISTS task_assignment_history (
    id              BIGSERIAL PRIMARY KEY,                    -- 历史记录内部唯一标识符，自增主键
    task_id         VARCHAR(100) NOT NULL,                    -- 任务业务ID，关联到tasks表的task_id
    device_id       VARCHAR(100),                             -- 设备业务ID，关联到devices表的device_id
    action          VARCHAR(20) NOT NULL,                     -- 操作动作：queued/assigned/reassigned/failed/completed/canceled
    previous_status VARCHAR(20),                              -- 之前状态，操作前的任务状态
    new_status      VARCHAR(20),                              -- 新状态，操作后的任务状态
    reason          VARCHAR(255),                             -- 操作原因，执行操作的理由
    details         JSONB,                                    -- 操作详情，操作的详细信息
    operation_source VARCHAR(50) DEFAULT 'system',            -- 操作来源：system/manual/api
    operator_id     VARCHAR(100),                             -- 操作者ID，执行操作的用户或系统
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 记录创建时间
    CONSTRAINT fk_task_assignment_history_task FOREIGN KEY (task_id) REFERENCES tasks(task_id) ON DELETE CASCADE,
    CONSTRAINT fk_task_assignment_history_device FOREIGN KEY (device_id) REFERENCES devices(device_id) ON DELETE SET NULL
);

-- 11. 设备任务队列表 (依赖 devices 和 tasks)
CREATE TABLE IF NOT EXISTS device_task_queue (
    id              BIGSERIAL PRIMARY KEY,                    -- 队列记录内部唯一标识符，自增主键
    device_id       VARCHAR(100) NOT NULL,                    -- 设备业务ID，关联到devices表的device_id
    task_id         VARCHAR(100) NOT NULL,                    -- 任务业务ID，关联到tasks表的task_id
    queue_priority  INTEGER NOT NULL DEFAULT 5,               -- 队列优先级，1-10，数字越大优先级越高
    original_priority INTEGER,                                -- 原始优先级，任务创建时的初始优先级
    is_manual_priority BOOLEAN NOT NULL DEFAULT FALSE,        -- 是否手动优先级，是否由用户手动调整的优先级
    queue_position  INTEGER,                                  -- 队列位置，任务在设备队列中的位置
    is_manual_position BOOLEAN NOT NULL DEFAULT FALSE,        -- 是否手动位置，是否由用户手动调整的位置
    status          VARCHAR(20) NOT NULL DEFAULT 'queued',    -- 队列状态：queued/running/completed/failed/canceled
    estimated_start_time TIMESTAMP,                           -- 预估开始时间，预计开始执行的时间
    estimated_duration INTEGER,                               -- 预估时长（秒），预计执行需要的时间
    actual_start_time TIMESTAMP,                              -- 实际开始时间，任务实际开始执行的时间
    actual_end_time TIMESTAMP,                                -- 实际结束时间，任务实际完成的时间
    max_retry_count INTEGER NOT NULL DEFAULT 3,               -- 最大重试次数，允许的最大重试次数
    current_retry   INTEGER NOT NULL DEFAULT 0,               -- 当前重试次数，已经重试的次数
    timeout_seconds INTEGER NOT NULL DEFAULT 3600,            -- 超时时间（秒），任务执行的最大允许时间
    depends_on_task_ids VARCHAR(100)[],                       -- 依赖任务ID列表，需要先完成的任务ID
    blocks_task_ids VARCHAR(100)[],                           -- 阻塞任务ID列表，被当前任务阻塞的任务ID
    requeue_count   INTEGER NOT NULL DEFAULT 0,               -- 重新入队次数，任务重新进入队列的次数
    last_requeue_at TIMESTAMP,                                -- 最后重新入队时间，最近一次重新入队的时间
    is_requeued     BOOLEAN NOT NULL DEFAULT FALSE,           -- 是否重新入队，任务是否被重新放入队列
    cancel_reason   TEXT,                                     -- 取消原因，任务被取消的理由
    last_priority_change_at TIMESTAMP,                        -- 最后优先级变更时间，最近一次优先级调整的时间
    last_position_change_at TIMESTAMP,                        -- 最后位置变更时间，最近一次位置调整的时间
    queued_by       BIGINT,                                   -- 入队者ID，将任务加入队列的用户ID
    last_modified_by BIGINT,                                  -- 最后修改者ID，最后修改记录的用户ID
    last_action     VARCHAR(50),                              -- 最后操作，最近一次对任务的操作
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 记录创建时间
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 记录更新时间
    UNIQUE(device_id, task_id),                               -- 设备任务唯一约束，同一设备同一任务只能有一条记录
    UNIQUE(device_id, queue_position),                        -- 设备位置唯一约束，同一设备同一位置只能有一个任务
    CONSTRAINT fk_device_task_queue_device FOREIGN KEY (device_id) REFERENCES devices(device_id) ON DELETE CASCADE,
    CONSTRAINT fk_device_task_queue_task FOREIGN KEY (task_id) REFERENCES tasks(task_id) ON DELETE CASCADE
);

-- 12. 设备队列操作历史表 (依赖 devices 和 tasks)
CREATE TABLE IF NOT EXISTS device_queue_operation_history (
    id              BIGSERIAL PRIMARY KEY,                    -- 操作历史记录内部唯一标识符，自增主键
    device_id       VARCHAR(100) NOT NULL,                    -- 设备业务ID，关联到devices表的device_id
    task_id         VARCHAR(100),                             -- 任务业务ID，关联到tasks表的task_id
    operation_type  VARCHAR(30) NOT NULL,                     -- 操作类型：enqueue/dequeue/priority_change/position_change/cancel/restart
    operation_by    VARCHAR(100),                             -- 操作者，执行操作的用户或系统
    operation_time  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 操作时间，执行操作的时间
    old_priority    INTEGER,                                  -- 旧优先级，操作前的任务优先级
    new_priority    INTEGER,                                  -- 新优先级，操作后的任务优先级
    old_position    INTEGER,                                  -- 旧位置，操作前的队列位置
    new_position    INTEGER,                                  -- 新位置，操作后的队列位置
    old_status      VARCHAR(20),                              -- 旧状态，操作前的任务状态
    new_status      VARCHAR(20),                              -- 新状态，操作后的任务状态
    reason          VARCHAR(255),                             -- 操作原因，执行操作的理由
    notes           TEXT,                                     -- 操作备注，操作的详细说明
    operation_source VARCHAR(20) NOT NULL DEFAULT 'system',   -- 操作来源：system/manual/api
    batch_id        VARCHAR(50),                              -- 批次ID，批量操作的批次标识
    is_batch_operation BOOLEAN NOT NULL DEFAULT FALSE,        -- 是否批量操作，是否为批量操作的一部分
    CONSTRAINT fk_device_queue_operation_history_device FOREIGN KEY (device_id) REFERENCES devices(device_id) ON DELETE CASCADE,
    CONSTRAINT fk_device_queue_operation_history_task FOREIGN KEY (task_id) REFERENCES tasks(task_id) ON DELETE SET NULL
);

-- 创建索引
CREATE INDEX idx_tasks_status ON tasks(status);               -- 任务状态索引，用于按状态查询任务
CREATE INDEX idx_tasks_priority ON tasks(priority);           -- 任务优先级索引，用于按优先级查询任务
CREATE INDEX idx_tasks_execute_time ON tasks(execute_time);   -- 任务执行时间索引，用于按执行时间查询任务
CREATE INDEX idx_tasks_device_id ON tasks(device_id);         -- 任务设备ID索引，用于按设备查询任务
CREATE INDEX idx_task_executions_task_id ON task_executions(task_id); -- 任务执行记录任务ID索引，用于查询任务执行历史
CREATE INDEX idx_devices_status ON devices(status);           -- 设备状态索引，用于按状态查询设备
CREATE INDEX idx_device_heartbeats_device_id ON device_heartbeats(device_id); -- 设备心跳设备ID索引，用于查询设备心跳历史
CREATE INDEX idx_device_logs_device_id ON device_logs(device_id); -- 设备日志设备ID索引，用于查询设备日志
CREATE INDEX idx_device_commands_device_id ON device_commands(device_id); -- 设备命令设备ID索引，用于查询设备命令历史
CREATE INDEX idx_task_assignment_queue_task_id ON task_assignment_queue(task_id); -- 任务分配队列任务ID索引，用于查询任务分配状态
CREATE INDEX idx_device_load_monitor_device_id ON device_load_monitor(device_id); -- 设备负载监控设备ID索引，用于查询设备负载信息
CREATE INDEX idx_device_task_queue_device_id ON device_task_queue(device_id); -- 设备任务队列设备ID索引，用于查询设备任务队列
CREATE INDEX idx_device_queue_operation_history_device_id ON device_queue_operation_history(device_id); -- 设备队列操作历史设备ID索引，用于查询设备操作历史

-- 创建触发器
CREATE TRIGGER update_tasks_updated_at                       -- 任务表更新时间触发器，自动更新updated_at字段
    BEFORE UPDATE ON tasks
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_devices_updated_at                      -- 设备表更新时间触发器，自动更新updated_at字段
    BEFORE UPDATE ON devices
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_task_assignment_queue_updated_at        -- 任务分配队列表更新时间触发器，自动更新updated_at字段
    BEFORE UPDATE ON task_assignment_queue
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_device_load_monitor_updated_at          -- 设备负载监控表更新时间触发器，自动更新updated_at字段
    BEFORE UPDATE ON device_load_monitor
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_device_task_queue_updated_at            -- 设备任务队列表更新时间触发器，自动更新updated_at字段
    BEFORE UPDATE ON device_task_queue
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 创建视图
CREATE OR REPLACE VIEW device_queue_status AS                 -- 设备队列状态视图，提供设备任务队列的汇总信息
SELECT
    d.device_id,                                        -- 设备业务 ID
    d.name,                                             -- 设备名称
    COUNT(dtq.id) AS total_tasks,                       -- 所有关联任务总数
    COUNT(dtq.id) FILTER (WHERE dtq.status = 'queued')    AS pending_tasks,    -- 待调度任务数
    COUNT(dtq.id) FILTER (WHERE dtq.status = 'running')   AS running_tasks,    -- 正在执行任务数
    COUNT(dtq.id) FILTER (WHERE dtq.status = 'completed') AS completed_tasks,  -- 已完成任务数
    COUNT(dtq.id) FILTER (WHERE dtq.status = 'failed')    AS failed_tasks,     -- 失败任务数
    MAX(dtq.updated_at) AS last_queue_update             -- 最近一次任务状态变化时间
FROM devices d
LEFT JOIN device_task_queue dtq
  ON dtq.device_id = d.device_id                       -- 业务主键关联
GROUP BY
    d.device_id,
    d.name;