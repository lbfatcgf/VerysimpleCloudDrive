package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"github.com/gin-gonic/gin"
)
func init(){
	openOrCreateDir("vscd")
}
func main() {
	app := gin.Default()
	app.Static("/mirrors", "./vscd")
	app.LoadHTMLGlob("view/*")
	app.GET("/show/:p", func(c *gin.Context) {
		p := c.Param("p")
		var cd string

		if p == "vscd" {
			cd = ""
		} else if strings.Index(p, "vscd/") == 0 {
			lg := strings.LastIndex(p, "/")
			cd = p[0:lg]
		} else {
			c.Redirect(http.StatusMovedPermanently, "show/vscd")
			return
		}

		list := flist("./" + p)

		fl := make([]vscdfile, 0)

		for _, v := range list {
			fl = append(fl,
				newvscdfile(v.Name(), filetype(v.IsDir()), fmt.Sprintf("%v%v", p, v.Name()), v.Size()),
			)
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"cd":   cd,
			"list": fl,
		})
	})

	app.GET("/",func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "show/vscd")
	})
	app.Run()
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
	f.Link = strings.Replace(l, "vscd", "mirrors", 1)
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
