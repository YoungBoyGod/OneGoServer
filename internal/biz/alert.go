package biz

import (
	"errors"
	"time"
)

// AlertRepo 告警仓储接口
type AlertRepo interface {
	CreateAlert(alert *Alert) error
	GetAlertByID(id uint) (*Alert, error)
	UpdateAlert(alert *Alert) error
	DeleteAlert(id uint) error
	ListAlerts(offset, limit int, filter *AlertFilter) ([]*Alert, int64, error)
}

// AlertFilter 告警查询过滤器
type AlertFilter struct {
	DeviceID *uint
	UserID   *uint
	Type     *AlertType
	Level    *AlertLevel
	Status   *AlertStatus
}

// AlertUsecase 告警用例
type AlertUsecase struct {
	repo AlertRepo
}

// NewAlertUsecase 创建告警用例
func NewAlertUsecase(repo AlertRepo) *AlertUsecase {
	return &AlertUsecase{repo: repo}
}

// CreateAlert 创建告警
func (uc *AlertUsecase) CreateAlert(alertType AlertType, level AlertLevel, title, message string, deviceID, userID uint, metadata map[string]interface{}) (*Alert, error) {
	// 验证输入
	if err := uc.validateAlertInput(title, message); err != nil {
		return nil, err
	}

	// 创建告警实体
	alert := &Alert{
		Title:     title,
		Message:   message,
		Type:      alertType,
		Level:     level,
		Status:    AlertStatusOpen,
		DeviceID:  deviceID,
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 处理元数据
	if metadata != nil {
		// TODO: 将metadata序列化为JSON字符串
		// alert.Metadata = jsonString
	}

	// 保存到数据库
	if err := uc.repo.CreateAlert(alert); err != nil {
		return nil, err
	}

	return alert, nil
}

// GetAlert 获取告警信息
func (uc *AlertUsecase) GetAlert(id uint) (*Alert, error) {
	return uc.repo.GetAlertByID(id)
}

// AcknowledgeAlert 确认告警
func (uc *AlertUsecase) AcknowledgeAlert(id, userID uint) error {
	alert, err := uc.repo.GetAlertByID(id)
	if err != nil {
		return err
	}

	if alert.Status == AlertStatusClosed {
		return errors.New("告警已关闭，无法确认")
	}

	now := time.Now()
	alert.Status = AlertStatusAcknowledged
	alert.AckedBy = &userID
	alert.AckedAt = &now
	alert.UpdatedAt = now

	return uc.repo.UpdateAlert(alert)
}

// ResolveAlert 解决告警
func (uc *AlertUsecase) ResolveAlert(id uint) error {
	alert, err := uc.repo.GetAlertByID(id)
	if err != nil {
		return err
	}

	if alert.Status == AlertStatusClosed {
		return errors.New("告警已关闭，无法解决")
	}

	now := time.Now()
	alert.Status = AlertStatusResolved
	alert.ResolvedAt = &now
	alert.UpdatedAt = now

	return uc.repo.UpdateAlert(alert)
}

// CloseAlert 关闭告警
func (uc *AlertUsecase) CloseAlert(id uint) error {
	alert, err := uc.repo.GetAlertByID(id)
	if err != nil {
		return err
	}

	alert.Status = AlertStatusClosed
	alert.UpdatedAt = time.Now()

	// 如果还未解决，设置解决时间
	if alert.ResolvedAt == nil {
		now := time.Now()
		alert.ResolvedAt = &now
	}

	return uc.repo.UpdateAlert(alert)
}

// ListAlerts 告警列表
func (uc *AlertUsecase) ListAlerts(page, pageSize int, filter *AlertFilter) ([]*Alert, int64, error) {
	offset := (page - 1) * pageSize
	return uc.repo.ListAlerts(offset, pageSize, filter)
}

// ListDeviceAlerts 获取设备的告警列表
func (uc *AlertUsecase) ListDeviceAlerts(deviceID uint, page, pageSize int) ([]*Alert, int64, error) {
	filter := &AlertFilter{DeviceID: &deviceID}
	offset := (page - 1) * pageSize
	return uc.repo.ListAlerts(offset, pageSize, filter)
}

// ListUserAlerts 获取用户的告警列表
func (uc *AlertUsecase) ListUserAlerts(userID uint, page, pageSize int) ([]*Alert, int64, error) {
	filter := &AlertFilter{UserID: &userID}
	offset := (page - 1) * pageSize
	return uc.repo.ListAlerts(offset, pageSize, filter)
}

// GetAlertStats 获取告警统计信息
func (uc *AlertUsecase) GetAlertStats(userID *uint) (map[string]int64, error) {
	filter := &AlertFilter{}
	if userID != nil {
		filter.UserID = userID
	}

	// 获取全部告警进行统计
	alerts, _, err := uc.repo.ListAlerts(0, 10000, filter)
	if err != nil {
		return nil, err
	}

	stats := make(map[string]int64)
	stats["total"] = int64(len(alerts))
	stats["open"] = 0
	stats["acknowledged"] = 0
	stats["resolved"] = 0
	stats["closed"] = 0

	// 按级别统计
	stats["info"] = 0
	stats["warning"] = 0
	stats["error"] = 0
	stats["critical"] = 0

	for _, alert := range alerts {
		// 按状态统计
		switch alert.Status {
		case AlertStatusOpen:
			stats["open"]++
		case AlertStatusAcknowledged:
			stats["acknowledged"]++
		case AlertStatusResolved:
			stats["resolved"]++
		case AlertStatusClosed:
			stats["closed"]++
		}

		// 按级别统计
		switch alert.Level {
		case AlertLevelInfo:
			stats["info"]++
		case AlertLevelWarning:
			stats["warning"]++
		case AlertLevelError:
			stats["error"]++
		case AlertLevelCritical:
			stats["critical"]++
		}
	}

	return stats, nil
}

// 私有方法

// validateAlertInput 验证告警输入
func (uc *AlertUsecase) validateAlertInput(title, message string) error {
	if len(title) == 0 {
		return errors.New("告警标题不能为空")
	}
	if len(message) == 0 {
		return errors.New("告警消息不能为空")
	}
	return nil
}
