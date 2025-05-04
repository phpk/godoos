import * as files from './files'
import * as fs from '@tauri-apps/plugin-fs'
import * as ps from '@tauri-apps/api/path';
interface OsFileInfo {
    Name: string;
    Path: string;
    OldPath: string;
    ParentPath: string;
    Content: string;
    Ext: string;
    Title: string;
    ID: number;
    IsFile: boolean;
    IsDir: boolean;
    IsSymlink: boolean;
    Size: number;
    ModTime: Date;
    AccessTime: Date;
    CreateTime: Date;
    Mode: number;
}

interface RootObject {
    UpdateTime: Date;
    Desktop: OsFileInfo[];
    Menulist: OsFileInfo[];
}
const RootAppList = [
    { name: "computer", icon: "diannao", position: "Desktop,Menulist" },
    { name: "appstore", icon: "store", position: "Desktop,Menulist" },
    { name: "localchat", icon: "chat", position: "Desktop" },
    { name: "document", icon: "word", position: "Desktop" },
    { name: "excel", icon: "excel", position: "Desktop" },
    { name: "markdown", icon: "markdown", position: "Desktop" },
    { name: "mindmap", icon: "mindexe", position: "Desktop" },
    { name: "ppt", icon: "pptexe", position: "Desktop" },
    { name: "fileEditor", icon: "editorbt", position: "Desktop" },
    { name: "board", icon: "kanban", position: "Desktop" },
    { name: "whiteBoard", icon: "baiban", position: "Desktop" },
    { name: "piceditor", icon: "picedit", position: "Desktop" },
    { name: "gantt", icon: "gant", position: "Desktop" },
    { name: "browser", icon: "brower", position: "Desktop,Menulist" },
    { name: "setting", icon: "setting", position: "Menulist" },
    { name: "system.version", icon: "info", position: "Menulist" },
    { name: "process.title", icon: "progress", position: "Menulist" },
    { name: "calculator", icon: "calculator", position: "Menulist" },
    { name: "calendar", icon: "calendar", position: "Menulist" },
    { name: "musicStore", icon: "music", position: "Menulist" },
    { name: "gallery", icon: "gallery", position: "Menulist" },
    { name: "aiHelper", icon: "aiassistant", position: "Desktop" },
    { name: "aiModule", icon: "aidown", position: "Desktop" },
    { name: "aiSetting", icon: "aisetting", position: "Menulist" },
];
const InitPaths = [
    'Desktop',
    'Menulist',
    'Documents',
    'Downloads',
    'Music',
    'Pictures',
    'Videos',
    'Schedule',
    'Reciv',
];
const InitDocPath = [
    'Word',
    'Markdown',
    'PPT',
    'Baiban',
    'Kanban',
    'Excel',
    'Mind',
    'Screenshot',
    'ScreenRecording',
];
export async function initOsSystem(): Promise<void> {
    try {
        const basePath = await files.appPath(); // 需要你自己实现 libs.getOsDir()
        const osCpath = await ps.join(basePath, 'C');

        if (!(await fs.exists(osCpath))) {
            await fs.mkdir(osCpath, { recursive: true });
        }

        const baseOsDir = ['D', 'E', 'B'];
        for (const dir of baseOsDir) {
            const dirPath = await ps.join(basePath, dir);
            if (!(await fs.exists(dirPath))) {
                await fs.mkdir(dirPath, { recursive: true });
            }
        }

        const systemPath = await ps.join(osCpath, 'System');
        if (!(await fs.exists(systemPath))) {
            await fs.mkdir(systemPath, { recursive: true });
        }

        const userPath = await ps.join(osCpath, 'Users');
        if (!(await fs.exists(userPath))) {
            await fs.mkdir(userPath, { recursive: true });


            for (const dir of InitPaths) {
                const dirPath = await ps.join(userPath, dir);
                if (!(await fs.exists(dirPath))) {
                    await fs.mkdir(dirPath, { recursive: true });
                }
            }
            const docpath = await ps.join(userPath, 'Documents');
            for (const dir of InitDocPath) {
                const dirPath = await ps.join(docpath, dir);
                if (!(await fs.exists(dirPath))) {
                    await fs.mkdir(dirPath, { recursive: true });
                }
            }

            const applist = getInitRootList();
            const desktopPath = await ps.join(userPath, 'Desktop');
            for (const app of applist.Desktop) {
                const appPath = await ps.join(desktopPath, app.Name);
                if (!(await fs.exists(appPath))) {
                    await fs.writeFile(appPath, new TextEncoder().encode(app.Content));
                }
            }

            const menulistPath = await ps.join(userPath, 'Menulist');
            for (const app of applist.Menulist) {
                const appPath = await ps.join(menulistPath, app.Name);
                if (!(await fs.exists(appPath))) {
                    await fs.writeFile(appPath, new TextEncoder().encode(app.Content));
                }
            }
        }
    } catch (err) {
        throw new Error(`InitOsSystem error: ${err}`);
    }
}



function getInitRootList(): RootObject {
    const nowtime = new Date();
    let id = 1;

    const desktopApps: OsFileInfo[] = [];
    const menulistApps: OsFileInfo[] = [];

    for (const app of RootAppList) {
        const positions = app.position.split(",");
        const content = `link::Desktop::${app.name}::${app.icon}`;

        for (const pos of positions) {
            const baseProps = {
                Name: `${app.name}.exe`,
                OldPath: `/C/Users/${pos}/${app.name}.exe`,
                ParentPath: `/C/Users/${pos}`,
                Content: content,
                Ext: "exe",
                Title: app.name,
                ID: id++,
                IsFile: true,
                IsDir: false,
                IsSymlink: false,
                Size: content.length,
                ModTime: nowtime,
                AccessTime: nowtime,
                CreateTime: nowtime,
                Mode: 511,
            };

            if (pos === "Desktop") {
                desktopApps.push({
                    ...baseProps,
                    Path: `/C/Users/Desktop/${app.name}.exe`,
                });
            } else if (pos === "Menulist") {
                menulistApps.push({
                    ...baseProps,
                    Path: `/C/Users/Menulist/${app.name}.exe`,
                });
            }
        }
    }

    return {
        UpdateTime: nowtime,
        Desktop: desktopApps,
        Menulist: menulistApps,
    };
}