<template>
    <el-select v-model="receiverId" remote :remote-method="handleSearch" filterable multiple placeholder="选择人员"
        popper-class="custom-header" value-key="id" @change="checkUsers">
        <template #header>
            <el-checkbox v-model="checkAll" @change="handleCheckAll">
                全选
            </el-checkbox>
        </template>
        <el-option v-for="item in userList" :key="item.id" :label="item.nickname" :value="item.id"
            style="display: flex;align-items: center;">
            <el-avatar size="small" :src="item.avatar || defaultAvatar" style="margin-right: 10px"></el-avatar>
            {{ item.nickname }}
        </el-option>
        <template #footer v-if="total > 10">
            <el-pagination size="small" background layout="prev, pager, next" :total="total" :page-size="pageSize"
                :current-page="currentPage" @current-change="handlePageChange" style="margin-top: 10px" />
        </template>
    </el-select>
</template>
<script setup lang="ts">
import { selectUserList, searchSelectUsers } from "@/api/share";
import { onMounted, ref, watch } from "vue";
import defaultAvatar from '/logo.png';
import { useChatStore } from "@/stores/chat";

const store = useChatStore()
const userList = ref<any[]>([])
const currentPage = ref(1);
const pageSize = 10;
const checkAll = ref(false);
const total = ref(0);
const nickname = ref("");
const receiverId = ref<any[]>([]);

const props = defineProps({
    ids: Array,
});
onMounted(() => {
    initData()
})
watch(() => store.groupChatInvitedDialogVisible, () => {
    receiverId.value = []
    store.departmentName = ''
})
const initData = () => {
    selectUserList(currentPage.value, nickname.value).then((res) => {
        resList(res.list);
        total.value = res.total;
    });
    if (props?.ids && props?.ids.length > 0) {
        searchSelectUsers(props?.ids).then((res) => receiverId.value = res);
    }
};
const resList = (list: any) => {
    if (!list || list.length === 0) return;
    userList.value = list.map((d: any) => {
        return {
            id: d.id,
            nickname: d.nickname,
            avatar: d.avatar
        };
    });
}
const handleCheckAll = (val: any) => {
    if (val) {
        receiverId.value = userList.value.map((d: any) => d.id);
    } else {
        receiverId.value = [];
    }
};
const checkUsers = (val: any) => {
    const res: any = [];
    val.forEach((item: any) => {
        if (item) {
            res.push(item);
        }
    });
    receiverId.value = res;
};
const handlePageChange = async (page: number) => {
    currentPage.value = page;
    await initData()
};
const handleSearch = async (val: any) => {
    currentPage.value = 1;
    nickname.value = val;
    await initData()
};

defineExpose({
    receiverId
});
</script>
<style scoped></style>