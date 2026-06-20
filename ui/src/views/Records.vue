<template>
  <main class="records-page">
    <section class="page-hero reveal-block">
      <div>
        <p class="eyebrow">{{ t('records.hero.eyebrow') }}</p>
        <h1>{{ t('records.hero.title') }}</h1>
        <p>{{ t('records.hero.subtitle') }}</p>
      </div>
    </section>

    <section class="filters-card reveal-block delay-1">
      <div class="filter-card-head">
        <div>
          <p class="eyebrow">{{ t('records.filters.title') }}</p>
          <h2>{{ t('records.filters.subtitle') }}</h2>
        </div>
        <div class="filter-card-actions">
          <button class="ghost-action" @click="resetAndLoad">{{ t('common.actions.refresh') }}</button>
          <button class="export-action" type="button" @click="handleExportRecords">{{ t('records.actions.exportRecords') }}</button>
        </div>
      </div>

      <div v-if="filters.categoryGroupId" class="filter-banner">
        <span class="banner-label">{{ t('records.filters.categoryGroup') }}</span>
        <strong>{{ filters.categoryGroupName || `#${filters.categoryGroupId}` }}</strong>
        <button class="banner-clear" type="button" @click="clearCategoryGroup">×</button>
      </div>

      <div class="filter-row">
        <span class="filter-label">{{ t('records.filters.timeRange') }}</span>
        <div class="filter-chip-rail" role="group" :aria-label="t('records.filters.timeRange')">
          <button v-for="preset in rangePresets" :key="preset.value" type="button" class="filter-chip" :class="{ active: filters.rangeMode === 'preset' && filters.rangePreset === preset.value }" @click="selectPresetRange(preset.value)">{{ t(preset.labelKey) }}</button>
          <button type="button" class="filter-chip" :class="{ active: filters.rangeMode === 'custom' }" @click="changeRangeMode('custom')">{{ t('records.rangeModes.custom') }}</button>
        </div>
        <el-date-picker v-if="filters.rangeMode === 'custom'" v-model="filters.customDateRange" type="daterange" value-format="YYYY-MM-DD" :start-placeholder="t('records.filters.startDate')" :end-placeholder="t('records.filters.endDate')" class="date-range-control" @change="resetAndLoad" />
      </div>

      <div class="filter-row">
        <span class="filter-label">{{ t('records.filters.eventType') }}</span>
        <div class="filter-chip-rail" role="group" :aria-label="t('records.filters.eventType')">
          <button type="button" class="filter-chip" :class="{ active: !filters.eventType }" @click="selectEventType('')">{{ t('records.filters.all') }}</button>
          <button v-for="item in eventFilterTypes" :key="item" type="button" class="filter-chip" :class="{ active: filters.eventType === item }" @click="selectEventType(item)">{{ eventTypeLabel(item) }}</button>
        </div>
      </div>

      <div class="filter-search-row">
        <el-input v-model="filters.keyword" :placeholder="t('records.filters.keyword')" class="keyword-control" @keyup.enter="resetAndLoad" />
        <button class="ghost-action" type="button" @click="advancedFiltersOpen = !advancedFiltersOpen">{{ advancedFiltersOpen ? t('records.filters.collapseAdvanced') : t('records.filters.advanced') }}</button>
      </div>

      <div v-if="advancedFiltersOpen" class="advanced-filter-panel">
        <el-select v-model="filters.currency" clearable :placeholder="t('records.filters.currency')" class="filter-control" @change="resetAndLoad">
          <el-option v-for="item in currencyOptions" :key="item.code" :label="getCurrencyLabel(item.code, config.locale)" :value="item.code" />
        </el-select>
        <el-select v-model="filters.categoryId" clearable :placeholder="t('records.filters.category')" class="filter-control" @change="resetAndLoad">
          <el-option v-for="item in categories" :key="item.id" :label="item.name" :value="item.id" />
        </el-select>
        <el-select v-model="filters.includeInStatistics" clearable :placeholder="t('records.filters.includeStats')" class="filter-control" @change="resetAndLoad">
          <el-option :label="t('record.preview.included')" :value="true" />
          <el-option :label="t('record.preview.excluded')" :value="false" />
        </el-select>
      </div>
    </section>

    <section class="records-layout reveal-block delay-2">
      <div v-loading="loading" class="records-list">
        <article v-for="event in events" :key="event.groupKey || event.id" class="record-group">
          <div class="record-row group-main-row" :class="[{ active: selectedEvent?.id === event.id }, event.eventType]" @click="selectEvent(event)">
            <button v-if="event.groupSize > 1" class="group-expand" type="button" @click.stop="toggleGroup(event)">{{ isGroupExpanded(event) ? '−' : '+' }}</button>
            <div v-else class="group-expand-placeholder"></div>
            <div class="type-dot"></div>
            <div class="record-main">
              <strong>{{ event.description || eventTypeLabel(event.eventType) }}</strong>
              <span>{{ event.eventTime }} · {{ eventTypeLabel(event.eventType) }}<template v-if="event.groupSize > 1"> · {{ t('records.group.count', { count: event.groupSize }) }}</template></span>
            </div>
            <div class="record-amount">{{ event.currency }} {{ formatAmount(event.displayAmount || event.amount) }}</div>
          </div>
          <div v-if="event.groupSize > 1 && isGroupExpanded(event)" class="group-children">
            <button v-for="child in event.children || []" :key="child.id" type="button" class="group-child-row" :class="[child.eventType, { active: selectedDetail?.id === child.id }]" @click="selectEvent(child)">
              <span class="child-marker">{{ groupEmoji(child.eventType) }}</span>
              <span>{{ child.description || eventTypeLabel(child.eventType) }}</span>
              <em>{{ child.currency }} {{ formatAmount(child.amount) }}</em>
            </button>
          </div>
        </article>
        <div v-if="!loading && !events.length" class="empty-state">
          <p>{{ t('records.detail.empty') }}</p>
        </div>
        <div class="pagination-wrap">
          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :total="pagination.total"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handlePageSizeChange"
            @current-change="handlePageChange"
          />
        </div>
      </div>

      <aside class="detail-panel">
        <template v-if="selectedDetail">
          <div class="section-head">
            <div>
              <p class="eyebrow">{{ t('records.detail.title') }}</p>
              <h2>{{ selectedDetail.description || eventTypeLabel(selectedDetail.eventType) }}</h2>
            </div>
          </div>

          <div class="event-card" :class="selectedDetail.eventType">
            <span>{{ t('records.detail.event') }}</span>
            <strong>{{ selectedDetail.currency }} {{ formatAmount(selectedDetail.amount) }}</strong>
            <p>{{ eventTypeLabel(selectedDetail.eventType) }} · {{ selectedDetail.includeInStatistics ? t('record.preview.included') : t('record.preview.excluded') }}</p>
          </div>

          <div class="detail-grid">
            <span>ID</span><strong>#{{ selectedDetail.id }}</strong>
            <span>{{ t('record.fields.eventTime') }}</span><strong>{{ selectedDetail.eventTime }}</strong>
            <span>{{ t('record.fields.category') }}</span><strong>{{ categoryName(selectedDetail.categoryId) }}</strong>
            <span>{{ t('record.fields.remark') }}</span><strong>{{ selectedDetail.remark || '-' }}</strong>
            <template v-if="selectedDetail.relatedFinancialEventId">
              <span>{{ t('record.fields.relatedEvent') }}</span>
              <button class="related-link" type="button" @click="focusEventById(selectedDetail.relatedFinancialEventId)">#{{ selectedDetail.relatedFinancialEventId }} →</button>
            </template>
          </div>

          <div v-if="groupEvents.length > 1" class="group-block">
            <h3>{{ t('records.detail.group') }} <span class="group-count">{{ t('records.detail.groupCount', { count: groupEvents.length }) }}</span></h3>
            <div class="group-list">
              <button v-for="g in groupEvents" :key="g.id" type="button" class="group-row" :class="[g.eventType, { current: g.isCurrent }]" @click="focusEventById(g.id)">
                <span class="group-marker">{{ groupEmoji(g.eventType) }}</span>
                <div class="group-main">
                  <strong>{{ g.description || eventTypeLabel(g.eventType) }}</strong>
                  <small>{{ formatDate(g.eventTime) }} · {{ eventTypeLabel(g.eventType) }}</small>
                </div>
                <div class="group-amount">
                  <strong>{{ g.currency }} {{ formatAmount(g.amount) }}</strong>
                  <small v-if="hasOutstandingMeta(g)" class="group-meta" :class="{ settled: Number(g.outstandingAmount) === 0 }">
                    {{ Number(g.outstandingAmount) === 0 ? t('records.detail.groupSettled') : t('records.detail.groupRemaining', { amount: formatAmount(g.outstandingAmount) }) }}
                  </small>
                </div>
              </button>
            </div>
          </div>

          <div class="entries-block">
            <h3>{{ t('records.detail.entries') }}</h3>
            <div v-for="entry in selectedDetail.ledgerEntries || []" :key="entry.id" class="entry-row">
              <div>
                <strong>{{ bucketDisplay(entry.bucketId) }}</strong>
                <span>{{ entryRoleLabel(entry.entryRole) }} · {{ entry.balanceAfter }}</span>
              </div>
              <em :class="{ positive: Number(entry.amount) > 0, negative: Number(entry.amount) < 0 }">{{ entry.currency }} {{ formatAmount(entry.amount) }}</em>
            </div>
          </div>

          <div class="detail-actions" v-if="canMutate(selectedDetail)">
            <button class="primary-action" @click="openEdit(selectedDetail)">{{ t('records.actions.edit') }}</button>
            <button class="danger-action" @click="handleDelete(selectedDetail)">{{ t('records.actions.delete') }}</button>
          </div>
          <p v-else class="readonly-note">{{ t('records.messages.unsupported') }}</p>
        </template>
        <div v-else class="empty-state compact">
          <p>{{ t('records.detail.empty') }}</p>
        </div>
      </aside>
    </section>

    <el-dialog v-model="dialogVisible" :title="t('records.dialog.editTitle')" width="560px" class="marmot-dialog">
      <el-form label-position="top">
        <el-form-item :label="t('record.fields.amount')"><el-input v-model="form.amount" /></el-form-item>
        <el-form-item :label="t('record.fields.currency')"><el-select v-model="form.currency" class="full-width"><el-option v-for="item in currencyOptions" :key="item.code" :label="getCurrencyLabel(item.code, config.locale)" :value="item.code" /></el-select></el-form-item>
        <el-form-item v-if="form.scenario !== 'transfer'" :label="t('record.fields.bucket')"><el-select v-model="form.bucketId" class="full-width" filterable><el-option v-for="bucket in buckets" :key="bucket.id" :label="bucketLabel(bucket)" :value="bucket.id" /></el-select></el-form-item>
        <template v-else>
          <el-form-item :label="t('record.fields.fromBucket')"><el-select v-model="form.fromBucketId" class="full-width" filterable><el-option v-for="bucket in buckets" :key="bucket.id" :label="bucketLabel(bucket)" :value="bucket.id" /></el-select></el-form-item>
          <el-form-item :label="t('record.fields.toBucket')"><el-select v-model="form.toBucketId" class="full-width" filterable><el-option v-for="bucket in buckets" :key="bucket.id" :label="bucketLabel(bucket)" :value="bucket.id" /></el-select></el-form-item>
        </template>
        <el-form-item v-if="form.scenario !== 'transfer'" :label="t('record.fields.category')"><el-select v-model="form.categoryId" clearable class="full-width"><el-option v-for="item in filteredCategories" :key="item.id" :label="item.name" :value="item.id" /></el-select></el-form-item>
        <el-form-item :label="t('record.fields.description')"><el-input v-model="form.description" /></el-form-item>
        <el-form-item :label="t('record.fields.eventTime')"><el-input v-model="form.eventTime" /></el-form-item>
        <el-form-item :label="t('record.fields.remark')"><el-input v-model="form.remark" type="textarea" /></el-form-item>
      </el-form>
      <template #footer>
        <button class="ghost-action" @click="dialogVisible = false">{{ t('common.actions.cancel') }}</button>
        <button class="primary-action" @click="submitEdit">{{ t('common.actions.save') }}</button>
      </template>
    </el-dialog>
  </main>
</template>

<script setup>
import { computed, onActivated, onMounted, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useConfigStore } from '@/stores/config'
import { listFinancialEvents, getFinancialEvent, getEventGroup } from '@/api/financialEvent/financialEvent'
import { updateRecord, deleteRecord } from '@/api/record/record'
import { listBuckets } from '@/api/bucket/bucket'
import { listCategories } from '@/api/category/category'
import { currencyOptions, getCurrencyLabel } from '@/utils/currency'
import { getBucketEmoji } from '@/utils/bucketEmoji'
import { exportRecords } from '@/api/export/export'

const { t, te } = useI18n()
const route = useRoute()
const router = useRouter()
const config = useConfigStore()
const events = ref([])
const selectedEvent = ref(null)
const selectedDetail = ref(null)
const groupEvents = ref([])
const expandedGroupKeys = ref([])
const buckets = ref([])
const categories = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const advancedFiltersOpen = ref(false)
const editingId = ref(0)
const editableTypes = ['income', 'expense', 'transfer', 'refund']
const eventFilterTypes = ['income', 'expense', 'refund', 'transfer', 'family_transfer', 'exchange', 'receivable_create', 'receivable_collect', 'deposit_create', 'deposit_refund', 'loan_out', 'loan_collect', 'investment_buy', 'investment_sell', 'investment_income', 'investment_revalue', 'balance_adjustment']
const filters = reactive({ eventType: '', categoryId: '', categoryGroupId: '', categoryGroupName: '', currency: '', keyword: '', includeInStatistics: '', rangeMode: 'preset', rangePreset: 'thisMonth', customDateRange: [] })
const pagination = reactive({ page: 1, pageSize: 10, total: 0 })
const form = reactive({ scenario: 'expense', bucketId: '', fromBucketId: '', toBucketId: '', categoryId: '', amount: '', currency: 'CNY', description: '', eventTime: '', relatedFinancialEventId: '', remark: '' })

const filteredCategories = computed(() => categories.value.filter((item) => form.scenario === 'income' ? item.type === 'income' : item.type === 'expense'))
const rangePresets = [
  { value: 'recentMonth', labelKey: 'statistics.rangePresets.recentMonth' },
  { value: 'thisMonth', labelKey: 'statistics.rangePresets.thisMonth' },
  { value: 'recentYear', labelKey: 'statistics.rangePresets.recentYear' },
  { value: 'thisYear', labelKey: 'statistics.rangePresets.thisYear' },
]

function enumLabel(prefix, value) {
  const key = `${prefix}.${value}`
  return value && te(key) ? t(key) : value
}
function eventTypeLabel(type) { return enumLabel('record.scenarios', type) }
function entryRoleLabel(role) { return enumLabel('record.entryRoles', role) }
function formatAmount(value) { return new Intl.NumberFormat(config.locale, { minimumFractionDigits: 2, maximumFractionDigits: 2 }).format(Number(value || 0)) }
function bucketName(id) { return buckets.value.find((item) => item.id === Number(id))?.name || `#${id}` }
function bucketDisplay(id) {
  const bucket = buckets.value.find((item) => item.id === Number(id))
  if (!bucket) return `#${id}`
  return `${getBucketEmoji(bucket.bucketType)} ${bucket.name}`
}
function bucketLabel(bucket) { return `${getBucketEmoji(bucket.bucketType)} ${bucket.name} · ${bucket.currency} ${formatAmount(bucket.balance)}` }
function categoryName(id) { return categories.value.find((item) => item.id === Number(id))?.name || '-' }
function canMutate(event) { return editableTypes.includes(event.eventType) }
function padDatePart(value) { return String(value).padStart(2, '0') }
function formatDateTime(date, endOfDay = false) {
  return `${date.getFullYear()}-${padDatePart(date.getMonth() + 1)}-${padDatePart(date.getDate())} ${endOfDay ? '23:59:59' : '00:00:00'}`
}
function lastDayOfMonth(year, monthIndex) { return new Date(year, monthIndex + 1, 0).getDate() }
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
function presetRangeParams() {
  const now = new Date()
  let start
  let end
  switch (filters.rangePreset) {
    case 'recentYear': start = subtractYearsClamped(now, 1); end = now; break
    case 'thisYear': start = new Date(now.getFullYear(), 0, 1); end = new Date(now.getFullYear(), 11, 31); break
    case 'recentMonth': start = subtractMonthsClamped(now, 1); end = now; break
    case 'thisMonth':
    default: start = new Date(now.getFullYear(), now.getMonth(), 1); end = new Date(now.getFullYear(), now.getMonth(), lastDayOfMonth(now.getFullYear(), now.getMonth())); break
  }
  return { startTime: formatDateTime(start), endTime: formatDateTime(end, true) }
}
function recordTimeParams() {
  if (filters.rangeMode === 'custom') {
    if (!Array.isArray(filters.customDateRange) || filters.customDateRange.length !== 2) return {}
    return { startTime: `${filters.customDateRange[0]} 00:00:00`, endTime: `${filters.customDateRange[1]} 23:59:59` }
  }
  return presetRangeParams()
}
function changeRangeMode(mode) {
  if (filters.rangeMode === mode) return
  filters.rangeMode = mode
  resetAndLoad()
}
function selectPresetRange(preset) {
  filters.rangeMode = 'preset'
  filters.rangePreset = preset
  resetAndLoad()
}
function selectEventType(type) {
  filters.eventType = type
  resetAndLoad()
}

function clearCategoryGroup() {
  filters.categoryGroupId = ''
  filters.categoryGroupName = ''
  resetAndLoad()
}

function applyRouteQuery(query) {
  if (!query || typeof query !== 'object') return false
  let touched = false
  if (query.rangeMode === 'preset' || query.rangeMode === 'custom') {
    filters.rangeMode = query.rangeMode
    touched = true
  }
  if (query.rangePreset) {
    filters.rangePreset = String(query.rangePreset)
    filters.rangeMode = 'preset'
    touched = true
  }
  if (Array.isArray(query.customDateRange) && query.customDateRange.length === 2) {
    filters.customDateRange = [String(query.customDateRange[0]), String(query.customDateRange[1])]
    filters.rangeMode = 'custom'
    touched = true
  }
  if (typeof query.eventType === 'string') {
    filters.eventType = query.eventType
    touched = true
  }
  if (typeof query.currency === 'string') {
    filters.currency = query.currency
    touched = true
  }
  if (query.categoryId) {
    filters.categoryId = String(query.categoryId)
    touched = true
  }
  if (query.categoryGroupId) {
    filters.categoryGroupId = String(query.categoryGroupId)
    filters.categoryGroupName = typeof query.categoryGroupName === 'string' ? query.categoryGroupName : ''
    touched = true
  }
  if (typeof query.keyword === 'string') {
    filters.keyword = query.keyword
    touched = true
  }
  if (query.includeInStatistics === 'true') { filters.includeInStatistics = true; touched = true }
  else if (query.includeInStatistics === 'false') { filters.includeInStatistics = false; touched = true }
  if (touched) {
    pagination.page = 1
    if (filters.categoryId || filters.categoryGroupId || filters.includeInStatistics !== '' || filters.currency) {
      advancedFiltersOpen.value = true
    }
  }
  return touched
}

async function loadBase() {
  const [bucketRes, incomeCats, expenseCats] = await Promise.all([listBuckets(), listCategories({ type: 'income', isActive: true }), listCategories({ type: 'expense', isActive: true })])
  if (bucketRes.success) buckets.value = bucketRes.data || []
  categories.value = [...(incomeCats.success ? incomeCats.data || [] : []), ...(expenseCats.success ? expenseCats.data || [] : [])]
}

async function loadEvents() {
  loading.value = true
  try {
    const params = { page: pagination.page, pageSize: pagination.pageSize, groupMode: true }
    if (filters.eventType) params.eventType = filters.eventType
    if (filters.categoryId) params.categoryId = filters.categoryId
    if (filters.categoryGroupId) params.categoryGroupId = filters.categoryGroupId
    if (filters.currency) params.currency = filters.currency
    Object.assign(params, recordTimeParams())
    if (filters.keyword) params.keyword = filters.keyword
    if (filters.includeInStatistics !== '') params.includeInStatistics = filters.includeInStatistics
    const res = await listFinancialEvents(params)
    if (res.success) {
      events.value = res.data?.list || []
      pagination.total = res.data?.total || 0
    } else ElMessage.error(res.error || t('records.messages.loadFailed'))
  } finally { loading.value = false }
}

function resetAndLoad() {
  pagination.page = 1
  loadEvents()
}

function exportParams() {
  const params = {}
  if (filters.eventType) params.eventType = filters.eventType
  if (filters.categoryId) params.categoryId = filters.categoryId
  if (filters.categoryGroupId) params.categoryGroupId = filters.categoryGroupId
  if (filters.currency) params.currency = filters.currency
  if (filters.keyword) params.keyword = filters.keyword
  if (filters.includeInStatistics !== '') params.includeInStatistics = filters.includeInStatistics
  Object.assign(params, recordTimeParams())
  return params
}

async function handleExportRecords() {
  try { await exportRecords(exportParams()) }
  catch (err) { console.warn(err); ElMessage.error(t('records.messages.exportFailed')) }
}
function handlePageChange(page) {
  pagination.page = page
  loadEvents()
}

function handlePageSizeChange(pageSize) {
  pagination.pageSize = pageSize
  pagination.page = 1
  loadEvents()
}

async function selectEvent(event) {
  selectedEvent.value = event
  const res = await getFinancialEvent(event.id)
  if (res.success) {
    selectedDetail.value = res.data
    await loadGroupEvents()
  }
}

async function loadGroupEvents() {
  const groupId = selectedDetail.value?.eventGroupId
  if (!groupId) {
    groupEvents.value = []
    return
  }
  const res = await getEventGroup(groupId, selectedDetail.value.id)
  if (res.success) groupEvents.value = res.data || []
  else groupEvents.value = []
}

function groupEmoji(eventType) {
  if (eventType === 'income') return '↙'
  if (eventType === 'expense') return '↗'
  if (eventType === 'refund') return '↩'
  if (eventType === 'transfer') return '⇄'
  if (eventType === 'family_transfer') return '⇄'
  if (eventType === 'exchange') return 'FX'
  if (eventType === 'receivable_create') return '📥'
  if (eventType === 'receivable_collect') return '✅'
  if (eventType === 'deposit_create') return '🔒'
  if (eventType === 'deposit_refund') return '✅'
  if (eventType === 'loan_out') return '🤝'
  if (eventType === 'loan_collect') return '✅'
  if (eventType.startsWith('investment')) return '📈'
  return '•'
}
function formatDate(eventTime) { return (eventTime || '').split(' ')[0] }
function hasOutstandingMeta(g) {
  return ['receivable_create', 'deposit_create', 'loan_out'].includes(g.eventType)
}

function groupKeyOf(event) {
  return Number(event.groupKey || event.id)
}
function isGroupExpanded(event) {
  return expandedGroupKeys.value.includes(groupKeyOf(event))
}
function toggleGroup(event) {
  const key = groupKeyOf(event)
  if (expandedGroupKeys.value.includes(key)) {
    expandedGroupKeys.value = expandedGroupKeys.value.filter((item) => item !== key)
  } else {
    expandedGroupKeys.value = [...expandedGroupKeys.value, key]
  }
}

function openEdit(event) {
  if (!canMutate(event)) return ElMessage.warning(t('records.messages.unsupported'))
  editingId.value = event.id
  const entries = event.ledgerEntries || []
  Object.assign(form, { scenario: event.eventType, bucketId: '', fromBucketId: '', toBucketId: '', categoryId: event.categoryId || '', amount: String(event.amount || ''), currency: event.currency, description: event.description, eventTime: event.eventTime, relatedFinancialEventId: event.relatedFinancialEventId || '', remark: event.remark || '' })
  if (event.eventType === 'transfer') {
    form.fromBucketId = entries.find((item) => item.entryRole === 'transfer_out')?.bucketId || ''
    form.toBucketId = entries.find((item) => item.entryRole === 'transfer_in')?.bucketId || ''
  } else {
    form.bucketId = entries[0]?.bucketId || ''
  }
  dialogVisible.value = true
}

async function submitEdit() {
  const payload = { ...form, bucketId: Number(form.bucketId || 0), fromBucketId: Number(form.fromBucketId || 0), toBucketId: Number(form.toBucketId || 0), categoryId: Number(form.categoryId || 0), relatedFinancialEventId: Number(form.relatedFinancialEventId || 0), amount: String(form.amount) }
  const res = await updateRecord(editingId.value, payload)
  if (res.success) { ElMessage.success(t('records.messages.updated')); dialogVisible.value = false; await loadEvents(); await selectEvent(res.data.financialEvent) }
  else ElMessage.error(res.error || t('records.messages.updateFailed'))
}

async function handleDelete(event) {
  try {
    await ElMessageBox.confirm(t('records.delete.confirm', { name: event.description || eventTypeLabel(event.eventType) }), t('records.delete.title'), { confirmButtonText: t('common.actions.delete'), cancelButtonText: t('common.actions.cancel'), type: 'warning', customClass: 'calm-marmot-message-box calm-marmot-delete-box', confirmButtonClass: 'calm-marmot-danger-confirm', cancelButtonClass: 'calm-marmot-soft-cancel' })
    const res = await deleteRecord(event.id)
    if (res.success) { ElMessage.success(t('records.messages.deleted')); selectedEvent.value = null; selectedDetail.value = null; await loadEvents() }
    else ElMessage.error(res.error || t('records.messages.deleteFailed'))
  } catch (err) { if (err !== 'cancel') console.warn(err) }
}

async function focusEventById(id) {
  if (!id) return
  const numericId = Number(id)
  if (!Number.isFinite(numericId) || numericId <= 0) return
  const existing = events.value.find((item) => Number(item.id) === numericId)
  if (existing) {
    await selectEvent(existing)
    return
  }
  const res = await getFinancialEvent(numericId)
  if (res.success) {
    selectedDetail.value = res.data
    selectedEvent.value = { id: numericId }
    await loadGroupEvents()
  }
}

async function refreshAll() {
  applyRouteQuery(route.query)
  await loadBase()
  await loadEvents()
  if (route.query?.eventId) {
    await focusEventById(route.query.eventId)
    router.replace({ path: route.path, query: {} })
  } else if (Object.keys(route.query || {}).length) {
    router.replace({ path: route.path, query: {} })
  }
}

watch(
  () => route.query,
  (next) => {
    if (route.path !== '/records') return
    if (!next || !Object.keys(next).length) return
    refreshAll()
  },
)
onMounted(refreshAll)
onActivated(refreshAll)
</script>

<style scoped>
.records-page{max-width:1200px;margin:0 auto;color:#1e293b}.reveal-block{animation:revealUp 480ms cubic-bezier(.16,1,.3,1) both}.delay-1{animation-delay:90ms}.delay-2{animation-delay:160ms}.page-hero,.filters-card,.records-list,.detail-panel{background:#fff;border-radius:16px;box-shadow:0 1px 3px rgba(15,23,42,.1),0 12px 30px rgba(15,23,42,.04)}.page-hero{margin-bottom:18px;padding:26px;background:linear-gradient(135deg,#fffaf0 0%,#fff 70%)}.eyebrow{margin:0 0 8px;color:#2f7d5c;font-size:12px;font-weight:800;letter-spacing:.08em;text-transform:uppercase}.page-hero h1{margin:0;font-size:30px;line-height:1.16;letter-spacing:-.022em}.page-hero p:last-child{max-width:720px;margin:12px 0 0;color:#64748b;line-height:1.7}.filters-card{display:grid;gap:14px;margin-bottom:18px;padding:18px;background:linear-gradient(135deg,#fffaf0 0%,#fff 72%)}.filter-card-head{display:flex;justify-content:space-between;gap:14px;align-items:flex-start}.filter-card-head h2{margin:0;font-size:18px;letter-spacing:-.012em}.filter-card-actions{display:flex;flex-wrap:wrap;gap:8px;justify-content:flex-end}.export-action{min-height:40px;border:0;border-radius:12px;padding:0 14px;background:rgba(47,125,92,.10);color:#2f7d5c;font-size:12px;font-weight:900;cursor:pointer;box-shadow:inset 0 0 0 1px rgba(47,125,92,.22);transition:transform 160ms,box-shadow 160ms,background-color 160ms}.export-action:hover{background:rgba(47,125,92,.16);box-shadow:inset 0 0 0 1px rgba(47,125,92,.34)}.export-action:active{transform:scale(.96)}.filter-banner{display:flex;align-items:center;gap:10px;padding:10px 14px;border-radius:12px;background:rgba(47,125,92,.1);box-shadow:inset 0 0 0 1px rgba(47,125,92,.2)}.filter-banner .banner-label{color:#2f7d5c;font-size:11px;font-weight:900;letter-spacing:.06em;text-transform:uppercase}.filter-banner strong{color:#1e293b;font-size:14px}.filter-banner .banner-clear{margin-left:auto;display:grid;place-items:center;width:26px;height:26px;border:0;border-radius:999px;background:#fff;color:#2f7d5c;font-weight:900;cursor:pointer;box-shadow:inset 0 0 0 1px rgba(47,125,92,.18)}.filter-banner .banner-clear:hover{background:#eef7f0}.filter-banner .banner-clear:active{transform:scale(.92)}.filter-row{display:grid;grid-template-columns:74px minmax(0,1fr);gap:12px;align-items:start}.filter-label{padding-top:8px;color:#64748b;font-size:12px;font-weight:900;letter-spacing:.04em}.filter-chip-rail{display:flex;gap:8px;flex-wrap:wrap}.filter-chip{min-height:34px;border:0;border-radius:999px;padding:0 13px;background:#fff;color:#64748b;font-size:12px;font-weight:900;cursor:pointer;box-shadow:inset 0 0 0 1px rgba(100,116,139,.14);transition-property:transform,box-shadow,background-color,color;transition-duration:160ms}.filter-chip:active{transform:scale(.96)}.filter-chip.active{background:#1f2933;color:rgba(255,255,255,.92);box-shadow:0 8px 16px rgba(31,41,51,.14)}.currency-chip{font-family:'SF Mono','Fira Code',monospace}.filter-search-row{display:flex;gap:12px;align-items:center}.keyword-control{max-width:360px}.advanced-filter-panel{display:flex;flex-wrap:wrap;gap:12px;padding:14px;border-radius:14px;background:#fff;box-shadow:inset 0 0 0 1px rgba(100,116,139,.12)}.filter-control{width:180px}.date-range-control{width:280px}.records-layout{display:grid;grid-template-columns:minmax(0,1fr) 380px;gap:18px}.records-list,.detail-panel{padding:16px}.record-group{margin-bottom:4px}.record-row{display:grid;grid-template-columns:28px 10px minmax(0,1fr) auto;gap:12px;align-items:center;padding:14px;border-radius:14px;cursor:pointer;transition-property:transform,box-shadow,background-color;transition-duration:160ms}.record-row.active{background:#f8faf7;box-shadow:0 0 0 2px rgba(47,125,92,.16)}.group-expand,.group-expand-placeholder{display:grid;place-items:center;width:26px;height:26px;border:0;border-radius:999px;background:rgba(47,125,92,.10);color:#2f7d5c;font-weight:900}.group-expand{cursor:pointer}.group-expand:active{transform:scale(.92)}.group-expand-placeholder{background:transparent}.group-children{display:grid;gap:6px;margin:2px 0 8px 40px;padding:8px;border-left:2px solid rgba(47,125,92,.18);border-radius:0 12px 12px 0;background:#fffaf0}.group-child-row{display:grid;grid-template-columns:26px minmax(0,1fr) auto;gap:10px;align-items:center;width:100%;border:0;border-radius:10px;background:#fff;padding:9px 10px;color:inherit;text-align:left;cursor:pointer}.group-child-row.active{box-shadow:inset 0 0 0 2px rgba(47,125,92,.28);background:#f8faf7}.child-marker{display:grid;place-items:center;width:22px;height:22px;border-radius:999px;background:rgba(47,125,92,.10);color:#2f7d5c;font-size:12px;font-weight:900}.group-child-row span:nth-child(2){overflow:hidden;text-overflow:ellipsis;white-space:nowrap;font-weight:800}.group-child-row em{font-family:'SF Mono','Fira Code',monospace;font-style:normal;font-weight:800}.type-dot{width:10px;height:10px;border-radius:999px;background:#3b82f6}.record-row.income .type-dot{background:#ef4444}.record-row.expense .type-dot{background:#f97316}.record-row.exchange .type-dot,.record-row.transfer .type-dot,.record-row.family_transfer .type-dot{background:#3b82f6}.record-row.receivable_create .type-dot,.record-row.receivable_collect .type-dot,.record-row.deposit_create .type-dot,.record-row.deposit_refund .type-dot,.record-row.loan_out .type-dot,.record-row.loan_collect .type-dot{background:#2f7d5c}.record-row.investment_buy .type-dot,.record-row.investment_sell .type-dot,.record-row.investment_income .type-dot,.record-row.investment_revalue .type-dot{background:#7c3aed}.record-main strong,.record-main span{display:block;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}.record-main span{margin-top:4px;color:#64748b;font-size:12px}.record-amount,.event-card strong,.entry-row em{font-family:'SF Mono','Fira Code',monospace;font-variant-numeric:tabular-nums}.section-head{margin-bottom:18px}.section-head h2{margin:0;letter-spacing:-.012em}.event-card{padding:16px;border-radius:14px;background:#f8faf7}.event-card.income{background:rgba(239,68,68,.1)}.event-card.expense{background:rgba(249,115,22,.1)}.event-card.exchange,.event-card.transfer,.event-card.family_transfer{background:rgba(59,130,246,.1)}.event-card.receivable_create,.event-card.receivable_collect,.event-card.deposit_create,.event-card.deposit_refund,.event-card.loan_out,.event-card.loan_collect{background:rgba(47,125,92,.1)}.event-card.investment_buy,.event-card.investment_sell,.event-card.investment_income,.event-card.investment_revalue{background:rgba(124,58,237,.1)}.event-card strong{display:block;margin:8px 0;font-size:24px}.detail-grid{display:grid;grid-template-columns:120px minmax(0,1fr);gap:10px;margin-top:18px}.detail-grid span{color:#64748b}.related-link{appearance:none;border:0;background:transparent;color:#2f7d5c;font-weight:900;cursor:pointer;text-align:left;padding:0}.related-link:hover{color:#1f5f44;text-decoration:underline}.group-block{margin-top:20px}.group-block h3{margin:0 0 10px;display:flex;align-items:center;gap:8px}.group-block .group-count{color:#64748b;font-size:12px;font-weight:700}.group-list{display:grid;gap:8px}.group-row{display:grid;grid-template-columns:32px minmax(0,1fr) auto;gap:12px;align-items:center;padding:12px 14px;border:0;border-radius:12px;background:#fff;color:inherit;text-align:left;cursor:pointer;box-shadow:inset 0 0 0 1px rgba(100,116,139,.10);transition:transform 160ms,box-shadow 160ms,background-color 160ms}.group-row:hover{background:#f8faf7;box-shadow:inset 0 0 0 1px rgba(47,125,92,.24)}.group-row:active{transform:scale(.99)}.group-row.current{background:rgba(47,125,92,.10);box-shadow:inset 0 0 0 2px rgba(47,125,92,.36)}.group-marker{display:grid;place-items:center;width:28px;height:28px;border-radius:999px;background:rgba(47,125,92,.10);color:#2f7d5c;font-size:14px;font-weight:900}.group-row.expense .group-marker{background:rgba(249,115,22,.14);color:#f97316}.group-row.income .group-marker{background:rgba(239,68,68,.14);color:#ef4444}.group-row.transfer .group-marker,.group-row.exchange .group-marker{background:rgba(59,130,246,.14);color:#3b82f6}.group-main{display:flex;flex-direction:column;gap:2px;min-width:0}.group-main strong{color:#1e293b;font-size:13px;font-weight:800;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}.group-main small{color:#64748b;font-size:11px;font-weight:700}.group-amount{display:flex;flex-direction:column;gap:2px;align-items:flex-end}.group-amount strong{color:#1e293b;font-family:'SF Mono','Fira Code',monospace;font-variant-numeric:tabular-nums;font-size:14px;font-weight:800}.group-amount .group-meta{color:#f97316;font-size:11px;font-weight:800}.group-amount .group-meta.settled{color:#2f7d5c}.entries-block{margin-top:20px}.entry-row{display:flex;justify-content:space-between;gap:12px;padding:12px;border-radius:12px;background:#f8faf7;margin-bottom:8px}.entry-row span{display:block;color:#64748b;font-size:12px}.positive{color:#ef4444}.negative{color:#f97316}.detail-actions{display:flex;gap:10px;margin-top:20px}.primary-action,.ghost-action,.danger-action{min-height:40px;border:0;border-radius:12px;padding:0 16px;font-weight:800;cursor:pointer;transition-property:transform,box-shadow,background-color,color;transition-duration:160ms}.primary-action:active,.ghost-action:active,.danger-action:active,.record-row:active{transform:scale(.96)}.primary-action{background:#3b82f6;color:#fff;box-shadow:0 10px 24px rgba(59,130,246,.22)}.ghost-action{background:#f8faf7;color:#1e293b}.danger-action{background:rgba(239,68,68,.1);color:#ef4444}.empty-state{display:grid;place-items:center;min-height:240px;color:#64748b}.empty-state.compact{min-height:260px}.pagination-wrap{display:flex;justify-content:flex-end;margin-top:16px;padding:12px;border-radius:14px;background:#fffaf0;box-shadow:inset 0 0 0 1px rgba(100,116,139,.1)}.readonly-note{color:#64748b;line-height:1.7}.full-width{width:100%}@media(hover:hover){.record-row:hover{transform:translateY(-2px);box-shadow:0 10px 22px rgba(15,23,42,.08)}.filter-chip:hover{background:#f8faf7;color:#2f7d5c}}@media(max-width:980px){.filter-card-head,.filter-search-row{flex-direction:column;align-items:stretch}.filter-card-actions{justify-content:flex-start}.filter-row{grid-template-columns:1fr}.filter-label{padding-top:0}.filter-control,.keyword-control,.date-range-control{width:100%;max-width:none}.records-layout{grid-template-columns:1fr}}@media(max-width:640px){.page-hero,.filters-card,.records-list,.detail-panel{border-radius:14px}.record-row{grid-template-columns:26px 10px minmax(0,1fr);gap:10px}.record-amount{grid-column:3;justify-self:start;margin-top:4px}.group-children{margin-left:28px}.group-child-row{grid-template-columns:24px minmax(0,1fr)}.group-child-row em{grid-column:2}.detail-panel{padding:14px}.detail-grid{grid-template-columns:1fr}.entry-row{flex-direction:column}.pagination-wrap{justify-content:flex-start;overflow-x:auto}.export-action,.ghost-action{width:100%}}@keyframes revealUp{from{opacity:0;transform:translateY(12px);filter:blur(4px)}to{opacity:1;transform:translateY(0);filter:blur(0)}}
</style>
