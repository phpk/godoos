<script setup lang="ts">
import { ref, onMounted, inject } from 'vue';
import { getSystemKey } from "@/system/config";
import {t} from "@/i18n";
import { BrowserWindow } from '@/system/window/BrowserWindow';
const win = inject<BrowserWindow>('browserWindow');
const processList = ref([]);
const apiUrl = getSystemKey("apiUrl");
let timerId: ReturnType<typeof setTimeout>; // 用于保存定时器ID

onMounted(() => {
    fetchProcesses();
});
win?.on("close", () => {
    //console.log(timerId)
    clearTimeout(timerId);
})
const fetchProcesses = async () => {
    try {
        const response = await fetch(apiUrl + '/store/listport');
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        const data = await response.json();
        processList.value = data.processes;
        timerId = setTimeout(() => {
            fetchProcesses()
        }, 3000);
    } catch (error) {
        console.error('Error fetching processes:', error);
    }
};

const killProcess = async (name: string) => {
    try {
        const response = await fetch(apiUrl + '/store/killport?name=' + name);
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        fetchProcesses(); // 刷新列表
    } catch (error) {
        console.error('Error killing process:', error);
    }
};
</script>
<template>
    <el-table :data="processList" stripe style="width: 96%;margin:auto;border:none" max-height="550">
        <el-table-column prop="port" :label="t('process.port')" sortable />
        <el-table-column prop="pid" :label="t('process.pid')" sortable />
        <el-table-column prop="name" :label="t('process.name')" />
        <el-table-column prop="proto" :label="t('process.proto')" />
        <el-table-column :label="t('process.action')">
            <template #default="{ row }">
                <el-button circle icon="Remove" @click="killProcess(row.name)" style="width: 35px;height: 35px;" />
            </template>
        </el-table-column>
    </el-table>
</template>