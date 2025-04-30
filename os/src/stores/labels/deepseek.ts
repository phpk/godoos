export const deepseekLabels = {
    name: "deepseek",
    family: "deepseek",
    action: "chat",
    models: [
        {
            model: "deepseek-r1:1.5b",
            params: {
                "stop": [
                    "<｜begin▁of▁sentence｜>",
                    "<｜end▁of▁sentence｜>",
                    "<｜User｜>",
                    "<｜Assistant｜>"
                ]
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "1.1GB",
                desk: "1.5GB",
                cpu: "8GB",
                gpu: "6GB",
                quant: "q4"
            }
        },
        {
            model: "deepseek-r1:7b",
            params: {
                "stop": [
                    "<｜begin▁of▁sentence｜>",
                    "<｜end▁of▁sentence｜>",
                    "<｜User｜>",
                    "<｜Assistant｜>"
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
    zhdesc: "DeepSeek-R1，是幻方量化旗下AI公司深度求索（DeepSeek）研发的推理模型。DeepSeek-R1采用强化学习进行后训练，旨在提升推理能力，尤其擅长数学、代码和自然语言推理等复杂任务。",
    endesc: "DeepSeek-R1， It is an inference model developed by DeepSeek, an AI company under the umbrella of Magic Square Quantification. DeepSeek-R1 uses reinforcement learning for post training, aiming to improve reasoning ability, especially in complex tasks such as mathematics, code, and natural language reasoning."
}