<template>
  <span>{{ displayValue }}</span>
</template>

<script>
// 轻量 Vue 3 数字滚动组件，替代已停更的 vue-count-to（兼容其常用 props）
export default {
  name: 'CountTo',
  props: {
    startVal: { type: Number, default: 0 },
    endVal: { type: Number, default: 0 },
    duration: { type: Number, default: 2000 },
    autoplay: { type: Boolean, default: true },
    decimals: { type: Number, default: 0 },
    decimal: { type: String, default: '.' },
    separator: { type: String, default: ',' },
    prefix: { type: String, default: '' },
    suffix: { type: String, default: '' }
  },
  data() {
    return { current: this.startVal, rafId: null }
  },
  computed: {
    displayValue() {
      return this.prefix + this.format(this.current) + this.suffix
    }
  },
  watch: {
    endVal() {
      if (this.autoplay) this.start()
    }
  },
  mounted() {
    if (this.autoplay) this.start()
  },
  beforeUnmount() {
    if (this.rafId) cancelAnimationFrame(this.rafId)
  },
  methods: {
    start() {
      const startTime = performance.now()
      const from = Number(this.startVal) || 0
      const to = Number(this.endVal) || 0
      const dur = this.duration
      const tick = (now) => {
        const progress = Math.min((now - startTime) / dur, 1)
        // easeOutQuad
        const eased = progress * (2 - progress)
        this.current = from + (to - from) * eased
        if (progress < 1) {
          this.rafId = requestAnimationFrame(tick)
        } else {
          this.current = to
        }
      }
      if (this.rafId) cancelAnimationFrame(this.rafId)
      this.rafId = requestAnimationFrame(tick)
    },
    format(num) {
      const n = Number(num).toFixed(this.decimals)
      const [int, dec] = n.split('.')
      const withSep = int.replace(/\B(?=(\d{3})+(?!\d))/g, this.separator)
      return dec ? withSep + this.decimal + dec : withSep
    }
  }
}
</script>
