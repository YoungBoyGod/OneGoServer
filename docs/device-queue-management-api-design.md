# 设备队列管理API接口设计

## 概述

针对用户需求："增加设备队列查询的表，需要查询设备上现有的队列情况，还需要考虑手动修改优先级和队列的情况"，设计了完整的设备队列管理系统。

## 核心功能

### 1. 设备队列查询
- 查询指定设备的任务队列
- 支持多种过滤和排序条件
- 提供队列统计信息

### 2. 手动队列管理
- 手动调整任务优先级
- 手动调整队列位置
- 批量队列操作
- 操作历史记录

### 3. 队列配置管理
- 设备队列配置
- 调度策略设置
- 资源限制配置

## API接口设计

### 1. 设备队列查询接口

#### 获取设备队列列表
```http
GET /api/v1/devices/queues
```

**查询参数**:
```typescript
interface DeviceQueueListParams {
  device_id?: number;        // 设备ID过滤
  device_esn?: string;       // 设备编号过滤
  status?: string[];         // 状态过滤: queued, executing, paused, completed, failed
  has_manual_adjustments?: boolean; // 是否有手动调整
  page?: number;             // 页码
  size?: number;             // 页大小
}
```

**响应结构**:
```json
{
  "devices": [
    {
      "device_info": {
        "id": 1,
        "device_esn": "DEV001",
        "name": "生产设备01",
        "type": "production",
        "status": "online"
      },
      "total_queued_tasks": 5,
      "executing_tasks": 1,
      "manual_adjustments": 2,
      "next_task_start_time": "2024-01-15T10:30:00Z",
      "estimated_duration": 3600,
      "last_queue_update": "2024-01-15T10:25:00Z"
    }
  ],
  "summary": {
    "total_devices": 10,
    "devices_with_tasks": 6,
    "total_queued_tasks": 25,
    "total_executing_tasks": 8,
    "total_manual_adjustments": 12,
    "avg_queue_length": 2.5
  }
}
```

#### 获取指定设备的详细队列
```http
GET /api/v1/devices/{deviceId}/queue
```

**查询参数**:
```typescript
interface DeviceQueueParams {
  status?: string[];         // 状态过滤
  priority?: number;         // 优先级过滤
  is_manual_priority?: boolean; // 是否手动调整优先级
  is_manual_position?: boolean; // 是否手动调整位置
  sort_by?: string;          // 排序字段: queue_position, queue_priority, created_at
  sort_order?: string;       // 排序方向: asc, desc
  page?: number;
  size?: number;
}
```

**响应结构**:
```json
{
  "device_info": {
    "id": 1,
    "device_esn": "DEV001",
    "name": "生产设备01",
    "type": "production",
    "status": "online"
  },
  "queue_config": {
    "max_queue_size": 100,
    "max_concurrent_tasks": 2,
    "auto_start_tasks": true,
    "priority_scheduling": true,
    "scheduling_strategy": "priority_first"
  },
  "queue_items": [
    {
      "id": 101,
      "task_id": 501,
      "device_id": 1,
      "device_esn": "DEV001",
      "queue_priority": 1,
      "queue_position": 1,
      "original_priority": 3,
      "is_manual_priority": true,
      "is_manual_position": false,
      "status": "queued",
      "estimated_start_time": "2024-01-15T10:30:00Z",
      "estimated_duration": 1800,
      "last_action": "priority_changed",
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:25:00Z",
      "task": {
        "id": 501,
        "name": "数据备份任务",
        "type": "backup",
        "priority": 3
      }
    }
  ],
  "statistics": {
    "device_id": 1,
    "device_esn": "DEV001",
    "total_tasks": 5,
    "queued_tasks": 4,
    "executing_tasks": 1,
    "manual_priority_tasks": 2,
    "manual_position_tasks": 1,
    "by_status": {
      "queued": 4,
      "executing": 1
    },
    "by_priority": {
      "1": 2,
      "2": 1,
      "3": 2
    },
    "avg_wait_time": 600.5,
    "total_estimated_duration": 7200,
    "next_task_start_time": "2024-01-15T10:30:00Z"
  },
  "pagination": {
    "page": 1,
    "size": 20,
    "total": 5,
    "total_pages": 1
  }
}
```

### 2. 手动队列管理接口

#### 调整任务优先级
```http
PUT /api/v1/devices/{deviceId}/queue/{taskId}/priority
```

**请求体**:
```json
{
  "new_priority": 1,
  "reason": "紧急任务，需要优先执行"
}
```

**响应**:
```json
{
  "success": true,
  "message": "任务优先级调整成功",
  "data": {
    "task_id": 501,
    "old_priority": 3,
    "new_priority": 1,
    "new_position": 1,
    "operation_time": "2024-01-15T10:25:00Z",
    "operation_by": 1001
  }
}
```

#### 调整队列位置
```http
PUT /api/v1/devices/{deviceId}/queue/{taskId}/position
```

**请求体**:
```json
{
  "new_position": 2,
  "reason": "按计划调整执行顺序"
}
```

#### 批量队列操作
```http
POST /api/v1/devices/{deviceId}/queue/batch
```

**请求体**:
```json
{
  "operation": "priority_change",
  "task_ids": [501, 502, 503],
  "new_priority": 2,
  "reason": "批量调整优先级",
  "batch_id": "batch_20240115_001"
}
```

**批量位置调整**:
```json
{
  "operation": "position_change",
  "task_ids": [501, 502, 503],
  "new_positions": [1, 2, 3],
  "reason": "重新排序队列",
  "batch_id": "batch_20240115_002"
}
```

#### 队列操作历史
```http
GET /api/v1/devices/{deviceId}/queue/history
```

**查询参数**:
```typescript
interface QueueHistoryParams {
  task_id?: number;          // 任务ID过滤
  operation_type?: string[]; // 操作类型过滤
  operation_by?: number;     // 操作人过滤
  start_time?: string;       // 开始时间
  end_time?: string;         // 结束时间
  batch_id?: string;         // 批量操作ID
  page?: number;
  size?: number;
}
```

**响应结构**:
```json
{
  "history": [
    {
      "id": 1001,
      "device_id": 1,
      "device_esn": "DEV001",
      "task_id": 501,
      "operation_type": "priority_change",
      "operation_by": 1001,
      "operation_time": "2024-01-15T10:25:00Z",
      "old_priority": 3,
      "new_priority": 1,
      "reason": "紧急任务，需要优先执行",
      "operation_source": "manual",
      "batch_id": null,
      "is_batch_operation": false
    }
  ],
  "pagination": {
    "page": 1,
    "size": 20,
    "total": 15,
    "total_pages": 1
  }
}
```

### 3. 队列配置管理接口

#### 获取设备队列配置
```http
GET /api/v1/devices/{deviceId}/queue/config
```

#### 更新设备队列配置
```http
PUT /api/v1/devices/{deviceId}/queue/config
```

**请求体**:
```json
{
  "max_queue_size": 50,
  "max_concurrent_tasks": 3,
  "auto_start_tasks": true,
  "priority_scheduling": true,
  "scheduling_strategy": "priority_first",
  "load_balancing": true,
  "work_start_time": "08:00:00",
  "work_end_time": "22:00:00",
  "timezone": "Asia/Shanghai",
  "max_cpu_usage": 75.0,
  "max_memory_usage": 85.0,
  "min_free_disk": 2147483648,
  "notify_on_completion": false,
  "notify_on_failure": true,
  "notification_webhook": "https://api.example.com/webhook"
}
```

### 4. 队列控制接口

#### 暂停设备队列
```http
POST /api/v1/devices/{deviceId}/queue/pause
```

#### 恢复设备队列
```http
POST /api/v1/devices/{deviceId}/queue/resume
```

#### 清空设备队列
```http
POST /api/v1/devices/{deviceId}/queue/clear
```

**请求体**:
```json
{
  "reason": "设备维护，清空队列",
  "preserve_executing": true,  // 是否保留正在执行的任务
  "notify_users": true         // 是否通知相关用户
}
```

#### 重新平衡队列
```http
POST /api/v1/devices/{deviceId}/queue/rebalance
```

**请求体**:
```json
{
  "strategy": "priority_first", // 重平衡策略
  "reset_manual_adjustments": false, // 是否重置手动调整
  "reason": "系统优化队列顺序"
}
```

## 前端界面设计

### 1. 设备队列概览页面

```vue
<template>
  <div class="device-queue-overview">
    <!-- 全局统计卡片 -->
    <div class="stats-cards">
      <StatCard 
        title="设备总数" 
        :value="summary.total_devices" 
        icon="device"
      />
      <StatCard 
        title="有任务设备" 
        :value="summary.devices_with_tasks" 
        icon="active-device"
      />
      <StatCard 
        title="队列中任务" 
        :value="summary.total_queued_tasks" 
        icon="queue"
      />
      <StatCard 
        title="执行中任务" 
        :value="summary.total_executing_tasks" 
        icon="running"
      />
    </div>

    <!-- 设备列表 -->
    <div class="device-list">
      <DeviceQueueCard 
        v-for="device in devices" 
        :key="device.device_info.id"
        :device="device"
        @view-detail="showDeviceQueueDetail"
        @manage-queue="showQueueManagement"
      />
    </div>
  </div>
</template>
```

### 2. 设备队列详情页面

```vue
<template>
  <div class="device-queue-detail">
    <!-- 设备信息 -->
    <DeviceInfoCard :device="deviceInfo" :config="queueConfig" />
    
    <!-- 队列操作工具栏 -->
    <div class="queue-toolbar">
      <Button @click="showBatchOperation">批量操作</Button>
      <Button @click="pauseQueue">暂停队列</Button>
      <Button @click="rebalanceQueue">重新平衡</Button>
      <Button @click="clearQueue">清空队列</Button>
    </div>

    <!-- 队列任务列表 -->
    <div class="queue-task-list">
      <QueueTaskItem 
        v-for="(task, index) in queueItems" 
        :key="task.id"
        :task="task"
        :position="index + 1"
        :draggable="true"
        @drag-end="handlePositionChange"
        @priority-change="handlePriorityChange"
        @remove="handleRemoveTask"
      />
    </div>

    <!-- 分页 -->
    <Pagination 
      :current="pagination.page"
      :total="pagination.total"
      :page-size="pagination.size"
      @change="handlePageChange"
    />

    <!-- 队列统计 -->
    <QueueStatistics :statistics="statistics" />
  </div>
</template>
```

### 3. 队列任务项组件

```vue
<template>
  <div class="queue-task-item" :class="taskStatusClass">
    <!-- 拖拽手柄 -->
    <div class="drag-handle">
      <Icon name="drag" />
    </div>

    <!-- 队列位置 -->
    <div class="queue-position">
      <span class="position-number">{{ position }}</span>
      <Input 
        v-if="editingPosition" 
        v-model="newPosition"
        size="small"
        @blur="savePosition"
        @keyup.enter="savePosition"
      />
      <Button 
        v-else
        size="small" 
        type="text"
        @click="editPosition"
      >
        调整
      </Button>
    </div>

    <!-- 任务信息 -->
    <div class="task-info">
      <div class="task-name">{{ task.task.name }}</div>
      <div class="task-meta">
        <Tag :color="taskTypeColor">{{ task.task.type }}</Tag>
        <span>预计: {{ formatDuration(task.estimated_duration) }}</span>
        <span v-if="task.estimated_start_time">
          开始: {{ formatTime(task.estimated_start_time) }}
        </span>
      </div>
    </div>

    <!-- 优先级 -->
    <div class="priority-section">
      <PrioritySelector 
        :value="task.queue_priority"
        :original="task.original_priority"
        :is-manual="task.is_manual_priority"
        @change="handlePriorityChange"
      />
    </div>

    <!-- 状态 -->
    <div class="status-section">
      <TaskStatus :status="task.status" />
      <div v-if="task.is_manual_priority || task.is_manual_position" class="manual-indicators">
        <Tag v-if="task.is_manual_priority" color="orange">手动优先级</Tag>
        <Tag v-if="task.is_manual_position" color="blue">手动位置</Tag>
      </div>
    </div>

    <!-- 操作按钮 -->
    <div class="actions">
      <Dropdown>
        <Button type="text">
          <Icon name="more" />
        </Button>
        <template #overlay>
          <Menu>
            <MenuItem @click="viewTaskDetail">查看详情</MenuItem>
            <MenuItem @click="editTask">编辑任务</MenuItem>
            <MenuItem @click="pauseTask">暂停任务</MenuItem>
            <MenuItem @click="removeFromQueue">移出队列</MenuItem>
          </Menu>
        </template>
      </Dropdown>
    </div>
  </div>
</template>
```

### 4. 批量操作对话框

```vue
<template>
  <Modal v-model:visible="visible" title="批量队列操作" width="600px">
    <div class="batch-operation-form">
      <!-- 选择操作类型 -->
      <div class="operation-type">
        <Radio.Group v-model:value="operationType">
          <Radio value="priority_change">调整优先级</Radio>
          <Radio value="position_change">调整位置</Radio>
          <Radio value="pause">暂停任务</Radio>
          <Radio value="cancel">取消任务</Radio>
        </Radio.Group>
      </div>

      <!-- 选择任务 -->
      <div class="task-selection">
        <h4>选择任务</h4>
        <CheckboxGroup v-model:value="selectedTasks">
          <Checkbox 
            v-for="task in availableTasks" 
            :key="task.id"
            :value="task.id"
          >
            {{ task.task.name }} (位置: {{ task.queue_position }})
          </Checkbox>
        </CheckboxGroup>
      </div>

      <!-- 操作参数 -->
      <div v-if="operationType === 'priority_change'" class="operation-params">
        <Form.Item label="新优先级">
          <Select v-model:value="newPriority">
            <Option v-for="i in 10" :key="i" :value="i">{{ i }}</Option>
          </Select>
        </Form.Item>
      </div>

      <div v-if="operationType === 'position_change'" class="operation-params">
        <Form.Item label="新位置">
          <span>拖拽下方任务列表重新排序</span>
          <DraggableList 
            v-model:items="taskOrder" 
            @change="handleOrderChange"
          />
        </Form.Item>
      </div>

      <!-- 操作原因 -->
      <Form.Item label="操作原因">
        <Input.TextArea 
          v-model:value="reason" 
          placeholder="请输入操作原因"
          :rows="3"
        />
      </Form.Item>
    </div>

    <template #footer>
      <Button @click="visible = false">取消</Button>
      <Button type="primary" @click="executeBatchOperation">执行操作</Button>
    </template>
  </Modal>
</template>
```

## 实施计划

### 阶段一：基础设施 (已完成)
- ✅ 数据库表设计
- ✅ 数据库迁移脚本
- 🔄 GORM模型定义 (需要修复类型冲突)

### 阶段二：后端API (2-3天)
- Repository层实现
- Service层业务逻辑
- Controller层API接口
- 权限和验证

### 阶段三：前端界面 (3-4天)
- 设备队列概览页面
- 队列详情管理页面
- 批量操作功能
- 拖拽排序功能

### 阶段四：高级功能 (2-3天)
- 实时队列状态更新
- 队列操作通知
- 性能监控和优化
- 操作日志审计

这个设计完整解决了设备队列查询和手动管理的需求，提供了灵活强大的队列管理能力。 