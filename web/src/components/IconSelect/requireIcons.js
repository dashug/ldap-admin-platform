// Vite：用 import.meta.glob 列出所有 svg 图标名（替代 webpack 的 require.context）
const modules = import.meta.glob('../../icons/svg/*.svg')
const re = /\/([^/]+)\.svg$/

const icons = Object.keys(modules).map(path => {
  const m = path.match(re)
  return m ? m[1] : ''
}).filter(Boolean)

export default icons
