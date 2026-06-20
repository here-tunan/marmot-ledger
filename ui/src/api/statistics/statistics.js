import request from '@/api/request'

export function getStatisticsSummary(params = {}) {
  return request.get('/statistics/summary', { params })
}

export function getStatisticsCategoryGroups(params = {}) {
  return request.get('/statistics/category-group', { params })
}

export function getStatisticsSummaries(params = {}) {
  return request.get('/statistics/summaries', { params })
}

export function getStatisticsCategoryGroupsByCurrency(params = {}) {
  return request.get('/statistics/category-groups', { params })
}

export function getStatisticsTrend(params = {}) {
  return request.get('/statistics/trend', { params })
}

export function getInvestmentSummaries(params = {}) {
  return request.get('/statistics/investment', { params })
}

export function getNetWorthTrend(params = {}) {
  return request.get('/statistics/net-worth-trend', { params })
}
