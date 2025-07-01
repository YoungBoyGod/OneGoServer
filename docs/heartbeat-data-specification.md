# 心跳数据传输规范

## 🎯 心跳机制概述

心跳(Heartbeat)是设备与服务器之间的**定期通信机制**，用于：
- 🔍 **状态监控** - 确认设备在线状态
- 📊 **数据采集** - 收集设备运行指标  
- ⚙️ **配置同步** - 下发配置和命令
- 🚨 **故障检测** - 及时发现设备异常

---

## 📤 Client设备 → Server服务器 (心跳数据)

### 1. 基础心跳数据包
```json
{
  "heartbeat_time": "2025-01-02T10:30:00Z",
  "device_id": "TEMP_SENSOR_001",
  "sequence": 12345,
  "status": "online",
  "ip_address": "192.168.1.100",
  "response_time": 15
}
```

### 2. 完整心跳数据包 (温度传感器)
```json
{
  // 基础信息
  "heartbeat_time": "2025-01-02T10:30:00Z",
  "device_id": "TEMP_SENSOR_001", 
  "device_name": "办公室温度传感器",
  "sequence": 12345,
  "status": "online",
  "ip_address": "192.168.1.100",
  "response_time": 15,
  
  // 系统指标
  "system_info": {
    "cpu_usage": 25.5,
    "memory_usage": 45.2,
    "disk_usage": 60.0,
    "uptime": 86400,
    "temperature": 42.5,
    "battery_level": 85
  },
  
  // 业务指标
  "metrics": {
    "current_temperature": 23.5,
    "humidity": 65.2,
    "sample_count": 1440,
    "last_sample_time": "2025-01-02T10:29:30Z",
    "calibration_status": "ok"
  },
  
  // 配置版本
  "config_version": "v1.2.3",
  "firmware_version": "2.1.0",
  
  // 错误和告警
  "errors": [],
  "warnings": [
    {
      "code": "TEMP_HIGH", 
      "message": "温度接近上限",
      "level": "warning",
      "timestamp": "2025-01-02T10:25:00Z"
    }
  ]
}
```

### 3. 智能摄像头心跳数据
```json
{
  "heartbeat_time": "2025-01-02T10:30:00Z",
  "device_id": "CAMERA_A1",
  "status": "online",
  "ip_address": "192.168.1.200",
  
  "system_info": {
    "cpu_usage": 65.0,
    "memory_usage": 78.5,
    "disk_usage": 45.0,
    "temperature": 38.2
  },
  
  "metrics": {
    "recording_status": "active",
    "video_quality": "1080p",
    "frame_rate": 30,
    "motion_detected": false,
    "storage_remaining": "500GB",
    "night_vision": false
  },
  
  "warnings": []
}
```

### 4. 工业控制器心跳数据
```json
{
  "heartbeat_time": "2025-01-02T10:30:00Z",
  "device_id": "PLC_001",
  "status": "running",
  "ip_address": "10.0.1.50",
  
  "system_info": {
    "cpu_usage": 15.5,
    "memory_usage": 32.1,
    "uptime": 2592000
  },
  
  "metrics": {
    "production_line_status": "running",
    "current_speed": 150,
    "target_speed": 150,
    "products_count": 1250,
    "quality_rate": 99.2,
    "last_maintenance": "2024-12-15T08:00:00Z"
  },
  
  "io_status": {
    "digital_inputs": [1, 0, 1, 1, 0, 0, 1, 1],
    "digital_outputs": [1, 1, 0, 1, 0, 0, 0, 1],
    "analog_inputs": [4.2, 3.8, 2.1, 0.0],
    "analog_outputs": [5.0, 3.3]
  }
}
```

---

## 📥 Server服务器 → Client设备 (响应数据)

### 1. 基础确认响应
```json
{
  "timestamp": "2025-01-02T10:30:01Z",
  "status": "ack",
  "sequence": 12345,
  "server_time": "2025-01-02T10:30:01Z",
  "next_heartbeat": 30
}
```

### 2. 配置更新响应
```json
{
  "timestamp": "2025-01-02T10:30:01Z",
  "status": "ack",
  "sequence": 12345,
  
  // 配置更新
  "config_update": {
    "sample_rate": 60,
    "alert_threshold": 25.0,
    "reporting_interval": 300
  },
  
  // 时间同步
  "time_sync": {
    "server_time": "2025-01-02T10:30:01Z",
    "timezone": "Asia/Shanghai"
  }
}
```

### 3. 命令下发响应
```json
{
  "timestamp": "2025-01-02T10:30:01Z", 
  "status": "ack",
  "sequence": 12345,
  
  // 控制命令
  "commands": [
    {
      "command_id": "cmd_001",
      "type": "restart",
      "execute_time": "2025-01-02T11:00:00Z",
      "parameters": {}
    },
    {
      "command_id": "cmd_002", 
      "type": "update_config",
      "parameters": {
        "sample_rate": 30
      }
    }
  ]
}
```

### 4. 固件更新响应
```json
{
  "timestamp": "2025-01-02T10:30:01Z",
  "status": "ack", 
  "sequence": 12345,
  
  // 固件更新
  "firmware_update": {
    "available": true,
    "version": "2.2.0",
    "download_url": "https://firmware.onego.com/v2.2.0.bin",
    "checksum": "sha256:abc123...",
    "mandatory": false,
    "schedule_time": "2025-01-02T02:00:00Z"
  }
}
```

---

## 🔄 心跳频率和策略

### 1. 动态心跳频率
```json
// 设备状态不同，心跳频率不同
{
  "normal": 30,        // 正常状态：30秒
  "warning": 10,       // 告警状态：10秒  
  "error": 5,          // 错误状态：5秒
  "maintenance": 60,   // 维护状态：60秒
  "offline": 0         // 离线状态：停止心跳
}
```

### 2. 智能心跳策略
```json
{
  "adaptive_interval": {
    "base_interval": 30,
    "max_interval": 300,
    "min_interval": 5,
    "factors": {
      "network_quality": 0.8,  // 网络质量影响因子
      "device_load": 1.2,      // 设备负载影响因子  
      "error_rate": 2.0        // 错误率影响因子
    }
  }
}
```

---

## 💾 在OneGo系统中的存储

### 1. device_heartbeats表结构
```sql
CREATE TABLE device_heartbeats (
    id              BIGSERIAL PRIMARY KEY,
    device_id       BIGINT NOT NULL,
    heartbeat_time  TIMESTAMP NOT NULL,
    status          VARCHAR(20) NOT NULL,
    metrics         JSONB,           -- 完整指标数据
    system_info     JSONB,           -- 系统信息
    ip_address      INET,
    response_time   INTEGER,
    FOREIGN KEY (device_id) REFERENCES devices(id)
);
```

### 2. JSONB数据示例
```sql
-- 插入心跳记录
INSERT INTO device_heartbeats (
    device_id, heartbeat_time, status, metrics, system_info, ip_address, response_time
) VALUES (
    1, 
    '2025-01-02 10:30:00',
    'online',
    '{
        "current_temperature": 23.5,
        "humidity": 65.2,
        "sample_count": 1440
    }'::jsonb,
    '{
        "cpu_usage": 25.5,
        "memory_usage": 45.2,
        "uptime": 86400
    }'::jsonb,
    '192.168.1.100',
    15
);
```

---

## 🛠️ Repository层处理

### 1. 创建心跳记录
```go
func (r *DeviceRepository) CreateHeartbeat(ctx context.Context, heartbeat *device.DeviceHeartbeat) error {
    // 自动更新设备最后活跃时间
    if err := r.db.WithContext(ctx).
        Model(&device.Device{}).
        Where("id = ?", heartbeat.DeviceID).
        Update("last_seen", heartbeat.HeartbeatTime).Error; err != nil {
        return err
    }
    
    // 创建心跳记录
    return r.db.WithContext(ctx).Create(heartbeat).Error
}
```

### 2. 分析心跳数据
```go
func (r *DeviceRepository) AnalyzeHeartbeat(heartbeat *device.DeviceHeartbeat) {
    // 检查CPU使用率
    if cpuUsage, ok := heartbeat.SystemInfo["cpu_usage"].(float64); ok {
        if cpuUsage > 90.0 {
            r.createAlert("HIGH_CPU", heartbeat.DeviceID)
        }
    }
    
    // 检查内存使用率
    if memUsage, ok := heartbeat.SystemInfo["memory_usage"].(float64); ok {
        if memUsage > 85.0 {
            r.createAlert("HIGH_MEMORY", heartbeat.DeviceID)
        }
    }
    
    // 检查业务指标
    if metrics := heartbeat.Metrics; metrics != nil {
        r.analyzeBusinessMetrics(heartbeat.DeviceID, metrics)
    }
}
```

---

## 📊 心跳数据应用场景

### 1. 设备健康度评估
```go
func CalculateHealthScore(heartbeats []DeviceHeartbeat) int {
    score := 100
    
    for _, hb := range heartbeats {
        // CPU使用率影响
        if cpu, ok := hb.SystemInfo["cpu_usage"].(float64); ok {
            if cpu > 80 { score -= 10 }
        }
        
        // 响应时间影响
        if hb.ResponseTime > 1000 { score -= 15 }
        
        // 错误数量影响
        if len(hb.Errors) > 0 { score -= 20 }
    }
    
    if score < 0 { score = 0 }
    return score
}
```

### 2. 异常检测和告警
```go
func DetectAnomalies(deviceID int64, heartbeat *DeviceHeartbeat) []Alert {
    var alerts []Alert
    
    // 响应时间异常
    if heartbeat.ResponseTime > 5000 {
        alerts = append(alerts, Alert{
            Type: "SLOW_RESPONSE",
            Message: "设备响应时间过长",
            Level: "warning",
        })
    }
    
    // 系统资源异常
    if sysInfo := heartbeat.SystemInfo; sysInfo != nil {
        if temp, ok := sysInfo["temperature"].(float64); ok && temp > 70 {
            alerts = append(alerts, Alert{
                Type: "HIGH_TEMPERATURE", 
                Message: "设备温度过高",
                Level: "critical",
            })
        }
    }
    
    return alerts
}
```

### 3. 预测性维护
```go
func PredictMaintenance(deviceID int64) *MaintenancePrediction {
    // 分析历史心跳数据
    heartbeats := getRecentHeartbeats(deviceID, 30) // 最近30天
    
    var avgCPU, avgTemp float64
    var errorCount int
    
    for _, hb := range heartbeats {
        if cpu, ok := hb.SystemInfo["cpu_usage"].(float64); ok {
            avgCPU += cpu
        }
        if temp, ok := hb.SystemInfo["temperature"].(float64); ok {
            avgTemp += temp
        }
        errorCount += len(hb.Errors)
    }
    
    // 基于趋势预测维护需求
    riskLevel := calculateRiskLevel(avgCPU, avgTemp, errorCount)
    
    return &MaintenancePrediction{
        DeviceID: deviceID,
        RiskLevel: riskLevel,
        RecommendedDate: calculateMaintenanceDate(riskLevel),
        Reasons: generateMaintenanceReasons(avgCPU, avgTemp, errorCount),
    }
}
```

---

## 🎯 心跳数据总结

### 📤 设备发送的核心数据
- ✅ **基础状态** - 设备ID、时间戳、在线状态、IP地址
- ✅ **系统指标** - CPU/内存/磁盘使用率、温度、电量
- ✅ **业务数据** - 传感器读数、执行状态、配置版本
- ✅ **错误告警** - 异常代码、告警级别、故障信息

### 📥 服务器响应的核心数据  
- ✅ **确认响应** - 接收确认、时间同步、下次心跳间隔
- ✅ **配置更新** - 参数配置、采样频率、阈值设置
- ✅ **命令下发** - 控制指令、任务分配、重启命令
- ✅ **固件管理** - 版本更新、下载链接、升级调度

通过这套完整的心跳机制，OneGo服务器可以实现设备的**实时监控、远程控制、故障预警、预测性维护**等全面的IoT设备管理功能！🚀 