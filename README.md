# RSSCache

一个用来帮助你缓存 RSS 内容的小工具，有些站点的 RSS 可能只会在整点出现几分钟，之后关闭，这让 RSS 阅读器或者 BT 下载器的日志会很难看

因此我开发了这个小工具，它能在 RSS 可用的时候缓存可用的 XML 内容，并且将其提供给用户，当 RSS 更新的时候自动更新缓存

## 部署

你需要自己配置好 Golang 的基础环境：

```bash
git clone https://github.com/liueic/RSSCache.git
go mod tidy
go build
./RSSCache -url "https://example.com/rss.xml" -port 8080
```

## 开发

```bash
git clone https://github.com/liueic/RSSCache.git
go mod tidy
go run main.go
```