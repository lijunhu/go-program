package test

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"sort"
	"strings"
	"testing"
	"time"
)

type TestConfig struct {
	Uri       string
	MatchType string
}

type TestConfigs []TestConfig

func (s TestConfigs) Len() int      { return len(s) }
func (s TestConfigs) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type SortTestConfigs struct {
	TestConfigs
}

func (s SortTestConfigs) Less(i, j int) bool {

	if s.TestConfigs[i].MatchType != "regex" {
		return true
	}
	uriILen := len(strings.Split(s.TestConfigs[i].Uri, "/"))
	uriJLen := len(strings.Split(s.TestConfigs[j].Uri, "/"))
	if uriILen == uriJLen {
		return len(s.TestConfigs[i].Uri) > len(s.TestConfigs[j].Uri)
	}
	return uriILen > uriJLen
}

func externalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network?")
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}
func TestAppend(t *testing.T) {
	a := 0
	a |=1<<3 + 1<<5
	fmt.Println(a)

	format := "2006-01-02 15:04:05.000"
	fmt.Println(strings.Replace(time.Now().Format(format),"."," ",1))

	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte("YourDataHere")); err != nil {
		panic(err)
	}
	if err := gz.Flush(); err != nil {
		panic(err)
	}
	if err := gz.Close(); err != nil {
		panic(err)
	}
	fmt.Println(b)

	localIp, _ := externalIP()
	fmt.Println(localIp.String())
	fmt.Println(os.Hostname())
	fmt.Println(time.Now().UnixNano() / 1e6)
	testConfigs := []TestConfig{{"/webleonid/aa", "regex"}, {"/webleonid/bb", "regex"},
		{"/webleonid/aa/cc", "regex"}, {"/webleonid/bb/dd", "regex"},
		{"/webleonid/", "regex"}, {"/webleonid", ""}, {"/webleonid", ""}}
	sort.Sort(SortTestConfigs{testConfigs})

	fmt.Println(len(testConfigs))

}


func TestGzipDecompress(t *testing.T) {

	test :="adasdsads"
	result ,_:=GzipCompress([]byte(test))
	result,_=GzipDecompress(result)
fmt.Println(string(result))
}

func GzipDecompress(params []byte) (result []byte, err error) {
	b := bytes.NewReader(params)
	r, err := gzip.NewReader(b)
	if err != nil {
		return nil, err
	} else {
		defer r.Close()
		var decompressData []byte
		decompressData, err = ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		return decompressData, nil
	}
}

func GzipCompress(params []byte) (result []byte, err error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	defer gz.Close()

	if _, err = gz.Write(params); err != nil {
		return
	}
	if err = gz.Close(); err != nil {
		return
	}
	return b.Bytes(), nil
}