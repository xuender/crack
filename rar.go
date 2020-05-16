package gocrack

import (
	"bytes"
	"io"
	"os"

	"github.com/nwaples/rardecode"
)

func rar(rarfile *os.File, pass string) bool {
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
				// fmt.Println(err)
				ok = false
			}
			break
		}
		return ok
	}
	return false
}
