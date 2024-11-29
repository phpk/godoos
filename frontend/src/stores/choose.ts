import { defineStore } from "pinia"
import { BrowserWindow } from "@/system";
import { ref } from 'vue';
export const useChooseStore = defineStore('chooseStore', () => {
    const win:any = ref()
    const path:any = ref([])
    const ifShow = ref(false)
    let savePath = ref({
      path: '',
      name: ''
    })
    const saveFileContent = ref<any[]>([])
    const select = (title = '选择文件', fileExt:any) => {
      
       win.value = new BrowserWindow({
            title,
            content: "Computer",
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
        ifShow.value = true
    }
    const saveFile = (title: string, fileExt: any, componentID: string, eventData: any) => {
      // console.log('保存文件');
      win.value = new BrowserWindow({
        title,
        content: "Computer",
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
        footer: true,
        componentID
      });
      win.value.show()
      saveFileContent.value.push({
        componentID,
        eventData,
        filePath: '',
        fileName: ''
      })
      return savePath
      // ifShow.value = true
    }
    const close = () => {
        ifShow.value = false
        win.value.close()
    }
    
    return {
        win,
        path,
        ifShow,
        savePath,
        saveFileContent,
        select,
        close,
        saveFile
    }

})