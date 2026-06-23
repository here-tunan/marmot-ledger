<template>
  <main class="channel-page">
    <section class="channel-workbench reveal-block">
      <ManagementPageHeader :eyebrow="t('channels.hero.eyebrow')" :title="t('channels.hero.title')" :subtitle="t('channels.hero.subtitle')">
        <template #actions>
          <span class="channel-count-pill">{{ t('channels.summary.total', { count: allChannels.length }) }}</span>
          <button class="management-primary-action" type="button" @click="openCreate">{{ t('channels.actions.new') }}</button>
        </template>
      </ManagementPageHeader>

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
              <button class="management-ghost-action" type="button" @click="loadChannels">{{ t('common.actions.refresh') }}</button>
            </div>
          </div>

          <div v-if="channels.length" class="channel-list">
            <article v-for="item in channels" :key="item.id" class="channel-row" :class="{ inactive: !item.isActive }">
              <div class="channel-symbol" :style="channelSymbolStyle(item)">{{ item.icon || channelTypeIcon(item.channelType) }}</div>
              <div class="channel-copy">
                <div class="channel-title-line">
                  <h2>{{ item.name }}</h2>
                  <span :class="['management-status-tag', { active: item.isActive }]">{{ item.isActive ? t('common.status.enabled') : t('common.status.disabled') }}</span>
                </div>
                <div class="channel-meta-line">
                  <span class="channel-type-tag">{{ typeLabel(item.channelType) }}</span>
                  <span class="channel-meta-tag">{{ item.providerCode || t('channels.summary.noProvider') }}</span>
                  <span class="channel-meta-tag">{{ t('channels.fields.sort') }} {{ item.sort || 0 }}</span>
                </div>
                <div class="event-pill-row">
                  <span v-for="event in splitEvents(item.supportedEventTypes)" :key="`${item.id}-${event}`" class="channel-event-tag">{{ eventLabel(event) }}</span>
                </div>
              </div>
              <div class="channel-actions">
                <button class="management-text-action" type="button" @click="openEdit(item)">{{ t('common.actions.edit') }}</button>
                <button class="management-danger-action" type="button" @click="handleDelete(item)">{{ t('common.actions.delete') }}</button>
              </div>
            </article>
          </div>

          <ManagementEmptyState v-else-if="!loading" :image="marmotOne" :alt="t('channels.empty.alt')" :title="t('channels.empty.title')" :text="t('channels.empty.text')">
            <template #action>
              <button class="management-primary-action" type="button" @click="openCreate">{{ t('channels.actions.new') }}</button>
            </template>
          </ManagementEmptyState>
        </section>
      </div>
    </section>

    <el-dialog v-model="dialogVisible" :title="editingId ? t('channels.dialog.editTitle') : t('channels.dialog.createTitle')" width="680px" class="marmot-dialog channel-dialog">
      <TemplateQuickSelect
        v-if="!editingId"
        :title="t('channels.templates.title')"
        :items="templates"
        :active-key="selectedTemplate?.id"
        :empty-text="t('channels.templates.empty')"
        @select="selectTemplate"
      >
        <template #chip="{ item }">
          <span class="template-quick-chip-icon">{{ item.icon || channelTypeIcon(item.channelType) }}</span>
          <strong>{{ item.name }}</strong>
        </template>
      </TemplateQuickSelect>

      <el-form :model="form" label-position="top" class="channel-form-grid">
        <el-form-item :label="t('channels.fields.name')"><el-input v-model="form.name" /></el-form-item>
        <el-form-item :label="t('channels.fields.type')"><el-select v-model="form.channelType" class="full-width"><el-option v-for="item in channelTypes" :key="item.value" :label="item.label.value || item.label" :value="item.value" /></el-select></el-form-item>
        <el-form-item :label="t('channels.fields.providerCode')"><el-input v-model="form.providerCode" /></el-form-item>
        <el-form-item :label="t('channels.fields.icon')">
          <IconColorPicker v-model:icon-value="form.icon" icon-label="" :show-color="false" />
        </el-form-item>
        <el-form-item class="wide" :label="t('channels.fields.supportedEvents')"><el-select v-model="supportedEvents" multiple class="full-width"><el-option v-for="item in eventTypes" :key="item.value" :label="item.label.value || item.label" :value="item.value" /></el-select></el-form-item>
        <el-form-item :label="t('channels.fields.sort')"><el-input-number v-model="form.sort" :min="0" /></el-form-item>
        <StatusSwitchField v-if="editingId" v-model="form.isActive" />
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
import { ElMessage } from 'element-plus'
import { createChannel, deleteChannel, listChannels, updateChannel, checkChannelUsage } from '@/api/channel/channel'
import { listChannelTemplates } from '@/api/channelTemplate'
import ManagementPageHeader from '@/components/management/ManagementPageHeader.vue'
import ManagementEmptyState from '@/components/management/ManagementEmptyState.vue'
import ManagementDialogFooter from '@/components/management/ManagementDialogFooter.vue'
import StatusSwitchField from '@/components/management/StatusSwitchField.vue'
import TemplateQuickSelect from '@/components/management/TemplateQuickSelect.vue'
import IconColorPicker from '@/components/IconColorPicker.vue'
import { confirmDelete, isCancelError } from '@/components/management/confirmDelete'
import marmotOne from '../../../img/marmot-ledger-1.png'

const { t } = useI18n()
const channels = ref([])
const allChannels = ref([])
const templates = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const editingId = ref(0)
const selectedTemplate = ref(null)
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
    await confirmDelete({ message, title: t('common.actions.delete'), confirmText: t('common.actions.delete'), cancelText: t('common.actions.cancel') })
    const res = await deleteChannel(item.id)
    if (res.success) {
      ElMessage.success(t('channels.messages.deleted'))
      await loadChannels()
    } else ElMessage.error(res.error || t('channels.messages.deleteFailed'))
  } catch (err) {
    if (!isCancelError(err)) console.warn(err)
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
  border-radius: 16px;
  background: #fff;
  box-shadow: 0 1px 3px rgba(15, 23, 42, .1), 0 12px 30px rgba(15, 23, 42, .04);
}

.channel-rail {
  padding: 10px;
  position: sticky;
  top: 18px;
}

.rail-item {
  width: 100%;
  min-height: 40px;
  border: 0;
  border-radius: 10px;
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
.rail-item.active { background: #dce9df; color: #245f48; }
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
  border-radius: 16px;
  background: #fff;
  box-shadow: 0 1px 3px rgba(15, 23, 42, .1), 0 12px 30px rgba(15, 23, 42, .04);
  transition-property: transform, box-shadow;
  transition-duration: 180ms;
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
.channel-meta-line { display: flex; gap: 6px; flex-wrap: wrap; margin-top: 7px; color: #8a7a67; font-size: 13px; }
.event-pill-row { display: flex; gap: 6px; flex-wrap: wrap; margin-top: 8px; }

.channel-type-tag,
.channel-meta-tag,
.channel-event-tag {
  display: inline-flex;
  align-items: center;
  min-height: 20px;
  border-radius: 999px;
  padding: 0 8px;
  font-size: 11px;
  font-weight: 700;
  line-height: 1.4;
}

.channel-type-tag {
  background: rgba(47, 125, 92, 0.10);
  color: #2f7d5c;
}

.channel-meta-tag {
  background: #f7f3eb;
  color: #7c6c5a;
}

.channel-event-tag {
  background: rgba(100, 116, 139, 0.08);
  color: #64748b;
}
.channel-actions {
  display: flex;
  gap: 6px;
  align-items: center;
}

.channel-form-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); column-gap: 14px; }
.channel-form-grid .wide { grid-column: 1 / -1; }

@media (hover: hover) {
  .rail-item:hover { background: #f8faf7; }
  .rail-item.active:hover { background: #dce9df; }
  .channel-row:hover {
    transform: translateY(-2px);
    box-shadow: 0 3px 8px rgba(15, 23, 42, .12), 0 16px 34px rgba(15, 23, 42, .06);
  }
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
