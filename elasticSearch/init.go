package elasticSearch

import (
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

type (
	config struct {
		Address []string `json:"address"`
	}
)

var (
	// 实例对象
	esClient = &elastic.Client{}
	// 配置对象
	cf = config{}
)

// es配置信息初始化
func ConfigInit(c []byte) error {
	// 外部传入json字符串配置
	err := json.Unmarshal(c, &cf)
	if err != nil {
		// 初始化失败
		return err
	}
	return nil
}

// get instance
func GetESClient() (*elastic.Client, error) {
	var err error
	if esClient == nil {
		// init
		esClient, err = initClient()
	}
	return esClient, err
}

// init client
func initClient() (*elastic.Client, error) {
	httpClient := &http.Client{}
	httpClient.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100, // maximum number of idle (keep-alive)
		MaxIdleConnsPerHost:   100, //the maximum idle (keep-alive) connections to keep per-host
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	return elastic.NewClient(
		elastic.SetHttpClient(httpClient),
		elastic.SetURL(cf.Address...),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetGzip(false),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetTraceLog(log.New(os.Stdout, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)))
}
