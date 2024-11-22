export const bakllavaLabels = {
    name: "bakllava",
    family: "llama",
    action: "img2txt",
    models: [
        {
            model: "bakllava:7b",
            params: {
                top_p: 0.95,
                stream: true,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                "num_ctx": 4096,
                "stop": [
                    "</s>",
                    "USER:"
                ]
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "4.7GB",
                desk: "5GB",
                cpu: "16GB",
                gpu: "8GB",
                quant: "q4"
            }
        },
    ],
    zhdesc: "BakLLaVA是一个多模式模型，由Mistral 7B基础模型和LLaVA架构组成。",
    endesc: "BakLLaVA is a multimodal model consisting of the Mistral 7B base model augmented with the LLaVA architecture. "
}