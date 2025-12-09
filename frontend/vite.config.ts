import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@guandan/sdk-ts': path.resolve(__dirname, '../packages/sdk-ts/src'),
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    host: true, // 监听所有网络接口，Docker 中需要
    port: 5173,
    strictPort: true,
    watch: {
      usePolling: true, // Docker 中需要轮询以检测文件变化
    },
    proxy: {
      // 代理所有 /api 请求到后端
      '/api': {
        target: process.env.VITE_API_URL || 'http://localhost:8080',
        changeOrigin: true,
        secure: false,
      },
      // 代理 WebSocket 连接
      '/ws': {
        target: process.env.VITE_WS_URL || 'ws://localhost:8080',
        ws: true,
        changeOrigin: true,
      }
    }
  }
})
