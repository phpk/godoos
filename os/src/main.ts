import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'virtual:uno.css'
import { createPinia, PiniaPlugin } from 'pinia'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'
import { createApp } from 'vue'
// import SlideVerify from 'vue-monoplasty-slide-verify'
// @ts-ignore
import screenShort from "vue-web-screen-shot"
import App from './App.vue'
import { i18n } from "./i18n/index.ts"
import router from './router'
// import { registerComponents } from './utils/regComps'
// GoCaptcha
import "go-captcha-vue/dist/style.css"
import GoCaptcha from "go-captcha-vue"


const app = createApp(App)
const pinia = createPinia()
pinia.use(piniaPluginPersistedstate as PiniaPlugin)
// 调用注册函数
// registerComponents(app);
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}
app.use(screenShort, { enableWebRtc: true })
app.use(pinia)
app.use(ElementPlus)
app.use(router)
app.use(i18n)
app.use(GoCaptcha)
app.mount('#app')
