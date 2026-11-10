import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
const apiTarget = process.env.VITE_DEV_API_TARGET || 'http://127.0.0.1:9090/'

const proxyConfig = {
  '/api': {
    target: apiTarget,
    changeOrigin: true,
    rewrite: (path: string) => path.replace(/^\/api/, '')
  }
}

export default defineConfig({
  base: '/',
  build: { manifest: true },
  server: {
    hmr: true,
    proxy: proxyConfig
  },
  preview: {
    proxy: proxyConfig
  },
  plugins: [vue()],
})
