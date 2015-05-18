package ActiveHttpProxy

import (
	"net"
	"net/rpc"
	"testing"
	"time"
	//	"time"
)

func Test_ClusterChannle(t *testing.T) {
	cc := &ClusterChannle{"11123", nil}
	cc.Begin()
	time.Sleep(1 * time.Second)

	address, err := net.ResolveTCPAddr("tcp", "192.168.1.129:11123")
	if err != nil {
		panic(err)
	}

	conn, _ := net.DialTCP("tcp", nil, address)
	defer conn.Close()

	client := rpc.NewClient(conn)
	defer client.Close()

	var res string
	err = client.Call("ClusterChannle.FlushLog", "haha", &res)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(res)
	}
	//	client.
	//	err = client.Call("Reload", args, &reply)
	//	if err != nil {
	//		log.Fatal("arith error:", err)
	//	}
	//	log.Println(reply)

	time.Sleep(20 * time.Second)

	cc.Stop()
}
