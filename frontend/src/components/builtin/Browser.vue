<template>
  <div class="outer win11-theme">
    <div class="controls">
      <button class="change-button" @click="goBack" :disabled="!canGoBack">
        <svg
          t="1632984723698"
          class="icon"
          viewBox="0 0 1024 1024"
          version="1.1"
          xmlns="http://www.w3.org/2000/svg"
          p-id="10100"
          width="20"
          height="20"
        >
          <path
            d="M862.485 481.154H234.126l203.3-203.3c12.497-12.497 12.497-32.758 0-45.255s-32.758-12.497-45.255 0L135.397 489.373c-12.497 12.497-12.497 32.758 0 45.254l256.774 256.775c6.249 6.248 14.438 9.372 22.627 9.372s16.379-3.124 22.627-9.372c12.497-12.497 12.497-32.759 0-45.255l-203.3-203.301h628.36c17.036 0 30.846-13.81 30.846-30.846s-13.81-30.846-30.846-30.846z"
            fill=""
            p-id="10101"
          ></path>
        </svg>
      </button>
      <button class="change-button" @click="goForward" :disabled="!canGoForward">
        <svg
          t="1632984737821"
          class="icon"
          viewBox="0 0 1024 1024"
          version="1.1"
          xmlns="http://www.w3.org/2000/svg"
          p-id="10249"
          width="20"
          height="20"
        >
          <path
            d="M885.113 489.373L628.338 232.599c-12.496-12.497-32.758-12.497-45.254 0-12.497 12.497-12.497 32.758 0 45.255l203.3 203.3H158.025c-17.036 0-30.846 13.811-30.846 30.846 0 17.036 13.811 30.846 30.846 30.846h628.36L583.084 746.147c-12.497 12.496-12.497 32.758 0 45.255 6.248 6.248 14.438 9.372 22.627 9.372s16.379-3.124 22.627-9.372l256.775-256.775a31.999 31.999 0 0 0 0-45.254z"
            fill=""
            p-id="10250"
          ></path>
        </svg>
      </button>
      <button class="change-button" @click="refreshPage">
        <svg
          t="1632984867128"
          class="icon"
          viewBox="0 0 1024 1024"
          version="1.1"
          xmlns="http://www.w3.org/2000/svg"
          p-id="1857"
          width="15"
          height="15"
        >
          <path
            d="M927.999436 531.028522a31.998984 31.998984 0 0 0-31.998984 31.998984c0 51.852948-10.147341 102.138098-30.163865 149.461048a385.47252 385.47252 0 0 1-204.377345 204.377345c-47.32295 20.016524-97.6081 30.163865-149.461048 30.163865s-102.138098-10.147341-149.461048-30.163865a385.47252 385.47252 0 0 1-204.377345-204.377345c-20.016524-47.32295-30.163865-97.6081-30.163865-149.461048s10.147341-102.138098 30.163865-149.461048a385.47252 385.47252 0 0 1 204.377345-204.377345c47.32295-20.016524 97.6081-30.163865 149.461048-30.163865a387.379888 387.379888 0 0 1 59.193424 4.533611l-56.538282 22.035878A31.998984 31.998984 0 1 0 537.892156 265.232491l137.041483-53.402685a31.998984 31.998984 0 0 0 18.195855-41.434674L639.723197 33.357261a31.998984 31.998984 0 1 0-59.630529 23.23882l26.695923 68.502679a449.969005 449.969005 0 0 0-94.786785-10.060642c-60.465003 0-119.138236 11.8488-174.390489 35.217667a449.214005 449.214005 0 0 0-238.388457 238.388457c-23.361643 55.252253-35.22128 113.925486-35.22128 174.390489s11.8488 119.138236 35.217668 174.390489a449.214005 449.214005 0 0 0 238.388457 238.388457c55.252253 23.368867 113.925486 35.217667 174.390489 35.217667s119.138236-11.8488 174.390489-35.217667A449.210393 449.210393 0 0 0 924.784365 737.42522c23.368867-55.270316 35.217667-113.925486 35.217667-174.390489a31.998984 31.998984 0 0 0-32.002596-32.006209z"
            fill=""
            p-id="1858"
          ></path>
        </svg>
      </button>

      <input 
        class="url-input" 
        v-model="urlinput" 
        @keydown.enter="changeUrl"
        placeholder="Enter URL"
      />
    </div>

    <iframe class="webframe" ref="iframeRef" :src="urlsrc"></iframe>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, onUnmounted } from 'vue';

const urlinput = ref('https://wap.baidu.com/');
const urlsrc = ref('https://wap.baidu.com/');
const iframeRef:any = ref(null);

let canGoBack = false;
let canGoForward = false;

function updateNavigationButtons() {
  canGoBack = window.history.length > 1;
  canGoForward = window.history.length < window.history.state?.length;
}


function goBack() {
  if (canGoBack) {
    window.history.back();
  }
}

function goForward() {
  if (canGoForward) {
    window.history.forward();
  }
}

function refreshPage() {
  iframeRef.value.contentWindow?.location.reload();
}

function changeUrl() {
  const protocolAdded = /^(http|https):\/\//.test(urlinput.value) ? urlinput.value : `https://${urlinput.value}`;
  urlsrc.value = protocolAdded;
}

onMounted(() => {
  window.addEventListener('popstate', updateNavigationButtons);
  updateNavigationButtons();
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
}
</style>