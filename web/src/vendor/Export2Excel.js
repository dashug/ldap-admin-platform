/* eslint-disable */
// 使用 exceljs 生成 Excel，替代存在原型污染/ReDoS 漏洞且官方 npm 包无修复版本的 SheetJS(xlsx)。
// 仅保留项目实际使用的 export_json_to_excel。
import { saveAs } from 'file-saver'
import ExcelJS from 'exceljs'

// 依据内容估算列宽（中文字符按 2 计），与旧实现的 autoWidth 行为保持一致
function colWidthFor(rows) {
  const widths = []
  rows.forEach(row => {
    row.forEach((val, c) => {
      let w = 10
      if (val != null) {
        const s = val.toString()
        w = s.charCodeAt(0) > 255 ? s.length * 2 : s.length
      }
      widths[c] = Math.max(widths[c] || 10, w)
    })
  })
  return widths
}

export async function export_json_to_excel({
  multiHeader = [],
  header,
  data,
  filename,
  merges = [],
  autoWidth = true,
  bookType = 'xlsx'
} = {}) {
  filename = filename || 'excel-list'
  const rows = [...data]
  rows.unshift(header)
  for (let i = multiHeader.length - 1; i > -1; i--) {
    rows.unshift(multiHeader[i])
  }

  const wb = new ExcelJS.Workbook()
  const ws = wb.addWorksheet('Sheet1')
  rows.forEach(r => ws.addRow(r))

  if (autoWidth) {
    colWidthFor(rows).forEach((w, i) => { ws.getColumn(i + 1).width = w })
  }

  // merges 为 "A1:B2" 形式的区间字符串，exceljs 直接接受；非法区间忽略
  merges.forEach(range => { try { ws.mergeCells(range) } catch (e) { /* ignore */ } })

  let buf
  let mime = 'application/octet-stream'
  if (bookType === 'csv') {
    buf = await wb.csv.writeBuffer()
    mime = 'text/csv;charset=utf-8'
  } else {
    buf = await wb.xlsx.writeBuffer()
  }
  saveAs(new Blob([buf], { type: mime }), `${filename}.${bookType}`)
}
