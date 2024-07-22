<template>
  <div class="dialog">
    <div class="dialog-content">
      <div class="dialog-icon">
        <!-- <img class="dialog-icon_img" :src="iconMap[win?.config.option.type]" alt="" /> -->
        <svg class="icon" aria-hidden="true" style="font-size: 1.2em">
          <use :xlink:href="'#icon-' + iconMap[win?.config.option.type]"></use>
        </svg>
      </div>
      {{ win?.config.option.message }}
    </div>
    <div class="dialog-button">
      <WinButton
        v-for="(item, index) in win?.config.option.buttons"
        :key="item"
        @click="handleClick(index)"
        >{{ item }}</WinButton
      >
    </div>
  </div>
</template>
<script setup lang="ts">
import { inject } from "vue";
import { BrowserWindow } from "@/system/window/BrowserWindow";
// import errorIcon from '@/assets/error-icon.png';
// import infoIcon from '@/assets/info-icon.png';
// import questionIcon from '@/assets/question-icon.png';
// import warningIcon from '@/assets/warning-icon.png';

const iconMap: {
  [key: string]: string;
} = {
  error: "error",
  info: "info",
  question: "question",
  warning: "warning",
};

const win: BrowserWindow | undefined = inject<BrowserWindow>("browserWindow");

function handleClick(index: number) {
  if (win) {
    win.config.res({
      response: index,
    });
    win.close();
  }
}
</script>
<style lang="scss" scoped>
.dialog {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  background-color: var(--theme-main-color); // 使用主题颜色作为对话框的背景
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1); // 添加柔和的阴影效果，类似Windows 11

  .dialog-content {
    width: 90%;
    max-width: 500px; // 添加最大宽度，让对话框看起来更像Windows 11的弹窗
    height: auto;
    font-size: var(--ui-font-size);
    display: flex;
    justify-content: center;
    align-items: center;
    background-color: var(--theme-main-color-light); // 使用较亮的主题颜色作为内容区域背景
    border-radius: 8px; // 添加圆角，符合Windows 11的设计风格
    padding: 20px; // 添加内边距，提高可读性

    .dialog-icon {
      width: 100px;
      height: 60px;
      display: flex;
      justify-content: center;
      align-items: center;

      .dialog-icon_img {
        width: 34px;
        height: 34px;
        fill: var(--color-blue); // 使用Windows 11的蓝色
      }
    }
  }

  .dialog-button {
    width: 100%;
    height: 40px;
    display: flex;
    flex-shrink: 0;
    justify-content: flex-end;
    align-items: center;
    background-color: var(--color-gray); // 使用较深的灰色作为按钮背景

    .win-button {
      margin: 0 10px;
      padding: 0 15px;
      height: 100%;
      border-radius: 4px; // 添加圆角
      font-size: var(--ui-font-size);
      color: var(--color-white); // 使用白色文字
      border: none;
      cursor: pointer;
      transition: background-color 0.2s, color 0.2s; // 添加过渡效果

      &:hover {
        background-color: var(--color-gray-hover);
        color: var(--color-white);
      }

      &:active {
        background-color: var(--color-gray-active);
        color: var(--color-white);
      }
    }
  }
}
</style>
