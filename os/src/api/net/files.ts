import { base64ToBuffer, isBase64 } from '@/utils/file'
import { get, post } from '@/utils/request'
export async function read(path: string, pwd?: string) {
  const res = await get(`net/files/read`, { path, pwd })
  return res
}
export async function readFile(path: string, pwd?: string) {
  const res = await get(`net/files/readfile`, { path, pwd })
  return res
}
export async function stat(path: string) {
  const res = await get(`net/files/stat`, { path })
  return res.data
}
export async function desktop() {
  const res = await get(`net/files/desktop`)
  if (res && res.success) {
    return res.data;
  }
  return false
}
export async function exists(path: string) {
  const res = await get(`net/files/exists`, { path })
  return res.data

}
export function mkdir(dirPath: string) {
  if (dirPath.length < 2 || dirPath.charAt(1) == 'B') {
    return false;
  }
  return post(`net/files/mkdir`, {}, { dirPath }).then(res => res.success)
}
export function rmdir(dirPath: string) {
  if (dirPath.length < 3) {
    return false;
  }
  const ext = dirPath.split('.').pop();
  if (ext == 'exe') {
    return false;
  }
  return get(`net/files/rmdir`, { dirPath }).then(res => res.success)
}
export function restore(dirPath: string) {
  if (dirPath.length < 2) {
    return false;
  }

  return get(`net/files/restore`, { dirPath }).then(res => res.success)
}
export function favorite(path: string) {
  if (path.length < 2) {
    return false;
  }
  return get(`net/files/favorite`, { path }).then(res => res.success)
}
export function pwd(path: string, pwd: string) {
  if (path.length < 2) {
    return false;
  }
  return get(`net/files/pwd`, { path, pwd }).then(res => res.success)
}
export function unpwd(path: string, pwd: string) {
  if (path.length < 2) {
    return false;
  }
  return get(`net/files/unpwd`, { path, pwd }).then(res => res.success)
}
export function rename(oldPath: string, newPath: string) {
  if (oldPath.length < 2) {
    return false;
  }
  return get(`net/files/rename`, { oldPath, newPath }).then(res => res.success)

}
export function clear() {
  return get(`net/files/clear`).then(res => res.success)
}
export function search(path: string, query: string) {
  return get(`net/files/search`, { path, query }).then(res => res.data)
}
export function copy(srcPath: string, dstPath: string) {
  if (dstPath.length < 2) {
    return false;
  }
  return get(`net/files/copyfile`, { srcPath, dstPath }).then(res => res.success)
}


export function unlink(path: string) {
  if (path.length < 2) {
    return false;
  }
  const ext = path.split('.').pop();
  if (ext == 'exe') {
    return false;
  }
  return get(`net/files/unlink`, { path }).then(res => res.success)
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
export function writeFile(path: string, data: any, pwd?: string) {

  if (path.length < 2) {
    return false;
  }
  const formData = getFormData(data);
  if (!formData) {
    return false;
  }

  return post(`net/files/writefile`, formData, { path, pwd }).then(res => {
    //console.log(res)
    return res.success
  })
}
export function appendFile(path: string, data: any) {
  const formData = getFormData(data);
  if (!formData) {
    return false;
  }
  return post(`net/files/appendfile`, formData, { path }).then(res => res.success)
}
export function zip(path: string, ext: string) {
  if (path.length < 2) {
    return false;
  }
  return get(`net/files/zip`, { path, ext }).then(res => res.data)
}
export function unzip(path: string) {
  if (path.length < 2) {
    return false;
  }
  return get(`net/files/unzip`, { path }).then(res => res.data)
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
