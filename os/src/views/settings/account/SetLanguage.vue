<template>
	<div class="language-select-container">
		<el-select class="language-select" v-model="selectedLanguage" placeholder="选择语言">
			<el-option v-for="lang in languages" :key="lang.value" :label="lang.label" :value="lang.value">
			</el-option>
		</el-select>
		<!-- 添加确认按钮 -->
		<el-button class="confirm-button" type="primary" @click="confirmLanguageChange">确认</el-button>
	</div>
</template>

<script setup>
import { ref } from "vue";
import { successMsg } from "@/utils/msg";
import { languages, getLang } from "@/i18n";
import { eventBus } from '@/interfaces/event';
const selectedLanguage = ref(getLang()); // 当前选择的语言


function confirmLanguageChange() {
	//setLang(selectedLanguage.value); // 更新语言设置
	eventBus.emit('setLanguages',selectedLanguage.value);
	successMsg("设置成功");
}
</script>

<style scoped>
.language-select {
	width: 250px;
	/* 设置选择菜单的宽度 */
	margin-bottom: 20px;
	/* 添加下边距 */
}

.confirm-button {
	width: 250px;
	/* 设置按钮的宽度与选择菜单一致 */
	height: 40px;
	/* 设置按钮的高度 */
}

.language-select-container {
	display: flex;
	flex-direction: column;
	align-items: center;
	width: 100%;
	/* 调整容器宽度 */
	padding: 20px;
	/* 添加内边距 */
}
</style>
