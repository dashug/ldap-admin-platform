<template>
  <div class="login-container">
    <!-- 左侧品牌栏（窄屏隐藏） -->
    <aside class="login-brand">
      <div class="brand-deco brand-deco--1" />
      <div class="brand-deco brand-deco--2" />
      <div class="brand-inner">
        <div class="brand-logo">
          <el-icon :size="26"><Connection /></el-icon>
        </div>
        <h1 class="brand-title">LDAP 管理平台</h1>
        <p class="brand-tagline">面向 OpenLDAP / Active Directory 的统一身份与目录管理</p>
        <ul class="brand-features">
          <li><el-icon><Select /></el-icon> OpenLDAP / AD 用户与组织统一管理</li>
          <li><el-icon><Select /></el-icon> 基于 Casbin 的 RBAC 与操作审计</li>
          <li><el-icon><Select /></el-icon> 钉钉 / 企业微信 / 飞书 一键同步</li>
        </ul>
      </div>
      <p class="brand-footer">Powered by LDAP 管理平台 · 前后端一体</p>
    </aside>

    <!-- 右侧登录表单 -->
    <section class="login-panel">
      <div class="login-card">
        <div class="login-header">
          <h2 class="login-title">欢迎回来 👋</h2>
          <p class="login-desc">请登录以继续使用管理后台</p>
        </div>
        <el-form ref="loginForm" :model="loginForm" :rules="loginRules" class="login-form" autocomplete="on" label-position="top" @submit.prevent="handleLogin">
          <el-form-item prop="username">
            <el-input
              ref="username"
              v-model="loginForm.username"
              size="large"
              placeholder="请输入用户名"
              name="username"
              type="text"
              tabindex="1"
              autocomplete="on"
              prefix-icon="User"
            />
          </el-form-item>
          <el-tooltip v-model:visible="capsTooltip" content="Caps Lock 已打开" placement="top" manual>
            <el-form-item prop="password">
              <el-input
                :key="passwordType"
                ref="password"
                v-model="loginForm.password"
                size="large"
                :type="passwordType"
                placeholder="请输入密码（至少 6 位）"
                name="password"
                tabindex="2"
                autocomplete="on"
                show-password
                prefix-icon="Lock"
                @keyup="checkCapslock"
                @blur="capsTooltip = false"
                @keyup.enter="handleLogin"
              />
            </el-form-item>
          </el-tooltip>
          <div class="login-actions">
            <span class="link-forget" @click="changePass">忘记密码？</span>
            <el-button :loading="loading" type="primary" size="large" class="login-btn" @click.prevent="handleLogin">登 录</el-button>
          </div>
        </el-form>
      </div>
    </section>
  </div>
</template>

<script>
import JSEncrypt from 'jsencrypt'
import { getPublicKey } from '@/api/system/base'

export default {
  name: 'Login',
  data() {
    const validatePassword = (rule, value, callback) => {
      if (value.length < 6) {
        callback(new Error('The password can not be less than 6 digits'))
      } else {
        callback()
      }
    }
    return {
      loginForm: {
        username: '',
        password: ''
      },
      loginRules: {
        username: [{ required: true, trigger: 'blur' }],
        password: [{ required: true, trigger: 'blur', validator: validatePassword }]
      },
      passwordType: 'password',
      publicKey: import.meta.env.VITE_APP_PUBLIC_KEY || '',
      capsTooltip: false,
      loading: false,
      publicKeyLoading: true,
      redirect: undefined,
      otherQuery: {}
    }
  },
  watch: {
    $route: {
      handler: function(route) {
        const query = route.query
        if (query) {
          this.redirect = query.redirect
          this.otherQuery = this.getOtherQuery(query)
        }
      },
      immediate: true
    }
  },
  created() {
    this.fetchPublicKeyIfNeeded()
  },
  mounted() {
    if (this.loginForm.username === '') {
      this.$refs.username.focus()
    } else if (this.loginForm.password === '') {
      this.$refs.password.focus()
    }
  },
  destroyed() {
    // window.removeEventListener('storage', this.afterQRScan)
  },
  methods: {
    fetchPublicKeyIfNeeded() {
      if (this.publicKey && this.publicKey.trim()) {
        this.publicKeyLoading = false
        return
      }
      getPublicKey()
        .then(res => {
          const key = res && res.data
          if (key && typeof key === 'string' && key.trim()) {
            this.publicKey = key.trim()
          }
        })
        .catch(() => {})
        .finally(() => {
          this.publicKeyLoading = false
        })
    },
    checkCapslock(e) {
      const { key } = e
      this.capsTooltip = key && key.length === 1 && (key >= 'A' && key <= 'Z')
    },
    showPwd() {
      if (this.passwordType === 'password') {
        this.passwordType = ''
      } else {
        this.passwordType = 'password'
      }
      this.$nextTick(() => {
        this.$refs.password.focus()
      })
    },
    handleLogin() {
      this.$refs.loginForm.validate(valid => {
        if (valid) {
          if (!this.publicKey || !this.publicKey.trim()) {
            this.$message.error('正在获取加密公钥，请稍候再试')
            this.fetchPublicKeyIfNeeded()
            return
          }
          this.loading = true
          const encryptor = new JSEncrypt()
          encryptor.setPublicKey(this.publicKey)
          const encPassword = encryptor.encrypt(this.loginForm.password)
          if (!encPassword) {
            this.loading = false
            this.$message.error('密码加密失败，请刷新页面重试')
            return
          }
          const encLoginForm = { username: this.loginForm.username, password: encPassword }
          this.$store.dispatch('user/login', encLoginForm)
            .then(() => {
              this.$router.push({ path: this.redirect || '/', query: this.otherQuery })
              this.loading = false
            })
            .catch(() => {
              this.loading = false
            })
        } else {
          return false
        }
      })
    },
    changePass() {
      // window.location.href='/changePass'
      this.$router.push({ path: '/changePass' })
    },
    getOtherQuery(query) {
      return Object.keys(query).reduce((acc, cur) => {
        if (cur !== 'redirect') {
          acc[cur] = query[cur]
        }
        return acc
      }, {})
    }
  }
}
</script>

<style lang="scss" scoped>
@import "@/styles/variables.scss";

.login-container {
  min-height: 100vh;
  width: 100%;
  display: flex;
  background: $slate50;
}

/* 左侧品牌栏 */
.login-brand {
  position: relative;
  flex: 1 1 0;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 64px 56px;
  overflow: hidden;
  color: #fff;
  background: linear-gradient(150deg, $themePrimaryDark 0%, $themePrimary 55%, $themePrimaryLight 130%);
}
.brand-deco {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.08);
  pointer-events: none;
  &--1 { width: 420px; height: 420px; top: -140px; right: -120px; }
  &--2 { width: 280px; height: 280px; bottom: -100px; left: -80px; background: rgba(255, 255, 255, 0.06); }
}
.brand-inner { position: relative; max-width: 460px; }
.brand-logo {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 56px;
  height: 56px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.16);
  backdrop-filter: blur(4px);
  margin-bottom: 28px;
}
.brand-title {
  margin: 0 0 14px;
  font-size: 34px;
  font-weight: 700;
  letter-spacing: -0.02em;
}
.brand-tagline {
  margin: 0 0 40px;
  font-size: 16px;
  line-height: 1.7;
  color: rgba(255, 255, 255, 0.85);
}
.brand-features {
  list-style: none;
  margin: 0;
  padding: 0;
  li {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 18px;
    font-size: 15px;
    color: rgba(255, 255, 255, 0.92);
    .el-icon {
      flex: none;
      width: 22px;
      height: 22px;
      border-radius: 50%;
      background: rgba(255, 255, 255, 0.18);
      font-size: 13px;
    }
  }
}
.brand-footer {
  position: relative;
  margin: 0;
  font-size: 13px;
  color: rgba(255, 255, 255, 0.6);
}

/* 右侧表单 */
.login-panel {
  flex: 0 0 480px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
  background: #fff;
}
.login-card {
  width: 100%;
  max-width: 360px;
}
.login-header {
  margin-bottom: 36px;
  .login-title {
    margin: 0 0 10px;
    font-size: 28px;
    font-weight: 700;
    color: $slate900;
    letter-spacing: -0.02em;
  }
  .login-desc {
    margin: 0;
    font-size: 15px;
    color: $slate500;
  }
}
.login-form {
  :deep(.el-form-item) { margin-bottom: 22px; }
  :deep(.el-input__wrapper) { border-radius: 12px; }
}
.login-actions {
  display: flex;
  flex-direction: column;
  align-items: stretch;
  gap: 16px;
  margin-top: 28px;
  .link-forget {
    align-self: flex-end;
    font-size: 14px;
    color: $slate500;
    cursor: pointer;
    &:hover { color: $themePrimary; }
  }
  .login-btn {
    width: 100%;
    height: 48px;
    font-size: 16px;
    font-weight: 600;
    border-radius: 12px;
  }
}

/* 窄屏：隐藏品牌栏，表单铺满 */
@media (max-width: 860px) {
  .login-brand { display: none; }
  .login-panel { flex: 1 1 auto; }
}
</style>
