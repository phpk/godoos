// import * as fs from "@tauri-apps/plugin-fs";
// import { BaseDirectory } from '@tauri-apps/plugin-fs';
// import { convertFileSrc } from "@tauri-apps/api/core";
// console.log(convertFileSrc("data"))
// const basePath = "/Users/rt008/www/osData";
// const baseFilePath = "/Users/rt008/www/osData/"
// export const useFile = () => {
//     return {
//         async readdir(path) {
//             return fs.readDir(basePath + path);
//         },
//         async stat(path) {
//             return fs.stat(basePath + path);
//         },
//         async exists(path) {
//             return fs.exists(basePath + path);
//         },
//         async readFile(path) {
//             return fs.readFile(basePath + path);
//         },
//         async unlink(path) {
//             return fs.remove(basePath + path);
//         },
//         async rename(oldPath, newPath) {
//             return fs.rename(basePath + oldPath, basePath + newPath);
//         },
//         async rmdir(path) {
//             return fs.remove(basePath + path);
//         },
//         async mkdir(path) {
//             return fs.mkdir(basePath + path);
//         },
//         async copyFile(src, dest) {
//             return fs.copyFile(basePath + src, basePath + dest);
//         },
//         async writeFile(path, content) {
//             return fs.writeFile(basePath + path, content);
//         },
//         async appendFile(path, content) {
//             return fs.writeFile(basePath + path, content, { append: true });
//         },
//         serializeFileSystem() {
//             return Promise.reject('fs 不支持序列化');
//         },
//         deserializeFileSystem(files) {
//             return Promise.reject('fs 不支持序列化');
//         },
//         removeFileSystem() {
//             return Promise.reject('fs 不支持清空');
//         },
//         registerWatcher(path, callback) {
//             return fs.watch(basePath + path, (event) => {
//                 callback(path);
//             });
//         },
//     }
// };