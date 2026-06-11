import request from '@/api/request'

export function createFamily(data) {
  return request.post('/family', data)
}

export function listFamilies() {
  return request.get('/family')
}

export function getFamily(id) {
  return request.get(`/family/${id}`)
}

export function listFamilyMembers(id, params = {}) {
  return request.get(`/family/${id}/members`, { params })
}

export function inviteFamilyMember(id, data) {
  return request.post(`/family/${id}/invitations`, data)
}

export function listFamilyInvitations() {
  return request.get('/family/invitations')
}

export function acceptFamilyInvitation(invitationId) {
  return request.post(`/family/invitations/${invitationId}/accept`)
}

export function rejectFamilyInvitation(invitationId) {
  return request.post(`/family/invitations/${invitationId}/reject`)
}

export function getFamilyStatisticsSummary(id, params = {}) {
  return request.get(`/family/${id}/statistics/summary`, { params })
}

export function getFamilyStatisticsCategoryGroups(id, params = {}) {
  return request.get(`/family/${id}/statistics/category-group`, { params })
}
