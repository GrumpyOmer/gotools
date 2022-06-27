package redis

import (
	"testing"
)

func TestConfigInit(t *testing.T) {
	var (
		err error
	)
	err = ConfigInit([]byte(`{
		"master":{
			"host":"127.0.0.1",
			"port":"6379",
			"auth":"",
			"user":"",
			"pass":"",
			"db":0,
			"network":"tcp",
			"max_idle":1,
			"max_active":1,
			"idle_timeout":1
		},
		"slave":[
			{
				"host":"127.0.0.1",
				"port":"6379",
				"user":"",
				"pass":"",
				"auth":"",
				"db":0,
				"network":"tcp",
				"max_idle":1,
				"max_active":1,
				"idle_timeout":1
			}
		]
	}`))
	if err != nil {
		t.Fatal(err)
	}
	master, err := Client().GetMaster()
	if err != nil {
		t.Fatal(err)
	}
	defer master.Close()
	res, err := master.Do("GET", "1111")
	t.Log(res)
	t.Log(err)
	slave, err := Client().GetSlave()
	if err != nil {
		t.Fatal(err)
	}
	defer slave.Close()
	res, err = slave.Do("GET", "1111")
	t.Log(res)
	t.Log(err)
}
