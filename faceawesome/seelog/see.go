package seelog

import (
	"log"
	"time"
)

func See(filepath string, port int) {
	if !checkParam(filepath, port) {
		return
	}

	//start socketmanager
	go manager.start()
	go monitor(filepath)
	// start httpserver
	go server(port)
	time.Sleep(200 * time.Millisecond)
}

//参数验证
func checkParam(filepath string, port int) bool {
	if filepath == "" {
		log.Println("filepath 不能为空")
		return false
	}
	if port == 0 {
		log.Println("port 不能为空")
		return false
	}
	return true
}
