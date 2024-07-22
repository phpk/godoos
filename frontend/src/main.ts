import { createApp } from "vue";
import App from "./App.vue";
import pinia from './stores/index.ts'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import './assets/windows10.scss'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import router from './system/router'
import screenShort from "vue-web-screen-shot"
const app = createApp(App)
app.use(router)
app.use(ElementPlus)
app.use(pinia)
app.use(screenShort, { enableWebRtc: true })

for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}
app.mount("#app");
