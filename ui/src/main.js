import { createApp } from 'vue'
import { createPinia } from 'pinia'

import ElementPlus from 'element-plus'

import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import 'element-plus/dist/index.css'
import './styles/management.css'

import App from './App.vue'
import router from './router'
import { i18n } from './i18n'

import './permission'


const app = createApp(App)

// 注册element的所有图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}

app.use(createPinia())
app.use(ElementPlus)
app.use(i18n)
app.use(router)

app.mount('#app')
