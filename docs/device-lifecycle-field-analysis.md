# 设备生命周期字段交互分析

## 概述

本文档详细分析设备从注册到注销的整个生命周期中，各个阶段需要交互的数据字段，为系统实现提供完整的数据交互指南。

## 设备生命周期阶段

### 1. 设备注册阶段 (Registration)

#### 1.1 必需字段
```json
{
  "device_esn": "设备序列号 - 唯一标识",
  "name": "设备名称",
  "type": "设备类型 (普通pc/服务器等)",
  "ip_address": "设备IP地址",
  "protocol": "通信协议 (ssh/http)",
  "reg_type": "认证类型 (密码/密钥)",
  "created_by": "创建者"
}
```

#### 1.2 可选字段
```json
{
  "model": "设备型号",
  "manufacturer": "设备厂商",
  "port": "设备端口",
  "endpoint": "设备端点",
  "login_username": "登录用户名",
  "login_port": "登录端口",
  "login_public_key": "公钥",
  "reg_token": "注册token",
  "metadata": "设备元数据",
  "tags": "设备标签"
}
```

#### 1.3 系统自动生成字段
```json
{
  "id": "系统内部ID",
  "status": "unknown",
  "health_score": 100,
  "reg_time": "注册时间",
  "created_at": "创建时间",
  "updated_at": "更新时间"
}
```

### 2. 设备首次上线阶段 (First Online)

#### 2.1 状态更新字段
```json
{
  "status": "online",
  "first_online_time": "首次上线时间",
  "last_online_time": "最后上线时间",
  "ip_address": "当前IP地址(可能变化)"
}
```

#### 2.2 心跳初始化字段
```json
{
  "device_esn": "设备序列号",
  "heartbeat_time": "心跳时间",
  "status": "online",
  "metrics": "设备指标数据",
  "ip_address": "设备IP",
  "response_time": "响应时间"
}
```

### 3. 设备运行阶段 (Runtime)

#### 3.1 心跳数据字段
```json
{
  "device_esn": "设备序列号",
  "heartbeat_time": "心跳时间",
  "status": "设备状态",
  "metrics": {
    "cpu_usage": "CPU使用率",
    "memory_usage": "内存使用率",
    "disk_usage": "磁盘使用率",
    "network_io": "网络IO",
    "system_load": "系统负载"
  },
  "current_task_id": "当前执行任务ID",
  "response_time": "响应时间"
}
```

#### 3.2 任务执行相关字段
```json
{
  "total_tasks": "总任务数",
  "total_success_tasks": "成功任务数",
  "total_failed_tasks": "失败任务数",
  "total_canceled_tasks": "取消任务数",
  "total_pending_tasks": "待执行任务数",
  "total_running_tasks": "执行中任务数",
  "total_completed_tasks": "已完成任务数"
}
```

#### 3.3 日志记录字段
```json
{
  "device_esn": "设备序列号",
  "level": "日志级别 (DEBUG/INFO/WARN/ERROR)",
  "category": "日志分类",
  "message": "日志消息",
  "details": "详细信息",
  "source": "日志来源",
  "correlation_id": "关联ID",
  "log_time": "日志时间"
}
```

#### 3.4 命令执行字段
```json
{
  "device_esn": "设备序列号",
  "command_id": "命令ID",
  "command_type": "命令类型",
  "command_data": "命令数据",
  "status": "执行状态",
  "sent_time": "发送时间",
  "executed_time": "执行时间",
  "completed_time": "完成时间",
  "response_data": "响应数据",
  "error_message": "错误信息"
}
```

### 4. 设备状态监控阶段 (Monitoring)

#### 4.1 健康状态字段
```json
{
  "health_score": "健康度评分 (0-100)",
  "status": "设备状态",
  "last_online_time": "最后上线时间",
  "last_offline_time": "最后下线时间",
  "total_heartbeats": "总心跳数"
}
```

#### 4.2 统计信息字段
```json
{
  "uptime_hours": "在线时长",
  "total_online_duration": "总在线时长",
  "total_offline_duration": "总离线时长",
  "total_alerts": "总告警数"
}
```

### 5. 设备下线阶段 (Offline)

#### 5.1 状态更新字段
```json
{
  "status": "offline",
  "last_offline_time": "下线时间",
  "total_online_duration": "累计在线时长",
  "updated_at": "更新时间"
}
```

### 6. 设备维护阶段 (Maintenance)

#### 6.1 维护状态字段
```json
{
  "status": "maintenance",
  "updated_by": "维护操作者",
  "updated_at": "维护时间"
}
```

### 7. 设备注销阶段 (Deregistration)

#### 7.1 注销记录字段
```json
{
  "status": "deregistered",
  "updated_by": "注销操作者",
  "updated_at": "注销时间"
}
```

## 字段交互矩阵

| 生命周期阶段 | 必需字段 | 可选字段 | 系统字段 | API端点 |
|-------------|---------|---------|---------|---------|
| 注册 | device_esn, name, type | model, manufacturer | id, status | POST /devices |
| 上线 | device_esn | ip_address | first_online_time | POST /devices/online |
| 心跳 | device_esn, status | metrics, current_task_id | heartbeat_time | POST /devices/{id}/heartbeat |
| 命令执行 | command_type, command_data | - | command_id, status | POST /devices/{id}/command |
| 状态查询 | device_esn | - | last_seen, health_score | GET /devices/{id}/status |
| 日志记录 | level, message | category, details | log_time | POST /devices/{id}/logs |
| 下线 | device_esn | - | last_offline_time | POST /devices/offline |
| 维护 | device_esn | - | updated_by | PUT /devices/{id} |
| 注销 | device_esn | - | updated_at | DELETE /devices/{id} |

## 数据流向图

```
设备注册 → 设备上线 → 心跳监控 → 任务执行 → 状态更新 → 设备下线 → 设备维护 → 设备注销
    ↓         ↓         ↓         ↓         ↓         ↓         ↓         ↓
基础信息   连接信息   运行指标   执行结果   状态变更   离线统计   维护记录   注销记录
```

## 关键字段说明

### 设备标识字段
- **device_esn**: 设备序列号，全局唯一标识符，在整个生命周期中不变
- **id**: 数据库内部ID，自增主键

### 状态字段
- **status**: 设备当前状态 (unknown/online/offline/maintenance/error)
- **health_score**: 健康度评分，0-100范围，综合评估设备健康状况

### 时间字段
- **reg_time**: 注册时间
- **first_online_time**: 首次上线时间
- **last_online_time**: 最后上线时间
- **last_offline_time**: 最后下线时间
- **created_at/updated_at**: 审计时间

### 统计字段
- **total_*****: 各类统计计数器
- **uptime_hours**: 累计在线时长
- **total_heartbeats**: 总心跳数

## 字段验证规则

### 必需字段验证
- device_esn: 非空，唯一，长度1-100
- name: 非空，长度1-255
- type: 非空，枚举值验证

### 格式验证
- ip_address: IPv4/IPv6格式验证
- port: 1-65535范围验证
- email: 邮箱格式验证

### 业务规则验证
- health_score: 0-100范围
- status: 状态转换规则验证
- 时间字段: 逻辑顺序验证

## 数据安全考虑

### 敏感字段加密
- login_public_key: 公钥信息
- reg_token: 注册令牌
- credentials: 认证凭据

### 访问控制
- 管理员：所有字段读写权限
- 设备：部分字段读写权限
- 普通用户：只读权限

## 性能优化建议

### 索引优化
- device_esn: 唯一索引
- status + type: 复合索引
- heartbeat_time: 时间索引

### 数据归档
- 心跳数据：保留30天
- 日志数据：保留90天
- 统计数据：永久保留

## 总结

设备生命周期涉及70+个字段，涵盖8个主要阶段。通过合理的字段设计和数据流管理，可以实现完整的设备生命周期管理，支撑IoT设备的全生命周期运维。 