// import { useSystemStore } from "./system.ts"
// const systemStore = useSystemStore();
// const API_BASE_URL = systemStore.getFileUrl()
import { getSystemKey } from "../config.ts";
const API_BASE_URL = getSystemKey('apiUrl') + "/file"
import { OsFileMode } from '../core/FileMode';
export async function handleReadDir(path: string): Promise<any> {
    // if(window.go){
    //     return await window.go.app.App.ReadDir(path)
    // }else{

    // }
    const res = await fetch(`${API_BASE_URL}/read?path=${encodeURIComponent(path)}`);
    if (!res.ok) {
        return false;
    }
    return await res.json();

}

export async function handleStat(path: string): Promise<any> {
    const res = await fetch(`${API_BASE_URL}/stat?path=${encodeURIComponent(path)}`);
    if (!res.ok) {
        return false;
    }
    return await res.json();
}
export async function handleChmod(path: string, mode: string): Promise<any> {
    const res = await fetch(`${API_BASE_URL}/chmod`, {
        method: 'POST',
        body: JSON.stringify({ path, mode }),
    });
    if (!res.ok) {
        return false;
    }
    return await res.json();
}
function osFileModeToOctal(mode: OsFileMode): string {
    switch (mode) {
        case OsFileMode.Read:
            return "400";
        case OsFileMode.Write:
            return "200";
        case OsFileMode.Execute:
            return "100";
        case OsFileMode.ReadWrite:
            return "600";
        case OsFileMode.ReadExecute:
            return "500";
        case OsFileMode.WriteExecute:
            return "300";
        case OsFileMode.ReadWriteExecute:
            return "700";
        default:
            throw new Error("Invalid OsFileMode");
    }
}
export async function handleExists(path: string): Promise<any> {
    const res = await fetch(`${API_BASE_URL}/exists?path=${encodeURIComponent(path)}`);
    if (!res.ok) {
        return false;
    }
    return await res.json();
}

export async function handleReadFile(path: string): Promise<any> {
    const res = await fetch(`${API_BASE_URL}/readfile?path=${encodeURIComponent(path)}`);
    if (!res.ok) {
        return false;
    }
    return await res.json();
}

export async function handleUnlink(path: string): Promise<any> {
    const res = await fetch(`${API_BASE_URL}/unlink?path=${encodeURIComponent(path)}`);
    if (!res.ok) {
        return false;
    }
    return await res.json();
}

export async function handleClear(): Promise<any> {
    const res = await fetch(`${API_BASE_URL}/clear`);
    if (!res.ok) {
        return false;
    }
    return await res.json();
}

export async function handleRename(oldPath: string, newPath: string): Promise<any> {
    const res = await fetch(`${API_BASE_URL}/rename?oldPath=${encodeURIComponent(oldPath)}&newPath=${encodeURIComponent(newPath)}`);
    if (!res.ok) {
        return false;
    }
    return await res.json();
}

export async function handleMkdir(dirPath: string): Promise<any> {
    const res = await fetch(`${API_BASE_URL}/mkdir?dirPath=${encodeURIComponent(dirPath)}`, { method: 'POST' });
    if (!res.ok) {
        return false;
    }
    return await res.json();
}

export async function handleRmdir(dirPath: string): Promise<any> {
    const res = await fetch(`${API_BASE_URL}/rmdir?dirPath=${encodeURIComponent(dirPath)}`);
    if (!res.ok) {
        return false;
    }
    return await res.json();
}


export async function handleCopyFile(srcPath: string, dstPath: string): Promise<any> {
    const res = await fetch(`${API_BASE_URL}/copyfile?srcPath=${encodeURIComponent(srcPath)}&dstPath=${encodeURIComponent(dstPath)}`);
    if (!res.ok) {
        return false;
    }
    return await res.json();
}
export function getFormData(content: any) {
    let blobContent: Blob;
    if (typeof content === 'string') {
        // 如果content是字符串，转换为Blob并指定文本类型
        blobContent = new Blob([content], { type: 'text/plain;charset=utf-8' });
    } else if ('data' in content && Array.isArray(content.data)) {
        // 假设data属性是一个字节数组，将其转换为ArrayBuffer
        const arrayBuffer = new Uint8Array(content.data).buffer;
        //console.log(arrayBuffer)
        blobContent = new Blob([arrayBuffer]);
    } else if (content instanceof ArrayBuffer) {
        // 如果已经是ArrayBuffer，直接使用
        blobContent = new Blob([content]);
    } else if (content instanceof Blob) {
        // 如果是Blob，直接使用
        blobContent = content;
    }
    else if (content instanceof Array || content instanceof Object) {
        // 如果是数组
        blobContent = new Blob([JSON.stringify(content)], { type: 'text/plain;charset=utf-8' });
    } else {
        throw new Error('Unsupported content format');
    }

    const formData = new FormData();
    formData.append('content', blobContent); // 可以自定义文件名
    return formData
}
export async function handleWriteFile(filePath: string, content: any): Promise<any> {
    const formData = getFormData(content);

    const url = `${API_BASE_URL}/writefile?filePath=${encodeURIComponent(filePath)}`;
    const res = await fetch(url, { method: 'POST', body: formData });
    if (!res.ok) {
        return false;
    }
    return await res.json();
}

export async function handleAppendFile(filePath: string, content: string | Blob): Promise<any> {
    const formData = getFormData(content);
    const url = `${API_BASE_URL}/appendfile?filePath=${encodeURIComponent(filePath)}`;
    const res = await fetch(url, { method: 'POST', body: formData });
    if (!res.ok) {
        return false;
    }
    return await res.json();
}
export function handleWatch(path: string, callback: any, errback: any) {
    if (typeof EventSource !== "undefined") {
        const source = new EventSource(`${API_BASE_URL}/watch?path=${encodeURIComponent(path)}`);
        source.onmessage = function (event) {
            callback(event)
        }
        // 当与服务器的连接打开时触发
        source.onopen = function () {
            console.log("Connection opened.");
        };

        // 当与服务器的连接关闭时触发
        source.onerror = function (event) {
            console.log("Connection closed.");
            errback(event)
        };
    } else {
        errback("Your browser does not support server-sent events.")
    }
}
export const useOsFile = () => {
    return {
        async readdir(path: string) {
            const response = await handleReadDir(path);
            if (response && response.data) {
                return response.data; // 假设返回的JSON包含名为"data"的字段，存储了文件内容
            }
            return [];
        },
        async stat(path: string) {
            const response = await handleStat(path);
            if (response && response.data) {
                return response.data; // 假设返回的JSON包含名为"data"的字段，存储了文件内容
            }
            return false;
        },
        async chmod(path: string, mode: OsFileMode) {
            const modes = osFileModeToOctal(mode)
            const response = await handleChmod(path, modes);
            if (response) {
                return response; // 假设返回的JSON包含操作状态或其他相关信息
            }
            return false;
        },
        async exists(path: string) {
            //return handleExists(path);
            const response = await handleExists(path);
            if (response && response.data) {
                return response.data; // 假设返回的JSON包含名为"data"的字段，存储了文件内容
            }
            return false;
        },
        async readFile(path: string) {
            // 注意：handleReadFile返回的是JSON，但通常readFile期望返回Buffer或字符串
            // 根据你的后端API实际返回的数据类型调整此方法
            const response = await handleReadFile(path);
            if (response && response.data) {
                return response.data; // 假设返回的JSON包含名为"data"的字段，存储了文件内容
            }
            return false;
        },
        async unlink(path: string) {
            //return handleUnlink(path);
            const response = await handleUnlink(path);
            if (response) {
                return response; // 假设返回的JSON包含名为"data"的字段，存储了文件内容
            }
            return false;
        },
        async rename(oldPath: string, newPath: string) {
            //return handleRename(oldPath, newPath);

            const response = await handleRename(oldPath, newPath);
            if (response) {
                return response; // 假设返回的JSON包含名为"data"的字段，存储了文件内容
            }
            return false;
        },
        async rmdir(path: string) {
            //return handleRmdir(path);
            const response = await handleRmdir(path);
            if (response) {
                return response; // 假设返回的JSON包含名为"data"的字段，存储了文件内容
            }
            return false;
        },
        async mkdir(path: string) {
            //return handleMkdir(path);
            const response = await handleMkdir(path);
            if (response) {
                return response; // 假设返回的JSON包含名为"data"的字段，存储了文件内容
            }
            return false;
        },
        async copyFile(srcPath: string, dstPath: string) {
            //return handleCopyFile(srcPath, dstPath);
            const response = await handleCopyFile(srcPath, dstPath);
            if (response) {
                return response; // 假设返回的JSON包含名为"data"的字段，存储了文件内容
            }
            return false;
        },
        async writeFile(path: string, content: string | Blob) {
            //return handleWriteFile(path, content);
            const response = await handleWriteFile(path, content);
            if (response) {
                return response; // 假设返回的JSON包含名为"data"的字段，存储了文件内容
            }
            return false;
        },
        async appendFile(path: string, content: string | Blob) {
            //return handleAppendFile(path, content);
            const response = await handleAppendFile(path, content);
            if (response) {
                return response; // 假设返回的JSON包含名为"data"的字段，存储了文件内容
            }
            return false;
        },
        async search(path: string, options: any) {
            console.log(path, options)
            return true;
        },
        serializeFileSystem() {
            return Promise.reject('fs 不支持序列化');
        },
        deserializeFileSystem() {
            return Promise.reject('fs 不支持序列化');
        },
        removeFileSystem() {
            return handleClear();
        },
        registerWatcher(path: string, callback: any, errback: any) {
            handleWatch(path, callback, errback);
        },
    }
};