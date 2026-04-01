import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import VueLazyload from 'vue-lazyload'

const app = createApp(App)
app.use(router)
app.use(VueLazyload, {
  loading: '',
  error: '',
  attempt: 1,
})
app.mount('#app')
