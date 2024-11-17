<script lang="ts" setup>
import {db} from "@/stores/db"
import { useModelStore } from "@/stores/model";
import {notifyError} from "@/util/msg";
import { getSystemKey } from "@/system/config";
import { t } from "@/i18n/index";
import { watchEffect, ref, toRaw } from "vue";
const modelStore = useModelStore();
const props = defineProps(['labelId']);

const emit = defineEmits(["closeFn","refreshFn"]);
// function close() {
//   emit("closeFn", false);
// }
const labelData:any = ref({
  name : "",
  zhdesc : "",
  endesc : "",
  family: "",
  engine:"llm",
  action:[],
})
async function save() {
  const saveData = toRaw(labelData.value)
  if(saveData.name == "") {
    notifyError(t('common.inputTitle'))
    return;
  }
  if(saveData.type == "") {
    notifyError(t('model.selectEngine'))
    return;
  }
  if(saveData.action.length == 0){
    notifyError(t('model.selectCategory'))
    return;
  }
  if(props.labelId > 0) {
    await db.update("modelslabel", props.labelId, saveData)
  }else{
    saveData.models = []
    saveData.chanel = getSystemKey('currentChanel')
    console.log(saveData)
    await db.addOne("modelslabel", saveData)
  }
  emit("closeFn", false);
  emit("refreshFn", true);
}
watchEffect(async () => {
  if(props.labelId > 0) {
    labelData.value = await db.getOne("modelslabel", props.labelId)
  }else{
    labelData.value = {
      name : "",
      zhdesc : "",
      endesc : "",
      family: "",
      engine:"llm",
      action:[],
    }
  }

})
</script>
<template>
  <el-form label-width="100px" style="margin-top:12px">
    <el-form-item :label="t('model.labelName')">
      <el-input
        v-model="labelData.name"
        :placeholder="t('model.labelName')"
        prefix-icon="House"
        clearable
        resize="none"
      ></el-input>
    </el-form-item>
    <el-form-item :label="t('model.family')">
      <el-input
        v-model="labelData.family"
        :placeholder="t('model.family')"
        prefix-icon="HomeFilled"
        clearable
        resize="none"
      ></el-input>
    </el-form-item>
    <el-form-item :label="t('model.category')">
      <el-select v-model="labelData.action" :multiple="true" :placeholder="t('model.selectCategory')">
      <el-option
          v-for="(item, key) in modelStore.cateList"
          :key="key"
          :label="t('model.'+item)"
          :value="item"
        />
    </el-select>
    </el-form-item>
    <el-form-item :label="t('model.engine')">
      <el-select v-model="labelData.engine"  :placeholder="t('model.selectEngine')">
      <el-option
        v-for="item,key in modelStore.modelEngines"
        :key="key"
        :label="item.name"
        :value="item.name"
      />
    </el-select>
    </el-form-item>
    
    <el-form-item :label="t('model.chineseDescription')">
      <el-input
        :placeholder="t('model.chineseDescription')"
        v-model="labelData.zhdesc"
      ></el-input>
    </el-form-item>
    <el-form-item :label="t('model.englishDescription')">
      <el-input
        :placeholder="t('model.englishDescription')"
        :row="3"
        v-model="labelData.endesc"
      ></el-input>
    </el-form-item>
    <el-form-item>
      <el-button type="primary" icon="CirclePlus" @click="save">{{ t('common.save') }}</el-button>
    </el-form-item>
  </el-form>
</template>