import { OsFileWithoutContent, BrowserWindow, Notify } from '@/system';
import { t } from "@/i18n"
export function useAppMenu(item: OsFileWithoutContent, sys: any, props: any) {
  let menuArr: any = [];
  const ext: any = item.name.split(".").pop();
  const picExt = ['jpg', 'jpeg', 'png', 'gif', 'bmp', 'webp', 'tiff'];
  if (picExt.includes(ext)) {
    menuArr.push({
      label: t('piceditor'),
      click: () => {
        const photoWindow = new BrowserWindow({
          width: 800,
          height: 600,
          icon: "picedit",
          center: true,
          title: t('piceditor'),
          url: "/picedit/index.html",
          config: item
        });
        photoWindow.show()
      },
    },)
  }
  const zipExt = ['zip', 'tar', 'gz', 'bz2']
  const unzipSucess = (res: any) => {
    //console.log(res)
    if (!res || res.code < 0) {
      new Notify({
        title: t('tips'),
        content: t('error'),
      });
    } else {
      props.onRefresh();
      new Notify({
        title: t('tips'),
        content: t('file.unzip.success'),
      });
    }

  };
  if (zipExt.includes(ext)) {
    menuArr.push({
      label: t('unzip'),
      click: () => {
        //console.log(item.path)
        sys.fs.unzip(item.path).then((res: any) => {
          unzipSucess(res)
        })
      },
    },)
  }
  return menuArr;
}
//export { useAppMenu };

