import request from '@/api/request'

export function createChannel(data) {
  return request.post('/channel', data)
}

export function listChannels(params = {}) {
  return request.get('/channel', { params })
}

export function getChannel(id) {
  return request.get(`/channel/${id}`)
}

export function updateChannel(id, data) {
  return request.put(`/channel/${id}`, data)
}

export function deleteChannel(id) {
  return request.delete(`/channel/${id}`)
}

export function checkChannelUsage(id) {
  return request.get(`/channel/${id}/usage`)
}

export function importChannelTemplates(templateIds) {
  return request.post('/channel/import', { templateIds })
}
