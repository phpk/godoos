import { OsFileSystem, OsFileWithoutContent } from '@/system/core/FileSystem';

import { System } from '@/system';
const unknownicon = 'unknow'
export function dealExeIcon(content: string | null | undefined) {
  
  if (!content) return unknownicon;
  const exeContent = content.split('::');
  //console.log(exeContent)
  const iconImg = exeContent.slice(3).join('::');
  //console.log(iconImg)
  if (iconImg != 'undefined' && iconImg != '' && iconImg != null && iconImg) {
    return iconImg;
  } else {
    return unknownicon;
  }
}
export async function dealIcon(
  file: OsFileWithoutContent | null | undefined,
  system: System,
  stopCircle = false
) {
  if (!file) return unknownicon;
  if (file.isDirectory && file.parentPath === '/') {
    // 是挂载在根目录的卷
    if (system.fs instanceof OsFileSystem) {
      if (system.fs.checkVolumePath(file.path)) {
        //return volumeNetIcon;
        return 'disk';
      } else {
        //return volumeLocalIcon;
        return 'disknet';
      }
    }
  }
  if (file.isDirectory) {
    //return foldericon;
    return 'folder';
  }
  const ext = file.ext;
  if (ext === 'exe' || ext === 'url') {
    const content = await system?.fs.readFile(file.path);
    return dealExeIcon(content);
  }
  // if (ext === '.png') {
  //   if (system.fs instanceof OsFileSystem) {
  //     if (system.fs.checkVolumePath(file.path)) {
  //       return imageicon;
  //     } else {
  //       // 只有当非自定义文件系统和非自定义卷才直接展示icon
  //       const content = await system.fs.readFile(file.path);
  //       return 'data:image/png;base64,' + content || unknownicon;
  //     }
  //   } else {
  //     return imageicon;
  //   }
  // }
  if (ext === 'ln') {
    if (stopCircle) {
      return unknownicon;
    } else {
      const target = await system.fs.readFile(file.path);
      //console.log(target)
      if (target) {
        return dealIcon(await system.fs.stat(target), system, true);
      } else {
        return unknownicon;
      }
    }
  }
  return file.ext;
  // const exts = ['mp3','mp4','txt','js','json','xml','svg','html','css','bmp','gif','jpg','png','tiff','webp']
  // if(exts.includes(file.ext)){
  //   return '/image/ext/' + file.ext + '.svg'
  // }
  // if (ext === '.mp3') return audioicon;
  // if (ext === '.mp4') return videoicon;

  // if (ext === '.txt') return txtIcon;
  // if (ext === '.js') return jsIcon;
  // if (ext === '.json') return jsonIcon;
  // if (ext === '.xml') return xmlIcon;
  // if (ext === '.svg') return svgIcon;
  // if (ext === '.html') return htmlIcon;
  // if (ext === '.css') return cssIcon;

  // if (ext === '.bmp') return bmpIcon;
  // if (ext === '.gif') return gifIcon;
  // if (ext === '.jpg') return jpgIcon;
  // if (ext === '.png') return pngIcon;
  // if (ext === '.tiff') return tiffIcon;
  // if (ext === '.webp') return webpIcon;

  //return system.getOpener(ext)?.icon || unknownicon;
}
