<template>
  <div class="management-dialog-footer">
    <button class="management-ghost-action" type="button" @click="$emit('cancel')">{{ cancelLabel || t('common.actions.cancel') }}</button>
    <button class="management-primary-action" type="button" :disabled="loading" @click="$emit('submit')">
      {{ submitLabel || t('common.actions.save') }}
    </button>
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

defineProps({
  submitLabel: { type: String, default: '' },
  cancelLabel: { type: String, default: '' },
  loading: { type: Boolean, default: false },
})

defineEmits(['submit', 'cancel'])
</script>

<style scoped>
.management-dialog-footer {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
  align-items: center;
  flex-wrap: nowrap;
}

.management-dialog-footer > :deep(.management-primary-action),
.management-dialog-footer > :deep(.management-ghost-action) {
  flex: 0 0 auto;
}
</style>
