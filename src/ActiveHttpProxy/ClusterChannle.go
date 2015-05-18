package ActiveHttpProxy

import (
	"fmt"
	"net/rpc"
	//	"net/rpc"
	log "RollLoger"
	"StoppableListener"
	"net"
)

type ClusterChannle struct {
	Port     string
	sl       *StoppableListener.StoppableListener
	Partners []string
}

func (p *ClusterChannle) FlushLog(req *string, res *string) error {
	fmt.Println("FlushLog")
	*res = "ok"
	return nil
}

func (p *ClusterChannle) Reload(req *string, res *string) error {
	fmt.Println("Reload")
	*res = "ok"
	return nil
}

func (p *ClusterChannle) CallPartner(cmd string) {

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

}

func (p *ClusterChannle) Begin() {

	go func() {
		rpc.Register(p)

		originalListener, err := net.Listen("tcp", ":"+p.Port)
		if err != nil {
			log.Fatal(err)
		}
		log.Debug("start clusterChannle tcp ", p.Port)
		p.sl, err = StoppableListener.New(originalListener)
		if err != nil {
			log.Fatal(err)
		}

		for {
			conn, err := p.sl.Accept()
			if err != nil {
				log.Error("rpc accept error ", err.Error())
				return
			}
			go rpc.ServeConn(conn)
		}
	}()

}

func (p *ClusterChannle) Stop() {
	p.sl.Stop()
}
