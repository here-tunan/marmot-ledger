<template>
  <section class="management-toolbar management-surface">
    <slot />
  </section>
</template>

<style scoped>
.management-toolbar {
  display: flex;
  gap: 12px;
  align-items: center;
  margin-bottom: 18px;
  padding: 14px;
  flex-wrap: wrap;
}

@media (max-width: 600px) {
  .management-toolbar {
    flex-direction: column;
    align-items: stretch;
  }

  .management-toolbar :deep(.filter-control),
  .management-toolbar :deep(.channel-filter) {
    width: 100%;
  }
}
</style>
