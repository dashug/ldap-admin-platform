<template>
  <div class="dashboard-bar-chart" :style="{ height: height, width: width }" />
</template>

<script>
import echarts from 'echarts'
require('echarts/theme/macarons')
import { debounce } from '@/utils'

export default {
  name: 'DashboardBar',
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
      if (!this.chart || !Array.isArray(this.list) || this.list.length === 0) return
      const names = this.list.map(item => item.dataName || item.dataType)
      const values = this.list.map(item => Number(item.dataCount) || 0)
      const colors = ['#5470c6', '#91cc75', '#fac858', '#ee6666', '#73c0de', '#3ba272']
      this.chart.setOption({
        title: { text: '各模块数量', left: 'center', top: 8, textStyle: { fontSize: 14 } },
        tooltip: { trigger: 'axis' },
        grid: { top: 40, left: '3%', right: '3%', bottom: '12%', containLabel: true },
        xAxis: {
          type: 'category',
          data: names,
          axisLabel: { interval: 0, rotate: names.length > 4 ? 20 : 0 }
        },
        yAxis: { type: 'value', axisTick: { show: false } },
        series: [{
          name: '数量',
          type: 'bar',
          barWidth: '56%',
          data: values.map((v, i) => ({ value: v, itemStyle: { color: colors[i % colors.length] } }))
        }]
      })
    }
  }
}
</script>
