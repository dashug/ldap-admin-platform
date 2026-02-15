-- 为已有数据库补充「API 密钥」菜单（若已存在可忽略）
-- 执行后需重新登录前端，或清除 token 后再次登录，侧栏才会刷新

-- 1. 插入 API 密钥 菜单（id=12，父级为系统管理 id=4）
INSERT IGNORE INTO menus (id, name, title, icon, path, redirect, component, sort, status, hidden, no_cache, always_show, breadcrumb, parent_id, creator, created_at, updated_at, deleted_at)
VALUES (
  12,
  'ApiKey',
  'API 密钥',
  'tree',
  'apiKey',
  '',
  '/system/apiKey/index',
  16,
  1,
  2,
  2,
  2,
  1,
  4,
  '系统',
  NOW(),
  NOW(),
  NULL
);

-- 2. 把该菜单挂到管理员角色（role_id=1）
INSERT IGNORE INTO role_menus (role_id, menu_id)
SELECT 1, 12
WHERE NOT EXISTS (SELECT 1 FROM role_menus WHERE role_id = 1 AND menu_id = 12);
