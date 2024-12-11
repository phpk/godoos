<template>
    <DocumentEditor id="docEditor" :documentServerUrl="config.onlyoffice.url" :config="editorConfig"
        :events_onDocumentReady="onDocumentReady" :onLoadComponentError="onLoadComponentError" />
</template>

<script lang="ts" setup>
import { DocumentEditor } from "@onlyoffice/document-editor-vue";
import { getSystemConfig, getSplit } from "@/system/config";
import { BrowserWindow, Dialog, Notify, System } from "@/system";
import { ref, onMounted, inject } from "vue";
import { url } from "inspector";

const config = getSystemConfig();
const sys: any = inject<System>("system");
const win: any = inject<BrowserWindow>("browserWindow");

const props = defineProps({
    src: {
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

const editorConfig:any = ref({});
const SP = getSplit();

const getBlobUrl = (content: any) => {
    let blobUrl: string;

    if (!content || content === "") {
        const emptyBlob = new Blob([], { type: 'application/octet-stream' });
        blobUrl = URL.createObjectURL(emptyBlob);
    } else {
        try {
            // 将 base64 转换为 Blob URL
            const byteCharacters = atob(content);
            const byteNumbers = new Array(byteCharacters.length);
            for (let i = 0; i < byteCharacters.length; i++) {
                byteNumbers[i] = byteCharacters.charCodeAt(i);
            }
            const byteArray = new Uint8Array(byteNumbers);
            const blob = new Blob([byteArray], { type: 'application/octet-stream' });
            blobUrl = URL.createObjectURL(blob);
        } catch (error) {
            console.error("Failed to decode base64 content:", error);
            const emptyBlob = new Blob([], { type: 'application/octet-stream' });
            blobUrl = URL.createObjectURL(emptyBlob);
        }
    }
    return blobUrl;
}

onMounted(() => {
    const path = win?.config?.path;
    const title = path.split(SP).pop();
    let content = win?.config?.content;
    // console.log("Path:", path);
    // console.log("Title:", title);
    // console.log("Content:", content);
    // console.log("ext:",props.ext)
    const blobUrl = getBlobUrl(content);
    //console.log("Blob URL:", blobUrl);
    //console.log("onlyoffice:", config.apiUrl + "/file/onlyoffice");
    editorConfig.value = {
        document: {
            fileType: props.ext, // 根据文件扩展名动态设置 fileType
            key: "docx" + Math.random(),
            
            title: title,
            url: blobUrl
            //url: config.storenet.url + "/file/readfile?stream=1&path="+path,
        },
        documentType: "word",
        editorConfig: {
            lang: "zh",
            //createUrl:config.apiUrl + "file/writefile?path="+path,
            callbackUrl: config.storenet.url + "/file/onlyoffice",
        },
        events: {
            onError: (e:any) =>{
                console.log("onError:", e);
            },
            onDocumentStateChange: (e:any) =>{
                console.log("onDocumentStateChange:", e);
            },
            onDocumentReady: (e:any) =>{
                console.log("onDocumentReady:", e);
            },
            onRequestViewRights: (e:any) =>{
                console.log("onRequestViewRights:", e);
            },
            onSave: (e:any)=>{
                console.log("onSave:", e);
            }
        }
    }
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