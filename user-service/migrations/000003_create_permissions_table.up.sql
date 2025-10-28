-- 1. Schema
CREATE TABLE permissions (id BIGSERIAL PRIMARY KEY, code TEXT NOT NULL UNIQUE);
CREATE TABLE users_permissions (user_id BIGINT, permission_id BIGINT, PRIMARY KEY (user_id, permission_id));

-- 2. Seed (separate file, high migration number)
INSERT INTO permissions (code) VALUES ('products:read'), ('products:write')
ON CONFLICT (code) DO NOTHING;