<template>
  <div class="container">
    <div class="nav">
      <ul>
        <li v-for="(item, index) in items" :key="index" @click="selectItem(index)"
          :class="{ active: index === activeIndex }">
          {{ item }}
        </li>
      </ul>
    </div>
    <div class="setting">
      <div v-if="1 === activeIndex">
        <div class="setting-item" style="margin-top: 60px">
          <label>存储方式</label>
          <el-select v-model="config.storeType">
            <el-option v-for="(item, key) in storeList" :key="key" :label="item.title" :value="item.value" />
          </el-select>
        </div>
        <div class="setting-item" v-if="config.storeType === 'local'">
          <label>存储地址</label>
          <el-input v-model="config.storePath" @click="selectFile()" placeholder="可为空，为空则取系统默认存储地址" />
        </div>
        <template v-if="config.storeType === 'net'">
          <div class="setting-item">
            <label>服务器地址</label>
            <el-input v-model="config.storenet.url" placeholder="http://192.168.1.16:56780 不要加斜杠" />
          </div>
        </template>
        <template v-if="config.storeType === 'webdav'">
          <div class="setting-item">
            <label>服务器地址</label>
            <el-input v-model="config.webdavClient.url" placeholder="https://godoos.com/webdav 不要加斜杠" />
          </div>
          <div class="setting-item">
            <label>登陆用户名</label>
            <el-input v-model="config.webdavClient.username" />
          </div>
          <div class="setting-item">
            <label>登陆密码</label>
            <el-input v-model="config.webdavClient.password" type="password" />
          </div>
        </template>

        <div class="setting-item">
          <label></label>
          <el-button @click="submitOsInfo" type="primary">
            {{ t("confirm") }}
          </el-button>
        </div>
        <SetFilePwd  v-if="config.userType === 'person'"></SetFilePwd>
      </div>

      <div v-if="2 === activeIndex">
        <div class="setting-item">
          <h1 class="setting-title">备份</h1>
        </div>
        <div class="setting-item">
          <label></label>
          <el-button @click="exportBackup" type="primary"> 导出 </el-button>
        </div>
        <div class="setting-item">
          <h1 class="setting-title">还原</h1>
        </div>
        <div class="setting-item">
          <label></label>
          <el-button @click="selectZipfile" type="primary"> 导入 </el-button>
          <input type="file" accept=".zip" style="display: none" ref="zipFileInput" />
        </div>
      </div>
      <div v-if="0 === activeIndex">
        <div class="setting-item" style="margin-top: 60px">
          <label>用户角色</label>
          <el-select v-model="config.userType">
            <el-option v-for="(item, key) in userTypes" :key="key" :label="item.title" :value="item.value" />
          </el-select>
        </div>
        <template v-if="config.userType === 'member'">
          <div class="setting-item">
            <label>服务器地址</label>
            <el-input v-model="config.userInfo.url" placeholder="网址或域名，例子：https://godoos.com 不要加斜杠" />
          </div>
          <div class="setting-item">
            <label>用户名</label>
            <el-input v-model="config.userInfo.username" placeholder="登录用户名" />
          </div>
          <div class="setting-item">
            <label>密码</label>
            <el-input v-model="config.userInfo.password" type="password" placeholder="登录密码" />
          </div>
        </template>
        <div class="setting-item">
          <label></label>
          <el-button @click="saveUserInfo" type="primary">
            {{ t("confirm") }}
          </el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { inject, ref } from "vue";
import { Dialog, join, System, t } from "@/system";
import JSZip from "jszip";
import FileSaver from "file-saver";
import { getSystemConfig, setSystemConfig,getClientId } from "@/system/config";
import { OpenDirDialog, RestartApp } from "@/util/goutil";
import { notifyError, notifySuccess } from "@/util/msg";
const config = ref(getSystemConfig());
const sys = inject<System>("system")!;
let zipFile: File | undefined = undefined;
const zipFileInput = ref();
const fileName: any = ref("");
const storeList = [
  // {
  //   title: "浏览器存储",
  //   value: "browser",
  // },
  {
    title: "本地存储",
    value: "local",
  },
  {
    title: "远程存储",
    value: "net",
  },
  {
    title: "webdav",
    value: "webdav",
  },
];

const items = ["用户角色","个人存储", "备份还原"];
const urlRegex = /^(https?:\/\/)/;
const userTypes = [
  {
    title: "独立用户",
    value: "person",
  },
  {
    title: "企业用户",
    value: "member",
  },
]
const activeIndex = ref(0);

const selectItem = (index: number) => {
  activeIndex.value = index;
};
function selectFile() {
  OpenDirDialog().then((res: string) => {
    config.value.storePath = res;
  });
}

function submitOsInfo() {
  const saveData = toRaw(config.value);
  const postData: any = {
    //name: "osPath",
    type: saveData.storeType,
  };
  if (saveData.storeType === "browser") {
    setSystemConfig(saveData);
    RestartApp();
    return;
  }
  if (saveData.storeType === "local") {
    if (saveData.storePath === "") {
      setSystemConfig(saveData);
      RestartApp();
      return;
    }
    postData.name = "osPath";
    postData.value = saveData.storePath;
    const postUrl = config.value.apiUrl + "/system/setting";
    fetch(postUrl, {
      method: "POST",
      body: JSON.stringify([postData]),
    })
      .then((res) => res.json())
      .then((res) => {
        if (res.code === 0) {
          setSystemConfig(saveData);
          RestartApp();
        } else {
          Dialog.showMessageBox({
            message: res.message,
            type: "error",
          });
        }
      });
  }
  if (saveData.storeType === "net") {
    if (saveData.storenet.url === "") {
      Dialog.showMessageBox({
        message: "服务器地址不能为空",
        type: "error",
      });
      return;
    }
    const urlRegex = /^(https?:\/\/)[^\/]+$/;
    if (!urlRegex.test(saveData.storenet.url)) {
      Dialog.showMessageBox({
        message: "服务器地址格式错误",
        type: "error",
      });
      return;
    }
    setSystemConfig(saveData);
    RestartApp();
  }
  if (saveData.storeType === "webdav") {
    const urlRegex = /^(https?:\/\/)/;
    if (!urlRegex.test(saveData.webdavClient.url.trim())) {
      Dialog.showMessageBox({
        message: "服务器地址格式错误",
        type: "error",
      });
      return;
    }
    if (saveData.webdavClient.username === "" || saveData.webdavClient.password === "") {
      Dialog.showMessageBox({
        message: "用户名或密码不能为空",
        type: "error",
      });
      return;
    }
    postData.name = "webdavClient";
    postData.value = saveData.webdavClient;
    const postUrl = config.value.apiUrl + "/system/setting";
    fetch(postUrl, {
      method: "POST",
      body: JSON.stringify([postData]),
    })
      .then((res) => res.json())
      .then((res) => {
        if (res.code === 0) {
          setSystemConfig(saveData);
          RestartApp();
        } else {
          Dialog.showMessageBox({
            message: res.message,
            type: "error",
          });
        }
      });
  }
}
async function saveUserInfo() {
  const saveData = toRaw(config.value);
  if (saveData.userType == 'person') {
    setSystemConfig(saveData);
    notifySuccess("保存成功");
    RestartApp();
    return
  }
  if (!urlRegex.test(saveData.userInfo.url.trim())) {
    Dialog.showMessageBox({
      message: "服务器地址格式错误",
      type: "error",
    });
    return;
  }
  if (saveData.userInfo.username === "" || saveData.userInfo.password === "") {
    Dialog.showMessageBox({
      message: "用户名或密码不能为空",
      type: "error",
    });
    return;
  }
  const password = saveData.userInfo.password
  const serverUrl = saveData.userInfo.url + '/member/login'
  const res = await fetch(serverUrl, {
    method: "POST",
    body: JSON.stringify({
      username: saveData.userInfo.username,
      password: password,
      clientId: getClientId(),
    }),
  });
  if (res.status === 200) {
    const data = await res.json();
    if (data.success) {
      data.data.url = saveData.userInfo.url
      data.data.password = password
      saveData.userInfo = data.data
      setSystemConfig(saveData);
      notifySuccess("登录成功");
      RestartApp();
      return
    }else{
      notifyError(data.message)
    }
  }else{
    notifyError("登录失败");
  }
  return
}
async function exportBackup() {
  const { setProgress } = Dialog.showProcessDialog({
    message: `正在打包`,
  });
  try {
    const zip = new JSZip();
    await dfsPackage("/", zip, setProgress);
    zip.generateAsync({ type: "blob" }).then(function (content) {
      FileSaver.saveAs(content, "backup.zip");
      setProgress(100);
    });
  } catch (error) {
    //console.log(error);
    Dialog.showMessageBox({
      message: "打包失败",
      type: "error",
    }).finally(() => {
      zipFile = undefined;
      fileName.value = "";
      setProgress(100);
    });
  }
}
async function dfsPackage(path: string, zip: JSZip, setProgress: any) {
  const dir = await sys.fs.readdir(path);

  for (let i = 0; i < dir.length; i++) {
    const item = dir[i];
    const stat = await sys.fs.stat(item.path);
    setProgress(Math.max((i / dir.length) * 100 - 0.1, 0.1));
    if (stat) {
      if (stat.isDirectory) {
        await dfsPackage(item.path, zip, setProgress);
      } else {
        const content = await sys.fs.readFile(item.path);
        try {
          atob(content || "");
          zip.file(item.path, content || "");
        } catch (error) {
          zip.file(item.path, content || "");
        }
      }
    }
  }
}
function selectZipfile() {
  zipFileInput.value.click();
  zipFileInput.value.onchange = (e: any) => {
    const tar: any = e.target as HTMLInputElement;
    if (tar.files?.[0]) {
      zipFile = tar.files[0];
      fileName.value = zipFile?.name;
      importBackup();
    } else {
      zipFile = undefined;
      fileName.value = "";
    }
  };
}
async function importBackup(path = "") {
  if (!zipFile) {
    return;
  }
  const { setProgress } = Dialog.showProcessDialog({
    message: "正在恢复备份",
  });
  try {
    const unziped = await JSZip.loadAsync(zipFile);
    const unzipArray: Array<JSZip.JSZipObject> = [];
    unziped.forEach((_, zipEntry) => {
      unzipArray.push(zipEntry);
    });
    for (let i = 0; i < unzipArray.length; i++) {
      const zipEntry = unzipArray[i];
      setProgress((i / unzipArray.length) * 100);
      if (zipEntry.dir) {
        await sys.fs.mkdir(join(path, zipEntry.name));
      } else {
        let fileC: any = await zipEntry.async("string");
        //console.log(fileC);
        if (!fileC.startsWith("link::")) {
          fileC = await zipEntry.async("arraybuffer");
        }
        sys.fs.writeFile(join(path, zipEntry.name), fileC);
      }
    }
    setProgress(100);
    selectItem(1);
  } catch (e) {
    //console.log(e);
    setTimeout(() => {
      Dialog.showMessageBox({
        message: "恢复失败",
        type: "error",
      }).finally(() => {
        zipFile = undefined;
        fileName.value = "";
        setProgress(100);
      });
    }, 100);
  }
}


</script>
<style scoped>
@import "./setStyle.css";

.ctrl {
  width: 100px;
}

.setting-item {
  display: flex;
  align-items: center;
}
</style>
