<template>
  <el-breadcrumb class="app-breadcrumb" separator="/">
    <transition-group name="breadcrumb">
      <el-breadcrumb-item v-for="(item,index) in levelList" :key="item.path">
        <span v-if="item.redirect==='noRedirect'||index==levelList.length-1" class="no-redirect">{{ item.meta.title }}</span>
        <a v-else @click.prevent="handleLink(item)">{{ item.meta.title }}</a>
      </el-breadcrumb-item>
    </transition-group>
  </el-breadcrumb>
</template>

<script>
import { compile } from 'path-to-regexp'

export default {
  data() {
    return {
      levelList: null
    }
  },
  watch: {
    $route(route) {
      if (route.path.startsWith('/redirect/')) {
        return
      }
      this.getBreadcrumb()
    }
  },
  created() {
    this.getBreadcrumb()
  },
  methods: {
    getBreadcrumb() {
      let matched = this.$route.matched.filter(item => item.meta && item.meta.title)
      const first = matched[0]

      if (!this.isDashboard(first)) {
        matched = [{ path: '/dashboard', meta: { title: '首页' }}].concat(matched)
      }

      this.levelList = matched.filter(item => item.meta && item.meta.title && item.meta.breadcrumb !== false)
    },
    isDashboard(route) {
      const name = route && route.name
      if (!name) {
        return false
      }
      return name.trim().toLocaleLowerCase() === 'Dashboard'.toLocaleLowerCase()
    },
    pathCompile(path) {
      // path-to-regexp v6+ 改为命名导出 compile；对无法编译的路径兜底返回原始 path
      try {
        const { params } = this.$route
        return compile(path)(params)
      } catch (e) {
        return path
      }
    },
    handleLink(item) {
      const { redirect, path } = item
      if (redirect) {
        this.$router.push(redirect)
        return
      }
      this.$router.push(this.pathCompile(path))
    }
  }
}
</script>

<style lang="scss" scoped>
@import "@/styles/variables.scss";

.app-breadcrumb.el-breadcrumb {
  display: inline-block;
  font-size: 13px;
  line-height: 56px;
  margin-left: 4px;

  .el-breadcrumb__item {
    .el-breadcrumb__inner {
      color: #64748b;
      font-weight: 400;
      &:hover { color: $themePrimary; }
    }
    &:last-child .el-breadcrumb__inner {
      color: #1e293b;
      font-weight: 500;
    }
  }

  .no-redirect {
    color: #1e293b;
    cursor: text;
    font-weight: 500;
  }
}
</style>
