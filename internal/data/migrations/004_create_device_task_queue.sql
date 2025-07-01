-- 创建设备任务队列表
CREATE TABLE IF NOT EXISTS device_task_queue (
    id              BIGSERIAL PRIMARY KEY,
    device_id       BIGINT NOT NULL,
    device_esn      VARCHAR(100) NOT NULL,
    task_id         BIGINT NOT NULL,
    
    -- 队列管理
    queue_priority  INTEGER NOT NULL DEFAULT 5, -- 任务优先级(1-10, 数字越小优先级越高)
    queue_position  INTEGER NOT NULL, -- 在该设备上的队列位置(1, 2, 3...)
    original_priority INTEGER, -- 原始优先级(用于重置)
    is_manual_priority BOOLEAN DEFAULT FALSE, -- 是否手动调整过优先级
    is_manual_position BOOLEAN DEFAULT FALSE, -- 是否手动调整过位置
    
    -- 状态信息
    status          VARCHAR(20) NOT NULL DEFAULT 'queued', -- queued, executing, paused, completed, failed, canceled
    estimated_start_time TIMESTAMP, -- 预估开始时间
    estimated_duration INTEGER, -- 预估执行时长(秒)
    actual_start_time TIMESTAMP, -- 实际开始时间
    actual_end_time TIMESTAMP, -- 实际结束时间
    
    -- 执行配置
    max_retry_count INTEGER DEFAULT 3, -- 最大重试次数
    current_retry   INTEGER DEFAULT 0, -- 当前重试次数
    timeout_seconds INTEGER DEFAULT 3600, -- 超时时间(秒)
    
    -- 依赖关系
    depends_on_task_ids BIGINT[], -- 依赖的任务ID列表
    blocks_task_ids    BIGINT[], -- 阻塞的任务ID列表
    
    -- 操作记录
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    queued_by       BIGINT, -- 入队操作者
    last_modified_by BIGINT, -- 最后修改者
    last_action     VARCHAR(50), -- 最后操作 (added, priority_changed, position_changed, started, paused, etc.)
    
    -- 外键约束
    FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE,
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    FOREIGN KEY (queued_by) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (last_modified_by) REFERENCES users(id) ON DELETE SET NULL,
    
    -- 唯一约束
    UNIQUE(device_id, task_id), -- 同一设备上的同一任务只能有一条记录
    UNIQUE(device_id, queue_position) -- 同一设备上队列位置唯一
);

-- 创建设备队列操作历史表
CREATE TABLE IF NOT EXISTS device_queue_operation_history (
    id              BIGSERIAL PRIMARY KEY,
    device_id       BIGINT NOT NULL,
    device_esn      VARCHAR(100) NOT NULL,
    task_id         BIGINT,
    
    -- 操作信息
    operation_type  VARCHAR(30) NOT NULL, -- add, remove, priority_change, position_change, start, pause, resume, cancel
    operation_by    BIGINT, -- 操作者ID
    operation_time  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- 变更详情
    old_priority    INTEGER, -- 变更前优先级
    new_priority    INTEGER, -- 变更后优先级
    old_position    INTEGER, -- 变更前位置
    new_position    INTEGER, -- 变更后位置
    old_status      VARCHAR(20), -- 变更前状态
    new_status      VARCHAR(20), -- 变更后状态
    
    -- 操作原因和备注
    reason          VARCHAR(255), -- 操作原因
    notes           TEXT, -- 操作备注
    operation_source VARCHAR(20) DEFAULT 'manual', -- manual, system, api, scheduler
    
    -- 批量操作支持
    batch_id        VARCHAR(50), -- 批量操作ID
    is_batch_operation BOOLEAN DEFAULT FALSE, -- 是否为批量操作
    
    -- 外键约束
    FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE,
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    FOREIGN KEY (operation_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 创建设备队列配置表
CREATE TABLE IF NOT EXISTS device_queue_config (
    id              BIGSERIAL PRIMARY KEY,
    device_id       BIGINT NOT NULL UNIQUE,
    device_esn      VARCHAR(100) NOT NULL,
    
    -- 队列配置
    max_queue_size  INTEGER DEFAULT 100, -- 最大队列长度
    max_concurrent_tasks INTEGER DEFAULT 1, -- 最大并发任务数
    auto_start_tasks BOOLEAN DEFAULT TRUE, -- 是否自动开始任务
    priority_scheduling BOOLEAN DEFAULT TRUE, -- 是否启用优先级调度
    
    -- 调度策略
    scheduling_strategy VARCHAR(20) DEFAULT 'priority_first', -- priority_first, fifo, lifo, weighted
    load_balancing BOOLEAN DEFAULT TRUE, -- 是否启用负载均衡
    
    -- 时间窗口配置
    work_start_time TIME, -- 工作开始时间
    work_end_time   TIME, -- 工作结束时间
    timezone        VARCHAR(50) DEFAULT 'UTC', -- 时区
    
    -- 资源限制
    max_cpu_usage   NUMERIC(5,2) DEFAULT 80.0, -- 最大CPU使用率
    max_memory_usage NUMERIC(5,2) DEFAULT 80.0, -- 最大内存使用率
    min_free_disk   BIGINT DEFAULT 1073741824, -- 最小可用磁盘空间(字节)
    
    -- 通知配置
    notify_on_completion BOOLEAN DEFAULT FALSE, -- 完成时通知
    notify_on_failure   BOOLEAN DEFAULT TRUE, -- 失败时通知
    notification_webhook VARCHAR(255), -- 通知webhook地址
    
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE
);

-- 创建索引
-- 设备任务队列表索引
CREATE INDEX IF NOT EXISTS idx_device_task_queue_device_id ON device_task_queue(device_id);
CREATE INDEX IF NOT EXISTS idx_device_task_queue_device_esn ON device_task_queue(device_esn);
CREATE INDEX IF NOT EXISTS idx_device_task_queue_task_id ON device_task_queue(task_id);
CREATE INDEX IF NOT EXISTS idx_device_task_queue_status ON device_task_queue(status);
CREATE INDEX IF NOT EXISTS idx_device_task_queue_priority ON device_task_queue(queue_priority);
CREATE INDEX IF NOT EXISTS idx_device_task_queue_position ON device_task_queue(queue_position);
CREATE INDEX IF NOT EXISTS idx_device_task_queue_start_time ON device_task_queue(estimated_start_time);

-- 复合索引
CREATE INDEX IF NOT EXISTS idx_device_task_queue_device_status ON device_task_queue(device_id, status);
CREATE INDEX IF NOT EXISTS idx_device_task_queue_device_priority ON device_task_queue(device_id, queue_priority DESC);
CREATE INDEX IF NOT EXISTS idx_device_task_queue_device_position ON device_task_queue(device_id, queue_position);
CREATE INDEX IF NOT EXISTS idx_device_task_queue_device_manual ON device_task_queue(device_id, is_manual_priority, is_manual_position);

-- 操作历史表索引
CREATE INDEX IF NOT EXISTS idx_device_queue_history_device_id ON device_queue_operation_history(device_id);
CREATE INDEX IF NOT EXISTS idx_device_queue_history_task_id ON device_queue_operation_history(task_id);
CREATE INDEX IF NOT EXISTS idx_device_queue_history_operation_type ON device_queue_operation_history(operation_type);
CREATE INDEX IF NOT EXISTS idx_device_queue_history_operation_time ON device_queue_operation_history(operation_time);
CREATE INDEX IF NOT EXISTS idx_device_queue_history_operation_by ON device_queue_operation_history(operation_by);
CREATE INDEX IF NOT EXISTS idx_device_queue_history_batch_id ON device_queue_operation_history(batch_id);

-- 设备队列配置表索引
CREATE INDEX IF NOT EXISTS idx_device_queue_config_device_id ON device_queue_config(device_id);
CREATE INDEX IF NOT EXISTS idx_device_queue_config_device_esn ON device_queue_config(device_esn);

-- 创建触发器函数：自动更新队列位置
CREATE OR REPLACE FUNCTION auto_update_device_queue_position()
RETURNS TRIGGER AS $$
BEGIN
    -- 插入新任务时，自动分配队列位置
    IF TG_OP = 'INSERT' THEN
        -- 如果没有指定位置，自动分配到队列末尾
        IF NEW.queue_position IS NULL OR NEW.queue_position = 0 THEN
            SELECT COALESCE(MAX(queue_position), 0) + 1 
            INTO NEW.queue_position 
            FROM device_task_queue 
            WHERE device_id = NEW.device_id AND status = 'queued' AND id != NEW.id;
        ELSE
            -- 如果指定了位置，需要调整其他任务的位置
            UPDATE device_task_queue 
            SET queue_position = queue_position + 1,
                updated_at = CURRENT_TIMESTAMP
            WHERE device_id = NEW.device_id 
              AND queue_position >= NEW.queue_position 
              AND id != NEW.id
              AND status = 'queued';
        END IF;
    END IF;
    
    -- 更新时处理位置变更
    IF TG_OP = 'UPDATE' AND OLD.queue_position != NEW.queue_position THEN
        -- 记录这是手动调整的位置
        NEW.is_manual_position := TRUE;
        NEW.last_action := 'position_changed';
        
        -- 调整其他任务位置
        IF NEW.queue_position > OLD.queue_position THEN
            -- 向后移动，前面的任务位置减1
            UPDATE device_task_queue 
            SET queue_position = queue_position - 1,
                updated_at = CURRENT_TIMESTAMP
            WHERE device_id = NEW.device_id 
              AND queue_position > OLD.queue_position 
              AND queue_position <= NEW.queue_position
              AND id != NEW.id
              AND status = 'queued';
        ELSE
            -- 向前移动，后面的任务位置加1
            UPDATE device_task_queue 
            SET queue_position = queue_position + 1,
                updated_at = CURRENT_TIMESTAMP
            WHERE device_id = NEW.device_id 
              AND queue_position >= NEW.queue_position 
              AND queue_position < OLD.queue_position
              AND id != NEW.id
              AND status = 'queued';
        END IF;
    END IF;
    
    -- 优先级变更
    IF TG_OP = 'UPDATE' AND OLD.queue_priority != NEW.queue_priority THEN
        NEW.is_manual_priority := TRUE;
        NEW.last_action := 'priority_changed';
    END IF;
    
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 创建队列位置自动更新触发器
CREATE TRIGGER auto_update_device_queue_position_trigger 
    BEFORE INSERT OR UPDATE ON device_task_queue 
    FOR EACH ROW 
    EXECUTE FUNCTION auto_update_device_queue_position();

-- 创建操作历史记录触发器函数
CREATE OR REPLACE FUNCTION log_device_queue_operation()
RETURNS TRIGGER AS $$
BEGIN
    -- 插入操作历史记录
    IF TG_OP = 'INSERT' THEN
        INSERT INTO device_queue_operation_history (
            device_id, device_esn, task_id, operation_type,
            new_priority, new_position, new_status,
            reason, operation_source
        ) VALUES (
            NEW.device_id, NEW.device_esn, NEW.task_id, 'add',
            NEW.queue_priority, NEW.queue_position, NEW.status,
            'Task added to device queue', 'system'
        );
        RETURN NEW;
    END IF;
    
    IF TG_OP = 'UPDATE' THEN
        -- 优先级变更
        IF OLD.queue_priority != NEW.queue_priority THEN
            INSERT INTO device_queue_operation_history (
                device_id, device_esn, task_id, operation_type,
                old_priority, new_priority, operation_source
            ) VALUES (
                NEW.device_id, NEW.device_esn, NEW.task_id, 'priority_change',
                OLD.queue_priority, NEW.queue_priority,
                CASE WHEN NEW.is_manual_priority THEN 'manual' ELSE 'system' END
            );
        END IF;
        
        -- 位置变更
        IF OLD.queue_position != NEW.queue_position THEN
            INSERT INTO device_queue_operation_history (
                device_id, device_esn, task_id, operation_type,
                old_position, new_position, operation_source
            ) VALUES (
                NEW.device_id, NEW.device_esn, NEW.task_id, 'position_change',
                OLD.queue_position, NEW.queue_position,
                CASE WHEN NEW.is_manual_position THEN 'manual' ELSE 'system' END
            );
        END IF;
        
        -- 状态变更
        IF OLD.status != NEW.status THEN
            INSERT INTO device_queue_operation_history (
                device_id, device_esn, task_id, operation_type,
                old_status, new_status, operation_source
            ) VALUES (
                NEW.device_id, NEW.device_esn, NEW.task_id, NEW.status,
                OLD.status, NEW.status, 'system'
            );
        END IF;
        
        RETURN NEW;
    END IF;
    
    IF TG_OP = 'DELETE' THEN
        INSERT INTO device_queue_operation_history (
            device_id, device_esn, task_id, operation_type,
            old_priority, old_position, old_status,
            reason, operation_source
        ) VALUES (
            OLD.device_id, OLD.device_esn, OLD.task_id, 'remove',
            OLD.queue_priority, OLD.queue_position, OLD.status,
            'Task removed from device queue', 'system'
        );
        RETURN OLD;
    END IF;
    
    RETURN NULL;
END;
$$ language 'plpgsql';

-- 创建操作历史记录触发器
CREATE TRIGGER log_device_queue_operation_trigger 
    AFTER INSERT OR UPDATE OR DELETE ON device_task_queue 
    FOR EACH ROW 
    EXECUTE FUNCTION log_device_queue_operation();

-- 为队列配置表创建自动更新触发器
CREATE TRIGGER update_device_queue_config_updated_at 
    BEFORE UPDATE ON device_queue_config 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- 为设备队列表创建自动更新触发器  
CREATE TRIGGER update_device_task_queue_updated_at 
    BEFORE UPDATE ON device_task_queue 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- 创建设备队列状态视图
CREATE OR REPLACE VIEW device_queue_status AS
SELECT 
    d.id as device_id,
    d.device_id as device_esn,
    d.name as device_name,
    COUNT(dtq.id) as total_queued_tasks,
    COUNT(CASE WHEN dtq.status = 'queued' THEN 1 END) as pending_tasks,
    COUNT(CASE WHEN dtq.status = 'executing' THEN 1 END) as executing_tasks,
    COUNT(CASE WHEN dtq.is_manual_priority = true THEN 1 END) as manual_priority_tasks,
    COUNT(CASE WHEN dtq.is_manual_position = true THEN 1 END) as manual_position_tasks,
    MIN(dtq.estimated_start_time) as next_task_start_time,
    SUM(dtq.estimated_duration) as total_estimated_duration,
    MAX(dtq.updated_at) as last_queue_update
FROM devices d
LEFT JOIN device_task_queue dtq ON d.id = dtq.device_id
GROUP BY d.id, d.device_id, d.name; 