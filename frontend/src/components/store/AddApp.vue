<script setup lang="ts">
import { ref, defineProps } from "vue";
import { OpenDirDialog } from "@/util/goutil";
import { getSystemKey, parseJson } from "@/system/config";
const apiUrl = getSystemKey("apiUrl");
import { notifySuccess, notifyError } from "@/util/msg";
import { t } from "@/i18n";
const props = defineProps({
    install: {
        type: Function,
        required: true,
    },
});
const formData = ref({
    importType: 'download',
    url: '',
    devPath: ''
})
const progress = ref(0)
const addType = [
    {
        'name': '远程下载',
        'type': 'download'
    },
    // {
    //     'name': '本地导入',
    //     'type': 'local'
    // },
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
    //console.log(item)
    if (formData.value.url == "") {
        notifyError(t("store.urlEmpty"))
        return;
    }
    const completion = await fetch(apiUrl + '/store/download?url=' + formData.value.url)
    if (!completion.ok) {
        notifyError(t("store.downloadError"))
    }
    //console.log(completion)
    const reader: any = completion.body?.getReader();
    if (!reader) {
        notifyError(t("store.cantStream"));
    }
    while (true) {
        const { done, value } = await reader?.read();
        if (done) {
            break;
        }
        // console.log(value)
        const json = await new TextDecoder().decode(value);
        //console.log(json)
        const res = parseJson(json)
        //console.log(res)
        if (res) {
            if (res.progress) {
                progress.value = res.progress
            }
            if (res.done) {
                notifySuccess(t("store.downloadSuccess"))
                progress.value = 0
                const pluginName = res.path.split("/").pop().split(".")[0]
                await installPlugin(pluginName)
                break;
            }
        }
    }
}
async function addAppByImport() {

}
async function addAppByDev() {
    if(formData.value.devPath && formData.value.devPath != ""){
        await installPlugin(formData.value.devPath)
    } 
}
async function installPlugin(pluginName: string) {
    const completion = await fetch(apiUrl + '/store/installInfo?name=' + pluginName)
    if (!completion.ok) {
        notifyError(t("store.installError"))
        return
    }
    const item = await completion.json()
    //console.log(item)
    item.isOut = true
    await props.install(item.data)
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
            <el-form-item label="下载进度" v-if="progress > 0">
                <el-progress :text-inside="true" :stroke-width="20" :percentage="progress" />
            </el-form-item>
            <el-button type="primary" @click="addAppByDownload()">添加</el-button>
        </template>
        <template v-if="formData.importType == 'local'">
            <el-button type="primary" @click="addAppByImport()">导入</el-button>
        </template>
        <template v-if="formData.importType == 'dev'">
            <el-form-item label="本地路径">
                <el-input v-model="formData.devPath" style="width: 280px;" placeholder="用户目录下.godoos/run，只填目录名称"
                    @click="selectFile()" />
            </el-form-item>
            <el-button type="primary" @click="addAppByDev()">添加</el-button>
        </template>
    </el-form>
</template>
<style scoped>
.el-button {
    margin: 20px 100px
}
</style>