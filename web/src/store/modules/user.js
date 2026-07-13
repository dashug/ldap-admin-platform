import { login, logout, refreshToken } from '@/api/system/base'
import { getInfo } from '@/api/system/user'
import { getToken, setToken, removeToken } from '@/utils/auth'
import { resetRouter } from '@/router'

const state = {
  token: getToken(),
  name: '',
  avatar: '',
  introduction: '',
  roles: [],
  // 是否为管理员（角色ID为1），用于路由守卫与侧边栏入口可见性
  isAdmin: false
}

const mutations = {
  SET_TOKEN: (state, token) => {
    state.token = token
  },
  SET_INTRODUCTION: (state, introduction) => {
    state.introduction = introduction
  },
  SET_NAME: (state, name) => {
    state.name = name
  },
  SET_AVATAR: (state, avatar) => {
    state.avatar = avatar
  },
  SET_ROLES: (state, roles) => {
    state.roles = roles
  },
  SET_IS_ADMIN: (state, isAdmin) => {
    state.isAdmin = isAdmin
  }
}

const actions = {
  
  // user login
  login({ commit }, userInfo) {

    const { username, password, otp } = userInfo
    return new Promise((resolve, reject) => {
      login({ username: username.trim(), password: password, otp: otp || '' }).then(response => {
        const { data } = response
        commit('SET_TOKEN', data.token)
        setToken(data.token)
        resolve()
      }).catch(error => {
        //  this.$router.push({ path: '/dashboard'})
        //  return false;
        reject(error)
      })
    })
  },

  // get user info
  getInfo({ commit, state }) {
    return new Promise((resolve, reject) => {
      getInfo(state.token).then(response => {
        const { data } = response
        if (!data) {
          reject('Verification failed, please Login again.')
        }

        const userInfo = data
        const { roles, username, avatar, introduction } = userInfo
        // roles must be a non-empty array
        if (!roles || roles.length <= 0) {
          reject('getInfo: roles must be a non-null array!')
        }
        const rolesNames = []
        roles.forEach(x => { rolesNames.push(x.name) })
        commit('SET_ROLES', rolesNames)
        // 管理员判定与 permission.js 保持一致：角色ID为1
        commit('SET_IS_ADMIN', roles.some(role => role.ID === 1))
        commit('SET_NAME', username)
        commit('SET_AVATAR', avatar)
        commit('SET_INTRODUCTION', introduction)
        resolve(userInfo)
      }).catch(error => {
        reject(error)
      })
    })
  },

  // user logout
  logout({ commit, state, dispatch }) {
    return new Promise((resolve, reject) => {
      logout(state.token).then(() => {
        commit('SET_TOKEN', '')
        commit('SET_ROLES', [])
        commit('SET_IS_ADMIN', false)
        removeToken()
        resetRouter()

        // reset visited views and cached views
        // to fixed https://github.com/PanJiaChen/vue-element-admin/issues/2485
        dispatch('tagsView/delAllViews', null, { root: true })

        resolve()
      }).catch(error => {
        reject(error)
      })
    })
  },

  // refresh token
  async refreshToken({ commit }) {
    // 刷新token
    const { data } = await refreshToken()
    commit('SET_TOKEN', data.token)
    setToken(data.token)
  },

  // remove token
  resetToken({ commit }) {
    return new Promise(resolve => {
      commit('SET_TOKEN', '')
      commit('SET_ROLES', [])
      commit('SET_IS_ADMIN', false)
      removeToken()
      resolve()
    })
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
