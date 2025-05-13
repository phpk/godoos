import { defineStore } from 'pinia'
import { ref, Ref, computed } from 'vue'
import * as fs from '@/api/net/files'
import { useLoginStore } from './login'

export const useDesktopStore = defineStore('desktop', () => {
  // 桌面图标
  const icons: any = ref([])
  const loginStore = useLoginStore()
  // 活动窗口列表
  const activeWindows: Ref<string[]> = ref([])

  // 菜单列表
  const menuList: any = ref([])
  const mobileicons = computed(() => {
    const map = new Map()
    ;[...icons.value, ...menuList.value].forEach((item: any) => {
      if (!map.has(item.name)) {
        map.set(item.name, item)
      }
    })
    return [...map.values()]
  })

  // 开始菜单是否打开
  const isStartMenuOpen: Ref<boolean> = ref(false)
  const screenshotStatus: Ref<boolean> = ref(false)
  const initDesktop = async () => {
    const res: any = await fs.desktop()
    if (!res) {
      loginStore.isLoginState = false
      return false
    }
    icons.value = res.apps
    menuList.value = res.menulist
    //console.log(menuList)
    loginStore.isLoginState = true
    return true
  }

  // 切换开始菜单状态
  const toggleStartMenu = () => {
    isStartMenuOpen.value = !isStartMenuOpen.value
  }

  // 打开窗口
  const openWindow = (windowName: string) => {
    if (!activeWindows.value.includes(windowName)) {
      activeWindows.value.push(windowName)
    }
  }

  // 关闭窗口
  const closeWindow = (windowName: string) => {
    activeWindows.value = activeWindows.value.filter((w) => w !== windowName)
  }
  const destroyScreenshotComponent = (status: boolean) => {
    screenshotStatus.value = status
  }
  return {
    icons,
    activeWindows,
    menuList,
    isStartMenuOpen,
    screenshotStatus,
    destroyScreenshotComponent,
    toggleStartMenu,
    openWindow,
    closeWindow,
    initDesktop,
    mobileicons,
  }
})
