<template>
  <main class="settings-page">
    <section class="profile-card reveal-block">
      <img :src="marmotTwo" :alt="t('dashboard.images.brandAlt')" width="118" height="118" />
      <p class="eyebrow">{{ t('user.title') }}</p>
      <h1>{{ form.name || userStore.username || 'Marmot' }}</h1>
      <p>{{ form.desc || t('user.subtitle') }}</p>
      <div class="account-chip">{{ t('user.accountId') }} · {{ form.account || '-' }}</div>
    </section>

    <section class="settings-stack">
      <div class="settings-card reveal-block delay-1">
        <div class="section-head">
          <div>
            <p class="eyebrow">{{ t('user.profile') }}</p>
            <h2>{{ t('user.saveProfile') }}</h2>
          </div>
        </div>
        <el-form :model="form" :rules="rules" ref="ruleRef" label-position="top">
          <el-form-item :label="t('user.accountId')">
            <el-input v-model="form.account" disabled />
          </el-form-item>
          <el-form-item :label="t('user.nickname')" prop="name">
            <el-input v-model="form.name" />
          </el-form-item>
          <el-form-item :label="t('user.description')" prop="desc">
            <el-input v-model="form.desc" type="textarea" :autosize="{ minRows: 4, maxRows: 5 }" />
          </el-form-item>
          <button class="primary-action" :disabled="savingProfile" type="button" @click="submitForm(ruleRef)">
            {{ savingProfile ? '...' : t('user.saveProfile') }}
          </button>
        </el-form>
      </div>

      <div class="settings-card reveal-block delay-2">
        <div class="section-head">
          <div>
            <p class="eyebrow">{{ t('user.security') }}</p>
            <h2>{{ t('user.updatePassword') }}</h2>
          </div>
        </div>
        <el-form :model="passwordForm" :rules="passwordRules" ref="passwordRuleRef" label-position="top">
          <el-form-item :label="t('user.newPassword')" prop="first">
            <el-input v-model="passwordForm.first" type="password" show-password />
          </el-form-item>
          <el-form-item :label="t('user.confirmPassword')" prop="second">
            <el-input v-model="passwordForm.second" type="password" show-password />
          </el-form-item>
          <button class="primary-action" :disabled="savingPassword" type="button" @click="submitPasswordForm(passwordRuleRef)">
            {{ savingPassword ? '...' : t('user.updatePassword') }}
          </button>
        </el-form>
      </div>

      <div class="settings-card preferences-card reveal-block delay-3">
        <div>
          <p class="eyebrow">{{ t('user.preferences') }}</p>
          <h2>{{ t('user.language') }}</h2>
        </div>
        <div class="locale-switch" role="group" aria-label="language switch">
          <button :class="{ active: config.locale === 'zh-CN' }" @click="config.setLocale('zh-CN')">中文</button>
          <button :class="{ active: config.locale === 'en-US' }" @click="config.setLocale('en-US')">English</button>
        </div>
      </div>
    </section>
  </main>
</template>
<script setup>

import { computed, reactive, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { useUserStore } from '@/stores/user';
import { useConfigStore } from '@/stores/config';
import { putUser } from '@/api/user/user';
import { ElMessage } from 'element-plus';
import marmotTwo from '../../../img/marmot-ledger-2.png';

const { t } = useI18n();
const userStore = useUserStore();
const config = useConfigStore();
const savingProfile = ref(false);
const savingPassword = ref(false);

const form = reactive({
  account: userStore.account,
  name: userStore.username,
  desc: userStore.desc
})

const rules = computed(() => ({
  name: [
    { required: true, message: t('user.nameRequired'), trigger: 'blur' },
    { min: 0, max: 20, message: t('user.nameLength'), trigger: 'blur' },
  ],
  desc: [
    { min: 0, max: 300, message: t('user.descLength'), trigger: 'blur' },
  ],
}))

const ruleRef = ref()

const passwordForm = reactive({
  first: '',
  second: ''
})

const checkSamePassword = (rule, value, callback) => {
  if (passwordForm.first !== value) {
    return callback(new Error(t('user.passwordMismatch')))
  }
  callback()
}

const passwordRules = computed(() => ({
  first: [
    { required: true, message: t('user.passwordRequired'), trigger: 'blur' },
  ],
  second: [
    { required: true, message: t('user.passwordRequired'), trigger: 'blur' },
    { validator: checkSamePassword, trigger: 'blur' },
  ],
}))

const passwordRuleRef = ref()

watch(userStore, () => {
  form.account = userStore.account;
  form.name = userStore.username;
  form.desc = userStore.desc;
})

const submitForm = async (formEl) => {
  if (!formEl) return
  await formEl.validate(async (valid) => {
    if (!valid) return
    savingProfile.value = true
    try {
      const res = await putUser({ name: form.name, desc: form.desc })
      if (res.success) {
        ElMessage.success(t('user.profileSaved'))
        userStore.username = res.data.name;
        userStore.desc = res.data.desc;
      } else {
        ElMessage.error(res.error || t('messages.genericFailed'))
      }
    } finally {
      savingProfile.value = false
    }
  })
}

const submitPasswordForm = async (formEl) => {
  if (!formEl) return
  await formEl.validate(async (valid) => {
    if (!valid) return
    savingPassword.value = true
    try {
      const res = await putUser({ password: passwordForm.first })
      if (res.success) {
        ElMessage.success(t('user.passwordUpdated'))
        passwordForm.first = ''
        passwordForm.second = ''
      } else {
        ElMessage.error(res.error || t('messages.genericFailed'))
      }
    } finally {
      savingPassword.value = false
    }
  })
}
</script>

<style scoped>
.settings-page {
  max-width: 1200px;
  margin: 0 auto;
  display: grid;
  grid-template-columns: 340px minmax(0, 1fr);
  gap: 24px;
  color: #1e293b;
}

.reveal-block {
  animation: revealUp 500ms cubic-bezier(0.16, 1, 0.3, 1) both;
}

.delay-1 { animation-delay: 80ms; }
.delay-2 { animation-delay: 150ms; }
.delay-3 { animation-delay: 220ms; }

.profile-card,
.settings-card {
  background: #ffffff;
  border-radius: 22px;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.1), 0 16px 36px rgba(15, 23, 42, 0.05);
}

.profile-card {
  align-self: start;
  position: sticky;
  top: 24px;
  padding: 28px;
  background: linear-gradient(135deg, #fffaf0 0%, #ffffff 68%, #eef7f0 100%);
}

.profile-card img {
  border-radius: 28px;
  box-shadow: 0 16px 34px rgba(47, 125, 92, 0.18);
}

.eyebrow {
  margin: 22px 0 8px;
  color: #2f7d5c;
  font-size: 12px;
  font-weight: 800;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.profile-card h1 {
  margin: 0;
  font-size: 30px;
  line-height: 1.12;
  letter-spacing: -0.022em;
}

.profile-card p:not(.eyebrow) {
  margin: 12px 0 0;
  color: #64748b;
  line-height: 1.7;
}

.account-chip {
  display: inline-flex;
  margin-top: 20px;
  min-height: 34px;
  align-items: center;
  border-radius: 999px;
  padding: 0 14px;
  background: rgba(31, 41, 51, 0.08);
  color: #1f2933;
  font-family: 'SF Mono', 'Fira Code', monospace;
  font-size: 12px;
  font-weight: 800;
}

.settings-stack {
  display: grid;
  gap: 18px;
}

.settings-card {
  padding: 24px;
}

.section-head {
  display: flex;
  justify-content: space-between;
  margin-bottom: 18px;
}

.section-head .eyebrow,
.preferences-card .eyebrow {
  margin-top: 0;
}

.section-head h2,
.preferences-card h2 {
  margin: 0;
  font-size: 20px;
  letter-spacing: -0.012em;
}

.primary-action,
.locale-switch button {
  min-height: 40px;
  border: 0;
  border-radius: 12px;
  padding: 0 18px;
  font-weight: 800;
  cursor: pointer;
  transition-property: transform, box-shadow, background-color, color;
  transition-duration: 160ms;
  transition-timing-function: cubic-bezier(0.16, 1, 0.3, 1);
  touch-action: manipulation;
}

.primary-action:active,
.locale-switch button:active {
  transform: scale(0.96);
}

.primary-action {
  background: #2f7d5c;
  color: #ffffff;
  box-shadow: 0 10px 24px rgba(47, 125, 92, 0.20);
}

.primary-action:disabled {
  cursor: not-allowed;
  opacity: 0.62;
}

.preferences-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 18px;
}

.locale-switch {
  display: flex;
  gap: 4px;
  padding: 4px;
  border-radius: 999px;
  background: #f7f5ef;
}

.locale-switch button {
  border-radius: 999px;
  background: transparent;
  color: #64748b;
}

.locale-switch button.active {
  background: #1f2933;
  color: rgba(255, 255, 255, 0.9);
}

@media (max-width: 900px) {
  .settings-page {
    grid-template-columns: 1fr;
  }

  .profile-card {
    position: static;
  }
}

@media (max-width: 540px) {
  .profile-card,
  .settings-card {
    padding: 20px;
  }

  .preferences-card {
    align-items: flex-start;
    flex-direction: column;
  }
}

@media (prefers-reduced-motion: reduce) {
  .reveal-block,
  .primary-action,
  .locale-switch button {
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