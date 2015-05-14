package ActiveHttpProxy

import (
	_ "code.google.com/p/odbc"
	"database/sql"
	//	"encoding/xml"
	//	"fmt"
	"testing"
)

//type Servers struct {
//	XMLName xml.Name `xml:"servers"`
//	Version string   `xml:"version,attr"`
//	Svs     []server `xml:"server"`
//}

//type server struct {
//	ServerName string `xml:"serverName"`
//	ServerIP   string `xml:"serverIP"`
//}

func Test_ServiceDBLog(t *testing.T) {

	as := &ArrayOfServiceGatewayLog{}
	for i := 0; i < 3; i++ {
		l := &ServiceGatewayLog{}
		l.ReqUrl = "/123"
		l.RouteToUrl = "wert"
		l.RspStatus = 1

		l.RouteInfo = "sdf"
		l.LogDateTime = "2015-01-01"
		l.UsedTime = 20
		l.PublishID = 5
		l.Host = "sdfs"

		as.Svs = append(as.Svs, l)
	}
	dbstr := "driver={sql server};server=192.168.1.100;port=1433;uid=sa;pwd=654321;database=ADC3LogDB"
	conn, err := sql.Open("odbc", dbstr)
	if err != nil {
		t.Fatal(err)
	}
	err = as.saveToDB(conn)
	if err != nil {
		t.Fatal(err)
	}
}
