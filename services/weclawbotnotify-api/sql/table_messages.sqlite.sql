-- Message 表：存储推送通知消息记录
-- 注意：此文件为 SQLite 语法，用于应用运行时建表

CREATE TABLE IF NOT EXISTS messages (
    id             INTEGER      PRIMARY KEY AUTOINCREMENT,              -- 主键，自增
    application_id INT          NOT NULL,                               -- 关联应用 ID
    title          VARCHAR(200) NOT NULL,                               -- 消息标题，最长200字符
    message        TEXT         NOT NULL,                               -- 消息内容
    status         INT          NOT NULL DEFAULT 100,                   -- 状态：100=待推送，200=成功，300=失败，301=微信API错误，302=超时，303=Client过期，304=部分成功，-1=已删除
    date           INT          NOT NULL DEFAULT (strftime('%s', 'now')) -- 创建时间(Unix时间戳)
);
