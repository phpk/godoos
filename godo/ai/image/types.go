/*
 * GodoAI - A software focused on localizing AI applications
 * Copyright (C) 2024 https://godoos.com
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 2.1 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */
package image

type CLIConfig struct {
	Action           string  `json:"action,omitempty"`              // M, --mode
	Threads          int     `json:"threads,omitempty"`             // -t, --threads
	Model            string  `json:"model,omitempty"`               // -m, --model
	FileName         string  `json:"file_name,omitempty"`           // file_name
	VAE              string  `json:"vae,omitempty"`                 // --vae
	TAESD            string  `json:"taesd,omitempty"`               // --taesd
	ControlNet       string  `json:"control_net,omitempty"`         // --control-net
	EmbeddingDir     string  `json:"embedding_dir,omitempty"`       // --embd-dir
	StackedIDEmbDir  string  `json:"stacked_id_embd_dir,omitempty"` // --stacked-id-embd-dir
	InputIDImagesDir string  `json:"input_id_images_dir,omitempty"` // --input-id-images-dir
	NormalizeInput   bool    `json:"normalize_input,omitempty"`     // --normalize-input
	UpscaleModel     string  `json:"upscale_model,omitempty"`       // --upscale-model
	UpscaleRepeats   int     `json:"upscale_repeats,omitempty"`     // --upscale-repeats
	Type             string  `json:"type,omitempty"`                // --type
	LoraModelDir     string  `json:"lora_model_dir,omitempty"`      // --lora-model-dir
	InitImg          string  `json:"img_path,omitempty"`            // -i, --init-img
	ControlImage     string  `json:"control_image,omitempty"`       // --control-image
	Output           string  `json:"output,omitempty"`              // -o, --output
	Prompt           string  `json:"prompt,omitempty"`              // -p, --prompt
	NegativePrompt   string  `json:"negative_prompt,omitempty"`     // -n, --negative-prompt
	CFGScale         float32 `json:"cfg_scale,omitempty"`           // --cfg-scale
	Strength         float32 `json:"strength,omitempty"`            // --strength
	StyleRatio       int     `json:"style_ratio,omitempty"`         // --style-ratio
	CtrlStrength     float32 `json:"control_strength,omitempty"`    // --control-strength
	Height           int     `json:"height,omitempty"`              // -H, --height
	Width            int     `json:"width,omitempty"`               // -W, --width
	SamplingMethod   string  `json:"sampling_method,omitempty"`     // --sampling-method
	Steps            int     `json:"steps,omitempty"`               // --steps
	RNG              string  `json:"rng,omitempty"`                 // --rng
	Seed             int     `json:"seed,omitempty"`                // -s, --seed
	BatchCount       int     `json:"num,omitempty"`                 // -b, --batch-count
	Schedule         string  `json:"schedule,omitempty"`            // --schedule
	ClipSkip         int     `json:"clip_skip,omitempty"`           // --clip-skip
	VaeTiling        bool    `json:"vae_tiling,omitempty"`          // --vae-tiling
	CtrlNetCPU       bool    `json:"ctrl_net_cpu,omitempty"`        // --control-net-cpu
	Canny            bool    `json:"canny,omitempty"`               // --canny
	Color            bool    `json:"color,omitempty"`               // --color
	Verbose          bool    `json:"verbose,omitempty"`             // -v, --verbose
}

// func ApplyDefaults(config *CLIConfig) ([]string, error) {
// 	params := []string{}
// 	if config.Verbose {
// 		params = append(params, "-v")
// 	}
// 	// Set default mode if not provided txt2img or img2img
// 	if config.Mode == "" {
// 		config.Mode = "txt2img"
// 	}
// 	if config.Mode != "txt2img" && config.Mode != "img2img" {
// 		return params, fmt.Errorf("invalid mode: %s", config.Mode)
// 	}
// 	params = append(params, "-m", config.Mode)
// 	if config.Model != "" {
// 		modelPath, err := GetModelPath(config.Model, config.FileName)
// 		if err != nil {
// 			return params, fmt.Errorf("error get modelpath: %s", config.Mode)
// 		}
// 		params = append(params, "-m", modelPath)
// 		if config.Model == "sd_xl_turbo_1.0" || config.Model == "sd_xl_base_1.0" {
// 			aesPath, err := GetModelPath(config.Model, "sdxl_vae.safetensors")
// 			if err != nil {
// 				return params, fmt.Errorf("error get aespath: %s", config.Mode)
// 			}
// 			config.VAE = aesPath
// 		}
// 		if config.Model == "stable-diffusion-v1-5" {
// 			tsaesdPath, err := GetModelPath(config.Model, "diffusion_pytorch_model.safetensors")
// 			if err != nil {
// 				return params, fmt.Errorf("error get aespath: %s", config.Mode)
// 			}
// 			config.TAESD = tsaesdPath
// 		}
// 		// 添加判断逻辑
// 		lowerModelPath := strings.ToLower(config.Model)
// 		if strings.Contains(lowerModelPath, "f16") {
// 			config.Type = "f16"
// 		}
// 		if strings.Contains(lowerModelPath, "q8_0") {
// 			config.Type = "q8_0"
// 		}
// 		if strings.Contains(lowerModelPath, "q5_0") {
// 			config.Type = "q5_0"
// 		}
// 		if strings.Contains(lowerModelPath, "q4_0") {
// 			config.Type = "q4_0"
// 		}
// 	}
// 	if config.Type == "" {
// 		config.Type = "f32"
// 	}
// 	params = append(params, "--type", config.Type)
// 	// Set threads to CPU core count if not provided or <= 0
// 	if config.Threads == 0 {
// 		config.Threads = runtime.NumCPU()
// 	}
// 	if config.Threads > 0 {
// 		params = append(params, "-t", fmt.Sprintf("%d", config.Threads))
// 	}
// 	if config.BatchCount == 0 {
// 		config.BatchCount = 1
// 	}
// 	params = append(params, "-b", fmt.Sprintf("%d", config.BatchCount))
// 	// Default output path
// 	if config.Output == "" {
// 		config.Output = "./output.png"
// 	}
// 	params = append(params, "-o", config.Output)
// 	// Default prompt and negative prompt
// 	if config.Prompt == "" {
// 		config.Prompt = "a beautiful landscape"
// 	}
// 	params = append(params, "-p", fmt.Sprintf(`"%s"`, config.Prompt))
// 	if config.NegativePrompt != "" {
// 		params = append(params, "-n", fmt.Sprintf(`"%s"`, config.NegativePrompt))
// 	}

// 	// Other defaults
// 	if config.CFGScale == 0 {
// 		config.CFGScale = 7.0
// 	}
// 	params = append(params, "--cfg-scale", fmt.Sprintf("%0.1f", config.CFGScale))
// 	if config.Strength == 0 {
// 		config.Strength = 0.75
// 	}
// 	params = append(params, "--strength", fmt.Sprintf("%0.2f", config.Strength))
// 	if config.StyleRatio == 0 {
// 		config.StyleRatio = 20
// 	}
// 	params = append(params, "--style-ratio", fmt.Sprintf("%d", config.StyleRatio))
// 	if config.CtrlStrength == 0 {
// 		config.CtrlStrength = 0.9
// 	}
// 	params = append(params, "--control-strength", fmt.Sprintf("%0.1f", config.CtrlStrength))
// 	if config.Height == 0 {
// 		config.Height = 512
// 	}
// 	params = append(params, "-H", fmt.Sprintf("%d", config.Height))
// 	if config.Width == 0 {
// 		config.Width = 512
// 	}
// 	params = append(params, "-W", fmt.Sprintf("%d", config.Width))
// 	if config.SamplingMethod == "" {
// 		config.SamplingMethod = "euler_a"
// 	}
// 	params = append(params, "--sampling-method", config.SamplingMethod)
// 	if config.Steps == 0 {
// 		config.Steps = 20
// 	}
// 	params = append(params, "--steps", fmt.Sprintf("%d", config.Steps))
// 	if config.RNG == "" {
// 		config.RNG = "cuda"
// 	}
// 	params = append(params, "--rng", config.RNG)
// 	if config.Seed == 0 {
// 		config.Seed = 42
// 	}
// 	params = append(params, "-s", fmt.Sprintf("%d", config.Seed))

// 	if config.Schedule == "" {
// 		config.Schedule = "discrete"
// 	}
// 	params = append(params, "--schedule", config.Schedule)
// 	// Ensure paths are absolute if provided

// 	// Convert paths to absolute if they are not empty
// 	if config.VAE != "" {
// 		if !libs.PathExists(config.VAE) {
// 			return params, fmt.Errorf("error: config.VAE path does not exist")
// 		}
// 		params = append(params, "--vae", config.VAE)
// 	}

// 	if config.TAESD != "" {
// 		if !libs.PathExists(config.TAESD) {
// 			return params, fmt.Errorf("error: config.TAESD path does not exist")
// 		}
// 		params = append(params, "--taesd", config.TAESD)
// 	}
// 	if config.ControlNet != "" {
// 		if !libs.PathExists(config.ControlNet) {
// 			return params, fmt.Errorf("error: config.ControlNet path does not exist")
// 		}
// 		params = append(params, "--control-net", config.ControlNet)
// 	}
// 	// Continue for EmbeddingDir, StackedIDEmbDir, InputIDImagesDir, LoraModelDir, InitImg, ControlImage similarly...

// 	// Handle --type default (assuming f32 as a common default for weights)

// 	// Handle --clip-skip default logic (assuming SD2.x default if not specified)
// 	if config.ClipSkip == 0 {
// 		config.ClipSkip = -1 // SD2.x default
// 	}
// 	params = append(params, "--clip-skip", fmt.Sprintf("%d", config.ClipSkip))
// 	// Convert more paths to absolute if they are not empty
// 	if config.EmbeddingDir != "" {
// 		if !libs.PathExists(config.EmbeddingDir) {
// 			return params, fmt.Errorf("error: config.EmbeddingDir path does not exist")
// 		}
// 		params = append(params, "--embd-dir", config.EmbeddingDir)
// 	}
// 	if config.StackedIDEmbDir != "" {
// 		if !libs.PathExists(config.StackedIDEmbDir) {
// 			return params, fmt.Errorf("error: config.StackedIDEmbDir path does not exist")
// 		}
// 		params = append(params, "--stacked-id-embd-dir", config.StackedIDEmbDir)
// 	}
// 	if config.InputIDImagesDir != "" {
// 		if !libs.PathExists(config.InputIDImagesDir) {
// 			return params, fmt.Errorf("error: config.InputIDImagesDir path does not exist")
// 		}
// 		params = append(params, "--input-id-images-dir", config.InputIDImagesDir)
// 	}
// 	if config.LoraModelDir != "" {
// 		if !libs.PathExists(config.LoraModelDir) {
// 			return params, fmt.Errorf("error: config.LoraModelDir path does not exist")
// 		}
// 		params = append(params, "--lora-model-dir", config.LoraModelDir)
// 	}
// 	//path to the input image, required by img2img
// 	if config.InitImg != "" {
// 		if !libs.PathExists(config.InitImg) {
// 			return params, fmt.Errorf("error: config.InitImg path does not exist")
// 		}
// 		params = append(params, "-i", config.InitImg)
// 	}
// 	if config.ControlImage != "" {
// 		if !libs.PathExists(config.ControlImage) {
// 			return params, fmt.Errorf("error: config.ControlImage path does not exist")
// 		}
// 		params = append(params, "--control-image", config.ControlImage)
// 	}

// 	// Handle --normalize-input default
// 	if config.NormalizeInput {
// 		params = append(params, "--normalize-input")
// 	}

// 	// Handle --upscale-model default (assuming no upscale model by default)
// 	if config.UpscaleModel != "" {
// 		if !libs.PathExists(config.UpscaleModel) {
// 			return params, fmt.Errorf("error: config.UpscaleModel path does not exist")
// 		}
// 		params = append(params, "--upscale-model", config.UpscaleModel)
// 	}
// 	// Handle --upscale-repeats default
// 	if config.UpscaleRepeats == 0 {
// 		config.UpscaleRepeats = 1
// 	}
// 	params = append(params, "--upscale-repeats", fmt.Sprintf("%d", config.UpscaleRepeats))

// 	// Handle --vae-tiling default
// 	if config.VaeTiling {
// 		params = append(params, "--vae-tiling")
// 	}

// 	// Handle --control-net-cpu default
// 	if config.CtrlNetCPU {
// 		params = append(params, "--control-net-cpu")
// 	}

// 	// Handle --canny default
// 	if config.Canny {
// 		params = append(params, "--canny")
// 	}

// 	// Handle --color default
// 	if config.Color {
// 		params = append(params, "--color")
// 	}
// 	return params, nil
// }
