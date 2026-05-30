<template>
  <div class="nav">
    <!-- 品牌 -->
    <div class="nav__brand">
      <div class="nav__logo"><el-icon :size="18"><Connection /></el-icon></div>
      <span v-if="!isCollapse" class="nav__brandname">LDAP 管理平台</span>
    </div>

    <!-- 搜索 / ⌘K -->
    <button class="nav__search" @click="openCmdk">
      <el-icon><Search /></el-icon>
      <template v-if="!isCollapse">
        <span class="nav__search-text">搜索…</span>
        <kbd class="nav__kbd">⌘K</kbd>
      </template>
    </button>

    <!-- 导航 -->
    <div class="nav__scroll">
      <template v-for="group in groups" :key="group.key">
        <div v-if="group.label && !isCollapse" class="nav__section">{{ group.label }}</div>
        <router-link
          v-for="it in group.items"
          :key="it.path"
          :to="it.path"
          class="nav__item"
          :class="{ 'is-active': isActive(it.path) }"
          :title="it.title"
        >
          <svg-icon v-if="it.icon && !it.icon.includes('el-icon')" :icon-class="it.icon" class="nav__icon" />
          <el-icon v-else class="nav__icon"><Document /></el-icon>
          <span v-if="!isCollapse" class="nav__label">{{ it.title }}</span>
        </router-link>
      </template>

      <!-- 设置（常驻可见的配置入口，提升可发现性） -->
      <div v-if="!isCollapse" class="nav__section">设置</div>
      <a class="nav__item" title="目录配置" @click="openSetting('directory')">
        <el-icon class="nav__icon"><Setting /></el-icon>
        <span v-if="!isCollapse" class="nav__label">目录配置</span>
      </a>
      <a class="nav__item" title="平台对接" @click="openSetting('thirdparty')">
        <el-icon class="nav__icon"><Connection /></el-icon>
        <span v-if="!isCollapse" class="nav__label">平台对接</span>
      </a>
      <a class="nav__item" title="通知设置" @click="openSetting('notification')">
        <el-icon class="nav__icon"><Bell /></el-icon>
        <span v-if="!isCollapse" class="nav__label">通知设置</span>
      </a>
    </div>

    <!-- 用户 -->
    <el-dropdown class="nav__user" trigger="click" placement="top-start" @command="onUserCommand">
      <div class="nav__user-row">
        <img :src="avatar" class="nav__avatar">
        <template v-if="!isCollapse">
          <div class="nav__user-meta">
            <div class="nav__user-name">{{ name || 'admin' }}</div>
            <div class="nav__user-sub">管理后台</div>
          </div>
          <el-icon class="nav__user-caret"><ArrowDown /></el-icon>
        </template>
      </div>
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item command="profile" icon="User">个人中心</el-dropdown-item>
          <el-dropdown-item command="logout" icon="SwitchButton" divided>退出登录</el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  name: 'AppSidebar',
  computed: {
    ...mapGetters(['permission_routes', 'sidebar', 'avatar', 'name']),
    isCollapse() {
      return !this.sidebar.opened
    },
    // 扁平分组：[{label, items:[{title,path,icon}]}]
    groups() {
      const out = []
      ;(this.permission_routes || []).forEach((route, ri) => {
        if (route.hidden) return
        const base = route.path
        const children = (route.children || []).filter(c => !c.hidden)
        const isGroup = route.alwaysShow || children.length > 1
        if (isGroup && children.length) {
          out.push({
            key: 'g' + ri,
            label: route.meta && route.meta.title,
            items: children.map(c => this.toItem(c, base))
          })
        } else {
          // 单页：取唯一子项，否则取自身
          const leaf = children.length === 1 ? children[0] : route
          const base2 = children.length === 1 ? base : ''
          out.push({ key: 'i' + ri, label: '', items: [this.toItem(leaf, base2)] })
        }
      })
      return out
    }
  },
  methods: {
    toItem(route, base) {
      return {
        title: (route.meta && route.meta.title) || route.name || '',
        path: this.resolvePath(base, route.path),
        icon: (route.meta && route.meta.icon) || ''
      }
    },
    resolvePath(base, path) {
      if (!path) return base || '/'
      if (path.startsWith('/')) return path
      return ((base || '').replace(/\/$/, '') + '/' + path).replace(/\/+/g, '/')
    },
    isActive(path) {
      const cur = this.$route.path
      return cur === path || cur.startsWith(path + '/')
    },
    openCmdk() {
      window.dispatchEvent(new Event('open-cmdk'))
    },
    openSetting(type) {
      // 跳到用户页并自动打开对应配置弹窗（带时间戳确保 keep-alive 下重复点击也生效）
      this.$router.push({ path: '/personnel/user', query: { openConfig: type, t: Date.now() } }).catch(() => {})
    },
    onUserCommand(cmd) {
      if (cmd === 'profile') this.$router.push('/profile/index').catch(() => {})
      else if (cmd === 'logout') this.logout()
    },
    async logout() {
      await this.$store.dispatch('user/logout')
      this.$router.push(`/login?redirect=${this.$route.fullPath}`)
    }
  }
}
</script>

<style lang="scss" scoped>
@import "@/styles/variables.scss";

.nav {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #fff;
  padding: 14px 12px 12px;
}

.nav__brand {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 4px 8px 14px;
  .nav__logo {
    width: 32px;
    height: 32px;
    flex: none;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 9px;
    color: #fff;
    background: linear-gradient(135deg, $themePrimary, $themePrimaryLight);
  }
  .nav__brandname {
    font-size: 15px;
    font-weight: $fontWeightSemibold;
    color: $slate800;
    white-space: nowrap;
  }
}

.nav__search {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  height: 36px;
  margin-bottom: 12px;
  padding: 0 10px;
  border: 1px solid $slate200;
  border-radius: 9px;
  background: $slate50;
  color: $slate400;
  font-size: 13px;
  cursor: pointer;
  transition: all $transitionBase;
  &:hover { background: #fff; border-color: $slate300; }
  .nav__search-text { flex: 1; text-align: left; }
  .nav__kbd {
    font-family: $fontFamilyBase;
    font-size: 11px;
    color: $slate500;
    background: #fff;
    border: 1px solid $slate200;
    border-radius: 5px;
    padding: 0 5px;
  }
}

.nav__scroll {
  flex: 1;
  overflow-y: auto;
  margin: 0 -4px;
  padding: 0 4px;
  &::-webkit-scrollbar { width: 4px; }
  &::-webkit-scrollbar-thumb { background: $slate200; border-radius: 4px; }
}

.nav__section {
  font-size: 11px;
  font-weight: $fontWeightSemibold;
  letter-spacing: 0.04em;
  color: $slate400;
  padding: 14px 10px 5px;
  text-transform: uppercase;
}

.nav__item {
  display: flex;
  align-items: center;
  gap: 11px;
  height: 38px;
  padding: 0 10px;
  margin: 1px 0;
  border-radius: 8px;
  color: $slate600;
  font-size: 14px;
  text-decoration: none;
  cursor: pointer;
  transition: background $transitionBase, color $transitionBase;
  .nav__icon { font-size: 17px; flex: none; }
  .nav__label { white-space: nowrap; overflow: hidden; }
  &:hover { background: $slate100; color: $slate900; }
  &.is-active {
    background: $menuActiveBg;
    color: $themePrimary;
    font-weight: $fontWeightSemibold;
  }
}

.nav__user {
  margin-top: 8px;
  padding-top: 10px;
  border-top: 1px solid $slate100;
  display: block;
  width: 100%;
}
.nav__user-row {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 7px 8px;
  border-radius: 9px;
  cursor: pointer;
  transition: background $transitionBase;
  &:hover { background: $slate100; }
  .nav__avatar { width: 30px; height: 30px; border-radius: 8px; object-fit: cover; flex: none; }
  .nav__user-meta { flex: 1; min-width: 0; }
  .nav__user-name { font-size: 13px; font-weight: $fontWeightMedium; color: $slate800; line-height: 1.2; }
  .nav__user-sub { font-size: 11px; color: $slate400; }
  .nav__user-caret { color: $slate400; font-size: 12px; }
}
</style>
