export const bilibiliLabels = {
    name: "bilibili",
    family: "llama",
    action: "chat",
    models: [
        {
            model: "milkey/bilibili-index:latest",
            params: {
                stream: true,
                "repeat_penalty": 1.1,
                "stop": [
                    "reserved_0",
                    "reserved_1",
                    "</s>",
                    "<unk>"
                ],
                "temperature": 0.3,
                "top_k": 5,
                "top_p": 0.8
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "2.3GB",
                desk: "3GB",
                cpu: "8GB",
                gpu: "6GB",
                quant: "q8"
            }
        },
    ],
    zhdesc: "由哔哩哔哩自主研发的大语言模型",
    endesc: "A large language model independently developed by Bilibili"
}