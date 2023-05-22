package main

import (
	"flag"
	"fmt"

	"github.com/xuri/excelize/v2"

	"github.com/mozillazg/go-pinyin"
)

func main() {
	commandStr()
}

func commandStr() {
	var fileNameFull string
	flag.StringVar(&fileNameFull, "fileNameFull", "", "fileNameFull")
	flag.Parse()
	readFile(fileNameFull)
}

func readFile(fileNameFull string) {
	rowsFull, err := findFull(fileNameFull, "Sheet1")
	if err != nil {
		fmt.Println(fileNameFull, "主文件读取错误")
	}
	for _, row1 := range rowsFull {
		py := Pinyin(row1[0])
		fmt.Printf("%s\t%s\n", row1[0], py)
	}
}

func Pinyin(input string) string {
	py := pinyin.NewArgs()
	res := pinyin.Pinyin(input, py)
	tmp := ""
	for _, val := range res {
		tmp += val[0]
	}
	return tmp
}

func findFull(fileName, sheetName string) (rows [][]string, err error) {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	rows, err = f.GetRows(sheetName)
	if err != nil {
		return
	}
	return
}
