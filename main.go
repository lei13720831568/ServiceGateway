// ServiceGateway project main.go
package main

import (
	//	"fmt"
	//"flag"
	//"log"
	//"time"
	//"io"
	"RollLoger"
	//	"os"
)

func main() {

	//RollLoger.InitRollLogger(1024, "log", "./log")
	defer RollLoger.Close()
	//	logfile, err := RollLoger.NewRollFileWriter(1024, "testlog", "./log")

	//	if err != nil {
	//		fmt.Println("%srn", err.Error())
	//		os.Exit(-1)
	//	}
	//	log.SetOutput(logfile)

	//	//	defer logfile.Close()

	//	//	logger := log.New(logfile, "", log.Ldate|log.Ltime|log.Llongfile)
	//	//	logger.SetPrefix("error")
	for i := 0; i < 10000; i++ {
		RollLoger.Debug("safasdfa终于")
		//.Println("safasdfa终于")
	}

	//	logger.Println()

}
