<template>
  <div class="login-container">
    <div class="login-bg" />
    <div class="login-card">
      <div class="login-header">
        <h1 class="login-title">LDAP 管理平台</h1>
        <p class="login-desc">统一身份与目录管理</p>
      </div>
      <el-form ref="loginForm" :model="loginForm" :rules="loginRules" class="login-form" autocomplete="on" label-position="top" @submit.native.prevent="handleLogin">
        <el-form-item prop="username">
          <el-input
            ref="username"
            v-model="loginForm.username"
            placeholder="请输入用户名"
            name="username"
            type="text"
            tabindex="1"
            autocomplete="on"
            prefix-icon="el-icon-user"
          />
        </el-form-item>
        <el-tooltip v-model="capsTooltip" content="Caps Lock 已打开" placement="top" manual>
          <el-form-item prop="password">
            <el-input
              :key="passwordType"
              ref="password"
              v-model="loginForm.password"
              :type="passwordType"
              placeholder="请输入密码（至少 6 位）"
              name="password"
              tabindex="2"
              autocomplete="on"
              show-password
              prefix-icon="el-icon-lock"
              @keyup.native="checkCapslock"
              @blur="capsTooltip = false"
              @keyup.enter.native="handleLogin"
            />
          </el-form-item>
        </el-tooltip>
        <div class="login-actions">
          <span class="link-forget" @click="changePass">忘记密码</span>
          <el-button :loading="loading" type="primary" class="login-btn" @click.native.prevent="handleLogin">登 录</el-button>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script>
import JSEncrypt from 'jsencrypt'

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
      publicKey: process.env.VUE_APP_PUBLIC_KEY,
      capsTooltip: false,
      loading: false,
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
    // window.addEventListener('storage', this.afterQRScan)
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
          this.loading = true
          // 密码RSA加密处理
          const encryptor = new JSEncrypt()
          // 设置公钥
          encryptor.setPublicKey(this.publicKey)
          // 加密密码
          const encPassword = encryptor.encrypt(this.loginForm.password)
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
@import "~@/styles/variables.scss";

.login-container {
  min-height: 100vh;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
}

.login-bg {
  position: absolute;
  inset: 0;
  /* 本地图：将图片放到 src/assets/backgd-image/ 后改为 url("~@/assets/backgd-image/你的图片.jpg") */
  background: linear-gradient(145deg, #0f172a 0%, #1e293b 40%, #334155 100%);
  background-image: url("https://images.unsplash.com/photo-1557683316-973673baf926?w=1920&q=80");
  background-size: cover;
  background-position: center;
  &::before {
    content: "";
    position: absolute;
    inset: 0;
    background: linear-gradient(160deg, rgba(15, 23, 42, 0.82) 0%, rgba(30, 41, 59, 0.78) 50%, rgba(51, 65, 81, 0.85) 100%);
    pointer-events: none;
  }
  &::after {
    content: "";
    position: absolute;
    inset: 0;
    background: radial-gradient(ellipse 80% 50% at 50% 120%, rgba(79, 70, 229, 0.15) 0%, transparent 50%);
    pointer-events: none;
  }
}

.login-card {
  position: relative;
  width: 420px;
  max-width: calc(100vw - 32px);
  padding: 52px 44px;
  background: #fff;
  border-radius: 20px;
  box-shadow: 0 32px 64px rgba(0, 0, 0, 0.2), 0 0 0 1px rgba(255, 255, 255, 0.05) inset;
}

.login-header {
  text-align: center;
  margin-bottom: 40px;
  .login-title {
    margin: 0 0 10px;
    font-size: 26px;
    font-weight: 700;
    color: #0f172a;
    letter-spacing: -0.02em;
  }
  .login-desc {
    margin: 0;
    font-size: 15px;
    color: #64748b;
    font-weight: 400;
  }
}

.login-form {
  ::v-deep .el-form-item {
    margin-bottom: 24px;
  }
  ::v-deep .el-input__inner {
    height: 50px;
    line-height: 50px;
    border-radius: 12px;
    border-color: #e2e8f0;
    &:focus {
      border-color: $themePrimary;
    }
  }
}

.login-actions {
  display: flex;
  flex-direction: column;
  align-items: stretch;
  gap: 18px;
  margin-top: 32px;
  .link-forget {
    text-align: right;
    font-size: 14px;
    color: #64748b;
    cursor: pointer;
    &:hover {
      color: $themePrimary;
    }
  }
  .login-btn {
    height: 50px;
    font-size: 16px;
    font-weight: 600;
    border-radius: 12px;
    background: $themePrimary;
    border-color: $themePrimary;
    &:hover, &:focus {
      background: $themePrimaryDark;
      border-color: $themePrimaryDark;
    }
  }
}
</style>
