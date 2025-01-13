<template>
	<template v-if="mode === 'detail'">
		<div class="file-item file-bar mode-detail">
			<div class="file-item_img"></div>
			<div class="file-item_title"></div>
			<div class="file-item_type">
				<span>{{ t("size") }}</span>
			</div>
			<div class="file-item_type">
				<span>{{ t("creation.time") }}</span>
			</div>
			<div class="file-item_type">
				<span>{{ t("modification.time") }}</span>
			</div>
			<div class="file-item_type">
				<span>{{ t("permission") }}</span>
			</div>
		</div>
	</template>
	<div draggable="true" class="file-item" :class="{
		chosen: chosenIndexs.includes(index),
		'no-chosen': !chosenIndexs.includes(index),
		'mode-icon': mode === 'icon',
		'mode-list': mode === 'list',
		'mode-big': mode === 'big',
		'mode-middle': mode === 'middle',
		'mode-detail': mode === 'detail',
		'drag-over': hoverIndex === index,
	}" :style="{
		'--theme-color': theme === 'light' ? '#ffffff6b' : '#3bdbff3d',
	}" v-for="(item, index) in fileList" :key="item.path" @dblclick="handleOnOpen(item)"
		@touchstart.passive="doubleTouch($event, item)"
		@contextmenu.stop.prevent="handleRightClick($event, item, index)" @drop="hadnleDrop($event, item.path)"
		@dragenter.prevent="handleDragEnter(index)" @dragover.prevent @dragleave="handleDragLeave()"
		@dragstart.stop="startDragApp($event, item)" @click="handleClick(index, item)" @mousedown.stop :ref="(ref: any) => {
			if (ref) {
				appPositions[index] = markRaw(ref as Element);
			}
		}
			">
		<!-- {{item.name}} -->
		<div class="file-item_img">
			<FileIcon :file="item" />
		</div>
		<span v-if="editIndex !== index" class="file-item_title">
			{{ getName(item) }}
		</span>
		<textarea autofocus draggable="false" @dragover.stop @dragstart.stop @dragenter.stop @mousedown.stop
			@dblclick.stop @click.stop @keydown.enter.prevent="onEditNameEnd" v-if="editIndex === index"
			class="file-item_title file-item_editing" v-model="editName"></textarea>
		<template v-if="mode === 'detail'">
			<div class="file-item_type">
				<span>{{ item.isDirectory ? "-" : dealSize(item.size) }}</span>
			</div>
			<div class="file-item_type">
				<span>{{ item.birthtime.toLocaleString() }}</span>
			</div>
			<div class="file-item_type">
				<span>{{ item.mtime.toLocaleString() }}</span>
			</div>
			<div class="file-item_type">
				<span>{{ item.mode?.toString?.(8) || "unknow" }}</span>
			</div>
		</template>
	</div>
</template>
<script lang="ts" setup>
import { useAppMenu } from "@/hook/useAppMenu";
import { useContextMenu } from "@/hook/useContextMenu.ts";
import { useFileDrag } from "@/hook/useFileDrag";
import { Rect } from "@/hook/useRectChosen";
import { dealSystemName, t } from "@/i18n";
import { useChooseStore } from "@/stores/choose";
import { getSystemKey } from "@/system/config";
import { emitEvent, mountEvent } from "@/system/event";
import { addKnowledge } from "@/hook/useAi";
import {
	basename,
	BrowserWindow,
	Notify,
	OsFileWithoutContent,
	useSystem,
} from "@/system/index.ts";
import { Menu } from "@/system/menu/Menu";
import { throttle } from "@/util/debounce";
import { dealSize } from "@/util/file";
import { markRaw, onMounted, ref } from "vue";
import { notifyError, notifySuccess } from "@/util/msg";
const { openPropsWindow, copyFile, createLink, deleteFile } =
	useContextMenu();
const sys = useSystem();
const { startDrag, folderDrop } = useFileDrag(sys);
const choose = useChooseStore();
const props = defineProps({
	onChosen: {
		type: Function,
		required: true,
	},
	onOpen: {
		type: Function,
		default: () => {
			//
		},
	},
	onRefresh: {
		type: Function,
		default: () => {
			//
		},
	},
	fileList: {
		type: Array<OsFileWithoutContent>,
		default: () => [],
	},
	theme: {
		type: String || Object,
		default: "light",
	},
	mode: {
		type: String,
		default: "icon",
	},
});

function getName(item: any) {
	const name = dealSystemName(basename(item.path));
	// console.log(name)
	// console.log(item.path)
	if (name.endsWith(".exe")) {
		return t(name.replace(".exe", ""));
	} else {
		return name;
	}
}

function handleOnOpen(item: OsFileWithoutContent) {
	// props.onOpen(item);
	// emitEvent('desktop.app.open');
	chosenIndexs.value = [];
	if (choose.ifShow && !item.isDirectory) {
		choose.path.push(item.path);
		choose.close();
	} else {
		// console.log(' file list:',props.fileList);

		props.onOpen(item);
		emitEvent("desktop.app.open");
	}
}
function hadnleDrop(mouse: DragEvent, path: string) {
	hoverIndex.value = -1;
	folderDrop(mouse, path);
	chosenIndexs.value = [];
}
let expired: number | null = null;
function doubleTouch(e: TouchEvent, item: OsFileWithoutContent) {
	if (e.touches.length === 1) {
		if (!expired) {
			expired = e.timeStamp + 400;
		} else if (e.timeStamp <= expired) {
			// remove the default of this event ( Zoom )
			handleOnOpen(item);
			e.preventDefault();
			// then reset the variable for other "double Touches" event
			expired = null;
		} else {
			// if the second touch was expired, make it as it's the first
			expired = e.timeStamp + 400;
		}
	}
}

const editIndex = ref<number>(-1);
const editName = ref<string>("");

async function onEditNameEnd() {
	const editEndName = editName.value.trim();
	if (editIndex.value >= 0) {
		const currentItem = props.fileList[editIndex.value];
		const success = await renameFile(currentItem, editEndName);
		if (!success) {
			editName.value = basename(currentItem.path);
		}
	}
	editIndex.value = -1;
}

mountEvent("edit.end", () => {
	onEditNameEnd();
});

const hoverIndex = ref<number>(-1);
const appPositions = ref<Array<Element>>([]);

const chosenIndexs = ref<Array<number>>([]);
import { isMobileDevice } from "@/util/device";
const isMobile = ref<boolean>(false);
function handleClick(index: number, item: OsFileWithoutContent) {
	chosenIndexs.value = [index];
	if (isMobile.value) {
		handleOnOpen(item);
	}
}
onMounted(() => {
	isMobile.value = isMobileDevice();
	chosenIndexs.value = [];
	props.onChosen(
		throttle((rect: Rect) => {
			const tempChosen: number[] = [];
			appPositions.value.forEach((el, index) => {
				const rect2 = el.getBoundingClientRect();
				const rect2Center = {
					x: rect2.left + rect2.width / 2,
					y: rect2.top + rect2.height / 2,
				};
				if (
					rect2Center.x > rect.left &&
					rect2Center.x < rect.left + rect.width &&
					rect2Center.y > rect.top &&
					rect2Center.y < rect.top + rect.height
				) {
					tempChosen.push(index);
				}
			});
			chosenIndexs.value = tempChosen;
		}, 100)
	);
});

function startDragApp(mouse: DragEvent, item: OsFileWithoutContent) {
	if (chosenIndexs.value.length) {
		startDrag(
			mouse,
			chosenIndexs.value.map((index) => {
				return props.fileList[index];
			}),
			() => {
				chosenIndexs.value = [];
			}
		);
	} else {
		startDrag(mouse, [item], () => {
			chosenIndexs.value = [];
		});
	}
}
document.addEventListener("copy", function () {
	const files = chosenIndexs.value.map((index) => props.fileList[index]);
	if (files.length > 0) {
		copyFile(files);
		chosenIndexs.value = [];
	}
});
// document.addEventListener('keydown', function (event) {
// 	// 检测 Control + C
// 	if (event.ctrlKey && event.key === 'c') {
// 		//console.log('Control + C 被按下');
// 		copyFile(
// 			chosenIndexs.value.map(
// 				(index) => props.fileList[index]
// 			)
// 		);
// 		chosenIndexs.value = [];
// 	}
// });
function handleRightClick(
	mouse: MouseEvent,
	item: OsFileWithoutContent,
	index: number
) {
	if (chosenIndexs.value.length <= 1) {
		chosenIndexs.value = [
			props.fileList.findIndex((app) => app.path === item.path),
		];
	}
	const ext = item.name.split(".").pop();

	const zipSucess = (res: any) => {
		if (!res || res.code < 0) {
			new Notify({
				title: t("tips"),
				content: t("error"),
			});
		} else {
			props.onRefresh();
			new Notify({
				title: t("tips"),
				content: t("file.zip.success"),
			});
			if (item.parentPath == "/C/Users/Desktop") {
				sys.refershAppList();
			}
		}
	};
	// eslint-disable-next-line prefer-const
	let menuArr: any = [
		{
			label: t("open"),
			click: () => {
				chosenIndexs.value = [];
				props.onOpen(item);
			},
		},
		// {
		//   label: t('open.with'),
		//   click: () => {
		//     chosenIndexs.value = [];
		//     openWith(item);
		//   },
		// },
	];
	if (item.isDirectory) {
		if (getSystemKey("storeType") == "local") {
			menuArr.push({
				label: t("zip"),
				submenu: [
					{
						label: "zip",
						click: () => {
							sys.fs
								.zip(item.path, "zip")
								.then((res: any) => {
									zipSucess(res);
								});
						},
					},
					{
						label: "tar",
						click: () => {
							sys.fs
								.zip(item.path, "tar")
								.then((res: any) => {
									zipSucess(res);
								});
						},
					},
					{
						label: "gz",
						click: () => {
							sys.fs.zip(item.path, "gz").then((res: any) => {
								zipSucess(res);
							});
						},
					},
				],
			});
		} else {
			menuArr.push({
				label: t("zip"),
				click: () => {
					sys.fs.zip(item.path, "zip").then((res: any) => {
						zipSucess(res);
					});
				},
			});
		}
		if (!item.knowledgeId) {
			menuArr.push({
				label: "加入知识库",
				click: () => {
					console.log(item);
					addKnowledge(item.path).then((res: any) => {
						console.log(res);
						if (res.code != 0) {
							new Notify({
								title: t("tips"),
								content: res.message,
							});
						} else {
							new Notify({
								title: t("tips"),
								content: "添加成功",
							});
							props.onRefresh();
						}
					});
				},
			});
		} else {
			menuArr.push({
				label: "对话知识库",
				click: () => {
					const win = new BrowserWindow({
						title: "对话知识库",
						content: "KnowledgeChat",
						config: {
							path: item.path,
							knowledgeId: item.knowledgeId,
						},
						width: 700,
						height: 600,
						center: true,
					});
					win.show();
				},
			});
		}
	}
	if (choose.ifShow) {
		menuArr.push({
			label: "选中发送",
			click: () => {
				const paths: any = [];
				chosenIndexs.value.forEach((index) => {
					const item = props.fileList[index];
					paths.push(item.path);
				});
				if (paths.length > 0) {
					choose.path = paths;
					choose.close();
				}
				chosenIndexs.value = [];
			},
		});
	}
	// eslint-disable-next-line prefer-const
	let extMenus = useAppMenu(item, sys, props);
	if (extMenus && extMenus.length > 0) {
		// eslint-disable-next-line prefer-spread
		menuArr.push.apply(menuArr, extMenus);
	}
	if (ext != "exe") {
		const fileMenus = [
			{
				label: t("rename"),
				click: () => {
					editIndex.value = index;
					editName.value = basename(item.path);
					chosenIndexs.value = [];
				},
			},
			{
				label: t("copy"),
				click: () => {
					//if(["/","/B"].includes(item.path)) return;
					copyFile(
						chosenIndexs.value.map(
							(index) => props.fileList[index]
						)
					);
					chosenIndexs.value = [];
				},
			},
		];
		if (item.path.indexOf("/F") < 0) {
			fileMenus.push({
				label: t("delete"),
				click: async () => {
					for (let i = 0; i < chosenIndexs.value.length; i++) {
						await deleteFile(
							props.fileList[chosenIndexs.value[i]]
						);
					}
					chosenIndexs.value = [];
					props.onRefresh();
					if (item.path.indexOf("Desktop") > -1) {
						sys.refershAppList();
					}
				},
			});
			fileMenus.push({
				label: t("create.shortcut"),
				click: () => {
					createLink(item.path)?.then(() => {
						chosenIndexs.value = [];
						props.onRefresh();
					});
				},
			});
		}
		const userType = sys.getConfig("userType");
		if (userType == "member") {
			menuArr.push({
				label: "分享给...",
				click: () => {
					const win = new BrowserWindow({
						title: "分享",
						content: "ShareFiles",
						config: {
							path: item.path,
						},
						width: 450,
						height: 280,
						center: true,
					});
					win.show();
				},
			});
		}
		if (!item.isDirectory) {
			menuArr.push({
				label: "文件加密",
				click: () => {
					const win = new BrowserWindow({
						title: "文件加密",
						content: "FilePwd",
						config: {
							path: item.path,
						},
						width: 400,
						height: 200,
						center: true,
					});
					win.show();
				},
			});
		}

		menuArr.push.apply(menuArr, fileMenus);
	}
	const sysEndMenu = [
		{
			label: t("props"),
			click: () => {
				chosenIndexs.value.forEach((index) => {
					openPropsWindow(props.fileList[index].path);
					chosenIndexs.value = [];
				});
			},
		},
	];
	// eslint-disable-next-line prefer-spread
	menuArr.push.apply(menuArr, sysEndMenu);
	//console.log(item)

	//console.log(ext)

	Menu.buildFromTemplate(menuArr).popup(mouse);
}

function handleDragEnter(index: number) {
	hoverIndex.value = index;
}

function handleDragLeave() {
	hoverIndex.value = -1;
}

// function dealtName(name: string) {
//   return name;
// }

function getPathSeparator(path: string): string {
	return path.indexOf("/") === -1 ? "\\" : "/";
}

function buildNewPath(editpath: string, editEndName: string): string {
	const sp = getPathSeparator(editpath);
	const pathParts = editpath.split(sp);
	pathParts.pop();
	pathParts.push(editEndName);
	return pathParts.join(sp);
}

async function renameFile(
	currentItem: OsFileWithoutContent,
	newName: string
) {
	const currentName = basename(currentItem.path);
	if (newName && newName !== currentName) {
		const editpath = currentItem.path.toString();
		const newPath = buildNewPath(editpath, newName);

		if (await sys?.fs.exists(newPath)) {
			notifyError("文件名已存在，请选择其他名称。");
			return false;
		}

		try {
			await sys?.fs.rename(editpath, newPath);
			props.onRefresh();
			if (newPath.includes("Desktop")) {
				sys.refershAppList();
			}
			notifySuccess("重命名成功");
			return true;
		} catch (error) {
			console.error("重命名失败:", error);
			notifyError("重命名失败，请重试。");
			return false;
		}
	}
	return false;
}
</script>
<style lang="scss" scoped>
.file-item {
	position: relative;
	display: flex;
	flex-direction: column;
	justify-content: flex-start;
	align-items: center;
	width: var(--desk-item-size);
	height: var(--desk-item-size);
	font-size: var(--ui-font-size);
	color: var(--icon-title-color);
	padding-top: 4px;
	border: 1px solid transparent;
	margin: 6px;

	.file-item_img {
		width: 60%;
		height: 60%;
		pointer-events: none;
	}

	.file-item_type {
		display: none;
	}

	.file-item_title {
		pointer-events: none;
	}

	.file-item_editing {
		display: inline-block !important;
		outline: none;
		pointer-events: all;
		padding: 6px 10px;
		margin: 0;
		min-width: 120px;
		max-width: 100%;
		height: auto !important;
		width: auto !important;
		resize: none;
		border-radius: 6px;
		border: 1px solid var(--border-color, #ccc);
		background-color: var(--input-bg-color, #f9f9f9);
		color: var(--input-text-color, #333);
		box-shadow: 0 3px 6px rgba(0, 0, 0, 0.15);
		transition: border-color 0.3s, box-shadow 0.3s,
			background-color 0.3s;
		overflow-x: auto;
		white-space: nowrap;
	}

	.file-item_editing:focus {
		border-color: var(--focus-border-color, #0056b3);
		box-shadow: 0 3px 8px rgba(0, 86, 179, 0.3);
		background-color: var(--focus-bg-color, #e6f7ff);
	}
}

.file-item:hover {
	background-color: #b1f1ff4c;
}

.chosen {
	border: 1px dashed #3bdbff3d;
	// background-color: #ffffff6b;
	background-color: var(--theme-color);

	.file-item_title {
		overflow: hidden;
		text-overflow: ellipsis;
		word-break: break-all;
		display: -webkit-box;
		-webkit-box-orient: vertical;
		-webkit-line-clamp: 2;
	}
}

.no-chosen {
	.file-item_title {
		overflow: hidden;
		text-overflow: ellipsis;
		word-break: break-all;
		display: -webkit-box;
		-webkit-box-orient: vertical;
		-webkit-line-clamp: 2;
	}
}

.drag-over {
	border: 1px dashed #3bdbff3d;
	// background-color: #ffffff6b;
	background-color: var(--theme-color);
}

.mode-icon {
	.file-item_img {
		width: 60%;
		height: calc(0.6 * var(--desk-item-size));
		margin: 0px auto;
		user-select: none;
		flex-shrink: 0;
	}

	.file-item_title {
		// color: var(--color-ui-desk-item-title);
		// height: calc(0.4 * var(--desk-item-size));
		// display: block;
		text-align: center;
		word-break: break-all;
		flex-grow: 0;
	}
}

.mode-list {
	display: flex;
	flex-direction: row;
	justify-content: flex-start;
	height: var(--menulist-item-height);
	width: var(--menulist-width);

	.file-item_img {
		width: var(--menulist-item-height);
		height: calc(0.6 * var(--menulist-item-height));

		flex-shrink: 0;
		user-select: none;
	}

	.file-item_title {
		height: min-content;
		word-break: break-all;
	}
}

.mode-icon {
	width: var(--desk-item-size);
	height: var(--desk-item-size);
}

.mode-big {
	width: calc(var(--desk-item-size) * 2.5);
	height: calc(var(--desk-item-size) * 2.5);
}

.mode-middle {
	width: calc(var(--desk-item-size) * 1.5);
	height: calc(var(--desk-item-size) * 1.5);
}

.mode-detail {
	display: flex;
	flex-direction: row;
	justify-content: flex-start;
	height: var(--menulist-item-height);
	width: 100%;
	margin: 2px;

	.file-item_img {
		width: 30px;
	}

	.file-item_title {
		width: 40%;
		display: flex;
		align-items: center;
		word-break: break-all;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.file-item_type {
		display: block;
		color: var(--color-dark-hover);
		width: 20%;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}
}

.file-bar:hover {
	background-color: unset;
	user-select: none;
}

@media screen and (max-width: 768px) {
	.file-item {
		height: vh(70);
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		overflow: hidden;

		&:hover {
			background-color: transparent;
		}

		.icon {
			width: vw(40);
			height: vh(40);
		}

		.title {
			font-size: vw(12);
			text-align: center;
		}
	}

	.chosen {
		background-color: transparent;
		border: 1px solid transparent;
	}
}
</style>
