export const vitsLabel = {
    name: "vits",
    family: "vits",
    engine: "voice",
    from:"network",
    action: ["tts"],
    models:[
        {
            model: "vits-zh-aishell3",
            file_name: "vits-aishell3.int8.onnx",
            url: [
                "https://hf-mirror.com/csukuangfj/vits-zh-aishell3/resolve/main/vits-aishell3.int8.onnx",
                "https://hf-mirror.com/csukuangfj/vits-zh-aishell3/resolve/main/tokens.txt",
                "https://hf-mirror.com/csukuangfj/vits-zh-aishell3/resolve/main/phone.fst",
                "https://hf-mirror.com/csukuangfj/vits-zh-aishell3/resolve/main/number.fst",
                "https://hf-mirror.com/csukuangfj/vits-zh-aishell3/resolve/main/lexicon.txt",
                "https://hf-mirror.com/csukuangfj/vits-zh-aishell3/resolve/main/date.fst"
            ],
            params : {
                type:"vits",
                model:"vits-aishell3.int8.onnx",
                token:"tokens.txt",
                lexicon:"lexicon.txt",
                ruleFsts:["date.fst","phone.fst","number.fst"]
            },
            info: {
                size: "121MB",
                desk: "200MB",
                cpu: "8GB",
                gpu: "6GB",
                quant: "q8"
            }
        }
    ],
    zhdesc: "VITS (Voice, Intent, and Text Space) 是一个用于端到端文本到语音合成的模型。",
    endesc: "VITS (Voice, Intent, and Text Space is a model used for end-to-end text to speech synthesis."
}