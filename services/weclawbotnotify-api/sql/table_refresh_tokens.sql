-- 刷新令牌表：存储 Refresh Token 的 JTI，用于轮换验证和过期清理
-- 注意：此文件为 MySQL 语法，仅用于 goctl model mysql ddl 生成 model

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id            INT          NOT NULL AUTO_INCREMENT,                -- 主键，自增
    user_id       INT          NOT NULL,                               -- 关联用户 ID
    jti           VARCHAR(36)  NOT NULL,                               -- JWT ID (UUID)，用于轮换验证
    revoked       TINYINT      NOT NULL DEFAULT 0,                     -- 是否已撤销：0=有效，1=已撤销
    expires_at    INT          NOT NULL,                               -- 过期时间(Unix时间戳)
    created_at    INT          NOT NULL DEFAULT (UNIX_TIMESTAMP()),     -- 创建时间(Unix时间戳)
    PRIMARY KEY (id),                                                    -- 主键约束
    UNIQUE KEY uk_jti (jti),                                            -- JTI 唯一索引，用于轮换查询
    KEY idx_user_id (user_id),                                          -- 用户 ID 索引，用于修改密码时批量撤销
    KEY idx_expires_at (expires_at)                                     -- 过期时间索引，用于清理过期令牌
);
