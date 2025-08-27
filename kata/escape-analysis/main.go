package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"
)

type Usage struct {
	Lock               sync.RWMutex
	SecretPath         *string
	Alias              *string
	Version            *string
	Count              *int64
	Issuer             *string
	UsageSuccessReport bool
}

func WithIssuer(usage *Usage, data ...interface{}) {
	if usage == nil {
		return
	}

	if len(data) != 1 {
		return
	}

	issuer, ok := data[0].(string)
	if !ok {
		return
	}

	atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&usage.Issuer)), unsafe.Pointer(&issuer))
}

func foo() int {
	x := 10  // 只在函数内部使用
	return x // 值返回 -> 栈上
}

func bar() *int {
	y := 20
	return &y // 返回指针 -> 堆上
}

func baz() func() int {
	z := 30
	return func() int {
		z++ // 被闭包捕获 -> 堆上
		return z
	}
}

func main() {
	fmt.Println("foo:", foo())

	p := bar()
	fmt.Println("bar:", *p)

	f := baz()
	fmt.Println("baz first:", f())
	fmt.Println("baz second:", f())

	usage := &Usage{}
	WithIssuer(usage, "test")
	fmt.Println("usage:", *usage.Issuer)
}
