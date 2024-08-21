package test

import (
	"bufio"
	"fmt"
	"github.com/go-pg/pg/v10"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func Test_Export_Sql_File(t *testing.T) {

	fPath := "D:\\work\\urlsaas"

	//fileNamePrefixList := []string{
	//	"t_black_regs",
	//}

	fileNamePrefixList := []string{
		"t_tag",
	}
	suffix := ".sql"
	outsuffix := ".txt"

	pg := pg.Connect(&pg.Options{
		Addr:         "pg-d23cbfaa-test.avlyun.org:5432",
		User:         "urlsaas_urls",
		Password:     "x7CsQvjPq6",
		Database:     "urlsaas_urls",
		PoolSize:     10,
		MaxRetries:   3,
		MinIdleConns: 10,
		DialTimeout:  30 * time.Second,
	}).WithTimeout(30 * time.Second)

	for _, prefix := range fileNamePrefixList {
		filename := filepath.Join(fPath, prefix+suffix)
		f, err := os.Open(filename)
		if err != nil {
			t.Error(err)
			return
		}
		sc := bufio.NewScanner(f)

		outfile := filepath.Join(fPath, prefix+outsuffix)
		_, err = os.Stat(outfile)
		if err == nil {
			os.Remove(outfile)
		}
		of, err := os.Create(outfile)
		if err != nil {
			t.Error(err)
			return
		}

		writer := bufio.NewWriter(of)

		for sc.Scan() {
			line := sc.Text()
			fields := strings.Split(line, "\t")
			value := "insert into " + prefix + "(id,tag_id,value) values (" + fields[0] + "," + fields[1] + ",'" + fields[2] + "');\n"
			fmt.Println(value)
			_, err = pg.Exec(value)
			if err != nil {
				t.Error(err)
				return
			}
			_, err = writer.WriteString(value)
			if err != nil {
				t.Error(err)
				return
			}
		}
		_ = writer.Flush()

	}

}

type Tag struct {
	TagID int `json:"tag_id" orm:"tag_id"`
}

func Test_Query_Table(t *testing.T) {

	pg := pg.Connect(&pg.Options{
		Addr:         "pg-d23cbfaa-test.avlyun.org:5432",
		User:         "urlsaas_urls",
		Password:     "x7CsQvjPq6",
		Database:     "urlsaas_urls",
		PoolSize:     10,
		MaxRetries:   3,
		MinIdleConns: 10,
		DialTimeout:  30 * time.Second,
	}).WithTimeout(30 * time.Second)
	var tags []*Tag
	_, err := pg.Query(&tags, "select DISTINCT(tag_id) from t_black_domain;")
	if err != nil {
		t.Error(err)
		return
	}

	var targetTags []*Tag
	_, err = pg.Query(&targetTags, "select DISTINCT(tag_id) from t_tag;")
	if err != nil {
		t.Error(err)
		return
	}
	targetTagMap := make(map[int]*Tag, len(targetTags))
	for _, tag := range targetTags {
		targetTagMap[tag.TagID] = tag
	}

	var missTagID []int
	for _, tag := range tags {
		if _, ok := targetTagMap[tag.TagID]; !ok {
			missTagID = append(missTagID, tag.TagID)
		}
	}
	t.Log(missTagID)
}
