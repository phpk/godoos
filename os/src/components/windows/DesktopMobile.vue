<template>
  <div class="desktop-area" @click.stop="handleDesktopClick($event)">
    <swiper :modules="modules" :loop="true" :slides-per-view="1" :pagination="{ type: 'bullets', clickable: true }"
      class="swiperBox">
      <swiper-slide class="swiper-slide" v-for="(item, i) in swiperCount" :key="item">
        <div v-for="icon in iconsToShow(i)" :key="icon.id" style="padding-top: 30px;" class="desktop-icon"
          @touchstart="gtouchstart($event, icon)" @touchmove="gtouchmove" @touchend="gtouchend($event, icon)">
          <icon :name="dealIcon(icon)" :size="50" />
          <span class="icon-name">{{ t(icon.title) }}</span>
        </div>
      </swiper-slide>
    </swiper>
    <div class="bottom-bar">
      <div v-for="icon in bottomBarData" :key="icon.id" class="desktop-icon" @click="handleClick($event, icon)">
        <icon :name="dealIcon(icon)" :size="50" />
        <span class="icon-name">{{ t(icon.title) }}</span>
      </div>
      <div class="desktop-icon" @click.stop="handleLogOut">
        <el-icon :size="50">
          <SwitchButton />
        </el-icon>
        <span class="icon-name">退出</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue';
import { t } from '@/i18n';
import { Swiper, SwiperSlide } from 'swiper/vue';
import { dealIcon } from '@/utils/icon';
import { Pagination } from 'swiper/modules';
import { useClickingStore } from "@/stores/clicking";
import { useContextMenuStore } from "@/stores/contextmenu";
import { useFileSystemStore } from "@/stores/filesystem";
import { useDesktopStore } from "@/stores/desktop";
import { useLoginStore } from '@/stores/login';
import { useRouter } from "vue-router";
import 'swiper/swiper-bundle.css';

import { eventBus } from "@/interfaces/event";


const router = useRouter();
const desktopStore = useDesktopStore();
const contextMenuStore = useContextMenuStore();
const clickingStore = useClickingStore();
const fileSystemStore = useFileSystemStore();
const loginStore = useLoginStore();

const props = defineProps({
  fileList: {
    type: Array as () => any[],
    required: true
  },
  isDesktop: {
    type: Boolean,
    default: true,
  },
  currentPath: {
    type: String,
    default: "/C/Users/Desktop",
  },
})
const emit = defineEmits(["navigateTo", "refeshList"]);
const currentPath = ref(props.currentPath);

// 底部栏数据
const bottomBarData = computed(() => {
  return props.fileList.filter(item => item.name.includes('computer') || item.name.includes('workchat'))
})
// 在modules加入要使用的模块
const modules = [Pagination];

// 计算swiper页数 页最多显示20个
const swiperCount = computed(() => {
  return Math.ceil(props.fileList.length / 20);
});

// 截取的图标数组
const iconsToShow = (i: number) => {
  return props.fileList.slice(20 * i, 20 * (i + 1));
};

const handleClick = (event: TouchEvent | MouseEvent, item: any) => {
  event.preventDefault();
  contextMenuStore.hideContextMenu(); // 隐藏上下文菜单
  clickingStore.clickedIcons = [];
  if (item.isDirectory) {
    if (props.isDesktop) {
      router.push({ path: "/computer", query: { path: item.path } });
    } else {
      currentPath.value = item.path;
      emit("navigateTo", item.path);
    }
  } else {
    console.log(item)
    fileSystemStore.openFile(item);
  }
}

const handleDesktopClick = (event: MouseEvent) => {
  event.preventDefault();

  fileSystemStore.initFolder(currentPath.value);
  contextMenuStore.hideContextMenu(); // 隐藏上下文菜单
  clickingStore.resetClickedIcons(event);
}

// 长按事件
//长按事件（起始）
let timer: any
const gtouchstart = (event: TouchEvent, item: any) => {
  timer = setTimeout(function () {
    longPress(event, item);
  }, 800);
  return false;
}

//手释放，如果在800毫秒内就释放，则取消长按事件，此时可以执行click应该执行的事件
const gtouchend = (event: TouchEvent, item: any) => {
  clearTimeout(timer); //清除定时器
  if (timer != 0) {
    handleClick(event, item)
  }
  return false;
}

//如果手指有移动，则取消所有事件，此时说明用户只是要移动而不是长按
const gtouchmove = () => {
  clearTimeout(timer); //清除定时器
  timer = 0;
}

//真正长按后应该执行的内容
const longPress = (event: TouchEvent, item: any) => {
  timer = 0;
  //执行长按要执行的内容
  contextMenuStore.showContextMenu(item, event)
}

const handleRefreshDesktop = () => {
  contextMenuStore.hideContextMenu();
  clickingStore.clickedIcons = [];
  if (props.isDesktop) {
    desktopStore.initDesktop();
  } else {
    emit("refeshList");
  }
};

const handleLogOut = () => {
  desktopStore.isStartMenuOpen = false
  loginStore.loginOut()
}

onMounted(() => {
  currentPath.value = props.currentPath;
  fileSystemStore.initFolder(currentPath.value);
  eventBus.on("refreshDesktop", handleRefreshDesktop);
  clickingStore.addEvents(currentPath.value, props.fileList);
});
watch(
  () => props.currentPath,
  (newPath) => {
    currentPath.value = newPath;
    clickingStore.addEvents(currentPath.value, props.fileList);
  }
);
watch(
  () => props.fileList,
  (newFileList) => {
    clickingStore.addEvents(currentPath.value, newFileList);
  }
);
onUnmounted(() => {
  clickingStore.removeEvents();
});

</script>

<style lang="scss" scoped>
@media screen and (max-width: 768px) {
  .desktop-area {
    width: 100%;
    height: 100%;
    position: relative;
  }

  .swiperBox {
    width: 100vw;
    padding-top: vh(10);
    height: vh(560);

    .swiper-slide {
      display: grid;
      grid-template-columns: repeat(4, 1fr);
      grid-template-rows: repeat(5, vh(100));
    }
  }

  .desktop-icon {
    width: auto;
    display: flex;
    flex-direction: column;
    align-items: center;
    cursor: pointer;
    border-radius: 4px;
  }

  .icon-name {
    margin-top: 5px;
    font-size: 13px;
    box-sizing: border-box;
    padding: 0 10px;
    width: 100%;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
    word-break: break-all;
    text-align: center;
    // color: transparent;
    font-family: STSong;
    // text-shadow: 0 0 0.5px white, 0 0 1px black, 0 0 2px rgba(0, 0, 0, 2);
  }

  .bottom-bar {
    display: flex;
    justify-content: space-evenly;
    align-items: center;
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
  }
}
</style>