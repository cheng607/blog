import { createApp } from 'vue'
import App from './App.vue'
import './style.css'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import router from './router/router'
import { formatTime } from './utils/format'

const app = createApp(App)

app.config.globalProperties.$formatTime = formatTime;
app.use(ElementPlus)
app.use(router)
app.mount('#app')

