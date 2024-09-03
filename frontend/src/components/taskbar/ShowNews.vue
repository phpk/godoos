<template>
    <el-dialog v-model="store.showNews" :title="store.newsDetail.title" show-close="true" width="600">
    <div class="new-content">
      <div class="news-detail" v-html="store.newsDetail.content"></div>
      <p class="new-time"><strong>发布时间:</strong> {{ store.newsDetail.showtime }}</p>
    </div>
  </el-dialog>
</template>
<script setup lang="ts">
import { useNotifyStore } from '@/stores/notify';
const store:any = useNotifyStore()
</script>
<style lang="scss" scoped>

.new-content {
  padding: 20px;
  background-color: #F5F5F5;
  border: 1px solid #E5E5E5;
  border-radius: 4px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;

  .new-time {
    margin-bottom: 10px;
    font-size: 14px;
    color: #666;

    strong {
      font-weight: bold;
      color: #333;
    }
  }

  .news-detail {
    font-size: 16px;
    line-height: 1.5;
    color: #333;
    white-space: pre-wrap;
  }
  .news-detail:deep(img){
    max-width: 90%;
  }

}
</style>