import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'
import AutoImport from 'unplugin-auto-import/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'
import Components from 'unplugin-vue-components/vite'
import { defineConfig } from 'vite'
import UnoCSS from 'unocss/vite'
import { API_URL } from './src/stores/config';
// import vueDevTools from 'vite-plugin-vue-devtools'
export default defineConfig(({ mode }) => {
  const serverConfig: any = {
    host: '0.0.0.0',
    port: 8868,
  }

  if (mode === 'development') {
    serverConfig.proxy = {
      // 代理规则示例
      '/api/': {
        target: API_URL, // 目标服务器地址
        changeOrigin: true, // 是否改变请求的源
      },
    }
  }
  return {
    plugins: [
      vue(),
      UnoCSS(),
      // vueDevTools(),
      AutoImport({
        resolvers: [ElementPlusResolver()],
      }),
      Components({
        resolvers: [ElementPlusResolver()],
      }),
    ],
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url)),
      },
    },
    server: {
      ...serverConfig,
      sourcemap: 'inline', // 使用内联源映射以减少外部文件请求
    },
    build: {
      sourcemap: false, // 禁用构建时的源映射
      outDir: '../server/public/os',
    },
    //base: '/os/',
  }
})
//export default defineConfig()
