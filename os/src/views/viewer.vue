<template>
	<div
		ref="container"
		class="viewer-container"
	>
		<div v-if="fileType === 'pdf'">
			<iframe
				:src="prop.content"
				type="application/pdf"
				class="doc-container"
				frameborder="0"
			></iframe>
		</div>
		<div v-else-if="fileType === 'pic'">
			<el-image
				:src="'data:image/png;base64,'+prop.content"
				alt=""
				style="margin: auto"
			/>
		</div>
		<div v-else-if="fileType === 'excel'">
			<ViewExcel
				:src="prop.content"
				:ext="prop.ext"
				style="height: 100vh"
			/>
		</div>
		<div
			v-else-if="fileType === 'ppt'"
			v-loading="fileLoading"
		>
			<div
				id="pptx-wrapper"
				ref="pptxWrapper"
			>
				<div
					v-for="(slide, index) in prop.content.slides"
					:key="index"
					class="slide"
				>
					<div
						v-for="(element, elementIndex) in slide.elements"
						:key="elementIndex"
						class="element"
						:style="{
							left: element.left + 'px',
							top: element.top + 'px',
							width: element.width + 'px',
							height: element.height + 'px',
						}"
					>
						<div
							v-if="element.type === 'text'"
							v-html="element.content"
						></div>
						<img
							v-else-if="element.type === 'image'"
							:src="element.src"
							alt="image"
							:style="{
								width: element.width + 'px',
								height: element.height + 'px',
							}"
						/>
						<!-- 添加其他类型的元素渲染 -->
					</div>
				</div>
			</div>
		</div>
		<div
			v-else-if="fileType === 'word'"
			style="height: 100%"
		>
			<iframe
				ref="fileIframe"
				:style="{
					width: '100%',
					height: isMobileDevice() ? '100%' : '88vh',
				}"
				src="/os/word/index.html"
				frameborder="0"
				@load="handleIframeLoad"
				v-loading="fileLoading"
			></iframe>
		</div>
		<div v-else-if="fileType === 'md'">
			<div
				v-html="prop.content"
				style="padding: 20px"
			></div>
		</div>
		<div v-else>
			<p>Unsupported file type</p>
		</div>
	</div>
</template>

<script setup lang="ts">
	//import { readFile } from "@/api/files";
	import { useFileSystemStore } from "@/stores/filesystem";
	import {
		base64ToBlobPdfUrl,
		base64ToBuffer,
		decodeBase64,
		determineFileType,
	} from "@/utils/file";
	import { renderMarkdown } from "@/utils/markdown";
	import { addWatermark } from "@/utils/watermark";
	import { onMounted, ref } from "vue";
	// import { init } from 'pptx-preview';
	import { isMobileDevice } from "@/utils/device";
	//import { parse } from "pptxtojson";
	const fileSystemStore = useFileSystemStore();
	const props = defineProps({
		win: {
			type: Object,
			required: true,
		},
	});
	const prop = ref(props.win.props);

	const fileType = ref("");
	const fileIframe = ref<HTMLIFrameElement | null>(null);
	const pptxWrapper = ref<HTMLElement | null>(null);
	const container: any = ref<HTMLElement | null>(null);

	const fileLoading = ref(true);
	onMounted(async () => {
		//if (prop.value.content === '') {
		//prop.value.content = await readFile(prop.value.path);
		//}
		prop.value.content = await fileSystemStore.handleReadFile(
			prop.value.path
		);
		console.log(prop.value.content)
		await loadFileContent();
		addWatermark("GodoOS", container);
	});

	// function transformData(data:any){
	//     console.log('transformData', data);
	//     return data;
	// }

	async function renderPPTX() {
		const buffer: any = base64ToBuffer(prop.value.content);
		//const arrayBuffer = buffer.buffer.slice(buffer.byteOffset, buffer.byteOffset + buffer.byteLength);
		// const json = await parse(buffer);
		// prop.value.content = json; // 将解析后的 JSON 数据赋值给 prop.value.content
		// fileLoading.value = false;
		// if (pptxWrapper.value) {
		//   const pptxPrviewer = init(pptxWrapper.value, {
		//     width: window.innerWidth * 0.9,
		//     height: window.innerHeight - 30
		//   });
		//   //console.log(prop.value.content)
		//   const buffer = base64ToBuffer(prop.value.content);
		//   pptxPrviewer.preview(buffer);
		//   fileLoading.value = false;
		// }
	}

	const loadFileContent = async () => {
		try {
			fileType.value = determineFileType(prop.value.ext);
			//console.log(fileType.value);
			switch (fileType.value) {
				case "excel":
					//prop.value.content = base64ToBuffer(prop.value.content);
					break;
				case "pdf":
					prop.value.content = base64ToBlobPdfUrl(prop.value.content);
					break;
				case "pic":
					prop.value.content =
						`data:image/${prop.value.ext};base64,` +
						prop.value.content;
					break;
				case "word":
					//prop.value.content = decodeBase64(prop.value.content)
					break;
				case "ppt":
					setTimeout(() => {
						renderPPTX();
					}, 1000);

					break;
				case "md":
					//const content = decodeBase64(prop.value.content);
					prop.value.content = renderMarkdown(prop.value.content);
					break;
				default:
					console.warn("Unsupported file type");
			}
		} catch (error) {
			console.error(`Error loading ${fileType.value} file:`, error);
		}
	};

	const handleIframeLoad = () => {
		if (fileIframe.value && fileIframe.value.contentWindow) {
			const postData = {
				content: prop.value.content,
			};
			fileIframe.value.contentWindow.postMessage(postData, "*");
		}
		fileLoading.value = false;
	};
</script>

<style scoped>
	.viewer-container {
		height: 100%;
		padding: 0;
	}

	.doc-container {
		width: 100%;
		height: calc(100vh - 50px);
	}

	.slide {
		position: relative;
		width: 960pt;
		/* 根据实际幻灯片宽度调整 */
		height: 540pt;
		/* 根据实际幻灯片高度调整 */
		margin-bottom: 20px;
	}

	.element {
		position: absolute;
		box-sizing: border-box;
	}
</style>
