<template>
  <div class="dashboard-pie-chart" :style="{ height: height, width: width }" />
</template>

<script>
import echarts from 'echarts'
require('echarts/theme/macarons')
import { debounce } from '@/utils'

export default {
  name: 'DashboardPie',
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
      const data = this.list.map(item => ({
        value: Number(item.dataCount) || 0,
        name: item.dataName || item.dataType
      }))
      this.chart.setOption({
        title: { text: '数据占比', left: 'center', top: 8, textStyle: { fontSize: 14 } },
        tooltip: { trigger: 'item', formatter: '{b}：{c}（{d}%）' },
        legend: { left: 'center', bottom: 8, data: data.map(d => d.name) },
        series: [{
          name: '数量',
          type: 'pie',
          radius: ['40%', '70%'],
          center: ['50%', '45%'],
          data,
          emphasis: { itemStyle: { shadowBlur: 10, shadowOffsetX: 0, shadowColor: 'rgba(0,0,0,0.2)' } }
        }]
      })
    }
  }
}
</script>
