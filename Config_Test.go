package main

import (
	"encoding/json"
	"testing"
)

func Test_ConfigJson(t *testing.T) {
	a := &AppConfigJson{"db1", "db2", &NodeInfo{"192.168.1.100", "12345"}}
	strbyte, err := json.Marshal(a)
	if err == nil {
		t.Log(string(strbyte))
	}
}
