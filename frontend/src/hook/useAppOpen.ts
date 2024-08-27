import { OsFileWithoutContent } from '../system/core/FileSystem';
import { useSystem } from '../system';

function useAppOpen(type: 'apps' | 'magnet' | 'menulist') {
  const system = useSystem();
  const rootState = system._rootState;
  
  const appList = rootState[type];
  //console.log(appList)

  //console.log(appList)
  function openapp(item: OsFileWithoutContent) {
    //console.log(item)
    system?.openFile(item.path);
  }
  return {
    appList,
    openapp,
  };
}
export { useAppOpen };
