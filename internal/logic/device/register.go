package device

import (
	"OneGfServer/internal/model/device"
	"context"

	"go.uber.org/zap"
)

//   ===============================
//   设备注册业务逻辑
//   ===============================

// 运行逻辑：

// 1. 解析设备注册数据
// 2. 验证设备注册数据
// 3. 生成设备ID
// 4. 设置设备注册时间
// 5. 初始化设备指标
// 6. 复制处理后的数据

// HandleDeviceRegistration 处理设备注册
func (s *sDevice) HandleDeviceRegistration(ctx context.Context, input *device.HandleDeviceRegistrationInput) (*device.HandleDeviceRegistrationOutput, error) {
	//  初始化设备注册输出
	output := &device.HandleDeviceRegistrationOutput{
		DeviceData: make(map[string]interface{}),
		IsSuccess:  false,
	}

	// 获取设备注册数据
	deviceData := input.DeviceData
	zap.L().Debug("设备注册数据", zap.Any("deviceData", deviceData))
	// //   验证设备注册数据
	// validateInput := &device.ValidateDeviceRegistrationInput{
	// 	Name:       s.getStringValue(input.DeviceData, "name", ""),
	// 	MacAddress: s.getStringValue(input.DeviceData, "mac_address", ""),
	// 	IPAddress:  s.getStringValue(input.DeviceData, "ip_address", ""),
	// 	DeviceType: s.getStringValue(input.DeviceData, "device_type", ""),
	// 	Model:      s.getStringValue(input.DeviceData, "model", ""),
	// 	Protocol:   s.getStringValue(input.DeviceData, "protocol", ""),
	// 	Port:       s.getIntValue(input.DeviceData, "port", 0),
	// }

	// validateOutput, err := s.ValidateDeviceRegistration(ctx, validateInput)
	// if err != nil {
	// 	return nil, err
	// }

	// if !validateOutput.IsValid {
	// 	output.Message = "设备注册验证失败: " + validateOutput.Message
	// 	return output, nil
	// }

	// //   生成设备ID（如果未提供）
	// if _, ok := input.DeviceData["device_id"]; !ok {
	// 	input.DeviceData["device_id"] = s.generateDeviceId()
	// }

	// //   设置注册时间
	// input.DeviceData["registered_at"] = gtime.Now().Format("2006-01-02 15:04:05")
	// input.DeviceData["last_heartbeat"] = gtime.Now().Format("2006-01-02 15:04:05")
	// input.DeviceData["status"] = "online"

	// //   初始化设备指标
	// input.DeviceData["cpu_usage"] = 0.0
	// input.DeviceData["memory_usage"] = 0.0
	// input.DeviceData["disk_usage"] = 0.0
	// input.DeviceData["network_latency"] = 0.0
	// input.DeviceData["running_task_count"] = 0
	// input.DeviceData["error_count"] = 0

	// //   复制处理后的数据
	// for key, value := range input.DeviceData {
	// 	output.DeviceData[key] = value
	// }

	// output.DeviceID = s.getStringValue(input.DeviceData, "device_id", "")
	// output.Message = "设备注册成功"
	// output.IsSuccess = true

	return output, nil

}
