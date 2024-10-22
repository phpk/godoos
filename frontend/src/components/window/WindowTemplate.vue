<template>
  <div
    class="wintmp_outer dragwin"
    :class="{
      topwin: istop,
      max: windowInfo.state == WindowStateEnum.maximize,
      min: windowInfo.state == WindowStateEnum.minimize,
      fullscreen: windowInfo.state == WindowStateEnum.fullscreen,
      noframe: !windowInfo.frame,
      disable: windowInfo.disable,
    }"
    :style="customerStyle"
    @touchstart.passive="onFocus"
    @mousedown="onFocus"
    ref="$win_outer"
    v-dragable
  >
    <!-- 窗口标题栏  -->
    <div class="wintmp_uper" @contextmenu.prevent>
      <MenuBar :browser-window="browserWindow"></MenuBar>
    </div>
    
    <div
      class="wintmp_main"
      :class="{ resizeing: resizemode != 'null' }"
      @mousedown.stop="predown"
      @touchstart.stop.passive="predown"
      @contextmenu.stop.prevent
    >
      <div
        class="content-mask"
        v-if="!istop && typeof browserWindow.content === 'string'"
      ></div>
      <WindowInner :win="browserWindow"></WindowInner>
    </div>
    <!-- 使用 v-for 生成拖拽边界 -->
    <div
      v-for="border in dragBorders"
      :key="border.type"
      :class="[
        border.class,
        'win_drag_border',
        { isChoseMode: resizemode == border.type },
        border.cursorClass
      ]"
      v-if="resizable"
      draggable="false"
      @mousedown.stop.prevent="startScale($event, border.type)"
      @touchstart.stop.passive="startScale($event, border.type)"
    ></div>
  </div>
</template>
<script lang="ts" setup>
import { onUnmounted, provide, ref } from "vue";
import { onMounted, computed, UnwrapNestedRefs } from "vue";
import { WindowStateEnum } from "@/system/window/BrowserWindow";

import { ScaleElement } from "@/system/window/dom/ScaleElement";
import { BrowserWindow } from "@/system/window/BrowserWindow";
import { emitEvent } from "@/system/event";
import { useSystem } from "@/system";
import { vDragable } from "@/system/window/MakeDragable";
const sys = useSystem();
const props = defineProps<{
  browserWindow: UnwrapNestedRefs<BrowserWindow>;
}>();

const browserWindow = props.browserWindow;
const windowInfo = browserWindow.windowInfo;
// 传递windowid
provide("browserWindow", browserWindow);
provide("system", sys);

function predown() {
  browserWindow.moveTop();
  emitEvent("window.content.click", browserWindow);
}

const customerStyle = ref<NonNullable<unknown>>({});

function onFocus(e: MouseEvent | TouchEvent): void {
  browserWindow?.moveTop();
  if (windowInfo.state === WindowStateEnum.maximize) {
    if (e instanceof MouseEvent) {
      e.preventDefault();
      e.stopPropagation();
    }
  }
}

const istop = computed(() => windowInfo.istop);

onMounted(() => {
  customerStyle.value = {
    width: computed(() => windowInfo.width + "px"),
    height: computed(() => windowInfo.height + "px"),
    left: computed(() => windowInfo.x + "px"),
    top: computed(() => windowInfo.y + "px"),

    zIndex: computed(() => {
      if (windowInfo.alwaysOnTop) {
        return 9999;
      }
      return windowInfo.zindex;
    }),
    backgroundColor: computed(() => windowInfo.backgroundColor),
  };
});

/*
挂载缩放事件
*/
const resizable = ref(windowInfo.resizable);
const resizemode = ref("null");
let scaleAble: ScaleElement;
onMounted(() => {
  scaleAble = new ScaleElement(
    resizemode,
    windowInfo.width,
    windowInfo.height,
    windowInfo.x,
    windowInfo.y
  );
  scaleAble.onResize((width: number, height: number, x: number, y: number) => {
    windowInfo.width = width || windowInfo.width;
    windowInfo.height = height || windowInfo.height;
    windowInfo.x = x || windowInfo.x;
    windowInfo.y = y || windowInfo.y;
    browserWindow.emit("resize", windowInfo.width, windowInfo.height);
  });
});
function startScale(e: MouseEvent | TouchEvent, dire: string) {
  console.log(e);
  if (windowInfo.disable) {
    return;
  }
  scaleAble?.startScale(
    e,
    dire,
    windowInfo.x,
    windowInfo.y,
    windowInfo.width,
    windowInfo.height
  );
}

onUnmounted(() => {
  scaleAble.unMount();
  // dragAble.unMount();
});

// 定义拖拽边界类型
const dragBorders = [
  { type: 'r', class: 'right_border' },
  { type: 'b', class: 'bottom_border' },
  { type: 'l', class: 'left_border' },
  { type: 't', class: 'top_border' },
  { type: 'rb', class: 'right_bottom_border' },
  { type: 'lb', class: 'left_bottom_border' },
  { type: 'lt', class: 'left_top_border' },
  { type: 'rt', class: 'right_top_border' },
];
</script>
<style>
.dragwin {
  position: absolute;
  width: 100%;
  height: 100%;
}
</style>
<style scoped lang="scss">
.wintmp_outer {
  position: absolute;
  padding: 0;
  margin: 0;
  // left: 0;
  // top: 0;
  width: max-content;
  height: max-content;
  background-color: #fff;
  border: var(--window-border);
  display: flex;
  flex-direction: column;
  box-shadow: var(--window-box-shadow);
  border-radius: var(--window-border-radius);
  .wintmp_main {
    position: relative;
    width: 100%;
    height: 100%;
    // background-color: rgb(255, 255, 255);
    overflow: hidden;
    contain: content;
  }
}

.topwin {
  border: 1px solid #0078d7;
  box-shadow: var(--window-top-box-shadow);
}

.icon {
  width: 12px;
  height: 12px;
}

.max {
  position: absolute;
  left: 0 !important;
  top: 0 !important;
  width: 100% !important;
  height: 100% !important;
  transition: left 0.1s ease-in-out, top 0.1s ease-in-out, width 0.1s ease-in-out,
    height 0.1s ease-in-out;
}

.disable {
  .wintmp_uper,
  .wintmp_main {
    pointer-events: none;
    user-select: none;
    box-shadow: none;
  }
}
.min {
  visibility: hidden;
  display: none;
}

.fullscreen {
  // 将声明移动到嵌套规则之上
  position: fixed;
  left: 0 !important;
  top: 0 !important;
  width: 100% !important;
  height: 100% !important;
  z-index: 205 !important;
  border: none;
  .wintmp_uper {
    display: none;
  }
}

.noframe {
  border: none;
  box-shadow: none;
  .wintmp_uper {
    display: none;
  }
}

.transparent {
  background-color: transparent;

  .wintmp_main {
    background-color: transparent;
  }

  .wintmp_uper {
    background-color: rgba(255, 255, 255, 0.774);
  }
}

.win_drag_border {
  position: absolute;
  background-color: rgba(0, 0, 0, 0);
}

.right_border {
  cursor: ew-resize;
  right: -12px;
  width: 16px;
  height: calc(100% - 4px);
}

.bottom_border {
  cursor: ns-resize;
  bottom: -12px;
  width: calc(100% - 4px);
  height: 16px;
}

.left_border {
  cursor: ew-resize;
  left: -12px;
  width: 16px;
  height: calc(100% - 4px);
}

.top_border {
  cursor: ns-resize;
  top: -12px;
  width: calc(100% - 4px);
  height: 16px;
}

.left_top_border {
  cursor: nwse-resize;
  left: -12px;
  top: -12px;
  width: 16px;
  height: 16px;
}

.right_top_border {
  cursor: nesw-resize;
  right: -12px;
  top: -12px;
  width: 16px;
  height: 16px;
}

.left_bottom_border {
  cursor: nesw-resize;
  left: -12px;
  bottom: -12px;
  width: 16px;
  height: 16px;
}

.right_bottom_border {
  cursor: nwse-resize;
  right: -12px;
  bottom: -12px;
  width: 16px;
  height: 16px;
}

.isChoseMode {
  width: 100vw;
  height: 100vh;
  position: fixed;
  left: 0;
  top: 0;
}

.resizeing {
  user-select: none;
  pointer-events: none;
}

.content-mask {
  position: absolute;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0);
  z-index: 100;
}
</style>