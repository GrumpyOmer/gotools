package redis

import (
	"strings"
	"testing"
)

func TestGetMasterReinitializesAfterConfigChange(t *testing.T) {
	if err := ConfigInit([]byte(`{
		"master":{"host":"127.0.0.1","port":"65011","network":"tcp"},
		"slave":[]
	}`)); err != nil {
		t.Fatal(err)
	}

	conn, err := Client().GetMaster()
	if conn != nil {
		conn.Close()
	}
	if err == nil || !strings.Contains(err.Error(), "65011") {
		t.Fatalf("expected error mentioning first port, got %v", err)
	}

	if err := ConfigInit([]byte(`{
		"master":{"host":"127.0.0.1","port":"65012","network":"tcp"},
		"slave":[]
	}`)); err != nil {
		t.Fatal(err)
	}

	conn, err = Client().GetMaster()
	if conn != nil {
		conn.Close()
	}
	if err == nil || !strings.Contains(err.Error(), "65012") {
		t.Fatalf("expected error mentioning second port after reconfig, got %v", err)
	}
}

func TestGetSlaveReinitializesAfterConfigChange(t *testing.T) {
	if err := ConfigInit([]byte(`{
		"master":{"host":"127.0.0.1","port":"65011","network":"tcp"},
		"slave":[{"host":"127.0.0.1","port":"65013","network":"tcp"}]
	}`)); err != nil {
		t.Fatal(err)
	}

	conn, err := Client().GetSlave()
	if conn != nil {
		conn.Close()
	}
	if err == nil || !strings.Contains(err.Error(), "65013") {
		t.Fatalf("expected error mentioning first slave port, got %v", err)
	}

	if err := ConfigInit([]byte(`{
		"master":{"host":"127.0.0.1","port":"65011","network":"tcp"},
		"slave":[{"host":"127.0.0.1","port":"65014","network":"tcp"}]
	}`)); err != nil {
		t.Fatal(err)
	}

	conn, err = Client().GetSlave()
	if conn != nil {
		conn.Close()
	}
	if err == nil || !strings.Contains(err.Error(), "65014") {
		t.Fatalf("expected error mentioning second slave port after reconfig, got %v", err)
	}
}
