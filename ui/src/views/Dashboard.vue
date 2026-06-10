<template>
  <div class="dashboard">
    <div class="dashboard-header">
      <h1 class="page-title">Marmot Ledger</h1>
      <p class="page-subtitle">土拨鼠账本基础框架已就绪，下一步将基于新的财务领域模型重建账本能力。</p>
    </div>

    <el-row :gutter="24" class="model-cards">
      <el-col v-for="model in models" :key="model.name" :xs="24" :sm="12" :lg="6">
        <div class="model-card">
          <div class="model-icon">{{ model.icon }}</div>
          <h3>{{ model.name }}</h3>
          <p>{{ model.desc }}</p>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
const models = [
  {
    icon: '📌',
    name: 'Financial Event',
    desc: '记录发生了什么财务事件。',
  },
  {
    icon: '📒',
    name: 'Ledger Entry',
    desc: '记录事件引起的余额变化。',
  },
  {
    icon: '🪣',
    name: 'Bucket',
    desc: '资产、负债或虚拟资金容器。',
  },
  {
    icon: '🏦',
    name: 'Account',
    desc: '平台、机构或账户分组。',
  },
]
</script>

<style scoped>
.dashboard {
  max-width: 1200px;
  margin: 0 auto;
}

.dashboard-header {
  margin-bottom: 32px;
  text-align: center;
  padding: 40px 24px;
  background: linear-gradient(135deg, #ffffff 0%, #fafbfc 100%);
  border-radius: 20px;
  border: 1px solid #e1e8ed;
  box-shadow:
    0 2px 8px rgba(0, 0, 0, 0.03),
    0 8px 24px rgba(0, 0, 0, 0.02);
  position: relative;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.dashboard-header::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, transparent, rgba(59, 130, 246, 0.3), transparent);
  border-radius: 20px 20px 0 0;
}

.page-title {
  font-size: 36px;
  font-weight: 700;
  color: #1e293b;
  margin: 0 0 12px 0;
  background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.page-subtitle {
  font-size: 16px;
  color: #64748b;
  margin: 0;
  font-weight: 400;
}

.model-cards {
  margin-top: 0;
}

.model-card {
  min-height: 180px;
  padding: 24px;
  background: #ffffff;
  border: 1px solid #e1e8ed;
  border-radius: 18px;
  box-shadow: 0 8px 24px rgba(15, 23, 42, 0.04);
}

.model-icon {
  margin-bottom: 16px;
  font-size: 32px;
}

.model-card h3 {
  margin: 0 0 12px 0;
  color: #1e293b;
  font-size: 18px;
}

.model-card p {
  margin: 0;
  color: #64748b;
  line-height: 1.6;
}

@media (max-width: 768px) {
  .dashboard-header {
    margin-bottom: 24px;
  }

  .page-title {
    font-size: 28px;
  }

  .page-subtitle {
    font-size: 14px;
  }
}
</style>
