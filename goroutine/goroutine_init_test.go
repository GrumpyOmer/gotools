package goroutine

import (
	"fmt"
	"github.com/pkg/errors"
	"sync"
	"testing"
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
