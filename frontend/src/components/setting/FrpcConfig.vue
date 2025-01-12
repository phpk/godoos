<script lang="ts" setup>
import { onMounted, reactive, ref } from "vue";
import { useProxyStore } from "@/stores/proxy";
import { notifyError, notifySuccess } from "@/util/msg";
import { FormInstance, FormRules } from "element-plus";
const store = useProxyStore();
type FrpConfig = {
    isOn: boolean;
    serverAddr: string;
    serverPort: number;
    authMethod: string;
    authToken: string;
    user: string;
};

const defaultFormData = ref<FrpConfig>({
    isOn: true,
    serverAddr: "",
    serverPort: 7000,
    authMethod: "null",
    authToken: "",
    user: ""
});

const formData = ref<FrpConfig>(defaultFormData.value);



const rules = reactive<FormRules>({
    serverAddr: [
        { required: true, message: "请输入服务端地址", trigger: "blur" },
        {
            pattern: /^[\w-]+(\.[\w-]+)+$/,
            message: "请输入正确的服务端地址",
            trigger: "blur"
        }
    ],
    serverPort: [
        { required: true, message: "请输入服务器端口", trigger: "blur" }
    ]
});


const formRef = ref<FormInstance>();
const handleSubmit = () => {
    if (!formRef.value) return;
    formRef.value.validate(valid => {
        if (valid) {
            const data = toRaw(formData.value);
            store.setConfig(data).then(res => {
                if (res.code == 0) {
                    notifySuccess("修改成功");
                } else {
                    notifyError("修改失败");
                }
            });
        } else {
            notifyError("请填写完整");
        }
    });

}

onMounted(async () => {
    const res = await store.getConfig();
    if (res.code == 0) {
        formData.value = res.data;
    }
});
</script>
<template>
    <el-form :model="formData" :rules="rules" label-position="right" ref="formRef" label-width="150">
        <el-form-item label="服务器地址：" prop="serverAddr">
            <el-input v-model="formData.serverAddr" placeholder="127.0.0.1"></el-input>
        </el-form-item>
        <el-form-item label="服务器端口：" prop="serverPort">
            <el-input-number placeholder="7000" v-model="formData.serverPort" :min="0" :max="65535"
                controls-position="right" class="!w-full"></el-input-number>
        </el-form-item>
        <el-form-item label="跟随启动：" prop="isOn">
            <el-switch active-text="开" inline-prompt inactive-text="关" v-model="formData.isOn" />
        </el-form-item>
        <el-form-item label="验证方式：" prop="authMethod">
            <el-select v-model="formData.authMethod" placeholder="请选择验证方式" clearable>
                <el-option label="无" value="null"></el-option>
                <el-option label="令牌（token）" value="token"></el-option>
            </el-select>
        </el-form-item>
        <el-form-item label="令牌：" prop="authToken" v-if="formData.authMethod === 'token'">
            <el-input placeholder="token" type="password" v-model="formData.authToken" :show-password="true" />
        </el-form-item>
        <el-row justify="center">
            <el-button type="primary" @click="handleSubmit" style="width: 100px;margin-top: 15px;">
                保存
            </el-button>
        </el-row>
    </el-form>

</template>

<style lang="scss" scoped>
.button-input {
    width: calc(100% - 68px);
}
</style>