package main

import (
	"files_export/svc/checkfile"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

/*
	**

	func main() {
		info, err := checkfile.GetFileStat("file.yaml")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		fmt.Println(info)

}
*/
var mutex sync.Mutex // 用于保护并发访问metrics的锁

func main() {
	// 添加命令行参数
	configPath := flag.String("metrics.config", "config.yaml", "Path to the metrics configuration file")
	port := flag.Int("metrics.port", 9103, "Port to expose metrics on")
	flag.Parse()

	// 创建一个新的registry
	reg := prometheus.NewRegistry()

	// 定义Gauge向量
	fileCountGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "backup_file_monitor",
			Name:      "today_files",
			Help:      "The number of files matching the criteria for each job.",
		},
		[]string{"backup_name"},
	)
	reg.MustRegister(fileCountGauge)

	// 设置HTTP服务器
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		mutex.Lock() // 锁住，防止并发冲突
		defer mutex.Unlock()

		// 加载配置并检查文件
		checkLists, err := checkfile.GetFileStat(*configPath)
		if err != nil {
			log.Printf("Failed to get file stat: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// 清除旧的metrics值，避免累积
		fileCountGauge.Reset()

		// 更新metrics
		for _, cl := range *checkLists {
			fileCountGauge.WithLabelValues(cl.Name).Set(float64(cl.Count))
		}

		// 使用自定义registry处理请求
		h := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)
	})

	fmt.Printf("Starting server on port %d...\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
