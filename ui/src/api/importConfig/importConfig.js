import request from '@/api/request'

export function createImportConfig(data) {
  return request.post('/import-config', data)
}

export function listImportConfigs(params = {}) {
  return request.get('/import-config', { params })
}

export function getImportConfig(id) {
  return request.get(`/import-config/${id}`)
}

export function updateImportConfig(id, data) {
  return request.put(`/import-config/${id}`, data)
}

export function deleteImportConfig(id) {
  return request.delete(`/import-config/${id}`)
}

export function previewImport(configId, file, defaultBucketId) {
  const formData = new FormData()
  formData.append('file', file)
  if (defaultBucketId) formData.append('defaultBucketId', String(defaultBucketId))
  return request.post(`/import-config/${configId}/preview`, formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}

export function commitImport(configId, rows) {
  return request.post(`/import-config/${configId}/commit`, { rows })
}
