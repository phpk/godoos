import { useSettingsStore } from "@/stores/settings";
import * as netFile from "@/api/net/files";
//import * as localAuth from "@/api/local/service/auth";
import * as memberFile from "@/api/member/files";
let localFile: any = null;

async function getLocalFile() {
  if (!localFile) {
    //localFile = await import('@/api/local/service/files');
  }
  return localFile;
}
export const Files = async() => {     
    try {
        const settingsStore = useSettingsStore();
        if (settingsStore.config.system.userType === "person") {
            if(settingsStore.config.system.storeType === "local"){
                if(settingsStore.systemInfo.isWeb){
                    return netFile;
                }else{
                    return await getLocalFile();
                }
                
            }else{
                return netFile;
            }
        }
        else{
            return memberFile;
        }
    }catch (error) {
        console.error(error);
    }
    // if (settingsStore.config.system.userType === "person") {
    //     if(settingsStore.config.system.storeType === "local"){
    //         if(settingsStore.systemInfo.isWeb){
    //             return netFile;
    //         }else{
    //             return await getLocalFile();
    //         }
            
    //     }else{
    //         return netFile;
    //     }
    // }
    // else {
    //     return memberFile;
    // }
}