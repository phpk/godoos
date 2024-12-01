<script lang="ts" setup>
import { ref } from "vue";
import { t } from "@/i18n/index";
const activeName = ref("system");
</script>
<template>
  <div>
    <el-tabs v-model="activeName" class="setting-tabs" style="margin: 12px">
      <el-tab-pane :label="t('aisetting.modelSetting')" name="system">
        <ai-setting-api />
      </el-tab-pane>
      <el-tab-pane :label="t('aisetting.defModel')" name="modelDef">
        <ai-setting-def />
      </el-tab-pane>
      <el-tab-pane :label="t('aisetting.chatSetting')" name="chatSetting">
        <ai-setting-conf />
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
  height: 75vh;
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
