package main

import (
	"fmt"
	"sync"
	"time"
)

// 好的！我们来从一个简单的例子理解锁的概念。锁的目的是防止多个线程或 Goroutine 同时修改同一个资源，导致数据错误或程序崩溃。
//
// 场景描述：一个银行账户的存取款问题
//
// 假设有一个银行账户，初始余额是 100 元。多个 Goroutine 同时对这个账户进行存取操作。我们来看看锁是如何保护数据一致性的。

var balance = 100 // 银行账户余额

func deposit(amount int) {
	b := balance // 读取余额，这里为什么要写 b := balance ,而不是直接 balance = balance + amount
	time.Sleep(1 * time.Millisecond)
	balance = b + amount
}

func withdraw(amount int) {
	b := balance
	time.Sleep(1 * time.Millisecond)
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
