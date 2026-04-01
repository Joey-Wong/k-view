/// <reference types="vite/client" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<Record<string, unknown>, Record<string, unknown>, unknown>
  export default component
}

// 扩展 Window 对象
interface Window {
  $message: import('naive-ui').MessageApiInjection
  wails?: unknown
}
