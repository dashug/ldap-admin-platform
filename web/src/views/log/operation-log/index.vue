<template>
  <div>
    <page-header title="操作日志" subtitle="记录关键操作，便于审计与排查">
      <template #actions>
        <el-button plain type="danger" icon="Delete" @click="handleClean">清空日志</el-button>
      </template>
    </page-header>

    <el-card class="container-card" shadow="never">
      <div class="filter-bar">
        <el-input v-model.trim="params.username" prefix-icon="Search" clearable placeholder="请求人" style="width: 150px;" @keyup.enter="search" @clear="search" />
        <el-input v-model.trim="params.ip" clearable placeholder="IP 地址" style="width: 140px;" @keyup.enter="search" @clear="search" />
        <el-input v-model.trim="params.path" clearable placeholder="请求路径" style="width: 180px;" @keyup.enter="search" @clear="search" />
        <el-select v-model="params.method" placeholder="请求方式" clearable style="width: 120px;" @change="search" @clear="search">
          <el-option v-for="item in RequestList" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
        <el-input v-model.trim="params.status" clearable placeholder="状态码" style="width: 110px;" @keyup.enter="search" @clear="search" />
        <el-button :loading="loading" icon="Search" @click="search">查询</el-button>
        <div class="filter-bar__spacer" />
        <el-button :disabled="multipleSelection.length === 0" :loading="loading" plain type="danger" icon="Delete" @click="batchDelete">批量删除</el-button>
      </div>

      <el-table
        v-loading="loading"
        :data="tableData"
        border
        stripe
        style="width: 100%"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" align="center" />
        <el-table-column show-overflow-tooltip sortable prop="username" label="请求人" />
        <el-table-column show-overflow-tooltip sortable prop="ip" label="IP地址" />
        <el-table-column show-overflow-tooltip sortable prop="path" label="请求路径" />
        <el-table-column show-overflow-tooltip sortable prop="method" label="请求方式" align="center">
          <template #default="scope">
            <el-tag v-if="scope.row.method === 'GET'" type="success">GET</el-tag>
            <el-tag v-else-if="scope.row.method === 'POST'" type="warning">POST</el-tag>
            <el-tag v-else-if="scope.row.method === 'PUT'" type="primary">PUT</el-tag>
            <el-tag v-else-if="scope.row.method === 'DELETE'" type="danger">DELETE</el-tag>
            <el-tag v-else type="info">{{ scope.row.method }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip sortable prop="status" label="请求状态" align="center">
          <template #default="scope">
            <el-tag size="small" :type="statusTagFilter(scope.row.status)" disable-transitions>{{ scope.row.status
            }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip sortable prop="startTime" label="发起时间">
          <!-- <template #default="scope">
            {{ parseGoTime(scope.row.startTime) }}
          </template> -->
        </el-table-column>
        <el-table-column show-overflow-tooltip sortable prop="timeCost" label="请求耗时(ms)" align="center">
          <template #default="scope">
            <el-tag size="small" :type="timeCostTagFilter(scope.row.timeCost)" disable-transitions>{{ scope.row.timeCost
            }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column show-overflow-tooltip sortable prop="desc" label="说明" />
        <el-table-column fixed="right" label="操作" align="center" width="80">
          <template #default="scope">
            <el-tooltip content="删除" effect="dark" placement="top">
              <el-popconfirm title="确定删除吗？" @confirm="singleDelete(scope.row.ID)">
                <template #reference><el-button size="small" icon="Delete" circle type="danger"  /></template>
              </el-popconfirm>
            </el-tooltip>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        :current-page="params.pageNum"
        :page-size="params.pageSize"
        :total="total"
        :page-sizes="[1, 5, 10, 30]"
        layout="total, prev, pager, next, sizes"
        background
        style="margin-top: 10px;float:right;margin-bottom: 10px;"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </el-card>
  </div>
</template>

<script>
import { getOperationLogs, batchDeleteOperationLogByIds, CleanOperationLog } from '@/api/log/operationLog'
import { parseGoTime } from '@/utils/index'
import { ElMessage as Message } from 'element-plus'
import PageHeader from '@/components/PageHeader/index.vue'

export default {
  name: 'OperationLog',
  components: { PageHeader },
  data() {
    return {
      // 查询参数
      params: {
        username: '',
        ip: '',
        path: '',
        status: '',
        pageNum: 1,
        pageSize: 10
      },
      // 表格数据
      tableData: [],
      total: 0,
      loading: false,

      // 删除按钮弹出框
      popoverVisible: false,
      // 表格多选
      multipleSelection: [],
      RequestList: [{
        value: 'GET',
        label: 'GET'
      }, {
        value: 'POST',
        label: 'POST'
      }, {
        value: 'DELETE',
        label: 'DELETE'
      }, {
        value: 'PUT',
        label: 'PUT'
      }]
    }
  },
  created() {
    this.getTableData()
  },
  methods: {
    statusTagFilter(val) {
      if (val === 200) {
        return 'success'
      } else if (val === 400) {
        return 'warning'
      } else if (val === 401 || val === 403 || val === 500) {
        return 'danger'
      } else {
        return 'info'
      }
    },
    timeCostTagFilter(val) {
      if (val <= 200) {
        return 'success'
      } else if (val > 200 && val <= 1000) {
        return ''
      } else if (val > 1000 && val <= 2000) {
        return 'warning'
      } else {
        return 'danger'
      }
    },
    parseGoTime,
    // 查询
    search() {
      this.params.pageNum = 1
      this.getTableData()
    },

    // 获取表格数据
    async getTableData() {
      this.loading = true
      try {
        const { data } = await getOperationLogs(this.params)
        this.tableData = data.logs
        this.total = data.total
      } finally {
        this.loading = false
      }
    },

    // 判断结果
    judgeResult(res) {
      if (res.code === 0) {
        Message({
          showClose: true,
          message: '操作成功',
          type: 'success'
        })
      }
    },
    // 清空日志
    handleClean() {
      this.$confirm('此操作将永久删除, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async res => {
        this.loading = true
        try {
          await CleanOperationLog().then(res => {
            this.judgeResult(res)
          })
        } finally {
          this.loading = false
        }
        this.getTableData()
      }).catch(() => {
        Message({
          showClose: true,
          type: 'info',
          message: '已取消删除'
        })
      })
    },

    // 批量删除
    batchDelete() {
      this.$confirm('此操作将永久删除, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async res => {
        this.loading = true
        const operationLogIds = []
        this.multipleSelection.forEach(x => {
          operationLogIds.push(x.ID)
        })
        try {
          await batchDeleteOperationLogByIds({ operationLogIds: operationLogIds }).then(res => {
            this.judgeResult(res)
          })
        } finally {
          this.loading = false
        }
        this.getTableData()
      }).catch(() => {
        Message({
          showClose: true,
          type: 'info',
          message: '已取消删除'
        })
      })
    },

    // 表格多选
    handleSelectionChange(val) {
      this.multipleSelection = val
    },

    // 单个删除
    async singleDelete(Id) {
      this.loading = true
      try {
        await batchDeleteOperationLogByIds({ operationLogIds: [Id] }).then(res => {
          this.judgeResult(res)
        })
      } finally {
        this.loading = false
      }
      this.getTableData()
    },

    // 分页
    handleSizeChange(val) {
      this.params.pageSize = val
      this.getTableData()
    },
    handleCurrentChange(val) {
      this.params.pageNum = val
      this.getTableData()
    }
  }
}
</script>

<style scoped>
.container-card {
  margin: 10px;
  margin-bottom: 100px;
}

.delete-popover {
  margin-left: 10px;
}
</style>
