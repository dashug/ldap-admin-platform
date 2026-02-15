<template>
  <el-row :gutter="40" class="panel-group">
    <el-col v-for="item in listData" :key="item.dataType" :xs="8" :sm="8" :lg="8" class="card-panel-col">
      <a :href="item.path">
      <div class="card-panel" @click="handleSetLineChartData(item.dataType)">
        <div class="card-panel-icon-wrapper icon-people">
          <svg-icon :icon-class="item.icon" class-name="card-panel-icon" />
        </div>
        <div class="card-panel-description">
          <div class="card-panel-text">
              {{ item.dataName }}
          </div>
          <count-to :start-val="0" :end-val="item.dataCount" :duration="2600" class="card-panel-num" />
        </div>
      </div>
    </a>
    </el-col>
  </el-row>
</template>

<script>
import CountTo from 'vue-count-to'
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
      localData: []
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
@import "~@/styles/variables.scss";

.panel-group {
  margin-top: 0;
  margin-bottom: 24px;

  .card-panel-col {
    margin-bottom: 20px;
  }

  .card-panel {
    height: 120px;
    cursor: pointer;
    font-size: $fontSizeSmall;
    position: relative;
    overflow: hidden;
    color: $slate600;
    background: #fff;
    border-radius: $cardRadius;
    border: 1px solid $borderColor;
    box-shadow: $cardShadow;
    transition: box-shadow $transitionBase, transform $transitionBase, border-color $transitionBase;

    &:hover {
      box-shadow: $cardShadowHover;
      transform: translateY(-4px);
      border-color: $slate300;

      .card-panel-icon-wrapper {
        color: #fff;
      }
      .icon-people { background: $themePrimary; }
      .icon-message { background: #0ea5e9; }
      .icon-money { background: $themeDanger; }
      .icon-shopping { background: $themeSuccess; }
    }

    .icon-people { color: $themePrimary; }
    .icon-message { color: #0ea5e9; }
    .icon-money { color: $themeDanger; }
    .icon-shopping { color: $themeSuccess; }

    .card-panel-icon-wrapper {
      float: left;
      margin: 20px 0 0 20px;
      padding: 18px;
      transition: all $transitionBase;
      border-radius: 12px;
    }

    .card-panel-icon {
      float: left;
      font-size: 44px;
    }

    .card-panel-description {
      float: right;
      font-weight: $fontWeightSemibold;
      margin: 28px 24px 28px 0;

      .card-panel-text {
        line-height: 1.4;
        color: $slate500;
        font-size: 15px;
        margin-bottom: 10px;
      }

      .card-panel-num {
        font-size: 24px;
        font-weight: 700;
        color: $slate800;
        letter-spacing: -0.02em;
      }
    }
  }
}

@media (max-width:550px) {
  .card-panel-description {
    display: none;
  }

  .card-panel-icon-wrapper {
    float: none !important;
    width: 100%;
    height: 100%;
    margin: 0 !important;

    .svg-icon {
      display: block;
      margin: 14px auto !important;
      float: none !important;
    }
  }
}
</style>
