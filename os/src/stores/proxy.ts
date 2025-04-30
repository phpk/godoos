import { defineStore } from "pinia";
import { ref } from "vue";

export const useProxyStore = defineStore("proxy", () => {
  const addShow = ref(false);
  const isEditor = ref(false);

  // 定义代理数据的类型
  type ProxyData = {
    name: string;
    type: string;
    localIp: string;
    localPort: number | null;
    customDomains: string[];
    subdomain: string;
    basicAuth: boolean;
    httpUser: string;
    httpPassword: string;
    https2http: boolean;
    https2httpCaFile: string;
    https2httpKeyFile: string;
    remotePort: number | null;
    visitedName: string;
    bindAddr: string;
    bindPort: number | null;
    stcpModel: string;
    secretKey: string;
    fallbackTo: string;
    fallbackTimeoutMs: number | null;
    keepAlive: boolean;
  };

  const proxies = ref<ProxyData[]>([]); // 为 proxies 指定类型

  const proxyData = ref<ProxyData>({
    name: "",
    type: "",
    localIp: "",
    localPort: null,
    customDomains: [],
    subdomain: "",
    basicAuth: false,
    httpUser: "",
    httpPassword: "",
    https2http: false,
    https2httpCaFile: "",
    https2httpKeyFile: "",
    remotePort: null,
    visitedName: "",
    bindAddr: "",
    bindPort: null,
    stcpModel: "",
    secretKey: "",
    fallbackTo: "",
    fallbackTimeoutMs: null,
    keepAlive: false,
  });

  const customDomains = ref([]);

  const addCustomDomain = () => {
    customDomains.value.push("");
  };

  const removeCustomDomain = (index: number) => {
    customDomains.value.splice(index, 1);
  };

  const resetProxyData = () => {
    proxyData.value = {
      name: "",
      type: "http",
      localIp: "",
      localPort: null,
      customDomains: [],
      subdomain: "",
      basicAuth: false,
      httpUser: "",
      httpPassword: "",
      https2http: false,
      https2httpCaFile: "",
      https2httpKeyFile: "",
      remotePort: null,
      visitedName: "",
      bindAddr: "",
      bindPort: null,
      stcpModel: "",
      secretKey: "",
      fallbackTo: "",
      fallbackTimeoutMs: null,
      keepAlive: false,
    };
  };

  const createFrpc = async () => {
    // 模拟创建代理配置的API调用
    try {
      // 假设API返回一个对象，包含code和message
      const response = { code: 0, message: "创建成功" };
      return response;
    } catch (error) {
      return { code: 1, message: "创建失败" };
    }
  };

  const updateProxy = async (data: any) => {

    try {

      const response = { code: 0, message: "更新成功" };
      return response;
    } catch (error) {
      return { code: 1, message: "更新失败" };
    }
  };

  return {
    proxyData,
    addShow,
    resetProxyData,
    isEditor,
    proxies,
    customDomains,
    addCustomDomain,
    removeCustomDomain,
    createFrpc,
    updateProxy,
  };
});
