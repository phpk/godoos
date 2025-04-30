<template>
	<div class="container">
		<div class="setting">
			<div>
				<div class="setting-item">
					<el-select v-model="bgConfig.type">
						<el-option v-for="(item, key) in desktopOptions" :key="key" :label="item.label"
							:value="item.value" />
					</el-select>
				</div>
				<template v-if="bgConfig.type === 'color'">
					<div class="setting-item" style="margin-top: 10px;display: flex;justify-content: center;">
						<ColorPicker v-model:modelValue="bgConfig.color" @update:modelValue="onColorChange">
						</ColorPicker>
					</div>
				</template>
				<template v-if="bgConfig.type === 'image'">
					<div class="setting-item">
						<ul class="image-gallery">
							<li
								v-for="(item, index) in bgConfig
									.imageList"
								:key="index"
								:class="
									bgConfig.url === item
										? 'selected'
										: ''
								"
								@click="setBg(item)"
							>
								<img :src="item" />
							</li>
						</ul>
					</div>
					<div class="setting-item">
						<label> </label>
					</div>
				</template>
			</div>
		</div>
	</div>
</template>

<script lang="ts" setup>
import { useSettingsStore } from "@/stores/settings";
import ColorPicker from "./ColorPicker.vue";
import { ref } from "vue";
const store = useSettingsStore();

const bgConfig: any = ref(store.config.background);
const desktopOptions = [
	{
		label: "图片",
		value: "image",
	},
	{
		label: "颜色",
		value: "color",
	},
];

function setBg(item: any) {
	bgConfig.value.url = item;
	bgConfig.value.type = "image";
	store.setConfig('background',bgConfig.value);
}
function onColorChange(color: string) {
	bgConfig.value.color = color;
	bgConfig.value.type = "color";
	store.setConfig('background',bgConfig.value);
}
</script>
<style lang="scss" scoped>
@use "@/styles/imglist.scss";
.image-gallery{
	margin:1vw auto;
}
</style>
