package download

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func waitDownloadDone(t *testing.T, wc *writeCounter) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		wc.Lock()
		inUse := wc.inUse
		wc.Unlock()
		if !inUse {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("download did not finish in time")
}

func TestDownloadFileSuccess(t *testing.T) {
	previous := httpGet
	httpGet = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Header:     http.Header{"Content-Length": []string{"11"}},
			Body:       io.NopCloser(strings.NewReader("hello world")),
		}, nil
	}
	defer func() { httpGet = previous }()

	wc := NewWc()
	target := filepath.Join(t.TempDir(), "test.txt")

	if err := wc.DownloadFile(target, "http://example.com/file"); err != nil {
		t.Fatal(err)
	}
	waitDownloadDone(t, wc)

	success, err := wc.DownloadRes()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !success {
		t.Fatal("expected download success")
	}
	content, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "hello world" {
		t.Fatalf("unexpected content: %q", string(content))
	}
}

func TestDownloadFileHTTPErrorCleansTempFile(t *testing.T) {
	previous := httpGet
	httpGet = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Status:     "404 Not Found",
			Header:     http.Header{},
			Body:       io.NopCloser(strings.NewReader("missing")),
		}, nil
	}
	defer func() { httpGet = previous }()

	wc := NewWc()
	target := filepath.Join(t.TempDir(), "test.txt")

	if err := wc.DownloadFile(target, "http://example.com/missing"); err != nil {
		t.Fatal(err)
	}
	waitDownloadDone(t, wc)

	success, err := wc.DownloadRes()
	if err == nil {
		t.Fatal("expected download error")
	}
	if success {
		t.Fatal("expected download failure")
	}
	if !strings.Contains(err.Error(), "404") {
		t.Fatalf("expected 404 error, got %v", err)
	}
	if _, statErr := os.Stat(target); !os.IsNotExist(statErr) {
		t.Fatalf("expected target file to be absent, stat err: %v", statErr)
	}
	if _, statErr := os.Stat(target + ".tmp"); !os.IsNotExist(statErr) {
		t.Fatalf("expected temp file cleanup, stat err: %v", statErr)
	}
}
