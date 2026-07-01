<template>
  <main class="import-page">
    <section class="import-workbench reveal-block">
      <ManagementPageHeader :eyebrow="t('import.hero.eyebrow')" :title="t('import.hero.title')" :subtitle="t('import.hero.subtitle')" />

      <section class="import-surface import-setup" v-loading="loadingRefs">
        <div class="setup-grid">
          <div class="setup-field">
            <label>{{ t('import.fields.config') }}</label>
            <el-select v-model="configId" filterable class="full-width" :placeholder="t('import.fields.configPlaceholder')" @change="onConfigChange">
              <el-option v-for="item in configs" :key="item.id" :label="`${item.name} · ${item.fileType.toUpperCase()}`" :value="item.id" />
            </el-select>
          </div>
          <div class="setup-field">
            <label>{{ t('import.fields.fallbackBucket') }}<span class="optional-tag">{{ t('import.fields.optional') }}</span></label>
            <el-select v-model="fallbackBucketId" filterable clearable class="full-width" :placeholder="t('import.fields.fallbackBucketPlaceholder')">
              <el-option v-for="item in buckets" :key="item.id" :label="bucketLabel(item)" :value="item.id" />
            </el-select>
          </div>
          <div class="setup-field setup-file-field">
            <label>{{ t('import.fields.file') }}</label>
            <el-upload ref="uploadRef" :auto-upload="false" :show-file-list="false" accept=".xlsx,.csv" :on-change="onFileChange" class="file-uploader">
              <button class="file-trigger" type="button">
                <span class="file-trigger-icon">📎</span>
                <span class="file-trigger-text">{{ fileName || t('import.fields.chooseFile') }}</span>
                <span v-if="fileName" class="file-trigger-clear" role="button" :aria-label="t('common.actions.delete')" @click.stop="clearFile">✕</span>
              </button>
            </el-upload>
          </div>
        </div>
        <div class="setup-footer">
          <p class="setup-hint">{{ t('import.fields.fallbackBucketHint') }}</p>
          <button class="management-primary-action setup-preview-btn" type="button" :disabled="!canPreview" @click="runPreview">{{ t('import.actions.preview') }}</button>
        </div>
      </section>

      <section v-if="warnings.length" class="import-surface import-warnings">
        <strong>{{ t('import.warnings.title') }}</strong>
        <ul>
          <li v-for="(w, i) in warnings" :key="i">{{ w }}</li>
        </ul>
      </section>

      <section v-if="preview" class="import-surface import-preview-panel">
        <div class="preview-head">
          <div>
            <strong>{{ t('import.table.title') }}</strong>
            <span>{{ t('import.table.summary', { total: preview.totalRows, errors: preview.errorRowCount, filtered: preview.filteredRowCount }) }}</span>
          </div>
          <div class="preview-actions">
            <button class="management-primary-action" type="button" :disabled="!canCommit" @click="onConfirmCommit">
              {{ committing ? t('import.actions.committing') : t('import.actions.confirm') }}
            </button>
          </div>
        </div>

        <div class="row-filter-bar">
          <div class="row-filter-chips">
            <button type="button" class="row-filter-chip" :class="{ active: rowFilter === 'all' }" @click="setRowFilter('all')">
              {{ t('import.filter.all', { count: previewRows.length }) }}
            </button>
            <button type="button" class="row-filter-chip danger" :class="{ active: rowFilter === 'errors' }" :disabled="errorRowCount === 0" @click="setRowFilter('errors')">
              {{ t('import.filter.errorsOnly', { count: errorRowCount }) }}
            </button>
          </div>
          <button type="button" class="management-ghost-action small" :disabled="errorRowCount === 0" @click="jumpToNextError">
            {{ t('import.filter.jumpNext') }}
          </button>
        </div>

        <div v-if="committing || commitProgress.processed > 0" class="commit-progress">
          <el-progress :percentage="progressPercent" :status="progressStatus" :stroke-width="10" />
          <p class="commit-progress-text">
            {{ t('import.progress.summary', { processed: commitProgress.processed, total: commitProgress.total, success: commitProgress.success, failed: commitProgress.failed }) }}
          </p>
        </div>

        <el-table :data="paginatedRows" class="preview-table" border size="small" :row-class-name="rowClassName">
          <el-table-column :label="t('import.table.columns.row')" width="56" align="center" prop="rowIndex" />
          <el-table-column :label="t('import.table.columns.scenario')" width="120">
            <template #default="{ row }">
              <el-select v-model="row.scenario" size="small" class="cell-select" @change="revalidateRow(row)">
                <el-option v-for="opt in scenarioOptions" :key="opt" :label="t(`record.scenarios.${opt}`)" :value="opt" />
              </el-select>
            </template>
          </el-table-column>
          <el-table-column :label="t('import.table.columns.amount')" width="120">
            <template #default="{ row }"><el-input v-model="row.amount" size="small" @blur="revalidateRow(row)" /></template>
          </el-table-column>
          <el-table-column :label="t('import.table.columns.currency')" width="110">
            <template #default="{ row }">
              <el-select v-model="row.currency" size="small" filterable class="cell-select" @change="revalidateRow(row)">
                <el-option v-for="item in currencyOptions" :key="item.code" :label="getCurrencyLabel(item.code, config.locale)" :value="item.code" />
              </el-select>
            </template>
          </el-table-column>
          <el-table-column :label="t('import.table.columns.eventTime')" width="180">
            <template #default="{ row }"><el-input v-model="row.eventTime" size="small" @blur="revalidateRow(row)" /></template>
          </el-table-column>
          <el-table-column :label="t('import.table.columns.category')" width="150">
            <template #default="{ row }">
              <el-select v-if="isEditing(row, 'category')" v-model="row.categoryId" size="small" filterable clearable class="cell-select" automatic-dropdown @visible-change="onEditorVisibleChange(row, 'category', $event)" @change="revalidateRow(row)">
                <el-option v-for="item in categoryOptions(row.scenario)" :key="item.id" :label="item.name" :value="item.id" />
              </el-select>
              <button v-else class="cell-view" type="button" @click="beginEdit(row, 'category')">{{ categoryLabel(row.categoryId) || t('import.table.emptyValue') }}</button>
            </template>
          </el-table-column>
          <el-table-column :label="t('import.table.columns.channel')" width="150">
            <template #default="{ row }">
              <el-select v-if="isEditing(row, 'channel')" v-model="row.channelId" size="small" filterable clearable class="cell-select" automatic-dropdown @visible-change="onEditorVisibleChange(row, 'channel', $event)">
                <el-option v-for="item in channels" :key="item.id" :label="`${item.icon || '🔗'} ${item.name}`" :value="item.id" />
              </el-select>
              <button v-else class="cell-view" type="button" @click="beginEdit(row, 'channel')">{{ channelLabel(row.channelId) || t('import.table.emptyValue') }}</button>
            </template>
          </el-table-column>
          <el-table-column :label="t('import.table.columns.bucket')" width="180">
            <template #default="{ row }">
              <template v-if="row.scenario === 'transfer'">
                <div class="transfer-bucket-row">
                  <span class="transfer-tag">{{ t('import.table.transferFrom') }}</span>
                  <el-select v-if="isEditing(row, 'bucket')" v-model="row.bucketId" size="small" filterable class="cell-select" automatic-dropdown @visible-change="onEditorVisibleChange(row, 'bucket', $event)" @change="revalidateRow(row)">
                    <el-option v-for="item in buckets" :key="item.id" :label="item.name" :value="item.id" />
                  </el-select>
                  <button v-else class="cell-view" type="button" @click="beginEdit(row, 'bucket')">{{ bucketLabelById(row.bucketId) || t('import.table.emptyValue') }}</button>
                </div>
                <div class="transfer-bucket-row">
                  <span class="transfer-tag transfer-tag-to">{{ t('import.table.transferTo') }}</span>
                  <el-select v-if="isEditing(row, 'toBucket')" v-model="row.toBucketId" size="small" filterable class="cell-select" automatic-dropdown @visible-change="onEditorVisibleChange(row, 'toBucket', $event)" @change="revalidateRow(row)">
                    <el-option v-for="item in buckets" :key="item.id" :label="item.name" :value="item.id" />
                  </el-select>
                  <button v-else class="cell-view" type="button" @click="beginEdit(row, 'toBucket')">{{ bucketLabelById(row.toBucketId) || t('import.table.emptyValue') }}</button>
                </div>
              </template>
              <template v-else>
                <el-select v-if="isEditing(row, 'bucket')" v-model="row.bucketId" size="small" filterable class="cell-select" automatic-dropdown @visible-change="onEditorVisibleChange(row, 'bucket', $event)" @change="revalidateRow(row)">
                  <el-option v-for="item in buckets" :key="item.id" :label="item.name" :value="item.id" />
                </el-select>
                <button v-else class="cell-view" type="button" @click="beginEdit(row, 'bucket')">{{ bucketLabelById(row.bucketId) || t('import.table.emptyValue') }}</button>
              </template>
            </template>
          </el-table-column>
          <el-table-column :label="t('import.table.columns.description')" min-width="160">
            <template #default="{ row }"><el-input v-model="row.description" size="small" /></template>
          </el-table-column>
          <el-table-column :label="t('import.table.columns.remark')" min-width="130">
            <template #default="{ row }"><el-input v-model="row.remark" size="small" /></template>
          </el-table-column>
          <el-table-column :label="t('import.table.columns.status')" width="180">
            <template #default="{ row }">
              <span v-if="row.filtered" class="row-filtered" :title="row.filterReason">{{ t('import.table.filtered') }}</span>
              <span v-else-if="row.errors && row.errors.length" class="row-errors">
                <em v-for="(e, i) in row.errors" :key="i">{{ e }}</em>
              </span>
              <span v-else class="row-ok">{{ t('import.table.ok') }}</span>
            </template>
          </el-table-column>
        </el-table>
        <div v-if="visibleRows.length > pageSize" class="preview-pagination">
          <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[50, 100]"
            :total="visibleRows.length"
            layout="total, sizes, prev, pager, next, jumper"
            background
            small
          />
        </div>
      </section>

      <section v-else-if="previewError" class="import-surface import-empty">
        <p class="empty-text">{{ previewError }}</p>
      </section>

      <section v-else class="import-surface import-empty">
        <ManagementEmptyState :image="marmotOne" :alt="t('import.empty.alt')" :title="t('import.empty.title')" :text="t('import.empty.text')" />
      </section>
    </section>
  </main>
</template>

<script setup>
import { computed, onActivated, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { listImportConfigs, previewImport, commitImport } from '@/api/importConfig/importConfig'
import { listCategories } from '@/api/category/category'
import { listChannels } from '@/api/channel/channel'
import { listBuckets } from '@/api/bucket/bucket'
import { currencyOptions, getCurrencyLabel } from '@/utils/currency'
import { useConfigStore } from '@/stores/config'
import ManagementPageHeader from '@/components/management/ManagementPageHeader.vue'
import ManagementEmptyState from '@/components/management/ManagementEmptyState.vue'
import marmotOne from '../../../img/marmot-ledger-1.png'

const { t } = useI18n()
const config = useConfigStore()

const configs = ref([])
const buckets = ref([])
const categories = ref([])
const channels = ref([])
const loadingRefs = ref(false)

const configId = ref(null)
const fallbackBucketId = ref(null)
const fileName = ref('')
const rawFile = ref(null)
const uploadRef = ref(null)
const preview = ref(null)
const previewRows = ref([])
const previewError = ref('')
const currentPage = ref(1)
const rowFilter = ref('all') // 'all' | 'errors'
const pageSize = ref(50)
const committing = ref(false)
const COMMIT_CHUNK_SIZE = 100
const commitProgress = ref({ total: 0, processed: 0, success: 0, failed: 0, aborted: false })

const scenarioOptions = ['income', 'expense', 'refund', 'transfer']
const warnings = computed(() => preview.value?.warnings || [])
const canPreview = computed(() => !!configId.value && !!rawFile.value)
const pendingRows = computed(() => previewRows.value.filter((r) => !r.filtered))
const blockingRow = computed(() => {
  for (const r of pendingRows.value) {
    if (r.errors && r.errors.length) return r
    if (r.scenario === 'transfer') {
      if (!r.bucketId || !r.toBucketId) return r
    } else if (!r.bucketId) {
      return r
    }
  }
  return null
})
const canCommit = computed(() => !!preview.value && !committing.value && pendingRows.value.length > 0 && !blockingRow.value)
const progressPercent = computed(() => {
  const total = commitProgress.value.total || 1
  return Math.min(100, Math.round((commitProgress.value.processed / total) * 100))
})
const progressStatus = computed(() => {
  if (commitProgress.value.aborted) return 'exception'
  if (commitProgress.value.failed > 0 && !committing.value) return 'warning'
  if (progressPercent.value === 100 && !committing.value) return 'success'
  return ''
})
const visibleRows = computed(() => {
  const base = previewRows.value.filter((r) => !r.filtered)
  if (rowFilter.value === 'errors') return base.filter((r) => r.errors && r.errors.length)
  return base
})
const errorRowCount = computed(() => previewRows.value.reduce((n, r) => n + ((!r.filtered && r.errors && r.errors.length) ? 1 : 0), 0))
function setRowFilter(kind) {
  rowFilter.value = kind
  currentPage.value = 1
}
function jumpToNextError() {
  const base = previewRows.value.filter((r) => !r.filtered)
  if (!base.length) return
  // 以 visibleRows 中当前页尾行为锚，找到 base 里下一个错行
  const pageEnd = currentPage.value * pageSize.value
  const anchorRow = paginatedRows.value[paginatedRows.value.length - 1] || paginatedRows.value[0]
  const anchorIndex = anchorRow ? base.findIndex((r) => r.rowIndex === anchorRow.rowIndex) : -1
  // 从锚之后循环找；如无则回环到开头
  const startFrom = anchorIndex + 1
  const nextIdx = findErrorIndex(base, startFrom) ?? findErrorIndex(base, 0)
  if (nextIdx == null) return
  // 若当前处于 errors-only 模式，直接翻页；否则也直接翻到目标行所在页
  const positionInBase = nextIdx
  const filteredList = rowFilter.value === 'errors' ? base.filter((r) => r.errors && r.errors.length) : base
  const targetRow = base[positionInBase]
  const positionInVisible = filteredList.indexOf(targetRow)
  if (positionInVisible < 0) return
  currentPage.value = Math.floor(positionInVisible / pageSize.value) + 1
}
function findErrorIndex(list, from) {
  for (let i = from; i < list.length; i++) {
    if (list[i].errors && list[i].errors.length) return i
  }
  return null
}
const paginatedRows = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  return visibleRows.value.slice(start, start + pageSize.value)
})

function bucketLabel(bucket) { return `${bucket.name} · ${bucket.currency}` }
function bucketLabelById(id) {
  if (!id) return ''
  const b = buckets.value.find((item) => item.id === Number(id))
  return b ? b.name : `#${id}`
}
function categoryLabel(id) {
  if (!id) return ''
  const c = categories.value.find((item) => item.id === Number(id))
  return c ? c.name : `#${id}`
}
function channelLabel(id) {
  if (!id) return ''
  const c = channels.value.find((item) => item.id === Number(id))
  return c ? `${c.icon || '🔗'} ${c.name}` : `#${id}`
}
function categoryOptions(scenario) {
  if (scenario === 'income') return categories.value.filter((item) => item.type === 'income')
  return categories.value.filter((item) => item.type === 'expense')
}

// ------- 单行合规校验：跟后端 record pipeline 的规则保持一致 -------
const VALID_SCENARIOS = new Set(['income', 'expense', 'refund', 'transfer'])
const EVENT_TIME_PATTERN = /^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$/
function isFiniteAmount(value) {
  if (value == null || value === '') return false
  const cleaned = String(value).replace(/[¥$€,，\s]/g, '')
  if (!cleaned) return false
  const n = Number(cleaned)
  return Number.isFinite(n)
}
function findBucket(id) { return id ? buckets.value.find((b) => b.id === Number(id)) : null }
function findCategory(id) { return id ? categories.value.find((c) => c.id === Number(id)) : null }
function findChannel(id) { return id ? channels.value.find((c) => c.id === Number(id)) : null }

// 单行校验：返回 errors 数组；不修改 row，让调用者决定何时写入 row.errors。
function computeRowErrors(row) {
  if (row.filtered) return []
  const errors = []
  const scenario = row.scenario || ''
  if (!VALID_SCENARIOS.has(scenario)) {
    errors.push(t('import.errors.scenarioInvalid', { value: scenario || '' }))
  }
  // scenario ↔ bucket
  if (scenario === 'transfer') {
    if (!row.bucketId) errors.push(t('import.errors.fromBucketRequired'))
    if (!row.toBucketId) errors.push(t('import.errors.toBucketRequired'))
    if (row.bucketId && row.toBucketId && row.bucketId === row.toBucketId) errors.push(t('import.errors.transferSameBucket'))
  } else if (!row.bucketId) {
    errors.push(t('import.errors.bucketRequired'))
  }
  // amount
  if (!row.amount) errors.push(t('import.errors.amountEmpty'))
  else if (!isFiniteAmount(row.amount)) errors.push(t('import.errors.amountInvalid', { value: row.amount }))
  // eventTime
  if (!row.eventTime) errors.push(t('import.errors.eventTimeEmpty'))
  else if (!EVENT_TIME_PATTERN.test(row.eventTime)) errors.push(t('import.errors.eventTimeInvalid', { value: row.eventTime }))
  // currency + bucket 币种一致性
  const currency = (row.currency || '').toUpperCase()
  if (!currency) errors.push(t('import.errors.currencyEmpty'))
  const bucket = findBucket(row.bucketId)
  if (bucket && currency && bucket.currency && bucket.currency !== currency) {
    errors.push(t('import.errors.currencyBucketMismatch', { bucket: bucket.name, bucketCurrency: bucket.currency, rowCurrency: currency }))
  }
  if (scenario === 'transfer') {
    const toBucket = findBucket(row.toBucketId)
    if (toBucket && currency && toBucket.currency && toBucket.currency !== currency) {
      errors.push(t('import.errors.currencyBucketMismatch', { bucket: toBucket.name, bucketCurrency: toBucket.currency, rowCurrency: currency }))
    }
  }
  // category type 必须与 scenario 匹配（只对 income/expense/refund 有意义，transfer 不带 category）
  if (row.categoryId) {
    const cat = findCategory(row.categoryId)
    if (!cat) {
      errors.push(t('import.errors.categoryNotFound', { id: row.categoryId }))
    } else if (scenario === 'income' && cat.type !== 'income') {
      errors.push(t('import.errors.categoryTypeMismatchIncome', { name: cat.name }))
    } else if ((scenario === 'expense' || scenario === 'refund') && cat.type !== 'expense') {
      errors.push(t('import.errors.categoryTypeMismatchExpense', { name: cat.name }))
    }
  }
  // channel 存在性（可选字段，非零时校验）
  if (row.channelId && !findChannel(row.channelId)) {
    errors.push(t('import.errors.channelNotFound', { id: row.channelId }))
  }
  return errors
}

// 单行重算：在 change 事件里触发；只跑这一行，不遍历全表。
function revalidateRow(row) {
  row.errors = computeRowErrors(row)
  refreshErrorCount()
}
function refreshErrorCount() {
  if (!preview.value) return
  let errCount = 0
  for (const r of previewRows.value) {
    if (!r.filtered && r.errors && r.errors.length) errCount++
  }
  preview.value.errorRowCount = errCount
  // 若停在 "仅异常" 视图但已经没有错行，自动回退到 "全部"，避免用户看到空表迷惑
  if (rowFilter.value === 'errors' && errCount === 0) {
    rowFilter.value = 'all'
    currentPage.value = 1
  }
}
// 全量校验：仅在预览刚返回、reference 数据可用时跑一次。
function revalidateAll() {
  for (const row of previewRows.value) row.errors = computeRowErrors(row)
  refreshErrorCount()
}

// 懒编辑：默认渲染为静态文本，点击后临时挂载 el-select，失焦回到文本态
const editing = ref({ rowIndex: 0, field: '' })
function isEditing(row, field) { return editing.value.rowIndex === row.rowIndex && editing.value.field === field }
function beginEdit(row, field) { editing.value = { rowIndex: row.rowIndex, field } }
function onEditorVisibleChange(row, field, visible) {
  if (!visible && isEditing(row, field)) editing.value = { rowIndex: 0, field: '' }
}
function rowClassName({ row }) {
  if (row.filtered) return 'filtered-row'
  if (row.errors && row.errors.length) return 'error-row'
  return ''
}

function onConfigChange() {
  preview.value = null
  previewRows.value = []
  currentPage.value = 1
  previewError.value = ''
}
function onFileChange(file) {
  // 每次选择都覆盖上一次；不用 el-upload 的 :limit，因为 limit=1 会阻止重新选择。
  rawFile.value = file.raw
  fileName.value = file.name
  preview.value = null
  previewRows.value = []
  currentPage.value = 1
  previewError.value = ''
  // 保留 el-upload 内部只留最新一条，避免其 fileList 无限增长。
  if (uploadRef.value && uploadRef.value.uploadFiles) {
    uploadRef.value.uploadFiles = uploadRef.value.uploadFiles.slice(-1)
  }
}
function clearFile() {
  rawFile.value = null
  fileName.value = ''
  preview.value = null
  previewRows.value = []
  currentPage.value = 1
  previewError.value = ''
  if (uploadRef.value && typeof uploadRef.value.clearFiles === 'function') {
    uploadRef.value.clearFiles()
  }
}

async function loadReferences() {
  loadingRefs.value = true
  try {
    const [cfgRes, bucketRes, incomeCats, expenseCats, channelRes] = await Promise.all([
      listImportConfigs({ isActive: true }),
      listBuckets(),
      listCategories({ type: 'income', isActive: true }),
      listCategories({ type: 'expense', isActive: true }),
      listChannels({ isActive: true }),
    ])
    configs.value = cfgRes.success ? cfgRes.data || [] : []
    buckets.value = bucketRes.success ? bucketRes.data || [] : []
    categories.value = [...(incomeCats.success ? incomeCats.data || [] : []), ...(expenseCats.success ? expenseCats.data || [] : [])]
    channels.value = channelRes.success ? channelRes.data || [] : []
  } finally { loadingRefs.value = false }
}

async function runPreview() {
  if (!canPreview.value) return
  preview.value = null
  previewRows.value = []
  rowFilter.value = 'all'
  currentPage.value = 1
  previewError.value = ''
  try {
    const res = await previewImport(configId.value, rawFile.value, fallbackBucketId.value)
    if (res.success) {
      preview.value = res.data
      previewRows.value = (res.data.rows || []).map((r) => ({ ...r, toBucketId: r.toBucketId || 0 }))
      revalidateAll()
      if (!previewRows.value.length) ElMessage.warning(t('import.messages.noRows'))
    } else {
      previewError.value = res.error || t('import.messages.previewFailed')
      ElMessage.error(previewError.value)
    }
  } catch (err) {
    previewError.value = t('import.messages.previewFailed')
    ElMessage.error(previewError.value)
  }
}

function buildCommitRows() {
  const list = []
  for (const row of previewRows.value) {
    if (row.filtered) continue
    if (row.errors && row.errors.length) continue
    const scenario = row.scenario
    const req = {
      scenario,
      bucketId: 0,
      fromBucketId: 0,
      toBucketId: 0,
      categoryId: Number(row.categoryId || 0),
      channelId: Number(row.channelId || 0),
      amount: String(row.amount || '0'),
      currency: row.currency || 'CNY',
      description: row.description || '',
      remark: row.remark || '',
      eventTime: row.eventTime || '',
      source: 'import',
    }
    if (scenario === 'transfer') {
      req.fromBucketId = Number(row.bucketId || 0)
      req.toBucketId = Number(row.toBucketId || 0)
    } else {
      req.bucketId = Number(row.bucketId || 0)
    }
    list.push({ rowIndex: row.rowIndex, request: req })
  }
  return list
}

function applyChunkResult(chunk, data) {
  const successSet = new Set(data.successRowIndexes || [])
  const failMap = new Map((data.failures || []).map((f) => [f.rowIndex, f.error]))
  // 从 previewRows 精准移除成功行
  if (successSet.size) {
    previewRows.value = previewRows.value.filter((r) => !successSet.has(r.rowIndex))
  }
  // 给失败行追加原因
  for (const row of previewRows.value) {
    if (failMap.has(row.rowIndex)) {
      row.errors = [...(row.errors || []), `${t('import.messages.commitFailed')}: ${failMap.get(row.rowIndex)}`]
    }
  }
  // 顶部统计
  if (preview.value) {
    preview.value.totalRows = previewRows.value.length
    preview.value.filteredRowCount = previewRows.value.filter((r) => r.filtered).length
    preview.value.errorRowCount = previewRows.value.filter((r) => !r.filtered && r.errors && r.errors.length).length
  }
  // 若停在 "仅异常" 视图但已经没有错行，自动回退到 "全部"
  if (rowFilter.value === 'errors' && preview.value && preview.value.errorRowCount === 0) {
    rowFilter.value = 'all'
    currentPage.value = 1
  }
  // 分页边界：基于当前可见行数
  const visibleLen = visibleRows.value.length
  if (currentPage.value > 1 && (currentPage.value - 1) * pageSize.value >= visibleLen) {
    currentPage.value = Math.max(1, Math.ceil(visibleLen / pageSize.value))
  }
}

async function onConfirmCommit() {
  if (!canCommit.value) {
    const bad = blockingRow.value
    if (bad) ElMessage.warning(t('import.messages.commitBlocked', { row: bad.rowIndex }))
    return
  }
  const allRows = buildCommitRows()
  if (allRows.length === 0) return
  committing.value = true
  commitProgress.value = { total: allRows.length, processed: 0, success: 0, failed: 0, aborted: false }
  let failedAt = -1
  try {
    for (let start = 0; start < allRows.length; start += COMMIT_CHUNK_SIZE) {
      const chunk = allRows.slice(start, start + COMMIT_CHUNK_SIZE)
      let res
      try {
        res = await commitImport(configId.value, chunk)
      } catch (err) {
        // 网络错误或超时：这一片状态未知，中止后续片，让用户先去 Records 核对
        failedAt = start
        commitProgress.value.aborted = true
        break
      }
      if (!res.success) {
        failedAt = start
        commitProgress.value.aborted = true
        break
      }
      applyChunkResult(chunk, res.data)
      commitProgress.value.processed = Math.min(start + chunk.length, allRows.length)
      commitProgress.value.success += res.data.successCount || 0
      commitProgress.value.failed += res.data.failedCount || 0
    }
    // 结果消息
    const p = commitProgress.value
    if (p.aborted) {
      ElMessage.error(t('import.messages.commitAborted', { processed: p.processed, total: p.total, from: failedAt + 1 }))
    } else if (p.failed > 0) {
      ElMessage.warning(t('import.messages.commitPartial', { success: p.success, failed: p.failed }))
    } else {
      ElMessage.success(t('import.messages.commitDone', { success: p.success }))
    }
  } finally {
    committing.value = false
  }
}

onMounted(loadReferences)
onActivated(loadReferences)
</script>

<style scoped>
.import-page { padding: 24px; min-height: 100%; background: #f6f0e6; }
.import-workbench { max-width: 1280px; margin: 0 auto; }
.import-surface { border-radius: 16px; background: #fff; box-shadow: 0 1px 3px rgba(15, 23, 42, .1), 0 12px 30px rgba(15, 23, 42, .04); margin-bottom: 18px; }
.import-setup { padding: 20px; }
.setup-grid { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 14px; align-items: start; }
.setup-field label { display: flex; align-items: center; gap: 6px; margin-bottom: 6px; color: #6f6254; font-size: 12px; font-weight: 700; }
.setup-field { min-width: 0; }
.optional-tag { color: #a8997f; font-weight: 600; font-size: 11px; }
.setup-file-field { display: flex; flex-direction: column; }
.file-uploader { width: 100%; }
.file-uploader :deep(.el-upload) { width: 100%; display: block; }
.file-trigger {
  width: 100%; min-height: 32px; display: flex; align-items: center; gap: 8px;
  border: 1px dashed #d8cab2; border-radius: 8px; padding: 0 12px;
  background: #faf6ee; color: #6f6254; font-size: 13px; font-weight: 600; cursor: pointer;
  transition: border-color 160ms ease, background-color 160ms ease;
}
.file-trigger:hover { border-color: #2f7d5c; background: #f3f8f3; }
.file-trigger-icon { font-size: 14px; }
.file-trigger-text { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; flex: 1 1 auto; }
.file-trigger-clear {
  flex: 0 0 auto; display: inline-grid; place-items: center; width: 20px; height: 20px;
  margin-left: auto; border-radius: 6px; background: rgba(239, 68, 68, 0.10); color: #ef4444;
  font-weight: 800; font-size: 11px; cursor: pointer;
}
.file-trigger-clear:hover { background: rgba(239, 68, 68, 0.18); }
.setup-footer { display: flex; align-items: center; justify-content: space-between; gap: 16px; margin-top: 16px; padding-top: 14px; border-top: 1px solid #efe7d8; }
.setup-hint { margin: 0; color: #857462; font-size: 13px; line-height: 1.6; }
.setup-preview-btn { flex: 0 0 auto; }

.import-warnings { padding: 16px 20px; background: linear-gradient(135deg, #fff8e6 0%, #fff 70%); }
.import-warnings strong { color: #b8860b; font-size: 13px; }
.import-warnings ul { margin: 8px 0 0; padding-left: 18px; color: #9a7b3a; font-size: 13px; line-height: 1.7; }

.import-preview-panel { padding: 16px; }
.preview-head { display: flex; align-items: center; justify-content: space-between; gap: 12px; margin-bottom: 12px; }
.preview-head strong { color: #21362d; font-size: 18px; }
.preview-head span { margin-left: 10px; color: #857462; font-size: 13px; }
.preview-table { width: 100%; }
.preview-pagination { display: flex; justify-content: flex-end; margin-top: 12px; }
.cell-view {
  width: 100%; min-height: 24px; padding: 2px 6px; border: 0; border-radius: 6px;
  background: transparent; color: #21362d; font: inherit; text-align: left; cursor: pointer;
  overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
  transition: background-color 120ms ease;
}
.cell-view:hover { background: rgba(47, 125, 92, 0.08); }
.cell-view:focus { outline: 2px solid rgba(47, 125, 92, 0.35); outline-offset: -2px; }

.transfer-bucket-row { display: flex; align-items: center; gap: 6px; min-width: 0; }
.transfer-bucket-row + .transfer-bucket-row { margin-top: 4px; }
.transfer-bucket-row .cell-view, .transfer-bucket-row .cell-select { flex: 1 1 auto; min-width: 0; }
.transfer-tag {
  flex: 0 0 auto; display: inline-flex; align-items: center; min-height: 18px;
  padding: 0 6px; border-radius: 4px; font-size: 10px; font-weight: 800; letter-spacing: 0.02em;
  background: rgba(249, 115, 22, 0.12); color: #c2410c;
}
.transfer-tag-to { background: rgba(47, 125, 92, 0.12); color: #2f7d5c; }
.cell-select { width: 100%; }

.preview-table :deep(.error-row) { background: #fef2f2; }
.preview-table :deep(.error-row:hover > td) { background: #fee2e2; }
.preview-table :deep(.filtered-row) { background: #f5f5f4; opacity: 0.55; }
.preview-table :deep(.filtered-row:hover > td) { background: #e7e5e4; }
.row-errors { display: flex; flex-direction: column; gap: 2px; }
.row-errors em { font-style: normal; color: #ef4444; font-size: 12px; }
.row-ok { color: #2f7d5c; font-size: 12px; font-weight: 700; }
.row-filtered { color: #78716c; font-size: 12px; font-weight: 700; }
.preview-actions button:disabled { opacity: 0.5; cursor: not-allowed; }
.row-filter-bar { display: flex; align-items: center; justify-content: space-between; gap: 12px; margin: 0 0 12px; }
.row-filter-chips { display: flex; gap: 8px; }
.row-filter-chip {
  min-height: 30px; padding: 0 12px; border: 0; border-radius: 999px;
  background: #f7f3eb; color: #6f6254; font-size: 12px; font-weight: 800; cursor: pointer;
  transition: background-color 160ms ease, color 160ms ease, box-shadow 160ms ease;
}
.row-filter-chip:hover { background: #ece2d2; }
.row-filter-chip.active { background: #21362d; color: #fff; }
.row-filter-chip.danger.active { background: #ef4444; }
.row-filter-chip:disabled { opacity: 0.4; cursor: not-allowed; }
.row-filter-chip:disabled:hover { background: #f7f3eb; }
.commit-progress { margin: 4px 0 14px; padding: 12px 14px; border-radius: 12px; background: #fffaf0; box-shadow: inset 0 0 0 1px rgba(100, 116, 139, 0.10); }
.commit-progress-text { margin: 8px 0 0; color: #857462; font-size: 12px; font-weight: 700; }

.import-empty { padding: 24px; }
.empty-text { color: #ef4444; }

@media (max-width: 980px) {
  .setup-grid { grid-template-columns: 1fr; }
  .setup-footer { flex-direction: column; align-items: stretch; }
  .setup-preview-btn { width: 100%; }
}
</style>
