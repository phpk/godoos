export const mistralLabels = {
    name: "mistral",
    family: "llama",
    action: "chat",
    models: [
        {
            model: "mistral:7b",
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
                size: "4.1GB",
                desk: "5GB",
                cpu: "16GB",
                gpu: "8GB",
                quant: "q4"
            }
        },
    ],
    zhdesc: "Mistral是一个7.3B参数模型，使用Apache许可证进行分发。它有指令（指令如下）和文本完成两种形式。",
    endesc: "Mistral is a 7.3B parameter model, distributed with the Apache license. It is available in both instruct (instruction following) and text completion. "
}