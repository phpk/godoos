<script setup lang="ts">
import { useModelStore } from "@/stores/model";
import { notifyError } from "@/util/msg";
import { ref, toRaw, computed } from "vue";
import { t } from "@/i18n/index";
const modelStore = useModelStore();

const formInit = {
  model: "",
  labelId: "",
  ip: "",
  info: {
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
  type: "",
};
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
const emit = defineEmits(["closeFn", "saveFn"]);
const localModels: any = ref([]);
async function getLocalModel() {
  const ipv4Pattern = /^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;
  console.log(formData.value.ip);
  if (!ipv4Pattern.test(formData.value.ip)) {
    notifyError(t('model.invalidIp'));
    return;
  }
  const url = `http://${formData.value.ip}:56780/ai/tags`;
  try {
    const res = await fetch(url);
    if (!res.ok) {
      notifyError(t('model.fetchFailed'));
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
    return item.model === formData.value.model;
  });
  if (!modelData) {
    notifyError(t('model.invalidModel'));
    return;
  }
  modelData = toRaw(modelData);
  const urls: any = [];
  const url = `http://${formData.value.ip}:56780/ai/server?path=`;
  modelData.info.path.forEach((item: any) => {
    urls.push(url + item);
  });
  //formData.value.info.url = urls;
  modelData.info.url = urls;
  formData.value.info = modelData.info;
  //formData.value.file_name = modelData.file_name;
  //formData.value.engine = modelData.engine;
  // if (modelData.engine == "ollama") {
  //   formData.value.type = "local";
  // }
}
async function download() {
  const saveData: any = toRaw(formData.value);
  // console.log(saveData)
  // return
  if (saveData.labelId == "") {
    notifyError(t('model.selectLabel'));
    return;
  }

  if (saveData.from == "ollama") {
    if (saveData.model == "") {
      notifyError(t('model.labelNameEmpty'));
      return;
    }
    if (saveData.model.indexOf(":") === -1) {
      saveData.model = saveData.model + ":latest";
    }
  }

  if (saveData.from == "local") {
    if (!saveData.url || saveData.url.length == 0) {
      notifyError(t('model.invalidModel'));
      return;
    }
  }
  if (saveData.from == "network") {
    if (isNaN(saveData.context_length) || saveData.context_length < 1) {
      notifyError(t('model.invalidContextLength'));
      return;
    }
    saveData.context_length = saveData.context_length * 1;

    if (saveData.url == "") {
      notifyError(t('model.invalidModelUrl'));
      return;
    }

    if (saveData.url != "" && typeof saveData.url === "string") {
      saveData.url = saveData.url.split("\n");
    } else {
      saveData.url = [];
    }
    if (saveData.engine == "ollama") {
      saveData.type = 'llm'
      saveData.params = {
        top_p: 0.95,
        stream: true,
        num_keep: 5,
        num_predict: 1,
        top_k: 40,
        temperature: 0.7,

      };
      if (saveData.parameters != "" && typeof saveData.parameters === "string") {
        saveData.parameters = saveData.parameters.split("\n");
      } else {
        saveData.parameters = [];
      }
      saveData.info = {
        quant: saveData.quant,
        context_length: saveData.context_length,
        template: saveData.template,
        parameters: saveData.parameters,
        pb: saveData.pb.toUpperCase(),
      };
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
      if (saveData.model.indexOf(":") === -1) {
        saveData.model = saveData.model + ":latest";
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
    <el-form-item :label="t('model.selectEngine')">
      <el-select v-model="formData.info.engine" :placeholder="t('model.selectEngine')" @change="setFrom">
        <el-option v-for="item, key in modelStore.modelEngines" :key="key" :label="item.name" :value="item.name" />
      </el-select>
    </el-form-item>
    <el-form-item :label="t('model.selectSource')">
      <el-select v-model="formData.info.from">
        <el-option v-for="(item, key) in fromSource" :key="key" :label="item.label" :value="item.value" />
      </el-select>
    </el-form-item>
    <el-form-item :label="t('model.modelName')" v-if="formData.info.from !== 'local'">
      <el-input v-model="formData.model" prefix-icon="House" clearable
        :placeholder="t('model.enterModelName')"></el-input>
    </el-form-item>

    <template v-if="formData.info.from === 'local'">
      <el-form-item :label="t('model.oppositeIpAddress')">
        <el-input v-model="formData.ip" prefix-icon="Key" clearable placeholder="192.168.1.66"
          @blur="getLocalModel"></el-input>
      </el-form-item>
      <el-form-item :label="t('model.selectModel')" v-if="localModels.length > 0">
        <el-select v-model="formData.model" @change="setLocalInfo">
          <el-option v-for="(item, key) in localModels" :key="key" :label="item.model" :value="item.model" />
        </el-select>
      </el-form-item>
    </template>
    <template v-if="formData.info.from === 'network'">
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
