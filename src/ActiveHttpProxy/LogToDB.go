package ActiveHttpProxy

import (
	log "RollLoger"
	_ "code.google.com/p/odbc"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	//	"io"
	//	"sync"
	"time"
)

func NewServiceGatewayLogger(connstr string) *ServiceGatewayLogger {
	ler := &ServiceGatewayLogger{}
	ler.dbconnstr = connstr
	ler.saveh = make(chan *ServiceGatewayLog, 9999)
	ler.ch = make(chan bool)
	ler.exit = make(chan bool)
	ler.save = make(chan bool)
	ler.saveExit = make(chan bool)
	ler.logs = &ArrayOfServiceGatewayLog{}
	ler.StartLogToDB()
	return ler
}

type ServiceGatewayLogger struct {
	dbconnstr string
	saveh     chan *ServiceGatewayLog
	ch        chan bool
	exit      chan bool
	logs      *ArrayOfServiceGatewayLog
	save      chan bool
	saveExit  chan bool
}

func (p *ServiceGatewayLogger) StartLogToDB() error {

	go func() {

		for {
			select {
			case <-p.ch:
				for i := 0; i < len(p.saveh); i++ {
					lo := <-p.saveh
					p.logs.Svs = append(p.logs.Svs, lo)
				}
				p.logs.saveToDB(p.dbconnstr)
				p.exit <- true
				return
			case <-p.save:
				p.logs.saveToDB(p.dbconnstr)
				p.saveExit <- true
			case lo := <-p.saveh:
				p.logs.Svs = append(p.logs.Svs, lo)
				if len(p.logs.Svs) > 200 {
					p.logs.saveToDB(p.dbconnstr)
				}
			case <-time.After(1 * time.Minute):
				p.logs.saveToDB(p.dbconnstr)
			}
		}

	}()

	return nil

}

func (p *ServiceGatewayLogger) StopLogToDB() {
	p.ch <- true
	<-p.exit
}

func (lo *ServiceGatewayLogger) AddLog(ar *ArRoute, requrl string, toUrl string, begin time.Time, end time.Time, rspStatus int, errorinfo string, host string) {
	p := &ServiceGatewayLog{}
	p.ReqUrl = requrl
	if ar != nil {
		p.PublishID = ar.PublishID
		rinfo, err := json.Marshal(ar)
		if err == nil {
			p.RouteInfo = string(rinfo)
		} else {
			p.RouteInfo = err.Error()
		}
	}

	p.RouteToUrl = toUrl
	p.LogDateTime = time.Now().Format("2006-01-02 15:04:05")
	p.BeginInvokeDate = begin.Format("2006-01-02 15:04:05")
	p.EndInvokeDate = end.Format("2006-01-02 15:04:05")
	p.UsedTime = end.Sub(begin).Seconds()
	p.RspStatus = rspStatus

	p.ErrorInfo = errorinfo
	lo.saveh <- p

}

func (lo *ServiceGatewayLogger) FlushLog() {
	lo.save <- true
	<-lo.saveExit
}

type ArrayOfServiceGatewayLog struct {
	XMLName xml.Name             `xml:"ArrayOfServiceGatewayLog"`
	Svs     []*ServiceGatewayLog `xml:"ServiceGatewayLog"`
}

func (p *ArrayOfServiceGatewayLog) saveToDB(dbstr string) error {

	if len(p.Svs) == 0 {
		return nil
	}
	xmlstr, err := xml.Marshal(p)
	if err != nil {
		log.Error("Marshal logs Error ", err.Error())
		return err
	}

	conn, err := sql.Open("odbc", dbstr)
	defer conn.Close()
	if err != nil {
		log.Error("Connecting Error ", err.Error())
		return err
	}

	_, err = conn.Exec("BatchInsertLog_ServiceGatewayLog ?", string(xmlstr))
	if err != nil {
		log.Error("conn exec Error", err.Error())
		return err
	} else {
		p.Svs = []*ServiceGatewayLog{}
	}

	p.Svs = []*ServiceGatewayLog{}

	return nil
}

type ServiceGatewayLog struct {
	ReqUrl          string
	RouteToUrl      string
	RspStatus       int
	RouteInfo       string
	LogDateTime     string
	BeginInvokeDate string
	EndInvokeDate   string
	UsedTime        float64
	PublishID       int
	Host            string
	ErrorInfo       string
}
