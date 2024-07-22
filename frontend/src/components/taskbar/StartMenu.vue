<template>
  <div class="startmenuicon" @click.stop="emitClick">
    <svg class="icon startmenuicon-img" aria-hidden="true">
      <use xlink:href="#icon-windows"></use>
    </svg>
  </div>
  <Transition name="startmenu">
    <div class="startmenu" @click="handleClick" v-if="isStartmenuShow">
      <div class="startmenu-item">
        <StartOption></StartOption>
      </div>
      <div class="startmenu-item">
        <Magnet></Magnet>
      </div>
    </div>
  </Transition>
</template>
<script lang="ts" setup>
import { emitEvent, mountEvent } from "@/system/event";
import { ref } from "vue";

const isStartmenuShow = ref(false);
mountEvent("startmenu.changeVisible", function () {
  isStartmenuShow.value = !isStartmenuShow.value;
});
mountEvent("startmenu.hidden", function () {
  isStartmenuShow.value = false;
});
function emitClick(e: MouseEvent) {
  emitEvent("taskbar.startmenu.leftClick", e);
}
function handleClick(e: MouseEvent) {
  emitEvent("startMenu.click", e);
}
</script>
<style lang="scss" scoped>
.startmenuicon {
  width: var(--startmenu-icon-size);
  height: var(--task-bar-height);
  background-color: var(--theme-main-color);
  user-select: none;
  display: flex;
  justify-content: center;
  align-items: center;
}

.startmenuicon-img {
  // width: 50%;
  // height: 110%;
  font-size: 1.5rem;
}

.startmenu {
  position: absolute;
  bottom: var(--task-bar-height);
  left: 0;
}

.startmenuicon:hover {
  background-color: var(--color-gray-hover);

  .startmenuicon-img {
    opacity: 0.5;
  }
}

.startmenu-enter-active {
  transition: all 0.3s var(--aniline);
}

.startmenu-leave-active {
  transition: all 0.05s;
}

.startmenu-enter-from {
  transform: translateY(100px);
  opacity: 0;
}

.startmenu-leave-to {
  transform: translateY(350px);
  opacity: 0;
}
.startmenu {
  --startmenu-width: 600px;
  --startmenu-height: 400px;
  height: var(--startmenu-height);
  background-color: rgba(242, 242, 242, 0.95); /* 更改为接近Win10开始菜单的背景色 */
  backdrop-filter: blur(10px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1); /* 添加轻微阴影效果 */
  border-radius: 4px;
  z-index: 40;
  display: flex;
  overflow: hidden;
}
.startmenu-item {
  position: relative;
}
</style>
