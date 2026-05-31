<template>
  <transition name="cmdk-fade">
    <div v-if="visible" class="cmdk-mask" @click.self="close">
      <div class="cmdk-panel" role="dialog" aria-modal="true">
        <div class="cmdk-search">
          <el-icon class="cmdk-search__icon"><Search /></el-icon>
          <input
            ref="input"
            v-model="query"
            class="cmdk-search__input"
            placeholder="搜索页面、执行操作…"
            @keydown.down.prevent="move(1)"
            @keydown.up.prevent="move(-1)"
            @keydown.enter.prevent="run()"
            @keydown.esc.prevent="close"
          >
          <kbd class="cmdk-kbd">Esc</kbd>
        </div>
        <div ref="list" class="cmdk-list">
          <template v-for="(group, gi) in grouped" :key="group.name">
            <div class="cmdk-group">{{ group.name }}</div>
            <div
              v-for="item in group.items"
              :key="item.key"
              class="cmdk-item"
              :class="{ 'is-active': item.idx === active }"
              @mouseenter="active = item.idx"
              @click="run(item)"
            >
              <el-icon class="cmdk-item__icon"><component :is="item.icon" /></el-icon>
              <span class="cmdk-item__title">{{ item.title }}</span>
              <span v-if="item.hint" class="cmdk-item__hint">{{ item.hint }}</span>
            </div>
            <div v-if="gi < grouped.length - 1" class="cmdk-sep" />
          </template>
          <div v-if="flat.length === 0" class="cmdk-empty">无匹配结果</div>
        </div>
        <div class="cmdk-footer">
          <span><kbd>↑</kbd><kbd>↓</kbd> 选择</span>
          <span><kbd>↵</kbd> 打开</span>
          <span><kbd>⌘</kbd><kbd>K</kbd> 唤起</span>
        </div>
      </div>
    </div>
  </transition>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  name: 'CommandPalette',
  data() {
    return { visible: false, query: '', active: 0 }
  },
  computed: {
    ...mapGetters(['permission_routes', 'isAdmin']),
    // 导航项：扁平化路由
    navItems() {
      const out = []
      const walk = (routes, basePath = '') => {
        (routes || []).forEach(r => {
          if (r.hidden) return
          const full = this.resolvePath(basePath, r.path)
          const hasChildren = r.children && r.children.length
          if (r.meta && r.meta.title && (!hasChildren || r.alwaysShow)) {
            out.push({ type: 'nav', title: r.meta.title, path: full, icon: 'Document', hint: '页面' })
          }
          if (hasChildren) walk(r.children, full)
        })
      }
      walk(this.permission_routes)
      // 设置分区为隐藏路由，walk 不会收集；管理员手动补充，提升可发现性
      if (this.isAdmin) {
        out.push({ type: 'nav', title: '目录配置', path: '/settings/directory', icon: 'Setting', hint: '设置' })
        out.push({ type: 'nav', title: '平台对接', path: '/settings/thirdparty', icon: 'Connection', hint: '设置' })
        out.push({ type: 'nav', title: '通知设置', path: '/settings/notification', icon: 'Bell', hint: '设置' })
        out.push({ type: 'nav', title: '定时同步', path: '/settings/sync', icon: 'Refresh', hint: '设置' })
        out.push({ type: 'nav', title: '登录安全', path: '/settings/security', icon: 'Lock', hint: '设置' })
      }
      // 去重
      const seen = new Set()
      return out.filter(i => (seen.has(i.path) ? false : seen.add(i.path)))
    },
    // 快捷操作
    actionItems() {
      return [
        { type: 'action', title: '个人中心', icon: 'User', hint: '操作', run: () => this.go('/profile/index') },
        { type: 'action', title: '退出登录', icon: 'SwitchButton', hint: '操作', run: () => this.logout() }
      ]
    },
    flat() {
      const q = this.query.trim().toLowerCase()
      const match = i => !q || i.title.toLowerCase().includes(q)
      const items = [...this.navItems.filter(match), ...this.actionItems.filter(match)]
      return items.map((it, idx) => ({ ...it, idx, key: (it.type + (it.path || it.title)) }))
    },
    grouped() {
      const nav = this.flat.filter(i => i.type === 'nav')
      const act = this.flat.filter(i => i.type === 'action')
      const g = []
      if (nav.length) g.push({ name: '页面', items: nav })
      if (act.length) g.push({ name: '操作', items: act })
      return g
    }
  },
  watch: {
    query() { this.active = 0 },
    visible(v) {
      if (v) this.$nextTick(() => this.$refs.input && this.$refs.input.focus())
    }
  },
  mounted() {
    window.addEventListener('keydown', this.onKeydown)
    window.addEventListener('open-cmdk', this.open)
  },
  beforeUnmount() {
    window.removeEventListener('keydown', this.onKeydown)
    window.removeEventListener('open-cmdk', this.open)
  },
  methods: {
    onKeydown(e) {
      if ((e.metaKey || e.ctrlKey) && (e.key === 'k' || e.key === 'K')) {
        e.preventDefault()
        this.toggle()
      }
    },
    toggle() {
      this.visible = !this.visible
      if (this.visible) { this.query = ''; this.active = 0 }
    },
    open() { this.visible = true; this.query = ''; this.active = 0 },
    close() { this.visible = false },
    move(delta) {
      const n = this.flat.length
      if (!n) return
      this.active = (this.active + delta + n) % n
      this.$nextTick(() => {
        const el = this.$refs.list && this.$refs.list.querySelector('.cmdk-item.is-active')
        if (el) el.scrollIntoView({ block: 'nearest' })
      })
    },
    run(item) {
      const target = item || this.flat[this.active]
      if (!target) return
      this.close()
      if (target.type === 'action' && target.run) target.run()
      else if (target.path) this.go(target.path)
    },
    go(path) {
      if (this.$route.path !== path) this.$router.push(path).catch(() => {})
    },
    async logout() {
      await this.$store.dispatch('user/logout')
      this.$router.push(`/login?redirect=${this.$route.fullPath}`)
    },
    resolvePath(basePath, path) {
      if (!path) return basePath || '/'
      if (path.startsWith('/')) return path
      return (basePath.replace(/\/$/, '') + '/' + path).replace(/\/+/g, '/')
    }
  }
}
</script>

<style lang="scss" scoped>
@import "@/styles/variables.scss";

.cmdk-mask {
  position: fixed;
  inset: 0;
  z-index: 3000;
  background: rgba(15, 23, 42, 0.45);
  backdrop-filter: blur(2px);
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding-top: 12vh;
}
.cmdk-panel {
  width: 600px;
  max-width: calc(100vw - 32px);
  background: #fff;
  border-radius: 14px;
  box-shadow: 0 24px 64px rgba(15, 23, 42, 0.28);
  overflow: hidden;
  border: 1px solid $slate200;
}
.cmdk-search {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 16px;
  border-bottom: 1px solid $slate100;
  &__icon { font-size: 18px; color: $slate400; }
  &__input {
    flex: 1;
    border: none;
    outline: none;
    font-size: 15px;
    color: $slate800;
    background: transparent;
    &::placeholder { color: $slate400; }
  }
}
.cmdk-kbd, .cmdk-footer kbd {
  font-family: $fontFamilyBase;
  font-size: 11px;
  color: $slate500;
  background: $slate100;
  border: 1px solid $slate200;
  border-radius: 5px;
  padding: 1px 6px;
}
.cmdk-list {
  max-height: 56vh;
  overflow-y: auto;
  padding: 8px;
}
.cmdk-group {
  font-size: 11px;
  font-weight: $fontWeightSemibold;
  color: $slate400;
  letter-spacing: 0.04em;
  padding: 8px 10px 4px;
}
.cmdk-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 9px 10px;
  border-radius: 8px;
  cursor: pointer;
  &__icon { font-size: 16px; color: $slate500; }
  &__title { flex: 1; font-size: 14px; color: $slate800; }
  &__hint { font-size: 12px; color: $slate400; }
  &.is-active {
    background: $menuActiveBg;
    .cmdk-item__icon, .cmdk-item__title { color: $themePrimary; }
  }
}
.cmdk-sep { height: 1px; background: $slate100; margin: 6px 4px; }
.cmdk-empty { padding: 28px; text-align: center; color: $slate400; font-size: 14px; }
.cmdk-footer {
  display: flex;
  gap: 18px;
  padding: 10px 16px;
  border-top: 1px solid $slate100;
  font-size: 12px;
  color: $slate400;
  kbd { margin-right: 2px; }
}

.cmdk-fade-enter-active, .cmdk-fade-leave-active { transition: opacity 0.15s ease; }
.cmdk-fade-enter-from, .cmdk-fade-leave-to { opacity: 0; }
.cmdk-fade-enter-active .cmdk-panel { transition: transform 0.15s ease; }
.cmdk-fade-enter-from .cmdk-panel { transform: translateY(-8px); }
</style>
