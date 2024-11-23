export const mixtralLabels = {
    name: "mixtral",
    family: "llama",
    action: "chat",
    models: [
        {
            model: "mixtral:8x7b",
            params: {
                top_p: 0.95,
                stream: true,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                "stop": [
                    "[INST]",
                    "[/INST]"
                ]
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "26GB",
                desk: "27GB",
                cpu: "32GB",
                gpu: "12GB",
                quant: "q4"
            }
        },
    ],
    zhdesc: "Mistral AI在8x7b和8x22b参数大小下的一组具有开放权重的专家混合（MoE）模型。",
    endesc: "A set of Mixture of Experts (MoE) model with open weights by Mistral AI in 8x7b and 8x22b parameter sizes. "
}