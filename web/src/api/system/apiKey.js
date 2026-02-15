import request from '@/utils/request'

// API 密钥列表（分页）
export function getApiKeyList(params) {
  return request({
    url: '/api/apiKey/list',
    method: 'get',
    params
  })
}

// 创建 API 密钥（返回的 key 仅此一次）
export function createApiKey(data) {
  return request({
    url: '/api/apiKey/create',
    method: 'post',
    data
  })
}

// 删除 API 密钥
export function deleteApiKey(data) {
  return request({
    url: '/api/apiKey/delete',
    method: 'post',
    data
  })
}
