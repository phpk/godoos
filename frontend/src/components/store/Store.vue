<script setup lang="ts">
import { onMounted } from "vue";
import { notifySuccess, notifyError } from "@/util/msg";
import { useStoreStore } from "@/stores/store";
import { parseJson } from "@/system/config";
import { t } from "@/i18n";
const store = useStoreStore()

onMounted(async () => {
  await store.getList();
});

async function install(item: any) {
  //console.log(item)
  if (item.isOut) {
    item.progress = 0
    const flag = await store.addOutList(item)
    if (!flag) {
      notifyError(t("store.installError"))
    }
  }
  if ((item.pkg != "" || item.url != "") && !item.isDev) {
    const postData = {
      name: item.name,
      pkg: item.pkg,
      url: item.url,
      needInstall: item.needInstall,
    }
    const url = store.apiUrl + '/store/download'
    const res = await download(item, url, postData)
    //console.log(res)
    item.progress = 0
    if(res.code < 0){
      if(res.data && res.data.dependencies && res.data.dependencies.length > 0){
        const names:any = []
        res.data.dependencies.forEach((element:any) => {
          names.push(element.name)
        });
        notifyError("you need install these plugins first" + names.join(","))
      }else{
        notifyError(t("store.installError"))
      }
      return
    }
    item = res.data
    item.progress = 0
    
    await store.addDesktop(item);
    await store.checkProgress()
    notifySuccess(t("store.installSuccess"))
  }
}

async function uninstall(item: any) {
  await store.removeDesktop(item);
  if (item.needInstall) {
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
async function download(item: any, url: string, postData: any) {
  //console.log(item)
  if (item.progress && item.progress > 0) return;
  const completion = await fetch(url, {
    method: 'POST',
    body: JSON.stringify(postData),
  })
  if (!completion.ok) {
    notifyError(t("store.downloadError"))
  }
  //console.log(completion)
  const reader: any = completion.body?.getReader();
  if (!reader) {
    notifyError(t("store.cantStream"));
  }
  let res:any
  while (true) {
    const { done, value } = await reader?.read();
    if (done) {
      break;
    }
    // console.log(value)
    const json = await new TextDecoder().decode(value);
    res = parseJson(json)
    //console.log(res)
    if (res) {
      if (res.progress) {
        item.progress = res.progress
      }
    }
  }
  return res
}
async function pauseApp(item: any) {
  const res: any = await fetch(store.apiUrl + '/store/stop/' + item.name)
  if (!res.ok) {
    const msg = await res.text()
    notifyError(msg)
    //return
  }
  setTimeout(async () => {
    await store.checkProgress()
  }, 3000)

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
              <div v-for="item in store.storeList" v-if="store.currentCate != 'add'" class="store-item"
                :key="item.name">
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
