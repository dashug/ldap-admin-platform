<template>
  <div>
    <page-header title="平台对接" subtitle="对接飞书 / 企微 / 钉钉，配置同步凭证" />

    <el-card class="settings-card" shadow="never">
      <el-alert
        title="建议流程：选择平台 -> 填写凭证 -> 先测试连接 -> 再保存配置。"
        type="info"
        :closable="false"
        show-icon
        class="settings-alert"
      />
      <el-tabs v-model="thirdPartyTab" v-loading="loading">
        <el-tab-pane label="钉钉" name="dingtalk">
          <el-form ref="dingtalkFormRef" size="small" :model="dingtalkForm" :rules="dingtalkRules" label-width="130px">
            <el-form-item label="平台标识" prop="flag"><el-input v-model.trim="dingtalkForm.flag" placeholder="默认 dingtalk" /></el-form-item>
            <el-form-item label="AppKey" prop="appKey"><el-input v-model.trim="dingtalkForm.appKey" /></el-form-item>
            <el-form-item label="AppSecret" prop="appSecret"><el-input v-model.trim="dingtalkForm.appSecret" show-password placeholder="留空表示不修改" /></el-form-item>
            <el-form-item label="AgentId" prop="agentId"><el-input v-model.trim="dingtalkForm.agentId" /></el-form-item>
            <el-form-item label="启用同步"><el-switch v-model="dingtalkForm.enableSync" /></el-form-item>
          </el-form>
        </el-tab-pane>
        <el-tab-pane label="企业微信" name="wecom">
          <el-form ref="wecomFormRef" size="small" :model="wecomForm" :rules="wecomRules" label-width="130px">
            <el-form-item label="平台标识" prop="flag"><el-input v-model.trim="wecomForm.flag" placeholder="默认 wecom" /></el-form-item>
            <el-form-item label="CorpId" prop="corpId"><el-input v-model.trim="wecomForm.corpId" /></el-form-item>
            <el-form-item label="CorpSecret" prop="corpSecret"><el-input v-model.trim="wecomForm.corpSecret" show-password placeholder="留空表示不修改" /></el-form-item>
            <el-form-item label="AgentId" prop="weComAgentId"><el-input-number v-model="wecomForm.weComAgentId" :min="1" style="width: 100%" /></el-form-item>
            <el-form-item label="启用同步"><el-switch v-model="wecomForm.enableSync" /></el-form-item>
          </el-form>
        </el-tab-pane>
        <el-tab-pane label="飞书" name="feishu">
          <el-form ref="feishuFormRef" size="small" :model="feishuForm" :rules="feishuRules" label-width="130px">
            <el-form-item label="平台标识" prop="flag"><el-input v-model.trim="feishuForm.flag" placeholder="默认 feishu" /></el-form-item>
            <el-form-item label="AppId" prop="appId"><el-input v-model.trim="feishuForm.appId" /></el-form-item>
            <el-form-item label="AppSecret" prop="appSecret"><el-input v-model.trim="feishuForm.appSecret" show-password placeholder="留空表示不修改" /></el-form-item>
            <el-form-item label="启用同步"><el-switch v-model="feishuForm.enableSync" /></el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
      <div class="settings-footer">
        <el-button size="small" type="warning" :loading="testing" @click="handleTest">测试连接</el-button>
        <el-button size="small" type="primary" :loading="saving" @click="handleSave">保 存</el-button>
      </div>
    </el-card>
  </div>
</template>

<script>
import { getConfig, testThirdPartyConfig, updateThirdPartyConfig } from '@/api/system/base'
import PageHeader from '@/components/PageHeader/index.vue'
import { ElMessage as Message } from 'element-plus'

export default {
  name: 'SettingsThirdParty',
  components: { PageHeader },
  data() {
    return {
      loading: false,
      testing: false,
      saving: false,
      thirdPartyTab: 'dingtalk',
      dingtalkForm: {
        platform: 'dingtalk',
        flag: 'dingtalk',
        appKey: '',
        appSecret: '',
        agentId: '',
        enableSync: false
      },
      wecomForm: {
        platform: 'wecom',
        flag: 'wecom',
        corpId: '',
        corpSecret: '',
        weComAgentId: 1,
        enableSync: false
      },
      feishuForm: {
        platform: 'feishu',
        flag: 'feishu',
        appId: '',
        appSecret: '',
        enableSync: false
      },
      dingtalkRules: {
        appKey: [{ required: true, message: '请输入 AppKey', trigger: 'blur' }]
      },
      wecomRules: {
        corpId: [{ required: true, message: '请输入 CorpId', trigger: 'blur' }],
        weComAgentId: [{ required: true, message: '请输入 AgentId', trigger: 'change' }]
      },
      feishuRules: {
        appId: [{ required: true, message: '请输入 AppId', trigger: 'blur' }]
      },
      // 离开页面时的脏检查基线（每次拉取/保存后刷新，覆盖三个平台表单）
      savedSnapshot: ''
    }
  },
  created() {
    this.fetchConfig()
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
      return JSON.stringify({ dingtalk: this.dingtalkForm, wecom: this.wecomForm, feishu: this.feishuForm })
    },
    isDirty() {
      return this.savedSnapshot !== '' && this.snapshot() !== this.savedSnapshot
    },
    async fetchConfig() {
      this.loading = true
      try {
        const { data } = await getConfig()
        this.dingtalkForm = {
          platform: 'dingtalk',
          flag: data.dingTalkFlag || 'dingtalk',
          appKey: data.dingTalkAppKey || '',
          appSecret: '',
          agentId: data.dingTalkAgentId || '',
          enableSync: !!data.dingTalkEnableSync
        }
        this.wecomForm = {
          platform: 'wecom',
          flag: data.weComFlag || 'wecom',
          corpId: data.weComCorpId || '',
          corpSecret: '',
          weComAgentId: data.weComAgentId || 1,
          enableSync: !!data.weComEnableSync
        }
        this.feishuForm = {
          platform: 'feishu',
          flag: data.feiShuFlag || 'feishu',
          appId: data.feiShuAppId || '',
          appSecret: '',
          enableSync: !!data.feiShuEnableSync
        }
      } catch (error) {
        Message.error('获取平台配置失败')
      } finally {
        this.loading = false
        this.savedSnapshot = this.snapshot()
      }
    },
    getCurrentForm() {
      if (this.thirdPartyTab === 'wecom') {
        return this.wecomForm
      }
      if (this.thirdPartyTab === 'feishu') {
        return this.feishuForm
      }
      return this.dingtalkForm
    },
    getCurrentRefName() {
      if (this.thirdPartyTab === 'wecom') {
        return 'wecomFormRef'
      }
      if (this.thirdPartyTab === 'feishu') {
        return 'feishuFormRef'
      }
      return 'dingtalkFormRef'
    },
    handleTest() {
      const refName = this.getCurrentRefName()
      const form = this.getCurrentForm()
      this.$refs[refName].validate(async valid => {
        if (!valid) return
        this.testing = true
        try {
          await testThirdPartyConfig(form)
          Message.success('连接测试成功')
        } finally {
          this.testing = false
        }
      })
    },
    handleSave() {
      const refName = this.getCurrentRefName()
      const form = this.getCurrentForm()
      this.$refs[refName].validate(async valid => {
        if (!valid) return
        this.saving = true
        try {
          await updateThirdPartyConfig(form)
          Message.success('平台配置已保存')
          this.fetchConfig()
        } finally {
          this.saving = false
        }
      })
    }
  }
}
</script>

<style scoped>
.settings-card {
  margin: 10px;
  margin-bottom: 100px;
}
.settings-alert {
  margin-bottom: 16px;
}
.settings-footer {
  text-align: right;
  margin-top: 8px;
}
.settings-footer .el-button {
  margin-left: 10px;
}
</style>
