package WorkPool

import (
	//	log "RollLoger"
	//	"encoding/json"
	//	"net/http"
	"errors"
	//	"strconv"
	"sync"
	"time"
)

type Worker interface {
	PHandle() error
}

type WPool struct {
	wg        *sync.WaitGroup
	workQueue chan bool
	PoolName  string
	mu        *sync.Mutex
	Maxworks  int
}

func NewWorkPool(name string) *WPool {
	wp := &WPool{}
	wp.Maxworks = 0
	wp.PoolName = name
	wp.wg = &sync.WaitGroup{}
	wp.workQueue = make(chan bool, 9999)
	wp.mu = &sync.Mutex{}

	return wp
}

func (pw *WPool) work(w Worker) error {

	defer func() {
		pw.workQueue <- true
	}()

	return w.PHandle()

}

//动态调整worker数量到 newmax，newmax取值访问0-9999，如果是减少操作则会产生阻塞，直到释放足够多的的worker
func (pw *WPool) SetMax(newmax int) error {
	pw.mu.Lock()
	defer pw.mu.Unlock()
	if pw.Maxworks != newmax {
		if (newmax) > 9999 {
			return errors.New("over MakWorker max=9999")
		}

		if (newmax) < 0 {
			return errors.New("over MakWorker min=0")
		}
		ws := newmax - pw.Maxworks //求出差异量

		//调整连接池
		if ws >= 0 {
			for i := 0; i < ws; i++ {
				pw.workQueue <- true
			}
		} else {
			for i := 0; i < ws; i++ {
				<-pw.workQueue
			}
		}
		pw.Maxworks = newmax
	}
	return nil
}

//使用当前协程进行work
func (pw *WPool) PutWork(w Worker, waitTime time.Duration) error {

	select {
	case <-pw.workQueue:
		return pw.work(w)
	case <-time.After(waitTime):
		return errors.New("put waiting TimeOut")
	}

	//		select {
	//		case pw.workQueue <- w:
	//			return nil
	//		case <-time.After(pw.WaitTimeOut):
	//			return errors.New("put waiting TimeOut")
	//		}
}

//使用新协程work，参数c是回调的通信通道
func (pw *WPool) AsyncPutWork(w Worker, waitTime time.Duration, c chan<- error) {

	go func() {
		err := pw.PutWork(w, waitTime)
		if c != nil {
			c <- err
		}

	}()

}

func (pw *WPool) GetStat() int {
	pw.mu.Lock()
	defer pw.mu.Unlock()
	return pw.Maxworks - len(pw.workQueue)
	//return len(pw.workQueue)
}
