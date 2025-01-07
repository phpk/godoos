package search

import (
	"encoding/json"
	"testing"
)

func TestBaidu(t *testing.T) {
	engine := &Baidu{Req: Req{Q: "godoos"}}
	s := engine.Search()
	marshal, _ := json.Marshal(s)
	t.Log(string(marshal))
}
