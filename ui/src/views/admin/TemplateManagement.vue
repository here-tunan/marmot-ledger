<template>
  <div class="template-management-page">
    <div class="page-hero reveal-block">
      <div>
        <p class="eyebrow">SYSTEM TEMPLATES</p>
        <h1>模板管理</h1>
        <p>管理账户、渠道、分类系统模板</p>
      </div>
    </div>

    <div class="content-card reveal-block delay-1">
      <div class="content-tabs">
        <button :class="['tab-btn', { active: activeTab === 'account' }]" @click="activeTab = 'account'">
          账户模板
        </button>
        <button :class="['tab-btn', { active: activeTab === 'channel' }]" @click="activeTab = 'channel'">
          渠道模板
        </button>
        <button :class="['tab-btn', { active: activeTab === 'category' }]" @click="activeTab = 'category'">
          分类模板
        </button>
      </div>

      <!-- 账户模板 -->
      <div v-if="activeTab === 'account'" class="tab-content">
        <div class="tab-actions">
          <button class="primary-action" @click="openAccountDialog()">新增模板</button>
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
                {{ item.enabled ? '启用' : '禁用' }}
              </span>
            </div>
            <div class="template-actions">
              <button class="text-action tiny" @click="openAccountDialog(item)">编辑</button>
            </div>
          </div>
        </div>
      </div>

      <!-- 渠道模板 -->
      <div v-if="activeTab === 'channel'" class="tab-content">
        <div class="tab-actions">
          <button class="primary-action" @click="openChannelDialog()">新增模板</button>
        </div>
        <div class="template-grid">
          <div v-for="item in channelTemplates" :key="item.id" class="template-card">
            <div class="template-icon" style="background: #8b5cf6">
              🔗
            </div>
            <div class="template-info">
              <h3>{{ item.name }}</h3>
              <p>{{ item.channelCode }} · {{ item.providerCode }}</p>
            </div>
            <div class="template-status">
              <span :class="['management-status-tag', { active: item.enabled }]">
                {{ item.enabled ? '启用' : '禁用' }}
              </span>
            </div>
            <div class="template-actions">
              <button class="text-action tiny" @click="openChannelDialog(item)">编辑</button>
            </div>
          </div>
        </div>
      </div>

      <!-- 分类模板 -->
      <div v-if="activeTab === 'category'" class="tab-content">
        <div class="tab-actions">
          <button class="primary-action" @click="openCategoryDialog()">新增模板</button>
        </div>
        <div class="template-grid">
          <div v-for="item in categoryTemplates" :key="item.id" class="template-card">
            <div class="template-icon" :style="{ background: item.color || (item.type === 'income' ? '#22c55e' : '#ef4444') }">
              {{ item.icon || (item.type === 'income' ? '💰' : '💸') }}
            </div>
            <div class="template-info">
              <h3>{{ item.name }}</h3>
              <p>{{ item.templateCode }} · {{ item.type === 'income' ? '收入' : '支出' }}</p>
            </div>
            <div class="template-status">
              <span :class="['management-status-tag', { active: item.enabled }]">
                {{ item.enabled ? '启用' : '禁用' }}
              </span>
            </div>
            <div class="template-actions">
              <button class="text-action tiny" @click="openCategoryDialog(item)">编辑</button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 账户模板弹窗 -->
    <el-dialog v-model="accountDialog.visible" :title="accountDialog.data?.id ? '编辑账户模板' : '新增账户模板'" width="500px">
      <el-form :model="accountDialog.data" label-width="100px">
        <el-form-item label="模板代码">
          <el-input v-model="accountDialog.data.providerCode" :disabled="!!accountDialog.data?.id" />
        </el-form-item>
        <el-form-item label="模板名称">
          <el-input v-model="accountDialog.data.name" />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="accountDialog.data.type" style="width: 100%">
            <el-option label="现金" value="cash" />
            <el-option label="钱包" value="wallet" />
            <el-option label="银行" value="bank" />
            <el-option label="信用卡" value="credit" />
            <el-option label="投资" value="investment" />
            <el-option label="负债" value="liability" />
            <el-option label="其他" value="other" />
          </el-select>
        </el-form-item>
        <el-form-item label="图标">
          <el-select v-model="accountDialog.data.icon" style="width: 100%">
            <el-option v-for="item in accountIconOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="颜色">
          <div class="admin-color-row">
            <button v-for="color in accountColors" :key="color" type="button" :class="['admin-color-dot', { active: accountDialog.data.color === color }]" :style="{ background: color }" @click="accountDialog.data.color = color"></button>
            <el-input v-model="accountDialog.data.color" style="width: 130px" />
          </div>
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="accountDialog.data.sort" :min="0" />
        </el-form-item>
        <el-form-item label="状态" v-if="accountDialog.data?.id">
          <el-switch v-model="accountDialog.data.enabled" active-text="启用" inactive-text="禁用" />
        </el-form-item>
      </el-form>
      <template #footer>
        <button class="ghost-action" @click="accountDialog.visible = false">取消</button>
        <button class="primary-action" @click="saveAccountTemplate">保存</button>
      </template>
    </el-dialog>

    <!-- 渠道模板弹窗 -->
    <el-dialog v-model="channelDialog.visible" :title="channelDialog.data?.id ? '编辑渠道模板' : '新增渠道模板'" width="500px">
      <el-form :model="channelDialog.data" label-width="100px">
        <el-form-item label="渠道代码">
          <el-input v-model="channelDialog.data.channelCode" :disabled="!!channelDialog.data?.id" />
        </el-form-item>
        <el-form-item label="渠道名称">
          <el-input v-model="channelDialog.data.name" />
        </el-form-item>
        <el-form-item label="渠道类型">
          <el-input v-model="channelDialog.data.channelType" />
        </el-form-item>
        <el-form-item label="所属平台">
          <el-input v-model="channelDialog.data.providerCode" />
        </el-form-item>
        <el-form-item label="支持事件">
          <el-input v-model="channelDialog.data.supportedEventTypes" placeholder="如: income,expense" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="channelDialog.data.sort" :min="0" />
        </el-form-item>
        <el-form-item label="状态" v-if="channelDialog.data?.id">
          <el-switch v-model="channelDialog.data.enabled" active-text="启用" inactive-text="禁用" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="channelDialog.visible = false">取消</el-button>
        <el-button type="primary" @click="saveChannelTemplate">保存</el-button>
      </template>
    </el-dialog>

    <!-- 分类模板弹窗 -->
    <el-dialog v-model="categoryDialog.visible" :title="categoryDialog.data?.id ? '编辑分类模板' : '新增分类模板'" width="500px">
      <el-form :model="categoryDialog.data" label-width="100px">
        <el-form-item label="模板代码">
          <el-input v-model="categoryDialog.data.templateCode" :disabled="!!categoryDialog.data?.id" />
        </el-form-item>
        <el-form-item label="分类名称">
          <el-input v-model="categoryDialog.data.name" />
        </el-form-item>
        <el-form-item label="收支类型">
          <el-select v-model="categoryDialog.data.type" style="width: 100%">
            <el-option label="收入" value="income" />
            <el-option label="支出" value="expense" />
          </el-select>
        </el-form-item>
        <el-form-item label="图标与颜色">
          <IconColorPicker
            v-model:icon-value="categoryDialog.data.icon"
            v-model:color-value="categoryDialog.data.color"
            icon-label="图标"
            color-label="颜色"
          />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="categoryDialog.data.sort" :min="0" />
        </el-form-item>
        <el-form-item label="状态" v-if="categoryDialog.data?.id">
          <el-switch v-model="categoryDialog.data.enabled" active-text="启用" inactive-text="禁用" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="categoryDialog.visible = false">取消</el-button>
        <el-button type="primary" @click="saveCategoryTemplate">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { listAccountTemplatesAdmin, createAccountTemplate, updateAccountTemplate } from '@/api/accountTemplate'
import { listChannelTemplatesAdmin, createChannelTemplate, updateChannelTemplate } from '@/api/channelTemplate'
import { listCategoryTemplatesAdmin, createCategoryTemplate, updateCategoryTemplate } from '@/api/categoryTemplate'
import IconColorPicker from '@/components/IconColorPicker.vue'

const activeTab = ref('account')
const accountTemplates = ref([])
const channelTemplates = ref([])
const categoryTemplates = ref([])

const loading = reactive({
  account: false,
  channel: false,
  category: false
})

const accountIconOptions = [
  { label: '钱包', value: 'Wallet' },
  { label: '银行卡', value: 'CreditCard' },
  { label: '现金', value: 'Money' },
  { label: '投资', value: 'TrendCharts' },
  { label: '银行', value: 'OfficeBuilding' },
  { label: '家庭', value: 'House' },
  { label: '集合', value: 'Collection' },
  { label: '其他', value: 'More' },
]
const accountColors = ['#2f7d5c', '#2f7d5c', '#1f2933', '#f59e0b', '#ef4444', '#f97316', '#22c55e', '#06b6d4', '#8b5cf6', '#ec4899', '#1677ff']

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
    ElMessage.error('加载账户模板失败')
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
    ElMessage.error('加载渠道模板失败')
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
    ElMessage.error('加载分类模板失败')
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
  const labels = { cash: '现金', wallet: '钱包', bank: '银行', credit: '信用卡', investment: '投资', liability: '负债', other: '其他' }
  return labels[String(type || '').toLowerCase()] || type || '其他'
}

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
      ElMessage.success('更新成功')
    } else {
      await createAccountTemplate(accountDialog.data)
      ElMessage.success('创建成功')
    }
    accountDialog.visible = false
    loadAccountTemplates()
  } catch (e) {
    ElMessage.error(e.message || '操作失败')
  }
}

async function saveChannelTemplate() {
  try {
    if (channelDialog.data.id) {
      await updateChannelTemplate(channelDialog.data.id, channelDialog.data)
      ElMessage.success('更新成功')
    } else {
      await createChannelTemplate(channelDialog.data)
      ElMessage.success('创建成功')
    }
    channelDialog.visible = false
    loadChannelTemplates()
  } catch (e) {
    ElMessage.error(e.message || '操作失败')
  }
}

async function saveCategoryTemplate() {
  try {
    if (categoryDialog.data.id) {
      await updateCategoryTemplate(categoryDialog.data.id, categoryDialog.data)
      ElMessage.success('更新成功')
    } else {
      await createCategoryTemplate(categoryDialog.data)
      ElMessage.success('创建成功')
    }
    categoryDialog.visible = false
    loadCategoryTemplates()
  } catch (e) {
    ElMessage.error(e.message || '操作失败')
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
