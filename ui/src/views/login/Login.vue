<template>
  <main class="login-page">
    <button class="locale-switch" @click="config.toggleLocale">
      {{ config.locale === 'zh-CN' ? 'EN' : '中文' }}
    </button>

    <section class="brand-panel reveal-block">
      <div class="brand-images">
        <img :src="marmotOne" :alt="t('dashboard.images.brandAlt')" width="156" height="156" />
        <img :src="marmotTwo" :alt="t('dashboard.images.brandAlt')" width="104" height="104" />
      </div>
      <p class="eyebrow">{{ t('auth.cockpit') }}</p>
      <h1>{{ t('auth.title') }}</h1>
      <p>{{ t('auth.subtitle') }}</p>
      <div class="promise-grid">
        <span>Account</span>
        <span>Bucket</span>
        <span>Ledger Entry</span>
      </div>
    </section>

    <section class="login-card reveal-block delay-1">
      <p class="eyebrow">{{ t('auth.welcome') }}</p>
      <h2>{{ t('app.name') }}</h2>
      <div class="form-stack" @keyup.enter="handleLogin">
        <label>
          <span>{{ t('auth.username') }}</span>
          <el-input v-model="username" :placeholder="t('auth.usernamePlaceholder')" size="large" />
        </label>
        <label>
          <span>{{ t('auth.password') }}</span>
          <el-input v-model="password" :placeholder="t('auth.passwordPlaceholder')" type="password" show-password size="large" />
        </label>
        <button class="primary-action" :disabled="loading" @click="handleLogin">
          {{ loading ? '...' : t('auth.login') }}
        </button>
      </div>
    </section>
  </main>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { login } from '@/api/auth/auth'
import { useConfigStore } from '@/stores/config'
import marmotOne from '../../../../img/marmot-ledger-1.png'
import marmotTwo from '../../../../img/marmot-ledger-2.png'

const { t } = useI18n()
const config = useConfigStore()
const router = useRouter()
const username = ref('')
const password = ref('')
const loading = ref(false)

const handleLogin = async () => {
  if (!username.value || !password.value) {
    ElMessage.warning(t('auth.emptyError'))
    return
  }
  loading.value = true
  try {
    const res = await login(username.value, password.value)
    if (res.success) {
      localStorage.setItem('token', res.data.token)
      localStorage.setItem('validTime', String(new Date().getTime() + parseInt(res.data.expiredTime) * 1000))
      ElMessage.success(t('auth.success'))
      router.push('/')
    } else {
      ElMessage.error(res.error || t('messages.genericFailed'))
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: grid;
  grid-template-columns: minmax(0, 1.1fr) minmax(360px, 0.72fr);
  gap: 28px;
  align-items: center;
  padding: 48px;
  background:
    radial-gradient(circle at 14% 16%, rgba(220, 233, 223, 0.95), transparent 32%),
    linear-gradient(135deg, #f7f5ef 0%, #ffffff 52%, #eef7f0 100%);
  color: #1e293b;
}

.locale-switch,
.primary-action {
  border: 0;
  cursor: pointer;
  touch-action: manipulation;
  transition-property: transform, box-shadow, background-color, color;
  transition-duration: 160ms;
  transition-timing-function: cubic-bezier(0.16, 1, 0.3, 1);
}

.locale-switch:active,
.primary-action:active {
  transform: scale(0.96);
}

.locale-switch {
  position: fixed;
  top: 24px;
  right: 24px;
  z-index: 2;
  min-width: 62px;
  height: 36px;
  border-radius: 999px;
  background: #1f2933;
  color: rgba(255, 255, 255, 0.9);
  font-weight: 800;
  box-shadow: 0 10px 22px rgba(31, 41, 51, 0.16);
}

.brand-panel,
.login-card {
  animation: revealUp 560ms cubic-bezier(0.16, 1, 0.3, 1) both;
}

.delay-1 {
  animation-delay: 120ms;
}

.brand-panel {
  max-width: 760px;
}

.brand-images {
  position: relative;
  height: 210px;
  margin-bottom: 28px;
}

.brand-images img:first-child {
  position: absolute;
  left: 0;
  top: 18px;
  border-radius: 34px;
  box-shadow: 0 18px 42px rgba(47, 125, 92, 0.18);
}

.brand-images img:last-child {
  position: absolute;
  left: 138px;
  top: 92px;
  border-radius: 26px;
  box-shadow: 0 14px 32px rgba(31, 41, 51, 0.16);
}

.eyebrow {
  margin: 0 0 10px;
  color: #2f7d5c;
  font-size: 12px;
  font-weight: 800;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.brand-panel h1 {
  max-width: 720px;
  margin: 0;
  color: #1f2933;
  font-size: 46px;
  line-height: 1.05;
  letter-spacing: -0.022em;
  text-wrap: balance;
}

.brand-panel p:not(.eyebrow) {
  max-width: 620px;
  margin: 18px 0 0;
  color: #64748b;
  font-size: 16px;
  line-height: 1.75;
  text-wrap: pretty;
}

.promise-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 28px;
}

.promise-grid span {
  min-height: 34px;
  display: inline-flex;
  align-items: center;
  border-radius: 999px;
  padding: 0 14px;
  background: rgba(255, 255, 255, 0.78);
  color: #1e293b;
  font-family: 'SF Mono', 'Fira Code', monospace;
  font-size: 12px;
  font-weight: 800;
  box-shadow: inset 0 0 0 1px rgba(100, 116, 139, 0.12);
}

.login-card {
  padding: 32px;
  border-radius: 22px;
  background: rgba(255, 255, 255, 0.94);
  box-shadow: 0 2px 8px rgba(15, 23, 42, 0.08), 0 22px 60px rgba(31, 41, 51, 0.12);
}

.login-card h2 {
  margin: 0 0 26px;
  font-size: 30px;
  letter-spacing: -0.022em;
}

.form-stack {
  display: grid;
  gap: 18px;
}

.form-stack label span {
  display: block;
  margin-bottom: 8px;
  color: #1e293b;
  font-weight: 800;
}

.primary-action {
  height: 44px;
  border-radius: 12px;
  background: #3b82f6;
  color: #ffffff;
  font-weight: 800;
  box-shadow: 0 10px 24px rgba(59, 130, 246, 0.22);
}

.primary-action:disabled {
  cursor: not-allowed;
  opacity: 0.62;
}

@media (max-width: 900px) {
  .login-page {
    grid-template-columns: 1fr;
    padding: 32px 18px;
  }

  .brand-panel h1 {
    font-size: 34px;
  }

  .brand-images {
    height: 150px;
  }

  .brand-images img:first-child {
    width: 108px;
    height: 108px;
  }

  .brand-images img:last-child {
    left: 94px;
    top: 72px;
    width: 78px;
    height: 78px;
  }
}

@media (max-width: 520px) {
  .login-card {
    padding: 22px;
  }

  .brand-panel h1 {
    font-size: 28px;
  }
}

@media (prefers-reduced-motion: reduce) {
  .brand-panel,
  .login-card,
  .locale-switch,
  .primary-action {
    animation: none;
    transition: none;
  }
}

@keyframes revealUp {
  from {
    opacity: 0;
    transform: translateY(12px);
    filter: blur(4px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
    filter: blur(0);
  }
}
</style>
