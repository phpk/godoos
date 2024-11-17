export const telespeechLabel = {
    name: "telespeech",
    family: "telespeech",
    engine: "voice",
    from:"network",
    action: ["audio"],
    models: [
        // {
        //     model: "telespeech",
        //     file_name: "model.onnx",
        //     url: [
        //         "https://hf-mirror.com/csukuangfj/sherpa-onnx-telespeech-ctc-zh-2024-06-04/blob/main/model.onnx",
        //         "https://hf-mirror.com/csukuangfj/sherpa-onnx-telespeech-ctc-zh-2024-06-04/blob/main/tokens.txt"
        //     ],
        //     params:{
        //         type:"telespeech",
        //         model:"model.onnx",
        //         token:"tokens.txt",
        //     },
        //     info: {
        //         size: "341MB",
        //         desk: "400MB",
        //         cpu: "8GB",
        //         gpu: "6GB",
        //         quant: "q8"
        //     }
        // },
        {
            model: "telespeech-int8",
            file_name: "model.int8.onnx",
            url: [
                "https://hf-mirror.com/csukuangfj/sherpa-onnx-telespeech-ctc-int8-zh-2024-06-04/blob/main/model.int8.onnx",
                "https://hf-mirror.com/csukuangfj/sherpa-onnx-telespeech-ctc-int8-zh-2024-06-04/blob/main/tokens.txt"
            ],
            params:{
                type:"telespeech",
                model:"model.int8.onnx",
                token:"tokens.txt",
            },
            info: {
                size: "341MB",
                desk: "400MB",
                cpu: "8GB",
                gpu: "6GB",
                quant: "q8"
            }
        }
    ],
    zhdesc: "星辰语义大模型是由中电信人工智能科技有限公司研发训练的大语言模型，采用1.5万亿 Tokens中英文高质量语料进行训练。",
    endesc: "The Star Semantic Big Model is a large language model developed and trained by China Telecom Artificial Intelligence Technology Co., Ltd. It uses high-quality corpus of 1.5 trillion tokens in both Chinese and English for training."
}