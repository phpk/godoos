<script setup lang="ts">
import { ref, onMounted, watchEffect, watch } from 'vue'

import { useChatStore } from "@/stores/chat";
import { isMobileDevice } from "@/utils/device";
import { Search } from "@element-plus/icons-vue";
const store = useChatStore();
// const workUrl = getWorkflowUrl();
onMounted(() => {
    store.initChat();
});
const chooseUserRef = ref()
const drawer = ref(false)
</script>
<template>
    <el-container class="container">
        <!--菜单-->
        <el-aside class="menu">
            <workbench-menu />
        </el-aside>
        <el-container class="side" v-if="store.currentNavId < 4">
            <!--搜索栏-->
            <div v-if="store.currentNavId < 2">
                <el-header class="search">
                    <el-input placeholder="搜索" :prefix-icon="Search" class="search-input" v-model="store.searchInput" />
                    <!-- 邀请群聊 -->
                    <button class="inviteGroupChats" @click="store.setGroupChatInvitedDialogVisible(true)">
                        <el-icon>
                            <Plus />
                        </el-icon>
                    </button>
                </el-header>
                <!--好友列表-->
                <el-main class="list">
                    <el-scrollbar>
                        <chat-msg-list v-if="store.currentNavId == 0" />
                        <chat-user-list v-if="store.currentNavId == 1" />
                    </el-scrollbar>
                </el-main>
            </div>
            <ai-chat-left v-if="store.currentNavId == 3" />
            <work-table-left v-if="store.currentNavId == 2"></work-table-left>
        </el-container>
        <!-- 手机端侧边栏 -->
        <el-drawer v-if="isMobileDevice()" v-model="drawer" :with-header="false" direction="ltr" size="50%">
            <div v-if="store.currentNavId < 2">
                <el-header class="search">
                    <el-input placeholder="搜索" :prefix-icon="Search" class="search-input" v-model="store.searchInput" />
                    <button class="inviteGroupChats" @click="store.setGroupChatInvitedDialogVisible(true)">
                        <el-icon>
                            <Plus />
                        </el-icon>
                    </button>
                </el-header>
                <el-main class="list">
                    <el-scrollbar>
                        <chat-msg-list v-if="store.currentNavId == 0" />
                        <chat-user-list v-if="store.currentNavId == 1" />
                    </el-scrollbar>
                </el-main>
            </div>
            <ai-chat-left v-if="store.currentNavId == 3" />
            <work-table-left v-if="store.currentNavId == 2"></work-table-left>
        </el-drawer>
        <el-container class="chat-box">
            <el-button v-if="isMobileDevice()" icon="Menu" size="small" @click="drawer = !drawer"></el-button>
            <chat-box-main v-if="store.currentNavId < 1" />
            <chat-user-info v-if="store.currentNavId == 1"></chat-user-info>
            <ai-chat-main v-if="store.currentNavId == 3" />
            <work-table-main v-if="store.currentNavId == 2"></work-table-main>
        </el-container>
        <!-- <el-container class="chat-setting" v-if="store.currentNavId == 2">
            <iframe class="workflow" :src="workUrl"></iframe>
        </el-container> -->
        <!-- 邀请群聊弹窗 -->
        <el-dialog v-model="store.groupChatInvitedDialogVisible" title="创建群聊" align-center
            :fullscreen="isMobileDevice() ? true : false" :show-close="isMobileDevice() ? false : true" draggable>
            <template #header v-if="isMobileDevice()">
                <div class="dialog-header">
                    <el-button style="background-color: #0078d4; color: #fff"
                        @click="store.groupChatInvitedDialogVisible = false">取消</el-button>

                    <div class="dialog-title">创建群聊</div>
                    <el-button style="background-color: #0078d4; color: #fff"
                        @click="store.createGroupChat(chooseUserRef.receiverId)">确定</el-button>
                </div>
            </template>
            <div class="dialog-body">
                <!-- 添加输入部门名的输入框 -->
                <div>
                    <el-form label-position="top">
                        <el-form-item label="群聊名称:">
                            <el-input maxlength="8" show-word-limit class="department-name"
                                v-model="store.departmentName" placeholder="请输入群聊名称"></el-input>
                        </el-form-item>
                        <!-- 选择成员 -->
                        <el-form-item label="选择成员:">
                            <choose-user ref="chooseUserRef"></choose-user>
                        </el-form-item>
                    </el-form>
                </div>
            </div>

            <template #footer v-if="!isMobileDevice()">
                <span>
                    <el-button style="background-color: #0078d4; color: #fff"
                        @click="store.groupChatInvitedDialogVisible = false">取消</el-button>
                    <el-button style="background-color: #0078d4; color: #fff"
                        @click="store.createGroupChat(chooseUserRef.receiverId)">确定</el-button>
                </span>
            </template>
        </el-dialog>
    </el-container>
</template>

<style lang="scss" scoped>
@use '@/styles/workbench.scss';
</style>
