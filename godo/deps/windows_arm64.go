//go:build windows && arm64

package deps

import (
	_ "embed"
)

//go:embed windows/arm64.zip
var embeddedZip []byte
