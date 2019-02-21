package main

import (
	"flag"
	"net/http"
	"net/http/httputil"
	"net/url"
	"proxy/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var Env string
var Port int

func main() {
	flag.IntVar(&Port, "port", 8892, "http server port")
	flag.StringVar(&Env, "env", "development", "go run with environment")
	flag.Parse()

	engine := gin.New()

	if Env == "development" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine.Use(CorsMiddleware())
	engine.Use(ProxyRouter())
	// engine.Use(CorsMiddleware())

	engine.GET("/", func(c *gin.Context) {
		c.String(200, `index`)
	})

	engine.Run(":" + strconv.Itoa(Port))
}

func ProxyRouter() gin.HandlerFunc {
	return func(c *gin.Context) {

		WithHeader(c)

	}
}

var simpleHostProxy = httputil.ReverseProxy{
	Director: func(req *http.Request) {
		var pathArr []string
		path := req.URL.EscapedPath()

		if strings.Contains(path, "%2F%2F") {
			decodeurl, _ := url.QueryUnescape(path)
			arr := strings.Split(decodeurl, "/")

			for _, value := range arr {
				if value != "" {
					pathArr = append(pathArr, value)
				}
			}

			protocol := strings.Split(pathArr[0], ":")[0]
			host := pathArr[1]
			url := strings.Join(pathArr[2:], "/")

			utils.Log("protocol:", protocol, "\nhost:", host, "\nurl:", url)

			req.URL.Scheme = protocol
			req.URL.Host = host
			req.Host = host
			if strings.Contains(url, "?") {
				splistLast := strings.Index(url, "?")
				req.URL.Path = "/" + url[0:splistLast]
				req.URL.ForceQuery = true
				req.URL.RawQuery = url[splistLast+1:]
			} else {
				req.URL.Path = "/" + url[0:]
			}

			utils.Log(req.URL.String())

			// Opaque     string    // encoded opaque data
			// User       *Userinfo // username and password information
			// Host       string    // host or host:port
			// Path       string    // path (relative paths may omit leading slash)
			// RawPath    string    // encoded path hint (see EscapedPath method)
			// ForceQuery bool      // append a query ('?') even if RawQuery is empty
			// RawQuery   string    // encoded query values, without '?'
			// Fragment   string    // fragment for references, without '#'

		}

	},
}

func WithHeader(ctx *gin.Context) {
	simpleHostProxy.ServeHTTP(ctx.Writer, ctx.Request)
}

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		//var filterHost = [...]string{"*.xxx.com"}
		// filterHost 做过滤器，防止不合法的域名访问
		// var isAccess = false
		// for _, v := range filterHost {
		// 	match, _ := regexp.MatchString(v, origin)
		// 	if match {
		// 		isAccess = true
		// 	}
		// }
		// fmt.Println(origin, isAccess)
		// if isAccess {
		// 核心处理方式
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		c.Header("Access-Control-Allow-Methods", "GET, POST")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Set("content-type", "application/json")
		// }
		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}

		c.Next()
	}
}
