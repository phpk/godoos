<template>
  <Transition name="fade">
    <div v-if="isPopShow" class="message scroll-bar" @mousedown.stop="">
      <div class="notify-center">
        <div class="message-title">
          共{{ store.page.total }}条通知
        </div>
        <div class="message-group scroll-bar">
          <div class="message-item" v-for="notify in store.notifyList" :key="notify.id">
            <div class="message-item-title" @click="store.viewContent(notify)">
              <span>{{ notify.title }}</span>
            </div>
          </div>
          <el-row justify="center" style="margin-top: 15px" v-if="store.page.pages > 1">
            <el-pagination background layout="prev, pager, next" v-model:current-page="store.page.current"
              v-model:page-size="store.page.size" :total="store.page.total"
              @current-change="(val: any) => store.pageClick(val)" />
          </el-row>
        </div>
      </div>
    </div>
  </Transition>
  <ShowNews />
</template>
<script setup lang="ts">
import { useNotifyStore } from '@/stores/notify';
import { ref } from 'vue';
import { mountEvent } from '@/system/event';
const store = useNotifyStore()
const isPopShow = ref(false);

mountEvent('messagecenter.show', () => {
  isPopShow.value = !isPopShow.value;
});
mountEvent('messagecenter.hidden', () => {
  isPopShow.value = false;
});
onMounted(async () => {
  // store.init();
  await store.getList()
});
</script>
<style lang="scss" scoped>
@import '@/assets/main.scss';

.message {
  position: absolute;
  top: 0;
  right: 0;
  width: 300px;
  height: 100%;
  z-index: 400;
  background-color: #F5F5F5;
  /* 更接近Win10的背景颜色 */
  border-left: 1px solid #E5E5E5;
  /* 边框颜色调整 */
  overflow-y: auto;
  user-select: none;
  box-sizing: content-box;
  display: flex;
  flex-direction: column;

  .message-title {
    padding: 10px 20px;
    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
    /* 使用Win10默认字体 */
    font-size: 14px;
    /* 字体大小调整 */
    font-weight: bold;
    border-bottom: 1px solid #E5E5E5;
    /* 边框颜色调整 */
    display: flex;
    justify-content: space-between;
    align-items: flex-end;

    .allclear {
      font-size: 12px;
      /* 字体大小调整 */
      cursor: pointer;

      &:hover {
        color: #0078D4;
        /* Win10主题蓝色 */
      }
    }
  }

  .notify-center {
    height: 100%;
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }

  .message-group {
    height: 100%;
    overflow: auto;
  }

  .message-item {
    padding: 15px 16px;
    background-color: #FFFFFF;
    /* 消息项背景色调整 */
    width: 100%;
    /* 设置为100%，以便内容自适应 */
    overflow: hidden;
    margin-bottom: 8px;
    /* 增加间距 */
    border: 1px solid #E5E5E5;
    /* 边框颜色调整 */
    border-radius: 4px;
    /* 添加圆角 */
    transition: all 0.2s ease;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
    /* 添加阴影效果 */

    .message-item-title {
      font-size: 14px;
      font-weight: bold;
      margin-bottom: 2px;
      text-overflow: ellipsis;
      overflow: hidden;
    }

    .message-item-body {
      font-size: 13px;
      /* 字体大小调整 */
    }

    &:hover {
      border-color: #C7C7C7;
      /* 鼠标悬停时边框颜色变淡 */
    }
  }
}

.fade-enter-active,
.fade-leave-active {
  transition: all 0.3s ease;
  /* 修改过渡效果 */
}

.fade-enter-to,
.fade-leave-from {
  opacity: 1;
  transform: translateX(0);
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateX(100%);
}

</style>