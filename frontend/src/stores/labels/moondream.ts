export const moondreamLabels = {
    name: "moondream",
    family: "moondream",
    action: "img2txt",
    models: [
        {
            model: "moondream:latest",
            params: {
                top_p: 0.95,
                stream: true,
                num_predict: 1,
                top_k: 40,
                temperature: 0,
                "stop": [
                    "<|endoftext|>",
                    "Question:"
                ],
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "1.7GB",
                desk: "2GB",
                cpu: "8GB",
                gpu: "6GB",
                quant: "q8"
            }
        },
    ],
    zhdesc: "moonvdream2是一个小型视觉语言模型，设计用于在边缘设备上高效运行。",
    endesc: "moondream2 is a small vision language model designed to run efficiently on edge devices. "
}