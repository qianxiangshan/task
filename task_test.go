package task

import (
	"testing"
	"time"
)

func TestInit(t *testing.T) {

	var testcases = []struct {
		testname   string
		channum    int
		routinenum int
		runningnum int64
		isdestroy  bool
		isrunning  bool
		exittime   time.Duration
	}{
		{"test1", 1, 1, 0, false, false, time.Second},
	}

	for _, testcase := range testcases {
		tmptestcase := testcase
		t.Run(testcase.testname, func(t *testing.T) {
			manage := InitPreAcllocTaskerManager(tmptestcase.channum, tmptestcase.routinenum)
			if manage.chanlen != tmptestcase.channum {
				t.Fatal("channel num not right")
			}
			if manage.TaskerChan == nil {
				t.Fatal("channel  not init")
				t.Fail()
			}
			if cap(manage.TaskerChan) != manage.chanlen {
				t.Fatal("cap(channel)  not ")
				t.Fail()
			}
			if manage.exitTime != time.Second {
				t.Fail()
			}

		})
	}
}
