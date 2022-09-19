package logCenter

import (
	"testing"
	"time"
)

func TestConfigInit(t *testing.T) {
	Add("omer test1")
	time.Sleep(2 * time.Second)
}
