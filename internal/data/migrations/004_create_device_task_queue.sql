-- ============================================================================
-- 设备任务队列表优化设计 (v2.0)
-- 特性：简化设计、重发机制、模块化触发器、数据类型统一
-- ============================================================================

-- 1. 通用触发器函数：自动更新 updated_at
CREATE OR REPLACE FUNCTION fn_update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = CURRENT_TIMESTAMP;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 2. 主队列表：device_task_queue
CREATE TABLE IF NOT EXISTS device_task_queue (
  id                       BIGSERIAL   PRIMARY KEY,
  device_id                BIGINT      NOT NULL,            -- 关联 devices(id)，应用层校验
  task_id                  VARCHAR(100) NOT NULL,           -- 业务层唯一任务标识
  queue_priority           INTEGER     NOT NULL DEFAULT 5,  -- 1–100，数字越小优先级越高
  original_priority        INTEGER,                         -- 初始优先级
  is_manual_priority       BOOLEAN     NOT NULL DEFAULT FALSE,
  queue_position           INTEGER,                         -- 1,2,3...，队列位置
  is_manual_position       BOOLEAN     NOT NULL DEFAULT FALSE,
  status                   VARCHAR(20) NOT NULL DEFAULT 'queued', -- 任务状态 （如：pending queued running completed failed canceled）
  estimated_start_time     TIMESTAMP, -- 预估开始时间
  estimated_duration       INTEGER,                         -- 预估执行时长(秒)
  actual_start_time        TIMESTAMP, -- 实际开始时间
  actual_end_time          TIMESTAMP, -- 实际结束时间

  max_retry_count          INTEGER     NOT NULL DEFAULT 3, -- 最大重试次数
  current_retry            INTEGER     NOT NULL DEFAULT 0, -- 当前重试次数
  timeout_seconds          INTEGER     NOT NULL DEFAULT 3600, -- 超时时间(秒)

  depends_on_task_ids      VARCHAR(100)[],                        -- 依赖任务ID列表
  blocks_task_ids          VARCHAR(100)[],                        -- 阻塞任务ID列表

  -- 重发 / 取消
  requeue_count            INTEGER     NOT NULL DEFAULT 0, -- 重发次数
  last_requeue_at          TIMESTAMP, -- 最后重发时间
  is_requeued              BOOLEAN     NOT NULL DEFAULT FALSE, -- 是否重发
  cancel_reason            TEXT, -- 取消原因

  -- 手动调整时间戳
  last_priority_change_at  TIMESTAMP, -- 最后优先级调整时间
  last_position_change_at  TIMESTAMP, -- 最后位置调整时间

  -- 审计
  queued_by                BIGINT, -- 入队者
  last_modified_by         BIGINT, -- 最后修改者
  last_action              VARCHAR(50), -- 最后操作

  created_at               TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at               TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,

  UNIQUE(device_id, task_id), -- 同一设备上任务ID唯一
  UNIQUE(device_id, queue_position) -- 同一设备上队列位置唯一
);

-- 3. 操作历史：device_queue_operation_history
CREATE TABLE IF NOT EXISTS device_queue_operation_history (
  id                   BIGSERIAL   PRIMARY KEY,
  device_id            BIGINT      NOT NULL,
  task_id              VARCHAR(100),
  operation_type       VARCHAR(30) NOT NULL,          -- add/remove/priority_change/position_change/requeue/status_change
  operation_by         BIGINT, -- 操作人
  operation_time       TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 操作时间
  old_priority         INTEGER, -- 旧优先级
  new_priority         INTEGER, -- 新优先级
  old_position         INTEGER, -- 旧位置
  new_position         INTEGER, -- 新位置
  old_status           VARCHAR(20), -- 旧状态
  new_status           VARCHAR(20), -- 新状态
  reason               VARCHAR(255), -- 操作原因
  notes                TEXT, -- 备注
  operation_source     VARCHAR(20) NOT NULL DEFAULT 'system', -- 操作来源
  batch_id             VARCHAR(50), -- 批量ID
  is_batch_operation   BOOLEAN     NOT NULL DEFAULT FALSE -- 是否为批量操作
);


-- 4. 索引
CREATE INDEX IF NOT EXISTS idx_dtq_device        ON device_task_queue(device_id);
CREATE INDEX IF NOT EXISTS idx_dtq_task          ON device_task_queue(task_id);
CREATE INDEX IF NOT EXISTS idx_dtq_status        ON device_task_queue(status);
CREATE INDEX IF NOT EXISTS idx_dtq_priority      ON device_task_queue(queue_priority);
CREATE INDEX IF NOT EXISTS idx_dtq_position      ON device_task_queue(queue_position);
CREATE INDEX IF NOT EXISTS idx_dtq_requeued      ON device_task_queue(is_requeued);
CREATE INDEX IF NOT EXISTS idx_dtq_requeue_time  ON device_task_queue(last_requeue_at);
CREATE INDEX IF NOT EXISTS idx_dtq_estimated_start ON device_task_queue(estimated_start_time);

-- 复合索引
CREATE INDEX IF NOT EXISTS idx_dtq_device_status     ON device_task_queue(device_id, status);
CREATE INDEX IF NOT EXISTS idx_dtq_device_priority   ON device_task_queue(device_id, queue_priority DESC);
CREATE INDEX IF NOT EXISTS idx_dtq_device_position   ON device_task_queue(device_id, queue_position);

-- 历史表索引
CREATE INDEX IF NOT EXISTS idx_hist_device       ON device_queue_operation_history(device_id);
CREATE INDEX IF NOT EXISTS idx_hist_task         ON device_queue_operation_history(task_id);
CREATE INDEX IF NOT EXISTS idx_hist_type         ON device_queue_operation_history(operation_type);
CREATE INDEX IF NOT EXISTS idx_hist_time         ON device_queue_operation_history(operation_time);
CREATE INDEX IF NOT EXISTS idx_hist_batch        ON device_queue_operation_history(batch_id);

-- 5. 触发器：自动分配/调整 queue_position 与 标记手动 priority／position 
CREATE OR REPLACE FUNCTION fn_auto_position_and_manual_flags()
RETURNS TRIGGER AS $$
BEGIN
  -- INSERT 时自动分配位置
  IF TG_OP = 'INSERT' THEN
    IF NEW.queue_position IS NULL OR NEW.queue_position <= 0 THEN
      SELECT COALESCE(MAX(queue_position),0) + 1
      INTO NEW.queue_position
      FROM device_task_queue
      WHERE device_id = NEW.device_id
        AND status = 'queued'
        AND id != NEW.id;
    ELSE
      -- 如果指定了位置，需要调整其他任务的位置
      UPDATE device_task_queue
      SET queue_position = queue_position + 1,
          updated_at = CURRENT_TIMESTAMP
      WHERE device_id = NEW.device_id
        AND queue_position >= NEW.queue_position
        AND status = 'queued'
        AND id != NEW.id;
    END IF;
    
    -- 设置初始优先级
    IF NEW.original_priority IS NULL THEN
      NEW.original_priority := NEW.queue_priority;
    END IF;
  END IF;

  -- UPDATE 时处理手动优先级/位置调整
  IF TG_OP = 'UPDATE' THEN
    -- 位置变动
    IF OLD.queue_position IS DISTINCT FROM NEW.queue_position THEN
      NEW.is_manual_position := TRUE;
      NEW.last_action := 'position_changed';
      NEW.last_position_change_at := CURRENT_TIMESTAMP;
      
      -- 调整其他行
      IF NEW.queue_position > OLD.queue_position THEN
        UPDATE device_task_queue
        SET queue_position = queue_position - 1,
            updated_at = CURRENT_TIMESTAMP
        WHERE device_id = NEW.device_id
          AND queue_position > OLD.queue_position
          AND queue_position <= NEW.queue_position
          AND id != NEW.id
          AND status = 'queued';
      ELSE
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

    -- 优先级变动
    IF OLD.queue_priority IS DISTINCT FROM NEW.queue_priority THEN
      NEW.is_manual_priority := TRUE;
      NEW.last_action := 'priority_changed';
      NEW.last_priority_change_at := CURRENT_TIMESTAMP;
    END IF;
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_auto_position ON device_task_queue;
CREATE TRIGGER trg_auto_position
  BEFORE INSERT OR UPDATE ON device_task_queue
  FOR EACH ROW EXECUTE FUNCTION fn_auto_position_and_manual_flags();

-- 6. 触发器：处理重发（requeue）
CREATE OR REPLACE FUNCTION fn_handle_requeue()
RETURNS TRIGGER AS $$
BEGIN
  IF TG_OP = 'UPDATE'
     AND OLD.status IN ('canceled','failed')
     AND NEW.status = 'queued'
  THEN
    NEW.requeue_count   := OLD.requeue_count + 1;
    NEW.is_requeued     := TRUE;
    NEW.last_requeue_at := CURRENT_TIMESTAMP;
    NEW.last_action     := 'requeued';
    NEW.current_retry   := 0; -- 重置重试次数
    NEW.actual_start_time := NULL; -- 清除之前的开始时间
    NEW.actual_end_time := NULL;   -- 清除之前的结束时间
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_requeue ON device_task_queue;
CREATE TRIGGER trg_requeue
  BEFORE UPDATE ON device_task_queue
  FOR EACH ROW EXECUTE FUNCTION fn_handle_requeue();

-- 7. 触发器：写入操作历史
CREATE OR REPLACE FUNCTION fn_log_device_queue_op()
RETURNS TRIGGER AS $$
BEGIN
  IF TG_OP = 'INSERT' THEN
    INSERT INTO device_queue_operation_history(
      device_id, task_id, operation_type,
      new_priority, new_position, new_status,
      operation_source
    ) VALUES (
      NEW.device_id, NEW.task_id, 'add',
      NEW.queue_priority, NEW.queue_position, NEW.status,
      'system'
    );

  ELSIF TG_OP = 'UPDATE' THEN
    -- priority_change
    IF OLD.queue_priority != NEW.queue_priority THEN
      INSERT INTO device_queue_operation_history(
        device_id, task_id, operation_type,
        old_priority, new_priority, operation_source
      ) VALUES (
        NEW.device_id, NEW.task_id, 'priority_change',
        OLD.queue_priority, NEW.queue_priority,
        CASE WHEN NEW.is_manual_priority THEN 'manual' ELSE 'system' END
      );
    END IF;
    
    -- position_change
    IF OLD.queue_position != NEW.queue_position THEN
      INSERT INTO device_queue_operation_history(
        device_id, task_id, operation_type,
        old_position, new_position, operation_source
      ) VALUES (
        NEW.device_id, NEW.task_id, 'position_change',
        OLD.queue_position, NEW.queue_position,
        CASE WHEN NEW.is_manual_position THEN 'manual' ELSE 'system' END
      );
    END IF;
    
    -- requeue
    IF NEW.last_action = 'requeued' THEN
      INSERT INTO device_queue_operation_history(
        device_id, task_id, operation_type,
        old_status, new_status, operation_source,
        reason
      ) VALUES (
        NEW.device_id, NEW.task_id, 'requeue',
        OLD.status, NEW.status, 'system',
        'Task requeued after failure'
      );
    END IF;
    
    -- status_change (其他状态变更)
    IF OLD.status != NEW.status AND COALESCE(NEW.last_action,'') != 'requeued' THEN
      INSERT INTO device_queue_operation_history(
        device_id, task_id, operation_type,
        old_status, new_status, operation_source
      ) VALUES (
        NEW.device_id, NEW.task_id, 'status_change',
        OLD.status, NEW.status, 'system'
      );
    END IF;

  ELSIF TG_OP = 'DELETE' THEN
    INSERT INTO device_queue_operation_history(
      device_id, task_id, operation_type,
      old_priority, old_position, old_status, operation_source,
      reason
    ) VALUES (
      OLD.device_id, OLD.task_id, 'remove',
      OLD.queue_priority, OLD.queue_position, OLD.status, 'system',
      'Task removed from device queue'
    );
  END IF;

  RETURN CASE WHEN TG_OP = 'DELETE' THEN OLD ELSE NEW END;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_log_queue_op ON device_task_queue;
CREATE TRIGGER trg_log_queue_op
  AFTER INSERT OR UPDATE OR DELETE ON device_task_queue
  FOR EACH ROW EXECUTE FUNCTION fn_log_device_queue_op();

-- 8. 自动更新时间触发器
CREATE TRIGGER trg_update_dtq_updated_at
  BEFORE UPDATE ON device_task_queue
  FOR EACH ROW EXECUTE FUNCTION fn_update_updated_at();

-- 9. 聚合视图：device_queue_status
CREATE OR REPLACE VIEW device_queue_status AS
SELECT
  d.id                                                          AS device_id,
  d.device_id                                                   AS device_esn,
  d.name                                                        AS device_name,
  COUNT(dtq.id)                                                 AS total_tasks,
  COUNT(dtq.id) FILTER (WHERE dtq.status = 'queued')           AS pending_tasks,
  COUNT(dtq.id) FILTER (WHERE dtq.status = 'executing')        AS executing_tasks,
  COUNT(dtq.id) FILTER (WHERE dtq.status = 'completed')        AS completed_tasks,
  COUNT(dtq.id) FILTER (WHERE dtq.status = 'failed')           AS failed_tasks,
  COUNT(dtq.id) FILTER (WHERE dtq.status = 'canceled')         AS canceled_tasks,
  COUNT(dtq.id) FILTER (WHERE dtq.is_requeued = true)          AS requeued_tasks,
  COUNT(dtq.id) FILTER (WHERE dtq.is_manual_priority = true)   AS manual_priority_tasks,
  COUNT(dtq.id) FILTER (WHERE dtq.is_manual_position = true)   AS manual_position_tasks,
  MAX(dtq.last_requeue_at)                                      AS last_requeue_time,
  MIN(dtq.estimated_start_time) FILTER (WHERE dtq.status = 'queued') AS next_task_start_time,
  SUM(dtq.estimated_duration) FILTER (WHERE dtq.status = 'queued')   AS total_estimated_duration,
  MAX(dtq.updated_at)                                           AS last_queue_update,
  AVG(dtq.queue_priority) FILTER (WHERE dtq.status = 'queued') AS avg_queue_priority
FROM devices d
LEFT JOIN device_task_queue dtq ON d.id = dtq.device_id
GROUP BY d.id, d.device_id, d.name;

-- 10. 有用的查询视图：队列详情
CREATE OR REPLACE VIEW device_queue_details AS
SELECT 
  dtq.id,
  dtq.device_id,
  d.device_id as device_esn,
  d.name as device_name,
  dtq.task_id,
  dtq.queue_priority,
  dtq.original_priority,
  dtq.queue_position,
  dtq.status,
  dtq.estimated_start_time,
  dtq.estimated_duration,
  dtq.actual_start_time,
  dtq.actual_end_time,
  dtq.current_retry,
  dtq.max_retry_count,
  dtq.requeue_count,
  dtq.is_requeued,
  dtq.is_manual_priority,
  dtq.is_manual_position,
  dtq.last_action,
  dtq.created_at,
  dtq.updated_at,
  -- 计算等待时间
  CASE 
    WHEN dtq.status = 'queued' AND dtq.estimated_start_time IS NOT NULL 
    THEN GREATEST(0, EXTRACT(EPOCH FROM (dtq.estimated_start_time - CURRENT_TIMESTAMP)))
    ELSE NULL 
  END as estimated_wait_seconds,
  -- 计算执行时长
  CASE 
    WHEN dtq.actual_start_time IS NOT NULL AND dtq.actual_end_time IS NOT NULL
    THEN EXTRACT(EPOCH FROM (dtq.actual_end_time - dtq.actual_start_time))
    ELSE NULL 
  END as actual_duration_seconds
FROM device_task_queue dtq
JOIN devices d ON d.id = dtq.device_id;

-- ============================================================================
-- 队列管理函数示例
-- ============================================================================

-- 重新排队函数
CREATE OR REPLACE FUNCTION requeue_failed_task(
  p_device_id BIGINT,
  p_task_id VARCHAR(100),
  p_reason TEXT DEFAULT 'Manual requeue'
)
RETURNS BOOLEAN AS $$
DECLARE
  v_count INTEGER;
BEGIN
  UPDATE device_task_queue 
  SET status = 'queued',
      cancel_reason = p_reason,
      last_modified_by = 0 -- 系统操作
  WHERE device_id = p_device_id 
    AND task_id = p_task_id 
    AND status IN ('failed', 'canceled');
    
  GET DIAGNOSTICS v_count = ROW_COUNT;
  RETURN v_count > 0;
END;
$$ LANGUAGE plpgsql;

-- 批量调整优先级函数
CREATE OR REPLACE FUNCTION batch_update_priority(
  p_device_id BIGINT,
  p_task_ids VARCHAR(100)[],
  p_new_priority INTEGER,
  p_operator_id BIGINT DEFAULT NULL
)
RETURNS INTEGER AS $$
DECLARE
  v_count INTEGER;
  v_batch_id VARCHAR(50);
BEGIN
  v_batch_id := 'batch_' || EXTRACT(epoch FROM CURRENT_TIMESTAMP)::bigint;
  
  -- 记录批量操作历史
  INSERT INTO device_queue_operation_history(
    device_id, operation_type, operation_by, 
    batch_id, is_batch_operation, reason
  ) VALUES (
    p_device_id, 'batch_priority_change', p_operator_id,
    v_batch_id, true, 'Batch priority update for ' || array_length(p_task_ids, 1) || ' tasks'
  );
  
  -- 批量更新优先级
  UPDATE device_task_queue 
  SET queue_priority = p_new_priority,
      last_modified_by = p_operator_id
  WHERE device_id = p_device_id 
    AND task_id = ANY(p_task_ids)
    AND status = 'queued';
    
  GET DIAGNOSTICS v_count = ROW_COUNT;
  RETURN v_count;
END;
$$ LANGUAGE plpgsql;
