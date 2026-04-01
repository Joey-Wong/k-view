import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/mv',
  },
  {
    path: '/cfg',
    name: 'Cfg',
    component: () => import('@/views/cfg.vue'),
  },
  {
    path: '/mv',
    name: 'MV',
    component: () => import('@/views/mv.vue'),
  },
  {
    path: '/img',
    name: 'Img',
    component: () => import('@/views/img.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
