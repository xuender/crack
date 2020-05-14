package main

import (
	"context"
	"fmt"

	"github.com/xuender/crack"
)

// TODO https://www.jianshu.com/p/f9cf46a4de0e

func main() {
	inChan := make(chan string, 1)
	outChan := make(chan string, 1)
	ctx, cancel := context.WithCancel(context.Background())

	for i := 0; i < 4; i++ {
		go crack.Crack(ctx, "a.rar", inChan, outChan)
	}
	go crack.Pass(inChan, outChan)
	select {
	case r := <-outChan:
		if r == "" {
			fmt.Println("密码未找到")
		} else {
			fmt.Printf("密码是:%s\n", r)
		}
		cancel()
	}
}
