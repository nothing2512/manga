package main

import (
	"fmt"
	"main/controllers"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	pth, _ := os.Executable()
	// pth = ""
	pth = strings.ReplaceAll(pth, "manga/manga", "manga/")
	c := controllers.NewInstance()

	backends := map[string]http.HandlerFunc{
		"/api/add":    c.Add,
		"/api/update": c.Update,
		"/api/list":   c.List,
		"/api/read":   c.Read,
		"/api/delete": c.Delete,
	}

	http.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
		expires := time.Unix(0, 0)
		w.Header().Set("Expires", expires.Format(http.TimeFormat))
		w.Header().Set("Pragma", "no-cache")
		if handler, ok := backends[r.URL.Path]; ok {
			handler.ServeHTTP(w, r)
			return
		}
		if strings.Contains(r.URL.Path, ".") {
			http.ServeFile(w, r, pth+"static/"+r.URL.Path)
		} else if r.URL.Path[len(r.URL.Path)-1] == '/' {
			http.ServeFile(w, r, pth+"html/"+r.URL.Path+"index.html")
		} else {
			http.ServeFile(w, r, pth+"html/"+r.URL.Path+".html")
		}
	}))

	// c.RouteView("", "test")
	// c.RouteView("/t2", "test2")

	ln, err := net.Listen("tcp", "0.0.0.0:3333")
	if err != nil {
		panic(err)
	}

	fmt.Println("Listening On 0.0.0.0:3333")
	if err = http.Serve(ln, nil); err != nil {
		panic(err)
	}
}
