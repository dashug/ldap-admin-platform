import { createApp } from 'vue'

import Cookies from 'js-cookie'

import 'normalize.css/normalize.css' // a modern alternative to CSS resets

import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox, ElNotification, ElLoading } from 'element-plus'

import '@/styles/element-variables.scss' // 主题色（CSS 变量覆盖）
import '@/styles/index.scss' // global css

import App from './App.vue'
import store from './store'
import router from './router'

import 'virtual:svg-icons-register' // svg sprite
import './permission' // permission control
import setupErrorLog from './utils/error-log'
import SvgIcon from '@/components/SvgIcon/index.vue'

const app = createApp(App)

// 全局注册 Element Plus 图标（用法：<el-icon><Edit /></el-icon>）
for (const [key, comp] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, comp)
}
// 全局 svg-icon 组件
app.component('SvgIcon', SvgIcon)

app.use(store)
app.use(router)
app.use(ElementPlus, {
  locale: zhCn,
  size: Cookies.get('size') || 'default'
})

// 兼容 Options API 中 element-ui 风格的 this.$message / this.$confirm 等调用
app.config.globalProperties.$message = ElMessage
app.config.globalProperties.$msgbox = ElMessageBox
app.config.globalProperties.$alert = ElMessageBox.alert
app.config.globalProperties.$confirm = ElMessageBox.confirm
app.config.globalProperties.$prompt = ElMessageBox.prompt
app.config.globalProperties.$notify = ElNotification
app.config.globalProperties.$loading = ElLoading.service

setupErrorLog(app)

app.mount('#app')
