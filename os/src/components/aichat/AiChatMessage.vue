<script setup lang="ts">
import { renderMarkdown } from "@/utils/markdown.ts";
import moment from "moment";
import { computed, ref } from "vue";
import { ElMessageBox } from "element-plus";
import { useFileSystemStore } from "@/stores/filesystem.ts";
const fs = useFileSystemStore()

interface Props {
  content: string;
  role?: string;
  link?: any;
  createdAt: number | string;
  doc: any[];
  web_search: any[];
}

const props = defineProps<Props>();
//console.log(props);
const sourceExpand = ref(false);
const openFile = (filePath: string) => {
  fs.openFile(filePath);
};
const openlink = (item: any) => {
  if (!item.link) return ElMessageBox.alert(
    item.content,
    {
      confirmButtonText: '确定',
    }
  )
  window.open(item.link)
}
const getDateTime = (t: any) => {
  return moment(t).format("MM-DD HH:mm");
};
const showLink = ref(false);
const showImg = ref(false);
function handleFlag() {
  showLink.value = !showLink.value;
  if (showLink.value) {
    showImg.value = false;
  }
}
function handleLink() {
  showImg.value = !showImg.value;
  if (showImg.value) {
    showLink.value = false;
  }
}
const refbox = ref(false);
const detailTxt = ref("");
const detailImg = ref("");
const imgList = computed<any[]>(() => {
  if (props.link && props.link.length > 0) {
    const list: any = [];
    const url = "user/chat/showimage?path=";
    props.link.forEach((item: any) => {
      if (item.metadata.type === "image") {
        list.push(url + item.metadata.file);
      }
    });
    return list;
  }
  return [];
});
async function showDetail(file: string, type: string) {
  const filePath = "user/chat/filedetail?path=" + file
  if (type != "image") {
    const reponse = await fetch(filePath);
    if (reponse.ok) {
      const data = await reponse.text();
      refbox.value = true;
      detailTxt.value = data;
      detailImg.value = "";
    }
  } else {
    refbox.value = true;
    detailTxt.value = '';
    detailImg.value = filePath;
  }

}

</script>
<template>
  <el-dialog v-model="refbox" width="80%">
    <el-scrollbar style="height: 70vh;padding:20px;">
      <el-text v-if="detailTxt != ''">{{ detailTxt }}</el-text>
      <img v-if="detailImg != ''" :src="detailImg" alt="Image" style="width: 100%" />
    </el-scrollbar>
  </el-dialog>
  <div>
    <!-- 助手信息 -->
    <div v-if="props.role === 'assistant'" class="assistant align-start spacing">
      <div class="assistant-avatar">
        <div class="icon-container">
          <el-icon>
            <GoldMedal />
          </el-icon>
        </div>
      </div>
      <div class="message assistant-message rounded-al">
        <div v-html="renderMarkdown(props.content)"></div>
        <el-row type="flex" justify="space-between" style="border-bottom:none">
          <div class="text-grey">
            {{ getDateTime(props.createdAt) }}
          </div>
          <div>
            <el-button icon="InfoFilled" circle size="small" v-if="props.link" @click="handleFlag" />
            <el-button icon="PictureFilled" circle size="small" v-if="imgList.length > 0" @click="handleLink" />
          </div>
        </el-row>
        <el-card v-if="showLink" style="color: black">
          <el-row type="flex" justify="space-between" align="middle" :gutter="24" style="margin-top: 10px"
            v-for="(item, key) in props.link">
            <el-col :span="2">
              <el-avatar size="small">
                {{ key + 1 }}
              </el-avatar>
            </el-col>
            <el-col :span="22" @click="showDetail(item.metadata.file, item.metadata.type)">
              {{ item.metadata.category }}
            </el-col>
          </el-row>
        </el-card>
        <el-card v-if="showImg" style="color: black">
          <div class="image-grid">
            <div v-for="(item, index) in imgList" :key="index" class="image-item">
              <el-image :src="item" alt="Image" :zoom-rate="1.2" :max-scale="7" :min-scale="0.2"
                :preview-src-list="imgList" :initial-index="index" fit="cover" class="image" />
            </div>
          </div>
        </el-card>

        <!-- 参考源 -->
        <div class="reference-source" v-if="props.doc.length || props.web_search.length">
          <div class="source-title" @click="sourceExpand = !sourceExpand" style="cursor: pointer;">
            <el-button type="default" style="margin-right: 5px; width: 20px;height: 20px;" icon="Document" circle />
            <span style="margin-right: 5px">{{ props.doc.length || props.web_search.length }}篇资料作为参考</span>
            <ArrowDown v-if="sourceExpand" class="icon"></ArrowDown>
            <ArrowRight v-else class="icon"></ArrowRight>
          </div>
          <!-- doucument -->
          <div v-if="props.doc.length">
            <ul class="source-list" v-if="sourceExpand" v-for="(item, i) in props.doc" :key="i">
              <li class="source-item" @click="openFile(item.file_path)">
                {{ item.file_name }}</li>
            </ul>
          </div>
          <!-- web_search -->
          <div v-else>
            <ul class="source-list" v-if="sourceExpand" v-for="item in props.web_search" :key="item.title">
              <li class="source-item" @click="openlink(item)">
                <img :src="item.icon" v-if="item.icon" style="width: 15px;height: 15px; margin-right: 5px;">{{
                  item.title
                }}
              </li>
            </ul>
          </div>
        </div>
      </div>
    </div>
    <!-- 用户信息 -->
    <div v-else-if="props.role === 'user'" class="user align-center spacing">
      <div class="message-container">
        <div class="message user-message rounded-xl" style="max-height: 80px; overflow-y: auto">
          <div>{{ props.content }}</div>
        </div>
      </div>
      <div class="avatar">
        <div class="icon-container">
          <el-icon>
            <UserFilled />
          </el-icon>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
// Windows 10 配色
$win10-blue: #0078d7;
$win10-light-blue: #c7e8ff;
$win10-grey: #afafaf;
$win10-light-grey: #f2f2f2;

.message-container {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  width: 99%;
}

.message {
  margin-top: 10px;
  padding: 8px;
  max-width: 80%;
  font-size: 1em;
  line-height: 1.5em;
  color: $win10-light-grey;
  text-align: left;
  margin: 8px;
}

.text-grey {
  color: $win10-grey;
}

.user-message {
  background-color: $win10-blue;
  border-radius: 12px 0 0 12px;
  font-weight: 600;
  font-size: 14px;
  color: $win10-light-blue;
}

.assistant-message {
  flex: 1;
  background-color: #fff;
  border-radius: 0 12px 12px 0;
  font-weight: 600;
  font-size: 14px;
  color: #333;
  box-shadow: rgba(100, 100, 111, 0.2) 0px 7px 29px 0px;
  padding: 10px 20px;
}

.avatar {
  display: flex;
  align-items: center;
  width: 50px;
  height: 50px;
}

.assistant-avatar {
  display: flex;
  align-items: center;
  width: 35px;
  height: 35px;
  padding-left: 10px;
}

.icon-container {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  border-radius: 50%;
  overflow: hidden;
  border: 1px solid $win10-grey;
}

.rounded-xl {
  border-radius: 12px 0px 12px 12px;
}

.rounded-al {
  border-radius: 0px 12px 12px 12px;
}

.align-start,
.align-center {
  display: flex;
  gap: 3px;
}

.align-center {
  justify-content: flex-end;
}

.image-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
}

.image-item {
  position: relative;
  border: 1px solid white;
  /* 白边效果 */
  overflow: hidden;
}

.image {
  width: 100%;
  height: auto;
}

.reference-source {
  font-size: 13px;
  color: #9b9b9b;
  margin-top: 10px;

  .source-title {
    display: flex;
    align-items: center;
  }

  .source-list {
    list-style: disc inside none;
    margin-top: 5px;

    .source-item {
      background-color: #ecebeb;
      height: 30px;
      padding: 0 10px;
      line-height: 30px;
      border-radius: 15px;
      cursor: pointer;
      overflow: hidden;
    }
  }
}
</style>
