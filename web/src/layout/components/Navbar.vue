<template>
  <div class="navbar">
    <hamburger id="hamburger-container" :is-active="sidebar.opened" class="hamburger-container" @toggleClick="toggleSideBar" />

    <breadcrumb id="breadcrumb-container" class="breadcrumb-container" />

    <div class="right-menu">
      <template v-if="device!=='mobile'">
        <el-tooltip content="搜索" effect="dark" placement="bottom">
          <search id="header-search" class="right-menu-item" />
        </el-tooltip>

        <error-log class="errLog-container right-menu-item hover-effect" />

        <el-tooltip content="全屏" effect="dark" placement="bottom">
          <screenfull id="screenfull" class="right-menu-item hover-effect" />
        </el-tooltip>

        <el-tooltip content="布局大小" effect="dark" placement="bottom">
          <size-select id="size-select" class="right-menu-item hover-effect" />
        </el-tooltip>

        <el-tooltip content="GitHub" effect="dark" placement="bottom">
          <el-link style="font-size: 23px;" class="iconfont icon-github right-menu-item" href="https://github.com/dashug/go-ldap-admin" :underline="false" target="_blank" />
        </el-tooltip>

      </template>

      <el-dropdown class="avatar-container right-menu-item hover-effect" trigger="click">
        <div class="avatar-wrapper">
          <img :src="avatar+'?imageView2/1/w/80/h/80'" class="user-avatar">
          <i class="el-icon-caret-bottom" />
        </div>
        <el-dropdown-menu slot="dropdown">
          <router-link to="/profile/index">
            <el-dropdown-item>个人中心</el-dropdown-item>
          </router-link>
          <el-dropdown-item divided @click.native="logout">
            <span style="display:block;">退出登陆</span>
          </el-dropdown-item>
        </el-dropdown-menu>
      </el-dropdown>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import Breadcrumb from '@/components/Breadcrumb'
import Hamburger from '@/components/Hamburger'
import ErrorLog from '@/components/ErrorLog'
import Screenfull from '@/components/Screenfull'
import SizeSelect from '@/components/SizeSelect'
import Search from '@/components/HeaderSearch'
import '@/assets/iconfont/font/iconfont.css'

export default {
  components: {
    Breadcrumb,
    Hamburger,
    ErrorLog,
    Screenfull,
    SizeSelect,
    Search
  },
  computed: {
    ...mapGetters([
      'sidebar',
      'avatar',
      'device'
    ])

  },
  methods: {
    toggleSideBar() {
      this.$store.dispatch('app/toggleSideBar')
    },
    async logout() {
      await this.$store.dispatch('user/logout')
      this.$router.push(`/login?redirect=${this.$route.fullPath}`)
    }
  }
}
</script>

<style lang="scss" scoped>
@import "~@/styles/variables.scss";

.head-github {
  cursor: pointer;
  font-size: 18px;
  vertical-align: middle;
}

.navbar {
  height: $headerHeight;
  overflow: hidden;
  position: relative;
  background: #fff;
  border-bottom: 1px solid $borderColor;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.04);
  display: flex;
  align-items: center;
  padding: 0 24px;

  .hamburger-container {
    height: 100%;
    display: flex;
    align-items: center;
    padding: 0 12px;
    margin: 0 -12px 0 -16px;
    cursor: pointer;
    transition: background 0.2s ease;
    -webkit-tap-highlight-color: transparent;
    border-radius: 8px;

    &:hover {
      background: rgba(0, 0, 0, 0.04);
    }
  }

  .breadcrumb-container {
    flex: 1;
    margin-left: 8px;
  }

  .errLog-container {
    display: inline-block;
    vertical-align: middle;
  }

  .right-menu {
    display: flex;
    align-items: center;
    height: 100%;
    gap: 4px;

    &:focus {
      outline: none;
    }

    .right-menu-item {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      padding: 0 12px;
      height: 40px;
      min-width: 40px;
      font-size: 18px;
      color: $slate500;
      border-radius: 10px;
      transition: color $transitionBase, background $transitionBase;

      &.hover-effect {
        cursor: pointer;
        &:hover {
          color: $slate700;
          background: $slate100;
        }
      }
    }

    .avatar-container {
      margin-left: 12px;
      padding-left: 16px;
      border-left: 1px solid $borderColor;

      .avatar-wrapper {
        display: flex;
        align-items: center;
        gap: 8px;
        cursor: pointer;
        padding: 6px 10px;
        border-radius: 12px;
        transition: background $transitionBase;

        &:hover {
          background: $slate100;
        }

        .user-avatar {
          width: 36px;
          height: 36px;
          border-radius: 10px;
          object-fit: cover;
        }

        .el-icon-caret-bottom {
          font-size: 12px;
          color: $slate400;
        }
      }
    }
  }
}
</style>
