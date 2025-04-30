<template>
    <div style="width: 100%;height: 100%;">
        <my-task v-if="store.currentTab === '0'" @show="showDrawer"></my-task>
        <my-apply v-if="store.currentTab === '1'" @show="showDrawer"></my-apply>
        <task-list v-if="store.currentTab === '2'" @show="showDrawer"></task-list>
        <apply-list :reload="reloadApplyList" v-if="store.currentTab === '3'" @show="showDrawer"></apply-list>
        <copy-from-me v-if="store.currentTab === '4'" @show="showDrawer"></copy-from-me>
        <my-check :reload="reloadCheck" v-if="store.currentTab === '5'" @show="showDrawer"></my-check>
        <copy-to-me v-if="store.currentTab === '6'" @show="showDrawer"></copy-to-me>
        <el-drawer v-model="isShow" size="95%" :title="title">
            <detail :id="id" v-if="showDetail"></detail>
            <iframe v-else :src="iframeScr" frameborder="0" style="width: 100%;height: 95%;"></iframe>
        </el-drawer>
    </div>
</template>

<script setup lang="ts">
import { useWorkTableStore } from '@/stores/worktable';
import { ref, onMounted, watch } from 'vue';
const store = useWorkTableStore()
const isShow = ref(false)
const showDetail = ref(false)
const iframeScr = ref('')
const title = ref('')
const id = ref('')

watch(isShow, (newVal) => {
    if (!newVal) {
        iframeScr.value = ''
        showDetail.value = false
    }
})
const showDrawer = (obj: any) => {
    isShow.value = true
    if (obj.type == 'userview') {
        title.value = '查看数据-' + obj.data.flowName
        showDetail.value = true
        id.value = obj.data.id
        // iframeScr.value = '/static/form/index.html?act=userview&dataId=' + obj.data.id
    } else if (obj.type == 'useradd') {
        title.value = '新增数据-' + obj.data.name
        iframeScr.value = '/static/form/index.html?act=useradd&flowId=' + obj.data.id
    } else if (obj.type == 'check') {
        title.value = '审批数据-' + obj.data.flowName
        iframeScr.value = '/views/desktop/mycheckaction.html?formdataId=' + obj.data.dataId + '&flowcheckId=' + obj.data.id + '&id=' + obj.data.id
    } else if (obj.type == 'useredit') {
        title.value = '编辑数据-' + obj.data.flowName
        iframeScr.value = '/static/form/index.html?act=useredit&dataId=' + obj.data.id
    }
}

const reloadCheck = ref(false)
const reloadApplyList = ref(false)
onMounted(() => {
    window.addEventListener('message', (e) => {
        reloadCheck.value = false
        reloadApplyList.value = false
        if (e.data == 'reloadCheck') {
            reloadCheck.value = true
        } else if (e.data == 'reloadApplyList') {
            reloadApplyList.value = true
        }
        if (e.data == 'closeDialog') {
            isShow.value = false
        }
    }, false)
})
</script>

<style lang="scss" scoped>
:deep(.el-drawer__header) {
    border-bottom: 1px solid #ccc;
    margin-bottom: 0;
    padding: 10px;
}

:deep(.el-drawer__body) {
    padding: 0;
}
</style>