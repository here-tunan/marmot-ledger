<template>
  <div class="template-management-page">
    <div class="page-hero reveal-block">
      <div>
        <p class="eyebrow">{{ t('templateManagement.hero.eyebrow') }}</p>
        <h1>{{ t('templateManagement.hero.title') }}</h1>
        <p>{{ t('templateManagement.hero.subtitle') }}</p>
      </div>
    </div>

    <div class="content-card reveal-block delay-1">
      <div class="content-tabs">
        <button :class="['tab-btn', { active: activeTab === 'account' }]" @click="activeTab = 'account'">
          {{ t('templateManagement.tabs.account') }}
        </button>
        <button :class="['tab-btn', { active: activeTab === 'channel' }]" @click="activeTab = 'channel'">
          {{ t('templateManagement.tabs.channel') }}
        </button>
        <button :class="['tab-btn', { active: activeTab === 'category' }]" @click="activeTab = 'category'">
          {{ t('templateManagement.tabs.category') }}
        </button>
      </div>

      <!-- Account templates -->
      <div v-if="activeTab === 'account'" class="tab-content">
        <div class="tab-actions">
          <button class="management-primary-action" @click="openAccountDialog()">{{ t('templateManagement.actions.new') }}</button>
        </div>
        <div class="template-grid">
          <div v-for="item in accountTemplates" :key="item.id" class="template-card">
            <div class="template-icon account-template-icon" :style="accountTemplateIconStyle(item)">
              <el-icon><component :is="resolveAccountTemplateIcon(item)" /></el-icon>
            </div>
            <div class="template-info">
              <h3>{{ item.name }}</h3>
              <p>{{ item.providerCode }} · {{ accountTypeLabel(item.type) }}</p>
            </div>
            <div class="template-status">
              <span :class="['management-status-tag', { active: item.enabled }]">
                {{ item.enabled ? t('common.status.enabled') : t('common.status.disabled') }}
              </span>
            </div>
            <div class="template-actions">
              <button class="management-text-action" @click="openAccountDialog(item)">{{ t('common.actions.edit') }}</button>
            </div>
          </div>
        </div>
      </div>

      <!-- Channel templates -->
      <div v-if="activeTab === 'channel'" class="tab-content">
        <div class="tab-actions">
          <button class="management-primary-action" @click="openChannelDialog()">{{ t('templateManagement.actions.new') }}</button>
        </div>
        <div class="template-grid">
          <div v-for="item in channelTemplates" :key="item.id" class="template-card">
            <div class="template-icon channel-template-icon">
              {{ item.icon || '🔗' }}
            </div>
            <div class="template-info">
              <h3>{{ item.name }}</h3>
              <p>{{ item.channelCode }} · {{ item.providerCode }}</p>
            </div>
            <div class="template-status">
              <span :class="['management-status-tag', { active: item.enabled }]">
                {{ item.enabled ? t('common.status.enabled') : t('common.status.disabled') }}
              </span>
            </div>
            <div class="template-actions">
              <button class="management-text-action" @click="openChannelDialog(item)">{{ t('common.actions.edit') }}</button>
            </div>
          </div>
        </div>
      </div>

      <!-- Category templates -->
      <div v-if="activeTab === 'category'" class="tab-content">
        <div class="tab-actions">
          <button class="management-primary-action" @click="openCategoryDialog()">{{ t('templateManagement.actions.new') }}</button>
        </div>
        <div class="template-grid">
          <div v-for="item in categoryTemplates" :key="item.id" class="template-card">
            <div class="template-icon" :style="{ background: item.color || (item.type === 'income' ? '#ef4444' : '#f97316') }">
              {{ item.icon || (item.type === 'income' ? '💰' : '💸') }}
            </div>
            <div class="template-info">
              <h3>{{ item.name }}</h3>
              <p>{{ item.templateCode }} · {{ item.type === 'income' ? t('domain.income') : t('domain.expense') }}</p>
            </div>
            <div class="template-status">
              <span :class="['management-status-tag', { active: item.enabled }]">
                {{ item.enabled ? t('common.status.enabled') : t('common.status.disabled') }}
              </span>
            </div>
            <div class="template-actions">
              <button class="management-text-action" @click="openCategoryDialog(item)">{{ t('common.actions.edit') }}</button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Account template dialog -->
    <el-dialog v-model="accountDialog.visible" :title="accountDialog.data?.id ? t('templateManagement.dialog.account.editTitle') : t('templateManagement.dialog.account.createTitle')" width="560px" class="marmot-dialog template-dialog">
      <el-form :model="accountDialog.data" label-position="top" class="template-form">
        <el-form-item :label="t('templateManagement.fields.providerCode')">
          <el-input v-model="accountDialog.data.providerCode" :disabled="!!accountDialog.data?.id" />
        </el-form-item>
        <el-form-item :label="t('templateManagement.fields.name')">
          <el-input v-model="accountDialog.data.name" />
        </el-form-item>
        <el-form-item :label="t('templateManagement.fields.type')">
          <el-select v-model="accountDialog.data.type" style="width: 100%">
            <el-option v-for="opt in accountTypeOptions" :key="opt.value" :label="opt.label.value || opt.label" :value="opt.value" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('templateManagement.fields.icon')">
          <el-select v-model="accountDialog.data.icon" style="width: 100%">
            <el-option v-for="item in accountIconOptions" :key="item.value" :label="item.label.value || item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('templateManagement.fields.color')">
          <div class="admin-color-row">
            <button v-for="color in accountColors" :key="color" type="button" :class="['admin-color-dot', { active: accountDialog.data.color === color }]" :style="{ background: color }" @click="accountDialog.data.color = color"></button>
            <el-input v-model="accountDialog.data.color" style="width: 130px" />
          </div>
        </el-form-item>
        <el-form-item :label="t('templateManagement.fields.sort')">
          <el-input-number v-model="accountDialog.data.sort" :min="0" />
        </el-form-item>
        <StatusSwitchField v-if="accountDialog.data?.id" v-model="accountDialog.data.enabled" />
      </el-form>
      <template #footer>
        <ManagementDialogFooter @cancel="accountDialog.visible = false" @submit="saveAccountTemplate" />
      </template>
    </el-dialog>

    <!-- Channel template dialog -->
    <el-dialog v-model="channelDialog.visible" :title="channelDialog.data?.id ? t('templateManagement.dialog.channel.editTitle') : t('templateManagement.dialog.channel.createTitle')" width="560px" class="marmot-dialog template-dialog">
      <el-form :model="channelDialog.data" label-position="top" class="template-form">
        <el-form-item :label="t('templateManagement.fields.channelCode')">
          <el-input v-model="channelDialog.data.channelCode" :disabled="!!channelDialog.data?.id" />
        </el-form-item>
        <el-form-item :label="t('templateManagement.fields.channelName')">
          <el-input v-model="channelDialog.data.name" />
        </el-form-item>
        <el-form-item :label="t('templateManagement.fields.channelType')">
          <el-input v-model="channelDialog.data.channelType" />
        </el-form-item>
        <el-form-item :label="t('templateManagement.fields.providerLabel')">
          <el-input v-model="channelDialog.data.providerCode" />
        </el-form-item>
        <el-form-item :label="t('templateManagement.fields.supportedEvents')">
          <el-input v-model="channelDialog.data.supportedEventTypes" :placeholder="t('templateManagement.placeholders.supportedEvents')" />
        </el-form-item>
        <el-form-item :label="t('templateManagement.fields.sort')">
          <el-input-number v-model="channelDialog.data.sort" :min="0" />
        </el-form-item>
        <StatusSwitchField v-if="channelDialog.data?.id" v-model="channelDialog.data.enabled" />
      </el-form>
      <template #footer>
        <ManagementDialogFooter @cancel="channelDialog.visible = false" @submit="saveChannelTemplate" />
      </template>
    </el-dialog>

    <!-- Category template dialog -->
    <el-dialog v-model="categoryDialog.visible" :title="categoryDialog.data?.id ? t('templateManagement.dialog.category.editTitle') : t('templateManagement.dialog.category.createTitle')" width="560px" class="marmot-dialog template-dialog">
      <el-form :model="categoryDialog.data" label-position="top" class="template-form">
        <el-form-item :label="t('templateManagement.fields.templateCode')">
          <el-input v-model="categoryDialog.data.templateCode" :disabled="!!categoryDialog.data?.id" />
        </el-form-item>
        <el-form-item :label="t('templateManagement.fields.categoryName')">
          <el-input v-model="categoryDialog.data.name" />
        </el-form-item>
        <el-form-item :label="t('templateManagement.fields.eventType')">
          <el-select v-model="categoryDialog.data.type" style="width: 100%">
            <el-option :label="t('domain.income')" value="income" />
            <el-option :label="t('domain.expense')" value="expense" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('templateManagement.fields.iconAndColor')">
          <IconColorPicker
            v-model:icon-value="categoryDialog.data.icon"
            v-model:color-value="categoryDialog.data.color"
            :icon-label="t('templateManagement.fields.icon')"
            :color-label="t('templateManagement.fields.color')"
          />
        </el-form-item>
        <el-form-item :label="t('templateManagement.fields.sort')">
          <el-input-number v-model="categoryDialog.data.sort" :min="0" />
        </el-form-item>
        <StatusSwitchField v-if="categoryDialog.data?.id" v-model="categoryDialog.data.enabled" />
      </el-form>
      <template #footer>
        <ManagementDialogFooter @cancel="categoryDialog.visible = false" @submit="saveCategoryTemplate" />
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { listAccountTemplatesAdmin, createAccountTemplate, updateAccountTemplate } from '@/api/accountTemplate'
import { listChannelTemplatesAdmin, createChannelTemplate, updateChannelTemplate } from '@/api/channelTemplate'
import { listCategoryTemplatesAdmin, createCategoryTemplate, updateCategoryTemplate } from '@/api/categoryTemplate'
import IconColorPicker from '@/components/IconColorPicker.vue'
import ManagementDialogFooter from '@/components/management/ManagementDialogFooter.vue'
import StatusSwitchField from '@/components/management/StatusSwitchField.vue'

const { t } = useI18n()
const activeTab = ref('account')
const accountTemplates = ref([])
const channelTemplates = ref([])
const categoryTemplates = ref([])

const loading = reactive({
  account: false,
  channel: false,
  category: false
})

const accountTypeOptions = [
  { value: 'cash', label: computed(() => t('templateManagement.accountTypes.cash')) },
  { value: 'wallet', label: computed(() => t('templateManagement.accountTypes.wallet')) },
  { value: 'bank', label: computed(() => t('templateManagement.accountTypes.bank')) },
  { value: 'credit', label: computed(() => t('templateManagement.accountTypes.credit')) },
  { value: 'investment', label: computed(() => t('templateManagement.accountTypes.investment')) },
  { value: 'liability', label: computed(() => t('templateManagement.accountTypes.liability')) },
  { value: 'other', label: computed(() => t('templateManagement.accountTypes.other')) },
]

const accountIconOptions = [
  { label: computed(() => t('templateManagement.iconOptions.wallet')), value: 'Wallet' },
  { label: computed(() => t('templateManagement.iconOptions.card')), value: 'CreditCard' },
  { label: computed(() => t('templateManagement.iconOptions.money')), value: 'Money' },
  { label: computed(() => t('templateManagement.iconOptions.investment')), value: 'TrendCharts' },
  { label: computed(() => t('templateManagement.iconOptions.office')), value: 'OfficeBuilding' },
  { label: computed(() => t('templateManagement.iconOptions.house')), value: 'House' },
  { label: computed(() => t('templateManagement.iconOptions.collection')), value: 'Collection' },
  { label: computed(() => t('templateManagement.iconOptions.other')), value: 'More' },
]
const accountColors = ['#2f7d5c', '#1f2933', '#f59e0b', '#ef4444', '#f97316', '#22c55e', '#06b6d4', '#8b5cf6', '#ec4899', '#1677ff']

const accountDialog = reactive({
  visible: false,
  data: null
})

const channelDialog = reactive({
  visible: false,
  data: null
})

const categoryDialog = reactive({
  visible: false,
  data: null
})

onMounted(() => {
  loadAccountTemplates()
  loadChannelTemplates()
  loadCategoryTemplates()
})

async function loadAccountTemplates() {
  loading.account = true
  try {
    const res = await listAccountTemplatesAdmin()
    if (res.success) {
      accountTemplates.value = res.data
    }
  } catch (e) {
    ElMessage.error(t('templateManagement.messages.loadAccountFailed'))
  } finally {
    loading.account = false
  }
}

async function loadChannelTemplates() {
  loading.channel = true
  try {
    const res = await listChannelTemplatesAdmin()
    if (res.success) {
      channelTemplates.value = res.data
    }
  } catch (e) {
    ElMessage.error(t('templateManagement.messages.loadChannelFailed'))
  } finally {
    loading.channel = false
  }
}

async function loadCategoryTemplates() {
  loading.category = true
  try {
    const res = await listCategoryTemplatesAdmin()
    if (res.success) {
      categoryTemplates.value = res.data
    }
  } catch (e) {
    ElMessage.error(t('templateManagement.messages.loadCategoryFailed'))
  } finally {
    loading.category = false
  }
}

function openAccountDialog(row = null) {
  accountDialog.data = row ? { ...row } : { providerCode: '', name: '', type: 'wallet', icon: 'Wallet', color: '#2f7d5c', sort: 0, enabled: true }
  accountDialog.data.type = String(accountDialog.data.type || 'wallet').toLowerCase()
  accountDialog.data.icon = accountDialog.data.icon || defaultAccountIcon(accountDialog.data.type)
  accountDialog.data.color = accountDialog.data.color || defaultAccountColor(accountDialog.data.type)
  accountDialog.visible = true
}

function openChannelDialog(row = null) {
  channelDialog.data = row ? { ...row } : { channelCode: '', name: '', channelType: '', providerCode: '', supportedEventTypes: '', sort: 0, enabled: true }
  channelDialog.visible = true
}

function openCategoryDialog(row = null) {
  categoryDialog.data = row ? { ...row } : { templateCode: '', name: '', type: 'expense', icon: '', color: '', sort: 0, enabled: true }
  categoryDialog.visible = true
}

function accountTypeLabel(type) {
  const key = String(type || '').toLowerCase()
  return t(`templateManagement.accountTypes.${key in tCheck ? key : 'other'}`)
}

const tCheck = { cash: 1, wallet: 1, bank: 1, credit: 1, investment: 1, liability: 1, other: 1 }

function defaultAccountIcon(type) {
  switch (String(type || '').toLowerCase()) {
    case 'cash': return 'Money'
    case 'bank':
    case 'credit': return 'CreditCard'
    case 'investment': return 'TrendCharts'
    case 'liability': return 'Warning'
    default: return 'Wallet'
  }
}

function defaultAccountColor(type) {
  switch (String(type || '').toLowerCase()) {
    case 'cash': return '#f59e0b'
    case 'bank': return '#2f7d5c'
    case 'credit': return '#ef4444'
    case 'investment': return '#1f2933'
    case 'liability': return '#f97316'
    case 'wallet': return '#22c55e'
    default: return '#2f7d5c'
  }
}

function resolveAccountTemplateIcon(item) {
  return item.icon || defaultAccountIcon(item.type)
}

function accountTemplateIconStyle(item) {
  const color = item.color || defaultAccountColor(item.type)
  return {
    color,
    background: `${color}18`,
    boxShadow: `inset 0 0 0 1px ${color}26`,
  }
}

async function saveAccountTemplate() {
  try {
    accountDialog.data.type = String(accountDialog.data.type || 'wallet').toLowerCase()
    accountDialog.data.icon = accountDialog.data.icon || defaultAccountIcon(accountDialog.data.type)
    accountDialog.data.color = accountDialog.data.color || defaultAccountColor(accountDialog.data.type)
    if (accountDialog.data.id) {
      await updateAccountTemplate(accountDialog.data.id, accountDialog.data)
      ElMessage.success(t('templateManagement.messages.updated'))
    } else {
      await createAccountTemplate(accountDialog.data)
      ElMessage.success(t('templateManagement.messages.created'))
    }
    accountDialog.visible = false
    loadAccountTemplates()
  } catch (e) {
    ElMessage.error(e.message || t('templateManagement.messages.saveFailed'))
  }
}

async function saveChannelTemplate() {
  try {
    if (channelDialog.data.id) {
      await updateChannelTemplate(channelDialog.data.id, channelDialog.data)
      ElMessage.success(t('templateManagement.messages.updated'))
    } else {
      await createChannelTemplate(channelDialog.data)
      ElMessage.success(t('templateManagement.messages.created'))
    }
    channelDialog.visible = false
    loadChannelTemplates()
  } catch (e) {
    ElMessage.error(e.message || t('templateManagement.messages.saveFailed'))
  }
}

async function saveCategoryTemplate() {
  try {
    if (categoryDialog.data.id) {
      await updateCategoryTemplate(categoryDialog.data.id, categoryDialog.data)
      ElMessage.success(t('templateManagement.messages.updated'))
    } else {
      await createCategoryTemplate(categoryDialog.data)
      ElMessage.success(t('templateManagement.messages.created'))
    }
    categoryDialog.visible = false
    loadCategoryTemplates()
  } catch (e) {
    ElMessage.error(e.message || t('templateManagement.messages.saveFailed'))
  }
}
</script>

<style scoped>
.template-management-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 0 24px 0;
}

.reveal-block {
  animation: revealUp 480ms cubic-bezier(.16, 1, .3, 1) both;
}

.delay-1 {
  animation-delay: 90ms;
}

.page-hero,
.content-card {
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 1px 3px rgba(15, 23, 42, .1), 0 12px 30px rgba(15, 23, 42, .04);
}

.page-hero {
  padding: 26px;
  margin-bottom: 18px;
  background: linear-gradient(135deg, #f0fdf4 0%, #fff 70%);
}

.eyebrow {
  margin: 0 0 8px;
  color: #2f7d5c;
  font-size: 12px;
  font-weight: 800;
  letter-spacing: .08em;
  text-transform: uppercase;
}

.page-hero h1 {
  margin: 0;
  font-size: 30px;
  line-height: 1.16;
  letter-spacing: -.022em;
  color: #1e293b;
}

.page-hero p {
  color: #64748b;
  line-height: 1.7;
  margin: 12px 0 0;
}

.content-card {
  padding: 22px;
}

.template-dialog :deep(.el-dialog__body) {
  padding: 18px 24px 8px;
}

.template-dialog :deep(.el-dialog__footer) {
  padding: 12px 24px 22px;
}

.template-form :deep(.el-form-item) {
  margin-bottom: 16px;
}

.template-form :deep(.el-form-item__label) {
  padding: 0 0 6px;
  line-height: 1.4;
  font-size: 13px;
  font-weight: 700;
  color: #334155;
  white-space: normal;
}

.template-form :deep(.el-form-item__content) {
  line-height: 1.4;
}

.template-form :deep(.el-input__wrapper),
.template-form :deep(.el-select),
.template-form :deep(.el-input-number) {
  width: 100%;
}

.content-tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 24px;
  border-bottom: 1px solid #e2e8f0;
  padding-bottom: 16px;
}

.tab-btn {
  padding: 8px 20px;
  border-radius: 10px;
  border: none;
  background: transparent;
  color: #64748b;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.tab-btn:hover {
  background: #f8faf7;
  color: #475569;
}

.tab-btn.active {
  background: #2f7d5c;
  color: white;
  box-shadow: 0 4px 12px rgba(47, 125, 92, .2);
}

.tab-content {
  animation: fadeIn 200ms ease;
}

.tab-actions {
  margin-bottom: 20px;
  text-align: right;
}

.template-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
}

.template-card {
  display: grid;
  grid-template-columns: 48px 1fr auto auto;
  gap: 12px;
  align-items: center;
  padding: 14px;
  background: #f8faf7;
  border-radius: 14px;
  transition: all 0.2s ease;
}

.template-card:hover {
  box-shadow: 0 2px 8px rgba(15, 23, 42, .08);
}

.template-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
}

.account-template-icon .el-icon {
  font-size: 22px;
}

.admin-color-row {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
}

.admin-color-dot {
  width: 28px;
  height: 28px;
  border: 0;
  border-radius: 999px;
  cursor: pointer;
  box-shadow: inset 0 0 0 2px rgba(255, 255, 255, 0.9), 0 0 0 1px rgba(15, 23, 42, 0.12);
}

.admin-color-dot.active {
  box-shadow: inset 0 0 0 2px rgba(255, 255, 255, 0.9), 0 0 0 3px rgba(47, 125, 92, 0.26);
}

.template-info h3 {
  margin: 0 0 4px 0;
  font-size: 15px;
  font-weight: 600;
  color: #1e293b;
}

.template-info p {
  margin: 0;
  font-size: 13px;
  color: #64748b;
}

.primary-action,
.ghost-action {
  min-height: 40px;
  border: 0;
  border-radius: 12px;
  padding: 0 16px;
  font-weight: 800;
  cursor: pointer;
  transition: all 0.2s ease;
}

.primary-action {
  background: #2f7d5c;
  color: #fff;
  box-shadow: 0 10px 24px rgba(59, 130, 246, .22);
}

.primary-action:hover {
  background: #256f53;
  transform: translateY(-1px);
}

.ghost-action {
  background: #ffffff;
  color: #334155;
  box-shadow: inset 0 0 0 1px rgba(100, 116, 139, .16);
}

.text-action {
  min-height: 32px;
  border: 0;
  border-radius: 10px;
  padding: 0 12px;
  background: #f4efe6;
  color: #6b5b49;
  font-size: 13px;
  font-weight: 700;
  cursor: pointer;
  transition-property: transform, color;
  transition-duration: 160ms;
  transition-timing-function: cubic-bezier(0.16, 1, 0.3, 1);
}

.text-action:active {
  transform: scale(0.96);
}

.ghost-action:hover {
  background: #f0f4f0;
}

.tiny {
  min-height: 28px;
  padding: 0 10px;
  font-size: 12px;
}

@media (hover: hover) {
  .text-action:hover {
    background: #ece2d2;
    color: #4b3f33;
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

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(4px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (max-width: 900px) {
  .template-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 640px) {
  .template-grid {
    grid-template-columns: 1fr;
  }
}
</style>
