import request from '@/api/request'

export function listFinancialEvents(params = {}) {
  return request.get('/financial-event', { params })
}

export function getFinancialEvent(id) {
  return request.get(`/financial-event/${id}`)
}

export function listOutstandingForBucket(params = {}) {
  return request.get('/financial-event/outstanding', { params })
}

export function getOutstandingSummary() {
  return request.get('/outstanding/summary')
}

export function getEventGroup(groupId, currentId) {
  return request.get(`/financial-event/group/${groupId}`, { params: { currentId } })
}
