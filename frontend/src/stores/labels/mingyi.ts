export const mingyiLabels = {
    name: "ming",
    family: "qwen",
    action: "chat",
    models: [
        {
            model: "ming:1.8B",
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
                from: "network",
                url: ["https://hf-mirror.com/capricornstone/MING-1.8B-Q8_0-GGUF/blob/main/ming-1.8b-q8_0.gguf"],
                size: "1.96GB",
                desk: "2GB",
                cpu: "8GB",
                gpu: "6GB",
                quant: "q8",
                "parameters": "stop                           \"<|im_start|>\"\nstop                           \"<|im_end|>\"",
                "context_length": 32768,
                "embedding_length": 1024,
                "template": "{{ if .System }}<|im_start|>system\n{{ .System }}<|im_end|>{{ end }}<|im_start|>user\n{{ .Prompt }}<|im_end|>\n<|im_start|>assistant\n",

            }
        },
    ],
    zhdesc: "明医 (MING)：中文医疗问诊大模型",
    endesc: "MING: A Chinese Medical Consultation Model"
}