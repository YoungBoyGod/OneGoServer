-- 增强版任务分配队列表 (可选优化)
-- 基于用户流程需求，增加一些实用字段

-- 创建任务分配队列表
CREATE TABLE IF NOT EXISTS task_assignment_queue (
    id              BIGSERIAL PRIMARY KEY,
    task_id         VARCHAR(100) NOT NULL UNIQUE,  -- 任务ID 
    priority        INTEGER NOT NULL DEFAULT 5, -- 队列优先级
    queue_status    VARCHAR(20) NOT NULL DEFAULT 'queued', -- 队列状态
    
    -- 分配条件
    required_device_type VARCHAR(50), -- 需要的设备类型   
    required_capabilities JSONB, -- 设备能力要求 （如：{"min_cpu": 2, "min_memory": 4096, "supported_protocols": ["http", "mqtt"]}）
    preferred_device_ids BIGINT[], -- 首选设备ID列表 (白名单)
    excluded_device_ids BIGINT[], -- 排除设备ID列表 (黑名单)
    
    -- 分配策略 (新增)
    assignment_strategy VARCHAR(50) DEFAULT 'load_balance', -- 分配策略: load_balance, round_robin, affinity, manual
    affinity_rules JSONB, -- 亲和性规则 {"prefer_same_device": true, "prefer_same_type": false}
    
    -- 分配结果
    assigned_device_id BIGINT, -- 分配的设备ID
    assigned_at      TIMESTAMP, -- 分配时间
    assignment_score NUMERIC(5,2), -- 分配得分(算法评估)
    
    -- 队列信息
    queue_position   INTEGER, -- 队列位置
    estimated_wait_time INTEGER, -- 预估等待时间(秒)
    retry_count      INTEGER DEFAULT 0, -- 分配重试次数
    max_retries      INTEGER DEFAULT 3, -- 最大重试次数
    
    -- 优先级调整历史 (新增)
    original_priority INTEGER, -- 原始优先级
    last_priority_change_at TIMESTAMP, -- 最后优先级调整时间
    priority_change_reason VARCHAR(255), -- 优先级调整原因
    priority_boost_reason VARCHAR(100), -- 自动提升原因
    
    -- 操作来源 (新增)
    operation_source VARCHAR(50) DEFAULT 'system', -- 操作来源: system, manual, api, scheduler
    
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
    
    -- 性能指标 (新增)
    avg_task_duration NUMERIC(10,2), -- 平均任务执行时间(秒)
    last_task_completion TIMESTAMP, -- 最后任务完成时间
    
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    
    -- 注意：移除外键约束，改为应用层维护数据一致性
);

-- 创建分配历史记录表
CREATE TABLE IF NOT EXISTS task_assignment_history (
    id              BIGSERIAL PRIMARY KEY,
    task_id         VARCHAR(100) NOT NULL,
    device_id       BIGINT,
    
    action          VARCHAR(20) NOT NULL, -- queued, assigned, reassigned, failed, completed, priority_changed
    previous_status VARCHAR(20), -- 前一个状态
    new_status      VARCHAR(20), -- 新状态
    reason          VARCHAR(255), -- 操作原因
    details         JSONB, -- 详细信息 （如：{"algorithm_score": 85.5, "retry_count": 2, "device_selection_reason": "best_load"}）
    
    -- 操作来源 (新增)
    operation_source VARCHAR(50) DEFAULT 'system', -- 操作来源: system, manual, api, scheduler
    operator_id     VARCHAR(100), -- 操作者ID
    
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

-- 新增索引
CREATE INDEX IF NOT EXISTS idx_task_assignment_queue_strategy ON task_assignment_queue(assignment_strategy);
CREATE INDEX IF NOT EXISTS idx_task_assignment_queue_source ON task_assignment_queue(operation_source);
CREATE INDEX IF NOT EXISTS idx_task_assignment_queue_priority_change ON task_assignment_queue(last_priority_change_at);

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

-- 新增历史索引
CREATE INDEX IF NOT EXISTS idx_task_assignment_history_source ON task_assignment_history(operation_source);
CREATE INDEX IF NOT EXISTS idx_task_assignment_history_operator ON task_assignment_history(operator_id);

-- 创建JSON字段索引 (PostgreSQL特性)
CREATE INDEX IF NOT EXISTS idx_task_assignment_queue_capabilities_gin ON task_assignment_queue USING GIN(required_capabilities);
CREATE INDEX IF NOT EXISTS idx_task_assignment_queue_affinity_gin ON task_assignment_queue USING GIN(affinity_rules);
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
    
    -- 记录原始优先级
    IF NEW.original_priority IS NULL THEN
        NEW.original_priority := NEW.priority;
    END IF;
    
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 创建队列位置更新触发器
CREATE TRIGGER update_queue_position_trigger 
    BEFORE INSERT OR UPDATE ON task_assignment_queue 
    FOR EACH ROW 
    EXECUTE FUNCTION update_queue_position();

-- 创建优先级调整历史记录函数
CREATE OR REPLACE FUNCTION log_priority_change()
RETURNS TRIGGER AS $$
BEGIN
    -- 当优先级发生变化时，记录历史
    IF OLD.priority != NEW.priority THEN
        INSERT INTO task_assignment_history (
            task_id, device_id, action, previous_status, new_status,
            reason, details, operation_source, operator_id
        ) VALUES (
            NEW.task_id, NEW.assigned_device_id, 'priority_changed',
            OLD.priority::VARCHAR, NEW.priority::VARCHAR,
            NEW.priority_change_reason,
            jsonb_build_object(
                'old_priority', OLD.priority,
                'new_priority', NEW.priority,
                'boost_reason', NEW.priority_boost_reason
            ),
            NEW.operation_source, 'system'
        );
        
        NEW.last_priority_change_at := CURRENT_TIMESTAMP;
    END IF;
    
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 创建优先级调整触发器
CREATE TRIGGER log_priority_change_trigger 
    BEFORE UPDATE ON task_assignment_queue 
    FOR EACH ROW 
    EXECUTE FUNCTION log_priority_change();

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

-- 创建动态优先级提升函数 (基于等待时间)
CREATE OR REPLACE FUNCTION auto_boost_waiting_tasks()
RETURNS INTEGER AS $$
DECLARE
    boosted_count INTEGER := 0;
BEGIN
    -- 等待超过30分钟的高优先级任务自动提升
    UPDATE task_assignment_queue 
    SET priority = GREATEST(priority - 1, 1),
        last_priority_change_at = CURRENT_TIMESTAMP,
        priority_boost_reason = 'auto_waiting_time_boost',
        operation_source = 'system'
    WHERE queue_status = 'queued' 
      AND priority <= 5  -- 仅对高优先级任务生效
      AND queued_at < CURRENT_TIMESTAMP - INTERVAL '30 minutes'
      AND (last_priority_change_at IS NULL 
           OR last_priority_change_at < CURRENT_TIMESTAMP - INTERVAL '10 minutes');
    
    GET DIAGNOSTICS boosted_count = ROW_COUNT;
    RETURN boosted_count;
END;
$$ LANGUAGE plpgsql;

-- 创建设备选择评分函数 (增强版)
CREATE OR REPLACE FUNCTION calculate_device_assignment_score(
    p_device_id BIGINT,
    p_task_type VARCHAR(50),
    p_priority INTEGER,
    p_estimated_duration INTEGER
) RETURNS NUMERIC AS $$
DECLARE
    load_score NUMERIC;
    affinity_score NUMERIC;
    history_score NUMERIC;
    final_score NUMERIC;
BEGIN
    -- 负载评分 (40%)
    SELECT dlm.load_score INTO load_score
    FROM device_load_monitor dlm 
    WHERE dlm.device_id = p_device_id;
    
    -- 亲和性评分 (30%) - 设备类型匹配度
    SELECT CASE 
        WHEN d.type = p_task_type THEN 100
        WHEN d.capabilities ? p_task_type THEN 80
        ELSE 50
    END INTO affinity_score
    FROM devices d WHERE d.id = p_device_id;
    
    -- 历史成功率评分 (30%)
    SELECT COALESCE(
        (COUNT(*) FILTER (WHERE status = 'completed') * 100.0 / COUNT(*)), 
        50
    ) INTO history_score
    FROM device_task_queue 
    WHERE device_id = p_device_id 
      AND created_at > CURRENT_TIMESTAMP - INTERVAL '7 days';
    
    -- 综合评分计算
    final_score := (100 - load_score) * 0.4 + affinity_score * 0.3 + history_score * 0.3;
    
    RETURN final_score;
END;
$$ LANGUAGE plpgsql;

-- 创建监控视图
CREATE OR REPLACE VIEW assignment_queue_monitor AS
SELECT 
    COUNT(*) as total_pending,
    COUNT(*) FILTER (WHERE priority <= 3) as high_priority_pending,
    AVG(EXTRACT(EPOCH FROM (CURRENT_TIMESTAMP - queued_at))/60) as avg_wait_minutes,
    MAX(EXTRACT(EPOCH FROM (CURRENT_TIMESTAMP - queued_at))/60) as max_wait_minutes,
    COUNT(DISTINCT required_device_type) as device_types_needed,
    COUNT(DISTINCT assignment_strategy) as strategies_used
FROM task_assignment_queue 
WHERE queue_status = 'queued';

-- 创建设备负载均衡监控视图
CREATE OR REPLACE VIEW device_load_balance_monitor AS
SELECT 
    d.type as device_type,
    COUNT(d.id) as total_devices,
    COUNT(d.id) FILTER (WHERE d.status = 'online') as online_devices,
    AVG(dlm.load_score) as avg_load_score,
    STDDEV(dlm.load_score) as load_balance_score,  -- 标准差越小负载越均衡
    SUM(dlm.current_tasks) as total_running_tasks
FROM devices d
LEFT JOIN device_load_monitor dlm ON d.id = dlm.device_id
GROUP BY d.type; 