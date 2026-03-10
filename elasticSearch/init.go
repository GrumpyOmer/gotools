package elasticSearch

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/olivere/elastic/v7"
)

type (
	config struct {
		Address []string `json:"address"`
	}
	client struct {
		es   *elastic.Client
		mu   sync.Mutex
		done bool
	}
)

var (
	// 实例对象
	esClient = client{}
	// 配置对象
	cf = config{}
)

// ConfigInit es配置信息初始化
func ConfigInit(c []byte) error {
	// 外部传入json字符串配置
	err := json.Unmarshal(c, &cf)
	if err != nil {
		// 初始化失败
		return err
	}
	return nil
}

// GetESClient 获取客户端实例
func GetESClient() (*elastic.Client, error) {
	esClient.mu.Lock()
	defer esClient.mu.Unlock()

	// 如果已经初始化成功，直接返回
	if esClient.es != nil {
		return esClient.es, nil
	}

	// 初始化客户端
	var err error
	esClient.es, err = initClient()
	if err != nil {
		// 初始化失败，重置状态以便下次重试
		esClient.es = nil
		return nil, err
	}

	esClient.done = true
	return esClient.es, nil
}

// 初始化实例
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
