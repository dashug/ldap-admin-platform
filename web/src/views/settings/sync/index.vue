<template>
  <div>
    <page-header title="定时同步" subtitle="按计划自动从 AD / 钉钉 / 企微 / 飞书同步用户与部门" />

    <el-card class="settings-card" shadow="never">
      <el-alert
        type="info"
        :closable="false"
        show-icon
        class="settings-alert"
        title="开启后将按 cron 周期自动同步所有「已启用同步」的来源（部门 + 用户）。各来源的启用开关在「平台对接 / 目录配置」中设置。"
      />
      <el-form v-loading="loading" size="small" label-width="130px">
        <el-form-item label="启用定时同步">
          <el-switch v-model="form.autoSyncEnabled" />
        </el-form-item>
        <el-form-item label="同步周期 (cron)">
          <el-input v-model.trim="form.autoSyncCron" placeholder="6 段含秒，如 0 0 */6 * * *（每 6 小时）" style="max-width: 320px;" :disabled="!form.autoSyncEnabled" />
          <div class="form-tip">
            常用：<el-link type="primary" :underline="false" @click="form.autoSyncCron = '0 0 */6 * * *'">每 6 小时</el-link> ·
            <el-link type="primary" :underline="false" @click="form.autoSyncCron = '0 0 2 * * *'">每天 02:00</el-link> ·
            <el-link type="primary" :underline="false" @click="form.autoSyncCron = '0 */30 * * * *'">每 30 分钟</el-link>
          </div>
        </el-form-item>
        <el-form-item label="当前已启用来源">
          <template v-if="enabledSources.length">
            <el-tag v-for="s in enabledSources" :key="s.key" type="success" effect="plain" class="src-tag">{{ s.label }}</el-tag>
          </template>
          <span v-else class="form-tip">暂无已启用同步的来源，请先在「目录配置 / 平台对接」中开启。</span>
        </el-form-item>
      </el-form>
      <div class="settings-footer">
        <el-button size="small" type="primary" :loading="saving" @click="save">保 存</el-button>
      </div>
    </el-card>

    <el-card class="settings-card" shadow="never">
      <div class="sync-run-head">
        <span class="sync-run-title">立即同步</span>
      </div>
      <div class="run-now-bar">
        <el-button v-for="s in allSources" :key="s.key" size="small" :disabled="!s.enabled" @click="runNow(s)">
          立即同步{{ s.label }}
        </el-button>
        <span class="form-tip">同步在后台执行，稍后点「刷新」查看结果。</span>
      </div>
    </el-card>

    <el-card class="settings-card" shadow="never">
      <div class="sync-run-head">
        <span class="sync-run-title">同步运行记录</span>
        <el-button size="small" icon="Refresh" :loading="runsLoading" @click="fetchRuns">刷新</el-button>
      </div>
      <el-table :data="runs" v-loading="runsLoading" size="small" border style="width: 100%">
        <el-table-column label="时间" width="170">
          <template #default="s">{{ formatTime(s.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column label="来源" width="110">
          <template #default="s">{{ sourceLabel(s.row.source) }}</template>
        </el-table-column>
        <el-table-column label="触发" width="80" align="center">
          <template #default="s">{{ s.row.trigger === 'auto' ? '定时' : '手动' }}</template>
        </el-table-column>
        <el-table-column label="结果" width="80" align="center">
          <template #default="s">
            <el-tag :type="s.row.success ? 'success' : 'danger'" size="small">{{ s.row.success ? '成功' : '失败' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="耗时" width="90" align="center">
          <template #default="s">{{ (s.row.duration / 1000).toFixed(1) }}s</template>
        </el-table-column>
        <el-table-column prop="message" label="信息" min-width="220" show-overflow-tooltip />
        <template #empty>暂无同步记录</template>
      </el-table>
      <div class="pagination-wrap">
        <el-pagination
          :current-page="runParams.pageNum"
          :page-size="runParams.pageSize"
          :total="runsTotal"
          layout="total, prev, pager, next"
          background
          @current-change="handleRunPage"
        />
      </div>
    </el-card>
  </div>
</template>

<script>
import { getConfig, updateSyncConfig, runSyncNow, getSyncRuns } from '@/api/system/base'
import PageHeader from '@/components/PageHeader/index.vue'
import { ElMessage as Message } from 'element-plus'

export default {
  name: 'SettingsSync',
  components: { PageHeader },
  data() {
    return {
      loading: false,
      saving: false,
      runsLoading: false,
      form: {
        autoSyncEnabled: false,
        autoSyncCron: '0 0 */6 * * *'
      },
      // 各来源是否已启用同步（来自 getConfig）
      sourceEnabled: { ldap: false, dingtalk: false, wecom: false, feishu: false },
      runs: [],
      runsTotal: 0,
      runParams: { pageNum: 1, pageSize: 10 },
      savedSnapshot: ''
    }
  },
  computed: {
    allSources() {
      return [
        { key: 'ldap', label: 'LDAP/AD', enabled: this.sourceEnabled.ldap },
        { key: 'dingtalk', label: '钉钉', enabled: this.sourceEnabled.dingtalk },
        { key: 'wecom', label: '企业微信', enabled: this.sourceEnabled.wecom },
        { key: 'feishu', label: '飞书', enabled: this.sourceEnabled.feishu }
      ]
    },
    enabledSources() {
      return this.allSources.filter(s => s.enabled)
    }
  },
  created() {
    this.fetchConfig()
    this.fetchRuns()
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
    sourceLabel(key) {
      const m = { ldap: 'LDAP/AD', dingtalk: '钉钉', wecom: '企业微信', feishu: '飞书' }
      return m[key] || key
    },
    formatTime(v) {
      if (!v) return '-'
      return String(v).replace('T', ' ').slice(0, 19)
    },
    async fetchConfig() {
      this.loading = true
      try {
        const { data } = await getConfig()
        this.form.autoSyncEnabled = !!data.autoSyncEnabled
        this.form.autoSyncCron = data.autoSyncCron || '0 0 */6 * * *'
        this.sourceEnabled = {
          ldap: !!data.ldapEnableSync,
          dingtalk: !!data.dingTalkEnableSync,
          wecom: !!data.weComEnableSync,
          feishu: !!data.feiShuEnableSync
        }
      } catch (e) {
        Message.error('获取配置失败')
      } finally {
        this.loading = false
        this.savedSnapshot = this.snapshot()
      }
    },
    async fetchRuns() {
      this.runsLoading = true
      try {
        const { data } = await getSyncRuns(this.runParams)
        this.runs = (data && data.list) || []
        this.runsTotal = (data && data.total) || 0
      } catch (e) {
        // 错误由请求拦截器统一提示
      } finally {
        this.runsLoading = false
      }
    },
    handleRunPage(page) {
      this.runParams.pageNum = page
      this.fetchRuns()
    },
    async save() {
      this.saving = true
      try {
        await updateSyncConfig(this.form)
        Message.success('定时同步配置已保存')
        this.fetchConfig()
      } catch (e) {
        // 错误由请求拦截器统一提示
      } finally {
        this.saving = false
      }
    },
    async runNow(s) {
      try {
        await runSyncNow({ source: s.key })
        Message.success(`已触发${s.label}同步，稍后刷新查看结果`)
        setTimeout(() => this.fetchRuns(), 1500)
      } catch (e) {
        // 错误由请求拦截器统一提示
      }
    }
  }
}
</script>

<style scoped>
.settings-card {
  margin: 10px;
}
.settings-card:last-child {
  margin-bottom: 100px;
}
.settings-alert {
  margin-bottom: 16px;
}
.form-tip {
  color: #909399;
  font-size: 12px;
  line-height: 1.6;
}
.src-tag {
  margin-right: 8px;
}
.sync-run-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}
.sync-run-title {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
}
.run-now-bar {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 10px;
}
.settings-footer {
  text-align: right;
  margin-top: 8px;
}
.pagination-wrap {
  display: flex;
  justify-content: flex-end;
  margin-top: 12px;
}
</style>
