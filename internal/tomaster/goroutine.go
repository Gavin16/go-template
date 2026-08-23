package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"
)

// 模拟机场安检
// 安检每条通道的执行流程为: 身份检查 -> 人身检查 -> X光检查
// 采用三条通道并发进行检查
// 每个检查环节耗时定义为 idCheckTmCost: 60, bodyCheckTmCost: 120, xRayCheckTmCost: 180
// 程序为每个通道的每个环节单独起 go routine

const (
	idCheckTmCost   = 60
	bodyCheckTmCost = 120
	xRayCheckTmCost = 180
)

func idCheck(id string) int {
	time.Sleep(idCheckTmCost * time.Millisecond)
	//fmt.Printf("%s id check ok\n", id)
	return idCheckTmCost
}

func bodyCheck(id string) int {
	time.Sleep(bodyCheckTmCost * time.Millisecond)
	//fmt.Printf("%s, body check ok\n", id)
	return bodyCheckTmCost
}

func xRayCheck(id string) int {
	time.Sleep(xRayCheckTmCost * time.Millisecond)
	//fmt.Printf("%s xRay check ok\n", id)
	return xRayCheckTmCost
}

func start(id string, f func(string) int, next chan<- struct{}) (
	chan<- struct{}, chan<- struct{}, <-chan int) {
	queue := make(chan struct{}, 50)
	quit := make(chan struct{})
	result := make(chan int)

	go func() {
		total := 0
		for {
			select {
			case <-quit:
				result <- total
				return
			case v := <-queue:
				total += f(id)
				if next != nil {
					next <- v
				}
			}
		}
	}()
	return queue, quit, result
}

func AirportSecurityCheckChannel(id string, queue chan struct{}) {
	go func(id string) {
		fmt.Printf("goroutine-%s : airportSecurityCheckChannel is ready...\n", id)

		queue3, quit3, result3 := start(id, xRayCheck, nil)
		queue2, quit2, result2 := start(id, bodyCheck, queue3)
		queue1, quit1, result1 := start(id, idCheck, queue2)

		for {
			select {
			case v, ok := <-queue:
				if !ok {
					close(quit1)
					close(quit2)
					close(quit3)
					total := max(<-result1, <-result2, <-result3)
					fmt.Printf("goroutine-%s : airportSecurityCheckChannel time cost: %d\n", id, total)
					fmt.Printf("goroutine-%s : airportSecurityCheckChannel closed\n", id)
					return
				}
				queue1 <- v
			}
		}
	}(id)
}

// -- goroutine调度原理
func dummy() {
	add(1, 2)
}

func add(x, y int) int {
	return x + y
}

func deadLoop() {
	for {
		dummy()
	}
}

// -- goroutine 退出模式中的 join 模式
func worker(args ...interface{}) {
	if len(args) == 0 {
		return
	}

	interval, ok := args[0].(int)
	if !ok {
		return
	}

	time.Sleep(time.Second * time.Duration(interval))
}

func spawnGroup(n int, f func(args ...interface{}), args ...interface{}) chan struct{} {
	c := make(chan struct{})
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			name := fmt.Sprintf("worker-%d", i)
			f(args...)
			println(name, "done")
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
		c <- struct{}{}
	}()
	return c
}

// main goroutine 通知关闭 goroutine group
func spawnGroup2(n int, f func(args ...interface{}), args ...interface{}) chan struct{} {
	quit := make(chan struct{})
	job := make(chan int, 1)
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		job <- i // 这里通道写入任意数据
		go func(i int) {
			defer wg.Done() // 确保所有Done 都在goroutine 退出前被执行
			name := fmt.Sprintf("worker-%d", i)
			for {
				_, ok := <-job
				if !ok {
					println(name, "done")
					return
				}
				f(args)
			}
		}(i)
	}

	go func() {
		<-quit
		close(job)
		wg.Wait()
		quit <- struct{}{}
	}()
	return quit
}

// GracefulShutdown 关闭goroutine组, 通过接口实现goroutine 优雅退出
type GracefulShutdown interface {
	Shutdown(waitTimeout time.Duration) error
}

type ShutdownFunc func(time.Duration) error

func (f ShutdownFunc) Shutdown(waitTimeout time.Duration) error {
	return f(waitTimeout)
}

func shutdownMaker(processTime time.Duration) func(time.Duration) error {
	return func(duration time.Duration) error {
		time.Sleep(processTime)
		return nil
	}
}

// 并发的关闭goroutine
func concurrentShutdown(waitTimeout time.Duration, ds []time.Duration, shutdowns ...GracefulShutdown) error {
	c := make(chan struct{})
	go func() {
		var wg sync.WaitGroup
		for i, shutdown := range shutdowns {
			wg.Add(1)
			shutdownWaitTm := ds[i]
			go func(s GracefulShutdown) {
				defer wg.Done()
				_ = s.Shutdown(shutdownWaitTm)
			}(shutdown)
		}
		wg.Wait()
		c <- struct{}{}
	}()

	timer := time.NewTimer(waitTimeout)
	defer timer.Stop()

	select {
	case <-c:
		close(c)
		return nil
	case <-timer.C:
		return errors.New("wait timeout")
	}
}

// -- 管道模式
// 对生成的序列过滤出偶数，并取平方
func newNumGenerator(start, count int) <-chan int {
	c := make(chan int)

	go func() {
		for i := 0; i < count; i++ {
			c <- i
		}
		close(c)
	}()
	return c
}

func filterOdd(in int) (int, bool) {
	if in%2 != 0 {
		return 0, false
	}
	return in, true
}

func square(in int) (int, bool) {
	return in * in, true
}

func spawn2(f func(int) (int, bool), in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for v := range in {
			r, ok := f(v)
			if ok {
				out <- r
			}
		}
		close(out)
	}()
	return out
}

// -- 管道的扩展模式: 扇出模式 和 扇入模式
func spawnGroup3(name string, num int, f func(int) (int, bool), in <-chan int) <-chan int {
	var outSlice []chan int
	// 读取输入通道数据，扇出到 outSlice 切片
	for i := 0; i < num; i++ {
		out := make(chan int)
		go func(i int) {
			name := fmt.Sprintf("%s-%d", name, i)
			fmt.Printf("%s begin to work...\n", name)

			for v := range in {
				r, ok := f(v)
				if ok {
					out <- r
				}
			}
			close(out)
			fmt.Printf("%s work done\n", name)
		}(i)

		outSlice = append(outSlice, out)
	}

	groupOut := make(chan int)

	// 读取outSlice 切片中各通道数据
	// 扇入到 groupOut 通道
	go func() {
		var wg sync.WaitGroup
		for _, out := range outSlice {
			wg.Add(1)
			go func(out <-chan int) {
				defer wg.Done()
				for v := range out {
					groupOut <- v
				}
			}(out)
		}
		wg.Wait()
		close(groupOut)
	}()
	return groupOut
}

// -- goroutine 超时与取消模式
// 经常遇到的一种场景是: 需要再go应用中向服务端发起多个请求 并保存应答结果
// 以下以向气象中心发起三个查询请求为例
type result struct {
	value string
}

func first(servers ...*httptest.Server) (result, error) {
	c := make(chan result, len(servers))
	queryFunc := func(server *httptest.Server) {
		url := server.URL
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("http get error: %s\n", err)
			return
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		c <- result{string(body)}
	}
	for _, server := range servers {
		go queryFunc(server)
	}
	// 返回第一个调用成功的请求
	return <-c, nil
}

func fakeWeatherServer(name string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s received a weather request\n", name)
		time.Sleep(1 * time.Second)
		_, _ = w.Write([]byte(name + ":ok"))
	}))
}

// 超时之后取消正在查询的请求
func first2(servers ...*httptest.Server) (result, error) {
	c := make(chan result)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	queryFunc := func(i int, server *httptest.Server) {
		url := server.URL
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Printf("query goroutine-%d:http NewRequest error:%s\n", i, err)
			return
		}
		req := request.WithContext(ctx)

		log.Printf("query goroutine-%d: send request...\n", i)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("query goroutine-%d: get return error:%s\n", i, err)
			return
		}
		body, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		log.Printf("query goroutine-%d: get response -> %v", i, string(body))
		c <- result{string(body)}
		return
	}

	for i, serv := range servers {
		go queryFunc(i, serv)
	}

	select {
	case r := <-c:
		return r, nil
	case <-time.After(500 * time.Millisecond):
		return result{}, errors.New(" timeout")
	}
}

func fakeWeatherServerWithTimeCost(name string, costTm int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s received a weather request\n", name)
		time.Sleep(time.Duration(costTm) * time.Millisecond)
		_, _ = w.Write([]byte(name + ":ok"))
	}))
}

// ** 无缓冲channel使用
// -- 用于通过一组goroutine 开始工作
func worker4(i int) {
	fmt.Printf("worker-%d: is working...\n", i)
	time.Sleep(1 * time.Second)
	fmt.Printf("worker-%d: works done\n", i)
}

type signal struct{}

func spawnGroup4(f func(i int), num int, groupSignal <-chan signal) <-chan signal {
	c := make(chan signal)
	var wg sync.WaitGroup
	for i := 0; i < num; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-groupSignal
			fmt.Printf("worker %d: start to work...\n", i)
			f(i)
		}(i)
	}
	go func() {
		wg.Wait()
		c <- signal{}
	}()
	return c
}

// -- 使用无缓冲channel 实现锁的功能
type counter struct {
	c chan int
	i int
}

var cnt counter

func initCounter() {
	cnt = counter{make(chan int), 0}

	go func() {
		for {
			cnt.i++
			cnt.c <- cnt.i
		}
	}()
	fmt.Println("counter init ok!")
}

func Increase() int {
	return <-cnt.c
}

// -- 使用带缓冲channel 实现信号量
func useChannelAsSemaphore() {
	active := make(chan struct{}, 3)
	jobs := make(chan int, 10)

	go func() {
		for i := 0; i < 10; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	wg := sync.WaitGroup{}
	for j := range jobs {
		wg.Add(1)
		go func(k int) {
			defer wg.Done()
			// active channel 最多写入3个struct{}
			active <- struct{}{}
			log.Printf("handler job: %d\n", j)
			time.Sleep(1 * time.Second)
			// active channel 读取一个相当于释放一个semaphore
			<-active
		}(j)
	}

	time.Sleep(5 * time.Second)
}

// -- 采用乐观锁方式访问 带缓冲channel 定义的semaphore
func tryReceive(c <-chan int) (int, bool) {
	select {
	case i := <-c:
		return i, true
	default:
		return 0, false
	}
}
func trySend(c chan<- int, i int) bool {
	select {
	case c <- i:
		return true
	default:
		return false
	}
}

func producer(c chan<- int) {
	var i int = 1
	for {
		time.Sleep(2 * time.Second)
		ok := trySend(c, i)
		if ok {
			fmt.Printf("[producer]: send[%d] to channel\n", i)
			i++
			if i > 5 {
				fmt.Printf("[producer]: exit.\n")
				return
			}
			continue
		}
		fmt.Printf("[producer]: try send [%d],but channel is full\n", i)

	}
}

func consumer(c <-chan int) {
	for {
		i, ok := tryReceive(c)
		if !ok {
			fmt.Println("[consumer]: try to receive from channel, but the channel is empty")
			time.Sleep(1 * time.Second)
			continue
		}
		fmt.Printf("[consumer]: receive [%d] from channel\n", i)

		if i >= 5 {
			fmt.Println("[consumer]: exit..")
			return
		}
	}
}

// -- channel close 之后，如果出现在select case语句中 将每次都会返回0(chan int)
// 但是对于 nil channel ,则会在读取或者写入时 一直阻塞

func nilChannelUsage() {
	c1, c2 := make(chan int), make(chan int)

	go func() {
		time.Sleep(3 * time.Second)
		c1 <- 3
		close(c1)
	}()

	go func() {
		time.Sleep(5 * time.Second)
		c2 <- 5
		close(c2)
	}()

	for {
		select {
		case i, ok := <-c1:
			if !ok {
				c1 = nil
			} else {
				fmt.Println(i)
			}
		case i, ok := <-c2:
			if !ok {
				c2 = nil
			} else {
				fmt.Println(i)
			}
		}
		if c1 == nil && c2 == nil {
			break
		}
	}
	fmt.Println("program end")
}
