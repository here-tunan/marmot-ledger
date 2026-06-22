<template>
  <main class="family-page">
    <ManagementPageHeader :eyebrow="t('family.hero.eyebrow')" :title="t('family.hero.title')" :subtitle="t('family.hero.subtitle')">
      <template #actions>
        <button class="management-primary-action" @click="openCreate">{{ t('family.actions.new') }}</button>
      </template>
    </ManagementPageHeader>

    <section class="family-layout reveal-block delay-1">
      <aside class="family-list">
        <h2>{{ t('family.sections.myFamilies') }}</h2>
        <button v-for="item in families" :key="item.id" class="family-item" :class="{ active: selectedFamily?.id === item.id }" @click="selectFamily(item)">
          <strong>{{ item.name }}</strong>
          <span>{{ item.role }}</span>
        </button>
        <p v-if="!families.length" class="muted">{{ t('family.empty.noFamilies') }}</p>

        <h2>{{ t('family.sections.invitations') }}</h2>
        <div v-for="item in invitations" :key="item.id" class="invite-card">
          <strong>{{ item.familyName }}</strong>
          <span>{{ item.displayName || item.account }}</span>
          <div>
            <button class="primary-action small" @click="acceptInvite(item)">{{ t('family.actions.accept') }}</button>
            <button class="ghost-action small" @click="rejectInvite(item)">{{ t('family.actions.reject') }}</button>
          </div>
        </div>
        <p v-if="!invitations.length" class="muted">{{ t('family.empty.noInvitations') }}</p>
      </aside>

      <section class="family-detail">
        <template v-if="selectedFamily">
          <!-- Tab 切换 -->
          <div class="detail-tabs">
            <button :class="['tab-btn', { active: activeTab === 'members' }]" @click="activeTab = 'members'">
              {{ t('family.sections.members') }}
            </button>
            <button :class="['tab-btn', { active: activeTab === 'categoryGroups' }]" @click="activeTab = 'categoryGroups'">
              {{ t('routes.familyCategoryGroups') }}
            </button>
          </div>

          <!-- 成员 Tab -->
          <div v-if="activeTab === 'members'" class="tab-content">
            <div class="section-head">
              <div>
                <p class="eyebrow">{{ t('family.sections.members') }}</p>
                <h2>{{ selectedFamily.name }}</h2>
              </div>
            </div>

            <div class="invite-form">
              <el-input v-model="inviteForm.account" :placeholder="t('family.fields.account')" />
              <el-input v-model="inviteForm.displayName" :placeholder="t('family.fields.displayName')" />
              <button class="primary-action invite-action" @click="inviteMember">{{ t('family.actions.invite') }}</button>
            </div>

            <h3>{{ t('family.sections.members') }}</h3>
            <div class="member-grid">
              <div v-for="member in members" :key="member.id" class="member-card">
                <strong>{{ member.displayName || member.name || member.account }}</strong>
                <span>{{ member.role }} · {{ member.status }}</span>
              </div>
            </div>
          </div>

          <!-- 分类组 Tab -->
          <div v-if="activeTab === 'categoryGroups'" class="tab-content">
            <div class="section-head">
              <div>
                <p class="eyebrow">{{ t('dashboard.categoryGroup.eyebrow') }}</p>
                <h2>{{ selectedFamily.name }} - {{ t('routes.familyCategoryGroups') }}</h2>
              </div>
              <button class="primary-action" @click="openCreateGroup">{{ t('dashboard.categoryGroup.actions.new') }}</button>
            </div>

            <!-- 收支类型切换 -->
            <div class="type-tabs">
              <button :class="['type-tab', { active: groupTypeFilter === 'expense' }]" @click="groupTypeFilter = 'expense'">
                {{ t('domain.expense') }}
              </button>
              <button :class="['type-tab', { active: groupTypeFilter === 'income' }]" @click="groupTypeFilter = 'income'">
                {{ t('domain.income') }}
              </button>
            </div>

            <!-- 分组列表 -->
            <div class="group-grid" v-loading="loadingGroups">
              <div v-for="group in filteredGroups" :key="group.id" class="group-card">
                <div class="group-header">
                  <div class="group-info">
                    <span class="group-icon" :style="{ background: group.color || '#2f7d5c' }">
                      {{ group.icon || '📁' }}
                    </span>
                    <div>
                      <h3>{{ group.name }}</h3>
                      <span class="management-meta-tag">{{ (groupMembers[group.id] || []).length }} 个分类</span>
                    </div>
                  </div>
                  <div class="group-actions">
                    <button class="text-action" @click="openEditGroup(group)">{{ t('common.actions.edit') }}</button>
                    <button class="danger-action" @click="handleDeleteGroup(group)">{{ t('common.actions.delete') }}</button>
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

              <div v-if="filteredGroups.length === 0 && !loadingGroups" class="empty-state-inline">
                <el-icon size="40" color="#cbd5e1"><Folder /></el-icon>
                <p>{{ t('dashboard.categoryGroup.empty.text') }}</p>
                <button class="primary-action small" @click="openCreateGroup">{{ t('dashboard.categoryGroup.actions.new') }}</button>
              </div>
            </div>
          </div>
        </template>
        <div v-else class="empty-state"><p>{{ t('family.empty.noFamilies') }}</p></div>
      </section>
    </section>

    <!-- 创建家庭弹窗 -->
    <el-dialog v-model="dialogVisible" :title="t('family.actions.new')" width="520px" class="marmot-dialog">
      <el-form label-position="top">
        <el-form-item :label="t('family.fields.name')">
          <el-input v-model="familyForm.name" />
        </el-form-item>
      </el-form>
      <template #footer>
        <button class="ghost-action" @click="dialogVisible = false">{{ t('common.actions.cancel') }}</button>
        <button class="primary-action" @click="createNewFamily">{{ t('common.actions.create') }}</button>
      </template>
    </el-dialog>

    <!-- 创建/编辑分组弹窗 -->
    <el-dialog v-model="groupDialog.visible" :title="groupDialog.data?.id ? t('dashboard.categoryGroup.dialog.editTitle') : t('dashboard.categoryGroup.dialog.createTitle')" width="480px" class="marmot-dialog">
      <el-form :model="groupDialog.data" label-position="top">
        <el-form-item label="分组名称" prop="name">
          <el-input v-model="groupDialog.data.name" placeholder="例如：餐饮娱乐" />
        </el-form-item>
        <el-form-item label="收支类型" prop="type">
          <el-select v-model="groupDialog.data.type" style="width: 100%">
            <el-option label="支出" value="expense" />
            <el-option label="收入" value="income" />
          </el-select>
        </el-form-item>
        <el-form-item label="图标与颜色">
          <IconColorPicker
            v-model:icon-value="groupDialog.data.icon"
            v-model:color-value="groupDialog.data.color"
            icon-label="图标"
            color-label="颜色"
          />
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
    <el-dialog v-model="memberDialog.visible" title="添加分类到分组" width="520px" class="marmot-dialog">
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
  </main>
</template>

<script setup>
import { onActivated, onMounted, reactive, ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Folder } from '@element-plus/icons-vue'
import { acceptFamilyInvitation, createFamily, inviteFamilyMember, listFamilies, listFamilyInvitations, listFamilyMembers, rejectFamilyInvitation } from '@/api/family/family'
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
import IconColorPicker from '@/components/IconColorPicker.vue'
import ManagementPageHeader from '@/components/management/ManagementPageHeader.vue'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()

// 家庭基础
const families = ref([])
const invitations = ref([])
const selectedFamily = ref(null)
const members = ref([])
const dialogVisible = ref(false)
const familyForm = reactive({ name: '' })
const inviteForm = reactive({ account: '', displayName: '' })

// Tab
const activeTab = ref('members')

// 分类组
const loadingGroups = ref(false)
const groupTypeFilter = ref('expense')
const groups = ref([])
const categories = ref([])
const groupMembers = reactive({})
const selectedMemberIds = ref([])

const groupDialog = reactive({
  visible: false,
  data: null
})

const memberDialog = reactive({
  visible: false,
  groupId: 0
})

const filteredGroups = computed(() => {
  return groups.value.filter(g => g.type === groupTypeFilter.value)
})

const availableCategories = computed(() => {
  return categories.value.filter(c => c.type === groupTypeFilter.value)
})

function openCreate() {
  familyForm.name = ''
  dialogVisible.value = true
}

async function loadFamilies() {
  const res = await listFamilies()
  if (res.success) {
    families.value = res.data || []
    if (!selectedFamily.value && families.value.length) await selectFamily(families.value[0])
  }
}

async function loadInvitations() {
  const res = await listFamilyInvitations()
  if (res.success) invitations.value = res.data || []
}

async function selectFamily(item) {
  selectedFamily.value = item
  const memberRes = await listFamilyMembers(item.id, { includeInvited: true })
  if (memberRes.success) members.value = memberRes.data || []

  // 切换家庭时加载分类数据
  if (item) {
    await Promise.all([loadGroups(item.id), loadAllCategories()])
  }
}

async function createNewFamily() {
  const res = await createFamily(familyForm)
  if (res.success) {
    ElMessage.success(t('family.messages.created'))
    dialogVisible.value = false
    selectedFamily.value = res.data
    await loadFamilies()
  } else ElMessage.error(res.error || t('family.messages.actionFailed'))
}

async function inviteMember() {
  if (!selectedFamily.value) return
  const res = await inviteFamilyMember(selectedFamily.value.id, inviteForm)
  if (res.success) {
    ElMessage.success(t('family.messages.invited'))
    inviteForm.account = ''
    inviteForm.displayName = ''
    await selectFamily(selectedFamily.value)
  } else ElMessage.error(res.error || t('family.messages.actionFailed'))
}

async function acceptInvite(item) {
  const res = await acceptFamilyInvitation(item.id)
  if (res.success) { ElMessage.success(t('family.messages.accepted')); await refreshAll() }
  else ElMessage.error(res.error || t('family.messages.actionFailed'))
}

async function rejectInvite(item) {
  const res = await rejectFamilyInvitation(item.id)
  if (res.success) { ElMessage.success(t('family.messages.rejected')); await refreshAll() }
  else ElMessage.error(res.error || t('family.messages.actionFailed'))
}

// 分类组方法
async function loadGroups(familyId) {
  loadingGroups.value = true
  try {
    const res = await listFamilyCategoryGroups(familyId, {})
    if (res.success) {
      groups.value = res.data || []
      // 加载每个组的成员
      for (const g of groups.value) {
        loadGroupMembers(familyId, g.id)
      }
    }
  } catch (e) {
    console.error('加载分组失败', e)
  } finally {
    loadingGroups.value = false
  }
}

async function loadGroupMembers(familyId, groupId) {
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

async function loadAllCategories() {
  try {
    const res = await listCategories()
    if (res.success) {
      categories.value = res.data || []
    }
  } catch (e) {
    console.error('加载分类失败', e)
  }
}

function openCreateGroup() {
  groupDialog.data = {
    name: '',
    type: groupTypeFilter.value,
    icon: '📁',
    color: '#f97316',
    sort: 0
  }
  groupDialog.visible = true
}

function openEditGroup(group) {
  groupDialog.data = { ...group }
  groupDialog.visible = true
}

async function saveGroup() {
  if (!selectedFamily.value) return
  try {
    if (groupDialog.data.id) {
      await updateFamilyCategoryGroup(selectedFamily.value.id, groupDialog.data.id, groupDialog.data)
      ElMessage.success('更新成功')
    } else {
      await createFamilyCategoryGroup(selectedFamily.value.id, groupDialog.data)
      ElMessage.success('创建成功')
    }
    groupDialog.visible = false
    loadGroups(selectedFamily.value.id)
  } catch (e) {
    ElMessage.error(e.message || '操作失败')
  }
}

async function handleDeleteGroup(group) {
  if (!selectedFamily.value) return
  try {
    await ElMessageBox.confirm(`确定删除分组「${group.name}」吗？`, '提示', { type: 'warning' })
    await deleteFamilyCategoryGroup(selectedFamily.value.id, group.id)
    ElMessage.success('删除成功')
    loadGroups(selectedFamily.value.id)
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
  if (!selectedFamily.value || selectedMemberIds.value.length === 0) return
  try {
    await addCategoriesToGroup(selectedFamily.value.id, memberDialog.groupId, selectedMemberIds.value)
    ElMessage.success('添加成功')
    memberDialog.visible = false
    loadGroupMembers(selectedFamily.value.id, memberDialog.groupId)
  } catch (e) {
    ElMessage.error(e.message || '添加失败')
  }
}

async function removeMember(groupId, categoryId) {
  if (!selectedFamily.value) return
  try {
    await removeCategoryFromGroup(selectedFamily.value.id, groupId, categoryId)
    ElMessage.success('移除成功')
    loadGroupMembers(selectedFamily.value.id, groupId)
  } catch (e) {
    ElMessage.error(e.message || '移除失败')
  }
}

// 处理从分类页面跳转过来的参数
function handleRouteParams() {
  // 切换到指定的Tab
  if (route.query.tab === 'categoryGroups') {
    activeTab.value = 'categoryGroups'
  }

  // 选中指定的家庭
  if (route.query.familyId) {
    const familyId = parseInt(route.query.familyId)
    const fam = families.value.find(f => f.id === familyId)
    if (fam) {
      selectFamily(fam)
    }
  }

  // 自动打开创建分组的弹窗
  if (route.query.createNew === 'true') {
    // 延迟一下确保数据加载完成
    setTimeout(() => {
      openCreateGroup()
      // 清除 query 参数，避免刷新页面时重复打开
      router.replace({ query: {} })
    }, 300)
  }
}

async function refreshAll() {
  await Promise.all([loadFamilies(), loadInvitations()])
  handleRouteParams()
}

// 监听路由变化，处理从其他页面跳转过来的情况
watch(() => route.query, handleRouteParams, { immediate: false })

onMounted(refreshAll)
onActivated(refreshAll)
</script>

<style scoped>
.family-page {
  max-width: 1200px;
  margin: 0 auto;
  color: #1e293b;
}

.reveal-block {
  animation: revealUp 480ms cubic-bezier(.16, 1, .3, 1) both;
}

.delay-1 {
  animation-delay: 90ms;
}

.family-list,
.family-detail {
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 1px 3px rgba(15, 23, 42, .1), 0 12px 30px rgba(15, 23, 42, .04);
}

.eyebrow {
  margin: 0 0 8px;
  color: #2f7d5c;
  font-size: 12px;
  font-weight: 800;
  letter-spacing: .08em;
  text-transform: uppercase;
}

.family-layout {
  display: grid;
  grid-template-columns: 320px minmax(0, 1fr);
  gap: 18px;
}

.family-list,
.family-detail {
  padding: 22px;
}

.family-list h2,
.family-detail h2 {
  margin: 0 0 14px;
}

.family-item,
.invite-card,
.member-card,
.stat-tile {
  display: block;
  width: 100%;
  border: 0;
  border-radius: 14px;
  padding: 14px;
  background: #f8faf7;
  margin-bottom: 10px;
  text-align: left;
}

.family-item {
  cursor: pointer;
}

.family-item.active {
  box-shadow: 0 0 0 2px rgba(47, 125, 92, .18);
}

.family-item strong,
.family-item span,
.member-card strong,
.member-card span,
.invite-card strong,
.invite-card span {
  display: block;
}

.family-item span,
.member-card span,
.invite-card span,
.muted {
  color: #64748b;
  font-size: 13px;
}

.primary-action,
.ghost-action,
.text-action,
.danger-action {
  min-height: 32px;
  border: 0;
  border-radius: 10px;
  padding: 0 12px;
  font-size: 13px;
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
.danger-action:active {
  transform: scale(0.96);
}

.primary-action {
  min-height: 36px;
  background: #2f7d5c;
  color: #fff;
  box-shadow: 0 8px 20px rgba(47, 125, 92, .18);
}

.invite-action,
.section-head .primary-action {
  min-height: 40px;
}

.ghost-action {
  background: #ffffff;
  color: #334155;
  box-shadow: inset 0 0 0 1px rgba(100, 116, 139, .16);
}

.text-action {
  background: #f4efe6;
  color: #6b5b49;
  box-shadow: none;
  padding: 0 12px;
}

.danger-action {
  background: rgba(239, 68, 68, .1);
  color: #ef4444;
}

.small {
  min-height: 32px;
  margin-right: 8px;
}

.tiny {
  min-height: 28px;
  padding: 0 10px;
  font-size: 12px;
}

.invite-form {
  display: grid;
  grid-template-columns: 1fr 1fr auto;
  gap: 10px;
  margin-bottom: 20px;
}

.member-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.full-width {
  width: 100%;
}

/* Tab 样式 */
.detail-tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 20px;
  border-bottom: 1px solid #e2e8f0;
  padding-bottom: 16px;
}

.tab-btn {
  padding: 8px 20px;
  border: 0;
  border-radius: 10px;
  background: transparent;
  color: #64748b;
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

.section-head {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
}

.section-head h2 {
  margin: 0;
}

/* 分类组样式 */
.type-tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 20px;
}

.type-tab {
  padding: 8px 18px;
  border-radius: 10px;
  border: none;
  background: #f8faf7;
  color: #64748b;
  cursor: pointer;
  font-size: 14px;
  font-weight: 600;
  transition: all 0.2s ease;
}

.type-tab:hover {
  background: #f0f4f0;
}

.type-tab.active {
  background: #2f7d5c;
  color: white;
  box-shadow: 0 4px 12px rgba(47, 125, 92, .15);
}

.group-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

@media (hover: hover) {
  .text-action:hover {
    background: #ece2d2;
    color: #4b3f33;
  }
}

.group-card {
  display: grid;
  gap: 0;
  background: #fff;
  border-radius: 16px;
  padding: 16px;
  box-shadow: 0 1px 3px rgba(15, 23, 42, .1), 0 12px 30px rgba(15, 23, 42, .04);
  transition-property: transform, box-shadow;
  transition-duration: 180ms;
}

@media (hover: hover) {
  .group-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 3px 8px rgba(15, 23, 42, .12), 0 16px 34px rgba(15, 23, 42, .06);
  }
}

.group-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 14px;
}

.group-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.group-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
}

.group-info h3 {
  font-size: 15px;
  font-weight: 600;
  margin: 0 0 4px 0;
}


.group-actions {
  display: flex;
  gap: 6px;
}

.group-members {
  min-height: 36px;
  padding: 12px 0;
  border-top: 1px solid #e8ede6;
  border-bottom: 1px solid #e8ede6;
  margin-bottom: 12px;
}

.member-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.empty-members {
  color: #94a3b8;
  font-size: 13px;
  text-align: center;
  padding: 8px 0;
}

.group-footer {
  display: flex;
  justify-content: flex-end;
}


.custom-icon-input {
  display: flex;
  align-items: center;
  gap: 10px;
}

.icon-preview {
  font-size: 24px;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f8faf7;
  border-radius: 8px;
  margin-top: 10px;
}

.category-selector {
  max-height: 360px;
  overflow-y: auto;
}

.category-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.category-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border-radius: 10px;
  background: #f8faf7;
  cursor: pointer;
  transition: background 0.2s;
}

.category-item:hover {
  background: #f0f4f0;
}

.category-name {
  flex: 1;
  font-size: 14px;
  color: #334155;
}

.empty-state-inline {
  grid-column: 1 / -1;
  text-align: center;
  padding: 40px 20px;
  background: #f8faf7;
  border-radius: 14px;
}

.empty-state-inline p {
  color: #64748b;
  margin: 12px 0;
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
  .family-layout {
    grid-template-columns: 1fr;
  }

  .invite-form {
    grid-template-columns: 1fr;
  }

  .page-hero {
    flex-direction: column;
  }

  .group-grid {
    grid-template-columns: 1fr;
  }
}
</style>