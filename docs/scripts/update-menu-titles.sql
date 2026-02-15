-- 更新菜单显示名称（与 init_mysql_data 中新文案一致）
-- 执行后重新登录前端即可看到新菜单名

UPDATE menus SET title = '组织与用户', sort = 10 WHERE id = 1;
UPDATE menus SET title = '用户', sort = 11 WHERE id = 2;
UPDATE menus SET title = '部门', sort = 12 WHERE id = 3;
UPDATE menus SET title = '同步字段映射', sort = 13 WHERE id = 4;
UPDATE menus SET title = '系统', sort = 20 WHERE id = 5;
UPDATE menus SET title = '角色与权限', sort = 21 WHERE id = 6;
UPDATE menus SET title = '菜单', sort = 22 WHERE id = 7;
UPDATE menus SET title = '接口', sort = 23 WHERE id = 8;
UPDATE menus SET title = '审计', sort = 30 WHERE id = 9;
UPDATE menus SET title = '操作日志', sort = 31 WHERE id = 10;
UPDATE menus SET title = '系统信息', sort = 24 WHERE id = 11;
UPDATE menus SET title = 'API 密钥', sort = 25 WHERE id = 12;
