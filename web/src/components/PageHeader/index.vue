<template>
  <div class="page-header">
    <div class="page-header__main">
      <h1 class="page-header__title">{{ title }}</h1>
      <p v-if="subtitle" class="page-header__subtitle">{{ subtitle }}</p>
    </div>
    <div v-if="$slots.actions" class="page-header__actions">
      <slot name="actions" />
    </div>
  </div>
</template>

<script>
export default {
  name: 'PageHeader',
  props: {
    title: { type: String, default: '' },
    subtitle: { type: String, default: '' }
  }
}
</script>

<style lang="scss" scoped>
@import "@/styles/variables.scss";

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 18px;
  flex-wrap: wrap;

  &__title {
    margin: 0;
    font-size: 22px;
    font-weight: $fontWeightBold;
    color: $slate900;
    letter-spacing: -0.02em;
  }
  &__subtitle {
    margin: 4px 0 0;
    font-size: 13px;
    color: $slate500;
  }
  &__actions {
    display: flex;
    align-items: center;
    gap: 10px;
  }
}
</style>
