import axios from 'axios'
import { ElMessageBox as MessageBox, ElMessage as Message } from 'element-plus'
import store from '@/store'
// import router from '@/router'
import { getToken } from '@/utils/auth'
import router from '@/router'

// create an axios instance
const service = axios.create({
  baseURL: import.meta.env.PROD ? (import.meta.env.VITE_APP_BASE_API || '/') : '/', // api 的 base_url
  // withCredentials: true, // send cookies when cross-domain requests
  // timeout: 5000 // request timeout
})

// request interceptor
service.interceptors.request.use(
  config => {
    // do something before request is sent
    if (store.getters.token) {
      // let each request carry token
      // ['X-Token'] is a custom headers key
      // please modify it according to the actual situation
      config.headers['Authorization'] = 'Bearer ' + getToken()
      // config.headers['Content-Type'] = 'application/json'
    }
    return config
  },
  error => {
    // do something with request error
    return Promise.reject(error)
  }
)

// response interceptor
service.interceptors.response.use(
  /**
   * If you want to get http information such as headers or status
   * Please return  response => response
  */

  /**
   * Determine the request status by custom code
   * Here is just an example
   * You can also judge the status by HTTP Status Code
   */
  response => {
    const res = response.data
    if (res.code === 0 || res.code === 200) {
      return res
    }
    const msg = res.msg || res.message || '请求失败'
    Message({
      message: msg,
      type: 'error'
    })
    return Promise.reject(new Error(msg))
  },
  error => {
    const response = error.response
    const dataMsg = response && response.data && (response.data.message || response.data.msg)
    const message = dataMsg || error.message || '网络异常或服务不可用'

    if (!response) {
      Message({
        showClose: true,
        message,
        type: 'error',
        duration: 5 * 1000
      })
      return Promise.reject(error)
    }

    if (response.status === 401) {
      // 登录时要求动态验证码：交由登录页内联处理，不弹通用「重新登录」框
      if (typeof dataMsg === 'string' && dataMsg.indexOf('动态验证码') !== -1) {
        return Promise.reject(error)
      }
      if (typeof dataMsg === 'string' && dataMsg.indexOf('JWT认证失败') !== -1) {
        MessageBox.confirm(
          '登录失败,用户名或密码错误,重新登录或继续停留在当前页？',
          '登录状态已失效',
          {
            confirmButtonText: '重新登录',
            cancelButtonText: '继续停留',
            type: 'warning'
          }
        ).then(() => {
          store.dispatch('user/logout').then(() => {
            location.reload()
          })
        })
      } else {
        Message({
          showClose: true,
          message,
          type: 'error',
          duration: 5 * 1000
        })
      }
      return Promise.reject(error)
    }
    if (response.status === 403) {
      router.push({ path: '/401' })
      return Promise.reject(error)
    }

    Message({
      showClose: true,
      message,
      type: 'error',
      duration: 5 * 1000
    })
    return Promise.reject(error)
  }
)

export default service
