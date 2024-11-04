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
    if (!currentPath) return;
    if (currentPath.startsWith('search:')) {
      const keyword = currentPath.substr(7);
      const result = await adpater.search(keyword);
      adpater.setFileList(result);
      return;
    }
    if (!(await isVia(currentPath))) return;
    // console.log('use computer refresh:', currentPath);
    let result
    if (currentPath === '/F/myshare' || currentPath === '/F/othershare') {
      result = await adpater.sharedir(currentPath)
    } else if (currentPath.indexOf('/F') === 0) {
      // result = await adpater.readShareDir(file?.path || currentPath)
      result = await adpater.readShareDir(currentPath)
    } else {
      result = await adpater.readdir(currentPath);
    }
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
      // const path = file?.isShare ? file.titleName : file.path
      const path = file.path
      adpater.setRouter(path);
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
