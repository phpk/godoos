export const phiLabels = {
    name: "phi",
    family: "phi",

    action: "chat",
    models: [
        {
            model: "phi3:mini",
            params: {
                top_p: 0.95,
                stream: true,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                "stop": [
                    "<|end|>",
                    "<|user|>",
                    "<|assistant|>"
                ]
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "2.4GB",
                desk: "3GB",
                cpu: "8GB",
                gpu: "6GB",
                quant: "q4"
            }
        },
        {
            model: "phi3:medium",
            params: {
                top_p: 0.95,
                stream: true,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                "stop": [
                    "<|end|>",
                    "<|user|>",
                    "<|assistant|>"
                ]
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "7.9GB",
                desk: "8GB",
                cpu: "16GB",
                gpu: "8GB",
                quant: "q4"
            }
        },
    ],
    zhdesc: "Phi是微软开发的一系列开放式人工智能模型。",
    endesc: "Phi is a family of open AI models developed by Microsoft."
}