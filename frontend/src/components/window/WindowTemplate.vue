<template>
	<div class="wintmp_outer dragwin" :class="{
		topwin: istop,
		max: windowInfo.state == WindowStateEnum.maximize,
		min: windowInfo.state == WindowStateEnum.minimize,
		fullscreen: windowInfo.state == WindowStateEnum.fullscreen,
		noframe: !windowInfo.frame,
		disable: windowInfo.disable,
	}" :style="customerStyle" @touchstart.passive="onFocus" @mousedown="onFocus" ref="$win_outer" v-dragable>
		<!-- 窗口标题栏  -->
		<div class="wintmp_uper" @contextmenu.prevent v-if="!isMobileDevice()">
			<MenuBar :browser-window="browserWindow"></MenuBar>
		</div>
		<van-nav-bar v-else left-text="返回" left-arrow @click-left="onClickLeft" fixed>
			<template #title>
				<div class="title">
					<FileIcon :icon="browserWindow.windowInfo.icon" />
					{{ browserWindow.windowInfo.title }}
				</div>
			</template>
		</van-nav-bar>

		<div class="wintmp_main"
			:class="{ resizeing: resizemode != 'null', 'saveFileMain': browserWindow.windowInfo.footer }"
			@mousedown.stop="predown" @touchstart.stop.passive="predown" @contextmenu.stop.prevent>
			<div class="content-mask" v-if="!istop && typeof browserWindow.content === 'string'"></div>
			<WindowInner :win="browserWindow"></WindowInner>
		</div>
		<!-- 使用 v-for 生成拖拽边界 -->
		<div v-for="border in dragBorders" :key="border.type" :class="[
			border.class,
			'win_drag_border',
			{ isChoseMode: resizemode == border.type },
			border.cursorClass,
		]" v-if="resizable" draggable="false" @mousedown.stop.prevent="startScale($event, border.type)"
			@touchstart.stop.passive="startScale($event, border.type)"></div>
		<div class="wintmp_footer" v-if="browserWindow.windowInfo.footer">
			<MenuFooter :browser-window="browserWindow" @translateSavePath="translateSavePath"></MenuFooter>
		</div>
	</div>
</template>
<script lang="ts" setup>
import { useSystem } from "@/system";
import { emitEvent } from "@/system/event"; ""
import {
	BrowserWindow,
	WindowStateEnum,
} from "@/system/window/BrowserWindow";
import { ScaleElement } from "@/system/window/dom/ScaleElement";
import { vDragable } from "@/system/window/MakeDragable";
import {
	computed,
	onMounted,
	onUnmounted,
	provide,
	ref,
	UnwrapNestedRefs,
} from "vue";
import { useChooseStore } from "@/stores/choose";
import eventBus from '@/system/event/eventBus'
import { isMobileDevice } from "@/util/device";

const onClickLeft = () => {
	if (choose.isExist(props.browserWindow.windowInfo.componentID)) {
		choose.closeSaveFile(props.browserWindow.windowInfo.componentID)
	}
	props.browserWindow.destroy();
}

const sys = useSystem();
const props = defineProps<{
	browserWindow: UnwrapNestedRefs<BrowserWindow>;
}>();


const browserWindow = props.browserWindow;
if (isMobileDevice()) browserWindow.maximize()
const windowInfo = browserWindow.windowInfo;
console.log(windowInfo, 'windowInfo')
// const temp = reactive(browserWindow)
provide("browserWindow", browserWindow);
provide("system", sys);
const choose = useChooseStore()

function translateSavePath(path: string, name?: string) {
	if (browserWindow.windowInfo.footer) {
		const pos = choose.saveFileContent.findIndex((item: any) => {
			return item.componentID == browserWindow.windowInfo.componentID
		})
		if (pos == -1) return
		if (path && path !== '') {
			choose.saveFileContent[pos].filePath = path
		} else if (name && name !== '') {
			choose.saveFileContent[pos].fileName = name
			eventBus.emit('saveFile', choose.saveFileContent[pos])
		}
	}
}
provide('translateSavePath', translateSavePath)

function predown() {
	browserWindow.moveTop();
	emitEvent("window.content.click", browserWindow);
}

const customerStyle = ref<NonNullable<unknown>>({});

function onFocus(e: MouseEvent | TouchEvent): void {
	browserWindow?.moveTop();
	if (windowInfo.state === WindowStateEnum.maximize) {
		if (e instanceof MouseEvent) {
			e.preventDefault();
			e.stopPropagation();
		}
	}
}

const istop = computed(() => windowInfo.istop);

onMounted(() => {
	customerStyle.value = {
		width: computed(() => windowInfo.width + "px"),
		height: computed(() => windowInfo.height + "px"),
		left: computed(() => windowInfo.x + "px"),
		top: computed(() => windowInfo.y + "px"),
		zIndex: computed(() => {
			if (windowInfo.alwaysOnTop) {
				return 9999;
			}
			return windowInfo.zindex;
		}),
		backgroundColor: computed(() => windowInfo.backgroundColor),
	};
});

const resizable = ref(windowInfo.resizable);
const resizemode = ref("null");
let scaleAble: ScaleElement;

onMounted(() => {
	scaleAble = new ScaleElement(
		resizemode,
		windowInfo.width,
		windowInfo.height,
		windowInfo.x,
		windowInfo.y
	);
	scaleAble.onResize(
		(width: number, height: number, x: number, y: number) => {
			windowInfo.width = width || windowInfo.width;
			windowInfo.height = height || windowInfo.height;
			windowInfo.x = x || windowInfo.x;
			windowInfo.y = y || windowInfo.y;
			browserWindow.emit(
				"resize",
				windowInfo.width,
				windowInfo.height
			);
		}
	);
});

function startScale(e: MouseEvent | TouchEvent, dire: string) {
	console.log(e);
	if (windowInfo.disable) {
		return;
	}
	scaleAble?.startScale(
		e,
		dire,
		windowInfo.x,
		windowInfo.y,
		windowInfo.width,
		windowInfo.height
	);
}

onUnmounted(() => {
	scaleAble.unMount();
});

const dragBorders = [
	{ type: "r", class: "right_border", cursorClass: "ew-resize" },
	{ type: "b", class: "bottom_border", cursorClass: "ns-resize" },
	{ type: "l", class: "left_border", cursorClass: "ew-resize" },
	{ type: "t", class: "top_border", cursorClass: "ns-resize" },
	{
		type: "rb",
		class: "right_bottom_border",
		cursorClass: "nwse-resize",
	},
	{ type: "lb", class: "left_bottom_border", cursorClass: "nesw-resize" },
	{ type: "lt", class: "left_top_border", cursorClass: "nwse-resize" },
	{ type: "rt", class: "right_top_border", cursorClass: "nesw-resize" },
];
</script>
<style>
.dragwin {
	position: absolute;
	width: 100%;
	height: 100%;
}
</style>
<style scoped lang="scss">
.wintmp_outer {
	position: absolute;
	padding: 0;
	margin: 0;
	// left: 0;
	// top: 0;
	// min-width: 800px;
	width: max-content;
	height: max-content;
	// min-height: 650px;
	border-radius: 10px;
	overflow: hidden;
	background-color: #fff;
	border: var(--window-border);
	display: flex;
	flex-direction: column;

	// box-shadow: var(--window-box-shadow);
	// border-radius: var(--window-border-radius);
	.wintmp_main {
		position: relative;
		width: 100%;
		height: 100%;
		// background-color: rgb(255, 255, 255);
		overflow: hidden;
		contain: content;
	}

	.wintmp_footer {
		position: relative;
	}
}

.saveFileMain {
	:deep(.main) {
		height: calc(100% - 60px);
	}
}

.topwin {
	// border: 1px solid #0078d7;
	// box-shadow: var(--window-top-box-shadow);
}

.icon {
	width: 12px;
	height: 12px;
}

.max {
	position: absolute;
	left: 0 !important;
	top: 0 !important;
	width: 100% !important;
	height: 100% !important;
	transition: left 0.1s ease-in-out, top 0.1s ease-in-out,
		width 0.1s ease-in-out, height 0.1s ease-in-out;
}

.disable {

	.wintmp_footer,
	.wintmp_uper,
	.wintmp_main {
		pointer-events: none;
		user-select: none;
		box-shadow: none;
	}
}

.min {
	visibility: hidden;
	display: none;
}

.fullscreen {
	// 将声明移动到嵌套规则之上
	position: fixed;
	left: 0 !important;
	top: 0 !important;
	width: 100% !important;
	height: 100% !important;
	z-index: 205 !important;
	border: none;

	.wintmp_uper {
		display: none;
	}
}

.noframe {
	border: none;
	box-shadow: none;

	.wintmp_uper {
		display: none;
	}
}

.transparent {
	background-color: transparent;

	.wintmp_main {
		background-color: transparent;
	}

	.wintmp_uper {
		background-color: rgba(255, 255, 255, 0.774);
	}
}

.win_drag_border {
	position: absolute;
	background-color: rgba(0, 0, 0, 0);
}

.right_border {
	cursor: ew-resize;
	right: -12px;
	width: 16px;
	height: calc(100% - 4px);
}

.bottom_border {
	cursor: ns-resize;
	bottom: -12px;
	width: calc(100% - 4px);
	height: 16px;
}

.left_border {
	cursor: ew-resize;
	left: -12px;
	width: 16px;
	height: calc(100% - 4px);
}

.top_border {
	cursor: ns-resize;
	top: -12px;
	width: calc(100% - 4px);
	height: 16px;
}

.left_top_border {
	cursor: nwse-resize;
	left: -12px;
	top: -12px;
	width: 16px;
	height: 16px;
}

.right_top_border {
	cursor: nesw-resize;
	right: -12px;
	top: -12px;
	width: 16px;
	height: 16px;
}

.left_bottom_border {
	cursor: nesw-resize;
	left: -12px;
	bottom: -12px;
	width: 16px;
	height: 16px;
}

.right_bottom_border {
	cursor: nwse-resize;
	right: -12px;
	bottom: -12px;
	width: 16px;
	height: 16px;
}

.isChoseMode {
	width: 100vw;
	height: 100vh;
	position: fixed;
	left: 0;
	top: 0;
}

.resizeing {
	user-select: none;
	pointer-events: none;
}

.content-mask {
	position: absolute;
	width: 100%;
	height: 100%;
	background-color: rgba(0, 0, 0, 0);
	z-index: 100;
}

@media screen and (max-width: 768px) {
	.wintmp_outer {
		height: vh(100);
		padding-top: vh(40);
	}

	.title {
		height: vh(46);
		display: flex;
		align-items: center;

		.icon {
			width: 1em;
			height: 1em;
		}
	}
}
</style>
