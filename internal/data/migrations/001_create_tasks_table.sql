-- 创建任务主表
CREATE TABLE IF NOT EXISTS tasks (
    id              BIGSERIAL PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    type            VARCHAR(50) NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'pending',
    priority        INTEGER NOT NULL DEFAULT 5,
    
    -- 执行配置
    execute_time    TIMESTAMP,
    timeout         INTEGER DEFAULT 300,
    retry_count     INTEGER DEFAULT 0,
    max_retries     INTEGER DEFAULT 3,
    is_urgent       BOOLEAN DEFAULT FALSE,
    
    -- 任务参数和结果 (JSON格式)
    parameters      JSONB,
    result          JSONB,
    error_message   TEXT,
    
    -- 执行信息
    executor_type   VARCHAR(50),
    executor_id     VARCHAR(100),
    device_id       BIGINT,
    
    -- 审计字段
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by      BIGINT,
    updated_by      BIGINT
);

-- 创建任务执行记录表
CREATE TABLE IF NOT EXISTS task_executions (
    id              BIGSERIAL PRIMARY KEY,
    task_id         BIGINT NOT NULL,
    execution_id    VARCHAR(100) NOT NULL UNIQUE,
    
    -- 执行状态
    status          VARCHAR(20) NOT NULL DEFAULT 'started',
    start_time      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    end_time        TIMESTAMP,
    duration        INTEGER, -- 执行耗时(秒)
    
    -- 执行详情
    executor_info   JSONB,
    logs            TEXT,
    metrics         JSONB,
    output          JSONB,
    error_details   JSONB,
    
    -- 资源使用
    cpu_usage       NUMERIC(5,2),
    memory_usage    NUMERIC(10,2),
    io_operations   BIGINT,
    
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
);

-- 创建基础查询索引
CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
CREATE INDEX IF NOT EXISTS idx_tasks_type ON tasks(type);
CREATE INDEX IF NOT EXISTS idx_tasks_priority ON tasks(priority);
CREATE INDEX IF NOT EXISTS idx_tasks_execute_time ON tasks(execute_time);
CREATE INDEX IF NOT EXISTS idx_tasks_device_id ON tasks(device_id);
CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at);

-- 创建复合索引
CREATE INDEX IF NOT EXISTS idx_tasks_status_priority ON tasks(status, priority DESC);
CREATE INDEX IF NOT EXISTS idx_tasks_type_status ON tasks(type, status);

-- 创建JSON字段索引 (PostgreSQL特性)
CREATE INDEX IF NOT EXISTS idx_tasks_parameters_gin ON tasks USING GIN(parameters);
CREATE INDEX IF NOT EXISTS idx_tasks_result_gin ON tasks USING GIN(result);

-- 创建任务执行记录索引
CREATE INDEX IF NOT EXISTS idx_task_executions_task_id ON task_executions(task_id);
CREATE INDEX IF NOT EXISTS idx_task_executions_status ON task_executions(status);
CREATE INDEX IF NOT EXISTS idx_task_executions_start_time ON task_executions(start_time);
CREATE UNIQUE INDEX IF NOT EXISTS idx_task_executions_execution_id ON task_executions(execution_id);

-- 创建更新时间触发器函数
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 为tasks表创建自动更新触发器
CREATE TRIGGER update_tasks_updated_at 
    BEFORE UPDATE ON tasks 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column(); 