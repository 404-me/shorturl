package main

import (
	"net/http"
	"shorturl/config"
	"shorturl/database"
	"shorturl/handler"
	"shorturl/model"

	"github.com/gin-gonic/gin"
)

func main() {
	DB := database.InitializeDatabase()

	DB.AutoMigrate(&model.URL{})

	r := gin.Default()
	// 加载HTML模板
	r.LoadHTMLGlob("templates/*")

	// 展示首页
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "短链接",
		})
	})
	r.GET("/shorten", func(c *gin.Context) {
		// 复用该html的模板加载到前端中展示新的html
		c.HTML(http.StatusOK, "result.html", gin.H{
			"title": "短链接生成",
		})
	})
	r.GET("/query", func(c *gin.Context) {
		c.HTML(http.StatusOK, "query.html", gin.H{
			"title": "短链接查询",
		})
	})
	// 生成短链接
	r.POST("/generate", func(c *gin.Context) {
		handler.Generate(c, DB, "result.html")
	})
	// 查询短链接
	r.POST("/:search", func(c *gin.Context) {
		handler.Search(c, DB, "query.html")
	})

	// 重定向到源网址
	r.GET("/:shortcode", func(c *gin.Context) {
		handler.Redirect(c, DB, "notfound.html")
	})
	// 启动服务器
	r.Run(config.ListenAddr)
}
