import { sdLabel } from './stable-diffusion.ts'
import { whisperLabel } from './whisper.ts'
import { nemoLabel } from './nemo.ts'
import { zipformerLabel } from './zipformer.ts'
import { paraformerLabel } from './paraformer.ts'
import { telespeechLabel } from './telespeech.ts'
import { vitsLabel } from './vits.ts'
export const aiLabels = [
    {
        name: "qwen",
        family: "llama",
        engine: "ollama",
        from:"ollama",
        models: [
            {
                model: "qwen2:0.5b",
                params: {
                    top_p: 0.95,
                    stream: true,
                    num_keep: 5,
                    num_predict: 1,
                    top_k: 40,
                    temperature: 0.7,
                    stop: [
                        "<|im_start|>",
                        "<|im_end|>"
                    ]
                },
                info: {
                    size: "352MB",
                    desk: "1GB",
                    cpu: "8GB",
                    gpu: "6GB",
                    quant: "q4"
                }
            },
            {
                model: "qwen2:1.5b",
                params: {
                    top_p: 0.95,
                    stream: true,
                    num_keep: 5,
                    num_predict: 1,
                    top_k: 40,
                    temperature: 0.7,
                    stop: [
                        "<|im_start|>",
                        "<|im_end|>"
                    ]
                },
                info: {
                    size: "935MB",
                    desk: "1.5GB",
                    cpu: "8GB",
                    gpu: "6GB",
                    quant: "q4"
                }
            },
            {
                model: "qwen2:7b",
                params: {
                    top_p: 0.95,
                    stream: true,
                    num_keep: 5,
                    num_predict: 1,
                    top_k: 40,
                    temperature: 0.7,
                    stop: [
                        "<|im_start|>",
                        "<|im_end|>"
                    ]
                },
                info: {
                    size: "4.4GB",
                    desk: "6GB",
                    cpu: "16GB",
                    gpu: "8GB",
                    quant: "q4"
                }
            },
            {
                model: "qwen:0.5b",
                params: {
                    top_p: 0.95,
                    stream: true,
                    num_keep: 5,
                    num_predict: 1,
                    top_k: 40,
                    temperature: 0.7,
                    stop: [
                        "<|im_start|>",
                        "<|im_end|>"
                    ]
                },
                info: {
                    size: "395MB",
                    desk: "395MB",
                    cpu: "8GB",
                    gpu: "6GB",
                    quant: "q4"
                }
            },
        ],
        action: "chat",
        zhdesc: "Qwen是阿里云基于transformer的一系列大型语言模型，在大量数据上进行预训练，包括网络文本、书籍、代码等。",
        endesc: "Qwen is a series of transformer-based large language models by Alibaba Cloud, pre-trained on a large volume of data, including web texts, books, code, etc."
    },
    {
        name: "gemma",
        family: "gemma",
        engine: "ollama",
        from:"ollama",
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
                    size: "5.0GB",
                    desk: "6GB",
                    cpu: "16GB",
                    gpu: "8GB",
                    quant: "q4"
                }
            },
        ],
        action: "chat",
        zhdesc: "Gemma是由谷歌及其DeepMind团队开发的一个新的开放模型。",
        endesc: "Gemma is a new open model developed by Google and its DeepMind team."
    },

    {
        name: "llama",
        family: "llama",
        engine: "ollama",
        from:"ollama",
        action: "chat",
        models: [
            {
                model: "llama3.2:1b",
                params: {
                    top_p: 0.95,
                    stream: true,
                    num_predict: 1,
                    top_k: 40,
                    temperature: 0.7,
                    "num_keep": 24,
                    "stop": [
                        "<|start_header_id|>",
                        "<|end_header_id|>",
                        "<|eot_id|>"
                    ]
                },
                info: {
                    size: "1.3GB",
                    desk: "2GB",
                    cpu: "8GB",
                    gpu: "6GB",
                    quant: "q4"
                }
            },
            {
                model: "llama3.2:3b",
                params: {
                    top_p: 0.95,
                    stream: true,
                    num_predict: 1,
                    top_k: 40,
                    temperature: 0.7,
                    "num_keep": 24,
                    "stop": [
                        "<|start_header_id|>",
                        "<|end_header_id|>",
                        "<|eot_id|>"
                    ]
                },
                info: {
                    size: "3.2GB",
                    desk: "4GB",
                    cpu: "16GB",
                    gpu: "8GB",
                    quant: "q4"
                }
            },
            {
                model: "llama3:8b",
                params: {
                    top_p: 0.95,
                    stream: true,
                    num_predict: 1,
                    top_k: 40,
                    temperature: 0.7,
                    "num_keep": 24,
                    "stop": [
                        "<|start_header_id|>",
                        "<|end_header_id|>",
                        "<|eot_id|>"
                    ]
                },
                info: {
                    size: "4.7GB",
                    desk: "5GB",
                    cpu: "16GB",
                    gpu: "8GB",
                    quant: "q4"
                }
            },
            {
                model: "llama3-chatqa:8b",
                params: {
                    top_p: 0.95,
                    stream: true,
                    num_predict: 1,
                    top_k: 40,
                    temperature: 0.3,
                    "num_keep": 24,
                    "stop": [
                        "<|start_header_id|>",
                        "<|end_header_id|>",
                        "<|eot_id|>"
                    ]
                },
                info: {
                    size: "4.7GB",
                    desk: "5GB",
                    cpu: "16GB",
                    gpu: "8GB",
                    quant: "q4"
                }
            },
            {
                model: "llama3chinese:8b",
                url:["https://hf-mirror.com/shenzhi-wang/Llama3-8B-Chinese-Chat-GGUF-8bit/resolve/v2/Llama3-8B-Chinese-Chat-q8-v2.gguf"],
                type:"llm",
                from:"network",
                params: {
                    top_p: 0.95,
                    stream: true,
                    num_predict: 1,
                    top_k: 40,
                    temperature: 0.7,
                    "num_keep": 24,
                    "stop": [
                        "<|start_header_id|>",
                        "<|end_header_id|>",
                        "<|eot_id|>"
                    ]
                },
                info: {
                    size: "8.54GB",
                    desk: "9GB",
                    cpu: "16GB",
                    gpu: "8GB",
                    quant: "q8",
                    "parameters": "num_keep                       24\nstop                           \"<|start_header_id|>\"\nstop                           \"<|end_header_id|>\"\nstop                           \"<|eot_id|>\"",
                    "template": "{{ if .System }}<|start_header_id|>system<|end_header_id|>\n\n{{ .System }}<|eot_id|>{{ end }}{{ if .Prompt }}<|start_header_id|>user<|end_header_id|>\n\n{{ .Prompt }}<|eot_id|>{{ end }}<|start_header_id|>assistant<|end_header_id|>\n\n{{ .Response }}<|eot_id|>"
                }
            },
            {
                model: "llama2:7b",
                params: {
                    top_p: 0.95,
                    stream: true,
                    num_predict: 1,
                    top_k: 40,
                    temperature: 0.7,
                    "num_keep": 24,
                    "stop": [
                        "[INST]",
                        "[/INST]",
                        "<<SYS>>",
                        "<</SYS>>"
                    ]
                },
                info: {
                    size: "3.8GB",
                    desk: "4GB",
                    cpu: "16GB",
                    gpu: "8GB",
                    quant: "q4"
                }
            },
            {
                model: "llama2:13b",
                params: {
                    top_p: 0.95,
                    stream: true,
                    num_predict: 1,
                    top_k: 40,
                    temperature: 0.7,
                    "num_keep": 24,
                    "stop": [
                        "[INST]",
                        "[/INST]",
                        "<<SYS>>",
                        "<</SYS>>"
                    ]
                },
                info: {
                    size: "7.4GB",
                    desk: "8GB",
                    cpu: "32GB",
                    gpu: "12GB",
                    quant: "q4"
                }
            },
            {
                model: "llama2-chinese:7b",
                params: {
                    top_p: 0.95,
                    stream: true,
                    num_predict: 1,
                    top_k: 40,
                    temperature: 0.7,
                    "num_keep": 24,
                    "stop": [
                        "Name:",
                        "Assistant:"
                    ]
                },
                info: {
                    size: "3.8GB",
                    desk: "4GB",
                    cpu: "16GB",
                    gpu: "8GB",
                    quant: "q4"
                }
            },
        ],
        zhdesc: "Llama由Meta Platforms发布,在通用基准测试上优于许多可用的开源聊天模型。",
        endesc: "Llama is released by Meta Platforms, Inc.Llama 3 instruction-tuned models are fine-tuned and optimized for dialogue/chat use cases and outperform many of the available open-source chat models on common benchmarks."
    },
    {
        name: "internlm",
        family: "internlm",
        engine: "ollama",
        from:"ollama",
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
    },
    {
        name: "ming",
        family: "qwen",
        engine: "ollama",
        
        action: "chat",
        models: [
            {
                model: "ming:1.8B",
                url:["https://hf-mirror.com/capricornstone/MING-1.8B-Q8_0-GGUF/blob/main/ming-1.8b-q8_0.gguf"],
                type:"llm",
                from:"network",
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
    },
    {
        name: "mindchat",
        family: "qwen",
        engine: "ollama",
        from:"network",
        action: "chat",
        models: [
            {
                model: "MindChat-Qwen2:4b",
                url:["https://hf-mirror.com/v8karlo/MindChat-Qwen2-4B-Q5_K_M-GGUF/blob/main/mindchat-qwen2-4b-q5_k_m.gguf"],
                type:"llm",
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
    },
    {
        name: "llava",
        family: "llama",
        engine: "ollama",
        from:"ollama",
        action: "img2txt",
        models: [
            {
                model: "llava:7b",
                params: {
                    top_p: 0.95,
                    stream: true,
                    num_predict: 1,
                    top_k: 40,
                    temperature: 0.7,
                    "num_keep": 24,
                    "stop": [
                        "[INST]",
                        "[/INST]"
                    ]
                },
                info: {
                    size: "4.7GB",
                    desk: "5GB",
                    cpu: "16GB",
                    gpu: "8GB",
                    quant: "q4"
                }
            },
            {
                model: "llava:13b",
                params: {
                    top_p: 0.95,
                    stream: true,
                    num_predict: 1,
                    top_k: 40,
                    temperature: 0.7,
                    "num_keep": 24,
                    "stop": [
                        "[INST]",
                        "[/INST]"
                    ]
                },
                info: {
                    size: "8GB",
                    desk: "9GB",
                    cpu: "32GB",
                    gpu: "12GB",
                    quant: "q4"
                }
            },
        ],
        zhdesc: "LLaVA是一种新颖的端到端训练的大型多模式模型，它结合了视觉编码器和Vicuna，用于通用视觉和语言理解。",
        endesc: "LLaVA is a novel end-to-end trained large multimodal model that combines a vision encoder and Vicuna for general-purpose visual and language understanding. "
    },
    {
        name: "bakllava",
        family: "llama",
        engine: "ollama",
        from:"ollama",
        action: "img2txt",
        models: [
            {
                model: "bakllava:7b",
                params: {
                    top_p: 0.95,
                    stream: true,
                    num_predict: 1,
                    top_k: 40,
                    temperature: 0.7,
                    "num_ctx": 4096,
                    "stop": [
                        "</s>",
                        "USER:"
                    ]
                },
                info: {
                    size: "4.7GB",
                    desk: "5GB",
                    cpu: "16GB",
                    gpu: "8GB",
                    quant: "q4"
                }
            },
        ],
        zhdesc: "BakLLaVA是一个多模式模型，由Mistral 7B基础模型和LLaVA架构组成。",
        endesc: "BakLLaVA is a multimodal model consisting of the Mistral 7B base model augmented with the LLaVA architecture. "
    },
    {
        name: "minicpm",
        family: "llama",
        engine: "ollama",
        from: "ollama",
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
    },
    {
        name: "moondream",
        family: "moondream",
        engine: "ollama",
        from: "ollama",
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
    },
    {
        name: "phi",
        family: "phi",
        engine: "ollama",
        from:"ollama",
        action: "chat",
        models: [
            {
                model: "phi3:mini",
                params: {
                    top_p: 0.95,
                    stream: true,
                    num_predict: 1,
                    top_k: 40,
                    temperature: 0.7,
                    "stop": [
                        "<|end|>",
                        "<|user|>",
                        "<|assistant|>"
                    ]
                },
                info: {
                    size: "2.4GB",
                    desk: "3GB",
                    cpu: "8GB",
                    gpu: "6GB",
                    quant: "q4"
                }
            },
            {
                model: "phi3:medium",
                params: {
                    top_p: 0.95,
                    stream: true,
                    num_predict: 1,
                    top_k: 40,
                    temperature: 0.7,
                    "stop": [
                        "<|end|>",
                        "<|user|>",
                        "<|assistant|>"
                    ]
                },
                info: {
                    size: "7.9GB",
                    desk: "8GB",
                    cpu: "16GB",
                    gpu: "8GB",
                    quant: "q4"
                }
            },
        ],
        zhdesc: "Phi是微软开发的一系列开放式人工智能模型。",
        endesc: "Phi is a family of open AI models developed by Microsoft."
    },
    {
        name: "openchat",
        family: "llama",
        engine: "ollama",
        from:"ollama",
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
    },
    {
        name: "aya",
        family: "llama",
        engine: "ollama",
        from:"ollama",
        action: "translation",
        models: [
            {
                model: "aya-expanse:8b",
                params: {
                    top_p: 0.95,
                    stream: true,
                    num_predict: 1,
                    top_k: 40,
                    temperature: 0.7,
                    "stop": [
                        "<|START_OF_TURN_TOKEN|>",
                        "<|END_OF_TURN_TOKEN|>"
                    ]
                },
                info: {
                    size: "8.1GB",
                    desk: "9GB",
                    cpu: "16GB",
                    gpu: "8GB",
                    quant: "q4"
                }
            },
            {
                model: "aya:8b",
                params: {
                    top_p: 0.95,
                    stream: true,
                    num_predict: 1,
                    top_k: 40,
                    temperature: 0.7,
                    "stop": [
                        "<|START_OF_TURN_TOKEN|>",
                        "<|END_OF_TURN_TOKEN|>"
                    ]
                },
                info: {
                    size: "4.8GB",
                    desk: "5GB",
                    cpu: "16GB",
                    gpu: "8GB",
                    quant: "q4"
                }
            },
        ],
        zhdesc: "Aya 23可以流利地说23种语言。",
        endesc: "Aya 23 can talk upto 23 languages fluently."
    },
    {
        name: "bilibili",
        family: "llama",
        engine: "ollama",
        from: "ollama",
        action: "chat",
        models: [
            {
                model: "milkey/bilibili-index:latest",
                params: {
                    stream: true,
                    "repeat_penalty": 1.1,
                    "stop": [
                        "reserved_0",
                        "reserved_1",
                        "</s>",
                        "<unk>"
                    ],
                    "temperature": 0.3,
                    "top_k": 5,
                    "top_p": 0.8
                },
                info: {
                    size: "2.3GB",
                    desk: "3GB",
                    cpu: "8GB",
                    gpu: "6GB",
                    quant: "q8"
                }
            },
        ],
        zhdesc: "由哔哩哔哩自主研发的大语言模型",
        endesc: "A large language model independently developed by Bilibili"
    },
    {
        name: "yi",
        family: "yi",
        engine: "ollama",
        from:"ollama",
        action: "chat",
        models: [
            {
                model: "yi:6b",
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
                    size: "3.5GB",
                    desk: "4GB",
                    cpu: "16GB",
                    gpu: "8GB",
                    quant: "q4"
                }
            },
            {
                model: "yi:9b",
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
                    size: "5GB",
                    desk: "6GB",
                    cpu: "32GB",
                    gpu: "12GB",
                    quant: "q4"
                }
            },
        ],
        zhdesc: "Yi是一系列大型语言模型，在3万亿个高质量的语料库上训练，支持英语和汉语。",
        endesc: "Yi is a series of large language models trained on a high-quality corpus of 3 trillion tokens that support both the English and Chinese languages."
    },
    {
        name: "wizardlm2",
        family: "llama",
        engine: "ollama",
        from:"ollama",
        action: "chat",
        models: [
            {
                model: "wizardlm2:7b",
                params: {
                    top_p: 0.95,
                    stream: true,
                    num_predict: 1,
                    top_k: 40,
                    temperature: 0.7,
                    "stop": [
                        "USER:",
                        "ASSISTANT:"
                    ]
                },
                info: {
                    size: "4.1GB",
                    desk: "5GB",
                    cpu: "16GB",
                    gpu: "8GB",
                    quant: "q4"
                }
            },
        ],
        zhdesc: "微软人工智能的最先进的大型语言模型，在复杂聊天、多语言、推理和代理用例方面的性能有所提高。",
        endesc: "State of the art large language model from Microsoft AI with improved performance on complex chat, multilingual, reasoning and agent use cases. "
    },
    {
        name: "mistral",
        family: "llama",
        engine: "ollama",
        from:"ollama",
        action: "chat",
        models: [
            {
                model: "mistral:7b",
                params: {
                    top_p: 0.95,
                    stream: true,
                    num_predict: 1,
                    top_k: 40,
                    temperature: 0.7,
                    "stop": [
                        "[INST]",
                        "[/INST]"
                    ]
                },
                info: {
                    size: "4.1GB",
                    desk: "5GB",
                    cpu: "16GB",
                    gpu: "8GB",
                    quant: "q4"
                }
            },
        ],
        zhdesc: "Mistral是一个7.3B参数模型，使用Apache许可证进行分发。它有指令（指令如下）和文本完成两种形式。",
        endesc: "Mistral is a 7.3B parameter model, distributed with the Apache license. It is available in both instruct (instruction following) and text completion. "
    },
    {
        name: "mixtral",
        family: "llama",
        engine: "ollama",
        from:"ollama",
        action: "chat",
        models: [
            {
                model: "mixtral:8x7b",
                params: {
                    top_p: 0.95,
                    stream: true,
                    num_predict: 1,
                    top_k: 40,
                    temperature: 0.7,
                    "stop": [
                        "[INST]",
                        "[/INST]"
                    ]
                },
                info: {
                    size: "26GB",
                    desk: "27GB",
                    cpu: "32GB",
                    gpu: "12GB",
                    quant: "q4"
                }
            },
        ],
        zhdesc: "Mistral AI在8x7b和8x22b参数大小下的一组具有开放权重的专家混合（MoE）模型。",
        endesc: "A set of Mixture of Experts (MoE) model with open weights by Mistral AI in 8x7b and 8x22b parameter sizes. "
    },
    {
        name: "h2o",
        family: "llama",
        engine: "ollama",
        from:"ollama",
        action: "chat",
        models: [
            {
                model: "cas/h2o-danube2-1.8b-chat:latest",
                params: {
                    top_p: 0.95,
                    stream: true,
                    num_predict: 1,
                    top_k: 40,
                    temperature: 0.7,
                    "stop": [
                        "<|prompt|>",
                        "</s>",
                        "<|answer|>"
                    ]
                },
                info: {
                    size: "1.1GB",
                    desk: "2GB",
                    cpu: "8GB",
                    gpu: "6GB",
                    quant: "q4"
                }
            },
        ],
        zhdesc: "H2O.ai在Apache v2.0下发布了其最新的开放权重小语言模型H2O-Danube2-1.8B。",
        endesc: "H2O.ai just released its latest open-weight small language model, H2O-Danube2-1.8B, under Apache v2.0."
    },
    {
        name: "zephyr",
        family: "llama",
        engine: "ollama",
        from:"ollama",
        action: "chat",
        models: [
            {
                model: "zephyr:7b",
                params: {
                    top_p: 0.95,
                    stream: true,
                    num_predict: 1,
                    top_k: 40,
                    temperature: 0.7,
                    "stop": [
                        "<|system|>",
                        "<|user|>",
                        "<|assistant|>",
                        "</s>"
                    ]
                },
                info: {
                    size: "4.1GB",
                    desk: "5GB",
                    cpu: "16GB",
                    gpu: "8GB",
                    quant: "q4"
                }
            },
        ],
        zhdesc: "Zephyr是Mistral和Mixtral模型的一系列微调版本，经过训练，可以充当有用的助手。",
        endesc: "Zephyr is a series of fine-tuned versions of the Mistral and Mixtral models that are trained to act as helpful assistants.  "
    },
    {
        name: "solar",
        family: "llama",
        engine: "ollama",
        from:"ollama",
        action: "chat",
        models: [
            {
                model: "solar:10.7b",
                params: {
                    top_p: 0.95,
                    stream: true,
                    num_predict: 1,
                    top_k: 40,
                    temperature: 0.7,
                    "num_ctx": 4096,
                    "stop": [
                        "</s>",
                        "### System:",
                        "### User:",
                        "### Assistant:"
                    ]
                },
                info: {
                    size: "6.1GB",
                    desk: "7GB",
                    cpu: "32GB",
                    gpu: "12GB",
                    quant: "q4"
                }
            },
        ],
        zhdesc: "一个紧凑而强大的10.7B大型语言模型，专为单回合对话而设计。",
        endesc: "A compact, yet powerful 10.7B large language model designed for single-turn conversation. "
    },
    {
        name: "codegemma",
        family: "gemma",
        engine: "ollama",
        from:"ollama",
        action: "code",
        models: [
            {
                model: "codegemma:7b-instruct",
                params: {
                    top_p: 0.95,
                    stream: true,
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
                    size: "5GB",
                    desk: "5GB",
                    cpu: "16GB",
                    gpu: "8GB",
                    quant: "q4"
                }
            },
            {
                model: "codegemma:7b-code",
                params: {
                    top_p: 0.95,
                    stream: true,
                    num_predict: 1,
                    top_k: 40,
                    temperature: 0.7,
                    "penalize_newline": false,

                    "repeat_penalty": 1,
                    "stop": [
                        "<|fim_prefix|>",
                        "<|fim_suffix|>",
                        "<|fim_middle|>",
                        "<|file_separator|>"
                    ]
                },
                info: {
                    size: "5GB",
                    desk: "6GB",
                    cpu: "16GB",
                    gpu: "8GB",
                    quant: "q4"
                }
            },
        ],
        zhdesc: "CodeGemma是一个功能强大、轻量级的模型集合，可以执行各种编码任务，如填充中间代码完成、代码生成、自然语言理解、数学推理和指令遵循。",
        endesc: "CodeGemma is a collection of powerful, lightweight models that can perform a variety of coding tasks like fill-in-the-middle code completion, code generation, natural language understanding, mathematical reasoning, and instruction following. "
    },
    {
        name: "codeqwen",
        family: "gemma",
        engine: "ollama",
        from:"ollama",
        action: "code",
        models: [
            {
                model: "codeqwen:7b-chat",
                params: {
                    top_p: 0.95,
                    stream: true,
                    num_predict: 1,
                    top_k: 40,
                    temperature: 0.7,
                    "penalize_newline": false,
                    "repeat_penalty": 1,
                    "stop": [
                        "<|im_start|>",
                        "<|im_end|>"
                    ]
                },
                info: {
                    size: "4.2GB",
                    desk: "5GB",
                    cpu: "12GB",
                    gpu: "8GB",
                    quant: "q4"
                }
            },
            {
                model: "codeqwen:7b-code",
                params: {
                    top_p: 0.95,
                    stream: true,
                    num_predict: 1,
                    top_k: 40,
                    temperature: 0.7,
                    "penalize_newline": false,
                    "repeat_penalty": 1,
                    "stop": [
                        "<|im_start|>",
                        "<|im_end|>"
                    ]
                },
                info: {
                    size: "4.2GB",
                    desk: "5GB",
                    cpu: "12GB",
                    gpu: "8GB",
                    quant: "q4"
                }
            },
        ],
        zhdesc: "CodeQwen是一个在大量代码数据上预训练的大型语言模型。",
        endesc: "CodeQwen is a large language model pretrained on a large amount of code data. "
    },
    {
        name: "codellama",
        family: "llama",
        engine: "ollama",
        from:"ollama",
        action: "code",
        models: [
            {
                model: "codellama:7b-instruct",
                params: {
                    top_p: 0.95,
                    stream: true,
                    num_predict: 1,
                    top_k: 40,
                    temperature: 0.7,
                    "rope_frequency_base": 1000000,
                    "stop": [
                        "[INST]",
                        "[/INST]",
                        "<<SYS>>",
                        "<</SYS>>"
                    ]
                },
                info: {
                    size: "3.8GB",
                    desk: "4GB",
                    cpu: "16GB",
                    gpu: "8GB",
                    quant: "q4"
                }
            },
            {
                model: "codellama:7b-code",
                params: {
                    top_p: 0.95,
                    stream: true,
                    num_predict: 1,
                    top_k: 40,
                    temperature: 0.7,
                    "rope_frequency_base": 1000000,
                    "stop": [
                        "[INST]",
                        "[/INST]",
                        "<<SYS>>",
                        "<</SYS>>"
                    ]
                },
                info: {
                    size: "3.8GB",
                    desk: "4GB",
                    cpu: "12GB",
                    gpu: "8GB",
                    quant: "q4"
                }
            },
        ],
        zhdesc: "Code Llama是一个用于生成和讨论代码的模型，构建在Llama 2之上。它旨在使开发人员的工作流程更快、更高效，并使人们更容易学习如何编码。它可以生成代码和关于代码的自然语言。Code Llama支持当今使用的许多最流行的编程语言，包括Python、C++、Java、PHP、Typescript（Javascript）、C#、Bash等。",
        endesc: "Code Llama is a model for generating and discussing code, built on top of Llama 2. It’s designed to make workflows faster and efficient for developers and make it easier for people to learn how to code. It can generate both code and natural language about code. Code Llama supports many of the most popular programming languages used today, including Python, C++, Java, PHP, Typescript (Javascript), C#, Bash and more."
    },
    {
        name: "deepseek-coder",
        family: "gemma",
        engine: "ollama",
        from:"ollama",
        action: "code",
        models: [
            {
                model: "deepseek-coder:1.3b-instruct-q8_0",
                params: {
                    top_p: 0.95,
                    stream: true,
                    num_predict: 1,
                    top_k: 40,
                    temperature: 0.7,
                    "stop": []
                },
                info: {
                    size: "1.4GB",
                    desk: "2GB",
                    cpu: "8GB",
                    gpu: "6GB",
                    quant: "q4"
                }
            },
            {
                model: "deepseek-coder:1.3b-base-q8_0",
                params: {
                    top_p: 0.95,
                    stream: true,
                    num_predict: 1,
                    top_k: 40,
                    temperature: 0.7,
                    "stop": []
                },
                info: {
                    size: "1.4GB",
                    desk: "2GB",
                    cpu: "8GB",
                    gpu: "6GB",
                    quant: "q4"
                }
            },
            {
                model: "deepseek-coder:6.7b-instruct-q8_0",
                params: {
                    top_p: 0.95,
                    stream: true,
                    num_predict: 1,
                    top_k: 40,
                    temperature: 0.7,
                    "stop": []
                },
                info: {
                    size: "7.2GB",
                    desk: "8GB",
                    cpu: "16GB",
                    gpu: "8GB",
                    quant: "q4"
                }
            },
            {
                model: "deepseek-coder:6.7b-base-q8_0",
                params: {
                    top_p: 0.95,
                    stream: true,
                    num_predict: 1,
                    top_k: 40,
                    temperature: 0.7,
                    "stop": []
                },
                info: {
                    size: "7.2GB",
                    desk: "8GB",
                    cpu: "16GB",
                    gpu: "8GB",
                    quant: "q4"
                }
            },
            {
                model: "deepseek-coder-v2:16b",
                params: {
                    top_p: 0.95,
                    stream: true,
                    num_predict: 1,
                    top_k: 40,
                    temperature: 0.7,
                    "stop": []
                },
                info: {
                    size: "8.9GB",
                    desk: "9GB",
                    cpu: "32GB",
                    gpu: "12GB",
                    quant: "f16"
                }
            },
        ],
        zhdesc: "DeepSeek Coder是一个基于两万亿代码和自然语言标记的强大编码模型。",
        endesc: "DeepSeek Coder is a capable coding model trained on two trillion code and natural language tokens. "
    },
    {
        name: "starcoder2",
        family: "starcoder2",
        engine: "ollama",
        from:"ollama",
        action: "code",
        models: [
            {
                model: "starcoder2:3b",
                params: {
                    top_p: 0.95,
                    stream: true,
                    num_predict: 1,
                    top_k: 40,
                    temperature: 0.7,
                    "stop": []
                },
                info: {
                    size: "1.7GB",
                    desk: "2GB",
                    cpu: "8GB",
                    gpu: "6GB",
                    quant: "q4"
                }
            },
            {
                model: "starcoder2:7b",
                params: {
                    top_p: 0.95,
                    stream: true,
                    num_predict: 1,
                    top_k: 40,
                    temperature: 0.7,
                    "stop": []
                },
                info: {
                    size: "4GB",
                    desk: "5GB",
                    cpu: "16GB",
                    gpu: "8GB",
                    quant: "q4"
                }
            },
        ],
        zhdesc: "StarCoder2是下一代经过透明训练的开放代码LLM，有三种大小：3B、7B和15B参数。",
        endesc: "StarCoder2 is the next generation of transparently trained open code LLMs that comes in three sizes: 3B, 7B and 15B parameters. "
    },
    {
        name: "duckdb-nsql",
        family: "llama",
        engine: "ollama",
        from:"ollama",
        action: "code",
        models: [
            {
                model: "duckdb-nsql:7b",
                params: {
                    top_p: 0.95,
                    stream: true,
                    num_predict: 1,
                    top_k: 40,
                    temperature: 0.7,
                    "stop": []
                },
                info: {
                    size: "3.8GB",
                    desk: "4GB",
                    cpu: "16GB",
                    gpu: "8GB",
                    quant: "q4"
                }
            },
        ],
        zhdesc: "由MotherDuck和Numbers Station制作的7B参数文本到SQL模型。",
        endesc: "7B parameter text-to-SQL model made by MotherDuck and Numbers Station. "
    },
    {
        name: "bge-large-zh-v1.5",
        family: "bge",
        engine: "ollama",
        from:"ollama",
        action: "embeddings",
        models: [
            {
                model: "quentinz/bge-large-zh-v1.5:latest",
                params: {
                    "num_ctx": 512
                },
                info: {
                    size: "651MB",
                    desk: "1GB",
                    cpu: "8GB",
                    gpu: "6GB",
                    quant: "f16"
                }
            },
            {
                model: "quentinz/bge-base-zh-v1.5:latest",
                params: {
                    "num_ctx": 512
                },
                info: {
                    size: "205MB",
                    desk: "1GB",
                    cpu: "8GB",
                    gpu: "6GB",
                    quant: "f16"
                }
            },
        ],
        zhdesc: "bge是BAAI开发的嵌入模型",
        endesc: "bge is an embedded model developed by BAAI"
    },
    {
        name: "dmeta-embedding-zh",
        family: "dmeta",
        engine: "ollama",
        from:"ollama",
        action: "embeddings",
        models: [
            {
                model: "herald/dmeta-embedding-zh:latest",
                params: {
                    "num_ctx": 1024
                },
                info: {
                    size: "205MB",
                    desk: "300MB",
                    cpu: "8GB",
                    gpu: "6GB",
                    quant: "f16"
                }
            },
        ],
        zhdesc: "Dmeta-embedding 是一款跨领域、跨任务、开箱即用的中文 Embedding 模型，适用于搜索、问答、智能客服、LLM+RAG 等各种业务场景",
        endesc: "Dmeta-embedding is a cross domain, cross task, and out of the box Chinese embedding model suitable for various business scenarios such as search, Q&A, intelligent customer service, LLM+RAG, etc"
    },
    {
        name: "nomic-embed-text",
        family: "nomic-bert",
        engine: "ollama",
        from:"ollama",
        action: "embeddings",
        models: [
            {
                model: "nomic-embed-text:latest",
                params: {
                    "num_ctx": 768
                },
                info: {
                    size: "274MB",
                    desk: "300MB",
                    cpu: "8GB",
                    gpu: "6GB",
                    quant: "f16"
                }
            },
        ],
        zhdesc: "一个具有大型令牌上下文窗口的高性能开放嵌入模型。",
        endesc: "A high-performing open embedding model with a large token context window. "
    },
    {
        name: "snowflake-arctic-embed",
        family: "bert",
        engine: "ollama",
        from:"ollama",
        action: "embeddings",
        models: [
            {
                model: "snowflake-arctic-embed:latest",
                params: {
                    "num_ctx": 1024
                },
                info: {
                    size: "669MB",
                    desk: "700MB",
                    cpu: "8GB",
                    gpu: "6GB",
                    quant: "f16"
                }
            },
        ],
        zhdesc: "snowflake-arctic-embed是一套文本嵌入模型，专注于创建针对性能优化的高质量检索模型。",
        endesc: "snowflake-arctic-embed is a suite of text embedding models that focuses on creating high-quality retrieval models optimized for performance. "
    },
    {
        name: "mxbai-embed-large",
        family: "bert",
        engine: "ollama",
        from:"ollama",
        action: "embeddings",
        models: [
            {
                model: "mxbai-embed-large:latest",
                params: {
                    "num_ctx": 512
                },
                info: {
                    size: "670MB",
                    desk: "700MB",
                    cpu: "8GB",
                    gpu: "6GB",
                    quant: "f16"
                }
            },
        ],
        zhdesc: "mixedbread.ai的最先进的大型嵌入模型",
        endesc: " State-of-the-art large embedding model from mixedbread.ai"
    },
    sdLabel,
    telespeechLabel,
    whisperLabel,
    nemoLabel,
    zipformerLabel,
    paraformerLabel,
    vitsLabel
]
