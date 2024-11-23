export const gemmaLabels = {
    name: "gemma",
    family: "gemma",
    action: "chat",
    zhdesc: "Gemma是由谷歌及其DeepMind团队开发的一个新的开放模型。",
    endesc: "Gemma is a new open model developed by Google and its DeepMind team.",
    models: [
        {
            model: "gemma2:9b",
            params: {
                top_p: 0.95,
                stream: true,
                num_keep: 5,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                "penalize_newline": false,
                "repeat_penalty": 1,
                "stop": [
                    "<start_of_turn>",
                    "<end_of_turn>"
                ]
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "5.5GB",
                desk: "6GB",
                cpu: "16GB",
                gpu: "8GB",
                quant: "q4"
            }
        },
        {
            model: "gemma:2b",
            params: {
                top_p: 0.95,
                stream: true,
                num_keep: 5,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                "penalize_newline": false,
                "repeat_penalty": 1,
                "stop": [
                    "<start_of_turn>",
                    "<end_of_turn>"
                ]
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
            model: "gemma:7b",
            params: {
                top_p: 0.95,
                stream: true,
                num_keep: 5,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                "penalize_newline": false,
                "repeat_penalty": 1,
                "stop": [
                    "<start_of_turn>",
                    "<end_of_turn>"
                ]
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "5.0GB",
                desk: "6GB",
                cpu: "16GB",
                gpu: "8GB",
                quant: "q4"
            }
        },
    ],
    
}