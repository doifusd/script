package main

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func main() {
	var fileNameFull, fileNameSub string
	flag.StringVar(&fileNameFull, "fileNameFull", "", "fileNameFull")
	flag.StringVar(&fileNameSub, "fileNameSub", "", "fileNameSub")
	flag.Parse()
	compare(fileNameFull, fileNameSub)
}

func compare(fileNameFull, fileNameSub string) {

	rowsFull, err := findFull(fileNameFull)
	if err != nil {
		fmt.Println(fileNameFull, "主文件读取错误")
	}
	rowsSub, err := findSub(fileNameSub)
	if err != nil {
		fmt.Println(fileNameSub, "子文件读取错误")
	}
	sum := 0.00
	for _, row1 := range rowsFull {
		for _, row2 := range rowsSub {
			if row1[0] == row2[0] {
				fmt.Println("find_res:", row1[1])
				tmp, _ := strconv.ParseFloat(row1[1], 64)
				sum += tmp
			}
		}
	}
	fmt.Println("res:", sum)
}

func findFull(fileName string) (rows [][]string, err error) {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	rows, err = f.GetRows("Sheet1")
	if err != nil {
		return
	}
	return
}

func findSub(fileName string) (rows [][]string, err error) {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	rows, err = f.GetRows("Sheet1")
	if err != nil {
		return
	}
	return
}
