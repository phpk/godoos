import { OsFileWithoutContent } from '../system/core/FileSystem';
export function turnServePath(path: string): string {
    const replaceUrl = path.indexOf('/F/myshare') === 0 ? '/F/myshare' : '/F/othershare'
    return path.replace(replaceUrl, 'data/userData')
}
//将浏览器地址转化为本地地址
export function turnLocalPath(path: string, newTemp: string, type?: number): string {
    if (type && type === 1) {
        const arr = newTemp.split('/')
        newTemp = `/${arr[1]}/${arr[2]}`
    }
    return path.replace('data/userData', newTemp)
}
export function isShareFile(path: string) : boolean {
    // console.log('是否是共享文件：',path,path.indexOf('/F/myshare') === 0 || path.indexOf('/F/othershare') === 0);
    return path.indexOf('/F/myshare') === 0 || path.indexOf('/F/othershare') === 0
}
//是否是分享根目录
export function isRootShare(path: string) : boolean {
    // console.log('根路径：',path === '/F/myshare' || path === '/F/othershare');
    return path === '/F/myshare' || path === '/F/othershare'
}
//返回路径
export function turnFilePath(file: OsFileWithoutContent) : string {
    const arr = file.path.split('/')
    return `/${arr[1]}/${arr[2]}/${arr[arr.length - 1]}`
}
//返回是我的分享还是接收分享的文件路径
export function getFileRootPath(path: string) : string {
    const rootPath = path.indexOf('/F/myshare') === 0 ? '/F/myshare' : '/F/othershare'
    return rootPath
}