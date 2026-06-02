<template>
  <div class="dashboard-radar-chart" :style="{ height: height, width: width }" />
</template>

<script>
import * as echarts from 'echarts'
import { debounce } from '@/utils'

export default {
  name: 'DashboardRadar',
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
      if (!this.chart || !Array.isArray(this.list) || this.list.length === 0) {
        return
      }
      const maxVal = Math.max(...this.list.map(i => Number(i.dataCount) || 0), 1)
      const indicator = this.list.map(item => ({
        name: item.dataName || item.dataType,
        max: maxVal
      }))
      const value = this.list.map(item => Number(item.dataCount) || 0)
      this.chart.setOption({
        tooltip: { trigger: 'item' },
        radar: {
          radius: '66%',
          center: ['50%', '52%'],
          indicator,
          axisName: { color: '#64748b', fontSize: 12 },
          splitNumber: 4,
          splitLine: { lineStyle: { color: '#e2e8f0' } },
          splitArea: { areaStyle: { color: ['rgba(99,102,241,0.03)', 'rgba(99,102,241,0.07)'] } },
          axisLine: { lineStyle: { color: '#e2e8f0' } }
        },
        series: [{
          name: '数量',
          type: 'radar',
          data: [{ value, name: '当前' }],
          symbol: 'circle',
          symbolSize: 5,
          lineStyle: { color: '#6366f1', width: 2 },
          itemStyle: { color: '#6366f1' },
          areaStyle: { color: 'rgba(99,102,241,0.18)' }
        }]
      })
    }
  }
}
</script>
