package main

import (
	"backup_files_exporter/svc/checkfile"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	mutex       sync.Mutex // 用于保护并发访问metrics的锁
	checkLists  *[]checkfile.CheckList
	lastUpdated time.Time
	//expiration  = 5 * time.Minute // 缓存过期时间设置为5分钟
)

func main() {
	// 添加命令行参数
	configPath := flag.String("metrics.config", "config.yaml", "Path to the metrics configuration file")
	port := flag.Int("metrics.port", 9103, "Port to expose metrics on")
	interval := flag.Int("metrics.interval", 30, "Interval to check metrics in minute")
	flag.Parse()

	// 创建一个新的registry
	reg := prometheus.NewRegistry()

	// 定义Gauge向量
	fileCountGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "backup_file_monitor",
			Name:      "today_files",
			Help:      "The number of DB backup files.",
		},
		[]string{"backup_name"},
	)
	fileSizeGauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "backup_file_monitor",
			Name:      "today_files_size",
			Help:      "The size of DB backup files.",
		},
		[]string{"backup_name"},
	)
	reg.MustRegister(fileCountGauge)
	reg.MustRegister(fileSizeGauge)
	// 获取缓存数据过期时间
	expiration := time.Duration(*interval) * time.Minute
	// 设置HTTP服务器
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		mutex.Lock() // 锁住，防止并发冲突
		defer mutex.Unlock()

		// 检查是否需要更新数据,防止频繁更新备份文件状态
		fmt.Println(lastUpdated)
		fmt.Println(time.Since(lastUpdated))
		fmt.Println(expiration)
		if time.Since(lastUpdated) > expiration {
			var err error
			checkLists, err = checkfile.GetFileStat(*configPath)
			if err != nil {
				log.Printf("Failed to get file stat: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			lastUpdated = time.Now()
		}
		// 清除旧的metrics值，避免累积
		fileCountGauge.Reset()

		// 更新metrics
		for _, cl := range *checkLists {
			fileCountGauge.WithLabelValues(cl.Name).Set(float64(cl.Count))
			fileSizeGauge.WithLabelValues(cl.Name).Set(float64(cl.Size))
		}

		// 使用自定义registry处理请求
		h := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)
	})

	fmt.Printf("Starting server on port %d...\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
