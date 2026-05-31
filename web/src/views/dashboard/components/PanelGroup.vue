<template>
  <el-row :gutter="16" class="panel-group">
    <el-col v-for="(item, i) in listData" :key="item.dataType" :xs="12" :sm="8" :lg="4" class="card-panel-col">
      <div class="stat-card" @click="handleSetLineChartData(item.dataType)">
        <div class="stat-card__icon" :style="{ background: palette[i % palette.length].bg, color: palette[i % palette.length].fg }">
          <svg-icon :icon-class="item.icon" class-name="stat-card__svg" />
        </div>
        <div class="stat-card__body">
          <div class="stat-card__label">{{ item.dataName }}</div>
          <count-to :start-val="0" :end-val="item.dataCount" :duration="2000" class="stat-card__num" />
        </div>
      </div>
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
        { bg: 'rgba(79, 70, 229, 0.10)', fg: '#4f46e5' },
        { bg: 'rgba(14, 165, 233, 0.10)', fg: '#0ea5e9' },
        { bg: 'rgba(16, 185, 129, 0.12)', fg: '#10b981' },
        { bg: 'rgba(245, 158, 11, 0.14)', fg: '#f59e0b' },
        { bg: 'rgba(139, 92, 246, 0.12)', fg: '#8b5cf6' },
        { bg: 'rgba(244, 63, 94, 0.10)', fg: '#f43f5e' }
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
    handleSetLineChartData(type) {
      this.$emit('handleSetLineChartData', type)
    }
  }
}
</script>

<style lang="scss" scoped>
@import "@/styles/variables.scss";

.panel-group {
  margin-bottom: 8px;

  .card-panel-col {
    margin-bottom: 16px;
  }
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 18px;
  background: #fff;
  border-radius: 14px;
  border: 1px solid $borderColor;
  box-shadow: $cardShadow;
  cursor: pointer;
  transition: box-shadow $transitionBase, transform $transitionBase, border-color $transitionBase;

  &:hover {
    box-shadow: $cardShadowHover;
    transform: translateY(-3px);
    border-color: transparent;
  }

  &__icon {
    flex: none;
    width: 48px;
    height: 48px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 13px;
  }
  &__svg { font-size: 24px; }

  &__body { min-width: 0; }
  &__label {
    font-size: 13px;
    color: $slate500;
    font-weight: $fontWeightMedium;
    margin-bottom: 4px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  &__num {
    font-size: 26px;
    font-weight: $fontWeightBold;
    color: $slate900;
    letter-spacing: -0.02em;
    font-variant-numeric: tabular-nums;
    line-height: 1.1;
  }
}
</style>
