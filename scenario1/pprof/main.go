package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof" // pprof is auto-registered on /debug/pprof
	"runtime"
	"time"
)

func cpuIntensiveWork() {
	for i := 0; i < 1e7; i++ {
		_ = i * i
	}
}

func memoryIntensiveWork() {
	data := make([][]byte, 0)
	for i := 0; i < 50; i++ {
		// Allocate ~10MB each time
		b := make([]byte, 10*1024*1024)
		data = append(data, b)
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("Memory-intensive work done, holding data in memory...")
	_ = data // keep reference so it’s not GC’d
}

func goroutineLeak() {
	ch := make(chan int)
	// Goroutines blocked forever (leak simulation)
	for i := 0; i < 10; i++ {
		go func() {
			for {
				ch <- 1
			}
		}()
	}
}

func main() {
	// Expose pprof at :6060/debug/pprof
	go func() {
		fmt.Println("pprof available at http://localhost:6060/debug/pprof/")
		http.ListenAndServe(":6060", nil)
	}()

	// Simulate workloads
	go goroutineLeak()
	go memoryIntensiveWork()

	for {
		cpuIntensiveWork()
		runtime.Gosched()
	}
}
