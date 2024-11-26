<template>
  <div class="outer win11-theme">
    <div class="controls">
      <el-icon class="change-button" @click="goBack" :disabled="!canGoBack">
        <ArrowLeft />
      </el-icon>
      <el-icon class="change-button" @click="goForward" :disabled="!canGoForward">
        <ArrowRight />
      </el-icon>
      <el-icon class="change-button" @click="refreshPage">
        <Refresh />
      </el-icon>

      <input class="url-input" v-model="urlinput" @keydown.enter="changeUrl" placeholder="Enter URL" />
    </div>
    <!-- <div class="webframe" v-html="webContent"></div> -->
    <iframe
      class="webframe"
      :srcdoc="webContent"
      frameborder="0"
      allowfullscreen
    ></iframe>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, onUnmounted } from 'vue';
import { getSystemKey } from '@/system/config';

const urlinput = ref('https://godoos.com/');
const webContent = ref('');
const canGoBack = ref(false);
const canGoForward = ref(false);

async function callApi(endpoint: string) {
  const apiUrl = await getSystemKey('apiUrl');
  let fetchUrl = `${apiUrl}/ie/${endpoint}?url=${encodeURIComponent(urlinput.value)}`
  const response = await fetch(fetchUrl);
  if (response.ok) {
    const data = await response.text();
    webContent.value = data;
  } else {
    console.log('Network response was not ok');
  }
}

function updateNavigationButtons() {
  canGoBack.value = window.history.length > 1;
  canGoForward.value = window.history.length < window.history.state?.length;
}

function goBack() {
  if (canGoBack.value) {
    callApi('back');
  }
}

function goForward() {
  if (canGoForward.value) {
    callApi('forward');
  }
}

function refreshPage() {
  callApi('refresh');
}

function changeUrl() {
  const protocolAdded = /^(http|https):\/\//.test(urlinput.value) ? urlinput.value : `https://${urlinput.value}`;
  urlinput.value = protocolAdded;
  callApi('navigate');
}

onMounted(() => {
  window.addEventListener('popstate', updateNavigationButtons);
  updateNavigationButtons();
  callApi('navigate');
});

onUnmounted(() => {
  window.removeEventListener('popstate', updateNavigationButtons);
});
</script>

<style scoped>
.win11-theme {
  display: flex;
  flex-direction: column;
  height: 100%;
  width: 100%;
  background-color: #f3f4f6;
}

.controls {
  display: flex;
  align-items: center;
  padding: 10px;
  background-color: #ffffff;
  border-bottom: 1px solid #e5e5e5;
}

.change-button {
  margin: 0 5px;
  border-radius: 8px;
  border: none;
  background-color: transparent;
  cursor: pointer;
  padding: 4px;
  transition: background-color 0.2s;
}

.change-button:hover {
  background-color: rgba(0, 0, 0, 0.05);
}

.url-input {
  flex-grow: 1;
  margin: 0 10px;
  padding: 8px 12px;
  border-radius: 8px;
  border: 1px solid #dee2e6;
  background-color: #ffffff;
  font-size: 14px;
  outline: none;
  transition: border-color 0.2s, background-color 0.2s;
}

.url-input:focus {
  border-color: #0078d4;
  background-color: #ffffff;
}

.webframe {
  flex: 1;
  width: 100%;
  border: none;
  overflow-y: auto; /* 添加滚动条 */
}
</style>