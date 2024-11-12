<script lang="ts" setup>
	import { useHistoryStore } from "@/stores/history";
	import { BrowserWindow, System } from "@/system";
	import { getSplit, getSystemConfig } from "@/system/config";
	import { decodeBase64, isBase64 } from "@/util/file";
	import { notifyError, notifySuccess } from "@/util/msg";
	import { getMdOption } from "@/util/vditor";
	import { saveAs } from "file-saver";
	import moment from "moment";
	import Vditor from "vditor";
	import "vditor/dist/index.css";
	import { inject, onMounted, onUnmounted, ref, toRaw } from "vue";

	const historyStore = useHistoryStore();

	const drawerBox = ref(false);
	const baseTitle = ref("");
	const content = ref("");
	const sys: any = inject<System>("system");
	const win: any = inject<BrowserWindow>("browserWindow");
	const isSaveing = ref(false);
	const vditor = ref();
	const filepath: any = ref("");
	const fileInput: any = ref(null);
	const SP = getSplit();
	console.log(SP);

	const debouncedHandleKeyDown = (event: KeyboardEvent) => {
		// 确保仅在我们的按钮获得焦点时处理快捷键
		if (
			(event.metaKey || event.ctrlKey) &&
			event.key.toLowerCase() === "s"
		) {
			event.stopPropagation(); // 先阻止事件冒泡
			event.preventDefault(); // 再阻止默认行为
			if (!isSaveing.value) {
				saveData();
			}
		}
	};

	const getDateTime = (t: any) => {
		return moment(t).format("MM-DD HH:mm");
	};
	function getTitle() {
		let title = win.getTitle();
		console.log(title);
		if (title.indexOf(SP) > -1) {
			title = title.split(SP).pop();
			title = title.split(".");
			title.pop();
			return title.join(".");
		} else {
			return "";
		}
	}
	onMounted(async () => {
		const editorOptions: any = getMdOption();
		editorOptions.input = (val: any) => {
			content.value = val;
		};

		vditor.value = new Vditor("vditorContainer", editorOptions);
		baseTitle.value = getTitle();
		if (win.config && win.config.content) {
			let winContent = toRaw(win.config.content);
			if (winContent && isBase64(winContent)) {
				winContent = decodeBase64(winContent);
			}
			//console.log(winContent);
			setTimeout(() => {
				if (winContent && winContent != "") {
					//console.log(winContent)
					vditor.value.setValue(winContent);
				}
			}, 1000);
		}
		document.addEventListener("keydown", debouncedHandleKeyDown);
	});
	// 清理监听器，确保组件卸载时移除，避免内存泄漏
	onUnmounted(() => {
		document.removeEventListener("keydown", debouncedHandleKeyDown);
	});

	async function saveData() {
		if (isSaveing.value) {
			return;
		}
		isSaveing.value = true;
		if (baseTitle.value == "") {
			notifyError("请输入标题");
			isSaveing.value = false;
			return;
		}
		//console.log(win.config.path)
		if (!filepath.value || filepath.value == "") {
			filepath.value = `${SP}C${SP}Users${SP}Desktop${SP}${baseTitle.value}.md`;
		}
		let refreshDesktop = false;
		if (win.config.path) {
			filepath.value = win.config.path;
			let fileTitleArr = filepath.value.split(SP).pop().split(".");
			fileTitleArr.pop();
			const oldTitle = fileTitleArr.join(".");
			if (oldTitle != baseTitle.value) {
				filepath.value = filepath.value.replace(
					oldTitle,
					baseTitle.value
				);
				refreshDesktop = true;
			}
		} else {
			refreshDesktop = true;
		}
		const file = await sys?.fs.getShareInfo(filepath.value);
		const isWrite =
			file.fs.sender === getSystemConfig().userInfo.id
				? 1
				: file.fs.is_write;
		//console.log(path)
		const res = file.fi.isShare
			? await sys?.fs.writeShareFile(
					filepath.value,
					vditor.value.getValue(),
					isWrite
			  )
			: await sys?.fs.writeFile(filepath.value, vditor.value.getValue());
		if (res.success) {
			notifySuccess(res.message || "保存成功！");
		} else {
			notifyError(res.message || "保存失败");
		}

		isSaveing.value = false;
		if (refreshDesktop) {
			sys.refershAppList();
		}
		historyStore.addList("markdown", {
			title: toRaw(baseTitle.value),
			path: filepath.value,
			time: Date.now(),
		});
	}
	function uploadFile(event: any) {
		const file = event.target.files[0];
		if (!file) {
			return;
		}
		const reader = new FileReader();
		reader.onload = (e: any) => {
			baseTitle.value = file.name.substring(
				0,
				file.name.lastIndexOf(".")
			);
			vditor.value.setValue(e.target.result);
		};
		reader.readAsText(file);
	}
	function importMd() {
		fileInput.value.click();
	}
	function download() {
		if (baseTitle.value == "") {
			notifyError("标题不能为空");
			return;
		}
		const contentData = vditor.value.getValue();
		if (contentData == "") {
			notifyError("内容不能为空");
			return;
		}
		let blob = new Blob([contentData], {
			type: "text/plain;charset=utf-8",
		});
		saveAs(blob, baseTitle.value + ".md");
	}
</script>
<template>
	<el-drawer
		v-model="drawerBox"
		direction="ltr"
		style="height: 100vh"
		:show-close="false"
		:with-header="false"
	>
		<div
			class="list-item"
			v-for="(item, index) in historyStore.getList('markdown')"
			:key="index"
		>
			<div class="list-title">
				<el-tooltip
					class="box-item"
					effect="dark"
					:content="item.path"
					placement="top-start"
				>
					{{ item.title }}
				</el-tooltip>
			</div>
			<div class="list-time">
				{{ getDateTime(item.time) }}
			</div>
		</div>
	</el-drawer>
	<el-row
		justify="space-between"
		:gutter="20"
		:span="24"
		style="margin: 10px 20px"
	>
		<el-col :span="5">
			<el-button
				@click.stop="drawerBox = !drawerBox"
				icon="Menu"
				circle
			/>
			<el-button
				@click.stop="importMd"
				icon="Upload"
				circle
			/>

			<el-button
				@click.stop="download"
				icon="Download"
				circle
			/>
			<input
				type="file"
				ref="fileInput"
				accept=".md"
				style="display: none"
				@change="uploadFile"
			/>
		</el-col>
		<el-col :span="15">
			<el-input
				v-model="baseTitle"
				placeholder="输入标题"
			/>
		</el-col>
		<el-col :span="2">
			<el-button
				@click.stop="saveData()"
				icon="Finished"
				:loading="isSaveing"
				circle
			/>
		</el-col>
	</el-row>
	<div
		id="vditorContainer"
		style="margin: 10px 20px"
	></div>
</template>
<style scoped>
	@import "@/assets/left.scss";
</style>
