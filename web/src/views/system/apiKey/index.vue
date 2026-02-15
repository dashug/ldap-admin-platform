<template>
  <div>
    <el-card class="container-card" shadow="always">
      <el-form size="mini" :inline="true" class="demo-form-inline">
        <el-form-item>
          <el-button :loading="loading" icon="el-icon-plus" type="primary" @click="openCreate">新建密钥</el-button>
        </el-form-item>
      </el-form>

      <el-table v-loading="loading" :data="tableData" border stripe style="width: 100%">
        <el-table-column show-overflow-tooltip prop="name" label="名称" />
        <el-table-column show-overflow-tooltip prop="keyPrefix" label="密钥前缀" width="200">
          <template slot-scope="scope">
            <code>{{ scope.row.keyPrefix }}…</code>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="180">
          <template slot-scope="scope">
            {{ scope.row.createdAt | formatTime }}
          </template>
        </el-table-column>
        <el-table-column fixed="right" label="操作" align="center" width="100">
          <template slot-scope="scope">
            <el-popconfirm title="删除后该密钥将立即失效，确定删除？" @onConfirm="doDelete(scope.row.id)">
              <el-button slot="reference" size="mini" type="danger" icon="el-icon-delete">删除</el-button>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        :current-page="params.pageNum"
        :page-size="params.pageSize"
        :total="total"
        :page-sizes="[10, 20, 50]"
        layout="total, prev, pager, next, sizes"
        background
        style="margin-top: 10px; float: right; margin-bottom: 10px;"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />

      <!-- 新建：输入名称 -->
      <el-dialog title="新建 API 密钥" :visible.sync="createDialogVisible" width="500px">
        <el-form ref="createForm" :model="createForm" :rules="createRules" label-width="名称">
          <el-form-item label="名称" prop="name">
            <el-input v-model.trim="createForm.name" placeholder="例如：第三方系统、脚本" maxlength="64" show-word-limit />
          </el-form-item>
        </el-form>
        <span slot="footer" class="dialog-footer">
          <el-button @click="createDialogVisible = false">取消</el-button>
          <el-button :loading="createLoading" type="primary" @click="submitCreate">确定</el-button>
        </span>
      </el-dialog>

      <!-- 创建成功：仅此一次显示密钥 -->
      <el-dialog title="请妥善保存密钥" :visible.sync="keyResultVisible" width="560px" :close-on-click-modal="false">
        <el-alert type="warning" :closable="false" show-icon style="margin-bottom: 12px;">
          <span>密钥仅显示一次，关闭后将无法再次查看，请复制保存后再关闭。</span>
        </el-alert>
        <el-input v-model="createdKey" type="textarea" :rows="3" readonly>
          <template slot="prepend">X-API-Key</template>
        </el-input>
        <span slot="footer" class="dialog-footer">
          <el-button type="primary" @click="copyKey">复制</el-button>
          <el-button @click="closeKeyResult">已保存，关闭</el-button>
        </span>
      </el-dialog>
    </el-card>
  </div>
</template>

<script>
import { getApiKeyList, createApiKey, deleteApiKey } from '@/api/system/apiKey'
import { Message } from 'element-ui'

export default {
  name: 'ApiKey',
  filters: {
    formatTime(val) {
      if (!val) return '-'
      const d = new Date(val)
      return isNaN(d.getTime()) ? val : d.toLocaleString('zh-CN')
    }
  },
  data() {
    return {
      params: {
        pageNum: 1,
        pageSize: 10
      },
      tableData: [],
      total: 0,
      loading: false,

      createDialogVisible: false,
      createLoading: false,
      createForm: { name: '' },
      createRules: {
        name: [{ required: true, message: '请输入名称', trigger: 'blur' }, { min: 1, max: 64, message: '1–64 个字符', trigger: 'blur' }]
      },

      keyResultVisible: false,
      createdKey: ''
    }
  },
  mounted() {
    this.fetchList()
  },
  methods: {
    async fetchList() {
      this.loading = true
      try {
        const res = await getApiKeyList(this.params)
        this.tableData = (res.data && res.data.items) || []
        this.total = (res.data && res.data.total) || 0
      } catch (e) {
        // 权限类错误已在 request 拦截器中提示，不再重复弹出「获取列表失败」
        const msg = (e && e.message) || ''
        if (msg !== '没有权限' && msg.indexOf('权限') === -1) {
          Message.error('获取列表失败')
        }
      } finally {
        this.loading = false
      }
    },
    handleSizeChange(val) {
      this.params.pageSize = val
      this.params.pageNum = 1
      this.fetchList()
    },
    handleCurrentChange(val) {
      this.params.pageNum = val
      this.fetchList()
    },
    openCreate() {
      this.createForm.name = ''
      this.createDialogVisible = true
      this.$nextTick(() => {
        this.$refs.createForm && this.$refs.createForm.clearValidate()
      })
    },
    submitCreate() {
      this.$refs.createForm.validate(async valid => {
        if (!valid) return
        this.createLoading = true
        try {
          const res = await createApiKey({ name: this.createForm.name })
          this.createDialogVisible = false
          this.createdKey = (res.data && res.data.key) || ''
          this.keyResultVisible = true
          this.fetchList()
        } catch (_) {
          Message.error('创建失败')
        } finally {
          this.createLoading = false
        }
      })
    },
    copyKey() {
      try {
        const el = document.createElement('textarea')
        el.value = this.createdKey
        document.body.appendChild(el)
        el.select()
        document.execCommand('copy')
        document.body.removeChild(el)
        Message.success('已复制到剪贴板')
      } catch (_) {
        Message.error('复制失败')
      }
    },
    closeKeyResult() {
      this.keyResultVisible = false
      this.createdKey = ''
    },
    async doDelete(id) {
      try {
        await deleteApiKey({ id })
        Message.success('已删除')
        this.fetchList()
      } catch (_) {
        Message.error('删除失败')
      }
    }
  }
}
</script>

<style scoped>
code {
  font-size: 12px;
  color: #606266;
}
</style>
