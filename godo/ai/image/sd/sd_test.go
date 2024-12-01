package sd

import (
	"io"
	"os"
	"testing"
)

func TestNewStableDiffusionAutoModelPredict(t *testing.T) {
	options := DefaultOptions
	t.Log(options)
	model, err := NewAutoModel(options)
	if err != nil {
		t.Error(err)
		return
	}
	defer model.Close()
	model.SetLogCallback(func(level LogLevel, msg string) {
		t.Log(msg)
	})
	err = model.LoadFromFile("./models/miniSD.ckpt")
	if err != nil {
		t.Error(err)
		return
	}
	var writers []io.Writer
	filenames := []string{
		"./assets/love_cat2.png",
	}
	for _, filename := range filenames {
		file, err := os.Create(filename)
		if err != nil {
			t.Error(err)
			return
		}
		defer file.Close()
		writers = append(writers, file)
	}

	params := DefaultFullParams
	params.BatchCount = 1
	params.Width = 256
	params.Height = 256
	params.NegativePrompt = ""
	err = model.Predict("british short hair cat, high quality", params, writers)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestModel_ROCm(t *testing.T) {
	options := DefaultOptions
	options.GpuEnable = true
	t.Log(options)
	model, err := NewAutoModel(options)
	if err != nil {
		t.Error(err)
		return
	}
	defer model.Close()
	model.SetLogCallback(func(level LogLevel, msg string) {
		t.Log(msg)
	})
	err = model.LoadFromFile("./models/miniSD.ckpt")
	if err != nil {
		t.Error(err)
		return
	}
	var writers []io.Writer
	filenames := []string{
		"./assets/love_cat2.png",
	}
	for _, filename := range filenames {
		file, err := os.Create(filename)
		if err != nil {
			t.Error(err)
			return
		}
		defer file.Close()
		writers = append(writers, file)
	}

	params := DefaultFullParams
	params.BatchCount = 1
	params.Width = 256
	params.Height = 256
	params.NegativePrompt = ""
	err = model.Predict("british short hair cat, high quality", params, writers)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestNewStableDiffusionAutoModelImagePredict(t *testing.T) {
	options := DefaultOptions
	options.VaeDecodeOnly = false
	t.Log(options)
	model, err := NewAutoModel(options)
	if err != nil {
		t.Error(err)
		return
	}
	defer model.Close()
	model.SetLogCallback(func(level LogLevel, msg string) {
		t.Log(msg)
	})
	err = model.LoadFromFile("./models/mysafetensors")
	if err != nil {
		t.Error(err)
		return
	}
	inFile, err := os.Open("./assets/love_cat0.png")
	if err != nil {
		t.Error(err)
		return
	}
	defer inFile.Close()

	var writers []io.Writer
	filenames := []string{
		"./assets/love_cat0_m.png",
		//"./assets/love_cat1_m.png",
		//"./assets/love_cat5.png",
		//"./assets/love_cat6.png"
	}
	for _, filename := range filenames {
		file, err := os.Create(filename)
		if err != nil {
			t.Error(err)
			return
		}
		defer file.Close()
		writers = append(writers, file)
	}
	params := DefaultFullParams
	params.BatchCount = 1
	params.Width = 256
	params.Height = 256
	params.NegativePrompt = ""
	err = model.ImagePredict(inFile, "dogs", params, writers)
	if err != nil {
		t.Error(err)
		return
	}
}
