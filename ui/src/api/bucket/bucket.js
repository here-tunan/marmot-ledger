import request from '@/api/request'

export function createBucket(data) {
  return request.post('/bucket', data)
}

export function listBuckets(params = {}) {
  return request.get('/bucket', { params })
}

export function getBucket(id) {
  return request.get(`/bucket/${id}`)
}

export function updateBucket(id, data) {
  return request.put(`/bucket/${id}`, data)
}

export function listBucketLedgerEntries(id) {
  return request.get(`/bucket/${id}/ledger-entry`)
}
