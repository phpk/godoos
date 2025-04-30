<template>
  <div 
    v-if="contextMenuStore.isDesktopContextMenuVisible" 
    class="context-menu"
    :style="{ top: contextMenuStore.desktopContextMenuTop + 'px', left: contextMenuStore.desktopContextMenuLeft + 'px' }"
    ref="menuRef"
  >
    <ul>
      <li v-if="showDeleteMenu" @click="delFiles">
        <el-icon :size="16" color="#333">
          <Delete />
        </el-icon>
        删除(Del)
      </li>
      <li @click="newDir">
        <el-icon :size="16" color="#333">
          <FolderAdd />
        </el-icon>
        新建文件夹
      </li>
      <li @mouseenter="contextMenuStore.toggleSubMenu('newFile')" @mouseleave="contextMenuStore.toggleSubMenu(null)"
        class="has-submenu">
        <el-icon :size="16" color="#333">
          <CirclePlus />
        </el-icon>
        新建文件
        <div v-show="contextMenuStore.activeSubMenu === 'newFile'" class="sub-menu" style="top:-120px">
          <ul>
            <li @click="newFile('未命名文档', '.docx')">
              <icon name="word" :size="16" />
              Word文档
            </li>
            <li @click="newFile('未命名表格', '.xlsx')">
              <icon name="excel" :size="16" />
              数据表格
            </li>
            <li @click="newFile('未命名演示文稿', '.pptd')">
              <icon name="pptexe" :size="16" />
              演示文稿
            </li>
            <li @click="newFile('未命名文稿', '.md')">
              <icon name="markdown" :size="16" />
              Markdown
            </li>
            <li @click="newFile('未命名文件', '.txt')">
              <icon name="editorbt" :size="16" />
              文本文件
            </li>
            <li @click="newFile('未命名思维导图', '.mind')">
              <icon name="mindexe" :size="16" />
              思维导图
            </li>
            <li @click="newFile('未命名甘特图', '.gant')">
              <icon name="gant" :size="16" />
              甘特图
            </li>
            <li @click="newFile('未命名看板', '.kb')">
              <icon name="kanban" :size="16" />
              看板
            </li>
            <li @click="newFile('未命名白板', '.bb')">
              <icon name="baiban" :size="16" />
              白板
            </li>
            <li @click="newFile('未命名图片', '.pic')">
              <icon name="pic" :size="16" />
              图片
            </li>
          </ul>
        </div>
      </li>
      <li @click="fileSystemStore.refreshPaths()">
        <el-icon :size="16" color="#333">
          <Refresh />
        </el-icon>
        刷新
      </li>
      <li @click="selectAll()">
        <el-icon :size="16" color="#333">
          <Finished />
        </el-icon>
        全选(A)
      </li>
      <li @click="pasteFile" v-if="checkPaste()">
        <el-icon :size="16" color="#333">
          <Files />
        </el-icon>
        粘贴(V)
      </li>
      <li @click="copyFile" v-if="showDeleteMenu">
        <el-icon :size="16" color="#333">
          <CopyDocument />
        </el-icon>
        复制(C)
      </li>
      <li @click="cutFile" v-if="showDeleteMenu">
        <el-icon :size="16" color="#333">
          <Scissor />
        </el-icon>
        剪切(T)
      </li>
      <li @click="uploadFile">
        <el-icon :size="16" color="#333">
          <UploadFilled />
        </el-icon>
        上传(Q)
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { useContextMenuStore } from '@/stores/contextmenu'
import { useFileSystemStore } from '@/stores/filesystem'
import { useClickingStore } from '@/stores/clicking';
import { computed, watch, ref, nextTick } from 'vue';

const contextMenuStore = useContextMenuStore()
const fileSystemStore = useFileSystemStore()
const clickingStore = useClickingStore();

const showDeleteMenu = computed(() => {
  return clickingStore.clickedIcons.length > 0;
});

const menuRef = ref<HTMLElement | null>(null);

const adjustContextMenuPosition = () => {
  if (!menuRef.value) return;

  const menuRect = menuRef.value.getBoundingClientRect();
  const windowWidth = window.innerWidth;
  const windowHeight = window.innerHeight - 48;

  let top = contextMenuStore.desktopContextMenuTop;
  let left = contextMenuStore.desktopContextMenuLeft;

  // 调整顶部位置
  if (menuRect.bottom > windowHeight) {
    top -= (menuRect.bottom - windowHeight);
  }

  // 调整左侧位置
  if (menuRect.right > windowWidth) {
    left -= (menuRect.right - windowWidth);
  }

  // 确保顶部和左侧不为负数
  top = Math.max(0, top);
  left = Math.max(0, left);

  contextMenuStore.desktopContextMenuTop = top;
  contextMenuStore.desktopContextMenuLeft = left;
};

watch(() => contextMenuStore.isDesktopContextMenuVisible, (val) => {
  if (val) {
    // 确保在下一个 tick 中计算菜单位置
    nextTick(() => {
      adjustContextMenuPosition();
    });
  }
});

const newFile = async (name: string, ext: string) => {
  await fileSystemStore.handleNewFile(name, ext)
  contextMenuStore.isDesktopContextMenuVisible = false
}

const newDir = async () => {
  await fileSystemStore.handleNewDir('新建文件夹')
  contextMenuStore.isDesktopContextMenuVisible = false
}

const delFiles = async () => {
  contextMenuStore.isDesktopContextMenuVisible = false
  await fileSystemStore.handleDeleteFiles(clickingStore.clickedIcons)
}

const checkPaste = () => {
  return contextMenuStore.currentPath == fileSystemStore.currentPath && clickingStore.checkPaste()
}

const pasteFile = () => {
  contextMenuStore.isDesktopContextMenuVisible = false
  clickingStore.pasteFiles(fileSystemStore.currentPath)
}

const copyFile = () => {
  contextMenuStore.isDesktopContextMenuVisible = false
  clickingStore.copiedIcons = [...clickingStore.clickedIcons]
}

const cutFile = () => {
  contextMenuStore.isDesktopContextMenuVisible = false
  clickingStore.cutedIcons = [...clickingStore.clickedIcons]
}

const selectAll = async () => {
  contextMenuStore.isDesktopContextMenuVisible = false
  const fileList = await fileSystemStore.getFilesInPath(fileSystemStore.currentPath)
  clickingStore.clickedIcons = fileList.map((file: any) => file.path)
}
const uploadFile = () => {
  contextMenuStore.isDesktopContextMenuVisible = false
  fileSystemStore.uploadFile()
}
</script>

<style lang="scss" scoped>
@use '@/styles/contextmenu.scss';
</style>