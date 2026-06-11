<template>
  <main class="ledger-page">
    <section class="page-hero reveal-block">
      <div>
        <p class="eyebrow">{{ t('categories.hero.eyebrow') }}</p>
        <h1>{{ t('categories.hero.title') }}</h1>
        <p>{{ t('categories.hero.subtitle') }}</p>
      </div>
      <button class="primary-action" @click="openCreate">{{ t('categories.actions.new') }}</button>
    </section>

    <section class="toolbar reveal-block delay-1">
      <el-select v-model="filters.type" clearable :placeholder="t('categories.fields.type')" class="filter-control" @change="loadCategories">
        <el-option :label="t('domain.expense')" value="expense" />
        <el-option :label="t('domain.income')" value="income" />
      </el-select>
      <el-select v-model="filters.isActive" clearable :placeholder="t('common.filters.enabledStatus')" class="filter-control" @change="loadCategories">
        <el-option :label="t('common.status.enabled')" :value="true" />
        <el-option :label="t('common.status.disabled')" :value="false" />
      </el-select>
      <button class="ghost-action" @click="loadCategories">{{ t('common.actions.refresh') }}</button>
    </section>

    <section v-loading="loading" class="category-grid reveal-block delay-2">
      <article v-for="item in categories" :key="item.id" class="category-card">
        <div class="category-dot" :style="{ background: item.categoryGroupColor || (item.type === 'income' ? '#ef4444' : '#10b981') }"></div>
        <div>
          <h2>{{ item.name }}</h2>
          <p>{{ item.categoryGroupName || item.categoryGroupCode }}</p>
          <div class="category-meta">
            <span :class="item.type">{{ item.type === 'income' ? t('domain.income') : t('domain.expense') }}</span>
            <span>{{ item.isActive ? t('common.status.enabled') : t('common.status.disabled') }}</span>
          </div>
        </div>
        <div class="card-actions">
          <button class="text-action" @click="openEdit(item)">{{ t('common.actions.edit') }}</button>
          <button class="danger-action" @click="handleDelete(item)">{{ t('common.actions.delete') }}</button>
        </div>
      </article>

      <div v-if="!loading && !categories.length" class="empty-state">
        <img :src="marmotOne" :alt="t('categories.empty.alt')" width="112" height="112" />
        <h2>{{ t('categories.empty.title') }}</h2>
        <p>{{ t('categories.empty.text') }}</p>
        <button class="primary-action" @click="openCreate">{{ t('categories.actions.new') }}</button>
      </div>
    </section>

    <el-dialog v-model="dialogVisible" :title="editingId ? t('categories.dialog.editTitle') : t('categories.dialog.createTitle')" width="520px" class="marmot-dialog">
      <el-form ref="formRef" :model="form" label-position="top">
        <el-form-item :label="t('categories.fields.name')">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item :label="t('categories.fields.type')">
          <el-select v-model="form.type" class="full-width" @change="handleTypeChange">
            <el-option :label="t('domain.expense')" value="expense" />
            <el-option :label="t('domain.income')" value="income" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('categories.fields.categoryGroup')">
          <el-select v-model="form.categoryGroupId" class="full-width">
            <el-option v-for="group in filteredGroups" :key="group.id" :label="group.name" :value="group.id" />
          </el-select>
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
import { createCategory, deleteCategory, listCategories, updateCategory } from '@/api/category/category'
import { listCategoryGroups } from '@/api/categoryGroup/categoryGroup'
import marmotOne from '../../../img/marmot-ledger-1.png'

const { t } = useI18n()
const categories = ref([])
const groups = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const editingId = ref(0)
const formRef = ref()
const filters = reactive({ type: '', isActive: '' })
const form = reactive({ name: '', type: 'expense', categoryGroupId: '', isActive: true })
const filteredGroups = computed(() => groups.value.filter((item) => item.type === form.type))

async function loadGroups() {
  const res = await listCategoryGroups({ enabled: true })
  if (res.success) groups.value = res.data || []
  else ElMessage.error(res.error || t('categories.messages.loadGroupsFailed'))
}

async function loadCategories() {
  loading.value = true
  try {
    const params = {}
    if (filters.type) params.type = filters.type
    if (filters.isActive !== '') params.isActive = filters.isActive
    const res = await listCategories(params)
    if (res.success) categories.value = res.data || []
    else ElMessage.error(res.error || t('categories.messages.loadFailed'))
  } finally { loading.value = false }
}

async function refreshAll() { await Promise.all([loadGroups(), loadCategories()]) }
function openCreate() { editingId.value = 0; Object.assign(form, { name: '', type: 'expense', categoryGroupId: '', isActive: true }); dialogVisible.value = true }
function openEdit(item) { editingId.value = item.id; Object.assign(form, { name: item.name, type: item.type, categoryGroupId: item.categoryGroupId, isActive: item.isActive !== false }); dialogVisible.value = true }
function handleTypeChange() { form.categoryGroupId = '' }

async function submitForm() {
  if (!form.name) return ElMessage.warning(t('categories.validation.nameRequired'))
  if (!form.type) return ElMessage.warning(t('categories.validation.typeRequired'))
  if (!form.categoryGroupId) return ElMessage.warning(t('categories.validation.groupRequired'))
  const payload = { ...form, categoryGroupId: Number(form.categoryGroupId), isActive: form.isActive !== false }
  const res = editingId.value ? await updateCategory(editingId.value, payload) : await createCategory(payload)
  if (res.success) { ElMessage.success(editingId.value ? t('categories.messages.updated') : t('categories.messages.created')); dialogVisible.value = false; await loadCategories() }
  else ElMessage.error(res.error || t('categories.messages.saveFailed'))
}

async function handleDelete(item) {
  try {
    await ElMessageBox.confirm(t('categories.delete.confirm', { name: item.name }), t('categories.delete.title'), { confirmButtonText: t('common.actions.delete'), cancelButtonText: t('common.actions.cancel'), type: 'warning', customClass: 'calm-marmot-message-box calm-marmot-delete-box', confirmButtonClass: 'calm-marmot-danger-confirm', cancelButtonClass: 'calm-marmot-soft-cancel' })
    const res = await deleteCategory(item.id)
    if (res.success) { ElMessage.success(t('categories.messages.deleted')); await loadCategories() }
    else ElMessage.error(res.error || t('categories.messages.deleteFailed'))
  } catch (err) { if (err !== 'cancel') console.warn(err) }
}

onMounted(refreshAll)
onActivated(refreshAll)
</script>

<style scoped>
.ledger-page{max-width:1200px;margin:0 auto;color:#1e293b}.reveal-block,.category-card{animation:revealUp 480ms cubic-bezier(.16,1,.3,1) both}.delay-1{animation-delay:90ms}.delay-2{animation-delay:160ms}.page-hero,.toolbar,.category-card,.empty-state{background:#fff;border-radius:16px;box-shadow:0 1px 3px rgba(15,23,42,.1),0 12px 30px rgba(15,23,42,.04)}.page-hero{display:flex;justify-content:space-between;gap:24px;align-items:flex-start;margin-bottom:18px;padding:26px;background:linear-gradient(135deg,#fffaf0 0%,#fff 70%)}.eyebrow{margin:0 0 8px;color:#2f7d5c;font-size:12px;font-weight:800;letter-spacing:.08em;text-transform:uppercase}.page-hero h1{max-width:720px;margin:0;font-size:30px;line-height:1.16;letter-spacing:-.022em}.page-hero p:last-child{max-width:680px;margin:12px 0 0;color:#64748b;line-height:1.7}.toolbar{display:flex;gap:12px;align-items:center;margin-bottom:18px;padding:14px}.filter-control{width:180px}.category-grid{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:16px}.category-card{display:grid;grid-template-columns:10px minmax(0,1fr) auto;gap:18px;overflow:hidden;transition-property:transform,box-shadow;transition-duration:180ms}.category-dot{min-height:100%}.category-card>div:nth-child(2){padding:20px 0}.category-card h2{margin:0;font-size:20px;letter-spacing:-.012em}.category-card p{margin:6px 0 0;color:#64748b}.category-meta{display:flex;gap:8px;margin-top:14px}.category-meta span{display:inline-flex;align-items:center;min-height:28px;border-radius:999px;padding:0 10px;background:rgba(31,41,51,.08);font-size:12px;font-weight:800}.category-meta .income{color:#ef4444;background:rgba(239,68,68,.1)}.category-meta .expense{color:#10b981;background:rgba(16,185,129,.1)}.card-actions{display:flex;gap:8px;align-items:center;padding-right:20px}.primary-action,.ghost-action,.text-action,.danger-action{min-height:40px;border:0;border-radius:12px;padding:0 16px;font-weight:800;cursor:pointer;transition-property:transform,box-shadow,background-color,color;transition-duration:160ms}.primary-action:active,.ghost-action:active,.text-action:active,.danger-action:active{transform:scale(.96)}.primary-action{background:#3b82f6;color:#fff;box-shadow:0 10px 24px rgba(59,130,246,.22)}.ghost-action,.text-action{background:#f8faf7;color:#1e293b}.danger-action{background:rgba(239,68,68,.1);color:#ef4444}.empty-state{grid-column:1/-1;display:grid;place-items:center;gap:12px;padding:40px 24px;text-align:center;color:#64748b}.empty-state img{border-radius:22px}.empty-state h2{margin:0;color:#1e293b}.full-width{width:100%}@media(hover:hover){.category-card:hover{transform:translateY(-2px);box-shadow:0 3px 8px rgba(15,23,42,.12),0 16px 34px rgba(15,23,42,.06)}}@media(max-width:820px){.page-hero,.toolbar{flex-direction:column;align-items:stretch}.category-grid{grid-template-columns:1fr}.filter-control{width:100%}.category-card{grid-template-columns:10px minmax(0,1fr)}}@media(max-width:520px){.card-actions{grid-column:2;justify-content:flex-start;padding:0 20px 18px}.page-hero h1{font-size:24px}}@keyframes revealUp{from{opacity:0;transform:translateY(12px);filter:blur(4px)}to{opacity:1;transform:translateY(0);filter:blur(0)}}
</style>
