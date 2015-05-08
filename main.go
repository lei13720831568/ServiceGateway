// ServiceGateway project main.go
package main

import (
	"encoding/json"
	"io/ioutil"
	//	"fmt"
	//"flag"
	//"log"
	//"time"
	//"io"
	"RollLoger"
	//	"os"
	"ActiveHttpProxy"
	//	"io"
)

func main() {

	RollLoger.Info("starting...")
	config, err := ioutil.ReadFile("Config.json")
	if err != nil {
		RollLoger.Error("read config.json failed ", err.Error())
		return
	}
	var appconfig AppConfigJson
	err = json.Unmarshal(config, &appconfig)
	if err != nil {
		RollLoger.Error("unmarshal config.json failed ", err.Error())
		return
	}

	r := &ActiveHttpProxy.ReaderFromDB{"driver={sql server};server=192.168.1.100;port=1433;uid=sa;pwd=654321;database=ADC3CoreDB"}
	arp := ActiveHttpProxy.NewArProxy("12001", r)
	arp.Start()
	//RollLoger.InitRollLogger(1024, "log", "./log")
	defer RollLoger.Close()
	//	//	logfile, err := RollLoger.NewRollFileWriter(1024, "testlog", "./log")

	//	//	if err != nil {
	//	//		fmt.Println("%srn", err.Error())
	//	//		os.Exit(-1)
	//	//	}
	//	//	log.SetOutput(logfile)

	//	//	//	defer logfile.Close()

	//	//	//	logger := log.New(logfile, "", log.Ldate|log.Ltime|log.Llongfile)
	//	//	//	logger.SetPrefix("error")
	//	for i := 0; i < 10000; i++ {
	//		RollLoger.Debug("safasdfa终于")
	//		//.Println("safasdfa终于")
	//	}

	//	logger.Println()

}
