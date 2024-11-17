<script lang="ts" setup>
import { useAssistantStore } from "@/stores/assistant";
import { t } from "@/i18n/index";
import { watchEffect, ref, toRaw } from "vue";
import { notifyError } from '@/util/msg'
const assistantStore = useAssistantStore()
const formData = ref({
  name: "",
  action: "",
  prompt: "",
  isdef: 0,
});
watchEffect(async () => {
  if(assistantStore.editId > 0) {
    const data = await assistantStore.getPromptById(assistantStore.editId);
    //console.log(data)
    formData.value = {
      name: data.name,
      action: data.action,
      prompt: data.prompt,
      isdef: data.isdef.toString(),
    };
    //console.log(formData.value)
  }else{
    formData.value = {
      name: "",
      action: "",
      prompt: "",
      isdef: 0,
    };
  }
});
const getCateList = () => {
  const ret: any = [];
  assistantStore.promptAction.forEach((item) => {
    const cate = {
      title: t("prompt." + item),
      value: item,
    };
    ret.push(cate);
  });
  return ret;
};
async function saveData() {
  const save = toRaw(formData.value)
  if(save.name == "") {
    notifyError("Name cannot be empty");
    return;
  }
  // if(save.prompt == "") {
  //   notifyError("Prompt cannot be empty");
  //   return;
  // }
  if(save.action == ""){
    notifyError("Action cannot be empty");
    return;
  }
  const flag = await assistantStore.savePromptData(save)
  if(!flag){
    notifyError("Save failed");
    return;
  }
}
</script>
<template>
  <el-form label-width="100px" style="margin-top: 12px">
    <el-form-item :label="t('common.title')">
        <el-input
          v-model="formData.name"
          :placeholder="t('common.inputTitle')"
          prefix-icon="Notification"
          clearable
          :autofocus="true"
        ></el-input>
      </el-form-item>
      <el-form-item :label="t('common.cate')">
        <el-select v-model="formData.action">
          <el-option
            v-for="(item, key) in getCateList()"
            :key="key"
            :label="item.title"
            :value="item.value"
          />
        </el-select>
      </el-form-item>
      <el-form-item :label="t('common.isDef')">
        <el-switch 
        v-model="formData.isdef" 
        active-value="1"
        inactive-value="0"
        />
      </el-form-item>
      <el-form-item :label="t('common.content')">
        <el-input
          type="textarea"
         v-model="formData.prompt"
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
