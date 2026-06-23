<template>
  <el-form-item :label="label || t('common.status.status')">
    <el-switch
      :model-value="modelValue"
      :active-text="t('common.status.enabled')"
      :inactive-text="t('common.status.disabled')"
      active-color="#2f7d5c"
      inactive-color="#cbd5e1"
      @update:model-value="$emit('update:modelValue', $event)"
    />
  </el-form-item>
</template>

<script setup>
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

defineProps({
  modelValue: { type: Boolean, default: true },
  label: { type: String, default: '' },
})

defineEmits(['update:modelValue'])
</script>
