package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	balance = 100      // 银行账户余额
	mu      sync.Mutex // 定义一个锁
)

func deposit(amount int) {
	mu.Lock()         // 加锁
	defer mu.Unlock() // 操作完成后解锁
	b := balance      // 读取余额，这里为什么要写 b := balance ,而不是直接 balance = balance + amount
	time.Sleep(5 * time.Second)
	balance = b + amount
}

func withdraw(amount int) {
	mu.Lock()         // 加锁
	defer mu.Unlock() // 操作完成后解锁

	b := balance
	time.Sleep(5 * time.Second)
	balance = b - amount
}

// 最终余额可能不是
func main() {
	var wg sync.WaitGroup

	// 启动多个 Goroutine，同时存取款
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			deposit(10) // 存 10 元
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			withdraw(10) // 取 10 元
		}()
	}

	wg.Wait() // 等待所有 Goroutine 完成
	fmt.Println("最终余额:", balance)
}
