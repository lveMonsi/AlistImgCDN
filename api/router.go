package api

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	router           = gin.Default()
	targetUrl string = os.Getenv("url") // alist域名
	expire    int    = 604800           // 过期时间(秒) Vercel云函数有内存释放，该选项无效

	expireTimeMap = make(map[string]string)
	fileMap       = make(map[string][]byte)
)

func init() {
	if strings.HasSuffix(targetUrl, "/") {
		targetUrl = targetUrl[:len(targetUrl)-1]
	}

	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "服务运行正常~")
	})

	router.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "Hello",
		})
	})

	// router路由
	router.GET("/img/*path", func(ctx *gin.Context) {
		path := ctx.Param("path")
		ok := cacheImg(path)
		if !ok {
			ctx.Status(http.StatusNotFound)
			return
		}
		// ctx.Request.Header.Set("Content-Type", "image/png")

		ctx.Data(http.StatusOK, "image/png", fileMap[path])
	})

	router.GET("/cachelist", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, expireTimeMap)
	})
}

// 缓存机制
func cacheImg(path string) bool {
	t, _ := time.Parse("2006-01-02 15:04", expireTimeMap[path])
	flag := t.After(time.Now())

	_, ok := expireTimeMap[path]
	if !ok {
		fmt.Println("新建缓存", path)
		expireTimeMap[path] = time.Now().Add(time.Duration(expire * 1000000000)).Format("2006-01-02 15:04")
		flag = true
	}

	if flag {
		// 创建一个自定义的 HTTP 客户端
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, // 跳过证书验证
				},
			},
		}

		resp, err := client.Get(targetUrl + path)

		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		if resp.Header.Get("Content-type") == "application/json; charset=utf-8" {
			fmt.Println("未找到该图片")
			delete(expireTimeMap, path)
			return false
		}

		fmt.Println("更新图片: ", targetUrl+path)

		tmp, _ := io.ReadAll(resp.Body)
		fileMap[path] = []byte(tmp)

		// 现在时间往后7天
		expireTimeMap[path] = time.Now().AddDate(0, 0, 7).Format("2006-01-02 15:04:05")
		fmt.Println("下次过期时间：", expireTimeMap[path])

	}
	return true
}

func Listen(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}
