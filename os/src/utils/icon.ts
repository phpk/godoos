//import { readFile } from '@/api/files';

export function dealExeIcon(content: string | null | undefined) {
  
  if (!content) return 'unknown';
  const exeContent = content.split('::');
  //console.log(exeContent)
  const iconImg = exeContent.slice(3).join('::');
  //console.log(iconImg)
  if (iconImg != 'undefined' && iconImg != '' && iconImg != null && iconImg) {
    //console.log(iconImg)
    return iconImg;
  } else {
    return 'unknown';
  }
}
export function dealIcon(file: any) {
  //console.log(file)
  if (!file) return 'unknown';
  if(typeof file.isDirectory === 'string'){
    if(file.isDirectory === 'true'){
      file.isDirectory = true;
    }else{
      file.isDirectory = false;
    }
  }
  if (file.isDirectory && file.path.length == 2) {
   // console.log(file.path)
    if (file.isSys) {
      return 'disk';
    } else {
      return 'disknet';
    }
  }
  if (file.isDirectory) {
    //return foldericon;
    return 'folder';
  }
  const ext = file.ext;
  if (ext === 'exe' || ext === 'url') {
    return dealExeIcon(file.content);
  }
  //console.log(ext)
  return ext;
}
