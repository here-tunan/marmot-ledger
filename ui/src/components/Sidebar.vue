<template>
  <aside class="sidebar" :class="{ collapsed: sidebar.collapse, show: sidebar.showMobile }">
    <div class="sidebar-header">
      <div class="brand-mark">
        <img :src="marmotOne" alt="Marmot Ledger" width="42" height="42" />
      </div>
      <div class="brand-copy" v-show="!sidebar.collapse">
        <strong>Marmot Ledger</strong>
        <span>{{ t('app.tagline') }}</span>
      </div>
    </div>

    <nav class="sidebar-menu">
      <button
        v-for="item in items"
        :key="item.index"
        class="menu-item"
        :class="{ active: isActive(item.index) }"
        @click="navigateTo(item.index)"
      >
        <span class="menu-icon">
          <el-icon><component :is="item.icon"></component></el-icon>
        </span>
        <span class="menu-text" v-show="!sidebar.collapse">{{ t(item.titleKey) }}</span>
        <span class="menu-tooltip" v-show="sidebar.collapse">{{ t(item.titleKey) }}</span>
      </button>
    </nav>

    <div class="sidebar-footer" v-show="!sidebar.collapse">
      <span>Calm Marmot</span>
      <strong>© 2024</strong>
    </div>
  </aside>
</template>

<script setup>
import { useSidebarStore } from '@/stores/sidebar';
import { useRoute, useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import marmotOne from '../../../img/marmot-ledger-1.png';

const { t } = useI18n();
const route = useRoute();
const router = useRouter();
const sidebar = useSidebarStore();

const items = [
  {
    icon: 'House',
    index: '/dashboard',
    titleKey: 'routes.dashboard',
  },
  {
    icon: 'EditPen',
    index: '/record',
    titleKey: 'routes.record',
  },
  {
    icon: 'Document',
    index: '/records',
    titleKey: 'routes.records',
  },
  {
    icon: 'Wallet',
    index: '/accounts',
    titleKey: 'routes.accounts',
  },
  {
    icon: 'CollectionTag',
    index: '/categories',
    titleKey: 'routes.categories',
  },
  {
    icon: 'Box',
    index: '/buckets',
    titleKey: 'routes.buckets',
  },
  {
    icon: 'HomeFilled',
    index: '/family',
    titleKey: 'routes.family',
  },
  {
    icon: 'User',
    index: '/user',
    titleKey: 'routes.userCenter',
  },
];

const isActive = (path) => route.path === path;

const navigateTo = (path) => {
  if (path && path.startsWith('/')) {
    router.push(path);
    sidebar.closeMobile();
  }
};
</script>

<style scoped>
.sidebar {
  display: flex;
  flex-direction: column;
  height: 100vh;
  width: 260px;
  background: linear-gradient(180deg, #1f2933 0%, #111827 100%);
  color: rgba(255, 255, 255, 0.88);
  transition-property: width, transform;
  transition-duration: 240ms;
  transition-timing-function: cubic-bezier(0.16, 1, 0.3, 1);
  position: relative;
  user-select: none;
  box-shadow: 8px 0 28px rgba(15, 23, 42, 0.18);
}

.sidebar.collapsed {
  width: 70px;
}

.sidebar-header {
  min-height: 88px;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 20px;
}

.brand-mark {
  flex: 0 0 auto;
  width: 46px;
  height: 46px;
  display: grid;
  place-items: center;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.1);
  box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.08), 0 12px 24px rgba(0, 0, 0, 0.18);
}

.brand-mark img {
  width: 38px;
  height: 38px;
  object-fit: cover;
  border-radius: 12px;
}

.brand-copy strong,
.brand-copy span {
  display: block;
}

.brand-copy strong {
  font-size: 16px;
  letter-spacing: -0.012em;
}

.brand-copy span {
  margin-top: 4px;
  color: rgba(255, 255, 255, 0.52);
  font-size: 12px;
  white-space: nowrap;
}

.sidebar-menu {
  flex: 1;
  display: grid;
  align-content: start;
  gap: 8px;
  padding: 12px;
  overflow-y: auto;
}

.menu-item {
  position: relative;
  display: flex;
  align-items: center;
  min-height: 46px;
  width: 100%;
  border: 0;
  border-radius: 12px;
  padding: 0 14px;
  background: transparent;
  color: rgba(255, 255, 255, 0.62);
  cursor: pointer;
  text-align: left;
  transition-property: transform, background-color, color, box-shadow;
  transition-duration: 160ms;
  touch-action: manipulation;
}

.menu-item:active {
  transform: scale(0.96);
}

.menu-item:hover {
  background: rgba(255, 255, 255, 0.07);
  color: rgba(255, 255, 255, 0.92);
}

.menu-item.active {
  background: rgba(255, 255, 255, 0.12);
  color: #ffffff;
  box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.08), 0 10px 22px rgba(0, 0, 0, 0.12);
}

.menu-item.active::before {
  content: '';
  width: 6px;
  height: 6px;
  border-radius: 999px;
  background: #dce9df;
  position: absolute;
  left: 8px;
}

.menu-icon {
  display: inline-flex;
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
  font-weight: 700;
  white-space: nowrap;
}

.collapsed .sidebar-header {
  padding: 18px 12px;
  justify-content: center;
}

.collapsed .menu-item {
  justify-content: center;
  padding: 0;
}

.collapsed .menu-icon {
  margin: 0;
}

.menu-tooltip {
  position: absolute;
  left: 76px;
  top: 50%;
  transform: translateY(-50%);
  background: #1f2933;
  color: rgba(255, 255, 255, 0.9);
  padding: 8px 12px;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 700;
  white-space: nowrap;
  opacity: 0;
  visibility: hidden;
  transition-property: opacity, transform;
  transition-duration: 160ms;
  z-index: 1000;
  pointer-events: none;
  box-shadow: 0 10px 22px rgba(0, 0, 0, 0.22);
}

.collapsed .menu-item:hover .menu-tooltip {
  opacity: 1;
  visibility: visible;
  transform: translateY(-50%) translateX(4px);
}

.sidebar-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  color: rgba(255, 255, 255, 0.46);
  font-size: 12px;
}

.sidebar-footer strong {
  color: rgba(255, 255, 255, 0.72);
}

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
}

@media (prefers-reduced-motion: reduce) {
  .sidebar,
  .menu-item,
  .menu-tooltip {
    transition: none;
  }
}
</style>
