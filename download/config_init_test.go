package download

import (
	"fmt"
	"testing"
	"time"
)

func TestConfigInit(t *testing.T) {
	test := NewWc()
	fmt.Println(test.DownloadFile("test.jpg", "http://e.hiphotos.baidu.com/image/pic/item/a1ec08fa513d2697e542494057fbb2fb4316d81e.jpg"))
	fmt.Println(test.DownloadFile("dasdas", "dsadas"))
	time.Sleep(5 * time.Second)
	t.Log(test.DownloadRes())
}
