import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'
import { createSvgIconsPlugin } from 'vite-plugin-svg-icons'

function resolve(dir) {
  return path.join(__dirname, dir)
}

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd())
  return {
    base: '/',
    resolve: {
      alias: {
        '@': resolve('src'),
        // 浏览器端 path 垫片（TagsView/SidebarItem/HeaderSearch 用到 path.resolve）
        path: 'path-browserify'
      },
      extensions: ['.mjs', '.js', '.vue', '.json', '.scss']
    },
    plugins: [
      vue(),
      // svg 雪碧图：保持 icon-[name] 的 symbolId 以兼容现有 <svg-icon icon-class="x">
      createSvgIconsPlugin({
        iconDirs: [resolve('src/icons/svg')],
        symbolId: 'icon-[name]'
      })
    ],
    css: {
      preprocessorOptions: {
        scss: {
          // 兼容旧版 dart-sass 写法，关闭新版弃用告警噪音
          silenceDeprecations: ['legacy-js-api', 'import', 'global-builtin', 'color-functions']
        }
      }
    },
    server: {
      port: 9527,
      open: true,
      proxy: {
        '/api': {
          target: env.VITE_APP_BASE_API || 'http://127.0.0.1:8888',
          changeOrigin: true
        }
      }
    },
    build: {
      outputDir: 'dist',
      sourcemap: false,
      chunkSizeWarningLimit: 4096
    }
  }
})
