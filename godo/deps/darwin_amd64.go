//go:build darwin && amd64

package deps

import (
	_ "embed"
)

//go:embed darwin/amd64.zip
var embeddedZip []byte
