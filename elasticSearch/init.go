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
)

var (
	// 实例对象
	esClient *elastic.Client
	esMu     sync.Mutex
	// 配置对象
	cf                = config{}
	httpClientFactory = defaultHTTPClient
)

// ConfigInit es配置信息初始化
func ConfigInit(c []byte) error {
	// 外部传入json字符串配置
	var next config
	err := json.Unmarshal(c, &next)
	if err != nil {
		// 初始化失败
		return err
	}
	esMu.Lock()
	cf = next
	esClient = nil
	esMu.Unlock()
	return nil
}

// GetESClient 获取客户端实例（并发安全）
func GetESClient() (*elastic.Client, error) {
	esMu.Lock()
	defer esMu.Unlock()
	if esClient != nil {
		return esClient, nil
	}
	client, err := initClient()
	if err != nil {
		return nil, err
	}
	esClient = client
	return esClient, nil
}

// 初始化实例
func initClient() (*elastic.Client, error) {
	httpClient := httpClientFactory()
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

func defaultHTTPClient() *http.Client {
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
	return httpClient
}
