package seelog

import (
	"fmt"
	"golang.org/x/net/websocket"
	"html/template"
	"log"
	"net/http"
	"path"
	"runtime"
	"time"
)

/**
start httpserver
*/
func server(port int) {
	/**
	catch panic
	*/
	defer func() {
		if err := recover(); err != nil {
			log.Printf("[seelog] error:%+v", err)
		}
	}()

	http.Handle("/ws", websocket.Handler(func(ws *websocket.Conn) {
		client := &client{time.Now().String(), ws, make(chan []byte, 1024)}
		manager.register <- client
		client.write()
	}))

	http.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
		//skip如果是0，返回当前调用Caller函数的函数名、文件、程序指针PC，1是上一层函数，以此类推
		_, currentfile, _, _ := runtime.Caller(0)
		filename := path.Join(path.Dir(currentfile), "index.html")
		t, err := template.ParseFiles(filename)
		if err != nil {
			log.Println(err)
		}
		t.Execute(writer, nil)
	})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
