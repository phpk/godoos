<template>
    <!-- 邀请群聊对话框 -->
    <el-dialog v-model="store.inviteFriendDialogVisible" title="邀请群聊" width="80%" style="height: 450px" align-center>
        <div>
            <el-transfer :titles="['可选项', '已选项']" filterable filter-placeholder="搜索用户名" style="height: 300px"
                v-model="users" :data="data" :props="{ key: 'key', label: 'label', avatar: 'avatar' }">
                <!-- 自定义穿梭框列表项模板 -->
                <template #default="{ option }">
                    <el-avatar :src="option.avatar" size="small" style="margin-right: 5px" />
                    <span>{{ option.label }}</span>
                </template>
            </el-transfer>
            <el-button class="chat-canel-button" @click="store.inviteFriendDialogVisible = false">取消</el-button>
            <el-button @click="
                store.inviteFriend(store.targetGroupInfo.group_id, users)
                " class="invite-group-button" style="">确定</el-button>
        </div>
    </el-dialog>
    <!-- 群成员抽屉 -->
    <el-drawer v-model="store.groupMemberDrawerVisible" direction="rtl" title="群成员" :with-header="true"
        :show-close="true" size="250px">
        <!-- 群成员列表 -->
        <div class="group-member-container">
            <div class="member-list" v-for="item in store.groupMembers" :key="item.id">
                <div class="member-item">
                    <!-- 头像和昵称 -->
                    <div class="avatar-container">
                        <el-avatar style="width: 100%; height: 100%" :src="item.avatar" />
                    </div>
                    <span>{{ item.nickname }}</span>
                </div>
            </div>
        </div>
    </el-drawer>
    <div class="chatbox-main" v-if="store.targetChatType === 'user' || store.targetChatType === 'group'">
        <!--聊天顶部区-->
        <el-header class="chat-header">
            <div class="header-title">
                <span v-if="store.targetUserInfo.displayName" class="header-title-name">{{
                    store.targetUserInfo.displayName
                    }}</span>
                <span v-else-if="store.targetGroupInfo.displayName" class="header-title-name">{{
                    store.targetGroupInfo.displayName }}</span>
                <div v-if="
                    store.targetGroupInfo &&
                    Object.keys(store.targetGroupInfo).length > 0
                ">
                    <div style="display: flex; gap: 10px">
                        <el-dropdown placement="bottom" style="border: none">
                            <el-icon class="chat-dropdown-icon">
                                <More />
                            </el-icon>
                            <template #dropdown>
                                <el-dropdown-menu>
                                    <el-dropdown-item @click="
                                        store.groupMemberDrawerVisible = true
                                        ">群成员</el-dropdown-item>

                                    <el-dropdown-item @click="openInviteGroupDialog()">邀请群聊</el-dropdown-item>
                                    <el-dropdown-item @click="
                                        store.quitGroup(
                                            store.targetGroupInfo.group_id
                                        )
                                        ">退出群聊</el-dropdown-item>
                                    <el-dropdown-item @click="clearMessages('group')">清空记录</el-dropdown-item>
                                </el-dropdown-menu>
                            </template>
                        </el-dropdown>
                    </div>
                </div>

                <div v-else class="chat-dropdown-menu">
                    <el-dropdown>
                        <el-icon class="chat-dropdown-icon">
                            <More />
                        </el-icon>
                        <template #dropdown>
                            <el-dropdown-menu>
                                <el-dropdown-item @click="clearMessages('user')">清空记录</el-dropdown-item>
                            </el-dropdown-menu>
                        </template>
                    </el-dropdown>
                </div>
            </div>
        </el-header>
        <!--聊天主体区-->
        <el-main class="msg-main">
            <!-- 聊天消息滚动区域 -->
            <div class="msg-container">
                <el-scrollbar ref="scrollbarRef">
                    <div ref="innerRef">
                        <ChatMessage />
                    </div>
                </el-scrollbar>
            </div>
        </el-main>

        <!--聊天输入区和发送按钮等-->
        <el-footer class="msg-footer">
            <!-- 输入主体 -->
            <div class="input-main">
                <el-input size="large" style="width: 100%; height: 100%" placeholder="请输入内容"
                    @keyup.enter.exact="store.sendMessage('text')" v-model="store.message">
                    <template #suffix>
                        <el-icon :size="20" class="input-option-icon" @click="selectImg">
                            <Picture />
                        </el-icon>
                        <el-icon :size="20" class="input-option-icon" @click="selectFile">
                            <Link />
                        </el-icon>
                        <el-icon @click="store.sendMessage('text')" :size="22" class="input-button">
                            <Promotion />
                        </el-icon>
                    </template>
                </el-input>
            </div>
        </el-footer>
    </div>
    <div class="chatbox-main" v-if="store.targetChatType === 'system'">
        <!--聊天顶部区-->
        <el-header class="chat-header">
            <div class="header-title">
                <span v-if="store.systemInfo" class="header-title-name">系统消息</span>

                <div class="chat-dropdown-menu">
                    <el-dropdown>
                        <el-icon class="chat-dropdown-icon">
                            <More />
                        </el-icon>
                        <template #dropdown>
                            <el-dropdown-menu>
                                <el-dropdown-item @click="clearMessages('system')">清空记录</el-dropdown-item>
                            </el-dropdown-menu>
                        </template>
                    </el-dropdown>
                </div>
            </div>
        </el-header>
        <!--聊天主体区-->
        <el-main class="system-main">
            <div class="msg-container">
                <el-scrollbar ref="scrollbarRef">
                    <div ref="innerRef">
                        <div class="msg-card" v-for="item in store.systemInfo" :key="item.id">
                            <div class="msg-time">{{ formatTimetoYMD(item.time) }}</div>
                            <el-card class="box-card">
                                {{ item.previewMessage }}
                            </el-card>
                        </div>
                    </div>
                </el-scrollbar>
            </div>
        </el-main>
    </div>
    <div class="empty-chatbox" v-if="!store.targetChatType">
        <el-empty description="请选择一个用户或群聊" />
    </div>
</template>

<script setup lang="ts">
import { useChatStore } from "@/stores/chat";
import { useFileSystemStore } from "@/stores/filesystem";
import { formatTimetoYMD } from "@/utils/time";
import { errMsg, successMsg } from "@/utils/msg";
import { ref, toRaw, watch, watchEffect } from "vue";

const store: any = useChatStore();
const fs = useFileSystemStore();
const imgExt = ["png", "jpg", "jpeg", "gif", "bmp", "webp", "svg"];
const choosetype = ref("");
const scrollbarRef = ref(null);
const innerRef = ref(null);


// 清空聊天记录
const clearMessages = (type: string) => {
    // 首先先删除发送的聊天记录，如果删除成功，再删除接收的聊天记录，如果发送消息删除失败也调用删除接收的聊天记录。如果删除对方的记录成功也弹出删除成功
    if (type === "user") {
        if (store.clearSentMessages()) {
            if (store.clearReceivedMessages()) {
                successMsg("删除成功");
                store.getSessionInfo(store.targetChatId, type);
            } else {
                errMsg("删除失败");
            }
        } else if (store.clearReceivedMessages()) {
            successMsg("删除成功");
        } else {
            errMsg("删除失败");
        }
    } else if (type === "group") {
        store.clearGroupMessages();
    } else if (type === 'system') {
        store.clearSystemMessages()
    }
};

const openInviteGroupDialog = () => {
    store.inviteFriendDialogVisible = true;
    store.getInviteUserList();
};

function scrollToBottom() {
    store.setScrollToBottom(innerRef, scrollbarRef);
}

const generateData = () => {
    return store.inviteUserList.map((user: any) => ({
        key: user.id,
        label: user.nickname,
        avatar: user.avatar, // 添加头像数据
    }));
};

const data = ref([]);

watchEffect(() => {
    if (store.inviteUserList && store.inviteUserList.length > 0) {
        data.value = generateData();
    }
});
// 声明 users 时指定类型为 any[]
const users = ref<any[]>([]);

watchEffect(() => {
    if (store.allUserList.length > 0) {
        data.value = generateData();
    }
});

// 监听store中的messageSendStatus.value = true，调用scrollToBottom
watch(
    () => store.messageSendStatus,
    (newVal, _) => {
        console.log("messageSendStatus 变化了:", newVal);
        if (newVal) {
            console.log("messageSendStatus 变化了:", newVal);
            scrollToBottom();
        }
    }
);

// 监听store中的messageReceiveStatus，调用scrollToBottom
watch(
    () => store.messageReceiveStatus,
    (newVal, _) => {
        if (newVal) {
            scrollToBottom();
        }
    }
);

// 监听store.drawerVisible
watch(
    () => store.drawerVisible,
    (newVal, _) => {
        if (newVal) {
            store.getGroupMember(store.targetGroupInfo.group_id);
        }
    }
);

function selectImg() {
    // choosetype.value = "image";
    console.log(choosetype.value);
    fs.chooseFiles(imgExt);
}
function selectFile() {
    // console.log("selectFile");
    choosetype.value = "applyfile";
    fs.chooseFiles(["*"]);
}
watch(
    () => fs.choose.paths,
    (newVal) => {
        if (newVal && newVal.length > 0) {
            store.sendInfo = toRaw(fs.choose.paths);
            store.sendMessage(choosetype.value);
            fs.clearChoose();
        }
    }
);

</script>

<style lang="scss" scoped>
@use "@/styles/chatbox.scss";
</style>
