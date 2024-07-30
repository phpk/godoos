<template>
  <div class="app-item">
    <img v-if="item?.icon" draggable="false" class="app-img" :src="item?.icon" alt="" />
    <div class="app-content">
      <div class="app-title">{{ item?.name }}</div>
      <div class="app-desc">
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
          <button @click="uninstall?.(item)" v-if="!item?.isRuning">卸载</button>
        </template>
        <template v-else>
          <button v-if="item?.url" @click="install?.(item)" :disabled="item?.progress">安装</button>
          
        </template>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
defineProps({
  item: Object,
  installedList: Array,
  install: Function,
  uninstall: Function,
  pause: Function,
  start: Function,
  restart: Function
});
</script>
<style scoped>
.app-item {
  width: 350px;
  height: 100px;
  margin: 10px;
  display: flex;
  flex-direction: row;
  justify-content: flex-start;
  align-items: center;
  padding: 10px;
  border-radius: 6px;
  border: 1px solid #cccccc79;
  transition: all 0.2s;
}

.app-img {
  padding: 10px;
  width: 32px;
  border-radius: 4px;
  border: 1px solid #cccccc7b;
  user-select: none;
}

.app-content {
  margin-left: 10px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: flex-start;
  height: 80px;
  width: 100%;
}

.app-title {
  font-size: 16px;
  font-weight: bold;
  padding-bottom: 6px;
  user-select: none;
}

.app-desc {
  font-size: 12px;
  color: #666;
  user-select: none;
}

.app-item:hover {
  box-shadow: 0 0 20px 2px #dddddd5b;
  transform: translateY(-2px);
}

.app-item:active {
  box-shadow: 0 0 0px #ccc;
  transform: translateY(0px);
}

.app-button {
  width: 100%;
  display: flex;
  flex-direction: row;
  justify-content: flex-end;
  align-items: center;
  gap: 10px;
}

.app-button button {
  width: 60px;
  height: 30px;
  border-radius: 4px;
  border: 1px solid #ccc;
  background-color: #fff;
  cursor: pointer;
  transition: all 0.2s;
}

.app-button button:hover {
  box-shadow: 0 0 10px #ccc;
  /* transform: translateY(-2px); */
}

.app-button button:active {
  box-shadow: 0 0 0px #ccc;
  transform: translateY(0px);
}
</style>
