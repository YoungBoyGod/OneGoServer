# OneGo服务器数据库设计分析

## 🎯 设计原则

### 数据库选择
- **主数据库**: PostgreSQL - 支持复杂查询、事务、JSON字段
- **缓存数据库**: Redis - 提升查询性能、会话管理
- **消息队列**: Kafka - 异步任务处理、事件通知

### 设计理念
1. **规范化设计**: 避免数据冗余，保持数据一致性
2. **扩展性考虑**: 预留扩展字段，支持业务发展
3. **性能优化**: 合理索引设计，支持高并发访问
4. **数据完整性**: 外键约束、字段验证、业务规则

---

## 📋 Task任务管理表设计

### 1. tasks (任务主表)

#### 表结构设计
```sql
CREATE TABLE tasks (
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
    updated_by      BIGINT,
    
    -- 外键约束
    FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE SET NULL,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
```

#### 字段说明
| 字段名 | 类型 | 说明 | 枚举值/示例 |
|--------|------|------|-------------|
| `id` | BIGSERIAL | 主键ID | 自增长 |
| `name` | VARCHAR(255) | 任务名称 | "数据备份任务" |
| `description` | TEXT | 任务描述 | 详细说明 |
| `type` | VARCHAR(50) | 任务类型 | backup, sync, monitor, custom |
| `status` | VARCHAR(20) | 任务状态 | pending, running, completed, failed, canceled |
| `priority` | INTEGER | 优先级(1-10) | 1=最低, 10=最高, 默认5 |
| `execute_time` | TIMESTAMP | 计划执行时间 | 可为空=立即执行 |
| `timeout` | INTEGER | 超时时间(秒) | 默认300秒 |
| `retry_count` | INTEGER | 当前重试次数 | 默认0 |
| `max_retries` | INTEGER | 最大重试次数 | 默认3次 |
| `parameters` | JSONB | 任务参数 | {"input": "data.csv", "output": "/backup/"} |
| `result` | JSONB | 执行结果 | {"status": "success", "files": 10} |
| `error_message` | TEXT | 错误信息 | 失败时的详细错误 |
| `executor_type` | VARCHAR(50) | 执行器类型 | local, remote, container, lambda |
| `executor_id` | VARCHAR(100) | 执行器ID | 具体执行器的标识 |
| `device_id` | BIGINT | 关联设备ID | 外键关联devices表 |

#### 索引设计
```sql
-- 基础查询索引
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_type ON tasks(type);
CREATE INDEX idx_tasks_priority ON tasks(priority);
CREATE INDEX idx_tasks_execute_time ON tasks(execute_time);
CREATE INDEX idx_tasks_device_id ON tasks(device_id);
CREATE INDEX idx_tasks_created_at ON tasks(created_at);

-- 复合索引
CREATE INDEX idx_tasks_status_priority ON tasks(status, priority DESC);
CREATE INDEX idx_tasks_type_status ON tasks(type, status);

-- JSON字段索引 (PostgreSQL特性)
CREATE INDEX idx_tasks_parameters_gin ON tasks USING GIN(parameters);
CREATE INDEX idx_tasks_result_gin ON tasks USING GIN(result);
```

### 2. task_executions (任务执行记录表)

#### 表结构设计
```sql
CREATE TABLE task_executions (
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
```

#### 索引设计
```sql
CREATE INDEX idx_task_executions_task_id ON task_executions(task_id);
CREATE INDEX idx_task_executions_status ON task_executions(status);
CREATE INDEX idx_task_executions_start_time ON task_executions(start_time);
CREATE UNIQUE INDEX idx_task_executions_execution_id ON task_executions(execution_id);
```

---

## 📱 Device设备管理表设计

### 1. devices (设备主表)

#### 表结构设计
```sql
CREATE TABLE devices (
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
    updated_by      BIGINT,
    
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
```

#### 字段说明
| 字段名 | 类型 | 说明 | 枚举值/示例 |
|--------|------|------|-------------|
| `device_id` | VARCHAR(100) | 设备唯一标识 | "SENSOR_001", "CAMERA_A1" |
| `type` | VARCHAR(50) | 设备类型 | sensor, camera, actuator, gateway |
| `status` | VARCHAR(20) | 设备状态 | online, offline, maintenance, error |
| `health_score` | INTEGER | 健康度 | 0-100，基于多个指标计算 |
| `ip_address` | INET | IP地址 | PostgreSQL INET类型 |
| `protocol` | VARCHAR(20) | 通信协议 | HTTP, MQTT, TCP, UDP |
| `auth_type` | VARCHAR(20) | 认证类型 | none, basic, token, certificate |
| `capabilities` | JSONB | 设备能力 | {"temperature": true, "humidity": true} |

#### 索引设计
```sql
-- 基础查询索引
CREATE UNIQUE INDEX idx_devices_device_id ON devices(device_id);
CREATE INDEX idx_devices_type ON devices(type);
CREATE INDEX idx_devices_status ON devices(status);
CREATE INDEX idx_devices_ip_address ON devices(ip_address);
CREATE INDEX idx_devices_last_seen ON devices(last_seen);

-- 复合索引
CREATE INDEX idx_devices_type_status ON devices(type, status);
CREATE INDEX idx_devices_status_last_seen ON devices(status, last_seen);
```

### 2. device_heartbeats (设备心跳表)

#### 表结构设计
```sql
CREATE TABLE device_heartbeats (
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
```

### 3. device_logs (设备日志表)

#### 表结构设计
```sql
CREATE TABLE device_logs (
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
```

#### 索引设计
```sql
CREATE INDEX idx_device_logs_device_id ON device_logs(device_id);
CREATE INDEX idx_device_logs_level ON device_logs(level);
CREATE INDEX idx_device_logs_log_time ON device_logs(log_time);
CREATE INDEX idx_device_logs_category ON device_logs(category);
```

### 4. device_commands (设备命令表)

#### 表结构设计
```sql
CREATE TABLE device_commands (
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
    
    FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL
);
```

---

## 🔗 关联关系设计

### 表关系图
```
users (用户表)
  ├── 1:N → tasks (创建的任务)
  ├── 1:N → devices (管理的设备)
  └── 1:N → device_commands (发送的命令)

tasks (任务表)
  ├── 1:N → task_executions (执行记录)
  └── N:1 → devices (执行设备)

devices (设备表)
  ├── 1:N → device_heartbeats (心跳记录)
  ├── 1:N → device_logs (设备日志)
  ├── 1:N → device_commands (设备命令)
  └── 1:N → tasks (分配的任务)
```

### 外键约束策略
- **CASCADE**: 主记录删除时，自动删除关联记录（logs, heartbeats）
- **SET NULL**: 主记录删除时，外键设为NULL（tasks.device_id）
- **RESTRICT**: 有关联记录时，禁止删除主记录（重要数据）

---

## 📊 性能优化策略

### 1. 分区表设计
```sql
-- 按时间分区的日志表
CREATE TABLE device_logs_y2025m01 PARTITION OF device_logs
FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');

-- 按设备类型分区的心跳表
CREATE TABLE device_heartbeats_sensor PARTITION OF device_heartbeats
FOR VALUES IN ('sensor');
```

### 2. 数据归档策略
- **热数据**: 近3个月的数据，保持在主表
- **温数据**: 3-12个月的数据，迁移到归档表
- **冷数据**: 12个月以上，压缩存储或迁移到对象存储

### 3. 索引优化
- **复合索引**: 根据查询模式设计多列索引
- **部分索引**: 只对活跃数据建索引
- **GIN索引**: JSON字段的全文检索

---

## 🔧 实施计划

### 阶段一：基础表创建 (1-2天)
1. 创建GORM模型文件
2. 编写数据库迁移脚本
3. 实现基础的Repository接口

### 阶段二：数据访问层 (3-5天)
1. 实现TaskRepository的CRUD操作
2. 实现DeviceRepository的CRUD操作
3. 添加复杂查询和统计功能

### 阶段三：业务逻辑层 (5-7天)
1. 实现TaskBiz的业务规则
2. 实现DeviceBiz的状态管理
3. 添加事务处理和错误处理

### 阶段四：服务编排层 (3-5天)
1. 实现TaskService的流程编排
2. 实现DeviceService的综合管理
3. 集成缓存和消息队列

这个设计为OneGo服务器提供了完整、可扩展的数据基础，支持高并发访问和复杂业务场景！🚀 