<script setup lang="ts">
import { t } from "@/i18n/index";
import { notifySuccess, notifyError } from "@/util/msg";
import {
  getSystemConfig,
  setSystemConfig,
} from "@/system/config";
import { OpenDirDialog } from "@/util/goutil";
const config: any = ref({});
const hoverTxt = {
  dataDir: t('aisetting.tips_dataDir'),
  apiUrl: t('aisetting.tips_apiUrl'),
};
onMounted(() => {
    config.value = getSystemConfig();
});
async function changeDir() {
  const path: any = await OpenDirDialog();
  //console.log(path)
  config.value.dataDir = path;
}
const saveConfig = async () => {
  try {
    await fetch(config.value.apiUrl, {
      method: "GET",
      mode: "no-cors",
    });
  } catch (error) {
    notifyError(t('common.urlError'));
    return;
  }
  let postData: any = []
  if (config.value.dataDir.trim() != "") {
    postData.push({
      name: "aiDir",
      value: config.value.dataDir.trim(),
    })
  }
  if (config.value.ollamaUrl.trim() != "") {
    postData.push({
      name: "ollamaUrl",
      value: config.value.ollamaUrl.trim(),
    })
  }
  if (config.value.openaiUrl.trim() != "") {
    postData.push({
      name: "openaiUrl",
      value: config.value.openaiUrl.trim(),
    })
  }
  if (postData.length > 0) {
    const postDatas = {
      method: "POST",
      body: JSON.stringify(postData),
    };
    const res: any = await fetch(config.value.apiUrl + "/system/setting", postDatas);
    //console.log(res)
    if (!res || !res.ok) {
      notifyError(t('common.saveError'));
      return;
    }
    const ret = await res.json();
    if (ret && ret.code != 0) {
      notifyError(ret.message);
      return;
    }
  }

  setSystemConfig(config.value);
  notifySuccess(t('common.saveSuccess'));
};
</script>
<template>
  <el-scrollbar class="scrollbarSettingHeight">
          <el-form label-width="150px" style="padding: 0 30px 50px 0">
            <el-form-item :label="t('aisetting.dataDir')">
              <div class="slider-container">
                <el-input v-model="config.dataDir" :placeholder="t('aisetting.localDirHolder')" prefix-icon="Menu"
                  @click="changeDir()" clearable></el-input>
                <el-popover placement="left" :width="400" trigger="click">
                  <template #reference>
                    <el-icon :size="22">
                      <InfoFilled />
                    </el-icon>
                  </template>
                  <template #default>
                    <div v-html="hoverTxt.dataDir"></div>
                  </template>
                </el-popover>
              </div>
            </el-form-item>
            <el-form-item :label="t('aisetting.serverUrl')">
              <div class="slider-container">
                <el-input v-model="config.aiUrl" :placeholder="t('aisetting.serverUrl')" prefix-icon="Notification"
                  clearable></el-input>
                <el-popover placement="left" :width="400" trigger="click">
                  <template #reference>
                    <el-icon :size="22">
                      <InfoFilled />
                    </el-icon>
                  </template>
                  <template #default>
                    <div v-html="hoverTxt.apiUrl"></div>
                  </template>
                </el-popover>
              </div>
            </el-form-item>
            <el-form-item :label="t('aisetting.ollamaUrl')">
              <div class="slider-container">
                <el-input v-model="config.ollamaUrl" :placeholder="t('aisetting.ollamaUrl')" prefix-icon="Notification"
                  clearable></el-input>
              </div>
            </el-form-item>
            <el-form-item label="OpenAI URL">
              <div class="slider-container">
                <el-input v-model="config.openaiUrl" placeholder="OpenAI URL" prefix-icon="Notification"
                  clearable></el-input>
              </div>
            </el-form-item>

            <el-form-item>
              <el-button @click="saveConfig" type="info" plain>
                {{ t("common.confim") }}
              </el-button>
            </el-form-item>
          </el-form>
        </el-scrollbar>
</template>