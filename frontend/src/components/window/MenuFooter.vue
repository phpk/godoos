<template>
  <div class="menubar">
    <div class="file-detail">
      <label for="">文件名称：</label>
      <el-input v-model="fileName"/>
    </div>
    <el-button @click="cancleSave"  type="primary" size="small">取消</el-button>
    <el-button @click="saveFile"  size="small">保存</el-button>
  </div>
</template>

<script lang="ts" setup>
import { BrowserWindow } from "@/system/window/BrowserWindow";
import { UnwrapNestedRefs, ref } from "vue";
// import { emitEvent } from "@/system/event";
const props = defineProps<{
  browserWindow: UnwrapNestedRefs<BrowserWindow>;
}>();
const emit = defineEmits(['translateSavePath'])
let fileName = ref('未命名文件')

function cancleSave() {
  props.browserWindow.close()
}
function saveFile() {
  emit('translateSavePath','',fileName.value)
  props.browserWindow.close()
}
</script>

<style lang="scss" scoped>
.menubar {
  // width: 100%;
  height: 70px;
  padding: 10px 30px;
  background-color: aliceblue;
  .file-detail {
    width: 100%;
    font-size: 14px;
    label {
      width: 80px; 
    }
    .el-input {
      width: calc(100% - 80px);
      min-width: 100px;
      height: 30px;
    }
  }
  .el-button {
    width: 60px;
    height: 30px;
    float: right;
    margin: 5px;
  }
}
</style>