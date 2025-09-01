<template>
  <el-dialog 
    v-model="props.showDialog" 
    @close="closeDialog" 
    draggable 
    title="账单导入" 
    width="1000px"
    top="6vh"
    class="import-dialog"
    :close-on-click-modal="false"
    center
    :modal="true"
    :append-to-body="true"
    :destroy-on-close="true"
  >
    <div class="import-container">
      <!-- 步骤指引 -->
      <div class="steps-container">
        <el-steps :active="currentStep" finish-status="success" align-center>
          <el-step title="选择平台" description="选择要导入的账单类型"></el-step>
          <el-step title="上传文件" description="选择对应格式的文件"></el-step>
          <el-step title="预览数据" description="检查并编辑导入内容"></el-step>
        </el-steps>
      </div>

      <!-- 主要内容区域 -->
      <div class="main-content">
        <!-- 步骤1: 平台选择 -->
        <div class="platform-section" v-show="currentStep === 0">
          <div class="section-header">
            <h3>选择账单平台</h3>
            <p>请选择要导入的账单类型，不同平台支持的文件格式不同</p>
          </div>
          
          <div class="platform-cards">
            <div 
              class="platform-card" 
              :class="{ active: selectedPlatform === 'wechat' }"
              @click="selectPlatform('wechat')"
            >
              <div class="card-icon">📱</div>
              <div class="card-title">微信账单</div>
              <div class="card-format">.xlsx 格式</div>
              <div class="card-desc">支持微信支付账单导入</div>
            </div>
            
            <div 
              class="platform-card" 
              :class="{ active: selectedPlatform === 'alipay' }"
              @click="selectPlatform('alipay')"
            >
              <div class="card-icon">💰</div>
              <div class="card-title">支付宝账单</div>
              <div class="card-format">.csv 格式</div>
              <div class="card-desc">支持支付宝账单导入</div>
            </div>
          </div>
        </div>

        <!-- 步骤2: 文件上传 -->
        <div class="upload-section" v-show="currentStep === 1">
          <div class="section-header">
            <h3>上传{{ platformName }}账单文件</h3>
            <p>{{ uploadTip }}</p>
          </div>
          
          <div class="upload-area">
            <el-upload
              drag
              action=""
              ref="upload"
              class="upload-component"
              :accept="currentAccept"
              :limit="1"
              :on-exceed="handleExceed"
              :auto-upload="false"
              :on-change="handleChange"
              :on-remove="handleRemove"
            >
              <div class="upload-content">
                <el-icon class="upload-icon">
                  <upload-filled/>
                </el-icon>
                <div class="upload-text">
                  <p>将文件拖拽到此处，或点击选择文件</p>
                  <p class="upload-format">支持格式：{{ currentAccept }}</p>
                </div>
              </div>
            </el-upload>
          </div>
        </div>

        <!-- 步骤3: 数据预览 -->
        <div class="preview-section" v-show="currentStep === 2">
          <div class="section-header">
            <h3>数据预览</h3>
            <p>请检查导入的数据，可以进行编辑后保存</p>
          </div>

          <!-- 数据表格 -->
          <div class="data-table-container">
            <vxe-table
              height="280"
              class="data-table"
              border
              stripe
              show-overflow
              :data="diaLogBoxTableData"
              :column-config="{resizable: true}"
              :edit-config="{trigger: 'click', mode: 'row'}"
              :tooltip-config="tooltipConfig"
            >
              <vxe-column type="seq" fixed="left" width="60" title="#"></vxe-column>
              <vxe-column field="description" title="描述" min-width="200" :edit-render="{}">
                <template #edit="{ row }">
                  <vxe-input v-model="row.description" type="text" placeholder="请输入描述"></vxe-input>
                </template>
              </vxe-column>

              <vxe-column field="amount" title="金额" width="120" :edit-render="{autofocus: '.vxe-input--inner'}">
                <template #edit="{ row }">
                  <vxe-input v-model="row.amount" type="float"></vxe-input>
                </template>
                <template #default="{ row }">
                  <span class="amount-text" :class="{ 'income': row.type === 1, 'expense': row.type === 2 }">
                    {{ row.type === 1 ? '+' : '-' }}{{ row.amount }}
                  </span>
                </template>
              </vxe-column>

              <vxe-column field="type" title="type" v-if="false"></vxe-column>
              <vxe-column field="typeName" width="100" title="收/支" :edit-render="{}">
                <template #default="{ row }">
                  <el-tag :type="row.type === 1 ? 'danger' : 'success'">
                    {{ getTransactionTypeNameById(row.type) }}
                  </el-tag>
                </template>
                <template #edit="{ row }">
                  <vxe-select v-model="row.type" transfer>
                    <vxe-option v-for="item in TRANSACTION_TYPES" :key="item.id" :value="item.id"
                                :label="item.name"></vxe-option>
                  </vxe-select>
                </template>
              </vxe-column>

              <vxe-column field="category" title="category" v-if="false"></vxe-column>
              <vxe-column field="categoryName" width="120" title="分类" :edit-render="{}">
                <template #default="{ row }">
                  <span>{{ getCategoryNameById(row.category, categories) }}</span>
                </template>
                <template #edit="{ row }">
                  <vxe-select v-model="row.category" transfer>
                    <vxe-option v-for="item in categories" :key="item.id" :value="item.id" :label="item.name"></vxe-option>
                  </vxe-select>
                </template>
              </vxe-column>

              <vxe-column field="account" title="account" v-if="false"></vxe-column>
              <vxe-column field="accountName" width="120" title="账户" :edit-render="{}">
                <template #default="{ row }">
                  <span>{{ getAccountNameById(row.account, accounts) }}</span>
                </template>
                <template #edit="{ row }">
                  <vxe-select v-model="row.account" transfer>
                    <vxe-option v-for="item in accounts" :key="item.id" :value="item.id" :label="item.name"></vxe-option>
                  </vxe-select>
                </template>
              </vxe-column>

              <vxe-column field="time" width="160" title="时间" :edit-render="{}">
                <template #edit="{ row }">
                  <vxe-input v-model="row.time" type="datetime" placeholder="请选择时间" transfer></vxe-input>
                </template>
              </vxe-column>

              <vxe-column title="操作" fixed="right" width="100">
                <template #default="{ row, column, rowIndex, columnIndex }">
                  <el-button type="danger" size="small" plain @click="dialogTableDelRow(rowIndex)">
                    删除
                  </el-button>
                </template>
              </vxe-column>
            </vxe-table>
          </div>
        </div>
      </div>

      <!-- 底部按钮区域 -->
      <div class="footer-actions">
        <div class="step-navigation">
          <el-button v-if="currentStep > 0" @click="prevStep" icon="ArrowLeft" class="nav-btn">
            上一步
          </el-button>
          
          <div class="action-buttons">
            <el-button 
              v-if="currentStep === 1 && curFile" 
              type="primary" 
              @click="submitUpload"
              :loading="uploading"
              class="main-btn"
            >
              解析文件
            </el-button>
            
            <el-button 
              v-if="currentStep === 2 && diaLogBoxTableData.length > 0" 
              type="success" 
              @click="dialogBoxSave"
              :loading="saving"
              class="main-btn"
            >
              保存数据 ({{ diaLogBoxTableData.length }}条)
            </el-button>
          </div>

          <el-button 
            v-if="currentStep < 2 && (currentStep === 0 ? selectedPlatform : (currentStep === 1 ? !curFile : curFile))" 
            type="primary" 
            @click="nextStep"
            icon="ArrowRight"
            class="nav-btn"
          >
            下一步
          </el-button>
        </div>
        
        <el-button @click="closeDialog" class="close-btn nav-btn">关闭</el-button>
      </div>
    </div>
  </el-dialog>
</template>

<script setup>
import {computed, onMounted, reactive, ref} from "vue";
import {UploadFilled} from "@element-plus/icons-vue";
import {ElMessage, ElMessageBox} from "element-plus";
import {postWeChatExcel, postAlipayExcel} from "@/api/money/transaction/transactionExcel";
import {getTransactionTypeNameById, TRANSACTION_TYPES} from "@/enums/transactionType";
import {getAccountNameById, getCategoryNameById} from "@/api/money/money";
import {allTransactionCategory} from "@/api/money/transaction/transactionCategory";
import {allTransactionAccount} from "@/api/money/transaction/transactionAccount";
import {batchPutTransaction} from "@/api/money/transaction/transaction";

const props = defineProps({
  showDialog: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['updateDialog']);

const upload = ref()
const curFile = ref()
const selectedPlatform = ref('')
const currentStep = ref(0)
const uploading = ref(false)
const saving = ref(false)

// 计算属性：根据选择的平台返回对应的accept类型
const currentAccept = computed(() => {
  if (selectedPlatform.value === 'wechat') {
    return '.xlsx'
  } else if (selectedPlatform.value === 'alipay') {
    return '.csv'
  }
  return ''
})

// 计算属性：平台名称
const platformName = computed(() => {
  if (selectedPlatform.value === 'wechat') {
    return '微信'
  } else if (selectedPlatform.value === 'alipay') {
    return '支付宝'
  }
  return ''
})

// 计算属性：上传提示
const uploadTip = computed(() => {
  if (selectedPlatform.value === 'wechat') {
    return '请选择微信导出的.xlsx格式账单文件'
  } else if (selectedPlatform.value === 'alipay') {
    return '请选择支付宝导出的.csv格式账单文件，确保编码为UTF-8'
  }
  return ''
})

// --- 弹出框表数据 START --- //

const diaLogBoxTableData = ref([])
// 光标移动展示内容
const tooltipConfig = reactive({
  showAll: true
})
// 下拉账户数据
const accounts = ref([])
// 下拉分类数据
const categories = ref([])

// --- 弹出框表数据 END --- //

// 初始化数据
onMounted(() => {
      allTransactionCategory().then(res => {
        categories.value = res.data
      });

      allTransactionAccount().then(res => {
        accounts.value = res.data
      })
    }
)


const closeDialog = () => {
  emit('updateDialog', false);
  // 重置所有状态
  resetDialog()
}

// 重置对话框状态
const resetDialog = () => {
  diaLogBoxTableData.value = []
  selectedPlatform.value = ''
  curFile.value = null
  currentStep.value = 0
  uploading.value = false
  saving.value = false
  if (upload.value) {
    upload.value.clearFiles()
  }
}

// 选择平台
const selectPlatform = (platform) => {
  selectedPlatform.value = platform
  // 清除之前选择的文件
  if (curFile.value) {
    curFile.value = null
    if (upload.value) {
      upload.value.clearFiles()
    }
  }
}

// 步骤导航
const nextStep = () => {
  if (currentStep.value < 2) {
    currentStep.value++
  }
}

const prevStep = () => {
  if (currentStep.value > 0) {
    currentStep.value--
  }
}

const handleChange = (file) => {
  if (file == null) {
    curFile.value = null
    return
  }
  let fileTemp = file.raw
  console.log("文件类型为：", fileTemp.type)
  curFile.value = fileTemp
}

const handleRemove = () => {
  curFile.value = null
}

const submitUpload = () => {
  if (!selectedPlatform.value) {
    ElMessage.warning("请先选择账单平台")
    return
  }
  if (curFile.value == null) {
    ElMessage.warning("请先选择一个文件上传")
    return
  }
  
  // 验证文件格式
  const fileName = curFile.value.name.toLowerCase()
  if (selectedPlatform.value === 'wechat' && !fileName.endsWith('.xlsx')) {
    ElMessage.error("微信账单仅支持XLSX格式文件")
    return
  }
  if (selectedPlatform.value === 'alipay' && !fileName.endsWith('.csv')) {
    ElMessage.error("支付宝账单仅支持CSV格式文件")
    return
  }
  
  uploading.value = true
  // 清空表格
  diaLogBoxTableData.value = []
  console.log("开始处理上传文件, 名称： " + curFile.value.name, curFile.value)
  let fileFormData = new FormData()
  fileFormData.append('file', curFile.value)
  
  // 根据平台选择不同的API
  const apiCall = selectedPlatform.value === 'wechat' 
    ? postWeChatExcel(fileFormData)
    : postAlipayExcel(fileFormData)
  
  apiCall.then(res => {
    if (res.success) {
      console.log(res.data)
      ElMessage.success(`${platformName.value}账单解析成功，共${res.data.length}条数据`)
      if (res.warnings && res.warnings.length > 0) {
        res.warnings.forEach(warning => {
          ElMessage.warning(warning)
        })
      }
      diaLogBoxTableData.value = res.data
      // 自动进入下一步
      nextStep()
    }
  }).catch(error => {
    ElMessage.error("上传解析失败：" + (error.response?.data?.error || error.message))
  }).finally(() => {
    uploading.value = false
  })
}

const handleExceed = () => {
  ElMessage.warning('最多只能上传一个文件！');
}

// 表格内的删除
const dialogTableDelRow = (rowIndex) => {
  diaLogBoxTableData.value.splice(rowIndex, 1)
}

// 保存按钮
const dialogBoxSave = () => {
  let params = diaLogBoxTableData.value
  if (params.length === 0) {
    ElMessage.warning('没有数据需要保存')
    return
  }

  let rowNum = 0
  for (let param of params) {
    rowNum++
    if (!param.account || !param.type || !param.category || !param.time) {
      ElMessage.error(`请完善第 ${rowNum} 行的数据`)
      return
    }
    param.amount = parseFloat(param.amount)
  }

  ElMessageBox.confirm(
    `确定要保存 ${params.length} 条交易记录吗？`,
    '确认保存',
    {
      confirmButtonText: '确定保存',
      cancelButtonText: '再检查一下',
      type: 'warning',
    }
  ).then(() => {
    saving.value = true
    // 批量保存
    batchPutTransaction(params).then(res => {
      if (!res.success) {
        ElMessage.error('保存失败，请重试')
        return
      }
      ElMessage.success(`保存成功！共保存 ${res.data} 条记录`)
      closeDialog()
    }).catch(error => {
      ElMessage.error('保存失败：' + (error.response?.data?.error || error.message))
    }).finally(() => {
      saving.value = false
    })
  }).catch(() => {
    // 用户取消，不做任何操作
  })
}

</script>

<style scoped>
/* 主容器 */
.import-container {
  padding: 0;
  background: #fafbfc;
  min-height: 420px;
  border-radius: 20px;
}

/* 步骤条 */
.steps-container {
  margin: 12px 16px 18px 16px;
  padding: 16px;
  background: #ffffff;
  border-radius: 16px;
  border: 1px solid #e2e8f0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.steps-container .el-steps {
  background: transparent;
}

.steps-container :deep(.el-step__title) {
  color: #1e293b !important;
  font-weight: 600;
}

.steps-container :deep(.el-step__description) {
  color: #64748b !important;
}

/* 主内容区域 */
.main-content {
  min-height: 280px;
  margin: 0 16px 16px 16px;
  background: #ffffff;
  border-radius: 16px;
  padding: 18px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  border: 1px solid #e2e8f0;
}

/* 节标题 */
.section-header {
  text-align: center;
  margin-bottom: 24px;
}

.section-header h3 {
  font-size: 20px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 6px 0;
}

.section-header p {
  color: #7f8c8d;
  font-size: 13px;
  margin: 0;
}

/* 平台选择卡片 */
.platform-cards {
  display: flex;
  justify-content: center;
  gap: 20px;
  flex-wrap: wrap;
}

.platform-card {
  flex: 0 0 220px;
  padding: 18px 16px;
  background: #ffffff;
  border: 2px solid #e2e8f0;
  border-radius: 16px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
}

.platform-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
  border-color: #3b82f6;
}

.platform-card.active {
  border-color: #3b82f6;
  background: #3b82f6;
  color: white;
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(59, 130, 246, 0.25);
}

.card-icon {
  font-size: 32px;
  margin-bottom: 8px;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.1));
}

.card-title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 6px;
}

.card-format {
  font-size: 14px;
  opacity: 0.8;
  margin-bottom: 12px;
  font-family: 'Monaco', 'Menlo', monospace;
}

.card-desc {
  font-size: 13px;
  opacity: 0.7;
  line-height: 1.4;
}

.platform-card.active .card-format,
.platform-card.active .card-desc {
  opacity: 0.9;
}

/* 上传区域 */
.upload-area {
  max-width: 600px;
  margin: 0 auto;
}

.upload-component {
  width: 100%;
}

.upload-component :deep(.el-upload-dragger) {
  border: 2px dashed #d1d5db;
  border-radius: 16px;
  background: #fafbfc;
  transition: all 0.3s ease;
  padding: 28px 20px;
}

.upload-component :deep(.el-upload-dragger:hover) {
  border-color: #3b82f6;
  background: #f8fafc;
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.1);
}

.upload-content {
  text-align: center;
}

.upload-icon {
  font-size: 40px;
  color: #3b82f6;
  margin-bottom: 12px;
  filter: drop-shadow(0 2px 8px rgba(59, 130, 246, 0.3));
}

.upload-text p {
  margin: 0 0 6px 0;
  font-size: 15px;
  color: #374151;
}

.upload-format {
  font-size: 13px;
  color: #6b7280;
  font-family: 'Monaco', 'Menlo', monospace;
  background: #f3f4f6;
  padding: 4px 12px;
  border-radius: 20px;
  display: inline-block;
}

/* 数据表格 */
.data-table-container {
  margin-top: 18px;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.data-table :deep(.vxe-table) {
  font-size: 14px;
}

.data-table :deep(.vxe-header--column) {
  background: #f8fafc;
  font-weight: 600;
  color: #374151;
}

.amount-text {
  font-weight: 600;
  font-family: 'Monaco', 'Menlo', monospace;
}

.amount-text.income {
  color: #ef4444;
}

.amount-text.expense {
  color: #10b981;
}

/* 底部操作区域 */
.footer-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 16px 12px 16px;
  border-top: 1px solid #e5e7eb;
  margin-top: 16px;
  background: #f9fafb;
  border-radius: 0 0 20px 20px;
}

.step-navigation {
  display: flex;
  align-items: center;
  gap: 20px;
  flex: 1;
}

.action-buttons {
  display: flex;
  justify-content: center;
  flex: 1;
}

.close-btn {
  margin-left: auto;
}

/* 按钮样式优化 */
.nav-btn {
  padding: 10px 20px;
  border-radius: 12px;
  font-weight: 500;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.main-btn {
  padding: 12px 32px;
  border-radius: 12px;
  font-weight: 600;
  font-size: 14px;
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.main-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(59, 130, 246, 0.4);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .import-container {
    padding: 16px;
  }
  
  .platform-cards {
    flex-direction: column;
    align-items: center;
  }
  
  .platform-card {
    flex: none;
    width: 100%;
    max-width: 320px;
  }
  
  .footer-actions {
    flex-direction: column;
    gap: 16px;
  }
  
  .step-navigation {
    width: 100%;
    justify-content: space-between;
  }
  
  .action-buttons {
    order: -1;
  }
}

/* 对话框样式 */
:deep(.import-dialog .el-dialog__body) {
  padding: 24px;
}

:deep(.import-dialog .el-dialog__header) {
  background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
  border: none;
  border-radius: 12px 12px 0 0;
  padding: 12px 20px;
  box-shadow: 0 2px 8px rgba(59, 130, 246, 0.2);
}

:deep(.import-dialog .el-dialog__title) {
  color: white;
  font-weight: 600;
  font-size: 16px;
}

:deep(.import-dialog .el-dialog__close) {
  color: rgba(255, 255, 255, 0.8);
  font-size: 18px;
}

:deep(.import-dialog .el-dialog__close:hover) {
  color: white;
  transform: scale(1.1);
}

:deep(.import-dialog .el-dialog) {
  border-radius: 20px;
  overflow: hidden;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
}

:deep(.import-dialog .el-dialog__body) {
  padding: 0;
  background: transparent;
}
</style>