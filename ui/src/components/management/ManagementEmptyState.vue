<template>
  <div class="management-empty-state management-surface">
    <img v-if="image" :src="image" :alt="alt" width="112" height="112" />
    <slot name="icon" />
    <h2>{{ title }}</h2>
    <p>{{ text }}</p>
    <slot name="action" />
  </div>
</template>

<script setup>
defineProps({
  image: { type: String, default: '' },
  alt: { type: String, default: '' },
  title: { type: String, required: true },
  text: { type: String, default: '' },
})
</script>

<style scoped>
.management-empty-state {
  grid-column: 1 / -1;
  display: grid;
  place-items: center;
  gap: 12px;
  padding: 40px 24px;
  text-align: center;
  color: #64748b;
}

.management-empty-state img {
  border-radius: 22px;
}

.management-empty-state h2 {
  margin: 0;
  color: #1e293b;
}

.management-empty-state p {
  max-width: 420px;
  margin: 0;
  line-height: 1.7;
}
</style>
