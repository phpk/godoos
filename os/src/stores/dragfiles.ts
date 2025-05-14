import { defineStore } from 'pinia';
import { join } from '@/api/files';
import { useClickingStore } from './clicking';
import { useFileSystemStore } from './filesystem';
import { ref } from 'vue';
export const useDragFilesStore = defineStore('dragfiles', () => {
  const clickingStore = useClickingStore();
  const fileStore = useFileSystemStore();
  const targetPath = ref('');
  const isDragFile = ref<boolean>(false); // 用于标记是否正在拖动文件
  function startDrag(ev: DragEvent) {
    //console.log(ev)
    isDragFile.value = true;
    //console.log(clickingStore.clickedIcons)
    if (clickingStore.clickedIcons.length < 1) return;
    ev?.dataTransfer?.setData('fromobj', 'os');
    ev?.dataTransfer?.setData('frompath', JSON.stringify(clickingStore.clickedIcons));
    //clickingStore.isFileDragging = true; // 设置文件拖放标志
  }
  async function folderDrop(ev: DragEvent, toPath: string) {
    const frompathArrStr = ev?.dataTransfer?.getData('frompath');
    //console.log(toPath)
    if (!frompathArrStr) return;
    const frompathArr = JSON.parse(frompathArrStr) as string[];
    if (frompathArr.length === 0) {
      return;
    }
    if (frompathArr.includes(toPath)) {
      return;
    }

    await fileStore.moveFiles(frompathArr, toPath);
    clickingStore.clickedIcons = [];
    
    //ev.stopPropagation();
    //clickingStore.isFileDragging = false;
  }
  async function outerFileDrop(path: string, list: any) {
    if (!list) return;
    const len = list.length
    // const { setProgress } = Dialog.showProcessDialog({
    //   message: '正在写入到文件夹中',
    // });
    for (let i = 0; i < len; i++) {
      //setProgress((i / len) * 100);
      await new Promise((resolve) => {
        const item = list?.[i];
        console.log(item)
        if (!item) return;
        // let oFile = null;
        const reader = new FileReader();
        //读取成功
        reader.onload = function () {
          // console.log(reader);
          //console.log(reader.result)
        };
        reader.onloadstart = function () {
          //console.log('读取开始');
        };
        reader.onloadend = async function () {
          //await writeFileToInner(path, item.name, reader.result as any, process);
          //await writeFile(join(path, item.name), reader.result as any)
          await fileStore.handleWriteFile(join(path, item.name), reader.result as any)
          resolve(true);
        };
        reader.onabort = function () {
          //console.log('中断');
        };
        reader.onerror = function () {
          //console.log('读取失败');
        };
        reader.onprogress = function (ev) {
          //const scale = ev.loaded / ev.total;
          //setProgress(scale);
        };
        reader.readAsArrayBuffer(item);
      });
    }
    //setProgress(101);
  }
  async function dragFileToDrop(ev: DragEvent, path: string) {
    ev.preventDefault();
    ev.stopPropagation(); // 阻止事件冒泡
    //console.log(path)
    const fromobj = ev?.dataTransfer?.getData('fromobj');
    console.log(fromobj)
    if (fromobj == 'os') {
      folderDrop(ev, path);
    } else {
      const oFileList: any = ev?.dataTransfer?.files;
      //console.log(oFileList);
      if (oFileList?.length > 0) {
        await outerFileDrop(path, oFileList)
      }
    }
    //ev.stopPropagation(); 
  }
  function handleDragEnter(path: string) {
    targetPath.value = path;
  }
  function handleDragOver(event: DragEvent) {
    event.preventDefault();
    event.stopPropagation();
}
  
  function handleDragLeave() {
    targetPath.value = ''
  }
  return { 
    startDrag, 
    isDragFile,
    dragFileToDrop,
    handleDragEnter,
    handleDragOver,
    handleDragLeave,
  };
})