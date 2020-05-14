package crack

import (
	"fmt"
	"strconv"
	"time"

	"github.com/xuender/oil/array"
)

// Pass 密码生成
func Pass(pass chan<- string, out chan<- string) {
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
	bs := []byte("htp:/w.lingocn.m")
	array.Combine(len(bs), func(is []int) error {

		return nil
	})
	for i := -1000; i < 1000; i++ {
		p := strconv.FormatInt(int64(i), 10)
		pass <- p
		old = p
		num++
	}
	out <- ""
}
