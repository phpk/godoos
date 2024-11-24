import { defineStore } from "pinia";
import { ref } from "vue";
import { db } from "./db.ts"
import { aiLabels } from "./labels/index.ts"
import { fetchGet, getSystemKey } from "@/system/config"
const modelEngines = [
  {
    name: "ollama",
    cpp: "llama.cpp",
    needQuant: true
  },
  {
    name: "sd",
    cpp: "stable-diffusion.cpp",
    needQuant: false
  },
  {
    name: "voice",
    cpp: "sherpa.cpp",
    needQuant: false
  }
]
const llamaQuant = [
  "q2_K",
  "q3_K",
  "q3_K_S",
  "q3_K_M",
  "q3_K_L",
  "q4_0",
  "q4_1",
  "q4_K",
  "q4_K_S",
  "q4_K_M",
  "q5_0",
  "q5_1",
  "q5_K",
  "q5_K_S",
  "q5_K_M",
  "q6_K",
  "q8_0",
  "f16",
]
export const useModelStore = defineStore('modelStore', () => {

  const labelList: any = ref([])
  const cateList: any = ["chat", "translation", "code", "img2txt", "image", "tts", "audio", "embeddings"]
  const modelList: any = ref([])
  const downList: any = ref([])
  const chatConfig: any = ref({
    chat: {
      key: "chat",
      contextLength: 10,
      num_keep: 5, //保留多少个最有可能的预测结果。这与top_k一起使用，决定模型在生成下一个词时考虑的词汇范围。
      num_predict: 3, //生成多少个预测结果
      top_p: 0.95,
      top_k: 40, //影响生成的随机性。较高的top_k值将使模型考虑更多的词汇
      temperature: 0.7, //影响生成的随机性。较低的温度产生更保守的输出，较高的温度产生更随机的输出。
    },
    translation: {
      key: "translation",
      num_keep: 5,
      num_predict: 1,
      top_k: 40,
      top_p: 0.95,
      temperature: 0.2,
    },
    creation: {
      key: "creation",
      num_keep: 3,
      num_predict: 1,
      top_k: 40,
      top_p: 0.95,
      temperature: 0.2,
    },
    knowledge: {
      key: "knowledge",
      contextLength: 10,
      num_keep: 5,
      num_predict: 1,
      top_k: 40,
      top_p: 0.95,
      temperature: 0.2,
    },
    spoken: {
      key: "spoken",
      contextLength: 10,
      num_keep: 5,
      num_predict: 1,
      top_k: 40,
      top_p: 0.95,
      temperature: 0.2,
    }
  })
  const aiUrl = getSystemKey("aiUrl")

  async function getLabelCate(cateName: string) {
    const list = await getLabelList()
    labelList.value = list.filter((d: any) => {
      if (cateName == 'all') {
        return true
      } else {
        return d.action == cateName
      }
    })
  }

  async function getLabelSearch(keyword: string) {
    const list = await getLabelList()
    if (!keyword || keyword == "") {
      labelList.value = list
    }
    labelList.value = list.filter((d: any) => d.name.toLowerCase().includes(keyword.toLowerCase()))
  }
  async function getLabelList() {
    return await db.getAll("modelslabel")
    //return await db.getByField("modelslabel", "chanel", getSystemKey("currentChanel"))
  }
  async function delLabel(id: number) {
    await db.delete("modelslabel", id)
    labelList.value = await getLabelList()
  }
  async function checkLabelData(data: any) {
    const labelData = await db.get("modelslabel", { name: data.label })
    if (!labelData) {
      return
    }
    if (labelData.models.find((d: any) => d.model == data.model)) {
      return
    }
    labelData.models.push(data)

    await db.update("modelslabel", labelData.id, labelData)

  }

  async function getModelList() {
    const res = await fetchGet(`${aiUrl}/ai/tags`)
    //console.log(res)
    if (res.ok) {
      await resetData(res)
    }
    return modelList.value
  }
  async function resetData(res: any) {
    const data = await res.json();
    // console.log(data);
    if (data && data.length > 0) {
        // 获取当前modelList中的模型名称
        const existingModels:any = [];
        const has = await db.getAll("modelslist");
        has.forEach((model: any) => {
          if(model.isdef && model.isdef > 0) {
            existingModels.push(model.model)
          }
        })
        data.forEach((d:any) => {
          if (existingModels.includes(d.model)) {
            d.isdef = 1
          }
        });
        await db.clear("modelslist");
        await db.addAll("modelslist", data);
        modelList.value = data;
    }
    // 重新获取所有模型列表
    
}
  async function refreshOllama() {
    const res = await fetchGet(`${aiUrl}/ai/refreshOllama`)
    //console.log(res)
    if (res.ok) {
      resetData(res)
    }
  }
  function getModelInfo(model: string) {
    return modelList.value.find((d: any) => d.model == model)
  }
  async function getModel(action: string) {
    const model = await db.get("modelslist", { action, isdef: 1 })
    if (!model) {
      return await db.addOne("modelslist", { action })
    } else {
      return model
    }
  }
  async function getList() {
    labelList.value = await getLabelList()
    await getModelList()
    downList.value.forEach((_: any, index: number) => {
      downList.value[index].isLoading = 0
    })
  }
  async function setCurrentModel(action: string, model?: string) {
    await db.modify("modelslist", "action", action, { isdef: 0 })
    //console.log(model)
    if (model !== "") {
      const data = await db.get("modelslist", { model })
      if (data) {
        return await db.update("modelslist", data.id, { isdef: 1 })
      }
    } else {
      const data = await db.get("modelslist", { action })
      if (data) {
        return await db.update("modelslist", data.id, { isdef: 1 })
      }
    }
  }
  async function setDefModel(action: string) {
    const has = await db.get("modelslist", { action, isdef: 1 })
    if (!has) {
      const data = await db.get("modelslist", { action })
      if (data) {
        return await db.update("modelslist", data.id, { isdef: 1 })
      }
    }
  }
  function getCurrentModelList(action: string) {
    //if (!modelList || modelList.length == 0) return
    return modelList.value.filter((d: any) => d.action == action)
  }


  async function deleteModelList(data: any) {
    //console.log(data)
    if (!data || !data.model) return
    const postData = {
      method: "POST",
      body: JSON.stringify(data),
    };
    const delUrl = aiUrl + "/ai/delete";
    const completion = await fetch(delUrl, postData);
    if (completion.status === 404) {
      return completion.statusText;
    }
    if (completion.status === 200) {
      modelList.value.forEach((d: any, index: number) => {
        if (d.model == data.model) {
          modelList.value.splice(index, 1);
        }
      });
      await db.deleteByField("modelslist", "model", data.model)
      if (data.isdef * 1 == 1) {
        await setCurrentModel(data.action, "")
      }
    }


    //await db.delete("modelslist", data.id)
    //await getModelList()
  }

  function checkDownload(name: string) {
    return modelList.value.find((d: any) => d.model === name);
  }
  function addDownload(data: any) {
    const has = downList.value.find((d: any) => d.model === data.model)
    if (!has) {
      downList.value.unshift(data)
    } else {
      updateDownload(data)
    }

    return data
  }
  function deleteDownload(model: string) {
    //console.log(model)
    downList.value.forEach((d: any, index: number) => {
      if (d.model == model) {
        downList.value.splice(index, 1);
      }
    });
  }
  async function updateDownload(modelData: any) {
    const index = downList.value.findIndex((d: any) => d.model === modelData.model);
    if (index !== -1) {
      // 或者使用splice方法替换对象
      downList.value.splice(index, 1, {
        ...downList.value[index],
        status: modelData.status,
        progress: modelData.progress,
        isLoading: modelData.isLoading ?? 0,
      });
      if (modelData.status === "success") {
        //await addDownList(modelData);
        await getModelList();
        await setDefModel(modelData.action);
        await checkLabelData(modelData);
      }
    }
  }
  function parseJson(str: string): any {
    try {
      return JSON.parse(str);
    } catch (e) {
      return undefined;
    }
  }
  function parseMsg(str: string) {
    const nres = { status: "" }
    try {
      //console.log(str)
      if (str == 'has done!') {
        return { status: 'success' }
      }
      const raw: any = str.split("\n")
      if (raw.length < 1) return nres
      // deno-lint-ignore no-explicit-any
      const rt: string[] = raw.filter((d: string) => d.trim() !== "");
      //console.log(rt)
      if (rt.length > 0) {
        let res: any[] = [];
        rt.forEach((d: string) => {
          const msg = parseJson(d);
          if (msg) {
            res.push(msg);
          }
        });
        if (res.length > 0) {
          return res[res.length - 1]
        } else {
          return nres
        }
      } else {
        return nres;
      }
    } catch (error) {
      console.log(error);
      return nres
    }
  }
  async function initModel() {
    await db.clear("modelslabel")
    await db.addAll("modelslabel", aiLabels);
  }


  return {
    cateList,
    labelList,
    modelList,
    downList,
    modelEngines,
    llamaQuant,
    chatConfig,
    getList,
    getModelList,
    getModelInfo,
    getModel,
    checkDownload,
    addDownload,
    deleteDownload,
    updateDownload,
    checkLabelData,
    getLabelCate,
    getLabelSearch,
    getLabelList,
    delLabel,
    //addDownList,
    deleteModelList,
    initModel,
    setCurrentModel,
    getCurrentModelList,
    parseMsg,
    refreshOllama
  }

}, {
  persist: {
    enabled: true,
    strategies: [
      {
        storage: localStorage,
        paths: [
          "downList",
          "chatConfig"
        ]
      }, // name 字段用localstorage存储
    ],
  }
})
