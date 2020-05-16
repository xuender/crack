package main

import (
	"context"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/xuender/gocrack"
	"github.com/xuender/oil/str"
)

// TODO https://www.jianshu.com/p/f9cf46a4de0e

var (
	_help bool
	_num  int
	_abc  string
)

func init() {
	flag.BoolVar(&_help, "help", false, "show this screen.")
	flag.IntVar(&_num, "goroutines", runtime.NumCPU(), "you can specify how many goroutines will be run, maximum 20")
	flag.StringVar(&_abc, "abc", "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", "the password may contain letters")
}
func main() {
	head()
	// doCrack("a.rar", 4)
	flag.Parse()
	if _help || len(flag.Args()) < 1 {
		usage()

	} else {
		for _, name := range flag.Args() {
			doCrack(name)
		}
	}
}

func head() {
	fmt.Println("GoCrack! 0.1 by Ender Xu (xuender@gmail.com)")
	fmt.Println()
}

func usage() {
	fmt.Println("Usage:")
	fmt.Println("  gocrack encrypted_archive.ext [-goroutines NUM]")
	fmt.Println()
	fmt.Println("Options:")
	flag.PrintDefaults()
	fmt.Println()
	fmt.Println("Info:")
	// fmt.Println("  This program supports only RAR, ZIP and 7Z encrypted archives.")
	fmt.Println("  This program supports only RAR encrypted archives.")
	fmt.Println("  GoCrack! usually detects the archive type.")
}

func doCrack(path string) {
	if err := gocrack.Check(path); err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}
	current := ""
	if file, err := os.Open(fmt.Sprintf("%s.xml", path)); err == nil {
		defer file.Close()
		if data, err := ioutil.ReadAll(file); err == nil {
			v := gocrack.Config{}
			if err = xml.Unmarshal(data, &v); err == nil {
				_abc = v.ABC
				if v.GoodPassword == "" {
					current = v.Current
				}
				// TODO GoodPassword
				// fmt.Println(v)
			}
		}
	}
	inChan := make(chan string, 1)
	outChan := make(chan string, 1)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	for i := 0; i < _num; i++ {
		go gocrack.Crack(ctx, path, inChan, outChan)
	}
	go gocrack.Pass(inChan, outChan, str.Union(_abc), current)
	select {
	case <-signalChan:
		fmt.Println("----保存")
		cancel()
	case r := <-outChan:
		if r == "" {
			fmt.Println("密码未找到")
		} else {
			fmt.Printf("GOOD: password cracked: '%s'\n", r)
		}
		cancel()
	}
}
