<!-- 首页页面 -->
<template>
  <div class="home-container" :class="{ 'sidebar-collapsed': sidebar.collapse }">
    <!-- 头部 -->
    <header class="app-header">
      <Header/>
    </header>
    
    <!-- 侧边栏 -->
    <aside class="app-sidebar" :class="{ 'collapsed': sidebar.collapse, 'show': sidebar.showMobile }">
      <Sidebar/>
    </aside>

    <!-- 主内容区域 -->
    <main class="app-main">
      <!-- 标签栏 -->
      <div class="tags-wrapper">
        <Tags/>
      </div>
      
      <!-- 页面内容 -->
      <div class="page-content">
        <router-view v-slot="{ Component }">
          <transition name="page-fade" mode="out-in">
            <keep-alive>
              <component :is="Component" class="page-component"></component>
            </keep-alive>
          </transition>
        </router-view>
      </div>
    </main>
  </div>
</template>
<script setup>
import Header from "@/components/Header.vue";
import Sidebar from "@/components/Sidebar.vue";
import Tags from "@/components/Tags.vue";

import {useSidebarStore} from '@/stores/sidebar'
import {useTagsStore} from '@/stores/tags';

const sidebar = useSidebarStore()
const tags = useTagsStore()

</script>

<style scoped>
/* 主容器 - 使用CSS Grid布局 */
.home-container {
  min-height: 100vh;
  display: grid;
  grid-template-areas: 
    "sidebar header"
    "sidebar main";
  grid-template-rows: auto 1fr;
  grid-template-columns: 260px 1fr;
  background: linear-gradient(135deg, #f8fafc 0%, #e2e8f0 100%);
  transition: grid-template-columns 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

/* 侧边栏收起时的布局 */
.home-container.sidebar-collapsed {
  grid-template-columns: 70px 1fr;
}

/* 头部区域 */
.app-header {
  grid-area: header;
  z-index: 1000;
}

/* 侧边栏区域 */
.app-sidebar {
  grid-area: sidebar;
  z-index: 999;
}

/* 主内容区域 */
.app-main {
  grid-area: main;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(10px);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  margin: 0;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

/* 标签栏包装器 */
.tags-wrapper {
  flex-shrink: 0;
}

/* 页面内容区域 */
.page-content {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
  background: #ffffff;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  max-height: calc(100vh - 118px);
}

/* 页面组件动画 */
.page-component {
  animation: slideInUp 0.4s ease-out;
}

/* 页面切换动画 */
.page-fade-enter-active,
.page-fade-leave-active {
  transition: all 0.3s ease;
}

.page-fade-enter-from {
  opacity: 0;
  transform: translateY(20px);
}

.page-fade-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}

/* 组件入场动画 */
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

/* 滚动条优化 */
.page-content::-webkit-scrollbar {
  width: 6px;
}

.page-content::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.1);
  border-radius: 3px;
}

.page-content::-webkit-scrollbar-thumb {
  background: var(--primary-color, #409eff);
  border-radius: 3px;
  transition: all 0.3s ease;
}

.page-content::-webkit-scrollbar-thumb:hover {
  background: var(--primary-dark, #337ecc);
}


/* 平板响应式 */
@media (max-width: 1024px) {
  .home-container {
    grid-template-columns: 260px 1fr;
  }
  
  .home-container.sidebar-collapsed {
    grid-template-columns: 70px 1fr;
  }
  
  .app-main {
    border-radius: 20px 0 0 0;
  }
  
  .page-content {
    padding: 20px;
    margin: 0 20px 20px 20px;
  }
}

/* 手机响应式 */
@media (max-width: 768px) {
  .home-container {
    grid-template-areas: 
      "header header"
      "main main";
    grid-template-rows: auto 1fr;
    grid-template-columns: 1fr;
  }
  
  .home-container.sidebar-collapsed {
    grid-template-columns: 1fr;
  }
  
  .app-sidebar {
    position: fixed;
    top: 0;
    left: 0;
    width: 260px;
    height: 100vh;
    transform: translateX(-100%);
    transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    z-index: 1001;
  }
  
  .app-sidebar.show {
    transform: translateX(0);
  }
  
  .app-main {
    border-radius: 16px 16px 0 0;
  }
  
  .page-content {
    padding: 16px;
    margin: 0 16px 16px 16px;
    border-radius: 0 0 12px 12px;
  }
}

/* 小屏手机响应式 */
@media (max-width: 480px) {
  .page-content {
    padding: 8px;
    margin: 0 8px 8px 8px;
  }
}

/* 加载状态 */
.page-content.loading {
  position: relative;
  overflow: hidden;
}

.page-content.loading::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.6), transparent);
  animation: loading 2s infinite;
}

@keyframes loading {
  0% { left: -100%; }
  100% { left: 100%; }
}
</style>