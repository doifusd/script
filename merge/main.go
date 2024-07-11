// messing with 读取excel，写入excel
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
)

func main() {
	do()
	fmt.Println("vim-go")
}

func do() {
	err := readPath()
	fmt.Println(err)
	fmt.Println(errors.Unwrap(err))
}

// 读取文件
func readPath() error {
	var dirName string
	flag.StringVar(&dirName, "dirName", "", "dirName")
	flag.Parse()
	if dirName == "" {
		return errors.New("请输入文件目录")
	}
	files, err := os.ReadDir(dirName)
	if err != nil {
		return errors.WithMessage(err, "无法打开文件夹")
	}
	fileName := make([]string, len(files))

	for _, f := range files {
		fileName = append(fileName, dirName+"/"+f.Name())
	}

	fileContent := make(chan []string)
	errReadCh := make(chan error, 1)
	for _, r := range fileName {
		go readExcel(r, fileContent, errReadCh)
	}
	for {
		select {
		case errData, ok := <-errReadCh:
			if ok {
				fmt.Errorf("%s", errData)
			}
		case fdata, ok := <-fileContent:
			if ok {
				go writeExcel("Sheet1", "./a.xlsx", fdata)

			}
		}
	}
	// 获取文件读取文件内容
	return nil
}

func readExcel(fileName string, fileContent chan []string, errs chan error) {
	// func readExcel(fileName string, errs error) {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		errs <- errors.WithMessage(err, fmt.Sprintf("%s:读取excel文件错误", fileName))
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			errs <- errors.WithMessage(err, fmt.Sprintf("%s:关闭文件错误", fileName))
			return
		}
	}()
	//todo 查看有哪些sheet
	//todo 遍历sheet获取内容
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		errs <- errors.WithMessage(err, fmt.Sprintf("%s:读取sheet内容错误", fileName))
		return
	}
	for _, val := range rows {
		fileContent <- val
	}
	return
}

func writeExcel(sheet, distFile string, data []string) {
	f := excelize.NewFile()
	rows, err := f.GetRows(sheet)
	if err != nil {
		fmt.Println("err:", err)
	}
	line := fmt.Sprintf("A%d", len(rows))
	err = f.SetSheetRow(sheet, line, &data)
	fmt.Println("err:", err)
	// 保存工作簿
	if err := f.SaveAs(distFile); err != nil {
		fmt.Println(err)
	}
}
