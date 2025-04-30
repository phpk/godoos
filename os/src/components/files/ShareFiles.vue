<template>
	<el-form :model="form" label-width="auto" style="max-width: 560px; margin-top: 20px; padding: 20px">
		<el-form-item label="分享给">
			<choose-user ref="chooseUserRef" />
		</el-form-item>
		<el-form-item label="编辑权限">
			<el-switch v-model="form.iswrite" active-value="1" inactive-value="0" />
		</el-form-item>
		<div class="btn-group">
			<el-button type="primary" @click="onSubmit">发布分享</el-button>
		</div>
	</el-form>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import { useContextMenuStore } from '@/stores/contextmenu'
import { useFileSystemStore } from '@/stores/filesystem'
import {errMsg,successMsg} from "@/utils/msg"
const contextMenuStore = useContextMenuStore()
const fileSystemStore = useFileSystemStore()
const chooseUserRef:any = ref(null);
const form: any = ref({
	receiverId: [],
	path: "",
	iswrite: "0",
});


const onSubmit = async () => {
	if (!chooseUserRef.value){
		return;
	}
	const receiverIds = chooseUserRef.value?.receiverId;
	if(receiverIds.length < 1){
		errMsg("请选择分享对象")
		return;
	}
	console.log(receiverIds)
	const postData = {
		path:contextMenuStore.currentFile?.path,
		receiverId: receiverIds,
		isWrite : form.value.iswrite * 1
	}
	//console.log(postData)
	const result:any = await fileSystemStore.handleShareFile(postData);
	if (result.success) {
		successMsg(result.message || "分享文件成功");
		contextMenuStore.isShareFile = false;
	} else {
		errMsg(result.message || "分享文件失败");
	}
};
</script>
<style scoped>
.btn-group {
	display: flex;
	justify-content: center;
}
</style>
