package ActiveHttpProxy //Reverse Porxy

import (
	"encoding/json"
	"fmt"
	"time"
)

type Watcher struct {
	ch chan bool
}

func NewWatcher() *Watcher {
	r := &Watcher{}
	r.ch = make(chan bool)
	return r
}

func (w *Watcher) StartWatch(p *ArProxy) {
	for {
		select {
		case <-w.ch:
			break
		case <-time.After(1 * time.Second):
			j, err := w.GetRouteJson()
			if err == nil {
				fmt.Println("begin reload")
				p.service_routes.RoadRoute(j)
			}
		}
	}
}

func (w *Watcher) StopWatch() {
	w.ch <- true
}

func (w *Watcher) GetRouteJson() (string, error) {

	arm := &ArRouteLoad{}
	arm.Routes = []*ArRoute{}

	max := int64(0)
	for i := 1; i < 14; i++ {
		ar := &ArRoute{}
		ar.ReqUrl = fmt.Sprintf("//testservice//%d.html", i)
		ar.Encrypt = "nil"
		ar.IpList = "nil"
		ar.MaxConnects = 10
		ar.Name = fmt.Sprintf("service%d", i)
		ar.ProxyToUrl = fmt.Sprintf("//destservice//%d.html", i)
		ar.PublishID = i
		ar.SecretType = "nil"
		ar.TimeOut = int64(10 * time.Second)
		ar.Ver = int64(i)
		if ar.Ver > max {
			max = ar.Ver
		}
		arm.Routes = append(arm.Routes, ar)
	}

	arm.MaxVer = max
	r, err := json.Marshal(arm)
	if err != nil {
		return "", err
	}

	return string(r), nil
}
