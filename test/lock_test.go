package test

import (
	"github.com/sirupsen/logrus"
	"sync"
	"testing"
)

var lock sync.Mutex

func TestLock(t *testing.T) {

	log := logrus.New()
	//log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.TextFormatter{DisableTimestamp:false})

	log.Info("1232131313")

}


func sumFunc(i int, sum *int, c chan int) {
	//lock.Lock()
	*sum = (*sum) + i
	//defer lock.Unlock()
}
