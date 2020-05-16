package gocrack

import (
	"fmt"
	"os"
)

// Check 文件检查
func Check(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("the specified file (%s) is not exists or you don't have a right permissions", path)
	}
	return err
	// TODO rar file check
}
