export const nemoLabel = {
    name: "nemo",
    family: "lstm",
    engine: "voice",
    from:"network",
    action: ["audio"],
    models:[
        {
            model:"nomo",
            file_name: "model.int8.onnx",
            url: [
                "https://hf-mirror.com/csukuangfj/sherpa-onnx-nemo-ctc-zh-citrinet-1024-gamma-0-25/resolve/main/model.int8.onnx",
                "https://hf-mirror.com/csukuangfj/sherpa-onnx-nemo-ctc-zh-citrinet-1024-gamma-0-25/resolve/main/tokens.txt"
            ],
            params:{
                type:"nomo",
                model:"model.int8.onnx",
                token:"tokens.txt",
            },
            info:{
                size:"147MB",
                desk: "200MB",
                cpu: "8GB",
                gpu: "6GB",
                quant: "q8"
            }
        }
    ],
    zhdesc: " NVIDIA NeMo是NVIDIA AI平台的一部分,是一个用于构建新的最先进对话式AI模型。",
    endesc: "NVIDIA NeMo is part of the NVIDIA AI platform and is used to build new state-of-the-art conversational AI models."
}