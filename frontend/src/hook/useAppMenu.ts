import { OsFileWithoutContent } from '@/system/core/FileSystem';
import { BrowserWindow } from '@/system/window/BrowserWindow';
export function useAppMenu(item: OsFileWithoutContent, _: number) {
  let menuArr: any = [];
  const ext:any = item.name.split(".").pop();
  const picExt = ['jpg', 'jpeg', 'png', 'gif', 'bmp', 'webp', 'tiff'];
  if (picExt.includes(ext)) {
    menuArr.push({
      label: '编辑',
      click: () => {
        const photoWindow = new BrowserWindow({
          width: 800,
          height: 600,
          icon: "picedit",
          center: true,
          title: '图像绘图',
          url: "/picedit/index.html",
          config: item
        });
        photoWindow.show()
      },
    },)
  } 
  // const fileExt = ['txt', 'html', 'json', 'xml', 'css', 'js','vue','go','php','java','py']
  // if (fileExt.includes(ext)) {
  //   menuArr.push({
  //     label: '文本编辑',
  //     click: () => {
  //       const textWindow = new BrowserWindow({
  //         width: 800,
  //         height: 600,
  //         icon:  "editorbt",
  //         center: true,
  //         title: '文件编辑',
  //         url: "/text/index.html",
  //         config: item
  //       });
  //       textWindow.show()
  //     },
  //   },)
  // }

  return menuArr;
}
//export { useAppMenu };

