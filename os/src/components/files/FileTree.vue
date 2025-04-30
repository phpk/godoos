<template>
    <div class="item-group" v-for="file in fileList" :key="file.id">
        <div class="file-item" :class="{ 'active': currentPath == file.path }"
            :style="{ paddingLeft: `${10 + level * 12}px` }" @click="openFile(file)" @mousedown.stop>
            <el-icon v-if="file.isOpen" class="arrow" style="margin-right: 10px;" @click.stop="onOpenArrow(file)">
                <ArrowDownBold />
            </el-icon>
            <el-icon v-else class="arrow inactive" style="margin-right: 10px;" @click.stop="onOpenArrow(file)">
                <ArrowRightBold />
            </el-icon>
            <icon :size="20" class="file-icon" :name="dealIcon(file)" />
            <span class="file-name">{{ dealSystemName(file.name) }}</span>
        </div>
        <div class="sub-tree">
            <FileTree v-if="file.isOpen" :file-list="file.subFileList" :level="level + 1" :navigate-to="navigateTo"
                :current-path="currentPath" :open-file="openFile" />
        </div>
    </div>
</template>

<script setup lang="ts">
import { dealIcon } from '@/utils/icon'
import { dealSystemName } from '@/i18n'
import { useFileSystemStore } from '@/stores/filesystem'
const fileSystemStore = useFileSystemStore();


const props = defineProps({
    fileList: {
        type: Array as () => any[],
        required: true
    },
    level: {
        type: Number,
        default: 0,
    },
    navigateTo: {
        type: Function,
        required: true
    },
    currentPath: {
        type: String,
        default: ''
    }
})
const openFile = (item: any) => {
    if (item.isDirectory) {
        props.navigateTo(item.path)
    }else{
        //props.openFile(item)
        fileSystemStore.openFile(item)
    }
}

// 打开子级
const onOpenArrow = async (item: any) => {
    if (item.isOpen && !item.subFileList?.length) {
        return;
    }
    item.isOpen = !item.isOpen;
    const res = await fileSystemStore.getFilesInPath(item.path) || [];
    //item.subFileList = res.filter((item: any) => item.isDirectory)
    item.subFileList = res;
}
</script>

<style lang="scss" scoped>
.file-item {
    display: flex;
    align-items: center;
    padding: 5px 5px 5px 10px;
    cursor: pointer;
    border: 1px solid transparent;
    transition: background-color 0.2s;

    &:hover {
        background-color: #e5f3ff;

        .arrow {
            opacity: 1;
        }
    }

    .arrow {
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

    .file-icon {
        margin-right: 8px;
    }

    .file-name {
        flex: 1;
    }
}

.active {
    background-color: #cce8ff !important;
    border: 1px solid transparent;

    .arrow {
        opacity: 1 !important;
    }

    &:hover {
        border: 1px solid #99d1ff;
    }
}
</style>