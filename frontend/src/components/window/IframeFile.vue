<template>
	<iframe class="setiframe" allow="fullscreen" ref="storeRef" :src="src"></iframe>
</template>
<script lang="ts" setup name="IframeFile">
//@ts-ignore
import { BrowserWindow, Dialog, Notify, System } from "@/system";
import { getSplit, getSystemConfig, setSystemKey } from "@/system/config";
import { base64ToBuffer, isBase64 } from "@/util/file";
import { isShareFile } from "@/util/sharePath.ts";
import { inject, onMounted, onUnmounted, ref, toRaw } from "vue";
const SP = getSplit();

const sys: any = inject<System>("system");
const win: any = inject<BrowserWindow>("browserWindow");
const props = defineProps({
	src: {
		type: String,
		default: "",
	},
	eventType: {
		type: String,
		default: "",
	},
	ext: {
		type: String,
		default: "md",
	},
});
//console.log('iframe: ', props);

//console.log(props);
//let path = win?.config?.path;
// let currentPath = ref('')
const storeRef = ref<HTMLIFrameElement | null>(null);
let hasInit = false;
const eventHandler = async (e: MessageEvent) => {
	const eventData = e.data;
  
	if (eventData.type == props.eventType) {
		let data = JSON.parse(eventData.data);
		let title = data.title;
		let path;
		let ext: any = props.ext;
		if (ext instanceof Array) {
			ext = ext[0];
		}
		if (data.ext) {
			ext = data.ext;
		}
		// console.log(ext)
		// console.log(data)
		if (win.config && win.config.path) {
			path = win.config.path;
			//去除重复文件名后的（1）
			let fileTitleArr = path.split(SP).pop().split(".");
			let oldExt = fileTitleArr.pop();
			let fileTitle = fileTitleArr.join(".");
			if (fileTitle != title) {
				path = path.replace(fileTitle, title);
			}
			if (oldExt != ext) {
				path = path.replace("." + oldExt, "." + ext);
			}
		} else {
			path = `${SP}C${SP}Users${SP}Desktop${SP}${title}.${ext}`;
		}
		//判断是否共享文件，以及编辑权限
		const isShare = ref(false);
		const isWrite = ref(0);
		if (isShareFile(path)) {
			const file = await sys?.fs.getShareInfo(path);
			isShare.value = true;
			isWrite.value = file.fs.sender === getSystemConfig().userInfo.id ? 1 : file.fs.is_write;
			if (
				!isWrite.value &&
				file.fs.sender !== getSystemConfig().userInfo.id
			) {
				new Notify({
					title: "提示",
					content: "该文件没有编辑权限",
				});
				return;
			}
		} else if (await sys?.fs.exists(path)) {
			let res = await Dialog.showMessageBox({
				type: "info",
				title: "提示",
				message: "存在相同的文件名-" + title,
				buttons: ["覆盖文件?", "取消"],
			});
			//console.log(res)
			if (res.response > 0) {
				return;
			}
		}
		if (typeof data.content === "string") {
			if (data.content.indexOf(";base64,") > -1) {
				const parts = data.content.split(";base64,");
				data.content = parts[1];
			}
			if (isBase64(data.content)) {
				data.content = base64ToBuffer(data.content);
				//console.log(data.content)
			}
		}

		const res = isShare.value
			? await sys?.fs.writeShareFile(
				path,
				data.content,
				isWrite.value
			)
			: await sys?.fs.writeFile(path, data.content);
		// console.log("编写文件：", res, isShare);
		new Notify({
			title: "提示",
			content: res.message
			// content: res.code === 0 ? "文件已保存" : res.message,
		});
		sys.refershAppList();
	} else if (eventData.type == "initSuccess") {
		if (hasInit) {
			return;
		}
		hasInit = true;
		let content = win?.config?.content;
		let title = win.getTitle();
		// console.log("win.config;", win?.config);
		// console.log(title);
		title = title.split(SP).pop();

		if (!content && win?.config.path) {
			const file = getSystemConfig().file;
			const header = {
				salt: file.salt,
				pwd: file.pwd,
			};
			content = await sys?.fs.readFile(win?.config.path, header);
		}
		content = toRaw(content);
		if (content && content !== "") {
			storeRef.value?.contentWindow?.postMessage(
				{
					type: "init",
					data: { content, title },
				},
				"*"
			);
		} else {
			storeRef.value?.contentWindow?.postMessage(
				{
					type: "start",
					title,
				},
				"*"
			);
		}
	}
	else if (eventData.type == "close") {
		// console.log("关闭");
		win.close();
	}
	else if (eventData.type == "saveMind") {
		// console.log("保存");
		const data = eventData.data;
		const path = win?.config?.path;
		//console.log(path,data)
		const winMind = new BrowserWindow({
            title:data.title,
			url: "/mind/index.html",
			frame: true,
            config: {
                ext: 'mind',
                path: path,
				content: data.content
            },
            icon: "gallery",
            width: 700,
            height: 500,
            x: 100,
            y: 100,
            //center: true,
            minimizable: false,
            resizable: true,
        });
        winMind.show()
		
	}
  else if (eventData.type == 'aiCreater') {
    // 模拟AI返回数据
    storeRef.value?.contentWindow?.postMessage(
      {
        type: 'aiReciver',
        data: '-------------经过AI处理后的数据-----------',
      },
      "*"
    );
  }
  else if (eventData.type == 'aiReciver') {
    storeRef.value?.contentWindow?.postMessage(
      {
        type: eventData.type,
        data: '----经过AI处理后的数据-----',
      },
      "*"
    );
  }
};
//删除本地暂存的文件密码
const delFileInputPwd = async () => {
	let fileInputPwd = getSystemConfig().fileInputPwd;
	const currentPath = win.config.path;
	const temp = fileInputPwd.filter(
		(item: any) => item.path !== currentPath
	);
	setSystemKey("fileInputPwd", temp);
};
onMounted(() => {
	window.addEventListener("message", eventHandler);
});

onUnmounted(async () => {
	await delFileInputPwd();
	window.removeEventListener("message", eventHandler);
});
</script>
<style scoped>
.setiframe {
	width: 100%;
	height: 100%;
	border: none;
}
</style>
