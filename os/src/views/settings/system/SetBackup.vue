<template>
	<div class="container">
		<el-form
			label-position="top"
			class="form-container"
		>
			<el-form-item>
				<el-button type="primary">导出</el-button>
			</el-form-item>
			<el-form-item>
				<el-button type="primary">导入</el-button>
				<input
					type="file"
					accept=".zip"
					style="display: none"
					ref="zipFileInput"
				/>
			</el-form-item>
		</el-form>
	</div>
</template>

<script lang="ts" setup></script>

<style scoped>
	.container {
		display: flex;
		justify-content: center;
		align-items: center;
		width: 100%;
		height: 100%;
	}

	.form-container {
		display: flex;
		justify-content: center;
		align-items: center;
		width: 22rem;
		gap: 1rem;
	}

	.setting-item {
		display: flex;
		align-items: center;
	}

	.setting-area {
		width: 100%;
		height: 90%;
	}
</style>
