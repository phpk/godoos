import { defineStore } from 'pinia';
import { ref } from 'vue';

export interface ProxyItem {
  port: string;
  domain: string;
  type: string;
  name: string;
  serverAddr: string;
  serverPort: string;
  httpAuth: boolean;
  authUsername: string;
  authPassword: string;
  https2http: boolean;
  https2httpCaFile: string;
  https2httpKeyFile: string;
  remotePort?: string;
  stcpModel?: string;
  secretKey?: string;
  visitedName?: string;
  bindAddr?: string;
  bindPort?: string;
  fallbackTo?: string;
  fallbackTimeoutMs?: number;
  keepAlive?: boolean;
}

export const useProxyStore = defineStore('proxyStore', () => {
  const localKey = "godoos_net_proxy";

  const proxies = ref<ProxyItem[]>(getProxies());
  const customDomains = ref<string[]>([""]);
  const proxyData = ref<ProxyItem>(createNewProxyData());

  function createNewProxyData(): ProxyItem {
    return {
      type: "http",
      name: "",
      port: "",
      domain: "",
      serverAddr: "",
      serverPort: "",
      httpAuth: false,
      authUsername: "",
      authPassword: "",
      https2http: false,
      https2httpCaFile: "",
      https2httpKeyFile: "",
      remotePort: "",
      stcpModel: "",
      secretKey: "",
      visitedName: "",
      bindAddr: "",
      bindPort: "",
      fallbackTo: "",
      fallbackTimeoutMs: 0,
      keepAlive: false,
    };
  }

  function getProxies(): ProxyItem[] {
    const proxies = localStorage.getItem(localKey);
    return proxies ? JSON.parse(proxies) : [];
  }

  function saveProxies(proxies: ProxyItem[]) {
    localStorage.setItem(localKey, JSON.stringify(proxies));
  }

  const addProxy = (proxy: ProxyItem) => {
    proxies.value.push(proxy);
    saveProxies(proxies.value);
  };

  const updateProxy = (updatedProxy: ProxyItem) => {
    const index = proxies.value.findIndex(p => p.id === updatedProxy.id);
    if (index !== -1) {
      proxies.value[index] = updatedProxy;
      saveProxies(proxies.value);
    }
  };

  const deleteProxy = (id: number) => {
    proxies.value = proxies.value.filter(p => p.id !== id);
    saveProxies(proxies.value);
  };

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

  const createFrpcConfig = async () => {
    const url = "http://localhost:56780/frpc/create";
    const requestData = {
      config: {
        ServerAddr: proxyData.value.serverAddr,
        ServerPort: Number(proxyData.value.serverPort),
        AuthMethod: "token",
        AuthToken: "your-auth-token",
        User: "your-username",
        MetaToken: "your-meta-token",
        TransportHeartbeatInterval: 30,
        TransportHeartbeatTimeout: 90,
        LogLevel: "info",
        LogMaxDays: 3,
        WebPort: 7500,
        TlsConfigEnable: true,
        TlsConfigCertFile: "/path/to/cert/file",
        TlsConfigKeyFile: "/path/to/key/file",
        TlsConfigTrustedCaFile: "/path/to/ca/file",
        TlsConfigServerName: "server-name",
        ProxyConfigEnable: true,
        ProxyConfigProxyUrl: "http://proxy.example.com"
      },
      proxies: [
        {
          Name: proxyData.value.name,
          Type: proxyData.value.type,
          LocalIp: proxyData.value.serverAddr,
          LocalPort: Number(proxyData.value.port),
          RemotePort: Number(proxyData.value.remotePort),
          CustomDomains: customDomains.value,
          Subdomain: proxyData.value.domain,
          BasicAuth: proxyData.value.httpAuth,
          HttpUser: proxyData.value.authUsername,
          HttpPassword: proxyData.value.authPassword,
          StcpModel: proxyData.value.stcpModel,
          ServerName: "server1",
          BindAddr: proxyData.value.bindAddr,
          BindPort: Number(proxyData.value.bindPort),
          FallbackTo: proxyData.value.fallbackTo,
          FallbackTimeoutMs: Number(proxyData.value.fallbackTimeoutMs),
          SecretKey: proxyData.value.secretKey
        }
      ]
    };

    try {
      const response = await fetch(url, {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify(requestData)
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      console.log("Response data:", data);
      // 可以在这里添加成功通知
    } catch (error) {
      console.error("Error:", error);
      // 可以在这里添加错误通知
    }
  };

  return {
    proxies,
    customDomains,
    proxyData,
    addProxy,
    updateProxy,
    deleteProxy,
    addCustomDomain,
    removeCustomDomain,
    resetProxyData,
    createFrpcConfig,
  };
});