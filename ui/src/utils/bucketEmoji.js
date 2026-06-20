const map = {
  cash: '💵',
  wallet: '📱',
  bank: '🏦',
  credit: '💳',
  investment_cash: '💰',
  investment_asset: '📈',
  receivable: '📥',
  deposit: '🔒',
  loan_out: '🤝',
  liability: '📉',
  virtual: '🪣',
}

export function getBucketEmoji(bucketType) {
  return map[bucketType] || '🪣'
}

export function decorateBucketName(bucket) {
  if (!bucket) return ''
  return `${getBucketEmoji(bucket.bucketType)} ${bucket.name}`
}
