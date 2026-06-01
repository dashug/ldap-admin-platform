<template>
  <div class="login">
    <!-- 动效背景：极光渐变 + 网格 + 漂浮光球 -->
    <div class="login__bg" aria-hidden="true">
      <span class="aurora aurora--1" />
      <span class="aurora aurora--2" />
      <span class="aurora aurora--3" />
      <span class="aurora aurora--4" />
      <div class="login__grid" />
      <span class="orb orb--1" />
      <span class="orb orb--2" />
      <span class="orb orb--3" />
    </div>

    <!-- 登录卡片 -->
    <div class="login__card">
      <!-- 互动机器人吉祥物：眼睛随用户名移动、输入密码时捂眼 -->
      <div
        class="mascot"
        :class="{ 'is-username': focusField === 'username', 'is-password': focusField === 'password', 'is-otp': focusField === 'otp' }"
        aria-hidden="true"
      >
        <svg class="mascot__svg" viewBox="0 0 220 200" xmlns="http://www.w3.org/2000/svg">
          <defs>
            <linearGradient id="mHead" x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stop-color="#7c80f2" />
              <stop offset="100%" stop-color="#4f46e5" />
            </linearGradient>
            <linearGradient id="mHand" x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stop-color="#8b8ff5" />
              <stop offset="100%" stop-color="#5b54e6" />
            </linearGradient>
            <radialGradient id="mTip" cx="50%" cy="40%" r="60%">
              <stop offset="0%" stop-color="#a5f3fc" />
              <stop offset="100%" stop-color="#22d3ee" />
            </radialGradient>
          </defs>

          <!-- 天线 -->
          <line class="mascot__antenna" x1="110" y1="46" x2="110" y2="30" stroke="#6366f1" stroke-width="4" stroke-linecap="round" />
          <circle class="mascot__tip" cx="110" cy="23" r="7" fill="url(#mTip)" />

          <!-- 头 -->
          <rect x="48" y="44" width="124" height="106" rx="32" fill="url(#mHead)" />
          <rect x="44" y="78" width="10" height="34" rx="5" fill="#4f46e5" />
          <rect x="166" y="78" width="10" height="34" rx="5" fill="#4f46e5" />
          <!-- 脸屏 -->
          <rect x="64" y="62" width="92" height="72" rx="26" fill="#2b2769" />

          <!-- 眼睛（含眨眼 + 瞳孔跟随） -->
          <g class="mascot__eye mascot__eye--l">
            <ellipse cx="94" cy="96" rx="13" ry="16" fill="#ffffff" />
            <g class="mascot__pupil" :style="pupilStyle">
              <circle cx="94" cy="98" r="6.5" fill="#1e1b4b" />
              <circle cx="96.4" cy="95.6" r="2.1" fill="#ffffff" />
            </g>
          </g>
          <g class="mascot__eye mascot__eye--r">
            <ellipse cx="126" cy="96" rx="13" ry="16" fill="#ffffff" />
            <g class="mascot__pupil" :style="pupilStyle">
              <circle cx="126" cy="98" r="6.5" fill="#1e1b4b" />
              <circle cx="128.4" cy="95.6" r="2.1" fill="#ffffff" />
            </g>
          </g>

          <!-- 腮红 -->
          <circle cx="78" cy="120" r="6" fill="#f472b6" opacity="0.55" />
          <circle cx="142" cy="120" r="6" fill="#f472b6" opacity="0.55" />
          <!-- 嘴 -->
          <path class="mascot__mouth" d="M100 119 Q110 128 120 119" stroke="#1e1b4b" stroke-width="3.4" stroke-linecap="round" fill="none" />

          <!-- 双手（输入密码时上抬捂眼） -->
          <g class="mascot__arm mascot__arm--l">
            <ellipse cx="74" cy="150" rx="20" ry="15" fill="url(#mHand)" />
          </g>
          <g class="mascot__arm mascot__arm--r">
            <ellipse cx="146" cy="150" rx="20" ry="15" fill="url(#mHand)" />
          </g>
        </svg>
      </div>

      <div class="login__head">
        <h1 class="login__title">LDAP 管理平台</h1>
        <p class="login__sub">{{ greeting }}</p>
      </div>

      <p v-if="isDemo" class="login__demo">演示账号 <b>admin</b> / <b>123456</b>　已预填，点「登录」即可</p>

      <el-form ref="loginForm" :model="loginForm" :rules="loginRules" class="login__form" autocomplete="on" label-position="top" @submit.prevent="handleLogin">
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
            @focus="onUserFocus"
            @input="onUserInput"
            @blur="onFieldBlur"
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
              @focus="onPwdFocus"
              @keyup="checkCapslock"
              @blur="onPwdBlur"
              @keyup.enter="handleLogin"
            />
          </el-form-item>
        </el-tooltip>
        <el-form-item v-if="mfaRequired" prop="otp">
          <el-input
            ref="otp"
            v-model.trim="loginForm.otp"
            size="large"
            maxlength="6"
            placeholder="请输入 6 位动态验证码"
            name="otp"
            tabindex="3"
            prefix-icon="Key"
            @focus="onOtpFocus"
            @blur="onFieldBlur"
            @keyup.enter="handleLogin"
          />
        </el-form-item>
        <div class="login__actions">
          <span class="login__forget" @click="changePass">忘记密码？</span>
          <el-button :loading="loading" type="primary" size="large" class="login__btn" @click.prevent="handleLogin">登 录</el-button>
        </div>
      </el-form>
    </div>
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
      // 吉祥物互动状态
      focusField: '',
      pupilX: 0,
      pupilY: 0,
      // 是否演示环境（托管域名）：预填账号并展示演示提示
      isDemo: false,
      mfaRequired: false,
      loginForm: {
        username: '',
        password: '',
        otp: ''
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
  computed: {
    pupilStyle() {
      return { transform: `translate(${this.pupilX}px, ${this.pupilY}px)` }
    },
    greeting() {
      if (this.focusField === 'password') return '已经闭眼啦，请放心输入密码 🙈'
      if (this.focusField === 'username') return '在看你输入用户名…'
      if (this.focusField === 'otp') return '请输入动态验证码 🔐'
      return '请登录以继续使用管理后台'
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
    // 演示环境（Render / Fly / Railway 等托管域名）自动预填管理员账号，方便一键体验。
    // 仅作用于这些演示域名，自有域名/本地/正式部署不受影响。
    this.isDemo = /onrender\.com$|\.fly\.dev$|\.up\.railway\.app$/.test(location.hostname)
    if (this.isDemo) {
      this.loginForm.username = 'admin'
      this.loginForm.password = '123456'
    }
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
    // ----- 吉祥物互动 -----
    onUserFocus() {
      this.focusField = 'username'
      this.trackPupil()
    },
    onUserInput() {
      // 同时兜底设置状态：避免 focus 事件因已聚焦而未触发时，问候语/表情不同步
      if (this.focusField !== 'password') this.focusField = 'username'
      this.trackPupil()
    },
    onPwdFocus() {
      this.focusField = 'password'
    },
    onPwdBlur() {
      this.capsTooltip = false
      if (this.focusField === 'password') this.focusField = ''
      this.resetPupil()
    },
    onOtpFocus() {
      this.focusField = 'otp'
      this.resetPupil()
    },
    onFieldBlur() {
      this.focusField = ''
      this.resetPupil()
    },
    trackPupil() {
      const len = (this.loginForm.username || '').length
      const ratio = Math.min(len / 16, 1) // 0..1
      this.pupilX = -3.5 + ratio * 7 // 随输入向右漂移
      this.pupilY = 3 // 略向下看着输入框
    },
    resetPupil() {
      this.pupilX = 0
      this.pupilY = 0
    },
    // ----- 原有登录逻辑（保持不变） -----
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
          const encLoginForm = { username: this.loginForm.username, password: encPassword, otp: this.loginForm.otp }
          this.$store.dispatch('user/login', encLoginForm)
            .then(() => {
              this.$router.push({ path: this.redirect || '/', query: this.otherQuery })
              this.loading = false
            })
            .catch((error) => {
              this.loading = false
              const msg = (error && error.response && error.response.data && (error.response.data.msg || error.response.data.message)) || (error && error.message) || ''
              if (msg.indexOf('动态验证码') !== -1) {
                const firstPrompt = !this.mfaRequired
                this.mfaRequired = true
                this.$nextTick(() => { this.$refs.otp && this.$refs.otp.focus() })
                if (firstPrompt || !this.loginForm.otp) {
                  this.$message.info('该账号已开启二次验证，请输入动态验证码')
                } else {
                  this.$message.error('动态验证码错误，请重试')
                }
              }
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

.login {
  position: relative;
  min-height: 100vh;
  min-height: 100dvh;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
  overflow: hidden;
  background: #080b18;
}

/* ---------- 动效背景 ---------- */
.login__bg { position: absolute; inset: 0; overflow: hidden; pointer-events: none; }
.aurora {
  position: absolute;
  border-radius: 50%;
  filter: blur(90px);
  opacity: 0.55;
  will-change: transform;
  &--1 { width: 46vw; height: 46vw; min-width: 360px; min-height: 360px; top: -12%; left: -6%;
    background: radial-gradient(circle, #6366f1, transparent 70%); animation: drift1 20s ease-in-out infinite; }
  &--2 { width: 42vw; height: 42vw; min-width: 320px; min-height: 320px; bottom: -16%; right: -6%;
    background: radial-gradient(circle, #8b5cf6, transparent 70%); animation: drift2 24s ease-in-out infinite; }
  &--3 { width: 34vw; height: 34vw; min-width: 280px; min-height: 280px; top: 24%; right: 16%; opacity: 0.42;
    background: radial-gradient(circle, #22d3ee, transparent 70%); animation: drift3 28s ease-in-out infinite; }
  &--4 { width: 30vw; height: 30vw; min-width: 240px; min-height: 240px; bottom: 8%; left: 20%; opacity: 0.38;
    background: radial-gradient(circle, #d946ef, transparent 70%); animation: drift4 32s ease-in-out infinite; }
}
@keyframes drift1 { 0%, 100% { transform: translate(0, 0) scale(1); } 50% { transform: translate(6%, 9%) scale(1.14); } }
@keyframes drift2 { 0%, 100% { transform: translate(0, 0) scale(1); } 50% { transform: translate(-7%, -6%) scale(1.18); } }
@keyframes drift3 { 0%, 100% { transform: translate(0, 0) scale(1); } 50% { transform: translate(-8%, 7%) scale(1.1); } }
@keyframes drift4 { 0%, 100% { transform: translate(0, 0) scale(1); } 50% { transform: translate(9%, -8%) scale(1.22); } }

.login__grid {
  position: absolute; inset: 0;
  background-image:
    linear-gradient(rgba(255, 255, 255, 0.04) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255, 255, 255, 0.04) 1px, transparent 1px);
  background-size: 46px 46px;
  -webkit-mask-image: radial-gradient(ellipse at center, #000 25%, transparent 72%);
  mask-image: radial-gradient(ellipse at center, #000 25%, transparent 72%);
}
.orb {
  position: absolute;
  border-radius: 50%;
  border: 1px solid rgba(255, 255, 255, 0.14);
  background: radial-gradient(circle at 30% 30%, rgba(255, 255, 255, 0.10), transparent 70%);
  will-change: transform;
  &--1 { width: 120px; height: 120px; top: 14%; left: 14%; animation: float 9s ease-in-out infinite; }
  &--2 { width: 64px; height: 64px; top: 72%; left: 22%; animation: float 11s ease-in-out infinite 0.6s; }
  &--3 { width: 90px; height: 90px; top: 24%; right: 12%; animation: float 13s ease-in-out infinite 1.2s; }
}
@keyframes float { 0%, 100% { transform: translateY(0); } 50% { transform: translateY(-20px); } }

/* ---------- 卡片 ---------- */
.login__card {
  position: relative;
  z-index: 1;
  width: min(420px, 100%);
  padding: 8px 40px 40px;
  border-radius: 24px;
  background: rgba(255, 255, 255, 0.95);
  border: 1px solid rgba(255, 255, 255, 0.6);
  box-shadow: 0 30px 80px -24px rgba(2, 6, 23, 0.85);
  backdrop-filter: blur(20px) saturate(140%);
  -webkit-backdrop-filter: blur(20px) saturate(140%);
  animation: cardIn 0.7s cubic-bezier(0.2, 0.7, 0.2, 1) both;
}
@keyframes cardIn { from { opacity: 0; transform: translateY(26px) scale(0.985); } to { opacity: 1; transform: none; } }

/* ---------- 吉祥物 ---------- */
.mascot {
  width: 168px;
  height: 152px;
  margin: 0 auto;
  animation: bob 4.5s ease-in-out infinite;
}
.mascot__svg { width: 100%; height: 100%; overflow: visible; display: block; }
@keyframes bob { 0%, 100% { transform: translateY(0); } 50% { transform: translateY(-6px); } }

.mascot__tip { animation: tip 2.2s ease-in-out infinite; transform-box: fill-box; transform-origin: center; }
@keyframes tip { 0%, 100% { opacity: 1; transform: scale(1); } 50% { opacity: 0.6; transform: scale(0.78); } }

.mascot__pupil { transition: transform 0.18s ease-out; }

/* 眨眼（捂眼时不影响） */
.mascot__eye { transform-box: fill-box; transform-origin: center; animation: blink 5.4s infinite; }
.mascot__eye--r { animation-delay: 0.08s; }
@keyframes blink { 0%, 92%, 100% { transform: scaleY(1); } 96% { transform: scaleY(0.12); } }

/* 双手：默认垂在头部下方，输入密码时上抬捂住眼睛 */
.mascot__arm { transition: transform 0.42s cubic-bezier(0.34, 1.4, 0.5, 1); transform-box: fill-box; transform-origin: center; }
.mascot.is-password .mascot__arm--l { transform: translate(20px, -52px) rotate(-12deg); }
.mascot.is-password .mascot__arm--r { transform: translate(-20px, -52px) rotate(12deg); }
/* 输入密码时嘴变成抿嘴 */
.mascot__mouth { transition: d 0.3s ease; }

/* ---------- 头部文字 ---------- */
.login__head {
  text-align: center;
  margin: 6px 0 26px;
  .login__title { margin: 0 0 8px; font-size: 24px; font-weight: 700; color: $slate900; letter-spacing: -0.01em; }
  .login__sub { margin: 0; font-size: 14px; color: $slate500; min-height: 20px; transition: color 0.2s ease; }
}
.login__demo {
  margin: 0 0 18px;
  padding: 9px 14px;
  text-align: center;
  font-size: 13px;
  color: $themePrimary;
  background: rgba(99, 102, 241, 0.08);
  border: 1px solid rgba(99, 102, 241, 0.22);
  border-radius: 10px;
  b { font-weight: 600; }
}

/* ---------- 表单 ---------- */
.login__form {
  :deep(.el-form-item) { margin-bottom: 20px; }
  :deep(.el-input__wrapper) {
    border-radius: 12px;
    box-shadow: 0 0 0 1px rgba(15, 23, 42, 0.12) inset;
    transition: box-shadow 0.2s ease;
  }
  :deep(.el-input__wrapper:hover) { box-shadow: 0 0 0 1px rgba(99, 102, 241, 0.45) inset; }
  :deep(.el-input__wrapper.is-focus) { box-shadow: 0 0 0 1px $themePrimary inset, 0 0 0 4px rgba(99, 102, 241, 0.15); }
}
.login__actions {
  display: flex;
  flex-direction: column;
  align-items: stretch;
  gap: 14px;
  margin-top: 24px;
}
.login__forget {
  align-self: flex-end;
  font-size: 14px;
  color: $slate500;
  cursor: pointer;
  transition: color 0.2s ease;
  &:hover { color: $themePrimary; }
}
.login__btn {
  width: 100%;
  height: 48px;
  font-size: 16px;
  font-weight: 600;
  border: none;
  border-radius: 12px;
  background: linear-gradient(135deg, $themePrimary, $themePrimaryLight);
  box-shadow: 0 8px 20px -6px rgba(79, 70, 229, 0.5);
  transition: transform 0.2s ease, box-shadow 0.2s ease, filter 0.2s ease;
  &:hover { transform: translateY(-1px); box-shadow: 0 12px 26px -6px rgba(79, 70, 229, 0.6); filter: brightness(1.05); }
  &:active { transform: translateY(0); }
}

@media (max-width: 480px) {
  .login__card { padding: 8px 26px 30px; }
}

/* 无障碍：减少动态效果 */
@media (prefers-reduced-motion: reduce) {
  .aurora, .orb, .mascot, .mascot__tip, .mascot__eye, .login__card { animation: none !important; }
}
</style>
