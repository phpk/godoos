<template>
	<div class="file-explorer-container">
		<div class="command-bar">
			<div class="command-group">
				<el-icon @click="goBack" :class="{ disabled: !canGoBack }">
					<Back />
				</el-icon>
				<el-icon @click="goForward" :class="{ disabled: !canGoBack }">
					<Right />
				</el-icon>
				<el-icon @click="goUp">
					<Top />
				</el-icon>
			</div>

			<div class="path" v-show="pathState == 'pendding'" @click="startInput">
				<span class="path-p" v-for="(p, index) in routePaths" :key="index" @click.stop="handlePathClick(index)">
					{{ p }}
					<el-icon>
						<ArrowRight />
					</el-icon>
				</span>
			</div>
			<el-input ref="myinput" v-show="pathState == 'inputing'" v-model="currentPath" class="path-input"
				placeholder="输入路径" @focusout="endInput" @keyup.enter="endInput">
				<template #prefix>
					<el-icon>
						<Folder />
					</el-icon>
				</template>
			</el-input>

			<el-input v-model="searchQuery" class="search-input" placeholder="搜索文件..." @keyup.enter="searchFile">
				<template #prefix>
					<el-icon @click="searchFile">
						<Search />
					</el-icon>
				</template>
			</el-input>

			<el-button @click="toggleViewMode" class="view-mode-button">
				<el-icon>
					<component :is="viewMode === 'grid' ? 'Grid' : 'List'" />
				</el-icon>
			</el-button>
		</div>
		<div class="main-content no-select">
			<div class="sidebar">
				<div class="navigation-pane">
					<div class="quick-access-section">
						<div class="section-header" @click="quickAccessShow = !quickAccessShow">
							<el-icon v-if="quickAccessShow" class="header-icon" style="margin-right: 10px">
								<ArrowDownBold />
							</el-icon>
							<el-icon v-else class="header-icon inactive" style="margin-right: 10px">
								<ArrowRightBold />
							</el-icon>
							<icon name="xingxing" :size="15" style="margin-right: 10px"></icon>
							快速访问
						</div>
						<div v-show="quickAccessShow" v-for="item in quickAccess" :key="item.id" class="navigation-item"
							:class="{ active: item.path === currentPath }" @click="navigateTo(item.path)">
							<icon :name="item.icon" :size="20" class="navigation-icon">
							</icon>
							<span class="navigation-text">{{ item.name }}</span>
							<icon name="tuding" :size="10"></icon>
						</div>
					</div>
				</div>
				<FileTreeSection :navigateTo="navigateTo" :currentPath="currentPath" />
			</div>
			<div class="file-list-area">
				<desktop-icons :isDesktop="false" :fileList="fileList"
					:showClass="viewMode === 'grid' ? 'file-grid' : 'file-list'" @navigateTo="navigateTo"
					@refeshList="refeshList" :currentPath="currentPath" :key="winId" />
			</div>
		</div>
		<SaveFile v-if="fileSystemStore.choose.ifShow" :winId="winId" :currentPath="currentPath" />
	</div>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, ref } from "vue";
//import { useDraggingStore } from '@/stores/dragging'
import { useFileSystemStore } from "@/stores/filesystem";
//import { eventBus } from '@/interfaces/event';

import { Folder, Search } from "@element-plus/icons-vue";
const currentPath = ref("/");
const history = ref<string[]>(["/"]);
const historyIndex = ref(0);
const searchQuery = ref("");
const viewMode = ref<"grid" | "list">("grid");
//const draggingStore = useDraggingStore();
const fileList: any = ref([]);
const fileSystemStore = useFileSystemStore();
const quickAccess = [
	{ id: 1, name: "桌面", path: "/C/Users/Desktop", icon: "zhuomian" },
	{ id: 2, name: "文档", path: "/C/Users/Documents", icon: "wendang" },
	{ id: 3, name: "下载", path: "/C/Users/Downloads", icon: "xiazai" },
	{ id: 4, name: "图片", path: "/C/Users/Pictures", icon: "tupian1" },
	{ id: 5, name: "分享", path: "/F", icon: "ink" },
	{ id: 6, name: "收藏", path: "/G", icon: "favorite" },
	{ id: 7, name: "知识库", path: "/H", icon: "aiknow" },
];
const quickAccessShow = ref(true);

const pathState = ref("pendding");
const myinput = ref();
const routePaths = computed(() => {
	const sp = currentPath.value.charAt(0);
	const arr = currentPath.value.split(sp);
	arr[0] = "我的电脑";
	return arr.filter((item) => item !== "");
});
const handlePathClick = (index: number) => {
	const sp = currentPath.value.charAt(0);
	const arr = currentPath.value.split(sp);
	const newPath = arr.slice(0, index + 1).join(sp) || "/";
	navigateTo(newPath);
};
const startInput = () => {
	pathState.value = "inputing";
	nextTick(() => {
		myinput.value?.focus();
		myinput.value?.select(0, currentPath.value.length);
	});
};

const endInput = () => {
	pathState.value = "pendding";
};
const searchFile = async () => {
	if (searchQuery.value != "") {
		fileList.value = await fileSystemStore.handleSerach(
			currentPath.value,
			searchQuery.value
		);
	} else {
		fileList.value = await fileSystemStore.getFilesInPath(
			currentPath.value
		);
	}
};

const props = defineProps({
	win: {
		type: Object,
		required: true,
	},
});
const winId = ref(props.win?.id);
const getFilesInPath = async (path: string) => {
	currentPath.value = path;
	fileList.value = await fileSystemStore.getFilesInPath(path);
	console.log(fileList.value);
};
const refeshList = () => {
	getFilesInPath(currentPath.value);
};

onMounted(async () => {
	const path = (props.win.props.path as string) || "/";
	await getFilesInPath(path);
});

const canGoBack = computed(() => historyIndex.value > 0);
const canGoForward = computed(
	() => historyIndex.value < history.value.length - 1
);

function navigateTo(path: string) {
	getFilesInPath(path).then(() => {
		currentPath.value = path;
		history.value = history.value.slice(0, historyIndex.value + 1);
		history.value.push(path);
		historyIndex.value = history.value.length - 1;
	});
}

function goBack() {
	if (canGoBack.value) {
		historyIndex.value--;
		currentPath.value = history.value[historyIndex.value];
		navigateTo(currentPath.value);
	}
}

function goForward() {
	if (canGoForward.value) {
		historyIndex.value++;
		currentPath.value = history.value[historyIndex.value];
		navigateTo(currentPath.value);
	}
}

function goUp() {
	navigateTo("/");
}

function toggleViewMode() {
	viewMode.value = viewMode.value === "grid" ? "list" : "grid";
}
</script>

<style lang="scss" scoped>
@use "@/styles/computer.scss";
</style>
