package utils

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"strconv"
)

type AutoScalingBloomFilter struct {
	rdb        *redis.Client
	baseKey    string
	errorRate  float64
	capacity   int
	currentIdx int
}

// NewAutoScalingBloomFilter 初始化自动扩展布隆过滤器
func NewAutoScalingBloomFilter(ctx context.Context, rdb *redis.Client, baseKey string, errorRate float64, capacity int) (*AutoScalingBloomFilter, error) {
	var err error
	bf := &AutoScalingBloomFilter{
		rdb:       rdb,
		baseKey:   baseKey,
		errorRate: errorRate,
		capacity:  capacity,
	}

	err = bf.createNewFilter(ctx)
	return bf, err
}

// createNewFilter 创建新的布隆过滤器
func (bf *AutoScalingBloomFilter) createNewFilter(ctx context.Context) error {
	err := bf.currIndexIncr(ctx)
	if err != nil {
		return err
	}
	key := bf.getCurrentKey()
	err = bf.rdb.Do(ctx, "BF.RESERVE", key, bf.errorRate, bf.capacity).Err()
	if err != nil {
		return err
	}
	return nil
}

func (bf *AutoScalingBloomFilter) currIndexIncr(ctx context.Context) error {
	index, err := bf.rdb.Incr(ctx, bf.baseKey+"index").Result()
	if err != nil {
		return err
	}
	bf.currentIdx = int(index)
	return nil
}

func (bf *AutoScalingBloomFilter) getCurrIndex(ctx context.Context) (int, error) {
	ret, err := bf.rdb.Get(ctx, bf.baseKey+"index").Result()
	if err != nil {
		return 0, err
	}
	index, err := strconv.Atoi(ret)
	if err != nil {
		return 0, err
	}
	return index, nil
}

// getCurrentKey 获取当前布隆过滤器的键
func (bf *AutoScalingBloomFilter) getCurrentKey() string {
	return bf.baseKey + ":" + strconv.Itoa(bf.currentIdx)
}

var BloomFilterFullErr = errors.New("ERR item exists")

// Add 向布隆过滤器添加元素
func (bf *AutoScalingBloomFilter) Add(ctx context.Context, element string) error {
	key := bf.getCurrentKey()
	_, err := bf.rdb.Do(ctx, "BF.ADD", key, element).Result()
	if err != nil {
		// 如果错误是因为达到容量限制，创建新的布隆过滤器
		if errors.Is(err, BloomFilterFullErr) {
			err = bf.createNewFilter(ctx)
			if err != nil {
				return err
			}
			key = bf.getCurrentKey()
			_, err = bf.rdb.Do(ctx, "BF.ADD", key, element).Result()
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

// Exists 检查元素是否存在于布隆过滤器中
func (bf *AutoScalingBloomFilter) Exists(ctx context.Context, element string) bool {
	for i := 1; i <= bf.currentIdx; i++ {
		key := bf.baseKey + ":" + strconv.Itoa(i)
		result, err := bf.rdb.Do(ctx, "BF.EXISTS", key, element).Result()
		if err != nil {
			return false
		}
		if result.(int64) == 1 {
			return true
		}
	}
	return false
}
