// ServiceGateway project main.go
package main

import (
	log "RollLoger"
	"flag"
	//	"fmt"
	"strings"
	//	"time"
	//	"net/http"
	_ "net/http/pprof"
	"winsvc/service"
)

var logger service.Logger

// Program structures.
//  Define Start and Stop methods.
type program struct {
	exit chan struct{}
}

func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		logger.Info("Running in terminal.")
	} else {
		logger.Info("Running under service manager.")
	}
	p.exit = make(chan struct{})

	go p.run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	// Any work in Stop should be quick, usually a few seconds at most.
	logger.Info("Stopping!")
	close(p.exit)
	return nil
}

// Service setup.
//   Define service config.
//   Create the service.
//   Setup the logger.
//   Handle service controls (optional).
//   Run the service.
func main() {
	defer log.Close()

	//性能监控
	//	go func() {
	//		http.ListenAndServe(":6060", nil)
	//	}()

	svcFlag := flag.String("service", "", "支持命令start, stop, restart, install,uninstall")
	flag.Parse()

	svcConfig := &service.Config{
		Name:        "ServiceGateway",
		DisplayName: "ServiceGateway",
		Description: "服务网关",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	errs := make(chan error, 5)

	logger, err = s.Logger(errs)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			err := <-errs
			if err != nil {
				log.Error(err.Error())
			}
		}
	}()

	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			log.Error("Valid actions: ", strings.Join(service.ControlAction[:], ","))
			log.Fatal(err)
		}
		return
	}

	//log.Info("begin run....")
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}

}
