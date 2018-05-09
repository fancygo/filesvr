package main

import (
	"fmt"
	"github.com/FancyGo/svrreg"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

var mux map[string]func(http.ResponseWriter, *http.Request)

type Myhandler struct{}
type home struct {
	Title string
}

const (
	Template_Dir = "./view/"
	File_Dir     = "./file/"
	Js_Dir       = "./js/"
	Css_Dir      = "./css/"
)

func main() {
	//filesvr服务
	server := http.Server{
		Addr:        ":" + strconv.Itoa(FILE_SVR_PORT),
		Handler:     &Myhandler{},
		ReadTimeout: 100 * time.Second,
	}
	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["/"] = index
	mux["/upload"] = upload
	mux["/js"] = jsFile
	mux["/file"] = commonFile
	mux["/css"] = cssFile
	mux["/health"] = health
	fmt.Println("Hello, this is FancyGo filesvr!")
	go server.ListenAndServe()

	//服务注册
	regCfg := &svrreg.RegCfg{
		LocalSvrID:       FILE_SVR_ID,
		LocalSvrName:     FILE_SVR_NAME,
		LocalSvrDNS:      FILE_SVR_DNS,
		LocalSvrPort:     FILE_SVR_PORT,
		CoreSvrDNS:       CORE_SVR_DNS,
		CoreSvrPort:      CORE_SVR_PORT,
		SvrCheckTimeout:  SVR_CHECK_TIMEOUT,
		SvrCheckInterval: SVR_CHECK_INTERVAL,
	}

	register := svrreg.NewRegConsul()
	if ok := register.SvrRegInit(regCfg); !ok {
		return
	}

	if ok := register.RegSvr(); !ok {
		return
	}
	/*
		if ok := svrreg.Reginit(regConsul, regCfg); !ok {
			return
		}

		if ok := svrreg.Reg(regConsul); !ok {
			return
		}
	*/

	//设置sigint信号
	close := make(chan os.Signal, 1)
	signal.Notify(close, os.Interrupt, os.Kill)
	<-close

	if ok := register.UnregSvr(); !ok {
		return
	}
	fmt.Println("Bye, FancyGo filesvr close")
}

func (*Myhandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.String()]; ok {
		h(w, r)
		return
	}
	if ok, _ := regexp.MatchString("/css/", r.URL.String()); ok {
		http.StripPrefix("/css/", http.FileServer(http.Dir("./css/"))).ServeHTTP(w, r)
	} else if ok, _ := regexp.MatchString("/js/", r.URL.String()); ok {
		http.StripPrefix("/js/", http.FileServer(http.Dir("./js/"))).ServeHTTP(w, r)
	} else if ok, _ := regexp.MatchString("/file/", r.URL.String()); ok {
		http.StripPrefix("/file/", http.FileServer(http.Dir("./file/"))).ServeHTTP(w, r)
	}
}

func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles(Template_Dir + "file.html")
		t.Execute(w, "上传文件")
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Printf("upload err = %v", err)
			fmt.Fprintf(w, "%v", "上传错误")
			return
		}
		fileext := filepath.Ext(handler.Filename)
		if check(fileext) == false {
			fmt.Printf("upload typ err = %v", err)
			fmt.Fprintf(w, "%v", "类型错误")
			return
		}
		f, _ := os.OpenFile(File_Dir+handler.Filename, os.O_CREATE|os.O_WRONLY, 0660)
		_, err = io.Copy(f, file)
		if err != nil {
			fmt.Printf("upload fail err = %v", err)
			fmt.Fprintf(w, "%v", "上传失败")
			return
		}
		fmt.Fprintf(w, "%v", handler.Filename+"上传完成!")
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", "check health ok!")
}

func index(w http.ResponseWriter, r *http.Request) {
	title := home{Title: "首页"}
	t, _ := template.ParseFiles(Template_Dir + "index.html")
	t.Execute(w, title)
}

func commonFile(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/file/", http.FileServer(http.Dir("./file/"))).ServeHTTP(w, r)
}

func jsFile(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/js/", http.FileServer(http.Dir("./js/"))).ServeHTTP(w, r)
}

func cssFile(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/css/", http.FileServer(http.Dir("./css/"))).ServeHTTP(w, r)
}

func check(name string) bool {
	ext := []string{".js", ".exe"}

	for _, v := range ext {
		if v == name {
			return false
		}
	}
	return true
}
