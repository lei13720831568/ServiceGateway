package main

import (
//"encoding/json"
//"testing"
)

type AppConfigJson struct {
	DBConnStr string
	SelfNode  *NodeInfo
}

type NodeInfo struct {
	Ip   string
	Port string
}
