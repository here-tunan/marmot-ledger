import request from '@/api/request'

export function listCategoryGroups(params = {}) {
  return request.get('/category-group', { params })
}

export function getCategoryGroup(id) {
  return request.get(`/category-group/${id}`)
}
