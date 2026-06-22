<template>
  <main class="record-page">
    <section class="page-hero reveal-block">
      <div>
        <p class="eyebrow">{{ t('record.hero.eyebrow') }}</p>
        <h1>{{ t('record.hero.title') }}</h1>
        <p>{{ t('record.hero.subtitle') }}</p>
      </div>
      <img :src="marmotOne" :alt="t('dashboard.images.brandAlt')" width="92" height="92" />
    </section>

    <section class="record-layout reveal-block delay-1">
      <aside class="scenario-panel">
        <template v-for="group in groupedScenarios" :key="group.value">
          <p v-if="group.value === 'core' || moreScenariosOpen" class="scenario-group-label">{{ t(group.labelKey) }}</p>
          <template v-if="group.value === 'core' || moreScenariosOpen">
            <button
              v-for="item in group.items"
              :key="item.value"
              class="scenario-card"
              :class="[{ active: form.scenario === item.value }, item.value, item.group, { disabled: item.disabled }]"
              :disabled="item.disabled"
              @click="selectScenario(item.value)"
            >
              <span>{{ item.icon }}</span>
              <strong>{{ t(item.labelKey) }}</strong>
              <small v-if="item.disabled">{{ t('record.scenarios.comingSoon') }}</small>
            </button>
          </template>
        </template>

        <button class="more-scenarios-toggle" type="button" @click="moreScenariosOpen = !moreScenariosOpen">
          {{ moreScenariosOpen ? t('record.actions.collapseScenarios') : t('record.actions.moreScenarios') }}
          <span class="more-scenarios-icon">{{ moreScenariosOpen ? '▲' : '▼' }}</span>
        </button>
      </aside>

      <section class="form-panel">
        <div class="section-head">
          <div>
            <p class="eyebrow">{{ t(selectedScenario.labelKey) }}</p>
            <h2>{{ t('record.actions.submit') }}</h2>
          </div>
        </div>

        <el-form ref="formRef" :model="form" label-position="top">
          <template v-if="isExchange">
            <el-form-item :label="t('record.fields.fromBucket')">
              <el-select v-model="form.fromBucketId" class="full-width" filterable @change="syncFromCurrency">
                <el-option v-for="bucket in buckets" :key="bucket.id" :label="bucketLabel(bucket)" :value="bucket.id" />
              </el-select>
            </el-form-item>
            <el-row :gutter="12">
              <el-col :span="14">
                <el-form-item :label="t('record.fields.fromAmount')">
                  <el-input v-model="form.amount" placeholder="100.00" />
                </el-form-item>
              </el-col>
              <el-col :span="10">
                <el-form-item :label="t('record.fields.fromCurrency')">
                  <el-select v-model="form.currency" class="full-width" disabled>
                    <el-option v-for="item in currencyOptions" :key="item.code" :label="getCurrencyLabel(item.code, config.locale)" :value="item.code" />
                  </el-select>
                </el-form-item>
              </el-col>
            </el-row>

            <el-form-item :label="t('record.fields.toBucket')">
              <el-select v-model="form.toBucketId" class="full-width" filterable @change="syncToCurrency">
                <el-option v-for="bucket in buckets" :key="bucket.id" :label="bucketLabel(bucket)" :value="bucket.id" />
              </el-select>
            </el-form-item>
            <el-row :gutter="12">
              <el-col :span="14">
                <el-form-item :label="t('record.fields.toAmount')">
                  <el-input v-model="form.toAmount" placeholder="720.00" />
                </el-form-item>
              </el-col>
              <el-col :span="10">
                <el-form-item :label="t('record.fields.toCurrency')">
                  <el-select v-model="form.toCurrency" class="full-width" disabled>
                    <el-option v-for="item in currencyOptions" :key="item.code" :label="getCurrencyLabel(item.code, config.locale)" :value="item.code" />
                  </el-select>
                </el-form-item>
              </el-col>
            </el-row>

            <p class="exchange-rate-hint">{{ rateHint }}</p>
          </template>

          <template v-else-if="isSplit">
            <el-row :gutter="12">
              <el-col :span="14">
                <el-form-item :label="t('record.fields.totalAmount')">
                  <el-input v-model="splitForm.totalAmount" placeholder="320.00" />
                </el-form-item>
              </el-col>
              <el-col :span="10">
                <el-form-item :label="t('record.fields.currency')">
                  <el-select v-model="form.currency" class="full-width">
                    <el-option v-for="item in currencyOptions" :key="item.code" :label="getCurrencyLabel(item.code, config.locale)" :value="item.code" />
                  </el-select>
                </el-form-item>
              </el-col>
            </el-row>

            <el-form-item :label="t('record.fields.cashFrom')">
              <el-select v-model="splitForm.cashBucketId" class="full-width" filterable :placeholder="t('record.fields.bucket')">
                <el-option v-for="bucket in cashBuckets" :key="bucket.id" :label="bucketLabel(bucket)" :value="bucket.id" />
              </el-select>
            </el-form-item>

            <el-form-item :label="t('record.fields.counterparty.receivable')">
              <el-select v-if="counterpartyBuckets.length" v-model="splitForm.receivableBucketId" class="full-width" filterable :placeholder="t('record.fields.bucket')">
                <el-option v-for="bucket in counterpartyBuckets" :key="bucket.id" :label="bucketLabel(bucket)" :value="bucket.id" />
              </el-select>
              <button v-else class="counterparty-empty" type="button" @click="goCreateBucket">
                <span class="counterparty-empty-icon">+</span>
                <span class="counterparty-empty-text">
                  <strong>{{ t('record.messages.noCounterpartyBucket') }}</strong>
                  <small>{{ t('record.actions.goCreateBucket') }} →</small>
                </span>
              </button>
            </el-form-item>

            <el-form-item :label="t('record.fields.description')">
              <el-input v-model="form.description" :placeholder="t('record.placeholders.splitDescription')" />
            </el-form-item>

            <el-form-item v-if="hasSelfShare" :label="t('record.fields.category')">
              <el-select v-model="form.categoryId" class="full-width" filterable :placeholder="t('record.placeholders.selectCategory')">
                <el-option v-for="category in categories" :key="category.id" :label="categoryLabel(category)" :value="category.id" />
              </el-select>
            </el-form-item>

            <div class="split-shares">
              <p class="split-shares-label">{{ t('record.fields.shares') }}</p>
              <div v-for="(share, idx) in splitForm.shares" :key="idx" class="split-share-row" :class="{ self: share.isSelf }">
                <div class="split-share-head">
                  <strong>{{ share.isSelf ? t('record.fields.shareSelf') : t('record.fields.shareFriend') }}</strong>
                  <button v-if="!share.isSelf" class="split-share-remove" type="button" @click="removeShare(idx)">×</button>
                </div>
                <el-input v-if="!share.isSelf" v-model="share.description" :placeholder="t('record.placeholders.shareDescription')" class="split-share-name" />
                <el-input v-model="share.amount" placeholder="0.00" class="split-share-amount" />
              </div>

              <div class="split-actions">
                <button class="split-chip primary" type="button" @click="addFriendShare">
                  <span class="split-chip-icon">+</span>
                  {{ t('record.actions.addShare') }}
                </button>
                <button class="split-chip" type="button" :disabled="!hasSelfShare || !splitForm.shares.some((s) => !s.isSelf)" @click="splitEvenly">
                  <span class="split-chip-icon">⊟</span>
                  {{ t('record.actions.splitEvenly') }}
                </button>
                <button class="split-chip" type="button" :disabled="!splitForm.shares.length" @click="splitEvenlyAll">
                  <span class="split-chip-icon">⊞</span>
                  {{ t('record.actions.splitEvenlyAll') }}
                </button>
                <button class="split-chip subtle" type="button" @click="toggleSelfShare">
                  <span class="split-chip-icon">{{ hasSelfShare ? '−' : '+' }}</span>
                  {{ hasSelfShare ? t('record.actions.removeSelf') : t('record.actions.addSelf') }}
                </button>
              </div>

              <div class="split-summary" :class="{ balanced: Math.abs(splitBalanceDiff) <= 0.001, off: Math.abs(splitBalanceDiff) > 0.001 }">
                <span>{{ t('record.fields.allocated') }} <strong>{{ form.currency }} {{ formatAmount(splitAllocated) }}</strong> / {{ form.currency }} {{ formatAmount(splitForm.totalAmount) }}</span>
                <span v-if="Math.abs(splitBalanceDiff) > 0.001" class="split-diff">{{ splitBalanceDiff < 0 ? t('record.messages.splitOver', { diff: formatAmount(splitDiffAbs) }) : t('record.messages.splitUnder', { diff: formatAmount(splitDiffAbs) }) }}</span>
                <span v-else class="split-balanced">✓</span>
              </div>
            </div>
          </template>

          <template v-else-if="isTransfer">
            <el-form-item :label="t('record.fields.fromBucket')">
              <el-select v-model="form.fromBucketId" class="full-width" filterable>
                <el-option v-for="bucket in buckets" :key="bucket.id" :label="bucketLabel(bucket)" :value="bucket.id" />
              </el-select>
            </el-form-item>
            <el-form-item :label="t('record.fields.toBucket')">
              <el-select v-model="form.toBucketId" class="full-width" filterable>
                <el-option v-for="bucket in buckets" :key="bucket.id" :label="bucketLabel(bucket)" :value="bucket.id" />
              </el-select>
            </el-form-item>
          </template>

          <template v-else-if="isPaired">
            <el-form-item :label="t(isCashToCounter ? 'record.fields.cashFrom' : 'record.fields.cashTo')">
              <el-select v-model="form.fromBucketId" class="full-width" filterable :placeholder="t('record.fields.bucket')">
                <el-option v-for="bucket in cashBuckets" :key="bucket.id" :label="bucketLabel(bucket)" :value="bucket.id" />
              </el-select>
            </el-form-item>
            <el-form-item :label="t(`record.fields.counterparty.${counterpartyType}`)">
              <el-select v-if="counterpartyBuckets.length" v-model="form.toBucketId" class="full-width" filterable :placeholder="t('record.fields.bucket')">
                <el-option v-for="bucket in counterpartyBuckets" :key="bucket.id" :label="bucketLabel(bucket)" :value="bucket.id" />
              </el-select>
              <button v-else class="counterparty-empty" type="button" @click="goCreateBucket">
                <span class="counterparty-empty-icon">+</span>
                <span class="counterparty-empty-text">
                  <strong>{{ t('record.messages.noCounterpartyBucket') }}</strong>
                  <small>{{ t('record.actions.goCreateBucket') }} →</small>
                </span>
              </button>
            </el-form-item>

            <div v-if="!isCashToCounter && form.toBucketId && outstandingItems.length" class="outstanding-list">
              <p class="outstanding-label">{{ t('record.fields.outstanding') }}</p>
              <button type="button" class="outstanding-row" :class="{ active: !form.relatedFinancialEventId }" @click="selectOutstanding(null)">
                <span class="outstanding-mark">{{ !form.relatedFinancialEventId ? '●' : '○' }}</span>
                <span class="outstanding-main">
                  <strong>{{ t('record.fields.outstandingFreeForm') }}</strong>
                  <small>{{ t('record.fields.outstandingFreeFormHint') }}</small>
                </span>
              </button>
              <button v-for="item in outstandingItems" :key="item.id" type="button" class="outstanding-row" :class="{ active: form.relatedFinancialEventId === item.id }" @click="selectOutstanding(item)">
                <span class="outstanding-mark">{{ form.relatedFinancialEventId === item.id ? '●' : '○' }}</span>
                <span class="outstanding-main">
                  <strong>{{ item.description || t('record.scenarios.' + item.eventType) }}</strong>
                  <small>{{ item.eventTime.split(' ')[0] }} · {{ t('record.fields.outstandingAmount') }} {{ item.currency }} {{ formatAmount(item.outstandingAmount) }}</small>
                </span>
                <span class="outstanding-original">{{ item.currency }} {{ formatAmount(item.amount) }}</span>
              </button>
            </div>
          </template>

          <template v-else-if="isInvestmentBuy">
            <el-form-item :label="t('record.fields.cashFrom')">
              <el-select v-model="form.fromBucketId" class="full-width" filterable :placeholder="t('record.fields.bucket')">
                <el-option v-for="bucket in investmentCashBuckets" :key="bucket.id" :label="bucketLabel(bucket)" :value="bucket.id" />
              </el-select>
            </el-form-item>
            <el-form-item :label="t('record.fields.toBucket')">
              <el-select v-model="form.toBucketId" class="full-width" filterable :placeholder="t('record.fields.bucket')">
                <el-option v-for="bucket in investmentBuckets" :key="bucket.id" :label="bucketLabel(bucket)" :value="bucket.id" />
              </el-select>
            </el-form-item>
          </template>

          <template v-else-if="isInvestmentSell">
            <el-form-item :label="t('record.fields.investmentBucketSell')">
              <el-select v-model="sellForm.investBucketId" class="full-width" filterable :placeholder="t('record.fields.bucket')">
                <el-option v-for="bucket in investmentBuckets" :key="bucket.id" :label="bucketLabel(bucket)" :value="bucket.id" />
              </el-select>
            </el-form-item>
            <el-form-item :label="t('record.fields.cashTo')">
              <el-select v-model="sellForm.cashBucketId" class="full-width" filterable :placeholder="t('record.fields.bucket')">
                <el-option v-for="bucket in investmentCashBuckets" :key="bucket.id" :label="bucketLabel(bucket)" :value="bucket.id" />
              </el-select>
            </el-form-item>

            <div class="sell-mode-switch" role="group">
              <button type="button" :class="{ active: sellForm.fullSell }" @click="sellForm.fullSell = true">{{ t('record.actions.sellAll') }}</button>
              <button type="button" :class="{ active: !sellForm.fullSell }" @click="sellForm.fullSell = false">{{ t('record.actions.sellPartial') }}</button>
            </div>

            <el-row :gutter="12">
              <el-col :span="14">
                <el-form-item :label="t('record.fields.sellReceivedAmount')">
                  <el-input v-model="sellForm.received" placeholder="10300.00" />
                </el-form-item>
              </el-col>
              <el-col :span="10">
                <el-form-item :label="t('record.fields.currency')">
                  <el-select v-model="form.currency" class="full-width">
                    <el-option v-for="item in currencyOptions" :key="item.code" :label="getCurrencyLabel(item.code, config.locale)" :value="item.code" />
                  </el-select>
                </el-form-item>
              </el-col>
            </el-row>

            <el-form-item v-if="!sellForm.fullSell" :label="t('record.fields.remainingMarketValue')">
              <el-input v-model="sellForm.remaining" placeholder="5150.00" />
            </el-form-item>

            <p class="exchange-rate-hint" :class="sellPnlClass">{{ sellPnlHint }}</p>
          </template>

          <template v-else-if="isInvestmentIncome">
            <el-form-item :label="t('record.fields.cashTo')">
              <el-select v-model="form.bucketId" class="full-width" filterable :placeholder="t('record.fields.bucket')">
                <el-option v-for="bucket in investmentCashBuckets" :key="bucket.id" :label="bucketLabel(bucket)" :value="bucket.id" />
              </el-select>
            </el-form-item>
          </template>

          <template v-else-if="isFamilyTransfer">
            <el-form-item :label="t('record.fields.family')">
              <el-select v-model="familyTransferForm.familyId" class="full-width" @change="loadFamilyAssetsForTransfer">
                <el-option v-for="family in families" :key="family.id" :label="family.name" :value="family.id" />
              </el-select>
            </el-form-item>
            <el-form-item :label="t('record.fields.fromBucket')">
              <el-select v-model="form.fromBucketId" class="full-width" filterable>
                <el-option v-for="bucket in buckets.filter((b) => !form.currency || b.currency === form.currency)" :key="bucket.id" :label="bucketLabel(bucket)" :value="bucket.id" />
              </el-select>
            </el-form-item>
            <el-form-item :label="t('record.fields.toMember')">
              <el-select v-model="familyTransferForm.toUserId" class="full-width" filterable>
                <el-option v-for="member in familyTransferMembers" :key="member.userId" :label="member.displayName || member.name || member.account" :value="member.userId" />
              </el-select>
            </el-form-item>
            <el-form-item :label="t('record.fields.toBucket')">
              <el-select v-model="familyTransferForm.toBucketId" class="full-width" filterable>
                <el-option v-for="bucket in familyTransferTargetBuckets" :key="bucket.id" :label="bucketLabel(bucket)" :value="bucket.id" />
              </el-select>
            </el-form-item>
          </template>

          <el-form-item v-if="!isExchange && !isSplit && !isTransfer && !isPaired && !isInvestmentScenario && !isFamilyTransfer" :label="t('record.fields.bucket')">
            <el-select v-model="form.bucketId" class="full-width" filterable>
              <el-option v-for="bucket in buckets" :key="bucket.id" :label="bucketLabel(bucket)" :value="bucket.id" />
            </el-select>
          </el-form-item>

          <el-form-item v-if="needsCategory" :label="t('record.fields.category')">
            <el-select v-model="form.categoryId" class="full-width" filterable :placeholder="t('record.placeholders.selectCategory')">
              <el-option v-for="category in categories" :key="category.id" :label="categoryLabel(category)" :value="category.id" />
            </el-select>
          </el-form-item>

          <el-form-item v-if="showChannelSelect" :label="t('record.fields.channel')">
            <el-select v-model="form.channelId" clearable class="full-width" filterable>
              <el-option v-for="channel in filteredChannels" :key="channel.id" :label="channelLabel(channel)" :value="channel.id" />
            </el-select>
          </el-form-item>

          <el-row v-if="!isExchange && !isSplit && !isInvestmentSell" :gutter="12">
            <el-col :span="12">
              <el-form-item :label="t('record.fields.amount')">
                <el-input v-model="form.amount" placeholder="36.50" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item :label="t('record.fields.currency')">
                <el-select v-model="form.currency" class="full-width">
                  <el-option v-for="item in currencyOptions" :key="item.code" :label="getCurrencyLabel(item.code, config.locale)" :value="item.code" />
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>

          <el-form-item :label="t('record.fields.description')">
            <el-input v-model="form.description" />
          </el-form-item>
          <el-form-item :label="t('record.fields.eventTime')">
            <el-input v-model="form.eventTime" placeholder="YYYY-MM-DD HH:mm:ss" />
          </el-form-item>
          <el-form-item v-if="form.scenario === 'refund'" :label="t('record.fields.relatedEvent')">
            <el-input v-model="form.relatedFinancialEventId" />
          </el-form-item>
          <el-form-item :label="t('record.fields.remark')">
            <el-input v-model="form.remark" type="textarea" :autosize="{ minRows: 2, maxRows: 4 }" />
          </el-form-item>

          <button class="primary-action" type="button" :disabled="submitting" @click="submitRecord">
            {{ submitting ? '...' : t('record.actions.submit') }}
          </button>
        </el-form>
      </section>

      <aside class="preview-panel">
        <div class="section-head">
          <div>
            <p class="eyebrow">{{ t('record.preview.title') }}</p>
            <h2>{{ t('record.preview.event') }}</h2>
          </div>
        </div>

        <div class="event-preview" :class="form.scenario">
          <span>{{ eventTypeLabel(form.scenario) }}</span>
          <strong>{{ form.currency }} {{ formatAmount(form.amount) }}</strong>
          <p>{{ includeInStatistics ? t('record.preview.included') : t('record.preview.excluded') }}</p>
        </div>

        <div class="entry-preview">
          <h3>{{ t('record.preview.entries') }}</h3>
          <div v-if="previewEntries.length" class="entry-list">
            <div v-for="entry in previewEntries" :key="entry.role + entry.bucket" class="entry-row">
              <span>{{ entry.bucket }}</span>
              <strong :class="{ negative: entry.amount.startsWith('-'), positive: entry.amount.startsWith('+') }">{{ entry.amount }}</strong>
            </div>
          </div>
          <p v-else>{{ t('record.preview.selectScenario') }}</p>
        </div>

        <div class="recent-events">
          <h3>{{ t('dashboard.events.title') }}</h3>
          <div v-for="event in recentEvents" :key="event.id" class="recent-row">
            <span>{{ eventTypeLabel(event.eventType) }}</span>
            <strong>{{ event.currency }} {{ formatAmount(event.amount) }}</strong>
          </div>
        </div>
      </aside>
    </section>
  </main>
</template>

<script setup>
import { computed, onActivated, onMounted, reactive, ref, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { useConfigStore } from '@/stores/config'
import { listBuckets } from '@/api/bucket/bucket'
import { listFinancialEvents, listOutstandingForBucket } from '@/api/financialEvent/financialEvent'
import { createRecord } from '@/api/record/record'
import { listCategories } from '@/api/category/category'
import { listChannels } from '@/api/channel/channel'
import { getFamilyAssets, listFamilies } from '@/api/family/family'
import marmotOne from '../../../img/marmot-ledger-1.png'
import { currencyOptions, getCurrencyLabel } from '@/utils/currency'
import { getBucketEmoji } from '@/utils/bucketEmoji'

const { t, te } = useI18n()
const router = useRouter()
const route = useRoute()
const config = useConfigStore()
const buckets = ref([])
const categories = ref([])
const channels = ref([])
const recentEvents = ref([])
const families = ref([])
const familyAssets = ref(null)
const submitting = ref(false)
const formRef = ref()
const moreScenariosOpen = ref(false)

const form = reactive({
  scenario: 'expense',
  bucketId: '',
  fromBucketId: '',
  toBucketId: '',
  categoryId: '',
  channelId: '',
  amount: '',
  currency: 'CNY',
  toAmount: '',
  toCurrency: 'CNY',
  description: '',
  eventTime: formatDateTime(new Date()),
  relatedFinancialEventId: '',
  remark: '',
})

const splitForm = reactive({
  totalAmount: '',
  cashBucketId: '',
  receivableBucketId: '',
  shares: [
    { isSelf: true, amount: '', description: '' },
    { isSelf: false, amount: '', description: '' },
  ],
})

const sellForm = reactive({
  investBucketId: '',
  cashBucketId: '',
  fullSell: true,
  received: '',
  remaining: '',
})

const familyTransferForm = reactive({ familyId: '', toUserId: '', toBucketId: '' })

const outstandingItems = ref([])
const relatedCreateTypeBySceneFront = {
  receivable_collect: 'receivable_create',
  deposit_refund: 'deposit_create',
  loan_collect: 'loan_out',
}

const scenarios = [
  { value: 'income', labelKey: 'record.scenarios.income', icon: '↙', group: 'core', disabled: false },
  { value: 'expense', labelKey: 'record.scenarios.expense', icon: '↗', group: 'core', disabled: false },
  { value: 'transfer', labelKey: 'record.scenarios.transfer', icon: '⇄', group: 'core', disabled: false },
  { value: 'refund', labelKey: 'record.scenarios.refund', icon: '↩', group: 'core', disabled: false },
  { value: 'exchange', labelKey: 'record.scenarios.exchange', icon: 'FX', group: 'core', disabled: false },
  { value: 'receivable_create', labelKey: 'record.scenarios.receivable_create', icon: 'R+', group: 'paired', disabled: false },
  { value: 'receivable_collect', labelKey: 'record.scenarios.receivable_collect', icon: 'R-', group: 'paired', disabled: false },
  { value: 'deposit_create', labelKey: 'record.scenarios.deposit_create', icon: 'D+', group: 'paired', disabled: false },
  { value: 'deposit_refund', labelKey: 'record.scenarios.deposit_refund', icon: 'D-', group: 'paired', disabled: false },
  { value: 'loan_out', labelKey: 'record.scenarios.loan_out', icon: 'L+', group: 'paired', disabled: false },
  { value: 'loan_collect', labelKey: 'record.scenarios.loan_collect', icon: 'L-', group: 'paired', disabled: false },
  { value: 'split', labelKey: 'record.scenarios.split', icon: '⊟', group: 'paired', disabled: false },
  { value: 'investment_buy', labelKey: 'record.scenarios.investment_buy', icon: '📈+', group: 'investment', disabled: false },
  { value: 'investment_sell', labelKey: 'record.scenarios.investment_sell', icon: '📈-', group: 'investment', disabled: false },
  { value: 'investment_income', labelKey: 'record.scenarios.investment_income', icon: '💵', group: 'investment', disabled: false },
  { value: 'family_transfer', labelKey: 'record.scenarios.family_transfer', icon: '⇄', group: 'family', disabled: false },
]

const scenarioGroups = [
  { value: 'core', labelKey: 'record.scenarios.groups.core' },
  { value: 'paired', labelKey: 'record.scenarios.groups.paired' },
  { value: 'investment', labelKey: 'record.scenarios.groups.investment' },
  { value: 'family', labelKey: 'record.scenarios.groups.family' },
]

const groupedScenarios = computed(() => scenarioGroups.map((g) => ({
  ...g,
  items: scenarios.filter((item) => item.group === g.value),
})))

const selectedScenario = computed(() => scenarios.find((item) => item.value === form.scenario) || scenarios[1])
const isTransfer = computed(() => form.scenario === 'transfer')
const isExchange = computed(() => form.scenario === 'exchange')
const isSplit = computed(() => form.scenario === 'split')
const pairedScenarios = ['receivable_create', 'receivable_collect', 'deposit_create', 'deposit_refund', 'loan_out', 'loan_collect']
const cashToCounterScenarios = ['receivable_create', 'deposit_create', 'loan_out']
const isPaired = computed(() => pairedScenarios.includes(form.scenario))
const isCashToCounter = computed(() => cashToCounterScenarios.includes(form.scenario))
const counterpartyType = computed(() => {
  if (form.scenario.startsWith('receivable') || form.scenario === 'split') return 'receivable'
  if (form.scenario.startsWith('deposit')) return 'deposit'
  if (form.scenario === 'loan_out' || form.scenario === 'loan_collect') return 'loan_out'
  return ''
})
const counterEntryRoleByScenario = {
  receivable_create: 'receivable_create',
  receivable_collect: 'receivable_collect',
  deposit_create: 'deposit_create',
  deposit_refund: 'deposit_refund',
  loan_out: 'loan_out',
  loan_collect: 'loan_collect',
}
const cashSideBucketTypes = ['cash', 'wallet', 'bank', 'virtual']
const investCashSideBucketTypes = ['cash', 'wallet', 'bank', 'virtual', 'investment_cash']
const investmentBucketTypeList = ['investment_cash', 'investment_asset']
const isInvestmentBuy = computed(() => form.scenario === 'investment_buy')
const isInvestmentSell = computed(() => form.scenario === 'investment_sell')
const isInvestmentIncome = computed(() => form.scenario === 'investment_income')
const isInvestmentScenario = computed(() => isInvestmentBuy.value || isInvestmentSell.value || isInvestmentIncome.value)
const isFamilyTransfer = computed(() => form.scenario === 'family_transfer')
const currentUserId = computed(() => buckets.value[0]?.userId || 0)
const familyTransferMembers = computed(() => (familyAssets.value?.members || []).filter((member) => Number(member.userId) !== Number(currentUserId.value)))
const familyTransferTargetBuckets = computed(() => {
  const member = familyTransferMembers.value.find((m) => Number(m.userId) === Number(familyTransferForm.toUserId))
  return (member?.accounts || []).flatMap((a) => a.buckets || []).filter((b) => !form.currency || b.currency === form.currency)
})
const cashBuckets = computed(() => buckets.value.filter((b) => cashSideBucketTypes.includes(b.bucketType) && b.bucketNature === 'asset' && (!form.currency || b.currency === form.currency)))
const counterpartyBuckets = computed(() => buckets.value.filter((b) => b.bucketType === counterpartyType.value && b.bucketNature === 'asset' && (!form.currency || b.currency === form.currency)))
const investmentCashBuckets = computed(() => buckets.value.filter((b) => investCashSideBucketTypes.includes(b.bucketType) && b.bucketNature === 'asset' && (!form.currency || b.currency === form.currency)))
const investmentBuckets = computed(() => buckets.value.filter((b) => investmentBucketTypeList.includes(b.bucketType) && b.bucketNature === 'asset' && (!form.currency || b.currency === form.currency)))
const sellInvestBucket = computed(() => buckets.value.find((b) => b.id === Number(sellForm.investBucketId)))
const needsCategory = computed(() => form.scenario === 'income' || form.scenario === 'expense' || form.scenario === 'refund')
const showChannelSelect = computed(() => ['income', 'expense', 'refund', 'transfer', 'exchange'].includes(form.scenario))
const filteredChannels = computed(() => channels.value.filter((item) => !item.supportedEventTypes || item.supportedEventTypes.split(',').includes(form.scenario)))
const includeInStatistics = computed(() => form.scenario === 'income' || form.scenario === 'expense')

const splitAllocated = computed(() => splitForm.shares.reduce((sum, s) => sum + (Number(s.amount) || 0), 0))
const splitBalanceDiff = computed(() => Number(splitForm.totalAmount || 0) - splitAllocated.value)
const splitDiffAbs = computed(() => Math.abs(splitBalanceDiff.value))
const hasSelfShare = computed(() => splitForm.shares.some((s) => s.isSelf))
const cashBucketName = computed(() => {
  const b = buckets.value.find((x) => x.id === Number(splitForm.cashBucketId))
  return b ? `${getBucketEmoji(b.bucketType)} ${b.name}` : '-'
})
const receivableBucketName = computed(() => {
  const b = buckets.value.find((x) => x.id === Number(splitForm.receivableBucketId))
  return b ? `${getBucketEmoji(b.bucketType)} ${b.name}` : '-'
})

const rateHint = computed(() => {
  if (!isExchange.value) return ''
  const fromAmount = Number(form.amount || 0)
  const toAmount = Number(form.toAmount || 0)
  if (!fromAmount || !toAmount || !form.currency || !form.toCurrency) return t('record.preview.ratePending')
  const rate = (toAmount / fromAmount).toFixed(4)
  return t('record.preview.rateLine', { from: form.currency, rate, to: form.toCurrency })
})

const sellPnl = computed(() => {
  if (!isInvestmentSell.value) return null
  const bucket = sellInvestBucket.value
  if (!bucket) return null
  const received = Number(sellForm.received || 0)
  const remaining = sellForm.fullSell ? 0 : Number(sellForm.remaining || 0)
  if (!received) return null
  const targetMV = received + remaining
  return targetMV - Number(bucket.balance)
})
const sellPnlHint = computed(() => {
  const pnl = sellPnl.value
  if (pnl === null) return t('record.preview.sellPnlPending')
  if (Math.abs(pnl) < 0.001) return t('record.preview.sellPnlEqual')
  if (pnl > 0) return t('record.preview.sellPnlGain', { amount: formatAmount(pnl) })
  return t('record.preview.sellPnlLoss', { amount: formatAmount(-pnl) })
})
const sellPnlClass = computed(() => {
  const pnl = sellPnl.value
  if (pnl === null || Math.abs(pnl) < 0.001) return 'neutral'
  return pnl > 0 ? 'gain' : 'loss'
})

const previewEntries = computed(() => {
  if (isSplit.value) {
    const lines = []
    splitForm.shares.forEach((share) => {
      const amount = Number(share.amount || 0)
      if (!amount) return
      lines.push({ role: share.isSelf ? 'expense' : 'cash_leg', bucket: cashBucketName.value, amount: `-${form.currency} ${formatAmount(amount)}` })
      if (!share.isSelf) {
        lines.push({ role: 'receivable_create', bucket: receivableBucketName.value, amount: `+${form.currency} ${formatAmount(amount)}` })
      }
    })
    return lines
  }
  if (isInvestmentBuy.value) {
    const a = Number(form.amount || 0)
    if (!a) return []
    return [
      { role: 'cash_leg', bucket: bucketDisplay(form.fromBucketId), amount: `-${form.currency} ${formatAmount(a)}` },
      { role: 'investment_buy', bucket: bucketDisplay(form.toBucketId), amount: `+${form.currency} ${formatAmount(a)}` },
    ]
  }
  if (isInvestmentSell.value) {
    const lines = []
    const received = Number(sellForm.received || 0)
    const remaining = sellForm.fullSell ? 0 : Number(sellForm.remaining || 0)
    const bucket = sellInvestBucket.value
    if (!received || !bucket) return []
    const targetMV = received + remaining
    const delta = targetMV - Number(bucket.balance)
    if (Math.abs(delta) > 0.001) {
      const role = delta > 0 ? 'revaluation_gain' : 'revaluation_loss'
      lines.push({ role, bucket: bucketDisplay(sellForm.investBucketId), amount: `${delta > 0 ? '+' : '-'}${form.currency} ${formatAmount(Math.abs(delta))}` })
    }
    lines.push({ role: 'investment_sell', bucket: bucketDisplay(sellForm.investBucketId), amount: `-${form.currency} ${formatAmount(received)}` })
    lines.push({ role: 'cash_leg', bucket: bucketDisplay(sellForm.cashBucketId), amount: `+${form.currency} ${formatAmount(received)}` })
    return lines
  }
  if (isInvestmentIncome.value) {
    const a = Number(form.amount || 0)
    if (!a) return []
    return [{ role: 'investment_income', bucket: bucketDisplay(form.bucketId), amount: `+${form.currency} ${formatAmount(a)}` }]
  }
  if (isFamilyTransfer.value) {
    const a = Number(form.amount || 0)
    if (!a) return []
    return [
      { role: 'family_transfer_out', bucket: bucketDisplay(form.fromBucketId), amount: `-${form.currency} ${formatAmount(a)}` },
      { role: 'family_transfer_in', bucket: familyTransferTargetBuckets.value.find((b) => Number(b.id) === Number(familyTransferForm.toBucketId))?.name || '-', amount: `+${form.currency} ${formatAmount(a)}` },
    ]
  }
  if (isExchange.value) {
    const fromAmount = Number(form.amount || 0)
    const toAmount = Number(form.toAmount || 0)
    if (!fromAmount || !toAmount) return []
    return [
      { role: 'exchange_out', bucket: bucketDisplay(form.fromBucketId), amount: `-${form.currency} ${formatAmount(fromAmount)}` },
      { role: 'exchange_in', bucket: bucketDisplay(form.toBucketId), amount: `+${form.toCurrency} ${formatAmount(toAmount)}` },
    ]
  }
  if (isPaired.value) {
    const amount = Number(form.amount || 0)
    if (!amount) return []
    const cashSign = isCashToCounter.value ? '-' : '+'
    const counterSign = isCashToCounter.value ? '+' : '-'
    return [
      { role: 'cash_leg', bucket: bucketDisplay(form.fromBucketId), amount: `${cashSign}${form.currency} ${formatAmount(amount)}` },
      { role: counterEntryRoleByScenario[form.scenario], bucket: bucketDisplay(form.toBucketId), amount: `${counterSign}${form.currency} ${formatAmount(amount)}` },
    ]
  }
  const amount = Number(form.amount || 0)
  if (!amount) return []
  if (isTransfer.value) {
    return [
      { role: 'transfer_out', bucket: bucketDisplay(form.fromBucketId), amount: `-${form.currency} ${formatAmount(amount)}` },
      { role: 'transfer_in', bucket: bucketDisplay(form.toBucketId), amount: `+${form.currency} ${formatAmount(amount)}` },
    ]
  }
  const sign = form.scenario === 'expense' ? '-' : '+'
  return [{ role: form.scenario, bucket: bucketDisplay(form.bucketId), amount: `${sign}${form.currency} ${formatAmount(amount)}` }]
})

function formatDateTime(date) {
  const pad = (value) => String(value).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`
}

function enumLabel(prefix, value) {
  const key = `${prefix}.${value}`
  return value && te(key) ? t(key) : value
}

function eventTypeLabel(type) {
  return enumLabel('record.scenarios', type)
}

function formatAmount(value) {
  const number = Number(value || 0)
  return new Intl.NumberFormat(config.locale, { minimumFractionDigits: 2, maximumFractionDigits: 2 }).format(number)
}

function bucketLabel(bucket) {
  return `${getBucketEmoji(bucket.bucketType)} ${bucket.name} · ${bucket.currency} ${formatAmount(bucket.balance)}`
}

function bucketName(id) {
  return buckets.value.find((item) => item.id === Number(id))?.name || '-'
}

function bucketDisplay(id) {
  const bucket = buckets.value.find((item) => item.id === Number(id))
  if (!bucket) return '-'
  return `${getBucketEmoji(bucket.bucketType)} ${bucket.name}`
}

function categoryLabel(category) {
  return category.categoryGroupName ? `${category.name} · ${category.categoryGroupName}` : category.name
}
function channelLabel(channel) {
  return `${channel.icon || '🔗'} ${channel.name}`
}

function selectScenario(value) {
  form.scenario = value
  form.categoryId = ''
  form.relatedFinancialEventId = ''
  outstandingItems.value = []
  const picked = scenarios.find((s) => s.value === value)
  if (picked && picked.group !== 'core') moreScenariosOpen.value = true
  loadCategories()
  refreshOutstanding()
}

async function refreshOutstanding() {
  if (!isPaired.value || isCashToCounter.value || !form.toBucketId) {
    outstandingItems.value = []
    return
  }
  const eventType = relatedCreateTypeBySceneFront[form.scenario]
  if (!eventType) {
    outstandingItems.value = []
    return
  }
  const res = await listOutstandingForBucket({ bucketId: form.toBucketId, eventType })
  if (res.success) outstandingItems.value = res.data || []
  else outstandingItems.value = []
}

function selectOutstanding(item) {
  if (!item) {
    form.relatedFinancialEventId = ''
    return
  }
  form.relatedFinancialEventId = item.id
  form.amount = String(item.outstandingAmount)
  if (item.description) form.description = item.description
}

watch(() => form.toBucketId, refreshOutstanding)

function syncFromCurrency() {
  const bucket = buckets.value.find((item) => item.id === Number(form.fromBucketId))
  if (bucket) form.currency = bucket.currency
}

function syncToCurrency() {
  const bucket = buckets.value.find((item) => item.id === Number(form.toBucketId))
  if (bucket) form.toCurrency = bucket.currency
}

async function loadBuckets() {
  const res = await listBuckets({ isActive: true })
  if (res.success) buckets.value = res.data || []
}

async function loadChannels() {
  const res = await listChannels({ isActive: true, eventType: form.scenario })
  if (res.success) channels.value = res.data || []
}

async function loadCategories() {
  if (!needsCategory.value && form.scenario !== 'split') {
    categories.value = []
    form.categoryId = ''
    return
  }
  const type = form.scenario === 'income' ? 'income' : 'expense'
  const res = await listCategories({ type, isActive: true })
  if (res.success) categories.value = res.data || []
  else ElMessage.error(res.error || t('record.messages.loadCategoriesFailed'))
}

async function loadRecentEvents() {
  const res = await listFinancialEvents({ page: 1, pageSize: 5 })
  if (res.success) recentEvents.value = res.data?.list || []
}

async function loadFamiliesForTransfer() {
  const res = await listFamilies()
  if (res.success) {
    families.value = res.data || []
    if (!familyTransferForm.familyId && families.value.length) familyTransferForm.familyId = families.value[0].id
    await loadFamilyAssetsForTransfer()
  }
}

async function loadFamilyAssetsForTransfer() {
  if (!familyTransferForm.familyId) {
    familyAssets.value = null
    return
  }
  const res = await getFamilyAssets(familyTransferForm.familyId)
  if (res.success) familyAssets.value = res.data
}

async function refreshAll() {
  await Promise.all([loadBuckets(), loadCategories(), loadChannels(), loadRecentEvents(), loadFamiliesForTransfer()])
}

function validateForm() {
  if (isSplit.value) {
    return validateSplitForm()
  }
  if (isInvestmentSell.value) {
    return validateInvestmentSellForm()
  }
  if (!Number(form.amount || 0) || Number(form.amount || 0) <= 0) {
    ElMessage.warning(t('record.messages.amountRequired'))
    return false
  }
  if (isInvestmentBuy.value) {
    if (!form.fromBucketId || !form.toBucketId) {
      ElMessage.warning(t('record.messages.selectInvestmentBuckets'))
      return false
    }
    return true
  }
  if (isInvestmentIncome.value) {
    if (!form.bucketId) {
      ElMessage.warning(t('record.messages.selectBucket'))
      return false
    }
    return true
  }
  if (isFamilyTransfer.value) {
    if (!familyTransferForm.familyId) {
      ElMessage.warning(t('record.messages.selectFamily'))
      return false
    }
    if (!familyTransferForm.toUserId) {
      ElMessage.warning(t('record.messages.selectFamilyMember'))
      return false
    }
    if (!form.fromBucketId || !familyTransferForm.toBucketId) {
      ElMessage.warning(t('record.messages.selectFamilyTransferBuckets'))
      return false
    }
    return true
  }
  if (isExchange.value) {
    if (!form.fromBucketId || !form.toBucketId) {
      ElMessage.warning(t('record.messages.selectExchangeBuckets'))
      return false
    }
    if (form.fromBucketId === form.toBucketId) {
      ElMessage.warning(t('record.messages.selectExchangeBuckets'))
      return false
    }
    if (!Number(form.toAmount || 0) || Number(form.toAmount || 0) <= 0) {
      ElMessage.warning(t('record.messages.toAmountRequired'))
      return false
    }
    if (!form.currency || !form.toCurrency || form.currency === form.toCurrency) {
      ElMessage.warning(t('record.messages.sameCurrencyExchange'))
      return false
    }
    return true
  }
  if (isPaired.value) {
    if (!form.fromBucketId || !form.toBucketId) {
      ElMessage.warning(t('record.messages.selectPairedBuckets'))
      return false
    }
    if (form.fromBucketId === form.toBucketId) {
      ElMessage.warning(t('record.messages.selectPairedBuckets'))
      return false
    }
    return true
  }
  if (isTransfer.value) {
    if (!form.fromBucketId || !form.toBucketId) {
      ElMessage.warning(t('record.messages.selectTransferBuckets'))
      return false
    }
    return true
  }
  if (!form.bucketId) {
    ElMessage.warning(t('record.messages.selectBucket'))
    return false
  }
  if (needsCategory.value && !form.categoryId) {
    ElMessage.warning(t('record.messages.selectCategory'))
    return false
  }
  return true
}

function goCreateBucket() {
  router.push('/buckets')
}

function validateSplitForm() {
  const total = Number(splitForm.totalAmount || 0)
  if (!total || total <= 0) {
    ElMessage.warning(t('record.messages.splitTotalRequired'))
    return false
  }
  if (!splitForm.cashBucketId || !splitForm.receivableBucketId) {
    ElMessage.warning(t('record.messages.selectPairedBuckets'))
    return false
  }
  const validShares = splitForm.shares.filter((s) => Number(s.amount || 0) > 0)
  if (!validShares.length) {
    ElMessage.warning(t('record.messages.amountRequired'))
    return false
  }
  if (!validShares.some((s) => !s.isSelf)) {
    ElMessage.warning(t('record.messages.splitNeedFriend'))
    return false
  }
  if (Math.abs(splitBalanceDiff.value) > 0.001) {
    ElMessage.warning(t('record.messages.splitMismatch'))
    return false
  }
  if (validShares.some((s) => s.isSelf) && !form.categoryId) {
    ElMessage.warning(t('record.messages.selectCategory'))
    return false
  }
  return true
}

function validateInvestmentSellForm() {
  if (!sellForm.investBucketId || !sellForm.cashBucketId) {
    ElMessage.warning(t('record.messages.selectInvestmentBuckets'))
    return false
  }
  if (sellForm.investBucketId === sellForm.cashBucketId) {
    ElMessage.warning(t('record.messages.selectInvestmentBuckets'))
    return false
  }
  const received = Number(sellForm.received || 0)
  if (!received || received <= 0) {
    ElMessage.warning(t('record.messages.investmentReceivedRequired'))
    return false
  }
  if (!sellForm.fullSell) {
    const remaining = Number(sellForm.remaining || 0)
    if (remaining < 0) {
      ElMessage.warning(t('record.messages.amountRequired'))
      return false
    }
  }
  return true
}

function addFriendShare() {
  splitForm.shares.push({ isSelf: false, amount: '', description: '' })
}
function removeShare(idx) {
  splitForm.shares.splice(idx, 1)
}
function splitEvenly() {
  const total = Number(splitForm.totalAmount) || 0
  if (!total) {
    ElMessage.warning(t('record.messages.splitTotalRequired'))
    return
  }
  const selfAmount = Number(splitForm.shares.find((s) => s.isSelf)?.amount || 0)
  const remaining = total - selfAmount
  const friendShares = splitForm.shares.filter((s) => !s.isSelf)
  if (!friendShares.length) return
  const each = (remaining / friendShares.length).toFixed(2)
  friendShares.forEach((s) => {
    s.amount = each
  })
}
function splitEvenlyAll() {
  const total = Number(splitForm.totalAmount) || 0
  if (!total) {
    ElMessage.warning(t('record.messages.splitTotalRequired'))
    return
  }
  const all = splitForm.shares
  if (!all.length) return
  const each = (total / all.length).toFixed(2)
  all.forEach((s) => {
    s.amount = each
  })
}
function toggleSelfShare() {
  if (hasSelfShare.value) {
    splitForm.shares = splitForm.shares.filter((s) => !s.isSelf)
  } else {
    splitForm.shares.unshift({ isSelf: true, amount: '', description: '' })
  }
}

async function submitRecord() {
  if (!validateForm()) return
  submitting.value = true
  try {
    let payload
    if (isSplit.value) {
      payload = {
        scenario: 'split',
        cashBucketId: Number(splitForm.cashBucketId),
        receivableBucketId: Number(splitForm.receivableBucketId),
        currency: form.currency,
        amount: String(splitForm.totalAmount),
        categoryId: Number(form.categoryId || 0),
        channelId: Number(form.channelId || 0),
        description: form.description,
        eventTime: form.eventTime,
        remark: form.remark,
        shares: splitForm.shares
          .filter((s) => Number(s.amount || 0) > 0)
          .map((s) => ({
            amount: String(s.amount),
            description: s.isSelf ? form.description : (s.description || form.description),
            isSelf: s.isSelf,
          })),
      }
    } else if (isInvestmentSell.value) {
      payload = {
        scenario: 'investment_sell',
        fromBucketId: Number(sellForm.investBucketId),
        toBucketId: Number(sellForm.cashBucketId),
        amount: String(sellForm.received),
        remainingMarketValue: sellForm.fullSell ? '0' : String(sellForm.remaining || '0'),
        currency: form.currency,
        description: form.description,
        eventTime: form.eventTime,
        remark: form.remark,
      }
    } else if (isFamilyTransfer.value) {
      payload = {
        scenario: 'family_transfer',
        familyId: Number(familyTransferForm.familyId),
        fromBucketId: Number(form.fromBucketId),
        toBucketId: Number(familyTransferForm.toBucketId),
        amount: String(form.amount),
        currency: form.currency,
        description: form.description,
        eventTime: form.eventTime,
        remark: form.remark,
      }
    } else {
      payload = {
        ...form,
        bucketId: Number(form.bucketId || 0),
        fromBucketId: Number(form.fromBucketId || 0),
        toBucketId: Number(form.toBucketId || 0),
        categoryId: Number(form.categoryId || 0),
        channelId: Number(form.channelId || 0),
        relatedFinancialEventId: Number(form.relatedFinancialEventId || 0),
        amount: String(form.amount),
        toAmount: form.toAmount === '' ? '0' : String(form.toAmount),
      }
    }
    const res = await createRecord(payload)
    if (res.success) {
      ElMessage.success(t('record.messages.created'))
      form.amount = ''
      form.toAmount = ''
      form.description = ''
      form.remark = ''
      form.relatedFinancialEventId = ''
      form.eventTime = formatDateTime(new Date())
      if (isSplit.value) {
        splitForm.totalAmount = ''
        splitForm.shares.forEach((s) => { s.amount = ''; s.description = '' })
      }
      if (isInvestmentSell.value) {
        sellForm.received = ''
        sellForm.remaining = ''
      }
      await refreshAll()
    } else {
      ElMessage.error(res.error || t('record.messages.createFailed'))
    }
  } finally {
    submitting.value = false
  }
}

async function applyRouteQuery(query) {
  if (!query || !Object.keys(query).length) return false
  let touched = false
  if (query.scenario && scenarios.find((s) => s.value === query.scenario)) {
    form.scenario = String(query.scenario)
    form.categoryId = ''
    form.relatedFinancialEventId = ''
    outstandingItems.value = []
    const picked = scenarios.find((s) => s.value === query.scenario)
    if (picked && picked.group !== 'core') moreScenariosOpen.value = true
    await loadCategories()
    touched = true
  }
  if (isPaired.value && !isCashToCounter.value) {
    if (query.currency) form.currency = String(query.currency)
    if (query.bucketId) form.toBucketId = Number(query.bucketId)
    if (query.amount) form.amount = String(query.amount)
    if (query.description) form.description = String(query.description)
    if (query.relatedId) form.relatedFinancialEventId = Number(query.relatedId)
    await refreshOutstanding()
    touched = true
  }
  return touched
}

onMounted(async () => {
  await refreshAll()
  const touched = await applyRouteQuery(route.query)
  if (touched) router.replace({ path: route.path, query: {} })
})
onActivated(async () => {
  await refreshAll()
  const touched = await applyRouteQuery(route.query)
  if (touched) router.replace({ path: route.path, query: {} })
})
</script>

<style scoped>
.record-page {
  max-width: 1200px;
  margin: 0 auto;
  color: #1e293b;
}

.reveal-block {
  animation: revealUp 500ms cubic-bezier(0.16, 1, 0.3, 1) both;
}

.delay-1 { animation-delay: 100ms; }

.page-hero,
.scenario-panel,
.form-panel,
.preview-panel {
  background: #ffffff;
  border-radius: 16px;
  box-shadow: 0 1px 3px rgba(15, 23, 42, 0.1), 0 12px 30px rgba(15, 23, 42, 0.04);
}

.page-hero {
  display: flex;
  justify-content: space-between;
  gap: 24px;
  align-items: center;
  margin-bottom: 18px;
  padding: 26px;
  background: linear-gradient(135deg, #fffaf0 0%, #ffffff 70%);
}

.page-hero img {
  border-radius: 24px;
  box-shadow: 0 14px 32px rgba(47, 125, 92, 0.16);
}

.eyebrow {
  margin: 0 0 8px;
  color: #2f7d5c;
  font-size: 12px;
  font-weight: 800;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.page-hero h1,
.section-head h2 {
  margin: 0;
  letter-spacing: -0.022em;
}

.page-hero h1 {
  max-width: 720px;
  font-size: 30px;
  line-height: 1.16;
}

.page-hero p:last-child {
  max-width: 680px;
  margin: 12px 0 0;
  color: #64748b;
  line-height: 1.7;
}

.record-layout {
  display: grid;
  grid-template-columns: 220px minmax(0, 1fr) 330px;
  gap: 18px;
  align-items: start;
}

.scenario-panel,
.form-panel,
.preview-panel {
  padding: 18px;
}

.form-panel,
.preview-panel {
  position: sticky;
  top: 18px;
  align-self: start;
  max-height: calc(100vh - 36px);
  overflow-y: auto;
}

.scenario-panel {
  display: grid;
  gap: 8px;
  align-content: start;
}

.scenario-card {
  min-height: 50px;
  border: 0;
  border-radius: 12px;
  display: grid;
  grid-template-columns: 28px minmax(0, 1fr);
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: #f8faf7;
  color: #1e293b;
  text-align: left;
  cursor: pointer;
  transition-property: transform, box-shadow, background-color, color;
  transition-duration: 160ms;
}

.scenario-card:active,
.primary-action:active {
  transform: scale(0.96);
}

.scenario-card span {
  display: grid;
  place-items: center;
  width: 28px;
  height: 28px;
  border-radius: 10px;
  background: rgba(31, 41, 51, 0.08);
  font-size: 13px;
  font-weight: 900;
}

.scenario-card strong,
.scenario-card small {
  display: block;
  font-size: 13px;
}

.scenario-card small {
  grid-column: 2;
  color: #94a3b8;
  font-size: 11px;
}

.more-scenarios-toggle {
  display: inline-flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  width: 100%;
  margin-top: 6px;
  padding: 9px 14px;
  border: 0;
  border-radius: 999px;
  background: #ffffff;
  color: #2f7d5c;
  font-size: 12px;
  font-weight: 900;
  letter-spacing: 0.02em;
  cursor: pointer;
  box-shadow: inset 0 0 0 1px rgba(47, 125, 92, 0.20);
  transition: transform 160ms, box-shadow 160ms, background-color 160ms;
}

.more-scenarios-toggle:hover {
  background: #f8faf7;
  box-shadow: inset 0 0 0 1px rgba(47, 125, 92, 0.36);
}

.more-scenarios-toggle:active {
  transform: scale(0.98);
}

.more-scenarios-icon {
  font-size: 10px;
  font-weight: 900;
}

.scenario-card.active.income {
  background: rgba(239, 68, 68, 0.12);
  box-shadow: 0 0 0 2px rgba(239, 68, 68, 0.18);
}

.scenario-card.active.expense {
  background: rgba(249, 115, 22, 0.12);
  box-shadow: 0 0 0 2px rgba(249, 115, 22, 0.18);
}

.scenario-card.active.transfer,
.scenario-card.active.refund,
.scenario-card.active.exchange {
  background: rgba(59, 130, 246, 0.12);
  box-shadow: 0 0 0 2px rgba(47, 125, 92, 0.18);
}

.scenario-card.active.paired {
  background: rgba(47, 125, 92, 0.12);
  box-shadow: 0 0 0 2px rgba(47, 125, 92, 0.20);
}

.scenario-card.active.investment {
  background: rgba(124, 58, 237, 0.12);
  box-shadow: 0 0 0 2px rgba(124, 58, 237, 0.22);
}

.scenario-group-label {
  margin: 8px 4px 4px;
  color: #94a3b8;
  font-size: 11px;
  font-weight: 900;
  letter-spacing: 0.06em;
  text-transform: uppercase;
}

.scenario-group-label:first-child {
  margin-top: 0;
}

.counterparty-row {
  display: flex;
  gap: 8px;
  width: 100%;
}

.counterparty-row .el-select {
  flex: 1;
}

.counterparty-empty {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  padding: 14px 16px;
  border: 0;
  border-radius: 14px;
  background: linear-gradient(135deg, #fffaf0 0%, #ffffff 72%);
  box-shadow: inset 0 0 0 1px rgba(47, 125, 92, 0.18);
  color: #1e293b;
  text-align: left;
  cursor: pointer;
  transition-property: transform, box-shadow, background-color;
  transition-duration: 160ms;
}

.counterparty-empty:hover {
  background: #f8faf7;
  box-shadow: inset 0 0 0 1px rgba(47, 125, 92, 0.32), 0 8px 18px rgba(47, 125, 92, 0.10);
}

.counterparty-empty:active {
  transform: scale(0.99);
}

.counterparty-empty-icon {
  display: grid;
  place-items: center;
  width: 36px;
  height: 36px;
  flex-shrink: 0;
  border-radius: 12px;
  background: rgba(47, 125, 92, 0.12);
  color: #2f7d5c;
  font-size: 22px;
  font-weight: 900;
  line-height: 1;
}

.counterparty-empty-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.counterparty-empty-text strong {
  color: #1e293b;
  font-size: 13px;
  font-weight: 800;
}

.counterparty-empty-text small {
  color: #2f7d5c;
  font-size: 12px;
  font-weight: 800;
}

.split-shares {
  margin: 12px 0 18px;
  padding: 16px;
  border-radius: 14px;
  background: #fffaf0;
  box-shadow: inset 0 0 0 1px rgba(100, 116, 139, 0.12);
}

.split-shares-label {
  margin: 0 0 12px;
  color: #64748b;
  font-size: 12px;
  font-weight: 900;
  letter-spacing: 0.04em;
}

.split-share-row {
  display: grid;
  grid-template-columns: minmax(0, 1.4fr) 130px;
  gap: 10px;
  align-items: center;
  margin-bottom: 10px;
  padding: 12px;
  border-radius: 12px;
  background: #ffffff;
  box-shadow: inset 0 0 0 1px rgba(100, 116, 139, 0.10);
}

.split-share-row.self {
  background: rgba(249, 115, 22, 0.06);
  box-shadow: inset 0 0 0 1px rgba(249, 115, 22, 0.20);
  grid-template-columns: minmax(0, 1fr) 130px;
}

.split-share-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  grid-column: 1 / -1;
  margin-bottom: -2px;
}

.split-share-head strong {
  color: #1e293b;
  font-size: 13px;
  font-weight: 800;
}

.split-share-remove {
  display: grid;
  place-items: center;
  width: 24px;
  height: 24px;
  border: 0;
  border-radius: 999px;
  background: rgba(100, 116, 139, 0.12);
  color: #64748b;
  font-size: 16px;
  font-weight: 900;
  cursor: pointer;
}

.split-share-remove:hover {
  background: rgba(239, 68, 68, 0.14);
  color: #ef4444;
}

.split-share-row.self .split-share-amount {
  grid-column: 2;
}

.split-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 14px;
}

.split-chip {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  min-height: 34px;
  padding: 0 14px;
  border: 0;
  border-radius: 999px;
  background: #ffffff;
  color: #64748b;
  font-size: 12px;
  font-weight: 900;
  letter-spacing: 0.01em;
  cursor: pointer;
  box-shadow: inset 0 0 0 1px rgba(100, 116, 139, 0.16);
  transition-property: transform, box-shadow, background-color, color;
  transition-duration: 160ms;
}

.split-chip:hover:not(:disabled) {
  background: #f8faf7;
  color: #2f7d5c;
  box-shadow: inset 0 0 0 1px rgba(47, 125, 92, 0.28);
}

.split-chip:active:not(:disabled) {
  transform: scale(0.96);
}

.split-chip:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.split-chip.primary {
  background: rgba(47, 125, 92, 0.10);
  color: #2f7d5c;
  box-shadow: inset 0 0 0 1px rgba(47, 125, 92, 0.32);
}

.split-chip.primary:hover:not(:disabled) {
  background: rgba(47, 125, 92, 0.16);
  box-shadow: inset 0 0 0 1px rgba(47, 125, 92, 0.40);
}

.split-chip.subtle {
  background: rgba(31, 41, 51, 0.06);
  color: #475569;
  box-shadow: none;
}

.split-chip.subtle:hover:not(:disabled) {
  background: rgba(31, 41, 51, 0.10);
  color: #1e293b;
  box-shadow: none;
}

.split-chip-icon {
  display: grid;
  place-items: center;
  width: 18px;
  height: 18px;
  border-radius: 999px;
  background: rgba(47, 125, 92, 0.14);
  color: #2f7d5c;
  font-size: 12px;
  font-weight: 900;
  line-height: 1;
}

.split-chip:hover:not(:disabled) .split-chip-icon {
  background: rgba(47, 125, 92, 0.22);
}

.split-chip.subtle .split-chip-icon {
  background: rgba(31, 41, 51, 0.10);
  color: #475569;
}

.split-summary {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-top: 14px;
  padding: 12px 14px;
  border-radius: 12px;
  font-family: 'SF Mono', 'Fira Code', monospace;
  font-size: 13px;
  font-weight: 800;
}

.split-summary.balanced {
  background: rgba(47, 125, 92, 0.10);
  color: #2f7d5c;
}

.split-summary.off {
  background: rgba(239, 68, 68, 0.08);
  color: #ef4444;
}

.split-summary strong {
  font-weight: 900;
}

.split-diff {
  font-weight: 900;
}

.split-balanced {
  color: #2f7d5c;
  font-weight: 900;
}

.scenario-card.active.split {
  background: rgba(249, 115, 22, 0.12);
  box-shadow: 0 0 0 2px rgba(249, 115, 22, 0.20);
}

.event-preview.split {
  background: rgba(249, 115, 22, 0.10);
}

.event-preview.investment_buy,
.event-preview.investment_sell,
.event-preview.investment_income,
.event-preview.investment_revalue {
  background: rgba(124, 58, 237, 0.10);
}

.sell-mode-switch {
  display: inline-flex;
  gap: 4px;
  margin-bottom: 14px;
  padding: 4px;
  border-radius: 999px;
  background: rgba(31, 41, 51, 0.06);
}

.sell-mode-switch button {
  min-height: 30px;
  border: 0;
  border-radius: 999px;
  padding: 0 14px;
  background: transparent;
  color: #64748b;
  font-size: 12px;
  font-weight: 900;
  cursor: pointer;
  transition: transform 160ms, background-color 160ms, color 160ms;
}

.sell-mode-switch button:active { transform: scale(0.96); }

.sell-mode-switch button.active {
  background: rgba(124, 58, 237, 0.14);
  color: #7c3aed;
}

.exchange-rate-hint.gain {
  background: rgba(47, 125, 92, 0.10);
  color: #2f7d5c;
}

.exchange-rate-hint.loss {
  background: rgba(239, 68, 68, 0.10);
  color: #ef4444;
}

.exchange-rate-hint.neutral {
  background: rgba(31, 41, 51, 0.06);
  color: #64748b;
}

.outstanding-list {
  display: grid;
  gap: 8px;
  margin: 4px 0 18px;
  padding: 14px;
  border-radius: 14px;
  background: #fffaf0;
  box-shadow: inset 0 0 0 1px rgba(100, 116, 139, 0.12);
}

.outstanding-label {
  margin: 0 0 4px;
  color: #64748b;
  font-size: 12px;
  font-weight: 900;
  letter-spacing: 0.04em;
}

.outstanding-row {
  display: grid;
  grid-template-columns: 22px minmax(0, 1fr) auto;
  gap: 12px;
  align-items: center;
  padding: 12px 14px;
  border: 0;
  border-radius: 12px;
  background: #ffffff;
  color: #1e293b;
  text-align: left;
  cursor: pointer;
  box-shadow: inset 0 0 0 1px rgba(100, 116, 139, 0.10);
  transition-property: transform, box-shadow, background-color;
  transition-duration: 160ms;
}

.outstanding-row:hover {
  background: #f8faf7;
  box-shadow: inset 0 0 0 1px rgba(47, 125, 92, 0.20);
}

.outstanding-row.active {
  background: rgba(47, 125, 92, 0.08);
  box-shadow: inset 0 0 0 2px rgba(47, 125, 92, 0.32);
}

.outstanding-row:active {
  transform: scale(0.99);
}

.outstanding-mark {
  font-size: 16px;
  font-weight: 900;
  color: #2f7d5c;
  text-align: center;
}

.outstanding-main {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.outstanding-main strong {
  color: #1e293b;
  font-size: 13px;
  font-weight: 800;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.outstanding-main small {
  color: #64748b;
  font-size: 11px;
  font-weight: 700;
}

.outstanding-original {
  color: #94a3b8;
  font-family: 'SF Mono', 'Fira Code', monospace;
  font-size: 12px;
  font-weight: 800;
  font-variant-numeric: tabular-nums;
}

.scenario-card.disabled {
  cursor: not-allowed;
  opacity: 0.48;
}

.section-head {
  display: flex;
  justify-content: space-between;
  margin-bottom: 18px;
}

.full-width {
  width: 100%;
}

.primary-action {
  min-height: 42px;
  border: 0;
  border-radius: 12px;
  padding: 0 18px;
  background: #2f7d5c;
  color: #ffffff;
  font-weight: 800;
  cursor: pointer;
  box-shadow: 0 10px 24px rgba(47, 125, 92, 0.20);
  transition-property: transform, box-shadow, background-color, color;
  transition-duration: 160ms;
}

.event-preview {
  padding: 16px;
  border-radius: 14px;
  background: #f8faf7;
}

.event-preview.income {
  background: rgba(239, 68, 68, 0.1);
}

.event-preview.expense {
  background: rgba(249, 115, 22, 0.1);
}

.event-preview.exchange,
.event-preview.transfer {
  background: rgba(59, 130, 246, 0.1);
}

.event-preview.receivable_create,
.event-preview.receivable_collect,
.event-preview.deposit_create,
.event-preview.deposit_refund,
.event-preview.loan_out,
.event-preview.loan_collect {
  background: rgba(47, 125, 92, 0.1);
}

.exchange-rate-hint {
  margin: 0 0 16px;
  padding: 10px 14px;
  border-radius: 12px;
  background: rgba(59, 130, 246, 0.08);
  color: #1d4ed8;
  font-family: 'SF Mono', 'Fira Code', monospace;
  font-weight: 800;
  font-size: 13px;
}

.event-preview span,
.event-preview p {
  color: #64748b;
  font-weight: 700;
}

.event-preview strong {
  display: block;
  margin: 8px 0;
  font-family: 'SF Mono', 'Fira Code', monospace;
  font-size: 24px;
}

.entry-preview,
.recent-events {
  margin-top: 18px;
}

.entry-preview h3,
.recent-events h3 {
  margin: 0 0 10px;
}

.entry-list,
.recent-events {
  display: grid;
  gap: 10px;
}

.entry-row,
.recent-row {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  padding: 12px;
  border-radius: 12px;
  background: #ffffff;
  box-shadow: inset 0 0 0 1px rgba(100, 116, 139, 0.12);
}

.entry-row strong,
.recent-row strong {
  font-family: 'SF Mono', 'Fira Code', monospace;
}

.negative {
  color: #f97316;
}

.positive {
  color: #ef4444;
}

@media (hover: hover) {
  .scenario-card:not(.disabled):hover {
    transform: translateY(-2px);
    box-shadow: 0 10px 22px rgba(15, 23, 42, 0.08);
  }
}

@media (max-width: 1080px) {
  .record-layout {
    grid-template-columns: 1fr;
  }

  .scenario-panel {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 600px) {
  .page-hero {
    align-items: flex-start;
    flex-direction: column;
  }

  .page-hero h1 {
    font-size: 24px;
  }

  .scenario-panel {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .scenario-group-label,
  .more-scenarios-toggle {
    grid-column: 1 / -1;
  }

  .scenario-card {
    min-height: 44px;
    grid-template-columns: 24px minmax(0, 1fr);
    padding: 8px 10px;
  }

  .scenario-card span {
    width: 24px;
    height: 24px;
  }

  .form-panel,
  .preview-panel {
    position: static;
    max-height: none;
    overflow: visible;
  }

  .split-share-row,
  .investment-bucket-row {
    grid-template-columns: 1fr;
  }
}

@media (prefers-reduced-motion: reduce) {
  .reveal-block,
  .scenario-card,
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
