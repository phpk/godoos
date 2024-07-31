<script setup lang="ts">
import { inject, onMounted } from "vue";
import { System } from "@/system";
import { notifySuccess, notifyError } from "@/util/msg";
import { useStoreStore } from "@/stores/store";
import { setSystemKey,parseJson } from "@/system/config";
import { t } from "@/i18n";
const sys: any = inject<System>("system");
const store = useStoreStore()

onMounted(async () => {
  await store.getList();
});

function setCache() {
  setSystemKey("intstalledPlugins", store.installedList);
  setTimeout(() => {
    sys.refershAppList();
  }, 1000);
}
async function install(item: any) {
  //console.log(item)
  if (item.isOut) {
    item.progress = 0
    const flag = await store.addOutList(item)
    if(!flag){
      notifyError(t("store.installError"))
    }
  }
  if (item.needDownload && !item.isDev) {
    await download(item)
  }
  
  if (item.needInstall) {
    item.progress = 0
    const completion = await fetch(store.apiUrl + '/store/install?name=' + item.name)
    if (!completion.ok) {
      notifyError(t("store.installError"))
      return
    }
    const res = await completion.json()
    if (res.code && res.code < 0) {
      notifyError(res.message)
      return
    }
    if (res.data) {
      item.icon = res.data.icon
    }
  }
  
  if (item.webUrl) {
    await sys.fs.writeFile(
      `${sys._options.userLocation}Desktop/${item.name}.url`,
      `link::url::${item.webUrl}::${item.icon}`
    );
  }

  notifySuccess(t("store.installSuccess"))
  store.installedList.push(item.name);
  await store.checkProgress()
  setCache();
}

async function uninstall(item: any) {
  if (item.webUrl) {
    await sys.fs.unlink(`${sys._options.userLocation}Desktop/${item.name}.url`);
  } 
  
  delete store.installedList[store.installedList.indexOf(item.name)];
  setCache();
  if (item.needInstall) {
    item.progress = 0
    const completion = await fetch(store.apiUrl + '/store/uninstall?name=' + item.name)
    if (!completion.ok) {
      notifyError(t("store.hasSameName"))
      return
    }
    const res = await completion.json()
    if (res.code && res.code < 0) {
      notifyError(res.message)
      return
    }
  }
  notifySuccess(t("store.uninstallSuccess"))
}
async function download(item: any) {
  //console.log(item)
  if (item.progress) return;
  const completion = await fetch(store.apiUrl + '/store/download?url=' + item.url)
  if (!completion.ok) {
    notifyError(t("store.downloadError"))
  }
  //console.log(completion)
  const reader: any = completion.body?.getReader();
  if (!reader) {
    notifyError(t("store.cantStream"));
  }
  while (true) {
    const { done, value } = await reader?.read();
    if (done) {
      break;
    }
    // console.log(value)
    const json = await new TextDecoder().decode(value);
    //console.log(json)
    const res = parseJson(json)
    //console.log(res)
    if (res) {
      if (res.progress) {
        item.progress = res.progress
      }
      if (res.done) {
        notifySuccess(t("store.downloadSuccess"))
        item.progress = 0
        break;
      }
    }
  }
}
async function pauseApp(item: any) {
  const res: any = await fetch(store.apiUrl + '/store/stop/' + item.name)
  if (!res.ok) {
    const msg = await res.text()
    notifyError(msg)
    return
  }
  setTimeout(async () => {
    await store.checkProgress()
  }, 1000)

}
async function restartApp(item: any) {
  await fetch(store.apiUrl + '/store/restart/' + item.name)
  setTimeout(async () => {
    await store.checkProgress()
  }, 1000)
}
async function startApp(item: any) {
  await fetch(store.apiUrl + '/store/start/' + item.name)
  setTimeout(async () => {
    await store.checkProgress()
  }, 1000)
}
</script>
<template>
  <div class="outer">
    <div class="main">
      <div class="left">
        <div class="left-icon" v-for="(item, key) in store.categoryList" @click="store.changeCate(key, item)">
          <div class="icon-derc" v-if="key == store.currentCateId"></div>
          <el-tooltip class="box-item" effect="dark" :content="t('store.' + item)" placement="right">
            <el-icon size="22">
              <component :is="store.categoryIcon[key]" />
            </el-icon>
          </el-tooltip>
        </div>
      </div>
      <div class="store">
        <div v-if="store.isready" class="store-top">
          <div class="right-main">
            <div class="main-title">
              <span class="sub-title">{{ store.currentTitle }} </span>
            </div>
            <div class="main-app">
              <div v-for="item in store.storeList" v-if="store.currentCate != 'add'" class="store-item" :key="item.name">
                <AppItem :item="item" :installed-list="store.installedList" :install="install" :uninstall="uninstall"
                  :pause="pauseApp" :start="startApp" :restart="restartApp" />
              </div>
              <AddApp v-else :install="install" />

            </div>
          </div>
        </div>
        <div v-else class="store-noready">
          <div id="wait">
            <div class="waitd" id="wait1"></div>
            <div class="waitd" id="wait2"></div>
            <div class="waitd" id="wait3"></div>
            <div class="waitd" id="wait4"></div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<style scoped>
@import "@/assets/store.scss";
</style>
