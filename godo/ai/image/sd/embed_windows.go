// Copyright (c) seasonjs. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

//go:build windows && amd64

package sd

import (
	_ "embed"
	"encoding/json"
	"log"
	"os/exec"
	"strings"
	"syscall"

	"golang.org/x/sys/cpu"
)

//go:embed deps/windows/sd-abi_avx2.dll
var libStableDiffusionAvx2 []byte

//go:embed deps/windows/sd-abi_avx.dll
var libStableDiffusionAvx []byte

//go:embed deps/windows/sd-abi_avx512.dll
var libStableDiffusionAvx512 []byte

//go:embed deps/windows/sd-abi_cuda12.dll
var libStableDiffusionCuda12 []byte

//go:embed deps/windows/sd-abi_rocm5.5.dll
var libStableDiffusionRocm5 []byte

var libName = "stable-diffusion-*.dll"

func getDl(gpu bool) []byte {
	if gpu {
		info, err := NewGPU()
		if err != nil {
			log.Println(err)
		}
		cuda := info.Cuda()
		rocm := info.ROCm()

		if cuda.Available() {
			log.Print("get gpu info: ", cuda.Name)
			log.Println("Use GPU CUDA instead.")
			return libStableDiffusionCuda12
		}

		if rocm.Available() {
			log.Print("get gpu info: ", cuda.Name)
			log.Println("Use GPU ROCm instead.")
			return libStableDiffusionRocm5
		}

		log.Println("GPU not support, use CPU instead.")
	}

	if cpu.X86.HasAVX512 {
		log.Println("Use CPU AVX512 instead.")
		return libStableDiffusionAvx512
	}

	if cpu.X86.HasAVX2 {
		log.Println("Use CPU AVX2 instead.")
		return libStableDiffusionAvx2
	}

	if cpu.X86.HasAVX {
		log.Println("Use CPU AVX instead.")
		return libStableDiffusionAvx
	}

	panic("Automatic loading of dynamic library failed, please use `NewRwkvModel` method load manually. ")
}

type Driver struct {
	Name                 string `json:"Name"`
	AdapterCompatibility string `json:"AdapterCompatibility"`
	AdapterRAM           int64  `json:"AdapterRAM"`
}

func (d *Driver) Available() bool {
	return d.Name != "" && d.AdapterCompatibility != "" && d.AdapterRAM != 0
}

// GPU 类用于管理显卡信息
type GPU struct {
	drivers []Driver
	cuda    []Driver
	rocm    []Driver
}

func NewGPU() (*GPU, error) {
	cmd := exec.Command("powershell", `
        $graphicsCards = Get-WmiObject Win32_VideoController
        $graphicsArray = @()
		$graphicsEmpty = @{
			'Name'                 = ''
			'AdapterCompatibility' = ''
			'AdapterRAM'           = 0
		}
		$graphicsArray += $graphicsEmpty
        foreach ($card in $graphicsCards) {
            $graphicsInfo = @{
                'Name'                 = $card.Caption
                'AdapterCompatibility' = $card.VideoProcessor
                'AdapterRAM'           = $card.AdapterRAM
            }
            $graphicsArray += $graphicsInfo
        }
        $graphicsArray | ConvertTo-Json
    `)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
		//https://learn.microsoft.com/en-us/windows/win32/procthread/process-creation-flags
		CreationFlags: 0x08000000,
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var drivers []Driver
	err = json.Unmarshal(output, &drivers)
	if err != nil {
		return nil, err
	}

	cudaSupport := make([]Driver, 0)
	rocmSupport := make([]Driver, 0)

	for _, driver := range drivers {
		if strings.Contains(strings.ToUpper(driver.Name), "NVIDIA") {
			cudaSupport = append(cudaSupport, driver)
		} else if strings.Contains(strings.ToUpper(driver.Name), "AMD") {
			rocmSupport = append(rocmSupport, driver)
		}
	}

	return &GPU{
		drivers: drivers,
		cuda:    cudaSupport,
		rocm:    rocmSupport,
	}, nil
}

func (g *GPU) Cuda() *Driver {
	if len(g.cuda) > 0 {
		return &g.cuda[0]
	}
	return &Driver{}
}

func (g *GPU) ROCm() *Driver {
	if len(g.rocm) > 0 {
		return &g.rocm[0]
	}
	return &Driver{}
}

func (g *GPU) Info() []Driver {
	return g.drivers
}
