import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  base: './',
  build: {
    outDir: '../../desktop/frontend/public/baiban', // 指定打包输出目录，默认就是'dist'，这里可以根据需要修改
  },
})
