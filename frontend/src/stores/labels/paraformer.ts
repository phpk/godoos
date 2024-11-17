export const paraformerLabel = {
    name: "paraformer",
    family: "paraformer",
    engine: "voice",
    from:"network",
    action: ["audio"],
    models: [
        {
            model: "paraformer",
            file_name: "model.int8.onnx",
            url: [
                "https://hf-mirror.com/csukuangfj/sherpa-onnx-paraformer-zh-2023-03-28/resolve/main/model.int8.onnx",
                "https://hf-mirror.com/csukuangfj/sherpa-onnx-paraformer-zh-2023-03-28/resolve/main/tokens.txt"
            ],
            params:{
                type:"paraformer",
                model:"model.int8.onnx",
                token:"tokens.txt",
            },
            info: {
                size: "223MB",
                desk: "300MB",
                cpu: "8GB",
                gpu: "6GB",
                quant: "q8"
            }
        }
    ],
    zhdesc: "Paraformer是通义实验室研发的新一代非自回归端到端语音识别模型,具有识别准确率高、推理效率高的特点。",
    endesc: "Paraformer is a new generation of non autoregressive end-to-end speech recognition model developed by Tongyi Laboratory, which has the characteristics of high recognition accuracy and high inference efficiency."
}