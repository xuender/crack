package gocrack

import "encoding/xml"

// Config 配置文件
type Config struct {
	XMLName      xml.Name `xml:"rarcrack"`
	ABC          string   `xml:"abc"`           // 字母
	Current      string   `xml:"current"`       // 当前
	GoodPassword string   `xml:"good_password"` // 密码
}
