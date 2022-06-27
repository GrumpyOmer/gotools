package mysql

import (
	"testing"
)

func TestConfigInit(t *testing.T) {
	var (
		err error
	)
	err = ConfigInit([]byte(`{
    "master": {
        "user": "root",
        "pass": "123456",
        "ip": "localhost",
        "port": "3306",
        "db_name": "hq",
        "max_igle_conn": 0,
        "max_open_conn": 0,
        "conn_max_life_time": 0
        },
    "slave": [
        {
			"user": "root",
			"pass": "123456",
			"ip": "localhost",
			"port": "3306",
			"db_name": "hq",
			"max_igle_conn": 1,
			"max_open_conn": 1,
			"conn_max_life_time": 0
        }
    ]
 }`))
	if err != nil {
		t.Fatal(err)
	}
	_,err =Client().GetMaster()
	if err != nil {
		t.Fatal(err)
	}
	_,err =Client().GetSlave()
	if err != nil {
		t.Fatal(err)
	}
}
