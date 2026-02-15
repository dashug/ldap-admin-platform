<template>
  <div class="dashboard-radar-chart" :style="{ height: height, width: width }" />
</template>

<script>
import echarts from 'echarts'
require('echarts/theme/macarons')
import { debounce } from '@/utils'

export default {
  name: 'DashboardRadar',
  props: {
    list: {
      type: Array,
      default: () => []
    },
    height: { type: String, default: '300px' },
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
    this.chart = echarts.init(this.$el, 'macarons')
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
        title: { text: '数据分布', left: 'center', top: 8, textStyle: { fontSize: 14 } },
        tooltip: { trigger: 'item' },
        radar: {
          radius: '62%',
          center: ['50%', '48%'],
          indicator
        },
        series: [{
          name: '数量',
          type: 'radar',
          data: [{ value, name: '当前' }],
          areaStyle: { opacity: 0.3 }
        }]
      })
    }
  }
}
</script>
