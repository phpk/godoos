export const internlmLabels = {
    name: "internlm",
    family: "internlm",
    action: "chat",
    models: [
        {
            model: "internlm2:7b",
            params: {
                top_p: 0.95,
                stream: true,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                "stop": [
                    "<|im_start|>",
                    "<|im_end|>"
                ]
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "4.5GB",
                desk: "5GB",
                cpu: "16GB",
                gpu: "8GB",
                quant: "q4"
            }
        },
    ],
    zhdesc: "上海人工智能实验室与商汤科技联合香港中文大学、复旦大学发布的新一代大语言模型书生·浦语",
    endesc: "The new generation of big language model internlm, jointly released by Shanghai Artificial Intelligence Laboratory and Shangtang Technology, the Chinese University of Hong Kong and Fudan University"
}