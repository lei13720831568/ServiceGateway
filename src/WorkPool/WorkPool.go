package WorkPool

import (
	log "RollLoger"
	//	"encoding/json"
	//	"net/http"
	"errors"
	"strconv"
	"sync"
	"time"
)

type Worker interface {
	PHandle()
}

type WPool struct {
	wg          *sync.WaitGroup
	workQueue   chan Worker
	PoolName    string
	status      bool
	Maxworks    int
	WaitTimeOut time.Duration //在队列满载情况下等待的时间
}

func NewWorkPool(name string, max int, waitTimeout time.Duration) *WPool {
	wp := &WPool{}
	wp.Maxworks = max
	wp.PoolName = name
	wp.status = false
	wp.wg = &sync.WaitGroup{}
	wp.WaitTimeOut = waitTimeout
	return wp
}

func (pw *WPool) Start() error {
	if pw.status == false {
		pw.workQueue = make(chan Worker, pw.Maxworks)

		for i := 0; i < pw.Maxworks; i++ {
			go pw.work()
		}
		pw.status = true
		log.Info("start works name:", pw.PoolName, " max:", strconv.Itoa(pw.Maxworks))

		return nil
	} else {
		return errors.New("serviceName works already started")
	}

}

func (pw *WPool) work() {
	pw.wg.Add(1)
	defer pw.wg.Done()
	for {
		r, cl := <-pw.workQueue
		if r != nil {
			r.PHandle() //处理
		}
		if cl == false {
			return //退出事件
		}
	}
	log.Info("close worker")
}

func (pw *WPool) Stop() {
	log.Info("begin close works name:", pw.PoolName)
	if pw.status != false {
		pw.status = false
		close(pw.workQueue)
		pw.wg.Wait()
		log.Info("close works name:", pw.PoolName)
	}
}

func (pw *WPool) Put(w Worker) error {
	if pw.status == true {
		select {
		case pw.workQueue <- w:
			return nil
		case <-time.After(pw.WaitTimeOut):
			return errors.New("put waiting TimeOut")
		}

	}
	return errors.New("put Pool not start")
}

func (pw *WPool) GetStat() int {
	if pw.status {
		return len(pw.workQueue)
	}
	return 0
}
