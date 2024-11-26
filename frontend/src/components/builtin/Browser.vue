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
    <iframe
      ref="iframeRef"
      class="webframe"
      :srcdoc="webContent"
      frameborder="0"
      allowfullscreen
      @load="onIframeLoad"
    ></iframe>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, onUnmounted, Ref } from 'vue';
import { getSystemKey } from '@/system/config';

const urlinput = ref('https://godoos.com/');
const webContent = ref('');
const canGoBack = ref(false);
const canGoForward = ref(false);
const iframeRef: Ref<HTMLIFrameElement | null> = ref(null);

async function callApi(endpoint: string) {
  const apiUrl = await getSystemKey('apiUrl');
  let fetchUrl = `${apiUrl}/ie/${endpoint}?url=${encodeURIComponent(urlinput.value)}`;
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

function onIframeLoad() {
  if (iframeRef.value) {
    let currentUrl = '';
    if (iframeRef.value.srcdoc) {
      // 当使用 srcdoc 时，从外部获取实际 URL
      currentUrl = urlinput.value;
    } else {
      // 当使用 src 时，从 contentWindow 获取实际 URL
      currentUrl = iframeRef.value.contentWindow?.location.href || '';
    }
    console.log('Current URL:', currentUrl);
    //urlinput.value = currentUrl;
    callApi('navigate'); // 每次 iframe 加载时调用 callApi
    updateNavigationButtons();
  }
}
function handleIframeMessage(event: MessageEvent) {
  if (event.origin !== 'https://your-iframe-origin') return; // 替换为实际的 origin
  const currentUrl = event.data;
  console.log('Current URL from iframe:', currentUrl);
  urlinput.value = currentUrl;
  updateNavigationButtons();
}
onMounted(() => {
  window.addEventListener('popstate', updateNavigationButtons);
  window.addEventListener('message', handleIframeMessage);

  updateNavigationButtons();
  callApi('navigate');
});

onUnmounted(() => {
  window.removeEventListener('popstate', updateNavigationButtons);
  window.removeEventListener('message', handleIframeMessage);
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