package test

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

type Location struct {
	U   U    `json:"uri"`
	M   M    `json:"matchType"`
	Rs  []R  `json:"rewrite"`
	P   P    `json:"proxy_pass"`
	Ifs []If `json:"ifs"`
	Ss  []S  `json:"sets"`
}

type If struct {
	P P `json:"p"`
	R R `json:"r"`
}

type M string
type P string
type U string
type R Pairs
type S Pairs

type Pairs struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var locationTexts [][]string
var locationText []string
var inLocation bool
var inIf bool

func TestReadFile(T *testing.T) {
	filePath := "/Users/lijh/Work/move_config/"
	fileName := "www.ly.com.txt"
	err := ReadLine(filePath+fileName, processLine)
	if err != nil {
		fmt.Printf("err:%s", err)
	}

	var locations []Location
	for _, locationConfig := range locationTexts {
		locations = append(locations, parseLocationConfig(locationConfig))
	}
	dataBytes, _ := json.MarshalIndent(locations, "", "\t")
	_ = ioutil.WriteFile(filePath+"parse."+fileName, dataBytes, 0644)

	fmt.Println(string(dataBytes))
}

func locationStart(lineText string) bool {
	locationStart := strings.Index(lineText, "location")
	locationChunkStart := strings.Index(lineText, "{")
	if locationStart >= 0 && locationChunkStart >= 0 {
		return true
	}
	return false
}

func IfStart(lineText string) bool {
	start := strings.Index(lineText, "if")
	ifChunkStart := strings.Index(lineText, "{")
	if start >= 0 && ifChunkStart >= 0 {
		return true
	}
	return false
}

func chunkEnd(lineText string) bool {
	locationChunkEnd := strings.Index(lineText, "}")
	lineText = strings.Replace(lineText, " ", "", -1)
	lineText = strings.Replace(lineText, "	", "", -1)
	if locationChunkEnd >= 0 && lineText == "}" {
		return true
	}
	return false
}

func processLine(line string) {
	if locationStart(strings.Trim(line, "")) {
		inLocation = true
	}

	if IfStart(line) {
		inIf = true
	}
	if inLocation {
		locationText = append(locationText, line)
	}
	if chunkEnd(strings.Trim(string(line), "")) && inLocation && !inIf {
		if len(locationText) > 0 {
			locationTexts = append(locationTexts, locationText)
		}
		inLocation = false
		locationText = []string{}
	}
	if chunkEnd(strings.Trim(string(line), "")) && inIf {
		inIf = false
	}
}

var (
	spaceStr     = " "
	setStr       = "set"
	ifStr        = "if"
	proxyPassStr = "proxy_pass"
	rewriteStr   = "rewrite"
	sign         = "^"
)

func parseLocationConfig(locationConfig []string) (location Location) {

	var ss []S
	var ifs []If
	var rs []R
	var p P
	var inIf bool
	var u U
	var m M
	var f If
	var r R
	for _, lineText := range locationConfig {
		array := strings.Split(strings.Trim(lineText, spaceStr), spaceStr)
		array = FilterString(array, func(x string) bool {
			return len(x) > 0 && x != " "
		})
		arrayLen := len(array)
		setIndex := IndexString(array, setStr)
		ifIndex := IndexString(array, ifStr)
		proxyPassIndex := IndexString(array, proxyPassStr)
		rewriteIndex := IndexString(array, rewriteStr)

		if ifIndex >= 0 {
			f = If{}
			inIf = true
		}
		if chunkEnd(lineText) && inIf {
			ifs = append(ifs, f)
			f = If{}
			inIf = false
		}

		if inIf {
			if proxyPassIndex >= 0 {
				p = P(array[arrayLen-1])
				f.P = p
			}

			if rewriteIndex >= 0 {
				r = R(Pairs{
					Value: array[arrayLen-2] + array[arrayLen-1],
					Key:   array[1],
				})
				f.R = r
			}

			continue
		}

		if locationStart(lineText) {

			uri := array[arrayLen-2]
			if strings.HasPrefix(uri, sign) {
				u = U(uri[1:])
			}
			m = M(array[arrayLen-3])
		}
		if setIndex >= 0 {

			s := S(Pairs{
				Value: array[arrayLen-1],
				Key:   array[1],
			})
			ss = append(ss, s)
		}

		if proxyPassIndex >= 0 {
			p = P(array[arrayLen-1])
		}

		if rewriteIndex >= 0 {
			r = R(Pairs{
				Value: array[arrayLen-2] + array[arrayLen-1],
				Key:   array[1],
			})
			rs = append(rs, r)
		}

	}
	location.M = m
	location.U = u
	location.Ss = ss
	location.Rs = rs
	location.P = p
	location.Ifs = ifs
	return
}

func ReadLine(filePth string, hookfn func(string)) error {
	f, err := os.Open(filePth)
	if err != nil {
		return err
	}
	defer f.Close()

	bfRd := bufio.NewReader(f)
	for {
		line, err := bfRd.ReadBytes('\n')
		lineText := strings.Replace(string(line), "\r\n", "", -1)

		hookfn(lineText) //放在错误处理前面，即使发生错误，也会处理已经读取到的数据。
		if err != nil { //遇到任何错误立即返回，并忽略 EOF 错误信息
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
	return nil
}

func ContainsString(a []string, x string) bool {
	if len(a) == 0 {
		return false
	}
	for k := range a {
		if a[k] == x {
			return true
		}
	}
	return false
}

func IndexString(a []string, x string) int {
	if len(a) <= 0 {
		return -1
	}
	for k := range a {
		if a[k] == x {
			return k
		}
	}
	return -1
}

func FilterString(a []string, keep func(x string) bool) []string {
	if len(a) == 0 {
		return a
	}

	n := 0
	for _, v := range a {
		if keep(v) {
			a[n] = v
			n++
		}
	}

	return a[:n]
}
