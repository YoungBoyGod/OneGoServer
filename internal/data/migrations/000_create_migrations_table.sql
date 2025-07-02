-- 创建migrations表来记录执行历史
CREATE TABLE IF NOT EXISTS schema_migrations (
    version         VARCHAR(255) PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    applied_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    applied_by      VARCHAR(255) NOT NULL DEFAULT CURRENT_USER,
    checksum        VARCHAR(64) NOT NULL, -- 文件内容的MD5校验和
    execution_time  INTEGER NOT NULL, -- 执行时间(毫秒)
    is_success      BOOLEAN NOT NULL DEFAULT TRUE,
    error_message   TEXT,
    rollback_sql    TEXT -- 回滚SQL
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_migrations_applied_at ON schema_migrations(applied_at);
CREATE INDEX IF NOT EXISTS idx_migrations_success ON schema_migrations(is_success); 