export const ayaLabels = {
    name: "aya",
    family: "llama",
    action: "translation",
    models: [
        {
            model: "aya-expanse:8b",
            params: {
                top_p: 0.95,
                stream: true,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                "stop": [
                    "<|START_OF_TURN_TOKEN|>",
                    "<|END_OF_TURN_TOKEN|>"
                ]
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "8.1GB",
                desk: "9GB",
                cpu: "16GB",
                gpu: "8GB",
                quant: "q4"
            }
        },
        {
            model: "aya:8b",
            params: {
                top_p: 0.95,
                stream: true,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                "stop": [
                    "<|START_OF_TURN_TOKEN|>",
                    "<|END_OF_TURN_TOKEN|>"
                ]
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "4.8GB",
                desk: "5GB",
                cpu: "16GB",
                gpu: "8GB",
                quant: "q4"
            }
        },
    ],
    zhdesc: "Aya 23可以流利地说23种语言。",
    endesc: "Aya 23 can talk upto 23 languages fluently."
}