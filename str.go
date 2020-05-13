package crack

import "fmt"

func str(str string, max int) {
	m := map[byte]bool{}
	for _, b := range []byte(str) {
		m[b] = true
	}
	bs := make([]byte, len(m))
	i := 0
	for k := range m {
		bs[i] = k
		i++
	}
	for i := 1; i <= max; i++ {

	}
	fmt.Println(bs)
}
