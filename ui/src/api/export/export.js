import request from '@/api/request'

function downloadBlob(blob, filename) {
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = filename
  document.body.appendChild(a)
  a.click()
  a.remove()
  URL.revokeObjectURL(url)
}

export async function exportCsv(path, params, filename) {
  const blob = await request.get(path, { params, responseType: 'blob' })
  downloadBlob(blob, filename)
}

export function exportRecords(params = {}) {
  return exportCsv('/export/records.csv', params, `marmot-records-${Date.now()}.csv`)
}

export function exportBuckets() {
  return exportCsv('/export/buckets.csv', {}, `marmot-buckets-${Date.now()}.csv`)
}

export function exportOutstanding() {
  return exportCsv('/export/outstanding.csv', {}, `marmot-outstanding-${Date.now()}.csv`)
}
