package mysql

import (
	"strings"
	"testing"
)

func TestGetMasterReinitializesAfterConfigChange(t *testing.T) {
	if err := ConfigInit([]byte(`{
		"master":{"user":"root","pass":"123456","ip":"127.0.0.1","port":"65001","db_name":"hq"},
		"slave":[]
	}`)); err != nil {
		t.Fatal(err)
	}

	_, err := Client().GetMaster()
	if err == nil || !strings.Contains(err.Error(), "65001") {
		t.Fatalf("expected error mentioning first port, got %v", err)
	}

	if err := ConfigInit([]byte(`{
		"master":{"user":"root","pass":"123456","ip":"127.0.0.1","port":"65002","db_name":"hq"},
		"slave":[]
	}`)); err != nil {
		t.Fatal(err)
	}

	_, err = Client().GetMaster()
	if err == nil || !strings.Contains(err.Error(), "65002") {
		t.Fatalf("expected error mentioning second port after reconfig, got %v", err)
	}
}

func TestGetSlaveReinitializesAfterConfigChange(t *testing.T) {
	if err := ConfigInit([]byte(`{
		"master":{"user":"root","pass":"123456","ip":"127.0.0.1","port":"65001","db_name":"hq"},
		"slave":[{"user":"root","pass":"123456","ip":"127.0.0.1","port":"65003","db_name":"hq"}]
	}`)); err != nil {
		t.Fatal(err)
	}

	_, err := Client().GetSlave()
	if err == nil || !strings.Contains(err.Error(), "65003") {
		t.Fatalf("expected error mentioning first slave port, got %v", err)
	}

	if err := ConfigInit([]byte(`{
		"master":{"user":"root","pass":"123456","ip":"127.0.0.1","port":"65001","db_name":"hq"},
		"slave":[{"user":"root","pass":"123456","ip":"127.0.0.1","port":"65004","db_name":"hq"}]
	}`)); err != nil {
		t.Fatal(err)
	}

	_, err = Client().GetSlave()
	if err == nil || !strings.Contains(err.Error(), "65004") {
		t.Fatalf("expected error mentioning second slave port after reconfig, got %v", err)
	}
}
