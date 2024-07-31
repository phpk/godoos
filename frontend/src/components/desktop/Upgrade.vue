<template>
	<div class="upgrade-dialog">
		<el-dialog
			v-model="state.isUpgrade"
			width="300px"
			destroy-on-close
			:show-close="false"
			:close-on-click-modal="false"
			:close-on-press-escape="false"
		>
			<div class="upgrade-title">
				<div class="upgrade-title-warp">
					<span class="upgrade-title-warp-txt">{{ $t('upgrade.title') }}</span>
					<span class="upgrade-title-warp-version">v{{ upgradeStore.versionTag }}</span>
				</div>
			</div>
			<div class="upgrade-content">
				GodoOS {{ $t('upgrade.msg') }}
				<div class="mt5">
					<el-link type="primary" class="font12" href="https://gitee.com/lyt-top/vue-next-admin/blob/master/CHANGELOG.md" target="_black">
						CHANGELOG.md
					</el-link>
				</div>
				<div class="upgrade-content-desc mt5">{{ $t('upgrade.desc') }}</div>
				<div class="upgrade-content-desc" v-if="upgradeStore.progress > 0">
					<el-progress :text-inside="true" :stroke-width="20" :percentage="upgradeStore.progress" />
				</div>
			</div>
			<div class="upgrade-btn">
				<el-button round size="default" type="info" text @click="onCancel">{{ $t('upgrade.btnOne') }}</el-button>
				<el-button type="primary" round size="default" @click="onUpgrade" :loading="state.isLoading">{{ state.btnTxt }}</el-button>
			</div>
		</el-dialog>
	</div>
</template>

<script setup lang="ts" name="layoutUpgrade">
import { reactive, onMounted } from 'vue';
import { t } from '@/i18n';
import { useUpgradeStore } from '@/stores/upgrade';
const upgradeStore = useUpgradeStore();
const state = reactive({
	isUpgrade: false,
	isLoading: false,
	btnTxt: '',
});


// 残忍拒绝
const onCancel = () => {
	state.isUpgrade = false;
};
// 马上更新
const onUpgrade = () => {
	state.isLoading = true;
	state.btnTxt = t('upgrade.btnTwoLoading');
	setTimeout(async () => {
		await upgradeStore.update();
		state.isLoading = false;
	}, 2000);
};
// 延迟显示，防止刷新时界面显示太快
const delayShow = () => {
	setTimeout(() => {
		state.isUpgrade = true;
	}, 2000);
};
// 页面加载时
onMounted(() => {
	delayShow();
	setTimeout(() => {
		state.btnTxt = t('upgrade.btnTwo');
	}, 200);
});
</script>

<style scoped lang="scss">
.upgrade-dialog {
	:deep(.el-dialog) {
		padding: 0 !important;
		.el-dialog__body {
			padding: 0 !important;
		}
		.el-dialog__header {
			display: none !important;
		}
		.upgrade-title {
			text-align: center;
			height: 100px;
			display: flex;
			align-items: center;
			justify-content: center;
			position: relative;
			&::after {
				content: '';
				position: absolute;
				background-color: #2882dc;
				width: 100%;
				height: 100px;
				border-bottom-left-radius: 100%;
				border-bottom-right-radius: 100%;
			}
			.upgrade-title-warp {
				z-index: 1;
				position: relative;
				.upgrade-title-warp-txt {
					color: #ffffff;
					font-size: 22px;
					letter-spacing: 3px;
				}
				.upgrade-title-warp-version {
					color: #ffffff;
					background-color: #2882dc;
					font-size: 12px;
					position: absolute;
					display: flex;
					top: -2px;
					right: -50px;
					padding: 2px 4px;
					border-radius: 2px;
				}
			}
		}
		.upgrade-content {
			padding: 20px;
			line-height: 22px;
			.upgrade-content-desc {
				color: rgba(255, 255, 255, 0.7);
				font-size: 12px;
			}
		}
		.upgrade-btn {
			border-top: 1px solid #ebeef5;
			display: flex;
			justify-content: space-around;
			padding: 15px 20px;
			.el-button {
				width: 100%;
			}
		}
	}
}
</style>
