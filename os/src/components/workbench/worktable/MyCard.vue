<template>
    <!-- 申请列表 -->
    <el-card>
        <template #header>
            <slot name="title"></slot>
            <el-form :inline="true" v-if="showFilter">
                <slot name="filter"></slot>
                <!-- <el-form-item> -->
                <!-- <el-select style="width: 150px;" v-model="select">
                        <el-option label="全部" value="0"></el-option>
                        <el-option v-for="item in flowList" :key="item.id" :label="item.name"
                            :value="item.flowId"></el-option>
                    </el-select> -->
                <!-- </el-form-item> -->
                <el-form-item>
                    <el-button @click="search" icon="Search" type="primary" plain>查询</el-button>
                    <el-button @click="reset" icon="RefreshLeft" type="primary" plain>重置</el-button>
                </el-form-item>
            </el-form>
        </template>
        <el-scrollbar>
            <el-table v-loading="isLoading" :data="data" style="width: 100%;"
                :show-overflow-tooltip="isMobileDevice() ? false : true">
                <el-table-column v-for="item in columns" :key="item.prop" :prop="item.prop" :label="item.label"
                    align="center"></el-table-column>
                <el-table-column label="操作" align="center" v-if="showAction">
                    <template #default="{ row }">
                        <slot name="action" :scope="row"></slot>
                    </template>
                </el-table-column>
            </el-table>
        </el-scrollbar>

        <template #footer v-if="showPagination">
            <el-pagination :total="total" :current-page="page" layout="prev, pager, next" background size="small"
                @current-change="handleChange"></el-pagination>
        </template>
    </el-card>
</template>

<script setup lang="ts">
import { isMobileDevice } from '@/utils/device';
import { onMounted, ref } from 'vue';

const props = defineProps<{
    getData: (a?: any, b?: any) => any,
    statusList?: StatusType,
    showPagination?: boolean,
    showFilter?: boolean,
    filterDate?: any,
    showAction?: boolean,
    columns: any[]
}>()
const page = ref(1);
const data: any = ref([])
const total = ref(0)
const isLoading = ref(false)

type StatusType = { [key: number]: string }
// 处理时间戳函数
const handleTime = (time: number) => {
    const date = new Date(time * 1000)
    return date.toLocaleString()
}

const getTableData = async (page: number, query?: any) => {
    isLoading.value = true
    try {
        let res: any
        if (query) {
            console.log(query)
            res = await props.getData(page, query)
        } else {
            res = await props.getData(page)
        }
        if (res.datalist) {
            res.datalist.forEach((item: any) => {
                item.add_time = handleTime(item.add_time)
                if (props.statusList) {
                    item.status = props.statusList[item.status]
                }
            })
        }
        data.value = res.datalist || res.list
        total.value = res.total || 0
    } catch (error) {
        console.log(error)
    }
    isLoading.value = false
}

const handleChange = async (newPage: number) => {
    page.value = newPage
    if (props.filterDate) {
        await getTableData(page.value, props.filterDate)
    } else {
        await getTableData(page.value)
    }
}
const search = async () => {
    await getTableData(page.value, props.filterDate)
}
const emit = defineEmits(['reset'])
const reset = async () => {
    await emit('reset')
    page.value = 1
    await getTableData(page.value, props.filterDate)
}

const reload = async () => {
    if (props.filterDate) {
        await getTableData(page.value, props.filterDate)
    } else {
        await getTableData(page.value)
    }
}

onMounted(async () => {
    if (props.filterDate) {
        await getTableData(page.value, props.filterDate)
    } else {
        await getTableData(page.value)
    }
})
defineExpose({
    reload
})

</script>

<style lang="scss" scoped>
@use '@/styles/card.scss';
</style>