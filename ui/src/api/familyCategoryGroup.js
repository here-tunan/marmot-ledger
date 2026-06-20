import request from '@/api/request'

export function listFamilyCategoryGroups(familyId, params) {
  return request.get(`/family/${familyId}/category-group`, { params })
}

export function getFamilyCategoryGroup(familyId, groupId) {
  return request.get(`/family/${familyId}/category-group/${groupId}`)
}

export function createFamilyCategoryGroup(familyId, data) {
  return request.post(`/family/${familyId}/category-group`, data)
}

export function updateFamilyCategoryGroup(familyId, groupId, data) {
  return request.put(`/family/${familyId}/category-group/${groupId}`, data)
}

export function deleteFamilyCategoryGroup(familyId, groupId) {
  return request.delete(`/family/${familyId}/category-group/${groupId}`)
}

export function getGroupCategoryIds(familyId, groupId) {
  return request.get(`/family/${familyId}/category-group/${groupId}/members`)
}

export function addCategoriesToGroup(familyId, groupId, categoryIds) {
  return request.post(`/family/${familyId}/category-group/${groupId}/members`, { categoryIds })
}

export function removeCategoryFromGroup(familyId, groupId, categoryId) {
  return request.delete(`/family/${familyId}/category-group/${groupId}/members/${categoryId}`)
}
