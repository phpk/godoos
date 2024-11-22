export const openchatLabels = {
    name: "openchat",
    family: "llama",
    action: "chat",
    models: [
        {
            model: "openchat:7b",
            params: {
                top_p: 0.95,
                stream: true,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                "stop": [
                    "<|endoftext|>",
                    "<|end_of_turn|>"
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
    zhdesc: "OpenChat是一组开源语言模型，使用C-RLFT进行了微调：这是一种受离线强化学习启发的策略。",
    endesc: "OpenChat is set of open-source language models, fine-tuned with C-RLFT: a strategy inspired by offline reinforcement learning."
}