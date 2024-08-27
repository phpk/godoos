<template>
	<div class="ad-container" v-if="upgradeStore.hasAd" @click="onSponsorsClick">
		<el-carousel height="350px" indicator-position="none" :arrow="setCarouselShow" @change="onCarouselChange">
			<el-carousel-item v-for="(v, k) in upgradeStore.adList" :key="k">
				<img :src="v.img" v-if="v.img && v.img != ''" class="ad-img" />
				<div class="ad-text" v-html="v.desc"></div>
			</el-carousel-item>
		</el-carousel>
		<div class="ad-close">
            <el-icon :size="20"  @click.stop="onCloseSponsors">
            <CircleCloseFilled />
            </el-icon>
		</div>
	</div>
</template>

<script setup lang="ts" name="layoutSponsors">
import { reactive, computed } from 'vue';
import { useUpgradeStore } from '@/stores/upgrade';
const upgradeStore = useUpgradeStore();

// 定义变量内容
const state = reactive({
	ad: {
		isShow: false,
		index: 0,
	},
});

// 设置轮播图箭头显示
const setCarouselShow = computed(() => {
	return upgradeStore.adList.length <= 1 ? 'never' : 'hover';
});
// 关闭
const onCloseSponsors = () => {
	upgradeStore.hasAd = false;
};
// 轮播图改变时
const onCarouselChange = (e: number) => {
	state.ad.index = e;
};
// 当前项内容点击
const onSponsorsClick = () => {
    const link = upgradeStore.adList[state.ad.index]?.link
    if(link){
        window.open(link)
    }
};
</script>

<style scoped lang="scss">
.ad-container {
	position: fixed;
	top: calc(50vh - 175px);
    left: calc(50vw - 175px);
	z-index: 3;
	width: 350px;
    height: 350px;
	background-color: #ffffff;
	box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
	border-radius: 5px;
	overflow: hidden;
	cursor: pointer;
	.ad-img {
		width: 100%;
		height: 100px;
	}
	.ad-text {
		padding: 10px;
		color: #303133;
		font-size: 14px;
		text-align: left;
	}
	.ad-close {
		width: 60px;
		height: 60px;
		border-radius: 100%;
		position: absolute;
        background: rgba(0, 0, 0, 0.05);
		transition: all 0.3s ease;
		right: -30px;
		top: -30px;
		:deep(i) {
			position: absolute;
			left: 9px;
			bottom: 9px;
			color: #afafaf;
			transition: all 0.3s ease;
		}
		&:hover {
			transition: all 0.3s ease;
			:deep(i) {
				color: #409EFF;
				transition: all 0.3s ease;
			}
		}
	}
}
</style>
