import request from '@/api/request'

// 用户接口
export function listCategoryTemplates(params) {
  return request.get('/category-template', { params })
}

export function importCategoryTemplates(data) {
  return request.post('/category-template/import', data)
}

// 管理员接口
export function listCategoryTemplatesAdmin(params) {
  return request.get('/category-template/admin', { params })
}

export function getCategoryTemplateAdmin(id) {
  return request.get(`/category-template/admin/${id}`)
}

export function createCategoryTemplate(data) {
  return request.post('/category-template/admin', data)
}

export function updateCategoryTemplate(id, data) {
  return request.put(`/category-template/admin/${id}`, data)
}
