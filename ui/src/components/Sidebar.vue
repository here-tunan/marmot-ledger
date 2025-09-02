<template>
  <div class="sidebar" :class="{ 'collapsed': sidebar.collapse }">
    <!-- 侧边栏头部 -->
    <div class="sidebar-header">
      <div class="logo" v-show="!sidebar.collapse">
        <div class="logo-icon">📊</div>
        <div class="logo-text">My Life</div>
      </div>
      <div class="logo-mini" v-show="sidebar.collapse">
        <div class="logo-icon">📊</div>
      </div>
    </div>

    <!-- 菜单列表 -->
    <div class="sidebar-menu">
      <div class="menu-section">
        <template v-for="item in items" :key="item.index">
          <!-- 一级菜单项 -->
          <div 
            v-if="!item.subs" 
            class="menu-item"
            :class="{ 'active': isActive(item.index) }"
            @click="navigateTo(item.index)"
          >
            <div class="menu-icon">
              <el-icon><component :is="item.icon"></component></el-icon>
            </div>
            <div class="menu-text" v-show="!sidebar.collapse">{{ item.title }}</div>
            <div class="menu-tooltip" v-show="sidebar.collapse">{{ item.title }}</div>
          </div>

          <!-- 有子菜单的项目 -->
          <div v-else class="menu-group">
            <div 
              class="menu-group-title"
              :class="{ 'expanded': expandedGroups.includes(item.index) }"
              @click="toggleGroup(item.index)"
            >
              <div class="menu-icon">
                <el-icon><component :is="item.icon"></component></el-icon>
              </div>
              <div class="menu-text" v-show="!sidebar.collapse">
                {{ item.title }}
              </div>
              <div class="menu-arrow" v-show="!sidebar.collapse">
                <el-icon><ArrowRight /></el-icon>
              </div>
              <div class="menu-tooltip" v-show="sidebar.collapse">{{ item.title }}</div>
            </div>

            <!-- 子菜单 -->
            <div 
              class="submenu-wrapper"
              v-show="!sidebar.collapse && expandedGroups.includes(item.index)"
            >
              <div 
                v-for="subItem in item.subs" 
                :key="subItem.index"
                class="submenu-item"
                :class="{ 'active': isActive(subItem.index) }"
                @click="navigateTo(subItem.index)"
              >
                <div class="submenu-icon" v-if="subItem.icon">
                  <el-icon><component :is="subItem.icon"></component></el-icon>
                </div>
                <div class="submenu-text">{{ subItem.title }}</div>
              </div>
            </div>
          </div>
        </template>
      </div>
    </div>

    <!-- 侧边栏底部 -->
    <div class="sidebar-footer" v-show="!sidebar.collapse">
      <div class="footer-text">© 2024 My Life</div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useSidebarStore } from '@/stores/sidebar';
import { useRoute, useRouter } from 'vue-router';
import { ArrowRight } from '@element-plus/icons-vue';

const route = useRoute();
const router = useRouter();
const sidebar = useSidebarStore();

// 展开的菜单组
const expandedGroups = ref(['1', '2', '3']); // 默认展开所有组

const items = [
  {
    icon: 'House',
    index: '/dashboard',
    title: '系统首页',
  },
  {
    icon: 'Wallet',
    index: '1',
    title: '我的账本',
    subs: [
      {
        index: '/transaction-record',
        icon: 'EditPen',
        title: '我要记账',
      },
      {
        index: '/my-bill',
        icon: 'PieChart',
        title: '账单统计',
      },
      {
        index: '/management-panel',
        icon: 'Setting',
        title: '记账管理',
      }
    ],
  },
  {
    icon: 'TrendCharts',
    title: '健康生活',
    index: '2',
    subs: [
      {
        index: '/health',
        icon: 'DataAnalysis',
        title: '指标记录',
      },
      {
        index: '/health-board',
        icon: 'Monitor',
        title: '健康看板',
      },
      {
        index: '/sport',
        icon: 'Trophy',
        title: '运动健康',
      },
      {
        index: '/uuuu',
        icon: 'MoonNight',
        title: '睡眠健康',
      },
    ]
  },
  {
    icon: 'User',
    index: '3',
    title: '我的家庭',
    subs: [
      {
        index: '/family',
        icon: 'HomeFilled',
        title: '温馨小家',
      },
      {
        index: '/family-bill',
        icon: 'Files',
        title: '家庭账单',
      }
    ]
  },
];

// 判断菜单项是否激活
const isActive = (path) => {
  return route.path === path;
};

// 导航到指定路径
const navigateTo = (path) => {
  if (path && path.startsWith('/')) {
    router.push(path);
  }
};

// 切换菜单组展开状态
const toggleGroup = (groupIndex) => {
  if (sidebar.collapse) {
    // 折叠状态下点击直接展开侧边栏
    sidebar.handleCollapse();
    return;
  }
  
  const index = expandedGroups.value.indexOf(groupIndex);
  if (index > -1) {
    expandedGroups.value.splice(index, 1);
  } else {
    expandedGroups.value.push(groupIndex);
  }
};
</script>

<style scoped>
/* 侧边栏主容器 */
.sidebar {
  display: flex;
  flex-direction: column;
  height: 100vh;
  width: 260px;
  background: linear-gradient(180deg, 
    #1e293b 0%, 
    #0f172a 100%);
  border-right: 1px solid rgba(255, 255, 255, 0.1);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  user-select: none;
  box-shadow: 4px 0 15px rgba(0, 0, 0, 0.1);
}

.sidebar.collapsed {
  width: 70px;
}

/* 侧边栏头部 */
.sidebar-header {
  padding: 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  min-height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-mini {
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-icon {
  font-size: 24px;
  line-height: 1;
}

.logo-text {
  font-size: 18px;
  font-weight: 600;
  color: #f1f5f9;
}

/* 菜单区域 */
.sidebar-menu {
  flex: 1;
  padding: 16px 0;
  overflow-y: auto;
  overflow-x: hidden;
}

.sidebar-menu::-webkit-scrollbar {
  width: 4px;
}

.sidebar-menu::-webkit-scrollbar-track {
  background: transparent;
}

.sidebar-menu::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 2px;
}

.menu-section {
  padding: 0 12px;
}

/* 菜单项 */
.menu-item {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  margin: 2px 0;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
  position: relative;
  color: #94a3b8;
}

.menu-item:hover {
  background: rgba(255, 255, 255, 0.08);
  color: #f1f5f9;
  transform: translateX(4px);
}

.menu-item.active {
  background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
  color: #ffffff;
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
  transform: translateX(4px);
}

.menu-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  font-size: 18px;
  margin-right: 12px;
}

.menu-text {
  flex: 1;
  font-size: 14px;
  font-weight: 500;
  white-space: nowrap;
}

/* 菜单组 */
.menu-group {
  margin: 4px 0;
}

.menu-group-title {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  margin: 2px 0;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
  color: #94a3b8;
  position: relative;
}

.menu-group-title:hover {
  background: rgba(255, 255, 255, 0.08);
  color: #f1f5f9;
  transform: translateX(4px);
}

.menu-group-title.expanded {
  background: rgba(255, 255, 255, 0.05);
  color: #f1f5f9;
}

.menu-arrow {
  margin-left: auto;
  transition: transform 0.3s ease;
  font-size: 14px;
}

.menu-group-title.expanded .menu-arrow {
  transform: rotate(90deg);
}

/* 子菜单 */
.submenu-wrapper {
  margin: 4px 0;
  background: rgba(0, 0, 0, 0.1);
  border-radius: 8px;
  padding: 4px 0;
}

.submenu-item {
  display: flex;
  align-items: center;
  padding: 10px 20px 10px 48px;
  margin: 1px 8px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
  color: #94a3b8;
  font-size: 13px;
}

.submenu-item:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #f1f5f9;
  transform: translateX(2px);
}

.submenu-item.active {
  background: rgba(59, 130, 246, 0.2);
  color: #60a5fa;
  border-left: 3px solid #3b82f6;
}

.submenu-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  font-size: 14px;
  margin-right: 8px;
}

.submenu-text {
  flex: 1;
  white-space: nowrap;
}

/* 折叠状态 */
.collapsed .menu-item,
.collapsed .menu-group-title {
  justify-content: center;
  padding: 12px;
}

.collapsed .menu-icon {
  margin: 0;
  font-size: 20px;
}

.collapsed .submenu-wrapper {
  display: none;
}

/* 工具提示 */
.menu-tooltip {
  position: absolute;
  left: 75px;
  top: 50%;
  transform: translateY(-50%);
  background: rgba(0, 0, 0, 0.9);
  color: white;
  padding: 8px 12px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
  white-space: nowrap;
  opacity: 0;
  visibility: hidden;
  transition: all 0.3s ease;
  z-index: 1000;
  pointer-events: none;
}

.collapsed .menu-item:hover .menu-tooltip,
.collapsed .menu-group-title:hover .menu-tooltip {
  opacity: 1;
  visibility: visible;
}

/* 侧边栏底部 */
.sidebar-footer {
  padding: 16px 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  text-align: center;
}

.footer-text {
  font-size: 12px;
  color: #64748b;
  font-weight: 400;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .sidebar {
    position: fixed;
    top: 0;
    left: 0;
    z-index: 1000;
    transform: translateX(-100%);
    width: 260px;
  }
  
  .sidebar.show {
    transform: translateX(0);
  }
  
  .sidebar-header {
    padding: 16px;
    min-height: 70px;
  }
  
  .menu-item,
  .menu-group-title {
    padding: 10px 14px;
  }
  
  .submenu-item {
    padding: 8px 16px 8px 40px;
  }
}

@media (max-width: 480px) {
  .sidebar {
    width: 100vw;
  }
}

/* 动画效果 */
@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateX(-10px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

.menu-item,
.menu-group-title,
.submenu-item {
  animation: slideIn 0.3s ease-out;
}

/* 暗色主题适配 */
@media (prefers-color-scheme: dark) {
  .sidebar {
    background: linear-gradient(180deg, 
      #0f172a 0%, 
      #020617 100%);
  }
}
</style>
