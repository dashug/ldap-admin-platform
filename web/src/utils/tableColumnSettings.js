/**
 * 表格列设置：显示/隐藏、列宽记忆，持久化到 localStorage
 */

const STORAGE_KEY_USER = 'go_ldap_admin_user_table_columns'
const STORAGE_KEY_GROUP = 'go_ldap_admin_group_table_columns'

/** 用户列表默认列配置：prop, label, width(可选), minWidth(可选), visible */
export const defaultUserColumns = [
  { prop: 'username', label: '用户名', minWidth: 90, visible: true },
  { prop: 'nickname', label: '中文名', minWidth: 90, visible: true },
  { prop: 'givenName', label: '花名', minWidth: 80, visible: true },
  { prop: 'status', label: '状态', width: 70, visible: true, customSlot: true },
  { prop: 'mail', label: '邮箱', minWidth: 140, visible: true },
  { prop: 'mobile', label: '手机号', minWidth: 110, visible: true },
  { prop: 'jobNumber', label: '工号', minWidth: 80, visible: true },
  { prop: 'departments', label: '部门', minWidth: 120, visible: true },
  { prop: 'position', label: '职位', minWidth: 90, visible: true },
  { prop: 'creator', label: '创建人', minWidth: 80, visible: true },
  { prop: 'introduction', label: '说明', minWidth: 100, visible: true },
  { prop: 'userDn', label: 'DN', minWidth: 180, visible: true },
  { prop: 'CreatedAt', label: '创建时间', minWidth: 155, visible: true },
  { prop: 'UpdatedAt', label: '更新时间', minWidth: 155, visible: true },
  { prop: 'lastLoginAt', label: '最后登录时间', minWidth: 155, visible: true },
  { prop: 'expireAt', label: '过期日', minWidth: 110, visible: true }
]

/** 部门列表默认列配置 */
export const defaultGroupColumns = [
  { prop: 'groupName', label: '名称', minWidth: 140, visible: true },
  { prop: 'groupType', label: '类型', minWidth: 80, visible: true },
  { prop: 'groupDn', label: 'DN', minWidth: 200, visible: true },
  { prop: 'remark', label: '描述', minWidth: 120, visible: true },
  { prop: 'CreatedAt', label: '创建时间', minWidth: 155, visible: true },
  { prop: 'UpdatedAt', label: '更新时间', minWidth: 155, visible: true }
]

/**
 * 从 localStorage 读取列配置并合并默认列
 * @param {string} storageKey
 * @param {Array} defaultColumns
 * @returns {{ visible: Object, widths: Object, columns: Array }}
 */
export function loadTableColumnConfig(storageKey, defaultColumns) {
  let saved = {}
  try {
    const raw = localStorage.getItem(storageKey)
    if (raw) saved = JSON.parse(raw)
  } catch (_) {}
  const visible = { ...saved.visible }
  const widths = { ...saved.widths }
  const columns = defaultColumns.map(col => ({
    ...col,
    visible: visible[col.prop] !== undefined ? visible[col.prop] : col.visible,
    width: widths[col.prop] !== undefined ? widths[col.prop] : (col.width || null),
    minWidth: col.minWidth
  }))
  defaultColumns.forEach(col => {
    if (visible[col.prop] === undefined) visible[col.prop] = col.visible
    if (widths[col.prop] === undefined && col.width) widths[col.prop] = col.width
  })
  return { visible, widths, columns }
}

/**
 * 保存列配置到 localStorage
 * @param {string} storageKey
 * @param {{ visible: Object, widths: Object }} config
 */
export function saveTableColumnConfig(storageKey, config) {
  try {
    localStorage.setItem(storageKey, JSON.stringify({
      visible: config.visible,
      widths: config.widths
    }))
  } catch (_) {}
}

export const STORAGE_KEY_USER_TABLE = STORAGE_KEY_USER
export const STORAGE_KEY_GROUP_TABLE = STORAGE_KEY_GROUP
