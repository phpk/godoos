export const starcoder2Labels = {
    name: "starcoder2",
    family: "starcoder2",
    action: "code",
    models: [
        {
            model: "starcoder2:3b",
            params: {
                top_p: 0.95,
                stream: true,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                "stop": []
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "1.7GB",
                desk: "2GB",
                cpu: "8GB",
                gpu: "6GB",
                quant: "q4"
            }
        },
        {
            model: "starcoder2:7b",
            params: {
                top_p: 0.95,
                stream: true,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                "stop": []
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "4GB",
                desk: "5GB",
                cpu: "16GB",
                gpu: "8GB",
                quant: "q4"
            }
        },
    ],
    zhdesc: "StarCoder2是下一代经过透明训练的开放代码LLM，有三种大小：3B、7B和15B参数。",
    endesc: "StarCoder2 is the next generation of transparently trained open code LLMs that comes in three sizes: 3B, 7B and 15B parameters. "
}