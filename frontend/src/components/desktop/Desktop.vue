<template>
  <div
    class="desktop"
    @dragenter.prevent
    @dragover.prevent
    @drop="dragFileToDrop($event, `${system._options.userLocation}Desktop`)"
    :style="{
      filter: `brightness(${system._rootState.info.brightness * 2}%)`,
    }"
  >
    <div class="userarea" @contextmenu.stop="handleRightClick" @mousedown="userareaDown">
      <div @mousedown="backgroundDown">
        <DeskItem class="userarea-upper zhighher" :on-chosen="onChosen"></DeskItem>
        <DesktopBackground class="userarea-upper"></DesktopBackground>
      </div>
      <WindowGroup></WindowGroup>
      <NotifyGroup></NotifyGroup>
      <MessageCenterPop></MessageCenterPop>
      <Chosen></Chosen>
      <Notice v-if="upgradeStore.hasNotice"></Notice>
      <Upgrade v-if="upgradeStore.hasUpgrade"></Upgrade>
      <Ad v-if="upgradeStore.hasAd"></Ad>
    </div>
    <div class="bottom">
      <Taskbar></Taskbar>
    </div>
    <ContextMenu></ContextMenu>
  </div>
</template>
<script lang="ts" setup>
import { emitEvent } from "@/system/event";
import { useContextMenu } from "@/hook/useContextMenu";
import { useFileDrag } from "@/hook/useFileDrag";
import { Rect, useRectChosen } from "@/hook/useRectChosen";
import { useSystem } from "@/system";
import { onErrorCaptured } from "vue";
import { useUpgradeStore } from '@/stores/upgrade';

const { createDesktopContextMenu } = useContextMenu();
const { choseStart, chosing, choseEnd, getRect, Chosen } = useRectChosen();
const system = useSystem();
const upgradeStore = useUpgradeStore();
const { dragFileToDrop } = useFileDrag(system);

let chosenCallback: (rect: Rect) => void = () => {
  //
};
function onChosen(callback: (rect: Rect) => void) {
  chosenCallback = callback;
}
function userareaDown(e: MouseEvent) {
  emitEvent("desktop.background.leftClick", e);
  chosenCallback({
    left: e.clientX,
    top: e.clientY,
    width: 0,
    height: 0,
  });
}
function backgroundDown(e: MouseEvent) {
  choseStart(e);
  addEventListener("mousemove", backgroundMove);
  addEventListener("mouseup", backgroundUp);
}
function backgroundMove(e: MouseEvent) {
  emitEvent("desktop.background.leftMove", e);
  chosing(e);
  const rectValue = getRect();
  if (rectValue) {
    chosenCallback(rectValue);
  }
}
function backgroundUp(e: MouseEvent) {
  emitEvent("desktop.background.leftUp", e);
  choseEnd();
  const rectValue = getRect();
  if (rectValue) {
    chosenCallback(rectValue);
    emitEvent("desktop.background.rectChosen", rectValue);
  }
  removeEventListener("mousemove", backgroundMove);
  removeEventListener("mouseup", backgroundUp);
}

function handleRightClick(e: MouseEvent) {
  e.preventDefault();
  createDesktopContextMenu(e, `${system._options.userLocation}Desktop`, () => {
    system.initAppList();
  });
}

onErrorCaptured((err) => {
  system.emitError(err.message.toString());
});

</script>
<style lang="scss" scoped>
.desktop {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  flex-grow: 0;
  overflow: hidden;

  .zhighher {
    z-index: 2;
  }

  .userarea {
    flex: 1;
    position: relative;
    overflow: hidden;

    .userarea-upper {
      position: absolute;
      top: 0;
      left: 0;
    }
  }
}
</style>
