<template>
    <my-card ref="cardRef" :getData="getMyCheckList" :columns="columns" :showPagination="true" :showFilter="false"
        :showAction="true">
        <template #title>我的审核</template>
        <template #action="{ scope }">
            <el-button size="small" @click="view(scope)">开始审核</el-button>
        </template>
    </my-card>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { getMyCheckList } from '@/api/net/workflow'
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
    { prop: 'userName', label: '发起人' },
    { prop: 'flowName', label: '流程名' },
    { prop: 'add_time', label: '触发时间' },
]
const emit = defineEmits(['show'])
const view = (item: any) => {
    emit('show', { type: 'check', data: item })
}

</script>
