package logger

import (
	"fmt"
	"net/http"
)

// LogStatus 返回是否在记录日志
func (o *Logger) LogStatus() bool {
	return o.logged
}

// StartServer 启动一个web管理面板
func (o *Logger) StartServer(addr string) {
	srv := http.NewServeMux()
	srv.HandleFunc("/", o.root)
	go http.ListenAndServe(addr, srv)
}
func (o *Logger) root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	r.ParseForm()
	if r.FormValue("logged") == "false" {
		o.logged = false
	} else if r.FormValue("logged") == "true" {
		o.logged = true
	}

	if o.logged {
		fmt.Fprintln(w, `状态: 开启 <a href="/?logged=false">关闭</a>`)
	} else {
		fmt.Fprintln(w, `状态: 关闭 <a href="/?logged=true">开启</a>`)
	}
}
