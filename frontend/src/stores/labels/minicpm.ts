export const minicpmLabels = {
    name: "minicpm",
    family: "llama",
    action: "img2txt",
    models: [
        {
            model: "scomper/minicpm-v2.5:latest",
            params: {
                top_p: 0.95,
                stream: true,
                num_predict: 1,
                top_k: 40,
                temperature: 0.7,
                "num_ctx": 2048,
                "num_keep": 4,
                "stop": [
                    "<|start_header_id|>",
                    "<|end_header_id|>",
                    "<|eot_id|>"
                ]
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "8.5GB",
                desk: "9GB",
                cpu: "16GB",
                gpu: "8GB",
                quant: "q8"
            }
        },
    ],
    zhdesc: "MiniCPM-V是面向图文理解的端侧多模态大模型系列，该系列模型接受图像和文本输入，并提供高质量的文本输出。",
    endesc: "MiniCPM-V is an end-to-end multimodal large model series for text and image understanding. This series of models accepts image and text inputs and provides high-quality text output."
}