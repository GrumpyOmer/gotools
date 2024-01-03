package logCenter

import (
	"testing"
	"time"
)

func TestLogcenterInit(t *testing.T) {
	//SaveFSync(true)
	Add("omer test1")
	time.Sleep(2 * time.Second)
}
