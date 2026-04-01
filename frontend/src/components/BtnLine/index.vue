<template>
  <div :class="['row-wrap', noMarginBottom ? 'margin-bottom-2px' : 'margin-bottom-16px']">
    <span>{{ title }}</span>
    <n-button :disabled="disabled" :type="type" @click="emit('onClick')">
      {{ text }}
    </n-button>
  </div>
</template>

<script setup lang="ts">
import { NButton, type ButtonType } from 'naive-ui'

withDefaults(
  defineProps<{
    title?: string
    text?: string
    disabled?: boolean
    type?: ButtonType
    noMarginBottom?: boolean
  }>(),
  {
    title: '',
    text: '',
    disabled: false,
    type: 'default',
    noMarginBottom: false,
  }
)

const emit = defineEmits<{
  (e: 'onClick'): void
}>()
</script>

<style lang="less" scoped>
.margin-bottom-16px {
  margin-bottom: 16px;
}
.margin-bottom-2px {
  margin-bottom: 2px;
}
.row-wrap {
  display: flex;
  align-items: center;
  span {
    display: inline-block;
    width: 80px;
    text-align: left;
  }
}
</style>
