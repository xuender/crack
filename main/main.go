package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/nwaples/rardecode"
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
func main2() {
	if rarfile, err := os.Open("a.rar"); err == nil {
		for _, pass := range []string{"a", "1", "b", "33"} {
			ok := true
			rarfile.Seek(0, 0)
			rdr, e := rardecode.NewReader(rarfile, pass)
			if e != nil {
				ok = false
				fmt.Println(e)
				continue
			}
			for {
				header, err := rdr.Next()
				if err == io.EOF {
					break
				} else if err != nil {
					ok = false
					fmt.Println(err)
					break
				}
				if header.IsDir {
					continue
				}
				_, err = io.Copy(bytes.NewBufferString(""), rdr)
				if err != nil {
					ok = false
					fmt.Println(err)
				}
				break
			}
			if ok {
				fmt.Println("OK")
			}
		}
	} else {
		fmt.Println(err)
	}
}
