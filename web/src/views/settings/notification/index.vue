<template>
  <div>
    <page-header title="通知设置" subtitle="配置 Webhook 回调与账户通知邮件（SMTP）" />

    <el-card class="settings-card" shadow="never">
      <el-form ref="formRef" v-loading="loading" size="small" :model="form" :rules="rules" label-width="180px" class="notification-form">
        <el-form-item label="Webhook 回调地址">
          <el-input v-model.trim="form.webhookUrl" placeholder="用户/部门创建或同步后 POST 回调的 URL，留空则不发送" clearable />
          <div class="form-tip">填写后，在用户或部门创建、更新、同步完成时会向该地址发送 POST 请求（失败自动重试 3 次）</div>
        </el-form-item>
        <el-form-item label="Webhook 签名密钥">
          <el-input v-model.trim="form.webhookSecret" type="password" show-password clearable :placeholder="webhookSecretSet ? '已配置，留空表示不修改' : '可选：填写后回调将带 HMAC-SHA256 签名头'" />
          <div class="form-tip">接收端用同一密钥验签：sha256=hex(hmac_sha256(secret, 时间戳 + '.' + 原始body))，请求头 X-Webhook-Signature / X-Webhook-Timestamp。</div>
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
      <div class="settings-footer">
        <el-button size="small" :loading="testingWebhook" @click="testWebhook">测试 Webhook</el-button>
        <el-button size="small" :loading="testingEmail" @click="testEmail">发送测试邮件</el-button>
        <el-button size="small" type="primary" :loading="saving" @click="submit">保 存</el-button>
      </div>
    </el-card>

    <el-card class="settings-card webhook-log-card" shadow="never">
      <div class="webhook-log-head">
        <span class="webhook-log-title">Webhook 投递记录</span>
        <el-button size="small" icon="Refresh" :loading="deliveriesLoading" @click="fetchDeliveries">刷新</el-button>
      </div>
      <el-table :data="deliveries" v-loading="deliveriesLoading" size="small" border style="width: 100%">
        <el-table-column label="时间" width="170">
          <template #default="s">{{ formatTime(s.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column prop="event" label="事件" width="120" show-overflow-tooltip />
        <el-table-column label="结果" width="80" align="center">
          <template #default="s">
            <el-tag :type="s.row.success ? 'success' : 'danger'" size="small">{{ s.row.success ? '成功' : '失败' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="statusCode" label="状态码" width="80" align="center" />
        <el-table-column prop="attempts" label="尝试" width="64" align="center" />
        <el-table-column prop="url" label="地址" min-width="180" show-overflow-tooltip />
        <el-table-column prop="error" label="错误" min-width="160" show-overflow-tooltip />
        <template #empty>暂无投递记录</template>
      </el-table>
      <div class="pagination-wrap">
        <el-pagination
          :current-page="deliveryParams.pageNum"
          :page-size="deliveryParams.pageSize"
          :total="deliveriesTotal"
          layout="total, prev, pager, next"
          background
          @current-change="handleDeliveryPage"
        />
      </div>
    </el-card>
  </div>
</template>

<script>
import { getConfig, updateEmailConfig, testNotification, getWebhookDeliveries } from '@/api/system/base'
import PageHeader from '@/components/PageHeader/index.vue'
import { ElMessage as Message } from 'element-plus'

export default {
  name: 'SettingsNotification',
  components: { PageHeader },
  data() {
    return {
      loading: false,
      saving: false,
      testingEmail: false,
      testingWebhook: false,
      webhookSecretSet: false,
      // Webhook 投递记录
      deliveries: [],
      deliveriesTotal: 0,
      deliveriesLoading: false,
      deliveryParams: { pageNum: 1, pageSize: 10 },
      form: {
        webhookUrl: '',
        webhookSecret: '',
        sendUserCreationMail: false,
        smtpHost: '',
        smtpPort: '',
        smtpUser: '',
        smtpPass: '',
        smtpFrom: ''
      },
      // SMTP 字段仅在开启邮件通知时渲染，未渲染的项不会被 validate 校验，
      // 因此这里无条件声明必填即可（关闭开关时自动跳过）。smtpPass 留空表示不修改，不强制。
      rules: {
        smtpHost: [{ required: true, message: '请输入 SMTP 服务器', trigger: 'blur' }],
        smtpPort: [{ required: true, message: '请输入 SMTP 端口', trigger: 'blur' }],
        smtpUser: [{ required: true, message: '请输入发件人邮箱', trigger: 'blur' }]
      },
      // 离开页面时的脏检查基线（每次拉取/保存后刷新）
      savedSnapshot: ''
    }
  },
  created() {
    this.fetchConfig()
    this.fetchDeliveries()
  },
  beforeRouteLeave(to, from, next) {
    if (!this.isDirty()) return next()
    this.$confirm('有未保存的修改，确定离开吗？', '提示', {
      confirmButtonText: '离开',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(() => next()).catch(() => next(false))
  },
  methods: {
    snapshot() {
      return JSON.stringify(this.form)
    },
    isDirty() {
      return this.savedSnapshot !== '' && this.snapshot() !== this.savedSnapshot
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
        this.form.webhookSecret = ''
        this.webhookSecretSet = !!data.webhookSecretSet
      } catch (e) {
        Message.error('获取配置失败')
      } finally {
        this.loading = false
        this.savedSnapshot = this.snapshot()
      }
    },
    async fetchDeliveries() {
      this.deliveriesLoading = true
      try {
        const { data } = await getWebhookDeliveries(this.deliveryParams)
        this.deliveries = (data && data.list) || []
        this.deliveriesTotal = (data && data.total) || 0
      } catch (e) {
        // 错误由请求拦截器统一提示
      } finally {
        this.deliveriesLoading = false
      }
    },
    handleDeliveryPage(page) {
      this.deliveryParams.pageNum = page
      this.fetchDeliveries()
    },
    formatTime(v) {
      if (!v) return '-'
      return String(v).replace('T', ' ').slice(0, 19)
    },
    submit() {
      this.$refs.formRef.validate(async valid => {
        if (!valid) return
        this.saving = true
        try {
          await updateEmailConfig(this.form)
          Message.success('通知设置已保存')
          this.fetchConfig()
        } catch (e) {
          Message.error(e?.message || '保存失败')
        } finally {
          this.saving = false
        }
      })
    },
    // 发送测试邮件：弹窗输入收件人，使用当前表单的 SMTP 值（留空字段后端回落到已保存配置）
    testEmail() {
      this.$prompt('请输入测试邮件收件人邮箱', '发送测试邮件', {
        confirmButtonText: '发送',
        cancelButtonText: '取消',
        inputPattern: /^[^@\s]+@[^@\s]+\.[^@\s]+$/,
        inputErrorMessage: '邮箱格式不正确'
      }).then(async ({ value }) => {
        this.testingEmail = true
        try {
          await testNotification({
            target: 'email',
            mail: value,
            smtpHost: this.form.smtpHost,
            smtpPort: this.form.smtpPort,
            smtpUser: this.form.smtpUser,
            smtpPass: this.form.smtpPass,
            smtpFrom: this.form.smtpFrom
          })
          Message.success('测试邮件已发送，请查收')
        } catch (e) {
          // 失败信息由请求拦截器统一弹出
        } finally {
          this.testingEmail = false
        }
      }).catch(() => {})
    },
    // 测试 Webhook：向当前填写的回调地址发送一条测试回调
    async testWebhook() {
      if (!this.form.webhookUrl) {
        Message.warning('请先填写 Webhook 回调地址')
        return
      }
      this.testingWebhook = true
      try {
        await testNotification({ target: 'webhook', webhookUrl: this.form.webhookUrl })
        Message.success('测试回调已发送，请在接收端确认')
      } catch (e) {
        // 失败信息由请求拦截器统一弹出
      } finally {
        this.testingWebhook = false
        this.fetchDeliveries()
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
.notification-form :deep(.el-form-item__content) {
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
  line-height: 1.4;
}
.settings-footer {
  text-align: right;
  margin-top: 8px;
}
.settings-footer .el-button {
  margin-left: 10px;
}
.webhook-log-card {
  margin-top: 0;
}
.webhook-log-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}
.webhook-log-title {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
}
.pagination-wrap {
  display: flex;
  justify-content: flex-end;
  margin-top: 12px;
}
</style>
