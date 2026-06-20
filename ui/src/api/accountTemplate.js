import request from '@/api/request'

// 用户接口
export function listAccountTemplates(params) {
  return request.get('/account-template', { params })
}

export function instantiateAccountFromTemplate(data) {
  return request.post('/account-template/instantiate', data)
}

// 管理员接口
export function listAccountTemplatesAdmin(params) {
  return request.get('/account-template/admin', { params })
}

export function getAccountTemplateAdmin(id) {
  return request.get(`/account-template/admin/${id}`)
}

export function createAccountTemplate(data) {
  return request.post('/account-template/admin', data)
}

export function updateAccountTemplate(id, data) {
  return request.put(`/account-template/admin/${id}`, data)
}
