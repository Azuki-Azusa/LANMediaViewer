package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	showIP()
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	loadTemplates(r)
	route(r)
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}

// loat static file
func loadTemplates(r *gin.Engine) {
	r.Delims("{{{", "}}}")
	r.LoadHTMLGlob("templates/*")
	r.StaticFS("/public", http.Dir("./public"))
}

// route table
func route(r *gin.Engine) {
	r.GET("", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.GET("/video", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"video": getObject("public/video"),
		})
	})

	r.GET("/image", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"image": getObject("public/image"),
		})
	})

	r.GET("/music", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"music": getObject("public/music"),
		})
	})
}

// show ip of this pc
func showIP() {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	temp := "192.168.1.1"
	fmt.Println("Your IP addresses which start with '192.'：")
	for _, address := range addrs {

		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil && ipnet.IP.String()[:3] == "192" {
				temp = ipnet.IP.String()
				fmt.Println(ipnet.IP.String())
			}

		}
	}

	fmt.Println("Server started, please connect to the [ip:8080] via browsers on your devices(PC, Phone, PS4, etc.)")
	fmt.Println("For example: " + temp + ":8080")
}

// get files tree object by root path
func getObject(root string) gin.H {
	dir_list, e := ioutil.ReadDir(root)
	label := strings.Split(root, "/")[len(strings.Split(root, "/"))-1]
	path := root
	var children []gin.H
	if e != nil {
		fmt.Println("read dir error")
		return gin.H{}
	}
	for _, v := range dir_list {
		if v.IsDir() {
			children = append(children, getObject(root+"/"+v.Name()))
		} else {
			children = append(children, gin.H{
				"path":  path + "/" + v.Name(),
				"label": v.Name(),
			})
		}
	}
	return gin.H{
		"path":     path,
		"label":    label,
		"children": children,
	}
}
