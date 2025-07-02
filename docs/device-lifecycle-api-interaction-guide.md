# 设备生命周期API交互详细指南

## 概述

本文档详细描述设备生命周期各阶段的API交互细节，包括HTTP头文件、请求体格式、响应格式和认证要求。

## 通用HTTP头文件要求

### 基础头文件
```http
Content-Type: application/json
Accept: application/json
User-Agent: OneGo-Client/1.0
X-Request-ID: uuid-generated-by-client
```

### 认证头文件
```http
Authorization: Bearer <jwt-token>
X-Device-ESN: <device-serial-number>  // 设备专用
X-API-Key: <api-key>                  // API密钥认证（可选）
```

### 可选头文件
```http
X-Client-Version: 1.0.0
X-Timestamp: 2025-01-01T12:00:00Z
X-Signature: <request-signature>      // 请求签名验证
```

---

## 1. 设备注册阶段 (Device Registration)

### API端点: POST /api/v1/devices/register

#### 请求示例
```bash
curl -X POST http://localhost:8080/api/v1/devices/register \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "X-Request-ID: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Client-Version: 1.0.0" \
  -d '{
    "device_esn": "DEV-2025-001",
    "name": "生产环境服务器-01",
    "type": "server",
    "ip_address": "192.168.1.100",
    "protocol": "ssh",
    "reg_type": "key",
    "model": "Dell PowerEdge R740",
    "manufacturer": "Dell",
    "port": 22,
    "endpoint": "/api/v1/status",
    "login_username": "admin",
    "login_port": 22,
    "login_public_key": "ssh-rsa AAAAB3NzaC1yc...",
    "metadata": {
      "environment": "production",
      "location": "datacenter-A",
      "department": "IT"
      
    },
    "tags": {
      "critical": "yes",
      "backup": "enabled"
    },
    "created_by": "admin"
  }'
```

#### 成功响应 (201 Created)
```json
{
  "message": "设备注册成功",
  "data": {
    "id": 1,
    "device_esn": "DEV-2025-001",
    "name": "生产环境服务器-01",
    "type": "server",
    "status": "unknown",
    "health_score": 100,
    "reg_time": "2025-01-01T12:00:00Z",
    "created_at": "2025-01-01T12:00:00Z",
    "updated_at": "2025-01-01T12:00:00Z"
  }
}
```

#### 错误响应 (400 Bad Request)
```json
{
  "message": "设备注册失败",
  "error": {
    "code": 3002,
    "message": "设备已存在",
    "details": "设备序列号DEV-2025-001已被注册"
  }
}
```

---

## 2. 设备首次上线阶段 (Device First Online)

### API端点: POST /api/v1/devices/online

#### 请求示例
```bash
curl -X POST http://localhost:8080/api/v1/devices/online \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "X-Device-ESN: DEV-2025-001" \
  -H "X-Request-ID: 550e8400-e29b-41d4-a716-446655440001" \
  -d '{
    "device_esn": "DEV-2025-001",
    "ip_address": "192.168.1.100",
    "system_info": {
      "os": "Ubuntu 20.04",
      "kernel": "5.4.0-90-generic",
      "cpu_cores": 8,
      "memory_total": "32GB",
      "disk_total": "1TB"
    }
  }'
```

#### 成功响应 (200 OK)
```json
{
  "message": "设备上线成功",
  "data": {
    "device_esn": "DEV-2025-001",
    "status": "online",
    "first_online_time": "2025-01-01T12:05:00Z",
    "last_online_time": "2025-01-01T12:05:00Z",
    "ip_address": "192.168.1.100",
    "health_score": 100
  }
}
```

---

## 3. 心跳监控阶段 (Heartbeat Monitoring)

### API端点: POST /api/v1/devices/{id}/heartbeat

#### 请求示例 (设备主动上报)
```bash
curl -X POST http://localhost:8080/api/v1/devices/1/heartbeat \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "X-Device-ESN: DEV-2025-001" \
  -H "X-Request-ID: 550e8400-e29b-41d4-a716-446655440002" \
  -d '{
    "status": "online",
    "metrics": {
      "cpu_usage": 45.2,
      "memory_usage": 68.5,
      "disk_usage": 23.1,
      "network_io": {
        "rx_bytes": 1048576,
        "tx_bytes": 2097152
      },
      "system_load": 1.25,
      "process_count": 156
    },
    "current_task_id": "TASK-2025-001",
    "response_time": 15
  }'
```

#### 成功响应 (200 OK)
```json
{
  "message": "心跳更新成功",
  "data": {
    "device_esn": "DEV-2025-001",
    "heartbeat_time": "2025-01-01T12:10:30Z",
    "status": "online",
    "health_score": 95,
    "next_heartbeat": "2025-01-01T12:11:00Z"
  }
}
```

### API端点: GET /api/v1/devices/{id}/heartbeat

#### 请求示例 (获取心跳历史)
```bash
curl -X GET "http://localhost:8080/api/v1/devices/1/heartbeat?limit=10&hours=24" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "X-Request-ID: 550e8400-e29b-41d4-a716-446655440003"
```

#### 成功响应 (200 OK)
```json
{
  "message": "获取心跳记录成功",
  "data": {
    "device_esn": "DEV-2025-001",
    "latest_heartbeat": "2025-01-01T12:10:30Z",
    "heartbeat_interval": 30,
    "history": [
      {
        "heartbeat_time": "2025-01-01T12:10:30Z",
        "status": "online",
        "metrics": {
          "cpu_usage": 45.2,
          "memory_usage": 68.5
        },
        "response_time": 15
      }
    ],
    "statistics": {
      "total_heartbeats": 2880,
      "avg_response_time": 18.5,
      "uptime_percentage": 99.85
    }
  }
}
```

---

## 4. 命令执行阶段 (Command Execution)

### API端点: POST /api/v1/devices/{id}/command

#### 请求示例 (发送Shell命令)
```bash
curl -X POST http://localhost:8080/api/v1/devices/1/command \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "X-Request-ID: 550e8400-e29b-41d4-a716-446655440004" \
  -d '{
    "command_type": "shell",
    "command_data": {
      "command": "systemctl status nginx",
      "timeout": 30,
      "working_directory": "/home/admin",
      "environment": {
        "PATH": "/usr/bin:/bin"
      }
    },
    "priority": "normal",
    "async": false
  }'
```

#### 成功响应 (201 Created)
```json
{
  "message": "命令发送成功",
  "data": {
    "command_id": "CMD-2025-001",
    "device_esn": "DEV-2025-001",
    "command_type": "shell",
    "status": "pending",
    "sent_time": "2025-01-01T12:15:00Z",
    "estimated_completion": "2025-01-01T12:15:30Z"
  }
}
```

#### 请求示例 (文件传输命令)
```bash
curl -X POST http://localhost:8080/api/v1/devices/1/command \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "X-Request-ID: 550e8400-e29b-41d4-a716-446655440005" \
  -d '{
    "command_type": "file_transfer",
    "command_data": {
      "action": "upload",
      "source_path": "/local/file.txt",
      "destination_path": "/remote/file.txt",
      "checksum": "sha256:abc123...",
      "file_size": 1024
    }
  }'
```

---

## 5. 状态查询阶段 (Status Query)

### API端点: GET /api/v1/devices/{id}/status

#### 请求示例
```bash
curl -X GET http://localhost:8080/api/v1/devices/1/status \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "X-Request-ID: 550e8400-e29b-41d4-a716-446655440006"
```

#### 成功响应 (200 OK)
```json
{
  "message": "获取设备状态成功",
  "data": {
    "device_esn": "DEV-2025-001",
    "status": "online",
    "health_score": 95,
    "last_seen": "2025-01-01T12:10:30Z",
    "uptime_hours": 72.5,
    "online": true,
    "current_tasks": {
      "running": 2,
      "pending": 1,
      "total": 3
    },
    "resource_usage": {
      "cpu_usage": 45.2,
      "memory_usage": 68.5,
      "disk_usage": 23.1
    },
    "network_status": {
      "ip_address": "192.168.1.100",
      "latency": 15,
      "bandwidth_usage": "moderate"
    }
  }
}
```

---

## 6. 日志记录阶段 (Log Management)

### API端点: GET /api/v1/devices/{id}/logs

#### 请求示例
```bash
curl -X GET "http://localhost:8080/api/v1/devices/1/logs?level=ERROR&limit=50&start_time=2025-01-01T00:00:00Z" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "X-Request-ID: 550e8400-e29b-41d4-a716-446655440007"
```

#### 成功响应 (200 OK)
```json
{
  "message": "获取设备日志成功",
  "data": {
    "device_esn": "DEV-2025-001",
    "total_logs": 1250,
    "filtered_logs": 8,
    "logs": [
      {
        "id": 12345,
        "level": "ERROR",
        "category": "system",
        "message": "磁盘空间不足",
        "details": {
          "disk_path": "/var/log",
          "available_space": "1GB",
          "threshold": "5GB"
        },
        "source": "disk_monitor",
        "correlation_id": "ERR-2025-001",
        "log_time": "2025-01-01T11:45:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "size": 50,
      "total_pages": 1
    }
  }
}
```

---

## 7. 设备下线阶段 (Device Offline)

### API端点: POST /api/v1/devices/offline

#### 请求示例 (主动下线)
```bash
curl -X POST http://localhost:8080/api/v1/devices/offline \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "X-Device-ESN: DEV-2025-001" \
  -H "X-Request-ID: 550e8400-e29b-41d4-a716-446655440008" \
  -d '{
    "device_esn": "DEV-2025-001",
    "offline_reason": "planned_maintenance",
    "estimated_downtime": 3600,
    "final_status": {
      "cpu_usage": 12.5,
      "memory_usage": 35.2,
      "running_tasks": 0
    }
  }'
```

#### 成功响应 (200 OK)
```json
{
  "message": "设备下线成功",
  "data": {
    "device_esn": "DEV-2025-001",
    "status": "offline",
    "last_offline_time": "2025-01-01T15:30:00Z",
    "total_online_duration": 12600,
    "session_summary": {
      "tasks_completed": 45,
      "commands_executed": 128,
      "uptime_hours": 3.5,
      "avg_cpu_usage": 42.8
    }
  }
}
```

---

## 8. 设备维护阶段 (Device Maintenance)

### API端点: PUT /api/v1/devices/{id}

#### 请求示例 (进入维护模式)
```bash
curl -X PUT http://localhost:8080/api/v1/devices/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "X-Request-ID: 550e8400-e29b-41d4-a716-446655440009" \
  -d '{
    "status": "maintenance",
    "maintenance_info": {
      "reason": "hardware_upgrade",
      "description": "内存扩容和固件更新",
      "scheduled_duration": 7200,
      "technician": "张工程师",
      "contact": "zhang@company.com"
    },
    "updated_by": "admin"
  }'
```

#### 成功响应 (200 OK)
```json
{
  "message": "设备维护状态更新成功",
  "data": {
    "device_esn": "DEV-2025-001",
    "status": "maintenance",
    "maintenance_start_time": "2025-01-01T16:00:00Z",
    "estimated_completion": "2025-01-01T18:00:00Z",
    "updated_by": "admin",
    "updated_at": "2025-01-01T16:00:00Z"
  }
}
```

---

## 9. 设备注销阶段 (Device Deregistration)

### API端点: DELETE /api/v1/devices/{id}

#### 请求示例
```bash
curl -X DELETE http://localhost:8080/api/v1/devices/1 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "X-Request-ID: 550e8400-e29b-41d4-a716-446655440010" \
  -H "X-Reason: device_retired" \
  -H "X-Confirmation: true"
```

#### 成功响应 (200 OK)
```json
{
  "message": "设备注销成功",
  "data": {
    "device_esn": "DEV-2025-001",
    "status": "deregistered",
    "deregistration_time": "2025-01-01T18:30:00Z",
    "lifecycle_summary": {
      "registration_time": "2025-01-01T12:00:00Z",
      "total_lifetime_hours": 6.5,
      "total_tasks_completed": 45,
      "total_commands_executed": 128,
      "total_heartbeats": 780,
      "final_health_score": 95
    },
    "archived_data": {
      "logs_archived": true,
      "metrics_archived": true,
      "retention_period": "90_days"
    }
  }
}
```

---

## 认证和安全

### JWT Token 格式
```javascript
// Token Payload
{
  "sub": "user_id",
  "iat": 1640995200,
  "exp": 1641081600,
  "scope": ["device:read", "device:write", "device:admin"],
  "device_access": ["DEV-2025-001", "DEV-2025-002"],
  "role": "admin"
}
```

### API密钥认证 (可选)
```http
X-API-Key: ak_live_1234567890abcdef
X-API-Signature: sha256=abc123...
```

### 请求签名算法
```bash
# 签名计算
signature = HMAC-SHA256(api_secret, request_method + request_path + request_body + timestamp)
```

---

## 错误处理

### 标准错误响应格式
```json
{
  "message": "操作失败",
  "error": {
    "code": 4001,
    "message": "设备不存在",
    "details": "设备ID 999 未找到",
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "timestamp": "2025-01-01T12:00:00Z"
  }
}
```

### 常见错误码
| 错误码 | HTTP状态码 | 描述 |
|--------|------------|------|
| 3001 | 404 | 设备不存在 |
| 3002 | 409 | 设备已存在 |
| 3003 | 503 | 设备离线 |
| 3004 | 401 | 设备未授权 |
| 3005 | 409 | 设备忙碌 |
| 1003 | 401 | 未授权访问 |
| 1001 | 400 | 参数错误 |
| 1005 | 429 | 请求频率限制 |

---

## 限流和配额

### 请求频率限制
```http
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1640995260
Retry-After: 60
```

### 配额使用情况
```http
X-Quota-Used: 1250
X-Quota-Limit: 10000
X-Quota-Reset: 2025-02-01T00:00:00Z
```

---

## SDK和工具

### Python SDK示例
```python
from onego_client import OneGoClient

# 初始化客户端
client = OneGoClient(
    base_url="http://localhost:8080",
    token="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
)

# 注册设备
device = client.devices.register(
    device_esn="DEV-2025-001",
    name="生产环境服务器-01",
    type="server",
    ip_address="192.168.1.100"
)

# 设备上线
client.devices.online(device_esn="DEV-2025-001")

# 发送心跳
client.devices.heartbeat(
    device_id=1,
    status="online",
    metrics={"cpu_usage": 45.2}
)
```

### 命令行工具示例
```bash
# 安装OneGo CLI
pip install onego-cli

# 配置认证
onego config set-token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

# 注册设备
onego device register \
  --esn DEV-2025-001 \
  --name "生产环境服务器-01" \
  --type server \
  --ip 192.168.1.100

# 查看设备状态
onego device status DEV-2025-001

# 发送命令
onego device command DEV-2025-001 \
  --type shell \
  --command "systemctl status nginx"
```

---

## 监控和调试

### 请求追踪
```http
X-Request-ID: 550e8400-e29b-41d4-a716-446655440000
X-Trace-ID: trace_12345
X-Span-ID: span_67890
```

### 性能指标
```http
X-Response-Time: 150ms
X-Database-Time: 45ms
X-Cache-Hit: true
```

### 调试信息
```http
X-Debug-Info: {"sql_queries": 3, "cache_hits": 2}
X-Version: v1.2.3
```

通过以上详细的API交互指南，开发者可以完整地实现设备生命周期管理的各个阶段，确保系统的稳定性和可维护性。 