<template>
  <div class="dashboard-bar-chart" :style="{ height: height, width: width }" />
</template>

<script>
import * as echarts from 'echarts'
import { debounce } from '@/utils'

const BRAND = ['#6366f1', '#22d3ee', '#10b981', '#f59e0b', '#8b5cf6', '#f43f5e']

export default {
  name: 'DashboardBar',
  props: {
    list: {
      type: Array,
      default: () => []
    },
    height: { type: String, default: '280px' },
    width: { type: String, default: '100%' }
  },
  data() {
    return { chart: null }
  },
  watch: {
    list: {
      handler() {
        this.updateChart()
      },
      deep: true
    }
  },
  mounted() {
    this.chart = echarts.init(this.$el)
    this.updateChart()
    this.__resizeHandler = debounce(() => {
      if (this.chart) this.chart.resize()
    }, 100)
    window.addEventListener('resize', this.__resizeHandler)
  },
  beforeDestroy() {
    window.removeEventListener('resize', this.__resizeHandler)
    if (this.chart) {
      this.chart.dispose()
      this.chart = null
    }
  },
  methods: {
    updateChart() {
      if (!this.chart || !Array.isArray(this.list) || this.list.length === 0) return
      // 各模块数量跨度大（如日志 >> 角色），用「横向排行榜」呈现：
      // 按数量升序（ECharts 类目轴自下而上 → 最大在顶部），barMinHeight 保证小值也有可见短条，
      // 颜色按模块身份固定（与环形图一致），条末显示真实数值。
      const items = this.list
        .map((item, i) => ({
          name: item.dataName || item.dataType,
          value: Number(item.dataCount) || 0,
          color: BRAND[i % BRAND.length]
        }))
        .sort((a, b) => a.value - b.value)
      this.chart.setOption({
        tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
        grid: { top: 8, left: 4, right: 44, bottom: 4, containLabel: true },
        xAxis: {
          type: 'value',
          axisLine: { show: false },
          axisTick: { show: false },
          splitLine: { lineStyle: { color: '#f1f5f9' } },
          axisLabel: { color: '#94a3b8' }
        },
        yAxis: {
          type: 'category',
          data: items.map(i => i.name),
          axisTick: { show: false },
          axisLine: { lineStyle: { color: '#e2e8f0' } },
          axisLabel: { color: '#475569', fontSize: 12, fontWeight: 500 }
        },
        series: [{
          name: '数量',
          type: 'bar',
          barWidth: 14,
          barMinHeight: 6,
          label: { show: true, position: 'right', color: '#475569', fontSize: 11, fontWeight: 600 },
          data: items.map(i => ({
            value: i.value,
            itemStyle: { color: i.color, borderRadius: [0, 6, 6, 0] }
          }))
        }]
      })
    }
  }
}
</script>
