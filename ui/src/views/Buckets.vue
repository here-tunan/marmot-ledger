<template>
  <main class="ledger-page">
    <section class="page-hero reveal-block">
      <div>
        <p class="eyebrow">{{ t('buckets.hero.eyebrow') }}</p>
        <h1>{{ t('buckets.hero.title') }}</h1>
        <p>{{ t('buckets.hero.subtitle') }}</p>
      </div>
      <button class="primary-action" :disabled="!accounts.length" @click="openCreate">{{ t('buckets.actions.new') }}</button>
    </section>

    <section class="toolbar reveal-block delay-1">
      <el-select v-model="filters.accountId" clearable :placeholder="t('buckets.filters.accountPlaceholder')" class="filter-control" @change="loadBuckets">
        <el-option v-for="item in accounts" :key="item.id" :label="item.name" :value="item.id" />
      </el-select>
      <el-select v-model="filters.currency" clearable :placeholder="t('buckets.filters.currencyPlaceholder')" class="filter-control" @change="loadBuckets">
        <el-option v-for="item in currencyOptions" :key="item.code" :label="getCurrencyLabel(item.code, config.locale)" :value="item.code" />
      </el-select>
      <el-select v-model="filters.bucketNature" clearable :placeholder="t('buckets.filters.naturePlaceholder')" class="filter-control" @change="loadBuckets">
        <el-option :label="t('domain.asset')" value="asset" />
        <el-option :label="t('domain.liability')" value="liability" />
      </el-select>
      <button class="ghost-action" @click="refreshAll">{{ t('common.actions.refresh') }}</button>
      <button class="export-action" type="button" @click="handleExportBuckets">{{ t('records.actions.exportBuckets') }}</button>
    </section>

    <section class="bucket-layout reveal-block delay-2">
      <div v-loading="loading" class="bucket-list">
        <article v-for="(item, index) in buckets" :key="item.id" class="bucket-card" :class="{ active: selectedBucket?.id === item.id }" :style="{ animationDelay: `${index * 55}ms` }" @click="selectBucket(item)">
          <div class="bucket-topline">
            <div>
              <h2>{{ getBucketEmoji(item.bucketType) }} {{ item.name }}</h2>
              <p>{{ getAccountName(item.accountId) }}</p>
            </div>
            <span class="nature-pill" :class="item.bucketNature">{{ item.bucketNature === 'liability' ? t('domain.liability') : t('domain.asset') }}</span>
          </div>
          <div class="bucket-balance">
            <span>{{ getCurrencyDisplay(item.currency).icon }} {{ item.currency }}</span>
            <strong>{{ formatAmount(item.balance) }}</strong>
          </div>
          <div class="bucket-meta">
            <span>{{ getBucketTypeLabel(item.bucketType) }}</span>
            <span>{{ t('buckets.card.initial', { amount: formatAmount(item.initialBalance) }) }}</span>
          </div>
        </article>

        <div v-if="!loading && !buckets.length" class="empty-state">
          <img :src="marmotTwo" :alt="t('buckets.empty.alt')" width="112" height="112" />
          <h2>{{ t('buckets.empty.title') }}</h2>
          <p>{{ t('buckets.empty.text') }}</p>
          <button class="primary-action" :disabled="!accounts.length" @click="openCreate">{{ t('buckets.actions.new') }}</button>
        </div>
      </div>

      <aside class="ledger-panel">
        <div class="section-head">
          <div>
            <p class="eyebrow">{{ t('buckets.ledger.eyebrow') }}</p>
            <h2>{{ selectedBucket ? `${getBucketEmoji(selectedBucket.bucketType)} ${selectedBucket.name}` : t('buckets.ledger.selectBucketTitle') }}</h2>
          </div>
        </div>

        <div v-if="selectedBucket" class="selected-summary">
          <span>{{ getCurrencyDisplay(selectedBucket.currency).icon }} {{ selectedBucket.currency }}</span>
          <strong>{{ formatAmount(selectedBucket.balance) }}</strong>
          <p>{{ getAccountName(selectedBucket.accountId) }} · {{ getBucketTypeLabel(selectedBucket.bucketType) }}</p>
          <div class="selected-actions">
            <button class="edit-selected-action" type="button" @click="openEdit(selectedBucket)">{{ t('common.actions.edit') }}</button>
            <button v-if="isInvestmentBucket(selectedBucket)" class="adjust-selected-action" type="button" :disabled="selectedBucket.isActive === false" @click="openRevalue(selectedBucket)">{{ t('buckets.actions.revalue') }}</button>
            <button v-else class="adjust-selected-action" type="button" :disabled="selectedBucket.isActive === false" @click="openAdjust(selectedBucket)">{{ t('buckets.actions.adjustBalance') }}</button>
          </div>
        </div>

        <div v-if="ledgerEntries.length" class="entry-section">
          <div class="entry-list" :class="{ expanded: ledgerExpanded }">
            <div v-for="entry in visibleLedgerEntries" :key="entry.id" class="entry-row">
              <div>
                <strong>{{ entryRoleLabel(entry.entryRole) }}</strong>
                <span>{{ entry.createdAt }}</span>
              </div>
              <div class="entry-amount">{{ entry.currency }} {{ formatAmount(entry.amount) }}</div>
            </div>
          </div>
          <button v-if="ledgerEntries.length > ledgerPreviewCount" class="entry-toggle" type="button" @click="ledgerExpanded = !ledgerExpanded">
            {{ ledgerExpanded ? t('buckets.ledger.fold') : t('buckets.ledger.unfold', { count: ledgerEntries.length - ledgerPreviewCount }) }}
          </button>
        </div>
        <div v-else class="empty-state compact">
          <img :src="marmotOne" :alt="t('buckets.ledger.emptyAlt')" width="88" height="88" />
          <p>{{ selectedBucket ? t('buckets.ledger.emptyForSelected') : t('buckets.ledger.emptyNoSelection') }}</p>
        </div>
      </aside>
    </section>

    <el-dialog v-model="dialogVisible" :title="editingId ? t('buckets.dialog.editTitle') : t('buckets.dialog.createTitle')" width="560px" class="marmot-dialog">
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <el-form-item :label="t('buckets.fields.account')" prop="accountId">
          <el-select v-model="form.accountId" :placeholder="t('buckets.placeholders.selectAccount')" class="full-width" :disabled="Boolean(editingId)">
            <el-option v-for="item in accounts" :key="item.id" :label="item.name" :value="item.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('buckets.fields.bucketName')" prop="name">
          <el-input v-model="form.name" :placeholder="t('buckets.placeholders.name')" />
        </el-form-item>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item :label="t('common.fields.currency')" prop="currency">
              <el-select v-model="form.currency" class="full-width" :disabled="Boolean(editingId)">
                <el-option v-for="item in currencyOptions" :key="item.code" :label="getCurrencyLabel(item.code, config.locale)" :value="item.code" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="t('common.fields.initialBalance')" prop="initialBalance">
              <el-input v-model="form.initialBalance" :placeholder="t('buckets.placeholders.initialBalance')" :disabled="Boolean(editingId)" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item :label="t('buckets.fields.bucketType')" prop="bucketType">
              <el-select v-model="form.bucketType" class="full-width">
                <el-option v-for="item in bucketTypes" :key="item.value" :label="item.label.value || item.label" :value="item.value" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="t('buckets.fields.bucketNature')" prop="bucketNature">
              <el-select v-model="form.bucketNature" class="full-width">
                <el-option :label="t('domain.asset')" value="asset" />
                <el-option :label="t('domain.liability')" value="liability" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item v-if="editingId" :label="t('common.status.status')">
          <el-switch v-model="form.isActive" :active-text="t('common.status.enabled')" :inactive-text="t('common.status.disabled')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <button class="ghost-action" @click="dialogVisible = false">{{ t('common.actions.cancel') }}</button>
        <button class="primary-action" @click="submitForm">{{ editingId ? t('common.actions.save') : t('buckets.actions.create') }}</button>
      </template>
    </el-dialog>

    <el-dialog v-model="adjustDialogVisible" :title="t('buckets.adjustDialog.title')" width="520px" class="marmot-dialog">
      <el-form label-position="top">
        <div v-if="selectedBucket" class="adjust-summary">
          <span>{{ getBucketEmoji(selectedBucket.bucketType) }} {{ selectedBucket.name }}</span>
          <strong>{{ getCurrencyDisplay(selectedBucket.currency).icon }} {{ selectedBucket.currency }} {{ formatAmount(selectedBucket.balance) }}</strong>
        </div>
        <el-form-item :label="t('buckets.adjustDialog.targetBalance')">
          <el-input v-model="adjustForm.targetBalance" :placeholder="t('buckets.placeholders.targetBalance')" />
        </el-form-item>
        <div class="adjust-diff" :class="adjustmentDeltaClass">
          <span>{{ t('buckets.adjustDialog.difference') }}</span>
          <strong>{{ signedAdjustmentText }}</strong>
          <small>{{ adjustmentDeltaHint }}</small>
        </div>
        <el-form-item :label="t('buckets.adjustDialog.remark')">
          <el-input v-model="adjustForm.remark" type="textarea" :placeholder="t('buckets.adjustDialog.remarkPlaceholder')" :autosize="{ minRows: 2, maxRows: 4 }" />
        </el-form-item>
      </el-form>
      <template #footer>
        <button class="ghost-action" @click="adjustDialogVisible = false">{{ t('common.actions.cancel') }}</button>
        <button class="primary-action" @click="submitAdjustment">{{ t('common.actions.save') }}</button>
      </template>
    </el-dialog>

    <el-dialog v-model="revalueDialogVisible" :title="t('buckets.revalueDialog.title')" width="520px" class="marmot-dialog">
      <el-form label-position="top">
        <div v-if="selectedBucket" class="adjust-summary">
          <span>{{ getBucketEmoji(selectedBucket.bucketType) }} {{ selectedBucket.name }}</span>
          <strong>{{ getCurrencyDisplay(selectedBucket.currency).icon }} {{ selectedBucket.currency }} {{ formatAmount(selectedBucket.balance) }}</strong>
        </div>
        <el-form-item :label="t('buckets.revalueDialog.targetMarketValue')">
          <el-input v-model="revalueForm.targetMarketValue" :placeholder="t('buckets.placeholders.targetMarketValue')" />
        </el-form-item>
        <div class="adjust-diff" :class="revalueDeltaClass">
          <span>{{ t('buckets.revalueDialog.pnl') }}</span>
          <strong>{{ signedRevalueText }}</strong>
          <small>{{ revalueDeltaHint }}</small>
        </div>
        <el-form-item :label="t('buckets.revalueDialog.remark')">
          <el-input v-model="revalueForm.remark" type="textarea" :placeholder="t('buckets.revalueDialog.remarkPlaceholder')" :autosize="{ minRows: 2, maxRows: 4 }" />
        </el-form-item>
      </el-form>
      <template #footer>
        <button class="ghost-action" @click="revalueDialogVisible = false">{{ t('common.actions.cancel') }}</button>
        <button class="primary-action" @click="submitRevalue">{{ t('common.actions.save') }}</button>
      </template>
    </el-dialog>
  </main>
</template>

<script setup>
import { computed, onActivated, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useConfigStore } from '@/stores/config'
import { ElMessage } from 'element-plus'
import { listAccounts } from '@/api/account/account'
import { createBucket, listBucketLedgerEntries, listBuckets, updateBucket } from '@/api/bucket/bucket'
import { createRecord } from '@/api/record/record'
import { exportBuckets } from '@/api/export/export'
import marmotOne from '../../../img/marmot-ledger-1.png'
import marmotTwo from '../../../img/marmot-ledger-2.png'
import { currencyOptions, getCurrencyDisplay, getCurrencyLabel } from '@/utils/currency'
import { getBucketEmoji } from '@/utils/bucketEmoji'

const { t, te } = useI18n()
const config = useConfigStore()
const accounts = ref([])
const buckets = ref([])
const ledgerEntries = ref([])
const ledgerExpanded = ref(false)
const ledgerPreviewCount = 6
const selectedBucket = ref(null)
const loading = ref(false)
const dialogVisible = ref(false)
const adjustDialogVisible = ref(false)
const revalueDialogVisible = ref(false)
const editingId = ref(0)
const formRef = ref()
const filters = reactive({
  accountId: '',
  currency: '',
  bucketNature: '',
})
const form = reactive(createEmptyForm())
const adjustForm = reactive({ targetBalance: '', remark: '' })
const revalueForm = reactive({ targetMarketValue: '', remark: '' })

const bucketTypes = [
  { label: computed(() => t('buckets.types.cash')), value: 'cash' },
  { label: computed(() => t('buckets.types.wallet')), value: 'wallet' },
  { label: computed(() => t('buckets.types.bank')), value: 'bank' },
  { label: computed(() => t('buckets.types.credit')), value: 'credit' },
  { label: computed(() => t('buckets.types.investmentCash')), value: 'investment_cash' },
  { label: computed(() => t('buckets.types.investmentAsset')), value: 'investment_asset' },
  { label: computed(() => t('buckets.types.receivable')), value: 'receivable' },
  { label: computed(() => t('buckets.types.deposit')), value: 'deposit' },
  { label: computed(() => t('buckets.types.loanOut')), value: 'loan_out' },
  { label: computed(() => t('buckets.types.liability')), value: 'liability' },
  { label: computed(() => t('buckets.types.virtual')), value: 'virtual' },
]
const rules = {
  accountId: [{ required: true, message: t('buckets.validation.accountRequired'), trigger: 'change' }],
  name: [{ required: true, message: t('buckets.validation.nameRequired'), trigger: 'blur' }],
  currency: [{ required: true, message: t('buckets.validation.currencyRequired'), trigger: 'change' }],
  initialBalance: [{ required: true, message: t('buckets.validation.initialBalanceRequired'), trigger: 'blur' }],
  bucketType: [{ required: true, message: t('buckets.validation.typeRequired'), trigger: 'change' }],
  bucketNature: [{ required: true, message: t('buckets.validation.natureRequired'), trigger: 'change' }],
}

const adjustmentDelta = computed(() => {
  if (!selectedBucket.value || adjustForm.targetBalance === '') return 0
  const target = Number(adjustForm.targetBalance)
  if (!Number.isFinite(target)) return 0
  return roundMoney(target - Number(selectedBucket.value.balance || 0))
})

const adjustmentDeltaClass = computed(() => {
  if (!adjustmentDelta.value) return 'neutral'
  return adjustmentDelta.value > 0 ? 'positive' : 'negative'
})

const signedAdjustmentText = computed(() => {
  const delta = adjustmentDelta.value
  const sign = delta > 0 ? '+' : ''
  return `${sign}${formatAmount(delta)}`
})

const adjustmentDeltaHint = computed(() => {
  if (!adjustmentDelta.value) return t('buckets.adjustDialog.noChange')
  return adjustmentDelta.value > 0 ? t('buckets.adjustDialog.increase') : t('buckets.adjustDialog.decrease')
})

const revalueDelta = computed(() => {
  if (!selectedBucket.value || revalueForm.targetMarketValue === '') return 0
  const target = Number(revalueForm.targetMarketValue)
  if (!Number.isFinite(target)) return 0
  return target - Number(selectedBucket.value.balance)
})

const revalueDeltaClass = computed(() => {
  if (!revalueDelta.value) return 'neutral'
  return revalueDelta.value > 0 ? 'positive' : 'negative'
})

const signedRevalueText = computed(() => {
  const delta = revalueDelta.value
  if (!delta) return formatAmount(0)
  return (delta > 0 ? '+' : '−') + formatAmount(Math.abs(delta))
})

const revalueDeltaHint = computed(() => {
  if (!revalueDelta.value) return t('buckets.revalueDialog.pnlEqual')
  return revalueDelta.value > 0 ? t('buckets.revalueDialog.pnlGain') : t('buckets.revalueDialog.pnlLoss')
})

const visibleLedgerEntries = computed(() => {
  if (ledgerExpanded.value) return ledgerEntries.value
  return ledgerEntries.value.slice(0, ledgerPreviewCount)
})

function isInvestmentBucket(bucket) {
  if (!bucket) return false
  return bucket.bucketType === 'investment_asset' || bucket.bucketType === 'investment_cash'
}

function createEmptyForm() {
  return {
    accountId: '',
    name: '',
    currency: 'CNY',
    initialBalance: '0.0000',
    bucketType: 'cash',
    bucketNature: 'asset',
    isActive: true,
  }
}

function resetForm(data = createEmptyForm()) {
  Object.assign(form, createEmptyForm(), data)
}

function getAccountName(accountId) {
  const account = accounts.value.find((item) => item.id === accountId)
  return account?.name || t('buckets.accountFallback', { accountId })
}

function getBucketTypeLabel(type) {
  return bucketTypes.find((item) => item.value === type)?.label.value || type
}

function entryRoleLabel(role) {
  const key = `record.entryRoles.${role}`
  return role && te(key) ? t(key) : t('buckets.entryRoleFallback')
}

function formatAmount(value) {
  const number = Number(value || 0)
  return new Intl.NumberFormat(config.locale, {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(number)
}

function roundMoney(value) {
  return Math.round(Number(value || 0) * 10000) / 10000
}

function formatPlainAmount(value) {
  return roundMoney(value).toFixed(4)
}

async function loadAccounts() {
  const res = await listAccounts({ isActive: true })
  if (res.success) accounts.value = res.data || []
}

async function loadBuckets() {
  loading.value = true
  try {
    const params = {}
    if (filters.accountId) params.accountId = filters.accountId
    if (filters.currency) params.currency = filters.currency
    if (filters.bucketNature) params.bucketNature = filters.bucketNature
    const res = await listBuckets(params)
    if (res.success) {
      buckets.value = res.data || []
      if (selectedBucket.value && !buckets.value.some((item) => item.id === selectedBucket.value.id)) {
        selectedBucket.value = null
        ledgerEntries.value = []
        ledgerExpanded.value = false
      }
    } else {
      ElMessage.error(res.error || t('buckets.messages.loadFailed'))
    }
  } finally {
    loading.value = false
  }
}

async function refreshAll() {
  await loadAccounts()
  await loadBuckets()
}

async function handleExportBuckets() {
  try { await exportBuckets() }
  catch (err) { console.warn(err); ElMessage.error(t('records.messages.exportFailed')) }
}

function openCreate() {
  if (!accounts.value.length) {
    ElMessage.warning(t('buckets.messages.createAccountFirst'))
    return
  }
  editingId.value = 0
  resetForm()
  form.accountId = accounts.value[0]?.id || ''
  dialogVisible.value = true
}

function openEdit(item) {
  editingId.value = item.id
  resetForm({
    accountId: item.accountId,
    name: item.name,
    currency: item.currency,
    initialBalance: String(item.initialBalance || '0.0000'),
    bucketType: item.bucketType,
    bucketNature: item.bucketNature,
    isActive: item.isActive !== false,
  })
  dialogVisible.value = true
}

function openAdjust(item) {
  selectedBucket.value = item
  adjustForm.targetBalance = formatPlainAmount(item.balance)
  adjustForm.remark = ''
  adjustDialogVisible.value = true
}

function openRevalue(item) {
  selectedBucket.value = item
  revalueForm.targetMarketValue = formatPlainAmount(item.balance)
  revalueForm.remark = ''
  revalueDialogVisible.value = true
}

async function submitForm() {
  await formRef.value?.validate()
  if (editingId.value) {
    const payload = {
      name: form.name,
      bucketType: form.bucketType,
      bucketNature: form.bucketNature,
      isActive: form.isActive !== false,
    }
    const res = await updateBucket(editingId.value, payload)
    if (res.success) {
      ElMessage.success(t('buckets.messages.updated'))
      dialogVisible.value = false
      await loadBuckets()
      const updated = buckets.value.find((item) => item.id === editingId.value)
      if (updated) {
        selectedBucket.value = updated
      }
    } else {
      ElMessage.error(res.error || t('buckets.messages.updateFailed'))
    }
    return
  }

  const payload = {
    ...form,
    accountId: Number(form.accountId),
    initialBalance: String(form.initialBalance),
  }
  const res = await createBucket(payload)
  if (res.success) {
    ElMessage.success(t('buckets.messages.createdWithRefs', { eventId: res.data.initialFinancialEventId, entryId: res.data.initialLedgerEntryId }))
    dialogVisible.value = false
    await loadBuckets()
    const created = buckets.value.find((item) => item.id === res.data.id)
    if (created) await selectBucket(created)
  } else {
    ElMessage.error(res.error || t('buckets.messages.createFailed'))
  }
}

async function selectBucket(item) {
  selectedBucket.value = item
  ledgerExpanded.value = false
  const res = await listBucketLedgerEntries(item.id)
  if (res.success) {
    ledgerEntries.value = res.data || []
  } else {
    ledgerEntries.value = []
    ElMessage.error(res.error || t('buckets.messages.loadEntriesFailed'))
  }
}

async function submitAdjustment() {
  if (!selectedBucket.value) return
  if (adjustForm.targetBalance === '') {
    ElMessage.warning(t('buckets.validation.targetBalanceRequired'))
    return
  }
  const target = Number(adjustForm.targetBalance)
  if (!Number.isFinite(target) || target < 0) {
    ElMessage.warning(t('buckets.validation.targetBalanceInvalid'))
    return
  }
  const delta = adjustmentDelta.value
  if (!delta) {
    ElMessage.warning(t('buckets.validation.targetBalanceUnchanged'))
    return
  }

  const bucketId = selectedBucket.value.id
  const payload = {
    scenario: 'balance_adjustment',
    bucketId,
    amount: formatPlainAmount(delta),
    currency: selectedBucket.value.currency,
    description: t('buckets.adjustDialog.description', { from: formatAmount(selectedBucket.value.balance), to: formatAmount(target) }),
    remark: adjustForm.remark,
  }
  const res = await createRecord(payload)
  if (res.success) {
    ElMessage.success(t('buckets.messages.adjusted'))
    adjustDialogVisible.value = false
    await loadBuckets()
    const updated = buckets.value.find((item) => item.id === bucketId)
    if (updated) await selectBucket(updated)
  } else {
    ElMessage.error(res.error || t('buckets.messages.adjustFailed'))
  }
}

async function submitRevalue() {
  if (!selectedBucket.value) return
  if (revalueForm.targetMarketValue === '') {
    ElMessage.warning(t('buckets.validation.targetBalanceRequired'))
    return
  }
  const target = Number(revalueForm.targetMarketValue)
  if (!Number.isFinite(target) || target < 0) {
    ElMessage.warning(t('buckets.validation.targetBalanceInvalid'))
    return
  }
  const delta = revalueDelta.value
  if (!delta || Math.abs(delta) < 0.001) {
    ElMessage.warning(t('buckets.validation.targetBalanceUnchanged'))
    return
  }

  const bucketId = selectedBucket.value.id
  const payload = {
    scenario: 'investment_revalue',
    bucketId,
    amount: formatPlainAmount(delta),
    currency: selectedBucket.value.currency,
    description: t('buckets.revalueDialog.description', { from: formatAmount(selectedBucket.value.balance), to: formatAmount(target) }),
    remark: revalueForm.remark,
  }
  const res = await createRecord(payload)
  if (res.success) {
    ElMessage.success(t('buckets.messages.revalued'))
    revalueDialogVisible.value = false
    await loadBuckets()
    const updated = buckets.value.find((item) => item.id === bucketId)
    if (updated) await selectBucket(updated)
  } else {
    ElMessage.error(res.error || t('buckets.messages.revalueFailed'))
  }
}

onMounted(refreshAll)
onActivated(refreshAll)
</script>

<style scoped>
.ledger-page {
  max-width: 1200px;
  margin: 0 auto;
  color: #1e293b;
}

.reveal-block,
.bucket-card {
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
.bucket-card,
.ledger-panel,
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

.page-hero h1,
.section-head h2 {
  margin: 0;
  letter-spacing: -0.022em;
}

.page-hero h1 {
  max-width: 680px;
  font-size: 30px;
  line-height: 1.16;
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

.export-action {
  min-height: 40px;
  border: 0;
  border-radius: 12px;
  padding: 0 14px;
  background: rgba(47, 125, 92, 0.10);
  color: #2f7d5c;
  font-size: 12px;
  font-weight: 900;
  cursor: pointer;
  box-shadow: inset 0 0 0 1px rgba(47, 125, 92, 0.22);
}

.export-action:active { transform: scale(0.96); }

.filter-control {
  width: 180px;
}

.bucket-layout {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 360px;
  gap: 18px;
}

.bucket-list {
  display: grid;
  gap: 14px;
}

.bucket-card {
  padding: 20px;
  cursor: pointer;
  transition-property: transform, box-shadow, background-color;
  transition-duration: 180ms;
}

.bucket-card.active {
  background: #f8faf7;
  box-shadow: 0 0 0 2px rgba(47, 125, 92, 0.16), 0 14px 34px rgba(47, 125, 92, 0.1);
}

.bucket-topline,
.section-head,
.entry-row {
  display: flex;
  justify-content: space-between;
  gap: 14px;
}

.selected-actions {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
  margin-top: 14px;
}

.edit-selected-action,
.adjust-selected-action {
  min-height: 38px;
  border: 0;
  border-radius: 12px;
  font-weight: 800;
  cursor: pointer;
}

.edit-selected-action {
  background: #f4efe6;
  color: #6b5b49;
  box-shadow: none;
}

.adjust-selected-action {
  background: rgba(47, 125, 92, 0.1);
  color: #2f7d5c;
}

.adjust-selected-action:disabled {
  cursor: not-allowed;
  opacity: 0.5;
}

.edit-selected-action:active,
.adjust-selected-action:active:not(:disabled) {
  transform: scale(0.96);
}

.bucket-topline h2 {
  margin: 0;
  font-size: 20px;
  letter-spacing: -0.012em;
}

.bucket-topline p,
.selected-summary p {
  margin: 6px 0 0;
  color: #64748b;
}

.nature-pill,
.bucket-meta span {
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

.nature-pill.liability {
  background: rgba(31, 41, 51, 0.1);
  color: #1f2933;
}

.bucket-balance {
  display: flex;
  align-items: baseline;
  gap: 12px;
  margin-top: 18px;
}

.bucket-balance span,
.selected-summary span {
  color: #64748b;
  font-weight: 700;
}

.bucket-balance strong,
.selected-summary strong,
.entry-amount {
  font-family: 'SF Mono', 'Fira Code', monospace;
  font-variant-numeric: tabular-nums;
}

.bucket-balance strong {
  font-size: 30px;
  letter-spacing: -0.022em;
}

.bucket-meta {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  margin-top: 16px;
}

.ledger-panel {
  position: sticky;
  top: 24px;
  align-self: start;
  padding: 22px;
}

.selected-summary {
  margin: 18px 0;
  padding: 18px;
  border-radius: 14px;
  background: #f8faf7;
}

.selected-summary strong {
  display: block;
  margin-top: 6px;
  font-size: 28px;
}

.entry-section {
  margin-top: 18px;
}

.entry-list {
  display: grid;
  gap: 10px;
  max-height: 396px;
  overflow: hidden;
}

.entry-list.expanded {
  max-height: 520px;
  overflow-y: auto;
  padding-right: 4px;
  overscroll-behavior: contain;
}

.entry-row {
  align-items: center;
  padding: 12px;
  border-radius: 12px;
  background: #ffffff;
  box-shadow: inset 0 0 0 1px rgba(100, 116, 139, 0.12);
}

.entry-toggle {
  width: 100%;
  min-height: 36px;
  margin-top: 10px;
  border: 0;
  border-radius: 10px;
  background: #f8faf7;
  color: #1e293b;
  font-size: 13px;
  font-weight: 700;
  cursor: pointer;
  transition-property: transform, background-color, color;
  transition-duration: 160ms;
  transition-timing-function: cubic-bezier(0.16, 1, 0.3, 1);
  touch-action: manipulation;
}

.entry-toggle:active {
  transform: scale(0.96);
}

.entry-row strong,
.entry-row span {
  display: block;
}

.entry-row span {
  margin-top: 4px;
  color: #94a3b8;
  font-size: 12px;
}

.entry-amount {
  color: #1e293b;
  font-weight: 700;
}

.adjust-summary,
.adjust-diff {
  margin-bottom: 18px;
  padding: 16px;
  border-radius: 14px;
  background: #f8faf7;
}

.adjust-summary span,
.adjust-diff span,
.adjust-diff small {
  display: block;
  color: #64748b;
  font-size: 12px;
  font-weight: 800;
}

.adjust-summary strong,
.adjust-diff strong {
  display: block;
  margin-top: 8px;
  color: #1e293b;
  font-family: 'SF Mono', 'Fira Code', monospace;
  font-size: 22px;
}

.adjust-diff.positive strong {
  color: #ef4444;
}

.adjust-diff.negative strong {
  color: #f97316;
}

.adjust-diff.neutral strong {
  color: #64748b;
}

.adjust-diff small {
  margin-top: 6px;
  font-weight: 700;
}

.primary-action,
.ghost-action {
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
.bucket-card:active {
  transform: scale(0.96);
}

.primary-action {
  background: #2f7d5c;
  color: #ffffff;
  box-shadow: 0 10px 24px rgba(47, 125, 92, 0.20);
}

.primary-action:disabled {
  cursor: not-allowed;
  opacity: 0.55;
}

.ghost-action {
  background: #f8faf7;
  color: #1e293b;
}

.empty-state {
  display: grid;
  place-items: center;
  gap: 12px;
  padding: 40px 24px;
  text-align: center;
  color: #64748b;
}

.empty-state.compact {
  padding: 26px 12px;
  box-shadow: none;
  background: #f8faf7;
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

@media (hover: hover) {
  .edit-selected-action:hover {
    background: #ece2d2;
    color: #4b3f33;
  }

  .bucket-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 3px 8px rgba(15, 23, 42, 0.12), 0 16px 34px rgba(15, 23, 42, 0.06);
  }
}

@media (max-width: 980px) {
  .page-hero,
  .toolbar {
    flex-direction: column;
    align-items: stretch;
  }

  .bucket-layout {
    grid-template-columns: 1fr;
  }

  .ledger-panel {
    position: static;
  }

  .filter-control {
    width: 100%;
  }

  .export-action,
  .ghost-action,
  .primary-action {
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

  .bucket-topline,
  .entry-row {
    flex-direction: column;
  }
}

@media (prefers-reduced-motion: reduce) {
  .reveal-block,
  .bucket-card,
  .primary-action,
  .ghost-action {
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
