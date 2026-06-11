import request from '@/api/request'

export function getStatisticsSummary(params = {}) {
  return request.get('/statistics/summary', { params })
}

export function getStatisticsCategoryGroups(params = {}) {
  return request.get('/statistics/category-group', { params })
}
