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

func init() {
	openOrCreateDir("vscd")
}
func main() {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	port := fmt.Sprintf(":%v", listener.Addr().(*net.TCPAddr).Port)

	listener.Close()
	app := gin.Default()
	app.MaxMultipartMemory = 4 *1024*1024*1024
	app.Static("/mirrors", "./vscd")
	app.LoadHTMLGlob("view/*")
	app.GET("/", func(c *gin.Context) {
		p := c.Query("p")
		var cd string

		if p == "vscd" || p == "" {
			p = "vscd"
			cd = ""
		} else if strings.Index(p, "vscd-") == 0 {
			lg := strings.LastIndex(p, "-")
			cd = p[0:lg]
		} else {
			p = "vscd"
		}

		list := flist("./" + strings.ReplaceAll(p, "-", "/"))

		fl := make([]vscdfile, 0)

		for _, v := range list {
			fl = append(fl,
				newvscdfile(v.Name(), filetype(v.IsDir()), p, v.Size()),
			)
		}
		fmt.Println(fl)
		c.HTML(http.StatusOK, "index.html", gin.H{
			"cd":   cd,
			"list": fl,
		})
	})
	// app.POST("/upload", func(c *gin.Context) {
	// 	file, _ := c.FormFile("file")
	// 	c.SaveUploadedFile(file, "./vscd/"+file.Filename)
	// 	c.Status(http.StatusOK)
	// })
	app.Run(port)

}

type vscdfile struct {
	Name  string
	Ftype string
	Size  int64
	Link  string
}

func newvscdfile(n, t, l string, s int64) vscdfile {
	f := vscdfile{
		Name:  n,
		Ftype: t,
		Size:  s,
	}
	if t == "file" {
		f.Link = strings.Replace(l, "vscd", "mirrors", 1)
		f.Link = strings.ReplaceAll(f.Link, "-", "/") + "/" + n
	} else {
		f.Link = strings.ReplaceAll(l, "/", "-") + "-" + n
	}

	return f
}

func filetype(isdir bool) string {
	if isdir {
		return "dir"
	}
	return "file"

}

func flist(str string) []os.FileInfo {

	dir_list, e := ioutil.ReadDir(str)
	if e != nil {
		fmt.Println("read dir error")
		return nil
	}
	return dir_list
}

func openOrCreateDir(path string) {

	_, err := os.Stat(path)
	if err == nil {
	} else {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
