import {Files} from '@/api/files';
import { joinknowledge } from '@/api/net/knowledge';
import { shareCreate } from "@/api/net/share";
import { eventBus } from '@/interfaces/event';
import { getFileType } from '@/router/filemaplist';
import { errMsg, noticeMsg, promptMsg, promptPwd, successMsg } from '@/utils/msg';
import { defineStore } from 'pinia';
import { ref } from 'vue';
import { useRouter } from 'vue-router';
export const useFileSystemStore = defineStore('filesystem', () => {
  const currentPath: any = ref('')
  const currentShareFile: any = ref({});
  const router = useRouter()
  const initChoose = {
    ifShow: false,
    isSave: false,
    paths: [] as string[],
    exts: [] as string[],
    ext: '',
    defName: '',
    content: ''
  }
  const choose = ref({ ...initChoose })
  const pwdPathMap: any = ref({})

  // 加入知识库
  const joinKnowledge = async (path: string) => {
    const res = await joinknowledge(path)
    if (res.success) {
      successMsg(res.message);
    } else {
      errMsg(res.message);
    }
    refreshPaths()
  }

  // 读取知识库盘
  

  const setPwdPathMap = (path: string, pwd?: string) => {
    if (pwd && pwd !== '') {
      pwdPathMap.value[path] = pwd;
    }
  };

  const getPwdPathMap = (path: string) => {
    const pwd = pwdPathMap.value[path] || '';
    return pwd;
  };

  const getFilesInPath = async (path: string, pwd?: string) => {
    currentPath.value = path;
    setPwdPathMap(path, pwd);
    const fs = await Files();
    const res = await fs.read(path, getPwdPathMap(path));
    if (!res.success) {
      pwd = await promptPwd();
      if (!pwd) return;
      setPwdPathMap(path, pwd);
      const res = await fs.read(path, pwd);
      //console.log(res)
      if (!res.success) {
        noticeMsg('密码错误', '提示', 'error');
        return [];
      } else {
        return res.data;
      }
    }
    return res.data;
  };
  const initFolder = (path: string) => {
    currentPath.value = path;
  }
  const getUniqueFilePath = async (filePath: string, defname: string, ext: string, path: string, sp: string, i: number): Promise<string> => {
    const fs = await Files();
    const isExits = await fs.exists(filePath);
    //console.log(isExits)
    if (isExits) {
      return getUniqueFilePath(`${path}${sp}${defname}(${i})${ext}`, defname, ext, path, sp, i + 1);
    }
    return filePath;
  };
  const getUniquePath = async (path: string): Promise<string> => {
    const sp = path.charAt(0);
    const fs = await Files();
    const baseName = fs.basename(path);
    const arr = baseName.split('.');
    //console.log(arr)
    const defname = arr[0];
    const ext = arr[1] ? `.${arr[1]}` : '';
    const dir = arr[1] ? fs.dirname(path) : path;
    let newFilePath = arr[1] ? `${dir}${sp}${defname}${ext}` : `${path}`;
    //console.log(newFilePath)
    newFilePath = await getUniqueFilePath(newFilePath, defname, ext, dir, sp, 1);
    return newFilePath;
  };
  const handleNewFile = async (defname: string, ext: string) => {
    const path = currentPath.value;
    //console.log(path)
    const sp = path.charAt(0);
    let newFilePath = `${path}${sp}${defname}${ext}`;
    newFilePath = await getUniqueFilePath(newFilePath, defname, ext, path, sp, 1);
    //console.log(newFilePath);
    const fs = await Files();
    await fs.writeFile(newFilePath, '');
    refreshPaths()
  }
  const handleReadFile = async (filePath: string, pwd?: string) => {
    setPwdPathMap(filePath, pwd);
    const fs = await Files();
    const res = await fs.readFile(filePath, getPwdPathMap(filePath));
    console.log(res);
    if (!res.success) {
      pwd = await promptPwd();
      if (!pwd) return;
      setPwdPathMap(filePath, pwd);
      const res = await fs.readFile(filePath, pwd);
      if (!res.success) {
        noticeMsg('密码错误', '提示', 'error');
        return;
      }
      else {
        return res.data;
      }
    }
    return res.data;
  }
  const handlePwdFile = async (path: string) => {
    const pwd = await promptPwd();
    if (!pwd) return;
    const fs = await Files();
    const res = await fs.pwd(path, pwd);
    refreshPaths()
    return res;
  }
  const handleUnpwdFile = async (filePath: string) => {
    const pwd = await promptPwd();
    if (!pwd) return;
    const fs = await Files();
    const res = await fs.unpwd(filePath, pwd);
    refreshPaths()
    return res;
  }
  const handleWriteFile = async (filePath: string, data: any, pwd?: string) => {
    //console.log('data:', data)
    setPwdPathMap(filePath, pwd);
    const fs = await Files();
    let res = await fs.writeFile(filePath, data, getPwdPathMap(filePath));
    if (!res) {
      const pwd = await promptPwd();
      if (!pwd) return;
      setPwdPathMap(filePath, pwd);
      res = await fs.writeFile(filePath, data, pwd);
    }
    refreshPaths()
    return res
  }
  const refreshPaths = () => {
    eventBus.emit('refreshDesktop');
  }

  const handleSerach = async (path: string, query: string) => {
    const fs = await Files();
    const res = await fs.search(path, query);
    return res || [];
  }
  const handleDeleteFile = async (path: string) => {
    if (!path) return;
    const fs = await Files();
    await fs.rmdir(path)
    //console.log(dirname(path))
    refreshPaths()
  }
  const handleDeleteFiles = async (filePaths: string[]) => {
    const fs = await Files();
    for (const filePath of filePaths) {
      await fs.rmdir(filePath)
    }
    refreshPaths()
  }
  const handleNewDir = async (dirname: string) => {
    const path = currentPath.value;
    //console.log(path)
    const sp = path.charAt(0);
    let newDirPath = `${path}${sp}${dirname}`;
    const fs = await Files();
    newDirPath = await getUniqueFilePath(`${newDirPath}`, dirname, '', path, sp, 1);
    await fs.mkdir(newDirPath);
    refreshPaths()
  }
  const handleRenameFile = async (oldPath: string) => {
    if (!oldPath) return;
    //const sp = oldPath.charAt(0);
    const fs = await Files();
    const oldName = fs.basename(oldPath);
    const newName = await promptMsg('请输入新文件名', '重命名', oldName);
    //console.log(newName)
    if (!newName || newName == '') return;
    const parentpath = fs.getParentPath(oldPath);
    const newFilePath = fs.join(parentpath, newName);
    await fs.rename(oldPath, newFilePath);
    refreshPaths()
  }
  const handleShareFile = async (data: any) => {
    const res = await shareCreate(JSON.stringify(data));
    refreshPaths()
    return res;
  }

  const openfile = async (file: any) => {
    if (file.isPwd) {
      const pwd = await promptPwd();
      if (pwd) {
        setPwdPathMap(file.path, pwd);
      } else {
        errMsg('请输入密码');
        return;
      }
    }
    if (file.isDirectory) {
      router.push({ path: '/computer', query: { path: file.path } });
    } else {
      if (file.ext == 'exe') {
        router.push('/' + file.title);
      } else {
        //file.content = ''
        //console.log(file)
        const fileMap = getFileType(file.ext)
        if (fileMap) {
          file.editor = fileMap.editor;
          file.hasPrview = fileMap.hasPrview;
          file.eventType = fileMap.eventType;
          file.exts = fileMap.ext;
        }
        router.push({ path: '/viewer', query: file });
      }
    }
  }
  const openFile = async (file: any) => {
    const fs = await Files();
    if (typeof file == 'string') {
      const f = await fs.stat(file)
      await openfile(f)
    } else {
      await openfile(file)
    }
  }

  const openEditor = (file: any) => {
    file.content = ''
    file.action = 'edit';
    // router.push({ path: '/viewer', query: file });
    openFile(file)
  }
  const moveFiles = async (filePaths: string[], targetPath: string) => {
    const fs = await Files();
    const toFile = await fs.stat(targetPath);
    if (toFile?.isDirectory) {
      for (const filePath of filePaths) {
        //console.log(filePath)
        let newPath = fs.join(targetPath, fs.basename(filePath));
        //console.log(newPath)
        newPath = await getUniquePath(newPath);
        //console.log(newPath)
        await fs.rename(filePath, newPath);
      }
      refreshPaths()
    }
  }
  const copyFiles = async (filePaths: string[], targetPath: string) => {
    const fs = await Files();
    const toFile = await fs.stat(targetPath);
    if (toFile?.isDirectory) {
      for (const filePath of filePaths) {
        let newPath = fs.join(targetPath, fs.basename(filePath));
        newPath = await getUniquePath(newPath);
        await fs.copy(filePath, newPath);
      }
      refreshPaths()
    }
  }
  const handleUnzipFile = async (filePath: string) => {
    const fs = await Files();
    await fs.unzip(filePath);
    refreshPaths()
  }
  const handleFavorite = async (filePath: string) => {
    const fs = await Files();
    await fs.favorite(filePath);
    refreshPaths()
  }
  const handleZipFile = async (filePath: string, ext: string) => {
    const fs = await Files();
    await fs.zip(filePath, ext);
    refreshPaths()
    noticeMsg('压缩成功');
  }
  const handleFileStat = async (filePath: string) => {
    const fs = await Files();
    return await fs.stat(filePath);
  }
  const chooseFiles = (exts: string[] = []) => {
    choose.value.ifShow = true;
    choose.value.isSave = false;
    choose.value.exts = exts;
    router.push({ path: '/computer' });
  }
  const saveFiles = async (defname: string, ext: string, content: any) => {
    choose.value.ifShow = true;
    choose.value.isSave = true;
    choose.value.ext = ext;
    if (!defname || defname == '') {
      defname = await promptMsg('请输入文件名', '文件名', '');
    }
    //defname = defname ?? prompt('请输入文件名') ?? '';

    if (!defname) {
      return;
    }

    if (!defname.endsWith(ext)) {
      defname = `${defname}.${ext}`;
    }

    choose.value.defName = defname;
    choose.value.content = content;
    const extArr = {
      '/C/Users/Documents/Word': ['docx'],
      '/C/Users/Documents/PPT': ['pptx'],
      '/C/Users/Documents/Markdown': ['md'],
      '/C/Users/Documents/Execl': ['xlsx', 'xls'],
      '/C/Users/Documents/Mind': ['mind'],
      '/C/Users/Documents/Kanban': ['kb'],
      '/C/Users/Documents/Baiban': ['bb'],
      '/C/Users/Documents/Screenshot': ['screenshot'],
      '/C/Users/Documents/ScreenRecoding': ['screentRecording'],
      '/C/Users/Pictures': ['png', 'jpg', 'webp', 'gif', 'bmp', 'tiff'],
      '/C/Users/Music': ['mp3'],
      '/C/Users/Videos': ['mp4']
    }
    let savePath = '/D'
    if (ext !== '') {
      for (const [key, value] of Object.entries(extArr)) {
        if (value.includes(ext)) {
          savePath = key;
          break;
        }
      }
    }
    router.push({ path: '/computer', query: { path: savePath } });
  }
  const clearChoose = () => {
    choose.value = { ...initChoose };
  }
  const parserFormData = async (content: any, contentType: any) => {
    const fs = await Files();
    return fs.parserFormData(content, contentType);
  }
  const reStoreFiles = async (dirPath: string) => {
    const fs = await Files();
    await fs.restore(dirPath);
    refreshPaths()
  }
  function uploadFile(accept: string = '*/*') {
    const fileInput = ref<HTMLInputElement | null>(null);

    // 创建隐藏的文件输入元素
    const createFileInput = () => {
      const input = document.createElement('input');
      input.type = 'file';
      input.accept = accept;
      input.multiple = true;
      input.style.display = 'none';
      input.addEventListener('change', handleFileChange);
      document.body.appendChild(input);
      fileInput.value = input;
    };

    // 触发文件选择对话框
    const triggerFileInput = () => {
      if (fileInput.value) {
        fileInput.value.click();
      }
    };

    // 处理文件选择事件
    const handleFileChange = (event: Event) => {
      const target = event.target as HTMLInputElement;
      if (target.files && target.files.length > 0) {
        //const file = target.files[0];
        onFileSelected(target.files);
      }
      // 移除文件输入元素
      if (fileInput.value) {
        fileInput.value.remove();
        fileInput.value = null;
      }
    };

    // 创建并触发文件输入元素
    createFileInput();
    triggerFileInput();
  }
  async function onFileSelected(list: any) {
    if (!list) return;
    const fs = await Files();
    const len = list.length
    const path = fs.dirname(currentPath.value);
    for (let i = 0; i < len; i++) {
      //setProgress((i / len) * 100);
      await new Promise((resolve) => {
        const item = list?.[i];
        //console.log(item)
        if (!item) return;
        const reader = new FileReader();
        //读取成功
        reader.onload = function () {
        };
        reader.onloadstart = function () {
          //console.log('读取开始');
        };
        reader.onloadend = async function () {
          await handleWriteFile(fs.join(path, item.name), reader.result as any)
          resolve(true);
        };
        reader.onabort = function () {
          //console.log('中断');
        };
        reader.onerror = function () {
          //console.log('读取失败');
        };
        reader.onprogress = function () {
          //const scale = ev.loaded / ev.total;
          //setProgress(scale);
        };
        reader.readAsArrayBuffer(item);
      });
    }
    //setProgress(100);
  }
  async function downloadFile(filePath: string) {
    try {
      const fileContent = await handleReadFile(filePath);
      if (!fileContent) {
        noticeMsg('文件读取失败', '提示', 'error');
        return;
      }

      const blob = new Blob([fileContent], { type: 'application/octet-stream' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      const fs = await Files();
      a.download = fs.basename(filePath);
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
      //noticeMsg('文件下载成功', '提示', 'success');
    } catch (error) {
      noticeMsg('文件下载失败', '提示', 'error');
      console.error(error);
    }
  }
  return {
    // editFilePath,
    currentPath,
    choose,
    currentShareFile,
    fs:Files,
    initFolder,
    getFilesInPath,
    getUniquePath,
    handleNewFile,
    handleDeleteFiles,
    handleDeleteFile,
    handleWriteFile,
    handlePwdFile,
    handleUnpwdFile,
    handleRenameFile,
    handleShareFile,
    handleUnzipFile,
    handleZipFile,
    handleFavorite,
    handleFileStat,
    handleSerach,
    refreshPaths,
    reStoreFiles,
    handleNewDir,
    openFile,
    handleReadFile,
    openEditor,
    moveFiles,
    copyFiles,
    chooseFiles,
    clearChoose,
    saveFiles,
    parserFormData,
    uploadFile,
    downloadFile,
    joinKnowledge,
    
  };
});