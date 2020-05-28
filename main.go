package main

import (
	"fmt"
	"time"

	"github.com/devfabric/diskqueue/config"
	diskqueue "github.com/devfabric/diskqueue/go-diskqueue"
)

func NewTestLogger() diskqueue.AppLogFunc {
	return func(lvl diskqueue.LogLevel, f string, args ...interface{}) {
		// tbl.Log(fmt.Sprintf(lvl.String()+": "+f, args...))
		fmt.Println(fmt.Sprintf(lvl.String()+": "+f, args...))
	}
}

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

	err = dq.Put([]byte("test"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	msgOut := <-dq.ReadChan()
	fmt.Println(string(msgOut))
}
