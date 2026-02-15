<template>
  <div class="err-page err-404">
    <div class="err-content">
      <div class="err-code">404</div>
      <p class="err-title">页面不存在</p>
      <p class="err-desc">请检查地址是否正确，或返回首页</p>
      <router-link to="/" class="err-btn">返回首页</router-link>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Page404'
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
  color: $themePrimary;
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

.err-btn {
  display: inline-block;
  padding: 10px 24px;
  background: $themePrimary;
  color: #fff;
  border-radius: 8px;
  font-size: 14px;
  text-decoration: none;
  transition: opacity 0.2s;
  &:hover {
    opacity: 0.9;
    color: #fff;
  }
}
</style>
