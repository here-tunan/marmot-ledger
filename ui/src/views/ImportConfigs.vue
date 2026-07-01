<template>
  <main class="import-config-page">
    <section class="import-config-workbench reveal-block">
      <ManagementPageHeader :eyebrow="t('importConfigs.hero.eyebrow')" :title="t('importConfigs.hero.title')" :subtitle="t('importConfigs.hero.subtitle')">
        <template #actions>
          <span class="import-count-pill">{{ t('importConfigs.summary.total', { count: configs.length }) }}</span>
          <button class="management-primary-action" type="button" @click="openCreate">{{ t('importConfigs.actions.new') }}</button>
        </template>
      </ManagementPageHeader>

      <section class="import-main-panel" v-loading="loading">
        <div class="import-toolbelt">
          <div>
            <strong>{{ t('importConfigs.hero.title') }}</strong>
            <span>{{ t('importConfigs.summary.active', { count: activeConfigs.length }) }}</span>
          </div>
          <button class="management-ghost-action" type="button" @click="loadConfigs">{{ t('common.actions.refresh') }}</button>
        </div>

        <div v-if="configs.length" class="import-list">
          <article v-for="item in configs" :key="item.id" class="import-row" :class="{ inactive: !item.isActive }">
            <div class="import-symbol">{{ item.icon || (item.fileType === 'csv' ? '📄' : '📊') }}</div>
            <div class="import-copy">
              <div class="import-title-line">
                <h2>{{ item.name }}</h2>
                <span :class="['management-status-tag', { active: item.isActive }]">{{ item.isActive ? t('common.status.enabled') : t('common.status.disabled') }}</span>
              </div>
              <div class="import-meta-line">
                <span class="import-type-tag">{{ item.fileType.toUpperCase() }}</span>
                <span class="import-meta-tag">{{ t('importConfigs.summary.mappingCount', { count: (item.mappings || []).length }) }}</span>
                <span class="import-meta-tag">{{ t('importConfigs.fields.sort') }} {{ item.sort || 0 }}</span>
              </div>
            </div>
            <div class="import-actions">
              <button class="management-text-action" type="button" @click="openEdit(item)">{{ t('common.actions.edit') }}</button>
              <button class="management-danger-action" type="button" @click="handleDelete(item)">{{ t('common.actions.delete') }}</button>
            </div>
          </article>
        </div>

        <ManagementEmptyState v-else-if="!loading" :image="marmotOne" :alt="t('importConfigs.empty.alt')" :title="t('importConfigs.empty.title')" :text="t('importConfigs.empty.text')">
          <template #action>
            <button class="management-primary-action" type="button" @click="openCreate">{{ t('importConfigs.actions.new') }}</button>
          </template>
        </ManagementEmptyState>
      </section>
    </section>

    <el-dialog v-model="dialogVisible" :title="editingId ? t('importConfigs.dialog.editTitle') : t('importConfigs.dialog.createTitle')" width="860px" class="marmot-dialog import-config-dialog">
      <el-form :model="form" label-position="top">
        <div class="import-form-grid">
          <el-form-item :label="t('importConfigs.fields.name')"><el-input v-model="form.name" /></el-form-item>
          <el-form-item :label="t('importConfigs.fields.fileType')">
            <el-select v-model="form.fileType" class="full-width">
              <el-option label="XLSX" value="xlsx" />
              <el-option label="CSV" value="csv" />
            </el-select>
          </el-form-item>
          <el-form-item v-if="form.fileType === 'xlsx'" :label="t('importConfigs.fields.sheetName')">
            <el-input v-model="form.sheetName" :placeholder="t('importConfigs.fields.sheetNamePlaceholder')" />
          </el-form-item>
          <el-form-item :label="t('importConfigs.fields.headerRow')"><el-input-number v-model="form.headerRow" :min="1" /></el-form-item>
          <el-form-item :label="t('importConfigs.fields.icon')"><el-input v-model="form.icon" maxlength="4" /></el-form-item>
          <el-form-item :label="t('importConfigs.fields.sort')"><el-input-number v-model="form.sort" :min="0" /></el-form-item>
          <StatusSwitchField v-if="editingId" v-model="form.isActive" />
        </div>

        <div class="mapping-editor">
          <h3 class="mapping-editor-title">{{ t('importConfigs.mapping.title') }}</h3>
          <p class="mapping-editor-hint">{{ t('importConfigs.mapping.hint') }}</p>

          <div v-for="field in targetFields" :key="field.key" class="mapping-card">
            <div class="mapping-card-head">
              <strong>{{ targetFieldLabel(field.key) }}</strong>
              <span class="mapping-kind-tag">{{ field.direct ? t('importConfigs.mapping.directTag') : t('importConfigs.mapping.ruleTag') }}</span>
            </div>

            <div class="mapping-card-body">
              <div class="mapping-field">
                <label>{{ t('importConfigs.mapping.sourceColumn') }}</label>
                <el-input v-model="form.mappings[field.key].sourceColumn" :placeholder="sourceColumnPlaceholder(field.key)" />
              </div>
              <div class="mapping-field">
                <label>{{ t('importConfigs.mapping.defaultValue') }}</label>
                <component
                  :is="valueControl(field).is"
                  v-model="form.mappings[field.key].defaultValue"
                  v-bind="valueControl(field).bind"
                  class="full-width"
                >
                  <el-option v-for="opt in valueControl(field).options" :key="opt.value" :label="opt.label" :value="opt.value" />
                </component>
              </div>
            </div>

            <div v-if="!field.direct" class="mapping-rules">
              <div class="mapping-rules-head">
                <span>{{ t('importConfigs.mapping.rules') }}</span>
                <button class="management-ghost-action small" type="button" @click="addRule(field.key)">{{ t('importConfigs.mapping.addRule') }}</button>
              </div>
              <p v-if="!form.mappings[field.key].rules.length" class="mapping-rules-empty">{{ t('importConfigs.mapping.noRules') }}</p>
              <div v-for="(rule, idx) in form.mappings[field.key].rules" :key="idx" class="rule-row">
                <el-input v-model="rule.matchColumn" :placeholder="t('importConfigs.mapping.matchColumn')" class="rule-col" />
                <el-select v-model="rule.operator" class="rule-op">
                  <el-option v-for="op in operators" :key="op" :label="t(`importConfigs.operators.${op}`)" :value="op" />
                </el-select>
                <el-input v-model="rule.matchValue" :placeholder="matchValuePlaceholder(rule.operator)" class="rule-val" />
                <span class="rule-arrow">→</span>
                <component
                  :is="valueControl(field).is"
                  v-model="rule.resultValue"
                  v-bind="valueControl(field).bind"
                  class="rule-result"
                >
                  <el-option v-for="opt in valueControl(field).options" :key="opt.value" :label="opt.label" :value="opt.value" />
                </component>
                <div class="rule-ops">
                  <button class="rule-move" type="button" :disabled="idx === 0" :title="t('importConfigs.mapping.moveUp')" @click="moveRule(field.key, idx, -1)">↑</button>
                  <button class="rule-move" type="button" :disabled="idx === form.mappings[field.key].rules.length - 1" :title="t('importConfigs.mapping.moveDown')" @click="moveRule(field.key, idx, 1)">↓</button>
                  <button class="rule-remove" type="button" :title="t('common.actions.delete')" @click="removeRule(field.key, idx)">✕</button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="filter-editor">
          <div class="filter-editor-head">
            <div>
              <h3 class="filter-editor-title">{{ t('importConfigs.filters.title') }}</h3>
              <p class="filter-editor-hint">{{ t('importConfigs.filters.hint') }}</p>
            </div>
            <button class="management-ghost-action small" type="button" @click="addFilter">{{ t('importConfigs.filters.addFilter') }}</button>
          </div>
          <p v-if="!form.filters.length" class="filter-empty">{{ t('importConfigs.filters.empty') }}</p>
          <div v-for="(filter, idx) in form.filters" :key="idx" class="filter-row">
            <el-select v-model="filter.action" class="filter-action">
              <el-option v-for="a in filterActions" :key="a" :label="t(`importConfigs.filters.actions.${a}`)" :value="a" />
            </el-select>
            <span class="filter-when">{{ t('importConfigs.filters.when') }}</span>
            <el-input v-model="filter.matchColumn" :placeholder="t('importConfigs.mapping.matchColumn')" class="filter-col" />
            <el-select v-model="filter.operator" class="filter-op">
              <el-option v-for="op in operators" :key="op" :label="t(`importConfigs.operators.${op}`)" :value="op" />
            </el-select>
            <el-input v-model="filter.matchValue" :placeholder="matchValuePlaceholder(filter.operator)" class="filter-val" />
            <div class="rule-ops">
              <button class="rule-move" type="button" :disabled="idx === 0" :title="t('importConfigs.mapping.moveUp')" @click="moveFilter(idx, -1)">↑</button>
              <button class="rule-move" type="button" :disabled="idx === form.filters.length - 1" :title="t('importConfigs.mapping.moveDown')" @click="moveFilter(idx, 1)">↓</button>
              <button class="rule-remove" type="button" :title="t('common.actions.delete')" @click="removeFilter(idx)">✕</button>
            </div>
          </div>
        </div>
      </el-form>
      <template #footer>
        <ManagementDialogFooter @cancel="dialogVisible = false" @submit="submitForm" />
      </template>
    </el-dialog>
  </main>
</template>

<script setup>
import { computed, onActivated, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElInput, ElSelect, ElMessage } from 'element-plus'
import { createImportConfig, deleteImportConfig, listImportConfigs, updateImportConfig } from '@/api/importConfig/importConfig'
import { listCategories } from '@/api/category/category'
import { listChannels } from '@/api/channel/channel'
import { listBuckets } from '@/api/bucket/bucket'
import { currencyOptions, getCurrencyLabel } from '@/utils/currency'
import { useConfigStore } from '@/stores/config'
import ManagementPageHeader from '@/components/management/ManagementPageHeader.vue'
import ManagementEmptyState from '@/components/management/ManagementEmptyState.vue'
import ManagementDialogFooter from '@/components/management/ManagementDialogFooter.vue'
import StatusSwitchField from '@/components/management/StatusSwitchField.vue'
import { confirmDelete, isCancelError } from '@/components/management/confirmDelete'
import marmotOne from '../../../img/marmot-ledger-1.png'

const { t } = useI18n()
const config = useConfigStore()
const configs = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const editingId = ref(0)
const categories = ref([])
const channels = ref([])
const buckets = ref([])

// 目标字段：direct=直接映射列原值；其余可设默认值+条件规则。
// valueType 决定 defaultValue / resultValue 用哪种控件。
const targetFields = [
  { key: 'amount', direct: true, valueType: 'text' },
  { key: 'description', direct: true, valueType: 'text' },
  { key: 'remark', direct: true, valueType: 'text' },
  { key: 'eventTime', direct: true, valueType: 'text' },
  { key: 'currency', direct: false, valueType: 'currency' },
  { key: 'scenario', direct: false, valueType: 'scenario' },
  { key: 'category', direct: false, valueType: 'category' },
  { key: 'channel', direct: false, valueType: 'channel' },
  { key: 'bucket', direct: false, valueType: 'bucket' },
]
const operators = ['contains', 'equals', 'notContains', 'notEquals', 'containsAny', 'notContainsAny', 'equalsAny', 'notEqualsAny']
const multiValueOperators = ['containsAny', 'notContainsAny', 'equalsAny', 'notEqualsAny']
const filterActions = ['drop', 'keep']
const scenarioOptions = ['income', 'expense', 'refund', 'transfer']

const form = reactive(createEmptyForm())
const activeConfigs = computed(() => configs.value.filter((item) => item.isActive))

function createEmptyMappings() {
  const map = {}
  for (const field of targetFields) {
    map[field.key] = { sourceColumn: '', defaultValue: field.key === 'currency' ? 'CNY' : '', rules: [] }
  }
  return map
}
function createEmptyForm() {
  return { name: '', fileType: 'xlsx', sheetName: '', headerRow: 1, icon: '📊', sort: 0, isActive: true, mappings: createEmptyMappings(), filters: [] }
}
function resetForm(data) {
  Object.assign(form, createEmptyForm())
  if (!data) return
  form.name = data.name || ''
  form.fileType = data.fileType || 'xlsx'
  form.sheetName = data.sheetName || ''
  form.headerRow = data.headerRow || 1
  form.icon = data.icon || '📊'
  form.sort = data.sort || 0
  form.isActive = data.isActive !== false
  // 回填后端返回的 mappings 数组到按字段索引的对象
  for (const mapping of data.mappings || []) {
    if (form.mappings[mapping.targetField]) {
      form.mappings[mapping.targetField] = {
        sourceColumn: mapping.sourceColumn || '',
        defaultValue: mapping.defaultValue || '',
        rules: (mapping.rules || []).map((r) => ({ matchColumn: r.matchColumn || '', operator: r.operator || 'contains', matchValue: r.matchValue || '', resultValue: r.resultValue || '' })),
      }
    }
  }
  form.filters = (data.filters || []).map((f) => ({ matchColumn: f.matchColumn || '', operator: f.operator || 'contains', matchValue: f.matchValue || '', action: f.action || 'drop' }))
}

function targetFieldLabel(key) { return t(`importConfigs.targetFields.${key}`) }
function sourceColumnPlaceholder(key) { return t('importConfigs.mapping.sourceColumnExample', { example: t(`importConfigs.columnExamples.${key}`) }) }
function matchValuePlaceholder(operator) { return multiValueOperators.includes(operator) ? t('importConfigs.mapping.matchValueMulti') : t('importConfigs.mapping.matchValue') }

function valueControl(field) {
  switch (field.valueType) {
    case 'currency':
      return { is: ElSelect, bind: { clearable: true, filterable: true }, options: currencyOptions.map((item) => ({ value: item.code, label: getCurrencyLabel(item.code, config.locale) })) }
    case 'scenario':
      return { is: ElSelect, bind: { clearable: true }, options: scenarioOptions.map((value) => ({ value, label: t(`record.scenarios.${value}`) })) }
    case 'category':
      return { is: ElSelect, bind: { clearable: true, filterable: true }, options: categories.value.map((item) => ({ value: String(item.id), label: item.name })) }
    case 'channel':
      return { is: ElSelect, bind: { clearable: true, filterable: true }, options: channels.value.map((item) => ({ value: String(item.id), label: `${item.icon || '🔗'} ${item.name}` })) }
    case 'bucket':
      return { is: ElSelect, bind: { clearable: true, filterable: true }, options: buckets.value.map((item) => ({ value: String(item.id), label: item.name })) }
    default:
      return { is: ElInput, bind: { placeholder: t('importConfigs.mapping.defaultValuePlaceholder') }, options: [] }
  }
}

function addRule(key) {
  form.mappings[key].rules.push({ matchColumn: '', operator: 'contains', matchValue: '', resultValue: '' })
}
function removeRule(key, idx) {
  form.mappings[key].rules.splice(idx, 1)
}
function moveRule(key, idx, delta) {
  const rules = form.mappings[key].rules
  const target = idx + delta
  if (target < 0 || target >= rules.length) return
  const [moved] = rules.splice(idx, 1)
  rules.splice(target, 0, moved)
}

function addFilter() {
  form.filters.push({ matchColumn: '', operator: 'contains', matchValue: '', action: 'drop' })
}
function removeFilter(idx) {
  form.filters.splice(idx, 1)
}
function moveFilter(idx, delta) {
  const target = idx + delta
  if (target < 0 || target >= form.filters.length) return
  const [moved] = form.filters.splice(idx, 1)
  form.filters.splice(target, 0, moved)
}

function buildFiltersPayload() {
  return form.filters.map((f) => ({ matchColumn: f.matchColumn || '', operator: f.operator, matchValue: f.matchValue || '', action: f.action || 'drop' }))
}

function buildMappingsPayload() {
  return targetFields.map((field) => {
    const entry = form.mappings[field.key]
    return {
      targetField: field.key,
      sourceColumn: entry.sourceColumn || '',
      defaultValue: entry.defaultValue == null ? '' : String(entry.defaultValue),
      rules: field.direct ? [] : entry.rules.map((r) => ({ matchColumn: r.matchColumn || '', operator: r.operator, matchValue: r.matchValue || '', resultValue: r.resultValue == null ? '' : String(r.resultValue) })),
    }
  })
}

async function loadConfigs() {
  loading.value = true
  try {
    const res = await listImportConfigs()
    if (res.success) configs.value = res.data || []
    else ElMessage.error(res.error || t('importConfigs.messages.loadFailed'))
  } finally { loading.value = false }
}
async function loadReferences() {
  const [incomeCats, expenseCats, channelRes, bucketRes] = await Promise.all([
    listCategories({ type: 'income', isActive: true }),
    listCategories({ type: 'expense', isActive: true }),
    listChannels({ isActive: true }),
    listBuckets(),
  ])
  categories.value = [...(incomeCats.success ? incomeCats.data || [] : []), ...(expenseCats.success ? expenseCats.data || [] : [])]
  channels.value = channelRes.success ? channelRes.data || [] : []
  buckets.value = bucketRes.success ? bucketRes.data || [] : []
}
async function refreshAll() { await Promise.all([loadConfigs(), loadReferences()]) }

function openCreate() {
  editingId.value = 0
  resetForm()
  dialogVisible.value = true
}
function openEdit(item) {
  editingId.value = item.id
  resetForm(item)
  dialogVisible.value = true
}
async function submitForm() {
  if (!form.name) return ElMessage.warning(t('importConfigs.validation.nameRequired'))
  const payload = {
    name: form.name,
    fileType: form.fileType,
    sheetName: form.sheetName,
    headerRow: form.headerRow,
    icon: form.icon,
    sort: form.sort,
    isActive: form.isActive !== false,
    mappings: buildMappingsPayload(),
    filters: buildFiltersPayload(),
  }
  const res = editingId.value ? await updateImportConfig(editingId.value, payload) : await createImportConfig(payload)
  if (res.success) {
    ElMessage.success(editingId.value ? t('importConfigs.messages.updated') : t('importConfigs.messages.created'))
    dialogVisible.value = false
    await loadConfigs()
  } else {
    ElMessage.error(res.error || t('importConfigs.messages.saveFailed'))
  }
}
async function handleDelete(item) {
  try {
    await confirmDelete({ message: t('importConfigs.delete.confirm', { name: item.name }), title: t('common.actions.delete'), confirmText: t('common.actions.delete'), cancelText: t('common.actions.cancel') })
    const res = await deleteImportConfig(item.id)
    if (res.success) {
      ElMessage.success(t('importConfigs.messages.deleted'))
      await loadConfigs()
    } else ElMessage.error(res.error || t('importConfigs.messages.deleteFailed'))
  } catch (err) {
    if (!isCancelError(err)) console.warn(err)
  }
}

onMounted(refreshAll)
onActivated(refreshAll)
</script>

<style scoped>
.import-config-page { padding: 24px; min-height: 100%; background: #f6f0e6; }
.import-config-workbench { max-width: 1180px; margin: 0 auto; }
.import-count-pill {
  min-height: 36px; display: inline-flex; align-items: center; padding: 0 12px; border-radius: 10px;
  background: #fff; color: #64748b; box-shadow: 0 1px 3px rgba(15, 23, 42, 0.10); font-size: 13px; font-weight: 700;
}
.import-main-panel { border-radius: 16px; background: #fff; box-shadow: 0 1px 3px rgba(15, 23, 42, .1), 0 12px 30px rgba(15, 23, 42, .04); padding: 16px; min-height: 480px; }
.import-toolbelt { display: flex; align-items: center; justify-content: space-between; gap: 16px; padding: 8px 4px 16px; }
.import-toolbelt strong { display: block; color: #21362d; font-size: 20px; letter-spacing: -0.012em; }
.import-toolbelt span { color: #857462; font-size: 13px; }

.import-list { display: grid; gap: 10px; }
.import-row {
  display: grid; grid-template-columns: 58px minmax(0, 1fr) auto; gap: 14px; align-items: center; padding: 14px;
  border-radius: 16px; background: #fff; box-shadow: 0 1px 3px rgba(15, 23, 42, .1), 0 12px 30px rgba(15, 23, 42, .04);
  transition-property: transform, box-shadow; transition-duration: 180ms;
}
.import-row.inactive { opacity: 0.64; }
.import-symbol { width: 52px; height: 52px; border-radius: 18px; display: grid; place-items: center; font-size: 24px; background: rgba(47, 125, 92, 0.10); }
.import-copy { min-width: 0; }
.import-title-line { display: flex; gap: 10px; align-items: center; flex-wrap: wrap; }
.import-title-line h2 { margin: 0; color: #21362d; font-size: 18px; letter-spacing: -0.012em; }
.import-meta-line { display: flex; gap: 6px; flex-wrap: wrap; margin-top: 7px; color: #8a7a67; font-size: 13px; }
.import-type-tag, .import-meta-tag {
  display: inline-flex; align-items: center; min-height: 20px; border-radius: 999px; padding: 0 8px; font-size: 11px; font-weight: 700; line-height: 1.4;
}
.import-type-tag { background: rgba(47, 125, 92, 0.10); color: #2f7d5c; }
.import-meta-tag { background: #f7f3eb; color: #7c6c5a; }
.import-actions { display: flex; gap: 6px; align-items: center; }

.import-form-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); column-gap: 14px; }

.mapping-editor { margin-top: 8px; border-top: 1px solid #efe7d8; padding-top: 16px; }
.mapping-editor-title { margin: 0; color: #21362d; font-size: 16px; }
.mapping-editor-hint { margin: 6px 0 14px; color: #857462; font-size: 13px; line-height: 1.6; }
.mapping-card { border-radius: 14px; background: #faf6ee; padding: 14px; margin-bottom: 10px; }
.mapping-card-head { display: flex; align-items: center; gap: 8px; margin-bottom: 10px; }
.mapping-card-head strong { color: #21362d; font-size: 14px; }
.mapping-kind-tag { display: inline-flex; min-height: 18px; align-items: center; border-radius: 999px; padding: 0 8px; font-size: 11px; font-weight: 700; background: rgba(100, 116, 139, 0.10); color: #64748b; }
.mapping-card-body { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 12px; }
.mapping-field label { display: block; margin-bottom: 4px; color: #6f6254; font-size: 12px; font-weight: 700; }

.mapping-rules { margin-top: 12px; border-top: 1px dashed #e6dcc9; padding-top: 10px; }
.mapping-rules-head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 8px; color: #6f6254; font-size: 12px; font-weight: 700; }
.mapping-rules-empty { margin: 0 0 6px; color: #a8997f; font-size: 12px; }
.management-ghost-action.small { min-height: 30px; padding: 0 10px; font-size: 12px; }
.rule-row { display: grid; grid-template-columns: minmax(0, 1.2fr) 110px minmax(0, 1.2fr) 18px minmax(0, 1.4fr) 92px; gap: 8px; align-items: center; margin-bottom: 8px; }
.rule-arrow { text-align: center; color: #a8997f; font-weight: 800; }
.rule-ops { display: flex; gap: 4px; align-items: center; }
.rule-move { width: 26px; height: 26px; border: 0; border-radius: 8px; background: rgba(100, 116, 139, 0.10); color: #64748b; font-weight: 800; cursor: pointer; }
.rule-move:disabled { opacity: 0.4; cursor: not-allowed; }
.rule-move:not(:disabled):active { transform: scale(0.92); }
.rule-remove { width: 26px; height: 26px; border: 0; border-radius: 8px; background: rgba(239, 68, 68, 0.10); color: #ef4444; font-weight: 800; cursor: pointer; }
.rule-remove:active { transform: scale(0.92); }
.full-width { width: 100%; }

.filter-editor { margin-top: 18px; border-top: 1px solid #efe7d8; padding-top: 16px; }
.filter-editor-head { display: flex; align-items: flex-start; justify-content: space-between; gap: 12px; margin-bottom: 12px; }
.filter-editor-title { margin: 0; color: #21362d; font-size: 16px; }
.filter-editor-hint { margin: 6px 0 0; color: #857462; font-size: 13px; line-height: 1.6; max-width: 600px; }
.filter-empty { margin: 0; color: #a8997f; font-size: 13px; }
.filter-row { display: grid; grid-template-columns: 100px auto minmax(0, 1.2fr) 110px minmax(0, 1.4fr) 92px; gap: 8px; align-items: center; margin-bottom: 8px; }
.filter-when { color: #857462; font-size: 13px; font-weight: 700; text-align: center; }
.filter-action :deep(.el-input__inner) { font-weight: 700; }

@media (hover: hover) {
  .import-row:hover { transform: translateY(-2px); box-shadow: 0 3px 8px rgba(15, 23, 42, .12), 0 16px 34px rgba(15, 23, 42, .06); }
}
@media (max-width: 760px) {
  .import-config-page { padding: 14px; }
  .import-form-grid, .mapping-card-body { grid-template-columns: 1fr; }
  .import-row { grid-template-columns: 48px minmax(0, 1fr); }
  .import-actions { grid-column: 1 / -1; justify-content: flex-end; }
  .rule-row { grid-template-columns: 1fr 1fr; }
  .rule-arrow { display: none; }
}
</style>
