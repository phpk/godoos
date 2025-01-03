<template>
  <div v-if="!isMobileDevice()" class="desk-group">
    <FileList :on-chosen="props.onChosen" :on-open="openapp" :file-list="appList"></FileList>
  </div>
  <div v-else>
    <swiper :modules="modules" :loop="true" :slides-per-view="1" :space-between="20"
      :pagination="{ type: 'bullets', clickable: true }" class="swiperBox">
      <swiper-slide class="swiper-slide">
        <FileList :on-chosen="props.onChosen" :on-open="openapp" :file-list="appList"></FileList>
      </swiper-slide>
      <swiper-slide>
        <div @click.stop="handle(item)" class="magnet-item" :style="{
          animationDelay: `${Math.floor(index / 4) * 0.02}s`,
          animationDuration: `${Math.floor(index / 4) * 0.04 + 0.1}s`,
        }" v-for="(item, index) in menulist" v-glowing :key="basename(item.path)">
          <FileIcon class="magnet-item_img" :file="item" />
          <span class="magnet-item_title">{{ getName(item) }}</span>
        </div>
      </swiper-slide>
    </swiper>
  </div>
</template>
<script lang="ts" setup>
import { mountEvent } from "@/system/event";
import { useSystem, BrowserWindow } from "@/system/index.ts";
import { useAppOpen } from "@/hook/useAppOpen";
import { onMounted } from "vue";
import { isMobileDevice } from "@/util/device";
import { Swiper, SwiperSlide } from 'swiper/vue'
import { basename } from "@/system/core/Path";
import { emitEvent } from "@/system/event";
import { t } from "@/i18n";


// 引入swiper样式
import 'swiper/swiper-bundle.css';
// 引入swiper核心和所需模块
import { Pagination } from 'swiper/modules'
// 在modules加入要使用的模块
const modules = [Pagination]
const { openapp, appList } = useAppOpen("apps");
const { appList: menulist } = useAppOpen("menulist")
const props = defineProps({
  onChosen: {
    type: Function,
    required: true,
  },
});
onMounted(() => {
  mountEvent("file.props.edit", async () => {
    useSystem().initAppList();
  });
});
function handle(item: any) {
  emitEvent("magnet.item.click", item);
  const sys = useSystem();
  const winopt = sys._rootState.windowMap["Menulist"].get(item.title);
  if (winopt) {
    if (winopt._hasShow) {
      return;
    } else {
      winopt._hasShow = true;
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
</script>
<style lang="scss" scoped>
.desk-group {
  display: flex;
  flex-direction: column;
  flex-wrap: wrap;
  height: 100%;
  // 应用镂空效果
  color: transparent;
  /* 文字颜色设为透明 */
  text-shadow: 0 0 0.5px white, 0 0 1px black, 0 0 2px rgba(0, 0, 0, 2);
  /* 多层阴影 */

  // 重置子元素的默认样式
  >* {
    color: inherit;
    /* 继承颜色 */
    text-shadow: inherit;
    /* 继承描边效果 */
    font-size: 0.8rem;
  }
}

.swiperBox {
  width: 100vw;
  padding-top: vh(10);
  // position: absolute;
  // top:0;
  // left: 0;
  height: vh(560);

  .swiper-slide {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    grid-template-rows: repeat(5, vh(100))
  }
}

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
</style>
