# sqlite-vec-bindings-go

Go bindings for [`sqlite-vec`](https://github.com/asg017/sqlite-vec) project. In a separate repo to keep the original repo from getting too large.

There are two options when adding `sqlite-vec` to a Go project — A traditional CGO option, and a WASM-based option for the [ncruces/go-sqlite3](https://github.com/ncruces/go-sqlite3) project.

Both are available in this Go module, which can be installed with:

```bash
go get -u github.com/asg017/sqlite-vec-go-bindings
```

## CGO Bindings

For most SQLite Go libraries that use CGO, like [`mattn/go-sqlite3`](https://github.com/mattn/go-sqlite3), use the CGO portion of this Go module. It will compile the `sqlite-vec` libary from source and embed into your application.

```go
package main

import (
	"database/sql"
	"log"

	sqlite_vec "github.com/asg017/sqlite-vec-go-bindings/cgo"
	_ "github.com/mattn/go-sqlite3"
)

import "C"

func main() {
	sqlite_vec.Auto()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var sqliteVersion string
	var vecVersion string
	err = db.QueryRow("select sqlite_version(), vec_version()").Scan(&sqliteVersion, &vecVersion)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("sqlite_version=%s, vec_version=%s\n", sqliteVersion, vecVersion)
}
```

Use [`sqlite_vec.Auto()`](#TODO) before opening a connection to automatically register `sqlite-vec` functions. See [`simple-go-cgo/demo.go`](#TODO) for a larger example.

While this works with CGO SQLite/Go libraries like `mattn/go-sqlite3`, this will NOT work with other non-CGO library like [`modernc.org/sqlite`](https://pkg.go.dev/modernc.org/sqlite) or [ncruces/go-sqlite3](https://github.com/ncruces/go-sqlite3).

## `ncruces` WASM Bindings

If you are using the [ncruces/go-sqlite3](https://github.com/ncruces/go-sqlite3) library, then use the `ncruces` portion of this Go module.

```go
package main

import (
	_ "embed"
	"log"

	_ "github.com/asg017/sqlite-vec-go-bindings/ncruces"
	"github.com/ncruces/go-sqlite3"
)

func main() {
	db, err := sqlite3.Open(":memory:")
	if err != nil {
		log.Fatal(err)
	}

	stmt, _, err := db.Prepare(`SELECT sqlite_version(), vec_version()`)
	if err != nil {
		log.Fatal(err)
	}

	stmt.Step()

	log.Printf("sqlite_version=%s, vec_version=%s\n", stmt.ColumnText(0), stmt.ColumnText(1))
}

```

`"github.com/asg017/sqlite-vec-go-bindings/ncruces"` will automatically register a new SQLite WASM build that includes `sqlite-vec` functions by default. This replaces the `"github.com/ncruces/go-sqlite3/embed"` module in that project, so do NOT include them both.
