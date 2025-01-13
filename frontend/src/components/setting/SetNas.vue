<script setup lang="ts">
import { Plus } from '@element-plus/icons-vue'
import { ref, computed } from "vue";
interface ProxyItem {
    id: number;
    server_url: string;
    disk: string;
    username: string;
    password: string;
}
const localKey = "godoos_local_nasclient"
const getProxies = (): ProxyItem[] => {
    const proxies = localStorage.getItem(localKey);
    return proxies ? JSON.parse(proxies) : [];
};

const saveProxies = (proxies: ProxyItem[]) => {
    localStorage.setItem(localKey, JSON.stringify(proxies));
};

const proxies = ref<ProxyItem[]>(getProxies());
const initData = {
    id: Date.now(),
    server_url: "",
    disk: "",
    username: "",
    password: "",
}
const proxyData = ref<ProxyItem>(initData);
const proxyDialogShow = ref(false);
const isEditing = ref(false);
const pwdRef = ref<any>(null);

const addProxy = () => {
    if (pwdRef.value.validate()) {
        proxies.value.push({ ...proxyData.value });
        saveProxies(proxies.value);
        proxyDialogShow.value = false;
        proxyData.value = initData;
    }
};

const editNas = (proxy: ProxyItem) => {
    proxyData.value = { ...proxy };
    isEditing.value = true;
    proxyDialogShow.value = true;
};

const updateProxy = () => {
    if (pwdRef.value.validate()) {
        const index = proxies.value.findIndex(p => p.id === proxyData.value.id);
        if (index !== -1) {
            proxies.value[index] = { ...proxyData.value };
            saveProxies(proxies.value);
            proxyDialogShow.value = false;
            proxyData.value = initData;
            isEditing.value = false;
        }
    }
};
const deleteNas = (id: number) => {
    proxies.value = proxies.value.filter(p => p.id !== id);
    saveProxies(proxies.value);
};

const saveNas = () => {
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
    username: [
        { required: true, message: '请输入用户名', trigger: 'blur' }
    ],
    password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
        { min: 6, message: '密码长度至少为6位', trigger: 'blur' }
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
    <div style="padding:15px">
        <el-row justify="end">
            <el-button type="primary" :icon="Plus" circle @click="proxyDialogShow = true" />
        </el-row>
        <el-table :data="paginatedProxies" style="width: 98%;border:none">
            <el-table-column prop="disk" label="挂载盘符" width="180" />
            <el-table-column prop="server_url" label="类型" width="180" />
            <el-table-column label="操作">
                <template #default="scope">
                    <el-button size="small" @click="editNas(scope.row)">编辑</el-button>
                    <el-button size="small" type="danger" @click="deleteNas(scope.row.id)">删除</el-button>
                </template>
            </el-table-column>
        </el-table>
        <el-pagination v-if="totalPages > 1" layout="prev, pager, next" :total="getProxies().length"
            :page-size="pageSize" v-model:current-page="currentPage" @next-click="nextPage" @prev-click="prevPage" />
        <el-dialog v-model="proxyDialogShow" :title="isEditing ? '编辑Nas' : '添加Nas'" width="400px">
            <span>
                <el-form :model="proxyData" :rules="proxyRules" ref="pwdRef">
                    <el-form-item label="挂载" prop="disk">
                        <el-input v-model="proxyData.disk"/>
                    </el-form-item>
                    <el-form-item label="地址" prop="server_url">
                        <el-input v-model="proxyData.server_url" />
                    </el-form-item>
                    <el-form-item label="用户" prop="username">
                        <el-input v-model="proxyData.username" />
                    </el-form-item>
                    <el-form-item label="密码" prop="password">
                        <el-input v-model="proxyData.password" type="password" />
                    </el-form-item>
                    <el-form-item>
                        <el-button type="primary" @click="saveNas" style="margin: 0 auto;">
                            确认
                        </el-button>
                    </el-form-item>
                </el-form>
            </span>
        </el-dialog>
    </div>
</template>