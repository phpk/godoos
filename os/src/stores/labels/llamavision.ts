export const llamaLabels = {
    name: "llama3-vision",
    family: "llama",
    action: "img2txt",
    zhdesc: "Llama 3.2-Vision多模态大型语言模型（LLM）集合是11B和90B大小（文本+图像输入/文本输出）的指令调优图像推理生成模型集合。Llama 3.2-Vision指令调优模型针对视觉识别、图像推理、字幕和回答有关图像的一般问题进行了优化。在常见的行业基准上，这些模型的表现优于许多可用的开源和封闭式多模式模型。",
    endesc: "The Llama 3.2-Vision collection of multimodal large language models (LLMs) is a collection of instruction-tuned image reasoning generative models in 11B and 90B sizes (text + images in / text out). The Llama 3.2-Vision instruction-tuned models are optimized for visual recognition, image reasoning, captioning, and answering general questions about an image. The models outperform many of the available open source and closed multimodal models on common industry benchmarks.",
    models: [
        {
            model: "llama3.2-vision:11b",
            params: {
                "temperature": 0.6,
                "top_p": 0.9            
            },
            info: {
                engine: "ollama",
                from: "ollama",
                size: "6GB",
                desk: "6GB",
                cpu: "32GB",
                gpu: "12GB",
                quant: "q4"
            }
        }
    ],

}