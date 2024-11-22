export const duckdbnsqlLabels = {
    name: "duckdb-nsql",
    family: "llama",
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
                engine: "ollama",
                from: "ollama",
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
}