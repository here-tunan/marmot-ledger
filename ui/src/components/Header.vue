<!--首页的Header控件-->
<template>
  <header class="header-shell">
    <div class="header-left">
      <button class="collapse-btn" :aria-label="sidebar.collapse ? 'expand' : 'collapse'" @click="collapseChange">
        <el-icon v-if="sidebar.collapse"><Expand /></el-icon>
        <el-icon v-else><Fold /></el-icon>
      </button>
      <div class="route-title">
        <span>{{ t('header.currentWorkspace') }}</span>
        <strong>{{ routeTitle }}</strong>
      </div>
    </div>

    <div class="header-right">
      <div class="locale-switch" role="group" aria-label="language switch">
        <button :class="{ active: config.locale === 'zh-CN' }" @click="config.setLocale('zh-CN')">中文</button>
        <button :class="{ active: config.locale === 'en-US' }" @click="config.setLocale('en-US')">EN</button>
      </div>

      <el-dropdown class="user-name" trigger="click" @command="handleCommand" placement="bottom-end">
        <button class="user-trigger">
          <el-avatar class="user-avatar" :size="34" :src="avatarImg" />
          <span>{{ username || 'Marmot' }}</span>
          <el-icon><arrow-down /></el-icon>
        </button>
        <template #dropdown>
          <el-dropdown-menu>
            <a href="https://github.com/here-tunan/marmot-ledger" target="_blank">
              <el-dropdown-item>{{ t('header.repo') }}</el-dropdown-item>
            </a>
            <el-dropdown-item command="user">{{ t('header.profile') }}</el-dropdown-item>
            <el-dropdown-item divided command="loginOut">{{ t('header.logout') }}</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </header>
</template>

<script setup>
import { useSidebarStore } from '@/stores/sidebar';
import { ArrowDown, Expand, Fold } from '@element-plus/icons-vue';
import { computed, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { getLoginUserInfo } from '@/api/user/user';
import { removeToken } from '@/api/auth/auth';
import { ElMessage } from 'element-plus';
import { useUserStore } from '@/stores/user';
import { useConfigStore } from '@/stores/config';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const sidebar = useSidebarStore();
const config = useConfigStore();
const router = useRouter();
const route = useRoute();
const userStore = useUserStore();

const avatarImg = computed(() => userStore.avatar)
const username = computed(() => userStore.username)
const routeTitle = computed(() => t(route.meta.titleKey || 'routes.dashboard'))

onMounted(() => {
  if (document.body.clientWidth < 1000) {
    sidebar.collapse = true;
  }
});

onMounted(() => {
  getLoginUserInfo().then((res) => {
    if (res.success) {
      userStore.account = res.data.account
      userStore.username = res.data.name
      userStore.desc = res.data.desc
      userStore.avatar = res.data.avatar
      localStorage.setItem('avatar', res.data.avatar)
    } else {
      ElMessage.error(t('header.userInfoFailed'))
    }
  })
})

const handleCommand = (command) => {
  if (command === 'loginOut') {
    removeToken();
    localStorage.removeItem('avatar');
    ElMessage.success(t('header.logoutSuccess'))
    router.push('/login');
  } else if (command === 'user') {
    router.push('/user');
  }
};

const collapseChange = () => {
  if (window.innerWidth <= 768) {
    sidebar.toggleMobile();
    return;
  }
  sidebar.handleCollapse();
};
</script>

<style scoped>
.header-shell {
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-sizing: border-box;
  width: 100%;
  height: 70px;
  padding: 0 24px;
  color: #1e293b;
  background: rgba(255, 255, 255, 0.88);
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.1), 0 12px 28px rgba(15, 23, 42, 0.04);
}

.header-left,
.header-right,
.user-trigger,
.locale-switch {
  display: flex;
  align-items: center;
}

.header-left,
.header-right {
  gap: 14px;
}

.collapse-btn,
.user-trigger,
.locale-switch button {
  border: 0;
  cursor: pointer;
  touch-action: manipulation;
  transition-property: transform, background-color, color, box-shadow;
  transition-duration: 160ms;
  transition-timing-function: cubic-bezier(0.16, 1, 0.3, 1);
}

.collapse-btn:active,
.user-trigger:active,
.locale-switch button:active {
  transform: scale(0.96);
}

.collapse-btn {
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: 12px;
  background: #f7f5ef;
  color: #1f2933;
  font-size: 18px;
}

.route-title span {
  display: block;
  color: #64748b;
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.route-title strong {
  display: block;
  margin-top: 3px;
  font-size: 20px;
  letter-spacing: -0.012em;
}

.locale-switch {
  gap: 3px;
  padding: 4px;
  border-radius: 999px;
  background: #f7f5ef;
}

.locale-switch button {
  min-width: 44px;
  height: 30px;
  border-radius: 999px;
  background: transparent;
  color: #64748b;
  font-weight: 800;
}

.locale-switch button.active {
  background: #1f2933;
  color: rgba(255, 255, 255, 0.9);
  box-shadow: 0 8px 16px rgba(31, 41, 51, 0.14);
}

.user-trigger {
  gap: 10px;
  min-height: 42px;
  padding: 0 10px 0 4px;
  border-radius: 999px;
  background: #ffffff;
  color: #1e293b;
  box-shadow: inset 0 0 0 1px rgba(100, 116, 139, 0.14), 0 6px 16px rgba(15, 23, 42, 0.04);
  font-weight: 700;
}

.user-avatar {
  background: #dce9df;
}

@media (max-width: 640px) {
  .header-shell {
    padding: 0 14px;
  }

  .route-title span {
    display: none;
  }

  .route-title strong {
    font-size: 16px;
  }

  .user-trigger span {
    display: none;
  }
}

@media (prefers-reduced-motion: reduce) {
  .collapse-btn,
  .user-trigger,
  .locale-switch button {
    transition: none;
  }
}
</style>
