// src/stores/window.ts
import { defineStore } from 'pinia'
import { computed, Ref, ref } from 'vue'
import { dealIcon } from '@/utils/icon'
interface WindowState {
  id: string // 使用 string 类型的 id
  title: string
  icon: string
  component: any
  url?: string
  props?: Record<string, any>
  position: {
    x: number
    y: number
  }
  size: {
    width: number
    height: number
  }
  isMinimized: boolean
  isMaximized: boolean
  isVisible: boolean // 添加 isVisible 属性
  isDragging: boolean
  zIndex: number // 添加 zIndex 属性
}

// 生成唯一标识符的函数
const generateUniqueId = (): string => {
  const timestamp = Date.now().toString(36)
  const randomNum = Math.random().toString(36).substr(2, 5)
  return `${timestamp}-${randomNum}`
}
function checkIsFile(w: any) {
  if (w.props?.isFile) {
    w.icon = dealIcon(w.props)
    w.title = w.props?.title
  }
}
export const useWindowStore = defineStore('window', () => {
  // 窗口状态
  const windows: Ref<WindowState[]> = ref([])

  // 获取活动窗口
  const activeWindows = computed(() =>
    windows.value.filter((w) => !w.isMinimized && w.isVisible)
  )
  const currentZindex = ref(1)
  // 创建窗口
  const createWindow = (options: WindowState) => {
    options.position = calculateWindowPosition(
      options.size.width,
      options.size.height
    )
    currentZindex.value++
    options.zIndex = currentZindex.value // 初始化 zIndex
    windows.value.push(options)
    return options
  }

  // 封装创建窗口的逻辑
  const create = (options: {
    title: string
    icon: string
    component: any
    props?: Record<string, any>
    size: {
      width: number
      height: number
    }
    isMaximized: boolean
  }) => {
    const isOne = options.props?.isOne || false
    // console.log(options)
    if (isOne) {
      const exist = windows.value.find(
        (w) => w.props?.isOne === true && w.title === options.title
      )
      if (exist) {
        bringToFront(exist.id)
        updateWindowVisibility(exist.id, true)
        return
      }
    }
    checkIsFile(options)
    const windowId = generateUniqueId() // 使用生成唯一标识符的函数
    const win = createWindow({
      id: windowId, // 添加唯一标识符
      title: options.title,
      icon: options.icon,
      component: options.component,
      props: options.props || {},
      position: {
        x: 0,
        y: 0,
      },
      size: options.size,
      isMinimized: false, // 初始化时设置为未最小化
      isMaximized: options.isMaximized, // 初始化时设置为未最大化
      isVisible: true, // 初始化时设置为可见
      isDragging: false, // 初始化时设置为未拖拽
      zIndex: 1, // 初始化 zIndex
    })
    bringToFront(windowId)
    return win
  }

  // 动态计算窗口位置
  function calculateWindowPosition(
    width: number,
    height: number
  ): { x: number; y: number } {
    const windowWidth = window.innerWidth
    const windowHeight = window.innerHeight
    // 首先尝试将窗口居中靠上
    let x = (windowWidth - width) / 2
    let y = 10

    // 检查是否与现有窗口重叠
    for (let win of activeWindows.value) {
      if (
        x < win.position.x + win.size.width &&
        x + width > win.position.x &&
        y < win.position.y + win.size.height &&
        y + height > win.position.y
      ) {
        // 如果重叠，使用偏移逻辑
        x = 100 + windows.value.length * 30
        y = 50 + windows.value.length * 30
        break
      }
    }

    // 确保窗口在视窗可视范围内
    if (x + width > windowWidth) {
      x = windowWidth - width
    }
    if (y + height > windowHeight) {
      y = windowHeight - height
    }

    return { x, y }
  }

  // 关闭窗口
  const closeWindow = (id: string) => {
    const index = windows.value.findIndex((w) => w.id === id)
    if (index !== -1) {
      windows.value.splice(index, 1)
    }
  }

  // 最小化窗口
  const minimizeWindow = (id: string) => {
    const win = windows.value.find((w) => w.id === id)
    if (win) {
      win.isMinimized = !win.isMinimized
      win.isVisible = !win.isVisible
    }
  }

  // 切换窗口最大化状态
  const toggleMaximizeWindow = (id: string) => {
    const win = windows.value.find((w) => w.id === id)
    if (win) {
      win.isMaximized = !win.isMaximized
      if (win.isMaximized) {
        win.position = { x: 0, y: 0 }
        win.size = {
          width: window.innerWidth,
          height: window.innerHeight - 40, // Subtract taskbar height
        }
      } else {
        // Reset to default size when minimized
        win.size = {
          width: 800,
          height: 600,
        }
      }
    }
  }

  // 更新窗口位置
  const updateWindowPosition = (
    id: string,
    position: { x: number; y: number }
  ) => {
    const win = windows.value.find((w) => w.id === id)
    if (win && !win.isMaximized) {
      win.position = position
    }
  }
  const updateWindow = (win: any) => {
    const index = windows.value.findIndex((w) => w.id === win.id)
    if (index !== -1) {
      windows.value[index] = win
    }
  }
  // 更新窗口大小
  const updateWindowSize = (
    id: string,
    size: { width: number; height: number }
  ) => {
    const win = windows.value.find((w) => w.id === id)
    if (win && !win.isMaximized) {
      win.size = size
    }
  }

  // 将窗口带到最前面
  const bringToFront = (id: string) => {
    const win = windows.value.find((w) => w.id === id)
    if (win) {
      currentZindex.value++
      win.zIndex = currentZindex.value
    }
  }

  // 更新窗口可见性
  const updateWindowVisibility = (id: string, isVisible: boolean) => {
    const win = windows.value.find((w) => w.id === id)
    if (win) {
      win.isVisible = isVisible
    }
  }
  const setDragging = (id: string, isDragging: boolean) => {
    const window = windows.value.find((win) => win.id === id)
    if (window) {
      window.isDragging = isDragging
    }
  }
  // 根据图标分组窗口
  const groupedWindows = computed(() => {
    const grouped = new Map<string, WindowState[]>()

    windows.value.forEach((win) => {
      if (!grouped.has(win.icon)) {
        grouped.set(win.icon, [])
      }
      grouped.get(win.icon)?.push(win)
    })

    // 将 Map 转换为对象，保持插入顺序
    return Object.fromEntries(grouped.entries())
  })

  const minimizeAllWindows = () => {
    activeWindows.value.forEach((win) => {
      win.isMinimized = true
      win.isVisible = false
    })
  }

  return {
    windows,
    activeWindows,
    create,
    createWindow,
    closeWindow,
    updateWindow,
    minimizeWindow,
    minimizeAllWindows,
    toggleMaximizeWindow,
    updateWindowPosition,
    updateWindowSize,
    bringToFront,
    updateWindowVisibility,
    setDragging,
    generateUniqueId, // 导出生成唯一标识符的函数
    groupedWindows, // 添加分组窗口的计算属性
  }
})
