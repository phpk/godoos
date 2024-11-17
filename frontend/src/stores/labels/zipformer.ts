export const zipformerLabel = {
    name: "zipformer",
    family: "zipformer",
    engine: "voice",
    from:"network",
    action: ["audio"],
    models:[
        {
            model: "zipformer-en",
            file_name: "encoder-epoch-99-avg-1.int8.onnx",
            url: [
                "https://hf-mirror.com/csukuangfj/sherpa-onnx-zipformer-en-2023-06-26/resolve/main/encoder-epoch-99-avg-1.int8.onnx",
                "https://hf-mirror.com/csukuangfj/sherpa-onnx-zipformer-en-2023-06-26/resolve/main/decoder-epoch-99-avg-1.int8.onnx",
                "https://hf-mirror.com/csukuangfj/sherpa-onnx-zipformer-en-2023-06-26/resolve/main/joiner-epoch-99-avg-1.int8.onnx",
                "https://hf-mirror.com/csukuangfj/sherpa-onnx-zipformer-en-2023-06-26/resolve/main/tokens.txt"
            ],
            params:{
                type:"zipformer",
                decoder:"decoder-epoch-99-avg-1.int8.onnx",
                encoder:"encoder-epoch-99-avg-1.int8.onnx",
                joiner:"joiner-epoch-99-avg-1.int8.onnx",
                token:"tokens.txt",
            },
            info: {
                size: "76MB",
                desk: "100MB",
                cpu: "8GB",
                gpu: "6GB",
                quant: "q8"
            }
        }
    ],
    zhdesc: "Zipformer 模型是新一代 Kaldi 团队提出的新型声学建模架构。",
    endesc: "The Zipformer model is a new acoustic modeling architecture proposed by the new generation Kaldi team."
}