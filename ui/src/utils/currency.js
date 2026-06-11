export const currencyDisplayMap = {
  CNY: { icon: '🇨🇳', name: '人民币', englishName: 'Chinese Yuan', symbol: '¥' },
  USD: { icon: '🇺🇸', name: '美元', englishName: 'US Dollar', symbol: '$' },
  HKD: { icon: '🇭🇰', name: '港币', englishName: 'Hong Kong Dollar', symbol: 'HK$' },
  EUR: { icon: '🇪🇺', name: '欧元', englishName: 'Euro', symbol: '€' },
  JPY: { icon: '🇯🇵', name: '日元', englishName: 'Japanese Yen', symbol: '¥' },
  GBP: { icon: '🇬🇧', name: '英镑', englishName: 'British Pound', symbol: '£' },
  SGD: { icon: '🇸🇬', name: '新加坡元', englishName: 'Singapore Dollar', symbol: 'S$' },
  AUD: { icon: '🇦🇺', name: '澳大利亚元', englishName: 'Australian Dollar', symbol: 'A$' },
  NZD: { icon: '🇳🇿', name: '新西兰元', englishName: 'New Zealand Dollar', symbol: 'NZ$' },
}

export const currencyOptions = Object.keys(currencyDisplayMap).map((code) => ({
  code,
  ...currencyDisplayMap[code],
}))

export function getCurrencyDisplay(code) {
  return currencyDisplayMap[code] || { icon: '¤', name: code, englishName: code, symbol: code }
}

export function getCurrencyLabel(code, locale = 'zh-CN') {
  const display = getCurrencyDisplay(code)
  const name = locale === 'en-US' ? display.englishName : display.name
  return `${display.icon} ${code} · ${name}`
}
