<template>
  <article class="reference-card management-surface management-hover-lift" :class="{ inactive }">
    <div v-if="markerColor" class="reference-marker" :style="{ background: markerColor }"></div>
    <div class="reference-main">
      <div class="reference-head">
        <div class="reference-title-row">
          <span v-if="$slots.icon || icon" class="reference-icon" :style="iconStyle">
            <slot name="icon">{{ icon }}</slot>
          </span>
          <div class="reference-title-copy">
            <h2>{{ title }}</h2>
            <p v-if="subtitle">{{ subtitle }}</p>
          </div>
        </div>
        <slot name="status" />
      </div>
      <div v-if="$slots.meta" class="reference-meta"><slot name="meta" /></div>
      <div v-if="$slots.footer || $slots.actions" class="reference-footer">
        <slot name="footer" />
        <div v-if="$slots.actions" class="management-actions reference-actions"><slot name="actions" /></div>
      </div>
    </div>
  </article>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  title: { type: String, required: true },
  subtitle: { type: String, default: '' },
  icon: { type: String, default: '' },
  color: { type: String, default: '#2f7d5c' },
  markerColor: { type: String, default: '' },
  inactive: { type: Boolean, default: false },
})

const iconStyle = computed(() => ({
  color: props.color,
  background: `${props.color}18`,
  boxShadow: `inset 0 0 0 1px ${props.color}26`,
}))
</script>

<style scoped>
.reference-card {
  position: relative;
  display: grid;
  grid-template-columns: 6px minmax(0, 1fr);
  overflow: hidden;
  min-height: 132px;
  transition-property: transform, box-shadow;
  transition-duration: 180ms;
}

.reference-card:not(:has(.reference-marker)) {
  grid-template-columns: minmax(0, 1fr);
}

.reference-card.inactive {
  opacity: .66;
}

.reference-marker {
  height: 100%;
}

.reference-main {
  min-width: 0;
  padding: 14px 12px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.reference-head {
  display: flex;
  justify-content: space-between;
  gap: 10px;
  align-items: flex-start;
}

.reference-title-row {
  display: flex;
  min-width: 0;
  gap: 10px;
  align-items: center;
}

.reference-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex: 0 0 auto;
  font-size: 20px;
}

.reference-title-copy {
  min-width: 0;
}

.reference-title-copy h2 {
  margin: 0;
  font-size: 15px;
  font-weight: 600;
  color: #1e293b;
  letter-spacing: -.012em;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.reference-title-copy p {
  margin: 2px 0 0;
  color: #64748b;
  font-size: 12px;
}

.reference-meta {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  color: #64748b;
  font-size: 12px;
}

.reference-footer {
  margin-top: auto;
  display: flex;
  justify-content: space-between;
  gap: 10px;
  align-items: center;
}

.reference-actions {
  margin-left: auto;
}
</style>
