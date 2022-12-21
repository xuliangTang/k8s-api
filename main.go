package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/apps/v1"
	"k8sapi/core"
	"k8sapi/deployment"
	"k8sapi/lib"
	"log"
	"net/http"
)

type DepHandler struct{}

func (this *DepHandler) OnAdd(obj interface{}) {}
func (this *DepHandler) OnUpdate(oldObj, newObj interface{}) {
	if dep, ok := newObj.(*v1.Deployment); ok {
		fmt.Println(dep.Name)
	}
}
func (this *DepHandler) OnDelete(obj interface{}) {}

func main() {
	r := gin.New()
	r.Static("/static", "./static")
	r.LoadHTMLGlob("html/**/*")

	r.Use(func(context *gin.Context) {
		defer func() {
			if e := recover(); e != nil {
				log.Println(e)
				context.AbortWithStatusJSON(500, gin.H{"message": e})
			}
		}()
		context.Next()
	})

	deployment.RegHandlers(r)

	r.GET("/deployments", func(c *gin.Context) {
		c.HTML(http.StatusOK, "list.html",
			lib.DataBuilder().
				SetTitle("deployment列表").
				SetData("DepList", deployment.ListAll("default")))
	})

	r.GET("/deployment/:name", func(c *gin.Context) {
		c.HTML(http.StatusOK, "detail.html",
			lib.DataBuilder().
				SetTitle(c.Param("name")+"详情").
				SetData("DepDetail", deployment.DepDetail("default", c.Param("name"))))
	})

	r.GET("/create/deployment", func(c *gin.Context) {
		c.HTML(http.StatusOK, "deployment_create.html",
			lib.DataBuilder().
				SetTitle("deployment创建"))
	})

	core.InitDeployments()

	r.Run(":8080")
}
