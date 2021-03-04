package s3

import (
	"fmt"
	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/s3"
	"go-program/sword/utils"
	"gopkg.in/h2non/filetype.v1"
	"mime"
	"path/filepath"
	"time"
)

type Bucket struct {
	AccessKey  string
	SecretKey  string
	EndPoint   string
	Region     string
	BucketName string
	Bucket     *s3.Bucket
	hasReady   bool
	//因为遇到一个线下容量问题，导致超时了，导致第一版本资源挂起，所以在初始化的时候设置一下超时时间
	// A value of zero means no timeout.
	ConnectTimeout time.Duration

	// ReadTimeout is the maximum time a request attempt will wait
	// for an individual read to complete.
	//
	// A value of zero means no timeout.
	ReadTimeout time.Duration

	// WriteTimeout is the maximum time a request attempt will
	// wait for an individual write to complete.
	//
	// A value of zero means no timeout.
	WriteTimeout time.Duration

	// RequestTimeout is the maximum time a request attempt can
	// take before operations return a timeout error.
	//
	// This includes connection time, any redirects, and reading
	// the response body. The timer remains running after the request
	// is made so it can interrupt reading of the response data.
	//
	// A Timeout of zero means no timeout.
	RequestTimeout time.Duration
}

func (b *Bucket) authBucket() {

	//if has init ready
	if b.hasReady == true {
		return
	}

	auth := aws.Auth{AccessKey: b.AccessKey, SecretKey: b.SecretKey}
	client := s3.New(auth, aws.Region{Name: b.Region, S3Endpoint: b.EndPoint})

	client.ConnectTimeout = b.ConnectTimeout
	client.ReadTimeout = b.ReadTimeout
	client.WriteTimeout = b.WriteTimeout
	client.RequestTimeout = b.RequestTimeout

	b.Bucket = client.Bucket(b.BucketName)
	b.hasReady = true
}

//get bucket list
func (b *Bucket) GetList(prefix string, skip int, limit int) (listObj *s3.ListResp, err error) {
	b.authBucket()

	listObj, err = b.Bucket.List(prefix, "/", "", skip + limit)
	if err != nil {
		return
	}

	//skip some data
	listObj.Contents = listObj.Contents[skip:]
	return listObj, err

}

//get object
func (b *Bucket) GetObject(objectKey string) (data []byte, err error) {
	b.authBucket()
	return b.Bucket.Get(objectKey)
}

//del object
func (b *Bucket) DelObject(objectKey string) error {
	b.authBucket()
	return b.Bucket.Del(objectKey)
}

//exist object
func (b *Bucket) ExistObject(objectKey string) (bool, error) {
	b.authBucket()
	return b.Bucket.Exists(objectKey)
}

//put object
func (b *Bucket) PutObject(objectKey string, data []byte) (err error) {
	b.authBucket()
	//因为服务器的 mime.types　文件映射补全，这里补一个
	//mime.AddExtensionType(".mp3", "audio/mpeg")
	mimeType := mime.TypeByExtension(filepath.Ext(objectKey))
	if mimeType == "" {
		mType, err := filetype.Match(data)
		if err != nil {
			return fmt.Errorf("文件类型识别异常: %s", err.Error())
		}
		mimeType = mType.MIME.Type
	}

	md5Val := utils.Md5(data)
	meta := make(map[string][]string)
	meta["content-hash"] = []string{md5Val}
	s3Opt := s3.Options{
		Meta: meta,
	}

	err = b.Bucket.Put(objectKey, data, mimeType, s3.PublicRead, s3Opt)
	if err != nil {
		return
	}

	//set url expire
	b.Bucket.SignedURL(objectKey, time.Now().Add(1000000*time.Hour))
	//url := self.Bucket.SignedURL(objectKey, time.Now().Add(1000000*time.Hour))
	//fmt.Println("s3.PutObject| put url", url)
	return
}

//PUT bucket
func (b *Bucket) PutBucket() (err error) {
	b.authBucket()
	return b.Bucket.PutBucket(s3.BucketOwnerFull)
}

//DEL BUCKET
func (b *Bucket) DelBucket() (err error) {
	b.authBucket()
	return b.Bucket.DelBucket()
}
