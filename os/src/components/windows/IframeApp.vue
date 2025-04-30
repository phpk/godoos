<script setup lang="ts">
	import * as fs from "@/api/files";
	import { getExportType } from "@/router/filemaplist";
	import { useAiChatStore } from "@/stores/aichat";
	import { useFileSystemStore } from "@/stores/filesystem";
	import { noticeMsg } from "@/utils/msg";
	import {
		defineProps,
		onMounted,
		onUnmounted,
		ref,
		toRaw,
		watch,
	} from "vue";
	const storeRef = ref<HTMLIFrameElement | null>(null);
	const fileSystemStore = useFileSystemStore();
	const aiChatStore = useAiChatStore();
	let hasInit = false;
	const props = defineProps({
		win: {
			type: Object,
			required: true,
		},
	});
	let win = toRaw(props.win);
	let prop: any = win.props;
	async function eventHandler(event: any) {
		//if (event.origin !== win.origin) return;
		const eventData = event.data;
		const isShare = prop.isShare === "true";
		const fileMap = getExportType(eventData.type);
		if (fileMap && eventData.type == fileMap.eventType) {
			//eventData.type = fileMap.eventType;
			let contentData = JSON.parse(eventData.data);
			if (contentData.contentType) {
				contentData.content = fileSystemStore.parserFormData(
					contentData.content,
					contentData.contentType
				);
			}
			const ext = contentData.ext || fileMap.ext[0];
			if (!prop.path) {
				await fileSystemStore.saveFiles(
					contentData.title,
					ext,
					contentData.content
				);
				return;
			} else {
				if (!contentData.title.endsWith(ext)) {
					contentData.title = `${contentData.title}.${ext}`;
				}
				prop.path = fs.join(fs.dirname(prop.path), contentData.title);
			}
			//win.component = fileMap.editor;
			const res = await fileSystemStore.handleWriteFile(
				prop.path,
				contentData.content
			);
			if (!contentData.isShare) {
				if (res) {
					noticeMsg("文件保存成功！");
				} else {
					noticeMsg("文件保存失败！");
				}
			}
		}
		switch (eventData.type) {
			case "initSuccess":

				if (hasInit) {
					return;
				}
				hasInit = true;
				const postData: any = {
					title: prop.title,
					ext: prop.ext,
					content: prop.content,
					isShare: isShare,
				};
				if (prop.path) {
					const res = await fileSystemStore.handleReadFile(prop.path);
					if (res) {
						postData.content = res;
					}
				}
				storeRef.value?.contentWindow?.postMessage(
					{
						type: "init",
						data: postData,
					},
					"*"
				);
				break;
			case "aiCreater":
				// 模拟AI返回数据
				const res: any = await aiChatStore.askAi(eventData);
				storeRef.value?.contentWindow?.postMessage(
					{
						type: "aiReciver",
						data: res,
						action: eventData.action,
					},
					"*"
				);
				break;
			case "saveMind":
				const data = eventData.data;
				const file = {
					title: data.title,
					content: data.content,
					ext: "mind",
					isFile: true,
					isShare: false,
					isDirectory: false,
				};
				fileSystemStore.openFile(file);
				break;
			default:
				break;
		}
	}

	onMounted(() => {
		window.addEventListener("message", eventHandler);
	});
	onUnmounted(() => {
		window.removeEventListener("message", eventHandler);
	});
	watch(
		() => props.win,
		() => {
			window.addEventListener("message", eventHandler);
		}
	);
	watch(
		() => fileSystemStore.currentShareFile,
		async () => {
			if (!fileSystemStore.currentShareFile.path) return;
			if (!win.props.path) return;
			const fileStat: any = await fileSystemStore.handleFileStat(
				win.props.path
			);
			if (fileStat.isDirectory) return;
		
			if (fileStat.truePath != fileSystemStore.currentShareFile.path)
				return;
			const content = await fileSystemStore.handleReadFile(prop.path);
			storeRef.value?.contentWindow?.postMessage(
				{
					type: "init",
					data: {
						title: fileStat.title,
						ext: fileStat.ext,
						content: content,
						isShare: true,
						isRefresh: true,
					},
				},
				"*"
			);
			fileSystemStore.currentShareFile = {};
		}
	);
</script>
<template>
	<iframe
		class="setiframe"
		allow="fullscreen"
		ref="storeRef"
		:src="win.component"
	></iframe>
</template>
<style scoped>
	.setiframe {
		width: 100%;
		height: 100%;
		border: none;
	}
</style>
