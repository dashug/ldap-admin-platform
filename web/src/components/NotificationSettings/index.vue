<template>
  <el-dialog
    title="通知设置"
    :visible.sync="visible"
    width="520px"
    :close-on-click-modal="false"
    @open="fetchConfig"
  >
    <el-form v-loading="loading" size="small" label-width="180px" class="notification-form">
      <el-form-item label="Webhook 回调地址">
        <el-input v-model.trim="form.webhookUrl" placeholder="用户/部门创建或同步后 POST 回调的 URL，留空则不发送" clearable />
        <div class="form-tip">填写后，在用户或部门创建、更新、同步完成时会向该地址发送 POST 请求</div>
      </el-form-item>
      <el-form-item label="新建/同步用户时发送通知邮件">
        <div class="switch-row">
          <el-switch v-model="form.sendUserCreationMail" />
          <span class="form-tip">开启后会在创建或同步用户时向用户邮箱发送账户通知</span>
        </div>
      </el-form-item>
      <template v-if="form.sendUserCreationMail">
        <el-divider content-position="left">邮件服务器（SMTP）</el-divider>
        <el-form-item label="SMTP 服务器" prop="smtpHost">
          <el-input v-model.trim="form.smtpHost" placeholder="如 smtp.163.com" />
        </el-form-item>
        <el-form-item label="SMTP 端口" prop="smtpPort">
          <el-input v-model.trim="form.smtpPort" placeholder="如 465（SSL）" />
        </el-form-item>
        <el-form-item label="发件人邮箱" prop="smtpUser">
          <el-input v-model.trim="form.smtpUser" placeholder="登录邮箱账号" />
        </el-form-item>
        <el-form-item label="邮箱密码/授权码" prop="smtpPass">
          <el-input v-model.trim="form.smtpPass" type="password" show-password placeholder="留空表示不修改" />
        </el-form-item>
        <el-form-item label="发件人显示名称" prop="smtpFrom">
          <el-input v-model.trim="form.smtpFrom" placeholder="如：LDAP 管理平台" />
        </el-form-item>
      </template>
    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button size="small" @click="visible = false">取 消</el-button>
      <el-button size="small" type="primary" :loading="saving" @click="submit">保 存</el-button>
    </div>
  </el-dialog>
</template>

<script>
import { getConfig, updateEmailConfig } from '@/api/system/base'

export default {
  name: 'NotificationSettings',
  data() {
    return {
      visible: false,
      loading: false,
      saving: false,
      form: {
        webhookUrl: '',
        sendUserCreationMail: false,
        smtpHost: '',
        smtpPort: '',
        smtpUser: '',
        smtpPass: '',
        smtpFrom: ''
      }
    }
  },
  methods: {
    open() {
      this.visible = true
    },
    async fetchConfig() {
      this.loading = true
      try {
        const { data } = await getConfig()
        this.form.webhookUrl = data.webhookUrl || ''
        this.form.sendUserCreationMail = !!data.sendUserCreationMail
        this.form.smtpHost = data.smtpHost || ''
        this.form.smtpPort = data.smtpPort || ''
        this.form.smtpUser = data.smtpUser || ''
        this.form.smtpFrom = data.smtpFrom || ''
        this.form.smtpPass = ''
      } catch (e) {
        this.$message.error('获取配置失败')
      } finally {
        this.loading = false
      }
    },
    async submit() {
      this.saving = true
      try {
        await updateEmailConfig(this.form)
        this.$message.success('通知设置已保存')
        this.visible = false
        this.$emit('saved')
      } catch (e) {
        this.$message.error(e?.message || '保存失败')
      } finally {
        this.saving = false
      }
    }
  }
}
</script>

<style scoped>
.notification-form >>> .el-form-item__content {
  line-height: 1.5;
}
.switch-row {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
}
.form-tip {
  color: #909399;
  font-size: 12px;
  max-width: 280px;
  line-height: 1.4;
}
.dialog-footer {
  text-align: right;
}
.dialog-footer .el-button {
  margin-left: 10px;
}
</style>
