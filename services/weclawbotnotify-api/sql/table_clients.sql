-- Client 表：存储用户创建的客户端信息，用于微信登录和消息推送
-- 注意：此文件为 MySQL 语法，仅用于 goctl model mysql ddl 生成 model

CREATE TABLE IF NOT EXISTS clients (
    id            INT          NOT NULL AUTO_INCREMENT,                -- 主键，自增
    user_id       INT          NOT NULL,                               -- 关联用户 ID
    token         VARCHAR(64)  NOT NULL,                               -- 客户端 Token（随机字符串），唯一
    name          VARCHAR(100) NOT NULL,                               -- 客户端名称，最长100字符
    ilink_user_id VARCHAR(64)  DEFAULT '',                             -- iLink 用户 ID（绑定后填充）
    status        VARCHAR(20)  NOT NULL DEFAULT 'pending',             -- 状态：pending/bound/disconnected
    created_at    INT          NOT NULL DEFAULT (UNIX_TIMESTAMP()),     -- 创建时间(Unix时间戳)
    last_used_at  INT          DEFAULT NULL,                           -- 最后使用时间(Unix时间戳)，可为NULL
    PRIMARY KEY (id),                                                    -- 主键约束
    UNIQUE KEY uk_token (token),                                         -- Token 唯一索引
    KEY idx_user_id (user_id)                                            -- 用户 ID 索引
);
