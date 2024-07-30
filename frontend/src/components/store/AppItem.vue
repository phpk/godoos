<template>
  <div class="app-item">
    <img v-if="item?.icon" draggable="false" class="app-img" :src="item?.icon" alt="" />
    <div class="app-content">
      <div class="app-title">{{ item?.name }}</div>
      <div class="app-desc">
        <span>{{ item?.desc ?? item?.name }}</span>
        <el-progress :text-inside="true" style="margin-top: 3px;" v-if="item?.progress" :stroke-width="20"
          :percentage="item?.progress" />
      </div>
      <div class="app-button">

        <template v-if="installedList!.includes(item?.name)">
          <div v-if="item!.checkProgress">
            <button @click="pause?.(item)" v-if="item?.isRuning">暂停</button>
            <button @click="start?.(item)" v-else>启动</button>
          </div>
          <div v-if="item!.hasRestart && item?.isRuning">
            <button @click="restart?.(item)">重启</button>
          </div>
          <button @click="setting?.(item)" v-if="!item?.isRuning && item?.setting">配置</button>
          <button @click="uninstall?.(item)" v-if="!item?.isRuning">卸载</button>
        </template>
        <template v-else>
          <button @click="install?.(item)" :disabled="item?.progress">安装</button>

        </template>
      </div>
    </div>
  </div>

</template>
<script setup lang="ts">
import { BrowserWindow } from "@/system";
defineProps({
  item: Object,
  installedList: Array,
  install: Function,
  uninstall: Function,
  pause: Function,
  start: Function,
  restart: Function
});

function setting(item: any) {
  const win = new BrowserWindow({
    title: "配置",
    url: `http://localhost:56780/static/${item.name}/index.html`,
    icon: "gallery",
    width: 500,
    height: 500,
    center: true,
    minimizable: false,
    resizable: true,
  });
  win.show()
}
</script>
<style scoped>
.app-item {
  width: calc(100% - 20px);
  min-height: 100px;
  min-width:200px;
  margin: 10px;
  display: flex;
  flex-direction: row;
  justify-content: flex-start;
  align-items: center;
  padding: 10px;
  border-radius: 6px;
  border: 1px solid rgba(204, 204, 204, 0.79);
  transition: all 0.2s;
  background-color: white;
}

@media (max-width: 768px) {
  .app-item {
    flex-direction: column;
    width: calc(100% - 20px);
  }
}

.app-img {
  width: 32px;
  height: 32px;
  border-radius: 4px;
  border: 1px solid rgba(204, 204, 204, 0.79);
  user-select: none;
}

.app-content {
  margin-left: 10px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: flex-start;
  height: 100%;
}

.app-title {
  font-size: 16px;
  font-weight: bold;
  padding-bottom: 6px;
  user-select: none;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.app-desc {
  font-size: 12px;
  color: #666;
  user-select: none;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  padding: 5px; /* 添加内边距 */
  /* 可选：添加边框 */
  /* border: 1px solid #ddd; */
  /* 可选：添加背景色 */
  /* background-color: #f9f9f9; */
}
@media (min-width: 769px) {
  .app-desc {
    font-size: 14px; /* 在大屏幕上使用更大的字体大小 */
  }
}
.app-item:hover {
  box-shadow: 0 0 20px 2px rgba(221, 221, 221, 0.5);
  transform: translateY(-2px);
}

.app-item:active {
  box-shadow: none;
  transform: translateY(0);
}

.app-button {
  display: flex;
  flex-direction: row;
  justify-content: flex-end;
  align-items: center;
  gap: 10px;
  margin-top: 20px;
}

.app-button button {
  width: 60px;
  height: 30px;
  border-radius: 4px;
  border: 1px solid #ccc;
  background-color: #fff;
  cursor: pointer;
  transition: all 0.2s;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.app-button button:hover {
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

.app-button button:active {
  box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.15);
}
</style>
