import * as fs from '@tauri-apps/plugin-fs';
import * as ps from '@tauri-apps/api/path';
import { base64ToBuffer, isBase64 } from '@/utils/file'
import { getUsername } from '@/utils/request';
const basedir = fs.BaseDirectory.Home
export async function appPath() {
    const resourcePath = await ps.homeDir();
    //const resourcePath = await ps.appDataDir();
    const username:string = getUsername();
    const appDir = await ps.join(resourcePath, ".godoos", username);
    if (!await fs.exists(appDir, { baseDir: basedir })) {
        await fs.mkdir(appDir, { baseDir: basedir, recursive: true });
    }
    return appDir;
}
export async function resolvePath(path: string) {
    const appDir = await appPath();
    return await ps.join(appDir, path);
}
export async function parseDir(path:string,entries: any) {
    const appDir = await appPath();
    path = path.replace(appDir, '');

    entries.forEach(async(entry:any) => {
        entry.path = join(path, entry.name); 
        entry.isPwd = false;      
        if(entry.isFile){
            const ns = entry.name.split('.')
            entry.ext = ns[1];
            entry.title = ns[0];
            if(entry.ext == 'exe' || entry.ext == 'lnk'){
                const entryPath = join(appDir, entry.path)
                entry.content = await fs.readTextFile(entryPath, { baseDir: basedir })
            }
        }else{
            entry.title = entry.name;
        }
    });
    return {
        success: true,
        data: entries
    };
}
export async function read(path: string, pwd?: string) {
    path = await resolvePath(path);
    const entries = await fs.readDir(path, { baseDir: basedir });
    //processEntriesRecursive(path, entries);
    return parseDir(path,entries);
}
export async function readFile(path: string, pwd?: string) {
    path = await resolvePath(path);
    const data = await fs.readTextFile(path, { baseDir: basedir });
    return {
        success: true,
        data: data
    }
}
export async function stat(path: string) {
    path = await resolvePath(path);
    const response = await fs.stat(path, { baseDir: basedir })
    return response;
}
export async function desktop() {
    const apps = await read("C/Users/Desktop");
    //console.log(apps)
    const menulist = await read("C/Users/Menulist");
    return { apps : apps.data, menulist: menulist.data };
}
export async function exists(path: string) {
    path = await resolvePath(path);
    return await fs.exists(path, { baseDir: basedir });

}
export async function mkdir(path: string) {
    if (path.length < 2 || path.charAt(1) == 'B') {
        return false;
    }
    path = await resolvePath(path);
    return await fs.mkdir(path, { baseDir: basedir });
}
export async function rmdir(dirPath: string) {
    if (dirPath.length < 3) {
        return false;
    }
    const ext = dirPath.split('.').pop();
    if (ext == 'exe') {
        return false;
    }
    dirPath = await resolvePath(dirPath);
    return await fs.remove(dirPath, { baseDir: basedir });
}
export function restore(dirPath: string) {
    if (dirPath.length < 2) {
        return false;
    }

    return false;
}
export function favorite(path: string) {
    if (path.length < 2) {
        return false;
    }
    return false;
}
export function pwd(path: string, pwd: string) {
    if (path.length < 2) {
        return false;
    }
    return false;
}
export function unpwd(path: string, pwd: string) {
    if (path.length < 2) {
        return false;
    }
    return false;
}
export async function rename(oldPath: string, newPath: string) {
    if (oldPath.length < 2) {
        return false;
    }
    oldPath = await resolvePath(oldPath);
    newPath = await resolvePath(newPath);
    return await fs.rename(oldPath, newPath, { oldPathBaseDir: basedir, newPathBaseDir: basedir });

}
export async function clear() {
    const appDir = await appPath();
    return await fs.remove(appDir, { baseDir: basedir });
}
export function search(path: string, query: string) {
    return false
}
export async function copy(srcPath: string, dstPath: string) {
    srcPath = await resolvePath(srcPath);
    dstPath = await resolvePath(dstPath);
    return await fs.copyFile(srcPath, dstPath, { fromPathBaseDir: basedir, toPathBaseDir: basedir });
}


export async function unlink(path: string) {
    if (path.length < 2) {
        return false;
    }
    const ext = path.split('.').pop();
    if (ext == 'exe') {
        return false;
    }
    path = await resolvePath(path);
    return await fs.remove(path, { baseDir: basedir });
}
export function parserFormData(content: any, contentType: any) {
    if (!content || content == '') {
        return new Blob([], { type: 'text/plain;charset=utf-8' });
    }
    if (contentType == 'text') {
        return new Blob([content], { type: 'text/plain;charset=utf-8' });
    }
    else if (contentType == 'base64') {
        if (content.indexOf(";base64,") > -1) {
            const parts = content.split(";base64,");
            content = parts[1];
        }
        content = base64ToBuffer(content);
        return new Blob([content]);
    }
    else if (typeof content === 'object' && content !== null && 'data' in content && Array.isArray(content.data)) {
        return new Blob([new Uint8Array(content.data).buffer]);
    }
    else if (contentType == 'buffer') {
        return new Blob([content]);
    }
}
export function getFormData(content: any) {
    let blobContent: Blob;
    //console.log(content)
    if (typeof content === 'string') {
        if (content) { // 检查 content 是否为空
            if (isBase64(content)) {
                if (content.indexOf(";base64,") > -1) {
                    const parts = content.split(";base64,");
                    content = parts[1];
                }
                content = base64ToBuffer(content);
                blobContent = new Blob([content]);
            } else {
                //console.log(content)
                blobContent = new Blob([content], { type: 'text/plain;charset=utf-8' });
            }
        } else {
            // 处理 content 为空的情况
            blobContent = new Blob([], { type: 'text/plain;charset=utf-8' });
        }
    }
    else if (content instanceof Blob) {
        // 如果是Blob，直接使用
        blobContent = content;
    }
    else if ('data' in content && Array.isArray(content.data)) {
        // 假设data属性是一个字节数组，将其转换为ArrayBuffer
        const arrayBuffer = new Uint8Array(content.data).buffer;
        //console.log(arrayBuffer)
        blobContent = new Blob([arrayBuffer]);
    } else if (content instanceof ArrayBuffer) {
        // 如果已经是ArrayBuffer，直接使用
        blobContent = new Blob([content]);
    }
    else if (content instanceof Array || content instanceof Object) {
        // 如果是数组
        blobContent = new Blob([JSON.stringify(content)], { type: 'text/plain;charset=utf-8' });
    } else {
        throw new Error('Unsupported content format');
    }

    const formData = new FormData();
    formData.append('content', blobContent);
    return formData
}
export async function writeFile(path: string, data: any, pwd?: string) {

    if (path.length < 2) {
        return false;
    }
    path = await resolvePath(path);
    await fs.writeTextFile(path, data, { baseDir: basedir });
    return true
}
export async function appendFile(path: string, data: any) {
    if (path.length < 2) {
        return false;
    }
    path = await resolvePath(path);
    return await fs.writeFile(path, data, { baseDir: basedir, append: true });
}
export function zip(path: string, ext: string) {
    if (path.length < 2) {
        return false;
    }
    return false
}
export function unzip(path: string) {
    if (path.length < 2) {
        return false;
    }
    return false
}
export function isDesktop(path: string) {
    const sp = getSp(path)
    const arr = path.split(sp)
    return arr[1] === 'C' && arr[2] === 'Users' && arr[3] === 'Desktop' && arr.length === 4
}
export function join(path: string, ...paths: string[]) {
    const sp = getSp(path)
    if (path.endsWith(sp)) {
        return path + paths.join(sp)
    } else {
        return path + sp + paths.join(sp)
    }
}
export function basename(path: string): string {
    const sp = getSp(path)
    return path.split(sp).pop() || path
}
export function dirname(path: string): string {
    const sp = getSp(path)
    if (path.indexOf(".") > -1) {
        return path.split(sp).slice(0, -1).join(sp)
    } else {
        return path
    }
    //return path.split(sp).slice(0, -1).join(sp)
}
export function getSp(path: string): string {
    if (path.indexOf("\\") > -1) {
        return "\\"
    } else {
        return "/"
    }
}
export function getParentPath(path: string): string {
    const sp = getSp(path)
    const arr = path.split(sp);
    arr.pop();
    return arr.join(sp);
}
export function getTopPath(path: string) {
    const sp = getSp(path)
    const arr = path.split(sp);
    if (arr[0] == "") {
        return arr[1]
    } else {
        return arr[0]
    }
}
export function getExt(path: string) {
    return path.split('.').pop() || ''
}

