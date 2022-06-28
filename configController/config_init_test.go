package configController

import (
	"testing"
)

func TestConfigInit(t *testing.T) {
	t.Log(GetJsonField("test1"))
	SetPubDir(".")
	SetJsonConfigName("json_test.json")
	SetXmlConfigName("xml_test.xml")
	t.Log(GetJsonField("test1"))
	t.Log(GetXmlField("test1"))
}
