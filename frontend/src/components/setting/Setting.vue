<template>
	<div
		class="window-outer"
		:class="{
			focus: focusState && currentRouter !== 'main',
		}"
	>
		<div
			class="upbar"
			v-dragable
			v-if="!isMobileDevice()"
		>
			<div class="upbar-left">
				<div
					class="back-arr"
					v-if="currentRouter !== 'main'"
					@click="back"
				>
					←
				</div>
				<div class="upbar-text">
					{{ t("setting") }}
				</div>
			</div>
			<div class="upbar-right">
				<WinUpButtonGroup
					:browser-window="browserWindow"
				></WinUpButtonGroup>
			</div>
		</div>

		<div class="settings-container">
			<aside class="sidebar">
				<ul>
					<li
						v-for="item in setList"
						:key="item.key"
						@click="openSet(item.key)"
					>
						<ElTooltip
							v-if="!isMobileDevice()"
							:content="item.title"
							placement="right"
						>
							<svg
								class="icon"
								aria-hidden="true"
							>
								<use :xlink:href="'#icon-' + item.icon"></use>
							</svg>
						</ElTooltip>
						<div v-else class="icon-container">
							<svg
								class="icon"
								aria-hidden="true"
							>
								<use :xlink:href="'#icon-' + item.icon"></use>
							</svg>
							<div class="icon-title">{{ item.title }}</div>
						</div>
					</li>
				</ul>
			</aside>
			<main class="content">
				<Transition
					name="fade"
					appear
				>
					<component
						:is="stepComponent(currentContent)"
						v-if="currentContent"
					/>
				</Transition>
			</main>
		</div>
	</div>
</template>

<script lang="ts" setup>
	import { computed, inject, ref } from "vue";

	import { t } from "@/i18n";
	import { useSystem } from "@/system";
	import { BrowserWindow } from "@/system/window/BrowserWindow";
	import { vDragable } from "@/system/window/MakeDragable";
	import { isMobileDevice } from "@/util/device";
	import { stepComponent } from "@/util/stepComponent";
	import { ElTooltip } from "element-plus";

	const browserWindow = inject<BrowserWindow>("browserWindow")!;
	const sys = useSystem();
	const currentRouter = ref(browserWindow.config?.router || "system");

	const focusState = ref(false);
	browserWindow?.on("focus", () => {
		focusState.value = true;
	});
	browserWindow?.on("blur", () => {
		focusState.value = false;
	});
	function back() {
		currentRouter.value = "main";
	}
	function openSet(key: string) {
		currentRouter.value = key;
	}
	const setList = ref([
		{
			key: "system",
			title: t("system"),
			desc: "存储、备份还原、用户角色",
			icon: "system",
			content: "SetSystem",
		},
		{
			key: "custom",
			title: "代理",
			desc: "本地代理、远程代理",
			icon: "personal",
			content: "SetCustom",
		},
		{
			key: "nas",
			title: "NAS服务",
			desc: "NAS/webdav服务",
			icon: "disk",
			content: "SetNas",
		},
		{
			key: "account",
			title: "屏幕",
			desc: "壁纸/语言/锁屏/广告",
			icon: "account",
			content: "SetAccount",
		},
		...(sys._rootState.settings ? sys._rootState.settings : []),
	]);

	const currentContent = computed(() => {
		return setList.value.find((item) => item.key === currentRouter.value)
			?.content;
	});
</script>

<style lang="scss" scoped>
	@import "@/assets/main.scss";

	.window-outer {
		background-color: white;
		width: 100%;
		height: 100%;
		border: #0076d795 1px solid;
		box-sizing: border-box;
		transition: background-color 0.1s;
		overflow: hidden;
	}

	.window-outer.focus {
		background-color: rgba(255, 255, 255, 0.704);
		backdrop-filter: blur(10px);
	}

	.upbar {
		height: 40px;
		width: 100%;
		display: flex;
		justify-content: space-between;
		align-items: center;
		user-select: none;

		.upbar-left {
			display: flex;
			align-items: center;
			font-size: 12px;
			height: 100%;

			.back-arr {
				padding: 0 6px;
				margin-right: 3px;
				font-size: 14px;
				width: 40px;
				height: 100%;
				display: flex;
				justify-content: center;
				align-items: center;
			}

			.back-arr:hover {
				color: var(--color-ui-gray);
				background-color: var(--color-dark);
			}

			.upbar-text {
				font-size: 14px;
				margin-left: 10px;
			}
		}

		.upbar-right {
			width: calc(100% - 200px);
			display: flex;
			justify-content: flex-end;
			height: 100%;
			background-color: white;
		}
	}

	.settings-container {
		display: flex;
		height: calc(100% - 40px); // 减去上方工具栏的高度
	}

	.sidebar {
		box-sizing: border-box;

		ul {
			list-style: none;
			padding: 0;

			li {
				display: flex;
				align-items: center;
				justify-content: center;
				padding: 10px;
				cursor: pointer;
				transition: background-color 0.2s;

				&:hover {
					// background-color: #e0e0e0;
				}

				.icon {
					font-size: 1.5em;
				}
			}
		}
	}

	.content {
		flex: 1;
		padding: 20px;
		overflow-y: auto;
		box-sizing: border-box;
	}

	@media screen and (max-width: 768px) {
		.window-outer {
			width: 100%;
			height: 100%;
			border: none;
			box-sizing: border-box;
			transition: background-color 0.1s;
			overflow: hidden;
		}

		.upbar {
			flex-direction: column;
			align-items: flex-start;
			padding: 10px;
		}

		.upbar-right {
			width: 100%;
			justify-content: space-between;
		}

		.settings-container {
			flex-direction: column-reverse;
		}

		.sidebar {
			width: 100%;
			display: flex;
			overflow-x: auto;
			border-right: none;
			// background-color: #f0f0f0;
			// border-bottom: 1px solid #dcdcdc;
			justify-content: center;
		}

		.sidebar ul {
			display: flex;
			flex-direction: row;
			width: 100%;
			justify-content: space-around;
		}

		.content {
			padding: 10px;
			flex: 1;
		}
	}

	.icon-container {
		display: flex;
		flex-direction: column;
		align-items: center;
	}

	.icon-title {
		// margin-top: 5px;
		font-size: 0.8em;
		color: #333;
	}
</style>
