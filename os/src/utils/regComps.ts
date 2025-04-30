import SetAd from '@/components/settings/account/SetAd.vue';
import SetLanguage from '@/components/settings/account/SetLanguage.vue';
import SetLock from '@/components/settings/account/SetLock.vue';
import SetWallpaper from '@/components/settings/account/SetWallpaper.vue';
import SetNas from '@/components/settings/nas/SetNas.vue';
import SetLocal from '@/components/settings/proxy/SetLocal.vue';
import SetRemote from '@/components/settings/proxy/SetRemote.vue';
import SetBackup from '@/components/settings/system/SetBackup.vue';
import SetEditor from '@/components/settings/system/SetEditor.vue';
import SetPassword from '@/components/settings/system/SetPassword.vue';
import SetStorage from '@/components/settings/system/SetStorage.vue';
import SetUser from '@/components/settings/system/SetUser.vue';
export function registerComponents(app: { component: (arg0: string, arg1: any) => void; }) {
  app.component('SetUser', SetUser);
  app.component('SetStorage', SetStorage);
  app.component('SetEditor', SetEditor);
  app.component('SetBackup', SetBackup);
  app.component('SetPassword', SetPassword);
  app.component('SetLocal', SetLocal);
  app.component('SetRemote', SetRemote);
  app.component('SetNas', SetNas);
  app.component('SetWallpaper', SetWallpaper);
  app.component('SetLanguage', SetLanguage);
  app.component('SetLock', SetLock);
  app.component('SetAd', SetAd);
}
