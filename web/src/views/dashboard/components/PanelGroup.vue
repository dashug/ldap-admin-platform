<template>
  <el-row :gutter="16" class="panel-group">
    <el-col v-for="(item, i) in listData" :key="item.dataType" :xs="12" :sm="8" :lg="4" class="panel-col">
      <router-link :to="toPath(item.path)" class="stat" :style="accentVars(i)" :aria-label="`${item.dataName}：${item.dataCount}，点击查看详情`">
        <span class="stat__glow" aria-hidden="true" />
        <div class="stat__top">
          <span class="stat__icon"><svg-icon :icon-class="item.icon" class-name="stat__svg" /></span>
          <el-icon class="stat__arrow"><Right /></el-icon>
        </div>
        <count-to :start-val="0" :end-val="item.dataCount" :duration="1600" class="stat__num" />
        <div class="stat__label">{{ item.dataName }}</div>
      </router-link>
    </el-col>
  </el-row>
</template>

<script>
import CountTo from '@/components/CountTo/index.vue'
import { getDash } from '@/api/dashboards/dashboard'
export default {
  components: {
    CountTo
  },
  props: {
    dataInfo: {
      type: Array,
      default: null
    }
  },
  data() {
    return {
      localData: [],
      palette: [
        { fg: '#4f46e5', bg: 'rgba(79, 70, 229, 0.12)' },
        { fg: '#0ea5e9', bg: 'rgba(14, 165, 233, 0.12)' },
        { fg: '#10b981', bg: 'rgba(16, 185, 129, 0.14)' },
        { fg: '#f59e0b', bg: 'rgba(245, 158, 11, 0.16)' },
        { fg: '#8b5cf6', bg: 'rgba(139, 92, 246, 0.14)' },
        { fg: '#f43f5e', bg: 'rgba(244, 63, 94, 0.12)' }
      ]
    }
  },
  computed: {
    listData() {
      if (Array.isArray(this.dataInfo) && this.dataInfo.length) return this.dataInfo
      return this.localData || []
    }
  },
  created() {
    if (!Array.isArray(this.dataInfo) || this.dataInfo.length === 0) {
      this.getDashInfo()
    }
  },
  methods: {
    async getDashInfo() {
      try {
        const { data } = await getDash()
        this.localData = Array.isArray(data) ? data : []
      } finally {}
    },
    accentVars(i) {
      const c = this.palette[i % this.palette.length]
      return { '--accent': c.fg, '--accent-bg': c.bg }
    },
    toPath(p) {
      return (p || '').replace(/^#/, '') || '/'
    }
  }
}
</script>

<style lang="scss" scoped>
@import "@/styles/variables.scss";

.panel-group {
  margin-bottom: 8px;
}
.panel-col {
  margin-bottom: 16px;
}

.stat {
  position: relative;
  display: block;
  overflow: hidden;
  padding: 18px 18px 16px;
  background: #fff;
  border-radius: 16px;
  border: 1px solid $borderColor;
  box-shadow: $cardShadow;
  text-decoration: none;
  transition: box-shadow $transitionBase, transform $transitionBase, border-color $transitionBase;

  &__glow {
    position: absolute;
    top: -36px;
    right: -36px;
    width: 96px;
    height: 96px;
    border-radius: 50%;
    background: var(--accent-bg);
    opacity: 0.7;
    transition: transform $transitionBase, opacity $transitionBase;
    pointer-events: none;
  }
  &__top {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 14px;
  }
  &__icon {
    width: 44px;
    height: 44px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    border-radius: 12px;
    color: var(--accent);
    background: var(--accent-bg);
  }
  &__svg { font-size: 22px; }
  &__arrow {
    color: var(--accent);
    font-size: 16px;
    opacity: 0;
    transform: translateX(-4px);
    transition: opacity $transitionBase, transform $transitionBase;
  }
  &__num {
    position: relative;
    display: block;
    font-size: 28px;
    font-weight: $fontWeightBold;
    color: $slate900;
    letter-spacing: -0.02em;
    line-height: 1.1;
    font-variant-numeric: tabular-nums;
  }
  &__label {
    position: relative;
    margin-top: 4px;
    font-size: 13px;
    color: $slate500;
    font-weight: $fontWeightMedium;
  }

  &:hover {
    transform: translateY(-3px);
    box-shadow: $cardShadowHover;
    border-color: color-mix(in srgb, var(--accent) 35%, transparent);
    .stat__glow { transform: scale(1.25); opacity: 1; }
    .stat__arrow { opacity: 1; transform: translateX(0); }
  }
  // 键盘可达性：为聚焦的指标卡提供清晰焦点环
  &:focus-visible {
    outline: 2px solid var(--accent);
    outline-offset: 2px;
    box-shadow: $cardShadowHover;
  }
}

@media (prefers-reduced-motion: reduce) {
  .stat, .stat__glow, .stat__arrow { transition: none; }
  .stat:hover { transform: none; }
}
</style>
