package test

import (
	"fmt"
	"github.com/bits-and-blooms/bloom/v3"
	"math"
	"testing"
)

func Test_Bloom_Filter(t *testing.T) {

	bloom.EstimateParameters(1000000, 0.001)
}

func Test_Redis_Bloom_Filter(t *testing.T) {
	//redis_bloom_go.NewClientFromPool()

	n := 1000000000 // 预期要存储的元素数量
	p := 0.0001     // 目标误报率
	m, k, memoryBytes := calculateBloomFilterMemory(n, p)
	fmt.Printf("位数组大小 (m): %d 位\n", m)
	fmt.Printf("哈希函数数量 (k): %d\n", k)
	fmt.Printf("内存占用大小: %.2f 字节 (%.2f KB)\n,内存 (%.2f GB)\n", memoryBytes, memoryBytes/1024.0, memoryBytes/(1024.0*1024.0*1024.0))

}

func calculateBloomFilterMemory(n int, p float64) (int, int, float64) {
	m := -float64(n) * math.Log(p) / math.Pow(math.Ln2, 2)
	k := (m / float64(n)) * math.Ln2
	memoryBytes := m / 8.0
	return int(m), int(math.Ceil(k)), memoryBytes
}

func main() {

}
