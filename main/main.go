package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/xuender/gocrack"
)

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
	flag.Parse()
	if _help || len(flag.Args()) < 1 {
		usage()
	} else {
		if c, err := gocrack.New(flag.Arg(0), _abc); err == nil {
			c.Ctx, c.Cancel = context.WithCancel(context.Background())
			signalChan := make(chan os.Signal, 1)
			signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
			defer c.Close()
			go c.Run(_num)
			go probing(c)
			select {
			case <-c.Ctx.Done():
				if c.GoodPassword != "" {
					fmt.Printf("GOOD: password cracked: '%s'\n", c.GoodPassword)
				}
			case <-signalChan:
			}
		} else {
			fmt.Printf("ERROR: %s\n", err)
		}
	}
}

func probing(c *gocrack.Crack) {
	for {
		time.Sleep(time.Second * 1)
		fmt.Printf("Probing: '%s' [%d pwds/sec]\n", c.Current, c.Num)
		c.Num = 0
	}
}
func head() {
	fmt.Println("GoCrack! 1.0.1 by Ender Xu (xuender@gmail.com)")
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
