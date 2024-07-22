<template>
  <div
    :style="{
      top: y + 'px',
      left: x + 'px',
    }"
    class="contextmenu"
    v-if="rootState.contextMenu"
  >
    <div class="contextmenu-item" v-for="item in menuList" :key="item.label">
      <div class="option-title" @click="handleClick(item)">{{ item.label }}</div>
      <div class="icon-arrow" v-if="item.submenu?.length"></div>
      <div class="children-item" v-if="item.submenu?.length">
        <div class="contextmenu-item" v-for="citem in item.submenu" :key="citem.label">
          <div class="option-title" @click="handleClick(citem)">{{ citem.label }}</div>
        </div>
      </div>
    </div>
  </div>
</template>
<script lang="ts" setup>
import { ref, watch } from "vue";
import { mountEvent, emitEvent } from "@/system/event";
import { useSystem, MenuItem } from "@/system";
// import { MenuItem } from '@/menu/MenuItem';

const x = ref(0);
const y = ref(0);
const menuList = ref<MenuItem[]>([]);
const rootState = useSystem()._rootState;
watch(
  () => rootState.contextMenu,
  (contextMenu) => {
    // get window inner width and height
    const innerWidth = rootState.info.screenWidth;
    const innerHeight = rootState.info.screenHeight;
    // get contextmenu width
    const contextmenuWidth = 160;
    // get contextmenu height
    const contextmenuHeight = 24 * (contextMenu?.items.length || 0);
    // get mouse position
    const outer = useSystem()?.rootRef;
    const mouseX = (contextMenu?._mouse?.x || 0) - (outer?.offsetLeft || 0);
    const mouseY = (contextMenu?._mouse?.y || 0) - (outer?.offsetTop || 0);

    // get contextmenu position
    const contextmenuX =
      mouseX + contextmenuWidth > innerWidth ? mouseX - contextmenuWidth : mouseX;
    const contextmenuY =
      mouseY + contextmenuHeight > innerHeight ? mouseY - contextmenuHeight : mouseY;

    x.value = contextmenuX;
    y.value = contextmenuY;
    menuList.value = contextMenu?.items || [];
  }
);

mountEvent("contextMenu.hidden", () => {
  useSystem()._rootState.contextMenu = null;
});

function handleClick(item: MenuItem) {
  if (item?.click) {
    item?.click?.();
    emitEvent("contextMenu.hidden");
  }
}
</script>
<style lang="scss">
.contextmenu {
  position: absolute;
  top: 0;
  left: 0;
  width: var(--contextmenu-width);
  z-index: 100;
  background-color: #f0f0f0; /* 更改背景色以接近Windows 11的浅灰色 */
  border: 1px solid #e5e5e5; /* 边框颜色调整 */
  border-radius: 8px; /* 添加圆角 */
  padding: 4px 0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1); /* 添加阴影效果 */
  user-select: none;

  &::before {
    content: "";
    position: absolute;
    top: -4px;
    left: calc(50% - 4px);
    width: 0;
    height: 0;
    border-left: 4px solid transparent;
    border-right: 4px solid transparent;
    border-bottom: 4px solid #e5e5e5; /* 小三角形下拉提示 */
  }

  .contextmenu-item {
    height: var(--ui-list-item-height);
    font-size: var(--ui-font-size);
    padding: 4px 12px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: space-between;
    border-radius: 4px;
    transition: background-color 0.1s;

    .option-title {
      width: 100%;
    }

    .children-item {
      display: none;
      position: absolute;
      top: 0;
      left: 100%;
      width: var(--contextmenu-width);
      z-index: 101; /* 确保子菜单在父菜单之上 */
      background-color: #f0f0f0;
      border: 1px solid #e5e5e5;
      border-radius: 8px;
      padding: 4px 0;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
      user-select: none;
    }

    &:hover {
      background-color: #e5e5e5; /* 鼠标悬停时的背景色变淡 */
      .children-item {
        display: block;
      }
    }
  }

  .icon-arrow {
    display: inline-block;
    width: 8px;
    height: 8px;
    transform: translateY(2px) rotate(-45deg);
    border: 2px solid rgba(0, 0, 0, 0.25);
    border-left: none;
    border-top: none;
    transition: all 0.1s;
  }
}
</style>
