<script setup lang="ts">
import { onMounted, ref, toRaw } from "vue";
import { t, getLang } from "@/i18n/index.ts";
import { useModelStore } from "@/stores/model.ts";
import { notifySuccess, notifyError } from "@/util/msg.ts";
import { Vue3Lottie } from "vue3-lottie";
import { Search } from "@element-plus/icons-vue";
import {
  getSystemConfig
} from "@/system/config";
const currentLang = getLang()
const modelStore = useModelStore();
const config = getSystemConfig();
const downLeft = ref(false);
const downAdd = ref(false);
const labelEditor = ref(false);
const labelId = ref(0);
const searchKey = ref("");
const currentCate = ref("all");

const showDetail = ref(false);
const detailModel = ref("")
let downloadAbort: any = {};
onMounted(async () => {
  await modelStore.getList();
});

async function showCate(name: any) {
  await modelStore.getLabelCate(name.paneName);
}
async function showSearch() {
  await modelStore.getLabelSearch(searchKey.value);
}
function downAddUpdate(val: any) {
  downAdd.value = val;
}
async function downLabel(modelData: any, labelData: any) {
  labelData = toRaw(labelData);
  modelData = toRaw(modelData);
  //console.log(modelData, labelData)
  if(modelData.model){
    modelData.info.model = modelData.model
  }
  const saveData = {
    model: modelData.info.model,
    label: labelData.name,
    action: labelData.action,
    engine: modelData.info.engine,
    url: modelData.info.url ?? [],
    from: modelData.info.from,
    type: modelData.type ?? "",
    file_name: modelData.info.file_name ?? "",
    //options: modelData.options ?? {},
    params: modelData.params ?? {},
    info: modelData.info ?? {},
  };
  await download(saveData);
}
async function saveBox(modelData: any) {
  const labelData = modelStore.labelList.find((d: any) => d.id == modelData.labelId);
  if (!labelData) {
    notifyError(t('model.chooseLabel'));
    return;
  }
  //console.log(modelData)
  downLabel(modelData, labelData);
}
async function download(saveData: any) {
  saveData = toRaw(saveData);
  saveData.info = toRaw(saveData.info);
  //saveData.url = toRaw(saveData.url);
  saveData.params = toRaw(saveData.params);
  downAdd.value = false;
  downLeft.value = true;
  const has = modelStore.checkDownload(saveData.model);
  if (has) {
    notifyError(t('model.labelDown'));
    return;
  }
  //console.log(saveData);
  const downUrl = config.aiUrl + "/ai/download";

  try {
    const completion = await fetch(downUrl, {
      method: "POST",
      body: JSON.stringify(saveData),
    });
    downloadAbort[saveData.model] = false;
    if (!completion.ok) {
      const errorData = await completion.json();
      notifyError(errorData.message);
      return;
    }

    saveData.status = "loading";
    saveData.progress = 0;
    modelStore.addDownload(saveData);
    await handleDown(saveData, completion);
  } catch (error: any) {
    notifyError(error.message);
  }
}
function cancelDownload(model: string) {
  modelStore.deleteDownload(model);
  if (!downloadAbort[model]) {
    downloadAbort[model] = true;
    notifySuccess(t('model.downChanel'));
  }
}

async function handleDown(modelData: any, completion: any) {
  const reader: any = completion.body?.getReader();
  if (!reader) {
    notifyError(t("common.cantStream"));
  }

  while (true) {
    try {
      const { done, value } = await reader.read();
      if (done) {
        //console.log("has done!");
        reader.releaseLock();
        break;
      }
      if (downloadAbort[modelData.model]) {
        break;
      }
      const rawjson = new TextDecoder().decode(value);
      //console.log(rawjson);
      const msg = modelStore.parseMsg(rawjson);
      //console.log(msg)
      if (msg.message && msg.code) {
        notifyError(msg.message);
        break;
      }
      if (msg.status == "") {
        continue;
      }
      modelData.status = msg.status;

      if (msg.total && msg.completed && msg.total > 0) {
        modelData.isLoading = 1;
        modelData.progress = Math.ceil((msg.completed / msg.total) * 100);
        if (modelData.progress == 100 || msg.total == msg.completed) {
          msg.status = "success"
        }
      } else {
        modelData.progress = 0;
      }
      await modelStore.updateDownload(modelData);
    } catch (error) {
      console.error("An error occurred:", error);
      break;
    }
  }
}

async function deleteModel(modelData: any) {
  modelData = toRaw(modelData);
  //console.log(modelData)
  try {
    await modelStore.deleteModelList(modelData);
    notifySuccess(t('prompt.delSuccess'));
  } catch (error: any) {
    //console.log(error);
    notifyError(error.message);
  }
}
function labelShow(val: any) {
  labelId.value = val;
  labelEditor.value = true;
}
function closeLabel() {
  labelId.value = 0;
  labelEditor.value = false;
}
async function refreshList() {
  modelStore.labelList = await modelStore.getLabelList();
}
async function delLabel(id: number) {
  await modelStore.delLabel(id);
  notifySuccess("success!");
}
function getModelStatus(model: string) {
  let name = t('model.noDown');
  if (modelStore.modelList.find((item: any) => item.model === model)) {
    name = t('model.hasDown');
  }
  if (modelStore.downList.find((item: any) => item.model === model)) {
    name = t('model.downloading');
  }
  return name;
}
function showModel(model: string) {
  detailModel.value = model;
  showDetail.value = true;
}
async function refreshOllama() {
  try {
    await modelStore.refreshOllama();
    notifySuccess(t('model.refreshSuccess'));
  } catch (error) {
    notifyError(t('model.refreshFail'));
  }

}
</script>
<template>
  <el-dialog v-model="showDetail" width="600" append-to-body>
    <DownModelInfo :model="detailModel" />
  </el-dialog>
  <div class="app-container">
    <el-drawer v-model="downLeft" direction="ltr" :show-close="false" :with-header="false" :size="300">
      <div>
        <el-tag size="large" style="margin-bottom: 10px">{{ t('model.downloading') }}</el-tag>
        <div class="pa-2">
          <Vue3Lottie animationLink="/bot/search.json" :height="200" :width="200"
            v-if="modelStore.downList.length < 1" />
          <el-space direction="vertical" v-else>
            <el-card v-for="(val, key) in modelStore.downList" :key="key" class="box-card" style="width: 250px">
              <div class="card-header">
                <span>{{ val.model }}</span>
              </div>
              <div class="text item" v-if="val.progress && val.isLoading > 0">
                <el-progress :text-inside="true" :stroke-width="15" :percentage="val.progress" />
              </div>
              <div class="drawer-model-actions" style="margin-top: 10px">
                <el-tag size="small" v-if="val.isLoading > 0">{{ val.status }}</el-tag>
                <el-icon :size="18" color="red" @click="cancelDownload(val.model)">
                  <Delete />
                </el-icon>
                <el-icon :size="18" color="blue" v-if="val.isLoading < 1 && val.status != 'success'"
                  @click="download(toRaw(val))">
                  <VideoPlay />
                </el-icon>
              </div>
            </el-card>
          </el-space>
        </div>
        <el-tag size="large" style="margin: 10px auto">{{ t('model.hasDown') }}</el-tag>
        <div class="pa-2">
          <div class="list-item" v-for="(item, index) in modelStore.modelList" :key="index">
            <div class="list-title" @click="showModel(item.model)">
              {{ item.model }}
            </div>
            <el-button class="delete-btn" icon="Delete" size="small" @click.stop="deleteModel(item)" circle></el-button>
          </div>
        </div>
      </div>
    </el-drawer>
    <el-dialog v-model="labelEditor" width="600" :title="t('model.modelLabel')">
      <down-labeleditor @closeFn="closeLabel" @refreshFn="refreshList" :labelId="labelId" />
    </el-dialog>
    <el-dialog v-model="downAdd" width="600" :title="t('model.modelDown')">
      <down-addbox @closeFn="downAddUpdate" @saveFn="saveBox" />
    </el-dialog>
    <el-page-header icon="null">
      <template #title>
        <div></div>
      </template>
      <template #content>
        <el-button @click.stop="downLeft = !downLeft" icon="Menu" circle />
        <el-button @click.stop="downAdd = true" icon="Plus" circle />
        <el-button @click.stop="labelShow(0)" icon="CollectionTag" circle />
        <el-button @click.stop="refreshOllama" icon="RefreshRight" circle />

      </template>
      <template #extra>
        <el-space class="mr-10">
          <el-input :placeholder="t('model.search')" v-model="searchKey" v-on:keydown.enter="showSearch"
            style="width: 200px" :suffix-icon="Search" />
        </el-space>
      </template>
    </el-page-header>

    <div class="flex-fill ml-10 mr-10">
      <el-tabs v-model="currentCate" @tab-click="showCate">
        <el-tab-pane :label="t('model.all')" name="all" />
        <el-tab-pane :label="t('model.' + item)" :name="item" v-for="(item, key) in modelStore.cateList" :key="key" />
      </el-tabs>
    </div>

    <el-scrollbar class="scrollbarHeightList">
      <div class="model-list">
        <div v-for="item in modelStore.labelList" :key="item.name" class="model-item flex align-center pa-5">
          <div class="flex-fill mx-5">
            <div class="font-weight-bold">
              {{ item.name }}
            </div>
            <div class="desc">
              {{ currentLang == "zh-cn" ? item.zhdesc : item.endesc }}
            </div>
            <div></div>
          </div>
          <div class="drawer-model-actions">
            <el-popover placement="left" :width="300" trigger="click">
              <template #reference>
                <el-button icon="Download" circle />
              </template>
              <template #default>
                <div v-for="(el, index) in item.models" :key="index" :value="el.model" @click="downLabel(el, item)"
                  class="list-column">
                  <div class="list-column-title">
                    {{ el.model }}
                    <el-tag size="small" type="info">{{
                      getModelStatus(el.model)
                    }}</el-tag>
                  </div>
                  <div class="list-footer">
                    <el-tag size="small" type="primary">size:{{ el.info.size }}</el-tag>
                    <el-tag size="small" type="success">cpu:{{ el.info.cpu }}</el-tag>
                    <el-tag size="small" type="danger">gpu:{{ el.info.gpu }}</el-tag>
                  </div>
                </div>
              </template>
            </el-popover>
            <el-button icon="Edit" circle @click="labelShow(item.id)" />
            <el-button @click.stop="delLabel(item.id)" icon="Delete" v-if="item.models.length === 0" circle />
          </div>
        </div>
      </div>
    </el-scrollbar>
  </div>
</template>

<style scoped lang="scss">
@import "@/assets/list.scss";
@import "@/assets/left.scss";
</style>
