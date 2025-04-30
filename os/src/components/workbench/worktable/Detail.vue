<template>
    <el-tabs v-model="activeName" style="height: 100%;">
        <el-tab-pane label="内容" name="content">
            <iframe :src="contentScr" frameborder="0" style="width: 100%;height: 100%;"></iframe>
        </el-tab-pane>
        <el-tab-pane label="流程" name="flow">
            <iframe :src="flowSrc" frameborder="0" style="width: 100%;height: 100%;"></iframe>
        </el-tab-pane>
    </el-tabs>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
const props = defineProps<{
    id: any
}>()
const activeName = ref('content')
const contentScr = computed(() => {
    return '/static/form/index.html?act=userview&dataId=' + props.id
})
const flowSrc = computed(() => {
    return '/static/workflow/index.html/#/view?id=' + props.id
})

</script>

<style scoped lang="scss">
@use '@/styles/tabs.scss';
</style>