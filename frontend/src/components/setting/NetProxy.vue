<script setup lang="ts">
import { notifyError } from '@/util/msg';
import { Plus } from '@element-plus/icons-vue'
import { ref, computed } from "vue";

interface ProxyItem {
    id: number;
    port: string;
    domain: string;
    type:string;
    name:string;
}
const localKey = "godoos_net_proxy"
const getProxies = (): ProxyItem[] => {
    const proxies = localStorage.getItem(localKey);
    return proxies ? JSON.parse(proxies) : [];
};

const saveProxies = (proxies: ProxyItem[]) => {
    localStorage.setItem(localKey, JSON.stringify(proxies));
};

const proxies = ref<ProxyItem[]>(getProxies());
const proxyData = ref<ProxyItem>({
    id: Date.now(),
    type:"http",
    name:"",
    port: "",
    domain: "",
});
const proxyDialogShow = ref(false);
const isEditing = ref(false);
const pwdRef = ref<any>(null);
const types = ref([
    { label: 'HTTP', value: 'http' },
    { label: 'SOCKS5', value: 'socks5' },
])

const addProxy = () => {
    if (pwdRef.value.validate()) {
        const isNameDuplicate = proxies.value.some(p => p.name === proxyData.value.name);
        if (isNameDuplicate) {
            notifyError("代理名称已存在");
            return;
        }
        proxies.value.push({ ...proxyData.value });
        saveProxies(proxies.value);
        proxyDialogShow.value = false;
        proxyData.value = { id: Date.now(), port: "", domain: "",type:"http",name:"" };
    }
};

const editProxy = (proxy: ProxyItem) => {
    proxyData.value = { ...proxy };
    isEditing.value = true;
    proxyDialogShow.value = true;
};

const updateProxy = () => {
    if (pwdRef.value.validate()) {
        const index = proxies.value.findIndex(p => p.id === proxyData.value.id);
        if (index !== -1) {
            const isNameDuplicate = proxies.value.some((p, i) => p.name === proxyData.value.name && i !== index);
            if (isNameDuplicate) {
                notifyError("代理名称已存在");
                return;
            }
            proxies.value[index] = { ...proxyData.value };
            saveProxies(proxies.value);
            proxyDialogShow.value = false;
            proxyData.value = { id: Date.now(), port: "", domain: "",type:"http",name:"" };
            isEditing.value = false;
        }
    }
};

const deleteProxy = (id: number) => {
    proxies.value = proxies.value.filter(p => p.id !== id);
    saveProxies(proxies.value);
};

const saveProxy = () => {
    pwdRef.value.validate((valid: boolean) => {
        if (valid) {
            if (isEditing.value) {
                updateProxy();
            } else {
                addProxy();
            }
        } else {
            console.log('表单验证失败');
        }
    });
};

const proxyRules = {
    name: [
        { required: true, message: '代理名称不能为空', trigger: 'blur' },
        { min: 1, max: 20, message: '代理名称长度在 1 到 20 个字符', trigger: 'blur' },
        { pattern: /^[a-zA-Z]+$/, message: '代理名称只能包含英文字符', trigger: 'blur' }
    ],
    port: [
        { required: true, message: '端口不能为空', trigger: 'blur' },
        { pattern: /^(6553[0-5]|655[0-2][0-9]|65[0-4][0-9]{2}|6[0-4][0-9]{3}|[1-5][0-9]{4}|[1-9][0-9]{1,3}|[0-9])$/, message: '请输入有效的端口号（1-65535）', trigger: 'blur' }
    ],
    domain: [
        { required: true, message: '代理域名不能为空', trigger: 'blur' },
        { pattern: /^(https?:\/\/)?(?:[a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}(:\d{1,5})?(\/[^\s]*)?$/, message: '请输入有效的域名格式', trigger: 'blur' }
    ]
};

const pageSize = 10;
const currentPage = ref(1);

const paginatedProxies = computed(() => {
    const start = (currentPage.value - 1) * pageSize;
    const end = start + pageSize;
    return proxies.value.slice(start, end);
});

const totalPages = computed(() => Math.ceil(proxies.value.length / pageSize));

const nextPage = () => {
    if (currentPage.value < totalPages.value) {
        currentPage.value++;
    }
};

const prevPage = () => {
    if (currentPage.value > 1) {
        currentPage.value--;
    }
};
</script>
<template>
    <div>
        <el-row justify="end">
            <el-button type="primary" :icon="Plus" circle @click="proxyDialogShow = true" />
        </el-row>
        <el-table :data="paginatedProxies" style="width: 98%;border:none">
            <el-table-column prop="name" label="名称" width="100" />
            <el-table-column prop="type" label="类型" width="80" />
            <el-table-column prop="port" label="本地端口" width="80" />
            <el-table-column prop="domain" label="代理域名" width="180" />
            <el-table-column label="操作">
                <template #default="scope">
                    <el-button size="small" @click="editProxy(scope.row)">编辑</el-button>
                    <el-button size="small" type="danger" @click="deleteProxy(scope.row.id)">删除</el-button>
                </template>
            </el-table-column>
        </el-table>
        <el-pagination v-if="totalPages > 1" layout="prev, pager, next" :total="getProxies().length"
            :page-size="pageSize" v-model:current-page="currentPage" @next-click="nextPage" @prev-click="prevPage" />
        <el-dialog v-model="proxyDialogShow" :title="isEditing ? '编辑代理' : '添加代理'" width="400px">
            <span>
                <el-form :model="proxyData" :rules="proxyRules" ref="pwdRef">
                    <el-form-item label="代理类型" prop="type">
                        <el-select v-model="proxyData.type" placeholder="代理类型">
                            <el-option v-for="type in types" :key="type.value" :label="type.label" :value="type.value" />
                        </el-select>
                    </el-form-item>
                    <el-form-item label="代理名称" prop="name">
                        <el-input v-model="proxyData.name" />
                    </el-form-item>
                    
                    <el-form-item label="本地端口" prop="port">
                        <el-input v-model="proxyData.port" />
                    </el-form-item>
                    <el-form-item label="代理域名" prop="domain">
                        <el-input v-model="proxyData.domain" />
                    </el-form-item>
                    <el-form-item>
                        <el-button type="primary" @click="saveProxy" style="margin: 0 auto;">
                            确认
                        </el-button>
                    </el-form-item>
                </el-form>
            </span>
        </el-dialog>
    </div>
</template>