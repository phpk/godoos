<template>
	<div class="s-option">
		<div class="s-option-inner">
			<div
				v-if="sys.getConfig('userType') == 'member'"
				class="s-option-button"
				@click.stop="handleLogout"
				v-glowing
			>
				<div class="s-option-button_img">
					<svg
						class="icon"
						aria-hidden="true"
					>
						<use xlink:href="#icon-logout"></use>
					</svg>
				</div>
				<div class="s-option-button_title">
					{{ t("startMenu.logout") }}
				</div>
			</div>
			<div
				class="s-option-button"
				@click.stop="handleRecover"
				v-glowing
			>
				<div class="s-option-button_img">
					<svg
						class="icon"
						aria-hidden="true"
					>
						<use xlink:href="#icon-huifu"></use>
					</svg>
				</div>
				<div class="s-option-button_title">
					{{ t("startMenu.recover") }}
				</div>
			</div>
			<div
				class="s-option-button"
				@click.stop="($ev) => handleClick(1, $ev)"
				v-glowing
			>
				<div class="s-option-button_img">
					<svg
						class="icon"
						aria-hidden="true"
					>
						<use xlink:href="#icon-setting"></use>
					</svg>
				</div>
				<div class="s-option-button_title">
					{{ t("setting") }}
				</div>
			</div>
		</div>
	</div>
</template>

<script lang="ts" setup>
	import { BrowserWindow, Dialog, t, useSystem } from "@/system";
	import { fetchGet } from "@/system/config";
	import { emitEvent } from "@/system/event";
	import { vGlowing } from "@/util/glowingBorder";
	import { RestartApp } from "@/util/goutil";

	const sys = useSystem();

	function handleLogout() {
		const userInfo: any = sys.getConfig("userInfo");
		fetchGet(userInfo.url + "/member/loginout").then(() => {
			RestartApp();
		});
	}

	function handleRecover() {
		Dialog.showMessageBox({
			title: t("startMenu.recover"),
			message: t("is.recover"),
			buttons: [t("startMenu.recover"), t("cancel")],
		}).then((res) => {
			if (res.response === -1) {
				emitEvent("system.recover");
			}
		});
	}

	function handleClick(key: number, ev: MouseEvent) {
		if (key === 1) {
			emitEvent("startMenu.set.click", {
				mouse: ev,
			});
			const winopt = sys._rootState.windowMap["Menulist"].get("setting");

			if (winopt) {
				if (winopt._hasShow) {
					return;
				} else {
					winopt._hasShow = true;
					const win = new BrowserWindow(winopt.window);
					win.show();
					win.on("close", () => {
						winopt._hasShow = false;
					});
				}
			}
		}
	}
</script>
<style lang="scss" scoped>
	.s-option {
		position: relative;
		width: var(--startmenu-icon-size);
		height: 100%;
		z-index: 40;
		user-select: none;

		.s-option-inner {
			position: absolute;
			height: 100%;
			width: var(--startmenu-icon-size);
			background-color: #f5f5f5;
			// background-color: var(--theme-main-color-opacity);

			transition: width 0.1s ease-in-out, box-shadow 0.1s ease-in-out;
			transition-delay: 0s;
			display: flex;
			flex-direction: column-reverse;
			overflow: hidden;

			.s-option-button {
				height: var(--start-option-size);
				width: var(--startmenu-icon-size-hover);
				display: flex;
				justify-content: center;
				align-items: center;
				position: relative;
				z-index: 1;

				.s-option-button_img {
					height: var(--start-option-size);
					width: var(--start-option-size);
					display: flex;
					justify-content: center;
					align-items: center;

					svg {
						width: 40%;
					}
				}

				.s-option-button_title {
					width: 100%;
					height: var(--start-option-size);
					display: flex;
					justify-content: flex-start;
					align-items: center;
					opacity: 0.8;
				}
			}

			.s-option-button:hover {
				// background-color: var(--color-gray-hover);
			}
		}
	}

	.s-option:hover {
		.s-option-inner {
			transition-delay: 0.5s;
			width: var(--startmenu-icon-size-hover);
			box-shadow: 10px 0px 20px 0px rgba(0, 0, 0, 0.216);
		}
	}
</style>
