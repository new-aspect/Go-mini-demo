package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
)

//简单案例教程：构建一个 CLI 工具
//
//任务：创建一个简单的 CLI 工具，有两个命令
//
//	1.	hello：打印一段问候语。
//	2.	add：接收两个数字作为参数并返回它们的和。

func main() {
	var rootCmd = &cobra.Command{
		Use:   "hugo",
		Short: "Hugo is a very fast static site generator",
		Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at https://gohugo.io/documentation/`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
			fmt.Println("执行到这里面来了")
		},
	}

	var helloCmd = &cobra.Command{
		Use:   "hello",
		Short: "你好",
		Long:  "你好，很高兴见到你，和你相处是我的荣幸",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("这是执行到hello语句")
		},
	}

	var add = &cobra.Command{
		Use:   "add",
		Short: "计算两个数之和",
		Long:  "这是一个简单的加法计算器",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 2 {
				fmt.Println("请输出两个数字")
				return
			}
			num1, _ := strconv.Atoi(args[0])
			num2, _ := strconv.Atoi(args[1])
			fmt.Println("两数之和", num1+num2)
		},
	}

	rootCmd.AddCommand(helloCmd, add)

	err := rootCmd.Execute()
	if err != nil {
		panic(fmt.Errorf("报错，%s", err.Error()))
	}
	fmt.Println("执行成功")
}
