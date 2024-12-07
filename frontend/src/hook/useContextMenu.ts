import * as fspath from '../system/core/Path';
import { useSystem } from '../system';
import { BrowserWindow } from '../system/window/BrowserWindow';
import { OsFileWithoutContent } from '../system/core/FileSystem';
import { t } from '../i18n';
import { Dialog } from '../system/window/Dialog';
import { Menu } from '../system/menu/Menu';
import { uniqBy } from '../util/modash';
import { UnwrapNestedRefs } from 'vue';
// import { getSystemConfig } from "@/system/config";

export function createTaskbarIconContextMenu(e: MouseEvent, windowNode: UnwrapNestedRefs<BrowserWindow>) {
  Menu.buildFromTemplate([
    {
      label: t('close'),
      click: () => {
        windowNode.close();
      },
    },
    {
      label: t('maximize'),
      click: () => {
        windowNode.maximize();
      },
    },
    {
      label: t('minimize'),
      click: () => {
        windowNode.minimize();
      },
    },
  ]).popup(e);
}
function useContextMenu() {
  function createDesktopContextMenu(
    e: MouseEvent,
    path = `${useSystem()._options.userLocation}Desktop`,
    callback?: () => void
  ) {
    const system = useSystem();
    if (!system) return;

    const clickFunc = (ext: string, name: string) => {
      createNewFile(path, '.' + ext, name).then(() => {
        callback?.();
      });
    }

    //console.log(path)
    let menuArr: any = [
      {
        label: t('refresh'),
        click: () => {
          callback?.();
        },
      },
      {
        label: '新建文件',
        submenu: [

          {
            label: 'MarkDown文件',
            click: () => clickFunc('md', '未命名MD文件'),
          },
          {
            label: 'word文档(Docx)',
            click: () => clickFunc('docx', '未命名文档'),
          },
          {
            label: '演示文稿(PPT)',
            click: () => clickFunc('pptx', '未命名演示文稿'),
          },
          {
            label: '数据表格(Excel)',
            click: () => clickFunc('xlsx', '未命名数据表格'),
          },
          {
            label: '文本文件(txt)',
            click: () => clickFunc('txt', '未命名文本文件'),
          },
        ],
      },
      {
        label: '新建工作图',
        submenu: [
          {
            label: '思维导图(Mind)',
            click: () => clickFunc('mind', '根节点'),
          },
          {
            label: '看板(Kanban)',
            click: () => clickFunc('kb', '未命名看板'),
          },
          {
            label: '白板(Baiban)',
            click: () => clickFunc('bb', '未命名白板'),
          },
          {
            label: '甘特图(Gant)',
            click: () => clickFunc('gant', '未命名项目'),
          },
        ],
      },

      {
        label: t('paste'),
        click: () => {
          pasteFile(path).then(() => {
            callback?.();
          });
        },
      },
      {
        label: t('new.folder'),
        click: () => {
          createNewDir(path).then(() => {
            callback?.();
          });
        },
      },
    ];
    // const userInfo: any = system.getConfig('userInfo');
    // if (userInfo.user_auths && userInfo.user_auths.length > 0 && userInfo.user_auths != "") {
    //   menuArr.push(
    //     {
    //       label: '安排任务',
    //       click: () => {
    //         const win = new BrowserWindow({
    //           title: '安排任务',
    //           content: "PlanTasks",
    //           width: 600,
    //           height: 600,
    //           center: true,
    //         });
    //         win.show();
    //       }
    //     }
    //   )
    // }
    menuArr = [...menuArr, ...(system._rootState.options.contextMenus || [])]
    const menu = Menu.buildFromTemplate(
      uniqBy(
        menuArr,
        (val) => val.label
      )
    );
    menu.popup(e);
  }
  async function createNewFile(path: string, ext: string = '.txt', title: string = t('new.file')) {
    const system = useSystem();
    if (!system) return;
    if (["/", "/B"].includes(path)) return;
    let newFilePath = fspath.join(path, title + ext);
    if (await system.fs.exists(newFilePath)) {
      let i = 1;
      while (await system.fs.exists(fspath.join(path, `${title}(${i})${ext}`))) {
        i++;
      }
      newFilePath = fspath.join(path, `${title}(${i})${ext}`);
    }

    const content = '';
    // initPwd(newFilePath)
    return await system.fs.writeFile(newFilePath, content);
  }
  // 开源版新建文件初始化密码
  // function initPwd(path: string) {
  //   //console.log('新建路径：',path);
    
  //   if (getSystemConfig().userType == 'person' && getSystemConfig().file.isPwd == 1) {
  //     const win = new BrowserWindow({
  //       title: "初始化文件密码",
  //       content: "FilePwd",
  //       config: {
  //         path: path,
  //       },
  //       width: 400,
  //       height: 200,
  //       center: true,
  //     });
  //     win.show()
  //   }
  // }
  async function createNewDir(path: string) {
    const system = useSystem();
    if (!system) return;
    if (["/", "/B"].includes(path)) return;
    let newFilePath = fspath.join(path, t('new.folder'));
    if (await system.fs.exists(newFilePath)) {
      let i = 1;
      while (await system.fs.exists(fspath.join(path, `${t('new.folder')}(${i})`))) {
        i++;
      }
      newFilePath = fspath.join(path, `${t('new.folder')}(${i})`);
    }
    return await system.fs.mkdir(newFilePath);
  }
  function openPropsWindow(path: string) {
    new BrowserWindow({
      title: t('props'),
      content: "FileProps",
      config: {
        content: path,
      },
      width: 350,
      height: 400,
      resizable: false,
    }).show();
  }
  async function backFile(file: OsFileWithoutContent) {
    const system = useSystem();
    if (!system) return;
    const filePath = file.path;
    if (filePath === "/") return;

    const vol = filePath.charAt(1);
    if (vol != "B") return;
    let newPath: any = file.oldPath;
    //console.log(newPath)
    if (await system.fs.exists(newPath)) {

      if (file.isDirectory) {
        let i = 1;
        while (await system.fs.exists(`${newPath}(${i})`)) {
          i++;
        }
        newPath = `${newPath}(${i})`;

      } else {
        let i = 1;
        let parentPath = '/B' + file.parentPath.substring(2)
        while (await system.fs.exists(fspath.join(parentPath, `${file.title}(${i}).${file.ext}`))) {
          i++;
        }
        newPath = fspath.join(parentPath, `${file.title}(${i}).${file.ext}`)
      }
    }

    //console.log(newPath)
    return system?.fs.rename(file.path, newPath)


  }
  async function deleteFile(file: OsFileWithoutContent) {
    const system = useSystem();
    if (!system) return;
    if (!file) return;
    const filePath = file.path;
    if (filePath === "/") return;

    if (file.isDirectory) {
      if (["/", "/B", "/C", "/D", "/E", "/F"].includes(file.path)) return;
    }
    //console.log(file)
    const vol = filePath.charAt(1);
    //console.log(vol)
    if (file && file.id && file.id < 1) {
      if (["C", "D", "E"].includes(vol)) {
        //console.log('/B' + filePath.substring(2))
        let newPath = '/B' + filePath
        // if (file.parentPath == "/C/Users/Desktop" && !file.isDirectory) {
        //   newPath = '/B/' + file.title + "." + file.ext
        // }
        const fileDir = '/B' + (file.isDirectory ? file.path : file.path.replace(file.name, ""))
        //console.log(fileDir)
        if (!(await system.fs.exists(fileDir))) {
          await system.fs.mkdir(fileDir)
        }
        if (await system.fs.exists(newPath)) {

          if (file.isDirectory) {
            let i = 1;
            while (await system.fs.exists(`${newPath}(${i})`)) {
              i++;
            }
            newPath = `${newPath}(${i})`;

          } else {
            let i = 1;
            let parentPath = '/B' + file.parentPath
            while (await system.fs.exists(fspath.join(parentPath, `${file.title}(${i}).${file.ext}`))) {
              i++;
            }
            newPath = fspath.join(parentPath, `${file.title}(${i}).${file.ext}`)
          }
        }
        return system?.fs.rename(file.path, newPath)
      } else {
        if (file.isDirectory) {
          return system?.fs.rmdir(file.path);
        } else {
          return system?.fs.unlink(file.path);
        }
      }
    } else {
      if (file.isDirectory) {
        return system?.fs.rmdir(file.path);
      } else {
        return system?.fs.unlink(file.path);
      }
    }



  }
  async function openWith(file: OsFileWithoutContent) {
    const tempWin = new BrowserWindow({
      title: t('open.with'),
      content: "OpenWiteDialog",
      config: {
        content: file.path,
      },
      width: 400,
      height: 400,
      center: true,
    });
    tempWin.on('blur', () => {
      tempWin.close();
    });
    tempWin.show();
  }

  async function copyFile(files: OsFileWithoutContent[]) {
    const system = useSystem();
    const rootState = system._rootState;
    if (!system) return;
    if (rootState.clipboard) {
      rootState.clipboard = files.map((file) => file.path);
    }
  }
  async function pasteFile(path: string) {
    const system = useSystem();
    if (!system) return;
    if (["/", "/B"].includes(path)) return;
    const rootState = system._rootState;
    const clipLen = Object.keys(rootState.clipboard).length;
    if (clipLen) {
      const clipFiles = rootState.clipboard;

      if (!clipFiles.forEach) {
        return;
      }
      await clipFiles.forEach(async (clipFile: string) => {
        let tempName = fspath.filename(clipFile);
        const ext = fspath.extname(clipFile);

        if (await system.fs.exists(fspath.join(path, tempName) + ext)) {
          let i = 1;
          while (await system.fs.exists(fspath.join(path, `${tempName}(${i})`) + ext)) {
            i++;
          }
          tempName = `${tempName}(${i})`;
        }
        return system.fs.copyFile(clipFile, fspath.join(path, tempName) + ext);
      });
    } else {
      system.emitError('no file in clipboard');
    }
  }
  // 创建快捷方式
  async function createLink(path: string) {
    const system = useSystem();
    if (!system) return;
    if (["/", "/B"].includes(path)) return;
    if (fspath.extname(path) === '.ln') {
      Dialog.showMessageBox({
        title: t('error'),
        message: t('cannot.create.shortcut'),
        type: 'error',
      });
      return;
    }
    const parentPath = fspath.dirname(path);
    const baseName = fspath.basename(path);
    const linkPath = fspath.join(parentPath, baseName + '.ln');
    if (await system.fs.exists(linkPath)) {
      Dialog.showMessageBox({
        title: t('error'),
        message: t('shortcut.has.been.created'),
        type: 'error',
      });
      return;
    }
    return await system.fs.writeFile(linkPath, path);
  }
  return {
    createDesktopContextMenu,
    createNewFile,
    openPropsWindow,
    createNewDir,
    deleteFile,
    backFile,
    openWith,
    copyFile,
    pasteFile,
    createLink,
  };
}

export { useContextMenu };
