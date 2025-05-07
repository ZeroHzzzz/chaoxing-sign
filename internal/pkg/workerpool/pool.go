package workerpool

import (
	"chaoxing/internal/pkg/xerr"
	"context"
	"sync"
	"time"
)

// Task 表示要执行的任务
type Task struct {
	Handler func() error  // 任务处理函数
	Timeout time.Duration // 任务超时时间
	Retries int           // 重试次数
}

// Pool 表示线程池
type Pool struct {
	workers    int           // 工作线程数
	taskQueue  chan Task     // 任务队列
	done       chan struct{} // 关闭信号
	wg         sync.WaitGroup
	maxRetries int          // 最大重试次数
	closed     bool         // 池是否已关闭
	mux        sync.RWMutex // 保护 closed 字段
}

// NewPool 创建一个新的线程池
func NewPool(workers, queueSize int) *Pool {
	p := &Pool{
		workers:    workers,
		taskQueue:  make(chan Task, queueSize),
		done:       make(chan struct{}),
		maxRetries: 3,
	}
	p.Start()
	return p
}

// Start 启动工作线程
func (p *Pool) Start() {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker()
	}
}

// worker 工作线程的主循环
func (p *Pool) worker() {
	defer p.wg.Done()

	for {
		select {
		case task := <-p.taskQueue:
			p.executeTask(task)
		case <-p.done:
			return
		}
	}
}

// executeTask 执行任务，支持重试和超时
func (p *Pool) executeTask(task Task) {
	ctx, cancel := context.WithTimeout(context.Background(), task.Timeout)
	defer cancel()

	// 确定重试次数
	maxRetries := p.maxRetries
	if task.Retries > 0 {
		maxRetries = task.Retries
	}

	// 执行任务，支持重试
	var err error
	for retry := 0; retry < maxRetries; retry++ {
		err = p.runTask(ctx, task)
		if err == nil {
			break
		}
		// 如果是超时错误，不再重试
		if err == context.DeadlineExceeded {
			break
		}
		// 重试间隔随重试次数增加
		time.Sleep(time.Second * time.Duration(retry+1))
	}
}

// runTask 在上下文约束下运行任务
func (p *Pool) runTask(ctx context.Context, task Task) error {
	done := make(chan error, 1)
	go func() {
		done <- task.Handler()
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Submit 提交任务到线程池
func (p *Pool) Submit(task Task) error {
	p.mux.RLock()
	if p.closed {
		p.mux.RUnlock()
		return xerr.PoolClosedErr
	}
	p.mux.RUnlock()

	select {
	case p.taskQueue <- task:
		return nil
	default:
		return xerr.PoolFullErr
	}
}

// Stop 停止线程池
func (p *Pool) Stop() {
	p.mux.Lock()
	if !p.closed {
		p.closed = true
		close(p.done)
		p.wg.Wait()
		close(p.taskQueue)
	}
	p.mux.Unlock()
}

// SetMaxRetries 设置最大重试次数
func (p *Pool) SetMaxRetries(retries int) {
	p.maxRetries = retries
}
