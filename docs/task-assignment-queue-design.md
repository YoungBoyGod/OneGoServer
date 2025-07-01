# 任务分配队列设计方案

## 问题背景

用户提出核心问题：**"现在任务分配队列如何处理呢在表单中如何体现"**

这涉及到任务管理系统的核心功能：
- 任务创建后如何分配到合适的设备执行
- 如何处理任务排队和负载均衡
- 在前端界面中如何展示分配状态和队列信息

## 系统架构设计

### 任务分配流程

```
任务创建 → 进入分配队列 → 设备选择算法 → 分配到设备 → 执行监控 → 完成/失败处理
```

### 核心组件

1. **任务队列管理器** - 管理待分配任务
2. **设备负载监控器** - 监控设备状态和负载
3. **分配调度器** - 执行分配算法
4. **状态跟踪器** - 跟踪任务和设备状态

## 数据库设计

### 1. 扩展任务状态

```sql
-- 原有状态：pending, running, completed, failed, canceled
-- 新增分配相关状态：
-- queued        - 已入队待分配
-- assigning     - 分配中
-- assigned      - 已分配待执行
-- dispatching   - 派发中
-- executing     - 执行中
```

### 2. 任务分配队列表

```sql
CREATE TABLE IF NOT EXISTS task_assignment_queue (
    id              BIGSERIAL PRIMARY KEY,
    task_id         BIGINT NOT NULL UNIQUE,
    priority        INTEGER NOT NULL DEFAULT 5, -- 队列优先级
    queue_status    VARCHAR(20) NOT NULL DEFAULT 'queued', -- 队列状态
    
    -- 分配条件
    required_device_type VARCHAR(50), -- 需要的设备类型
    required_capabilities JSONB, -- 设备能力要求
    preferred_device_ids BIGINT[], -- 首选设备ID列表
    excluded_device_ids BIGINT[], -- 排除设备ID列表
    
    -- 分配结果
    assigned_device_id BIGINT, -- 分配的设备ID
    assigned_device_esn VARCHAR(100), -- 分配的设备编号
    assigned_at      TIMESTAMP, -- 分配时间
    assignment_score NUMERIC(5,2), -- 分配得分(算法评估)
    
    -- 队列信息
    queue_position   INTEGER, -- 队列位置
    estimated_wait_time INTEGER, -- 预估等待时间(秒)
    retry_count      INTEGER DEFAULT 0, -- 分配重试次数
    max_retries      INTEGER DEFAULT 3, -- 最大重试次数
    
    -- 时间戳
    queued_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    FOREIGN KEY (assigned_device_id) REFERENCES devices(id) ON DELETE SET NULL
);
```

### 3. 设备负载监控表

```sql
CREATE TABLE IF NOT EXISTS device_load_monitor (
    id              BIGSERIAL PRIMARY KEY,
    device_id       BIGINT NOT NULL,
    device_esn      VARCHAR(100) NOT NULL,
    
    -- 负载指标
    current_tasks   INTEGER DEFAULT 0, -- 当前任务数
    max_concurrent_tasks INTEGER DEFAULT 1, -- 最大并发任务数
    cpu_load        NUMERIC(5,2), -- CPU负载
    memory_usage    NUMERIC(5,2), -- 内存使用率
    disk_usage      NUMERIC(5,2), -- 磁盘使用率
    network_latency INTEGER, -- 网络延迟(ms)
    
    -- 设备状态
    status          VARCHAR(20) NOT NULL DEFAULT 'online', -- online, offline, busy, maintenance
    last_heartbeat  TIMESTAMP, -- 最后心跳时间
    load_score      NUMERIC(5,2), -- 负载评分(越低越好)
    
    -- 统计信息
    total_assigned  INTEGER DEFAULT 0, -- 累计分配任务数
    total_completed INTEGER DEFAULT 0, -- 累计完成任务数
    total_failed    INTEGER DEFAULT 0, -- 累计失败任务数
    success_rate    NUMERIC(5,2), -- 成功率
    
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE
);
```

### 4. 分配历史记录表

```sql
CREATE TABLE IF NOT EXISTS task_assignment_history (
    id              BIGSERIAL PRIMARY KEY,
    task_id         BIGINT NOT NULL,
    device_id       BIGINT,
    device_esn      VARCHAR(100),
    
    action          VARCHAR(20) NOT NULL, -- queued, assigned, reassigned, failed, completed
    previous_status VARCHAR(20), -- 前一个状态
    new_status      VARCHAR(20), -- 新状态
    reason          VARCHAR(255), -- 操作原因
    details         JSONB, -- 详细信息
    
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE SET NULL
);
```

## 分配算法设计

### 设备选择算法

```json
{
  "algorithm": "weighted_score",
  "factors": [
    {"name": "load_score", "weight": 0.4, "direction": "asc"},
    {"name": "success_rate", "weight": 0.3, "direction": "desc"},
    {"name": "network_latency", "weight": 0.2, "direction": "asc"},
    {"name": "capability_match", "weight": 0.1, "direction": "desc"}
  ],
  "constraints": [
    {"field": "status", "operator": "in", "values": ["online", "busy"]},
    {"field": "current_tasks", "operator": "<", "value": "max_concurrent_tasks"},
    {"field": "last_heartbeat", "operator": ">", "value": "NOW() - INTERVAL '5 minutes'"}
  ]
}
```

### 优先级处理

1. **紧急任务** (is_urgent=true): 立即分配，可抢占资源
2. **高优先级任务** (priority=1-3): 优先队列
3. **普通任务** (priority=4-7): 标准队列
4. **低优先级任务** (priority=8-10): 后台队列

## 前端界面设计

### 1. 任务列表界面

```typescript
interface TaskListItem {
  id: number;
  name: string;
  type: string;
  status: string;
  priority: number;
  
  // 分配信息
  queue_status?: string;
  queue_position?: number;
  estimated_wait_time?: number;
  assigned_device?: {
    id: number;
    name: string;
    esn: string;
    status: string;
  };
  assigned_at?: string;
  
  created_at: string;
  updated_at: string;
}
```

**界面元素**:
- 任务状态标签（不同颜色区分）
- 队列位置提示（"队列中第3位"）
- 分配设备信息
- 预估等待时间
- 分配操作按钮（重新分配、取消分配）

### 2. 队列管理界面

```vue
<template>
  <div class="queue-management">
    <!-- 队列统计 -->
    <div class="queue-stats">
      <StatCard title="队列总数" :value="queueStats.total" />
      <StatCard title="等待分配" :value="queueStats.queued" />
      <StatCard title="分配中" :value="queueStats.assigning" />
      <StatCard title="已分配" :value="queueStats.assigned" />
    </div>
    
    <!-- 队列列表 -->
    <div class="queue-list">
      <QueueItem 
        v-for="item in queueItems" 
        :key="item.id"
        :task="item"
        @assign="handleManualAssign"
        @cancel="handleCancelAssignment"
      />
    </div>
    
    <!-- 批量操作 -->
    <div class="batch-operations">
      <Button @click="batchAssign">批量分配</Button>
      <Button @click="clearQueue">清空队列</Button>
      <Button @click="rebalanceQueue">重新平衡</Button>
    </div>
  </div>
</template>
```

### 3. 设备负载监控界面

```vue
<template>
  <div class="device-monitor">
    <!-- 设备负载概览 -->
    <div class="load-overview">
      <DeviceCard 
        v-for="device in devices" 
        :key="device.id"
        :device="device"
        :load="device.load_info"
        @view-tasks="showDeviceTasks"
      />
    </div>
    
    <!-- 负载图表 -->
    <div class="load-charts">
      <LineChart 
        title="设备负载趋势"
        :data="loadTrendData"
        :options="chartOptions"
      />
    </div>
  </div>
</template>
```

### 4. 任务分配详情界面

```vue
<template>
  <div class="assignment-details">
    <!-- 分配时间线 -->
    <Timeline>
      <TimelineItem 
        v-for="history in assignmentHistory"
        :key="history.id"
        :timestamp="history.created_at"
        :status="history.action"
        :description="history.reason"
      />
    </Timeline>
    
    <!-- 分配算法结果 -->
    <div class="algorithm-result">
      <h3>分配算法评估</h3>
      <ScoreCard 
        :score="assignment.assignment_score"
        :factors="algorithmFactors"
      />
    </div>
  </div>
</template>
```

## API接口设计

### 队列管理接口

```typescript
// 获取队列状态
GET /api/v1/tasks/queue
Response: {
  total: number;
  by_status: {[status: string]: number};
  queue_items: QueueItem[];
  statistics: QueueStatistics;
}

// 手动分配任务
POST /api/v1/tasks/{taskId}/assign
Body: {
  device_id?: number;
  force?: boolean; // 强制分配
}

// 重新分配任务
POST /api/v1/tasks/{taskId}/reassign
Body: {
  reason: string;
  preferred_device_id?: number;
}

// 取消分配
DELETE /api/v1/tasks/{taskId}/assignment
Body: {
  reason: string;
}

// 队列批量操作
POST /api/v1/tasks/queue/batch
Body: {
  action: 'assign' | 'cancel' | 'rebalance';
  task_ids: number[];
  options?: any;
}
```

### 设备负载接口

```typescript
// 获取设备负载信息
GET /api/v1/devices/load
Response: {
  devices: DeviceLoadInfo[];
  summary: LoadSummary;
}

// 更新设备负载
PUT /api/v1/devices/{deviceId}/load
Body: {
  current_tasks: number;
  cpu_load: number;
  memory_usage: number;
  // ... 其他负载指标
}
```

## 状态流转图

```
pending → queued → assigning → assigned → dispatching → executing → completed
   ↓         ↓         ↓          ↓           ↓            ↓
failed ← failed ← failed ← failed ← failed ← failed
   ↓         ↓         ↓          ↓           ↓            ↓
canceled ← canceled ← canceled ← canceled ← canceled ← canceled
```

## 实施优先级

### 阶段一：基础分配 (当前)
- ✅ 添加device_esn字段
- 🔄 扩展任务状态定义
- 🔄 实现基础分配逻辑

### 阶段二：队列管理
- 创建队列管理表
- 实现队列API接口
- 开发队列管理界面

### 阶段三：智能分配
- 设备负载监控
- 分配算法优化
- 性能监控和调优

### 阶段四：高级功能
- 任务抢占机制
- 动态重平衡
- 预测性分配

这个设计方案完整解决了任务分配队列的处理问题，并提供了清晰的前端界面展示方案。 