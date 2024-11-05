import * as fspath from '../system/core/Path';
import { OsFileWithoutContent } from '../system/core/FileSystem';

export type RouterPath = string;
export const useComputer = (adpater: {
  setRouter: (path: RouterPath) => void;
  getRouter: () => RouterPath;
  setFileList: (list: OsFileWithoutContent[]) => void;
  openFile: (path: RouterPath) => void;
  rmdir: (path: RouterPath) => Promise<void>;
  mkdir: (path: RouterPath) => Promise<void>;
  readdir: (path: RouterPath) => Promise<OsFileWithoutContent[]>;
  sharedir: (path: RouterPath) => Promise<OsFileWithoutContent[]>;
  exists: (path: RouterPath) => Promise<boolean>;
  isDirectory: (file: OsFileWithoutContent) => boolean;
  notify: (title: string, content: string) => void;
  search: (keyword: string) => Promise<OsFileWithoutContent[]>;
  readShareDir: (path: RouterPath) => Promise<OsFileWithoutContent[]>;
}) => {
  const isVia = async (path: RouterPath) => {
    if (path === '') path = '/';
    else if (path === '/') path = '/';
    else if (path.endsWith('/')) path = path.substr(0, path.length - 1);

    if(path.substring(0,2) !== '/F') {
      const isExist = await adpater.exists(path);
      if (!isExist) {
        adpater.notify('路径不存在', path);
        return false;
      }
    }
    return true;
  };
  const refersh = async () => {
    const currentPath = adpater.getRouter();
    //console.log('refresh:', currentPath);
    
    if (!currentPath) return;
    if (currentPath.startsWith('search:')) {
      const keyword = currentPath.substr(7);
      const result = await adpater.search(keyword);
      adpater.setFileList(result);
      return;
    }
    if (!(await isVia(currentPath))) return;
    let result
    if (currentPath === '/F/myshare' || currentPath === '/F/othershare') {
      result = await adpater.sharedir(currentPath)
    } else {
      result = await adpater.readdir(currentPath);
    }
    // else if (currentPath.indexOf('/F') === 0) {
    //   //判断是否是回退
    //   // console.log('currentPath:', currentPath);
      
    //   if (offset && offset === -1) {
    //     const parentPath = useShareFile().CurrentFile?.parentPath
    //     //console.log('回退:',parentPath);
    //     result = await adpater.readShareDir('path', parentPath, currentPath)
    //   } else {
    //     result = await adpater.readShareDir('file',file)
    //   }
    //   // 查找文件回退的路径
    //   // const arr = useShareFile().getShareFile()
    //   // file = arr?.find(item => {
    //   //   return item.titleName === currentPath
    //   // })
    //   // console.log('file:',arr,currentPath, file);
    //   // if (!file) {
    //   //   const parentPath = useShareFile().CurrentFile?.parentPath
    //   //   console.log('回退:',parentPath);
    //   //   result = await adpater.readShareDir('path', parentPath, currentPath)
    //   // } else {
    //   //   result = await adpater.readShareDir('file',file)
    //   // }
    // } 
    if (result) adpater.setFileList(result);
  };
  const createFolder = (path: RouterPath) => {
    const currentPath = adpater.getRouter();
    adpater.mkdir(fspath.join(currentPath, path)).then(
      () => {
        refersh();
      },
      () => {
        // Notify
      }
    );
  };
  const backFolder = () => {
    const path = adpater.getRouter();
    if (path === '/') return;
    adpater.setRouter(fspath.join(path, '..'));
    refersh();
  };
  const openFolder = (file: OsFileWithoutContent) => {
    if (adpater.isDirectory(file)) {
      // console.log('文件：',file);
      // const path = file?.isShare ? turnLocalPath(file.titleName, file.root) : file.path
      // console.log('路径：', path);
      adpater.setRouter(file.path);
      refersh();
    } else {
      adpater.openFile(file.path);
    }
  };
  const onComputerMount = async () => {
    await refersh();
  };
  return {
    isVia,
    refersh,
    createFolder,
    backFolder,
    openFolder,
    onComputerMount,
  };
};
