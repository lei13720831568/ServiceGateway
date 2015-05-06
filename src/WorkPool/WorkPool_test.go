package WorkPool

import (
	"fmt"
	"math/rand"
	//	"runtime"
	"strconv"
	"testing"
	"time"
)

type test_dosoming struct {
	name string
}

func (d *test_dosoming) PHandle() error {
	time.Sleep(3 * time.Second)
	fmt.Println("do ", d.name)
	return nil

}

func goput(p *WPool, i int) {
	d := &test_dosoming{"t" + strconv.Itoa(i)}
	err := p.PutWork(d, 1*time.Second)
	if err != nil {
		fmt.Println("put error :", err)
	}
}

func randPut(p *WPool) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {

		t := r.Intn(5)
		<-time.After(time.Duration(t) * time.Second)
		goput(p, t)

	}
}

func Test_WorkPool(t *testing.T) {
	pool := NewWorkPool("工作池1")
	pool.SetMax(2)

	go func() {
		for {
			<-time.After(1 * time.Second)
			fmt.Println("stat ", pool.GetStat())
		}
	}()

	for i := 1; i < 14; i++ {
		go randPut(pool)
	}

	time.Sleep(10 * time.Second)

	pool.SetMax(20)

	time.Sleep(10 * time.Second)

	fmt.Println("begin stop ")

	pool.SetMax(0)

	fmt.Println("end stop ")

	time.Sleep(20 * time.Second)
}
