export const snowflakeLabels = {
    name: "snowflake-arctic-embed",
    family: "bert",
    action: "embeddings",
    models: [
        {
            model: "snowflake-arctic-embed:latest",
            params: {
                "num_ctx": 1024
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "669MB",
                desk: "700MB",
                cpu: "8GB",
                gpu: "6GB",
                quant: "f16"
            }
        },
    ],
    zhdesc: "snowflake-arctic-embed是一套文本嵌入模型，专注于创建针对性能优化的高质量检索模型。",
    endesc: "snowflake-arctic-embed is a suite of text embedding models that focuses on creating high-quality retrieval models optimized for performance. "
}