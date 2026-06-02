<template>
  <div class="dashboard">
    <!-- 页头 -->
    <div class="dash-header">
      <div class="dash-header__main">
        <h1 class="dash-header__title">概览</h1>
        <p class="dash-header__desc">LDAP 管理平台运行状态与最近动态</p>
      </div>
      <div class="dash-header__right">
        <span class="status-pill" :class="ldapStatus.connected ? 'is-ok' : 'is-warn'">
          <span class="status-pill__dot" />{{ ldapStatus.connected ? 'LDAP 已连接' : 'LDAP 未连接' }}
        </span>
        <el-button icon="Refresh" circle plain :loading="loading" title="刷新" @click="refreshAll" />
      </div>
    </div>

    <!-- 未连接提示 -->
    <div v-if="!ldapStatus.connected" class="ldap-banner">
      <el-icon class="ldap-banner__icon"><WarningFilled /></el-icon>
      <div class="ldap-banner__text">
        <div class="ldap-banner__title">LDAP 尚未连接</div>
        <div class="ldap-banner__desc">配置目录服务（OpenLDAP / AD）后即可同步并管理用户与组织。</div>
      </div>
      <el-button type="primary" size="small" @click="goConfig">去配置</el-button>
    </div>

    <!-- 首屏骨架屏 -->
    <template v-if="loading">
      <el-skeleton animated :throttle="300">
        <template #template>
          <el-row :gutter="16" style="margin-bottom: 20px;">
            <el-col v-for="i in 6" :key="i" :xs="12" :sm="8" :lg="4">
              <el-skeleton-item variant="image" style="height: 116px; border-radius: 16px;" />
            </el-col>
          </el-row>
          <el-row :gutter="20" style="margin-bottom: 20px;">
            <el-col :lg="8" :xs="24"><el-skeleton-item variant="image" style="height: 340px; border-radius: 16px;" /></el-col>
            <el-col :lg="16" :xs="24"><el-skeleton-item variant="image" style="height: 340px; border-radius: 16px;" /></el-col>
          </el-row>
        </template>
      </el-skeleton>
    </template>

    <template v-else>
      <!-- 指标卡 -->
      <panel-group :data-info="dashboardList" @handleSetLineChartData="handleSetLineChartData" />

      <!-- 图表区 -->
      <el-row :gutter="20" class="dash-row">
        <el-col :xs="24" :lg="8">
          <section class="dash-card">
            <header class="dash-card__head">
              <div class="dash-card__titles">
                <h3 class="dash-card__title">目录构成</h3>
                <span class="dash-card__cap">各模块数据占比</span>
              </div>
            </header>
            <dashboard-pie :list="dashboardList" />
          </section>
        </el-col>
        <el-col :xs="24" :lg="16">
          <section class="dash-card">
            <header class="dash-card__head">
              <div class="dash-card__titles">
                <h3 class="dash-card__title">各模块数量</h3>
                <span class="dash-card__cap">按数量排行</span>
              </div>
            </header>
            <dashboard-bar :list="dashboardList" />
          </section>
        </el-col>
      </el-row>

      <el-row :gutter="20" class="dash-row">
        <el-col :xs="24" :lg="8">
          <section class="dash-card">
            <header class="dash-card__head">
              <div class="dash-card__titles">
                <h3 class="dash-card__title">数据分布</h3>
                <span class="dash-card__cap">各维度雷达视图</span>
              </div>
            </header>
            <dashboard-radar :list="dashboardList" />
          </section>
        </el-col>
        <el-col :xs="24" :lg="16">
          <section class="dash-card">
            <header class="dash-card__head">
              <div class="dash-card__titles">
                <h3 class="dash-card__title">最近操作</h3>
                <span class="dash-card__cap">最新审计日志</span>
              </div>
              <router-link to="/log/operation-log" class="dash-card__more">查看全部</router-link>
            </header>
            <el-table v-loading="recentOpsLoading" :data="recentOps" size="small" style="width: 100%">
              <el-table-column show-overflow-tooltip prop="username" label="请求人" width="96" />
              <el-table-column show-overflow-tooltip prop="path" label="请求路径" min-width="180" />
              <el-table-column prop="method" label="方式" width="76" align="center">
                <template #default="scope">
                  <el-tag v-if="scope.row.method === 'GET'" type="success" size="small" effect="light">GET</el-tag>
                  <el-tag v-else-if="scope.row.method === 'POST'" type="warning" size="small" effect="light">POST</el-tag>
                  <el-tag v-else-if="scope.row.method === 'PUT'" type="primary" size="small" effect="light">PUT</el-tag>
                  <el-tag v-else-if="scope.row.method === 'DELETE'" type="danger" size="small" effect="light">DEL</el-tag>
                  <el-tag v-else size="small" effect="light">{{ scope.row.method }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="status" label="状态" width="72" align="center">
                <template #default="scope">
                  <el-tag :type="scope.row.status >= 200 && scope.row.status < 300 ? 'success' : 'danger'" size="small" effect="light">{{ scope.row.status }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column show-overflow-tooltip prop="startTime" label="时间" width="160" />
              <template #empty>暂无操作记录</template>
            </el-table>
          </section>
        </el-col>
      </el-row>
    </template>
  </div>
</template>

<script>
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
      loading: true,
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
  // 普通用户登录后跳转到个人中心
  beforeRouteEnter(to, from, next) {
    next(vm => {
      const roles = vm.$store.getters.roles
      if (roles.length > 0 && roles.includes('普通用户')) {
        vm.$router.push('/profile/index')
      }
    })
  },
  methods: {
    handleSetLineChartData() {
      // 保留事件以兼容 PanelGroup 点击
    },
    goConfig() {
      this.$router.push('/settings/directory')
    },
    refreshAll() {
      this.loading = true
      this.fetchDashboard()
      this.fetchLDAPStatus()
      this.fetchRecentOps()
    },
    async fetchDashboard() {
      try {
        const res = await getDash()
        this.dashboardList = Array.isArray(res.data) ? res.data : []
      } catch (_) {
        this.dashboardList = []
      } finally {
        this.loading = false
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
        const res = await getOperationLogs({ pageNum: 1, pageSize: 8 })
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

<style lang="scss" scoped>
@import "@/styles/variables.scss";

.dashboard { padding: 0; }

/* 页头 */
.dash-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 22px;
  flex-wrap: wrap;

  &__title {
    margin: 0 0 6px;
    font-size: 22px;
    font-weight: 700;
    color: $slate900;
    letter-spacing: -0.02em;
  }
  &__desc { margin: 0; font-size: 14px; color: $slate500; }
  &__right { display: flex; align-items: center; gap: 12px; }
}

.status-pill {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  height: 30px;
  padding: 0 13px;
  border-radius: 999px;
  font-size: 13px;
  font-weight: 500;
  border: 1px solid;

  &__dot { width: 7px; height: 7px; border-radius: 50%; }
  &.is-ok {
    color: #15803d; background: #f0fdf4; border-color: #bbf7d0;
    .status-pill__dot { background: #22c55e; box-shadow: 0 0 0 3px rgba(34, 197, 94, 0.18); }
  }
  &.is-warn {
    color: #b45309; background: #fffbeb; border-color: #fde68a;
    .status-pill__dot { background: #f59e0b; box-shadow: 0 0 0 3px rgba(245, 158, 11, 0.18); }
  }
}

/* 未连接提示 */
.ldap-banner {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 13px 18px;
  margin-bottom: 22px;
  border-radius: 12px;
  background: #fffbeb;
  border: 1px solid #fde68a;

  &__icon { font-size: 22px; flex: none; color: $themeWarning; }
  &__text { flex: 1; min-width: 0; }
  &__title { font-size: 14px; font-weight: 600; color: #92400e; }
  &__desc { font-size: 13px; margin-top: 2px; color: #b45309; }
}

/* 卡片 */
.dash-row { margin-bottom: 20px; }
.dash-card {
  height: 100%;
  background: #fff;
  border-radius: 16px;
  border: 1px solid $borderColor;
  box-shadow: $cardShadow;
  padding: 20px 22px 18px;
  transition: box-shadow $transitionBase;

  &:hover { box-shadow: $cardShadowHover; }

  &__head {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 12px;
    margin-bottom: 8px;
  }
  &__title { margin: 0; font-size: 16px; font-weight: 600; color: $slate800; letter-spacing: -0.01em; }
  &__cap { display: block; margin-top: 3px; font-size: 12px; color: $slate400; }
  &__more {
    flex: none;
    font-size: 13px;
    font-weight: 500;
    color: $themePrimary;
    text-decoration: none;
    &:hover { color: $themePrimaryDark; }
  }

  :deep(.el-table) { --el-table-border-color: #f1f5f9; }
  :deep(.el-table th.el-table__cell) { background: #f8fafc; color: $slate500; font-weight: 600; }
}

@media (max-width: 1024px) {
  .dash-card { padding: 16px; }
  .dash-card :deep(.dashboard-bar-chart),
  .dash-card :deep(.dashboard-pie-chart),
  .dash-card :deep(.dashboard-radar-chart) { height: 260px !important; }
}
</style>
