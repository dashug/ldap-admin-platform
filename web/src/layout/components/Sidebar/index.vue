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
          <el-icon class="nav__icon"><component :is="iconFor(it.title)" /></el-icon>
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
    // 统一为 Element Plus 线性图标（替代旧 svg 雪碧图，更精致一致）
    iconFor(title) {
      const m = {
        '首页': 'Odometer',
        '用户': 'User',
        '部门': 'OfficeBuilding',
        '同步字段映射': 'Switch',
        '角色与权限': 'UserFilled',
        '菜单': 'Menu',
        '接口': 'Share',
        '系统信息': 'Monitor',
        'API 密钥': 'Key',
        'API密钥': 'Key',
        '操作日志': 'Tickets',
        '分组成员': 'User',
        '个人中心': 'User'
      }
      return m[title] || 'Menu'
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
  background: linear-gradient(185deg, #1e293b 0%, #0f172a 60%, #0b1120 100%);
  padding: 16px 14px 14px;
}

.nav__brand {
  display: flex;
  align-items: center;
  gap: 11px;
  padding: 4px 8px 18px;
  .nav__logo {
    width: 34px;
    height: 34px;
    flex: none;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 10px;
    color: #fff;
    background: linear-gradient(135deg, $themePrimary, $themePrimaryLight);
    box-shadow: 0 4px 12px rgba(79, 70, 229, 0.45);
  }
  .nav__brandname {
    font-size: 15px;
    font-weight: $fontWeightSemibold;
    color: #f8fafc;
    white-space: nowrap;
  }
}

.nav__search {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  height: 38px;
  margin-bottom: 14px;
  padding: 0 11px;
  border: 1px solid rgba(255, 255, 255, 0.10);
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.05);
  color: $slate400;
  font-size: 13px;
  cursor: pointer;
  transition: all $transitionBase;
  &:hover { background: rgba(255, 255, 255, 0.09); border-color: rgba(255, 255, 255, 0.16); }
  .nav__search-text { flex: 1; text-align: left; }
  .nav__kbd {
    font-family: $fontFamilyBase;
    font-size: 11px;
    color: $slate400;
    background: rgba(255, 255, 255, 0.08);
    border: 1px solid rgba(255, 255, 255, 0.10);
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
  &::-webkit-scrollbar-thumb { background: rgba(255, 255, 255, 0.12); border-radius: 4px; }
}

.nav__section {
  font-size: 11px;
  font-weight: $fontWeightSemibold;
  letter-spacing: 0.06em;
  color: #64748b;
  padding: 16px 10px 6px;
  text-transform: uppercase;
}

.nav__item {
  display: flex;
  align-items: center;
  gap: 12px;
  height: 40px;
  padding: 0 11px;
  margin: 2px 0;
  border-radius: 9px;
  color: #cbd5e1;
  font-size: 14px;
  font-weight: $fontWeightMedium;
  text-decoration: none;
  cursor: pointer;
  transition: background $transitionBase, color $transitionBase, box-shadow $transitionBase;
  .nav__icon { font-size: 18px; width: 18px; flex: none; color: #64748b; transition: color $transitionBase; }
  .nav__label { white-space: nowrap; overflow: hidden; }
  &:hover {
    background: rgba(255, 255, 255, 0.06);
    color: #fff;
    .nav__icon { color: #cbd5e1; }
  }
  &.is-active {
    background: linear-gradient(135deg, rgba(99, 102, 241, 0.30), rgba(79, 70, 229, 0.18));
    color: #fff;
    font-weight: $fontWeightSemibold;
    box-shadow: 0 4px 12px rgba(79, 70, 229, 0.28), inset 0 0 0 1px rgba(129, 140, 248, 0.30);
    .nav__icon { color: #a5b4fc; }
  }
}

.nav__user {
  margin-top: 8px;
  padding-top: 10px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
  display: block;
  width: 100%;
}
.nav__user-row {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px;
  border-radius: 10px;
  cursor: pointer;
  transition: background $transitionBase;
  &:hover { background: rgba(255, 255, 255, 0.06); }
  .nav__avatar { width: 32px; height: 32px; border-radius: 9px; object-fit: cover; flex: none; }
  .nav__user-meta { flex: 1; min-width: 0; }
  .nav__user-name { font-size: 13px; font-weight: $fontWeightSemibold; color: #f1f5f9; line-height: 1.2; }
  .nav__user-sub { font-size: 11px; color: #64748b; }
  .nav__user-caret { color: #64748b; font-size: 12px; }
}
</style>
