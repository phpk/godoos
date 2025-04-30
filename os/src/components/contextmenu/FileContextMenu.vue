<template>
	<div v-show="contextMenuStore.isContextMenuVisible" class="context-menu"
		:style="isMobileDevice() ? contextMenuStore.mobileStyle : style">
		<ul>
			<li @click="openFile">
				<el-icon :size="16" color="#333">
					<FolderOpened />
				</el-icon>
				打开
			</li>
			<li v-if="checkEditor()" @click="editorFile">
				<el-icon :size="16" color="#333">
					<Edit />
				</el-icon>
				编辑
			</li>
			<li @click="delFile" v-if="checkDelete()">
				<el-icon :size="16" color="#333">
					<Delete />
				</el-icon>
				删除
			</li>
			<li @click="joinKnowledge" v-if="checkJoin()">
				<el-icon :size="16" color="#333">
					<Notebook />
				</el-icon>
				加入知识库
			</li>
			
			<!-- 对话 -->
			<li @click="dialogue" v-if="checkDialogue()">
				<el-icon :size="16" color="#333">
					<ChatDotRound />
				</el-icon>
				知识库对话
			</li>
			<!-- 取消加入知识库 -->
			<li @click="joinKnowledge" v-if="checkUnJoin()">
				<el-icon :size="16" color="#333">
					<SetUp />
				</el-icon>
				取消知识库
			</li>
			<li @click="renameDir" v-if="checkDelete()">
				<el-icon ::size="16" color="#333">
					<EditPen />
				</el-icon>
				重命名
			</li>
			<li @click="copyFile" v-if="checkDelete()">
				<el-icon ::size="16" color="#333">
					<CopyDocument />
				</el-icon>
				复制(C)
			</li>
			<li @click="cutFile" v-if="checkDelete()">
				<el-icon ::size="16" color="#333">
					<Scissor />
				</el-icon>
				剪切(T)
			</li>
			<li v-if="checkZip()" @click="zipFile('zip')" @mouseenter="contextMenuStore.toggleSubMenu('zipFile')"
				@mouseleave="contextMenuStore.toggleSubMenu(null)" class="has-submenu">
				<el-icon :size="16" color="#333">
					<Wallet />
				</el-icon>
				压缩
				<div v-show="contextMenuStore.activeSubMenu === 'zipFile'" class="sub-menu" style="top: -40px">
					<ul>
						<li @click="zipFile('zip')">
							<icon name="zip" :size="16" />
							zip格式
						</li>
						<li @click="zipFile('tar')">
							<icon name="tar" :size="16" />
							tar格式
						</li>
						<li @click="zipFile('gz')">
							<icon name="gz" :size="16" />
							gzip格式
						</li>
					</ul>
				</div>
			</li>
			<li v-if="checkUnZip()" @click="unZipFile">
				<el-icon :size="16" color="#333">
					<FolderAdd />
				</el-icon>
				解压
			</li>
			<li v-if="checkDownload()" @click="download">
				<el-icon :size="16" color="#333">
					<Download />
				</el-icon>
				下载
			</li>
			<li v-if="checkFavorite()" @click="addFavorite">
				<el-icon :size="16" color="#333">
					<StarFilled />
				</el-icon>
				取消收藏
			</li>
			<li v-else @click="addFavorite">
				<el-icon :size="16" color="#333">
					<Star />
				</el-icon>
				收藏
			</li>
			<li v-if="checkShare()" @click="addShare()">
				<el-icon :size="16" color="#333">
					<Share />
				</el-icon>
				分享给...
			</li>
			<li v-if="checkUnShare()" @click="unShare()">
				<el-icon :size="16" color="#333">
					<Brush />
				</el-icon>
				取消分享
			</li>
			<li v-if="checkPwd()" @click="addPwd">
				<el-icon :size="16" color="#333">
					<SuitcaseLine />
				</el-icon>
				加密
			</li>
			<li v-if="checkUnpwd()" @click="unPwd">
				<el-icon :size="16" color="#333">
					<Failed />
				</el-icon>
				去除加密
			</li>
			<li v-if="checkReback()" @click="reback">
				<el-icon :size="16" color="#333">
					<FirstAidKit />
				</el-icon>
				还原
			</li>
		</ul>
	</div>
	<el-dialog v-model="contextMenuStore.isShareFile" :show-close="true" width="35%">
		<share-files />
	</el-dialog>
</template>

<script setup lang="ts">
import { useClickingStore } from "@/stores/clicking";
import { useContextMenuStore } from "@/stores/contextmenu";
import { useFileSystemStore } from "@/stores/filesystem";
import { isMobileDevice } from "@/utils/device";

import { computed, watch } from "vue";
import { useRouter } from "vue-router";
const router = useRouter();
const contextMenuStore = useContextMenuStore();
const fileSystemStore = useFileSystemStore();
const clickingStore = useClickingStore();
const zipExts = ["zip", "tar", "gz"];

const style = computed(() => {
	return contextMenuStore.contextMenuTop
		? {
			top: contextMenuStore.contextMenuTop + "px",
			left: contextMenuStore.contextMenuLeft + "px",
		}
		: {
			bottom: contextMenuStore.contextMenuBottom + "px",
			left: contextMenuStore.contextMenuLeft + "px",
		};
});

const joinKnowledge = async () => {
	await fileSystemStore.joinKnowledge(contextMenuStore.currentFile?.path);
};

// 检查加入
const checkJoin = () => {
	return (
		contextMenuStore.currentFile?.isDirectory &&
		contextMenuStore.currentFile?.knowledgeId*1 < 1
	);
};

// 检查取消加入
const checkUnJoin = () => {
	return (
		contextMenuStore.currentFile?.isDirectory &&
		contextMenuStore.currentFile?.knowledgeId*1 > 0
	);
};

// 对话
const dialogue = () => {
	// console.log(contextMenuStore.currentFile.knowledgeId);
	// console.log(contextMenuStore.currentFile)
	contextMenuStore.hideContextMenu()
	router.push({
		path: "/knowledgechat",
		query: {
			knowledgeId: contextMenuStore.currentFile.knowledgeId,
		},
	})
};

// 检查对话
const checkDialogue = () => {
	return (
		!contextMenuStore.currentFile?.ext &&
		contextMenuStore.currentFile?.knowledgeId !== 0
	);
};

const delFile = async () => {
	await fileSystemStore.handleDeleteFile(
		contextMenuStore.currentFile?.path
	);
	contextMenuStore.isContextMenuVisible = false;
};
const openFile = async () => {
	await fileSystemStore.openFile(contextMenuStore.currentFile);
	contextMenuStore.isContextMenuVisible = false;
};
const checkDelete = () => {
	if (contextMenuStore.currentFile?.ext == "exe") {
		return false;
	}
	return true;
};

const copyFile = async () => {
	clickingStore.copiedIcons = [contextMenuStore.currentFile?.path];
	contextMenuStore.isContextMenuVisible = false;
};
const cutFile = async () => {
	clickingStore.cutedIcons = [contextMenuStore.currentFile?.path];
	contextMenuStore.isContextMenuVisible = false;
};
const checkUnZip = () => {
	if (zipExts.includes(contextMenuStore.currentFile?.ext)) {
		return true;
	}
	return false;
};
const checkZip = () => {
	if (contextMenuStore.currentFile?.isDirectory) {
		return true;
	}
	return false;
};
const unZipFile = async () => {
	contextMenuStore.isContextMenuVisible = false;
	await fileSystemStore.handleUnzipFile(
		contextMenuStore.currentFile?.path
	);
};
const zipFile = async (ext: string) => {
	contextMenuStore.isContextMenuVisible = false;
	await fileSystemStore.handleZipFile(
		contextMenuStore.currentFile?.path,
		ext
	);
};
const checkEditor = () => {
	if (
		contextMenuStore.currentFile?.isDirectory ||
		contextMenuStore.currentFile?.ext == "exe"
	) {
		return false;
	}
	return true;
};
const editorFile = () => {
	contextMenuStore.isContextMenuVisible = false;
	if (contextMenuStore.currentFile) {
		fileSystemStore.openEditor(contextMenuStore.currentFile);
	}
};
const renameDir = async () => {
	contextMenuStore.isContextMenuVisible = false;
	await fileSystemStore.handleRenameFile(
		contextMenuStore.currentFile?.path
	);
};
const checkFavorite = () => {
	return (
		contextMenuStore.currentFile?.isFavorite === true ||
		contextMenuStore.currentFile?.isFavorite === "true"
	);
};
const addFavorite = () => {
	contextMenuStore.isContextMenuVisible = false;
	fileSystemStore.handleFavorite(contextMenuStore.currentFile?.path);
};
const checkShare = () => {
	return (
		contextMenuStore.currentFile?.isDirectory === false &&
		contextMenuStore.currentFile?.ext !== "exe" &&
		contextMenuStore.currentFile?.isShare === false
	);
};
const addShare = () => {
	contextMenuStore.isContextMenuVisible = false;
	contextMenuStore.isShareFile = true;
};
const checkUnShare = () => {
	return (
		contextMenuStore.currentFile?.isShare === true ||
		contextMenuStore.currentFile?.isShare === "true"
	);
};
const unShare = async () => {
	contextMenuStore.isContextMenuVisible = false;
	const postData = {
		path: contextMenuStore.currentFile?.path,
	};
	await fileSystemStore.handleShareFile(postData);
};
const checkPwd = () => {
	return (
		contextMenuStore.currentFile?.isPwd === false &&
		contextMenuStore.currentFile?.ext !== "exe"
	);
};
const addPwd = () => {
	contextMenuStore.isContextMenuVisible = false;
	fileSystemStore.handlePwdFile(contextMenuStore.currentFile?.path);
};
const checkUnpwd = () => {
	return (
		contextMenuStore.currentFile?.isPwd === true ||
		contextMenuStore.currentFile?.isPwd === "true"
	);
};
const unPwd = async () => {
	contextMenuStore.isContextMenuVisible = false;
	await fileSystemStore.handleUnpwdFile(
		contextMenuStore.currentFile?.path
	);
};
const checkDownload = () => {
	return (
		contextMenuStore.currentFile?.isDirectory === false &&
		contextMenuStore.currentFile?.ext !== "exe"
	);
};
const download = () => {
	contextMenuStore.isContextMenuVisible = false;
	fileSystemStore.downloadFile(contextMenuStore.currentFile?.path);
};
const checkReback = () => {
	return (
		contextMenuStore.currentFile?.path &&
		fileSystemStore.fs.getTopPath(contextMenuStore.currentFile.path) ===
		"B"
	);
};
const reback = () => {
	contextMenuStore.isContextMenuVisible = false;
	fileSystemStore.reStoreFiles(contextMenuStore.currentFile?.path);
};
watch(
	() => contextMenuStore.isContextMenuVisible,
	(val) => {
		if (!val) {
			contextMenuStore.contextMenuTop = 0;
		}
	}
);
</script>

<style lang="scss" scoped>
@use "@/styles/contextmenu.scss";
</style>
