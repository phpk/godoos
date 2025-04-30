<template>
	<div class="desktop-area" 
	@click.stop="handleDesktopClick($event)" 
	@dblclick.prevent="handleDesktoDubleClick()"
	@contextmenu.prevent="hanldeContextMenu($event)" 
	@drop.stop="handleDrop($event)" 
	@dragleave.prevent="dragFilesStore.handleDragLeave()" 
	ref="desktopAreaRef">
		<div :class="props.showClass" class="no-select">
			<div v-for="icon in props.fileList" :key="icon.id" class="desktop-icon"
				@dblclick.prevent="handleIconDubleClick($event, icon)" 
				@click.prevent="handleIconClick($event, icon)"
				draggable="true" @contextmenu.stop.prevent="
					contextMenuStore.showContextMenu(icon, $event)
					" :class="{
					selected: clickingStore.clickedIcons.includes(icon.path),
				}" 
				@dragenter.prevent="dragFilesStore.handleDragEnter(icon.path)"
				@dragover.prevent
				@dragstart.stop="dragFilesStore.startDrag($event)" 
				@dragleave="dragFilesStore.handleDragLeave()"
				@drop.prevent="dragFilesStore.dragFileToDrop($event, icon.path)" 
				:data-id="icon.id">
				<!-- <FileIcon :file="icon" /> -->
				<div class="icon-container">
					<icon :name="dealIcon(icon)" :size="36" />
					<icon v-if="icon.isPwd" class="lock-img" name="pwdbox" :size="12"></icon>
					<icon v-if="icon.isFavorite" class="fav-img" name="xingxing" :size="12"></icon>
					<icon v-if="icon.isShare" class="share-img" name="gallery" :size="12"></icon>
					<icon v-if="icon.knowledgeId > 0" class="ln-img" name="ink" :size="12"></icon>
				</div>
				<span class="icon-name">{{ dealSystemName(icon.title) }}</span>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { useClickingStore } from "@/stores/clicking";
import { dealSystemName } from "@/i18n";
import { dealIcon } from "@/utils/icon";
import { useContextMenuStore } from "@/stores/contextmenu";
import { useDesktopStore } from "@/stores/desktop";
import { useDragFilesStore } from "@/stores/dragfiles";
import { useFileSystemStore } from "@/stores/filesystem";
import { onMounted, onUnmounted, ref, watch } from "vue";
import { useRouter } from "vue-router";

import { eventBus } from "@/interfaces/event";

const router = useRouter();
const desktopStore = useDesktopStore();
const contextMenuStore = useContextMenuStore();
const fileSystemStore = useFileSystemStore();
const dragFilesStore = useDragFilesStore();
const clickingStore = useClickingStore();

const desktopAreaRef = ref<HTMLDivElement | null>(null);
const props = defineProps({
	fileList: {
		type: Array as () => any[],
		required: true,
	},
	isDesktop: {
		type: Boolean,
		default: true,
	},
	showClass: {
		type: String,
		default: "desktop-icons",
	},
	currentPath: {
		type: String,
		default: "/C/Users/Desktop",
	},
});
const emit = defineEmits(["navigateTo", "refeshList"]);
// 当前路径
const currentPath = ref(props.currentPath);

// function checkSelectedIcons(path: any) {
//   return clickingStore.clickedIcons.includes(path);
// }

function handleIconDubleClick(event: MouseEvent, item: any) {
	event.preventDefault();
	contextMenuStore.hideContextMenu(); // 隐藏上下文菜单
	clickingStore.clickedIcons = [];
	if (item.isDirectory) {
		if (props.isDesktop) {
			router.push({ path: "/computer", query: { path: item.path } });
		} else {
			currentPath.value = item.path;
			emit("navigateTo", item.path);
		}
	} else {
		//console.log(item)
		fileSystemStore.openFile(item);
	}
}

function handleIconClick(event: MouseEvent, icon: any) {
	clickingStore.handleClick(event, icon);
}

function handleDesktopClick(event: MouseEvent) {
	event.preventDefault();

	fileSystemStore.initFolder(currentPath.value);
	contextMenuStore.hideContextMenu(); // 隐藏上下文菜单
	clickingStore.resetClickedIcons(event);
}

function handleDesktoDubleClick() {
	contextMenuStore.hideContextMenu(); // 隐藏上下文菜单
	//clickingStore.clickedIcons = [];
}
function hanldeContextMenu(event: MouseEvent) {
	fileSystemStore.initFolder(currentPath.value);
	contextMenuStore.showDesktopContextMenu(event, currentPath.value);
}

function handleDrop(event: DragEvent) {
	//event.preventDefault(); // 阻止默认行为
	dragFilesStore.dragFileToDrop(event, currentPath.value);
}

const handleRefreshDesktop = () => {
	contextMenuStore.hideContextMenu();
	clickingStore.clickedIcons = [];
	if (props.isDesktop) {
		desktopStore.initDesktop();
	} else {
		emit("refeshList");
	}
};

onMounted(async () => {
	currentPath.value = props.currentPath;
	fileSystemStore.initFolder(currentPath.value);
	eventBus.on("refreshDesktop", handleRefreshDesktop);
	clickingStore.addEvents(currentPath.value, props.fileList);
});
watch(
	() => props.currentPath,
	(newPath) => {
		currentPath.value = newPath;
		clickingStore.addEvents(currentPath.value, props.fileList);
	}
);
watch(
	() => props.fileList,
	(newFileList) => {
		clickingStore.addEvents(currentPath.value, newFileList);
	}
);
onUnmounted(() => {
	clickingStore.removeEvents();
});
</script>

<style lang="scss" scoped>
@use "@/styles/desktopicons.scss";
@use "@/styles/contextmenu.scss";
</style>
