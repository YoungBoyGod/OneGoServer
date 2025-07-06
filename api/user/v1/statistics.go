package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ===============================
// 用户统计分析相关API
// ===============================

// GetUserStatisticsReq 获取用户统计信息请求
type GetUserStatisticsReq struct {
	g.Meta     `path:"/user/statistics" method:"get" tags:"用户统计" summary:"获取用户统计信息"`
	TimeRange  string `json:"time_range" d:"30d" v:"in:24h,7d,30d,90d#时间范围无效"`
	Department string `json:"department,omitempty" v:"length:1,100#部门名称长度为1-100字符"`
	Role       string `json:"role,omitempty" v:"length:1,50#角色名称长度为1-50字符"`
	GroupBy    string `json:"group_by" d:"status" v:"in:status,role,department,registration_date#分组字段无效"`
}

type GetUserStatisticsRes struct {
	Statistics UserStatistics `json:"statistics"`
}

// GetUserBehaviorAnalysisReq 获取用户行为分析请求
type GetUserBehaviorAnalysisReq struct {
	g.Meta    `path:"/user/{userId}/behavior" method:"get" tags:"用户统计" summary:"获取用户行为分析"`
	UserId    string `json:"user_id" v:"required|max-length:50#用户ID不能为空"`
	TimeRange string `json:"time_range" d:"30d" v:"in:7d,30d,90d#时间范围无效"`
}

type GetUserBehaviorAnalysisRes struct {
	BehaviorAnalysis UserBehaviorAnalysis `json:"behavior_analysis"`
}
