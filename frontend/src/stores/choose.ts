import { defineStore } from "pinia"
import { BrowserWindow } from "@/system";
import { ref } from 'vue';
export const useChooseStore = defineStore('chooseStore', () => {
    const win:any = ref()
    const path = ref([])
    const select = (title = '选择文件', fileExt:any) => {
       win.value = new BrowserWindow({
            title,
            content: "FileList",
            config: {
                ext: fileExt,
                path: '/'
            },
            icon: "gallery",
            width: 700,
            height: 500,
            x: 100,
            y: 100,
            center: true,
            minimizable: false,
            resizable: true,
        });
        win.value.show()
    }
    const close = () => {
        win.value.close()
    }
    
    return {
        win,
        path,
        select,
        close
    }

})