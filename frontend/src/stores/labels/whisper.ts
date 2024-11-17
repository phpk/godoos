export const whisperLabel = {
    name: "whisper",
    family: "whisper",
    engine: "voice",
    from:"network",
    action: ["audio"],
    models:[
        {
            model: "whisper-tiny",
            file_name: "tiny-decoder.int8.onnx",
            url: [
                "https://hf-mirror.com/csukuangfj/sherpa-onnx-whisper-tiny/resolve/main/tiny-decoder.int8.onnx",
                "https://hf-mirror.com/csukuangfj/sherpa-onnx-whisper-tiny/resolve/main/tiny-encoder.int8.onnx",
                "https://hf-mirror.com/csukuangfj/sherpa-onnx-whisper-tiny/resolve/main/tiny-tokens.txt"
            ],
            params:{
                type:"whisper",
                decoder:"tiny-decoder.int8.onnx",
                encoder:"tiny-encoder.int8.onnx",
                token:"tiny-tokens.txt",
            },
            info: {
                size: "103MB",
                desk: "200MB",
                cpu: "8GB",
                gpu: "6GB",
                quant: "q8"
            }
        },
        {
            model: "whisper-base",
            file_name: "base-decoder.int8.onnx",
            url: [
                "https://hf-mirror.com/csukuangfj/sherpa-onnx-whisper-base/resolve/main/base-decoder.int8.onnx",
                "https://hf-mirror.com/csukuangfj/sherpa-onnx-whisper-base/resolve/main/base-encoder.int8.onnx",
                "https://hf-mirror.com/csukuangfj/sherpa-onnx-whisper-base/resolve/main/base-tokens.txt"
            ],
            params:{
                type:"whisper",
                decoder:"base-decoder.int8.onnx",
                encoder:"base-encoder.int8.onnx",
                token:"base-tokens.txt",
            },
            info: {
                size: "151MB",
                desk: "200MB",
                cpu: "8GB",
                gpu: "6GB",
                quant: "q8"
            }
        },
        {
            model: "whisper-small",
            file_name: "small-decoder.int8.onnx",
            url: [
                "https://hf-mirror.com/csukuangfj/sherpa-onnx-whisper-small/resolve/main/small-decoder.int8.onnx",
                "https://hf-mirror.com/csukuangfj/sherpa-onnx-whisper-small/resolve/main/small-encoder.int8.onnx",
                "https://hf-mirror.com/csukuangfj/sherpa-onnx-whisper-small/resolve/main/small-tokens.txt"
            ],
            params:{
                type:"whisper",
                decoder:"small-decoder.int8.onnx",
                encoder:"small-encoder.int8.onnx",
                token:"small-tokens.txt",
            },
            info: {
                size: "374MB",
                desk: "400MB",
                cpu: "8GB",
                gpu: "6GB",
                quant: "q8"
            }
        },
        {
            model: "whisper-medium",
            file_name: "medium-decoder.int8.onnx",
            url: [
                "https://hf-mirror.com/csukuangfj/sherpa-onnx-whisper-medium/resolve/main/medium-decoder.int8.onnx",
                "https://hf-mirror.com/csukuangfj/sherpa-onnx-whisper-medium/resolve/main/medium-encoder.int8.onnx",
                "https://hf-mirror.com/csukuangfj/sherpa-onnx-whisper-medium/resolve/main/medium-tokens.txt"
            ],
            params:{
                type:"whisper",
                decoder:"medium-decoder.int8.onnx",
                encoder:"mdeium-encoder.int8.onnx",
                token:"medium-tokens.txt",
            },
            info: {
                size: "945MB",
                desk: "1000MB",
                cpu: "8GB",
                gpu: "6GB",
                quant: "q8"
            }
        }
    ],
   
    zhdesc: "Whisper是一种通用的语音识别模型，由OpenAI研发并开源。它是在包含各种音频的大型数据集上训练的，可以执行多语言语音识别、语音翻译和语言识别。",
    endesc: "Whisper is a universal speech recognition model developed and open-source by OpenAI. It is trained on large datasets containing various types of audio and can perform multilingual speech recognition, speech translation, and language recognition."
}