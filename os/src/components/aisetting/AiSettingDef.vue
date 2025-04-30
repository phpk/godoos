<script setup lang="ts">
import { t } from "@/i18n/index";
import { useModelStore } from "@/stores/model.ts";
import { onMounted, ref, watch } from "vue";
const modelStore = useModelStore();
const currentsModel: any = ref({});
onMounted(async () => {
  await modelStore.getModelList();
  updateCurrentsModel();
});
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
watch(modelStore.modelList, () => {
  updateCurrentsModel();
});
</script>
<template>
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
</template>