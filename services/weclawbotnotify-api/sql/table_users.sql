-- 用户表：存储系统账户信息，注册/登录认证的基础表
-- 注意：此文件为 MySQL 语法，仅用于 goctl model mysql ddl 生成 model

CREATE TABLE IF NOT EXISTS users (
    id            INT          NOT NULL AUTO_INCREMENT,                -- 主键，自增
    username      VARCHAR(50)  NOT NULL,                               -- 用户名，最长50字符
    password_hash CHAR(60)     NOT NULL,                               -- bcrypt哈希，固定60字符
    created_at    INT          NOT NULL DEFAULT (UNIX_TIMESTAMP()),     -- 创建时间(Unix时间戳)
    updated_at    INT          NOT NULL DEFAULT (UNIX_TIMESTAMP()),     -- 更新时间(Unix时间戳，应用层维护)
    PRIMARY KEY (id),                                                    -- 主键约束
    UNIQUE KEY uk_username (username)                                    -- 用户名唯一索引
);
