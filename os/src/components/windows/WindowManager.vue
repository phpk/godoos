<template>
	<div class="window-manager">
		<template v-for="win in store.windows" :key="win.id">
			<div class="window" :class="{
				'window-active': win.id === activeWindowId,
				'window-maximized': win.isMaximized,
				'window-dragging': win.isDragging,
			}" v-show="win.isVisible" :style="{
				left: win.isMaximized ? '0px' : win.position.x + 'px',
				top: win.isMaximized ? '0px' : win.position.y + 'px',
				width: win.isMaximized ? '100%' : win.size.width + 'px',
				height: win.isMaximized ? '100%' : win.size.height + 'px',
				zIndex: win.zIndex,
			}" @mousedown="activateWindow(win.id)" :data-id="win.id">
				<div class="window-titlebar" @mousedown="startDrag($event, win)" @mousemove="handleMouseMove($event)"
					@mouseup="handleMouseUp">
					<!-- <div class="window-titlebar-back">←</div> -->
					<el-icon v-if="isMobileDevice()" class="back-icon" @click="closeWindow(win.id)">
						<ArrowLeftBold />
					</el-icon>

					<div class="title-area">
						<icon :name="win.icon" size="20" />
						<div class="window-title">
							{{ t(win.title) }}
						</div>
					</div>
					<WindowFiles :win="win" v-if="win.props?.isFile" @updateWin="updateWindow" />
					<div class="window-controls">
						<el-icon @click="minimizeWindow(win.id)">
							<Minus />
						</el-icon>
						<el-icon title="最大化" v-if="!win.isMaximized" @click="toggleMaximize(win.id)">
							<FullScreen />
						</el-icon>
						<el-icon title="还原" v-else @click="toggleMaximize(win.id)">
							<CopyDocument />
						</el-icon>
						<el-icon title="关闭" @click="closeWindow(win.id)">
							<Close />
						</el-icon>
					</div>
				</div>
				<!-- <WindowTabs :win="win" v-if="win.props?.isFile" /> -->
				<div class="window-content" :class="{ 'sidebar-hidden': !win.props?.showFiles }">
					<div class="left-sidebar" v-if="win.props?.showFiles">
						<!-- 左侧栏内容 -->
						<WindowLeft v-if="win.props?.showFiles" :win="win" />
					</div>
					<div class="right-content">
						<!-- 右侧栏内容 -->
						<IframeApp v-if="checkIfIframe(win)" :win="win" />
						<component v-else :is="win.component" :win="win" />
					</div>
				</div>
				<!-- <div class="window-content">
					<WindowLeft v-if="win.props?.showFiles" :win="win" />
					<IframeApp v-if="typeof win.component === 'string'" :win="win" />
					<component v-else :is="win.component" :win="win" />
				</div> -->
				<div class="window-resize-handle" @mousedown="startResize($event, win)"
					@mousemove="handlResizingMove($event)" @mouseup="handleMouseUp"></div>
			</div>
		</template>
	</div>
</template>

<script setup lang="ts">
import { t } from "@/i18n";
import { useWindowStore } from "@/stores/window";
import { onMounted, onUnmounted, ref } from "vue";
import IframeApp from "./IframeApp.vue";
import { isMobileDevice } from "@/utils/device";
const store = useWindowStore();
const activeWindowId = ref<string | null>(null);
const isDragging = ref(false);
const isResizing = ref(false);
const dragOffset = ref({ x: 0, y: 0 });
const initialSize = ref({ width: 0, height: 0 });
const initialPosition = ref({ x: 0, y: 0 });

function activateWindow(id: string) {
	activeWindowId.value = id;
	store.bringToFront(id);
}

function startDrag(event: MouseEvent, window: any) {
	if (
		event.target instanceof HTMLElement &&
		event.target.closest(".window-controls")
	) {
		return;
	}

	isDragging.value = true;
	const rect = (
		event.currentTarget as HTMLElement
	).getBoundingClientRect();
	dragOffset.value = {
		x: event.clientX - rect.left,
		y: event.clientY - rect.top,
	};
	activateWindow(window.id);
	store.setDragging(window.id, true);
}

function startResize(event: MouseEvent, window: any) {
	isResizing.value = true;
	initialSize.value = { ...window.size };
	initialPosition.value = {
		x: event.clientX,
		y: event.clientY,
	};
	activateWindow(window.id);
}

function handleMouseMove(event: MouseEvent) {
	if (isDragging.value && activeWindowId.value !== null) {
		const x = event.clientX - dragOffset.value.x;
		const y = event.clientY - dragOffset.value.y;
		store.updateWindowPosition(activeWindowId.value, { x, y });
	}

	if (isResizing.value && activeWindowId.value !== null) {
		const deltaX = event.clientX - initialPosition.value.x;
		const deltaY = event.clientY - initialPosition.value.y;
		store.updateWindowSize(activeWindowId.value, {
			width: Math.max(400, initialSize.value.width + deltaX),
			height: Math.max(300, initialSize.value.height + deltaY),
		});
	}
}
function handlResizingMove(event: MouseEvent) {
	if (isResizing.value && activeWindowId.value !== null) {
		const deltaX = event.clientX - initialPosition.value.x;
		const deltaY = event.clientY - initialPosition.value.y;
		store.updateWindowSize(activeWindowId.value, {
			width: Math.max(400, initialSize.value.width + deltaX),
			height: Math.max(300, initialSize.value.height + deltaY),
		});
	}
}
function handleMouseUp() {
	isDragging.value = false;
	isResizing.value = false;
	if (activeWindowId.value !== null) {
		// 移除拖动时的类
		store.setDragging(activeWindowId.value, false);
	}
}

function minimizeWindow(id: string) {
	//store.minimizeWindow(id)
	const windowElement = document.querySelector(
		`.window[data-id="${id}"]`
	);
	if (windowElement) {
		windowElement.classList.add("minimizing");
		windowElement.addEventListener(
			"animationend",
			() => {
				windowElement.classList.remove("minimizing");
				store.minimizeWindow(id);
			},
			{ once: true }
		);
	}
}

function toggleMaximize(id: string) {
	store.toggleMaximizeWindow(id);
}

function closeWindow(id: string) {
	store.closeWindow(id);
}

function updateWindow(updatedWin: any) {
	store.updateWindow(updatedWin);
}
function checkIfIframe(win: any) {
	return typeof win.component === "string"
}
onMounted(() => {
	window.addEventListener("mousemove", handleMouseMove);
	window.addEventListener("mouseup", handleMouseUp);
});

onUnmounted(() => {
	window.removeEventListener("mousemove", handleMouseMove);
	window.removeEventListener("mouseup", handleMouseUp);
});
</script>

<style lang="scss" scoped>
@use "@/styles/windows.scss";
</style>
