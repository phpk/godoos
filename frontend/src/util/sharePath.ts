import { OsFileWithoutContent } from '../system/core/FileSystem';
export function turnServePath(path: string): string {
 const replaceUrl = path.indexOf('/F/myshare') === 0 ? '/F/myshare' : '/F/othershare'
 return path.replace(replaceUrl, 'data/userData')
}
export function turnLocalPath(path: string, newTemp: string): string {
    return path.replace('data/userData', newTemp)
}
export function isShareFile(path: string) : boolean {
    // console.log('是否是共享文件：',path,path.indexOf('/F/myshare') === 0 || path.indexOf('/F/othershare') === 0);
    return path.indexOf('/F/myshare') === 0 || path.indexOf('/F/othershare') === 0
}
//返回路径
export function turnFilePath(file: OsFileWithoutContent) : string {
    // console.log('是否是共享文件：',path,path.indexOf('/F/myshare') === 0 || path.indexOf('/F/othershare') === 0);
    const arr = file.path.split('/')
    return `/${arr[1]}/${arr[2]}/${arr[arr.length - 1]}`
}