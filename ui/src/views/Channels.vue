<template>
  <main class="channel-page">
    <section class="channel-workbench reveal-block">
      <header class="channel-head">
        <div>
          <p class="eyebrow">{{ t('channels.hero.eyebrow') }}</p>
          <h1>{{ t('channels.hero.title') }}</h1>
          <p>{{ t('channels.hero.subtitle') }}</p>
        </div>
        <div class="channel-head-actions">
          <span class="channel-count-pill">{{ t('channels.summary.total', { count: allChannels.length }) }}</span>
          <button class="primary-action" type="button" @click="openCreate">{{ t('channels.actions.new') }}</button>
        </div>
      </header>

      <div class="channel-board">
        <aside class="channel-rail" :aria-label="t('channels.filters.typePlaceholder')">
          <button type="button" class="rail-item" :class="{ active: !filters.channelType }" @click="setTypeFilter('')">
            <span class="rail-icon">✦</span>
            <span>{{ t('channels.filters.allTypes') }}</span>
            <em>{{ allChannels.length }}</em>
          </button>
          <button v-for="item in channelTypes" :key="item.value" type="button" class="rail-item" :class="{ active: filters.channelType === item.value }" @click="setTypeFilter(item.value)">
            <span class="rail-icon">{{ channelTypeIcon(item.value) }}</span>
            <span>{{ item.label.value || item.label }}</span>
            <em>{{ typeCounts[item.value] || 0 }}</em>
          </button>
        </aside>

        <section class="channel-main-panel" v-loading="loading">
          <div class="channel-toolbelt">
            <div>
              <strong>{{ activeTypeTitle }}</strong>
              <span>{{ t('channels.summary.active', { count: activeChannels.length }) }}</span>
            </div>
            <div class="channel-filters">
              <el-select v-model="filters.eventType" clearable :placeholder="t('channels.filters.eventTypePlaceholder')" class="channel-filter" @change="loadChannels">
                <el-option v-for="item in eventTypes" :key="item.value" :label="item.label.value || item.label" :value="item.value" />
              </el-select>
              <el-select v-model="filters.isActive" clearable :placeholder="t('common.filters.enabledStatus')" class="channel-filter" @change="loadChannels">
                <el-option :label="t('common.status.enabled')" :value="true" />
                <el-option :label="t('common.status.disabled')" :value="false" />
              </el-select>
              <button class="ghost-action" type="button" @click="loadChannels">{{ t('common.actions.refresh') }}</button>
            </div>
          </div>

          <div v-if="channels.length" class="channel-list">
            <article v-for="item in channels" :key="item.id" class="channel-row" :class="{ inactive: !item.isActive }">
              <div class="channel-symbol" :style="channelSymbolStyle(item)">{{ item.icon || channelTypeIcon(item.channelType) }}</div>
              <div class="channel-copy">
                <div class="channel-title-line">
                  <h2>{{ item.name }}</h2>
                  <span :class="['status-tag', { active: item.isActive }]">{{ item.isActive ? t('common.status.enabled') : t('common.status.disabled') }}</span>
                </div>
                <div class="channel-meta-line">
                  <span>{{ typeLabel(item.channelType) }}</span>
                  <span>{{ item.providerCode || t('channels.summary.noProvider') }}</span>
                  <span>{{ t('channels.fields.sort') }} {{ item.sort || 0 }}</span>
                </div>
                <div class="event-pill-row">
                  <span v-for="event in splitEvents(item.supportedEventTypes)" :key="`${item.id}-${event}`" class="event-pill">{{ eventLabel(event) }}</span>
                </div>
              </div>
              <div class="channel-actions">
                <button class="text-action" type="button" @click="openEdit(item)">{{ t('common.actions.edit') }}</button>
                <button class="danger-action" type="button" @click="handleDelete(item)">{{ t('common.actions.delete') }}</button>
              </div>
            </article>
          </div>

          <div v-else-if="!loading" class="channel-empty">
            <img :src="marmotOne" :alt="t('channels.empty.alt')" width="112" height="112" />
            <h2>{{ t('channels.empty.title') }}</h2>
            <p>{{ t('channels.empty.text') }}</p>
            <button class="primary-action" type="button" @click="openCreate">{{ t('channels.actions.new') }}</button>
          </div>
        </section>
      </div>
    </section>

    <el-dialog v-model="dialogVisible" :title="editingId ? t('channels.dialog.editTitle') : t('channels.dialog.createTitle')" width="680px" class="marmot-dialog channel-dialog">
      <div v-if="!editingId" class="template-shelf">
        <button class="shelf-head" type="button" @click="templateCollapsed = !templateCollapsed">
          <span>⚡</span>
          <strong>{{ t('channels.templates.title') }}</strong>
          <em>{{ templateCollapsed ? t('channels.templates.expand') : t('channels.templates.collapse') }}</em>
        </button>
        <div v-if="!templateCollapsed" class="template-shelf-grid">
          <p v-if="!templates.length" class="template-empty">{{ t('channels.templates.empty') }}</p>
          <button v-for="tpl in templates" v-else :key="tpl.id" type="button" class="template-chip" :class="{ active: selectedTemplate?.id === tpl.id }" @click="selectTemplate(tpl)">
            <span>{{ tpl.icon || channelTypeIcon(tpl.channelType) }}</span>
            <strong>{{ tpl.name }}</strong>
          </button>
        </div>
      </div>

      <el-form :model="form" label-position="top" class="channel-form-grid">
        <el-form-item :label="t('channels.fields.name')"><el-input v-model="form.name" /></el-form-item>
        <el-form-item :label="t('channels.fields.type')"><el-select v-model="form.channelType" class="full-width"><el-option v-for="item in channelTypes" :key="item.value" :label="item.label.value || item.label" :value="item.value" /></el-select></el-form-item>
        <el-form-item :label="t('channels.fields.providerCode')"><el-input v-model="form.providerCode" /></el-form-item>
        <el-form-item :label="t('channels.fields.icon')"><el-input v-model="form.icon" placeholder="💬" /></el-form-item>
        <el-form-item class="wide" :label="t('channels.fields.supportedEvents')"><el-select v-model="supportedEvents" multiple class="full-width"><el-option v-for="item in eventTypes" :key="item.value" :label="item.label.value || item.label" :value="item.value" /></el-select></el-form-item>
        <el-form-item :label="t('channels.fields.sort')"><el-input-number v-model="form.sort" :min="0" /></el-form-item>
        <el-form-item v-if="editingId" :label="t('common.status.status')"><el-switch v-model="form.isActive" :active-text="t('common.status.enabled')" :inactive-text="t('common.status.disabled')" /></el-form-item>
      </el-form>
      <template #footer>
        <button class="ghost-action" type="button" @click="dialogVisible = false">{{ t('common.actions.cancel') }}</button>
        <button class="primary-action" type="button" @click="submitForm">{{ t('common.actions.save') }}</button>
      </template>
    </el-dialog>
  </main>
</template>

<script setup>
import { computed, onActivated, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { createChannel, deleteChannel, listChannels, updateChannel, checkChannelUsage } from '@/api/channel/channel'
import { listChannelTemplates } from '@/api/channelTemplate'
import marmotOne from '../../../img/marmot-ledger-1.png'

const { t } = useI18n()
const channels = ref([])
const allChannels = ref([])
const templates = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const editingId = ref(0)
const selectedTemplate = ref(null)
const templateCollapsed = ref(false)
const supportedEvents = ref(['income', 'expense'])
const filters = reactive({ channelType: '', eventType: '', isActive: '' })
const form = reactive(createEmptyForm())

const channelTypes = [
  { value: 'wallet', label: computed(() => t('channels.types.wallet')) },
  { value: 'bank', label: computed(() => t('channels.types.bank')) },
  { value: 'card', label: computed(() => t('channels.types.card')) },
  { value: 'cash', label: computed(() => t('channels.types.cash')) },
  { value: 'system', label: computed(() => t('channels.types.system')) },
  { value: 'other', label: computed(() => t('channels.types.other')) },
]
const eventTypes = ['income', 'expense', 'refund', 'transfer'].map((value) => ({ value, label: computed(() => t(`record.scenarios.${value}`)) }))
const activeChannels = computed(() => channels.value.filter((item) => item.isActive))
const typeCounts = computed(() => allChannels.value.reduce((acc, item) => {
  const key = item.channelType || 'other'
  acc[key] = (acc[key] || 0) + 1
  return acc
}, {}))
const activeTypeTitle = computed(() => filters.channelType ? typeLabel(filters.channelType) : t('channels.filters.allTypes'))

function createEmptyForm() {
  return { name: '', channelType: 'wallet', providerCode: '', supportedEventTypes: 'income,expense', icon: '🔗', sort: 0, isActive: true }
}
function resetForm(data = createEmptyForm()) {
  Object.assign(form, createEmptyForm(), data)
  supportedEvents.value = String(form.supportedEventTypes || '').split(',').filter(Boolean)
}
function typeLabel(type) { return channelTypes.find((item) => item.value === type)?.label.value || type }
function splitEvents(value) { return String(value || '').split(',').filter(Boolean) }
function eventLabel(value) { return t(`record.scenarios.${value}`) }
function channelTypeIcon(type) {
  return { wallet: '👛', bank: '🏦', card: '💳', cash: '💵', system: '🔁', other: '🔗' }[type] || '🔗'
}
function channelColor(type) {
  return { wallet: '#2f7d5c', bank: '#3b82f6', card: '#8b5cf6', cash: '#d97706', system: '#64748b', other: '#78716c' }[type] || '#2f7d5c'
}
function channelSymbolStyle(item) {
  const color = channelColor(item.channelType)
  return { background: `${color}18`, color, boxShadow: `inset 0 0 0 1px ${color}24` }
}
function setTypeFilter(type) {
  filters.channelType = type
  loadChannels()
}
function selectTemplate(tpl) {
  selectedTemplate.value = tpl
  resetForm({ name: tpl.name, channelType: tpl.channelType || 'wallet', providerCode: tpl.providerCode || '', supportedEventTypes: tpl.supportedEventTypes || 'income,expense', icon: tpl.icon || channelTypeIcon(tpl.channelType), sort: tpl.sort || 0, isActive: true })
}
async function loadTemplates() {
  const res = await listChannelTemplates()
  if (res.success) templates.value = res.data || []
}
async function loadChannels() {
  loading.value = true
  try {
    const [allRes, filteredRes] = await Promise.all([
      listChannels(),
      listChannels(channelFilterParams()),
    ])
    if (allRes.success) allChannels.value = allRes.data || []
    if (filteredRes.success) channels.value = filteredRes.data || []
    else ElMessage.error(filteredRes.error || t('channels.messages.loadFailed'))
  } finally { loading.value = false }
}

function channelFilterParams() {
  const params = {}
  if (filters.channelType) params.channelType = filters.channelType
  if (filters.eventType) params.eventType = filters.eventType
  if (filters.isActive !== '') params.isActive = filters.isActive
  return params
}
async function refreshAll() { await Promise.all([loadChannels(), loadTemplates()]) }
function openCreate() {
  editingId.value = 0
  selectedTemplate.value = null
  resetForm()
  dialogVisible.value = true
}
function openEdit(item) {
  editingId.value = item.id
  selectedTemplate.value = null
  resetForm({ ...item })
  dialogVisible.value = true
}
async function submitForm() {
  if (!form.name) return ElMessage.warning(t('channels.validation.nameRequired'))
  if (!form.channelType) return ElMessage.warning(t('channels.validation.typeRequired'))
  const payload = { ...form, supportedEventTypes: supportedEvents.value.join(','), isActive: form.isActive !== false }
  const res = editingId.value ? await updateChannel(editingId.value, payload) : await createChannel(payload)
  if (res.success) {
    ElMessage.success(editingId.value ? t('channels.messages.updated') : t('channels.messages.created'))
    dialogVisible.value = false
    await loadChannels()
  } else {
    ElMessage.error(res.error || t('channels.messages.saveFailed'))
  }
}
async function handleDelete(item) {
  try {
    const usageRes = await checkChannelUsage(item.id)
    const eventCount = usageRes.success ? usageRes.data.eventCount : 0
    const message = eventCount > 0 ? t('channels.delete.confirmUsed', { name: item.name, count: eventCount }) : t('channels.delete.confirm', { name: item.name })
    await ElMessageBox.confirm(message, t('common.actions.delete'), { confirmButtonText: t('common.actions.delete'), cancelButtonText: t('common.actions.cancel'), type: 'warning' })
    const res = await deleteChannel(item.id)
    if (res.success) {
      ElMessage.success(t('channels.messages.deleted'))
      await loadChannels()
    } else ElMessage.error(res.error || t('channels.messages.deleteFailed'))
  } catch (err) {
    if (err !== 'cancel') console.warn(err)
  }
}

onMounted(refreshAll)
onActivated(refreshAll)
</script>

<style scoped>
.channel-page {
  padding: 24px;
  min-height: 100%;
  background: #f6f0e6;
}

.channel-workbench {
  max-width: 1180px;
  margin: 0 auto;
}

.channel-head {
  display: flex;
  justify-content: space-between;
  gap: 24px;
  align-items: flex-end;
  margin-bottom: 18px;
}

.channel-head h1 {
  margin: 0;
  color: #21362d;
  font-size: clamp(28px, 4vw, 42px);
  line-height: 1;
  letter-spacing: -0.022em;
  text-wrap: balance;
}

.channel-head p:not(.eyebrow) {
  max-width: 620px;
  margin: 10px 0 0;
  color: #6f6254;
  text-wrap: pretty;
}

.channel-head-actions {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
  justify-content: flex-end;
}

.channel-count-pill {
  min-height: 36px;
  display: inline-flex;
  align-items: center;
  padding: 0 12px;
  border-radius: 10px;
  background: #ffffff;
  color: #64748b;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.10);
  font-size: 13px;
  font-weight: 700;
}

.channel-board {
  display: grid;
  grid-template-columns: 220px minmax(0, 1fr);
  gap: 18px;
  align-items: start;
}

.channel-rail,
.channel-main-panel {
  border-radius: 24px;
  background: #fffaf1;
  box-shadow: 0 18px 48px rgba(90, 68, 39, 0.11), 0 1px 3px rgba(68, 52, 32, 0.10);
}

.channel-rail {
  padding: 10px;
  position: sticky;
  top: 18px;
}

.rail-item {
  width: 100%;
  min-height: 46px;
  border: 0;
  border-radius: 16px;
  background: transparent;
  display: grid;
  grid-template-columns: 32px 1fr auto;
  gap: 8px;
  align-items: center;
  color: #6f6254;
  font-weight: 800;
  text-align: left;
  cursor: pointer;
  transition: transform 160ms cubic-bezier(0.16, 1, 0.3, 1), background-color 160ms ease, color 160ms ease;
  touch-action: manipulation;
}

.rail-item + .rail-item { margin-top: 4px; }
.rail-item:active { transform: scale(0.96); }
.rail-item.active { background: #e6f1ea; color: #245f48; }
.rail-icon { font-size: 18px; text-align: center; }
.rail-item em { font-style: normal; color: #9a8a76; font-variant-numeric: tabular-nums; }
.rail-item.active em { color: #2f7d5c; }

.channel-main-panel { padding: 16px; min-height: 520px; }

.channel-toolbelt {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 8px 4px 16px;
}

.channel-toolbelt strong { display: block; color: #21362d; font-size: 20px; letter-spacing: -0.012em; }
.channel-toolbelt span { color: #857462; font-size: 13px; }
.channel-filters { display: flex; gap: 10px; align-items: center; flex-wrap: wrap; justify-content: flex-end; }
.channel-filter { width: 170px; }

.channel-list { display: grid; gap: 10px; }
.channel-row {
  display: grid;
  grid-template-columns: 58px minmax(0, 1fr) auto;
  gap: 14px;
  align-items: center;
  padding: 14px;
  border-radius: 20px;
  background: #fffdf8;
  box-shadow: 0 1px 3px rgba(68, 52, 32, 0.10);
}
.channel-row.inactive { opacity: 0.64; }
.channel-symbol {
  width: 52px;
  height: 52px;
  border-radius: 18px;
  display: grid;
  place-items: center;
  font-size: 24px;
}
.channel-copy { min-width: 0; }
.channel-title-line { display: flex; gap: 10px; align-items: center; flex-wrap: wrap; }
.channel-title-line h2 { margin: 0; color: #21362d; font-size: 18px; letter-spacing: -0.012em; }
.channel-meta-line { display: flex; gap: 8px; flex-wrap: wrap; margin-top: 5px; color: #8a7a67; font-size: 13px; }
.channel-meta-line span:not(:last-child)::after { content: '·'; margin-left: 8px; color: #c4b49f; }
.event-pill-row { display: flex; gap: 6px; flex-wrap: wrap; margin-top: 8px; }
.event-pill { padding: 4px 8px; border-radius: 999px; background: #f3eadc; color: #6f6254; font-size: 12px; font-weight: 800; }
.channel-actions {
  display: flex;
  gap: 6px;
  align-items: center;
}

.channel-empty {
  min-height: 360px;
  display: grid;
  place-items: center;
  text-align: center;
  color: #6f6254;
}
.channel-empty h2 { margin: 10px 0 4px; color: #21362d; }
.channel-empty p { max-width: 360px; margin: 0 0 16px; }

.primary-action,
.ghost-action,
.text-action,
.danger-action,
.template-chip,
.shelf-head {
  min-height: 36px;
  border-radius: 10px;
  font-weight: 700;
  cursor: pointer;
  transition-property: transform, box-shadow, background-color, color;
  transition-duration: 160ms;
  transition-timing-function: cubic-bezier(0.16, 1, 0.3, 1);
  touch-action: manipulation;
}
.primary-action:active,
.ghost-action:active,
.text-action:active,
.danger-action:active,
.template-chip:active,
.shelf-head:active { transform: scale(0.96); }

.primary-action {
  min-height: 44px;
  border: 0;
  padding: 0 18px;
  background: #3b82f6;
  color: #ffffff;
  box-shadow: 0 8px 20px rgba(59, 130, 246, 0.18);
}

.ghost-action {
  border: 1px solid #e5e7eb;
  padding: 0 14px;
  background: #ffffff;
  color: #1e293b;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.08);
}

.text-action,
.danger-action {
  min-height: 32px;
  border: 0;
  padding: 0 12px;
  font-size: 13px;
}

.text-action {
  background: #f8faf7;
  color: #1e293b;
  box-shadow: none;
}

.danger-action {
  background: rgba(239, 68, 68, 0.10);
  color: #ef4444;
  box-shadow: none;
}

.template-shelf {
  border-radius: 20px;
  background: #fff8ec;
  padding: 12px;
  margin-bottom: 14px;
  box-shadow: inset 0 0 0 1px rgba(120, 92, 56, 0.08);
}
.shelf-head {
  width: 100%;
  min-height: 40px;
  border: 0;
  border-radius: 10px;
  background: transparent;
  display: flex;
  align-items: center;
  justify-content: space-between;
  color: #21362d;
  cursor: pointer;
}
.shelf-head em { font-style: normal; color: #857462; font-size: 13px; }
.template-shelf-grid { display: flex; flex-wrap: wrap; gap: 8px; padding-top: 10px; }
.template-chip {
  min-height: 36px;
  border: 1px solid #e5e7eb;
  border-radius: 10px;
  padding: 0 12px;
  background: #ffffff;
  color: #4b3f33;
  display: inline-flex;
  align-items: center;
  gap: 7px;
  cursor: pointer;
}
.template-chip.active {
  border-color: rgba(47, 125, 92, 0.28);
  background: #dce9df;
  color: #245f48;
}
.template-empty { margin: 0; color: #857462; }
.channel-form-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); column-gap: 14px; }
.channel-form-grid .wide { grid-column: 1 / -1; }

@media (hover: hover) {
  .rail-item:hover { background: #f4eddf; }
  .rail-item.active:hover { background: #e0ecdf; }
  .channel-row:hover { box-shadow: 0 10px 24px rgba(90, 68, 39, 0.10), 0 1px 3px rgba(68, 52, 32, 0.10); }
}

@media (max-width: 760px) {
  .channel-page { padding: 14px; }
  .channel-head { align-items: flex-start; flex-direction: column; }
  .channel-head-actions { justify-content: flex-start; }
  .channel-board { grid-template-columns: 1fr; }
  .channel-rail { position: static; display: flex; overflow-x: auto; gap: 8px; padding: 8px; }
  .rail-item { min-width: 152px; }
  .channel-toolbelt { align-items: stretch; flex-direction: column; }
  .channel-filters { justify-content: stretch; }
  .channel-filter { width: 100%; }
  .channel-row { grid-template-columns: 48px minmax(0, 1fr); }
  .channel-symbol { width: 44px; height: 44px; border-radius: 14px; }
  .channel-actions { grid-column: 1 / -1; justify-content: flex-end; }
  .channel-actions .text-action,
  .channel-actions .danger-action { min-height: 40px; }
  .channel-form-grid { grid-template-columns: 1fr; }
}
</style>
