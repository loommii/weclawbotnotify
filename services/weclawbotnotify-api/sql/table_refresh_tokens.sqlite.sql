-- 刷新令牌表：存储 Refresh Token 的 JTI，用于轮换验证和过期清理
-- 注意：此文件为 SQLite 语法，用于应用运行时建表

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id            INTEGER      PRIMARY KEY AUTOINCREMENT,              -- 主键，自增
    user_id       INT          NOT NULL,                               -- 关联用户 ID
    jti           VARCHAR(36)  NOT NULL UNIQUE,                        -- JWT ID (UUID)，用于轮换验证，唯一
    revoked       INT          NOT NULL DEFAULT 0,                     -- 是否已撤销：0=有效，1=已撤销
    expires_at    INT          NOT NULL,                               -- 过期时间(Unix时间戳)
    created_at    INT          NOT NULL DEFAULT (strftime('%s', 'now')) -- 创建时间(Unix时间戳)
);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens (user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens (expires_at);
