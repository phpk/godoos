<script lang="ts" setup>
import { useProxyStore } from "@/stores/proxy";
import { notifyError, notifySuccess } from "@/util/msg";
import { ref } from "vue";
import { OpenDirDialog } from "@/util/goutil";
const proxyStore = useProxyStore();
const {
    updateProxy,
} = proxyStore;

const proxyDialogShow = ref(false);
const pwdRef = ref<any>(null);

// 定义表单验证规则
const proxyRules = {
    name: [{ required: true, message: "请输入代理名称", trigger: "blur" }],
    // port: [
    //     { required: true, message: "请输入端口", trigger: "blur" },
    //     { type: "number", message: "端口必须是数字", trigger: "blur" },
    // ],
    domain: [{ required: false, message: "请输入域名", trigger: "blur" }],
    type: [{ required: true, message: "请选择类型", trigger: "blur" }],
    localIp: [
        { required: false, message: "请输入内网地址", trigger: "blur" },
        {
            pattern: /^[\w-]+(\.[\w-]+)+$/,
            message: "请输入正确的内网地址",
            trigger: "blur"
        }
    ],
    localPort: [
        { required: true, message: "请输入内网端口", trigger: "blur" },
        {
            pattern: /^(?:\d{1,5}|\d{1,5}-\d{1,5})(?:,(?:\d{1,5}|\d{1,5}-\d{1,5}))*$/,
            message: "请输入正确的端口",
            trigger: "blur"
        },
    ]
};

const addProxy = async () => {
    const res = await proxyStore.createFrpc()
    if (res.code !== 0) {
        notifyError(res.message);
        return;
    }
    proxyStore.addShow = false;
    proxyStore.resetProxyData();
    notifySuccess("代理配置已成功创建");
    await proxyStore.fetchProxies();
};

const update = async () => {
    await updateProxy(proxyStore.proxyData);
    notifySuccess("编辑成功");
    proxyStore.addShow = false;
    proxyStore.resetProxyData();
    proxyStore.isEditor = false;
    await proxyStore.fetchProxies();
};

const saveProxy = () => {
    pwdRef.value.validate((valid: boolean) => {
        if (valid) {
            if (proxyStore.isEditor) {
                update();
            } else {
                addProxy();
            }
        } else {
            console.log("表单验证失败");
        }
    });
};



const proxyTypes = ref([
    "http",
    "https",
    "tcp",
    "udp",
    "stcp",
    "xtcp",
    "sudp",
]);

const stcpModels = ref([
    { label: "访问者", value: "visitors" },
    { label: "被访问者", value: "visited" },
]);

const handleSelectFile = (mod: number) => {
    OpenDirDialog().then((res: string) => {
        if (mod == 1) {
            proxyStore.proxyData.https2httpCaFile = res;
        } else if (mod == 2) {
            proxyStore.proxyData.https2httpKeyFile = res;
        } else if (mod == 3) {
            proxyStore.proxyData.localPath = res;
        }
    });
};


</script>

<template>
    <el-form :model="proxyStore.proxyData" :rules="proxyRules" ref="pwdRef" label-position="top">
        <!-- 代理类型选择 -->
        <el-form-item label="代理类型：" prop="type">
            <el-radio-group v-model="proxyStore.proxyData.type">
                <el-radio-button v-for="type in proxyTypes" :key="type" :value="type">{{ type }}</el-radio-button>
            </el-radio-group>
        </el-form-item>

        <!-- HTTP/HTTPS模式 -->
        <template v-if="
            proxyStore.proxyData.type === 'http' ||
            proxyStore.proxyData.type === 'https' ||
            proxyStore.proxyData.type === 'tcp' ||
            proxyStore.proxyData.type === 'udp' ||
            proxyStore.proxyData.type === 'stcp' ||
            proxyStore.proxyData.type === 'xtcp' ||
            proxyStore.proxyData.type === 'sudp'
        ">
            <el-form-item label="代理名称：" prop="name">
                <el-input v-model="proxyStore.proxyData.name" placeholder="代理名称" />
            </el-form-item>
            <el-row v-if="
                proxyStore.proxyData.type === 'http' ||
                proxyStore.proxyData.type === 'https' ||
                proxyStore.proxyData.type === 'tcp' ||
                proxyStore.proxyData.type === 'udp'
            " :gutter="20">
                <el-col :span="12">
                    <el-form-item label="内网Ip：" prop="localIp">
                        <el-input v-model="proxyStore.proxyData.localIp" placeholder="内网Ip地址" />
                    </el-form-item>
                </el-col>
                <el-col :span="8">
                    <el-form-item label="端口：" prop="localPort">
                        <el-input-number v-model="proxyStore.proxyData.localPort" :min="0" :max="65535" />
                    </el-form-item>
                </el-col>
            </el-row>

            <el-form-item v-if="
                proxyStore.proxyData.type === 'http' ||
                proxyStore.proxyData.type === 'https'
            " label="自定义域名：" prop="customDomain">
                <el-row v-for="(_, index) in proxyStore.customDomains" :key="index" :gutter="24">
                    <el-col :span="12">
                        <el-input v-model="proxyStore.customDomains[index]" placeholder="example.com" />
                    </el-col>
                    <el-col :span="5">
                        <el-button type="primary" icon="Plus" style="width: 80px"
                            @click="proxyStore.addCustomDomain">添加</el-button>
                    </el-col>
                    <el-col :span="5">
                        <el-button type="primary" icon="Plus" style="width: 80px" @click="
                            proxyStore.removeCustomDomain(index)
                            ">删除</el-button>
                    </el-col>
                </el-row>
            </el-form-item>
            <el-form-item v-if="
                proxyStore.proxyData.type === 'http' ||
                proxyStore.proxyData.type === 'https'
            " label="子域名：" prop="domain">
                <el-input v-model="proxyStore.proxyData.domain" placeholder="子域名" />
            </el-form-item>
            <el-form-item v-if="
                proxyStore.proxyData.type === 'tcp' ||
                proxyStore.proxyData.type === 'udp'
            " label="外网端口：" prop="remotePort">
                <el-input-number v-model="proxyStore.proxyData.remotePort" :min="0" :max="65535" />
            </el-form-item>
            <template v-if="proxyStore.proxyData.type === 'http' ||
                proxyStore.proxyData.type === 'https' ||
                proxyStore.proxyData.type === 'tcp'">
                <el-form-item label="HTTP基本认证：" prop="httpAuth">
                    <el-switch v-model="proxyStore.proxyData.httpAuth" />
                </el-form-item>
                <el-form-item v-if="proxyStore.proxyData.httpAuth" label="认证用户名：" prop="authUsername">
                    <el-input v-model="proxyStore.proxyData.authUsername" placeholder="username" />
                </el-form-item>
                <el-form-item v-if="proxyStore.proxyData.httpAuth" label="认证密码：" prop="authPassword">
                    <el-input v-model="proxyStore.proxyData.authPassword" type="password" placeholder="password" />
                </el-form-item>
            </template>
            <el-form-item v-if="proxyStore.proxyData.type === 'https'" label="证书文件：" prop="https2httpCaFile">
                <el-input v-model="proxyStore.proxyData.https2httpCaFile" placeholder="点击选择证书文件"
                    @click="handleSelectFile(1)" />
            </el-form-item>
            <el-form-item v-if="proxyStore.proxyData.type === 'https'" label="密钥文件：" prop="https2httpKeyFile">
                <el-input v-model="proxyStore.proxyData.https2httpKeyFile" placeholder="点击选择密钥文件"
                    @click="handleSelectFile(2)" />
            </el-form-item>
        </template>
        <template v-if="proxyStore.proxyData.type === 'tcp'">
            <el-form-item label="文件访问：" prop="staticFile">
                <el-switch v-model="proxyStore.proxyData.staticFile" />
            </el-form-item>
            <el-form-item v-if="proxyStore.proxyData.staticFile" label="文件夹路径：" prop="localPath">
                <el-input v-model="proxyStore.proxyData.localPath" placeholder="点击选择文件夹路径"
                    @click="handleSelectFile(3)" />
            </el-form-item>
            <el-form-item v-if="proxyStore.proxyData.staticFile" label="URL的前缀：" prop="stripPrefix">
                <el-input v-model="proxyStore.proxyData.stripPrefix" placeholder="URL的前缀" />
            </el-form-item>
        </template>


        <!-- STCP/XTCP/SUDP模式 -->
        <template v-if="
            proxyStore.proxyData.type === 'stcp' ||
            proxyStore.proxyData.type === 'xtcp' ||
            proxyStore.proxyData.type === 'sudp'
        ">
            <el-row :gutter="22">
                <el-col :span="14">
                    <el-form-item label="STCP模式：" prop="stcpModel">
                        <el-radio-group v-model="proxyStore.proxyData.stcpModel">
                            <el-radio v-for="model in stcpModels" :key="model.value" :value="model.value">{{ model.label
                                }}</el-radio>
                        </el-radio-group>
                    </el-form-item>
                </el-col>
                <el-col :span="10">
                    <el-form-item label="共享密钥：" prop="secretKey">
                        <el-input v-model="proxyStore.proxyData.secretKey" type="password" placeholder="密钥" />
                    </el-form-item>
                </el-col>
            </el-row>

            <!-- 被访问者代理名称 -->
            <el-form-item v-if="
                proxyStore.proxyData.type === 'stcp' ||
                proxyStore.proxyData.type === 'xtcp' ||
                proxyStore.proxyData.type === 'sudp'
            " label="被访问者代理名称：" prop="visitedName">
                <el-input v-model="proxyStore.proxyData.visitedName" placeholder="被访问者代理名称" />
            </el-form-item>

            <template v-if="
                proxyStore.proxyData.type === 'stcp' ||
                proxyStore.proxyData.type === 'xtcp' ||
                proxyStore.proxyData.type === 'sudp'
            ">
                <el-row :gutter="20">
                    <el-col :span="10">
                        <el-form-item label="绑定地址：" prop="bindAddr">
                            <el-input v-model="proxyStore.proxyData.bindAddr" placeholder="127.0.0.1" />
                        </el-form-item>
                    </el-col>
                    <el-col :span="10">
                        <el-form-item label="绑定端口：" prop="bindPort">
                            <el-input-number v-model="proxyStore.proxyData.bindPort" :min="1" :max="65535" />
                        </el-form-item>
                    </el-col>
                </el-row>
            </template>
            <template v-if="proxyStore.proxyData.type === 'xtcp'">
                <el-row :gutter="20">
                    <el-col :span="10">
                        <el-form-item label="回退代理名称：" prop="fallbackTo">
                            <el-input v-model="proxyStore.proxyData.fallbackTo
                                " placeholder="回退代理名称" />
                        </el-form-item>
                    </el-col>
                    <el-col :span="10">
                        <el-form-item label="回退超时毫秒：" prop="fallbackTimeoutMs">
                            <el-input-number v-model="proxyStore.proxyData
                                .fallbackTimeoutMs
                                " :min="0" />
                        </el-form-item>
                    </el-col>
                </el-row>
                <!-- 保持隧道开启 -->
                <el-form-item label="保持隧道开启：" prop="keepAlive">
                    <el-switch v-model="proxyStore.proxyData.keepAlive" />
                </el-form-item>
            </template>
        </template>

        <!-- 保存和取消按钮 -->
        <el-row justify="start">
            <el-button type="primary" @click="saveProxy" style="width: 100px">
                {{ proxyStore.isEditor ? '编辑' : '保存' }}
            </el-button>
            <el-button type="primary" style="width: 100px" @click="proxyDialogShow = false">取消</el-button>
        </el-row>
    </el-form>
</template>