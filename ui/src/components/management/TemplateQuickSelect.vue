<template>
  <section class="template-quick-select">
    <header class="template-quick-head" :class="{ clickable: collapsible }" @click="toggle">
      <span class="template-quick-icon">⚡</span>
      <span class="template-quick-title">{{ title || t('common.templates.title') }}</span>
      <span v-if="collapsible" class="template-quick-toggle">{{ collapsed ? t('common.templates.expand') : t('common.templates.collapse') }}</span>
    </header>
    <p v-if="hint && !collapsed" class="template-quick-hint">{{ hint }}</p>
    <div v-if="!collapsed" class="template-quick-body">
      <p v-if="!items.length" class="template-quick-empty">{{ emptyText || t('common.templates.empty') }}</p>
      <div v-else class="template-quick-grid">
        <button
          v-for="item in items"
          :key="getKey(item)"
          type="button"
          class="template-quick-chip"
          :class="{ active: getKey(item) === activeKey }"
          @click="$emit('select', item)"
        >
          <slot name="chip" :item="item">
            <span class="template-quick-chip-icon">{{ item.icon || '📁' }}</span>
            <strong>{{ item.name }}</strong>
          </slot>
        </button>
      </div>
    </div>
  </section>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  title: { type: String, default: '' },
  hint: { type: String, default: '' },
  items: { type: Array, default: () => [] },
  itemKey: { type: [String, Function], default: 'id' },
  activeKey: { type: [Number, String, null], default: null },
  collapsible: { type: Boolean, default: true },
  defaultCollapsed: { type: Boolean, default: false },
  emptyText: { type: String, default: '' },
})

defineEmits(['select'])

const collapsed = ref(props.defaultCollapsed)

watch(() => props.defaultCollapsed, (value) => {
  collapsed.value = value
})

function toggle() {
  if (!props.collapsible) return
  collapsed.value = !collapsed.value
}

function getKey(item) {
  if (typeof props.itemKey === 'function') return props.itemKey(item)
  return item?.[props.itemKey]
}
</script>

<style scoped>
.template-quick-select {
  margin-bottom: 16px;
  padding: 14px;
  border-radius: 14px;
  background: #fffaf0;
  box-shadow: inset 0 0 0 1px rgba(120, 92, 56, 0.10);
}

.template-quick-head {
  display: flex;
  align-items: center;
  gap: 8px;
  user-select: none;
}

.template-quick-head.clickable {
  cursor: pointer;
}

.template-quick-icon {
  font-size: 16px;
}

.template-quick-title {
  font-size: 13px;
  font-weight: 800;
  color: #1e293b;
  letter-spacing: -.012em;
}

.template-quick-toggle {
  margin-left: auto;
  color: #857462;
  font-size: 12px;
  font-weight: 700;
}

.template-quick-hint {
  margin: 8px 0 0;
  color: #857462;
  font-size: 12px;
  line-height: 1.6;
}

.template-quick-body {
  margin-top: 10px;
}

.template-quick-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.template-quick-empty {
  margin: 0;
  color: #857462;
  font-size: 13px;
}

.template-quick-chip {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  min-height: 36px;
  border: 0;
  border-radius: 10px;
  padding: 0 12px;
  background: #ffffff;
  color: #4b3f33;
  font-size: 13px;
  font-weight: 700;
  cursor: pointer;
  box-shadow: inset 0 0 0 1px rgba(100, 116, 139, 0.16);
  transition-property: transform, box-shadow, background-color, color;
  transition-duration: 160ms;
  transition-timing-function: cubic-bezier(0.16, 1, 0.3, 1);
}

.template-quick-chip:active {
  transform: scale(0.96);
}

.template-quick-chip.active {
  background: #dce9df;
  color: #245f48;
  box-shadow: inset 0 0 0 1px rgba(47, 125, 92, 0.32);
}

.template-quick-chip-icon {
  font-size: 16px;
}

@media (hover: hover) {
  .template-quick-chip:hover {
    background: #f7f3eb;
  }
}
</style>
