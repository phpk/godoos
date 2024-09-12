import type { System } from "./index";
import { BrowserWindow } from "./window/BrowserWindow";
import { t } from "@/i18n";

import { dealIcon } from "../util/Icon";
import { basename } from "./core/Path";
import { Dialog } from "./window/Dialog";
import { Tray } from "./menu/Tary";
import { appList } from "./applist.ts";
import { getSystemConfig } from './config'
import { memberList } from "./member.ts";
const unknownIcon = "unknown";
export function initBuiltinApp(system: System) {
  const config = getSystemConfig();
  let sysList = appList;
  if(config.userType == 'member') {
    sysList = [...appList,...memberList]
  }
  sysList.forEach((d: any) => {

    let addSave = {
      name: d.name,
      icon: d.appIcon,
      window: {
        title: t(d.name),
        content: d.content,
        icon: d.appIcon,
        width: d.width,
        height: d.height,
        frame: d.frame,
        path: d.path ?? undefined,
        url: d.url ?? undefined,
        ext: d.ext ?? undefined,
        eventType: d.eventType ?? undefined,
        center: true,
        config: {
          path: d.dir ?? undefined
        }
      },
    };
    //console.log(d)
    if (d.isDeskTop) {
      system.addApp(addSave);
    }
    // if (d.isMagnet) {
    //   system.addMagnet(addSave);
    // }
    if (d.isMenuList) {
      system.addMenuList(addSave);
    }
    if (d.ext && d.ext.length > 0) {
      const ext = d.ext.map((e: string) => "." + e);
      //console.log(ext)
      system.registerFileOpener(ext, {
        icon: d.appIcon,
        func: (path: string, content: string) => {
          const win = new BrowserWindow({
            title: path,
            icon: d.appIcon,
            width: d.width,
            height: d.height,
            path: addSave.window.path,
            url: addSave.window.url,
            ext: addSave.window.ext,
            eventType: addSave.window.eventType,
            resizable: d.resizable === undefined ? true : d.resizable,
            center: true,
            content: d.content,
            config: {
              path: path,
              content: content,
            },
          });
          //win.maximize();
          //win.isResizable(true);
          win.show();
        },
      });
    }
    if (d.isContext) {
      system.setConfig("contextMenus", [
        {
          label: t(d.name),
          click() {
            new BrowserWindow({
              title: t(d.name),
              icon: d.appIcon,
              path: addSave.window.path,
              width: d.width,
              height: d.height,
              center: true,
              content: d.content,
              url: addSave.window.url,
              ext: addSave.window.ext,
              eventType: addSave.window.eventType,
              resizable: d.resizable,
            }).show();
          },
        },
      ]);
    }
  });

}
export function initBuiltinFileOpener(system: System) {
  if (system._options.builtinFeature?.includes("ExeOpener")) {
    system.registerFileOpener(".exe", {
      name: "可执行程序",
      icon: unknownIcon,
      hiddenInChosen: true,
      func: (path: string, content: string) => {
        console.log(path)
        // console.log(content)
        // console.log(JSON.stringify(system._rootState.windowMap))
        const exeContent = content.split("::");
        const winopt = system._rootState.windowMap[exeContent[1]].get(
          exeContent[2]
        );
        //console.log(winopt)
        if (winopt) {
          if (winopt.multiple ?? true) {
            const win = new BrowserWindow(winopt.window);
            win.show();
          } else {
            if (winopt._hasShow) {
              return;
            } else {
              winopt._hasShow = true;
              const win = new BrowserWindow(winopt.window);
              win.show();
              win.on("close", () => {
                winopt._hasShow = false;
              });
            }
          }
        }
      },
    });
  }


  system.registerFileOpener(".ln", {
    name: "快捷方式",
    icon: unknownIcon,
    hiddenInChosen: true,
    func: async (path, content) => {
      console.log(path)
      if (await system.fs.exists(content)) {
        try {
          system.openFile(content);
        } catch (e) {
          Dialog.showMessageBox({
            title: "错误",
            message: "无法打开快捷方式",
            type: "error",
          });
        }
      } else {
        Dialog.showMessageBox({
          title: "错误",
          message: "无法打开快捷方式，目标不存在",
          type: "error",
        });
      }
    },
  });



  system.registerFileOpener("dir", {
    name: "文件夹",
    icon: unknownIcon,
    hiddenInChosen: true,
    func: (path, content) => {
      const tempwindow = new BrowserWindow({
        width: 800,
        height: 600,
        center: true,
        title: t("computer"),
        content: "Computer",
        //icon: myComputerLogoIcon,
        icon: "computer",
        config: {
          content: content,
          path: path,
        },
      });
      tempwindow.show();
    },
  });

  system.registerFileOpener(".url", {
    name: "网址",
    icon: unknownIcon,
    func: async (path, content) => {
      const imgwindow = new BrowserWindow({
        width: 900,
        height: 600,
        icon: await dealIcon(await system.fs.stat(path), system),
        center: true,
        title: basename(path),
        content: "UrlBrowser",
        config: {
          content: content,
          path: path,
        },
      });
      imgwindow.show();
    },
  });


  const dateTimeT = new Tray({
    component: "DateTime",
  });
  dateTimeT.setContextMenu("DateTimePop", 320, 700);

}
