<template>
  <div class="file-tree-section">
    <div class="section-header" @click="fileTreeShow = !fileTreeShow">
      <el-icon v-if="fileTreeShow" class="header-icon" style="margin-right: 10px;">
        <ArrowDownBold />
      </el-icon>
      <el-icon v-else class="header-icon inactive" style="margin-right: 10px;">
        <ArrowRightBold />
      </el-icon>
      <icon name="diannao" :size="15" style="margin-right: 10px;"></icon>
      此电脑
    </div>
    <FileTree v-if="fileTreeShow" :currentPath="props.currentPath" :fileList="fileTreeList" 
      :navigateTo="props.navigateTo" />
  </div>
</template>
<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { useFileSystemStore } from '@/stores/filesystem'
const fileSystemStore = useFileSystemStore()
const props = defineProps({
  navigateTo: {
    type: Function,
    required: true
  },
  currentPath: {
    type: String,
    default: ''
  }
})
const fileTreeShow = ref(true)
const fileTreeList: any = ref([])
onMounted(async () => {
  fileTreeList.value = await fileSystemStore.getFilesInPath('/')
  //eventBus.on('refreshDesktop', handleRefreshDesktop);
})
</script>
<style lang="scss" scoped>
.file-tree-section {
  margin-bottom: 10px;

  .active {
    background-color: #cce8ff !important;
    border: 1px solid transparent;

    .header-icon {
      opacity: 1 !important;
    }

    &:hover {
      border: 1px solid #99d1ff;
    }
  }

  .section-header {
    display: flex;
    align-items: center;
    cursor: pointer;
    padding: 5px 0 5px 5px;
    border: 1px solid transparent;

    &:hover {
      background-color: #e5f3ff;

      .header-icon {
        opacity: 1;
      }
    }

    .header-icon {
      opacity: 0;
      font-size: 10px;
      transition: opacity 0.5s;

      &:hover {
        color: #3eccf8;
      }
    }

    .inactive {
      color: #8597a6;
    }
  }

  .navigation-item {
    display: flex;
    align-items: center;
    padding: 5px 5px 5px 30px;
    cursor: pointer;
    border: 1px solid transparent;
    transition: background-color 0.2s;

    &:hover {
      background-color: #e5f3ff;
    }

    .navigation-icon {
      margin-right: 8px;
    }

    .navigation-text {
      flex: 1;
    }
  }
}

</style>