<template>
  <main class="family-page">
    <section class="page-hero reveal-block">
      <div><p class="eyebrow">{{ t('family.hero.eyebrow') }}</p><h1>{{ t('family.hero.title') }}</h1><p>{{ t('family.hero.subtitle') }}</p></div>
      <button class="primary-action" @click="openCreate">{{ t('family.actions.new') }}</button>
    </section>
    <section class="family-layout reveal-block delay-1">
      <aside class="family-list">
        <h2>{{ t('family.sections.myFamilies') }}</h2>
        <button v-for="item in families" :key="item.id" class="family-item" :class="{active:selectedFamily?.id===item.id}" @click="selectFamily(item)"><strong>{{ item.name }}</strong><span>{{ item.role }}</span></button>
        <p v-if="!families.length" class="muted">{{ t('family.empty.noFamilies') }}</p>
        <h2>{{ t('family.sections.invitations') }}</h2>
        <div v-for="item in invitations" :key="item.id" class="invite-card"><strong>{{ item.familyName }}</strong><span>{{ item.displayName || item.account }}</span><div><button class="primary-action small" @click="acceptInvite(item)">{{ t('family.actions.accept') }}</button><button class="ghost-action small" @click="rejectInvite(item)">{{ t('family.actions.reject') }}</button></div></div>
        <p v-if="!invitations.length" class="muted">{{ t('family.empty.noInvitations') }}</p>
      </aside>
      <section class="family-detail">
        <template v-if="selectedFamily">
          <div class="section-head"><div><p class="eyebrow">{{ getCurrencyLabel(selectedFamily.baseCurrency, config.locale) }}</p><h2>{{ selectedFamily.name }}</h2></div></div>
          <div class="invite-form"><el-input v-model="inviteForm.account" :placeholder="t('family.fields.account')"/><el-input v-model="inviteForm.displayName" :placeholder="t('family.fields.displayName')"/><button class="primary-action" @click="inviteMember">{{ t('family.actions.invite') }}</button></div>
          <h3>{{ t('family.sections.members') }}</h3>
          <div class="member-grid"><div v-for="member in members" :key="member.id" class="member-card"><strong>{{ member.displayName || member.name || member.account }}</strong><span>{{ member.role }} · {{ member.status }}</span></div></div>
          <h3>{{ t('family.sections.statistics') }}</h3>
          <div class="stats-grid"><div class="stat-tile"><span>{{ t('statistics.income') }}</span><strong>{{ formatAmount(summary?.income) }}</strong></div><div class="stat-tile"><span>{{ t('statistics.expense') }}</span><strong>{{ formatAmount(summary?.netExpense || summary?.expense) }}</strong></div><div class="stat-tile"><span>{{ t('statistics.refund') }}</span><strong>{{ formatAmount(summary?.refund) }}</strong></div><div class="stat-tile"><span>{{ t('statistics.net') }}</span><strong>{{ formatAmount(summary?.net) }}</strong></div></div>
        </template>
        <div v-else class="empty-state"><p>{{ t('family.empty.noFamilies') }}</p></div>
      </section>
    </section>
    <el-dialog v-model="dialogVisible" :title="t('family.actions.new')" width="520px" class="marmot-dialog">
      <el-form label-position="top">
        <el-form-item :label="t('family.fields.name')">
          <el-input v-model="familyForm.name" />
        </el-form-item>
        <el-form-item :label="t('family.fields.baseCurrency')">
          <el-select v-model="familyForm.baseCurrency" class="full-width">
            <el-option v-for="currency in currencyOptions" :key="currency.code" :label="getCurrencyLabel(currency.code, config.locale)" :value="currency.code">
              <span class="currency-option"><span>{{ currency.icon }}</span><strong>{{ currency.code }}</strong><small>{{ config.locale === 'en-US' ? currency.englishName : currency.name }}</small></span>
            </el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <button class="ghost-action" @click="dialogVisible=false">{{ t('common.actions.cancel') }}</button>
        <button class="primary-action" @click="createNewFamily">{{ t('common.actions.create') }}</button>
      </template>
    </el-dialog>
  </main>
</template>
<script setup>
import { onActivated, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { useConfigStore } from '@/stores/config'
import { acceptFamilyInvitation, createFamily, getFamilyStatisticsSummary, inviteFamilyMember, listFamilies, listFamilyInvitations, listFamilyMembers, rejectFamilyInvitation } from '@/api/family/family'
import { currencyOptions, getCurrencyLabel } from '@/utils/currency'
const { t } = useI18n(); const config = useConfigStore(); const families=ref([]); const invitations=ref([]); const selectedFamily=ref(null); const members=ref([]); const summary=ref(null); const dialogVisible=ref(false); const familyForm=reactive({name:'',baseCurrency:'CNY'}); const inviteForm=reactive({account:'',displayName:''})
function formatAmount(v){return new Intl.NumberFormat(config.locale,{minimumFractionDigits:2,maximumFractionDigits:2}).format(Number(v||0))}
function openCreate(){familyForm.name='';familyForm.baseCurrency='CNY';dialogVisible.value=true}
async function loadFamilies(){const res=await listFamilies(); if(res.success){families.value=res.data||[]; if(!selectedFamily.value&&families.value.length) await selectFamily(families.value[0])}}
async function loadInvitations(){const res=await listFamilyInvitations(); if(res.success) invitations.value=res.data||[]}
async function selectFamily(item){selectedFamily.value=item; const [m,s]=await Promise.all([listFamilyMembers(item.id,{includeInvited:true}),getFamilyStatisticsSummary(item.id,{currency:item.baseCurrency||'CNY'})]); if(m.success)members.value=m.data||[]; if(s.success)summary.value=s.data}
async function createNewFamily(){const res=await createFamily(familyForm); if(res.success){ElMessage.success(t('family.messages.created')); dialogVisible.value=false; selectedFamily.value=res.data; await loadFamilies()} else ElMessage.error(res.error||t('family.messages.actionFailed'))}
async function inviteMember(){if(!selectedFamily.value)return; const res=await inviteFamilyMember(selectedFamily.value.id, inviteForm); if(res.success){ElMessage.success(t('family.messages.invited')); inviteForm.account=''; inviteForm.displayName=''; await selectFamily(selectedFamily.value)} else ElMessage.error(res.error||t('family.messages.actionFailed'))}
async function acceptInvite(item){const res=await acceptFamilyInvitation(item.id); if(res.success){ElMessage.success(t('family.messages.accepted')); await refreshAll()} else ElMessage.error(res.error||t('family.messages.actionFailed'))}
async function rejectInvite(item){const res=await rejectFamilyInvitation(item.id); if(res.success){ElMessage.success(t('family.messages.rejected')); await refreshAll()} else ElMessage.error(res.error||t('family.messages.actionFailed'))}
async function refreshAll(){await Promise.all([loadFamilies(),loadInvitations()])}
onMounted(refreshAll); onActivated(refreshAll)
</script>
<style scoped>.family-page{max-width:1200px;margin:0 auto;color:#1e293b}.reveal-block{animation:revealUp 480ms cubic-bezier(.16,1,.3,1) both}.delay-1{animation-delay:90ms}.page-hero,.family-list,.family-detail{background:#fff;border-radius:16px;box-shadow:0 1px 3px rgba(15,23,42,.1),0 12px 30px rgba(15,23,42,.04)}.page-hero{display:flex;justify-content:space-between;gap:24px;margin-bottom:18px;padding:26px;background:linear-gradient(135deg,#fffaf0 0%,#fff 70%)}.eyebrow{margin:0 0 8px;color:#2f7d5c;font-size:12px;font-weight:800;letter-spacing:.08em;text-transform:uppercase}.page-hero h1{margin:0;font-size:30px;line-height:1.16;letter-spacing:-.022em}.page-hero p:last-child{margin:12px 0 0;color:#64748b;line-height:1.7}.family-layout{display:grid;grid-template-columns:320px minmax(0,1fr);gap:18px}.family-list,.family-detail{padding:22px}.family-list h2,.family-detail h2{margin:0 0 14px}.family-item,.invite-card,.member-card,.stat-tile{display:block;width:100%;border:0;border-radius:14px;padding:14px;background:#f8faf7;margin-bottom:10px;text-align:left}.family-item{cursor:pointer}.family-item.active{box-shadow:0 0 0 2px rgba(47,125,92,.18)}.family-item strong,.family-item span,.member-card strong,.member-card span,.invite-card strong,.invite-card span{display:block}.family-item span,.member-card span,.invite-card span,.muted{color:#64748b;font-size:13px}.primary-action,.ghost-action{min-height:40px;border:0;border-radius:12px;padding:0 16px;font-weight:800;cursor:pointer}.primary-action{background:#3b82f6;color:#fff;box-shadow:0 10px 24px rgba(59,130,246,.22)}.ghost-action{background:#f8faf7;color:#1e293b}.small{min-height:32px;margin-right:8px}.invite-form{display:grid;grid-template-columns:1fr 1fr auto;gap:10px;margin-bottom:20px}.member-grid,.stats-grid{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:10px}.stat-tile span{color:#64748b;font-size:12px;font-weight:800}.stat-tile strong{display:block;margin-top:8px;font-family:'SF Mono','Fira Code',monospace;font-size:24px}.full-width{width:100%}.currency-option{display:flex;align-items:center;gap:8px}.currency-option strong{font-family:'SF Mono','Fira Code',monospace}.currency-option small{color:#64748b}@media(max-width:900px){.family-layout{grid-template-columns:1fr}.invite-form{grid-template-columns:1fr}.page-hero{flex-direction:column}}@keyframes revealUp{from{opacity:0;transform:translateY(12px);filter:blur(4px)}to{opacity:1;transform:translateY(0);filter:blur(0)}}</style>
