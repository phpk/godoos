package vec

// #cgo CFLAGS: -DSQLITE_CORE
// #include "sqlite-vec.h"
//
import "C"
import (
	"bytes"
	"encoding/binary"
)

// Once called, every future new SQLite3 connection created in this process
// will have the sqlite-vec extension loaded. It will persist until [Cancel] is
// called.
//
// Calls [sqlite3_auto_extension()] under the hood.
//
// [sqlite3_auto_extension()]: https://www.sqlite.org/c3ref/auto_extension.html
func Auto() {
	C.sqlite3_auto_extension( (*[0]byte) ((C.sqlite3_vec_init)) );
}

// "Cancels" any previous calls to [Auto]. Any new SQLite3 connections created
// will not have the sqlite-vec extension loaded.
//
// Calls sqlite3_cancel_auto_extension() under the hood.
func Cancel() {
	C.sqlite3_cancel_auto_extension( (*[0]byte) (C.sqlite3_vec_init) );
}

// Serializes a float32 list into a vector BLOB that sqlite-vec accepts.
func SerializeFloat32(vector []float32) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, vector)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

