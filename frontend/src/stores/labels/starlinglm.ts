export const starlinglmLabels = {
    name: "starling-lm",
    family: "llama",
    action: "chat",
    models: [
        {
            model: "starling-lm:latest",
            params: {
                stream: true,
                "stop": [
                    "<|endoftext|>",
                    "<|end_of_turn|>",
                    "Human:",
                    "Assistant:"                       
                ],
                "temperature": 0.7,
                "top_k": 5,
                "top_p": 0.8
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
    zhdesc: "Starling是一个大型语言模型，通过人工智能反馈的强化学习进行训练，专注于提高聊天机器人的有用性。",
    endesc: "Starling is a large language model trained by reinforcement learning from AI feedback focused on improving chatbot helpfulness. "
}