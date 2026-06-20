<template>
  <div class="category-group-page">
    <div class="page-header">
      <div>
        <p class="eyebrow">{{ t('dashboard.categoryGroup.eyebrow') }}</p>
        <h1>{{ t('dashboard.categoryGroup.title') }}</h1>
        <p>{{ t('dashboard.categoryGroup.subtitle') }}</p>
      </div>
      <button class="primary-action" @click="openCreateGroup">{{ t('dashboard.categoryGroup.actions.new') }}</button>
    </div>

    <!-- 收支类型切换 -->
    <div class="type-tabs">
      <button :class="['type-tab', { active: activeType === 'expense' }]" @click="activeType = 'expense'">
        {{ t('domain.expense') }}
      </button>
      <button :class="['type-tab', { active: activeType === 'income' }]" @click="activeType = 'income'">
        {{ t('domain.income') }}
      </button>
    </div>

    <!-- 分组列表 -->
    <div class="group-grid" v-loading="loading">
      <div v-for="group in filteredGroups" :key="group.id" class="group-card">
        <div class="group-header">
          <div class="group-info">
            <span class="group-icon" :style="{ background: group.color || '#3b82f6' }">
              {{ group.icon || '📁' }}
            </span>
            <div>
              <h3>{{ group.name }}</h3>
              <span class="group-count">{{ group.memberCount || 0 }} 个分类</span>
            </div>
          </div>
          <div class="group-actions">
            <button class="text-action" @click="openEditGroup(group)">编辑</button>
            <button class="text-action danger" @click="handleDeleteGroup(group)">删除</button>
          </div>
        </div>

        <div class="group-members">
          <div v-if="groupMembers[group.id]?.length > 0" class="member-tags">
            <el-tag v-for="cat in groupMembers[group.id]" :key="cat.id" size="small" closable @close="removeMember(group.id, cat.id)">
              {{ cat.name }}
            </el-tag>
          </div>
          <div v-else class="empty-members">
            <span>{{ t('dashboard.categoryGroup.noCategories') }}</span>
          </div>
        </div>

        <div class="group-footer">
          <button class="ghost-action small" @click="openAddMember(group)">
            + {{ t('dashboard.categoryGroup.actions.addCategory') }}
          </button>
        </div>
      </div>

      <div v-if="filteredGroups.length === 0 && !loading" class="empty-state">
        <el-icon size="48" color="#cbd5e1"><Folder /></el-icon>
        <h3>{{ t('dashboard.categoryGroup.empty.title') }}</h3>
        <p>{{ t('dashboard.categoryGroup.empty.text') }}</p>
        <button class="primary-action" @click="openCreateGroup">{{ t('dashboard.categoryGroup.actions.new') }}</button>
      </div>
    </div>

    <!-- 创建/编辑分组弹窗 -->
    <el-dialog v-model="groupDialog.visible" :title="groupDialog.data?.id ? t('dashboard.categoryGroup.dialog.editTitle') : t('dashboard.categoryGroup.dialog.createTitle')" width="480px">
      <el-form :model="groupDialog.data" label-width="100px">
        <el-form-item label="分组名称" prop="name">
          <el-input v-model="groupDialog.data.name" placeholder="例如：餐饮娱乐" />
        </el-form-item>
        <el-form-item label="收支类型" prop="type">
          <el-select v-model="groupDialog.data.type" style="width: 100%" disabled>
            <el-option label="支出" value="expense" />
            <el-option label="收入" value="income" />
          </el-select>
        </el-form-item>
        <el-form-item label="图标">
          <el-input v-model="groupDialog.data.icon" placeholder="emoji 图标" maxlength="4" />
        </el-form-item>
        <el-form-item label="颜色">
          <div class="color-palette">
            <button v-for="color in colors" :key="color" type="button" :class="['color-dot', { active: groupDialog.data.color === color }]" :style="{ background: color }" @click="groupDialog.data.color = color"></button>
          </div>
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="groupDialog.data.sort" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <button class="ghost-action" @click="groupDialog.visible = false">{{ t('common.actions.cancel') }}</button>
        <button class="primary-action" @click="saveGroup">{{ t('common.actions.save') }}</button>
      </template>
    </el-dialog>

    <!-- 添加成员弹窗 -->
    <el-dialog v-model="memberDialog.visible" title="添加分类到分组" width="520px">
      <div class="category-selector">
        <el-checkbox-group v-model="selectedMemberIds">
          <div class="category-list">
            <label v-for="cat in availableCategories" :key="cat.id" class="category-item">
              <el-checkbox :value="cat.id" :disabled="isCategoryInGroup(cat.id)" />
              <span class="category-name">{{ cat.name }}</span>
              <el-tag v-if="isCategoryInGroup(cat.id)" size="small" type="info">已添加</el-tag>
            </label>
          </div>
        </el-checkbox-group>
      </div>
      <template #footer>
        <button class="ghost-action" @click="memberDialog.visible = false">{{ t('common.actions.cancel') }}</button>
        <button class="primary-action" @click="addMembers" :disabled="selectedMemberIds.length === 0">{{ t('common.actions.save') }}</button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Folder } from '@element-plus/icons-vue'
import { listCategories } from '@/api/category/category'
import {
  listFamilyCategoryGroups,
  createFamilyCategoryGroup,
  updateFamilyCategoryGroup,
  deleteFamilyCategoryGroup,
  addCategoriesToGroup,
  removeCategoryFromGroup,
  getGroupCategoryIds
} from '@/api/familyCategoryGroup'

const { t } = useI18n()
const loading = ref(false)
const activeType = ref('expense')
const familyId = 1 // TODO: 从家庭选择或路由中获取

const groups = ref([])
const categories = ref([])
const groupMembers = reactive({}) // groupId -> [{id, name}]
const selectedMemberIds = ref([])

const colors = ['#3b82f6', '#2f7d5c', '#1f2933', '#f59e0b', '#ef4444', '#f97316', '#8b5cf6', '#ec4899']

const groupDialog = reactive({
  visible: false,
  data: null
})

const memberDialog = reactive({
  visible: false,
  groupId: 0
})

const filteredGroups = computed(() => {
  return groups.value.filter(g => g.type === activeType.value)
})

const availableCategories = computed(() => {
  return categories.value.filter(c => c.type === activeType.value)
})

onMounted(() => {
  loadGroups()
  loadCategories()
})

async function loadGroups() {
  loading.value = true
  try {
    const res = await listFamilyCategoryGroups(familyId, { type: activeType.value })
    if (res.success) {
      groups.value = res.data || []
      // 加载每个组的成员
      for (const g of groups.value) {
        loadGroupMembers(g.id)
      }
    }
  } catch (e) {
    ElMessage.error('加载分组失败')
  } finally {
    loading.value = false
  }
}

async function loadGroupMembers(groupId) {
  try {
    const res = await getGroupCategoryIds(familyId, groupId)
    if (res.success) {
      const categoryIds = res.data || []
      groupMembers[groupId] = categoryIds.map(id => {
        const cat = categories.value.find(c => c.id === id)
        return { id, name: cat?.name || `分类#${id}` }
      })
    }
  } catch (e) {
    console.error('加载组成员失败', e)
  }
}

async function loadCategories() {
  try {
    const res = await listCategories()
    if (res.success) {
      categories.value = res.data || []
    }
  } catch (e) {
    ElMessage.error('加载分类失败')
  }
}

function openCreateGroup() {
  groupDialog.data = {
    name: '',
    type: activeType.value,
    icon: '📁',
    color: colors[0],
    sort: 0
  }
  groupDialog.visible = true
}

function openEditGroup(group) {
  groupDialog.data = { ...group }
  groupDialog.visible = true
}

async function saveGroup() {
  try {
    if (groupDialog.data.id) {
      await updateFamilyCategoryGroup(familyId, groupDialog.data.id, groupDialog.data)
      ElMessage.success('更新成功')
    } else {
      const res = await createFamilyCategoryGroup(familyId, groupDialog.data)
      if (res.success) {
        ElMessage.success('创建成功')
      }
    }
    groupDialog.visible = false
    loadGroups()
  } catch (e) {
    ElMessage.error(e.message || '操作失败')
  }
}

async function handleDeleteGroup(group) {
  try {
    await ElMessageBox.confirm(`确定删除分组「${group.name}」吗？`, '提示', { type: 'warning' })
    await deleteFamilyCategoryGroup(familyId, group.id)
    ElMessage.success('删除成功')
    loadGroups()
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '删除失败')
    }
  }
}

function openAddMember(group) {
  memberDialog.groupId = group.id
  selectedMemberIds.value = []
  memberDialog.visible = true
}

function isCategoryInGroup(categoryId) {
  const members = groupMembers[memberDialog.groupId] || []
  return members.some(m => m.id === categoryId)
}

async function addMembers() {
  if (selectedMemberIds.value.length === 0) return
  try {
    await addCategoriesToGroup(familyId, memberDialog.groupId, selectedMemberIds.value)
    ElMessage.success('添加成功')
    memberDialog.visible = false
    loadGroupMembers(memberDialog.groupId)
  } catch (e) {
    ElMessage.error(e.message || '添加失败')
  }
}

async function removeMember(groupId, categoryId) {
  try {
    await removeCategoryFromGroup(familyId, groupId, categoryId)
    ElMessage.success('移除成功')
    loadGroupMembers(groupId)
  } catch (e) {
    ElMessage.error(e.message || '移除失败')
  }
}
</script>

<style scoped>
.category-group-page {
  padding: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
}

.eyebrow {
  font-size: 14px;
  color: #2f7d5c;
  font-weight: 500;
  margin: 0 0 8px 0;
}

.page-header h1 {
  font-size: 28px;
  font-weight: 600;
  color: #1e293b;
  margin: 0 0 8px 0;
}

.page-header p {
  color: #64748b;
  margin: 0;
}

.type-tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 24px;
}

.type-tab {
  padding: 8px 20px;
  border-radius: 8px;
  border: none;
  background: #f1f5f9;
  color: #64748b;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
}

.type-tab:hover {
  background: #e2e8f0;
}

.type-tab.active {
  background: #2f7d5c;
  color: white;
}

.group-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
  gap: 20px;
}

.group-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.08);
}

.group-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 16px;
}

.group-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.group-icon {
  width: 44px;
  height: 44px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
}

.group-info h3 {
  font-size: 16px;
  font-weight: 600;
  color: #1e293b;
  margin: 0 0 4px 0;
}

.group-count {
  font-size: 13px;
  color: #64748b;
}

.group-actions {
  display: flex;
  gap: 8px;
}

.group-members {
  min-height: 40px;
  padding: 12px 0;
  border-top: 1px solid #f1f5f9;
  border-bottom: 1px solid #f1f5f9;
}

.member-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.empty-members {
  color: #94a3b8;
  font-size: 13px;
  text-align: center;
  padding: 8px 0;
}

.group-footer {
  padding-top: 12px;
}

.color-palette {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.color-dot {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  border: 2px solid transparent;
  cursor: pointer;
  transition: all 0.2s;
}

.color-dot:hover,
.color-dot.active {
  border-color: #cbd5e1;
  transform: scale(1.1);
}

.category-selector {
  max-height: 400px;
  overflow-y: auto;
}

.category-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.category-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: 6px;
  transition: background 0.2s;
}

.category-item:hover {
  background: #f8fafc;
}

.category-name {
  flex: 1;
  font-size: 14px;
  color: #334155;
}

.empty-state {
  grid-column: 1 / -1;
  text-align: center;
  padding: 60px 20px;
  background: white;
  border-radius: 12px;
}

.empty-state h3 {
  font-size: 18px;
  color: #1e293b;
  margin: 16px 0 8px 0;
}

.empty-state p {
  color: #64748b;
  margin-bottom: 20px;
}

.small {
  padding: 6px 12px;
  font-size: 13px;
}

.danger {
  color: #ef4444;
}
</style>
