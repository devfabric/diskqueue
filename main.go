package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/devfabric/diskqueue/config"
	diskqueue "github.com/devfabric/diskqueue/go-diskqueue"
)

func NewTestLogger() diskqueue.AppLogFunc {
	return func(lvl diskqueue.LogLevel, f string, args ...interface{}) {
		// tbl.Log(fmt.Sprintf(lvl.String()+": "+f, args...))
		//fmt.Println(fmt.Sprintf(lvl.String()+": "+f, args...))
	}
}

const (
	GO_READ_THREAD  int = 100
	GO_WRITE_THREAD int = 1000
)

func main() {
	queueConfig, err := config.LoadQueueConfig("./")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(queueConfig)

	l := NewTestLogger()
	dq := diskqueue.New(queueConfig.Name, queueConfig.DataPath, queueConfig.MaxBytesPerFile,
		queueConfig.MinMsgSize, queueConfig.MaxMsgSize, queueConfig.SyncEvery,
		time.Duration(queueConfig.SyncTimeout)*time.Second, l)

	defer dq.Close()
	fmt.Println(dq.Depth())

	var synReadWait sync.WaitGroup
	for n := 0; n < GO_READ_THREAD; n++ {
		synReadWait.Add(1)
		go func(seq int) {
			for {
				msgOut := <-dq.ReadChan()
				fmt.Println("read seq:", seq, "--->", string(msgOut))
			}
		}(n)
	}

	var synWait sync.WaitGroup
	for n := 0; n < GO_READ_THREAD; n++ {
		synWait.Add(1)
		go func(seq int) {
			i := 0
			for {
				i++
				err = dq.Put([]byte(fmt.Sprintf("test-%d:%d", seq, i)))
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				fmt.Printf("Depth-%d:%d\n", seq, dq.Depth())
			}
		}(n)
	}

	synReadWait.Wait()
	synWait.Wait()
	select {}

}
