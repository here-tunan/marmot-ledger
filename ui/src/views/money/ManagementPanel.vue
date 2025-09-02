<template>
  <div class="management-panel">
    <!-- 页面头部 -->
    <div class="panel-header">
      <div class="header-content">
        <div class="header-title">
          <h1>记账管理</h1>
          <p>管理您的账户信息和账单类型，让记账更加便捷</p>
        </div>
        <div class="header-stats">
          <div class="stat-card">
            <div class="stat-number">{{ accountCount }}</div>
            <div class="stat-label">账户数量</div>
          </div>
          <div class="stat-card">
            <div class="stat-number">{{ categoryCount }}</div>
            <div class="stat-label">类型数量</div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 标签页内容 -->
    <div class="tab-section">
      <el-tabs 
        type="border-card" 
        v-model="activeTab"
        class="modern-tabs"
        @tab-change="handleTabChange"
      >
        <el-tab-pane name="account">
          <template #label>
            <div class="tab-label">
              <el-icon><Wallet /></el-icon>
              <span>账户管理</span>
            </div>
          </template>
          <div class="tab-content">
            <TransactionAccount @data-change="updateAccountCount" />
          </div>
        </el-tab-pane>
        
        <el-tab-pane name="category">
          <template #label>
            <div class="tab-label">
              <el-icon><Collection /></el-icon>
              <span>类型管理</span>
            </div>
          </template>
          <div class="tab-content">
            <TransactionCategory @data-change="updateCategoryCount" />
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Wallet, Collection } from '@element-plus/icons-vue'
import TransactionAccount from './TransactionAccount.vue'
import TransactionCategory from './TransactionCategory.vue'
import { allTransactionAccount } from '@/api/money/transaction/transactionAccount'
import { allTransactionCategory } from '@/api/money/transaction/transactionCategory'

const activeTab = ref('account')
const accountCount = ref(0)
const categoryCount = ref(0)

// 更新账户数量
const updateAccountCount = async () => {
  try {
    const res = await allTransactionAccount()
    if (res.success) {
      accountCount.value = res.data.length
    }
  } catch (error) {
    console.error('Failed to update account count:', error)
  }
}

// 更新类型数量
const updateCategoryCount = async () => {
  try {
    const res = await allTransactionCategory()
    if (res.success) {
      categoryCount.value = res.data.length
    }
  } catch (error) {
    console.error('Failed to update category count:', error)
  }
}

// 标签页切换处理
const handleTabChange = (tabName) => {
  console.log('切换到:', tabName)
}

// 组件挂载时初始化数据
onMounted(() => {
  updateAccountCount()
  updateCategoryCount()
})
</script>

<style scoped>
.management-panel {
  padding: 0;
  min-height: 100vh;
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
}

/* 页面头部样式 */
.panel-header {
  background: linear-gradient(135deg, #ffffff 0%, #f8fafc 100%);
  border-bottom: 1px solid #e2e8f0;
  padding: 24px 32px;
  margin-bottom: 0;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  max-width: 1200px;
  margin: 0 auto;
}

.header-title h1 {
  font-size: 28px;
  font-weight: 700;
  color: #1e293b;
  margin: 0 0 8px 0;
  background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.header-title p {
  color: #64748b;
  font-size: 16px;
  margin: 0;
  font-weight: 400;
}

.header-stats {
  display: flex;
  gap: 24px;
}

.stat-card {
  background: linear-gradient(135deg, #ffffff 0%, #f8fafc 100%);
  border: 1px solid #e2e8f0;
  border-radius: 16px;
  padding: 20px 24px;
  text-align: center;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  transition: all 0.3s ease;
  min-width: 100px;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
}

.stat-number {
  font-size: 32px;
  font-weight: 700;
  color: #3b82f6;
  margin-bottom: 4px;
  font-family: 'Monaco', 'Menlo', monospace;
}

.stat-label {
  font-size: 14px;
  color: #64748b;
  font-weight: 500;
}

/* 标签页样式 */
.tab-section {
  padding: 32px;
  max-width: 1200px;
  margin: 0 auto;
}

.modern-tabs {
  background: transparent;
  border: none;
  border-radius: 20px;
  overflow: hidden;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
}

:deep(.el-tabs--border-card) {
  border: 1px solid #e2e8f0;
  background: #ffffff;
  border-radius: 20px;
}

:deep(.el-tabs__header) {
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  border-bottom: 1px solid #e2e8f0;
  border-radius: 20px 20px 0 0;
  margin: 0;
  padding: 16px 24px;
}

:deep(.el-tabs__nav-wrap) {
  padding: 0;
}

:deep(.el-tabs__nav) {
  border: none;
  background: transparent;
  border-radius: 12px;
  overflow: hidden;
}

:deep(.el-tabs__item) {
  border: 1px solid #e2e8f0;
  background: #ffffff;
  color: #64748b;
  font-weight: 600;
  padding: 12px 24px;
  margin-right: 8px;
  border-radius: 12px;
  transition: all 0.3s ease;
  height: auto;
  line-height: 1.4;
}

:deep(.el-tabs__item:hover) {
  background: linear-gradient(135deg, #f0f9ff 0%, #e0f2fe 100%);
  color: #0369a1;
  border-color: #3b82f6;
  transform: translateY(-1px);
}

:deep(.el-tabs__item.is-active) {
  background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
  color: #ffffff;
  border-color: #3b82f6;
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
}

:deep(.el-tabs__item.is-active:hover) {
  transform: none;
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

:deep(.el-tabs__content) {
  padding: 0;
  border-radius: 0 0 20px 20px;
}

.tab-content {
  padding: 24px 32px;
  background: #ffffff;
  border-radius: 0 0 20px 20px;
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .panel-header {
    padding: 20px 24px;
  }
  
  .header-content {
    flex-direction: column;
    gap: 20px;
    text-align: center;
  }
  
  .header-stats {
    justify-content: center;
  }
  
  .tab-section {
    padding: 24px;
  }
  
  .tab-content {
    padding: 20px 24px;
  }
}

@media (max-width: 768px) {
  .panel-header {
    padding: 16px 20px;
  }
  
  .header-title h1 {
    font-size: 24px;
  }
  
  .header-title p {
    font-size: 14px;
  }
  
  .header-stats {
    gap: 16px;
  }
  
  .stat-card {
    padding: 16px 20px;
    min-width: 80px;
  }
  
  .stat-number {
    font-size: 24px;
  }
  
  .stat-label {
    font-size: 12px;
  }
  
  .tab-section {
    padding: 16px;
  }
  
  .tab-content {
    padding: 16px 20px;
  }
  
  :deep(.el-tabs__header) {
    padding: 12px 16px;
  }
  
  :deep(.el-tabs__item) {
    padding: 10px 16px;
    margin-right: 6px;
    font-size: 13px;
  }
  
  .tab-label {
    gap: 6px;
  }
  
  .tab-label .el-icon {
    font-size: 14px;
  }
}

/* 深色主题适配 */
@media (prefers-color-scheme: dark) {
  .management-panel {
    background: linear-gradient(135deg, #0f172a 0%, #020617 100%);
  }
  
  .panel-header {
    background: linear-gradient(135deg, #1e293b 0%, #0f172a 100%);
    border-bottom: 1px solid #334155;
  }
  
  .header-title h1 {
    color: #f1f5f9;
  }
  
  .header-title p {
    color: #94a3b8;
  }
  
  .stat-card {
    background: linear-gradient(135deg, #1e293b 0%, #0f172a 100%);
    border: 1px solid #334155;
  }
  
  .stat-number {
    color: #60a5fa;
  }
  
  .stat-label {
    color: #94a3b8;
  }
}
</style>