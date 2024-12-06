<script setup lang="ts">
import { t } from "@/i18n/index";
import { notifySuccess, notifyError } from "@/util/msg";
import {
  getSystemConfig,
  setSystemConfig,
} from "@/system/config";
import { OpenDirDialog, checkUrl } from "@/util/goutil";
const activeNames = ref([])
const hoverTxt = {
  dataDir: t('aisetting.tips_dataDir'),
  apiUrl: t('aisetting.tips_apiUrl'),
};
const formData: any = ref({
  aiDir: "",
  aiUrl: "",
  ollamaUrl: "",
  ollamaDir: "",
  openaiUrl: "",
  openaiSecret: "",
  giteeSecret: "",
  cloudflareUserId: "",
  cloudflareSecret: "",
  deepseekSecret: "",
  bigmodelSecret: "",
  volcesSecret:"",
  alibabaSecret: "",
  groqSecret: "",
  mistralSecret: "",
  anthropicSecret: "",
  llamafamilySecret: "",
  siliconflowSecret: ""
})
onMounted(() => {
  const config = getSystemConfig();
  for (const key in config) {
    if (formData.value.hasOwnProperty(key)) {
      formData.value[key] = config[key];
    }
  }
});
async function changeDir(name: any) {
  const path: any = await OpenDirDialog();
  //console.log(path)
  formData.value[name] = path;
}
const saveConfig = async () => {
  const config = getSystemConfig();
  const saveData = toRaw(formData.value)
  let postData: any = []
  for (const key in saveData) {
    saveData[key] = saveData[key].trim()
    config[key] = saveData[key];
    if (saveData[key] != "") {
      postData.push({
        name: key,
        value: saveData[key],
      })
    }
  }
  if (saveData.aiUrl != "") {
    if (!await checkUrl(saveData.aiUrl)) {
      notifyError("ai服务端地址有误");
      return;
    }
  }
  if (saveData.ollamaUrl != "") {
    if (!await checkUrl(saveData.ollamaUrl)) {
      notifyError("ollama服务端地址有误");
      return;
    }
  }

  if (postData.length > 0) {
    const postDatas = {
      method: "POST",
      body: JSON.stringify(postData),
    };
    const res: any = await fetch(config.apiUrl + "/system/setting", postDatas);
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

  setSystemConfig(config);
  notifySuccess(t('common.saveSuccess'));
};
</script>
<template>
  <el-scrollbar class="scrollbarSettingHeight">
    <el-form label-width="150px" style="padding: 0 30px 50px 0">
      <el-collapse v-model="activeNames">
        <el-collapse-item title="系统配置" name="system">
          <el-form-item :label="t('aisetting.dataDir')">
            <div class="slider-container">
              <el-input v-model="formData.aiDir" placeholder="系统模型下载本地地址" prefix-icon="Menu"
                @click="changeDir('aiDir')" clearable></el-input>
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
          <el-form-item label="AI服务端地址">
            <div class="slider-container">
              <el-input v-model="formData.aiUrl" placeholder="AI服务端地址" prefix-icon="Notification" clearable></el-input>
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
        </el-collapse-item>

        <el-collapse-item title="Ollama配置" name="ollama">
          <el-form-item :label="t('aisetting.ollamaUrl')">
            <div class="slider-container">
              <el-input v-model="formData.ollamaUrl" placeholder="ollama模型访问地址，为空则使用http://localhost:11434"
                prefix-icon="Notification" clearable></el-input>
            </div>
          </el-form-item>
          <!-- <el-form-item label="模型存储地址">
            <div class="slider-container">
              <el-input v-model="formData.ollamaDir" placeholder="ollama模型存储地址，设置后ollama需重启，为空则使用默认地址"
                prefix-icon="Menu" @click="changeDir('ollamaDir')" clearable></el-input>
            </div>
          </el-form-item> -->
        </el-collapse-item>

        <el-collapse-item title="OpenAI配置" name="openai">
          <el-form-item label="OpenAI地址">
            <div class="slider-container">
              <el-input v-model="formData.openaiUrl" placeholder="OpenAI地址，可设置为代理地址" prefix-icon="Notification"
                clearable></el-input>
            </div>
          </el-form-item>
          <el-form-item label="OpenAI私钥">
            <div class="slider-container">
              <el-input v-model="formData.openaiSecret" placeholder="OpenAI私钥" prefix-icon="Lock" clearable></el-input>
            </div>
          </el-form-item>
        </el-collapse-item>

        <el-collapse-item title="GiteeAI配置" name="gitee">
          <el-form-item label="GiteeAI私钥">
            <div class="slider-container">
              <el-input v-model="formData.giteeSecret" placeholder="GiteeAI私钥" prefix-icon="Lock" clearable></el-input>
            </div>
          </el-form-item>
        </el-collapse-item>

        <el-collapse-item title="Cloudflare Workers AI配置" name="cloudflare">
          <el-form-item label="Cloudflare用户ID">
            <div class="slider-container">
              <el-input v-model="formData.cloudflareUserId" placeholder="Cloudflare用户ID" prefix-icon="Notification"
                clearable></el-input>
            </div>
          </el-form-item>
          <el-form-item label="Cloudflare私钥">
            <div class="slider-container">
              <el-input v-model="formData.cloudflareSecret" placeholder="Cloudflare私钥" prefix-icon="Lock"
                clearable></el-input>
            </div>
          </el-form-item>
        </el-collapse-item>

        <el-collapse-item title="DeepSeek配置" name="deepseek">
          <el-form-item label="DeepSeek私钥">
            <div class="slider-container">
              <el-input v-model="formData.deepseekSecret" placeholder="DeepSeek私钥" prefix-icon="Lock" clearable></el-input>
            </div>
          </el-form-item>
        </el-collapse-item>

        <el-collapse-item title="智谱清言语BigModel配置" name="bigmodel">
          <el-form-item label="智谱清言语私钥">
            <div class="slider-container">
              <el-input v-model="formData.bigmodelSecret" placeholder="智谱清言语私钥" prefix-icon="Lock" clearable></el-input>
            </div>
          </el-form-item>
        </el-collapse-item>

        <el-collapse-item title="火山方舟配置" name="volces">
          <el-form-item label="火山方舟私钥">
            <div class="slider-container">
              <el-input v-model="formData.volcesSecret" placeholder="火山方舟私钥" prefix-icon="Lock" clearable></el-input>
            </div>
          </el-form-item>
        </el-collapse-item>

        <el-collapse-item title="阿里通义DashScope配置" name="alibaba">
          <el-form-item label="阿里通义私钥">
            <div class="slider-container">
              <el-input v-model="formData.alibabaSecret" placeholder="阿里通义DashScope私钥" prefix-icon="Lock" clearable></el-input>
            </div>
          </el-form-item>
        </el-collapse-item>

        <el-collapse-item title="Groq配置" name="groq">
          <el-form-item label="Groq私钥">
            <div class="slider-container">
              <el-input v-model="formData.groqSecret" placeholder="Groq私钥" prefix-icon="Lock" clearable></el-input>
            </div>
          </el-form-item>
        </el-collapse-item>

        <el-collapse-item title="Mistral配置" name="mistral">
          <el-form-item label="mistral私钥">
            <div class="slider-container">
              <el-input v-model="formData.mistralSecret" placeholder="mistral私钥" prefix-icon="Lock" clearable></el-input>
            </div>
          </el-form-item>
        </el-collapse-item>

        <el-collapse-item title="Anthropic配置" name="anthropic">
          <el-form-item label="anthropic私钥">
            <div class="slider-container">
              <el-input v-model="formData.anthropicSecret" placeholder="anthropic私钥" prefix-icon="Lock" clearable></el-input>
            </div>
          </el-form-item>
        </el-collapse-item>

        <el-collapse-item title="llama.family配置" name="llamafamily">
          <el-form-item label="llamafamily私钥">
            <div class="slider-container">
              <el-input v-model="formData.llamafamilySecret" placeholder="llamafamily私钥" prefix-icon="Lock" clearable></el-input>
            </div>
          </el-form-item>
        </el-collapse-item>

        <el-collapse-item title="硅基流动配置" name="siliconflow">
          <el-form-item label="硅基流动私钥">
            <div class="slider-container">
              <el-input v-model="formData.siliconflowSecret" placeholder="硅基流动siliconflow私钥" prefix-icon="Lock" clearable></el-input>
            </div>
          </el-form-item>
        </el-collapse-item>

      </el-collapse>




      <el-form-item style="margin-top: 10px;">
        <el-button @click="saveConfig" type="info" plain>
          {{ t("common.confim") }}
        </el-button>
      </el-form-item>
    </el-form>
  </el-scrollbar>
</template>