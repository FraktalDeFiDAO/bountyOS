import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'

const apiTarget = process.env.VITE_API_TARGET || 'http://localhost:12496'
const wsTarget =
  process.env.VITE_WS_TARGET ||
  apiTarget.replace(/^http/, 'ws')

export default defineConfig({
  plugins: [vue(), tailwindcss()],
  server: {
    port: 13440,
    proxy: {
      '/api': apiTarget,
      '/ws': {
        target: wsTarget,
        ws: true
      }
    }
  }
})
