package main

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestAirportSecurityCheckChannel(t *testing.T) {
	passengers := 150
	queue := make(chan struct{}, 150)

	AirportSecurityCheckChannel("channel-1", queue)
	AirportSecurityCheckChannel("channel-2", queue)
	AirportSecurityCheckChannel("channel-3", queue)

	time.Sleep(2 * time.Second)
	for i := 0; i < passengers; i++ {
		queue <- struct{}{}
	}
	time.Sleep(5 * time.Second)

	close(queue)
	time.Sleep(30 * time.Second)
}

// 需要再环境变量中配置 GODEBUG=schedtrace=1000
func TestRoutineScheduler(t *testing.T) {
	runtime.GOMAXPROCS(1)
	go deadLoop()
	for {
		time.Sleep(1 * time.Second)
		fmt.Println("main goroutine got scheduled!")
	}
}

func TestSpawnGroup(t *testing.T) {
	done := spawnGroup(5, worker, 2)
	println("spawn a group of workers...")
	<-done
	println("group worker done...")
}

func TestSpawnGroupWithOvertime(t *testing.T) {

	done := spawnGroup(5, worker, 10)
	println("spawn a group of workers...")

	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	select {
	case v := <-timer.C:
		fmt.Printf("timer returns value:%v\n", v)
		println("wait group workers timeout!")
	case <-done:
		println("group worker done!")
	}
}

func TestSpawnGroupNotifyMultiGoroutine(t *testing.T) {
	quit := spawnGroup2(5, worker, 3)
	println("spawn a group of workers...")
	time.Sleep(5 * time.Second)

	println("notify all worker to exit!")
	quit <- struct{}{}

	timer := time.NewTimer(time.Second * 5)
	defer timer.Stop()
	select {
	case <-timer.C:
		println("wait group workers timeout!")
	case <-quit:
		println("group worker done!")
	}
}

func TestConcurrentShutdown(t *testing.T) {
	d1 := 1 * time.Second
	d2 := 5 * time.Second
	f1 := shutdownMaker(d1)
	f2 := shutdownMaker(d2)

	err1 := concurrentShutdown(3*time.Second, []time.Duration{d1, d2}, ShutdownFunc(f1), ShutdownFunc(f2))
	if err1 == nil {
		t.Errorf("want timeout, actual nil")
		return
	}

	err := concurrentShutdown(6*time.Second, []time.Duration{d1, d2}, ShutdownFunc(f1), ShutdownFunc(f2))
	if err != nil {
		t.Errorf("want nil, actual %v", err)
		return
	}

}

func TestPipelineLikedChannel(t *testing.T) {
	in := newNumGenerator(1, 10)
	// spawn2(filterOdd, in) 产生的读通道,作为外层spawn2的入参
	// 从而实现了unix/linux 中的管道效果
	out := spawn2(square, spawn2(filterOdd, in))

	for v := range out {
		fmt.Println(v)
	}
}

// 测试管道模式的扩展模式: 扇出模式/扇入模式
func TestSpawnGroup3(t *testing.T) {
	in := newNumGenerator(1, 20)
	out := spawnGroup3("square", 2, square, spawnGroup3("filterOdd", 3, filterOdd, in))
	for v := range out {
		fmt.Println(v)
	}
}

func TestFakeWeatherRequest(t *testing.T) {
	r, err := first(fakeWeatherServer("server-1"),
		fakeWeatherServer("server-2"),
		fakeWeatherServer("server-3"))
	if err != nil {
		t.Errorf("invoke first error: %v", err)
		return
	}
	fmt.Println(r)
}

func TestFakeWeatherRequestWithTimeout(t *testing.T) {
	r, err := first2(fakeWeatherServerWithTimeCost("server-1", 200),
		fakeWeatherServerWithTimeCost("server-2", 600),
		fakeWeatherServerWithTimeCost("server-3", 1000))
	if err != nil {
		t.Errorf("invoke first2 error: %v", err)
		return
	}
	fmt.Println(r)
}

// 测试select case 语句执行方式
// -- 除非default 之外所有的case 都未就绪，否则不会选择default case，即使把default 语句放在最上面也是这样
func TestSelectCase(t *testing.T) {
	c := make(chan int, 12)

	for i := 0; i < 12; i++ {
		select {
		default:
			fmt.Println("default case selected..")
		case c <- i:
			fmt.Println("write data to channel selected..")
		case <-c:
			fmt.Println("read channel case selected")
		}
	}
}

func TestCancel(t *testing.T) {
	// 创建一个可取消的 context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // 确保在函数返回时取消

	// 启动一个 goroutine 在 2 秒后调用 cancel
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("Calling cancel()")
		cancel()
	}()

	// 使用 select 等待 ctx.Done() 或 5 秒超时
	select {
	case <-ctx.Done():
		fmt.Println("Received cancellation signal")
	case <-time.After(1 * time.Second):
		fmt.Println("Timed out")
	}
}

// 使用groupSignal 信号作为栓子 控制所有worker同时开始工作
// 各个goroutine 中的worker 启动时间可能不一致,但是都会一直阻塞在 groupSignal处，直到main goroutine关闭 groupSignal
func TestUseChannelAsLatch(t *testing.T) {
	fmt.Println("start a group of workers...")
	groupSignal := make(chan signal)
	c := spawnGroup4(worker4, 5, groupSignal)
	time.Sleep(5 * time.Second)
	fmt.Println("the group of workers start to work..")
	close(groupSignal)
	<-c
	fmt.Println("the group of workers work done!")
}

// -- 无缓冲channel 用作计数器
func TestUseChannelAsLock(t *testing.T) {
	initCounter()

	for i := 0; i < 10; i++ {
		go func(i int) {
			v := Increase()
			fmt.Printf("goroutine-%d current counter value is：%d\n", i, v)
		}(i)
	}
	time.Sleep(5 * time.Second)
	fmt.Printf("counter final result is: %d\n", cnt.i-1)
}

func TestUseChannelAsSemaphore(t *testing.T) {
	useChannelAsSemaphore()
}

func TestUsaChannelAsOptimisticLock(t *testing.T) {
	c := make(chan int, 3)
	go producer(c)
	go consumer(c)
	select {}
}

func TestNilChannelUsage(t *testing.T) {
	nilChannelUsage()
}

func TestSelectHeartBeat(t *testing.T) {
	// 每3秒钟写入一次 timer 通道
	timer := time.NewTicker(3 * time.Second)
	c := make(chan int)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			fmt.Println("process heart beat jobs..")
		case <-c:
			fmt.Println("process business..")
		}
	}

}
