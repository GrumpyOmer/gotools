package goroutine

import (
	"fmt"
	"github.com/pkg/errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestGoroutineInit(t *testing.T) {
	var (
		err1 error
		err2 error
	)
	SetGoroutineNumber(0)
	wg := sync.WaitGroup{}
	wg.Add(2)
	t1 := func() {
		fmt.Println("test1")
		err1 = errors.New("err1")
	}
	t2 := func() {
		fmt.Println("test2")
		err2 = errors.New("err2")
	}
	if err := MakeTask(t1, &wg); err != nil {
		wg.Done()
		t.Log(err)
	}
	if err := MakeTask(t2, &wg); err != nil {
		wg.Done()
		t.Log(err)
	}
	wg.Wait()
	fmt.Println(err1)
	fmt.Println(err2)

}

func TestQueuedTaskCanRunAfterSlotFrees(t *testing.T) {
	if err := SetGoroutineNumber(1); err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	firstStarted := make(chan struct{}, 1)
	releaseFirst := make(chan struct{})
	var ranSecond int32

	first := func() {
		firstStarted <- struct{}{}
		<-releaseFirst
	}
	second := func() {
		atomic.StoreInt32(&ranSecond, 1)
	}

	if err := MakeTask(first, &wg); err != nil {
		t.Fatal(err)
	}
	<-firstStarted
	if err := MakeTask(second, &wg); err != nil {
		t.Fatal(err)
	}

	select {
	case <-time.After(100 * time.Millisecond):
		if atomic.LoadInt32(&ranSecond) != 0 {
			t.Fatal("second task should still be queued")
		}
	}

	close(releaseFirst)
	wg.Wait()
	if atomic.LoadInt32(&ranSecond) != 1 {
		t.Fatal("expected queued task to run after first task finished")
	}
}
