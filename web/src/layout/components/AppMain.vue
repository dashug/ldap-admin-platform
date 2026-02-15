<template>
  <section class="app-main">
    <transition name="fade-transform" mode="out-in">
      <keep-alive :include="cachedViews">
        <router-view :key="key" />
      </keep-alive>
    </transition>
    <el-footer  class="footer-copyright">
      <div>
        <span>Since 2022 </span>
        <el-divider direction="vertical" />
        <span>Powered by </span>
          <span>
            <a href="https://github.com/dashug/ldap-admin-platform" target="_blank">LDAP 管理平台</a>
          </span>
        <el-divider direction="vertical" />
        <span>Copyright </span>
          <span>
            <a href="https://github.com/dashug/ldap-admin-platform" target="_blank">dashug</a>
          </span>
      </div>
    </el-footer>
  </section>
</template>

<script>
export default {
  name: 'AppMain',
  computed: {
    cachedViews() {
      return this.$store.state.tagsView.cachedViews
    },
    key() {
      return this.$route.path
    }
  }
}
</script>

<style lang="scss" scoped>
@import "~@/styles/variables.scss";

.app-main {
  min-height: calc(100vh - #{$headerHeight});
  width: 100%;
  position: relative;
  overflow-x: hidden;
  overflow-y: auto;
  padding: 20px 24px 72px;
}

.fixed-header + .app-main {
  padding-top: $headerHeight;
}

.footer-copyright {
  position: fixed;
  bottom: 0;
  left: $sideBarWidth;
  right: 0;
  height: 48px;
  line-height: 48px;
  text-align: center;
  font-size: 12px;
  color: #94a3b8;
  background: #fff;
  border-top: 1px solid $borderColor;
  z-index: 0;
  transition: left 0.28s ease;

  a {
    color: #64748b;
    &:hover {
      color: $themePrimary;
    }
  }
}

.hideSidebar .footer-copyright {
  left: $sidebarCollapseWidth;
}

.hasTagsView {
  .app-main {
    min-height: calc(100vh - 90px);
  }

  .fixed-header + .app-main {
    padding-top: 90px;
  }
}
</style>

<style lang="scss">
// fix css style bug in open el-dialog
.el-popup-parent--hidden {
  .fixed-header {
    padding-right: 15px;
  }
}
</style>
