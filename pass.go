package gocrack

import (
	"fmt"
	"time"

	"github.com/xuender/oil/array"
)

// Pass 密码生成
func Pass(pass chan string, out chan string, bs []byte, current string) {
	old := ""
	num := 0
	go func() {
		for {
			time.Sleep(time.Second * 1)
			fmt.Printf("Probing: '%s' [%d pwds/sec]\n", old, num)
			num = 0
		}
	}()

	for l := 1; l < 20; l++ {
		p, _ := array.NewProduct(len(bs), l)
		for p.Next() {
			is := p.Value()
			ps := make([]byte, l)
			for i, b := range is {
				ps[i] = bs[b]
			}
			old = string(ps)
			if current != "" {
				if current == old {
					current = ""
				}
				continue
			}
			// fmt.Println(old)
			pass <- old
			num++
		}
	}
	out <- ""
}
