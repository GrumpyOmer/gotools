package elasticSearch

import (
	"context"
	"testing"
)

type (
	aa struct {
		v int
	}
)

func TestConfigInit(t *testing.T) {
	var (
		err error
		ctx = context.Background()
	)
	err = ConfigInit([]byte(`{"address":["http://127.0.0.1:9200","http://127.0.0.1:9201","http://127.0.0.1:9202"]}`))
	if err != nil {
		t.Fatal(err)
	}
	instance, err := GetESClient()
	if err != nil {
		t.Fatal(err)
	}
	healthInfo, err := instance.CatHealth().Do(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(healthInfo)
}
