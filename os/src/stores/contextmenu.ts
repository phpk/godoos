// src/stores/contextmenu.ts
import { defineStore } from 'pinia'
import { computed, ref, Ref } from 'vue'

export const useContextMenuStore = defineStore('contextMenu', () => {
  // 上下文菜单是否可见
  const isContextMenuVisible: Ref<boolean> = ref(false)
  // 上下文菜单的位置
  const contextMenuTop: Ref<number> = ref(0)
  const contextMenuLeft: Ref<number> = ref(0)
  const contextMenuBottom: Ref<number> = ref(0)
  // 当前图标
  const currentFile: any = ref({})
  const currentPath: Ref<string> = ref('')

  const isShow = computed(() => {
    return isContextMenuVisible.value || isDesktopContextMenuVisible.value
  })
  const isShareFile = ref(false)
  // 桌面上下文菜单是否可见
  const isDesktopContextMenuVisible: Ref<boolean> = ref(false)
  // 桌面上下文菜单的位置
  const desktopContextMenuTop: Ref<number> = ref(0)
  const desktopContextMenuLeft: Ref<number> = ref(0)
  // 当前子菜单
  const activeSubMenu: Ref<string | null> = ref(null)

  // 显示上下文菜单
  const mobileStyle: any = ref()
  const showContextMenu = (icon: any, event: MouseEvent | TouchEvent) => {
    console.log(icon)
    if (isContextMenuVisible.value) return
    currentFile.value = icon
    isContextMenuVisible.value = true
    // contextMenuTop.value = event.clientY
    // contextMenuLeft.value = event.clientX
    isDesktopContextMenuVisible.value = false // 确保桌面上下文菜单不可见
    if ('clientX' in event && 'clientY' in event) {
      // 处理 MouseEvent
      const screenH = document.querySelector('.desktop')!.clientHeight
      if (event.clientY > screenH / 2) {
        contextMenuBottom.value = screenH - event.clientY
      } else {
        contextMenuTop.value = event.clientY
      }
      contextMenuLeft.value = event.clientX
    } else {
      // 处理 TouchEvent
      const screenW = window.screen.width //设备的宽度
      const screenH = document.querySelector('.desktop')!.clientHeight //窗口高度
      const touch = (event as TouchEvent).touches[0]
      if (touch) {
        if (touch.clientY > screenH / 2) {
          mobileStyle.value = `bottom:${screenH - touch.clientY}px;`
        } else {
          mobileStyle.value = `top:${touch.clientY}px;`
        }
        if (touch.clientX > screenW / 2) {
          mobileStyle.value += `right:${screenW - touch.clientX}px;`
        } else {
          mobileStyle.value += `left:${touch.clientX}px;`
        }
      }
    }
  }

  // 隐藏上下文菜单
  const hideContextMenu = () => {
    isContextMenuVisible.value = false
    isDesktopContextMenuVisible.value = false
    activeSubMenu.value = null
  }

  // 显示桌面上下文菜单
  const showDesktopContextMenu = (event: MouseEvent, currentpath: string) => {
    if (isDesktopContextMenuVisible.value) return
    currentPath.value = currentpath
    isContextMenuVisible.value = false // 确保文件上下文菜单不可见
    isDesktopContextMenuVisible.value = true
    desktopContextMenuTop.value = event.clientY
    desktopContextMenuLeft.value = event.clientX
  }

  // 切换子菜单
  const toggleSubMenu = (subMenuId: string | null) => {
    activeSubMenu.value = subMenuId
  }

  return {
    isShow,
    isContextMenuVisible,
    contextMenuTop,
    contextMenuLeft,
    currentFile,
    isDesktopContextMenuVisible,
    desktopContextMenuTop,
    desktopContextMenuLeft,
    activeSubMenu,
    currentPath,
    isShareFile,
    showContextMenu,
    hideContextMenu,
    showDesktopContextMenu,
    // handleContextMenuItemClick,
    // handleDesktopContextMenuItemClick,
    toggleSubMenu,
    mobileStyle,
    contextMenuBottom,
  }
})
