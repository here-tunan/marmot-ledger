<template>
  <main class="record-page">
    <section class="page-hero reveal-block">
      <div>
        <p class="eyebrow">{{ t('record.hero.eyebrow') }}</p>
        <h1>{{ t('record.hero.title') }}</h1>
        <p>{{ t('record.hero.subtitle') }}</p>
      </div>
      <img :src="marmotOne" :alt="t('dashboard.images.brandAlt')" width="92" height="92" />
    </section>

    <section class="record-layout reveal-block delay-1">
      <aside class="scenario-panel">
        <button
          v-for="item in scenarios"
          :key="item.value"
          class="scenario-card"
          :class="[{ active: form.scenario === item.value }, item.value, { disabled: item.disabled }]"
          :disabled="item.disabled"
          @click="selectScenario(item.value)"
        >
          <span>{{ item.icon }}</span>
          <strong>{{ t(item.labelKey) }}</strong>
          <small v-if="item.disabled">{{ t('record.scenarios.comingSoon') }}</small>
        </button>
      </aside>

      <section class="form-panel">
        <div class="section-head">
          <div>
            <p class="eyebrow">{{ t(selectedScenario.labelKey) }}</p>
            <h2>{{ t('record.actions.submit') }}</h2>
          </div>
        </div>

        <el-form ref="formRef" :model="form" label-position="top">
          <template v-if="isTransfer">
            <el-form-item :label="t('record.fields.fromBucket')">
              <el-select v-model="form.fromBucketId" class="full-width" filterable>
                <el-option v-for="bucket in buckets" :key="bucket.id" :label="bucketLabel(bucket)" :value="bucket.id" />
              </el-select>
            </el-form-item>
            <el-form-item :label="t('record.fields.toBucket')">
              <el-select v-model="form.toBucketId" class="full-width" filterable>
                <el-option v-for="bucket in buckets" :key="bucket.id" :label="bucketLabel(bucket)" :value="bucket.id" />
              </el-select>
            </el-form-item>
          </template>

          <el-form-item v-else :label="t('record.fields.bucket')">
            <el-select v-model="form.bucketId" class="full-width" filterable>
              <el-option v-for="bucket in buckets" :key="bucket.id" :label="bucketLabel(bucket)" :value="bucket.id" />
            </el-select>
          </el-form-item>

          <el-form-item v-if="needsCategory" :label="t('record.fields.category')">
            <el-select v-model="form.categoryId" class="full-width" filterable :placeholder="t('record.placeholders.selectCategory')">
              <el-option v-for="category in categories" :key="category.id" :label="categoryLabel(category)" :value="category.id" />
            </el-select>
          </el-form-item>

          <el-row :gutter="12">
            <el-col :span="12">
              <el-form-item :label="t('record.fields.amount')">
                <el-input v-model="form.amount" placeholder="36.50" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item :label="t('record.fields.currency')">
                <el-select v-model="form.currency" class="full-width">
                  <el-option v-for="item in currencyOptions" :key="item.code" :label="getCurrencyLabel(item.code, config.locale)" :value="item.code" />
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>

          <el-form-item :label="t('record.fields.description')">
            <el-input v-model="form.description" />
          </el-form-item>
          <el-form-item :label="t('record.fields.eventTime')">
            <el-input v-model="form.eventTime" placeholder="YYYY-MM-DD HH:mm:ss" />
          </el-form-item>
          <el-form-item v-if="form.scenario === 'refund'" :label="t('record.fields.relatedEvent')">
            <el-input v-model="form.relatedFinancialEventId" />
          </el-form-item>
          <el-form-item :label="t('record.fields.remark')">
            <el-input v-model="form.remark" type="textarea" :autosize="{ minRows: 2, maxRows: 4 }" />
          </el-form-item>

          <button class="primary-action" type="button" :disabled="submitting" @click="submitRecord">
            {{ submitting ? '...' : t('record.actions.submit') }}
          </button>
        </el-form>
      </section>

      <aside class="preview-panel">
        <div class="section-head">
          <div>
            <p class="eyebrow">{{ t('record.preview.title') }}</p>
            <h2>{{ t('record.preview.event') }}</h2>
          </div>
        </div>

        <div class="event-preview" :class="form.scenario">
          <span>{{ eventTypeLabel(form.scenario) }}</span>
          <strong>{{ form.currency }} {{ formatAmount(form.amount) }}</strong>
          <p>{{ includeInStatistics ? t('record.preview.included') : t('record.preview.excluded') }}</p>
        </div>

        <div class="entry-preview">
          <h3>{{ t('record.preview.entries') }}</h3>
          <div v-if="previewEntries.length" class="entry-list">
            <div v-for="entry in previewEntries" :key="entry.role + entry.bucket" class="entry-row">
              <span>{{ entry.bucket }}</span>
              <strong :class="{ negative: entry.amount.startsWith('-'), positive: entry.amount.startsWith('+') }">{{ entry.amount }}</strong>
            </div>
          </div>
          <p v-else>{{ t('record.preview.selectScenario') }}</p>
        </div>

        <div class="recent-events">
          <h3>{{ t('dashboard.events.title') }}</h3>
          <div v-for="event in recentEvents" :key="event.id" class="recent-row">
            <span>{{ eventTypeLabel(event.eventType) }}</span>
            <strong>{{ event.currency }} {{ formatAmount(event.amount) }}</strong>
          </div>
        </div>
      </aside>
    </section>
  </main>
</template>

<script setup>
import { computed, onActivated, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { useConfigStore } from '@/stores/config'
import { listBuckets } from '@/api/bucket/bucket'
import { listFinancialEvents } from '@/api/financialEvent/financialEvent'
import { createRecord } from '@/api/record/record'
import { listCategories } from '@/api/category/category'
import marmotOne from '../../../img/marmot-ledger-1.png'
import { currencyOptions, getCurrencyLabel } from '@/utils/currency'

const { t, te } = useI18n()
const config = useConfigStore()
const buckets = ref([])
const categories = ref([])
const recentEvents = ref([])
const submitting = ref(false)
const formRef = ref()

const form = reactive({
  scenario: 'expense',
  bucketId: '',
  fromBucketId: '',
  toBucketId: '',
  categoryId: '',
  amount: '',
  currency: 'CNY',
  description: '',
  eventTime: formatDateTime(new Date()),
  relatedFinancialEventId: '',
  remark: '',
})

const scenarios = [
  { value: 'income', labelKey: 'record.scenarios.income', icon: '↗', disabled: false },
  { value: 'expense', labelKey: 'record.scenarios.expense', icon: '↘', disabled: false },
  { value: 'transfer', labelKey: 'record.scenarios.transfer', icon: '⇄', disabled: false },
  { value: 'refund', labelKey: 'record.scenarios.refund', icon: '↩', disabled: false },
  { value: 'exchange', labelKey: 'record.scenarios.exchange', icon: 'FX', disabled: true },
  { value: 'deposit', labelKey: 'record.scenarios.deposit', icon: 'D', disabled: true },
  { value: 'receivable', labelKey: 'record.scenarios.receivable', icon: 'R', disabled: true },
  { value: 'loan', labelKey: 'record.scenarios.loan', icon: 'L', disabled: true },
  { value: 'investment', labelKey: 'record.scenarios.investment', icon: 'I', disabled: true },
]

const selectedScenario = computed(() => scenarios.find((item) => item.value === form.scenario) || scenarios[1])
const isTransfer = computed(() => form.scenario === 'transfer')
const needsCategory = computed(() => form.scenario === 'income' || form.scenario === 'expense' || form.scenario === 'refund')
const includeInStatistics = computed(() => form.scenario === 'income' || form.scenario === 'expense')

const previewEntries = computed(() => {
  const amount = Number(form.amount || 0)
  if (!amount) return []
  if (isTransfer.value) {
    return [
      { role: 'transfer_out', bucket: bucketName(form.fromBucketId), amount: `-${formatAmount(amount)}` },
      { role: 'transfer_in', bucket: bucketName(form.toBucketId), amount: `+${formatAmount(amount)}` },
    ]
  }
  const sign = form.scenario === 'expense' ? '-' : '+'
  return [{ role: form.scenario, bucket: bucketName(form.bucketId), amount: `${sign}${formatAmount(amount)}` }]
})

function formatDateTime(date) {
  const pad = (value) => String(value).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`
}

function enumLabel(prefix, value) {
  const key = `${prefix}.${value}`
  return value && te(key) ? t(key) : value
}

function eventTypeLabel(type) {
  return enumLabel('record.scenarios', type)
}

function formatAmount(value) {
  const number = Number(value || 0)
  return new Intl.NumberFormat(config.locale, { minimumFractionDigits: 2, maximumFractionDigits: 2 }).format(number)
}

function bucketLabel(bucket) {
  return `${bucket.name} · ${bucket.currency} ${formatAmount(bucket.balance)}`
}

function bucketName(id) {
  return buckets.value.find((item) => item.id === Number(id))?.name || '-'
}

function categoryLabel(category) {
  return category.categoryGroupName ? `${category.name} · ${category.categoryGroupName}` : category.name
}

function selectScenario(value) {
  form.scenario = value
  form.categoryId = ''
  loadCategories()
}

async function loadBuckets() {
  const res = await listBuckets({ isActive: true })
  if (res.success) buckets.value = res.data || []
}

async function loadCategories() {
  if (!needsCategory.value) {
    categories.value = []
    form.categoryId = ''
    return
  }
  const type = form.scenario === 'income' ? 'income' : 'expense'
  const res = await listCategories({ type, isActive: true })
  if (res.success) categories.value = res.data || []
  else ElMessage.error(res.error || t('record.messages.loadCategoriesFailed'))
}

async function loadRecentEvents() {
  const res = await listFinancialEvents({ page: 1, pageSize: 5 })
  if (res.success) recentEvents.value = res.data?.list || []
}

async function refreshAll() {
  await Promise.all([loadBuckets(), loadCategories(), loadRecentEvents()])
}

function validateForm() {
  if (!Number(form.amount || 0) || Number(form.amount || 0) <= 0) {
    ElMessage.warning(t('record.messages.amountRequired'))
    return false
  }
  if (isTransfer.value) {
    if (!form.fromBucketId || !form.toBucketId) {
      ElMessage.warning(t('record.messages.selectTransferBuckets'))
      return false
    }
    return true
  }
  if (!form.bucketId) {
    ElMessage.warning(t('record.messages.selectBucket'))
    return false
  }
  if (needsCategory.value && !form.categoryId) {
    ElMessage.warning(t('record.messages.selectCategory'))
    return false
  }
  return true
}

async function submitRecord() {
  if (!validateForm()) return
  submitting.value = true
  try {
    const payload = {
      ...form,
      bucketId: Number(form.bucketId || 0),
      fromBucketId: Number(form.fromBucketId || 0),
      toBucketId: Number(form.toBucketId || 0),
      categoryId: Number(form.categoryId || 0),
      relatedFinancialEventId: Number(form.relatedFinancialEventId || 0),
      amount: String(form.amount),
    }
    const res = await createRecord(payload)
    if (res.success) {
      ElMessage.success(t('record.messages.created'))
      form.amount = ''
      form.description = ''
      form.remark = ''
      form.relatedFinancialEventId = ''
      form.eventTime = formatDateTime(new Date())
      await refreshAll()
    } else {
      ElMessage.error(res.error || t('record.messages.createFailed'))
    }
  } finally {
    submitting.value = false
  }
}

onMounted(refreshAll)
onActivated(refreshAll)
</script>

<style scoped>
.record-page {
  max-width: 1200px;
  margin: 0 auto;
  color: #1e293b;
}

.reveal-block {
  animation: revealUp 500ms cubic-bezier(0.16, 1, 0.3, 1) both;
}

.delay-1 { animation-delay: 100ms; }

.page-hero,
.scenario-panel,
.form-panel,
.preview-panel {
  background: #ffffff;
  border-radius: 16px;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.1), 0 12px 30px rgba(15, 23, 42, 0.04);
}

.page-hero {
  display: flex;
  justify-content: space-between;
  gap: 24px;
  align-items: center;
  margin-bottom: 18px;
  padding: 26px;
  background: linear-gradient(135deg, #fffaf0 0%, #ffffff 70%);
}

.page-hero img {
  border-radius: 24px;
  box-shadow: 0 14px 32px rgba(47, 125, 92, 0.16);
}

.eyebrow {
  margin: 0 0 8px;
  color: #2f7d5c;
  font-size: 12px;
  font-weight: 800;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.page-hero h1,
.section-head h2 {
  margin: 0;
  letter-spacing: -0.022em;
}

.page-hero h1 {
  max-width: 720px;
  font-size: 30px;
  line-height: 1.16;
}

.page-hero p:last-child {
  max-width: 680px;
  margin: 12px 0 0;
  color: #64748b;
  line-height: 1.7;
}

.record-layout {
  display: grid;
  grid-template-columns: 220px minmax(0, 1fr) 330px;
  gap: 18px;
}

.scenario-panel,
.form-panel,
.preview-panel {
  padding: 18px;
}

.scenario-panel {
  display: grid;
  gap: 10px;
  align-content: start;
}

.scenario-card {
  min-height: 64px;
  border: 0;
  border-radius: 14px;
  display: grid;
  grid-template-columns: 34px minmax(0, 1fr);
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  background: #f8faf7;
  color: #1e293b;
  text-align: left;
  cursor: pointer;
  transition-property: transform, box-shadow, background-color, color;
  transition-duration: 160ms;
}

.scenario-card:active,
.primary-action:active {
  transform: scale(0.96);
}

.scenario-card span {
  display: grid;
  place-items: center;
  width: 34px;
  height: 34px;
  border-radius: 12px;
  background: rgba(31, 41, 51, 0.08);
  font-weight: 900;
}

.scenario-card strong,
.scenario-card small {
  display: block;
}

.scenario-card small {
  grid-column: 2;
  color: #94a3b8;
  font-size: 11px;
}

.scenario-card.active.income {
  background: rgba(239, 68, 68, 0.12);
  box-shadow: 0 0 0 2px rgba(239, 68, 68, 0.18);
}

.scenario-card.active.expense {
  background: rgba(16, 185, 129, 0.12);
  box-shadow: 0 0 0 2px rgba(16, 185, 129, 0.18);
}

.scenario-card.active.transfer,
.scenario-card.active.refund {
  background: rgba(59, 130, 246, 0.12);
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.18);
}

.scenario-card.disabled {
  cursor: not-allowed;
  opacity: 0.48;
}

.section-head {
  display: flex;
  justify-content: space-between;
  margin-bottom: 18px;
}

.full-width {
  width: 100%;
}

.primary-action {
  min-height: 42px;
  border: 0;
  border-radius: 12px;
  padding: 0 18px;
  background: #3b82f6;
  color: #ffffff;
  font-weight: 800;
  cursor: pointer;
  box-shadow: 0 10px 24px rgba(59, 130, 246, 0.22);
  transition-property: transform, box-shadow, background-color, color;
  transition-duration: 160ms;
}

.event-preview {
  padding: 16px;
  border-radius: 14px;
  background: #f8faf7;
}

.event-preview.income {
  background: rgba(239, 68, 68, 0.1);
}

.event-preview.expense {
  background: rgba(16, 185, 129, 0.1);
}

.event-preview span,
.event-preview p {
  color: #64748b;
  font-weight: 700;
}

.event-preview strong {
  display: block;
  margin: 8px 0;
  font-family: 'SF Mono', 'Fira Code', monospace;
  font-size: 24px;
}

.entry-preview,
.recent-events {
  margin-top: 18px;
}

.entry-preview h3,
.recent-events h3 {
  margin: 0 0 10px;
}

.entry-list,
.recent-events {
  display: grid;
  gap: 10px;
}

.entry-row,
.recent-row {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  padding: 12px;
  border-radius: 12px;
  background: #ffffff;
  box-shadow: inset 0 0 0 1px rgba(100, 116, 139, 0.12);
}

.entry-row strong,
.recent-row strong {
  font-family: 'SF Mono', 'Fira Code', monospace;
}

.negative {
  color: #10b981;
}

.positive {
  color: #ef4444;
}

@media (hover: hover) {
  .scenario-card:not(.disabled):hover {
    transform: translateY(-2px);
    box-shadow: 0 10px 22px rgba(15, 23, 42, 0.08);
  }
}

@media (max-width: 1080px) {
  .record-layout {
    grid-template-columns: 1fr;
  }

  .scenario-panel {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 600px) {
  .page-hero {
    align-items: flex-start;
    flex-direction: column;
  }

  .page-hero h1 {
    font-size: 24px;
  }

  .scenario-panel {
    grid-template-columns: 1fr;
  }
}

@media (prefers-reduced-motion: reduce) {
  .reveal-block,
  .scenario-card,
  .primary-action {
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
