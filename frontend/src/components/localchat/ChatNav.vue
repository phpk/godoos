<template>
  <el-space direction="vertical" :size="20" class="win11-chat-nav">
    <div class="nav-item" v-for="item in store.navList" :key="item.index">
      <el-button
        :icon="item.icon"
        :class="store.navId === item.index ? 'active' : ''"
        dark
        circle
        @click="store.handleSelect(item.index)"
      />
    </div>
  </el-space>
</template>

<script setup lang="ts">
import { useLocalChatStore } from "@/stores/localchat";
const store = useLocalChatStore();
</script>
<style lang="scss" scoped>
.win11-chat-nav {
  height: 100vh;
  background-color: #f8f8f8; /* 使用更亮的淡灰色，更接近Win11的背景色 */
  border-right: 1px solid rgba(230, 230, 230, 0.5); /* 更浅的边框颜色 */
  padding: 16px;
  box-shadow: 2px 0 4px rgba(0, 0, 0, 0.1); 
  overflow-y: auto;
  display: flex;
  flex-direction: column;

  /* 添加全局字体样式以匹配Win11 */
  font-family: "Segoe UI", Tahoma, Geneva, Verdana, sans-serif;
}

.win11-chat-nav .nav-item {
  /* 假定每个按钮的基类 */
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  padding: 3px;
  border-radius: 50%; /* 圆角 */
  background-color: white;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1); /* 渐变阴影 */
  transition: all 0.2s ease-in-out;

  &:hover {
    background-color: #f0f0f0; /* 鼠标悬停时的轻微颜色变化 */
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
  }

  &.active {
    /* 加强背景色对比，使用Win11的强调色或品牌色 */
    background-color: #4579a1; /* 示例颜色，可根据设计调整 */
    color: white; /* 文字颜色反转，确保可读性 */

    /* 增加外边框以进一步区分 */
    border: 2px solid #005a9c; /* 较深的强调色作为边框 */
    border-radius: 50%;
    /* 内发光效果，让按钮看起来更‘活跃’ */
    box-shadow: 0 5px 8px rgba(0, 120, 212, 0.5) inset;

    /* 微动效，当状态改变时给予用户反馈 */
    transform: scale(1.02);
    transition: transform 0.2s cubic-bezier(0.2, 0.4, 0.6, 1);

    /* 确保文字在按下时不会因按钮尺寸变化而偏移 */
    transition-property: background-color, box-shadow, transform, color;

    /* 如果按钮包含图标，可以考虑为图标也添加强调效果，例如改变颜色 */
    .el-icon {
      color: inherit; /* 或指定特定强调色 */
    }

    /* 为了平滑的过渡，确保所有相关属性都被包含在transition中 */
  }
  /* 考虑到按钮是圆形且使用了 `circle` 属性，确保图标和背景颜色调整得当 */
  &.active .el-button {
    background-color: transparent !important; /* 确保背景色不影响图标颜色 */
  }

  /* 图标颜色调整，确保在active状态下足够突出 */
  &.active .el-icon {
    color: #ffffff; /* 确保图标颜色与背景对比鲜明 */
  }

  /* 非活动状态的悬停效果，保持与.active状态的区分 */
  .nav-item:hover:not(.active) {
    /* 调整以与.active状态区分，例如使用较浅的颜色 */
    background-color: #eaeaea;
  }
}
</style>