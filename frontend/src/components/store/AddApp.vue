<script setup lang="ts">
import { ref, defineProps } from "vue";
import { OpenDirDialog } from "@/util/goutil";
import { getSystemKey, parseJson } from "@/system/config";
import { useStoreStore } from "@/stores/store";
import { notifySuccess, notifyError } from "@/util/msg";
import { t } from "@/i18n";
const apiUrl = getSystemKey("apiUrl");
const store = useStoreStore()
const fileInput = ref()
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
        name: t('store.remoteDownload'),
        type: 'download'
    },
    {
        name: t('store.localImport'),
        type: 'local'
    },
    {
        name: t('store.devMode'),
        type: 'dev'
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
    fileInput.value.click()
}
async function upload(e: any) {
    const file = e.target.files[0];
    if (file) {
        const formData = new FormData();
        formData.append('files', file);
        const completion = await fetch(apiUrl + '/store/upload', {
            method: 'POST',
            body: formData
        })
        if (!completion.ok) {
            notifyError(t("store.importError"))
            return
        }
        const item = await completion.json()
        //console.log(item)
        if (item.data) {
            await installPlugin(item.data)
        }
    }
}
async function addAppByDev() {
    if (formData.value.devPath && formData.value.devPath != "") {
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
    item.data.isOut = true
    await props.install(item.data)
    store.changeCate(0, 'hots')
}
</script>
<template>
    <el-form label-width="auto" style="max-width: 600px;padding-left: 30px;">
        <el-form-item :label="t('store.addMethod')">
            <el-select v-model="formData.importType" style="width: 280px;">
                <el-option v-for="item in addType" :key="item.type" :label="item.name" :value="item.type" />
            </el-select>
        </el-form-item>
        <template v-if="formData.importType == 'download'">
            <el-form-item :label="t('store.downloadUrl')">
                <el-input v-model="formData.url" style="width: 280px;" />
            </el-form-item>
            <el-form-item :label="t('store.downloadProgress')" v-if="progress > 0">
                <el-progress :text-inside="true" :stroke-width="20" :percentage="progress" />
            </el-form-item>
            <el-button type="primary" @click="addAppByDownload()">{{ t('store.add') }}</el-button>
        </template>
        <template v-if="formData.importType == 'local'">
            <input type="file" style="display: none;" ref="fileInput" accept=".zip,.tar,.tar.gz,.tar.bz2"
                @change="upload">
            <el-button type="primary" @click="addAppByImport()">{{ t('store.importZip') }}</el-button>
        </template>
        <template v-if="formData.importType == 'dev'">
            <el-form-item :label="t('store.localPath')">
                <el-input v-model="formData.devPath" style="width: 280px;" :placeholder="t('store.userDirectory')"
                    @click="selectFile()" />
            </el-form-item>
            <el-button type="primary" @click="addAppByDev()">{{ t('store.add') }}</el-button>
        </template>
    </el-form>
</template>
<style scoped>
.el-button {
    margin: 20px 100px
}
</style>