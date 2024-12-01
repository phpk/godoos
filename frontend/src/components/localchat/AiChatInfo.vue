<script setup lang="ts">
import { useAiChatStore } from "@/stores/aichat";
import { t } from "@/i18n/index";
import { notifyInfo,notifySuccess } from "@/util/msg.ts";
const chatStore = useAiChatStore();
function changPrompt(id:number) {
  const promptData:any = chatStore.promptList.find((item:any) => {
    return item.id == id;
  });
  //console.log(promptData)
  if(promptData) {
    chatStore.editInfo.prompt = promptData.prompt;
  }
}
const changeInfo = async () => {
  const info = chatStore.editInfo;
  console.log(info)
  if (!info.title) {
    notifyInfo(t("aichat.inputTitle"));
    return;
  }
  if (!info.model) {
    notifyInfo(t("aichat.selectModel"));
    return;
  }
  if (!info.prompt) {
    info.prompt = "";
  }
  delete info.id;

  if (chatStore.isEditor) {
    await chatStore.updateChat(info, chatStore.activeId);
    notifySuccess(t("aichat.editsuccess"));
  } else {
    const promptData = {
      prompt: info.prompt,
      promptId: info.promptId,
    };
    const modelData = chatStore.modelList.find((item:any) => {
      return item.model == info.model;
    });
    await chatStore.addChat(info.title, modelData, promptData, "");
    notifySuccess(t("aichat.addsuccess"));
  }
  await chatStore.getActiveChat();
  chatStore.showInfo = false;
};

</script>
<template>
    <el-form label-width="150px" style="margin-top: 12px">
      <el-form-item :label="t('common.title')">
        <el-input
          v-model="chatStore.editInfo.title"
          :placeholder="t('aichat.inputTitle')"
          prefix-icon="Notification"
          clearable
          :autofocus="true"
        ></el-input>
      </el-form-item>
      <el-form-item :label="t('aichat.model')">
        <el-select v-model="chatStore.editInfo.model">
          <el-option
            v-for="(item, key) in chatStore.modelList"
            :key="key"
            :label="item.model"
            :value="item.model"
          />
        </el-select>
      </el-form-item>
      <!-- <el-form-item :label="t('chat.refknow')">
        <el-select
          v-model="chatStore.editInfo.kid"
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
      </el-form-item> -->
      <el-form-item :label="t('aichat.role')">
        <el-select
          v-model="chatStore.editInfo.promptId"
          @change="changPrompt"
        >
          <el-option
            v-for="(item, key) in chatStore.promptList"
            :key="key"
            :label="item.name"
            :value="item.id"
          />
        </el-select>
      </el-form-item>
      
      <el-form-item>
        <el-input
          type="textarea"
          v-model="chatStore.editInfo.prompt"
          :rows="6"
        ></el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="changeInfo">{{
          t("common.confim")
        }}</el-button>
      </el-form-item>
    </el-form>
</template>