import request from '@/api/request'

export function createAccount(data) {
  return request.post('/account', data)
}

export function listAccounts(params = {}) {
  return request.get('/account', { params })
}

export function getAccount(id) {
  return request.get(`/account/${id}`)
}

export function updateAccount(id, data) {
  return request.put(`/account/${id}`, data)
}

export function deleteAccount(id) {
  return request.delete(`/account/${id}`)
}
