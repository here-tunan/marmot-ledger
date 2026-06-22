<template>
  <main class="ledger-page">
    <ManagementPageHeader :eyebrow="t('categories.hero.eyebrow')" :title="t('categories.hero.title')" :subtitle="t('categories.hero.subtitle')">
      <template #actions>
        <button class="management-primary-action" @click="openCreate">{{ t('categories.actions.new') }}</button>
      </template>
    </ManagementPageHeader>

    <ManagementToolbar class="reveal-block delay-1">
      <el-select v-model="filters.type" clearable :placeholder="t('categories.fields.type')" class="filter-control" @change="loadCategories">
        <el-option :label="t('domain.expense')" value="expense" />
        <el-option :label="t('domain.income')" value="income" />
      </el-select>
      <el-select v-model="filters.isActive" clearable :placeholder="t('common.filters.enabledStatus')" class="filter-control" @change="loadCategories">
        <el-option :label="t('common.status.enabled')" :value="true" />
        <el-option :label="t('common.status.disabled')" :value="false" />
      </el-select>
      <button class="management-ghost-action" @click="loadCategories">{{ t('common.actions.refresh') }}</button>
    </ManagementToolbar>

    <section v-loading="loading" class="category-grid management-grid reveal-block delay-2">
      <article v-for="item in categories" :key="item.id" class="category-card">
        <!-- 左侧颜色条 -->
        <div
          class="category-color-bar"
          :style="{ background: item.color || (item.type === 'income' ? '#ef4444' : '#f97316') }"
        ></div>

        <div class="category-content">
          <!-- 图标和标题 -->
          <div class="category-main">
            <div
              class="category-icon-wrapper"
              :style="{ background: (item.color || (item.type === 'income' ? '#ef4444' : '#f97316')) + '20' }"
            >
              <span class="category-icon">{{ item.icon || '📁' }}</span>
            </div>
            <div class="category-text">
              <h2 class="category-name">{{ item.name }}</h2>
              <p class="category-group">{{ item.type === 'income' ? t('domain.income') : t('domain.expense') }}</p>
            </div>
          </div>

          <!-- 底部标签 -->
          <div class="category-footer">
            <span :class="['management-type-tag', item.type]">
              {{ item.type === 'income' ? t('domain.income') : t('domain.expense') }}
            </span>
            <span :class="['management-status-tag', { active: item.isActive }]">
              {{ item.isActive ? t('common.status.enabled') : t('common.status.disabled') }}
            </span>
          </div>
        </div>

        <div class="card-actions">
          <button class="management-text-action" @click="openEdit(item)">{{ t('common.actions.edit') }}</button>
          <button class="management-danger-action" @click="handleDelete(item)">{{ t('common.actions.delete') }}</button>
        </div>
      </article>

      <ManagementEmptyState v-if="!loading && !categories.length" :image="marmotOne" :alt="t('categories.empty.alt')" :title="t('categories.empty.title')" :text="t('categories.empty.text')">
        <template #action>
          <button class="management-primary-action" @click="openCreate">{{ t('categories.actions.new') }}</button>
        </template>
      </ManagementEmptyState>
    </section>

    <el-dialog v-model="dialogVisible" :title="editingId ? t('categories.dialog.editTitle') : t('categories.dialog.createTitle')" width="640px" class="marmot-dialog category-dialog">
      <!-- 模板快速选择 -->
      <div v-if="!editingId" class="template-quick-select">
        <div class="section-header" @click="templateCollapsed = !templateCollapsed" style="cursor: pointer;">
          <span class="section-icon">⚡</span>
          <span class="section-title">快速选择模板</span>
          <span class="toggle-icon">{{ templateCollapsed ? '展开' : '收起' }} {{ templateCollapsed ? '▶' : '▼' }}</span>
        </div>
        <template v-if="!templateCollapsed">
          <div v-if="filteredTemplates.length === 0" class="no-templates">
            <span>暂无该类型的模板</span>
          </div>
          <div v-else class="template-grid-compact">
          <button
            v-for="tpl in filteredTemplates"
            :key="tpl.id"
            :class="['template-item', { active: selectedTemplate?.id === tpl.id }]"
            :style="{
              background: selectedTemplate?.id === tpl.id ? (tpl.color || '#2f7d5c') : '#f5f7fa',
              color: selectedTemplate?.id === tpl.id ? '#fff' : '#1e293b'
            }"
            @click="selectTemplate(tpl)"
          >
            <span class="item-icon">{{ tpl.icon || '📁' }}</span>
            <span class="item-name">{{ tpl.name }}</span>
          </button>
        </div>
        </template>
      </div>

      <el-form ref="formRef" :model="form" label-position="top">
        <el-form-item :label="t('categories.fields.name')">
          <el-input v-model="form.name" placeholder="请输入分类名称" />
        </el-form-item>
        <el-form-item :label="t('categories.fields.type')">
          <el-select v-model="form.type" class="full-width" @change="handleTypeChange">
            <el-option :label="t('domain.expense')" value="expense" />
            <el-option :label="t('domain.income')" value="income" />
          </el-select>
        </el-form-item>

        <!-- 使用通用图标颜色选择器 -->
        <el-form-item label="图标与颜色">
          <IconColorPicker
            v-model:icon-value="form.icon"
            v-model:color-value="form.color"
            icon-label="图标"
            color-label="颜色"
          />
        </el-form-item>

        <!-- 家庭分组选择（如果有家庭） -->
        <el-form-item v-if="families.length > 0" label="加入家庭分组" class="family-groups-section">
          <p class="hint-text">分类归你所有，可以选择性加入家庭分组，方便联合统计</p>

          <!-- 家庭卡片网格布局 -->
          <div class="families-grid">
            <div v-for="fam in families" :key="fam.id" class="family-card">
              <div class="family-card-header">
                <span class="family-name">{{ fam.name }}</span>
                <span class="family-count">{{ getGroupsByFamily(fam.id).length }} 个分组</span>
              </div>

              <div class="family-card-content">
                <!-- 有分组时显示列表 -->
                <div v-if="getGroupsByFamily(fam.id).length > 0" class="groups-list">
                  <label
                    v-for="group in getGroupsByFamily(fam.id)"
                    :key="group.id"
                    class="group-item"
                  >
                    <el-checkbox v-model="selectedGroupIds" :value="group.id" />
                    <span class="group-icon" :style="{ background: group.color || '#e2e8f0' }">
                      {{ group.icon || '📁' }}
                    </span>
                    <span class="group-name">{{ group.name }}</span>
                  </label>
                </div>

                <!-- 无分组时显示引导 -->
                <div v-else class="no-groups-placeholder">
                  <span class="placeholder-icon">📭</span>
                  <span class="placeholder-text">暂无{{ form.type === 'expense' ? '支出' : '收入' }}类型分组</span>
                  <el-button type="primary" size="small" @click="goToCreateGroup(fam.id)">
                    + 创建分组
                  </el-button>
                </div>
              </div>
            </div>
          </div>
        </el-form-item>

        <el-form-item v-if="editingId" :label="t('common.status.status')">
          <el-switch v-model="form.isActive" :active-text="t('common.status.enabled')" :inactive-text="t('common.status.disabled')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <button class="management-ghost-action" @click="dialogVisible = false">{{ t('common.actions.cancel') }}</button>
        <button class="management-primary-action" @click="submitForm">{{ t('common.actions.save') }}</button>
      </template>
    </el-dialog>
  </main>
</template>

<script setup>
import { computed, onActivated, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { createCategory, deleteCategory, listCategories, updateCategory, checkCategoryUsage } from '@/api/category/category'
import { listCategoryTemplates } from '@/api/categoryTemplate'
import { listFamilies } from '@/api/family/family'
import { listFamilyCategoryGroups } from '@/api/familyCategoryGroup'
import IconColorPicker from '@/components/IconColorPicker.vue'
import ManagementPageHeader from '@/components/management/ManagementPageHeader.vue'
import ManagementToolbar from '@/components/management/ManagementToolbar.vue'
import ManagementEmptyState from '@/components/management/ManagementEmptyState.vue'
import marmotOne from '../../../img/marmot-ledger-1.png'

const { t } = useI18n()
const router = useRouter()
const categories = ref([])
const templates = ref([])
const families = ref([])
const familyGroups = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const editingId = ref(0)
const formRef = ref()
const filters = reactive({ type: '', isActive: '' })
const form = reactive({ name: '', type: 'expense', icon: '📁', color: '', groupIds: [], isActive: true })
const selectedTemplate = ref(null)
const selectedGroupIds = ref([])
const templateCollapsed = ref(false)

const filteredTemplates = computed(() => templates.value.filter(tpl => tpl.type === form.type))

function selectTemplate(tpl) {
  selectedTemplate.value = tpl
  form.name = tpl.name
  form.icon = tpl.icon || '📁'
  form.color = tpl.color || ''
}

async function loadTemplates() {
  const res = await listCategoryTemplates({ enabled: true })
  if (res.success) templates.value = res.data || []
}

async function loadFamilies() {
  const res = await listFamilies()
  if (res.success) families.value = res.data || []
}

async function loadFamilyGroups() {
  if (!families.value.length) return
  const allGroups = []
  for (const fam of families.value) {
    const res = await listFamilyCategoryGroups(fam.id, { type: form.type })
    if (res.success && res.data) {
      allGroups.push(...res.data.map(g => ({ ...g, familyId: fam.id })))
    }
  }
  familyGroups.value = allGroups
}

async function loadCategories() {
  loading.value = true
  try {
    const params = {}
    if (filters.type) params.type = filters.type
    if (filters.isActive !== '') params.isActive = filters.isActive
    const res = await listCategories(params)
    if (res.success) categories.value = res.data || []
    else ElMessage.error(res.error || t('categories.messages.loadFailed'))
  } finally { loading.value = false }
}

async function refreshAll() {
  await Promise.all([loadCategories(), loadTemplates(), loadFamilies()])
  await loadFamilyGroups()
}

function openCreate() {
  editingId.value = 0
  selectedTemplate.value = null
  selectedGroupIds.value = []
  Object.assign(form, { name: '', type: 'expense', icon: '📁', color: '', groupIds: [], isActive: true })
  dialogVisible.value = true
}

function openEdit(item) {
  editingId.value = item.id
  selectedTemplate.value = null
  selectedGroupIds.value = item.groupIds ? [...item.groupIds] : []
  Object.assign(form, {
    name: item.name,
    type: item.type,
    icon: item.icon || '📁',
    color: item.color || '',
    groupIds: item.groupIds || [],
    isActive: item.isActive !== false
  })
  dialogVisible.value = true
}

function handleTypeChange() {
  selectedGroupIds.value = []
  loadFamilyGroups()
}

function getFamilyName(familyId) {
  const fam = families.value.find(f => f.id === familyId)
  return fam ? fam.name : '未知家庭'
}

function getGroupsByFamily(familyId) {
  return familyGroups.value.filter(g => g.familyId === familyId)
}

function goToCreateGroup(familyId) {
  // 关闭当前弹窗
  dialogVisible.value = false
  // 跳转到家庭页面，并传递选中的家庭ID和切换到分类组Tab
  router.push({
    name: 'family',
    query: { familyId: familyId, tab: 'categoryGroups', createNew: 'true' }
  })
}

async function submitForm() {
  if (!form.name) return ElMessage.warning(t('categories.validation.nameRequired'))
  if (!form.type) return ElMessage.warning(t('categories.validation.typeRequired'))

  const payload = {
    ...form,
    icon: form.icon || '',
    color: form.color || '',
    groupIds: selectedGroupIds.value,
    isActive: form.isActive !== false
  }

  const res = editingId.value ? await updateCategory(editingId.value, payload) : await createCategory(payload)
  if (res.success) {
    ElMessage.success(editingId.value ? t('categories.messages.updated') : t('categories.messages.created'))
    dialogVisible.value = false
    await loadCategories()
  } else {
    ElMessage.error(res.error || t('categories.messages.saveFailed'))
  }
}

async function handleDelete(item) {
  try {
    // 先检查使用情况
    const usageRes = await checkCategoryUsage(item.id)
    const eventCount = usageRes.success ? usageRes.data.eventCount : 0

    let confirmMessage = t('categories.delete.confirm', { name: item.name })
    if (eventCount > 0) {
      confirmMessage = `⚠️  分类「${item.name}」下有 ${eventCount} 条账单记录\n\n删除后，这些账单在统计时将归为「其他分类」\n\n确定要删除吗？`
    }

    await ElMessageBox.confirm(confirmMessage, t('categories.delete.title'), {
      confirmButtonText: t('common.actions.delete'),
      cancelButtonText: t('common.actions.cancel'),
      type: eventCount > 0 ? 'warning' : 'info',
      customClass: 'calm-marmot-message-box calm-marmot-delete-box',
      confirmButtonClass: 'calm-marmot-danger-confirm',
      cancelButtonClass: 'calm-marmot-soft-cancel'
    })

    const res = await deleteCategory(item.id)
    if (res.success) {
      if (res.data.affectedCount > 0) {
        ElMessage.success(`已删除，共 ${res.data.affectedCount} 条账单受影响`)
      } else {
        ElMessage.success(t('categories.messages.deleted'))
      }
      await loadCategories()
    } else {
      ElMessage.error(res.error || t('categories.messages.deleteFailed'))
    }
  } catch (err) {
    if (err !== 'cancel') console.warn(err)
  }
}

onMounted(refreshAll)
onActivated(refreshAll)
</script>

<style scoped>
.ledger-page {
  max-width: 1200px;
  margin: 0 auto;
  color: #1e293b;
}

.reveal-block,
.category-card {
  animation: revealUp 480ms cubic-bezier(.16, 1, .3, 1) both;
}

.delay-1 {
  animation-delay: 90ms;
}

.delay-2 {
  animation-delay: 160ms;
}

.category-card {
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 1px 3px rgba(15, 23, 42, .1), 0 12px 30px rgba(15, 23, 42, .04);
}

.filter-control {
  width: 180px;
}

.category-card {
  display: grid;
  grid-template-columns: 6px minmax(0, 1fr) auto;
  gap: 0;
  overflow: hidden;
  transition-property: transform, box-shadow;
  transition-duration: 180ms;
}

/* 左侧颜色条 */
.category-color-bar {
  height: 100%;
  border-radius: 12px 0 0 12px;
}

/* 中间内容区 */
.category-content {
  padding: 14px 12px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

/* 上半部分：图标和标题 */
.category-main {
  display: flex;
  align-items: center;
  gap: 10px;
}

.category-icon-wrapper {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.category-icon {
  font-size: 20px;
  line-height: 1;
}

.category-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.category-name {
  margin: 0;
  font-size: 15px;
  font-weight: 600;
  color: #1e293b;
  letter-spacing: -.012em;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.category-group {
  margin: 0;
  color: #64748b;
  font-size: 12px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 下半部分：标签 */
.category-footer {
  display: flex;
  gap: 6px;
  align-items: center;
}

.card-actions {
  display: flex;
  gap: 6px;
  align-items: center;
  padding-right: 14px;
}

.full-width {
  width: 100%;
}

@media (hover: hover) {
  .category-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 3px 8px rgba(15, 23, 42, .12), 0 16px 34px rgba(15, 23, 42, .06);
  }
}

@media (max-width: 900px) {
  .category-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 600px) {
  .category-grid {
    grid-template-columns: 1fr;
  }

  .page-hero,
  .toolbar {
    flex-direction: column;
    align-items: stretch;
  }

  .filter-control {
    width: 100%;
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

/* 分类弹窗宽度 */
.category-dialog :deep(.el-dialog__body) {
  padding: 20px 24px;
}

/* 模板快速选择样式 */
.template-quick-select {
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid #e2e8f0;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 10px;
  font-size: 13px;
  font-weight: 600;
  color: #64748b;
  user-select: none;
}

.section-icon {
  font-size: 14px;
}

.toggle-icon {
  margin-left: auto;
  font-size: 11px;
  color: #94a3b8;
  font-weight: 500;
}

.no-templates {
  padding: 14px;
  background: #f8fafc;
  border-radius: 10px;
  text-align: center;
  color: #94a3b8;
  font-size: 13px;
}

/* 紧凑网格布局 - 3列 */
.template-grid-compact {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 8px;
}

.template-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  border-radius: 10px;
  border: none;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.2s ease;
}

.template-item:hover {
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.template-item.active {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  font-weight: 600;
}

.item-icon {
  font-size: 16px;
  line-height: 1;
  flex-shrink: 0;
}

.item-name {
  line-height: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 家庭分组选择样式 - 卡片网格布局 */
.family-groups-section .hint-text {
  margin: 0 0 16px 0;
  color: #64748b;
  font-size: 13px;
}

/* 家庭卡片网格 */
.families-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.family-card {
  display: flex;
  flex-direction: column;
  background: #f8fafc;
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid #e2e8f0;
}

.family-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 14px;
  background: #fff;
  border-bottom: 1px solid #e2e8f0;
}

.family-name {
  font-size: 14px;
  font-weight: 600;
  color: #1e293b;
}

.family-count {
  font-size: 12px;
  color: #64748b;
  background: #f1f5f9;
  padding: 2px 8px;
  border-radius: 10px;
}

.family-card-content {
  flex: 1;
  min-height: 100px;
  padding: 8px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

/* 分组列表 */
.groups-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.group-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 10px;
  background: #fff;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.group-item:hover {
  background: #f0f4f0;
}

.group-icon {
  font-size: 14px;
  width: 26px;
  height: 26px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  flex-shrink: 0;
}

.group-name {
  flex: 1;
  font-size: 13px;
  color: #334155;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 无分组占位 */
.no-groups-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 16px 10px;
  text-align: center;
  height: 100%;
}

.placeholder-icon {
  font-size: 24px;
  opacity: 0.5;
}

.placeholder-text {
  font-size: 12px;
  color: #94a3b8;
}

.no-groups-placeholder .el-button {
  font-size: 12px;
  padding: 4px 12px;
}

/* 响应式：单列 */
@media (max-width: 520px) {
  .families-grid {
    grid-template-columns: 1fr;
  }
}
</style>
