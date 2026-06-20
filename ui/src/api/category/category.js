import request from '@/api/request'

export function createCategory(data) {
  return request.post('/category', data)
}

export function listCategories(params = {}) {
  return request.get('/category', { params })
}

export function getCategory(id) {
  return request.get(`/category/${id}`)
}

export function updateCategory(id, data) {
  return request.put(`/category/${id}`, data)
}

export function deleteCategory(id) {
  return request.delete(`/category/${id}`)
}

// 检查分类使用情况（删除前调用）
export function checkCategoryUsage(id) {
  return request.get(`/category/${id}/usage`)
}
