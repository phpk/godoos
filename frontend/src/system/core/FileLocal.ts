// import { exists, readDir, stat, readFile, remove, mkdir, rename, copyFile, writeFile, watch, BaseDirectory } from '@tauri-apps/plugin-fs';
// import { join, appDataDir } from '@tauri-apps/api/path';
// import { OsFileMode } from './FileMode';

// // async function processEntriesRecursive(parent:string, entries:any) {
// //     for (const entry of entries) {
// //       console.log(`Entry: ${entry.name}`);
// //       if (entry.isDirectory) {
// //          const dir = await join(parent, entry.name);
// //         processEntriesRecursive(dir, await readDir(dir, { baseDir: basedir }))
// //       }
// //     }
// //   }
// const basedir = BaseDirectory.AppData
// export async function appPath() {
//     const resourcePath = await appDataDir();
//     const appDir = await join(resourcePath, "os");
//     if (!await exists(appDir, { baseDir: basedir })) {
//         await mkdir(appDir, { baseDir: basedir });
//     }
//     return appDir;
// }
// export async function resolvePath(path: string) {
//     const appDir = await appPath();
//     return await join(appDir, path);
// }

// export const useLocalFile = () => {

//     return {
//         async readdir(path: string) {
//             path = await resolvePath(path);
//             const entries = await readDir(path, { baseDir: basedir });
//             //processEntriesRecursive(path, entries);
//             return entries;
//         },
//         async stat(path: string) {
//             path = await resolvePath(path);
//             const response = await stat(path, { baseDir: basedir })
//             return response;
//         },
//         async chmod(path: string, mode: OsFileMode) {
//             console.log(path, mode)
//             return true;
//         },
//         async search(path: string, options: any){
//             console.log(path, options)
//             return true;
//         },
//         async exists(path: string) {
//             path = await resolvePath(path);
//             return await exists(path, { baseDir: basedir });
//         },
//         async readFile(path: string) {
//             path = await resolvePath(path);
//             return await readFile(path, { baseDir: basedir });
//         },
//         async unlink(path: string) {
//             path = await resolvePath(path);
//             return await remove(path, { baseDir: basedir });
//         },
//         async rename(oldPath: string, newPath: string) {
//             oldPath = await resolvePath(oldPath);
//             newPath = await resolvePath(newPath);
//             return await rename(oldPath, newPath, { oldPathBaseDir: basedir, newPathBaseDir: basedir });
//         },
//         async rmdir(path: string) {
//             path = await resolvePath(path);
//             return await remove(path, { baseDir: basedir });
//         },
//         async mkdir(path: string) {
//             path = await resolvePath(path);
//             return await mkdir(path, { baseDir: basedir });
//         },
//         async copyFile(srcPath: string, dstPath: string) {
//             srcPath = await resolvePath(srcPath);
//             dstPath = await resolvePath(dstPath);
//             return await copyFile(srcPath, dstPath, { fromPathBaseDir: basedir, toPathBaseDir: basedir });
//         },
//         async writeFile(path: string, content: string | undefined) {
//             path = await resolvePath(path);
//             let encoder = new TextEncoder();
//             let data = encoder.encode(content);
//             console.log(data)
//             return await writeFile(path, data, { baseDir: basedir });
//         },
//         async appendFile(path: string, content: string | undefined) {
//             path = await resolvePath(path);
//             let encoder = new TextEncoder();
//             let data = encoder.encode(content);
//             return await writeFile(path, data, { baseDir: basedir, append: true });
//         },
//         serializeFileSystem() {
//             return Promise.reject('fs 不支持序列化');
//         },
//         deserializeFileSystem() {
//             return Promise.reject('fs 不支持序列化');
//         },
//         async removeFileSystem() {
//             const appDir = await appPath();
//             return await remove(appDir, { baseDir: basedir });
//             //return await remove(basedir, { baseDir: basedir });
//         },
//         registerWatcher(path: string, callback: any) {
//             watch(path, callback, { baseDir: basedir });
//         },
//     }
// };