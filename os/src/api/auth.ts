import { useSettingsStore } from "@/stores/settings";
import * as netAuth from "@/api/net/auth";
//import * as localAuth from "@/api/local/service/auth";
import * as memberAuth from "@/api/member/auth";
let localAuth: any = null;

async function getLocalAuth() {
  if (!localAuth) {
    localAuth = await import('@/api/local/service/auth');
  }
  return localAuth;
}
export const Auth = async() => {  
    const settingsStore = useSettingsStore();
    if (settingsStore.config.system.userType === "person") {
        if(settingsStore.config.system.storeType === "local"){
            if(settingsStore.systemInfo.isWeb){
                return netAuth;
            }else{
                return await getLocalAuth();
            }
            
        }else{
            return netAuth;
        }
    }
    else {
        return memberAuth;
    }
}