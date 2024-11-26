<script setup lang="ts">
import { useAiChatStore } from "@/stores/aichat";
const chatStore = useAiChatStore();
</script>
<template>
     <el-scrollbar>
      <el-space direction="vertical">
        <el-card
          v-for="(item, key) in chatStore.chatList"
          :key="key"
          class="box-card"
          style="width: 190px"
          :shadow="chatStore.activeId == item.id ? 'always' : 'hover'"
        >
          <el-row type="flex" justify="space-between">
            <el-button
              type="info"
              text
              @click.stop="chatStore.setActiveId(item.id)"
              size="small"
            >
              {{ item.title }}
            </el-button>
            <el-button
              icon="Delete"
              size="small"
              @click.stop="chatStore.deleteChat(item.id)"
              circle
            ></el-button>
          </el-row>
        </el-card>
      </el-space>
    </el-scrollbar>
</template>