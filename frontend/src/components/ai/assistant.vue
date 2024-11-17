<script lang="ts" setup>
import {useAssistantStore} from "@/stores/assistant";
import { t } from "@/i18n/index";
import { onMounted } from "vue";
import { notifyError, notifySuccess } from "@/util/msg"
const assistantStore:any = useAssistantStore()

onMounted(async () => {
  await assistantStore.getPromptList()
})
function add() {
  assistantStore.editId = 0
  assistantStore.showAdd = true
}
async function edit(item:any){
  assistantStore.editId = item.id
  assistantStore.showAdd = true
}
async function del(item:any){
 const flag = await assistantStore.deletePrompt(item.id)
 if(!flag){
  notifyError(t("prompt.contDelete"))
  return
 }else{
  notifySuccess(t("prompt.delSuccess"))
 }
}
</script>
<template>
  <el-dialog v-model="assistantStore.showAdd" width="600" :title="t('common.add')">
    <assistant-add />
  </el-dialog>
  <el-drawer
    v-model="assistantStore.showLeft"
    direction="ltr"
    style="height: 100vh"
    width="300px"
    :show-close="false"
    :with-header="false"
  >
    <el-space direction="vertical" alignment="flex-start">
      <el-button type="primary" style="width:150px">
      {{ t('prompt.cate') }}
    </el-button>
      <el-button 
      @click.stop="assistantStore.changeCate('all')"  
      text >
      {{ t('prompt.all') }}
    </el-button>
      <el-button 
      @click.stop="assistantStore.changeCate(item)" 
      v-for="item in assistantStore.promptAction" 
      text >
      {{ t('prompt.'+item) }}
    </el-button>
    </el-space>
  </el-drawer>
  <el-page-header icon="null">
    <template #title>
      <div></div>
    </template>
    <template #content>
      <div class="flex items-center">
        <el-button @click.stop="assistantStore.handlerLeft" icon="Menu" circle />
        <el-button @click.stop="add" icon="Plus" circle />
      </div>
    </template>
  </el-page-header>
  <el-scrollbar v-if="assistantStore.promptList.length > 0">
    <div style="padding: 15px 15px 50px 15px">
      <el-card
        class="model-item"
        v-for="(item, key) in assistantStore.promptList"
        :key="key"
        shadow="hover"
        style="margin: 8px 0"
      >
      <div class="card-header">
          <el-row type="flex" justify="space-between">
            <span style="font-size:14px">{{ item.name }}</span>
            <div class="card-header-right">
              <el-button icon="Edit" circle @click="edit(item)" />
              <el-button icon="Delete" circle @click="del(item)" />
            </div>
          </el-row>
        </div>
        <div style="display: flex; align-items: center; margin-top: 5px;">
          <el-tag size="small" style="margin-right: 5px;">{{ item.action }}</el-tag>
          <el-tag size="small" style="margin-right: 5px;">{{ item.lang }}</el-tag>
          <el-tag size="small" v-if="item.isdef > 0">default</el-tag>
        </div>
    </el-card>
    <el-row
        justify="center"
        style="margin-top: 15px"
        v-if="assistantStore.page.pages > 1"
      >
        <el-pagination
          background
          layout="prev, pager, next"
          v-model:current-page="assistantStore.page.current"
          v-model:page-size="assistantStore.page.size"
          :total="assistantStore.page.total"
          @current-change="(val:any) => assistantStore.pageClick(val)"
        />
      </el-row>
  </div>
  </el-scrollbar>
</template>

<style scoped lang="scss">
.main {
  flex: 1;
  width: 100%;
  height: 100%;
  background-color: #fff;
}

.model-list {
  height: 100%;
  overflow: scroll;
  .model-item {
    transition: all 0.3s;
    border-bottom: 1px solid #eee;
    &:hover {
      transition: all 0.3s;
      background-color: rgba(99, 99, 99, 0.2);
      box-shadow: rgba(99, 99, 99, 0.2) 0px 2px 8px 0px !important;
      cursor: pointer;
    }
  }
}


</style>



