// messing with 读取excel，写入excel
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
)

func main() {
	do()
	fmt.Println("vim-go")
}

func do() {
	var srcDir, dstName string
	flag.StringVar(&srcDir, "srcDir", "", "srcDir")
	flag.StringVar(&dstName, "dstName", "", "dstName")
	flag.Parse()
	if srcDir == "" {
		log.Fatal("文件源目录不存在")
		return
	}
	if dstName == "" {
		log.Fatal("结果文件参数不存在")
		return
	}
	err := readPath(srcDir, "sheet1", dstName)
	fmt.Println(err)
	fmt.Println(errors.Unwrap(err))
}

// 读取文件
func readPath(dirName, sheetName, dstFile string) error {
	files, err := os.ReadDir(dirName)
	if err != nil {
		return errors.WithMessage(err, "无法打开文件夹")
	}
	fileName := make([]string, len(files))

	for _, f := range files {
		if !strings.HasPrefix(f.Name(), ".") {
			fileName = append(fileName, dirName+"/"+f.Name())
		}
	}

	lineNum := 1
	fileIndex := 0
	f := excelize.NewFile()
	for _, r := range fileName {
		if r != "" {
			readData, err := readExcel(r, fileIndex)
			if err != nil {
				return err
			}
			if readData != nil {
				for _, input := range readData {
					row := make([]interface{}, 0)
					for _, cellValue := range input {
						row = append(row, cellValue)
					}
					writeExcel(f, sheetName, lineNum, row)
					lineNum++
				}
			}

		}
	}
	// 保存工作簿
	if err := f.SaveAs(dstFile); err != nil {
		return err
	}
	return nil
}

func readExcel(fileName string, fileIndex int) (data [][]string, err error) {
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		err = errors.WithMessage(err, fmt.Sprintf("%s:读取excel文件错误", fileName))
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			err = errors.WithMessage(err, fmt.Sprintf("%s:关闭文件错误", fileName))
			return
		}
	}()
	file, err := excelize.OpenReader(f)

	for _, sheetName := range file.GetSheetList() {
		rows, errs := file.GetRows(sheetName)
		if errs != nil {
			err = errors.WithMessage(errs, fmt.Sprintf("%s:读取sheet内容错误", fileName))
			return
		}
		for key, row := range rows {
			if fileIndex != 0 && key == 0 {
				continue
			}
			data = append(data, row)
		}
	}
	return
}

func writeExcel(f *excelize.File, sheet string, lineNum int, data []interface{}) (err error) {
	line := fmt.Sprintf("A%d", lineNum)
	err = f.SetSheetRow(sheet, line, &data)
	if err != nil {
		return
	}
	return
}

func setRowData(f *excelize.File, sheetName string, rowIndex int, data []string) (err error) {
	for colIndex, cellValue := range data {
		colLetter, errs := excelize.ColumnNumberToName(colIndex + 1)
		if errs != nil {
			err = errs
			return
		}
		cellName := fmt.Sprintf("%s%d", colLetter, rowIndex)
		err = f.SetCellValue(sheetName, cellName, cellValue)
		if err != nil {
			return
		}
	}
	return
}
