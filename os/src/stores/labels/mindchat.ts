export const mindchatLabels = {
    name: "mindchat",
    family: "qwen",
    action: "chat",
    models: [
        {
            model: "MindChat-Qwen2:4b",
            type: "llm",
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
                url: ["https://hf-mirror.com/v8karlo/MindChat-Qwen2-4B-Q5_K_M-GGUF/blob/main/mindchat-qwen2-4b-q5_k_m.gguf"],
                engine: "ollama",
                from: "network",
                size: "2.84GB",
                desk: "3GB",
                cpu: "16GB",
                gpu: "8GB",
                quant: "q5",
                "parameters": "stop                           \"<|im_start|>\"\nstop                           \"<|im_end|>\"",
                "context_length": 32768,
                "embedding_length": 896,
                "template": "{{ if .System }}<|im_start|>system\n{{ .System }}<|im_end|>\n{{ end }}{{ if .Prompt }}<|im_start|>user\n{{ .Prompt }}<|im_end|>\n{{ end }}<|im_start|>assistant\n{{ .Response }}<|im_end|>\n",

            }
        },
    ],
    zhdesc: "心理大模型——漫谈(MindChat)期望从心理咨询、心理评估、心理诊断、心理治疗四个维度帮助人们纾解心理压力与解决心理困惑",
    endesc: "MindChat aims to help people relieve psychological stress and solve psychological confusion from four dimensions: psychological counseling, psychological assessment, psychological diagnosis, and psychological therapy"
}