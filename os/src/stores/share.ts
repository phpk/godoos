import { defineStore } from "pinia";
import { ref } from "vue";
export const useShareStore = defineStore('shareStroe', () => {
    const fileList:any = ref([]);
    const currentFile:any = ref(null);
    const addFile = (file:any) => {
        if (fileList.value.indexOf(file) !== -1) return;
        fileList.value.push(file);
        currentFile.value = file;
    }
    const removeFile = (file:any) => {
        fileList.value.splice(fileList.value.indexOf(file), 1);
    }
    return { 
        fileList,
        currentFile,
        addFile,
        removeFile
     };
})