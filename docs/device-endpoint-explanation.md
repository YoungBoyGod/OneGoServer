# 设备Endpoint字段详解

## 🎯 Endpoint字段的作用

`endpoint` 字段存储设备提供的**API接口路径**，是OneGo服务器与IoT设备进行双向通信的关键信息。

## 📡 核心功能

### 1. 设备通信地址
```sql
-- 设备表中的endpoint字段
endpoint VARCHAR(500) -- 设备API端点路径
```

**作用**: 告诉OneGo服务器如何访问这个设备的API接口

### 2. 完整通信URL构建
```
完整URL = protocol://ip_address:port + endpoint
```

**示例**:
- 协议: `HTTP`
- IP地址: `192.168.1.100`  
- 端口: `8080`
- Endpoint: `/api/v1/sensor`
- **完整URL**: `http://192.168.1.100:8080/api/v1/sensor`

---

## 🛠️ 实际应用场景

### 场景1: 智能传感器设备
```json
{
  "device_id": "TEMP_SENSOR_001",
  "name": "办公室温度传感器",
  "type": "sensor",
  "ip_address": "192.168.1.100",
  "port": 8080,
  "protocol": "HTTP",
  "endpoint": "/api/v1/temperature"
}
```

**通信示例**:
```bash
# 获取温度数据
GET http://192.168.1.100:8080/api/v1/temperature/current

# 设置采样频率
POST http://192.168.1.100:8080/api/v1/temperature/config
{
  "sample_rate": 30
}
```

### 场景2: 智能摄像头
```json
{
  "device_id": "CAMERA_A1", 
  "name": "大门监控摄像头",
  "type": "camera",
  "ip_address": "192.168.1.200",
  "port": 80,
  "protocol": "HTTP", 
  "endpoint": "/api/v2/camera"
}
```

**通信示例**:
```bash
# 获取实时画面
GET http://192.168.1.200:80/api/v2/camera/stream

# 控制摄像头转向
POST http://192.168.1.200:80/api/v2/camera/control
{
  "action": "rotate",
  "angle": 45
}
```

### 场景3: 工业控制器
```json
{
  "device_id": "PLC_001",
  "name": "生产线控制器", 
  "type": "actuator",
  "ip_address": "10.0.1.50",
  "port": 502,
  "protocol": "HTTP",
  "endpoint": "/api/plc"
}
```

**通信示例**:
```bash
# 启动生产线
POST http://10.0.1.50:502/api/plc/start

# 获取运行状态  
GET http://10.0.1.50:502/api/plc/status
```

---

## 🔄 不同协议的Endpoint格式

### 1. HTTP REST API
```
endpoint: "/api/v1/device"
endpoint: "/sensors/temperature" 
endpoint: "/control/motor"
```

### 2. WebSocket连接
```
endpoint: "/ws/realtime"
endpoint: "/websocket/data"
endpoint: "/ws/commands"
```

### 3. MQTT主题路径
```
endpoint: "/devices/sensor001/data"
endpoint: "/commands/actuator002"  
endpoint: "/status/gateway"
```

### 4. 自定义协议
```
endpoint: "/modbus/read"
endpoint: "/custom/protocol/v1"
endpoint: "/proprietary/api"
```

---

## 💡 在OneGo系统中的使用

### 1. 发送设备命令
```go
// DeviceService 中使用endpoint发送命令
func (s *DeviceService) SendCommand(deviceID string, command *Command) error {
    device, err := s.deviceRepo.GetByDeviceID(ctx, deviceID)
    if err != nil {
        return err
    }
    
    // 构建完整URL
    url := fmt.Sprintf("%s://%s:%d%s/command", 
        device.Protocol, 
        device.IPAddress, 
        device.Port, 
        device.Endpoint)
    
    // 发送HTTP请求
    return s.httpClient.Post(url, command)
}
```

### 2. 获取设备数据
```go
// 定期从设备获取数据
func (s *DeviceService) CollectDeviceData(deviceID string) (*DeviceData, error) {
    device, err := s.deviceRepo.GetByDeviceID(ctx, deviceID)
    if err != nil {
        return nil, err
    }
    
    // 构建数据接口URL
    url := fmt.Sprintf("%s://%s:%d%s/data", 
        device.Protocol,
        device.IPAddress, 
        device.Port,
        device.Endpoint)
    
    // 获取设备数据
    return s.httpClient.Get(url)
}
```

### 3. 健康检查
```go
// 检查设备是否在线
func (s *DeviceService) HealthCheck(deviceID string) error {
    device, err := s.deviceRepo.GetByDeviceID(ctx, deviceID)
    if err != nil {
        return err
    }
    
    // 构建健康检查URL
    url := fmt.Sprintf("%s://%s:%d%s/health", 
        device.Protocol,
        device.IPAddress,
        device.Port, 
        device.Endpoint)
    
    // 执行健康检查
    return s.httpClient.Ping(url)
}
```

---

## 🔧 Endpoint配置最佳实践

### 1. 标准化路径设计
```
推荐格式: /api/version/resource
示例:
- /api/v1/sensor      # 传感器接口
- /api/v1/camera      # 摄像头接口  
- /api/v1/actuator    # 执行器接口
- /api/v2/gateway     # 网关接口(新版本)
```

### 2. 版本控制
```
支持API版本演进:
- /api/v1/device  # 旧版本API
- /api/v2/device  # 新版本API
- /api/beta/device # 测试版API
```

### 3. 功能模块化
```
按功能组织接口:
- /api/v1/data       # 数据接口
- /api/v1/control    # 控制接口
- /api/v1/config     # 配置接口
- /api/v1/status     # 状态接口
```

### 4. 安全考虑
```
支持认证和加密:
- /secure/api/v1/device    # 安全接口
- /auth/api/v1/device      # 需要认证
- /ssl/api/v1/device       # SSL加密
```

---

## 📊 数据库存储示例

```sql
-- 实际数据库记录示例
INSERT INTO devices (device_id, name, type, ip_address, port, protocol, endpoint) VALUES
('TEMP_001', '温度传感器1', 'sensor', '192.168.1.100', 8080, 'HTTP', '/api/v1/temperature'),
('CAM_A1', '大门摄像头', 'camera', '192.168.1.200', 80, 'HTTP', '/api/v2/camera'),
('MOTOR_01', '电机控制器', 'actuator', '10.0.1.50', 502, 'HTTP', '/api/plc/motor'),
('GATEWAY_1', '边缘网关', 'gateway', '192.168.1.1', 8883, 'MQTT', '/devices/gateway'),
('SENSOR_WS', 'WebSocket传感器', 'sensor', '192.168.1.150', 3000, 'WS', '/ws/realtime');
```

---

## 🎯 总结

`endpoint` 字段是IoT设备管理的核心组件，它：

✅ **定义通信路径** - 告诉系统如何访问设备API  
✅ **支持多协议** - HTTP、WebSocket、MQTT等  
✅ **版本控制** - 支持API版本演进  
✅ **模块化管理** - 按功能组织接口  
✅ **扩展性强** - 适应不同类型设备  

通过合理配置endpoint，OneGo服务器可以与各种IoT设备建立标准化的通信连接，实现统一的设备管理和控制！🚀 