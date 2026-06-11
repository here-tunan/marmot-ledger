import request from '@/api/request'

export function createRecord(data) {
  return request.post('/record', data)
}

export function updateRecord(id, data) {
  return request.put(`/record/${id}`, data)
}

export function deleteRecord(id) {
  return request.delete(`/record/${id}`)
}
