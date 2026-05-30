import { createStore } from 'vuex'
import getters from './getters'

// Vite：用 import.meta.glob 自动加载 ./modules 下的所有 vuex 模块
const modulesFiles = import.meta.glob('./modules/*.js', { eager: true })

const modules = Object.keys(modulesFiles).reduce((modules, modulePath) => {
  // './modules/app.js' => 'app'
  const moduleName = modulePath.replace(/^\.\/modules\/(.*)\.\w+$/, '$1')
  modules[moduleName] = modulesFiles[modulePath].default
  return modules
}, {})

const store = createStore({
  modules,
  getters
})

export default store
