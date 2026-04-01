<template>
  <NMessageProvider :max="1">
    <router-view />
  </NMessageProvider>
</template>

<script setup lang="ts">
import { NMessageProvider } from 'naive-ui'
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

onMounted(() => {
  // 检测是否在 Wails 环境中
  if (window.wails) {
    // 在 Wails 环境中，跳转到配置页面
    router.push('/cfg')
  }
})
</script>
