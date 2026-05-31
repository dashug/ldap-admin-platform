<template>
  <el-color-picker
    v-model="theme"
    :predefine="['#4f46e5', '#409EFF', '#1890ff', '#304156', '#11a983', '#13c2c2', '#6959CD', '#f5222d']"
    class="theme-picker"
    popper-class="theme-picker-dropdown"
  />
</template>

<script>
// Element Plus 通过 CSS 变量控制主题，这里在选色后动态设置 --el-color-primary 及其浅/深变体。
export default {
  data() {
    return { theme: '' }
  },
  computed: {
    defaultTheme() {
      return this.$store.state.settings.theme
    }
  },
  watch: {
    defaultTheme: {
      handler(val) {
        this.theme = val
      },
      immediate: true
    },
    theme(val) {
      if (typeof val !== 'string' || !val) return
      this.applyTheme(val)
      this.$emit('change', val)
    }
  },
  mounted() {
    if (this.theme) this.applyTheme(this.theme)
  },
  methods: {
    applyTheme(color) {
      const el = document.documentElement
      el.style.setProperty('--el-color-primary', color)
      for (let i = 1; i <= 9; i++) {
        el.style.setProperty(`--el-color-primary-light-${i}`, this.mix(color, '#ffffff', i / 10))
      }
      el.style.setProperty('--el-color-primary-dark-2', this.mix(color, '#000000', 0.2))
    },
    mix(color1, color2, weight) {
      const c1 = this.parse(color1)
      const c2 = this.parse(color2)
      const r = Math.round(c1[0] * (1 - weight) + c2[0] * weight)
      const g = Math.round(c1[1] * (1 - weight) + c2[1] * weight)
      const b = Math.round(c1[2] * (1 - weight) + c2[2] * weight)
      return `#${this.hex(r)}${this.hex(g)}${this.hex(b)}`
    },
    parse(c) {
      c = c.replace('#', '')
      if (c.length === 3) c = c.split('').map(x => x + x).join('')
      return [parseInt(c.slice(0, 2), 16), parseInt(c.slice(2, 4), 16), parseInt(c.slice(4, 6), 16)]
    },
    hex(n) {
      return n.toString(16).padStart(2, '0')
    }
  }
}
</script>

<style>
.theme-message,
.theme-picker-dropdown {
  z-index: 99999 !important;
}

.theme-picker .el-color-picker__trigger {
  height: 26px !important;
  width: 26px !important;
  padding: 2px;
}

.theme-picker-dropdown .el-color-dropdown__link-btn {
  display: none;
}
</style>
