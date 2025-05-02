<template>
    <el-tabs>
        <el-tab-pane label="全部">
            <div class="cat-container">
                <div class="cat-items" v-for="item in allApllyList" :key="item.id" @click="add(item)">
                    <el-icon class="cat-icon">
                        <Star />
                    </el-icon>
                    <span>{{ item.name }}</span>
                </div>
            </div>
        </el-tab-pane>
        <el-tab-pane v-for="item in catList" :key="item.id" :label="item.name">
            <div class="cat-container" v-for="w in item.workflow_list" :key="w.id" @click="add(w)">
                <div class="cat-items">
                    <el-icon class="cat-icon">
                        <Star />
                    </el-icon>
                    <span>{{ w.name }}</span>
                </div>
            </div>
        </el-tab-pane>
    </el-tabs>
</template>

<script setup lang="ts">
import { getMyApply } from '@/api/net/workflow';
import { onMounted, ref } from 'vue';

const catList: any = ref([])
const allApllyList: any = ref([])
const getData = async () => {
    try {
        const res = await getMyApply();
        if (res.list && res.list.length > 0) {
            res.catList.forEach((item: any) => {
                item.workflow_list = res.list.filter((i: any) => i.cateId === item.id);
            });
        }
        catList.value = res.catList || [];
        allApllyList.value = res.list || [];
    } catch (error) {
        console.log(error);
    }
}

const emit = defineEmits(['show'])
const add = (item: any) => {
    console.log(item)
    emit('show', { type: 'useradd', data: item })
}
onMounted(async () => {
    await getData();
})
</script>

<style scoped lang="scss">
@use '@/styles/tabs.scss';
</style>