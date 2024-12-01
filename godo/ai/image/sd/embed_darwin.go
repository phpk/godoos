// Copyright (c) seasonjs. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

//go:build darwin

package sd

import (
	_ "embed" // Needed for go:embed
)

//go:embed deps/darwin/libsd-abi.dylib
var libStableDiffusion []byte

var libName = "libstable-diffusion-*.dylib"

func getDl(gpu bool) []byte {
	if gpu {
		panic("Not support darwin. Push request is welcome.")
	}
	return libStableDiffusion
}
