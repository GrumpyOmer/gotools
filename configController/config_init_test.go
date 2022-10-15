package configController

import (
	"fmt"
	"testing"
	"time"
)

func TestConfigInit(t *testing.T) {
	time.Sleep(time.Second)
	fmt.Println(GetEnvField("XIXIXI"))
	SetPubDir(".")
	time.Sleep(time.Second)

	fmt.Println(GetEnvField("XIXIXI"))
	SetEnvConfigName(".env.example")
	time.Sleep(time.Second)
	fmt.Println(GetEnvField("XIXIXI"))

	SetXmlConfigName("xml_test.xml")
	time.Sleep(time.Second)
	t.Log(GetXmlField("test1"))
	SetXmlConfigName("xml_test1.xml")
	time.Sleep(time.Second)
	t.Log(GetXmlField("test1"))
}
