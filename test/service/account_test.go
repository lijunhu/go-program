package service

import (
	"fmt"
	"sync"
	"testing"
)

type Op interface {
	add(amount float64)
	sub(amount float64)
	query(name string) (amount float64)
}

type Person struct {
	money float64
	name string
	mutex sync.Mutex
}

func (p *Person) add(amount float64) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.money += amount
}
func (p *Person) sub(amount float64) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.money -=amount
}

func (p *Person) query(name string)(amount float64) {
	if p.name == name {
		amount = p.money
	}
	return
}

func TestPerson(t *testing.T)  {

	waitGroup := new(sync.WaitGroup)
	p :=&Person{
		money: 10000,
		name:  "张三",
	}

	waitGroup.Add(1)
	go func(p *Person) {
		defer waitGroup.Done()
		p.add(100000)
	}(p)
	waitGroup.Add(1)
	go func(p *Person) {
		defer waitGroup.Done()
		p.sub(1000)
	}(p)

	waitGroup.Add(1)
	go func(p *Person) {
		defer waitGroup.Done()
		fmt.Println(p.query("张三"))
	}(p)
	waitGroup.Wait()
	fmt.Println(p.query("张三"))
}
