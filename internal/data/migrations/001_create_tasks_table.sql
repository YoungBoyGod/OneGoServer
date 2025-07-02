-- 创建任务主表
CREATE TABLE IF NOT EXISTS tasks (
    id              BIGSERIAL PRIMARY KEY, -- 任务ID
    name            VARCHAR(255) NOT NULL, -- 任务名称
    description     TEXT, -- 任务描述
    type            VARCHAR(50) NOT NULL, -- 任务类型 （如：shell 文件传输 数据库操作） 
    status          VARCHAR(20) NOT NULL DEFAULT 'pending', -- 任务状态 （如：pending running completed failed canceled）
    priority        INTEGER NOT NULL DEFAULT 5, -- 任务优先级 （如：1-10）
    
    -- 执行配置
    execute_time    TIMESTAMP, -- 执行时间
    timeout         INTEGER DEFAULT 86400, -- 超时时间(秒) 默认24小时
    retry_count     INTEGER DEFAULT 0, -- 重试次数
    max_retries     INTEGER DEFAULT 3, -- 最大重试次数
    is_urgent       BOOLEAN DEFAULT FALSE, -- 是否紧急
    
    -- 任务参数和结果 (JSON格式)
    parameters      JSONB, -- 任务参数 （如：{"command": "ls -l", "args": ["-l"]}）
    result          JSONB, -- 任务结果 （如：{"output": "ls -l", "exit_code": 0}）
    error_message   TEXT, -- 错误信息 （如："command not found"）
    
    -- 执行信息
    executor_type   VARCHAR(50), -- 执行器类型 （如：shell 文件传输 数据库操作） shell python test_executor
    executor_id     VARCHAR(100), -- 执行器ID 
    device_id       BIGINT, -- 设备ID 
    
    -- 审计字段
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by      BIGINT,
    updated_by      BIGINT
);

-- 创建任务执行记录表
CREATE TABLE IF NOT EXISTS task_executions (
    id              BIGSERIAL PRIMARY KEY, -- 数据库主键ID
    task_id         BIGINT NOT NULL, -- 关联的任务ID
    execution_id    VARCHAR(100) NOT NULL UNIQUE, -- 业务执行记录唯一标识
    
    --分配信息
    device_esn      VARCHAR(100), -- 设备编号



    -- 执行状态
    status          VARCHAR(20) NOT NULL DEFAULT 'started', -- 执行状态 （如：pending running completed failed canceled）
   

    start_time      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 开始时间
    end_time        TIMESTAMP, -- 结束时间
    duration        INTEGER, -- 执行耗时(秒) 计算方式：end_time - start_time
    
    -- 执行详情
    executor_info   JSONB, -- 执行器信息 （如：{"executor_type": "shell", "executor_id": "123"}）
    logs            TEXT, -- 执行日志
    metrics         JSONB, -- 执行指标 （如：{"cpu_usage": 0.5, "memory_usage": 1024, "io_operations": 100}）
    output          JSONB, -- 执行输出 （如：{"output": "ls -l", "exit_code": 0}）
    error_details   JSONB, -- 错误详情 （如：{"error_message": "command not found", "error_code": 123}）
    
    -- 资源使用统计 (聚合数据)
    cpu_usage_avg    NUMERIC(5,2), -- CPU平均使用率(%)
    cpu_usage_peak   NUMERIC(5,2), -- CPU峰值使用率(%)
    memory_usage_avg NUMERIC(10,2), -- 内存平均使用量(MB)
    memory_usage_peak NUMERIC(10,2), -- 内存峰值使用量(MB)
    io_operations_total BIGINT, -- IO操作总次数 
    io_bytes_total   BIGINT, -- IO字节总数
    
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    
    -- 注意：移除外键约束，改为应用层维护数据一致性
);



-- 设备队列信息
CREATE TABLE IF NOT EXISTS device_queues (
    id              BIGSERIAL PRIMARY KEY, -- 数据库主键ID
    device_esn      VARCHAR(100) NOT NULL, -- 设备编号
    queue_name      VARCHAR(100) NOT NULL, -- 队列名称
    queue_priority  INTEGER NOT NULL, -- 队列优先级
    queue_size      INTEGER NOT NULL, -- 队列大小
    queue_remaining INTEGER NOT NULL, -- 队列剩余大小
    queue_position  INTEGER NOT NULL, -- 队列位置

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