package logCenter

import (
	"testing"
	"time"
)

func TestLogcenterInit(t *testing.T) {
	Add("omer test1")
	time.Sleep(2 * time.Second)
}
