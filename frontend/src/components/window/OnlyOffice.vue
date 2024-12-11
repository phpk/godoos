<template>
    <DocumentEditor id="docEditor" :documentServerUrl="config.onlyoffice.url" :config="editorConfig"
        :events_onDocumentReady="onDocumentReady" :onLoadComponentError="onLoadComponentError" />
</template>

<script lang="ts" setup>
import { DocumentEditor } from "@onlyoffice/document-editor-vue";
import { getSystemConfig } from "@/system/config";
import { BrowserWindow, Dialog, Notify, System } from "@/system";
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
const editorConfig: any = ref({})
onMounted(() => {
    const path = win?.config?.path;
    //const uniqueKey = fetchDocumentKey(path)
    const readUrl = config.apiUrl + "/file/readfile?stream=1&path=" + path
    editorConfig.value = {
        document: {
            fileType: "docx",
            key: "docx" + Math.random(),
            title: "Example Document Title.docx",
            url: readUrl
        },
        documentType: "word",
        editorConfig: {
            //callbackUrl: "https://example.com/url-to-callback.ashx",
            customization: {
                "anonymous": {
                    request: true,
                    label: "Guest",
                }
            },
        }
    }
})
const onDocumentReady = () => {
    console.log("Document is loaded");
}
const onLoadComponentError = (errorCode: any, errorDescription: any) => {
    switch (errorCode) {
        case -1: // Unknown error loading component
            console.log(errorDescription);
            break;

        case -2: // Error load DocsAPI from http://documentserver/
            console.log(errorDescription);
            break;

        case -3: // DocsAPI is not defined
            console.log(errorDescription);
            break;
    }
}

</script>