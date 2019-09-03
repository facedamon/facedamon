package seelog

import (
	"log"
	"os"
	"time"
)

//监控日志文件
func monitor(filepath string) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("[seelog] error:%+v", err)
		}
	}()

	var fileInfo os.FileInfo
	var err error

	for i := 1; i <= 10; i++ {
		if fileInfo, err = os.Stat(filepath); err != nil {
			log.Printf("[seelog] error:%v", err.Error())
			continue
		}
		break
	}

	offset := fileInfo.Size()
	for {
		if fileInfo, err = os.Stat(filepath); err != nil {
			log.Printf("[seelog] error:%v", err.Error())
			continue
		}
		newOffset := fileInfo.Size()
		if offset < newOffset {
			msg := make([]byte, newOffset-offset)
			file, err := os.Open(filepath)
			if err != nil {
				log.Printf("[seelog] error:%v", err.Error())
				continue
			}
			//读取二进制文件中的指定位置_
			if _, err = file.Seek(offset, 0); err != nil {
				log.Printf("[seelog] error:%v", err.Error())
			}
			if _, err = file.Read(msg); err != nil {
				log.Printf("[seelog] error:%v", err.Error())
			}

			manager.broadcast <- msg
			offset = newOffset
			file.Close()
		}
		offset = newOffset
		time.Sleep(200 * time.Millisecond)
	}
}
