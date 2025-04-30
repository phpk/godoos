<template>
  <div class="taskbar" @mouseleave="hideWindowList">
    <div class="start-button" @click="store.toggleStartMenu()">
      <icon name="windows" :size="20" />
    </div>
    <StartMenu />
    <div class="taskbar-icons">
      <div 
        v-for="(windows, icon) in groupedWindows" 
        :key="icon" 
        class="taskbar-icon" 
        @mouseenter="showWindowList(windows)"
      >
        <template v-if="Array.isArray(windows) && windows.length > 1">
          <icon :name="windows[0].icon" size="20" />
          <div class="window-count">{{ windows.length }}</div>
          <div v-if="showingWindowList && currentWindowList[0].icon == windows[0].icon" class="window-list">
            <div v-for="win in currentWindowList" :key="win.id" class="window-list-item" @click="restoreWindow(win.id)">
              {{ t(win.title) }}
            </div>
          </div>
        </template>
        <template v-else-if="Array.isArray(windows) && windows.length === 1">
          <icon :name="windows[0].icon" size="20" @click="restoreWindow(windows[0].id)" />
        </template>
        
      </div>
    </div>

    <div class="system-tray">
      <Screenshort />
      <ScreenRecorder />
      <TrayDate />
      <icon name="notices" :size="16"/>
      
    </div>
    <div class="vertical-separator" @click="windowStore.minimizeAllWindows"></div> 
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useDesktopStore } from '@/stores/desktop'
import { useWindowStore } from '@/stores/window'
import { t } from '@/i18n'

const store = useDesktopStore()
const windowStore = useWindowStore()

// 使用计算属性获取分组窗口
const groupedWindows = computed(() => {
  return windowStore.groupedWindows
})

const showingWindowList = ref(false)
const currentWindowList = ref<any>([])

const restoreWindow = (id: string) => {
  windowStore.minimizeWindow(id) // Toggle to restore
  windowStore.bringToFront(id)
}

const showWindowList = (windows: any) => {
  showingWindowList.value = true
  currentWindowList.value = windows
}

const hideWindowList = () => {
  showingWindowList.value = false
  currentWindowList.value = []
}
</script>

<style scoped>
.taskbar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  align-items: center;
  padding: 0 10px;
  z-index: 1000;
  align-items: center;
  background: #f0f0f0;
  color: #333333;
  height: 48px;
  padding: 0;
}

.start-button {
  border-radius: 50%;
  width: 40px;
  height: 40px;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  margin-left:8px;
  margin-right: 10px;
}

.start-button:hover {
  background-color: aliceblue;
}

.taskbar-icons {
  display: flex;
  gap: 16px;
}

.taskbar-icon {
  cursor: pointer;
  transition: transform 0.3s;
  position: relative;
}

.taskbar-icon:hover {
  transform: scale(1.1);
  color: #2980b9;
  background-color: aliceblue;
}

.window-count {
  position: absolute;
  top: -5px;
  right: -5px;
  background-color: red;
  color: white;
  border-radius: 50%;
  width: 18px;
  height: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
}

.system-tray {
  display: flex;
  margin-left: auto;
  align-items: center;
  gap: 8px;
}

.system-tray .el-icon {
  margin: 0 4px;
  cursor: pointer;
}

.system-tray .el-icon:hover {
  color: #2980b9;
}

.system-tray div {
  font-size: 14px;
}

.window-list {
  position: fixed;
  bottom: 30px;
  background-color: white;
  border: 1px solid #ccc;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  z-index: 1000;
  padding: 10px;
  border-radius: 4px;
  min-width: 120px; /* 确保最小宽度 */
  max-width: 200px; /* 确保最大宽度 */
}

.window-list-item {
  padding: 5px 10px;
  cursor: pointer;
  white-space: nowrap; /* 防止换行 */
  overflow: hidden; /* 隐藏溢出内容 */
  text-overflow: ellipsis; /* 使用省略号显示溢出内容 */
  border-radius: 4px; /* 圆角 */
  font-size: 12px;
  transition: background-color 0.3s; /* 添加过渡效果 */
}

.window-list-item:hover {
  background-color: #f0f0f0;
}
.vertical-separator {
  height: 48px;
  width:10px;
  border-left:1px solid #ccc;
  margin-left: 8px;
}
.vertical-separator:hover {
  background-color: aliceblue;
}
@media screen and (max-width: 768px) {
  .taskbar {
    display: none;
  }
}
</style>