<script setup lang="ts">
import { useKnowStore } from "@/stores/know.ts";
import {getModelList,getPromptList} from "@/stores/config";
import { onMounted,ref } from "vue";
import { t } from "@/i18n/index";
const knowStore = useKnowStore();
const props = defineProps<{
  dataInfo: any
}>()
const modelList:any = ref([])
const knowList:any = ref([])
const promptList:any = ref([])
const emit = defineEmits(['saveFn'])
function changPrompt(promptName:string) {
  const promptData:any = promptList.value.find((item:any) => {
    return item.name == promptName;
  });
  if(promptData) {
    props.dataInfo.prompt = promptData.prompt;
  }
  
}
function saveData(){
  emit('saveFn', props.dataInfo)
}
onMounted(async () => {
  modelList.value = await getModelList('chat')
  knowList.value = await knowStore.getKnowAll()
  promptList.value = await getPromptList('chat')
})
</script>
<template>
    <el-form label-width="150px" style="margin-top: 12px">
      <el-form-item :label="t('common.title')">
        <el-input
          v-model="props.dataInfo.title"
          :placeholder="t('chat.inputTitle')"
          prefix-icon="Notification"
          clearable
          :autofocus="true"
        ></el-input>
      </el-form-item>
      <el-form-item :label="t('chat.model')">
        <el-select v-model="props.dataInfo.model">
          <el-option
            v-for="(item, key) in modelList"
            :key="key"
            :label="item.model"
            :value="item.model"
          />
        </el-select>
      </el-form-item>
      <el-form-item :label="t('chat.refknow')">
        <el-select
          v-model="props.dataInfo.kid"
          :clearable="true"
          :filterable="true"
        >
          <el-option
            v-for="(item, key) in knowList"
            :key="key"
            :label="item.name"
            :value="item.uuid"
          />
        </el-select>
      </el-form-item>
      <el-form-item :label="t('chat.role')">
        <el-select
          v-model="props.dataInfo.promptName"
          @change="changPrompt"
        >
          <el-option
            v-for="(item, key) in promptList"
            :key="key"
            :label="item.name"
            :value="item.name"
          />
        </el-select>
      </el-form-item>
      
      <el-form-item>
        <el-input
          type="textarea"
          v-model="props.dataInfo.prompt"
          :rows="6"
        ></el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="saveData">{{
          t("common.confim")
        }}</el-button>
      </el-form-item>
    </el-form>
</template>