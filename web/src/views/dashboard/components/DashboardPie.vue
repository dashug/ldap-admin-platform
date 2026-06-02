<template>
  <div class="dashboard-pie-chart" :style="{ height: height, width: width }" />
</template>

<script>
import * as echarts from 'echarts'
import { debounce } from '@/utils'

const BRAND = ['#6366f1', '#22d3ee', '#10b981', '#f59e0b', '#8b5cf6', '#f43f5e']

export default {
  name: 'DashboardPie',
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
      const data = this.list.map(item => ({
        value: Number(item.dataCount) || 0,
        name: item.dataName || item.dataType
      }))
      const total = data.reduce((s, d) => s + d.value, 0)
      this.chart.setOption({
        color: BRAND,
        tooltip: { trigger: 'item', formatter: '{b}：{c}（{d}%）' },
        title: {
          text: String(total),
          subtext: '总计',
          left: '34%',
          top: '41%',
          textAlign: 'center',
          textStyle: { fontSize: 26, fontWeight: 700, color: '#0f172a' },
          subtextStyle: { fontSize: 12, color: '#94a3b8' }
        },
        legend: {
          orient: 'vertical',
          right: '4%',
          top: 'center',
          icon: 'circle',
          itemWidth: 8,
          itemHeight: 8,
          itemGap: 12,
          textStyle: { color: '#64748b', fontSize: 12 }
        },
        series: [{
          name: '数量',
          type: 'pie',
          radius: ['56%', '78%'],
          center: ['34%', '50%'],
          avoidLabelOverlap: true,
          itemStyle: { borderColor: '#fff', borderWidth: 3, borderRadius: 6 },
          label: { show: false },
          labelLine: { show: false },
          data,
          emphasis: { scale: true, scaleSize: 5 }
        }]
      })
    }
  }
}
</script>
