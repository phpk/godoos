//go:build windows && amd64

package deps

import (
	_ "embed"
)

//go:embed windows/amd64.zip
var embeddedZip []byte
