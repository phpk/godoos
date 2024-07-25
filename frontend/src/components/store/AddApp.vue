<script setup lang="ts">
import { ref } from "vue";
const formData = ref({
    importType: 'remote'
})
const addType = [
    {
        'name' : '远程下载',
        'type' : 'remote'
    },
    {
        'name' : '本地导入',
        'type' : 'local'
    },
    {
        'name' : '开发模式',
        'type' : 'dev'
    }
]

</script>
<template>
    <div class="setting">
        <div class="setting-item">
            <label>添加方式</label>
            <el-select v-model="formData.importType" style="width: 180px;">
                <el-option
                    v-for="item in addType"
                    :key="item.type"
                    :label="item.name"
                    :value="item.type"
                />
            </el-select>
        </div>
    </div>
</template>
<style scoped>

.setting-item {
  display: flex;
  flex-direction: row;
  justify-content: flex-start;
  align-items: flex-start;
  margin-bottom: 20px;
  gap: 20px;
}

.setting-item label {
  display: block;
  width: 100px;
  flex-shrink: 0;
  text-align: right;
}

.setting-item {
  display: flex;
  align-items: center;
}
</style>