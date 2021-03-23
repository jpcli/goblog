package router

import (
	"fmt"
	"goblog/handler"
	"goblog/utils/errors"
	"html/template"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

func dispWeb(r *gin.Engine) {
	r.GET("/p/:path/", handler.DispIndexPostList)
	r.GET("/", func(c *gin.Context) {
		c.Request.URL.Path = "/p/1/"
		r.HandleContext(c)
	})

	r.GET("/category/:slug/:path/", handler.DispCategoryPostList)
	r.GET("/category/:slug/", func(c *gin.Context) {
		c.Request.URL.Path = fmt.Sprintf("/category/%s/1/", c.Param("slug"))
		r.HandleContext(c)
	})

	r.GET("/tag/:slug/:path/", handler.DispTagPostList)
	r.GET("/tag/:slug/", func(c *gin.Context) {
		c.Request.URL.Path = fmt.Sprintf("/tag/%s/1/", c.Param("slug"))
		r.HandleContext(c)
	})

	r.GET("/post/:path", handler.DispPost)

	r.GET("/sitemap.xml", handler.Sitemap)
	r.GET("/page/archives.html", handler.Archives)
	r.GET("/page/declaration.html", handler.Declaration)

}

// display system view
func dispView(r *multitemplate.Renderer) {
	// 非模板嵌套
	htmls, err := filepath.Glob("./view/disp/htmls/*.html")
	if err != nil {
		panic(errors.WrapfErrorWithStack(err, "failed to load htmls in display system"))
	}
	for _, html := range htmls {
		(*r).AddFromGlob(fmt.Sprintf("disp-%s", filepath.Base(html)), html)
	}

	// 布局模板
	// layouts, err := filepath.Glob(templatesDir + "/layouts/*.html")
	// if err != nil {
	// 	panic(err.Error())
	// }

	// 嵌套的内容模板
	includes, err := filepath.Glob("./view/disp/includes/*.html")
	if err != nil {
		panic(errors.WrapfErrorWithStack(err, "failed to load includes in display system"))
	}

	// template自定义函数
	funcMap := dispFunMap()

	// 将主模板，include页面，layout子模板组合成一个完整的html页面
	for _, include := range includes {
		files := []string{}
		files = append(files, "./view/disp/frame.html", include)
		// files = append(files, layouts...)
		(*r).AddFromFilesFuncs(fmt.Sprintf("disp-%s", filepath.Base(include)), funcMap, files...)
	}
}

func dispFunMap() template.FuncMap {
	return template.FuncMap{
		"SafeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"GetPageNavUrl": func(format string, i uint32) string {
			// for the first page, delete the pageNum, for example,
			// change "/tag/slug/1/" to "tag/slug/", change "/p/%d/" to "/"
			if i == 1 {
				if format == "/p/%d/" {
					return "/"
				}
				return strings.Replace(format, "%d/", "", -1)
			}
			return fmt.Sprintf(format, i)
		},
		"TimestampToDate": func(t uint32) string {
			tm := time.Unix(int64(t), 0)
			// return tm.Format("2006年01月02日 15时04分")
			return tm.Format("2006年01月02日")
		},
		"StringToLower": func(str string) string {
			return strings.ToLower(str)
		},
	}
}
