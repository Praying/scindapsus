package util

import (
	"fmt"
	"testing"
)

func TestRandStringBytesMaskImprSrc(t *testing.T) {
	fmt.Println(RandStringBytesMaskImprSrc(24))
}

func TestOther(t *testing.T) {
	println("start main")
	ch := make(chan int)
	var result int
	go func() {
		println("come into goroutine1")
		var r int
		for i := 1; i <= 10; i++ {
			r += i
		}
		ch <- r
	}()
	go func() {
		println("come into goroutine2")
		var r int = 1
		for i := 1; i <= 10; i++ {
			r *= i
		}
		ch <- r
	}()
	go func() {
		println("come into goroutine3")
		ch <- 11
	}()
	for i := 0; i < 3; i++ {
		result += <-ch
	}
	close(ch)
	println("result is:", result)
	println("end main")
}
