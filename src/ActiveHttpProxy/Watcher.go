package ActiveHttpProxy //Reverse Porxy

import (
	//	log "RollLoger"
	//_ "code.google.com/p/odbc"
	//"database/sql"
	//	"encoding/json"
	//	"fmt"
	"time"
)

type Watcher struct {
	ch     chan bool
	reader RouteReader
}

func NewWatcher(wr RouteReader) *Watcher {
	r := &Watcher{}
	r.ch = make(chan bool)
	r.reader = wr
	return r
}

func (w *Watcher) StartWatch(p *ArProxy) {
	for {
		select {
		case <-w.ch:
			break
		case <-time.After(10 * time.Minute):
			j, err := w.reader.Read()
			if err == nil {
				p.service_routes.RoadRoute(j)
			}
		}
	}
}

func (w *Watcher) StopWatch() {

	w.ch <- true
}

//func (w *Watcher) ReadDBRoute() (*ArRouteLoad, error) {
//	connstr := "driver={sql server};server=192.168.1.100;port=1433;uid=sa;pwd=654321;database=ADC3CoreDB"
//	conn, err := sql.Open("odbc", connstr)

//	if err != nil {
//		log.Error("Connecting Error ", err.Error())
//		return nil, err
//	}
//	defer conn.Close()
//	stmt, err := conn.Prepare("select * from vwService_info")
//	if err != nil {
//		log.Error("Prepare Query Error ", err.Error())
//		return nil, err
//	}
//	defer stmt.Close()
//	row, err := stmt.Query()
//	if err != nil {
//		log.Error("stmt Query Error", err.Error())
//		return nil, err
//	}
//	arrl := &ArRouteLoad{}
//	arrl.MaxVer = int64(0)
//	arrl.Routes = *new([]*ArRoute)
//	defer row.Close()
//	for row.Next() {
//		ar := &ArRoute{}
//		if err := row.Scan(&ar.PublishID, &ar.Name, &ar.ReqUrl, &ar.ProxyToUrl, &ar.IpList, &ar.SecretType, &ar.Encrypt, &ar.MaxConnects, &ar.TimeOut, &ar.Ver, &ar.Status); err == nil {
//			arrl.Routes = append(arrl.Routes, ar)
//		} else {
//			log.Error("read fields Error", err.Error())
//			return nil, err
//		}
//	}
//	return arrl, nil
//}

//func (w *Watcher) GetRouteJson() (string, error) {

//	arm := &ArRouteLoad{}
//	arm.Routes = []*ArRoute{}

//	max := int64(0)
//	for i := 1; i < 14; i++ {
//		ar := &ArRoute{}
//		ar.ReqUrl = fmt.Sprintf("/testservice/%d.html", i)
//		ar.Encrypt = "nil"
//		ar.IpList = "nil"
//		ar.MaxConnects = 10
//		ar.Name = fmt.Sprintf("service%d", i)
//		ar.ProxyToUrl = fmt.Sprintf("http://192.168.1.129/test/test%d.html", i)
//		ar.PublishID = i
//		ar.SecretType = "nil"
//		ar.TimeOut = int64(10 * time.Second)
//		ar.Ver = int64(i)
//		if ar.Ver > max {
//			max = ar.Ver
//		}
//		arm.Routes = append(arm.Routes, ar)
//	}

//	arm.MaxVer = max
//	r, err := json.Marshal(arm)
//	if err != nil {
//		return "", err
//	}

//	return string(r), nil
//}
