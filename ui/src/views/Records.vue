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
      <el-select v-model="filters.eventType" clearable :placeholder="t('records.filters.eventType')" class="filter-control" @change="resetAndLoad">
        <el-option v-for="item in editableTypes" :key="item" :label="eventTypeLabel(item)" :value="item" />
        <el-option :label="eventTypeLabel('balance_adjustment')" value="balance_adjustment" />
      </el-select>
      <el-select v-model="filters.categoryId" clearable :placeholder="t('records.filters.category')" class="filter-control" @change="resetAndLoad">
        <el-option v-for="item in categories" :key="item.id" :label="item.name" :value="item.id" />
      </el-select>
      <el-input v-model="filters.keyword" :placeholder="t('records.filters.keyword')" class="keyword-control" @keyup.enter="resetAndLoad" />
      <el-select v-model="filters.includeInStatistics" clearable :placeholder="t('records.filters.includeStats')" class="filter-control" @change="resetAndLoad">
        <el-option :label="t('record.preview.included')" :value="true" />
        <el-option :label="t('record.preview.excluded')" :value="false" />
      </el-select>
      <button class="ghost-action" @click="resetAndLoad">{{ t('common.actions.refresh') }}</button>
    </section>

    <section class="records-layout reveal-block delay-2">
      <div v-loading="loading" class="records-list">
        <article v-for="event in events" :key="event.id" class="record-row" :class="[{ active: selectedEvent?.id === event.id }, event.eventType]" @click="selectEvent(event)">
          <div class="type-dot"></div>
          <div class="record-main">
            <strong>{{ event.description || eventTypeLabel(event.eventType) }}</strong>
            <span>{{ event.eventTime }} · {{ eventTypeLabel(event.eventType) }}</span>
          </div>
          <div class="record-amount">{{ event.currency }} {{ formatAmount(event.amount) }}</div>
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
          </div>

          <div class="entries-block">
            <h3>{{ t('records.detail.entries') }}</h3>
            <div v-for="entry in selectedDetail.ledgerEntries || []" :key="entry.id" class="entry-row">
              <div>
                <strong>{{ bucketName(entry.bucketId) }}</strong>
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
import { computed, onActivated, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useConfigStore } from '@/stores/config'
import { listFinancialEvents, getFinancialEvent } from '@/api/financialEvent/financialEvent'
import { updateRecord, deleteRecord } from '@/api/record/record'
import { listBuckets } from '@/api/bucket/bucket'
import { listCategories } from '@/api/category/category'
import { currencyOptions, getCurrencyLabel } from '@/utils/currency'

const { t, te } = useI18n()
const config = useConfigStore()
const events = ref([])
const selectedEvent = ref(null)
const selectedDetail = ref(null)
const buckets = ref([])
const categories = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const editingId = ref(0)
const editableTypes = ['income', 'expense', 'transfer', 'refund']
const filters = reactive({ eventType: '', categoryId: '', keyword: '', includeInStatistics: '' })
const pagination = reactive({ page: 1, pageSize: 10, total: 0 })
const form = reactive({ scenario: 'expense', bucketId: '', fromBucketId: '', toBucketId: '', categoryId: '', amount: '', currency: 'CNY', description: '', eventTime: '', relatedFinancialEventId: '', remark: '' })

const filteredCategories = computed(() => categories.value.filter((item) => form.scenario === 'income' ? item.type === 'income' : item.type === 'expense'))

function enumLabel(prefix, value) {
  const key = `${prefix}.${value}`
  return value && te(key) ? t(key) : value
}
function eventTypeLabel(type) { return enumLabel('record.scenarios', type) }
function entryRoleLabel(role) { return enumLabel('record.entryRoles', role) }
function formatAmount(value) { return new Intl.NumberFormat(config.locale, { minimumFractionDigits: 2, maximumFractionDigits: 2 }).format(Number(value || 0)) }
function bucketName(id) { return buckets.value.find((item) => item.id === Number(id))?.name || `#${id}` }
function bucketLabel(bucket) { return `${bucket.name} · ${bucket.currency} ${formatAmount(bucket.balance)}` }
function categoryName(id) { return categories.value.find((item) => item.id === Number(id))?.name || '-' }
function canMutate(event) { return editableTypes.includes(event.eventType) }

async function loadBase() {
  const [bucketRes, incomeCats, expenseCats] = await Promise.all([listBuckets(), listCategories({ type: 'income', isActive: true }), listCategories({ type: 'expense', isActive: true })])
  if (bucketRes.success) buckets.value = bucketRes.data || []
  categories.value = [...(incomeCats.success ? incomeCats.data || [] : []), ...(expenseCats.success ? expenseCats.data || [] : [])]
}

async function loadEvents() {
  loading.value = true
  try {
    const params = { page: pagination.page, pageSize: pagination.pageSize }
    if (filters.eventType) params.eventType = filters.eventType
    if (filters.categoryId) params.categoryId = filters.categoryId
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
  if (res.success) selectedDetail.value = res.data
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

async function refreshAll() { await loadBase(); await loadEvents() }
onMounted(refreshAll)
onActivated(refreshAll)
</script>

<style scoped>
.records-page{max-width:1200px;margin:0 auto;color:#1e293b}.reveal-block{animation:revealUp 480ms cubic-bezier(.16,1,.3,1) both}.delay-1{animation-delay:90ms}.delay-2{animation-delay:160ms}.page-hero,.filters-card,.records-list,.detail-panel{background:#fff;border-radius:16px;box-shadow:0 1px 3px rgba(15,23,42,.1),0 12px 30px rgba(15,23,42,.04)}.page-hero{margin-bottom:18px;padding:26px;background:linear-gradient(135deg,#fffaf0 0%,#fff 70%)}.eyebrow{margin:0 0 8px;color:#2f7d5c;font-size:12px;font-weight:800;letter-spacing:.08em;text-transform:uppercase}.page-hero h1{margin:0;font-size:30px;line-height:1.16;letter-spacing:-.022em}.page-hero p:last-child{max-width:720px;margin:12px 0 0;color:#64748b;line-height:1.7}.filters-card{display:flex;gap:12px;align-items:center;margin-bottom:18px;padding:14px}.filter-control{width:170px}.keyword-control{width:220px}.records-layout{display:grid;grid-template-columns:minmax(0,1fr) 380px;gap:18px}.records-list,.detail-panel{padding:16px}.record-row{display:grid;grid-template-columns:10px minmax(0,1fr) auto;gap:12px;align-items:center;padding:14px;border-radius:14px;cursor:pointer;transition-property:transform,box-shadow,background-color;transition-duration:160ms}.record-row.active{background:#f8faf7;box-shadow:0 0 0 2px rgba(47,125,92,.16)}.type-dot{width:10px;height:10px;border-radius:999px;background:#3b82f6}.record-row.income .type-dot{background:#ef4444}.record-row.expense .type-dot{background:#10b981}.record-main strong,.record-main span{display:block;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}.record-main span{margin-top:4px;color:#64748b;font-size:12px}.record-amount,.event-card strong,.entry-row em{font-family:'SF Mono','Fira Code',monospace;font-variant-numeric:tabular-nums}.section-head{margin-bottom:18px}.section-head h2{margin:0;letter-spacing:-.012em}.event-card{padding:16px;border-radius:14px;background:#f8faf7}.event-card.income{background:rgba(239,68,68,.1)}.event-card.expense{background:rgba(16,185,129,.1)}.event-card strong{display:block;margin:8px 0;font-size:24px}.detail-grid{display:grid;grid-template-columns:120px minmax(0,1fr);gap:10px;margin-top:18px}.detail-grid span{color:#64748b}.entries-block{margin-top:20px}.entry-row{display:flex;justify-content:space-between;gap:12px;padding:12px;border-radius:12px;background:#f8faf7;margin-bottom:8px}.entry-row span{display:block;color:#64748b;font-size:12px}.positive{color:#ef4444}.negative{color:#10b981}.detail-actions{display:flex;gap:10px;margin-top:20px}.primary-action,.ghost-action,.danger-action{min-height:40px;border:0;border-radius:12px;padding:0 16px;font-weight:800;cursor:pointer;transition-property:transform,box-shadow,background-color,color;transition-duration:160ms}.primary-action:active,.ghost-action:active,.danger-action:active,.record-row:active{transform:scale(.96)}.primary-action{background:#3b82f6;color:#fff;box-shadow:0 10px 24px rgba(59,130,246,.22)}.ghost-action{background:#f8faf7;color:#1e293b}.danger-action{background:rgba(239,68,68,.1);color:#ef4444}.empty-state{display:grid;place-items:center;min-height:240px;color:#64748b}.empty-state.compact{min-height:260px}.pagination-wrap{display:flex;justify-content:flex-end;margin-top:16px;padding-top:14px;border-top:1px solid rgba(100,116,139,.12)}.readonly-note{color:#64748b;line-height:1.7}.full-width{width:100%}@media(hover:hover){.record-row:hover{transform:translateY(-2px);box-shadow:0 10px 22px rgba(15,23,42,.08)}}@media(max-width:980px){.filters-card{flex-direction:column;align-items:stretch}.filter-control,.keyword-control{width:100%}.records-layout{grid-template-columns:1fr}}@keyframes revealUp{from{opacity:0;transform:translateY(12px);filter:blur(4px)}to{opacity:1;transform:translateY(0);filter:blur(0)}}
</style>
