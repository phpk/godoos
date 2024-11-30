import { defineStore } from "pinia"
import { BrowserWindow } from "@/system";
import { ref } from 'vue';
type ExtMap = {
  [key: string]: string | string[]
}
export const useChooseStore = defineStore('chooseStore', () => {
    const win:any = ref()
    const path:any = ref([])
    const ifShow = ref(false)
    const extArr: ExtMap = {
      'Documents/Word': 'docx',
      'Documents/PPT': 'pptx',
      'Documents/Markdown': 'md',
      'Documents/Excel': 'xlsx',
      'Documents/Mind': 'mind',
      'Documents/Kanban': 'kb',
      'Documents/Baiban': 'bb',
      'Documents/Screenshot': 'screenshot',
      'Documents/ScreenRecoding': 'screentRecording',
      'Pictures': ['png','jpg','webp','gif','bmp','tiff'],
      'Music': 'mp3',
      'Videos': 'mp4'
    }
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
    const saveFile = (title: string, fileExt: any, componentID: string, eventData: any, ext: string) => {
      // console.log('保存文件');
      // 判断是否已经存在
      if (isExist(componentID)) return
      // 打开弹窗
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
      // 默认路径
      let defaultPath = ''
      Object.keys(extArr).forEach( (item: string) => {
        if (extArr[item] && 
          ((Array.isArray(extArr[item]) && extArr[item].includes(ext)) || extArr[item] == ext)
        ) {
          defaultPath = item
        }
      })
      defaultPath == '' ? defaultPath = '/D' : defaultPath = '/C/Users/' + defaultPath
      // console.log('1111111:', defaultPath);
      saveFileContent.value.push({
        componentID,
        eventData,
        defaultPath,
        filePath: '',
        fileName: ''
      })
      // ifShow.value = true
    }
    const close = () => {
      ifShow.value = false
      win.value.close()
    }
    //  清除缓存
    const closeSaveFile = (componentID: string) => {
      saveFileContent.value.forEach((item, index) => {
        if (item.componentID == componentID) {
          saveFileContent.value.splice(index, 1)
        }
      })
      win.value.close()
    }
    // 判断弹窗是否存在
    const isExist = (componentID: string) => {
      saveFileContent.value.forEach(item => {
        if (item.componentID == componentID) {
          return true
        }
      })
      return false
    }
    return {
        win,
        path,
        ifShow,
        saveFileContent,
        select,
        close,
        saveFile,
        closeSaveFile,
        isExist
    }

})