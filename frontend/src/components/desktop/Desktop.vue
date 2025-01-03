<template>
  <div class="desktop" @dragenter.prevent @dragover.prevent
    @drop="dragFileToDrop($event, `${system._options.userLocation}Desktop`)" :style="{
      filter: `brightness(${system._rootState.info.brightness * 2}%)`,
    }">
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
    <div class="bottom-bar" v-if="isMobileDevice()">
      <div @click.stop="handle(item)" class="magnet-item" :style="{
        animationDelay: `${Math.floor(index / 4) * 0.02}s`,
        animationDuration: `${Math.floor(index / 4) * 0.04 + 0.1}s`,
      }" v-for="(item, index) in bottomApp" v-glowing :key="basename(item.path)">
        <FileIcon class="magnet-item_img" :file="item" />
        <span class="magnet-item_title">{{ getName(item) }}</span>
      </div>
    </div>
    <ContextMenu></ContextMenu>
  </div>
</template>
<script lang="ts" setup>
import { emitEvent } from "@/system/event";
import { useContextMenu } from "@/hook/useContextMenu";
import { useFileDrag } from "@/hook/useFileDrag";
import { Rect, useRectChosen } from "@/hook/useRectChosen";
import { useSystem, BrowserWindow } from "@/system";
import { onErrorCaptured } from "vue";
import { useUpgradeStore } from '@/stores/upgrade';
import { isMobileDevice } from "@/util/device";
import { basename } from "@/system/core/Path";
import { t } from "@/i18n";
import { useAppOpen } from "@/hook/useAppOpen";

const { openapp } = useAppOpen("menulist");
function handle(item: any) {
  if (item.name.includes('localchat')) {
    openapp(item);
    emitEvent("desktop.app.open");
  } else {
    emitEvent("magnet.item.click", item);
  }

  const sys = useSystem();
  const winopt = sys._rootState.windowMap["Menulist"].get(item.title);
  if (winopt) {
    if (winopt._hasShow) {
      return;
    } else {
      winopt._hasShow = true;
      winopt.window.fullscreen = true
      const win = new BrowserWindow(winopt.window);
      win.show();
      win.on("close", () => {
        winopt._hasShow = false;
      });
    }
  }
}
function getName(item: any) {
  const name = basename(item.path);
  if (name.endsWith(".exe")) {
    return t(name.replace(".exe", ""));
  } else {
    return name;
  }
}

const bottomApp = [
  {
    "isFile": true,
    "isDirectory": false,
    "isSymlink": false,
    "size": 32,
    "modTime": "2024-12-16T09:17:18.7214789+08:00",
    "atime": "2024-12-16T09:17:18.7214789+08:00",
    "birthtime": "2024-12-16T09:17:18.7214789+08:00",
    "mtime": "2024-12-16T09:17:18.7214789+08:00",
    "rdev": 0,
    "mode": 511,
    "name": "computer.exe",
    "path": "/C/Users/Menulist/computer.exe",
    "oldPath": "/C/Users/Menulist/computer.exe",
    "parentPath": "/C/Users/Menulist",
    "content": "link::Desktop::computer::diannao",
    "ext": "exe",
    "title": "computer",
    "id": 1,
    "isPwd": false
  },
  {
    "isFile": true,
    "isDirectory": false,
    "isSymlink": false,
    "size": 31,
    "modTime": "2024-12-16T09:17:18.7214789+08:00",
    "atime": "2024-12-16T09:17:18.7214789+08:00",
    "birthtime": "2024-12-16T09:17:18.7214789+08:00",
    "mtime": "2024-12-16T09:17:18.7214789+08:00",
    "rdev": 0,
    "mode": 511,
    "name": "setting.exe",
    "path": "/C/Users/Menulist/setting.exe",
    "oldPath": "/C/Users/Menulist/setting.exe",
    "parentPath": "/C/Users/Menulist",
    "content": "link::Desktop::setting::setting",
    "ext": "exe",
    "title": "setting",
    "id": 15,
    "isPwd": false
  },
  {
    "isFile": true,
    "isDirectory": false,
    "isSymlink": false,
    "size": 30,
    "modTime": "2024-12-16T09:17:18.7214789+08:00",
    "atime": "2024-12-16T09:17:18.7214789+08:00",
    "birthtime": "2024-12-16T09:17:18.7214789+08:00",
    "mtime": "2024-12-16T09:17:18.7214789+08:00",
    "rdev": 0,
    "mode": 511,
    "name": "localchat.exe",
    "path": "/C/Users/Desktop/localchat.exe",
    "oldPath": "/C/Users/Desktop/localchat.exe",
    "parentPath": "/C/Users/Desktop",
    "content": "link::Desktop::localchat::chat",
    "ext": "exe",
    "title": "localchat",
    "id": 3,
    "isPwd": false
  }
];
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

@media screen and (max-width: 768px) {
  .bottom {
    display: none;
  }

  .bottom-bar {
    display: flex;
    justify-content: space-evenly;
    position: absolute;
    width: vw(340);
    left: 50%;
    transform: translateX(-50%);
    bottom: vh(15);
    height: vh(80);
    border-radius: vw(50);
    color: #e5e3e3d5;
    background-color: rgba(255, 255, 255, .2);
    backdrop-filter: blur(15px);

    .magnet-item {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;

      .magnet-item_img {
        width: vw(40);
        height: vh(40);
      }
    }
  }
}
</style>
