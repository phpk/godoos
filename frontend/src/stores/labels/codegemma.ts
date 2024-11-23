export const codegemmaLabels = {
    name: "codegemma",
    family: "gemma",
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
                engine: "ollama",
                from: "ollama",
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
                engine: "ollama",
                from: "ollama",
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
}