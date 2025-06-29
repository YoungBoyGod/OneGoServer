-- 创建告警表迁移文件
-- 执行时间: 2025-07-01
-- 描述: 添加 alerts 表以支持告警功能

-- 创建告警表
CREATE TABLE IF NOT EXISTS alerts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    message VARCHAR(1000) NOT NULL,
    type INTEGER NOT NULL,
    level INTEGER NOT NULL,
    status INTEGER NOT NULL DEFAULT 1,
    source VARCHAR(100),
    metadata TEXT,
    device_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    acked_by INTEGER,
    acked_at TIMESTAMP,
    resolved_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_alerts_device_id ON alerts(device_id);
CREATE INDEX IF NOT EXISTS idx_alerts_user_id ON alerts(user_id);
CREATE INDEX IF NOT EXISTS idx_alerts_status ON alerts(status);
CREATE INDEX IF NOT EXISTS idx_alerts_level ON alerts(level);
CREATE INDEX IF NOT EXISTS idx_alerts_type ON alerts(type);
CREATE INDEX IF NOT EXISTS idx_alerts_created_at ON alerts(created_at);
CREATE INDEX IF NOT EXISTS idx_alerts_deleted_at ON alerts(deleted_at);

-- 添加外键约束
ALTER TABLE alerts ADD CONSTRAINT fk_alerts_device_id 
  FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE;

ALTER TABLE alerts ADD CONSTRAINT fk_alerts_user_id 
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE alerts ADD CONSTRAINT fk_alerts_acked_by 
  FOREIGN KEY (acked_by) REFERENCES users(id) ON DELETE SET NULL;

-- 添加检查约束
ALTER TABLE alerts ADD CONSTRAINT chk_alert_type 
  CHECK (type >= 1 AND type <= 5);

ALTER TABLE alerts ADD CONSTRAINT chk_alert_level 
  CHECK (level >= 1 AND level <= 4);

ALTER TABLE alerts ADD CONSTRAINT chk_alert_status 
  CHECK (status >= 1 AND status <= 4);

-- 添加注释
COMMENT ON TABLE alerts IS '告警表';
COMMENT ON COLUMN alerts.id IS '主键ID';
COMMENT ON COLUMN alerts.title IS '告警标题';
COMMENT ON COLUMN alerts.message IS '告警消息';
COMMENT ON COLUMN alerts.type IS '告警类型: 1-系统,2-设备,3-任务,4-安全,5-性能';
COMMENT ON COLUMN alerts.level IS '告警级别: 1-信息,2-警告,3-错误,4-严重';
COMMENT ON COLUMN alerts.status IS '告警状态: 1-开放,2-已确认,3-已解决,4-已关闭';
COMMENT ON COLUMN alerts.source IS '告警来源';
COMMENT ON COLUMN alerts.metadata IS '元数据(JSON格式)';
COMMENT ON COLUMN alerts.device_id IS '关联设备ID';
COMMENT ON COLUMN alerts.user_id IS '关联用户ID';
COMMENT ON COLUMN alerts.acked_by IS '确认人ID';
COMMENT ON COLUMN alerts.acked_at IS '确认时间';
COMMENT ON COLUMN alerts.resolved_at IS '解决时间'; 