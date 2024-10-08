package test

import (
	"sync"
	"testing"
	
)

func sortArr(arr []int)  {
	length := len(arr)

	for i:=0;i<length;i++{
		num := arr[i]
		index := i
		for j:=i;j<length;j++ {
			if num>arr[j] {
				num = arr[j]
				index = j
			}
		}

		if index != i {
			arr[index] = arr[i]
			arr[i] = num
		}
	}

}

func Test_sort(t *testing.T)  {

	a := []int{4,3,6,7,10,4,5,6}
	sortArr(a)
	t.Log(a)
}


func errorGroup(){
	var wait sync.WaitGroup

		wait.Add(1)
	wait.Wait()
}