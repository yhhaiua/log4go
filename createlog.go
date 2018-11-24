package log4go

import (
	"fmt"
	"os"
	"strings"
)

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
func createFileLog(fname string)  {
	var end, endL, endR int
	endL = strings.LastIndex(fname, "/")
	endR = strings.LastIndex(fname, "\\")
	if endL > endR {
		end = endL
	} else {
		end = endR
	}
	if end != -1 {
		folder := fname[0:end]
		res, err := pathExists(folder)
		if err != nil {
			fmt.Println("err:", err)
		}
		if !res {
			err := os.MkdirAll(folder, os.ModePerm)
			if err != nil {
				fmt.Println("err:", err)
			} else {
				fmt.Println("create directory success!")
			}
		}
	}
}
