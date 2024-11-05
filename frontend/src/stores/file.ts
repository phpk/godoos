import { defineStore } from 'pinia'
import { OsFileWithoutContent } from '@/system/core/FileSystem';
export const useShareFile = defineStore('shareFile', () =>{
    const currentShareFile = ref<OsFileWithoutContent[]>([]);
    let CurrentFile = ref<OsFileWithoutContent>();
    const setShareFile = (file: OsFileWithoutContent[]) => {
        currentShareFile.value = file
        //console.log('存储文件：', currentShareFile);
    }
    const setCurrentFile = (file: OsFileWithoutContent) => {
        CurrentFile.value = file
        //console.log('当前文件：', CurrentFile.value);
    }
    const getShareFile = (): OsFileWithoutContent[] | null =>{
        return currentShareFile.value
    }
    const findServerPath = (titleName: string): OsFileWithoutContent | string => {
        //console.log('查找文件：',getShareFile());
        let result = currentShareFile.value.find(item =>{
            return item.titleName === titleName
        })
        if (!result) {
            //console.log('回退茶砸后：',currentShareFile.value[0]);
            return currentShareFile.value[0].parentPath 
        }
        return result
    }
    return {
        setShareFile,
        setCurrentFile,
        getShareFile,
        findServerPath,
        currentShareFile,
        CurrentFile
    }
})
