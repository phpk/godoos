<script setup lang="ts">
import { useModelStore } from "@/stores/model";
import { errMsg } from "@/utils/msg";
import { ref, toRaw, computed } from "vue";
import { t } from "@/i18n/index";
const modelStore = useModelStore();

const formInit = {
  labelId: "",
  ip: "",
  info: {
    model: "",
    url: "",
    from: "ollama",
    file_name: "",
    context_length: "",
    engine: "ollama",
    template: "",
    parameters: "",
    quant: "q4_K_M",
    pb: "",
  },
  type: "local",
};
const engineType = [
  {
    name: "local",
    label: "本地引擎",
  },
  {
    name: "net",
    label: "网络引擎",
  },
]
const formData = ref(formInit);
const fromSource = computed(() => {
  if (formData.value.info.engine == "ollama") {
    return [
      {
        label: "ollama.com",
        value: "ollama",
      },
      {
        label: t("model.network"),
        value: "network",
      },
      {
        label: t("model.local"),
        value: "local",
      },
    ]
  } else {
    return [
      {
        label: t("model.network"),
        value: "network",
      },
      {
        label: t("model.local"),
        value: "local",
      },
    ]
  }
});
function setFrom(val: string) {
  if (val == "ollama") {
    formData.value.info.from = "ollama"
  } else {
    formData.value.info.from = "network"
  }
}
function setType(val: string) {
  if (val == "local") {
    formData.value.info.from = "ollama"
    formData.value.info.engine = "ollama"
  } else {
    formData.value.info.from = "network"
    formData.value.info.engine = "openai"
  }
}
const emit = defineEmits(["closeFn", "saveFn"]);
const localModels: any = ref([]);
async function getLocalModel() {
  const ipv4Pattern = /^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;
  console.log(formData.value.ip);
  if (!ipv4Pattern.test(formData.value.ip)) {
    errMsg(t('model.invalidIp'));
    return;
  }
  const url = `http://${formData.value.ip}:11434/api/tags`;
  try {
    const res = await fetch(url);
    if (!res.ok) {
      errMsg(t('model.fetchFailed'));
      return;
    }
    const data = await res.json();
    if (data && data.length > 0) {
      localModels.value = data;
    }
  } catch (e) {
    console.log(e);
    return;
  }
}
function setLocalInfo() {
  let modelData: any = localModels.value.find((item: any) => {
    return item.model === formData.value.info.model;
  });
  if (!modelData) {
    errMsg(t('model.invalidModel'));
    return;
  }
  modelData = toRaw(modelData);
  const urls: any = [];
  const url = `http://${formData.value.ip}:56780/ai/server?path=`;
  modelData.info.path.forEach((item: any) => {
    urls.push(url + item);
  });
  //formData.value.info.url = urls;
  modelData.info.url = urls.join("\n");
  modelData.info.engine = "ollama";
  modelData.info.from = "local";
  modelData.info.model = formData.value.info.model;
  formData.value.info = modelData.info;
}
async function download() {
  const saveData: any = toRaw(formData.value);
  // console.log(saveData)
  // return
  if (saveData.labelId == "") {
    errMsg(t('model.selectLabel'));
    return;
  }
  if(saveData.info.model == ""){
    errMsg(t('model.labelNameEmpty'));
    return;
  }
  if(saveData.type == "net"){
    saveData.info.url = [];
    saveData.info.context_length = 1024;
    emit("saveFn", saveData);
    return;
  }
  if (saveData.info.from == "ollama") {
    if (saveData.info.model.indexOf(":") === -1) {
      saveData.info.model = saveData.info.model + ":latest";
    }
    if (saveData.info.url == "") {
      saveData.info.url = []
    }
    saveData.info.context_length = 1024
  }

  if (saveData.info.from == "local") {
    if (!saveData.info.url || saveData.info.url.length == 0) {
      errMsg(t('model.invalidModel'));
      return;
    }
  }
  if (saveData.info.from == "network") {
    if (isNaN(saveData.context_length) || saveData.info.context_length < 1) {
      errMsg(t('model.invalidContextLength'));
      return;
    }
    saveData.info.context_length = saveData.info.context_length * 1;

    if (saveData.info.url == "") {
      errMsg(t('model.invalidModelUrl'));
      return;
    }

    if (saveData.info.url != "" && typeof saveData.url === "string") {
      saveData.info.url = saveData.info.url.split("\n");
    } else {
      saveData.info.url = [];
    }
    if (saveData.engine == "ollama") {
      saveData.type = 'local'
      saveData.info.params = {
        top_p: 0.95,
        stream: true,
        num_keep: 5,
        num_predict: 1,
        top_k: 40,
        temperature: 0.7,

      };
      if (saveData.info.parameters != "" && typeof saveData.info.parameters === "string") {
        saveData.info.parameters = saveData.info.parameters.split("\n");
      } else {
        saveData.info.parameters = [];
      }
      // saveData.info = {
      //   quant: saveData.quant,
      //   context_length: saveData.context_length,
      //   template: saveData.template,
      //   parameters: saveData.parameters,
      //   pb: saveData.pb.toUpperCase(),
      // };
      const lowerName = saveData.info.pb.replace("B", "") * 1;
      if (lowerName < 3) {
        saveData.info.cpu = "8GB";
        saveData.info.gpu = "6GB";
      }
      else if (lowerName < 9) {
        saveData.info.cpu = "16GB";
        saveData.info.gpu = "8GB";
      } else {
        saveData.info.cpu = "32GB";
        saveData.info.gpu = "12GB";
      }
      if (saveData.info.model.indexOf(":") === -1) {
        saveData.info.model = saveData.info.model + ":latest";
      }
    }
  }
  //console.log(saveData)
  emit("saveFn", saveData);
}
</script>
<template>
  <el-form ref="form" :model="formData" label-width="150px" style="margin-top: 15px">
    <el-form-item :label="t('model.selectLabel')">
      <el-select v-model="formData.labelId">
        <el-option v-for="(item, key) in modelStore.labelList" :key="key" :label="item.name" :value="item.id" />
      </el-select>
    </el-form-item>
    <el-form-item label="引擎类型">
      <el-select v-model="formData.type" placeholder="选择引擎类型" @change="setType">
        <el-option v-for="item, key in engineType" :key="key" :label="item.label" :value="item.name" />
      </el-select>
    </el-form-item>
    <el-form-item :label="t('model.selectEngine')">
      <el-select v-model="formData.info.engine" :placeholder="t('model.selectEngine')" @change="setFrom">
        <el-option v-for="item, key in modelStore.aiEngine" :key="key" :label="item.name" :value="item.cpp" />
      </el-select>
    </el-form-item>
    <el-form-item :label="t('model.selectSource')"  v-if="formData.type == 'local'">
      <el-select v-model="formData.info.from">
        <el-option v-for="(item, key) in fromSource" :key="key" :label="item.label" :value="item.value" />
      </el-select>
    </el-form-item>
    <el-form-item :label="t('model.modelName')" v-if="formData.info.from !== 'local'">
      <el-input v-model="formData.info.model" prefix-icon="House" clearable
        :placeholder="t('model.enterModelName')"></el-input>
    </el-form-item>

    <template v-if="formData.info.from === 'local' && formData.type == 'local'">
      <el-form-item :label="t('model.oppositeIpAddress')">
        <el-input v-model="formData.ip" prefix-icon="Key" clearable placeholder="192.168.1.66"
          @blur="getLocalModel"></el-input>
      </el-form-item>
      <el-form-item :label="t('model.selectModel')" v-if="localModels.length > 0">
        <el-select v-model="formData.info.model" @change="setLocalInfo">
          <el-option v-for="(item, key) in localModels" :key="key" :label="item.model" :value="item.model" />
        </el-select>
      </el-form-item>
    </template>
    <template v-if="formData.info.from === 'network' && formData.type == 'local'">
      <el-form-item :label="t('model.modelUrl')">
        <el-input type="textarea" :row="3" v-model="formData.info.url"
          :placeholder="t('model.enterModelUrl')"></el-input>
      </el-form-item>
      <!-- <el-form-item :label="t('model.selectEngine')">
        <el-select v-model="formData.info.engine">
          <el-option v-for="(item, key) in modelStore.modelEngines" :key="key" :label="item.name" :value="item.name" />
        </el-select>
      </el-form-item> -->
      <template v-if="formData.info.engine === 'ollama' && formData.info.from === 'network'">
        <el-form-item :label="t('model.template')">
          <el-input type="textarea" :row="3" v-model="formData.info.template"></el-input>
        </el-form-item>

        <el-form-item :label="t('model.contextLength')">
          <el-input type="number" v-model="formData.info.context_length" prefix-icon="Key" clearable
            :placeholder="t('model.enterContextLength')"></el-input>
        </el-form-item>
        <el-form-item :label="t('model.parameterSettings')">
          <el-input type="textarea" :row="3" :placeholder="t('model.onePerLine')"
            v-model="formData.info.parameters"></el-input>
        </el-form-item>
        <el-form-item :label="t('model.parameterSize')">
          <el-input type="number" v-model="formData.info.pb" prefix-icon="Key" clearable
            :placeholder="t('model.enterParameterSize')"></el-input>
        </el-form-item>
        <el-form-item :label="t('model.selectQuantization')">
          <el-select v-model="formData.info.quant">
            <el-option v-for="(item, key) in modelStore.llamaQuant" :key="key" :label="item" :value="item" />
          </el-select>
        </el-form-item>
      </template>
    </template>
    <el-form-item>
      <el-button type="primary" icon="Download" @click="download">{{ t('common.confim') }}</el-button>
    </el-form-item>
  </el-form>
</template>
