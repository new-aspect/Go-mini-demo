package main

import (
	"fmt"
	"sync"
	"time"
)

// 我是想要控制每个不同的项目对应不同的锁，进过调整，发现一个优雅便捷的方式
//

// 全局变量：存储projName对应的锁
//
// Example:
// // 获取该项目的锁
// lock := getProjectLock(projName)
// lock.Lock() // 加锁
// defer lock.Unlock()
var projectLocks sync.Map

// 锁管理函数：获取projName的锁
func getProjectLock(projName string) *sync.Mutex {
	// 检查projName对应的锁是否已经存在
	lock, exists := projectLocks.Load(projName)
	if !exists {
		// 如果不存在，创建一个新的锁并存储到sync.Map中
		lock = &sync.Mutex{}
		projectLocks.Store(projName, lock)
	}
	return lock.(*sync.Mutex)
}

func main() {
	// 共享计数器，记录每个项目的操作次数
	projectCounts := make(map[string]int)
	var mu sync.Mutex // 保证对 projectCounts 的线程安全操作

	// 并发测试参数
	numProjects := 100    // 不同的projName数量
	numGoroutines := 1000 // 总的并发goroutine数量

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	startTime := time.Now()

	// 启动多个goroutine，模拟并发操作
	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			defer wg.Done()

			// 模拟随机选择项目名称
			projName := fmt.Sprintf("Project-%d", i%numProjects)

			// 获取该项目的锁
			//lock := getProjectLock(projName)
			//lock.Lock() // 加锁
			//defer lock.Unlock()

			// 修改共享数据
			mu.Lock()
			projectCounts[projName]++
			mu.Unlock()

			// 模拟耗时操作
			time.Sleep(1 * time.Millisecond)
		}(i)
	}

	wg.Wait() // 等待所有goroutine完成
	duration := time.Since(startTime)

	// 验证结果
	fmt.Println("All goroutines completed in:", duration)
	fmt.Println("Final project counts:")

	for projName, count := range projectCounts {
		fmt.Printf("%s: %d\n", projName, count)
	}

	// 检查总计数是否正确
	expectedCount := numGoroutines
	actualCount := 0
	for _, count := range projectCounts {
		actualCount += count
	}

	if actualCount == expectedCount {
		fmt.Println("Test passed: All operations were successfully synchronized.")
	} else {
		fmt.Printf("Test failed: Expected total count %d, but got %d\n", expectedCount, actualCount)
	}
}
