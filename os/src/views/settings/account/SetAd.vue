<template>
	<div class="set-ad">
		<el-switch
			v-model="isAdEnabled"
			active-text="开启"
			inactive-text="关闭"
			@change="confirmAdChange"
		/>
	</div>
</template>

<script setup lang="ts">
	import { ref } from "vue";
	import { successMsg } from "@/utils/msg";

	const isAdEnabled = ref(false);

	function confirmAdChange() {
		successMsg("设置成功");
	}
</script>

<style scoped>
	.set-ad {
		display: flex;
		justify-content: center;
		align-items: center;
		width: 100%;
		height: 100%;
	}
</style>
