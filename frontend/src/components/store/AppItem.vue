<script setup lang="ts">
import { BrowserWindow } from "@/system";
import { t } from "@/i18n";
const downloading = ref(false)
const props = defineProps({
  item: Object,
  installedList: Array,
  install: {
        type: Function,
        required: true,
  },
  uninstall: Function,
  pause: Function,
  start: Function,
  restart: Function
});

function setting(item: any) {
  const win = new BrowserWindow({
    title: t('store.setting'),
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
async function installApp(item: any) {
  downloading.value = true;
  await props.install(item);
}
</script>
<template>
  <div class="app-item">
    <img v-if="item?.icon" draggable="false" class="app-img" :src="item?.icon" alt="" />
    <div class="app-content">
      <div class="app-title">{{ item?.name }}</div>
      <div class="app-desc">
        <span>{{ item?.desc ?? item?.name }}</span>
        <el-progress :text-inside="true" style="margin-top: 3px;" v-if="item?.progress > 0" :stroke-width="20"
          :percentage="item?.progress" />
      </div>
      <div class="app-button">

        <template v-if="installedList!.includes(item?.name)">
          <template v-if="item!.hasStart">
            <el-button @click="pause?.(item)" v-if="item?.isRuning">{{ t('store.stop') }}</el-button>
            <el-button @click="start?.(item)" v-else>{{ t('store.start') }}</el-button>
          </template>
          <el-button @click="restart?.(item)" v-if="item!.hasRestart && item?.isRuning">{{ t('store.restart') }}</el-button>
          <el-button @click="setting?.(item)" v-if="!item?.isRuning && item?.setting">{{ t('store.setting') }}</el-button>
          <el-button @click="uninstall?.(item)" v-if="!item?.isRuning">{{ t('store.uninstall') }}</el-button>
        </template>
        <template v-else>
          <el-button type="primary" @click="installApp(item)" :disabled="item?.progress" :loading="downloading">{{ t('store.install') }}</el-button> 
        </template>
      </div>
    </div>
  </div>

</template>

<style scoped>
.app-item {
  width: calc(100% - 20px);
  height: 120px;
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
  width: 100%;
  -webkit-box-orient: vertical;
  padding: 5px; /* 添加内边距 */
  /* 可选：添加边框 */
  /* border: 1px solid #ddd; */
  /* 可选：添加背景色 */
  /* background-color: #f9f9f9; */
}
.app-desc .el-progress {
  width: 90% !important; /* 或者任何你想要的具体宽度 */
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
  gap: 5px;
  margin-top: 10px;
}

.app-button button {
  width: 50px;
  height: 30px;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 14px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.app-button button:hover {
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

.app-button button:active {
  box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.15);
}
</style>
