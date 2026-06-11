<template>
  <main class="dashboard-shell">
    <section class="cockpit-panel reveal-block">
      <div class="brand-cluster">
        <div class="brand-orbit">
          <img :src="marmotOne" :alt="t('dashboard.images.brandAlt')" width="132" height="132" />
          <img :src="marmotTwo" :alt="t('dashboard.images.brandAlt')" width="86" height="86" />
        </div>
        <div>
          <p class="eyebrow">{{ t('dashboard.hero.eyebrow') }}</p>
          <h1>{{ t('dashboard.hero.title') }}</h1>
          <p class="cockpit-copy">{{ t('dashboard.hero.subtitle') }}</p>
          <div class="view-switch" role="group" aria-label="overview view mode">
            <button class="active">{{ t('dashboard.viewMode.personal') }}</button>
            <button disabled>{{ t('dashboard.viewMode.family') }} · {{ t('dashboard.viewMode.familySoon') }}</button>
          </div>
          <div class="cockpit-actions">
            <button class="primary-action" @click="router.push('/accounts')">{{ t('dashboard.hero.createAccount') }}</button>
            <button class="secondary-action" @click="router.push('/buckets')">{{ t('dashboard.hero.createBucket') }}</button>
          </div>
        </div>
      </div>
      <div class="signal-card">
        <span>{{ t('dashboard.signal.label') }}</span>
        <strong>{{ t('dashboard.signal.title') }}</strong>
        <p>{{ t('dashboard.signal.description') }}</p>
      </div>
    </section>

    <section class="metric-grid reveal-block delay-1">
      <div v-for="metric in metrics" :key="metric.label" class="metric-card">
        <span>{{ metric.label }}</span>
        <strong>{{ metric.value }}</strong>
        <p>{{ metric.hint }}</p>
      </div>
    </section>

    <section class="stats-panel reveal-block delay-2">
      <div class="section-head">
        <div>
          <p class="eyebrow">{{ t('statistics.eyebrow') }}</p>
          <h2>{{ t('statistics.title') }}</h2>
        </div>
      </div>
      <div class="stats-grid">
        <div class="stat-tile income"><span>{{ t('statistics.income') }}</span><strong>{{ formatAmount(statisticsSummary?.income) }}</strong></div>
        <div class="stat-tile expense"><span>{{ t('statistics.expense') }}</span><strong>{{ formatAmount(statisticsSummary?.netExpense || statisticsSummary?.expense) }}</strong><small>{{ t('statistics.grossExpense') }} {{ formatAmount(statisticsSummary?.grossExpense) }}</small></div>
        <div class="stat-tile"><span>{{ t('statistics.refund') }}</span><strong>{{ formatAmount(statisticsSummary?.refund) }}</strong></div>
        <div class="stat-tile"><span>{{ t('statistics.net') }}</span><strong>{{ formatAmount(statisticsSummary?.net) }}</strong></div>
        <div class="stat-tile"><span>{{ t('statistics.eventCount') }}</span><strong>{{ statisticsSummary?.eventCount || 0 }}</strong></div>
      </div>
      <div v-if="categoryGroupStats.length" class="category-rank">
        <h3>{{ t('statistics.categoryGroups') }}</h3>
        <div v-for="item in categoryGroupStats" :key="item.categoryGroupCode" class="rank-row">
          <span :style="{ background: item.color || '#10b981' }"></span>
          <strong>{{ item.categoryGroupName }}</strong>
          <em>{{ formatAmount(item.amount) }}</em>
        </div>
      </div>
    </section>

    <section class="workspace-grid reveal-block delay-2">
      <div class="surface-card asset-surface">
        <div class="section-head">
          <div>
            <p class="eyebrow">{{ t('dashboard.buckets.eyebrow') }}</p>
            <h2>{{ t('dashboard.buckets.title') }}</h2>
          </div>
          <button class="ghost-action" @click="router.push('/buckets')">{{ t('common.actions.viewAll') }}</button>
        </div>

        <div v-if="currencyTotals.length" class="currency-stack">
          <div v-for="item in currencyTotals" :key="item.currency" class="currency-row">
            <div>
              <span class="currency-code">{{ item.currency }}</span>
              <small>{{ t('dashboard.buckets.count', { count: item.count }) }}</small>
            </div>
            <strong>{{ formatAmount(item.total) }}</strong>
          </div>
        </div>
        <div v-else class="empty-state compact">
          <img :src="marmotTwo" :alt="t('dashboard.buckets.emptyAlt')" width="92" height="92" />
          <p>{{ t('dashboard.buckets.emptyText') }}</p>
        </div>
      </div>

      <div class="surface-card events-surface">
        <div class="section-head">
          <div>
            <p class="eyebrow">{{ t('dashboard.events.eyebrow') }}</p>
            <h2>{{ t('dashboard.events.title') }}</h2>
          </div>
          <button class="ghost-action" @click="refreshAll">{{ t('common.actions.refresh') }}</button>
        </div>

        <div v-if="recentEvents.length" class="event-list">
          <div v-for="event in recentEvents" :key="event.id" class="event-row">
            <div class="event-dot"></div>
            <div class="event-main">
              <strong>{{ event.description || eventTypeLabel(event.eventType) }}</strong>
              <span>{{ event.eventTime }} · {{ eventTypeLabel(event.eventType) }} · {{ event.remark || t('common.misc.noRemark') }}</span>
            </div>
            <div class="event-amount">{{ event.currency }} {{ formatAmount(event.amount) }}</div>
          </div>
        </div>
        <div v-else class="empty-state compact">
          <img :src="marmotOne" :alt="t('dashboard.events.emptyAlt')" width="92" height="92" />
          <p>{{ t('dashboard.events.emptyText') }}</p>
        </div>
      </div>
    </section>
  </main>
</template>

<script setup>
import { computed, onActivated, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useConfigStore } from '@/stores/config'
import { listAccounts } from '@/api/account/account'
import { listBuckets } from '@/api/bucket/bucket'
import { listFinancialEvents } from '@/api/financialEvent/financialEvent'
import { getStatisticsCategoryGroups, getStatisticsSummary } from '@/api/statistics/statistics'
import marmotOne from '../../../img/marmot-ledger-1.png'
import marmotTwo from '../../../img/marmot-ledger-2.png'

const router = useRouter()
const { t, te } = useI18n()
const config = useConfigStore()
const accounts = ref([])
const buckets = ref([])
const recentEvents = ref([])
const statisticsSummary = ref(null)
const categoryGroupStats = ref([])
const loading = ref(false)

const metrics = computed(() => [
  {
    label: t('dashboard.metrics.accounts.label'),
    value: accounts.value.length,
    hint: t('dashboard.metrics.accounts.hint'),
  },
  {
    label: t('dashboard.metrics.buckets.label'),
    value: buckets.value.length,
    hint: t('dashboard.metrics.buckets.hint'),
  },
  {
    label: t('dashboard.metrics.currencies.label'),
    value: currencyTotals.value.length,
    hint: t('dashboard.metrics.currencies.hint'),
  },
  {
    label: t('dashboard.metrics.initialEvents.label'),
    value: recentEvents.value.length,
    hint: t('dashboard.metrics.initialEvents.hint'),
  },
])

const currencyTotals = computed(() => {
  const map = new Map()
  buckets.value.forEach((item) => {
    if (item.bucketNature !== 'asset') return
    const current = map.get(item.currency) || { currency: item.currency, total: 0, count: 0 }
    current.total += Number(item.balance || 0)
    current.count += 1
    map.set(item.currency, current)
  })
  return Array.from(map.values()).sort((a, b) => b.total - a.total)
})

const eventTypeLabel = (type) => {
  const key = `record.scenarios.${type}`
  return type && te(key) ? t(key) : type
}

const formatAmount = (value) => {
  const number = Number(value || 0)
  return new Intl.NumberFormat(config.locale, {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(number)
}

const refreshAll = async () => {
  loading.value = true
  try {
    const [accountRes, bucketRes, eventRes, summaryRes, categoryStatsRes] = await Promise.all([
      listAccounts(),
      listBuckets(),
      listFinancialEvents({ page: 1, pageSize: 5 }),
      getStatisticsSummary({ currency: 'CNY' }),
      getStatisticsCategoryGroups({ currency: 'CNY', type: 'expense' }),
    ])
    if (accountRes.success) accounts.value = accountRes.data || []
    if (bucketRes.success) buckets.value = bucketRes.data || []
    if (eventRes.success) recentEvents.value = eventRes.data?.list || []
    if (summaryRes.success) statisticsSummary.value = summaryRes.data
    if (categoryStatsRes.success) categoryGroupStats.value = categoryStatsRes.data?.items || []
  } finally {
    loading.value = false
  }
}

onMounted(refreshAll)
onActivated(refreshAll)
</script>

<style scoped>
.dashboard-shell {
  max-width: 1200px;
  margin: 0 auto;
  color: #1e293b;
}

.reveal-block {
  animation: revealUp 520ms cubic-bezier(0.16, 1, 0.3, 1) both;
}

.delay-1 {
  animation-delay: 90ms;
}

.delay-2 {
  animation-delay: 180ms;
}

.cockpit-panel {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 280px;
  gap: 24px;
  margin-bottom: 24px;
  padding: 28px;
  background: linear-gradient(135deg, #fffaf0 0%, #ffffff 58%, #eef7f0 100%);
  border-radius: 22px;
  box-shadow: 0 2px 8px rgba(15, 23, 42, 0.06), 0 18px 48px rgba(47, 125, 92, 0.12);
}

.brand-cluster {
  display: grid;
  grid-template-columns: 210px minmax(0, 1fr);
  gap: 28px;
  align-items: center;
}

.brand-orbit {
  position: relative;
  min-height: 170px;
}

.brand-orbit img:first-child {
  position: absolute;
  left: 18px;
  top: 8px;
  border-radius: 28px;
  box-shadow: 0 16px 36px rgba(47, 125, 92, 0.18);
}

.brand-orbit img:last-child {
  position: absolute;
  right: 8px;
  bottom: 14px;
  border-radius: 22px;
  box-shadow: 0 12px 28px rgba(31, 41, 51, 0.16);
}

.eyebrow {
  margin: 0 0 8px;
  color: #2f7d5c;
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.cockpit-panel h1 {
  max-width: 680px;
  margin: 0;
  color: #1f2933;
  font-size: 34px;
  font-weight: 760;
  line-height: 1.12;
  letter-spacing: -0.022em;
  text-wrap: balance;
}

.cockpit-copy {
  max-width: 620px;
  margin: 16px 0 0;
  color: #64748b;
  font-size: 15px;
  line-height: 1.7;
  text-wrap: pretty;
}

.view-switch {
  display: inline-flex;
  gap: 4px;
  margin-top: 20px;
  padding: 4px;
  border-radius: 999px;
  background: rgba(31, 41, 51, 0.08);
}

.view-switch button {
  min-height: 34px;
  border: 0;
  border-radius: 999px;
  padding: 0 14px;
  background: transparent;
  color: #64748b;
  font-weight: 800;
  cursor: pointer;
  transition-property: transform, background-color, color, box-shadow;
  transition-duration: 160ms;
}

.view-switch button:active {
  transform: scale(0.96);
}

.view-switch button.active {
  background: #1f2933;
  color: rgba(255, 255, 255, 0.9);
  box-shadow: 0 8px 16px rgba(31, 41, 51, 0.14);
}

.view-switch button:disabled {
  cursor: not-allowed;
  opacity: 0.58;
}

.cockpit-actions {
  display: flex;
  gap: 12px;
  margin-top: 24px;
  flex-wrap: wrap;
}

.primary-action,
.secondary-action,
.ghost-action {
  min-height: 40px;
  border: 0;
  border-radius: 12px;
  padding: 0 18px;
  font-weight: 700;
  cursor: pointer;
  transition-property: transform, box-shadow, background-color, color;
  transition-duration: 160ms;
  transition-timing-function: cubic-bezier(0.16, 1, 0.3, 1);
  touch-action: manipulation;
}

.primary-action:active,
.secondary-action:active,
.ghost-action:active {
  transform: scale(0.96);
}

.primary-action {
  background: #3b82f6;
  color: #ffffff;
  box-shadow: 0 10px 24px rgba(59, 130, 246, 0.22);
}

.secondary-action,
.ghost-action {
  background: rgba(255, 255, 255, 0.72);
  color: #1e293b;
  box-shadow: inset 0 0 0 1px rgba(100, 116, 139, 0.18);
}

.signal-card {
  align-self: stretch;
  padding: 22px;
  border-radius: 18px;
  background: rgba(31, 41, 51, 0.92);
  color: rgba(255, 255, 255, 0.86);
  box-shadow: 0 18px 36px rgba(31, 41, 51, 0.18);
}

.signal-card span,
.metric-card span {
  display: block;
  color: #94a3b8;
  font-size: 12px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.08em;
}

.signal-card strong {
  display: block;
  margin-top: 18px;
  font-size: 22px;
  line-height: 1.25;
}

.signal-card p {
  margin: 14px 0 0;
  color: rgba(255, 255, 255, 0.64);
  line-height: 1.7;
}

.metric-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.stats-panel {
  margin-bottom: 24px;
  padding: 24px;
  background: #ffffff;
  border-radius: 16px;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.1), 0 12px 30px rgba(15, 23, 42, 0.04);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 12px;
}

.stat-tile {
  padding: 16px;
  border-radius: 14px;
  background: #f8faf7;
}

.stat-tile span {
  color: #64748b;
  font-size: 12px;
  font-weight: 800;
}

.stat-tile strong {
  display: block;
  margin-top: 8px;
  font-family: 'SF Mono', 'Fira Code', monospace;
  font-size: 24px;
}

.stat-tile small {
  display: block;
  margin-top: 6px;
  color: #94a3b8;
  font-size: 12px;
}

.stat-tile.income strong {
  color: #ef4444;
}

.stat-tile.expense strong {
  color: #10b981;
}

.category-rank {
  margin-top: 18px;
}

.category-rank h3 {
  margin: 0 0 12px;
}

.rank-row {
  display: grid;
  grid-template-columns: 12px minmax(0, 1fr) auto;
  gap: 10px;
  align-items: center;
  padding: 10px 0;
}

.rank-row span {
  width: 10px;
  height: 10px;
  border-radius: 999px;
}

.rank-row em {
  font-style: normal;
  font-family: 'SF Mono', 'Fira Code', monospace;
}

.metric-card,
.surface-card {
  background: #ffffff;
  border-radius: 16px;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.1), 0 10px 24px rgba(15, 23, 42, 0.04);
}

.metric-card {
  padding: 20px;
  transition-property: transform, box-shadow;
  transition-duration: 180ms;
}

.metric-card strong {
  display: block;
  margin-top: 12px;
  font-family: 'SF Mono', 'Fira Code', monospace;
  font-size: 28px;
  letter-spacing: -0.022em;
}

.metric-card p {
  margin: 8px 0 0;
  color: #64748b;
  font-size: 13px;
  line-height: 1.6;
}

.workspace-grid {
  display: grid;
  grid-template-columns: minmax(0, 0.9fr) minmax(0, 1.1fr);
  gap: 24px;
}

.surface-card {
  padding: 24px;
}

.section-head {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  align-items: flex-start;
  margin-bottom: 18px;
}

.section-head h2 {
  margin: 0;
  font-size: 20px;
  letter-spacing: -0.012em;
}

.currency-stack,
.event-list {
  display: grid;
  gap: 10px;
}

.currency-row,
.event-row {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px;
  border-radius: 12px;
  background: #f8faf7;
}

.currency-row {
  justify-content: space-between;
}

.currency-code {
  display: block;
  color: #1e293b;
  font-weight: 800;
}

.currency-row small,
.event-main span {
  color: #64748b;
  font-size: 12px;
}

.currency-row strong,
.event-amount {
  font-family: 'SF Mono', 'Fira Code', monospace;
  font-variant-numeric: tabular-nums;
}

.event-dot {
  width: 10px;
  height: 10px;
  border-radius: 999px;
  background: #2f7d5c;
  box-shadow: 0 0 0 5px rgba(47, 125, 92, 0.12);
}

.event-main {
  flex: 1;
  min-width: 0;
}

.event-main strong,
.event-main span {
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.event-amount {
  color: #1e293b;
  font-weight: 700;
}

.empty-state {
  display: grid;
  place-items: center;
  gap: 12px;
  padding: 24px;
  text-align: center;
  color: #64748b;
}

.empty-state img {
  border-radius: 20px;
}

@media (hover: hover) {
  .metric-card:hover,
  .surface-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 3px 8px rgba(15, 23, 42, 0.12), 0 16px 34px rgba(15, 23, 42, 0.06);
  }
}

@media (max-width: 1024px) {
  .cockpit-panel,
  .brand-cluster,
  .workspace-grid {
    grid-template-columns: 1fr;
  }

  .signal-card {
    max-width: none;
  }

  .metric-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 640px) {
  .dashboard-shell {
    padding-bottom: 24px;
  }

  .cockpit-panel,
  .surface-card {
    padding: 18px;
    border-radius: 18px;
  }

  .brand-cluster {
    gap: 16px;
  }

  .brand-orbit {
    min-height: 130px;
  }

  .brand-orbit img:first-child {
    width: 96px;
    height: 96px;
  }

  .brand-orbit img:last-child {
    width: 64px;
    height: 64px;
    left: 92px;
    right: auto;
  }

  .cockpit-panel h1 {
    font-size: 26px;
  }

  .metric-grid {
    grid-template-columns: 1fr;
  }

  .event-row,
  .currency-row {
    align-items: flex-start;
  }

  .event-row {
    flex-wrap: wrap;
  }

  .event-amount {
    width: 100%;
    padding-left: 24px;
  }
}

@media (prefers-reduced-motion: reduce) {
  .reveal-block,
  .metric-card,
  .surface-card,
  .primary-action,
  .secondary-action,
  .ghost-action {
    animation: none;
    transition: none;
  }
}

@keyframes revealUp {
  from {
    opacity: 0;
    transform: translateY(12px);
    filter: blur(4px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
    filter: blur(0);
  }
}
</style>
