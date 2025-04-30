export const dmetaLabels = {
    name: "dmeta-embedding-zh",
    family: "dmeta",
    action: "embeddings",
    models: [
        {
            model: "herald/dmeta-embedding-zh:latest",
            params: {
                "num_ctx": 1024
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "205MB",
                desk: "300MB",
                cpu: "8GB",
                gpu: "6GB",
                quant: "f16"
            }
        },
    ],
    zhdesc: "Dmeta-embedding 是一款跨领域、跨任务、开箱即用的中文 Embedding 模型，适用于搜索、问答、智能客服、LLM+RAG 等各种业务场景",
    endesc: "Dmeta-embedding is a cross domain, cross task, and out of the box Chinese embedding model suitable for various business scenarios such as search, Q&A, intelligent customer service, LLM+RAG, etc"
}