import { createApp } from 'vue'
import App from './App.vue'
import './style.css'
import './styles/responsive.css'
import 'element-plus/theme-chalk/base.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import router from './router/router'
import { formatTime } from './utils/format'
import { initTheme } from './utils/theme'

initTheme()

const app = createApp(App)

app.config.globalProperties.$formatTime = formatTime
app.use(router)
app.mount('#app')
