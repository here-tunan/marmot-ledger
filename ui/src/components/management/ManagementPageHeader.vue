<template>
  <section class="management-page-header management-surface">
    <div>
      <p class="eyebrow">{{ eyebrow }}</p>
      <h1>{{ title }}</h1>
      <p v-if="subtitle">{{ subtitle }}</p>
    </div>
    <div v-if="$slots.actions" class="management-page-header-actions">
      <slot name="actions" />
    </div>
  </section>
</template>

<script setup>
defineProps({
  eyebrow: { type: String, default: '' },
  title: { type: String, required: true },
  subtitle: { type: String, default: '' },
})
</script>

<style scoped>
.management-page-header {
  display: flex;
  justify-content: space-between;
  gap: 24px;
  align-items: flex-start;
  margin-bottom: 18px;
  padding: 26px;
  background: linear-gradient(135deg, #fffaf0 0%, #fff 70%);
}

.eyebrow {
  margin: 0 0 8px;
  color: #2f7d5c;
  font-size: 12px;
  font-weight: 800;
  letter-spacing: .08em;
  text-transform: uppercase;
}

.management-page-header h1 {
  max-width: 720px;
  margin: 0;
  font-size: 30px;
  line-height: 1.16;
  letter-spacing: -.022em;
  text-wrap: balance;
}

.management-page-header p:last-child {
  max-width: 680px;
  margin: 12px 0 0;
  color: #64748b;
  line-height: 1.7;
  text-wrap: pretty;
}

.management-page-header-actions {
  display: flex;
  gap: 10px;
  align-items: center;
  flex-wrap: wrap;
  justify-content: flex-end;
}

@media (max-width: 600px) {
  .management-page-header {
    flex-direction: column;
    align-items: stretch;
  }

  .management-page-header-actions {
    justify-content: flex-start;
  }
}
</style>
