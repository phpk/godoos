//go:build linux && amd64

package deps

import (
	_ "embed"
)

//go:embed linux/amd64.zip
var embeddedZip []byte
