package search

import (
	"encoding/json"
	"testing"
)

func TestBing(t *testing.T) {
	engine := &Bing{Req: Req{Q: "godoos"}}
	s := engine.Search()
	marshal, _ := json.Marshal(s)
	t.Log(string(marshal))
}
