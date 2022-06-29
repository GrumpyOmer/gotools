package configController

import (
	"testing"
	"time"
)

func TestConfigInit(t *testing.T) {
	SetPubDir(".")
	SetXmlConfigName("xml_test.xml")
	time.Sleep(time.Second)
	t.Log(GetXmlField("test1"))
	SetXmlConfigName("xml_test1.xml")
	time.Sleep(time.Second)
	t.Log(GetXmlField("test1"))
}
