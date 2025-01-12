//go:build darwin && arm64

package deps

import (
	_ "embed"
)

//go:embed darwin/arm64.zip
var embeddedZip []byte
