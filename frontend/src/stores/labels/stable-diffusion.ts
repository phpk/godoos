export const sdLabel = {
    name: "stable-diffusion",
    family: "stable-diffusion",
    engine: "sd",
    from:"network",
    action: ["image"],
    models: [
        {
            model: "stable-diffusion-v-1-4",
            file_name:"sd-v1-4.ckpt",
            url: ["https://hf-mirror.com/CompVis/stable-diffusion-v-1-4-original/resolve/main/sd-v1-4.ckpt"],
            info: {
                size: "4.27GB",
                desk: "5GB",
                cpu: "16GB",
                gpu: "8GB",
                quant: "f32"
            }
        },
        {
            model: "stable-diffusion-v1-5",
            file_name: "v1-5-pruned-emaonly.safetensors",
            url:[
                "https://hf-mirror.com/runwayml/stable-diffusion-v1-5/resolve/main/v1-5-pruned-emaonly.safetensors",
                "https://hf-mirror.com/madebyollin/taesd/blob/main/diffusion_pytorch_model.safetensors"
            ],
            info: {
                size: "4.27GB",
                desk: "5GB",
                cpu: "16GB",
                gpu: "8GB",
                quant: "f32"
            }
        },
        {
            model: "stable-diffusion-2-1",
            file_name:"v2-1_768-nonema-pruned.safetensors",
            url: ["https://hf-mirror.com/stabilityai/stable-diffusion-2-1/resolve/main/v2-1_768-nonema-pruned.safetensors"],
            info: {
                size: "5.21GB",
                desk: "6GB",
                cpu: "16GB",
                gpu: "8GB",
                quant: "f32"
            }
        },
        {
            model: "sd_xl_base_1.0",
            file_name: "sd_xl_base_1.0.safetensors",
            url:[
                "https://hf-mirror.com/stabilityai/stable-diffusion-xl-base-1.0/resolve/main/sd_xl_base_1.0.safetensors",
                "https://hf-mirror.com/madebyollin/sdxl-vae-fp16-fix/blob/main/sdxl_vae.safetensors"
            ],
            info: {
                size: "6.94GB",
                desk: "7GB",
                cpu: "16GB",
                gpu: "8GB",
                quant: "f32"
            }
        },
        {
            model: "sd_xl_turbo_1.0",
            file_name: "sd_xl_turbo_1.0_fp16.safetensors",
            url:[
                "https://hf-mirror.com/stabilityai/sdxl-turbo/resolve/main/sd_xl_turbo_1.0_fp16.safetensors",
                "https://hf-mirror.com/madebyollin/sdxl-vae-fp16-fix/blob/main/sdxl_vae.safetensors"
            ],
            info: {
                size: "6.94GB",
                desk: "7GB",
                cpu: "16GB",
                gpu: "8GB",
                quant: "f16"
            }
        }

    ],
    zhdesc: "Stable Diffusion 是一款支持由文本生成图像的 AI 绘画工具，它主要用于根据文本描述生成对应图像的任务",
    endesc: "Stable Diffusion is an AI drawing tool that supports generating images from text. It is mainly used to generate corresponding images based on text descriptions "
}