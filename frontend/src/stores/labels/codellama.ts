export const codellamaLabels = {
    name: "codellama",
    family: "llama",

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
                engine: "ollama",
                from: "ollama",
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
                engine: "ollama",
                from: "ollama",
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
}