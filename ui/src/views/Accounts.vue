<template>
  <main class="ledger-page">
    <section class="page-hero reveal-block">
      <div>
        <p class="eyebrow">{{ t('accounts.hero.eyebrow') }}</p>
        <h1>{{ t('accounts.hero.title') }}</h1>
        <p>{{ t('accounts.hero.subtitle') }}</p>
      </div>
      <button class="primary-action" @click="openCreate">{{ t('accounts.actions.new') }}</button>
    </section>

    <section class="toolbar reveal-block delay-1">
      <el-select v-model="filters.type" clearable :placeholder="t('accounts.filters.typePlaceholder')" class="filter-control" @change="loadAccounts">
        <el-option v-for="item in accountTypes" :key="item.value" :label="item.label.value || item.label" :value="item.value" />
      </el-select>
      <el-select v-model="filters.isActive" clearable :placeholder="t('common.filters.enabledStatus')" class="filter-control" @change="loadAccounts">
        <el-option :label="t('common.status.enabled')" :value="true" />
        <el-option :label="t('common.status.disabled')" :value="false" />
      </el-select>
      <button class="ghost-action" @click="loadAccounts">{{ t('common.actions.refresh') }}</button>
    </section>

    <section v-loading="loading" class="account-grid reveal-block delay-2">
      <article v-for="(item, index) in accounts" :key="item.id" class="account-card" :style="{ animationDelay: `${index * 55}ms` }">
        <div class="account-marker" :style="{ background: item.color || '#3b82f6' }"></div>
        <div class="account-main">
          <div class="account-head">
            <div>
              <h2>{{ item.displayName || item.name }}</h2>
              <p>{{ item.name }}</p>
            </div>
            <span class="status-pill" :class="{ inactive: !item.isActive }">{{ item.isActive ? t('common.status.enabled') : t('common.status.disabled') }}</span>
          </div>
          <div class="account-meta">
            <span>{{ getTypeLabel(item.type) }}</span>
          </div>
          <div class="card-actions">
            <button class="text-action" @click="openEdit(item)">{{ t('common.actions.edit') }}</button>
            <button class="danger-action" @click="handleDelete(item)">{{ t('common.actions.delete') }}</button>
          </div>
        </div>
      </article>

      <div v-if="!loading && !accounts.length" class="empty-state">
        <img :src="marmotOne" :alt="t('accounts.empty.alt')" width="112" height="112" />
        <h2>{{ t('accounts.empty.title') }}</h2>
        <p>{{ t('accounts.empty.text') }}</p>
        <button class="primary-action" @click="openCreate">{{ t('accounts.actions.new') }}</button>
      </div>
    </section>

    <el-dialog v-model="dialogVisible" :title="editingId ? t('accounts.dialog.editTitle') : t('accounts.dialog.createTitle')" width="520px" class="marmot-dialog">
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <el-form-item :label="t('accounts.fields.accountName')" prop="name">
          <el-input v-model="form.name" :placeholder="t('accounts.placeholders.name')" />
        </el-form-item>
        <el-form-item :label="t('accounts.fields.displayName')">
          <el-input v-model="form.displayName" :placeholder="t('accounts.placeholders.displayName')" />
        </el-form-item>
        <el-form-item :label="t('accounts.fields.accountType')" prop="type">
          <el-select v-model="form.type" :placeholder="t('accounts.placeholders.selectType')" class="full-width">
            <el-option v-for="item in accountTypes" :key="item.value" :label="item.label.value || item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('common.fields.color')">
          <div class="color-palette">
            <button v-for="color in colors" :key="color" type="button" :class="['color-dot', { active: form.color === color }]" :style="{ background: color }" @click="form.color = color" :aria-label="`${t('common.fields.color')} ${color}`"></button>
          </div>
        </el-form-item>
        <el-form-item v-if="editingId" :label="t('common.status.status')">
          <el-switch v-model="form.isActive" :active-text="t('common.status.enabled')" :inactive-text="t('common.status.disabled')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <button class="ghost-action" @click="dialogVisible = false">{{ t('common.actions.cancel') }}</button>
        <button class="primary-action" @click="submitForm">{{ t('common.actions.save') }}</button>
      </template>
    </el-dialog>
  </main>
</template>

<script setup>
import { computed, onActivated, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { createAccount, deleteAccount, listAccounts, updateAccount } from '@/api/account/account'
import marmotOne from '../../../img/marmot-ledger-1.png'

const { t } = useI18n()
const accounts = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const editingId = ref(0)
const formRef = ref()
const filters = reactive({
  type: '',
  isActive: '',
})
const form = reactive(createEmptyForm())

const accountTypes = [
  { label: computed(() => t('accounts.types.cash')), value: 'cash' },
  { label: computed(() => t('accounts.types.wallet')), value: 'wallet' },
  { label: computed(() => t('accounts.types.bank')), value: 'bank' },
  { label: computed(() => t('accounts.types.credit')), value: 'credit' },
  { label: computed(() => t('accounts.types.investment')), value: 'investment' },
  { label: computed(() => t('accounts.types.liability')), value: 'liability' },
  { label: computed(() => t('accounts.types.other')), value: 'other' },
]
const colors = ['#3b82f6', '#2f7d5c', '#1f2933', '#f59e0b', '#ef4444', '#10b981']
const rules = {
  name: [{ required: true, message: t('accounts.validation.nameRequired'), trigger: 'blur' }],
  type: [{ required: true, message: t('accounts.validation.typeRequired'), trigger: 'change' }],
}

function createEmptyForm() {
  return {
    name: '',
    type: '',
    providerCode: '',
    accountGroupKey: '',
    standardAccountId: 0,
    bankCode: '',
    displayName: '',
    icon: '',
    color: '#3b82f6',
    isActive: true,
  }
}

function resetForm(data = createEmptyForm()) {
  Object.assign(form, data)
}

function getTypeLabel(type) {
  return accountTypes.find((item) => item.value === type)?.label.value || type || t('common.misc.uncategorized')
}

async function loadAccounts() {
  loading.value = true
  try {
    const params = {}
    if (filters.type) params.type = filters.type
    if (filters.isActive !== '') params.isActive = filters.isActive
    const res = await listAccounts(params)
    if (res.success) {
      accounts.value = res.data || []
    } else {
      ElMessage.error(res.error || t('accounts.messages.loadFailed'))
    }
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editingId.value = 0
  resetForm()
  dialogVisible.value = true
}

function openEdit(item) {
  editingId.value = item.id
  resetForm({ ...createEmptyForm(), ...item, isActive: item.isActive !== false })
  dialogVisible.value = true
}

async function submitForm() {
  await formRef.value?.validate()
  const payload = { ...form, isActive: form.isActive !== false }
  const res = editingId.value ? await updateAccount(editingId.value, payload) : await createAccount(payload)
  if (res.success) {
    ElMessage.success(editingId.value ? t('accounts.messages.updated') : t('accounts.messages.created'))
    dialogVisible.value = false
    await loadAccounts()
  } else {
    ElMessage.error(res.error || t('accounts.messages.saveFailed'))
  }
}

async function handleDelete(item) {
  try {
    await ElMessageBox.confirm(t('accounts.delete.confirm', { name: item.displayName || item.name }), t('accounts.delete.title'), {
      confirmButtonText: t('common.actions.delete'),
      cancelButtonText: t('common.actions.cancel'),
      type: 'warning',
      customClass: 'calm-marmot-message-box calm-marmot-delete-box',
      confirmButtonClass: 'calm-marmot-danger-confirm',
      cancelButtonClass: 'calm-marmot-soft-cancel',
    })
    const res = await deleteAccount(item.id)
    if (res.success) {
      ElMessage.success(t('accounts.messages.deleted'))
      await loadAccounts()
    } else {
      ElMessage.error(res.error || t('accounts.messages.deleteFailed'))
    }
  } catch (err) {
    if (err !== 'cancel') console.warn(err)
  }
}

onMounted(loadAccounts)
onActivated(loadAccounts)
</script>

<style scoped>
.ledger-page {
  max-width: 1200px;
  margin: 0 auto;
  color: #1e293b;
}

.reveal-block,
.account-card {
  animation: revealUp 480ms cubic-bezier(0.16, 1, 0.3, 1) both;
}

.delay-1 {
  animation-delay: 90ms;
}

.delay-2 {
  animation-delay: 160ms;
}

.page-hero,
.toolbar,
.account-card,
.empty-state {
  background: #ffffff;
  border-radius: 16px;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.1), 0 12px 30px rgba(15, 23, 42, 0.04);
}

.page-hero {
  display: flex;
  justify-content: space-between;
  gap: 24px;
  align-items: flex-start;
  margin-bottom: 18px;
  padding: 26px;
  background: linear-gradient(135deg, #fffaf0 0%, #ffffff 70%);
}

.eyebrow {
  margin: 0 0 8px;
  color: #2f7d5c;
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.page-hero h1 {
  max-width: 680px;
  margin: 0;
  font-size: 30px;
  line-height: 1.16;
  letter-spacing: -0.022em;
  text-wrap: balance;
}

.page-hero p:last-child {
  max-width: 620px;
  margin: 12px 0 0;
  color: #64748b;
  line-height: 1.7;
}

.toolbar {
  display: flex;
  gap: 12px;
  align-items: center;
  margin-bottom: 18px;
  padding: 14px;
}

.filter-control {
  width: 180px;
}

.account-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.account-card {
  position: relative;
  display: grid;
  grid-template-columns: 10px minmax(0, 1fr);
  overflow: hidden;
  min-height: 176px;
  transition-property: transform, box-shadow;
  transition-duration: 180ms;
}

.account-marker {
  min-height: 100%;
}

.account-main {
  padding: 22px;
}

.account-head {
  display: flex;
  justify-content: space-between;
  gap: 12px;
}

.account-head h2 {
  margin: 0;
  font-size: 20px;
  letter-spacing: -0.012em;
}

.account-head p {
  margin: 6px 0 0;
  color: #64748b;
}

.status-pill,
.account-meta span {
  display: inline-flex;
  align-items: center;
  min-height: 28px;
  border-radius: 999px;
  padding: 0 10px;
  background: rgba(47, 125, 92, 0.1);
  color: #2f7d5c;
  font-size: 12px;
  font-weight: 700;
}

.status-pill.inactive {
  background: rgba(100, 116, 139, 0.12);
  color: #64748b;
}

.account-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 18px;
}

.card-actions {
  display: flex;
  gap: 8px;
  margin-top: 22px;
}

.primary-action,
.ghost-action,
.text-action,
.danger-action {
  min-height: 40px;
  border: 0;
  border-radius: 12px;
  padding: 0 16px;
  font-weight: 700;
  cursor: pointer;
  transition-property: transform, box-shadow, background-color, color;
  transition-duration: 160ms;
  transition-timing-function: cubic-bezier(0.16, 1, 0.3, 1);
  touch-action: manipulation;
}

.primary-action:active,
.ghost-action:active,
.text-action:active,
.danger-action:active,
.color-dot:active {
  transform: scale(0.96);
}

.primary-action {
  background: #3b82f6;
  color: #ffffff;
  box-shadow: 0 10px 24px rgba(59, 130, 246, 0.22);
}

.ghost-action,
.text-action {
  background: #f8faf7;
  color: #1e293b;
}

.danger-action {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
}

.empty-state {
  grid-column: 1 / -1;
  display: grid;
  place-items: center;
  gap: 12px;
  padding: 40px 24px;
  text-align: center;
  color: #64748b;
}

.empty-state img {
  border-radius: 22px;
}

.empty-state h2 {
  margin: 0;
  color: #1e293b;
}

.empty-state p {
  max-width: 420px;
  margin: 0;
  line-height: 1.7;
}

.full-width {
  width: 100%;
}

.color-palette {
  display: flex;
  gap: 10px;
}

.color-dot {
  width: 30px;
  height: 30px;
  border: 0;
  border-radius: 999px;
  cursor: pointer;
  box-shadow: inset 0 0 0 2px rgba(255, 255, 255, 0.9), 0 0 0 1px rgba(15, 23, 42, 0.12);
  transition-property: transform, box-shadow;
  transition-duration: 160ms;
}

.color-dot.active {
  box-shadow: inset 0 0 0 2px rgba(255, 255, 255, 0.9), 0 0 0 3px rgba(59, 130, 246, 0.26);
}

@media (hover: hover) {
  .account-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 3px 8px rgba(15, 23, 42, 0.12), 0 16px 34px rgba(15, 23, 42, 0.06);
  }
}

@media (max-width: 820px) {
  .page-hero,
  .toolbar {
    flex-direction: column;
    align-items: stretch;
  }

  .account-grid {
    grid-template-columns: 1fr;
  }

  .filter-control {
    width: 100%;
  }
}

@media (max-width: 520px) {
  .page-hero {
    padding: 20px;
  }

  .page-hero h1 {
    font-size: 24px;
  }

  .account-head {
    flex-direction: column;
  }
}

@media (prefers-reduced-motion: reduce) {
  .reveal-block,
  .account-card,
  .primary-action,
  .ghost-action,
  .text-action,
  .danger-action,
  .color-dot {
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
