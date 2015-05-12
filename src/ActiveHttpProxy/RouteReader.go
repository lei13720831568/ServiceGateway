package ActiveHttpProxy

import (
	log "RollLoger"
	_ "code.google.com/p/odbc"
	"database/sql"
	"encoding/json"
	//	"fmt"
	"io/ioutil"
	//	"os"
)

type RouteReader interface {
	Read() (*ArRouteLoad, error)
}

type ReaderFromDB struct {
	Dbconnstr string
}

func (p *ReaderFromDB) Read() (*ArRouteLoad, error) {
	connstr := p.Dbconnstr //"driver={sql server};server=192.168.1.100;port=1433;uid=sa;pwd=654321;database=ADC3CoreDB"
	conn, err := sql.Open("odbc", connstr)

	if err != nil {
		log.Error("Connecting Error ", err.Error())
		return nil, err
	}
	defer conn.Close()
	stmt, err := conn.Prepare("select * from vwService_info order by MatchMode desc,ReqUrl")
	if err != nil {
		log.Error("Prepare Query Error ", err.Error())
		return nil, err
	}
	defer stmt.Close()
	row, err := stmt.Query()
	if err != nil {
		log.Error("stmt Query Error", err.Error())
		return nil, err
	}
	arrl := &ArRouteLoad{}
	arrl.MaxVer = int64(0)
	arrl.Routes = *new([]*ArRoute)
	defer row.Close()
	for row.Next() {
		ar := &ArRoute{}
		if err := row.Scan(&ar.PublishID, &ar.Name, &ar.ReqUrl, &ar.ProxyToUrl, &ar.IpList, &ar.SecretType, &ar.Encrypt, &ar.MaxConnects, &ar.TimeOut, &ar.Ver, &ar.Status, &ar.SecKey, &ar.MatchMode); err == nil {
			arrl.Routes = append(arrl.Routes, ar)
			arrl.MaxVer = ar.Ver
		} else {
			log.Error("read fields Error", err.Error())
			return nil, err
		}
	}
	return arrl, nil
}

type ReaderFromJsonFile struct {
	FilePath string
}

func (p *ReaderFromJsonFile) Read() (*ArRouteLoad, error) {

	fc, err := ioutil.ReadFile(p.FilePath)
	if err != nil {
		return nil, err
	}

	var arrl ArRouteLoad
	err = json.Unmarshal(fc, &arrl)
	if err == nil {
		return &arrl, nil
	} else {
		return nil, err
	}
}
