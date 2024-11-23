export const h2oLabels = {
    name: "h2o",
    family: "llama",

    action: "chat",
    models: [
        {
            model: "cas/h2o-danube2-1.8b-chat:latest",
            params: {
                top_p: 0.95,
                stream: true,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                "stop": [
                    "<|prompt|>",
                    "</s>",
                    "<|answer|>"
                ]
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "1.1GB",
                desk: "2GB",
                cpu: "8GB",
                gpu: "6GB",
                quant: "q4"
            }
        },
    ],
    zhdesc: "H2O.ai在Apache v2.0下发布了其最新的开放权重小语言模型H2O-Danube2-1.8B。",
    endesc: "H2O.ai just released its latest open-weight small language model, H2O-Danube2-1.8B, under Apache v2.0."
}