package runner

import (
	"errors"
	"os"
	"os/signal"
	"time"
)

var ErrTimeOut = errors.New("执行超时错误")
var ErrInterrupt = errors.New("执行中断错误")

type runner struct {
	tasks     []func(int)      `任务列表，具体任务func`
	complete  chan error       `是否完成，error类型可以输出执行中断错误类型`
	timeout   <-chan time.Time `超时时间, 单向接收通道`
	interrupt chan os.Signal   `中断信号`
}

func RunnerNew(d time.Duration) *runner {
	// complete 无缓冲通道，控制整个程序的终止，要让main等待
	//interrupt 一个缓冲，确保至少有一个os中断信号被接收
	return &runner{complete: make(chan error), timeout: time.After(d), interrupt: make(chan os.Signal, 1)}
}

func (r *runner) Add(t ...func(int)) {
	r.tasks = append(r.tasks, t...)
}

func (r *runner) isInterrupt() bool {
	select {
	case <-r.interrupt:
		signal.Stop(r.interrupt)
		return true
	default:
		return false
	}
}

func (r *runner) Run() error {
	for id, task := range r.tasks {
		if r.isInterrupt() {
			return ErrInterrupt
		}
		task(id)
	}
	return nil
}

func (r *runner) Start() error {
	signal.Notify(r.interrupt, os.Interrupt)

	go func() { r.complete <- r.Run() }()

	select {
	// 阻塞等待
	case err := <-r.complete:
		return err
	case <-r.timeout:
		return ErrTimeOut
	}
}
