<template>
  <el-dialog title="批量导入用户" v-model="visible" width="760px" :close-on-click-modal="false" @closed="reset">
    <el-alert
      type="info"
      :closable="false"
      show-icon
      class="import-alert"
      title="按模板填写后上传（支持 .xlsx / .csv）。用户名、中文名、邮箱、工号为必填；手机号留空会自动生成。导入用户默认角色为「普通用户」，初始密码为目录配置中的默认密码，可后续编辑分配部门。"
    />

    <div class="import-toolbar">
      <el-button size="small" icon="Download" @click="downloadTemplate">下载模板</el-button>
      <el-upload
        action="#"
        :auto-upload="false"
        :show-file-list="false"
        :on-change="handleFile"
        accept=".xlsx,.csv"
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
import ExcelJS from 'exceljs'
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
    async downloadTemplate() {
      const headers = ['用户名', '中文名', '花名', '邮箱', '手机号', '工号', '职位']
      const example = ['zhangsan', '张三', '三哥', 'zhangsan@example.com', '13800138000', '1001', '工程师']
      const wb = new ExcelJS.Workbook()
      const ws = wb.addWorksheet('users')
      ws.addRow(headers)
      ws.addRow(example)
      const buf = await wb.xlsx.writeBuffer()
      saveAs(new Blob([buf], { type: 'application/octet-stream' }), '用户导入模板.xlsx')
    },
    handleFile(file) {
      const raw = file.raw || file
      this.fileName = file.name || ''
      this.result = null
      const name = (this.fileName || '').toLowerCase()
      const reader = new FileReader()
      reader.onload = async e => {
        try {
          let rows
          if (name.endsWith('.csv')) {
            rows = this.parseCsv(new TextDecoder('utf-8').decode(e.target.result))
          } else if (name.endsWith('.xlsx')) {
            rows = await this.parseXlsx(e.target.result)
          } else {
            Message.error('仅支持 .xlsx 或 .csv 文件（旧版 .xls 请先另存为 .xlsx）')
            return
          }
          this.users = this.parseRows(rows)
          if (!this.users.length) Message.warning('未解析到有效数据行')
        } catch (err) {
          Message.error('解析文件失败：' + (err.message || err))
        }
      }
      reader.readAsArrayBuffer(raw)
    },
    // 用 exceljs 读取 .xlsx，转成「行数组的数组」（与旧 sheet_to_json({header:1}) 等价）
    async parseXlsx(arrayBuffer) {
      const wb = new ExcelJS.Workbook()
      await wb.xlsx.load(arrayBuffer)
      const ws = wb.worksheets[0]
      const rows = []
      if (ws) {
        ws.eachRow((row) => {
          const vals = row.values // 1-indexed，下标 0 为空
          const arr = []
          for (let i = 1; i < vals.length; i++) arr.push(this.cellText(vals[i]))
          rows.push(arr)
        })
      }
      return rows
    },
    // 提取单元格显示文本（兼容富文本 / 超链接 / 公式结果 / 日期）
    cellText(v) {
      if (v == null) return ''
      if (typeof v === 'object') {
        if (v instanceof Date) return v.toISOString()
        if (v.text != null) return String(v.text)
        if (v.result != null) return String(v.result)
        if (Array.isArray(v.richText)) return v.richText.map(t => t.text).join('')
        return String(v)
      }
      return String(v)
    },
    // 简易 CSV 解析（支持双引号包裹与 "" 转义），去除 BOM
    parseCsv(text) {
      const rows = []
      const lines = text.replace(/^﻿/, '').split(/\r\n|\n|\r/)
      for (const line of lines) {
        if (line === '') continue
        rows.push(this.splitCsvLine(line))
      }
      return rows
    },
    splitCsvLine(line) {
      const out = []
      let cur = ''
      let inQ = false
      for (let i = 0; i < line.length; i++) {
        const ch = line[i]
        if (inQ) {
          if (ch === '"') {
            if (line[i + 1] === '"') { cur += '"'; i++ } else { inQ = false }
          } else cur += ch
        } else if (ch === '"') {
          inQ = true
        } else if (ch === ',') {
          out.push(cur); cur = ''
        } else {
          cur += ch
        }
      }
      out.push(cur)
      return out
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
