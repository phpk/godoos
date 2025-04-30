<script setup lang="ts">
import { useModelStore } from "@/stores/model.ts";
const modelStore = useModelStore();

</script>
<template>
    <el-form label-position="right">
      <el-collapse>
        <el-collapse-item :title="'&nbsp;&nbsp;' + item.name + '配置'" :name="item.cpp" v-for="item in modelStore.aiEngine" class="pd10">
          <el-form-item label="私钥配置" v-if="'apikey' in item" class="inline-layout">
              <el-input v-model="item.apikey" placeholder="私钥配置"
                prefix-icon="Notification" clearable size="large"></el-input>
          </el-form-item>
          <el-form-item label="访问地址" v-if="'url' in item">
              <el-input v-model="item.url" placeholder="访问地址配置"
                prefix-icon="Notification" clearable size="large"></el-input>
          </el-form-item>
        </el-collapse-item>

      </el-collapse>
    </el-form>
</template>
<style scoped lang="scss">
.pd10{
  padding: 10px;
}
.slider-container {
  width: 90%;
  margin: 0 5%;
}
</style>