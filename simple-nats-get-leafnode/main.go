package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/nats-io/nats.go"
)

// 定义节点和边的数据结构
type Node struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Info struct {
		ID   string `json:"id,omitempty"`
		IP   string `json:"ip,omitempty"`
		Port int    `json:"port"`
	} `json:"info"`
}

type Edge struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type Network struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

// 定义 NATS 响应结构
type NATSResponse struct {
	Server struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"server"`
	Data struct {
		Leafs []struct {
			Name     string `json:"name"`
			Account  string `json:"account"`
			IP       string `json:"ip"`
			Port     int    `json:"port"`
			RTT      string `json:"rtt"`
			SubCount int    `json:"subscriptions"`
		} `json:"leafs"`
	} `json:"data"`
}

type Config struct {
	Servers []string `json:"servers"`
}

func loadConfig(filePath string) (Config, error) {
	var config Config
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return config, err
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return config, err
	}

	return config, nil
}

// 向单个 NATS 服务器发送请求的函数
func requestNATS(serverURL string, subject string, wg *sync.WaitGroup, responses chan<- NATSResponse) {
	defer wg.Done()

	// 连接到 NATS 服务器
	nc, err := nats.Connect(serverURL)
	if err != nil {
		log.Printf("连接到 NATS 服务器失败: %s, 错误: %v\n", serverURL, err)
		return
	}
	defer nc.Close()

	// 发送请求
	msg, err := nc.Request(subject, []byte(""), nats.DefaultTimeout)
	if err != nil {
		log.Printf("请求失败: %s, 错误: %v\n", serverURL, err)
		return
	}

	// 解析响应数据
	var response NATSResponse
	if err := json.Unmarshal(msg.Data, &response); err != nil {
		log.Printf("解析响应失败: %v\n", err)
		return
	}

	// 将响应数据发送到 channel
	responses <- response
}

// 拼接 NATS 请求结果为 JSON 格式
func fetchNetworkData(servers []string) Network {
	subject := "$SYS.REQ.SERVER.PING.LEAFZ"

	// 使用 WaitGroup 同步 goroutines
	var wg sync.WaitGroup

	// 使用 channel 收集结果
	responses := make(chan NATSResponse, len(servers))

	// 启动 goroutines 处理每个 NATS 请求
	for _, server := range servers {
		wg.Add(1)
		go requestNATS(server, subject, &wg, responses)
	}

	// 等待所有 goroutines 完成
	wg.Wait()
	close(responses)

	// 拼接结果到 nodes 和 edges
	network := Network{}

	for response := range responses {
		// 添加主节点
		masterNode := Node{
			Name: response.Server.Name,
			Type: "master",
		}
		masterNode.Info.ID = response.Server.ID
		masterNode.Info.Port = 4222 // 假设主节点的端口是 4222
		network.Nodes = append(network.Nodes, masterNode)

		// 添加叶子节点和边
		for _, leaf := range response.Data.Leafs {
			leafNode := Node{
				Name: leaf.Name,
				Type: "node",
			}
			leafNode.Info.IP = leaf.IP
			leafNode.Info.Port = leaf.Port
			network.Nodes = append(network.Nodes, leafNode)

			// 添加边
			edge := Edge{
				Source: response.Server.Name,
				Target: leaf.Name,
			}
			network.Edges = append(network.Edges, edge)
		}
	}

	return network
}

func main() {
	config, err := loadConfig("/Users/zhaon/go/github/Go-mini-demo/simple-nats-get-leafnode/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v\n", err)
	}

	// 注册 HTTP 路由
	http.HandleFunc("/api/v1/nats/network", func(w http.ResponseWriter, r *http.Request) {
		networkData := fetchNetworkData(config.Servers)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(networkData); err != nil {
			http.Error(w, fmt.Sprintf("无法生成 JSON: %v", err), http.StatusInternalServerError)
		}
	})

	// 启动 HTTP 服务
	fmt.Println("服务启动，监听端口 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
