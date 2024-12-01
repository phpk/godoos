// Copyright (c) seasonjs. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package sd

import (
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"path/filepath"
)

type OutputsImageType string

const (
	PNG  OutputsImageType = "PNG"
	JPEG                  = "JPEG"
)

// Options 定义了模型配置的结构体，用于配置模型的各种参数和选项。
type Options struct {
	// VaePath 是变分自编码器(VAE)模型的路径。
	VaePath string
	// TaesdPath 是自编码器(TAE)模型的路径。
	TaesdPath string
	// LoraModelDir 是Lora模型目录的路径，Lora模型用于特定类型的模型训练和推理。
	LoraModelDir string
	// VaeDecodeOnly 指定VAE模型是否仅用于解码操作。
	VaeDecodeOnly bool
	// VaeTiling 指定是否在VAE模型中使用分块(tiling)处理。
	VaeTiling bool
	// FreeParamsImmediately 指定是否在不再需要参数时立即释放参数占用的内存。
	FreeParamsImmediately bool
	// Threads 定义了并行处理的线程数量。
	Threads int
	// Wtype 是权重类型，决定了模型权重的处理方式。
	Wtype WType
	// RngType 是随机数生成器类型，用于指定模型中随机数的生成方式。
	RngType RNGType
	// Schedule 定义了模型训练或推理的时间表安排。
	Schedule Schedule
	// GpuEnable 指定是否启用GPU加速计算。
	GpuEnable bool
}

// FullParams 包含了所有与生成图像相关的参数配置。
// 这些参数用于指导图像生成模型的条件和过程，以达到预期的图像生成效果。
type FullParams struct {
	// NegativePrompt 指定模型应避免的图像内容描述。
	// 通过提供负面提示，用户可以指导模型远离不期望的生成结果。
	NegativePrompt string

	// ClipSkip 表示在CLIP模型中跳过的层数。
	// 这可以影响生成图像的细节和质量。
	ClipSkip int

	// CfgScale 是配置模型的尺度因子。
	// 它决定了模型在生成图像时对提示的遵循程度。
	CfgScale float32

	// Width 和 Height 分别定义了生成图像的宽度和高度。
	// 这些参数允许用户自定义输出图像的尺寸。
	Width  int
	Height int

	// SampleMethod 指定用于生成图像的采样方法。
	// 不同的采样方法可以影响图像生成的速度和质量。
	SampleMethod SampleMethod

	// SampleSteps 是执行采样方法的步数。
	// 较多的步数可能产生更高质量的图像，但需要更长的生成时间。
	SampleSteps int

	// Strength 定义了在图像生成过程中，新信息引入的强度。
	// 该参数可以影响生成图像的多样性和细节。
	Strength float32

	// Seed 是用于随机数生成器的种子。
	// 相同的种子和参数将产生相同的图像，确保了结果的可重复性。
	Seed int64

	// BatchCount 表示一次性生成的图像数量。
	// 这允许用户一次性获得多个不同的生成结果。
	BatchCount int

	// OutputsImageType 指定输出图像的类型。
	// 这可以是不同的文件格式或质量设置。
	OutputsImageType OutputsImageType
}

var DefaultOptions = Options{
	Threads:               -1, // auto
	VaeDecodeOnly:         true,
	FreeParamsImmediately: true,
	RngType:               CUDA_RNG,
	Wtype:                 F32,
	Schedule:              DEFAULT,
}

var DefaultFullParams = FullParams{
	NegativePrompt:   "out of frame, lowers, text, error, cropped, worst quality, low quality, jpeg artifacts, ugly, duplicate, morbid, mutilated, out of frame, extra fingers, mutated hands, poorly drawn hands, poorly drawn face, mutation, deformed, blurry, dehydrated, bad anatomy, bad proportions, extra limbs, cloned face, disfigured, gross proportions, malformed limbs, missing arms, missing legs, extra arms, extra legs, fused fingers, too many fingers, long neck, username, watermark, signature",
	CfgScale:         7.0,
	Width:            512,
	Height:           512,
	SampleMethod:     EULER_A,
	SampleSteps:      20,
	Strength:         0.4,
	Seed:             42,
	BatchCount:       1,
	OutputsImageType: PNG,
}

type Model struct {
	ctx                *CStableDiffusionCtx
	options            *Options
	csd                CStableDiffusion
	isAutoLoad         bool
	dylibPath          string
	diffusionModelPath string
	esrganPath         string
	upscalerCtx        *CUpScalerCtx
}

func NewAutoModel(options Options) (*Model, error) {
	file, err := dumpSDLibrary(options.GpuEnable)
	if err != nil {
		return nil, err
	}

	if options.GpuEnable {
		log.Printf("If you want to try offload your model to the GPU. " +
			"Please confirm the size of your GPU memory to prevent memory overflow.")
	}
	dylibPath := file.Name()
	model, err := NewModel(dylibPath, options)
	if err != nil {
		return nil, err
	}
	model.isAutoLoad = true
	return model, nil
}

func NewModel(dylibPath string, options Options) (*Model, error) {
	csd, err := NewCStableDiffusion(dylibPath)
	if err != nil {
		return nil, err
	}

	return &Model{
		dylibPath: dylibPath,
		options:   &options,
		csd:       csd,
	}, nil
}

func (sd *Model) LoadFromFile(path string) error {
	if sd.ctx != nil {
		sd.csd.FreeCtx(sd.ctx)
		sd.ctx = nil
		log.Printf("model already loaded, free old model")
	}

	_, err := os.Stat(path)
	if err != nil {
		return errors.New("the system cannot find the model file specified")
	}

	if !filepath.IsAbs(path) {
		sd.diffusionModelPath, err = filepath.Abs(path)
		if err != nil {
			return err
		}
	} else {
		sd.diffusionModelPath = path
	}

	ctx := sd.csd.NewCtx(path,
		sd.options.VaePath,
		sd.options.TaesdPath,
		sd.options.LoraModelDir,
		sd.options.VaeDecodeOnly,
		sd.options.VaeTiling,
		sd.options.FreeParamsImmediately,
		sd.options.Threads,
		sd.options.Wtype,
		sd.options.RngType,
		sd.options.Schedule)
	sd.ctx = ctx
	return nil
}

func (sd *Model) SetOptions(options Options) {
	if sd.ctx != nil {
		sd.csd.FreeCtx(sd.ctx)
		sd.ctx = nil
		log.Printf("model already loaded, free old model and set new options")
	}
	sd.options = &options
	ctx := sd.csd.NewCtx(
		sd.diffusionModelPath,
		sd.options.VaePath,
		sd.options.TaesdPath,
		sd.options.LoraModelDir,
		sd.options.VaeDecodeOnly,
		sd.options.VaeTiling,
		sd.options.FreeParamsImmediately,
		sd.options.Threads,
		sd.options.Wtype,
		sd.options.RngType,
		sd.options.Schedule)
	sd.ctx = ctx
}

func (sd *Model) Predict(prompt string, params FullParams, writer []io.Writer) error {
	if len(writer) != params.BatchCount {
		return errors.New("writer count not match batch count")
	}
	if sd.ctx == nil {
		return errors.New("model not loaded")
	}

	if params.Width%8 != 0 || params.Height%8 != 0 {
		return errors.New("width and height must be multiples of 8")
	}

	images := sd.csd.PredictImage(
		sd.ctx,
		prompt,
		params.NegativePrompt,
		params.ClipSkip,
		params.CfgScale,
		params.Width,
		params.Height,
		params.SampleMethod,
		params.SampleSteps,
		params.Seed,
		params.BatchCount,
	)

	if images == nil || len(images) != params.BatchCount {
		return errors.New("predict failed")
	}

	for i, img := range images {
		outputsImage := bytesToImage(img.Data, int(img.Width), int(img.Height))

		err := imageToWriter(outputsImage, params.OutputsImageType, writer[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (sd *Model) ImagePredict(reader io.Reader, prompt string, params FullParams, writer []io.Writer) error {

	if len(writer) != params.BatchCount {
		return errors.New("writer count not match batch count")
	}

	if sd.ctx == nil {
		return errors.New("model not loaded")
	}

	decode, _, err := image.Decode(reader)
	if err != nil {
		return err
	}
	initImage := imageToBytes(decode)
	images := sd.csd.ImagePredictImage(
		sd.ctx,
		initImage,
		prompt,
		params.NegativePrompt,
		params.ClipSkip,
		params.CfgScale,
		params.Width,
		params.Height,
		params.SampleMethod,
		params.SampleSteps,
		params.Strength,
		params.Seed,
		params.BatchCount,
	)
	for i, img := range images {
		outputsImage := bytesToImage(img.Data, int(img.Width), int(img.Height))
		err = imageToWriter(outputsImage, params.OutputsImageType, writer[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (sd *Model) UpscaleImage(reader io.Reader, esrganPath string, upscaleFactor uint32, writer io.Writer) error {
	if sd.upscalerCtx == nil {
		sd.esrganPath = esrganPath
		sd.upscalerCtx = sd.csd.NewUpscalerCtx(esrganPath, sd.options.Threads, sd.options.Wtype)
	}

	if sd.esrganPath != esrganPath {
		if sd.upscalerCtx != nil {
			sd.csd.FreeUpscalerCtx(sd.upscalerCtx)
		}
		sd.upscalerCtx = sd.csd.NewUpscalerCtx(esrganPath, sd.options.Threads, sd.options.Wtype)
	}

	decode, _, err := image.Decode(reader)
	if err != nil {
		return err
	}
	initImage := imageToBytes(decode)
	img := sd.csd.UpscaleImage(sd.upscalerCtx, initImage, upscaleFactor)
	outputsImage := bytesToImage(img.Data, int(img.Width), int(img.Height))
	err = imageToWriter(outputsImage, PNG, writer)
	return err
}

func (sd *Model) SetLogCallback(cb CLogCallback) {
	sd.csd.SetLogCallBack(cb)
}

func (sd *Model) Close() error {
	if sd.ctx != nil {
		sd.csd.FreeCtx(sd.ctx)
		sd.ctx = nil
	}

	if sd.upscalerCtx != nil {
		sd.csd.FreeUpscalerCtx(sd.upscalerCtx)
		sd.upscalerCtx = nil
	}

	if sd.csd != nil {
		err := sd.csd.Close()
		if err != nil {
			return err
		}
	}

	if sd.isAutoLoad {
		err := os.Remove(sd.dylibPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func imageToBytes(decode image.Image) Image {
	bounds := decode.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y
	size := width * height * 3
	bytesImg := make([]byte, size)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			idx := (y*width + x) * 3
			r, g, b, _ := decode.At(x, y).RGBA()
			bytesImg[idx] = byte(r >> 8)
			bytesImg[idx+1] = byte(g >> 8)
			bytesImg[idx+2] = byte(b >> 8)
		}
	}
	return Image{
		Width:   uint32(width),
		Height:  uint32(height),
		Data:    bytesImg,
		Channel: 3,
	}
}

func bytesToImage(byteData []byte, width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			idx := (y*width + x) * 3
			img.Set(x, y, color.RGBA{
				R: byteData[idx],
				G: byteData[idx+1],
				B: byteData[idx+2],
				A: 255,
			})
		}
	}
	return img
}

func imageToWriter(image image.Image, imageType OutputsImageType, writer io.Writer) error {
	switch imageType {
	case PNG:
		err := png.Encode(writer, image)
		if err != nil {
			return err
		}
	case JPEG:
		err := jpeg.Encode(writer, image, &jpeg.Options{Quality: 100})
		if err != nil {
			return err
		}
	default:
		return errors.New("unknown image type")
	}
	return nil
}

//func chunkBytes(data []byte, chunks int) [][]byte {
//	length := len(data)
//	chunkSize := (length + chunks - 1) / chunks
//	result := make([][]byte, chunks)
//
//	for i := 0; i < chunks; i++ {
//		start := i * chunkSize
//		end := (i + 1) * chunkSize
//		if end > length {
//			end = length
//		}
//		result[i] = data[start:end:end]
//	}
//
//	return result
//}
