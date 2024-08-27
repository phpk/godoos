import { OsFileSystem, OsFileWithoutContent } from '@/system/core/FileSystem';

import { System } from '@/system';
import unknownicon from "@/assets/unknown.png"
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
  file: any,
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
    //console.log(file)
    if(!file.content){
      file.content = await system?.fs.readFile(file.path);
    }
    
    return dealExeIcon(file.content);
  }
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
}
