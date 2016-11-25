package main

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"

	"dana-tech.com/pg/task"
)

var index int64

type taske struct {
}

func (t *taske) Excute() {

	//tmp := index / (index - index)
	atomic.AddInt64(&index, 1)
	time.Sleep(time.Second)
}

func main() {

	manage := task.InitPreAcllocTaskerManager(1, 2)

	manage.Run()

	go func() {
		for {
			var tasker taske
			fmt.Println(manage.InputTask(&tasker, time.Millisecond*100))
			time.Sleep(time.Millisecond * 100)
		}
	}()

	go func() {
		for {
			fmt.Println(manage.Info(), index)
			fmt.Println(runtime.NumGoroutine())
			time.Sleep(time.Second * 2)
		}

	}()

	select {}
}
