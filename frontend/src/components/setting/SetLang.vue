<template>
	<div class="container">
		<div class="nav">
			<ul>
				<li
					v-for="(item, index) in items"
					:key="index"
					@click="selectItem(index)"
					:class="{ active: index === activeIndex }"
				>
					{{ item }}
				</li>
			</ul>
		</div>
		<div class="setting">
			<div v-if="0 === activeIndex">
				<div class="setting-item">
					<h1 class="setting-title">{{ t("language") }}</h1>
				</div>
				<div class="setting-item">
					<label></label>
					<el-select v-model="modelvalue">
						<el-option
							v-for="(item, key) in langList"
							:key="key"
							:label="item.label"
							:value="item.value"
						/>
					</el-select>
				</div>

				<div class="setting-item">
					<label></label>
					<el-button
						@click="submit"
						type="primary"
					>
						{{ t("confirm") }}
					</el-button>
				</div>
			</div>
		</div>
	</div>
</template>

<script lang="ts" setup>
	import { getLang, setLang, t } from "@/i18n";
	import { ElMessageBox } from "element-plus";
	import { ref } from "vue";
	import { useI18n } from "vue-i18n";
	const { locale } = useI18n();
	// const showDialog = ref(false);
	const langList = [
		{
			label: "中文",
			value: "zh-cn",
		},
		{
			label: "English",
			value: "en",
		},
	];
	const items = [t("language")];

	const activeIndex = ref(0);
	const currentLang = getLang();
	const modelvalue = ref(currentLang);

	const selectItem = (index: number) => {
		activeIndex.value = index;
	};

	async function submit() {
		setLang(modelvalue.value);
		locale.value = modelvalue.value;
		ElMessageBox.alert(t("save.success"), t("language"), {
			confirmButtonText: "OK",
		});
		//showDialog.value = false

		// Dialog.showMessageBox({
		//   message: t("save.success"),
		//   title: t("language"),
		//   type: "info",
		// });
	}
</script>
<style scoped>
	@import "./setStyle.css";
</style>
