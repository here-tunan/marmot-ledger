<template>
  <div class="statistic-container">
    <!-- 顶部控制栏 -->
    <div class="control-panel">
      <div class="control-card">
        <!-- 时间选择器区域 -->
        <div class="time-selector">
          <div class="selector-header">
            <el-icon class="header-icon"><Calendar /></el-icon>
            <span class="header-title">时间范围</span>
          </div>
          <div class="time-tabs-wrapper">
            <el-tabs v-model="activeTab" @tab-change="tabChange" class="time-tabs">
              <el-tab-pane label="月账单" name="month">
                <div class="picker-wrapper">
                  <el-date-picker
                      v-model="monthSelect"
                      value-format="YYYY-MM"
                      type="month"
                      placeholder="选择月份"
                      @change="monthSelectClick"
                      class="modern-picker"
                  />
                </div>
              </el-tab-pane>

              <el-tab-pane label="年账单" name="year">
                <div class="picker-wrapper">
                  <el-date-picker
                      v-model="yearSelect"
                      type="year"
                      value-format="YYYY"
                      placeholder="选择年份"
                      @change="yearSelectClick"
                      class="modern-picker"
                  />
                </div>
              </el-tab-pane>

              <el-tab-pane label="自定义" name="custom">
                <div class="picker-wrapper">
                  <el-date-picker
                      v-model="customDateSelect"
                      type="daterange"
                      unlink-panels
                      range-separator="至"
                      start-placeholder="开始日期"
                      end-placeholder="结束日期"
                      value-format="YYYY-MM-DD"
                      @change="customSelectClick"
                      class="modern-picker"
                  />
                </div>
              </el-tab-pane>
            </el-tabs>
          </div>
        </div>

        <!-- 类型选择和差额显示 -->
        <div class="type-selector">
          <div class="selector-header">
            <el-icon class="header-icon"><TrendCharts /></el-icon>
            <span class="header-title">交易类型</span>
          </div>
          <div class="type-buttons">
            <el-radio-group v-model="typeId" class="modern-radio-group">
              <el-radio-button
                  @change="typeChange"
                  v-for="item in transactionTypes"
                  :key="item.id"
                  :label="item.id"
                  :class="{ 'income-btn': item.id === 1, 'expense-btn': item.id === 2 }"
              >
                <el-icon v-if="item.id === 1" class="btn-icon"><Plus /></el-icon>
                <el-icon v-if="item.id === 2" class="btn-icon"><Minus /></el-icon>
                {{ item.name }}
              </el-radio-button>
            </el-radio-group>
          </div>
          <div class="balance-display" v-if="theDiff">
            <span class="balance-label">收支差额</span>
            <span class="balance-amount">{{ theDiff }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 主要内容区域 -->
    <div class="main-content">
      <div v-show="!noData" class="data-content">
        <!-- 统计图表卡片 -->
        <div class="chart-section">
          <div class="section-card chart-card">
            <div class="card-header">
              <div class="header-left">
                <el-icon class="header-icon"><PieChart /></el-icon>
                <h3 class="card-title">{{ getTransactionTypeNameById(typeId) }}概览</h3>
              </div>
              <div class="chart-type-selector">
                <el-radio-group v-model="chartType" size="small">
                  <el-radio-button label="pie">饼图</el-radio-button>
                  <el-radio-button label="bar">柱图</el-radio-button>
                </el-radio-group>
              </div>
            </div>
            <div class="chart-body">
              <div class="chart-container">
                <div class="chart-wrapper" :class="{ 'with-description': analysisDesc }">
                  <div class="chart-content">
                    <div class="month-chart" id="month-chart"></div>
                  </div>
                  <div class="chart-description" v-if="analysisDesc">
                    <div class="description-card">
                      <h4 class="desc-title">
                        <el-icon><InfoFilled /></el-icon>
                        数据洞察
                      </h4>
                      <p class="desc-content">{{ analysisDesc }}</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 排行榜卡片 -->
        <div class="ranking-section">
          <div class="section-card ranking-card">
            <div class="card-header">
              <div class="header-left">
                <el-icon class="header-icon"><Trophy /></el-icon>
                <h3 class="card-title">{{ getTransactionTypeNameById(typeId) }}排行榜</h3>
              </div>
              <div class="ranking-stats">
                <span class="stat-item">
                  <span class="stat-label">总计</span>
                  <span class="stat-value">{{ tableDataTotal }}条</span>
                </span>
              </div>
            </div>
            <div class="card-body">
              <div class="table-container">
                <el-table :data="rankTableData" class="modern-table" stripe>
                  <el-table-column type="index" label="#" width="60" class-name="rank-column">
                    <template #default="{ $index }">
                      <div class="rank-number" :class="getRankClass($index)">
                        <el-icon v-if="$index === 0"><Medal /></el-icon>
                        <el-icon v-else-if="$index === 1"><Medal /></el-icon>
                        <el-icon v-else-if="$index === 2"><Medal /></el-icon>
                        <span v-else>{{ $index + 1 }}</span>
                      </div>
                    </template>
                  </el-table-column>
                  <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip>
                    <template #default="{ row }">
                      <div class="description-cell">
                        <span class="description-text">{{ row.description }}</span>
                      </div>
                    </template>
                  </el-table-column>
                  <el-table-column prop="amount" label="金额" width="140" align="right">
                    <template #default="{ row }">
                      <div class="amount-cell">
                        <span class="amount-symbol">¥</span>
                        <span class="amount-value">{{ formatAmount(row.amount) }}</span>
                      </div>
                    </template>
                  </el-table-column>
                  <el-table-column prop="time" label="时间" width="180">
                    <template #default="{ row }">
                      <div class="time-cell">
                        <el-icon class="time-icon"><Clock /></el-icon>
                        <span>{{ formatTime(row.time) }}</span>
                      </div>
                    </template>
                  </el-table-column>
                  <el-table-column prop="userName" label="出账人" width="120" v-if="props.familyId !== 0">
                    <template #default="scope">
                      <div class="user-cell">
                        <el-avatar size="small" class="user-avatar">
                          {{ getUserNameById(scope.row.userId, users)?.charAt(0) }}
                        </el-avatar>
                        <span class="user-name">{{ getUserNameById(scope.row.userId, users) }}</span>
                      </div>
                    </template>
                  </el-table-column>
                  <el-table-column prop="categoryName" label="分类" width="120">
                    <template #default="scope">
                      <el-tag class="category-tag" :type="getCategoryTagType(scope.row.category)">
                        {{ getCategoryNameById(scope.row.category, categories) }}
                      </el-tag>
                    </template>
                  </el-table-column>
                  <el-table-column prop="accountName" label="账户" width="120">
                    <template #default="scope">
                      <div class="account-cell">
                        <el-icon class="account-icon"><Wallet /></el-icon>
                        <span>{{ getAccountNameById(scope.row.account, accounts) }}</span>
                      </div>
                    </template>
                  </el-table-column>
                </el-table>
              </div>
              <div class="pagination-container">
                <el-pagination
                    v-model:current-page="currentPage"
                    :page-sizes="[10, 20, 50, 100]"
                    v-model:page-size="pageSize"
                    background layout="total, sizes, prev, pager, next, jumper"
                    :total="tableDataTotal"
                    @size-change="handleSizeChange"
                    @current-change="handleCurrentChange"
                    class="modern-pagination"
                />
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 空数据状态 -->
      <div v-show="noData" class="empty-state">
        <div class="empty-card">
          <el-empty class="modern-empty" :image-size="160">
            <template #image>
              <el-icon class="empty-icon"><DocumentRemove /></el-icon>
            </template>
            <template #description>
              <div class="empty-description">
                <h4>暂无数据</h4>
                <p>当前时间范围内没有{{ getTransactionTypeNameById(typeId) }}记录</p>
              </div>
            </template>
          </el-empty>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import {nextTick, onMounted, reactive, ref, shallowRef, watch} from "vue";
import * as echarts from "echarts";
import {queryTransaction, transactionAnalysis} from "@/api/money/transaction/transaction";
import {getTransactionTypeNameById, TRANSACTION_TYPE_ID, TRANSACTION_TYPES} from "@/enums/transactionType";
import {allTransactionCategory} from "@/api/money/transaction/transactionCategory";
import {allTransactionAccount} from "@/api/money/transaction/transactionAccount";
import {getAccountNameById, getCategoryNameById} from "@/api/money/money";
import {getLoginUserInfo, getUserNameById} from "@/api/user/user";
import {getFamily} from "@/api/family/family";
import {
  Calendar, TrendCharts, Plus, Minus, PieChart, InfoFilled,
  Trophy, Medal, Clock, Wallet, DocumentRemove
} from '@element-plus/icons-vue';


// 父组件传来的参数
const props = defineProps({
  familyId: {
    type: BigInt,
    default: 0,
  }
});

// 收入/支出
const transactionTypes = ref([])
const theDiff = ref('')
// 下拉账户数据
const accounts = ref([])
// 下拉分类数据
const categories = ref([])
// 用户信息(家庭下的所有成员)
const users = ref([])

const typeId = ref(TRANSACTION_TYPE_ID.EXPEND)

const activeTab = ref('month')
const monthSelect = ref('')
const yearSelect = ref('')
const customDateSelect = ref('')

const isInitChart = ref(false)

// 柱状图右边的描述
const analysisDesc = ref('')

// 表格相关
const rankTableData = ref()
const pageSize = ref(10)
const currentPage = ref(1)
const tableDataTotal = ref(110)

const showChart = shallowRef()

const showChartOption = reactive({
  series: [
    {
      type: 'pie',
      data: [],
      label: {
        normal: {
          formatter: '{b}:{d}% ( {c} )'
        }
      }
    }
  ]
})

const noData = ref(true)
const chartType = ref('pie') // 新增图表类型选择

const monthSelectClick = () => {
  initTableParam()
  let param = buildQueryParams();
  transactionsAnalysis(param)
  transactionsDetail(param)
}
const yearSelectClick = () => {
  initTableParam()
  let param = buildQueryParams();
  transactionsAnalysis(param)
  transactionsDetail(param)
}
const customSelectClick = () => {
  initTableParam()
  let param = buildQueryParams();
  transactionsAnalysis(param)
  transactionsDetail(param)
}
const tabChange = () => {
  initTableParam()
  let param = buildQueryParams();
  transactionsAnalysis(param)
  transactionsDetail(param)
}

const typeChange = () => {
  // ElMessage.success("收入/支出改变：" + typeId.value)
  let param = buildQueryParams();
  transactionsAnalysis(param)
  transactionsDetail(param)
}

const initTableParam = () => {
  currentPage.value = 1
  pageSize.value = 10
  rankTableData.value = []
  tableDataTotal.value = 0
  theDiff.value = ''
}

onMounted(() => {
  // 时间初始化
  const now = new Date();
  yearSelect.value = now.getFullYear().toString()
  monthSelect.value = now.getFullYear() + "-" + ((now.getMonth() + 1).toString().padStart(2, '0'));

  // 初始化
  transactionTypes.value = TRANSACTION_TYPES

  allTransactionCategory().then(res => {
    categories.value = res.data
  });

  allTransactionAccount().then(res => {
    accounts.value = res.data
  })

  // 相关用户信息
  if (props.familyId === 0) {
    // 用户信息
    getLoginUserInfo().then((res) => {
      if (res.success) {
        let userList = []
        userList.push(res.data)
        users.value = userList
      }
      let param = buildQueryParams();
      // 查询聚合信息
      transactionsAnalysis(param)
      // 查询排行具体信息
      transactionsDetail(param)
    })
  } else {
    // 获取家庭用户信息
    getFamily().then((res) => {
      if (res.success) {
        let userList = []
        for (let member of res.data.members) {
          userList.push({
            id: member.userId,
            name: member.name,
          })
        }
        users.value = userList
        let param = buildQueryParams();
        // 查询聚合信息
        transactionsAnalysis(param)
        // 查询排行具体信息
        transactionsDetail(param)
      }
    })
  }
})

// 分页按钮
const handleSizeChange = () => {
  // ElMessage.success("当前 size：" + pageSize.value)
  transactionsDetail(buildQueryParams())
}
const handleCurrentChange = () => {
  // ElMessage.success("当前 page：" + currentPage.value)
  transactionsDetail(buildQueryParams())
}

// 查询聚合信息
const transactionsAnalysis = (param) => {
  if (param.startTime === '' || param.endTime === '') {
    return
  }

  transactionAnalysis(param).then((res) => {
        let isNoData
        if (res.success) {
          let data = []
          if (res.data.total > 0) {
            buildTheDiff(res.data)
            buildAnalysisDesc(res.data)
            for (let categoryResult of res.data.categoryAggregations) {
              data.push({
                name: categoryResult.categoryName,
                value: categoryResult.amount
              })
            }
            data.sort((a, b) => b.value - a.value);
            isNoData = false
            showChartOption.series[0].data = data
          } else {
            isNoData = true
          }
        } else {
          isNoData = true
          // ElMessage.error("数据查询失败！")
        }
        noData.value = isNoData
      }
  )
}

// 构建收入和支出的差值
const buildTheDiff = (data) => {
  // 先乘100 再除100 避免精度的问题
  let diff = parseFloat(data.income) * 100 - parseFloat(data.expenditure) * 100
  theDiff.value = data.income + '-' + data.expenditure + '=¥' + (diff / 100).toFixed(2)
}

// 构建柱形图右边的描述信息
const buildAnalysisDesc = (data) => {
  // 用户聚合
  let userAggregations = data.userAggregations
  userAggregations.sort(function (a, b) {
    return b.amount - a.amount
  });
  // 分类聚合
  let categoryAggregation = data.categoryAggregations
  categoryAggregation.sort(function (a, b) {
    return b.amount - a.amount
  });

  let desc = '在这段时间里，'
  if (props.familyId === 0) {
    desc += '您共' + getTransactionTypeNameById(typeId.value) + data.total + '元。\n'
    desc += '其中' + categoryAggregation[0].categoryName + '占比(' + categoryAggregation[0].percent + '%)最高，共' + categoryAggregation[0].amount + '元。';

  } else {
    desc += '您的家庭共' + getTransactionTypeNameById(typeId.value) + data.total + '元。\n'
    desc += '其中' + categoryAggregation[0].categoryName + '占比(' + categoryAggregation[0].percent + '%)最高，共' + categoryAggregation[0].amount + '元。\n';
    for (let userAggregation of userAggregations) {
      desc += '\n' + getUserNameById(userAggregation.userId, users.value) + ': ' + userAggregation.amount + '元，' + '占比：' + userAggregation.percent + '%'
    }
  }
  analysisDesc.value = desc
}

// 查询具体信息
const transactionsDetail = (param) => {
  if (param.startTime === '' || param.endTime === '') {
    return
  }
  // api 调用
  queryTransaction(param).then((res) => {
    if (res.success) {
      tableDataTotal.value = res.total
      rankTableData.value = res.data
    }
  })
}

// 更新图表配置的通用方法
const updateChartOption = () => {
  if (!showChart.value || !showChartOption.series[0].data.length) return
  
  if (chartType.value === 'bar') {
    // 柱状图配置
    const barOption = {
      title: {
        show: false
      },
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'shadow'
        },
        formatter: '{b}: ¥{c}'
      },
      xAxis: {
        type: 'category',
        data: showChartOption.series[0].data.map(item => item.name),
        axisLabel: {
          interval: 0,
          rotate: 45,
          textStyle: {
            fontSize: 12
          }
        }
      },
      yAxis: {
        type: 'value',
        axisLabel: {
          formatter: '¥{value}'
        }
      },
      series: [{
        type: 'bar',
        data: showChartOption.series[0].data.map(item => item.value),
        itemStyle: {
          color: function(params) {
            const colors = ['#5470c6', '#91cc75', '#fac858', '#ee6666', '#73c0de', '#3ba272', '#fc8452', '#9a60b4', '#ea7ccc']
            return colors[params.dataIndex % colors.length]
          }
        },
        label: {
          show: false
        }
      }]
    }
    showChart.value.setOption(barOption, true)
  } else {
    // 饼图配置
    const pieOption = {
      ...showChartOption,
      tooltip: {
        trigger: 'item',
        formatter: '{b}: ¥{c} ({d}%)'
      }
    }
    showChart.value.setOption(pieOption, true)
  }
}

// 当数据变动时就展示数据
watch(showChartOption, () => {
      // 未初始化，且有数据时
      if (isInitChart.value === false && !noData.value) {
        // chart 初始化
        const chartDom = document.getElementById('month-chart')
        showChart.value = echarts.init(chartDom);
        isInitChart.value = true

        // 若dom尺寸变化，则resize
        const chartObserver = new ResizeObserver(() => {
          // setTimeOut 解决闪烁问题
          setTimeout(() => {
            showChart.value.resize();
          }, 0)
        });
        chartObserver.observe(chartDom);
      }
      
      // 使用 nextTick 方法延迟调用 resize() 【确保在 showChart dom 渲染后执行】
      nextTick(() => {
        if (showChart.value) {
          showChart.value.resize();
          updateChartOption();
        }
      });
    }
)

// 组装当前页面查询参数信息
const buildQueryParams = () => {
  let startTime
  let endTime
  if (activeTab.value === 'month') {
    if (monthSelect === '') {
      noData.value = true
      return
    }
    // 2023-10 + -01
    startTime = monthSelect.value + '-01'
    const [year, month] = startTime.split('-');
    const lastDay = new Date(year, month, 0).getDate();
    endTime = `${year}-${month}-${lastDay}`;
  }
  if (activeTab.value === 'year') {
    if (yearSelect === '') {
      noData.value = true
      return;
    }
    // 2023 + -01-01
    startTime = yearSelect.value + '-01-01'
    const lastDay = new Date(yearSelect.value, 12, 0).getDate();
    endTime = `${yearSelect.value}-12-${lastDay}`;
  }
  if (activeTab.value === 'custom') {
    if (customDateSelect.value === '') {
      noData.value = true
      return;
    }
    startTime = customDateSelect.value[0]
    endTime = customDateSelect.value[1]
  }
  let param = {}
  if (props.familyId !== '') {
    param.familyId = parseInt(props.familyId)
  }
  param.startTime = startTime
  param.endTime = endTime

  // 收入or支出
  param.type = typeId.value

  // 用户
  if (users.value.length > 0) {
    let userIds = []
    for (let user of users.value) {
      userIds.push(user.id)
    }
    param.userIds = userIds
  }

  // 排序
  param.amountDesc = true

  // 分页
  param.pageSize = pageSize.value
  param.pageIndex = currentPage.value

  return param
}

// 格式化金额显示
const formatAmount = (amount) => {
  return parseFloat(amount).toLocaleString('zh-CN', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  })
}

// 格式化时间显示
const formatTime = (time) => {
  const date = new Date(time)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  })
}

// 获取排名样式类
const getRankClass = (index) => {
  if (index === 0) return 'rank-first'
  if (index === 1) return 'rank-second'
  if (index === 2) return 'rank-third'
  return 'rank-normal'
}

// 获取分类标签类型
const getCategoryTagType = (categoryId) => {
  const hash = categoryId % 4
  const types = ['', 'success', 'warning', 'danger']
  return types[hash]
}

// 监听图表类型变化，更新图表配置
watch(chartType, () => {
  if (showChart.value && showChartOption.series[0].data.length > 0) {
    updateChartOption()
  }
})

// // 图表随着容器大小而改变
window.onresize = function () {
  if (showChart.value) {
    showChart.value.resize();
  }
};

// 侧边栏改变
// watch(() => useSidebarStore().collapse, (newValue, oldValue) => {
//   // ElMessage.success(newValue)
//   nextTick(()=>{
//     setTimeout(() => {
//       console.log(document.getElementById('month-chart').getBoundingClientRect())
//       showChart.value.resize();
//     }, 0);
//   })
// })

</script>

<style scoped>
/* CSS变量定义 */
.statistic-container {
  --primary-color: #3b82f6;
  --success-color: #10b981;
  --warning-color: #f59e0b;
  --danger-color: #ef4444;
  --text-primary: #1f2937;
  --text-secondary: #6b7280;
  --text-tertiary: #9ca3af;
  --bg-primary: #ffffff;
  --bg-secondary: #f9fafb;
  --bg-tertiary: #f3f4f6;
  --border-color: #e5e7eb;
  --border-light: #f3f4f6;
  --shadow-sm: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
  --shadow-md: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
  --shadow-lg: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
  --shadow-xl: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
  --radius-sm: 8px;
  --radius-md: 12px;
  --radius-lg: 16px;
  --radius-xl: 20px;
}

/* 主容器 */
.statistic-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  background-attachment: fixed;
  padding: 24px;
}

/* 控制面板 */
.control-panel {
  margin-bottom: 32px;
}

.control-card {
  background: var(--bg-primary);
  border-radius: var(--radius-xl);
  padding: 32px;
  box-shadow: var(--shadow-xl);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 40px;
}

/* 选择器通用样式 */
.selector-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;
}

.header-icon {
  font-size: 20px;
  color: var(--primary-color);
  background: rgba(59, 130, 246, 0.1);
  padding: 8px;
  border-radius: var(--radius-sm);
}

.header-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
}

/* 时间选择器 */
.time-selector {
  position: relative;
}

.time-tabs {
  position: relative;
}

.time-tabs :deep(.el-tabs__header) {
  margin-bottom: 20px;
  position: relative;
}

.time-tabs :deep(.el-tabs__nav-wrap) {
  position: relative;
  background: var(--bg-secondary);
  border-radius: var(--radius-lg) var(--radius-lg) 0 0;
  padding: 8px;
  box-shadow: inset 0 1px 3px rgba(0, 0, 0, 0.05);
  border-bottom: none;
}

.time-tabs :deep(.el-tabs__nav-scroll) {
  position: relative;
}

.time-tabs :deep(.el-tabs__nav) {
  background: transparent;
  border: none;
  display: flex;
  gap: 4px;
}

.time-tabs :deep(.el-tabs__item) {
  border: none !important;
  border-radius: var(--radius-sm) !important;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  font-weight: 500;
  padding: 10px 16px !important;
  margin: 0 !important;
  position: relative;
  background: transparent;
  color: var(--text-secondary) !important;
  font-size: 14px;
  min-height: auto !important;
  line-height: 1.4 !important;
}

.time-tabs :deep(.el-tabs__item:hover) {
  background: rgba(59, 130, 246, 0.1) !important;
  color: var(--primary-color) !important;
  transform: translateY(-1px);
}

.time-tabs :deep(.el-tabs__item.is-active) {
  background: var(--primary-color) !important;
  color: white !important;
  box-shadow: 0 2px 8px rgba(59, 130, 246, 0.3);
  transform: translateY(-1px);
}

.time-tabs :deep(.el-tabs__active-bar) {
  display: none !important;
}

.time-tabs :deep(.el-tabs__content) {
  padding: 0;
  margin: 0;
  background: transparent;
}

.time-tabs :deep(.el-tab-pane) {
  position: relative;
}

.time-tabs-wrapper {
  background: var(--bg-primary);
  border-radius: var(--radius-lg);
  border: 1px solid var(--border-light);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
}

.picker-wrapper {
  padding: 24px;
  background: var(--bg-primary);
  position: relative;
}

.picker-wrapper::before {
  content: '';
  position: absolute;
  top: 0;
  left: 24px;
  right: 24px;
  height: 1px;
  background: var(--border-light);
}

.modern-picker {
  width: 100%;
  position: relative;
  z-index: 1;
}

.modern-picker :deep(.el-input) {
  height: 48px;
}

.modern-picker :deep(.el-input__wrapper) {
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-sm);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 2px solid var(--border-color);
  background: var(--bg-primary);
  min-height: 48px;
  padding: 0 16px;
}

.modern-picker :deep(.el-input__wrapper:hover) {
  box-shadow: var(--shadow-md);
  border-color: var(--primary-color);
  transform: translateY(-1px);
}

.modern-picker :deep(.el-input__wrapper.is-focus) {
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1), var(--shadow-md);
}

.modern-picker :deep(.el-input__inner) {
  height: auto;
  line-height: 1.5;
  font-size: 15px;
  color: var(--text-primary);
  font-weight: 500;
}

.modern-picker :deep(.el-input__inner::placeholder) {
  color: var(--text-tertiary);
  font-weight: 400;
}

.modern-picker :deep(.el-input__suffix) {
  color: var(--text-secondary);
}

.modern-picker :deep(.el-input__suffix-inner) {
  display: flex;
  align-items: center;
}

/* 日期范围选择器特殊样式 */
.modern-picker :deep(.el-range-separator) {
  color: var(--text-secondary);
  font-weight: 500;
}

.modern-picker :deep(.el-range-input) {
  background: transparent;
  color: var(--text-primary);
  font-weight: 500;
}

.modern-picker :deep(.el-range-input::placeholder) {
  color: var(--text-tertiary);
  font-weight: 400;
}

/* 类型选择器 */
.type-selector {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.type-buttons {
  margin-bottom: 24px;
}

.modern-radio-group :deep(.el-radio-button) {
  margin-right: 12px;
}

.modern-radio-group :deep(.el-radio-button__inner) {
  border: 2px solid var(--border-color);
  border-radius: var(--radius-md);
  padding: 12px 24px;
  font-weight: 600;
  font-size: 14px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  gap: 8px;
}

.modern-radio-group :deep(.el-radio-button__inner:hover) {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.modern-radio-group :deep(.el-radio-button.is-active .el-radio-button__inner) {
  box-shadow: var(--shadow-lg);
  transform: translateY(-2px);
}

.income-btn :deep(.el-radio-button__inner) {
  border-color: var(--danger-color);
  color: var(--danger-color);
}

.income-btn :deep(.el-radio-button.is-active .el-radio-button__inner) {
  background: var(--danger-color);
  border-color: var(--danger-color);
  color: white;
}

.expense-btn :deep(.el-radio-button__inner) {
  border-color: var(--success-color);
  color: var(--success-color);
}

.expense-btn :deep(.el-radio-button.is-active .el-radio-button__inner) {
  background: var(--success-color);
  border-color: var(--success-color);
  color: white;
}

.btn-icon {
  font-size: 16px;
}

/* 余额显示 */
.balance-display {
  background: linear-gradient(135deg, var(--bg-secondary) 0%, var(--bg-tertiary) 100%);
  padding: 16px 20px;
  border-radius: var(--radius-md);
  border: 1px solid var(--border-light);
}

.balance-label {
  display: block;
  font-size: 14px;
  color: var(--text-secondary);
  margin-bottom: 8px;
}

.balance-amount {
  font-size: 20px;
  font-weight: 700;
  color: var(--text-primary);
  font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', monospace;
}

/* 主要内容区域 */
.main-content {
  display: grid;
  gap: 32px;
}

.data-content {
  display: grid;
  gap: 32px;
}

/* 卡片通用样式 */
.section-card {
  background: var(--bg-primary);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-xl);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  overflow: hidden;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px 32px;
  border-bottom: 1px solid var(--border-light);
  background: linear-gradient(135deg, var(--bg-secondary) 0%, rgba(59, 130, 246, 0.05) 100%);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.card-title {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
}

/* 图表卡片 */
.chart-card .card-body,
.chart-body {
  padding: 0;
}

.chart-type-selector :deep(.el-radio-group) {
  background: var(--bg-tertiary);
  border-radius: var(--radius-sm);
  padding: 2px;
}

.chart-type-selector :deep(.el-radio-button__inner) {
  border: none;
  background: transparent;
  border-radius: var(--radius-sm);
  padding: 6px 12px;
  font-size: 12px;
}

.chart-type-selector :deep(.el-radio-button.is-active .el-radio-button__inner) {
  background: var(--bg-primary);
  color: var(--primary-color);
  box-shadow: var(--shadow-sm);
}

.chart-container {
  padding: 32px;
}

.chart-wrapper {
  display: grid;
  gap: 32px;
}

.chart-wrapper.with-description {
  grid-template-columns: 2fr 1fr;
}

.chart-content {
  min-height: 400px;
}

.month-chart {
  width: 100%;
  height: 400px;
}

.chart-description {
  display: flex;
  align-items: flex-start;
}

.description-card {
  background: var(--bg-secondary);
  border-radius: var(--radius-lg);
  padding: 24px;
  border: 1px solid var(--border-light);
  width: 100%;
}

.desc-title {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 0 0 16px 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.desc-content {
  margin: 0;
  line-height: 1.6;
  color: var(--text-secondary);
  white-space: pre-wrap;
  font-size: 14px;
}

/* 排行榜卡片 */
.ranking-stats {
  display: flex;
  gap: 24px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.stat-label {
  font-size: 12px;
  color: var(--text-secondary);
}

.stat-value {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.card-body {
  padding: 32px;
}

/* 表格样式 */
.table-container {
  margin-bottom: 24px;
}

.modern-table {
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
}

.modern-table :deep(.el-table__header-wrapper) {
  background: var(--bg-secondary);
}

.modern-table :deep(.el-table__header th) {
  background: var(--bg-secondary) !important;
  border-bottom: 1px solid var(--border-light);
  font-weight: 600;
  color: var(--text-primary);
  padding: 16px 12px;
}

.modern-table :deep(.el-table__body tr:hover) {
  background: var(--bg-secondary) !important;
}

.modern-table :deep(.el-table td) {
  padding: 16px 12px;
  border-bottom: 1px solid var(--border-light);
}

/* 排名列样式 */
.rank-number {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  font-weight: 600;
}

.rank-first {
  background: linear-gradient(135deg, #ffd700 0%, #ffed4e 100%);
  color: #7c2d12;
}

.rank-second {
  background: linear-gradient(135deg, #c0c0c0 0%, #e5e7eb 100%);
  color: #374151;
}

.rank-third {
  background: linear-gradient(135deg, #cd7f32 0%, #f59e0b 100%);
  color: white;
}

.rank-normal {
  background: var(--bg-tertiary);
  color: var(--text-secondary);
}

/* 表格单元格样式 */
.description-cell {
  display: flex;
  align-items: center;
}

.description-text {
  font-weight: 500;
  color: var(--text-primary);
}

.amount-cell {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 4px;
  font-family: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', monospace;
}

.amount-symbol {
  color: var(--text-secondary);
  font-size: 12px;
}

.amount-value {
  font-weight: 600;
  color: var(--text-primary);
  font-size: 14px;
}

.time-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--text-secondary);
  font-size: 14px;
}

.time-icon {
  font-size: 14px;
  color: var(--text-tertiary);
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.user-avatar {
  width: 24px;
  height: 24px;
  font-size: 12px;
  font-weight: 600;
}

.user-name {
  font-size: 14px;
  color: var(--text-primary);
}

.category-tag {
  border-radius: var(--radius-sm);
  font-size: 12px;
  font-weight: 500;
}

.account-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--text-secondary);
  font-size: 14px;
}

.account-icon {
  font-size: 14px;
  color: var(--text-tertiary);
}

/* 分页样式 */
.pagination-container {
  display: flex;
  justify-content: center;
  padding-top: 24px;
  border-top: 1px solid var(--border-light);
}

.modern-pagination :deep(.el-pagination) {
  display: flex;
  align-items: center;
  gap: 8px;
}

.modern-pagination :deep(.el-pagination .btn-prev),
.modern-pagination :deep(.el-pagination .btn-next),
.modern-pagination :deep(.el-pagination .el-pager li) {
  border-radius: var(--radius-sm);
  border: 1px solid var(--border-color);
  transition: all 0.3s ease;
}

.modern-pagination :deep(.el-pagination .btn-prev:hover),
.modern-pagination :deep(.el-pagination .btn-next:hover),
.modern-pagination :deep(.el-pagination .el-pager li:hover) {
  transform: translateY(-1px);
  box-shadow: var(--shadow-sm);
}

/* 空数据状态 */
.empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
}

.empty-card {
  background: var(--bg-primary);
  border-radius: var(--radius-xl);
  padding: 48px;
  box-shadow: var(--shadow-xl);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  text-align: center;
}

.modern-empty {
  padding: 0;
}

.empty-icon {
  font-size: 80px;
  color: var(--text-tertiary);
  margin-bottom: 16px;
}

.empty-description h4 {
  margin: 0 0 8px 0;
  font-size: 18px;
  color: var(--text-primary);
}

.empty-description p {
  margin: 0;
  color: var(--text-secondary);
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .control-card {
    grid-template-columns: 1fr;
    gap: 32px;
  }
  
  .chart-wrapper.with-description {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .statistic-container {
    padding: 16px;
  }
  
  .control-card,
  .card-header,
  .card-body {
    padding: 20px;
  }
  
  .card-header {
    flex-direction: column;
    gap: 16px;
    align-items: flex-start;
  }
  
  .modern-table :deep(.el-table td),
  .modern-table :deep(.el-table th) {
    padding: 12px 8px;
  }
  
  .chart-container {
    padding: 20px;
  }
  
  .month-chart {
    height: 300px;
  }
}

/* 动画效果 */
.section-card {
  animation: slideInUp 0.6s ease-out;
}

@keyframes slideInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 深色模式支持 */
@media (prefers-color-scheme: dark) {
  .statistic-container {
    --text-primary: #f9fafb;
    --text-secondary: #d1d5db;
    --text-tertiary: #9ca3af;
    --bg-primary: #1f2937;
    --bg-secondary: #374151;
    --bg-tertiary: #4b5563;
    --border-color: #4b5563;
    --border-light: #374151;
  }
}
</style>