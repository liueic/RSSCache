# RSSCache

一个用来帮助你缓存 RSS 内容的小工具

## 部署

在 Release 页面下载构建好的二进制文件

```bash
./RSSCache -url "https://example.com/rss.xml" -port 8080
```

## 开发

```bash
git clone https://github.com/liueic/RSSCache.git
go mod tidy
go run main.go
```