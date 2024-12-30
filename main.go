package main

import (
	"flag"
	"fmt"

	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func init() {
	openOrCreateDir("vscd")
}

var server_port string

func main() {
	flag.StringVar(&server_port, "port", "", "server port")
	flag.Parse()
	if server_port == "" {
		listener, err := net.Listen("tcp", ":0")
		if err != nil {
			panic(err)
		}
		listener.Close()
		server_port = fmt.Sprintf("%v", listener.Addr().(*net.TCPAddr).Port)
	}

	
	app := gin.Default()
	app.MaxMultipartMemory = 4 * 1024 * 1024 * 1024
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
				newvscdfile(v, p),
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
	app.Run(":"+server_port)

}

type vscdfile struct {
	Name  string
	Ftype string
	Size  int64
	Link  string
}

func newvscdfile(tfile os.DirEntry, l string) vscdfile {
	
	finfo, err := tfile.Info()
	if err != nil {
		return vscdfile{}
	}
	f := vscdfile{
		Name:  tfile.Name(),
		Ftype: filetype(tfile.IsDir()),
		Size:  finfo.Size(),
	}
	if f.Ftype == "file" {
		f.Link = strings.Replace(l, "vscd", "mirrors", 1)
		f.Link = strings.ReplaceAll(f.Link, "-", "/") + "/" + f.Name
	} else {
		f.Link = strings.ReplaceAll(l, "/", "-") + "-" + f.Name
	}

	return f
}

func filetype(isdir bool) string {
	if isdir {
		return "dir"
	}
	return "file"

}

func flist(str string) []os.DirEntry {

	dir_list, e := os.ReadDir(str)
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
