package system

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	system "OneGfServer/internal/model/system"
)

// ===============================
// 系统告警业务逻辑
// ===============================

// GetSystemAlerts 获取系统告警
func (s *sSystem) GetSystemAlerts(ctx context.Context, input *system.GetSystemAlertsInput) (*system.GetSystemAlertsOutput, error) {
	// 这里应该从数据库获取系统告警
	// 目前返回模拟数据
	generateInput := &system.GenerateSystemAlertsInput{
		Level:  input.Level,
		Status: input.Status,
		Page:   input.Page,
		Size:   input.Size,
	}
	alerts := s.generateSystemAlerts(generateInput)

	return &system.GetSystemAlertsOutput{
		List:  alerts.Alerts,
		Total: 100,
		Page:  input.Page,
		Size:  input.Size,
	}, nil
}

// AcknowledgeAlert 确认告警
func (s *sSystem) AcknowledgeAlert(ctx context.Context, input *system.AcknowledgeAlertInput) (*system.AcknowledgeAlertOutput, error) {
	// 这里应该更新数据库中的告警状态
	g.Log().Info(ctx, "确认告警", g.Map{
		"alert_id": input.AlertID,
		"comment":  input.Comment,
		"time":     gtime.Now(),
	})

	return &system.AcknowledgeAlertOutput{
		AlertID: input.AlertID,
		Status:  "acknowledged",
	}, nil
}

// ResolveAlert 解决告警
func (s *sSystem) ResolveAlert(ctx context.Context, input *system.ResolveAlertInput) (*system.ResolveAlertOutput, error) {
	// 这里应该更新数据库中的告警状态
	g.Log().Info(ctx, "解决告警", g.Map{
		"alert_id": input.AlertID,
		"comment":  input.Comment,
		"time":     gtime.Now(),
	})

	return &system.ResolveAlertOutput{
		AlertID: input.AlertID,
		Status:  "resolved",
	}, nil
}

// generateSystemAlerts 生成系统告警
func (s *sSystem) generateSystemAlerts(input *system.GenerateSystemAlertsInput) *system.GenerateSystemAlertsOutput {
	var alerts []map[string]interface{}
	levels := []string{"info", "warn", "error", "critical"}
	statuses := []string{"active", "resolved", "acknowledged"}

	for i := 0; i < input.Size; i++ {
		alertLevel := levels[i%len(levels)]
		alertStatus := statuses[i%len(statuses)]

		if input.Level != "" && alertLevel != input.Level {
			continue
		}
		if input.Status != "" && alertStatus != input.Status {
			continue
		}

		alert := map[string]interface{}{
			"id":             int64((input.Page-1)*input.Size + i + 1),
			"level":          alertLevel,
			"title":          "系统告警标题 " + string(rune(i+1)),
			"message":        "系统告警消息 " + string(rune(i+1)),
			"module":         "system",
			"status":         alertStatus,
			"source":         "monitor",
			"sourceId":       "source-" + string(rune(i+1)),
			"acknowledgedBy": "admin",
			"acknowledgedAt": gtime.Now().Add(-time.Duration(i) * time.Hour).Format("2006-01-02 15:04:05"),
			"resolvedBy":     "admin",
			"resolvedAt":     gtime.Now().Add(-time.Duration(i) * time.Hour).Format("2006-01-02 15:04:05"),
			"createdAt":      gtime.Now().Add(-time.Duration(i) * time.Hour).Format("2006-01-02 15:04:05"),
			"updatedAt":      gtime.Now().Add(-time.Duration(i) * time.Hour).Format("2006-01-02 15:04:05"),
		}
		alerts = append(alerts, alert)
	}

	return &system.GenerateSystemAlertsOutput{
		Alerts: alerts,
	}
}
