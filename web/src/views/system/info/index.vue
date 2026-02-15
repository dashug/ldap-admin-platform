<template>
  <div class="system-info-page">
    <el-card class="box-card" shadow="hover">
      <div slot="header" class="card-header">
        <span>系统信息</span>
        <el-button type="text" icon="el-icon-refresh" @click="fetch">刷新</el-button>
      </div>
      <el-descriptions v-loading="loading" :column="1" border>
        <el-descriptions-item label="版本号">{{ info.version && info.version.version ? info.version.version : '-' }}</el-descriptions-item>
        <el-descriptions-item label="Git 提交">{{ info.version && info.version.gitCommit ? info.version.gitCommit : '-' }}</el-descriptions-item>
        <el-descriptions-item label="构建时间">{{ info.version && info.version.buildTime ? info.version.buildTime : '-' }}</el-descriptions-item>
        <el-descriptions-item label="Go 版本">{{ info.version && info.version.goVersion ? info.version.goVersion : '-' }}</el-descriptions-item>
        <el-descriptions-item label="运行时长">{{ info.uptime || '-' }}</el-descriptions-item>
        <el-descriptions-item label="数据库驱动">{{ info.dbDriver || '-' }}</el-descriptions-item>
        <el-descriptions-item label="数据库状态">
          <el-tag :type="info.dbStatus === '正常' ? 'success' : 'danger'" size="small">{{ info.dbStatus || '-' }}</el-tag>
          <span v-if="info.dbMessage" class="db-msg">{{ info.dbMessage }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="账户过期策略">
          {{ info.inactiveDays > 0 ? (info.inactiveDays + ' 天未登录自动禁用') : '未开启' }}
        </el-descriptions-item>
      </el-descriptions>
    </el-card>
  </div>
</template>

<script>
import { getSystemInfo } from '@/api/system/base'

export default {
  name: 'SystemInfo',
  data() {
    return {
      loading: false,
      info: {
        version: {},
        uptime: '',
        dbDriver: '',
        dbStatus: '',
        dbMessage: '',
        inactiveDays: 0
      }
    }
  },
  mounted() {
    this.fetch()
  },
  methods: {
    async fetch() {
      this.loading = true
      try {
        const res = await getSystemInfo()
        if (res.data) {
          this.info = {
            version: res.data.version || {},
            uptime: res.data.uptime || '',
            dbDriver: res.data.dbDriver || '',
            dbStatus: res.data.dbStatus || '',
            dbMessage: res.data.dbMessage || '',
            inactiveDays: res.data.inactiveDays != null ? res.data.inactiveDays : 0
          }
        }
      } catch (_) {
        this.$message.error('获取系统信息失败')
      } finally {
        this.loading = false
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.system-info-page {
  padding: 0;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.db-msg {
  margin-left: 8px;
  color: #f56c6c;
  font-size: 12px;
}
</style>
