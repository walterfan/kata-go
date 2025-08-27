package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

var aString string

// ---------------- 错误用法 ----------------
func unsafeSwapDemo() {
	fmt.Println("=== Unsafe atomic.SwapPointer Demo ===")

	aString = "hello"

	var wg sync.WaitGroup
	wg.Add(2)

	// Writer goroutine
	go func() {
		defer wg.Done()
		words := []string{"aaa", "bbbbbbbbbbbb", "cccccccccccccccccccc"}
		for i := 0; i < 1e6; i++ {
			newStr := words[i%len(words)]
			atomic.SwapPointer(
				(*unsafe.Pointer)(unsafe.Pointer(&aString)),
				unsafe.Pointer(&newStr),
			)
		}
	}()

	// Reader goroutine
	go func() {
		defer wg.Done()
		for i := 0; i < 1e6; i++ {
			s := aString // 这里可能读到乱的 string
			if len(s) > 0 {
				_ = s[0] // 可能 panic: index out of range
			}
		}
	}()

	wg.Wait()
	fmt.Println("Unsafe demo finished (可能没报错，但有数据竞争风险)")
}

// ---------------- 正确用法 ----------------
func safeAtomicValueDemo() {
	fmt.Println("=== Safe atomic.Value Demo ===")

	var v atomic.Value
	v.Store("hello")

	var wg sync.WaitGroup
	wg.Add(2)

	// Writer
	go func() {
		defer wg.Done()
		words := []string{"aaa", "bbbbbbbbbbbb", "cccccccccccccccccccc"}
		for i := 0; i < 1e6; i++ {
			v.Store(words[i%len(words)])
		}
	}()

	// Reader
	go func() {
		defer wg.Done()
		for i := 0; i < 1e6; i++ {
			s := v.Load().(string) // 完整 string 原子替换
			if len(s) > 0 {
				_ = s[0] // 永远安全
			}
		}
	}()

	wg.Wait()
	fmt.Println("Safe demo finished (一定安全)")
}

func main() {
	unsafeSwapDemo()
	time.Sleep(500 * time.Millisecond)
	safeAtomicValueDemo()
}
