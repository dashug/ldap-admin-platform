<template>
  <div>
    <page-header title="目录配置" subtitle="配置 OpenLDAP / AD 连接、DN 与同步规则">
      <template #actions>
        <el-button size="small" @click="exportConfig">导出配置</el-button>
        <el-button size="small" @click="openConfigImport">导入配置</el-button>
        <input ref="configImportInput" type="file" accept=".json,application/json" style="display: none" @change="onConfigImportFile">
      </template>
    </page-header>

    <el-card class="settings-card" shadow="never">
      <!-- 当前（已保存配置的）LDAP 连接状态 -->
      <div class="ldap-status" :class="ldapStatus.connected ? 'is-ok' : 'is-bad'">
        <el-icon class="ldap-status__dot">
          <component :is="ldapStatus.connected ? 'CircleCheckFilled' : 'CircleCloseFilled'" />
        </el-icon>
        <span class="ldap-status__text">
          {{ ldapStatusLoading ? '正在检测当前连接…' : (ldapStatus.connected ? 'LDAP 连接正常' : 'LDAP 未连接') }}
          <template v-if="!ldapStatusLoading && ldapStatus.message">（{{ ldapStatus.message }}）</template>
        </span>
        <el-button class="ldap-status__refresh" link type="primary" size="small" icon="Refresh" :loading="ldapStatusLoading" @click="fetchLDAPStatus">刷新</el-button>
      </div>

      <el-alert
        title="先选目录类型，再填地址和 DN；管理员密码留空表示不修改。「测试连接」按当前填写的地址/DN 试连，不会保存。"
        type="info"
        :closable="false"
        show-icon
        class="settings-alert"
      />
      <el-form ref="directoryFormRef" v-loading="loading" size="small" :model="directoryForm" :rules="directoryRules" label-width="130px" class="directory-form">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="目录类型" prop="directoryType">
              <el-select v-model="directoryForm.directoryType" style="width: 100%">
                <el-option label="OpenLDAP" value="openldap" />
                <el-option label="Windows AD" value="ad" />
              </el-select>
            </el-form-item>
            <el-form-item label="LDAP 地址" prop="url">
              <el-input v-model.trim="directoryForm.url" placeholder="ldap://10.0.0.10:389" />
            </el-form-item>
            <el-form-item label="Base DN" prop="baseDN">
              <el-input v-model.trim="directoryForm.baseDN" placeholder="dc=example,dc=com" />
            </el-form-item>
            <el-form-item label="管理员 DN" prop="adminDN">
              <el-input v-model.trim="directoryForm.adminDN" placeholder="cn=admin,dc=example,dc=com" />
            </el-form-item>
            <el-form-item label="管理员密码" prop="adminPass">
              <el-input v-model.trim="directoryForm.adminPass" show-password placeholder="留空不修改" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="用户 OU DN" prop="userDN">
              <el-input v-model.trim="directoryForm.userDN" placeholder="ou=people,dc=example,dc=com" />
            </el-form-item>
            <el-form-item label="默认初始密码" prop="userInitPassword">
              <el-input v-model.trim="directoryForm.userInitPassword" show-password placeholder="新建用户默认密码" />
            </el-form-item>
            <el-form-item label="默认邮箱后缀" prop="defaultEmailSuffix">
              <el-input v-model.trim="directoryForm.defaultEmailSuffix" placeholder="example.com" />
            </el-form-item>
            <el-form-item label="启用 LDAP 同步">
              <el-switch v-model="directoryForm.ldapEnableSync" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-divider content-position="left">同步与 DN 规则</el-divider>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="用户名规则" prop="syncUsernameRule">
              <el-select v-model="directoryForm.syncUsernameRule" style="width: 100%" placeholder="用户名的生成方式">
                <el-option label="邮箱前段（@ 前部分）" value="email_prefix" />
                <el-option label="姓名拼音" value="pinyin" />
                <el-option label="工号" value="job_number" />
                <el-option label="按字段关联配置" value="field_relation" />
              </el-select>
            </el-form-item>
            <el-form-item label="部门名规则" prop="syncGroupNameRule">
              <el-select v-model="directoryForm.syncGroupNameRule" style="width: 100%" placeholder="部门名来源">
                <el-option label="中文名" value="name" />
                <el-option label="拼音" value="pinyin" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="用户 DN 的 RDN" prop="userRDNAttr">
              <el-select v-model="directoryForm.userRDNAttr" style="width: 100%" placeholder="uid 或 cn">
                <el-option label="uid（OpenLDAP）" value="uid" />
                <el-option label="cn（AD）" value="cn" />
              </el-select>
            </el-form-item>
            <el-form-item label="部门 DN 的 RDN" prop="groupRDNAttr">
              <el-select v-model="directoryForm.groupRDNAttr" style="width: 100%" placeholder="cn 或 ou">
                <el-option label="cn" value="cn" />
                <el-option label="ou" value="ou" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <div class="settings-footer">
        <el-button size="small" :loading="testing" @click="handleTestConnection">测试连接</el-button>
        <el-button size="small" type="primary" :loading="saving" @click="submitDirectoryConfig">保 存</el-button>
      </div>
    </el-card>
  </div>
</template>

<script>
import { getConfig, updateDirectoryConfig, importConfig, testDirectoryConfig, getLDAPStatus } from '@/api/system/base'
import PageHeader from '@/components/PageHeader/index.vue'
import { ElMessage as Message } from 'element-plus'

export default {
  name: 'SettingsDirectory',
  components: { PageHeader },
  data() {
    return {
      loading: false,
      saving: false,
      testing: false,
      // 当前已保存配置的连接状态（由 getLDAPStatus 探测）
      ldapStatus: { connected: false, message: '' },
      ldapStatusLoading: false,
      directoryForm: {
        directoryType: 'openldap',
        url: '',
        baseDN: '',
        adminDN: '',
        adminPass: '',
        userDN: '',
        userInitPassword: '',
        defaultEmailSuffix: '',
        ldapEnableSync: false,
        syncUsernameRule: 'email_prefix',
        syncGroupNameRule: 'name',
        userRDNAttr: 'uid',
        groupRDNAttr: 'cn'
      },
      directoryRules: {
        directoryType: [
          { required: true, message: '请选择目录类型', trigger: 'change' }
        ],
        url: [
          { required: true, message: '请输入 LDAP 地址', trigger: 'blur' }
        ],
        baseDN: [
          { required: true, message: '请输入 Base DN', trigger: 'blur' }
        ],
        adminDN: [
          { required: true, message: '请输入管理员 DN', trigger: 'blur' }
        ],
        userDN: [
          { required: true, message: '请输入用户 OU DN', trigger: 'blur' }
        ],
        userInitPassword: [
          { required: true, message: '请输入默认初始密码', trigger: 'blur' }
        ],
        defaultEmailSuffix: [
          { required: true, message: '请输入默认邮箱后缀', trigger: 'blur' }
        ]
      },
      // 离开页面时的脏检查基线（每次拉取/保存后刷新）
      savedSnapshot: ''
    }
  },
  created() {
    this.fetchConfig()
    this.fetchLDAPStatus()
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
      return JSON.stringify(this.directoryForm)
    },
    isDirty() {
      return this.savedSnapshot !== '' && this.snapshot() !== this.savedSnapshot
    },
    async fetchLDAPStatus() {
      this.ldapStatusLoading = true
      try {
        const { data } = await getLDAPStatus()
        this.ldapStatus = { connected: !!(data && data.connected), message: (data && data.message) || '' }
      } catch (e) {
        this.ldapStatus = { connected: false, message: '获取状态失败' }
      } finally {
        this.ldapStatusLoading = false
      }
    },
    async handleTestConnection() {
      this.testing = true
      try {
        const res = await testDirectoryConfig({
          url: this.directoryForm.url,
          adminDN: this.directoryForm.adminDN,
          adminPass: this.directoryForm.adminPass
        })
        const msg = res && res.data && res.data.message
        Message.success(msg ? ('连接成功：' + msg) : '连接成功')
      } catch (e) {
        // 失败信息由请求拦截器统一弹出
      } finally {
        this.testing = false
      }
    },
    async fetchConfig() {
      this.loading = true
      try {
        const { data } = await getConfig()
        const dirType = (data.directoryType || 'openldap').toLowerCase()
        this.directoryForm = {
          directoryType: dirType,
          url: data.url || '',
          baseDN: data.baseDN || '',
          adminDN: data.adminDN || '',
          adminPass: '',
          userDN: data.userDN || '',
          userInitPassword: data.userInitPassword || '',
          defaultEmailSuffix: data.defaultEmailSuffix || '',
          ldapEnableSync: !!data.ldapEnableSync,
          syncUsernameRule: data.syncUsernameRule || 'email_prefix',
          syncGroupNameRule: data.syncGroupNameRule || 'name',
          userRDNAttr: data.userRDNAttr || (dirType === 'ad' ? 'cn' : 'uid'),
          groupRDNAttr: data.groupRDNAttr || 'cn'
        }
      } catch (error) {
        Message.error('获取目录配置失败')
      } finally {
        this.loading = false
        this.savedSnapshot = this.snapshot()
      }
    },
    submitDirectoryConfig() {
      this.$refs.directoryFormRef.validate(async valid => {
        if (!valid) {
          return
        }
        this.saving = true
        try {
          await updateDirectoryConfig(this.directoryForm)
          Message.success('目录配置已保存')
          this.fetchConfig()
          this.fetchLDAPStatus()
        } finally {
          this.saving = false
        }
      })
    },
    async exportConfig() {
      try {
        const res = await getConfig()
        const data = res.data || {}
        const exportData = {
          directoryType: data.directoryType || 'openldap',
          url: data.url || '',
          baseDN: data.baseDN || '',
          adminDN: data.adminDN || '',
          adminPass: '',
          userDN: data.userDN || '',
          userInitPassword: data.userInitPassword || '',
          defaultEmailSuffix: data.defaultEmailSuffix || '',
          ldapEnableSync: !!data.ldapEnableSync,
          syncUsernameRule: data.syncUsernameRule || 'email_prefix',
          syncGroupNameRule: data.syncGroupNameRule || 'name',
          userRDNAttr: data.userRDNAttr || 'uid',
          groupRDNAttr: data.groupRDNAttr || 'cn'
        }
        const blob = new Blob([JSON.stringify(exportData, null, 2)], { type: 'application/json' })
        const a = document.createElement('a')
        a.href = URL.createObjectURL(blob)
        a.download = `go-ldap-admin-config-${new Date().toISOString().slice(0, 10)}.json`
        a.click()
        URL.revokeObjectURL(a.href)
        Message.success('配置已导出')
      } catch (e) {
        Message.error('导出失败')
      }
    },
    openConfigImport() {
      this.$refs.configImportInput && this.$refs.configImportInput.click()
    },
    async onConfigImportFile(e) {
      const file = e.target && e.target.files[0]
      if (!file) return
      e.target.value = ''
      try {
        await this.$confirm(`导入将用文件「${file.name}」中的内容覆盖现有目录与同步规则配置，确定继续吗？`, '导入配置', {
          confirmButtonText: '导入并覆盖',
          cancelButtonText: '取消',
          type: 'warning'
        })
      } catch (cancel) {
        // 用户取消导入
        return
      }
      try {
        const text = await new Promise((resolve, reject) => {
          const r = new FileReader()
          r.onload = () => resolve(r.result)
          r.onerror = reject
          r.readAsText(file)
        })
        const data = JSON.parse(text)
        await importConfig(data)
        Message.success('配置已导入')
        this.fetchConfig()
        this.fetchLDAPStatus()
      } catch (err) {
        Message.error('导入失败：' + (err.msg || err.message || '无效的 JSON'))
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
.settings-alert {
  margin-bottom: 16px;
}
.ldap-status {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  margin-bottom: 12px;
  border-radius: 6px;
  font-size: 13px;
  border: 1px solid transparent;
}
.ldap-status.is-ok {
  color: #67c23a;
  background: #f0f9eb;
  border-color: #e1f3d8;
}
.ldap-status.is-bad {
  color: #f56c6c;
  background: #fef0f0;
  border-color: #fde2e2;
}
.ldap-status__dot {
  font-size: 16px;
}
.ldap-status__text {
  flex: 1;
}
.directory-form .el-divider {
  margin: 20px 0 16px;
}
.settings-footer {
  text-align: right;
  margin-top: 8px;
}
.settings-footer .el-button {
  margin-left: 10px;
}
</style>
