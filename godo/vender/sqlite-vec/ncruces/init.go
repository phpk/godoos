// Package embeds a SQLite build that includes sqlite-vec into your application.
//
// Importing package embed initializes the [sqlite3.Binary] variable
// with an appropriate build of SQLite:
//
//	import _ "github.com/asg017/sqlite-vec-go-bindings/ncruces"
package embed

import (
	"bytes"
	_ "embed"
	"encoding/binary"

	"github.com/ncruces/go-sqlite3"
)

//go:embed sqlite3.wasm
var wasmBinary []byte

func init() {
	sqlite3.Binary = wasmBinary
}

func SerializeFloat32(vector []float32) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, vector)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
