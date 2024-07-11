package main

import (
	"fmt"
	"os"

	"github.com/xuri/excelize/v2"
)

func main() {
	readPath("../data/dataS")
}

// 读取单个excel文件
func Read(fileName, sheetName string) (contents [][]string) {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	contents, err = f.GetRows(sheetName)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

// 获取文件下文件
func readPath(filePath string) {
	files, err := os.ReadDir(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	sheetName := "Sheet1"
	for _, file := range files {
		if file.Name() == ".DS_Store" {
			continue
		}
		fileName := filePath + "/" + file.Name()
		res := Read(fileName, sheetName)
		Write("save1.xlsx", res)
	}
}

func Write(fileName string, contents [][]string) {
	lastRowId := GetWriteRows(fileName, "Sheet1")
	//fmt.Println("lastRowId", lastRowId)
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	streamWriter, err := f.NewStreamWriter("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}
	//将新加contents写进流式写入器
	for rowID := 0; rowID < len(contents); rowID++ {
		row := make([]interface{}, len(contents[0]))
		for colID := 0; colID < len(contents[0]); colID++ {
			row[colID] = contents[rowID][colID]
		}
		cell, _ := excelize.CoordinatesToCellName(1, rowID+lastRowId+1) //决定写入的位置
		if err := streamWriter.SetRow(cell, row); err != nil {
			fmt.Println(err)
		}
	}

	if err := streamWriter.Flush(); err != nil {
		fmt.Println(err)
		return
	}

	if err := f.SaveAs(fileName); err != nil {
		fmt.Println(err)
	}
	return
}

func GetWriteRows(fileName, sheetName string) (total_rows int) {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	contents, err := f.GetRows(sheetName)
	if err != nil {
		fmt.Println(err)
		return
	}
	total_rows = len(contents)
	return
}
