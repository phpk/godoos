import { md5 } from "js-md5";
import { fetchGet, fetchPost, getFileUrl } from "../config.ts";
import { OsFileMode } from '../core/FileMode';
const API_BASE_URL = getFileUrl()
// import { notifyError } from "@/util/msg";
export async function handleReadDir(path: any): Promise<any> {
  const res = await fetchGet(`${API_BASE_URL}/read?path=${encodeURIComponent(path)}`);
  if (!res.ok) {
    return false;
  }
  return await res.json();
}
export async function handleStat(path: string): Promise<any> {
  const res = await fetchGet(`${API_BASE_URL}/stat?path=${encodeURIComponent(path)}`);
  if (!res.ok) {
    return false;
  }
  return await res.json();
}

export async function handleChmod(path: string, mode: string): Promise<any> {
  const res = await fetchPost(`${API_BASE_URL}/chmod`, JSON.stringify({ path, mode }));
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
  const res = await fetchGet(`${API_BASE_URL}/exists?path=${encodeURIComponent(path)}`);
  if (!res.ok) {
    return false;
  }
  return await res.json();
}

export async function handleReadFile(path: string, header?: any): Promise<any> {
  //const userType = getSystemConfig().userType
  let head = {}
  if (header) {
    head = {
      pwd: header.pwd !== '' ? md5(header.pwd) : ''
    }
  }

  const res = await fetchGet(`${API_BASE_URL}/readfile?path=${encodeURIComponent(path)}`, head);
  if (!res.ok) {
    return false;
  }
  return await res.json();
}

export async function handleUnlink(path: string): Promise<any> {
  const res = await fetchGet(`${API_BASE_URL}/unlink?path=${encodeURIComponent(path)}`);
  if (!res.ok) {
    return false;
  }
  return await res.json();
}
export async function handleSearch(path: string,query:string): Promise<any> {
  const res = await fetchGet(`${API_BASE_URL}/search?path=${encodeURIComponent(path)}&query=${encodeURIComponent(query)}`);
  if (!res.ok) {
    return false;
  }
  return await res.json();
}

export async function handleClear(): Promise<any> {
  const res = await fetchGet(`${API_BASE_URL}/clear`);
  if (!res.ok) {
    return false;
  }
  return await res.json();
}

export async function handleRename(oldPath: string, newPath: string): Promise<any> {
  let params = `oldPath=${encodeURIComponent(oldPath)}&newPath=${encodeURIComponent(newPath)}`
  const res = await fetchGet(`${API_BASE_URL}/rename?${params}`);
  if (!res.ok) {
    return false;
  }
  return await res.json();
}

export async function handleMkdir(dirPath: string): Promise<any> {
  const res = await fetchPost(`${API_BASE_URL}/mkdir?dirPath=${encodeURIComponent(dirPath)}`, {});
  if (!res.ok) {
    return false;
  }
  return await res.json();
}

export async function handleRmdir(dirPath: string): Promise<any> {
  const res = await fetchGet(`${API_BASE_URL}/rmdir?dirPath=${encodeURIComponent(dirPath)}`);
  if (!res.ok) {
    return false;
  }
  return await res.json();
}

export async function handleCopyFile(srcPath: string, dstPath: string): Promise<any> {
  const res = await fetchGet(`${API_BASE_URL}/copyfile?srcPath=${encodeURIComponent(srcPath)}&dstPath=${encodeURIComponent(dstPath)}`);
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
  formData.append('content', blobContent);
  return formData
}
export async function handleWriteFile(filePath: string, content: any, header?: { [key: string]: any }): Promise<any> {
  //console.log(content)
  const formData = getFormData(content);
  const head: any = header ? { ...header } : {}
  if (head.pwd && head.pwd !== '') {
    head.pwd = md5(head.pwd)
  }
  const url = `${API_BASE_URL}/writefile?path=${encodeURIComponent(filePath)}`;
  // const url = `${API_BASE_URL}/writefile?filePath=${encodeURIComponent(filePath)}`;
  const res = await fetchPost(url, formData, head);
  if (!res.ok) {
    return false;
  }
  return await res.json();
}

export async function handleAppendFile(filePath: string, content: string | Blob): Promise<any> {
  const formData = getFormData(content);
  const url = `${API_BASE_URL}/appendfile?filePath=${encodeURIComponent(filePath)}`;
  const res = await fetchPost(url, formData);
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
export async function handleZip(path: string, ext: string): Promise<any> {
  const res = await fetchGet(`${API_BASE_URL}/zip?path=${encodeURIComponent(path)}&ext=${ext}`);
  if (!res.ok) {
    return false;
  }
  return await res.json();
}
export async function handleUnZip(path: string): Promise<any> {
  const res = await fetchGet(`${API_BASE_URL}/unzip?path=${encodeURIComponent(path)}`);
  if (!res.ok) {
    return false;
  }
  return await res.json();
}
export const useOsFile = () => {
  return {

    async readdir(path: string) {
      const response = await handleReadDir(path);
      if (response && response.data) {
        return response.data;
      }
      return [];
    },
    async stat(path: string) {
      const response = await handleStat(path);
      if (response && response.data) {
        return response.data;
      }
      return false;
    },
    async chmod(path: string, mode: OsFileMode) {
      const modes = osFileModeToOctal(mode)
      const response = await handleChmod(path, modes);
      if (response) {
        return response;
      }
      return false;
    },
    async exists(path: string) {
      const response = await handleExists(path);
      if (response && response.data) {
        return response.data;
      }
      return false;
    },
    async readFile(path: string, header?: { [key: string]: string }) {
      // 注意：handleReadFile返回的是JSON，但通常readFile期望返回Buffer或字符串
      const response = await handleReadFile(path, header);
      if (response && response.code === 0) {
        return response.data;
      }
      if (response && response.code == -1) {
        return response
      }
      return false;
    },
    async unlink(path: string) {
      // const fun = path.indexOf('/F') === 0 ? handleShareUnlink : handleUnlink
      const response = await handleUnlink(path);
      if (response) {
        return response;
      }
      return false;
    },
    async rename(oldPath: string, newPath: string) {
      const response = await handleRename(oldPath, newPath);
      if (response) {
        return response;
      }
      return false;
    },
    async rmdir(path: string) {
      // const fun = path.indexOf('/F') === 0 ? handleShareUnlink : handleRmdir
      const response = await handleRmdir(path);
      if (response) {
        return response;
      }
      return false;
    },
    async mkdir(path: string) {
      const response = await handleMkdir(path);
      if (response) {
        return response;
      }
      return false;
    },
    async copyFile(srcPath: string, dstPath: string) {
      const response = await handleCopyFile(srcPath, dstPath);
      if (response) {
        return response;
      }
      return false;
    },
    async writeFile(path: string, content: string | Blob, header?: { [key: string]: any }) {
      const response = await handleWriteFile(path, content, header);
      if (response) {
        return response;
      }
      return false;
    },
    async appendFile(path: string, content: string | Blob) {
      const response = await handleAppendFile(path, content);
      if (response) {
        return response;
      }
      return false;
    },
    async search(path: string, query: string) {
      const response = await handleSearch(path, query);
      if (response && response.data) {
        return response.data;
      }
      return false;
    },
    async zip(path: string, ext: string) {
      const response = await handleZip(path, ext);
      if (response) {
        return response;
      }
      return false;
    },
    async unzip(path: string) {
      const response = await handleUnZip(path);
      if (response) {
        return response;
      }
      return false;
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