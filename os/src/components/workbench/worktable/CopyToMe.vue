<template>
    <!-- <div>
        <iframe src="/views/desktop/copytome.html" frameborder="0" style="width: 100%; height: 100%;"></iframe>
        <el-form :inline="true">
            <el-form-item label="日期范围">
                <el-date-picker v-model="dateRange" value-format="YYYY-MM-DD" type="daterange" start-placeholder="开始日期"
                    end-placeholder="结束日期" style="width: 200px;" />
            </el-form-item>
        </el-form>
    </div> -->
    <my-card :getData="getMyCopyList" :filterDate="dateRange" :showAction="true" :showFilter="true" :columns="columns"
        :showPagination="true" @reset="resetDate">
        <template #filter>
            <el-form-item label="日期范围">
                <el-date-picker v-model="dateRange" value-format="YYYY-MM-DD" type="daterange" start-placeholder="开始日期"
                    end-placeholder="结束日期" style="width: 200px;" />
            </el-form-item>
        </template>
        <template #action="{ scope }">
            <el-button size="small" @click="view(scope)">查看流程信息</el-button>
        </template>
    </my-card>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { getMyCopyList } from '@/api/net/workflow';

const columns = [
    { prop: 'Id', label: 'ID' },
    { prop: 'from_nickname', label: '发送者' },
    { prop: 'flowName', label: '流程名' },
    { prop: 'add_time', label: '发送时间' },
]
// 获取本月-上月范围
const getRangeTime = () => {
    let year
    let month
    const date = new Date()
    const nowYear = date.getFullYear()
    const nowMonth = date.getMonth()
    const nowDate = date.getDate()
    if (nowMonth === 0) {
        year = nowYear - 1
        month = 12
    } else {
        year = nowYear
        month = nowMonth
    }
    return {
        thisMonth: `${nowYear}-${(nowMonth + 1) > 9 ? (nowMonth + 1) : "0" + (nowMonth + 1)}-${nowDate > 9 ? nowDate : "0" + nowDate}`,
        lastMonth: `${year}-${month > 9 ? month : "0" + month}-${nowDate > 9 ? nowDate : "0" + nowDate}`
    }
}
const dateRange = ref([getRangeTime().lastMonth, getRangeTime().thisMonth])
const resetDate = () => {
    dateRange.value = [getRangeTime().lastMonth, getRangeTime().thisMonth]
}

const emit = defineEmits(['show'])
const view = (item: any) => {
    emit('show', { type: 'userview', data: item })
}
</script>

<style scoped></style>