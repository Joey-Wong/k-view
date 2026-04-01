<template>
  <div :class="['row-wrap', noMarginBottom ? 'margin-bottom-2px' : 'margin-bottom-16px']">
    <span>{{ title }}</span>
    <slot />
  </div>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    title?: string
    noMarginBottom?: boolean
  }>(),
  {
    title: '',
    noMarginBottom: false,
  }
)
</script>

<style lang="less" scoped>
.row-wrap {
  display: flex;
  align-items: center;
  span {
    display: inline-block;
    width: 80px;
    text-align: left;
  }
}
.margin-bottom-16px {
  margin-bottom: 16px;
}
.margin-bottom-2px {
  margin-bottom: 2px;
}
</style>
