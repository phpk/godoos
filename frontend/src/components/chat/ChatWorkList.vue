<script setup>
import { ref } from "vue";
import { useChatStore } from "@/stores/chat";
const store = useChatStore()

</script>

<template>
    <el-menu default-active="1" class="chat-menu">
        <el-menu-item index="1">
          <el-icon><Guide /></el-icon>
          <span>我要申请</span>
        </el-menu-item>
        <el-menu-item index="2">
          <el-icon><Place /></el-icon>
          <span>我的申请</span>
        </el-menu-item>
        <el-menu-item index="3">
          <el-icon><DishDot /></el-icon>
          <span>我的任务</span>
        </el-menu-item>
        <!-- <el-menu-item index="4">
          <el-icon><SuitcaseLine /></el-icon>
          <span>安排任务</span>
        </el-menu-item> -->
        <el-menu-item index="5">
          <el-icon><Present /></el-icon>
          <span>我的审批</span>
        </el-menu-item>
    </el-menu>
</template>
<style scoped>
    .chat-menu{
        background: none;
        border: none;
    }
</style>