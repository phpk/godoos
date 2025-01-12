//go:build linux && arm64

package deps

import (
	_ "embed"
)

//go:embed linux/arm64.zip
var embeddedZip []byte
