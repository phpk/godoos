<script setup lang="ts">
import { ref } from "vue";
import { OpenDirDialog } from "@/util/goutil";
const formData = ref({
    importType: 'download',
    url: '',
    devPath:''
})
const addType = [
    {
        'name': '远程下载',
        'type': 'download'
    },
    {
        'name': '本地导入',
        'type': 'local'
    },
    {
        'name': '开发模式',
        'type': 'dev'
    }
]
function selectFile() {
  OpenDirDialog().then((res: string) => {
    formData.value.devPath = res;
  });
}
async function addAppByDownload() {

}
async function addAppByImport() {

}
async function addAppByDev() {

}
</script>
<template>
    <el-form label-width="auto" style="max-width: 600px;padding-left: 30px;">
        <el-form-item label="添加方式">
            <el-select v-model="formData.importType" style="width: 280px;">
                <el-option v-for="item in addType" :key="item.type" :label="item.name" :value="item.type" />
            </el-select>
        </el-form-item>
        <template v-if="formData.importType == 'download'">
            <el-form-item label="下载地址">
                <el-input v-model="formData.url" style="width: 280px;" />
            </el-form-item>
            <el-button type="primary" @click="addAppByDownload()">添加</el-button>
        </template>
        <template v-if="formData.importType == 'local'">
            <el-button type="primary" @click="addAppByImport()">导入</el-button>
        </template>
        <template v-if="formData.importType == 'dev'">
            <el-form-item label="本地路径">
                <el-input v-model="formData.devPath" style="width: 280px;" @click="selectFile()" />
            </el-form-item>
            <el-button type="primary" @click="addAppByDev()">添加</el-button>
        </template>
    </el-form>
</template>
<style scoped>
.el-button {
    margin:20px 100px
}
</style>