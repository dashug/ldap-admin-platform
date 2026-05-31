import store from '@/store'
import { isString, isArray } from '@/utils/validate'
import settings from '@/settings'

// you can set in settings.js
// errorLog:'production' | ['production', 'development']
const { errorLog: needErrorLog } = settings

function checkNeed() {
  const env = import.meta.env.MODE
  if (isString(needErrorLog)) {
    return env === needErrorLog
  }
  if (isArray(needErrorLog)) {
    return needErrorLog.includes(env)
  }
  return false
}

// Vue 3：错误处理器挂在 app 实例上，由 main.js 在创建 app 后调用
export default function setupErrorLog(app) {
  if (!checkNeed()) return
  app.config.errorHandler = function(err, vm, info) {
    store.dispatch('errorLog/addErrorLog', {
      err,
      vm,
      info,
      url: window.location.href
    })
    console.error(err, info)
  }
}
