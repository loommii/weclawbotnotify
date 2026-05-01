-- Application 表：存储用户创建的应用信息，用于推送通知的鉴权基础表
-- 注意：此文件为 SQLite 语法，用于应用运行时建表

CREATE TABLE IF NOT EXISTS applications (
    id            INTEGER      PRIMARY KEY AUTOINCREMENT,              -- 主键，自增
    user_id       INT          NOT NULL,                               -- 关联用户 ID
    token         VARCHAR(64)  NOT NULL UNIQUE,                        -- 应用 Token（随机字符串），唯一
    name          VARCHAR(100) NOT NULL,                               -- 应用名称，最长100字符
    description   VARCHAR(500) DEFAULT '',                             -- 应用描述，可选
    status        INT          NOT NULL DEFAULT 1,                     -- 状态：1=正常，2=禁用，-1=已删除（软删除）
    created_at    INT          NOT NULL DEFAULT (strftime('%s', 'now')), -- 创建时间(Unix时间戳)
    last_used_at  INT          DEFAULT NULL                            -- 最后使用时间(Unix时间戳)，可为NULL
);
