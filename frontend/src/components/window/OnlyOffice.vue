<template>
    <DocumentEditor id="docEditor" :documentServerUrl="config.onlyoffice.url" :config="editorConfig"
        :events_onDocumentReady="onDocumentReady" :onLoadComponentError="onLoadComponentError" />
</template>

<script lang="ts" setup>
import { DocumentEditor } from "@onlyoffice/document-editor-vue";
import { getSystemConfig, getSplit } from "@/system/config";
import { BrowserWindow } from "@/system";
import { ref, onMounted, inject } from "vue";

const config = getSystemConfig();
const win: any = inject<BrowserWindow>("browserWindow");

const props = defineProps({
    onlyType: {
        type: String,
        default: "",
    },
    eventType: {
        type: String,
        default: "",
    },
    ext: {
        type: String,
        default: "docx",
    },
});

const editorConfig: any = ref({});
const SP = getSplit();


onMounted(async () => {
    const path = win?.config?.path;
    let title = "未命名文档"
    if (path != "") {
        title = path.split(SP).pop();
    }
    const apiUrl = config.storenet.url || config.apiUrl
    const fileKey = "docx" + Math.random()
    editorConfig.value = {
        document: {
            fileType: props.ext, // 根据文件扩展名动态设置 fileType
            key:fileKey,
            title: title,
            url: apiUrl + "/file/readfile?stream="+fileKey+"&path=" + path,
            "info": {
                "owner": "GodoOS", 
            },
        },
        documentType: props.onlyType,
        editorConfig: {
            lang: "zh",
            callbackUrl: apiUrl + "/file/onlyofficecallback",
            "user": { //用户信息
                "id": "godoos", //用户ID
                "name": "写作员" //用户全名称
            },
        },
    }
    //console.log(editorConfig.value)
})

const onDocumentReady = () => {
    console.log("Document is loaded");
}

const onLoadComponentError = (errorCode: any, errorDescription: any) => {
    switch (errorCode) {
        case -1: // Unknown error loading component
            console.log("Unknown error loading component", errorDescription);
            break;

        case -2: // Error load DocsAPI from http://documentserver/
            console.log("Error load DocsAPI from http://documentserver/", errorDescription);
            break;

        case -3: // DocsAPI is not defined
            console.log("DocsAPI is not defined", errorDescription);
            break;
    }
}
</script>