<template>
	<div class="save-file">
		<template v-if="!fileSystemStore.choose.isSave">
			<div class="button-container">
				<el-button
					type="primary"
					@click="handleChoose"
					>选择</el-button
				>
			</div>
		</template>
		<template v-else>
			<el-input
				v-model="fileSystemStore.choose.defName"
				placeholder="请输入文件名"
			/>
			<div class="button-container">
				<el-button
					type="primary"
					@click="handleSave"
					>保存</el-button
				>
				<el-button @click="handleCancel">取消</el-button>
			</div>
		</template>
	</div>
</template>

<script lang="ts" setup>
	import { join } from "@/api/net/files";
	import { useClickingStore } from "@/stores/clicking";
	import { useFileSystemStore } from "@/stores/filesystem";
	import { useWindowStore } from "@/stores/window";
	import { noticeMsg } from "@/utils/msg";
	const props = defineProps({
		currentPath: {
			type: String,
			required: true,
		},
		winId: {
			type: String,
			required: true,
		},
	});
	const windowStore = useWindowStore();
	const fileSystemStore = useFileSystemStore();
	const clickingStore = useClickingStore();
	const handleChoose = () => {
		if (clickingStore.clickedIcons.length === 0) {
			noticeMsg("请选择一个文件或文件夹");
			return;
		}
		windowStore.closeWindow(props.winId);
		fileSystemStore.choose.paths = [...clickingStore.clickedIcons];
		clickingStore.clickedIcons = [];
		//fileSystemStore.clearChoose();
	};
	const handleSave = async () => {
		const savePath = join(
			fileSystemStore.currentPath,
			fileSystemStore.choose.defName
		);
    
		await fileSystemStore.handleWriteFile(
			savePath,
			fileSystemStore.choose.content
		);

		// 保存逻辑
		windowStore.closeWindow(props.winId);
		fileSystemStore.clearChoose();
		noticeMsg("文件保存成功！");
	};

	const handleCancel = () => {
		// 取消逻辑
		windowStore.closeWindow(props.winId);
		fileSystemStore.clearChoose();
	};
</script>
<style scoped>
	.save-file {
		border-top: 1px solid #ebeef5;
		background-color: rgba(255, 255, 255, 0.95);
		padding: 10px;
	}

	.button-container {
		margin-top: 10px;
		text-align: right;
	}
</style>
