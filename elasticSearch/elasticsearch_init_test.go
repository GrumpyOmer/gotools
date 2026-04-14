package elasticSearch

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestGetESClientCanRecoverAfterFailedInit(t *testing.T) {
	badPort := "19201"
	if err := ConfigInit([]byte(fmt.Sprintf(`{"address":["http://127.0.0.1:%s"]}`, badPort))); err != nil {
		t.Fatal(err)
	}
	if _, err := GetESClient(); err == nil || !strings.Contains(err.Error(), badPort) {
		t.Fatalf("expected initial error mentioning port %s, got %v", badPort, err)
	}

	previous := httpClientFactory
	httpClientFactory = func() *http.Client {
		return &http.Client{
			Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				body := `{
					"name":"test-node",
					"cluster_name":"test-cluster",
					"cluster_uuid":"test-cluster-uuid",
					"version":{"number":"7.17.0","build_flavor":"default","build_type":"tar","build_hash":"abc","build_date":"2021-01-01T00:00:00.000Z","build_snapshot":false,"lucene_version":"8.11.1","minimum_wire_compatibility_version":"6.8.0","minimum_index_compatibility_version":"6.0.0-beta1"},
					"tagline":"You Know, for Search"
				}`
				if r.Method == http.MethodHead {
					body = ""
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Status:     "200 OK",
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body:       io.NopCloser(strings.NewReader(body)),
					Request:    r,
				}, nil
			}),
		}
	}
	defer func() { httpClientFactory = previous }()

	if err := ConfigInit([]byte(`{"address":["http://example.com:9200"]}`)); err != nil {
		t.Fatal(err)
	}
	client, err := GetESClient()
	if err != nil {
		t.Fatal(err)
	}
	if client == nil {
		t.Fatal("expected client instance")
	}
}
