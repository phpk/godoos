<template>
    <my-card :getData="getTaskList" :statusList="statusList" :columns="columns" :showPagination="true"
        :showFilter="true" :showAction="true" @reset="resetSelect">
        <template #filter>
            <el-form-item>
                <el-select style="width: 150px;" v-model="select" placeholder="请选择" :empty-values="[null, undefined]">
                    <el-option value="" label="全部"></el-option>
                    <el-option v-for="item in flowList" :key="item.id" :label="item.name"
                        :value="item.flowId"></el-option>
                </el-select>
            </el-form-item>
        </template>
        <template #action="{ scope }">
            <el-button size="small" @click="view(scope)">查看</el-button>
        </template>
    </my-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { getTaskList } from '@/api/net/workflow';
const page = ref(1);
const flowList: any = ref([]);
const select = ref('');
const resetSelect = () => {
    select.value = ''
}

const statusList = {
    0: '未审批',
    1: '审批中',
    2: '审批通过',
    3: '驳回上一节点',
    4: '审批驳回初始阶段',
    9: '审批流程已全部通过',
}
const columns = [
    { prop: 'id', label: 'ID' },
    { prop: 'flowName', label: '任务表' },
    { prop: 'status', label: '状态' },
    { prop: 'add_time', label: '添加时间' },
]

const getFlowData = async (page: number) => {
    try {

        const res = await getTaskList(page)
        flowList.value = res.flowlist || []
    } catch (error) {
        console.log(error)
    }
}
const emit = defineEmits(['show'])
const view = (item: any) => {
    emit('show', { type: 'userview', data: item })
}
onMounted(async () => {
    await getFlowData(page.value)
})
</script>

<style scoped></style>