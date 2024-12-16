package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

// 这个案例主要是利用cobra的机制，写出来version的版本，并且知道git信息和打包日期，这些参数是从外面注入进来的，也就是shell脚本
// 这里面注意的细节是有一个version的命令在root命令里面
// 那么这个值具体来说是怎么设计传入进来的，
// 我发现是变异的时候带来的，就是 go build -ldflags "-X main.version=1.0 -X main.commit=$(git rev-parse HEAD) -X main.buildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o version-demo

var version, commit, buildDate = "0.1", "unknown", "unknown"

func main() {
	rootCmd := &cobra.Command{}

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "输出版本的信息",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("版本: %s \nGit标签: %s\n打包日期: %s\n", version, commit, buildDate)
		},
	}

	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
