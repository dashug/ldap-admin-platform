<template>
  <div class="err-page err-401">
    <div class="err-content">
      <div class="err-code">403</div>
      <p class="err-title">没有访问权限</p>
      <p class="err-desc">当前账号无权访问该页面，请联系管理员或返回</p>
      <div class="err-actions">
        <el-button icon="el-icon-arrow-left" @click="back">返回上页</el-button>
        <router-link to="/dashboard">
          <el-button type="primary">返回首页</el-button>
        </router-link>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Page401',
  methods: {
    back() {
      if (this.$route.query.noGoBack) {
        this.$router.push({ path: '/dashboard' })
      } else {
        this.$router.go(-1)
      }
    }
  }
}
</script>

<style lang="scss" scoped>
@import "~@/styles/variables.scss";

.err-page {
  min-height: calc(100vh - 56px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
}

.err-content {
  text-align: center;
}

.err-code {
  font-size: 72px;
  font-weight: 700;
  line-height: 1;
  color: #f59e0b;
  letter-spacing: -2px;
  margin-bottom: 16px;
}

.err-title {
  font-size: 18px;
  font-weight: 600;
  color: #1e293b;
  margin: 0 0 8px;
}

.err-desc {
  font-size: 14px;
  color: #64748b;
  margin: 0 0 24px;
}

.err-actions {
  display: flex;
  gap: 12px;
  justify-content: center;
  flex-wrap: wrap;

  a {
    text-decoration: none;
  }
}
</style>
