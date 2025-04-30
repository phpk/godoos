<template>
	<div class="buttons-controls">
		<el-button class="el-icon" :class="{ active: prop.action === 'file' }" :icon="Document"
			@click="handleFile" v-if="!isMobileDevice()">
			文件
		</el-button>
		<el-button class="el-icon" :class="{ active: prop.action === 'edit' }" :icon="Edit" @click="handleEdit">
			编辑
		</el-button>
		<el-button class="el-icon" :class="{ active: prop.action === 'preview' }" :icon="View"
			@click="handlePreview" v-if="win.props?.hasPrview == 'true'">
			预览
		</el-button>

		<el-button v-if="prop.isShare === 'true' && prop.action === 'edit'" class="el-icon" :icon="Share" @click="isShow = true">
			协同
		</el-button>

		<!-- <el-button class="el-icon" :icon="Share" @click="handleShare">
			分享
		</el-button>
		<el-button class="el-icon" :icon="Star" @click="handleFavorite">
			收藏
		</el-button> -->
		<!-- <el-button class="el-icon" :class="{ 'active': win.props?.action === 'readAloud' }" :icon="Microphone"
			@click="handleReadAloud">
			朗读
		</el-button> -->

		<!-- 协同 -->
		<el-dialog :modal="false" v-if="prop.isShare == 'true'" :close-on-click-modal="false" v-model="isShow" title="协同" width="500px">
			<Collaboration :path="prop.path" :truePath="prop.truePath" />
		</el-dialog>
	</div>
</template>

<script setup lang="ts">
import { Document, Edit, Share, View } from "@element-plus/icons-vue";
import { defineEmits, defineProps, onMounted, ref } from "vue";
//import { determineFileType } from "@/utils/file";
import { isMobileDevice } from "@/utils/device";
import Viewer from "@/views/viewer.vue";

// import { getFileType } from '@/router/filemaplist'
const props = defineProps({
	win: {
		type: Object,
		required: true,
	},
});
let win: any = props.win;
let prop: any = win.props;

const emit = defineEmits(["updateWin"]);
onMounted(() => {
	console.log(prop)
	if (prop.action == "edit") {
		//const fileType: any = getFileType(prop.ext);
		win.component = win.props.editor;
		emit("updateWin", win);
	}
});

const isShow = ref(false);
const handleFile = () => {
	// 添加文件按钮的逻辑
	win.props.showFiles = !win.props.showFiles;
	emit("updateWin", win);
};

const handleEdit = () => {
	if (prop.action === "edit") return;
	win.props.action = "edit";
	win.component = win.props.editor;
	emit("updateWin", win);
};

const handlePreview = () => {
	if (prop.action === "preview") return;
	win.props.action = "preview";
	win.component = Viewer;
	emit("updateWin", win);
};


// const handleFavorite = () => {
// 	// 添加收藏按钮的逻辑
// };

// const handleReadAloud = () => {
// 	console.log('朗读按钮点击');
// 	// 添加朗读按钮的逻辑
// };
</script>

<style scoped>
.buttons-controls {
	display: flex;
	gap: 5px;
	align-items: center;
	margin-right: auto;
}

.buttons-controls .el-icon {
	width: 65px;
	height: 100%;
	cursor: pointer;
	display: flex;
	align-items: center;
	justify-content: center;
	background-color: transparent;
	/* 浅色背景 */
	border: 1px solid rgba(0, 0, 0, 0.1);
	/* 细边框 */
	border-radius: 4px;
	/* 圆角 */
	transition: background-color 0.2s ease;
	/* 悬停效果过渡 */
}

.buttons-controls .el-icon:hover {
	background-color: rgba(0, 0, 0, 0.05);
	border-color: rgba(0, 0, 0, 0.2);
	color: #333;
	/* 悬停时文字颜色 */
}

.buttons-controls .el-icon.active {
	background-color: rgba(0, 0, 0, 0.1);
	/* 激活时的背景颜色 */
	border-color: rgba(0, 0, 0, 0.3);
	/* 激活时的边框颜色 */
	box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
	/* 激活时的阴影效果 */
}

.buttons-controls .el-icon i {
	color: #333;
	/* 图标颜色 */
}
</style>
