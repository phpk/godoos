<template>
    <el-scrollbar height="550px">
        <FileTreeSection :navigateTo="navigateTo" :currentPath="currentPath"/>
    </el-scrollbar>
    
</template>
<script setup lang="ts">
import { useFileSystemStore } from '@/stores/filesystem'
import { onMounted } from 'vue';
const fsStore = useFileSystemStore();
const props = defineProps({
    win: {
        type: Object,
        required: true,
    },
});
const win = props.win;
const navigateTo = (file: any) => {
    console.log(file.path);

};

const currentPath = fsStore.currentPath;
onMounted(() => {
    //tabsStore.addTab(win);
});
</script>
<style scoped>
.filetree-area{
    margin: 0px;
    padding: 0px;
    height: 100%;
    overflow-y: auto;
}
</style>