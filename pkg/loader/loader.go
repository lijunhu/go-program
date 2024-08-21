package loader

import (
	"fmt"
	"sync"
)

// Loader 初始化 .
type Loader interface {
	GetName() string
	RunLoad() error
}

var (
	internal []Loader   // 内部变量
	lock     sync.Mutex // 内部变量锁
)

func Register(loader Loader) {
	lock.Lock()
	defer lock.Unlock()
	internal = append(internal, loader) //nolint:wsl
}

func Load() error {
	for _, loader := range internal {
		if err := loader.RunLoad(); err != nil {
			return fmt.Errorf("%s run fail err:%w", loader.GetName(), err)
		}
	}
	return nil
}
