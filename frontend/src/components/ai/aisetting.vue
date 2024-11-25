<script lang="ts" setup>
import { onMounted, ref,watch } from "vue";
import {
  getSystemConfig,
  setSystemConfig,
} from "@/system/config";
import { notifySuccess, notifyError } from "@/util/msg";
import { t } from "@/i18n/index";
import { OpenDirDialog } from "@/util/goutil";
import { useModelStore } from "@/stores/model.ts";
const modelStore = useModelStore();
// 参数	描述	值类型	示例用法
// mirostat	启用Mirostat采样以控制困惑度。默认值：0，其中0=禁用，1=Mirostat，2=Mirostat 2.0	int	mirostat 0
// mirostat_eta	影响算法根据生成文本的反馈进行调整的速度。较低的学习率将导致调整速度变慢，而较高的学习率将使算法响应更迅速。默认值：0.1	float	mirostat_eta 0.1
// mirostat_tau	控制输出的连贯性和多样性之间的平衡。更低的值将产生更集中和连贯的文本。默认值：5.0	float	mirostat_tau 5.0
// num_ctx	设置用于生成下一个令牌的上下文窗口大小。默认值：2048	int	num_ctx 4096
// repeat_last_n	设置模型应回顾多远以防止重复。（默认：64，0=禁用，-1=num_ctx）	int	repeat_last_n 64
// repeat_penalty	设置对重复的惩罚程度。较高的值（例如，1.5）会更强地惩罚重复，而较低的值（例如，0.9）则较为宽容。默认值：1.1	float	repeat_penalty 1.1
// temperature	模型的温度。增加温度会使模型回答更具创造性。默认值：0.8	float	temperature 0.7
// seed	设置用于生成的随机数种子。将其设置为特定数字将使模型在相同的提示下生成相同的文本。默认值：0	int	seed 42
// stop	设置使用的停止序列。当遇到此模式时，LLM将停止生成文本并返回。可以在模型文件中通过指定多个单独的stop参数来设置多个停止模式。	string	stop "AI助手:"
// tfs_z	尾部自由采样用于减少输出中较不可能出现的令牌的影响。较高的值（例如，2.0）将更大程度地减少这种影响，而1.0的值则禁用此设置。默认值：1	float	tfs_z 1
// num_predict	生成文本时预测的最大令牌数量。（默认：128，-1=无限生成，-2=填充上下文）	int	num_predict 42
// top_k	减少生成无意义内容的概率。较高的值（例如，100）将给出更多样化的答案，而较低的值（例如，10）则更为保守。默认值：40	int	top_k 40
// top_p	与top-k一起工作。较高的值（例如，0.95）将产生更多样化的文本，而较低的值（例如，0.5）将生成更集中且保守的文本。默认值：0.9	float	top_p 0.9

const hoverTxt = {
  dataDir: t('aisetting.tips_dataDir'),
  apiUrl: t('aisetting.tips_apiUrl'),
  contextLength: t("setting.tips_contextLength"),
  top_k: t('aisetting.tips_top_k'),
  top_p: t('aisetting.tips_top_p'),
  temperature: t('aisetting.tips_temperature'),
  frequency_penalty: t('aisetting.tips_frequency_penalty'),
  presence_penalty: t('aisetting.tips_presence_penalty'),
  num_predict: t('aisetting.tips_num_predict'),
  num_keep: t('aisetting.tips_num_keep'),
};
// const systemStore = useSystemStore();
const config: any = ref({});
//const chatConfig: any = ref({});
const currentsModel: any = ref({});
import type { TabsPaneContext } from "element-plus";

const activeName = ref("system");
const activeModel = ref("chat");

const handleClick = (tab: TabsPaneContext, event: Event) => {
  console.log(tab, event);
};
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
      name: "dataDir",
      value: config.value.dataDir.trim(),
    })
  }
  if (config.value.ollamaUrl.trim() != "") {
    postData.push({
      name: "ollamaUrl",
      value: config.value.ollamaUrl.trim(),
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

  await changeConfig();
  notifySuccess(t('common.saveSuccess'));
};

const changeConfig = async () => {
  //console.log(v)
  //setSystemKey('llmType',v)
  setSystemConfig(config.value);
  await initConfig();
};
const initConfig = async () => {
  config.value = getSystemConfig();
  //currentsModel.value = getCurrents();
  //chatConfig.value = getChatConfig();
};
// 更新 currentsModel 的函数
function updateCurrentsModel() {
  modelStore.cateList.forEach((item:any) => {
    const currentModel = modelStore.modelList.find((el:any) => el.action === item && el.isdef === 1);
    if (currentModel) {
      currentsModel.value[item] = currentModel.model;
    } else {
      const firstModel = modelStore.getCurrentModelList(item)[0];
      currentsModel.value[item] = firstModel ? firstModel.model : '';
    }
  });
}
onMounted(async () => {
  await initConfig();
  await modelStore.getModelList();
  updateCurrentsModel();

});
watch(modelStore.modelList, () => {
  updateCurrentsModel();
});
async function changeDir() {
  const path: any = await OpenDirDialog();
  //console.log(path)
  config.value.dataDir = path;
}
</script>
<template>
  <div>
    <el-tabs v-model="activeName" class="setting-tabs" style="margin: 12px" @tab-click="handleClick">
      <el-tab-pane :label="t('aisetting.modelSetting')" name="system">
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

            <el-form-item>
              <el-button @click="saveConfig" type="info" plain>
                {{ t("common.confim") }}
              </el-button>
            </el-form-item>
          </el-form>
        </el-scrollbar>
      </el-tab-pane>
      <el-tab-pane :label="t('aisetting.defModel')" name="modelDef">
        <el-scrollbar class="scrollbarSettingHeight">
          <el-form label-width="150px" style="padding: 0 30px 50px 0">
            <el-form-item :label="t('model.' + item)" v-for="(item, index) in modelStore.cateList" :key="index">
              <el-select v-model="currentsModel[item]" @change="(val: any) => modelStore.setCurrentModel(item, val)">
                <el-option v-for="(el, key) in modelStore.getCurrentModelList(item)" :key="key" :label="el.model"
                  :value="el.model" />
              </el-select>
            </el-form-item>
          </el-form>
        </el-scrollbar>
      </el-tab-pane>
      <el-tab-pane :label="t('aisetting.chatSetting')" name="chatSetting">
        <el-tabs tab-position="left" v-model="activeModel" style="height: 490px" class="setting-tabs">
          <el-tab-pane :name="item.key" :label="t('model.' + item.key)" v-for="item in modelStore.chatConfig">
            <el-form label-width="100px" style="width: 500px">
              <el-form-item :label="t('aisetting.contextLength')" v-if="item.contextLength" class="inline-layout">
                <div class="slider-container">
                  <el-slider v-model="item.contextLength" :max="10" :min="1" :step="1" />
                  <el-popover placement="left" :width="400" trigger="click">
                    <template #reference>
                      <el-icon :size="22">
                        <InfoFilled />
                      </el-icon>
                    </template>
                    <template #default>
                      <div v-html="hoverTxt.contextLength"></div>
                    </template>
                  </el-popover>
                </div>
              </el-form-item>

              <el-form-item :label="t('aisetting.num_predict')" class="inline-layout">
                <div class="slider-container">
                  <el-slider v-model="item.num_predict" :max="5000" :min="1" />
                  <el-popover placement="left" :width="400" trigger="click">
                    <template #reference>
                      <el-icon :size="22">
                        <InfoFilled />
                      </el-icon>
                    </template>
                    <template #default>
                      <div v-html="hoverTxt.num_predict"></div>
                    </template>
                  </el-popover>
                </div>
              </el-form-item>
              <el-form-item :label="t('aisetting.num_keep')" class="inline-layout">
                <div class="slider-container">
                  <el-slider v-model="item.num_keep" :max="500" :min="1" />
                  <el-popover placement="left" :width="400" trigger="click">
                    <template #reference>
                      <el-icon :size="22">
                        <InfoFilled />
                      </el-icon>
                    </template>
                    <template #default>
                      <div v-html="hoverTxt.num_keep"></div>
                    </template>
                  </el-popover>
                </div>
              </el-form-item>
              <el-form-item :label="t('aisetting.top_k')" class="inline-layout">
                <div class="slider-container">
                  <el-slider v-model="item.top_k" :max="100" :min="1" />
                  <el-popover placement="left" :width="400" trigger="click">
                    <template #reference>
                      <el-icon :size="22">
                        <InfoFilled />
                      </el-icon>
                    </template>
                    <template #default>
                      <div v-html="hoverTxt.top_k"></div>
                    </template>
                  </el-popover>
                </div>
              </el-form-item>
              <el-form-item :label="t('aisetting.top_p')" class="inline-layout">
                <div class="slider-container">
                  <el-slider v-model="item.top_p" :max="1" :min="0.01" :step="0.01" />
                  <el-popover placement="left" :width="400" trigger="click">
                    <template #reference>
                      <el-icon :size="22">
                        <InfoFilled />
                      </el-icon>
                    </template>
                    <template #default>
                      <div v-html="hoverTxt.top_p"></div>
                    </template>
                  </el-popover>
                </div>
              </el-form-item>
              <el-form-item :label="t('aisetting.temperature')" class="inline-layout">
                <div class="slider-container">
                  <el-slider v-model="item.temperature" :max="0.99" :min="0.01" :step="0.01" />
                  <el-popover placement="left" :width="400" trigger="click">
                    <template #reference>
                      <el-icon :size="22">
                        <InfoFilled />
                      </el-icon>
                    </template>
                    <template #default>
                      <div v-html="hoverTxt.temperature"></div>
                    </template>
                  </el-popover>
                </div>
              </el-form-item>
            </el-form>
          </el-tab-pane>
        </el-tabs>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>
<style>
.setting-tabs>.el-tabs__content {
  padding: 0px;
  color: #6b778c;
  font-size: 32px;
  font-weight: 600;
}

.scrollbarSettingHeight {
  height: 80vh;
  padding-bottom: 30px;
}

.inline-layout {
  display: flex;
  align-items: center;
  /* 如果有必要，可以在这里添加额外的样式来调整表单项的整体表现 */
}

.slider-container {
  display: flex;
  align-items: center;
  width: 100%;
}

.slider-container .el-icon {
  cursor: pointer;
  margin-left: 20px;
}
</style>
