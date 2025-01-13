<script setup lang="ts">
import { Plus } from '@element-plus/icons-vue'
import { ref } from "vue";
import { OpenDirDialog } from "@/util/goutil";
import { getSystemConfig } from '@/system/config';
import { notifySuccess, notifyError } from '@/util/msg';
interface ProxyItem {
    id: number;
    port: number;
    proxyType: string;
    domain: string;
    path?: string;
    status: boolean;
}
const proxyInit = {
    id: Date.now(),
    port: 8080,
    proxyType: "http",
    domain: "",
    path : "",
    status: true,
}
const config = getSystemConfig();
const proxies = ref<ProxyItem[]>([]);
const page = ref({
    current: 1,
    size: 10,
    total: 0,
})
const fetchProxies = () => {
    fetch(config.apiUrl + "/proxy/local/list").then(res => res.json()).then(res => {
        if (res.code === 0) {
            const data = res.data;
            if (data.proxies && Array.isArray(data.proxies)) {
                proxies.value = data.proxies;
                page.value.total = data.total;
            } else {
                console.error('Invalid data format:', data);
            }
        }

    }).catch(err => {
        console.error('Error fetching data:', err);
    });
};
const createProxies = (data: ProxyItem) => {
    fetch(config.apiUrl + '/proxy/local/create', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    }).then(res => {
        if (!res.ok) {
            notifyError('添加代理失败');
        } else {
            notifySuccess('添加代理成功');
            proxyData.value = proxyInit;
            proxyDialogShow.value = false;
            fetchProxies()
        }
    });
};
const updateProxies = (data: ProxyItem) => {
    fetch(config.apiUrl + '/proxy/local/update', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    }).then(res => {
        if (!res.ok) {
            notifyError('保存代理失败');
        } else {
            notifySuccess('保存代理成功');
            proxyData.value = proxyInit;
            proxyDialogShow.value = false;
            fetchProxies()
        }
    });
};
const DeleteProxy = (id: number) => {
    fetch(config.apiUrl + '/proxy/local/delete?id=' + id).then(res => {
        if (!res.ok) {
            notifyError('删除代理失败');
        } else {
            notifySuccess('删除代理成功');
            fetchProxies()
        }
    });
}
const SetStatus = (id: number) => {
    fetch(config.apiUrl + '/proxy/local/status?id=' + id).then(res => {
        if (!res.ok) {
            notifyError('设置代理状态失败');
        } else {
            notifySuccess('设置代理状态成功');
            fetchProxies()
        }
    });
}
const changePage = (current: number) => {
    page.value.current = current;
    fetchProxies();
};
onMounted(() => {
    fetchProxies();
});


const proxyData = ref<ProxyItem>(proxyInit);
const types = ref([
    { label: 'HTTP', value: 'http' },
    { label: 'Udp', value: 'udp' },
    { label: '静态文件访问', value: 'file' },
])
const proxyDialogShow = ref(false);
const isEditing = ref(false);
const pwdRef = ref<any>(null);


const editProxy = (proxy: ProxyItem) => {
    proxyData.value = { ...proxy };
    isEditing.value = true;
    proxyDialogShow.value = true;
};
const addProxy = () => {
    proxyData.value = proxyInit;
    isEditing.value = false;
    proxyDialogShow.value = true;
};

function selectFile() {
    OpenDirDialog().then((res: string) => {
        proxyData.value.domain = res;
    });
}
const saveProxy = () => {
    pwdRef.value.validate((valid: boolean) => {
        if (valid) {
            proxyData.value.port = Number(proxyData.value.port)
            if (isEditing.value) {
                updateProxies(proxyData.value)
            } else {
                createProxies(proxyData.value)
            }
        } else {
            console.log('表单验证失败');
        }
    });
};

const proxyRules = {
    port: [
        { required: true, message: '端口不能为空', trigger: 'blur' },
        { pattern: /^(6553[0-5]|655[0-2][0-9]|65[0-4][0-9]{2}|6[0-4][0-9]{3}|[1-5][0-9]{4}|[1-9][0-9]{1,3}|[0-9])$/, message: '请输入有效的端口号（1-65535）', trigger: 'blur' }
    ],
    domain: [
    { required: true, message: '代理域名不能为空', trigger: 'blur' },
    {
        pattern: /^(https?:\/\/)?((?:[a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}|localhost)(:\d{1,5})?(\/[^\s]*)?$/,
        message: '请输入有效的域名格式',
        trigger: 'blur'
    }
]
};


</script>
<template>
    <div>
        <el-row justify="end">
            <el-button type="primary" :icon="Plus" circle @click="addProxy" />
        </el-row>
        <el-table :data="proxies" style="width: 98%;border:none">
            <el-table-column prop="proxyType" label="代理类型" width="80" />
            <el-table-column prop="port" label="本地端口" width="80" />
            <!-- <el-table-column prop="domain" label="代理域名" /> -->
            
            <el-table-column label="状态">
                <template #default="scope">
                    <!-- <el-switch v-model="scope.row.status" active-color="#ff4949" inactive-color="#13ce66" @change="SetStatus(scope.row.id)"></el-switch> -->
                    <el-button size="small"  @click="SetStatus(scope.row.id)">{{scope.row.status ? '启用' : '禁用'}}</el-button>
                </template>
            </el-table-column>
            <el-table-column label="操作">
                <template #default="scope">
                    <el-button size="small" circle icon="Edit" @click="editProxy(scope.row)"></el-button>
                    <el-button size="small" circle icon="Delete"  @click="DeleteProxy(scope.row.id)"></el-button>
                </template>
            </el-table-column>
        </el-table>
        <el-pagination v-if="page.total > page.size" layout="prev, pager, next" :total="page.total"
            :page-size="page.size" v-model:current-page="page.current" @current-change="changePage" />
        <el-dialog v-model="proxyDialogShow" :title="isEditing ? '编辑代理' : '添加代理'" width="400px">
            <span>
                <el-form :model="proxyData" :rules="proxyRules" ref="pwdRef">
                    <el-form-item label="代理类型" prop="type">
                        <el-select v-model="proxyData.proxyType" placeholder="代理类型">
                            <el-option v-for="type in types" :key="type.value" :label="type.label"
                                :value="type.value" />
                        </el-select>
                    </el-form-item>
                    <el-form-item label="本地端口" prop="port">
                        <el-input v-model="proxyData.port" />
                    </el-form-item>
                    <el-form-item label="状态" prop="status">
                        <el-switch v-model="proxyData.status" active-color="#13ce66" inactive-color="#ff4949"
                            active-text="启用" inactive-text="禁用" />
                    </el-form-item>
                    <el-form-item label="代理域名" prop="domain" v-if="proxyData.proxyType === 'http'">
                        <el-input v-model="proxyData.domain" />
                    </el-form-item>
                    <el-form-item label="文件路径" prop="path" v-if="proxyData.proxyType === 'file'">
                        <el-input v-model="proxyData.path" @click="selectFile()" />
                    </el-form-item>
                    <el-form-item label="转发IP+端口" prop="domain" v-if="proxyData.proxyType === 'udp'">
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