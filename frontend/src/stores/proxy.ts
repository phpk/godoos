import { defineStore } from 'pinia';
import { ref } from 'vue';
import { getSystemKey } from "@/system/config"
export interface ProxyItem {
  id?:number;
  port: number;
  domain: string;
  type: string;
  name: string;
  localPort: number;
  localIp: string;
  serverAddr: string;
  serverPort: number;
  httpAuth: boolean;
  authUsername: string;
  authPassword: string;
  https2http: boolean;
  https2httpCaFile: string;
  https2httpKeyFile: string;
  remotePort?: number;
  stcpModel?: string;
  secretKey?: string;
  visitedName?: string;
  bindAddr?: string;
  bindPort?: number;
  fallbackTo?: string;
  fallbackTimeoutMs?: number;
  keepAlive?: boolean;
  customDomains?: string;
  staticFile?:boolean;
  localPath?:string;
  stripPrefix?: string;

}

export const useProxyStore = defineStore('proxyStore', () => {

  const proxies = ref<ProxyItem[]>([]);
  const customDomains = ref<string[]>([]);
  const proxyData = ref<ProxyItem>(createNewProxyData());
  const apiUrl = getSystemKey("apiUrl")
  const isEditor = ref(false)
  const addShow = ref(false)
  const settingShow = ref(false)
  const status = ref(false)
  const page = ref({
    current: 1,
    size: 10,
    total: 0,
    pages: 0,
    visible: 5
  })
  function createNewProxyData(): ProxyItem {
    return {
      id:0,
      type: "http",
      name: "",
      port: 8000,
      domain: "",
      localPort: 56780,
      localIp: "127.0.0.1",
      serverAddr: "",
      serverPort: 0,
      httpAuth: false,
      authUsername: "",
      authPassword: "",
      https2http: false,
      https2httpCaFile: "",
      https2httpKeyFile: "",
      remotePort: 0,
      stcpModel: "",
      secretKey: "",
      visitedName: "",
      bindAddr: "",
      bindPort: 0,
      fallbackTo: "",
      fallbackTimeoutMs: 500,
      keepAlive: false,
      customDomains: "",
      staticFile:false,
      localPath:"",
      stripPrefix:""
    };
  }


  const addCustomDomain = () => {
    customDomains.value.push("");
  };

  const removeCustomDomain = (index: number) => {
    if (customDomains.value.length > 1) {
      customDomains.value.splice(index, 1);
    }
  };
  const resetProxyData = () => {
    proxyData.value = createNewProxyData();
  };

  const createFrpc = async () => {
    const url = `${apiUrl}/frpc/create`;
    const postData:any = toRaw(proxyData.value)
    if(customDomains.value.length > 0){
      postData.customDomains = customDomains.value.join(',');
    }
    postData.LocalPort = parseInt(postData.localPort)
    postData.remotePort = parseInt(postData.remotePort)
    postData.bindPort = parseInt(postData.bindPort)
    postData.FallbackTimeoutMs = parseInt(postData.FallbackTimeoutMs)
    return await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify(postData)
    }).then(response => response.json())
  };
  async function pageChange(val: number) {
    page.value.current = val
    await fetchProxies()
   }
  const fetchProxies = async () => {
    const url = `${apiUrl}/frpc/list?page=${page.value.current}&limit=${page.value.size}`;
    const response = await fetch(url);

    if (!response.ok) {
      return
    }
    const res = await response.json();
    if (res.code == 0) {
      proxies.value = res.data.list;
      page.value.total = res.data.total;
    } else {
      console.error("Failed to retrieve proxies:", res.message);
    }
    fetch(`${apiUrl}/frpc/status`).then(res => res.json()).then(res => {
      status.value = res.data
    })
  };

  const fetchProxy = async (id: number) => {
    fetch(`${apiUrl}/frpc/get?id=${id}`).then(res => res.json()).then(res => {
      if (res.code == 0) {
        proxyData.value = res.data;
        customDomains.value = proxyData.value.customDomains?.split(',') || [];
      }
    })
  };

  const deleteProxyById = async (id: number) => {
    const url = `${apiUrl}/frpc/delete?id=${id}`;
    const response = await fetch(url);

    if (!response.ok) {
      return false;
    }

    const data = await response.json();
    if (data.code === 0) {
      await fetchProxies();
      return true;
    } else {
      return false;
    }
  };

  const updateProxy = async (proxy: ProxyItem) => {
    const url = `${apiUrl}/frpc/update`;
    if(customDomains.value.length > 0){
      proxy.customDomains = customDomains.value.join(',');
    }
    fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify(proxy)
    }).then(res => res.json()).then(res => {
      if (res.code == 0) {
        fetchProxies();
      }
    })
  };
  const getConfig = async () => {
    return await fetch(`${apiUrl}/frpc/getconfig`).then(res => res.json())
  }
  const setConfig = async (config: any) => {
    return await fetch(`${apiUrl}/frpc/setconfig`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify(config)
    }).then(res => res.json())
  }
  const startFrpc = async () => {
    return await fetch(`${apiUrl}/frpc/start`).then(res => res.json()).then(res => {
      if (res.code == 0) {
        status.value = true
      }
    })
  }
  const stopFrpc = async () => {
    return await fetch(`${apiUrl}/frpc/stop`).then(res => res.json()).then(res => {
      if (res.code == 0) {
        status.value = false
      }
    })
  }
  const restartFrpc = async () => {
    return await fetch(`${apiUrl}/frpc/restart`).then(res => res.json()).then(res => {
      if (res.code == 0) {
        status.value = true
      }
    })
  }
  return {
    proxies,
    customDomains,
    proxyData,
    page,
    isEditor,
    addShow,
    settingShow,
    status,
    pageChange,
    addCustomDomain,
    removeCustomDomain,
    resetProxyData,
    createFrpc,
    fetchProxies,
    fetchProxy,
    deleteProxyById,
    updateProxy,
    getConfig,
    setConfig,
    startFrpc,
    stopFrpc,
    restartFrpc
  };
});