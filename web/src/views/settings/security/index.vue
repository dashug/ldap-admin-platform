<template>
  <div>
    <page-header title="登录安全" subtitle="为当前管理员账号开启二次验证（TOTP）" />

    <el-card class="settings-card" shadow="never" v-loading="loading">
      <div class="mfa-status">
        <el-icon class="mfa-status__icon" :class="enabled ? 'is-on' : 'is-off'">
          <component :is="enabled ? 'Lock' : 'Unlock'" />
        </el-icon>
        <div class="mfa-status__text">
          <div class="mfa-status__title">二次验证（MFA / TOTP）{{ enabled ? '已开启' : '未开启' }}</div>
          <div class="mfa-status__desc">
            开启后，登录时除密码外还需输入 Authenticator 应用（Google / 微软 / 1Password 等）生成的 6 位动态验证码。
          </div>
        </div>
        <el-button v-if="!enabled && !setup" type="primary" size="small" :loading="setupLoading" @click="startSetup">开启二次验证</el-button>
        <el-button v-if="enabled" type="danger" size="small" plain @click="openDisable">关闭二次验证</el-button>
      </div>

      <!-- 开启流程：扫码 + 输入验证码确认 -->
      <div v-if="setup && !enabled" class="mfa-setup">
        <el-divider content-position="left">绑定 Authenticator</el-divider>
        <div class="mfa-setup__body">
          <div class="mfa-qr">
            <img v-if="setup.qr" :src="setup.qr" alt="MFA 二维码" class="mfa-qr__img">
          </div>
          <div class="mfa-setup__steps">
            <p>1. 用 Authenticator 应用扫描左侧二维码；无法扫码时手动输入密钥：</p>
            <el-input :model-value="setup.secret" readonly size="small" class="mfa-secret">
              <template #append>
                <el-button icon="DocumentCopy" @click="copySecret">复制</el-button>
              </template>
            </el-input>
            <p style="margin-top: 14px;">2. 输入应用生成的 6 位验证码以完成绑定：</p>
            <div class="mfa-verify-row">
              <el-input v-model.trim="code" maxlength="6" placeholder="6 位验证码" style="width: 160px;" @keyup.enter="confirmEnable" />
              <el-button type="primary" size="default" :loading="verifyLoading" @click="confirmEnable">确认开启</el-button>
              <el-button size="default" @click="cancelSetup">取消</el-button>
            </div>
          </div>
        </div>
      </div>
    </el-card>

    <!-- 关闭确认：需输入验证码 -->
    <el-dialog title="关闭二次验证" v-model="disableVisible" width="380px">
      <p class="disable-tip">为确认是本人操作，请输入当前的 6 位动态验证码：</p>
      <el-input v-model.trim="disableCode" maxlength="6" placeholder="6 位验证码" @keyup.enter="confirmDisable" />
      <template #footer><div class="dialog-footer">
        <el-button size="small" @click="disableVisible = false">取 消</el-button>
        <el-button size="small" type="danger" :loading="disableLoading" @click="confirmDisable">确认关闭</el-button>
      </div></template>
    </el-dialog>
  </div>
</template>

<script>
import { getMfaStatus, mfaSetup, mfaVerify, mfaDisable } from '@/api/system/base'
import PageHeader from '@/components/PageHeader/index.vue'
import { ElMessage as Message } from 'element-plus'

export default {
  name: 'SettingsSecurity',
  components: { PageHeader },
  data() {
    return {
      loading: false,
      enabled: false,
      setup: null, // { secret, qr, otpauthUrl }
      setupLoading: false,
      code: '',
      verifyLoading: false,
      disableVisible: false,
      disableCode: '',
      disableLoading: false
    }
  },
  created() {
    this.fetchStatus()
  },
  methods: {
    async fetchStatus() {
      this.loading = true
      try {
        const { data } = await getMfaStatus()
        this.enabled = !!(data && data.enabled)
      } catch (e) {
        // 错误由请求拦截器统一提示
      } finally {
        this.loading = false
      }
    },
    async startSetup() {
      this.setupLoading = true
      try {
        const { data } = await mfaSetup()
        this.setup = data || null
        this.code = ''
      } catch (e) {
        // 错误由请求拦截器统一提示
      } finally {
        this.setupLoading = false
      }
    },
    cancelSetup() {
      this.setup = null
      this.code = ''
    },
    async confirmEnable() {
      if (this.code.length !== 6) {
        Message.warning('请输入 6 位验证码')
        return
      }
      this.verifyLoading = true
      try {
        await mfaVerify({ code: this.code })
        Message.success('二次验证已开启')
        this.enabled = true
        this.setup = null
        this.code = ''
      } catch (e) {
        // 错误由请求拦截器统一提示
      } finally {
        this.verifyLoading = false
      }
    },
    copySecret() {
      const ta = document.createElement('textarea')
      ta.value = this.setup ? this.setup.secret : ''
      document.body.appendChild(ta)
      ta.select()
      try {
        document.execCommand('copy')
        Message.success('密钥已复制')
      } catch (e) {
        Message.error('复制失败，请手动复制')
      }
      document.body.removeChild(ta)
    },
    openDisable() {
      this.disableCode = ''
      this.disableVisible = true
    },
    async confirmDisable() {
      if (this.disableCode.length !== 6) {
        Message.warning('请输入 6 位验证码')
        return
      }
      this.disableLoading = true
      try {
        await mfaDisable({ code: this.disableCode })
        Message.success('二次验证已关闭')
        this.enabled = false
        this.disableVisible = false
      } catch (e) {
        // 错误由请求拦截器统一提示
      } finally {
        this.disableLoading = false
      }
    }
  }
}
</script>

<style scoped>
.settings-card {
  margin: 10px;
  margin-bottom: 100px;
}
.mfa-status {
  display: flex;
  align-items: center;
  gap: 14px;
}
.mfa-status__icon {
  font-size: 30px;
}
.mfa-status__icon.is-on { color: #67c23a; }
.mfa-status__icon.is-off { color: #909399; }
.mfa-status__text { flex: 1; }
.mfa-status__title { font-size: 15px; font-weight: 600; color: #303133; }
.mfa-status__desc { font-size: 12px; color: #909399; margin-top: 4px; line-height: 1.5; max-width: 560px; }
.mfa-setup__body { display: flex; gap: 24px; align-items: flex-start; }
.mfa-qr__img { width: 200px; height: 200px; border: 1px solid #ebeef5; border-radius: 6px; }
.mfa-setup__steps { flex: 1; font-size: 13px; color: #606266; line-height: 1.6; }
.mfa-secret { max-width: 360px; }
.mfa-verify-row { display: flex; align-items: center; gap: 10px; }
.disable-tip { margin-bottom: 12px; color: #606266; }
.dialog-footer { text-align: right; }
.dialog-footer .el-button { margin-left: 10px; }
</style>
