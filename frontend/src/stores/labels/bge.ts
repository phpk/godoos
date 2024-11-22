export const bgeLabels = {
    name: "bge-large-zh-v1.5",
    family: "bge",
    action: "embeddings",
    models: [
        {
            model: "quentinz/bge-large-zh-v1.5:latest",
            params: {
                "num_ctx": 512
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "651MB",
                desk: "1GB",
                cpu: "8GB",
                gpu: "6GB",
                quant: "f16"
            }
        },
        {
            model: "quentinz/bge-base-zh-v1.5:latest",
            params: {
                "num_ctx": 512
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "205MB",
                desk: "1GB",
                cpu: "8GB",
                gpu: "6GB",
                quant: "f16"
            }
        },
    ],
    zhdesc: "bge是BAAI开发的嵌入模型",
    endesc: "bge is an embedded model developed by BAAI"
}