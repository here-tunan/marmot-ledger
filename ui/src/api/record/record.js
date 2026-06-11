import request from '@/api/request'

export function createRecord(data) {
  return request.post('/record', data)
}
