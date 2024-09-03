<template>
	<div class="notice-container" v-if="upgradeStore.hasNotice" @click="onnoticeClick">
		<el-carousel height="240px" indicator-position="none" :arrow="setCarouselShow" @change="onCarouselChange">
			<el-carousel-item v-for="(v, k) in upgradeStore.noticeList" :key="k">
				<img :src="v.img" v-if="v.img && v.img != ''" class="notice-img" />
				<div class="notice-text" v-html="v.desc"></div>
			</el-carousel-item>
		</el-carousel>
		<div class="notice-close">
            <el-icon :size="20"  @click.stop="onClosenotice">
            <CircleCloseFilled />
            </el-icon>
		</div>
	</div>
	<ShowNews />
</template>

<script setup lang="ts" name="layoutnotice">
import { reactive, computed } from 'vue';
import { useUpgradeStore } from '@/stores/upgrade';
import { useNotifyStore } from '@/stores/notify';
const upgradeStore = useUpgradeStore();
const notifyStore = useNotifyStore();
// 定义变量内容
const state = reactive({
	notice: {
		index: 0,
	},
});

// 设置轮播图箭头显示
const setCarouselShow = computed(() => {
	return upgradeStore.noticeList.length <= 1 ? 'never' : 'hover';
});
// 关闭
const onClosenotice = () => {
	upgradeStore.hasNotice = false;
};
// 轮播图改变时
const onCarouselChange = (e: number) => {
	state.notice.index = e;
};
// 当前项内容点击
const onnoticeClick = () => {
	const data = upgradeStore.noticeList[state.notice.index]
	const link = data?.link
    if(link){
        window.open(link)
    }
	if(data?.content){
		notifyStore.viewContent(data)
	}

};
</script>

<style scoped lang="scss">
.notice-container {
	position: fixed;
	right: 5px;
	bottom: 50px;
	z-index: 3;
	width: 200px;
	background-color: #ffffff;
	box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
	border-radius: 5px;
	overflow: hidden;
	cursor: pointer;
	.notice-img {
		width: 100%;
		height: 80px;
	}
	.notice-text {
		padding: 10px;
		color: #303133;
		font-size: 12px;
		text-align: left;
	}
	.notice-close {
		width: 60px;
		height: 60px;
		border-radius: 100%;
		position: absolute;
        background: rgba(0, 0, 0, 0.05);
		transition: all 0.3s ease;
		right: -30px;
		bottom: -30px;
		:deep(i) {
			position: absolute;
			left: 9px;
			top: 9px;
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
