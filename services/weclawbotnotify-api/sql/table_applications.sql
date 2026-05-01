-- Application 表：存储用户创建的应用信息，用于推送通知的鉴权基础表
-- 注意：此文件为 MySQL 语法，仅用于 goctl model mysql ddl 生成 model

CREATE TABLE IF NOT EXISTS applications (
    id            INT          NOT NULL AUTO_INCREMENT,                -- 主键，自增
    user_id       INT          NOT NULL,                               -- 关联用户 ID
    token         VARCHAR(64)  NOT NULL,                               -- 应用 Token（随机字符串），唯一
    name          VARCHAR(100) NOT NULL,                               -- 应用名称，最长100字符
    description   VARCHAR(500) DEFAULT '',                             -- 应用描述，可选
    created_at    INT          NOT NULL DEFAULT (UNIX_TIMESTAMP()),     -- 创建时间(Unix时间戳)
    last_used_at  INT          DEFAULT NULL,                           -- 最后使用时间(Unix时间戳)，可为NULL
    PRIMARY KEY (id),                                                    -- 主键约束
    UNIQUE KEY uk_token (token),                                         -- Token 唯一索引
    KEY idx_user_id (user_id)                                            -- 用户 ID 索引
);
