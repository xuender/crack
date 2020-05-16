package gocrack

import (
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/xuender/oil/array"
	"github.com/xuender/oil/str"
)

// Crack 文件破解
type Crack struct {
	XMLName      xml.Name           `xml:"rarcrack"`
	ABC          string             `xml:"abc"`           // 字母
	Current      string             `xml:"current"`       // 当前
	GoodPassword string             `xml:"good_password"` // 密码
	Num          int                `xml:"-"`             // 一秒生成密码次数
	Ctx          context.Context    `xml:"-"`
	Cancel       context.CancelFunc `xml:"-"`
	path         string
	chanPass     chan string
}

// Run 运行
func (c *Crack) Run(num int) {
	if c.GoodPassword != "" {
		rarfile, err := os.Open(c.path)
		if err != nil {
			fmt.Println(err)
			c.Cancel()
			return
		}
		if rar(rarfile, c.GoodPassword) {
			c.Cancel()
			return
		} else {
			c.GoodPassword = ""
		}
	}
	if num > 50 {
		num = 50
	}
	for i := 0; i < num; i++ {
		go c.crack()
	}

	bs := str.Union(c.ABC)
	pass := c.Current != ""
	for l := 1; l < 20; l++ {
		p, _ := array.NewProduct(len(bs), l)
		for p.Next() {
			is := p.Value()
			ps := make([]byte, l)
			for i, b := range is {
				ps[i] = bs[b]
			}
			str := string(ps)
			if pass {
				if str == c.Current {
					pass = false
				}
				continue
			}
			c.chanPass <- str
			c.Current = str
			c.Num++
		}
	}
	c.Cancel()
}

func (c *Crack) crack() {
	rarfile, err := os.Open(c.path)
	if err != nil {
		fmt.Println(err)
		c.Cancel()
		return
	}
	defer rarfile.Close()
	for {
		select {
		case <-c.Ctx.Done():
			return
		case pass := <-c.chanPass:
			if rar(rarfile, pass) {
				c.GoodPassword = pass

				c.Cancel()
				return
			}
		}
	}
}

// Close 关闭
func (c *Crack) Close() {
	output, err := xml.MarshalIndent(*c, "", "  ")
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	bs := append([]byte(xml.Header), output...)
	ioutil.WriteFile(fmt.Sprintf("%s.xml", c.path), bs, 0644)
}

// New 新建
func New(path string, abc string) (*Crack, error) {
	if err := check(path); err != nil {
		return nil, err
	}

	c := Crack{
		path:     path,
		chanPass: make(chan string, 1),
	}
	if file, err := os.Open(fmt.Sprintf("%s.xml", path)); err == nil {
		defer file.Close()
		if data, err := ioutil.ReadAll(file); err == nil {
			xml.Unmarshal(data, &c)
			fmt.Printf("INFO: cracking %s, status file: %s.xml\n", path, path)
			fmt.Printf("INFO: Resuming cracking from password: '%s'\n", c.Current)
		}
	}
	if c.ABC == "" {
		c.ABC = abc
	}
	return &c, nil
}
