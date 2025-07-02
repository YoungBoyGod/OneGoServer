-- 创建任务分配队列表
CREATE TABLE IF NOT EXISTS task_assignment_queue (
    id              BIGSERIAL PRIMARY KEY,
    task_id         VARCHAR(100) NOT NULL UNIQUE,  -- 任务ID 
    priority        INTEGER NOT NULL DEFAULT 5, -- 队列优先级
    queue_status    VARCHAR(20) NOT NULL DEFAULT 'queued', -- 队列状态
    
    -- 分配条件
    required_device_type VARCHAR(50), -- 需要的设备类型   
    required_capabilities JSONB, -- 设备能力要求 （如：{"min_cpu": 2, "min_memory": 4096, "supported_protocols": ["http", "mqtt"]}）
    preferred_device_ids BIGINT[], -- 首选设备ID列表
    excluded_device_ids BIGINT[], -- 排除设备ID列表
    
    -- 分配结果
    assigned_device_id BIGINT, -- 分配的设备ID
    assigned_at      TIMESTAMP, -- 分配时间
    assignment_score NUMERIC(5,2), -- 分配得分(算法评估)
    
    -- 队列信息
    queue_position   INTEGER, -- 队列位置
    estimated_wait_time INTEGER, -- 预估等待时间(秒)
    retry_count      INTEGER DEFAULT 0, -- 分配重试次数
    max_retries      INTEGER DEFAULT 3, -- 最大重试次数
    
    -- 时间戳
    queued_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    
    -- 注意：移除外键约束，改为应用层维护数据一致性
);

-- 创建设备负载监控表
CREATE TABLE IF NOT EXISTS device_load_monitor (
    id              BIGSERIAL PRIMARY KEY,
    device_id       BIGINT NOT NULL,
    
    -- 负载指标
    current_tasks   INTEGER DEFAULT 0, -- 当前任务数
    max_concurrent_tasks INTEGER DEFAULT 1, -- 最大并发任务数
    cpu_load        NUMERIC(5,2), -- CPU负载(%)
    memory_usage    NUMERIC(5,2), -- 内存使用率(%)
    disk_usage      NUMERIC(5,2), -- 磁盘使用率(%)
    network_latency INTEGER, -- 网络延迟(ms)
    
    -- 设备状态
    status          VARCHAR(20) NOT NULL DEFAULT 'online', -- online, offline, busy, maintenance
    last_heartbeat  TIMESTAMP, -- 最后心跳时间
    load_score      NUMERIC(5,2), -- 负载评分(越低越好)
    
    -- 统计信息
    total_assigned  INTEGER DEFAULT 0, -- 累计分配任务数
    total_completed INTEGER DEFAULT 0, -- 累计完成任务数
    total_failed    INTEGER DEFAULT 0, -- 累计失败任务数
    success_rate    NUMERIC(5,2), -- 成功率(%)
    
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    
    -- 注意：移除外键约束，改为应用层维护数据一致性
);

-- 创建分配历史记录表
CREATE TABLE IF NOT EXISTS task_assignment_history (
    id              BIGSERIAL PRIMARY KEY,
    task_id         VARCHAR(100) NOT NULL,
    device_id       BIGINT,
    
    action          VARCHAR(20) NOT NULL, -- queued, assigned, reassigned, failed, completed
    previous_status VARCHAR(20), -- 前一个状态
    new_status      VARCHAR(20), -- 新状态
    reason          VARCHAR(255), -- 操作原因
    details         JSONB, -- 详细信息 （如：{"algorithm_score": 85.5, "retry_count": 2, "device_selection_reason": "best_load"}）
    
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    
    -- 注意：移除外键约束，改为应用层维护数据一致性
);

-- 创建队列管理索引
CREATE INDEX IF NOT EXISTS idx_task_assignment_queue_task_id ON task_assignment_queue(task_id);
CREATE INDEX IF NOT EXISTS idx_task_assignment_queue_status ON task_assignment_queue(queue_status);
CREATE INDEX IF NOT EXISTS idx_task_assignment_queue_priority ON task_assignment_queue(priority DESC);
CREATE INDEX IF NOT EXISTS idx_task_assignment_queue_device_id ON task_assignment_queue(assigned_device_id);
CREATE INDEX IF NOT EXISTS idx_task_assignment_queue_position ON task_assignment_queue(queue_position);
CREATE INDEX IF NOT EXISTS idx_task_assignment_queue_queued_at ON task_assignment_queue(queued_at);

-- 创建复合索引
CREATE INDEX IF NOT EXISTS idx_task_assignment_queue_status_priority ON task_assignment_queue(queue_status, priority DESC);
CREATE INDEX IF NOT EXISTS idx_task_assignment_queue_device_type ON task_assignment_queue(required_device_type);

-- 创建设备负载监控索引
CREATE INDEX IF NOT EXISTS idx_device_load_monitor_device_id ON device_load_monitor(device_id);
CREATE INDEX IF NOT EXISTS idx_device_load_monitor_status ON device_load_monitor(status);
CREATE INDEX IF NOT EXISTS idx_device_load_monitor_load_score ON device_load_monitor(load_score);
CREATE INDEX IF NOT EXISTS idx_device_load_monitor_heartbeat ON device_load_monitor(last_heartbeat);

-- 创建复合索引
CREATE INDEX IF NOT EXISTS idx_device_load_monitor_status_load ON device_load_monitor(status, load_score);
CREATE INDEX IF NOT EXISTS idx_device_load_monitor_status_tasks ON device_load_monitor(status, current_tasks);

-- 创建分配历史索引
CREATE INDEX IF NOT EXISTS idx_task_assignment_history_task_id ON task_assignment_history(task_id);
CREATE INDEX IF NOT EXISTS idx_task_assignment_history_device_id ON task_assignment_history(device_id);
CREATE INDEX IF NOT EXISTS idx_task_assignment_history_action ON task_assignment_history(action);
CREATE INDEX IF NOT EXISTS idx_task_assignment_history_created_at ON task_assignment_history(created_at);

-- 创建JSON字段索引 (PostgreSQL特性)
CREATE INDEX IF NOT EXISTS idx_task_assignment_queue_capabilities_gin ON task_assignment_queue USING GIN(required_capabilities);
CREATE INDEX IF NOT EXISTS idx_task_assignment_history_details_gin ON task_assignment_history USING GIN(details);

-- 为队列表创建自动更新触发器
CREATE TRIGGER update_task_assignment_queue_updated_at 
    BEFORE UPDATE ON task_assignment_queue 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- 为设备负载监控表创建自动更新触发器
CREATE TRIGGER update_device_load_monitor_updated_at 
    BEFORE UPDATE ON device_load_monitor 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- 创建队列位置自动更新函数
CREATE OR REPLACE FUNCTION update_queue_position()
RETURNS TRIGGER AS $$
BEGIN
    -- 当队列状态变为queued时，自动分配队列位置
    IF NEW.queue_status = 'queued' AND (OLD.queue_status IS NULL OR OLD.queue_status != 'queued') THEN
        SELECT COALESCE(MAX(queue_position), 0) + 1 
        INTO NEW.queue_position 
        FROM task_assignment_queue 
        WHERE queue_status = 'queued' AND id != NEW.id;
    END IF;
    
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 创建队列位置更新触发器
CREATE TRIGGER update_queue_position_trigger 
    BEFORE INSERT OR UPDATE ON task_assignment_queue 
    FOR EACH ROW 
    EXECUTE FUNCTION update_queue_position();

-- 创建设备负载自动计算函数
CREATE OR REPLACE FUNCTION calculate_device_load_score()
RETURNS TRIGGER AS $$
BEGIN
    -- 计算负载评分：CPU(40%) + 内存(30%) + 任务数(20%) + 网络延迟(10%)
    NEW.load_score := (
        COALESCE(NEW.cpu_load, 0) * 0.4 +
        COALESCE(NEW.memory_usage, 0) * 0.3 +
        (CASE 
            WHEN NEW.max_concurrent_tasks > 0 THEN 
                (NEW.current_tasks::NUMERIC / NEW.max_concurrent_tasks * 100) * 0.2
            ELSE 0 
        END) +
        LEAST(COALESCE(NEW.network_latency, 0) / 10.0, 100) * 0.1
    );
    
    -- 计算成功率
    IF NEW.total_assigned > 0 THEN
        NEW.success_rate := (NEW.total_completed::NUMERIC / NEW.total_assigned * 100);
    ELSE
        NEW.success_rate := 100.0;
    END IF;
    
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 创建负载评分计算触发器
CREATE TRIGGER calculate_load_score_trigger 
    BEFORE INSERT OR UPDATE ON device_load_monitor 
    FOR EACH ROW 
    EXECUTE FUNCTION calculate_device_load_score(); 