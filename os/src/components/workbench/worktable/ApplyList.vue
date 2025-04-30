<template>
    <my-card ref="cardRef" :getData="getApplyList" :statusList="statusList" :columns="columns" :showPagination="true"
        :showFilter="true" :filterDate="select" :showAction="true" @reset="resetSelect">
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
            <el-button v-if="scope.nodeId == 1" size="small" @click="edit(scope)">编辑</el-button>
        </template>
    </my-card>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue';
import { getApplyList } from '@/api/workflow';
const page = ref(1);
const flowList: any = ref([])
const select: any = ref('')
const resetSelect = () => {
    select.value = ''
}

const cardRef = ref()
const props = defineProps<{
    reload: boolean
}>()
watch(() => props.reload, (newVal) => {
    if (newVal) {
        cardRef.value.reload()
    }
})
const columns = [
    { prop: 'id', label: 'ID' },
    { prop: 'flowName', label: '申请表' },
    { prop: 'status', label: '状态' },
    { prop: 'add_time', label: '发起时间' },
]
const statusList = {
    0: '未审批',
    1: '审批中',
    2: '审批通过',
    3: '驳回上一节点',
    4: '审批驳回初始阶段',
    9: '审批流程已全部通过',
}
const getFlowData = async (page: number) => {
    try {

        const res = await getApplyList(page)
        flowList.value = res.flowlist || []
    } catch (error) {
        console.log(error)
    }
}
const emit = defineEmits(['show'])
const view = (item: any) => {
    emit('show', { type: 'userview', data: item })
}
const edit = (item: any) => {
    emit('show', { type: 'useredit', data: item })
}
onMounted(async () => {
    await getFlowData(page.value)
})

</script>

<style lang="scss" scoped></style>