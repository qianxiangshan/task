//author qian
//time 2016-11-16
//提供一个异步的task管理，固定执行task的goroutine数量。限制同时执行的任务个数,提供失败通知。
package task

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

//task 接口
type Tasker interface {
	//任务执行的函数，不需要开启goroutine执行。task会启动新的routine执行task
	//不同的task会并行执行，使用者保证并发安全
	Excute()
}

type TaskerManage struct {
	//传递任务使用的chan
	//超时时间视自己业务和性能需要确定
	//向taskerchan中写入任务即可，chan造成是性能损失可以忽略
	TaskerChan chan Tasker
	chanlen    int
	routinenum int
	//正在运行的routine
	runningtroutine int64
	excutingroutine int64
	isrunning       bool
	//销毁标志
	isDestroy bool
	//退出任务处理goroutine前等待的时间。
	exitTime time.Duration
	mutex    sync.Mutex
}

type TaskerManager interface {
	Info() string
	Destroy()
	Run()
	TotalExcutingTaskNum() int64
	MaxRoutine() int
	InputTask(ts Tasker, timeout time.Duration) error
}

var (
	ErrTimeout = fmt.Errorf("timeout")
)

//返回使用情况
//allroutine； runningroutine，usage，bufftasknum
func (taskmanage *TaskerManage) Info() string {
	return fmt.Sprintf("allroutine %d, runningroutine %d ，usage %d ,bufftasknum %d cap taskchanbum %d ", taskmanage.routinenum, taskmanage.runningtroutine, int(float64(taskmanage.excutingroutine)/float64(taskmanage.routinenum)*100), len(taskmanage.TaskerChan), cap(taskmanage.TaskerChan))
}

//回收taskmanage的goroutine，可以保证task任务执行完毕。会延后不定期的时间销毁全部的routine
//非阻塞。单个对象不可多次启动退出。新的任务使用新的routine
// 怎么保证destroy后，chan写入阻塞
func (taskmanage *TaskerManage) Destroy() {
	taskmanage.isDestroy = true
}

func (taskmanage *TaskerManage) MaxRoutine() int {
	return taskmanage.routinenum
}
func (taskmanage *TaskerManage) TotalExcutingTaskNum() int64 {
	return atomic.LoadInt64(&taskmanage.excutingroutine)
}

func (taskmanage *TaskerManage) InputTask(ts Tasker, timeout time.Duration) error {
	select {
	case taskmanage.TaskerChan <- ts:
		return nil
	case <-time.After(timeout):
		return ErrTimeout
	}
}
