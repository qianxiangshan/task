//author qian
//time 2016-11-16
//提供一个异步的task管理，固定执行task的goroutine数量。限制同时执行的任务个数,提供失败通知。
package task

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"

	"dana-tech.com/wbw/logs"
)

type RealtimeAcllocManage struct {
	TaskerManage
}

//这个函数目前设计成，除非程序退出，否则一直存活的模式，
//目前加入task退出模式，保证task缓存中执行完后退出
func (taskmanage *RealtimeAcllocManage) taskexcute() {
	defer func() {
		panicerr := recover()
		//如果是panic，则重启一个goroutine，防止goroutine没了，阻塞成哥程序
		if panicerr != nil {
			go taskmanage.taskexcute()
		} else {
			//会退出的，收到退出信号。
		}
	}()

	var stask Tasker
	for {
		select {
		case stask = <-taskmanage.TaskerChan:
			if atomic.LoadInt64(&taskmanage.excutingroutine) < int64(taskmanage.routinenum) {
				go func(ts Tasker) {
					atomic.AddInt64(&taskmanage.excutingroutine, 1)
					//panic捕捉
					defer func() {
						atomic.AddInt64(&taskmanage.excutingroutine, -1)

						panicerr := recover()
						//如果是panic，则重启一个goroutine，防止goroutine没了，阻塞成哥程序
						if panicerr != nil {
							var stack string
							for i := 1; ; i++ {
								_, file, line, ok := runtime.Caller(i)
								if !ok {
									break
								}
								stack = stack + fmt.Sprintln(file, line)
							}
							logs.Logger.Errorf("%v\n%s", panicerr, stack)
						}
					}()

					ts.Excute()

				}(stask)
			}
		case <-time.After(taskmanage.exitTime):
			if taskmanage.isDestroy {
				//退出goroutine
				return
			}
		}
	}
}

//重写inputtask，一旦有输入task，能执行则立即启动一个routine执行，否则直接返回task已经满负荷运行的error
func (taskmanage *RealtimeAcllocManage) InputTask(ts Tasker, timeout time.Duration) error {
	if atomic.LoadInt64(&taskmanage.excutingroutine) < int64(taskmanage.routinenum) {
		go func(ts Tasker) {
			atomic.AddInt64(&taskmanage.excutingroutine, 1)
			//panic捕捉
			defer func() {
				atomic.AddInt64(&taskmanage.excutingroutine, -1)

				panicerr := recover()
				//如果是panic，则重启一个goroutine，防止goroutine没了，阻塞成哥程序
				if panicerr != nil {
					var stack string
					for i := 1; ; i++ {
						_, file, line, ok := runtime.Caller(i)
						if !ok {
							break
						}
						stack = stack + fmt.Sprintln(file, line)
					}
					logs.Logger.Errorf("%v\n%s", panicerr, stack)
				}
			}()

			ts.Excute()

		}(ts)
		return nil
	} else {
		return ErrTimeout
	}
}

// chanlen 任务缓冲区chan的长度
// goroutine 同时执行任务的task的个数，目前采用初始化时启动，后续动态启动，和删除
//保证只要一个taskmanage，否则很多goroutine在运行啊，一个taskmanage只能运行一组task的routine，这个怎么处理。
func InitRealtimeAcllocTaskerManager(chanlen int, goroutinenum int) TaskerManager {
	taskmanage := new(RealtimeAcllocManage)
	taskmanage.TaskerChan = make(chan Tasker, chanlen)
	taskmanage.chanlen = chanlen
	taskmanage.routinenum = goroutinenum
	taskmanage.runningtroutine = 0
	taskmanage.exitTime = time.Second
	taskmanage.isDestroy = false
	return taskmanage
}

//运行整个task任务,不用启动，只是初始化一下。
func (taskmanage *RealtimeAcllocManage) Run() {
	taskmanage.mutex.Lock()
	defer taskmanage.mutex.Unlock()
	if taskmanage.isrunning == false {
		taskmanage.isrunning = true
	}
}
