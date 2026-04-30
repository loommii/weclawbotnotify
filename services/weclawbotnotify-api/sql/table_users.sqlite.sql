-- 用户表：存储系统账户信息，注册/登录认证的基础表
-- 注意：此文件为 SQLite 语法，用于应用运行时建表

CREATE TABLE IF NOT EXISTS users (
    id            INTEGER      PRIMARY KEY AUTOINCREMENT,              -- 主键，自增
    username      VARCHAR(50)  NOT NULL UNIQUE,                        -- 用户名，最长50字符，唯一
    password_hash CHAR(60)     NOT NULL,                               -- bcrypt哈希，固定60字符
    created_at    INT          NOT NULL DEFAULT (strftime('%s', 'now')), -- 创建时间(Unix时间戳)
    updated_at    INT          NOT NULL DEFAULT (strftime('%s', 'now'))  -- 更新时间(Unix时间戳，应用层维护)
);
