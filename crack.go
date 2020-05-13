package crack

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/nwaples/rardecode"
)

// Crack 破解
func Crack(ctx context.Context, rar string, in <-chan string, out chan<- string) {
	rarfile, err := os.Open("a.rar")
	if err != nil {
		out <- ""
		return
	}
	defer rarfile.Close()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("协程终止")
			return
		case pass := <-in:
			rarfile.Seek(0, 0)
			if rdr, e := rardecode.NewReader(rarfile, pass); e == nil {
				ok := true
				for {
					header, err := rdr.Next()
					if err == io.EOF {
						break
					} else if err != nil {
						ok = false
						break
					}
					if header.IsDir {
						continue
					}
					_, err = io.Copy(bytes.NewBufferString(""), rdr)
					if err != nil {
						fmt.Println(err)
						ok = false
					}
					break
				}
				if ok {
					out <- pass
				}
			} else {
				fmt.Println(e)
			}
		}
	}
}
