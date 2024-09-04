package s3

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"
)

var bucket *Bucket = &Bucket{
	AccessKey:      "YL3GI7SV3FD4YUDM8GPG",
	SecretKey:      "XtbPvPVLQQkCG9F4PFS88AZDPKaoHpLh5u6kT5hY",
	EndPoint:       "http://tcstore1.17usoft.com",
	Region:         "ceph",
	BucketName:     "assets",
	ConnectTimeout: 10 * time.Second,
	ReadTimeout:    30 * time.Second,
	WriteTimeout:   30 * time.Second,
	RequestTimeout: 30 * time.Second,
}

func BenchmarkPutObject(b *testing.B) {

	lock := new(sync.Mutex)
	wg := new(sync.WaitGroup)

	go func(args interface{}) {

		wg.Add(1)
		lock.Lock()

		defer lock.Unlock()
	}("1")

	b.ResetTimer()
	file, err := os.Open("/Users/lijh/test/118.6debff56405727696eee.js")
	if err != nil {
		fmt.Print(err)
	}
	data, _ := io.ReadAll(file)
	for i := 0; i < b.N; i++ {
		startTime := time.Now()
		err = bucket.PutObject("/kylinfastapptest/118.6debff56405727696eee.js", data)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(time.Now().Sub(startTime).Milliseconds())
	}

}

func TestHttpClient(t *testing.T) {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (c net.Conn, err error) {
				c, err = net.DialTimeout(netw, addr, time.Second*10)
				if err != nil {
					return
				}

				var deadline time.Time
				deadline = time.Now().Add(time.Second * 10)
				c.SetDeadline(deadline)
				return
			},
		},
	}

	req := new(http.Request)
	resp := new(http.Response)
	var err error
	var data []byte

	req, err = http.NewRequest(http.MethodGet, "http://10.100.202.119/test/http/client", nil)
	resp, err = client.Do(req)

	if err != nil {
		fmt.Printf("body:%s,err:%s", "", err)
		return
	}
	defer resp.Body.Close()
	data, err = ioutil.ReadAll(resp.Body)

	fmt.Println(string(data))
}
