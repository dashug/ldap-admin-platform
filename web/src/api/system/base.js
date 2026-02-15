import request from '@/utils/request'
// 用户登录（已完成）
export function login(data) {
  return request({
    url: '/api/base/login',
    method: 'post',
    data
  })
}

export function refreshToken() {
  return request({
    url: '/api/base/refreshToken',
    method: 'post'
  })
}
// 用户退出接口（已完成）
export function logout() {
  return request({
    url: '/api/base/logout',
    method: 'post'
  })
}
// 获取登录用 RSA 公钥（未配置 VUE_APP_PUBLIC_KEY 时由后端获取，保证与后端一致）
export function getPublicKey() {
  return request({
    url: '/api/base/publicKey',
    method: 'get'
  })
}

// 获取配置信息
export function getConfig() {
  return request({
    url: '/api/base/config',
    method: 'get'
  })
}

// 更新目录服务配置
export function updateDirectoryConfig(data) {
  return request({
    url: '/api/base/directoryConfig',
    method: 'post',
    data
  })
}

// 导入配置（目录与同步规则 JSON）
export function importConfig(data) {
  return request({
    url: '/api/base/configImport',
    method: 'post',
    data
  })
}

// 更新第三方平台配置
export function updateThirdPartyConfig(data) {
  return request({
    url: '/api/base/thirdPartyConfig',
    method: 'post',
    data
  })
}

// 测试第三方平台配置
export function testThirdPartyConfig(data) {
  return request({
    url: '/api/base/thirdPartyConfig/test',
    method: 'post',
    data
  })
}

// 更新邮件通知配置（如是否发送用户创建通知）
export function updateEmailConfig(data) {
  return request({
    url: '/api/base/emailConfig',
    method: 'post',
    data
  })
}

// 获取版本信息
export function getVersion() {
  return request({
    url: '/api/base/version',
    method: 'get'
  })
}

// 获取 LDAP 连接状态
export function getLDAPStatus() {
  return request({
    url: '/api/base/ldapStatus',
    method: 'get'
  })
}

// 获取系统信息（版本、运行时长、数据库状态）
export function getSystemInfo() {
  return request({
    url: '/api/base/systemInfo',
    method: 'get'
  })
}
