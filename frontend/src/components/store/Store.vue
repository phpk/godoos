<script setup lang="ts">
import { inject, onMounted, ref } from "vue";
import { System } from "@/system";
import {notifySuccess, notifyError} from "@/util/msg";
import { t } from "@/i18n";
//import storeList from "@/assets/store.json";
import { getSystemKey, setSystemKey } from "@/system/config";
const sys: any = inject<System>("system");
const currentCateId = ref(0)
const currentTitle = ref(t("store.hots"))
const currentCate = ref('hots')
const categoryList = ['hots', 'work', 'development', 'games', 'education', 'news', 'shopping', 'social', 'utilities', 'others', 'add']
const categoryIcon = ['HomeFilled', 'Odometer', 'Postcard', 'TrendCharts', 'School', 'HelpFilled', 'ShoppingCart', 'ChatLineRound', 'MessageBox', 'Ticket', 'CirclePlusFilled']
const isready = ref(false);
const installed = getSystemKey("intstalledPlugins");
const apiUrl = getSystemKey("apiUrl");
const installedList: any = ref(installed);
const storeList: any = ref([])
onMounted(async() => {
  await getList();
  
});
async function getList(){
  if(currentCate.value == 'add')return;
  const storeUrl = apiUrl + '/system/storeList?cate=' + currentCate.value
  fetch(storeUrl).then(res => {
    res.json().then(data => {
      storeList.value = data
    })
  }).catch(() => {
    notifyError(t("store.errorList"))
  })
  isready.value = true;
}
async function changeCate(index: number, item: string) {
  currentCateId.value = index
  currentCate.value = item
  currentTitle.value = t("store." + item)
  await getList()
}
function setCache() {
  setSystemKey("intstalledPlugins", installedList.value);
  setTimeout(() => {
    sys.refershAppList();
  }, 1000);
}
function install(item: any) {
  //console.log(item)
  if(item.url){
    sys.fs.writeFile(
      `${sys._options.userLocation}Desktop/${item.name}.url`,
      `link::url::${item.url}::${item.icon}`
    );
  }
  
  notifySuccess(t("install.success"))
  installedList.value.push(item.name);
  setCache();
}

function uninstall(item: any) {
  sys.fs.unlink(`${sys._options.userLocation}Desktop/${item.name}.url`);
  notifySuccess(t("uninstall.success"))
  delete installedList.value[installedList.value.indexOf(item.name)];
  setCache();
}
</script>
<template>
  <div class="outer">
    <div class="main">
      <div class="left">
        <div class="left-icon" v-for="(item, key) in categoryList" @click="changeCate(key, item)">
          <div class="icon-derc" v-if="key == currentCateId"></div>
          <el-tooltip class="box-item" effect="dark" :content="t('store.' + item)" placement="right">
            <el-icon size="22">
              <component :is="categoryIcon[key]" />
            </el-icon>
          </el-tooltip>
        </div>
      </div>
      <div class="store">
        <div v-if="isready" class="store-top">
          <!-- <div class="left-bar"></div> -->
          <div class="right-main">
            <div class="main-title">
              <span class="sub-title">{{ currentTitle }} </span>
            </div>
            <!-- <div class="swiper">
              <div class="swiper-txt">主页</div>
              <div class="swiper-inner">
                <div class="swiper-tab">
                  <img src="/image/store/banner1.jpg" />
                </div>
                <div class="swiper-tab">
                  <img src="/image/store/banner2.jpg" />
                </div>
                <div class="swiper-tab">
                  <img src="/image/store/banner3.jpg" />
                </div>
              </div>
            </div> -->
            <div class="main-app">
                <div v-for="item in storeList"  v-if="currentCate != 'add'" class="store-item" :key="item.name">
                  <AppItem :item="item" :installed-list="installedList" :install="install" :uninstall="uninstall" />
                </div>
                  <AddApp v-else/>

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
