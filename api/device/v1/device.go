package v1

import "github.com/gogf/gf/v2/frame/g"

/*
基础设备管理：
	1. 设备主动注册
	2. 获取设备列表
	3. 设备白名单管理
	4. 删除设备

设备状态管理：
	1. 获取设备状态
	2. 更新设备状态

设备控制操作：
	1. 发送设备命令
	2. 获取设备心跳
	3. 更新设备心跳

设备信息查询：
	1. 获取设备详情
	2. 更新设备信息

设备任务管理：
	1. 获取设备任务列表
	2. 获取设备任务队列
	3. 获取设备任务详情


设备告警管理：
	1. 获取设备告警列表
	2. 获取设备告警详情
	3. 更新设备告警



*/

// RegisterDevice 注册设备
type RegisterDevice struct {
	g.Meta `path:"/device/register" method:"post" tags:"设备管理" summary:"注册设备"`
}
