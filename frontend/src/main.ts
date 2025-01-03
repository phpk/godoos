import * as ElementPlusIconsVue from "@element-plus/icons-vue"
//svg插件需要配置代码
import ElementPlus from "element-plus"
import "element-plus/dist/index.css"
import { createApp } from "vue"
import screenShort from "vue-web-screen-shot"
import App from "./App.vue"
import "./assets/windows10.scss"
import { i18n } from "./i18n/index.ts"
import pinia from "./stores/index.ts"
import router from "./system/router"
import { NavBar } from "vant"
import "vant/lib/index.css"
const app = createApp(App)

app.use(router)
app.use(ElementPlus)
app.use(pinia)
app.use(i18n)
app.use(NavBar)
app.use(screenShort, { enableWebRtc: true })

for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}
app.mount("#app")
