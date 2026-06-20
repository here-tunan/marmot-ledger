<template>
  <main class="outstanding-page">
    <section class="page-hero reveal-block">
      <div>
        <p class="eyebrow">{{ t('outstanding.hero.eyebrow') }}</p>
        <h1>{{ t('outstanding.hero.title') }}</h1>
        <p>{{ t('outstanding.hero.subtitle') }}</p>
      </div>
      <button class="export-action" type="button" @click="handleExportOutstanding">{{ t('records.actions.exportOutstanding') }}</button>
    </section>

    <section v-if="currencyTotals.length" class="summary-card reveal-block delay-1">
      <p class="eyebrow">{{ t('outstanding.summary.title') }}</p>
      <div class="summary-grid">
        <div v-for="row in currencyTotals" :key="row.currency" class="summary-row">
          <strong class="summary-currency">{{ row.currency }}</strong>
          <div class="summary-tile">
            <span>📥 {{ t('outstanding.summary.receivable') }}</span>
            <strong>{{ formatMoney(row.receivable, row.currency) }}</strong>
          </div>
          <div class="summary-tile">
            <span>🔒 {{ t('outstanding.summary.deposit') }}</span>
            <strong>{{ formatMoney(row.deposit, row.currency) }}</strong>
          </div>
          <div class="summary-tile">
            <span>🤝 {{ t('outstanding.summary.loanOut') }}</span>
            <strong>{{ formatMoney(row.loanOut, row.currency) }}</strong>
          </div>
          <div class="summary-tile total">
            <span>{{ t('outstanding.summary.total') }}</span>
            <strong>{{ formatMoney(row.total, row.currency) }}</strong>
          </div>
        </div>
      </div>
    </section>

    <section v-else class="empty-celebrate reveal-block delay-1">
      <p>🎉 {{ t('outstanding.empty.all') }}</p>
    </section>

    <section v-loading="loading" class="reveal-block delay-2">
      <div v-for="group in groupedSections" :key="group.kind" class="section-card">
        <div class="section-head">
          <div>
            <p class="eyebrow">{{ group.eyebrow }}</p>
            <h2>{{ group.title }} <span class="section-count">{{ t('outstanding.sections.count', { count: group.items.length }) }}</span></h2>
          </div>
        </div>

        <div v-if="group.items.length" class="bucket-groups">
          <div v-for="bucketGroup in groupByBucket(group.items)" :key="bucketGroup.bucketId" class="bucket-group">
            <p class="bucket-group-label">{{ bucketGroup.emoji }} {{ bucketGroup.bucketName || '#' + bucketGroup.bucketId }}</p>
            <button v-for="item in bucketGroup.items" :key="item.id" class="outstanding-row" type="button" @click="goCollect(item, group.kind)">
              <div class="row-main">
                <strong>{{ item.description || t('record.scenarios.' + item.eventType) }}</strong>
                <small>{{ formatDate(item.eventTime) }}</small>
              </div>
              <div class="row-amount">
                <strong>{{ item.currency }} {{ formatAmount(item.outstandingAmount) }}</strong>
                <small v-if="!eq(item.amount, item.outstandingAmount)">/ {{ formatAmount(item.amount) }}</small>
              </div>
              <span class="row-arrow">{{ group.kind === 'deposits' ? t('outstanding.actions.refund') : t('outstanding.actions.collect') }} →</span>
            </button>
          </div>
        </div>
        <div v-else class="empty-state">
          <p>{{ group.empty }}</p>
        </div>
      </div>
    </section>
  </main>
</template>

<script setup>
import { computed, onActivated, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { useConfigStore } from '@/stores/config'
import { getOutstandingSummary } from '@/api/financialEvent/financialEvent'
import { exportOutstanding } from '@/api/export/export'

const { t } = useI18n()
const router = useRouter()
const config = useConfigStore()

const loading = ref(false)
const summary = ref({ receivables: [], deposits: [], loansOut: [] })

const collectScenarioByKind = {
  receivables: 'receivable_collect',
  deposits: 'deposit_refund',
  loansOut: 'loan_collect',
}
const emojiByKind = {
  receivables: '📥',
  deposits: '🔒',
  loansOut: '🤝',
}

const groupedSections = computed(() => [
  {
    kind: 'receivables',
    eyebrow: 'Receivables',
    title: '📥 ' + t('outstanding.sections.receivables'),
    items: summary.value.receivables || [],
    empty: t('outstanding.empty.receivables'),
  },
  {
    kind: 'deposits',
    eyebrow: 'Deposits',
    title: '🔒 ' + t('outstanding.sections.deposits'),
    items: summary.value.deposits || [],
    empty: t('outstanding.empty.deposits'),
  },
  {
    kind: 'loansOut',
    eyebrow: 'Loans out',
    title: '🤝 ' + t('outstanding.sections.loansOut'),
    items: summary.value.loansOut || [],
    empty: t('outstanding.empty.loansOut'),
  },
])

const currencyTotals = computed(() => {
  const map = new Map()
  const add = (currency, key, amount) => {
    if (!currency) return
    if (!map.has(currency)) map.set(currency, { currency, receivable: 0, deposit: 0, loanOut: 0, total: 0 })
    const row = map.get(currency)
    row[key] += Number(amount || 0)
    row.total += Number(amount || 0)
  }
  ;(summary.value.receivables || []).forEach((it) => add(it.currency, 'receivable', it.outstandingAmount))
  ;(summary.value.deposits || []).forEach((it) => add(it.currency, 'deposit', it.outstandingAmount))
  ;(summary.value.loansOut || []).forEach((it) => add(it.currency, 'loanOut', it.outstandingAmount))
  return Array.from(map.values()).sort((a, b) => b.total - a.total)
})

function groupByBucket(items) {
  const map = new Map()
  items.forEach((item) => {
    const key = item.bucketId || 0
    if (!map.has(key)) {
      map.set(key, { bucketId: item.bucketId, bucketName: item.bucketName, emoji: bucketEmojiFor(item.eventType), items: [] })
    }
    map.get(key).items.push(item)
  })
  return Array.from(map.values())
}

function bucketEmojiFor(eventType) {
  if (eventType === 'receivable_create') return '📥'
  if (eventType === 'deposit_create') return '🔒'
  if (eventType === 'loan_out') return '🤝'
  return '🪣'
}

function formatAmount(value) {
  return new Intl.NumberFormat(config.locale, { minimumFractionDigits: 2, maximumFractionDigits: 2 }).format(Number(value || 0))
}
function formatMoney(value, currency) { return `${currency} ${formatAmount(value)}` }
function formatDate(eventTime) { return (eventTime || '').split(' ')[0] }
function eq(a, b) { return Number(a) === Number(b) }

async function load() {
  loading.value = true
  try {
    const res = await getOutstandingSummary()
    if (res.success) summary.value = res.data || { receivables: [], deposits: [], loansOut: [] }
    else ElMessage.error(res.error || t('outstanding.messages.loadFailed'))
  } finally {
    loading.value = false
  }
}

async function handleExportOutstanding() {
  try { await exportOutstanding() }
  catch (err) { console.warn(err); ElMessage.error(t('records.messages.exportFailed')) }
}

function goCollect(item, kind) {
  const scenario = collectScenarioByKind[kind]
  router.push({
    path: '/record',
    query: {
      scenario,
      relatedId: String(item.id),
      bucketId: String(item.bucketId),
      currency: item.currency,
      amount: String(item.outstandingAmount),
      description: item.description || '',
    },
  })
}

onMounted(load)
onActivated(load)
</script>

<style scoped>
.outstanding-page { max-width: 1200px; margin: 0 auto; color: #1e293b }
.reveal-block { animation: revealUp 480ms cubic-bezier(.16,1,.3,1) both }
.delay-1 { animation-delay: 90ms }
.delay-2 { animation-delay: 160ms }

.page-hero, .summary-card, .section-card, .empty-celebrate {
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 1px 3px rgba(15,23,42,.1), 0 12px 30px rgba(15,23,42,.04);
}

.page-hero { display:flex;justify-content:space-between;gap:16px;align-items:flex-start;margin-bottom: 18px; padding: 26px; background: linear-gradient(135deg, #fffaf0 0%, #fff 70%); }
.export-action{min-height:40px;border:0;border-radius:12px;padding:0 14px;background:rgba(47,125,92,.10);color:#2f7d5c;font-size:12px;font-weight:900;cursor:pointer;box-shadow:inset 0 0 0 1px rgba(47,125,92,.22);white-space:nowrap}.export-action:active{transform:scale(.96)}
.eyebrow { margin: 0 0 8px; color: #2f7d5c; font-size: 12px; font-weight: 800; letter-spacing: .08em; text-transform: uppercase }
.page-hero h1 { margin: 0; font-size: 30px; line-height: 1.16; letter-spacing: -.022em }
.page-hero p:last-child { max-width: 720px; margin: 12px 0 0; color: #64748b; line-height: 1.7 }

.summary-card { margin-bottom: 18px; padding: 20px 22px; }
.summary-grid { display: grid; gap: 12px; margin-top: 12px; }
.summary-row {
  display: grid;
  grid-template-columns: 80px repeat(4, minmax(0, 1fr));
  gap: 10px;
  align-items: center;
  padding: 14px;
  border-radius: 14px;
  background: #fffaf0;
  box-shadow: inset 0 0 0 1px rgba(100,116,139,.10);
}
.summary-currency { color: #1e293b; font-size: 15px; font-weight: 900; letter-spacing: .04em; }
.summary-tile { display: flex; flex-direction: column; gap: 2px; padding: 6px 10px; border-radius: 10px; background: #fff; }
.summary-tile span { color: #64748b; font-size: 11px; font-weight: 800 }
.summary-tile strong { color: #1e293b; font-family: 'SF Mono', 'Fira Code', monospace; font-variant-numeric: tabular-nums; font-size: 15px }
.summary-tile.total { background: rgba(47,125,92,.10); }
.summary-tile.total span { color: #2f7d5c; font-weight: 900 }
.summary-tile.total strong { color: #2f7d5c; }

.empty-celebrate { margin-bottom: 18px; padding: 32px; text-align: center; color: #2f7d5c; font-size: 18px; font-weight: 800; }

.section-card { margin-bottom: 18px; padding: 22px; }
.section-head { margin-bottom: 14px; }
.section-head h2 { margin: 0; display: flex; align-items: center; gap: 10px; font-size: 20px; letter-spacing: -.012em; }
.section-count { color: #64748b; font-size: 12px; font-weight: 700; letter-spacing: 0; text-transform: none; }

.bucket-groups { display: grid; gap: 14px; }
.bucket-group { padding: 14px; border-radius: 14px; background: #f8faf7; }
.bucket-group-label { margin: 0 0 10px; color: #2f7d5c; font-size: 13px; font-weight: 900; letter-spacing: .02em; }

.outstanding-row {
  display: grid;
  grid-template-columns: minmax(0, 1.6fr) minmax(0, 1fr) auto;
  gap: 14px;
  align-items: center;
  width: 100%;
  margin-bottom: 8px;
  padding: 14px 16px;
  border: 0;
  border-radius: 12px;
  background: #fff;
  color: inherit;
  text-align: left;
  cursor: pointer;
  box-shadow: inset 0 0 0 1px rgba(100,116,139,.10);
  transition-property: transform, box-shadow, background-color;
  transition-duration: 160ms;
}
.outstanding-row:last-child { margin-bottom: 0 }
.outstanding-row:hover { background: #fffaf0; box-shadow: inset 0 0 0 1px rgba(47,125,92,.32), 0 6px 14px rgba(15,23,42,.05) }
.outstanding-row:active { transform: scale(.99) }

.row-main { display: flex; flex-direction: column; gap: 3px; min-width: 0 }
.row-main strong { color: #1e293b; font-size: 14px; font-weight: 800; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.row-main small { color: #64748b; font-size: 12px; font-weight: 700 }

.row-amount { display: flex; flex-direction: column; gap: 2px; align-items: flex-end; }
.row-amount strong { color: #1e293b; font-family: 'SF Mono', 'Fira Code', monospace; font-variant-numeric: tabular-nums; font-size: 15px; font-weight: 800 }
.row-amount small { color: #94a3b8; font-family: 'SF Mono', 'Fira Code', monospace; font-size: 11px; }

.row-arrow {
  display: inline-flex;
  align-items: center;
  padding: 6px 12px;
  border-radius: 999px;
  background: rgba(47,125,92,.10);
  color: #2f7d5c;
  font-size: 12px;
  font-weight: 900;
  white-space: nowrap;
}
.outstanding-row:hover .row-arrow { background: rgba(47,125,92,.18); }

.empty-state { padding: 24px; text-align: center; color: #64748b }

@media (max-width: 720px) {
  .page-hero { flex-direction: column; }
  .export-action { width: 100%; }
  .summary-row { grid-template-columns: 1fr; }
  .outstanding-row { grid-template-columns: 1fr; gap: 8px; }
  .row-amount { align-items: flex-start; }
}

@keyframes revealUp {
  from { opacity: 0; transform: translateY(12px); filter: blur(4px) }
  to { opacity: 1; transform: translateY(0); filter: blur(0) }
}
</style>
