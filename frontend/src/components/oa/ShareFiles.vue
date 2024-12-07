<template>
	<el-form :model="form" label-width="auto" style="max-width: 560px; margin-top: 20px; padding: 20px">
		<el-form-item label="分享给">
			<el-select v-model="form.receiverId" remote :remote-method="handleSearch" filterable multiple clearable collapse-tags placeholder="选择人员"
				popper-class="custom-header" :max-collapse-tags="1" value-key="id" style="width: 240px"
				@change="checkUsers">
				<template #header>
					<el-checkbox v-model="checkAll" @change="handleCheckAll">
						全选
					</el-checkbox>
				</template>
				<el-option v-for="item in userList" :key="item.id" :label="item.nickname" :value="item.id" />
				<template #footer v-if="total > 10">
					<el-pagination small background layout="prev, pager, next" :total="total"
						:page-size="pageSize" :current-page="currentPage" @current-change="handlePageChange"
						style="margin-top: 10px" />
				</template>
			</el-select>
		</el-form-item>
		<el-form-item label="编辑权限">
			<el-switch v-model="form.iswrite" active-value="1" inactive-value="0" />
		</el-form-item>
		<div class="btn-group">
			<el-button type="primary" @click="onSubmit">发布分享</el-button>
		</div>
	</el-form>
</template>

<script lang="ts" setup>
import { BrowserWindow } from "@/system";
import { fetchPost, getSystemConfig, getUrl } from "@/system/config";
import { notifyError, notifySuccess } from "@/util/msg";
import { inject, onMounted, ref } from "vue";
const window: BrowserWindow | undefined = inject("browserWindow");
//const sys = useSystem();
//const userInfo: any = sys.getConfig("userInfo");
const userList: any = ref([]);
const checkAll = ref(false);
const form: any = ref({
	receiverId: [],
	path: "",
	iswrite: "0",
});
const currentPage = ref(1);
const pageSize = 10;
const total = ref(0);
const config = ref(getSystemConfig());
const pageUrl = ref("")
const getList = async() => {
	const apiUrl = getUrl(pageUrl.value, true)
	//console.log(apiUrl)
	const res: any = await fetch(apiUrl);
	//console.log(res)
	if (res.ok) {
		const rt = await res.json();
		//console.log(rt)
		userList.value = rt.data.list;
		total.value = rt.data.total;
		userList.value = userList.value.filter((item: any) => {
			return item?.id !== config.value.userInfo.id;
		});
		if(userList.value.length != rt.data.list.length){
			total.value = total.value - 1;
		}
	}
}
onMounted(async () => {
	pageUrl.value = '/user/sharelist?page='+currentPage.value
	await getList()
});
const handleCheckAll = (val: any) => {
	if (val) {
		form.value.receiverId = userList.value.map((d: any) => d.value);
	} else {
		form.value.receiverId = [];
	}
};
const checkUsers = (val: any) => {
	const res: any = [];
	val.forEach((item: any) => {
		if (item) {
			res.push(item);
		}
	});
	form.value.receiverId = res;
};
const onSubmit = async () => {
	const apiUrl = config.value.userInfo.url + "/files/share";
	const path = window?.config.path || "";
	const receiverIds = form.value.receiverId.map((item: any) => item*1);
	if(receiverIds.length < 1){
		notifyError("请选择分享对象")
		return;
	}
	const postData = {
		path: path,
		receiverId: receiverIds,
		isWrite : form.value.iswrite * 1
	}
	//console.log(postData)
	const res = await fetchPost(apiUrl, JSON.stringify(postData));
	const result = await res.json();
	if (res.ok && result.success) {
		notifySuccess(result.message || "分享文件成功");
		window?.close()
	} else {
		notifyError(result.message || "分享文件失败");
	}
};
const handleSearch = async (val:any) => {
	if(val !== ''){
		pageUrl.value = '/user/sharelist?page=1&nickname='+val
	}else{
		pageUrl.value = '/user/sharelist?page=1'
	}
	await getList()
};
const handlePageChange = async (page: number) => {
	currentPage.value = page;
	await getList()
};
// onMounted(() => {
// 	userList.value = userList.value.filter((item: any) => {
// 		return item?.id !== config.value.userInfo.id;
// 	});
// });
</script>
<style scoped>
.btn-group {
	display: flex;
	justify-content: center;
}
</style>
