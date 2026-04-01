# 项目长期记忆

## k-view 前端项目技术栈

- **框架**：Vue 3.2 + Composition API（`<script setup lang="ts">`）
- **语言**：TypeScript 5.4，严格模式（strict: true）
- **构建**：Vite + @vitejs/plugin-vue
- **UI 库**：Naive UI 2.44
- **路由**：Vue Router 4，History 模式
- **HTTP**：Axios 1.13
- **图片懒加载**：vue-lazyload 3
- **样式**：Less，路径别名 `@` → `src/`、`@wailsjs` → `wailsjs/`
- **桌面端**：Wails 2（Go + Vue 前端），wailsjs/ 目录为自动生成的桥接代码

## 项目规范

- 所有 `.vue` 文件使用 `<script setup lang="ts">`，泛型 defineProps/defineEmits
- 工具函数在 `src/utils/index.ts`，API 接口类型在 `src/services/index.ts`
- 全局类型扩展（Window.$message 等）在 `src/env.d.ts`
- 不保留任何旧版 `.js` 文件，已于 2026-04-01 清理完毕
- build 命令：`vue-tsc --noEmit && vite build`（先类型检查再构建）
