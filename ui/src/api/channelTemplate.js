import request from '@/api/request'

// 用户接口
export function listChannelTemplates(params) {
  return request.get('/channel-template', { params })
}

// 管理员接口
export function listChannelTemplatesAdmin(params) {
  return request.get('/channel-template/admin', { params })
}

export function getChannelTemplateAdmin(id) {
  return request.get(`/channel-template/admin/${id}`)
}

export function createChannelTemplate(data) {
  return request.post('/channel-template/admin', data)
}

export function updateChannelTemplate(id, data) {
  return request.put(`/channel-template/admin/${id}`, data)
}
