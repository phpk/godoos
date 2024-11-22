export const nomicLabels = {
    name: "nomic-embed-text",
    family: "nomic-bert",
    action: "embeddings",
    models: [
        {
            model: "nomic-embed-text:latest",
            params: {
                "num_ctx": 768
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "274MB",
                desk: "300MB",
                cpu: "8GB",
                gpu: "6GB",
                quant: "f16"
            }
        },
    ],
    zhdesc: "一个具有大型令牌上下文窗口的高性能开放嵌入模型。",
    endesc: "A high-performing open embedding model with a large token context window. "
}