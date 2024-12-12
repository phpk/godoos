<template>
    <DocumentEditor id="docEditor" :documentServerUrl="config.onlyoffice.url" :config="editorConfig"
        :events_onDocumentReady="onDocumentReady" :onLoadComponentError="onLoadComponentError" />
</template>

<script lang="ts" setup>
import { DocumentEditor } from "@onlyoffice/document-editor-vue";
import { getSystemConfig, getSplit } from "@/system/config";
import { BrowserWindow, Dialog, Notify, System } from "@/system";
import { ref, onMounted, inject } from "vue";

const config = getSystemConfig();
const sys: any = inject<System>("system");
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
        default: "md",
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
    const fileKey = "docx" + Math.random()
    editorConfig.value = {
        document: {
            fileType: props.ext, // 根据文件扩展名动态设置 fileType
            key:fileKey,
            title: title,
            //url: blobUrl
            url: config.storenet.url + "/file/readfile?stream="+fileKey+"&path=" + path,
            "info": {
                "owner": "王重阳", 
            },
            token: config.onlyoffice.sceret,
            //文档权限参数
            "permissions": {
                "edit": true, //（文件是否可以编辑，false时文件不可编辑）
                "fillForms": true, //定义是否能在文档中填充表单
                "print": true, //定义文档是否能打印
                "review": false, //第一是否显示审阅文档菜单
                "comment": true, //定义是否可以注释文档。如果注释权限设置为“ true”，则文档侧栏将包含“注释”菜单选项；只有将mode参数设置为edit时才生效，默认值与edit参数的值一致。
                "copy": true, //是否允许您将内容复制到剪贴板。默认值为true。
                "download": true, //定义是否可以下载文档或仅在线查看或编辑文档。如果下载权限设置为“false”下载为菜单选项将没有。默认值为true。
                "modifyContentControl": true, //定义是否可以更改内容控件设置。仅当mode参数设置为edit时，内容控件修改才可用于文档编辑器。默认值为true。
                "modifyFilter": true, //定义过滤器是否可以全局应用（true）影响所有其他用户，或局部应用（false），即仅适用于当前用户。如果将mode参数设置为edit，则过滤器修改仅对电子表格编辑器可用。默认值为true。
            }
        },
        documentType: props.onlyType,
        editorConfig: {
            lang: "zh",
            callbackUrl: config.storenet.url + "/file/onlyofficecallback",
            "canCoAuthoring": true,       
            "canHistoryClose": true,
            "canHistoryRestore": true,
            "canMakeActionLink": true,
            "canRename": true,
            "canRequestClose": true,
            "canRequestCompareFile": true,
            "canRequestCreateNew": true,
            "canRequestEditRights": true,
            "canRequestInsertImage": true,
            "canRequestMailMergeRecipients": true,
            "canRequestOpen": true,
            "canRequestReferenceData": true,
            "canRequestReferenceSource": true,
            "canRequestSaveAs": true,
            "canRequestSelectDocument": true,
            "canRequestSelectSpreadsheet": true,
            "canRequestSendNotify": true,
            "canRequestSharingSettings": true,
            "canRequestUsers": true,
            "canSaveDocumentToBinary": true,
            "canSendEmailAddresses": true,
            "canStartFilling": true,
            "canUseHistory": true,
            "user": { //用户信息
                "id": "godoos", //用户ID
                "name": "写作员" //用户全名称
            },
        },
        "events": {

        }
    }
    console.log(editorConfig.value)
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