export const mxbaiLabels = {
    name: "mxbai-embed-large",
    family: "bert",
    action: "embeddings",
    models: [
        {
            model: "mxbai-embed-large:latest",
            params: {
                "num_ctx": 512
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "670MB",
                desk: "700MB",
                cpu: "8GB",
                gpu: "6GB",
                quant: "f16"
            }
        },
    ],
    zhdesc: "mixedbread.ai的最先进的大型嵌入模型",
    endesc: " State-of-the-art large embedding model from mixedbread.ai"
}