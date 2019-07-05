package seelog

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

const (
	TestLog  = "test.log"
	TestPort = 9999
)

func TestSee(t *testing.T) {
	//测试
	See(TestLog, TestPort)
	if err := os.Remove(TestLog); err != nil {
		log.Printf(err.Error())
	}
	f, err := os.Create(TestLog)
	if err != nil {
		t.Log(err.Error())
		return
	}
	for i := 1; i <= 100; i++ {
		time.Sleep(1 * time.Second)
		message := fmt.Sprintf("「模拟日志」 [%s] 第[%d]行日志\n", time.Now().String(), i)
		if _, err := f.WriteString(message); err != nil {
			log.Println(err.Error())
		}
	}
}
