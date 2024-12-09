<template>
    <DocumentEditor id="docEditor" :documentServerUrl="config.onlyoffice.url" :config="editorConfig"
        :events_onDocumentReady="onDocumentReady" :onLoadComponentError="onLoadComponentError" />
</template>

<script lang="ts" setup>
import { DocumentEditor } from "@onlyoffice/document-editor-vue";
import { getSystemConfig } from "@/system/config";
import { BrowserWindow, Dialog, Notify, System } from "@/system";
import { generateRandomString } from "@/util/common";
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
// async function fetchDocumentKey(path: string): Promise<string> {
//     try {
//         const response:any = await fetch(`${config.onlyoffice.url}/get-document-key?path=${path}`);
//         console.log(response)
//         const res  = await response.json();
//         return res.data.key;
//     } catch (error) {
//         console.error("Failed to fetch document key:", error);
//         throw error;
//     }
// }
onMounted(() => {
    const path = win?.config?.path;
    const uniqueKey = generateRandomString(12);
    //const uniqueKey = fetchDocumentKey(path)
    const readUrl = config.apiUrl + "/file/readfile?stream=1&path=" + path
    editorConfig.value = {
        document: {
            fileType: "docx",
            //key: "ojR1OasBPnlIwF9WA80AW4NTrIWqs9",
            //"key": uniqueKey,
            key: "docx" + Math.random(),
            // "permissions": {
            //     "chat": true,
            //     "comment": true,
            //     "copy": true,
            //     "download": true,
            //     "edit": true,
            //     "fillForms": true,
            //     "modifyContentControl": true,
            //     "modifyFilter": true,
            //     "print": true,
            //     "review": true,
            //     "reviewGroups": null,
            //     "commentGroups": {},
            //     "userInfoGroups": null,
            //     "protect": true
            // },
            title: "Example Document Title.docx",
            url: readUrl
        },
        documentType: "word",
        editorConfig: {
            callbackUrl: "https://example.com/url-to-callback.ashx",
            // customization: {
            //     "anonymous": {
            //         request: true,
            //         label: "Guest",
            //     }
            // },
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