package crack

import (
	"fmt"
	"time"

	"github.com/xuender/oil/array"
)

// Pass 密码生成
func Pass(pass chan string, out chan string) {
	old := ""
	num := 0
	go func() {
		for {
			time.Sleep(time.Second * 1)
			fmt.Printf("Probing: '%s' [%d pwds/sec]\n", old, num)
			num = 0
		}
	}()

	// http://www.lingocn.com/
	bs := []byte("htp:/w.lingocn.m3")
	for l := 1; l < 20; l++ {
		p, _ := array.NewProduct(len(bs), l)
		for p.Next() {
			is := p.Value()
			ps := make([]byte, l)
			for i, b := range is {
				ps[i] = bs[b]
			}
			old = string(ps)
			// fmt.Println(old)
			pass <- old
			num++
		}
	}
	out <- ""
}
