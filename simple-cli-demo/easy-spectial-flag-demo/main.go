package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
	"time"
)

// 这是使用cobra的特定flag, 比如 --name ning --birthday 11.08
// 为了使用这样制定的flag，我们用到了cobra的rootCmd.Flags().StringVar()
// 同时我们也可以制定必填项，利用cobra.MakeFlagRequired
// 最后，我们执行 rootCmd.Execute()

var name, birthday = "", ""

func main() {
	rootCmd := &cobra.Command{
		Use:   "calculate birthday",
		Short: "这是计算距离生日还有多少天",
		Run:   run,
	}

	// 制定特定的flag
	rootCmd.Flags().StringVar(&name, "name", "askning", "这是输入名字")
	rootCmd.Flags().StringVar(&birthday, "birthday", "1.1", "这是输入生日")

	_ = rootCmd.MarkFlagRequired("birthday")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

// 这里输入现在距离生日有多少天，首先解析月和日，然后计算今年到这一天的月和日还有几天，如果是负数，则计算明天具体这一天的日期
func run(cmd *cobra.Command, args []string) {
	month, day, err := getMonthAndDay(birthday)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}

	// 计算距离生日的天数
	daysUntil := calculateDaysUntilBirthday(month, day)

	// 输出结果
	fmt.Printf("你好，%s！距离你的生日还有 %d 天。\n", name, daysUntil)
}

func getMonthAndDay(birthday string) (int, int, error) {
	// 按 "." 分隔字符串
	parts := strings.Split(birthday, ".")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("生日格式错误，正确格式为 MM.DD，例如 11.08")
	}

	// 分割后的数字转为整数，并且判断月份在1~12区间，判断天在1~31区间
	month, err := strconv.Atoi(parts[0])
	if err != nil || month < 1 || month > 12 {
		return 0, 0, fmt.Errorf("无效的月份: %s", parts[0])
	}

	day, err := strconv.Atoi(parts[1])
	if err != nil || day < 1 || day > 31 {
		return 0, 0, fmt.Errorf("无效的日期: %s", parts[1])
	}

	return month, day, nil
}

func calculateDaysUntilBirthday(month, day int) int {
	now := time.Now()
	currentYear := now.Year()

	// 今年的生日日期
	birthdayThisYear := time.Date(currentYear, time.Month(month), day, 0, 0, 0, 0, time.Local)

	var daysUntil int
	if birthdayThisYear.Before(now) {
		// 如果今年的生日已经过去，计算明年的生日
		birthdayNextYear := time.Date(currentYear+1, time.Month(month), day, 0, 0, 0, 0, time.Local)
		daysUntil = int(birthdayNextYear.Sub(now).Hours() / 24)
	} else {
		// 如果今年的生日还没到
		daysUntil = int(birthdayThisYear.Sub(now).Hours() / 24)
	}
	return daysUntil
}
