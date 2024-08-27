import { defineConfig } from "vite";
import path from 'path'
import vue from "@vitejs/plugin-vue";
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'


export default defineConfig(async () => ({
  plugins: [
    vue(),
    AutoImport({
      imports: [
        // presets
        'vue',
        'vue-router',
      ],
      resolvers: [ElementPlusResolver()],
    }),
    Components({
      resolvers: [ElementPlusResolver()],
    }),
  ],

  clearScreen: false,
  server: {
    port: 8215
  },
  resolve: {
    alias: {
      "@": path.resolve(__dirname, 'src'),
      "~": path.resolve(__dirname, 'wailsjs')
    },
  },
  build:{
    sourcemap: true,
    // 打包环境移除console.log，debugger
    terserOptions: {
      compress: {
        drop_console: true,
        drop_debugger: true
      }
    },
    outDir: '../godo/deps/dist',
  }
}));
