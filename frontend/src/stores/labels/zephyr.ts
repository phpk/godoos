export const zephyrLabels = {
    name: "zephyr",
    family: "llama",
    action: "chat",
    models: [
        {
            model: "zephyr:7b",
            params: {
                top_p: 0.95,
                stream: true,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                "stop": [
                    "<|system|>",
                    "<|user|>",
                    "<|assistant|>",
                    "</s>"
                ]
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "4.1GB",
                desk: "5GB",
                cpu: "16GB",
                gpu: "8GB",
                quant: "q4"
            }
        },
    ],
    zhdesc: "Zephyr是Mistral和Mixtral模型的一系列微调版本，经过训练，可以充当有用的助手。",
    endesc: "Zephyr is a series of fine-tuned versions of the Mistral and Mixtral models that are trained to act as helpful assistants.  "
}