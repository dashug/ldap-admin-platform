<template>
  <div class="sidebar-logo-container" :class="{'collapse':collapse}">
    <transition name="sidebarLogoFade">
      <router-link v-if="collapse" key="collapse" class="sidebar-logo-link" to="/">
        <img v-if="logo" :src="logo" class="sidebar-logo">
        <h1 v-else class="sidebar-title">{{ title }} </h1>
      </router-link>
      <router-link v-else key="expand" class="sidebar-logo-link" to="/">
        <img v-if="logo" :src="logo" class="sidebar-logo">
        <h1 class="sidebar-title">{{ title }} </h1>
      </router-link>
    </transition>
  </div>
</template>

<script>
export default {
  name: 'SidebarLogo',
  props: {
    collapse: {
      type: Boolean,
      required: true
    }
  },
  data() {
    return {
      title: 'LDAP 管理平台',
      logo: require('@/assets/sidebar-logo/logo-icon.svg')
    }
  }
}
</script>

<style lang="scss" scoped>
@import "~@/styles/variables.scss";

.sidebarLogoFade-enter-active {
  transition: opacity 0.25s ease;
}

.sidebarLogoFade-enter,
.sidebarLogoFade-leave-to {
  opacity: 0;
}

.sidebar-logo-container {
  position: relative;
  width: 100%;
  height: 60px;
  line-height: 60px;
  background: $menuBg;
  text-align: center;
  overflow: hidden;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  flex-shrink: 0;

  .sidebar-logo-link {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100%;
    width: 100%;
    padding: 0 20px;
    text-decoration: none;
  }

  .sidebar-logo {
    width: 38px;
    height: 38px;
    vertical-align: middle;
    margin-right: 12px;
    border-radius: 10px;
  }

  .sidebar-title {
    display: inline-block;
    margin: 0;
    color: #fff;
    font-weight: 600;
    font-size: 16px;
    letter-spacing: -0.02em;
    vertical-align: middle;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  &.collapse {
    .sidebar-logo {
      margin-right: 0;
    }
    .sidebar-title {
      display: none;
    }
  }
}
</style>
