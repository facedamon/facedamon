package runner

import (
	"log"
	"os"
	"testing"
	"time"
)

func TestRunner_Start(t *testing.T) {
	log.Println("...开始执行任务...")

	ti := 3 * time.Second
	r := RunnerNew(ti)
	r.Add(createTask(), createTask(), createTask())
	if err := r.Start(); err != nil {
		switch err {
		case ErrTimeOut:
			log.Println(err)
			os.Exit(1)
		case ErrInterrupt:
			log.Println(err)
			os.Exit(2)
		}
	}

	log.Println("执行任务结束")
}

func createTask() func(int) {
	return func(i int) {
		log.Printf("...正在执行任务%d...", i)
		time.Sleep(time.Duration(i) * time.Second)
	}
}
