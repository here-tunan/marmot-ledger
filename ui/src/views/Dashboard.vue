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
          <div class="view-switch" role="group" :aria-label="t('dashboard.viewMode.label')">
            <button :class="{ active: viewMode === 'personal' }" @click="switchMode('personal')">{{ t('dashboard.viewMode.personal') }}</button>
            <button :class="{ active: viewMode === 'family' }" @click="switchMode('family')">{{ t('dashboard.viewMode.family') }}</button>
          </div>
          <div v-if="isFamilyMode" class="family-workspace-card">
            <div>
              <span class="workspace-label">{{ t('dashboard.familySelector.workspaceTitle') }}</span>
              <strong>{{ selectedFamily?.name || t('dashboard.familySelector.empty') }}</strong>
              <small>{{ t('dashboard.familySelector.memberCount', { count: familyMembers.length }) }}</small>
            </div>
            <div class="workspace-actions">
              <span class="workspace-badge">{{ t('dashboard.familySelector.current') }}</span>
              <button class="workspace-manage" type="button" @click="router.push('/family')">{{ t('dashboard.familySelector.manage') }}</button>
            </div>
          </div>
          <div v-if="isFamilyMode && families.length" class="family-chip-rail" :aria-label="t('dashboard.familySelector.switchWorkspace')">
            <button v-for="item in families" :key="item.id" type="button" class="family-chip" :class="{ active: item.id === selectedFamilyId }" :aria-pressed="item.id === selectedFamilyId" @click="selectFamilyWorkspace(item.id)">
              {{ item.name }}
            </button>
          </div>
          <div class="cockpit-actions">
            <button v-if="!isFamilyMode" class="primary-action" @click="router.push('/accounts')">{{ t('dashboard.hero.createAccount') }}</button>
            <button v-if="!isFamilyMode" class="secondary-action" @click="router.push('/buckets')">{{ t('dashboard.hero.createBucket') }}</button>
          </div>
        </div>
      </div>
      <div class="signal-card">
        <span>{{ isFamilyMode ? t('dashboard.familyMembers.eyebrow') : t('dashboard.signal.label') }}</span>
        <strong>{{ isFamilyMode ? (selectedFamily?.name || t('dashboard.familySelector.empty')) : t('dashboard.signal.title') }}</strong>
        <p>{{ isFamilyMode ? t('dashboard.familyStats.hint') : t('dashboard.signal.description') }}</p>
      </div>
    </section>

    <section class="finance-overview-panel reveal-block delay-1">
      <div class="finance-overview-head">
        <div>
          <p class="eyebrow">{{ t(isFamilyMode ? 'dashboard.financeOverview.familyEyebrow' : 'dashboard.financeOverview.eyebrow') }}</p>
          <h2>{{ t(isFamilyMode ? 'dashboard.financeOverview.familyTitle' : 'dashboard.financeOverview.title') }}</h2>
          <p class="stats-hint">{{ t('dashboard.financeOverview.hint') }}</p>
        </div>
        <div class="stats-range-switch" role="group" :aria-label="t('dashboard.netWorthTrend.currencyLabel')">
          <button v-for="c in availableNetWorthCurrencies" :key="c" type="button" :class="{ active: netWorthCurrency === c }" :aria-pressed="netWorthCurrency === c" @click="switchNetWorthCurrency(c)">
            {{ c }}
          </button>
        </div>
      </div>

      <div class="metric-grid">
        <div v-for="metric in displayMetrics" :key="metric.label" class="metric-card" :class="metric.tone">
          <span>{{ metric.label }}</span>
          <strong>{{ metric.value }}</strong>
          <p>{{ metric.hint }}</p>
        </div>
      </div>

      <div class="finance-overview-grid">
        <div class="overview-subcard">
          <div class="section-head compact">
            <div>
              <p class="eyebrow">{{ t(isFamilyMode ? 'dashboard.netWorthTrend.familyEyebrow' : 'dashboard.netWorthTrend.eyebrow') }}</p>
              <h2>{{ t(isFamilyMode ? 'dashboard.netWorthTrend.familyTitle' : 'dashboard.netWorthTrend.title') }}</h2>
            </div>
            <div class="section-badge-row">
              <span class="report-range-badge">{{ activeStatsRangeLabel }}</span>
              <span class="chart-currency-badge">{{ netWorthCurrency }}</span>
            </div>
          </div>
          <ECharts v-if="netWorthTrendData.length" :options="netWorthTrendOptions" :height="'300px'" />
          <div v-else class="empty-state compact"><p>{{ t('dashboard.netWorthTrend.empty') }}</p></div>
        </div>

        <div class="overview-subcard">
          <div class="section-head compact">
            <div>
              <p class="eyebrow">{{ t(isFamilyMode ? 'dashboard.assetComposition.familyEyebrow' : 'dashboard.assetComposition.eyebrow') }}</p>
              <h2>{{ t('dashboard.assetComposition.currentTitle') }}</h2>
            </div>
            <div class="section-badge-row">
              <span class="snapshot-badge">{{ t('dashboard.reportContext.currentSnapshot') }}</span>
              <span class="chart-currency-badge">{{ netWorthCurrency }}</span>
            </div>
          </div>
          <ECharts v-if="selectedAssetComposition?.assetItems?.length" :options="assetCompositionOptions(selectedAssetComposition)" :height="'300px'" />
          <div v-else class="empty-state compact"><p>{{ t('dashboard.assetComposition.emptyAsset') }}</p></div>
          <div v-if="selectedAssetComposition?.liabilityItems?.length" class="liability-list">
            <p>{{ t('dashboard.assetComposition.liabilities') }}</p>
            <div v-for="item in selectedAssetComposition.liabilityItems" :key="item.key" class="liability-row">
              <span>{{ item.emoji }} {{ item.name }}</span>
              <strong>{{ formatMoney(item.amount, selectedAssetComposition.currency) }}</strong>
            </div>
          </div>
        </div>
      </div>
    </section>

    <section class="stats-panel reveal-block delay-2">
      <div class="section-head stats-head">
        <div>
          <p class="eyebrow">{{ t('statistics.eyebrow') }}</p>
          <h2>{{ isFamilyMode ? t('dashboard.familyStats.title') : t('statistics.title') }}</h2>
        </div>
        <div class="stats-range-switch" role="group" :aria-label="t('dashboard.reportContext.rangeAriaLabel')">
          <button v-for="preset in statsRangePresets" :key="preset.value" type="button" :class="{ active: statsRangePreset === preset.value }" :aria-pressed="statsRangePreset === preset.value" @click="switchStatsRangePreset(preset.value)">
            {{ t(preset.labelKey) }}
          </button>
        </div>
      </div>
      <div class="report-context-strip">
        <div class="report-context-main">
          <p class="report-context-kicker">{{ t('dashboard.reportContext.kicker') }}</p>
          <strong>{{ t('dashboard.reportContext.title') }} · {{ activeStatsRangeLabel }}</strong>
          <span>{{ t('dashboard.reportContext.description', { range: activeStatsRangeLabel }) }}</span>
        </div>
        <div class="report-context-modules" :aria-label="t('dashboard.reportContext.affectsLabel')">
          <small>{{ t('dashboard.reportContext.affectsLabel') }}</small>
          <span v-for="item in periodContextModules" :key="item">{{ item }}</span>
        </div>
      </div>
      <p class="stats-hint">{{ isFamilyMode ? t('dashboard.familyStats.hint') : t('statistics.originalCurrencyHint') }}</p>
      <div v-if="displayStatisticsSummaries.length" class="currency-stats-stack">
        <div v-for="summary in displayStatisticsSummaries" :key="summary.currency" class="currency-stats-card">
          <h3>{{ t('statistics.currencySection', { currency: summary.currency }) }}</h3>
          <div class="stats-grid">
            <component :is="isFamilyMode ? 'div' : 'button'" type="button" class="stat-tile income" :class="{ 'stat-tile-link': !isFamilyMode }" @click="!isFamilyMode && drillDownCurrency(summary.currency, 'income')"><span>{{ t('statistics.income') }}</span><strong>{{ formatMoney(summary.income, summary.currency) }}</strong></component>
            <component :is="isFamilyMode ? 'div' : 'button'" type="button" class="stat-tile expense" :class="{ 'stat-tile-link': !isFamilyMode }" @click="!isFamilyMode && drillDownCurrency(summary.currency, 'expense')"><span>{{ t('statistics.expense') }}</span><strong>{{ formatMoney(summary.netExpense || summary.expense, summary.currency) }}</strong><small>{{ t('statistics.grossExpense') }} {{ formatMoney(summary.grossExpense, summary.currency) }}</small></component>
            <component :is="isFamilyMode ? 'div' : 'button'" type="button" class="stat-tile" :class="{ 'stat-tile-link': !isFamilyMode }" @click="!isFamilyMode && drillDownCurrency(summary.currency, 'refund')"><span>{{ t('statistics.refund') }}</span><strong>{{ formatMoney(summary.refund, summary.currency) }}</strong></component>
            <div class="stat-tile"><span>{{ t('statistics.net') }}</span><strong>{{ formatMoney(summary.net, summary.currency) }}</strong></div>
            <component :is="isFamilyMode ? 'div' : 'button'" type="button" class="stat-tile" :class="{ 'stat-tile-link': !isFamilyMode }" @click="!isFamilyMode && drillDownCurrency(summary.currency, '')"><span>{{ t('statistics.eventCount') }}</span><strong>{{ summary.eventCount || 0 }}</strong></component>
          </div>
        </div>
      </div>
      <div v-else class="empty-state compact"><p>{{ isFamilyMode && !families.length ? t('dashboard.familySelector.empty') : t('statistics.emptyText') }}</p></div>
      <div class="category-mode-head">
        <div>
          <h3>{{ t('statistics.categoryGroups') }}</h3>
          <span class="report-range-badge inline">{{ activeStatsRangeLabel }}</span>
        </div>
        <div class="category-type-switch" role="group" :aria-label="t('statistics.categoryGroups')">
          <button type="button" :class="{ active: categoryStatsType === 'expense' }" @click="switchCategoryStatsType('expense')">{{ t('statistics.categoryTypes.netExpense') }}</button>
          <button type="button" :class="{ active: categoryStatsType === 'income' }" @click="switchCategoryStatsType('income')">{{ t('statistics.categoryTypes.income') }}</button>
        </div>
      </div>
      <div v-for="group in displayCategoryGroupStatsByCurrency" :key="group.currency" class="category-currency-card">
        <button class="category-currency-head" type="button" @click="toggleCategoryCurrency(group.currency)">
          <div>
            <strong>{{ t('statistics.categoryGroupsByCurrency', { currency: group.currency }) }}</strong>
            <small>{{ t('dashboard.categoryStats.collapsedHint', { count: (group.items || []).length, amount: formatMoney(categoryGroupTotal(group), group.currency) }) }}</small>
          </div>
          <span>{{ isCategoryCurrencyExpanded(group.currency) ? '▲' : '▼' }}</span>
        </button>
        <div v-if="isCategoryCurrencyExpanded(group.currency)" class="category-currency-body">
          <ECharts v-if="group.items?.length" :options="categoryDonutOptions(group)" :height="'260px'" />
          <component :is="isFamilyMode ? 'div' : 'button'" type="button" v-for="(item, idx) in visibleCategoryItems(group)" :key="`${group.currency}-${item.categoryGroupCode}`" class="rank-row" :class="{ 'rank-row-link': !isFamilyMode }" @click="!isFamilyMode && drillDownCategoryGroup(group.currency, item)">
            <span class="rank-dot" :style="`background:${categoryColorFor(item, idx)}`"></span>
            <strong>{{ item.categoryGroupName }}</strong>
            <em>{{ formatMoney(item.amount, group.currency) }}</em>
          </component>
          <button v-if="(group.items || []).length > 5" class="category-row-toggle" type="button" @click="toggleCategoryRows(group.currency)">
            {{ isCategoryRowsExpanded(group.currency) ? t('dashboard.categoryStats.collapseRows') : t('dashboard.categoryStats.expandAll') }}
          </button>
        </div>
      </div>
    </section>

    <section v-if="!isFamilyMode" class="trend-card reveal-block delay-2">
      <div class="section-head">
        <div>
          <p class="eyebrow">{{ t('dashboard.trend.eyebrow') }}</p>
          <h2>{{ t('dashboard.trend.title') }}</h2>
          <span class="report-range-badge inline">{{ activeStatsRangeLabel }}</span>
        </div>
        <div class="stats-range-switch" role="group" :aria-label="t('dashboard.trend.currencyLabel')">
          <button v-for="c in availableTrendCurrencies" :key="c" type="button" :class="{ active: trendCurrency === c }" :aria-pressed="trendCurrency === c" @click="switchTrendCurrency(c)">
            {{ c }}
          </button>
        </div>
      </div>
      <ECharts v-if="trendData.length" :options="trendOptions" :height="'320px'" />
      <div v-else class="empty-state compact">
        <p>{{ t('dashboard.trend.empty') }}</p>
      </div>
    </section>

    <section v-if="!isFamilyMode" class="investment-panel reveal-block delay-2">
      <div class="section-head">
        <div>
          <p class="eyebrow">{{ t('dashboard.investment.eyebrow') }}</p>
          <h2>{{ t('dashboard.investment.title') }}</h2>
          <span class="report-range-badge inline">{{ activeStatsRangeLabel }}</span>
        </div>
      </div>
      <p class="stats-hint">{{ t('dashboard.investment.hint') }}</p>
      <div v-if="investmentSummaries.length" class="investment-stack">
        <div v-for="summary in investmentSummaries" :key="summary.currency" class="investment-card" :class="{ gain: Number(summary.profitLoss) >= 0, loss: Number(summary.profitLoss) < 0 }">
          <div class="investment-card-head">
            <span>{{ summary.currency }}</span>
            <strong>{{ formatMoney(summary.profitLoss, summary.currency) }}</strong>
            <small>{{ t('dashboard.investment.profitLoss') }}</small>
          </div>
          <div class="investment-metrics">
            <div><span>{{ t('dashboard.investment.marketValue') }}</span><strong>{{ formatMoney(summary.marketValue, summary.currency) }}</strong></div>
            <div><span>{{ t('dashboard.investment.buyAmount') }}</span><strong>{{ formatMoney(summary.buyAmount, summary.currency) }}</strong></div>
            <div><span>{{ t('dashboard.investment.revalue') }}</span><strong>{{ formatMoney(summary.revalueAmount, summary.currency) }}</strong></div>
            <div><span>{{ t('dashboard.investment.income') }}</span><strong>{{ formatMoney(summary.incomeAmount, summary.currency) }}</strong></div>
          </div>
          <div v-if="summary.buckets?.length" class="investment-bucket-list">
            <div v-for="bucket in summary.buckets" :key="bucket.bucketId" class="investment-bucket-row">
              <span class="investment-bucket-name">{{ getBucketEmoji(bucket.bucketType) }} {{ bucket.bucketName }}</span>
              <div class="investment-bucket-metric">
                <small>{{ t('dashboard.investment.marketValue') }}</small>
                <strong>{{ formatMoney(bucket.marketValue, bucket.currency) }}</strong>
              </div>
              <div class="investment-bucket-metric">
                <small>{{ t('dashboard.investment.profitLoss') }}</small>
                <em :class="{ gain: Number(bucket.profitLoss) >= 0, loss: Number(bucket.profitLoss) < 0 }">{{ formatMoney(bucket.profitLoss, bucket.currency) }}</em>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div v-else class="empty-state compact"><p>{{ t('dashboard.investment.empty') }}</p></div>
    </section>

    <section class="workspace-grid reveal-block delay-2">
      <div v-if="!isFamilyMode" class="surface-card asset-surface">
        <div class="section-head">
          <div>
            <p class="eyebrow">{{ t('dashboard.buckets.eyebrow') }}</p>
            <h2>{{ t('dashboard.buckets.title') }}</h2>
            <span class="snapshot-badge inline">{{ t('dashboard.reportContext.currentSnapshot') }}</span>
          </div>
          <button class="ghost-action" @click="router.push('/buckets')">{{ t('common.actions.viewAll') }}</button>
        </div>
        <div v-if="currencyTotals.length" class="currency-stack">
          <div v-for="item in currencyTotals" :key="item.currency" class="currency-row">
            <div><span class="currency-code">{{ item.currency }}</span><small>{{ t('dashboard.buckets.count', { count: item.count }) }}</small></div>
            <strong>{{ formatAmount(item.total) }}</strong>
          </div>
        </div>
        <div v-else class="empty-state compact">
          <img :src="marmotTwo" :alt="t('dashboard.buckets.emptyAlt')" width="92" height="92" />
          <p>{{ t('dashboard.buckets.emptyText') }}</p>
        </div>
      </div>

      <div v-else class="surface-card family-asset-surface">
        <div class="section-head">
          <div>
            <p class="eyebrow">{{ t('dashboard.familyAssets.eyebrow') }}</p>
            <h2>{{ t('dashboard.familyAssets.title') }}</h2>
            <span class="snapshot-badge inline">{{ t('dashboard.reportContext.currentSnapshot') }}</span>
          </div>
          <button class="ghost-action" @click="router.push('/family')">{{ t('dashboard.familySelector.manage') }}</button>
        </div>
        <p class="stats-hint">{{ t('dashboard.familyAssets.hint') }}</p>
        <div v-if="familyAssetsMembers.length" class="family-assets-stack">
          <div v-for="member in familyAssetsMembers" :key="member.userId" class="member-asset-card">
            <button class="member-head-button" type="button" :aria-expanded="isMemberAssetExpanded(member)" :aria-controls="`member-assets-${member.userId}`" :aria-label="isMemberAssetExpanded(member) ? t('dashboard.familyAssets.collapseMemberAria', { name: member.displayName || member.name || member.account }) : t('dashboard.familyAssets.expandMemberAria', { name: member.displayName || member.name || member.account })" @click="toggleMemberAsset(member)">
              <div class="member-head-main">
                <strong>{{ member.displayName || member.name || member.account }}</strong>
                <span>{{ member.role }}</span>
              </div>
              <div class="member-summary-meta">
                <span>{{ t('dashboard.familyAssets.accountsCount', { count: memberAssetSummary(member).accountCount }) }}</span>
                <span>{{ t('dashboard.familyAssets.bucketsCount', { count: memberAssetSummary(member).bucketCount }) }}</span>
                <span>{{ t('dashboard.familyAssets.currenciesCount', { count: memberAssetSummary(member).currencyCount }) }}</span>
              </div>
              <span class="member-collapse-icon">{{ isMemberAssetExpanded(member) ? '−' : '+' }}</span>
            </button>
            <div v-if="member.totals?.length" class="asset-total-grid compact">
              <div v-for="total in member.totals" :key="`${member.userId}-${total.currency}`" class="asset-total-pill">
                <span>{{ total.currency }}</span>
                <strong>{{ t('dashboard.familyAssets.netWorth') }} {{ formatMoney(total.netWorth, total.currency) }}</strong>
                <small>{{ t('domain.asset') }} {{ formatMoney(total.asset, total.currency) }} · {{ t('domain.liability') }} {{ formatMoney(total.liability, total.currency) }}</small>
              </div>
            </div>
            <div v-if="isMemberAssetExpanded(member)" :id="`member-assets-${member.userId}`" class="member-asset-body">
              <div v-if="member.accounts?.length">
                <div v-for="account in member.accounts || []" :key="account.id" class="family-account-row">
                  <strong>{{ account.name }}</strong>
                  <div v-for="bucket in account.buckets || []" :key="bucket.id" class="family-bucket-row">
                    <span>{{ getBucketEmoji(bucket.bucketType) }} {{ bucket.name }} · {{ bucket.bucketNature === 'liability' ? t('domain.liability') : t('domain.asset') }}</span>
                    <em>{{ formatMoney(bucket.balance, bucket.currency) }}</em>
                  </div>
                </div>
              </div>
              <p v-else class="member-empty">{{ t('dashboard.familyAssets.emptyAccounts') }}</p>
            </div>
          </div>
        </div>
        <div v-else class="empty-state compact"><p>{{ t('dashboard.familyAssets.emptyText') }}</p></div>
      </div>

      <div class="surface-card events-surface">
        <div class="section-head">
          <div>
            <p class="eyebrow">{{ isFamilyMode ? t('dashboard.familyEvents.eyebrow') : t('dashboard.events.eyebrow') }}</p>
            <h2>{{ isFamilyMode ? t('dashboard.familyEvents.title') : t('dashboard.events.title') }}</h2>
            <span class="snapshot-badge inline latest">{{ t('dashboard.reportContext.latestSnapshot') }}</span>
          </div>
          <button class="ghost-action" @click="refreshAll">{{ t('common.actions.refresh') }}</button>
        </div>
        <p class="stats-hint">{{ t('dashboard.reportContext.latestEventsHint') }}</p>
        <div v-if="displayRecentEvents.length" class="event-list">
          <component :is="isFamilyMode ? 'div' : 'button'" type="button" v-for="event in displayRecentEvents" :key="event.id" class="event-row" :class="{ 'event-row-link': !isFamilyMode }" @click="!isFamilyMode && drillDownEvent(event)">
            <div class="event-dot"></div>
            <div class="event-main">
              <strong>{{ event.description || eventTypeLabel(event.eventType) }}</strong>
              <span>{{ event.eventTime }} · {{ eventTypeLabel(event.eventType) }} · {{ event.remark || t('common.misc.noRemark') }}</span>
            </div>
            <div class="event-amount">{{ formatMoney(event.amount, event.currency) }}</div>
          </component>
        </div>
        <div v-else class="empty-state compact">
          <img :src="marmotOne" :alt="t('dashboard.events.emptyAlt')" width="92" height="92" />
          <p>{{ isFamilyMode ? t('dashboard.familyEvents.emptyText') : t('dashboard.events.emptyText') }}</p>
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
import { getInvestmentSummaries, getNetWorthTrend, getStatisticsCategoryGroupsByCurrency, getStatisticsSummaries, getStatisticsTrend } from '@/api/statistics/statistics'
import { getFamilyAssets, getFamilyNetWorthTrend, getFamilyStatisticsCategoryGroupsByCurrency, getFamilyStatisticsSummaries, listFamilies, listFamilyFinancialEvents, listFamilyMembers } from '@/api/family/family'
import { getBucketEmoji } from '@/utils/bucketEmoji'
import ECharts from '@/components/ECharts.vue'
import marmotOne from '../../../img/marmot-ledger-1.png'
import marmotTwo from '../../../img/marmot-ledger-2.png'

const router = useRouter()
const { t, te } = useI18n()
const config = useConfigStore()
const viewMode = ref('personal')
const accounts = ref([])
const buckets = ref([])
const recentEvents = ref([])
const statisticsSummaries = ref([])
const categoryGroupStatsByCurrency = ref([])
const families = ref([])
const selectedFamilyId = ref('')
const familyMembers = ref([])
const familyRecentEvents = ref([])
const familyStatisticsSummaries = ref([])
const familyCategoryGroupStatsByCurrency = ref([])
const familyAssets = ref(null)
const expandedMemberId = ref('')
const categoryStatsType = ref('expense')
const statsRangePreset = ref('recentYear')
const trendCurrency = ref('CNY')
const trendData = ref([])
const netWorthCurrency = ref('CNY')
const netWorthTrendData = ref([])
const investmentSummaries = ref([])
const expandedCategoryCurrencies = ref([])
const expandedCategoryRows = ref({})
const loading = ref(false)

const statsRangePresets = [
  { value: 'recentMonth', labelKey: 'statistics.rangePresets.recentMonth' },
  { value: 'thisMonth', labelKey: 'statistics.rangePresets.thisMonth' },
  { value: 'recentYear', labelKey: 'statistics.rangePresets.recentYear' },
  { value: 'thisYear', labelKey: 'statistics.rangePresets.thisYear' },
]

const isFamilyMode = computed(() => viewMode.value === 'family')
const selectedFamily = computed(() => families.value.find((item) => item.id === selectedFamilyId.value))
const familyAssetsMembers = computed(() => familyAssets.value?.members || [])

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

const displayStatisticsSummaries = computed(() => isFamilyMode.value ? familyStatisticsSummaries.value : statisticsSummaries.value)
const displayCategoryGroupStatsByCurrency = computed(() => isFamilyMode.value ? familyCategoryGroupStatsByCurrency.value : categoryGroupStatsByCurrency.value)
const displayRecentEvents = computed(() => isFamilyMode.value ? familyRecentEvents.value : recentEvents.value)
const activeStatsRange = computed(() => statsRangePresets.find((item) => item.value === statsRangePreset.value) || statsRangePresets[0])
const activeStatsRangeLabel = computed(() => t(activeStatsRange.value.labelKey))
const periodContextModules = computed(() => [
  t('dashboard.reportContext.modules.statistics'),
  t('dashboard.reportContext.modules.categories'),
  t('dashboard.reportContext.modules.trends'),
  t('dashboard.reportContext.modules.netWorthTrend'),
  t('dashboard.reportContext.modules.investments'),
])

const displayMetrics = computed(() => {
  const currency = netWorthCurrency.value || availableNetWorthCurrencies.value[0] || 'CNY'
  const composition = assetCompositionByCurrency.value.find((item) => item.currency === currency) || { currency, asset: 0, liability: 0, netWorth: 0 }
  const stat = displayStatisticsSummaries.value.find((item) => item.currency === currency) || { net: 0, eventCount: 0 }
  return [
    { label: isFamilyMode.value ? t('dashboard.kpi.familyNetWorth') : t('dashboard.kpi.netWorth'), value: formatMoney(composition.netWorth, currency), hint: t('dashboard.kpi.netWorthHint'), tone: Number(composition.netWorth) >= 0 ? 'positive' : 'negative' },
    { label: t('dashboard.kpi.assets'), value: formatMoney(composition.asset, currency), hint: t('dashboard.kpi.assetsHint'), tone: 'positive' },
    { label: t('dashboard.kpi.liabilities'), value: formatMoney(composition.liability, currency), hint: t('dashboard.kpi.liabilitiesHint'), tone: Number(composition.liability) > 0 ? 'negative' : '' },
    { label: t('dashboard.kpi.periodNet'), value: formatMoney(stat.net || 0, currency), hint: t('dashboard.kpi.periodNetHint'), tone: Number(stat.net || 0) >= 0 ? 'positive' : 'negative' },
  ]
})

const eventTypeLabel = (type) => {
  const key = `record.scenarios.${type}`
  return type && te(key) ? t(key) : type
}

const formatAmount = (value) => new Intl.NumberFormat(config.locale, { minimumFractionDigits: 2, maximumFractionDigits: 2 }).format(Number(value || 0))
const formatMoney = (value, currency) => `${currency} ${formatAmount(value)}`

function memberAssetSummary(member) {
  const accounts = member.accounts || []
  const bucketCount = accounts.reduce((sum, account) => sum + Number(account.buckets?.length || 0), 0)
  return { accountCount: accounts.length, bucketCount, currencyCount: member.totals?.length || 0 }
}

function isMemberAssetExpanded(member) {
  return String(member.userId) === String(expandedMemberId.value)
}

function resetExpandedMember() {
  expandedMemberId.value = familyAssetsMembers.value[0]?.userId || ''
}

function toggleMemberAsset(member) {
  expandedMemberId.value = isMemberAssetExpanded(member) ? '' : member.userId
}

function padDatePart(value) {
  return String(value).padStart(2, '0')
}

function formatDateTime(date, endOfDay = false) {
  const hours = endOfDay ? '23' : '00'
  const minutes = endOfDay ? '59' : '00'
  const seconds = endOfDay ? '59' : '00'
  return `${date.getFullYear()}-${padDatePart(date.getMonth() + 1)}-${padDatePart(date.getDate())} ${hours}:${minutes}:${seconds}`
}

function lastDayOfMonth(year, monthIndex) {
  return new Date(year, monthIndex + 1, 0).getDate()
}

function subtractMonthsClamped(date, months) {
  const targetMonthIndex = date.getMonth() - months
  const target = new Date(date.getFullYear(), targetMonthIndex, 1)
  const day = Math.min(date.getDate(), lastDayOfMonth(target.getFullYear(), target.getMonth()))
  return new Date(target.getFullYear(), target.getMonth(), day)
}

function subtractYearsClamped(date, years) {
  const year = date.getFullYear() - years
  const month = date.getMonth()
  const day = Math.min(date.getDate(), lastDayOfMonth(year, month))
  return new Date(year, month, day)
}

function statsRangeParams() {
  const now = new Date()
  let start
  let end
  switch (statsRangePreset.value) {
    case 'recentYear':
      start = subtractYearsClamped(now, 1)
      end = now
      break
    case 'thisYear':
      start = new Date(now.getFullYear(), 0, 1)
      end = new Date(now.getFullYear(), 11, 31)
      break
    case 'recentMonth':
      start = subtractMonthsClamped(now, 1)
      end = now
      break
    case 'thisMonth':
    default:
      start = new Date(now.getFullYear(), now.getMonth(), 1)
      end = new Date(now.getFullYear(), now.getMonth(), lastDayOfMonth(now.getFullYear(), now.getMonth()))
      break
  }
  return { startTime: formatDateTime(start), endTime: formatDateTime(end, true) }
}

async function refreshPersonalStatistics() {
  const range = statsRangeParams()
  const [summaryRes, categoryStatsRes] = await Promise.all([
    getStatisticsSummaries(range),
    getStatisticsCategoryGroupsByCurrency({ ...range, type: categoryStatsType.value }),
  ])
  if (summaryRes.success) statisticsSummaries.value = summaryRes.data || []
  if (categoryStatsRes.success) {
    categoryGroupStatsByCurrency.value = categoryStatsRes.data || []
    ensureCategoryExpansion()
  }
}

async function refreshFamilyStatistics(familyId = selectedFamilyId.value) {
  if (!familyId) {
    familyStatisticsSummaries.value = []
    familyCategoryGroupStatsByCurrency.value = []
    return
  }
  const range = statsRangeParams()
  const [summaryRes, categoryStatsRes] = await Promise.all([
    getFamilyStatisticsSummaries(familyId, range),
    getFamilyStatisticsCategoryGroupsByCurrency(familyId, { ...range, type: categoryStatsType.value }),
  ])
  if (summaryRes.success) familyStatisticsSummaries.value = summaryRes.data || []
  if (categoryStatsRes.success) {
    familyCategoryGroupStatsByCurrency.value = categoryStatsRes.data || []
    ensureCategoryExpansion()
  }
}

async function refreshCurrentCategoryStats() {
  const range = statsRangeParams()
  if (isFamilyMode.value) {
    if (!selectedFamilyId.value) {
      familyCategoryGroupStatsByCurrency.value = []
      return
    }
    const res = await getFamilyStatisticsCategoryGroupsByCurrency(selectedFamilyId.value, { ...range, type: categoryStatsType.value })
    if (res.success) {
      familyCategoryGroupStatsByCurrency.value = res.data || []
      ensureCategoryExpansion()
    }
  } else {
    const res = await getStatisticsCategoryGroupsByCurrency({ ...range, type: categoryStatsType.value })
    if (res.success) {
      categoryGroupStatsByCurrency.value = res.data || []
      ensureCategoryExpansion()
    }
  }
}

async function refreshPersonalDashboard() {
  const [accountRes, bucketRes, eventRes] = await Promise.all([
    listAccounts(),
    listBuckets(),
    listFinancialEvents({ page: 1, pageSize: 5 }),
  ])
  if (accountRes.success) accounts.value = accountRes.data || []
  if (bucketRes.success) buckets.value = bucketRes.data || []
  if (eventRes.success) recentEvents.value = eventRes.data?.list || []
  await refreshPersonalStatistics()
}

async function refreshFamilyCatalog() {
  const res = await listFamilies()
  if (res.success) {
    families.value = res.data || []
    if (!selectedFamilyId.value && families.value.length) selectedFamilyId.value = families.value[0].id
  }
}

async function refreshFamilyDashboard(familyId = selectedFamilyId.value) {
  if (!familyId) {
    familyMembers.value = []
    familyRecentEvents.value = []
    familyStatisticsSummaries.value = []
    familyCategoryGroupStatsByCurrency.value = []
    familyAssets.value = null
    expandedMemberId.value = ''
    return
  }
  const [memberRes, eventRes, assetRes] = await Promise.all([
    listFamilyMembers(familyId, { includeInvited: false }),
    listFamilyFinancialEvents(familyId, { page: 1, pageSize: 5 }),
    getFamilyAssets(familyId),
  ])
  if (memberRes.success) familyMembers.value = memberRes.data || []
  if (eventRes.success) familyRecentEvents.value = eventRes.data?.list || []
  if (assetRes.success) {
    familyAssets.value = assetRes.data
    resetExpandedMember()
  }
  await refreshFamilyStatistics(familyId)
}

async function refreshAll() {
  loading.value = true
  try {
    await Promise.all([refreshPersonalDashboard(), refreshFamilyCatalog()])
    if (isFamilyMode.value) await refreshFamilyDashboard()
    if (!isFamilyMode.value) {
      const first = statisticsSummaries.value[0]?.currency
      if (first && !availableTrendCurrencies.value.includes(trendCurrency.value)) {
        trendCurrency.value = first
      } else if (first && trendCurrency.value === 'CNY' && !statisticsSummaries.value.some((s) => s.currency === 'CNY')) {
        trendCurrency.value = first
      }
      const firstNetWorthCurrency = availableNetWorthCurrencies.value[0]
      if (firstNetWorthCurrency && !availableNetWorthCurrencies.value.includes(netWorthCurrency.value)) netWorthCurrency.value = firstNetWorthCurrency
      await refreshTrend()
      await refreshNetWorthTrend()
      await refreshInvestmentSummaries()
    }
  } finally {
    loading.value = false
  }
}

async function switchMode(mode) {
  viewMode.value = mode
  if (mode === 'family') {
    await refreshFamilyCatalog()
    await refreshFamilyDashboard()
  }
  const firstNetWorthCurrency = availableNetWorthCurrencies.value[0]
  if (firstNetWorthCurrency && !availableNetWorthCurrencies.value.includes(netWorthCurrency.value)) netWorthCurrency.value = firstNetWorthCurrency
  await refreshNetWorthTrend()
}

async function selectFamilyWorkspace(familyId) {
  if (String(selectedFamilyId.value) === String(familyId)) return
  selectedFamilyId.value = familyId
  expandedMemberId.value = ''
  await refreshFamilyDashboard(familyId)
  const firstNetWorthCurrency = availableNetWorthCurrencies.value[0]
  if (firstNetWorthCurrency && !availableNetWorthCurrencies.value.includes(netWorthCurrency.value)) netWorthCurrency.value = firstNetWorthCurrency
  await refreshNetWorthTrend()
}

async function switchCategoryStatsType(type) {
  if (categoryStatsType.value === type) return
  categoryStatsType.value = type
  await refreshCurrentCategoryStats()
}

async function switchStatsRangePreset(preset) {
  statsRangePreset.value = preset
  if (isFamilyMode.value) {
    await refreshFamilyStatistics()
    await refreshNetWorthTrend()
  } else {
    await refreshPersonalStatistics()
    await refreshTrend()
    await refreshNetWorthTrend()
    await refreshInvestmentSummaries()
  }
}

const trendGranularity = computed(() => {
  if (statsRangePreset.value === 'recentMonth' || statsRangePreset.value === 'thisMonth') return 'week'
  return 'month'
})

const availableTrendCurrencies = computed(() => {
  const list = (statisticsSummaries.value || []).map((s) => s.currency).filter(Boolean)
  return list.length ? list : ['CNY']
})

async function refreshTrend() {
  if (!trendCurrency.value) {
    trendData.value = []
    return
  }
  const range = statsRangeParams()
  const res = await getStatisticsTrend({ ...range, currency: trendCurrency.value, granularity: trendGranularity.value })
  if (res.success) trendData.value = res.data || []
  else trendData.value = []
}

async function refreshInvestmentSummaries() {
  const range = statsRangeParams()
  const res = await getInvestmentSummaries(range)
  if (res.success) investmentSummaries.value = res.data || []
  else investmentSummaries.value = []
}

async function switchTrendCurrency(c) {
  if (trendCurrency.value === c) return
  trendCurrency.value = c
  await refreshTrend()
}

function formatPeriod(period) {
  if (!period) return ''
  if (period.includes('-W')) {
    const [, week] = period.split('-W')
    return `W${week}`
  }
  const [year, month] = period.split('-')
  return `${year}-${month}`
}

const trendOptions = computed(() => {
  const periods = trendData.value.map((p) => formatPeriod(p.period))
  const income = trendData.value.map((p) => Number(p.income || 0))
  const expense = trendData.value.map((p) => Number(p.netExpense || 0))
  const net = trendData.value.map((p) => Number(p.income || 0) - Number(p.netExpense || 0))
  return {
    grid: { left: 56, right: 24, top: 30, bottom: 40 },
    tooltip: { trigger: 'axis', backgroundColor: '#1f2933', textStyle: { color: '#fff', fontWeight: 700 }, borderWidth: 0 },
    legend: { data: [t('statistics.income'), t('statistics.expense'), t('statistics.net')], bottom: 0, textStyle: { color: '#64748b', fontWeight: 700 } },
    xAxis: { type: 'category', data: periods, boundaryGap: false, axisLine: { lineStyle: { color: '#e2e8f0' } }, axisLabel: { color: '#64748b', fontWeight: 700 } },
    yAxis: { type: 'value', axisLine: { show: false }, axisTick: { show: false }, axisLabel: { color: '#94a3b8', fontWeight: 700 }, splitLine: { lineStyle: { color: '#f1f5f9' } } },
    series: [
      { name: t('statistics.income'), type: 'line', smooth: true, data: income, lineStyle: { color: '#ef4444', width: 3 }, itemStyle: { color: '#ef4444' }, areaStyle: { color: 'rgba(239, 68, 68, 0.08)' }, symbol: 'circle', symbolSize: 6 },
      { name: t('statistics.expense'), type: 'line', smooth: true, data: expense, lineStyle: { color: '#f97316', width: 3 }, itemStyle: { color: '#f97316' }, areaStyle: { color: 'rgba(249, 115, 22, 0.08)' }, symbol: 'circle', symbolSize: 6 },
      { name: t('statistics.net'), type: 'line', smooth: true, data: net, lineStyle: { color: '#2f7d5c', width: 3, type: 'dashed' }, itemStyle: { color: '#2f7d5c' }, symbol: 'circle', symbolSize: 6 },
    ],
  }
})

const netWorthGranularity = computed(() => {
  if (statsRangePreset.value === 'recentMonth' || statsRangePreset.value === 'thisMonth') return 'week'
  return 'month'
})

const familyAssetBuckets = computed(() => familyAssetsMembers.value.flatMap((member) =>
  (member.accounts || []).flatMap((account) =>
    (account.buckets || []).map((bucket) => ({ ...bucket, memberName: member.displayName || member.name || member.account }))
  )
))

const allVisibleBuckets = computed(() => isFamilyMode.value ? familyAssetBuckets.value : buckets.value)

const availableNetWorthCurrencies = computed(() => {
  const set = new Set()
  ;(allVisibleBuckets.value || []).forEach((b) => { if (b.currency) set.add(b.currency) })
  return Array.from(set).sort()
})

function bucketTypeColor(type) {
  const colors = {
    cash: '#2f7d5c', wallet: '#3b82f6', bank: '#06b6d4', credit: '#ef4444',
    investment_cash: '#7c3aed', investment_asset: '#8b5cf6', receivable: '#14b8a6',
    deposit: '#f59e0b', loan_out: '#84cc16', liability: '#ef4444', virtual: '#94a3b8',
  }
  return colors[type] || '#94a3b8'
}

function getBucketTypeLabel(type) {
  const key = `buckets.types.${type?.replace(/_([a-z])/g, (_, c) => c.toUpperCase())}`
  return type && te(key) ? t(key) : type
}

const assetCompositionByCurrency = computed(() => {
  const map = new Map()
  ;(allVisibleBuckets.value || []).forEach((bucket) => {
    if (!bucket.currency) return
    if (!map.has(bucket.currency)) {
      map.set(bucket.currency, { currency: bucket.currency, asset: 0, liability: 0, netWorth: 0, assetItems: [], liabilityItems: [] })
    }
    const group = map.get(bucket.currency)
    const value = Number(bucket.balance || 0)
    const key = isFamilyMode.value ? bucket.bucketType : bucket.id
    const name = isFamilyMode.value ? getBucketTypeLabel(bucket.bucketType) : bucket.name
    const emoji = getBucketEmoji(bucket.bucketType)
    const target = bucket.bucketNature === 'liability' ? group.liabilityItems : group.assetItems
    let item = target.find((x) => x.key === key)
    if (!item) {
      item = { key, id: bucket.id, name, bucketType: bucket.bucketType, amount: 0, color: bucketTypeColor(bucket.bucketType), emoji }
      target.push(item)
    }
    item.amount += value
    if (bucket.bucketNature === 'liability') group.liability += value
    else group.asset += value
    group.netWorth = group.asset - group.liability
  })
  return Array.from(map.values()).sort((a, b) => b.netWorth - a.netWorth)
})

const selectedAssetComposition = computed(() => assetCompositionByCurrency.value.find((g) => g.currency === netWorthCurrency.value))

function assetCompositionOptions(group) {
  const data = (group.assetItems || []).filter((item) => item.amount > 0).map((item) => ({
    value: item.amount,
    name: `${item.emoji} ${item.name}`,
    itemStyle: { color: item.color },
  }))
  return {
    tooltip: { trigger: 'item', backgroundColor: '#1f2933', textStyle: { color: '#fff', fontWeight: 700 }, borderWidth: 0, formatter: '{b}<br/>{c} ({d}%)' },
    legend: { type: 'scroll', orient: 'horizontal', bottom: 0, textStyle: { color: '#64748b', fontWeight: 700 } },
    series: [{
      type: 'pie',
      radius: ['58%', '80%'],
      center: ['50%', '42%'],
      itemStyle: { borderRadius: 6, borderColor: '#fff', borderWidth: 2 },
      label: { show: false },
      labelLine: { show: false },
      data,
    }],
    graphic: [{
      type: 'group',
      left: 'center',
      top: '32%',
      children: [
        { type: 'text', left: 'center', top: 0, style: { text: t('dashboard.assetComposition.netWorth'), fontSize: 11, fill: '#94a3b8', fontWeight: 800 } },
        { type: 'text', left: 'center', top: 18, style: { text: `${group.currency} ${formatAmount(group.netWorth)}`, fontSize: 18, fill: '#1e293b', fontWeight: 900 } },
      ],
    }],
  }
}

async function refreshNetWorthTrend() {
  if (!netWorthCurrency.value) {
    netWorthTrendData.value = []
    return
  }
  const range = statsRangeParams()
  const params = { ...range, currency: netWorthCurrency.value, granularity: netWorthGranularity.value }
  const res = isFamilyMode.value && selectedFamilyId.value
    ? await getFamilyNetWorthTrend(selectedFamilyId.value, params)
    : await getNetWorthTrend(params)
  if (res.success) netWorthTrendData.value = res.data || []
  else netWorthTrendData.value = []
}

async function switchNetWorthCurrency(c) {
  if (netWorthCurrency.value === c) return
  netWorthCurrency.value = c
  await refreshNetWorthTrend()
}

const netWorthTrendOptions = computed(() => {
  const periods = netWorthTrendData.value.map((p) => formatPeriod(p.period))
  const asset = netWorthTrendData.value.map((p) => Number(p.asset || 0))
  const liability = netWorthTrendData.value.map((p) => Number(p.liability || 0))
  const netWorth = netWorthTrendData.value.map((p) => Number(p.netWorth || 0))
  return {
    grid: { left: 56, right: 24, top: 30, bottom: 40 },
    tooltip: { trigger: 'axis', backgroundColor: '#1f2933', textStyle: { color: '#fff', fontWeight: 700 }, borderWidth: 0 },
    legend: { data: [t('dashboard.netWorthTrend.asset'), t('dashboard.netWorthTrend.liability'), t('dashboard.netWorthTrend.netWorth')], bottom: 0, textStyle: { color: '#64748b', fontWeight: 700 } },
    xAxis: { type: 'category', data: periods, boundaryGap: false, axisLine: { lineStyle: { color: '#e2e8f0' } }, axisLabel: { color: '#64748b', fontWeight: 700 } },
    yAxis: { type: 'value', axisLine: { show: false }, axisTick: { show: false }, axisLabel: { color: '#94a3b8', fontWeight: 700 }, splitLine: { lineStyle: { color: '#f1f5f9' } } },
    series: [
      { name: t('dashboard.netWorthTrend.asset'), type: 'line', smooth: true, data: asset, lineStyle: { color: '#2f7d5c', width: 3 }, itemStyle: { color: '#2f7d5c' }, areaStyle: { color: 'rgba(47, 125, 92, 0.08)' }, symbol: 'circle', symbolSize: 6 },
      { name: t('dashboard.netWorthTrend.liability'), type: 'line', smooth: true, data: liability, lineStyle: { color: '#ef4444', width: 2 }, itemStyle: { color: '#ef4444' }, symbol: 'circle', symbolSize: 5 },
      { name: t('dashboard.netWorthTrend.netWorth'), type: 'line', smooth: true, data: netWorth, lineStyle: { color: '#3b82f6', width: 4 }, itemStyle: { color: '#3b82f6' }, symbol: 'circle', symbolSize: 7 },
    ],
  }
})

const categoryFallbackColors = ['#2f7d5c', '#3b82f6', '#ef4444', '#f59e0b', '#7c3aed', '#06b6d4', '#84cc16', '#ec4899', '#14b8a6', '#a855f7']

function categoryColorFor(item, idx) {
  return item.color || categoryFallbackColors[idx % categoryFallbackColors.length]
}

function categoryGroupTotal(group) {
  return (group.items || []).reduce((sum, item) => sum + Number(item.amount || 0), 0)
}

function ensureCategoryExpansion() {
  if (expandedCategoryCurrencies.value.length || !displayCategoryGroupStatsByCurrency.value.length) return
  expandedCategoryCurrencies.value = [displayCategoryGroupStatsByCurrency.value[0].currency]
}

function isCategoryCurrencyExpanded(currency) {
  return expandedCategoryCurrencies.value.includes(currency)
}

function toggleCategoryCurrency(currency) {
  if (isCategoryCurrencyExpanded(currency)) {
    expandedCategoryCurrencies.value = expandedCategoryCurrencies.value.filter((item) => item !== currency)
  } else {
    expandedCategoryCurrencies.value = [...expandedCategoryCurrencies.value, currency]
  }
}

function isCategoryRowsExpanded(currency) {
  return Boolean(expandedCategoryRows.value[currency])
}

function toggleCategoryRows(currency) {
  expandedCategoryRows.value = { ...expandedCategoryRows.value, [currency]: !expandedCategoryRows.value[currency] }
}

function visibleCategoryItems(group) {
  const items = group.items || []
  return isCategoryRowsExpanded(group.currency) ? items : items.slice(0, 5)
}

function categoryDonutOptions(group) {
  const items = group.items || []
  const total = items.reduce((s, it) => s + Number(it.amount || 0), 0)
  return {
    tooltip: { trigger: 'item', backgroundColor: '#1f2933', textStyle: { color: '#fff', fontWeight: 700 }, borderWidth: 0, formatter: '{b}<br/>{c} ({d}%)' },
    legend: { type: 'scroll', orient: 'horizontal', bottom: 0, textStyle: { color: '#64748b', fontWeight: 700 } },
    series: [{
      type: 'pie',
      radius: ['58%', '80%'],
      center: ['50%', '42%'],
      avoidLabelOverlap: true,
      itemStyle: { borderRadius: 6, borderColor: '#fff', borderWidth: 2 },
      label: { show: false },
      labelLine: { show: false },
      data: items.map((it, idx) => ({
        value: Number(it.amount || 0),
        name: it.categoryGroupName,
        itemStyle: { color: categoryColorFor(it, idx) },
      })),
      emphasis: { label: { show: true, fontSize: 14, fontWeight: 700 } },
    }],
    graphic: [{
      type: 'group',
      left: 'center',
      top: '32%',
      children: [
        { type: 'text', left: 'center', top: 0, style: { text: t(categoryStatsType.value === 'income' ? 'statistics.income' : 'statistics.categoryTypes.netExpense'), fontSize: 11, fill: '#94a3b8', fontWeight: 800 } },
        { type: 'text', left: 'center', top: 18, style: { text: `${group.currency} ${formatAmount(total)}`, fontSize: 18, fill: '#1e293b', fontWeight: 900 } },
      ],
    }],
  }
}

function drillDownCurrency(currency, eventType) {
  const query = { rangeMode: 'preset', rangePreset: statsRangePreset.value }
  if (currency) query.currency = currency
  if (eventType) query.eventType = eventType
  router.push({ path: '/records', query })
}

function drillDownCategoryGroup(currency, item) {
  if (!item) return
  const query = { rangeMode: 'preset', rangePreset: statsRangePreset.value, eventType: categoryStatsType.value }
  if (currency) query.currency = currency
  if (item.categoryGroupId) query.categoryGroupId = String(item.categoryGroupId)
  if (item.categoryGroupName) query.categoryGroupName = item.categoryGroupName
  router.push({ path: '/records', query })
}

function drillDownEvent(event) {
  if (!event?.id) return
  router.push({ path: '/records', query: { eventId: String(event.id), rangeMode: 'preset', rangePreset: statsRangePreset.value } })
}

onMounted(refreshAll)
onActivated(refreshAll)
</script>

<style scoped>
.dashboard-shell{max-width:1200px;margin:0 auto;color:#1e293b}.reveal-block{animation:revealUp 520ms cubic-bezier(.16,1,.3,1) both}.delay-1{animation-delay:90ms}.delay-2{animation-delay:180ms}.cockpit-panel{display:grid;grid-template-columns:minmax(0,1fr) 280px;gap:24px;margin-bottom:24px;padding:28px;background:linear-gradient(135deg,#fffaf0 0%,#fff 58%,#eef7f0 100%);border-radius:22px;box-shadow:0 2px 8px rgba(15,23,42,.06),0 18px 48px rgba(47,125,92,.12)}.brand-cluster{display:grid;grid-template-columns:210px minmax(0,1fr);gap:28px;align-items:center}.brand-orbit{position:relative;min-height:170px}.brand-orbit img:first-child{position:absolute;left:18px;top:8px;border-radius:28px;box-shadow:0 16px 36px rgba(47,125,92,.18)}.brand-orbit img:last-child{position:absolute;right:8px;bottom:14px;border-radius:22px;box-shadow:0 12px 28px rgba(31,41,51,.16)}.eyebrow{margin:0 0 8px;color:#2f7d5c;font-size:12px;font-weight:800;letter-spacing:.08em;text-transform:uppercase}.cockpit-panel h1{max-width:760px;margin:0;font-size:34px;line-height:1.12;letter-spacing:-.03em}.cockpit-copy{max-width:700px;margin:14px 0 0;color:#64748b;line-height:1.7}.view-switch{display:inline-flex;gap:4px;margin-top:20px;padding:4px;border-radius:999px;background:rgba(31,41,51,.08)}.view-switch button{min-height:34px;border:0;border-radius:999px;padding:0 14px;background:transparent;color:#64748b;font-weight:800;cursor:pointer;transition:transform 160ms,background-color 160ms,color 160ms}.view-switch button:active{transform:scale(.96)}.view-switch button.active{background:#1f2933;color:rgba(255,255,255,.9);box-shadow:0 8px 16px rgba(31,41,51,.14)}.family-workspace-card{display:flex;justify-content:space-between;gap:16px;align-items:center;margin-top:16px;padding:14px 16px;border-radius:16px;background:rgba(255,255,255,.74);box-shadow:inset 0 0 0 1px rgba(100,116,139,.14),0 10px 24px rgba(47,125,92,.08)}.family-workspace-card strong,.family-workspace-card small,.workspace-label{display:block}.workspace-label{color:#2f7d5c;font-size:12px;font-weight:900;letter-spacing:.08em;text-transform:uppercase}.family-workspace-card strong{margin-top:4px;color:#1e293b;font-size:18px}.family-workspace-card small{margin-top:3px;color:#64748b}.workspace-actions{display:flex;gap:8px;align-items:center}.workspace-badge{border-radius:999px;padding:6px 10px;background:rgba(47,125,92,.12);color:#2f7d5c;font-size:12px;font-weight:900}.workspace-manage{min-height:30px;border:0;border-radius:999px;padding:0 12px;background:#1f2933;color:rgba(255,255,255,.9);font-size:12px;font-weight:900;cursor:pointer}.workspace-manage:active{transform:scale(.96)}.family-chip-rail{display:flex;gap:8px;max-width:100%;margin-top:10px;padding-bottom:4px;overflow-x:auto}.family-chip{min-height:34px;white-space:nowrap;border:0;border-radius:999px;padding:0 14px;background:rgba(255,255,255,.72);color:#64748b;font-weight:900;cursor:pointer;box-shadow:inset 0 0 0 1px rgba(100,116,139,.16);transition:transform 160ms,background-color 160ms,color 160ms,box-shadow 160ms}.family-chip:active{transform:scale(.96)}.family-chip.active{background:#1f2933;color:rgba(255,255,255,.9);box-shadow:0 8px 16px rgba(31,41,51,.14)}.cockpit-actions{display:flex;gap:12px;margin-top:24px;flex-wrap:wrap}.primary-action,.secondary-action,.ghost-action{min-height:40px;border:0;border-radius:12px;padding:0 18px;font-weight:800;cursor:pointer;transition:transform 160ms,box-shadow 160ms,background-color 160ms,color 160ms}.primary-action:active,.secondary-action:active,.ghost-action:active{transform:scale(.96)}.primary-action{background:#3b82f6;color:#fff;box-shadow:0 10px 24px rgba(59,130,246,.22)}.secondary-action,.ghost-action{background:rgba(255,255,255,.72);color:#1e293b;box-shadow:inset 0 0 0 1px rgba(100,116,139,.18)}.signal-card{align-self:stretch;padding:22px;border-radius:18px;background:rgba(31,41,51,.92);color:rgba(255,255,255,.86);box-shadow:0 18px 36px rgba(31,41,51,.18)}.signal-card span,.metric-card span{display:block;color:#94a3b8;font-size:12px;font-weight:800;text-transform:uppercase;letter-spacing:.08em}.signal-card strong{display:block;margin-top:18px;font-size:22px;line-height:1.25}.signal-card p{margin:14px 0 0;color:rgba(255,255,255,.64);line-height:1.7}.finance-overview-panel{margin-bottom:24px;padding:24px;background:#fff;border-radius:16px;box-shadow:0 1px 3px rgba(15,23,42,.1),0 10px 24px rgba(15,23,42,.04)}.finance-overview-head{display:flex;justify-content:space-between;gap:16px;align-items:flex-start;margin-bottom:16px}.finance-overview-head h2{margin:0;font-size:20px;letter-spacing:-.012em}.finance-overview-grid{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:16px;margin-top:16px}.overview-subcard{min-width:0;padding:16px;border-radius:16px;background:#fffaf0;box-shadow:inset 0 0 0 1px rgba(100,116,139,.10)}.section-head.compact{margin-bottom:8px}.kpi-panel{margin-bottom:24px;padding:20px;background:#fff;border-radius:16px;box-shadow:0 1px 3px rgba(15,23,42,.1),0 10px 24px rgba(15,23,42,.04)}.kpi-head{display:flex;justify-content:space-between;gap:16px;align-items:flex-start;margin-bottom:16px}.kpi-head h2{margin:0;font-size:20px;letter-spacing:-.012em}.metric-grid{display:grid;grid-template-columns:repeat(4,minmax(0,1fr));gap:16px}.chart-currency-badge{display:inline-flex;align-items:center;min-height:30px;border-radius:999px;padding:0 12px;background:rgba(31,41,51,.08);color:#1e293b;font-size:12px;font-weight:900}.metric-card,.surface-card{background:#fff;border-radius:16px;box-shadow:0 1px 3px rgba(15,23,42,.1),0 10px 24px rgba(15,23,42,.04)}.trend-card,.networth-card,.asset-composition-panel{margin-bottom:24px;padding:24px;background:#fff;border-radius:16px;box-shadow:0 1px 3px rgba(15,23,42,.1),0 12px 30px rgba(15,23,42,.04)}.trend-card .section-head,.networth-card .section-head{margin-bottom:0}.composition-grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(320px,1fr));gap:14px}.composition-card{padding:16px;border-radius:16px;background:#fffaf0;box-shadow:inset 0 0 0 1px rgba(100,116,139,.10)}.composition-head{display:flex;flex-wrap:wrap;gap:8px;align-items:center;margin-bottom:10px}.composition-head strong{margin-right:auto;color:#1e293b;font-size:16px;font-weight:900}.composition-head span{border-radius:999px;padding:5px 9px;background:#fff;color:#64748b;font-size:12px;font-weight:800}.liability-list{display:grid;gap:8px;margin-top:10px;padding:12px;border-radius:12px;background:rgba(239,68,68,.06)}.liability-list p{margin:0;color:#ef4444;font-size:12px;font-weight:900}.liability-row{display:flex;justify-content:space-between;gap:12px;padding:8px 0;border-top:1px solid rgba(239,68,68,.12)}.liability-row:first-of-type{border-top:0}.liability-row span{color:#1e293b;font-weight:800}.liability-row strong{font-family:'SF Mono','Fira Code',monospace;color:#ef4444}.investment-panel{margin-bottom:24px;padding:24px;background:#fff;border-radius:16px;box-shadow:0 1px 3px rgba(15,23,42,.1),0 12px 30px rgba(15,23,42,.04)}.investment-stack{display:grid;gap:14px}.investment-card{display:grid;grid-template-columns:220px minmax(0,1fr);gap:16px;padding:16px;border-radius:16px;background:#fffaf0;box-shadow:inset 0 0 0 1px rgba(100,116,139,.10)}.investment-card.gain{box-shadow:inset 0 0 0 1px rgba(47,125,92,.18)}.investment-card.loss{box-shadow:inset 0 0 0 1px rgba(239,68,68,.18)}.investment-card-head{display:flex;flex-direction:column;justify-content:center;gap:4px;padding:14px;border-radius:14px;background:#fff}.investment-card-head span{color:#7c3aed;font-size:12px;font-weight:900;letter-spacing:.08em}.investment-card-head strong{font-family:'SF Mono','Fira Code',monospace;font-size:28px;letter-spacing:-.02em}.investment-card.gain .investment-card-head strong{color:#2f7d5c}.investment-card.loss .investment-card-head strong{color:#ef4444}.investment-card-head small{color:#64748b;font-size:12px;font-weight:800}.investment-metrics{display:grid;grid-template-columns:repeat(auto-fit,minmax(130px,1fr));gap:10px}.investment-metrics div{padding:12px;border-radius:12px;background:#fff}.investment-metrics span{display:block;color:#64748b;font-size:11px;font-weight:800}.investment-metrics strong{display:block;margin-top:6px;font-family:'SF Mono','Fira Code',monospace;color:#1e293b;font-size:16px}.investment-bucket-list{grid-column:1/-1;display:grid;gap:8px}.investment-bucket-row{display:grid;grid-template-columns:minmax(0,1fr) minmax(120px,auto) minmax(120px,auto);gap:14px;align-items:center;padding:12px 14px;border-radius:12px;background:#fff}.investment-bucket-name{color:#1e293b;font-weight:800;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}.investment-bucket-metric{display:flex;flex-direction:column;gap:2px;align-items:flex-end}.investment-bucket-metric small{color:#94a3b8;font-size:11px;font-weight:800}.investment-bucket-metric strong,.investment-bucket-metric em{font-family:'SF Mono','Fira Code',monospace;font-style:normal;font-weight:900}.investment-bucket-metric em.gain{color:#2f7d5c}.investment-bucket-metric em.loss{color:#ef4444}.metric-card{padding:20px;transition:transform 180ms,box-shadow 180ms}.metric-card strong{display:block;margin-top:12px;font-family:'SF Mono','Fira Code',monospace;font-size:24px;letter-spacing:-.022em;word-break:break-word}.metric-card.positive strong{color:#2f7d5c}.metric-card.negative strong{color:#ef4444}.metric-card p{margin:8px 0 0;color:#64748b;font-size:13px;line-height:1.6}.stats-panel{margin-bottom:24px;padding:24px;background:#fff;border-radius:16px;box-shadow:0 1px 3px rgba(15,23,42,.1),0 12px 30px rgba(15,23,42,.04)}.report-context-strip{display:flex;justify-content:space-between;gap:16px;align-items:center;margin:4px 0 16px;padding:16px;border-radius:16px;background:linear-gradient(135deg,rgba(220,233,223,.72),rgba(255,255,255,.96));box-shadow:inset 0 0 0 1px rgba(47,125,92,.14),0 10px 24px rgba(47,125,92,.06)}.report-context-main{min-width:0}.report-context-kicker{margin:0 0 6px;color:#2f7d5c;font-size:11px;font-weight:900;letter-spacing:.08em;text-transform:uppercase}.report-context-main strong{display:block;color:#1e293b;font-size:16px;letter-spacing:-.012em}.report-context-main span{display:block;margin-top:5px;color:#64748b;line-height:1.6}.report-context-modules{display:flex;flex-wrap:wrap;justify-content:flex-end;gap:8px;max-width:520px}.report-context-modules small{display:inline-flex;align-items:center;min-height:28px;color:#2f7d5c;font-size:12px;font-weight:900}.report-context-modules span,.report-range-badge,.snapshot-badge{display:inline-flex;align-items:center;min-height:28px;border-radius:999px;padding:0 10px;font-size:12px;font-weight:900;white-space:nowrap}.report-context-modules span{background:#fff;color:#1e293b;box-shadow:inset 0 0 0 1px rgba(100,116,139,.12)}.report-range-badge{background:rgba(47,125,92,.12);color:#2f7d5c}.snapshot-badge{background:rgba(31,41,51,.08);color:#1f2933}.snapshot-badge.latest{background:rgba(59,130,246,.12);color:#2563eb}.report-range-badge.inline,.snapshot-badge.inline{margin-top:8px}.section-badge-row{display:flex;flex-wrap:wrap;justify-content:flex-end;gap:8px;align-items:center}.section-head{display:flex;justify-content:space-between;gap:16px;align-items:flex-start;margin-bottom:18px}.section-head h2{margin:0;font-size:20px;letter-spacing:-.012em}.stats-hint{margin:0 0 16px;color:#64748b;line-height:1.7}.currency-stats-stack,.family-assets-stack,.event-list,.currency-stack{display:grid;gap:12px}.currency-stats-card,.member-asset-card{padding:16px;border-radius:16px;background:#fffaf0}.currency-stats-card h3{margin:0 0 12px}.stats-grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(150px,1fr));gap:12px}.stat-tile{padding:16px;border-radius:14px;background:#f8faf7}.stat-tile-link{appearance:none;border:0;width:100%;text-align:left;cursor:pointer;color:inherit;font:inherit;transition:transform 160ms,box-shadow 160ms,background-color 160ms}.stat-tile-link:hover{background:#eef7f0;box-shadow:0 8px 18px rgba(47,125,92,.12)}.stat-tile-link:active{transform:scale(.98)}.stat-tile span{color:#64748b;font-size:12px;font-weight:800}.stat-tile strong{display:block;margin-top:8px;font-family:'SF Mono','Fira Code',monospace;font-size:24px}.stat-tile small{display:block;margin-top:6px;color:#94a3b8;font-size:12px}.stat-tile.income strong{color:#ef4444}.stat-tile.expense strong{color:#f97316}.stats-range-switch{display:inline-flex;gap:4px;max-width:100%;padding:4px;border-radius:999px;background:rgba(31,41,51,.08);overflow-x:auto}.stats-range-switch button{min-height:30px;white-space:nowrap;border:0;border-radius:999px;padding:0 12px;background:transparent;color:#64748b;font-size:12px;font-weight:900;cursor:pointer}.stats-range-switch button.active{background:#1f2933;color:rgba(255,255,255,.9)}.stats-range-switch button:active{transform:scale(.96)}.category-mode-head{display:flex;justify-content:space-between;gap:12px;align-items:center;margin-top:18px}.category-mode-head h3{margin:0}.category-type-switch{display:inline-flex;gap:4px;padding:4px;border-radius:999px;background:rgba(31,41,51,.08)}.category-type-switch button{min-height:30px;border:0;border-radius:999px;padding:0 12px;background:transparent;color:#64748b;font-size:12px;font-weight:900;cursor:pointer}.category-type-switch button.active{background:#1f2933;color:rgba(255,255,255,.9)}.category-type-switch button:active{transform:scale(.96)}.category-rank{margin-top:14px}.category-rank h3{margin:0 0 12px}.category-currency-card{margin-top:12px;border-radius:14px;background:#fffaf0;box-shadow:inset 0 0 0 1px rgba(100,116,139,.10);overflow:hidden}.category-currency-head{display:flex;justify-content:space-between;align-items:center;gap:12px;width:100%;border:0;background:transparent;padding:14px 16px;text-align:left;cursor:pointer}.category-currency-head strong{display:block;color:#1e293b;font-weight:900}.category-currency-head small{display:block;margin-top:4px;color:#64748b;font-size:12px;font-weight:800}.category-currency-head>span{color:#2f7d5c;font-size:12px;font-weight:900}.category-currency-body{padding:0 14px 14px}.category-row-toggle{width:100%;min-height:34px;margin-top:8px;border:0;border-radius:999px;background:#fff;color:#2f7d5c;font-size:12px;font-weight:900;cursor:pointer;box-shadow:inset 0 0 0 1px rgba(47,125,92,.18)}.category-row-toggle:active{transform:scale(.98)}.rank-row{display:grid;grid-template-columns:12px minmax(0,1fr) auto;gap:10px;align-items:center;padding:10px 12px;border-radius:12px}.rank-row-link{appearance:none;border:0;width:100%;background:transparent;color:inherit;font:inherit;text-align:left;cursor:pointer;transition:transform 160ms,box-shadow 160ms,background-color 160ms}.rank-row-link:hover{background:#f8faf7;box-shadow:0 6px 14px rgba(15,23,42,.06)}.rank-row-link:active{transform:scale(.99)}.rank-row span{width:10px;height:10px;border-radius:999px}.rank-row em{font-style:normal;font-family:'SF Mono','Fira Code',monospace}.workspace-grid{display:grid;grid-template-columns:minmax(0,.9fr) minmax(0,1.1fr);gap:24px}.surface-card{padding:24px}.currency-row,.event-row{display:flex;align-items:center;gap:14px;padding:14px;border-radius:12px;background:#f8faf7;width:100%;box-sizing:border-box;min-width:0}.event-row-link{appearance:none;border:0;color:inherit;font:inherit;text-align:left;cursor:pointer;transition:transform 160ms,box-shadow 160ms,background-color 160ms}.event-row-link:hover{background:#eef7f0;box-shadow:0 8px 18px rgba(47,125,92,.12)}.event-row-link:active{transform:scale(.99)}.currency-row{justify-content:space-between}.currency-row>div{min-width:0}.currency-row>strong{min-width:0;text-align:right;overflow:hidden;text-overflow:ellipsis}.currency-code{display:block;color:#1e293b;font-weight:800}.currency-row small,.event-main span,.member-head span,.asset-total-pill small,.family-bucket-row span{color:#64748b;font-size:12px}.currency-row strong,.event-amount,.asset-total-pill strong,.family-bucket-row em{font-family:'SF Mono','Fira Code',monospace;font-variant-numeric:tabular-nums}.event-dot{width:10px;height:10px;border-radius:999px;background:#2f7d5c;box-shadow:0 0 0 5px rgba(47,125,92,.12)}.event-main{flex:1;min-width:0}.event-main strong,.event-main span{display:block;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}.event-amount{color:#1e293b;font-weight:700}.member-head-button{display:grid;grid-template-columns:minmax(0,1fr) auto 28px;gap:12px;align-items:center;width:100%;border:0;background:transparent;color:#1e293b;text-align:left;cursor:pointer}.member-head-button:active{transform:scale(.99)}.member-head-main strong,.member-head-main span{display:block}.member-head-main span{color:#64748b;font-size:12px}.member-summary-meta{display:flex;gap:6px;flex-wrap:wrap;justify-content:flex-end}.member-summary-meta span{border-radius:999px;padding:5px 8px;background:rgba(31,41,51,.06);color:#64748b;font-size:11px;font-weight:900}.member-collapse-icon{display:grid;place-items:center;width:28px;height:28px;border-radius:999px;background:#fff;color:#2f7d5c;font-weight:900}.asset-total-grid{display:grid;gap:8px;margin-top:12px;margin-bottom:12px}.asset-total-grid.compact{grid-template-columns:1fr}.asset-total-pill{padding:12px;border-radius:12px;background:#fff}.asset-total-pill span,.asset-total-pill strong,.asset-total-pill small{display:block}.asset-total-pill span{color:#2f7d5c;font-size:12px;font-weight:900}.asset-total-pill strong{margin-top:4px;color:#1e293b}.asset-total-pill small{margin-top:4px}.member-asset-body{margin-top:10px}.member-empty{margin:8px 0 0;color:#64748b;font-size:13px}.family-account-row{padding:12px;border-radius:12px;background:#fff;margin-top:8px}.family-account-row>strong{display:block;margin-bottom:8px}.family-bucket-row{display:flex;justify-content:space-between;gap:10px;padding:8px 0;border-top:1px solid rgba(100,116,139,.1)}.family-bucket-row em{font-style:normal}.empty-state{display:grid;place-items:center;gap:12px;padding:24px;text-align:center;color:#64748b}.empty-state img{border-radius:20px}@media(hover:hover){.metric-card:hover,.surface-card:hover{transform:translateY(-2px);box-shadow:0 3px 8px rgba(15,23,42,.12),0 16px 34px rgba(15,23,42,.06)}}@media(max-width:1024px){.cockpit-panel,.brand-cluster,.workspace-grid,.finance-overview-grid{grid-template-columns:1fr}.signal-card{max-width:none}.metric-grid{grid-template-columns:repeat(2,minmax(0,1fr))}}@media(max-width:720px){.finance-overview-head,.section-head,.category-mode-head,.report-context-strip{flex-direction:column;align-items:stretch}.report-context-modules{justify-content:flex-start;max-width:none}.finance-overview-head .stats-range-switch,.stats-range-switch,.category-type-switch{width:100%;overflow-x:auto}.metric-grid,.finance-overview-grid,.composition-grid,.investment-card,.workspace-grid{grid-template-columns:1fr}.investment-bucket-row,.summary-row{grid-template-columns:1fr}.investment-bucket-metric,.row-amount{align-items:flex-start}}@media(max-width:640px){.cockpit-panel,.surface-card,.finance-overview-panel,.stats-panel,.trend-card,.investment-panel,.asset-composition-panel{padding:18px;border-radius:18px}.brand-cluster{gap:16px}.brand-orbit{min-height:130px}.brand-orbit img:first-child{width:96px;height:96px}.brand-orbit img:last-child{width:64px;height:64px;left:92px;right:auto}.cockpit-panel h1{font-size:26px}.metric-grid{grid-template-columns:1fr}.event-row,.currency-row{align-items:flex-start}.event-row{flex-wrap:wrap}.family-workspace-card{align-items:flex-start;flex-direction:column}.stats-head{align-items:flex-start;flex-direction:column}.stats-range-switch{width:100%}.member-head-button{grid-template-columns:1fr}.member-summary-meta{justify-content:flex-start}.liability-row{flex-direction:column;align-items:flex-start}.composition-head{align-items:flex-start}.composition-head strong{width:100%}}@keyframes revealUp{from{opacity:0;transform:translateY(12px);filter:blur(4px)}to{opacity:1;transform:translateY(0);filter:blur(0)}}
</style>
