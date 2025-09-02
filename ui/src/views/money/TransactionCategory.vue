<template>
  <div class="category-management">
    <!-- 操作区域 -->
    <div class="action-section">
      <div class="action-header">
        <div class="section-title">
          <h3>类型管理</h3>
          <p>管理您的账单分类，收入和支出类型分别管理，更加清晰明了</p>
        </div>
        <div class="action-buttons">
          <el-button :icon="Refresh" @click="search" class="action-btn">
            刷新
          </el-button>
          <el-button type="primary" :icon="Plus" @click="addClick" class="action-btn">
            新增{{ currentTypeTab === 'income' ? '收入' : '支出' }}类型
          </el-button>
        </div>
      </div>
    </div>

    <!-- 分类标签页 -->
    <div class="category-tabs-section">
      <el-tabs 
        v-model="currentTypeTab" 
        type="border-card" 
        class="category-tabs"
        @tab-change="handleTypeTabChange"
      >
        <el-tab-pane name="expense">
          <template #label>
            <div class="tab-label expense-tab">
              <el-icon><Minus /></el-icon>
              <span>支出类型</span>
              <div class="tab-count">{{ expenseCount }}</div>
            </div>
          </template>
          
          <div class="tab-content-wrapper">
            <div class="table-container">
              <el-table 
                :data="filteredCategoryData" 
                stripe
                class="modern-table"
                :header-row-style="{ background: '#f8fafc' }"
                :empty-text="'暂无支出类型数据'"
                v-loading="loading"
              >
                <el-table-column prop="name" label="类型名称" min-width="150" show-overflow-tooltip>
                  <template #default="{ row }">
                    <div class="category-name">
                      <el-icon class="category-icon expense-icon"><Collection /></el-icon>
                      <span>{{ row.name }}</span>
                    </div>
                  </template>
                </el-table-column>
                
                <el-table-column prop="desc" label="类型描述" min-width="200" show-overflow-tooltip>
                  <template #default="{ row }">
                    <div class="category-desc">
                      {{ row.desc || '暂无描述' }}
                    </div>
                  </template>
                </el-table-column>
                
                <el-table-column prop="gmtCreate" label="创建时间" width="180" show-overflow-tooltip>
                  <template #default="{ row }">
                    <div class="create-time">
                      <el-icon class="time-icon"><Clock /></el-icon>
                      <span>{{ formatDate(row.gmtCreate) }}</span>
                    </div>
                  </template>
                </el-table-column>

                <el-table-column fixed="right" label="操作" width="160">
                  <template #default="scope">
                    <div class="table-actions">
                      <el-button link type="primary" :icon="Edit" @click="editClick(scope.row)" class="edit-btn">
                        编辑
                      </el-button>
                      <el-button link type="danger" :icon="Delete" @click="deleteClick(scope.row)" class="delete-btn">
                        删除
                      </el-button>
                    </div>
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </div>
        </el-tab-pane>
        
        <el-tab-pane name="income">
          <template #label>
            <div class="tab-label income-tab">
              <el-icon><TrendCharts /></el-icon>
              <span>收入类型</span>
              <div class="tab-count">{{ incomeCount }}</div>
            </div>
          </template>
          
          <div class="tab-content-wrapper">
            <div class="table-container">
              <el-table 
                :data="filteredCategoryData" 
                stripe
                class="modern-table"
                :header-row-style="{ background: '#f8fafc' }"
                :empty-text="'暂无收入类型数据'"
                v-loading="loading"
              >
                <el-table-column prop="name" label="类型名称" min-width="150" show-overflow-tooltip>
                  <template #default="{ row }">
                    <div class="category-name">
                      <el-icon class="category-icon income-icon"><Collection /></el-icon>
                      <span>{{ row.name }}</span>
                    </div>
                  </template>
                </el-table-column>
                
                <el-table-column prop="desc" label="类型描述" min-width="200" show-overflow-tooltip>
                  <template #default="{ row }">
                    <div class="category-desc">
                      {{ row.desc || '暂无描述' }}
                    </div>
                  </template>
                </el-table-column>
                
                <el-table-column prop="gmtCreate" label="创建时间" width="180" show-overflow-tooltip>
                  <template #default="{ row }">
                    <div class="create-time">
                      <el-icon class="time-icon"><Clock /></el-icon>
                      <span>{{ formatDate(row.gmtCreate) }}</span>
                    </div>
                  </template>
                </el-table-column>

                <el-table-column fixed="right" label="操作" width="160">
                  <template #default="scope">
                    <div class="table-actions">
                      <el-button link type="primary" :icon="Edit" @click="editClick(scope.row)" class="edit-btn">
                        编辑
                      </el-button>
                      <el-button link type="danger" :icon="Delete" @click="deleteClick(scope.row)" class="delete-btn">
                        删除
                      </el-button>
                    </div>
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>

    <!-- 编辑对话框 -->
    <el-dialog 
      v-model="dialogFormVisible" 
      width="580px"
      class="category-edit-dialog"
      :close-on-click-modal="false"
      :show-close="false"
      top="8vh"
    >
      <template #header>
        <div class="dialog-header">
          <div class="header-icon">
            <el-icon class="icon" :class="{ 'edit-mode': categoryInfo.id && categoryInfo.id !== '', 'add-mode': !categoryInfo.id || categoryInfo.id === '' }">
              <component :is="categoryInfo.id && categoryInfo.id !== '' ? 'Edit' : 'Plus'"/>
            </el-icon>
          </div>
          <div class="header-content">
            <h3 class="dialog-title">{{ categoryInfo.id && categoryInfo.id !== '' ? '编辑类型' : '新增类型' }}</h3>
            <p class="dialog-subtitle">{{ categoryInfo.id && categoryInfo.id !== '' ? '修改收支类型信息' : '创建新的收支类型' }}</p>
          </div>
          <el-button 
            @click="dialogBoxCancel" 
            class="close-btn" 
            :icon="Close" 
            circle 
            size="small"
          />
        </div>
      </template>

      <div class="dialog-body">
        <el-form 
          :model="categoryInfo" 
          :rules="formRules" 
          ref="formRef" 
          class="category-form"
          label-position="top"
        >
          <!-- 类型名称 -->
          <div class="form-section">
            <el-form-item label="类型名称" prop="name" class="name-field">
              <div class="input-wrapper">
                <el-input 
                  v-model="categoryInfo.name" 
                  placeholder="如：餐饮、交通、工资等"
                  class="modern-input"
                  size="large"
                  clearable
                >
                  <template #prefix>
                    <el-icon class="input-icon"><Collection /></el-icon>
                  </template>
                </el-input>
              </div>
            </el-form-item>
          </div>

          <!-- 收支类型选择 -->
          <div class="form-section type-section">
            <el-form-item label="收支类型" prop="type" class="type-field">
              <div class="type-selection-modern">
                <div class="type-options">
                  <div 
                    v-for="item in categoryTypes" 
                    :key="item.id"
                    class="type-option"
                    :class="{ 'active': categoryTypeValue.id === item.id, [`type-${item.id}`]: true }"
                    @click="selectType(item)"
                  >
                    <div class="option-icon-wrapper">
                      <el-icon class="option-icon">
                        <component :is="item.id === 1 ? 'TrendCharts' : 'Minus'"/>
                      </el-icon>
                    </div>
                    <div class="option-content">
                      <span class="option-title">{{ item.name }}</span>
                      <span class="option-desc">{{ item.id === 1 ? '资金流入' : '资金支出' }}</span>
                    </div>
                    <div class="option-check">
                      <el-icon v-if="categoryTypeValue.id === item.id"><Check /></el-icon>
                    </div>
                  </div>
                </div>
              </div>
            </el-form-item>
          </div>
          
          <!-- 类型描述 -->
          <div class="form-section">
            <el-form-item label="类型描述" prop="desc" class="desc-field">
              <div class="textarea-wrapper">
                <el-input 
                  v-model="categoryInfo.desc" 
                  type="textarea"
                  :rows="4"
                  placeholder="描述此类型的用途和说明（可选）"
                  maxlength="200"
                  show-word-limit
                  class="modern-textarea"
                  resize="none"
                />
              </div>
            </el-form-item>
          </div>
        </el-form>
      </div>
      
      <template #footer>
        <div class="dialog-footer-modern">
          <el-button 
            @click="dialogBoxCancel" 
            class="cancel-btn-modern"
            size="large"
            :icon="Close"
          >
            取消
          </el-button>
          <el-button 
            type="primary" 
            @click="dialogBoxConfirm" 
            :loading="saving" 
            class="confirm-btn-modern"
            size="large"
            :icon="saving ? '' : 'Check'"
          >
            <span>{{ saving ? '保存中...' : '确认保存' }}</span>
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import {reactive, ref, computed, defineEmits, onMounted} from "vue";
import {Plus, Refresh, Collection, Edit, Delete, Clock, TrendCharts, Minus, Close, Check} from '@element-plus/icons-vue';
import {allTransactionCategory, createTransactionCategory, updateTransactionCategory, deleteTransactionCategory} from "@/api/money/transaction/transactionCategory";
import {checkCategoryUsage} from "@/api/money/transaction/transaction";
import {ElMessage, ElMessageBox} from "element-plus";
import {TRANSACTION_TYPE_ID, TRANSACTION_TYPE_NAME} from "@/enums/transactionType";

const emit = defineEmits(['data-change'])

const categoryData = ref([])
const loading = ref(false)
const saving = ref(false)
const dialogFormVisible = ref(false)
const formRef = ref()
const currentTypeTab = ref('expense') // 默认显示支出类型

const categoryTypeValue = ref({
  id: '',
  name: ''
})

// 过滤后的数据（根据当前标签页显示）
const filteredCategoryData = computed(() => {
  const targetType = currentTypeTab.value === 'income' ? TRANSACTION_TYPE_ID.INCOME : TRANSACTION_TYPE_ID.EXPEND
  return categoryData.value.filter(item => item.type === targetType)
})

// 收入类型数量
const incomeCount = computed(() => {
  return categoryData.value.filter(item => item.type === TRANSACTION_TYPE_ID.INCOME).length
})

// 支出类型数量
const expenseCount = computed(() => {
  return categoryData.value.filter(item => item.type === TRANSACTION_TYPE_ID.EXPEND).length
})

const categoryInfo = reactive({
  id: '',
  name: '',
  desc: '',
  type: ''
})

// 表单验证规则
const formRules = {
  name: [
    { required: true, message: '请输入类型名称', trigger: 'blur' },
    { min: 1, max: 50, message: '类型名称长度在 1 到 50 个字符', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择收支类型', trigger: 'change' }
  ]
}

// 对话框标题
const dialogTitle = computed(() => {
  return categoryInfo.id ? '编辑类型' : '新增类型'
})

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return '未知'
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  })
}

const categoryTypes = [
  {id: TRANSACTION_TYPE_ID.EXPEND, name: TRANSACTION_TYPE_NAME[TRANSACTION_TYPE_ID.EXPEND]},
  {id: TRANSACTION_TYPE_ID.INCOME, name: TRANSACTION_TYPE_NAME[TRANSACTION_TYPE_ID.INCOME]}
]

const search = () => {
  doSearch()
}

const doSearch = async () => {
  try {
    loading.value = true
    const res = await allTransactionCategory();
    if (res.success) {
      let resList = res.data
      for (let one of resList) {
        one.typeName = TRANSACTION_TYPE_NAME[one.type] || "未知"
      }
      categoryData.value = resList
      updateCategoryCount()
    } else {
      categoryData.value = []
      ElMessage.error('查询类型失败')
    }
  } catch (error) {
    categoryData.value = []
    ElMessage.error('查询类型失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

const editClick = (row) => {
  categoryInfo.id = row.id
  categoryInfo.name = row.name
  categoryInfo.desc = row.desc
  categoryInfo.type = row.type
  
  // 从 categoryTypes 中找到对应的类型对象
  const typeObj = categoryTypes.find(type => type.id === row.type)
  categoryTypeValue.value = typeObj || categoryTypes[0]

  dialogFormVisible.value = true
}

const dialogBoxCancel = () => {
  closeDialogBox()
}

const dialogBoxConfirm = async () => {
  try {
    // 表单验证
    await formRef.value.validate()
    
    saving.value = true
    let param = {
      id: categoryInfo.id,
      name: categoryInfo.name.trim(),
      desc: categoryInfo.desc.trim(),
      type: parseInt(categoryInfo.type)
    }

    // 验证type是否有效
    if (!param.type || (param.type !== 1 && param.type !== 2)) {
      ElMessage.error('请选择有效的收支类型')
      return
    }

    if (param.id === '') {
      param.id = 0
    }
    
    let res;
    if (param.id === 0) {
      res = await createTransactionCategory(param);
    } else {
      res = await updateTransactionCategory(param);
    }
    
    if (res.success) {
      ElMessage.success('保存成功')
      closeDialogBox()
      await doSearch()
    } else {
      ElMessage.error('保存失败: ' + (res.error || '未知错误'))
    }
  } catch (error) {
    if (error !== false) { // 排除表单验证失败的情况
      console.error('保存类型失败:', error)
      let errorMsg = '未知错误'
      if (error.response && error.response.data && error.response.data.message) {
        errorMsg = error.response.data.message
      } else if (error.message) {
        errorMsg = error.message
      } else if (typeof error === 'string') {
        errorMsg = error
      }
      ElMessage.error('保存失败: ' + errorMsg)
    }
  } finally {
    saving.value = false
  }
}

const addClick = () => {
  categoryInfo.id = ''
  categoryInfo.name = ''
  categoryInfo.desc = ''
  categoryInfo.type = ''
  
  // 根据当前标签页设置默认类型
  const defaultType = currentTypeTab.value === 'income' 
    ? categoryTypes.find(t => t.id === TRANSACTION_TYPE_ID.INCOME)
    : categoryTypes.find(t => t.id === TRANSACTION_TYPE_ID.EXPEND)
  
  categoryTypeValue.value = defaultType || { id: '', name: '' }
  categoryInfo.type = categoryTypeValue.value.id
  dialogFormVisible.value = true
}

// 标签页切换处理
const handleTypeTabChange = (tabName) => {
  console.log('切换到:', tabName === 'income' ? '收入类型' : '支出类型')
}

const deleteClick = async (row) => {
  try {
    // 先检查分类是否被使用
    const usageRes = await checkCategoryUsage(row.id)
    if (!usageRes.success) {
      ElMessage.error('检查分类使用情况失败: ' + (usageRes.error || '未知错误'))
      return
    }
    
    let confirmMessage = '确定要删除此分类吗？'
    let confirmType = 'warning'
    
    if (usageRes.data.isUsed) {
      confirmMessage = `此分类已被 ${usageRes.data.count} 条交易记录使用，删除后这些记录将显示异常。确定要删除吗？`
      confirmType = 'error'
    }
    
    await ElMessageBox.confirm(confirmMessage, '确认删除', {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: confirmType,
    })
    
    const res = await deleteTransactionCategory(row.id)
    if (res.success) {
      ElMessage.success('删除成功')
      await doSearch()
    } else {
      ElMessage.error('删除失败: ' + (res.error || '未知错误'))
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除类型失败:', error)
      let errorMsg = '未知错误'
      if (error.response && error.response.data && error.response.data.message) {
        errorMsg = error.response.data.message
      } else if (error.message) {
        errorMsg = error.message
      } else if (typeof error === 'string') {
        errorMsg = error
      }
      ElMessage.error('删除失败: ' + errorMsg)
    }
  }
}

const closeDialogBox = () => {
  // 重置表单
  formRef.value?.resetFields()
  categoryInfo.id = ''
  categoryInfo.name = ''
  categoryInfo.desc = ''
  categoryInfo.type = ''
  categoryTypeValue.value = { id: '', name: '' }
  dialogFormVisible.value = false
  saving.value = false
}

// 处理类型选择变化
const handleTypeChange = (typeObj) => {
  categoryInfo.type = typeObj.id
}

// 选择类型方法
const selectType = (typeObj) => {
  categoryTypeValue.value = typeObj
  categoryInfo.type = typeObj.id
}

// 更新类型数据变化事件
const updateCategoryCount = () => {
  emit('data-change', categoryData.value.length)
}

// 组件挂载时获取数据
onMounted(() => {
  doSearch()
})


</script>

<style scoped>
/* 主容器样式 */
.category-management {
  background: transparent;
}

/* 操作区域样式 */
.action-section {
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  border: 1px solid #e2e8f0;
  border-radius: 16px;
  padding: 20px 24px;
  margin-bottom: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.action-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.section-title h3 {
  font-size: 20px;
  font-weight: 700;
  color: #1e293b;
  margin: 0 0 6px 0;
  background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.section-title p {
  font-size: 14px;
  color: #64748b;
  margin: 0;
}

.action-buttons {
  display: flex;
  gap: 12px;
}

.action-btn {
  border-radius: 10px;
  padding: 10px 20px;
  font-weight: 600;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.action-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}

/* 分类标签页样式 */
.category-tabs-section {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.06);
}

.category-tabs {
  background: transparent;
  border: none;
}

:deep(.category-tabs .el-tabs--border-card) {
  border: none;
  background: transparent;
  border-radius: 16px;
}

:deep(.category-tabs .el-tabs__header) {
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  border-bottom: 1px solid #e2e8f0;
  border-radius: 16px 16px 0 0;
  margin: 0;
  padding: 16px 24px;
}

:deep(.category-tabs .el-tabs__nav) {
  border: none;
  background: transparent;
}

:deep(.category-tabs .el-tabs__item) {
  border: 1px solid #e2e8f0;
  background: #ffffff;
  color: #64748b;
  font-weight: 600;
  padding: 12px 20px;
  margin-right: 8px;
  border-radius: 12px;
  transition: all 0.3s ease;
  height: auto;
  line-height: 1.4;
}

:deep(.category-tabs .el-tabs__item:hover) {
  background: linear-gradient(135deg, #f0f9ff 0%, #e0f2fe 100%);
  color: #0369a1;
  border-color: #3b82f6;
  transform: translateY(-1px);
}

:deep(.category-tabs .el-tabs__item.is-active) {
  background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
  color: #ffffff;
  border-color: #3b82f6;
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
}

:deep(.category-tabs .el-tabs__item.is-active:hover) {
  transform: none;
}

:deep(.category-tabs .el-tabs__content) {
  padding: 0;
}

.tab-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
}

.tab-label .el-icon {
  font-size: 16px;
}

.tab-count {
  background: rgba(255, 255, 255, 0.3);
  color: inherit;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 700;
  margin-left: 4px;
  min-width: 20px;
  text-align: center;
}

:deep(.category-tabs .el-tabs__item.is-active .tab-count) {
  background: rgba(255, 255, 255, 0.2);
  color: #ffffff;
}

.expense-tab {
  color: #16a34a;
}

.income-tab {
  color: #dc2626;
}

.tab-content-wrapper {
  padding: 24px 32px;
  background: #ffffff;
}

.table-container {
  overflow: hidden;
}

.modern-table {
  border-radius: 16px;
}

:deep(.el-table) {
  border-radius: 16px;
  overflow: hidden;
}

:deep(.el-table__header) {
  background: #f8fafc;
}

:deep(.el-table th) {
  background: #f8fafc !important;
  color: #1e293b;
  font-weight: 700;
  font-size: 14px;
  border-bottom: 1px solid #e2e8f0;
  padding: 16px 12px;
}

:deep(.el-table td) {
  padding: 16px 12px;
  border-bottom: 1px solid #f1f5f9;
}

:deep(.el-table__row:hover) {
  background: #f8fafc;
}

:deep(.el-table__row:hover td) {
  background: transparent;
}

/* 表格内容样式 */
.category-name {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: #1e293b;
}

.category-icon {
  font-size: 16px;
}

.income-icon {
  color: #ef4444;  /* 红色代表收入 */
}

.expense-icon {
  color: #10b981;  /* 绿色代表支出 */
}

.type-tag {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 8px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 600;
}

.type-tag.income {
  background: linear-gradient(135deg, #fee2e2 0%, #fecaca 100%);
  color: #ef4444;  /* 红色代表收入 */
  border: 1px solid #ef4444;
}

.type-tag.expense {
  background: linear-gradient(135deg, #dcfce7 0%, #bbf7d0 100%);
  color: #10b981;  /* 绿色代表支出 */
  border: 1px solid #10b981;
}

.type-icon {
  font-size: 14px;
}

.category-desc {
  color: #64748b;
  font-size: 14px;
}

.create-time {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #64748b;
}

.time-icon {
  font-size: 14px;
  color: #94a3b8;
}

.table-actions {
  display: flex;
  gap: 8px;
  justify-content: center;
}

.edit-btn {
  padding: 6px 12px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 600;
  transition: all 0.3s ease;
}

.edit-btn:hover {
  background: #3b82f6;
  color: #ffffff;
  transform: translateY(-1px);
}

.delete-btn {
  padding: 6px 12px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 600;
  transition: all 0.3s ease;
}

.delete-btn:hover {
  background: #ef4444;
  color: #ffffff;
  transform: translateY(-1px);
}

/* 新的现代化对话框样式 */
.category-edit-dialog {
  --primary-color: #3b82f6;
  --income-color: #ef4444;    /* 红色代表收入 */
  --expense-color: #10b981;   /* 绿色代表支出 */
  --warning-color: #f59e0b;
  --text-primary: #111827;
  --text-secondary: #6b7280;
  --bg-primary: #ffffff;
  --bg-secondary: #f9fafb;
  --border-color: #e5e7eb;
}

:deep(.category-edit-dialog .el-dialog) {
  border-radius: 24px;
  overflow: hidden;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
  border: 1px solid var(--border-color);
  margin: 0 auto;
  background: var(--bg-primary);
}

:deep(.category-edit-dialog .el-dialog__header) {
  padding: 0;
  margin: 0;
  border: none;
  background: transparent;
}

:deep(.category-edit-dialog .el-dialog__body) {
  padding: 0;
}

:deep(.category-edit-dialog .el-dialog__footer) {
  padding: 0;
  border: none;
  background: transparent;
}

/* 对话框头部样式 */
.dialog-header {
  display: flex;
  align-items: center;
  padding: 24px 32px 20px;
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  border-bottom: 1px solid var(--border-color);
  position: relative;
}

.header-icon {
  margin-right: 16px;
}

.header-icon .icon {
  width: 48px;
  height: 48px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  transition: all 0.3s ease;
}

.header-icon .icon.edit-mode {
  background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
  color: white;
  box-shadow: 0 8px 20px rgba(59, 130, 246, 0.3);
}

.header-icon .icon.add-mode {
  background: linear-gradient(135deg, #3b82f6 0%, #3b82f6 100%);
  color: white;
}

.header-content {
  flex: 1;
}

.dialog-title {
  font-size: 20px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 4px;
  line-height: 1.3;
}

.dialog-subtitle {
  font-size: 14px;
  color: var(--text-secondary);
  margin: 0;
  line-height: 1.4;
}

.close-btn {
  border: none;
  background: rgba(107, 114, 128, 0.1);
  color: var(--text-secondary);
  transition: all 0.3s ease;
}

.close-btn:hover {
  background: rgba(239, 68, 68, 0.1);
  transform: scale(1.1);
}

/* 对话框主体样式 */
.dialog-body {
  padding: 32px;
  background: var(--bg-primary);
}

.category-form {
  max-width: none;
}

.form-section {
  margin-bottom: 28px;
}

.form-section:last-child {
  margin-bottom: 0;
}

:deep(.category-form .el-form-item) {
  margin-bottom: 0;
}

:deep(.category-form .el-form-item__label) {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 8px;
  line-height: 1.4;
  padding: 0;
}

/* 输入框样式 */
.input-wrapper {
  position: relative;
}

:deep(.modern-input .el-input__wrapper) {
  border-radius: 16px;
  border: 2px solid var(--border-color);
  background: var(--bg-primary);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  padding: 14px 18px;
  min-height: 56px;
}

:deep(.modern-input .el-input__wrapper:hover) {
  border-color: var(--primary-color);
  box-shadow: 0 6px 20px rgba(59, 130, 246, 0.15);
  transform: translateY(-1px);
}

:deep(.modern-input .el-input__wrapper.is-focus) {
  border-color: var(--primary-color);
  box-shadow: 0 0 0 4px rgba(59, 130, 246, 0.15);
  transform: translateY(-1px);
}

:deep(.modern-input .el-input__inner) {
  font-size: 16px;
  font-weight: 500;
  color: var(--text-primary);
  line-height: 1.5;
}

:deep(.modern-input .el-input__inner::placeholder) {
  color: #9ca3af;
  font-weight: 400;
  font-size: 15px;
}

.input-icon {
  font-size: 18px;
  color: var(--text-secondary);
  margin-right: 8px;
}

/* 类型选择样式 */
.type-section {
  margin: 36px 0;
}

.type-selection-modern {
  border-radius: 16px;
  background: var(--bg-secondary);
  padding: 20px;
  border: 2px solid var(--border-color);
}

.type-options {
  display: flex;
  gap: 16px;
}

.type-option {
  flex: 1;
  padding: 20px 16px;
  border-radius: 16px;
  background: var(--bg-primary);
  border: 2px solid var(--border-color);
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  gap: 12px;
  position: relative;
  overflow: hidden;
}

.type-option.type-1:not(.active) {
  border-color: rgba(239, 68, 68, 0.2);
  background: rgba(239, 68, 68, 0.02);
}

.type-option.type-2:not(.active) {
  border-color: rgba(16, 185, 129, 0.2);
  background: rgba(16, 185, 129, 0.02);
}

.type-option::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, transparent 0%, rgba(59, 130, 246, 0.05) 100%);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.type-option:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
}

.type-option.type-1:hover:not(.active) {
  border-color: rgba(239, 68, 68, 0.4);
  background: rgba(239, 68, 68, 0.05);
  box-shadow: 0 8px 25px rgba(239, 68, 68, 0.2);
}

.type-option.type-2:hover:not(.active) {
  border-color: rgba(16, 185, 129, 0.4);
  background: rgba(16, 185, 129, 0.05);
  box-shadow: 0 8px 25px rgba(16, 185, 129, 0.2);
}

.type-option:hover::before {
  opacity: 1;
}

.type-option.type-1:hover:not(.active) .option-icon-wrapper {
  transform: scale(1.05);
}

.type-option.type-2:hover:not(.active) .option-icon-wrapper {
  background: rgba(16, 185, 129, 0.15);
  border-color: rgba(16, 185, 129, 0.3);
  transform: scale(1.05);
}

.type-option.active {
  border-color: var(--primary-color);
  box-shadow: 0 8px 25px rgba(59, 130, 246, 0.25);
  transform: translateY(-2px);
}

.type-option.active.type-1 {
  background: linear-gradient(135deg, rgba(239, 68, 68, 0.1) 0%, rgba(220, 38, 38, 0.05) 100%);
  border-color: var(--income-color);
  box-shadow: 0 8px 25px rgba(239, 68, 68, 0.25);
}

.type-option.active.type-2 {
  background: linear-gradient(135deg, rgba(16, 185, 129, 0.1) 0%, rgba(5, 150, 105, 0.05) 100%);
  border-color: var(--expense-color);
  box-shadow: 0 8px 25px rgba(16, 185, 129, 0.25);
}

.option-icon-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: 12px;
  background: var(--bg-secondary);
  transition: all 0.3s ease;
}

.type-option.type-1:not(.active) .option-icon-wrapper {
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.2);
}

.type-option.type-2:not(.active) .option-icon-wrapper {
  background: rgba(16, 185, 129, 0.1);
  border: 1px solid rgba(16, 185, 129, 0.2);
}

.type-option.active .option-icon-wrapper {
  transform: scale(1.1);
}

.type-option.active.type-1 .option-icon-wrapper {
  background: linear-gradient(135deg, var(--income-color) 0%, #dc2626 100%);
  color: white;
  box-shadow: 0 4px 12px rgba(239, 68, 68, 0.3);
}

.type-option.active.type-2 .option-icon-wrapper {
  background: linear-gradient(135deg, var(--expense-color) 0%, #059669 100%);
  color: white;
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.3);
}

.option-icon {
  font-size: 20px;
  color: var(--text-secondary);
}

.type-option:not(.active) .option-icon {
  color: #6b7280 !important;
}

.type-option.active .option-icon {
  color: white !important;
}

.type-option.type-1:not(.active) .option-icon {
  color: var(--income-color) !important;
}

.type-option.type-2:not(.active) .option-icon {
  color: var(--expense-color) !important;
}

/* 确保图标在所有状态下都有足够的对比度 */
:deep(.type-option .el-icon) {
  font-size: 20px !important;
}

:deep(.type-option.type-1:not(.active) .el-icon) {
  color: #ef4444 !important;
}

:deep(.type-option.type-2:not(.active) .el-icon) {
  color: #10b981 !important;
}

:deep(.type-option.active .el-icon) {
  color: white !important;
}

.option-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.option-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  line-height: 1.3;
}

.option-desc {
  font-size: 13px;
  color: var(--text-secondary);
  line-height: 1.4;
}

.type-option.type-1:not(.active) .option-title {
  color: var(--income-color);
}

.type-option.type-2:not(.active) .option-title {
  color: var(--expense-color);
}

.type-option.active .option-title {
  color: var(--text-primary);
}

.type-option.active .option-desc {
  color: var(--text-secondary);
}

.option-check {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  font-size: 14px;
  font-weight: bold;
  transition: all 0.3s ease;
}

.type-option.active.type-1 .option-check {
  background: var(--income-color);
  color: white;
}

.type-option.active.type-2 .option-check {
  background: var(--expense-color);
  color: white;
}

/* 文本域样式 */
.textarea-wrapper {
  position: relative;
}

:deep(.modern-textarea .el-textarea__inner) {
  border-radius: 16px;
  border: 2px solid var(--border-color);
  background: var(--bg-primary);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  padding: 18px;
  font-size: 15px;
  font-weight: 500;
  color: var(--text-primary);
  line-height: 1.6;
  min-height: 130px;
}

:deep(.modern-textarea .el-textarea__inner::placeholder) {
  color: #9ca3af;
  font-weight: 400;
  font-size: 14px;
}

:deep(.modern-textarea .el-textarea__inner:hover) {
  border-color: var(--primary-color);
  box-shadow: 0 6px 20px rgba(59, 130, 246, 0.15);
  transform: translateY(-1px);
}

:deep(.modern-textarea .el-textarea__inner:focus) {
  border-color: var(--primary-color);
  box-shadow: 0 0 0 4px rgba(59, 130, 246, 0.15);
  transform: translateY(-1px);
}

:deep(.modern-textarea .el-input__count) {
  background: rgba(107, 114, 128, 0.1);
  color: var(--text-secondary);
  border-radius: 8px;
  padding: 4px 8px;
  font-size: 12px;
  font-weight: 600;
  bottom: 8px;
  right: 12px;
}

/* 对话框底部样式 */
.dialog-footer-modern {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 16px;
  padding: 24px 32px;
  background: #ffffff;
  border-top: 1px solid var(--border-color);
}

.cancel-btn-modern {
  background: #f1f5f9 !important;
  border: 2px solid #cbd5e1 !important;
  color: #475569 !important;
  font-weight: 600;
  padding: 12px 24px;
  border-radius: 16px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  min-width: 100px;
  box-shadow: 0 3px 12px rgba(71, 85, 105, 0.15);
}

.cancel-btn-modern:focus {
  background: #f1f5f9 !important;
  border-color: #cbd5e1 !important;
  color: #475569 !important;
}

.cancel-btn-modern:hover {
  background: #e2e8f0 !important;
  border-color: #94a3b8 !important;
  color: #334155 !important;
  transform: translateY(-1px);
  box-shadow: 0 6px 20px rgba(71, 85, 105, 0.25);
}

.confirm-btn-modern {
  background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%) !important;
  border: 2px solid #3b82f6 !important;
  color: white !important;
  font-weight: 700;
  padding: 12px 32px;
  border-radius: 16px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  min-width: 140px;
  box-shadow: 0 6px 20px rgba(59, 130, 246, 0.4);
}

.confirm-btn-modern:focus {
  background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%) !important;
  border-color: #3b82f6 !important;
  color: white !important;
}

.confirm-btn-modern:hover {
  background: linear-gradient(135deg, #2563eb 0%, #1e40af 100%) !important;
  border-color: #2563eb !important;
  color: white !important;
  transform: translateY(-2px);
  box-shadow: 0 10px 32px rgba(59, 130, 246, 0.5);
}

.confirm-btn-modern:active {
  transform: translateY(-1px);
}

.confirm-btn-modern.is-loading {
  opacity: 0.8;
  transform: none;
}

.full-width {
  width: 100%;
}

:deep(.el-select) {
  width: 100%;
}

/* 响应式设计 - 新弹窗 */
@media (max-width: 768px) {
  :deep(.category-edit-dialog .el-dialog) {
    margin: 16px;
    width: calc(100vw - 32px) !important;
    max-height: calc(100vh - 32px);
  }

  .dialog-header {
    padding: 20px 24px 16px;
  }

  .header-icon .icon {
    width: 40px;
    height: 40px;
    font-size: 20px;
  }

  .dialog-title {
    font-size: 18px;
  }

  .dialog-subtitle {
    font-size: 13px;
  }

  .dialog-body {
    padding: 24px;
  }

  .form-section {
    margin-bottom: 20px;
  }

  .type-selection-modern {
    padding: 16px;
  }

  .type-options {
    flex-direction: column;
    gap: 12px;
  }

  .type-option {
    padding: 16px;
  }

  .option-icon-wrapper {
    width: 36px;
    height: 36px;
  }

  .option-title {
    font-size: 15px;
  }

  .option-desc {
    font-size: 12px;
  }

  .dialog-footer-modern {
    padding: 20px 24px;
    gap: 12px;
  }

  .cancel-btn-modern,
  .confirm-btn-modern {
    padding: 10px 20px;
    font-size: 15px;
    min-width: auto;
  }
}

/* 响应式设计 - 原有样式 */
@media (max-width: 768px) {
  .action-header {
    flex-direction: column;
    gap: 16px;
    text-align: center;
  }

  .action-buttons {
    justify-content: center;
  }

  .action-btn {
    padding: 8px 16px;
    font-size: 13px;
  }

  .section-title h3 {
    font-size: 18px;
  }

  .section-title p {
    font-size: 13px;
  }

  :deep(.el-table th),
  :deep(.el-table td) {
    padding: 12px 8px;
  }

  .type-tag {
    padding: 2px 6px;
    font-size: 12px;
  }

  .table-actions {
    flex-direction: column;
    gap: 4px;
  }

  .edit-btn, .delete-btn {
    padding: 4px 8px;
    font-size: 12px;
  }

  :deep(.el-dialog) {
    margin: 20px;
    width: calc(100% - 40px) !important;
  }

  .dialog-footer {
    padding: 16px 20px;
  }

  .cancel-btn, .confirm-btn {
    padding: 8px 16px;
    font-size: 14px;
  }
}

/* 暗色主题适配 */
@media (prefers-color-scheme: dark) {
  .action-section {
    background: linear-gradient(135deg, #1e293b 0%, #0f172a 100%);
    border: 1px solid #334155;
  }

  .section-title h3 {
    color: #f1f5f9;
  }

  .section-title p {
    color: #94a3b8;
  }

  .category-tabs-section {
    background: #1e293b;
    border: 1px solid #334155;
  }

  .tab-content-wrapper {
    background: #1e293b;
  }

  .category-name {
    color: #f1f5f9;
  }

  .category-desc {
    color: #94a3b8;
  }

  .create-time {
    color: #94a3b8;
  }

  .type-tag.income {
    background: linear-gradient(135deg, #7f1d1d 0%, #991b1b 100%);
    color: #f87171;
    border: 1px solid #ef4444;
  }

  .type-tag.expense {
    background: linear-gradient(135deg, #064e3b 0%, #065f46 100%);
    color: #34d399;
    border: 1px solid #10b981;
  }
}
</style>