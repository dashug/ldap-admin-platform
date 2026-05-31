<template>
  <el-dialog title="批量导入用户" v-model="visible" width="760px" :close-on-click-modal="false" @closed="reset">
    <el-alert
      type="info"
      :closable="false"
      show-icon
      class="import-alert"
      title="按模板填写后上传（支持 .xlsx / .xls / .csv）。用户名、中文名、邮箱、工号为必填；手机号留空会自动生成。导入用户默认角色为「普通用户」，初始密码为目录配置中的默认密码，可后续编辑分配部门。"
    />

    <div class="import-toolbar">
      <el-button size="small" icon="Download" @click="downloadTemplate">下载模板</el-button>
      <el-upload
        action="#"
        :auto-upload="false"
        :show-file-list="false"
        :on-change="handleFile"
        accept=".xlsx,.xls,.csv"
      >
        <el-button size="small" type="primary" icon="Upload">选择文件</el-button>
      </el-upload>
      <span v-if="fileName" class="import-filename">{{ fileName }}（共 {{ users.length }} 行）</span>
    </div>

    <el-table v-if="users.length" :data="users" size="small" border max-height="320" style="width: 100%; margin-top: 12px;">
      <el-table-column type="index" label="#" width="48" />
      <el-table-column prop="username" label="用户名" width="120" show-overflow-tooltip />
      <el-table-column prop="nickname" label="中文名" width="100" show-overflow-tooltip />
      <el-table-column prop="givenName" label="花名" width="90" show-overflow-tooltip />
      <el-table-column prop="mail" label="邮箱" min-width="160" show-overflow-tooltip />
      <el-table-column prop="mobile" label="手机号" width="120" show-overflow-tooltip />
      <el-table-column prop="jobNumber" label="工号" width="90" show-overflow-tooltip />
      <el-table-column prop="position" label="职位" width="100" show-overflow-tooltip />
    </el-table>

    <div v-if="result" class="import-result">
      <el-alert
        :type="result.failCount === 0 ? 'success' : 'warning'"
        :closable="false"
        show-icon
        :title="(result.dryRun ? '预检结果' : '导入结果') + `：成功 ${result.okCount} 条，失败 ${result.failCount} 条`"
      />
      <el-table v-if="failedRows.length" :data="failedRows" size="small" border max-height="220" style="width: 100%; margin-top: 8px;">
        <el-table-column prop="row" label="行号" width="64" align="center" />
        <el-table-column prop="username" label="用户名" width="140" show-overflow-tooltip />
        <el-table-column prop="message" label="原因" min-width="220" show-overflow-tooltip />
      </el-table>
    </div>

    <template #footer><div class="dialog-footer">
      <el-button size="small" @click="visible = false">取 消</el-button>
      <el-button size="small" :disabled="!users.length" :loading="checking" @click="submit(true)">预 检</el-button>
      <el-button size="small" type="primary" :disabled="!users.length" :loading="importing" @click="submit(false)">确认导入</el-button>
    </div></template>
  </el-dialog>
</template>

<script>
import XLSX from 'xlsx'
import { saveAs } from 'file-saver'
import { batchImportUsers } from '@/api/system/base'
import { ElMessage as Message } from 'element-plus'

const HEADER_MAP = {
  '用户名': 'username',
  '中文名': 'nickname',
  '花名': 'givenName',
  '邮箱': 'mail',
  '手机号': 'mobile',
  '工号': 'jobNumber',
  '职位': 'position'
}

export default {
  name: 'UserImportDialog',
  emits: ['done'],
  data() {
    return {
      visible: false,
      fileName: '',
      users: [],
      checking: false,
      importing: false,
      result: null
    }
  },
  computed: {
    failedRows() {
      if (!this.result || !this.result.results) return []
      return this.result.results.filter(r => !r.ok)
    }
  },
  methods: {
    open() {
      this.reset()
      this.visible = true
    },
    reset() {
      this.fileName = ''
      this.users = []
      this.result = null
      this.checking = false
      this.importing = false
    },
    downloadTemplate() {
      const headers = ['用户名', '中文名', '花名', '邮箱', '手机号', '工号', '职位']
      const example = ['zhangsan', '张三', '三哥', 'zhangsan@example.com', '13800138000', '1001', '工程师']
      const ws = XLSX.utils.aoa_to_sheet([headers, example])
      const wb = XLSX.utils.book_new()
      XLSX.utils.book_append_sheet(wb, ws, 'users')
      const out = XLSX.write(wb, { bookType: 'xlsx', type: 'array' })
      saveAs(new Blob([out], { type: 'application/octet-stream' }), '用户导入模板.xlsx')
    },
    handleFile(file) {
      const raw = file.raw || file
      this.fileName = file.name || ''
      this.result = null
      const reader = new FileReader()
      reader.onload = e => {
        try {
          const wb = XLSX.read(new Uint8Array(e.target.result), { type: 'array' })
          const ws = wb.Sheets[wb.SheetNames[0]]
          const rows = XLSX.utils.sheet_to_json(ws, { header: 1, defval: '' })
          this.users = this.parseRows(rows)
          if (!this.users.length) Message.warning('未解析到有效数据行')
        } catch (err) {
          Message.error('解析文件失败：' + (err.message || err))
        }
      }
      reader.readAsArrayBuffer(raw)
    },
    parseRows(rows) {
      if (!rows || rows.length < 2) return []
      const header = rows[0].map(h => String(h).trim())
      const idx = {}
      header.forEach((h, i) => {
        const key = HEADER_MAP[h]
        if (key) idx[key] = i
      })
      const out = []
      for (let i = 1; i < rows.length; i++) {
        const row = rows[i]
        if (!row || row.every(c => String(c).trim() === '')) continue
        const get = k => (idx[k] != null ? String(row[idx[k]] == null ? '' : row[idx[k]]).trim() : '')
        out.push({
          username: get('username'),
          nickname: get('nickname'),
          givenName: get('givenName'),
          mail: get('mail'),
          mobile: get('mobile'),
          jobNumber: get('jobNumber'),
          position: get('position')
        })
      }
      return out
    },
    async submit(dryRun) {
      if (!this.users.length) return
      if (dryRun) this.checking = true
      else this.importing = true
      try {
        const res = await batchImportUsers({ dryRun, users: this.users })
        this.result = res.data || null
        if (!dryRun && this.result) {
          if (this.result.failCount === 0) {
            Message.success(`导入完成：成功 ${this.result.okCount} 条`)
          } else {
            Message.warning(`导入完成：成功 ${this.result.okCount} 条，失败 ${this.result.failCount} 条`)
          }
          this.$emit('done')
        }
      } catch (e) {
        // 错误由请求拦截器统一提示
      } finally {
        this.checking = false
        this.importing = false
      }
    }
  }
}
</script>

<style scoped>
.import-alert {
  margin-bottom: 14px;
}
.import-toolbar {
  display: flex;
  align-items: center;
  gap: 12px;
}
.import-filename {
  font-size: 13px;
  color: #606266;
}
.import-result {
  margin-top: 14px;
}
.dialog-footer {
  text-align: right;
}
.dialog-footer .el-button {
  margin-left: 10px;
}
</style>
