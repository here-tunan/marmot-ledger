<template>
  <div class="dashboard">
    <div class="dashboard-header">
      <h1 class="page-title">系统首页</h1>
      <p class="page-subtitle">欢迎回到您的生活管理中心</p>
    </div>
    
    <el-row :gutter="24" class="dashboard-cards">
      <el-col :xs="24" :sm="12" :lg="12">
        <UserCard/>
      </el-col>
      <el-col :xs="24" :sm="12" :lg="12">
        <FamilyCard/>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import FamilyCard from "@/components/family/FamilyCard.vue";
import UserCard from "@/components/user/UserCard.vue";
import {ref} from "vue";

const hasFamily = ref(false)
</script>

<style scoped>
.dashboard {
  max-width: 1200px;
  margin: 0 auto;
}

.dashboard-header {
  margin-bottom: 32px;
  text-align: center;
  padding: 32px 24px;
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
  font-size: 32px;
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

.dashboard-cards {
  margin-top: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .dashboard-header {
    margin-bottom: 24px;
  }
  
  .page-title {
    font-size: 24px;
  }
  
  .page-subtitle {
    font-size: 14px;
  }
}

</style>

