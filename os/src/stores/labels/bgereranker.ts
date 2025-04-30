export const bgeRerankerLabels = {
    name: "bge-reranker",
    family: "bge",
    action: "reranker",
    models: [
        {
            model: "linux6200/bge-reranker-v2-m3",
            params: {
                "num_ctx": 4096,
                "temperature": 1            
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "1.2GB",
                desk: "2GB",
                cpu: "8GB",
                gpu: "6GB",
                quant: "f16"
            }
        },
    ],
    zhdesc: "bge-reranker是BAAI开发的排序模型",
    endesc: "bge is an reranker model developed by BAAI"
}