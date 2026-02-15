<template>
  <div class="dashboard-container">
    <div class="dashboard-editor-container">
      <!-- 页面标题 -->
      <div class="page-header">
        <h1 class="page-title">概览</h1>
        <p class="page-desc">LDAP 管理平台运行状态与最近动态</p>
      </div>
      <!-- LDAP 连接状态 -->
      <el-row class="ldap-status-row">
        <el-col :span="24">
          <el-alert
            :title="ldapStatus.connected ? 'LDAP 连接正常' : 'LDAP 未连接'"
            :description="ldapStatus.message"
            :type="ldapStatus.connected ? 'success' : 'warning'"
            :closable="false"
            show-icon
          />
        </el-col>
      </el-row>
      <panel-group :data-info="dashboardList" @handleSetLineChartData="handleSetLineChartData" />
      <!-- 图表看板（真实数据） -->
      <el-row :gutter="32">
        <el-col :xs="24" :sm="24" :lg="8">
          <div class="chart-wrapper">
            <dashboard-radar :list="dashboardList" />
          </div>
        </el-col>
        <el-col :xs="24" :sm="24" :lg="8">
          <div class="chart-wrapper">
            <dashboard-pie :list="dashboardList" />
          </div>
        </el-col>
        <el-col :xs="24" :sm="24" :lg="8">
          <div class="chart-wrapper">
            <dashboard-bar :list="dashboardList" />
          </div>
        </el-col>
      </el-row>
      <!-- 最近操作 -->
      <el-row class="recent-ops-row">
        <el-col :span="24">
          <div class="chart-wrapper">
            <div class="recent-ops-header">
              <span>最近操作</span>
              <router-link to="/log/operation-log" class="link-more">查看全部</router-link>
            </div>
            <el-table v-loading="recentOpsLoading" :data="recentOps" size="small" stripe style="width: 100%">
              <el-table-column show-overflow-tooltip prop="username" label="请求人" width="100" />
              <el-table-column show-overflow-tooltip prop="path" label="请求路径" min-width="180" />
              <el-table-column prop="method" label="方式" width="70" align="center">
                <template slot-scope="scope">
                  <el-tag v-if="scope.row.method === 'GET'" type="success" size="mini">GET</el-tag>
                  <el-tag v-else-if="scope.row.method === 'POST'" type="warning" size="mini">POST</el-tag>
                  <el-tag v-else-if="scope.row.method === 'PUT'" type="primary" size="mini">PUT</el-tag>
                  <el-tag v-else-if="scope.row.method === 'DELETE'" type="danger" size="mini">DEL</el-tag>
                  <el-tag v-else size="mini">{{ scope.row.method }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="status" label="状态" width="70" align="center">
                <template slot-scope="scope">
                  <el-tag :type="scope.row.status >= 200 && scope.row.status < 300 ? 'success' : 'danger'" size="mini">{{ scope.row.status }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column show-overflow-tooltip prop="startTime" label="时间" width="160" />
            </el-table>
          </div>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script>
// import GithubCorner from '@/components/GithubCorner'
import PanelGroup from './components/PanelGroup'
import DashboardRadar from './components/DashboardRadar'
import DashboardPie from './components/DashboardPie'
import DashboardBar from './components/DashboardBar'
import { getDash } from '@/api/dashboards/dashboard'
import { getLDAPStatus } from '@/api/system/base'
import { getOperationLogs } from '@/api/log/operationLog'
import { mapGetters } from 'vuex'

export default {
  name: 'Dashboard',
  components: {
    PanelGroup,
    DashboardRadar,
    DashboardPie,
    DashboardBar
  },
  computed: {
    ...mapGetters(['roles'])
  },
  data() {
    return {
      dashboardList: [],
      ldapStatus: { connected: false, message: '检测中...' },
      recentOps: [],
      recentOpsLoading: false
    }
  },
  mounted() {
    this.fetchDashboard()
    this.fetchLDAPStatus()
    this.fetchRecentOps()
  },
  //普通用户登录后跳转到个人中心
 beforeRouteEnter(to, from, next) {
    next(vm => {
      const roles = vm.$store.getters.roles;
      if (roles.length > 0 && roles.includes('普通用户')) {
        vm.$router.push('/profile/index');
      }
    });
  },

   methods: {
    handleSetLineChartData() {
      // 保留事件以兼容 PanelGroup 点击
    },
    async fetchDashboard() {
      try {
        const res = await getDash()
        this.dashboardList = Array.isArray(res.data) ? res.data : []
      } catch (_) {
        this.dashboardList = []
      }
    },
    async fetchLDAPStatus() {
      try {
        const res = await getLDAPStatus()
        if (res.data) {
          this.ldapStatus = { connected: res.data.connected, message: res.data.message || '' }
        }
      } catch (_) {
        this.ldapStatus = { connected: false, message: '获取状态失败' }
      }
    },
    async fetchRecentOps() {
      this.recentOpsLoading = true
      try {
        const res = await getOperationLogs({ pageNum: 1, pageSize: 10 })
        this.recentOps = (res.data && res.data.logs) ? res.data.logs : []
      } catch (_) {
        this.recentOps = []
      } finally {
        this.recentOpsLoading = false
      }
    }
  }
}
</script>

<style rel="stylesheet/scss" lang="scss" scoped>
@import "~@/styles/variables.scss";

.dashboard-container {
  padding: 0;
}

.page-header {
  margin-bottom: 24px;

  .page-title {
    margin: 0 0 8px 0;
    font-size: 22px;
    font-weight: 700;
    color: $slate800;
    letter-spacing: -0.02em;
  }

  .page-desc {
    margin: 0;
    font-size: 14px;
    color: $slate500;
  }
}

.ldap-status-row {
  margin-bottom: 20px;
  ::v-deep .el-alert {
    border-radius: 12px;
    padding: 14px 16px;
  }
}

.recent-ops-row {
  margin-top: 0;
}

.recent-ops-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  font-weight: 600;
  font-size: 16px;
  color: $slate800;

  .link-more {
    font-size: 14px;
    font-weight: 500;
    color: $themePrimary;
    text-decoration: none;
    &:hover { color: $themePrimaryDark; }
  }
}

.dashboard-editor-container {
  padding: 0;
  position: relative;
  min-height: 100%;

  .chart-wrapper {
    background: #fff;
    padding: 24px;
    margin-bottom: 24px;
    border-radius: $cardRadius;
    border: 1px solid $borderColor;
    box-shadow: $cardShadow;
    transition: box-shadow $transitionBase;

    &:hover {
      box-shadow: $cardShadowHover;
    }
  }
}

@media (max-width: 1024px) {
  .chart-wrapper {
    padding: 16px;
  }
}
</style>
