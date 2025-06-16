package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"shorturl/config"
	"shorturl/model"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Generate(c *gin.Context, DB *gorm.DB, urltemplate string) {
	initUrl := c.PostForm("init-url")

	// 检查源连接是否是合法的连接
	if checkUrl(c, initUrl, urltemplate) {
		return
	}

	// 数据库总存在数据则直接返回对应的短链接
	if findUrl(initUrl, DB, c, urltemplate) {
		return
	}

	// 生成短链接-ID 编码法（数据库自增ID + 编码)
	// 1.先将原始网址存入数据库
	url := model.URL{Original: initUrl}
	result := DB.Create(&url) // 自动填充 url.ID 为自增主键

	// 2.获取自增ID
	if result.Error != nil {
		c.HTML(http.StatusOK, urltemplate, gin.H{
			"url_check": "存储网址失败",
		})
		return
	}

	// 3.生成短链接CODE
	code := encodeBase62(int64(url.ID))

	// 4.短链接code存入数据库
	DB.Model(&url).Update("code", code)

	// 5.生成短链接
	shortUrl := config.Domain + code

	//  短链接输出前端
	c.HTML(http.StatusOK, urltemplate, gin.H{
		"url":          initUrl,
		"shortUrl":     shortUrl,
		"hidden_class": "hidden",
	})
}

func Search(c *gin.Context, DB *gorm.DB, queryTemplate string) {
	shortUrl := c.PostForm("short-url")
	// 检查源连接是否是合法的连接
	u, err := url.Parse(shortUrl)
	if err != nil || u.Scheme == "" || u.Host == "" {
		fmt.Println("错误展示")
		c.HTML(http.StatusOK, queryTemplate, gin.H{
			"url_check": "网址不正确",
			"title":     "短链接查询",
		})
		return
	}

	// u.Path 是 "/abc123"
	path := u.Path

	// 去掉开头的斜杠 "/"
	code := strings.TrimPrefix(path, "/")

	// 根据code查询数据库中的源连接
	URL, _ := findCode(code, DB, c, queryTemplate)

	// 返回到前端的源地址与短链接
	original := URL.Original
	c.HTML(http.StatusOK, queryTemplate, gin.H{
		"url":      original,
		"shortUrl": shortUrl,
		"title":    "短链接查询",
	})
	fmt.Println(code) // 输出: abc123

}

func Redirect(c *gin.Context, DB *gorm.DB, notFoundTemplate string) {
	code := c.Param("shortcode")

	URL, _ := findCode(code, DB, c, notFoundTemplate)
	// 重定向到原始网址
	c.Redirect(http.StatusFound, URL.Original)
}

// 检查url是否格式正确
func checkUrl(c *gin.Context, initUrl string, template string) bool {
	u, err := url.Parse(initUrl)
	if err != nil || u.Scheme == "" || u.Host == "" {
		c.HTML(http.StatusOK, template, gin.H{
			"url_check": "网址不正确",
		})
		return true
	}
	return false
}

// 检查数据中的连接是否与源连接重复
func findUrl(initUrl string, DB *gorm.DB, c *gin.Context, urltemplate string) bool {
	var existing model.URL
	result := DB.Where("original = ?", initUrl).First(&existing)
	if result.Error == nil {
		// 已存在，直接返回旧的短链接
		c.HTML(http.StatusOK, urltemplate, gin.H{
			"url":          existing.Original,
			"shortUrl":     config.Domain + existing.Code,
			"hidden_class": "hidden",
		})
		return true
	}
	return false
}

// 短链接生成
func encodeBase62(id int64) string {
	var base62chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if id == 0 {
		return "0"
	}
	var result []byte
	for id > 0 {
		remainder := id % 62
		result = append([]byte{base62chars[remainder]}, result...)
		id = id / 62
	}
	return string(result)
}

func findCode(code string, DB *gorm.DB, c *gin.Context, template string) (model.URL, error) {
	var URL model.URL
	result := DB.Where("code = ?", code).First(&URL) // 查询数据库获取原始网址

	if result.Error == gorm.ErrRecordNotFound {
		c.HTML(http.StatusNotFound, template, gin.H{
			"error": "短链接未找到",
		})
		return URL, result.Error
	}
	return URL, nil
}
