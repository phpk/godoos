<template>
    <div class="desktop" :style="{ backgroundColor: background }">
        <DesktopBackground class="userarea-upper"></DesktopBackground>
        <template v-if="backgroundType == 'image'">
        </template>
        <div class="dest-item userarea-upper">
             <!-- Swiper -->
        <swiper :modules="modules" :loop="true" :slides-per-view="1" :space-between="20"
            :pagination="{ type: 'bullets', clickable: true }" class="swiperBox" @swiper="onSwiper"
            @slideChange="onSlideChange">
            <swiper-slide class="swiper-slide">
                <MobileApp v-for="item in appList" :key="item.name" :item="item" class="item"></MobileApp>
            </swiper-slide>
            <swiper-slide></swiper-slide>
            <swiper-slide></swiper-slide>
        </swiper>
        </div>
        <BottomBar :list="appList"></BottomBar>
    </div>
</template>

<script lang="ts" setup>
import { Bios, System, useSystem } from '@/system';
import { ref, onMounted } from 'vue';
import { RootState } from '@/system/root.ts';
import { Swiper, SwiperSlide } from 'swiper/vue'
// 引入swiper样式
import 'swiper/swiper-bundle.css';
// 引入swiper核心和所需模块
import { Pagination } from 'swiper/modules'

import { useAppOpen } from "@/hook/useAppOpen";
const { openapp, appList } = useAppOpen("apps");
console.log(appList,'appList======')
console.log(openapp,'openapp=======')
// 在modules加入要使用的模块
const modules = [Pagination]
const onSwiper = (swiper: any) => {
    console.log(swiper);
};
// 更改当前活动swiper
const onSlideChange = (swiper: any) => {
    // swiper是当前轮播的对象，里面可以获取到当前swiper的所有信息，当前索引是activeIndex
    console.log(swiper.activeIndex)
}
const mobileref = ref();
const rootState = ref<RootState | undefined>(useSystem()?._rootState);
const backgroundType = ref("color");
const background: any = ref('#3A98CE');
// const
onMounted(() => {
    refershBack(rootState.value?.options.background);
});

function refershBack(val: string | undefined) {
    background.value = val || "#3A98CE";
    if (background.value || background.value.startsWith("/image/")) {
    backgroundType.value = "image";
  } else {
    backgroundType.value = "color";
  }
}

Bios.onOpen((system: System) => {
    rootState.value = system._rootState;
    system.rootRef = mobileref.value;
});
</script>

<style lang="scss" scoped>
*{
    padding: 0;
    margin: 0;
    box-sizing: border-box;
}
.desktop {
    position: relative;
    background-color: #3A98CE;
    background: url('@/pubic/image/bg/bg6.jpg');
    height: 100vh;
    overflow: hidden;
    .userarea-upper {
      position: absolute;
      top: 0;
      left: 0;
    }
    .dest-item{
        z-index: 2;
        .swiperBox {
        width: 100vw;
        padding:vw(20);
        // position: absolute;
        // top:0;
        // left: 0;
        height: vh(560);
        .swiper-slide{
            display: grid;
            grid-template-columns: repeat(4, 1fr);
            grid-template-rows: repeat(4,vh(100))
        }
    }
    }

}
</style>