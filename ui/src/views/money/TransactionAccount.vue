<template>
  <div class="account-management">
    <!-- 操作区域 -->
    <div class="action-section">
      <div class="action-header">
        <div class="section-title">
          <h3>账户管理</h3>
          <p>管理您的支付账户，如支付宝、微信支付、银行卡等</p>
        </div>
        <div class="action-buttons">
          <el-button :icon="Refresh" @click="search" class="action-btn">
            刷新
          </el-button>
          <el-button type="primary" :icon="Plus" @click="addClick" class="action-btn">
            新增账户
          </el-button>
        </div>
      </div>
    </div>

    <!-- 表格区域 -->
    <div class="table-section">
      <div class="table-container">
        <el-table 
          :data="accountData" 
          stripe
          class="modern-table"
          :header-row-style="{ background: '#f8fafc' }"
          empty-text="暂无账户数据"
          v-loading="loading"
        >
          <el-table-column prop="name" label="账户名称" min-width="150" show-overflow-tooltip>
            <template #default="{ row }">
              <div class="account-name">
                <el-icon class="account-icon"><Wallet /></el-icon>
                <span>{{ row.name }}</span>
              </div>
            </template>
          </el-table-column>
          
          <el-table-column prop="desc" label="账户描述" min-width="200" show-overflow-tooltip>
            <template #default="{ row }">
              <div class="account-desc">
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

    <!-- 编辑对话框 -->
    <el-dialog 
      v-model="dialogFormVisible" 
      width="560px"
      class="account-edit-dialog"
      :close-on-click-modal="false"
      :show-close="false"
      top="8vh"
    >
      <template #header>
        <div class="dialog-header">
          <div class="header-icon">
            <el-icon class="icon" :class="{ 'edit-mode': accountInfo.id && accountInfo.id !== '', 'add-mode': !accountInfo.id || accountInfo.id === '' }">
              <component :is="accountInfo.id && accountInfo.id !== '' ? 'Edit' : 'Plus'"/>
            </el-icon>
          </div>
          <div class="header-content">
            <h3 class="dialog-title">{{ accountInfo.id && accountInfo.id !== '' ? '编辑账户' : '新增账户' }}</h3>
            <p class="dialog-subtitle">{{ accountInfo.id && accountInfo.id !== '' ? '修改账户信息和描述' : '创建新的支付账户' }}</p>
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
          :model="accountInfo" 
          :rules="formRules" 
          ref="formRef" 
          class="account-form"
          label-position="top"
        >
          <!-- 账户名称 -->
          <div class="form-section">
            <el-form-item label="账户名称" prop="name" class="name-field">
              <div class="input-wrapper">
                <el-input 
                  v-model="accountInfo.name" 
                  placeholder="如：支付宝、微信、银行卡等"
                  class="modern-input"
                  size="large"
                  clearable
                >
                  <template #prefix>
                    <el-icon class="input-icon"><Wallet /></el-icon>
                  </template>
                </el-input>
              </div>
            </el-form-item>
          </div>
          
          <!-- 账户描述 -->
          <div class="form-section">
            <el-form-item label="账户描述" prop="desc" class="desc-field">
              <div class="textarea-wrapper">
                <el-input 
                  v-model="accountInfo.desc" 
                  type="textarea"
                  :rows="4"
                  placeholder="描述此账户的详细信息（可选）"
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
import {Plus, Refresh, Wallet, Edit, Delete, Clock, Close, Check} from '@element-plus/icons-vue';
import {allTransactionAccount, createTransactionAccount, updateTransactionAccount, deleteTransactionAccount} from "@/api/money/transaction/transactionAccount";
import {checkAccountUsage} from "@/api/money/transaction/transaction";
import {ElMessage, ElMessageBox} from "element-plus";

const emit = defineEmits(['data-change'])

const accountData = ref([])
const loading = ref(false)
const saving = ref(false)
const dialogFormVisible = ref(false)
const formRef = ref()

const accountInfo = reactive({
  id: '',
  name: '',
  desc: '',
})

// 表单验证规则
const formRules = {
  name: [
    { required: true, message: '请输入账户名称', trigger: 'blur' },
    { min: 1, max: 50, message: '账户名称长度在 1 到 50 个字符', trigger: 'blur' }
  ]
}

// 对话框标题
const dialogTitle = computed(() => {
  return accountInfo.id ? '编辑账户' : '新增账户'
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

const search = () => {
  doSearch()
}

const doSearch = async () => {
  try {
    loading.value = true
    const res = await allTransactionAccount();
    if (res.success) {
      accountData.value = res.data
      emit('data-change', res.data.length)
    } else {
      accountData.value = []
      ElMessage.error('查询账户失败')
    }
  } catch (error) {
    accountData.value = []
    ElMessage.error('查询账户失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

const editClick = (row) => {
  accountInfo.id = row.id
  accountInfo.name = row.name
  accountInfo.desc = row.desc

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
      id: accountInfo.id,
      name: accountInfo.name.trim(),
      desc: accountInfo.desc.trim(),
    }
    if (param.id === '') {
      param.id = 0
    }
    
    let res;
    if (param.id === 0) {
      res = await createTransactionAccount(param);
    } else {
      res = await updateTransactionAccount(param);
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
      ElMessage.error('保存失败: ' + error.message)
    }
  } finally {
    saving.value = false
  }
}

const addClick = () => {
  accountInfo.id = ''
  accountInfo.name = ''
  accountInfo.desc = ''
  dialogFormVisible.value = true
}

const deleteClick = async (row) => {
  try {
    // 先检查账户是否被使用
    const usageRes = await checkAccountUsage(row.id)
    if (!usageRes.success) {
      ElMessage.error('检查账户使用情况失败: ' + (usageRes.error || '未知错误'))
      return
    }
    
    let confirmMessage = '确定要删除此账户吗？'
    let confirmType = 'warning'
    
    if (usageRes.data.isUsed) {
      confirmMessage = `此账户已被 ${usageRes.data.count} 条交易记录使用，删除后这些记录将显示异常。确定要删除吗？`
      confirmType = 'error'
    }
    
    await ElMessageBox.confirm(confirmMessage, '确认删除', {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: confirmType,
    })
    
    const res = await deleteTransactionAccount(row.id)
    if (res.success) {
      ElMessage.success('删除成功')
      await doSearch()
    } else {
      ElMessage.error('删除失败: ' + (res.error || '未知错误'))
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败: ' + error.message)
    }
  }
}

const closeDialogBox = () => {
  // 重置表单
  formRef.value?.resetFields()
  accountInfo.id = ''
  accountInfo.name = ''
  accountInfo.desc = ''
  dialogFormVisible.value = false
  saving.value = false
}

// 组件挂载时获取数据
onMounted(() => {
  doSearch()
})


</script>

<style scoped>
/* 主容器样式 */
.account-management {
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

/* 表格区域样式 */
.table-section {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.06);
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
.account-name {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: #1e293b;
}

.account-icon {
  color: #3b82f6;
  font-size: 16px;
}

.account-desc {
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

/* 现代化账户对话框样式 */
.account-edit-dialog {
  --primary-color: #3b82f6;
  --wallet-color: #059669;
  --text-primary: #111827;
  --text-secondary: #6b7280;
  --bg-primary: #ffffff;
  --bg-secondary: #f9fafb;
  --border-color: #e5e7eb;
}

:deep(.account-edit-dialog .el-dialog) {
  border-radius: 24px;
  overflow: hidden;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
  border: 1px solid var(--border-color);
  margin: 0 auto;
  background: var(--bg-primary);
}

:deep(.account-edit-dialog .el-dialog__header) {
  padding: 0;
  margin: 0;
  border: none;
  background: transparent;
}

:deep(.account-edit-dialog .el-dialog__body) {
  padding: 0;
}

:deep(.account-edit-dialog .el-dialog__footer) {
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
  background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
  color: white;
  box-shadow: 0 8px 20px;
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
  color: #ef4444;
  transform: scale(1.1);
}

/* 对话框主体样式 */
.dialog-body {
  padding: 32px;
  background: var(--bg-primary);
}

.account-form {
  max-width: none;
}

.form-section {
  margin-bottom: 28px;
}

.form-section:last-child {
  margin-bottom: 0;
}

:deep(.account-form .el-form-item) {
  margin-bottom: 0;
}

:deep(.account-form .el-form-item__label) {
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
  border-color: var(--wallet-color);
  box-shadow: 0 6px 20px rgba(5, 150, 105, 0.15);
  transform: translateY(-1px);
}

:deep(.modern-input .el-input__wrapper.is-focus) {
  border-color: var(--wallet-color);
  box-shadow: 0 0 0 4px rgba(5, 150, 105, 0.15);
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
  color: var(--wallet-color);
  margin-right: 8px;
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
  border-color: var(--wallet-color);
  box-shadow: 0 6px 20px rgba(5, 150, 105, 0.15);
  transform: translateY(-1px);
}

:deep(.modern-textarea .el-textarea__inner:focus) {
  border-color: var(--wallet-color);
  box-shadow: 0 0 0 4px rgba(5, 150, 105, 0.15);
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
  color: white !important;
  font-weight: 700;
  padding: 12px 32px;
  border-radius: 16px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  min-width: 140px;
}

.confirm-btn-modern:focus {
  background: linear-gradient(135deg, var(--wallet-color) 0%, #047857 100%) !important;
  border-color: var(--wallet-color) !important;
  color: white !important;
}

.confirm-btn-modern:hover {
  color: white !important;
  transform: translateY(-2px);
}

.confirm-btn-modern:active {
  transform: translateY(-1px);
}

.confirm-btn-modern.is-loading {
  opacity: 0.8;
  transform: none;
}

/* 响应式设计 - 新弹窗 */
@media (max-width: 768px) {
  :deep(.account-edit-dialog .el-dialog) {
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

  .table-actions {
    flex-direction: column;
    gap: 4px;
  }

  .edit-btn, .delete-btn {
    padding: 4px 8px;
    font-size: 12px;
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

  .table-section {
    background: #1e293b;
    border: 1px solid #334155;
  }

  .account-name {
    color: #f1f5f9;
  }

  .account-desc {
    color: #94a3b8;
  }

  .create-time {
    color: #94a3b8;
  }
}
</style>