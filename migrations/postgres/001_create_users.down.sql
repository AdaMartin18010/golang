-- 删除索引
DROP INDEX IF EXISTS idx_users_created_at;
DROP INDEX IF EXISTS idx_users_email;

-- 删除表
DROP TABLE IF EXISTS users;
