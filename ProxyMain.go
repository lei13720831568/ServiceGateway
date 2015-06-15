package main

import (
	"ActiveHttpProxy"
	"RollLoger"
	"encoding/json"
	"io/ioutil"
	//	"os"
	"winsvc/osext"
)

func (p *program) run() error {
	RollLoger.Info("starting...")
	homedir, err := osext.ExecutableFolder()
	if err != nil {
		logger.Error(err)
		return err
	}
	config, err := ioutil.ReadFile(homedir + "/Config.json")
	if err != nil {
		RollLoger.Error("read config.json failed ", err.Error())
		return err
	}
	var appconfig AppConfigJson
	err = json.Unmarshal(config, &appconfig)
	if err != nil {
		RollLoger.Error("unmarshal config.json failed ", err.Error())
		return err
	}

	r := ActiveHttpProxy.NewReaderFromDB(appconfig.DBConnStr)
	dl := ActiveHttpProxy.NewServiceGatewayLogger(appconfig.LogDBConnStr)

	arp := ActiveHttpProxy.NewArProxy(appconfig.SelfNode.Port, r, dl, appconfig.SelfNode.Ip)
	arp.Start()

	for {
		select {
		case <-p.exit:
			arp.Stop()
			dl.StopLogToDB()
			return nil
		}
	}
	return nil
}

//func ProxyMain() {
//	RollLoger.Info("starting...")
//	config, err := ioutil.ReadFile("Config.json")
//	if err != nil {
//		RollLoger.Error("read config.json failed ", err.Error())
//		return
//	}
//	var appconfig AppConfigJson
//	err = json.Unmarshal(config, &appconfig)
//	if err != nil {
//		RollLoger.Error("unmarshal config.json failed ", err.Error())
//		return
//	}

//	r := &ActiveHttpProxy.ReaderFromDB{"driver={sql server};server=192.168.1.100;port=1433;uid=sa;pwd=654321;database=ADC3CoreDB"}
//	arp := ActiveHttpProxy.NewArProxy(appconfig.SelfNode.port, r)
//	arp.Start()

//}
