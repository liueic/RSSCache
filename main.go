package main

import (
	"flag" // 新增
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	rssCache   []byte
	cacheMutex sync.RWMutex

	rssURL string // 不再赋初值
	port   string // 新增端口参数
)

func initLogger() {
	// 创建或追加日志文件
	logFile, err := os.OpenFile("rss-proxy.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// 设置日志输出到文件和标准输出
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func fetchRSS() {
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", rssURL, nil)
	if err != nil {
		log.Printf("Failed to create request: %v\n", err)
		return
	}

	// 设置自定义 User-Agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (RSSFetcher)")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Fetch failed: %v\n", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	contentType := resp.Header.Get("Content-Type")
	if !(strings.Contains(contentType, "application/rss+xml") ||
		strings.Contains(contentType, "application/xml") ||
		strings.Contains(contentType, "text/xml")) {
		log.Printf("Invalid Content-Type: %s — skipping update\n", contentType)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response: %v\n", err)
		return
	}

	cacheMutex.Lock()
	rssCache = body
	cacheMutex.Unlock()
	log.Println("RSS cache updated successfully.")
}

func startFetcher() {
	fetchRSS() // 刚开始的时候就进行刷新
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		fetchRSS()
	}
}

func main() {
	// 命令行参数解析
	flag.StringVar(&rssURL, "url", "", "RSS 源链接 (必填)")
	flag.StringVar(&port, "port", "23451", "监听端口")
	flag.Parse()

	if rssURL == "" {
		log.Fatalln("请使用 -url 参数指定 RSS 源链接")
	}

	initLogger()

	// 启动定时抓取
	go startFetcher()

	// 初始化 Gin 服务
	r := gin.Default()

	r.GET("/rss", func(c *gin.Context) {
		cacheMutex.RLock()
		defer cacheMutex.RUnlock()

		if len(rssCache) == 0 {
			c.String(http.StatusServiceUnavailable, "RSS not available yet")
			return
		}

		c.Data(http.StatusOK, "application/xml", rssCache)
	})

	log.Printf("RSS proxy server started at :%s\n", port)
	err := r.Run(":" + port)
	if err != nil {
		return
	}
}
