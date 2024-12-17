package main

import (
	"fmt"
	"sync"
)

// 我们现在有10个人的账户，请你引导我怎么取写这个代码，注意，不能简单的用只要有操作就锁住所有人的账户，A读取写入时只锁A的账户，不影响B的读取和写入
type Account struct {
	balance int
	mu      sync.Mutex
}

// 全局账户表
var accounts sync.Map

// 创建账户
func createAccount(id string, initialBalance int) {
	account := &Account{balance: initialBalance}
	accounts.Store(id, account) // 存入 sync.Map
}

// 这个操作会增加钱
func deposit(id string, amount int) {
	account, exists := getAccount(id)
	if !exists {
		fmt.Printf("Account %s does not exist!\n", id)
		return
	}

	account.mu.Lock()         // 锁定该账户
	defer account.mu.Unlock() // 解锁账户
	account.balance += amount // 更新余额
	fmt.Printf("Deposited %d to account %s. New balance: %d\n", amount, id, account.balance)
}

// 这个操作会减少钱
func withdraw(id string, amount int) {
	account, exists := getAccount(id)
	if !exists {
		fmt.Printf("Account %s does not exist!\n", id)
		return
	}

	account.mu.Lock()         // 锁定该账户
	defer account.mu.Unlock() // 解锁账户

	account.mu.TryLock()

	if account.balance >= amount {
		account.balance -= amount // 更新余额
		fmt.Printf("Withdrew %d from account %s. New balance: %d\n", amount, id, account.balance)
	} else {
		fmt.Printf("Account %s has insufficient funds! Current balance: %d\n", id, account.balance)
	}
}

// 获取账户
func getAccount(id string) (*Account, bool) {
	value, exists := accounts.Load(id) // 从 sync.Map 获取账户
	if !exists {
		return nil, false
	}
	return value.(*Account), true
}

// 这里有10个人在循环5次写钱和读读钱的操作
func main() {
	// 创建多个账户
	for i := 1; i <= 10; i++ {
		createAccount(fmt.Sprintf("User%d", i), 1000)
	}

	var wg sync.WaitGroup
	for i := 1; i <= 10; i++ {
		userID := fmt.Sprintf("User%d", i)
		wg.Add(1) // 为什么这么写能执行到hello，写到函数里面却不行
		go func(uid string) {
			defer wg.Done()
			deposit(uid, 100) // 每个用户存款 100
		}(userID)

		wg.Add(1)
		go func(uid string) {
			defer wg.Done()
			withdraw(uid, 50) // 每个用户取款 50
		}(userID)
	}
	wg.Wait()

	// 打印最终余额
	for i := 1; i <= 10; i++ {
		userID := fmt.Sprintf("User%d", i)
		account, exists := getAccount(userID)
		if exists {
			fmt.Printf("Final balance for %s: %d\n", userID, account.balance)
		}
	}
	fmt.Println("执行结束")
}
