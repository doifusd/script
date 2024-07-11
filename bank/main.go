package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

func main() {
	readExcel()
}

func readExcel() {
	f, err := excelize.OpenFile("./data/data.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Get value from cell by given worksheet name and cell reference.
	//cell, err := f.GetCellValue("Sheet1", "A2")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(cell)
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, row := range rows {
		//f.GetCellValue("Sheet1", row[0])
		for key, colCell := range row {
			fmt.Print(key, colCell, "\t")

		}
		fmt.Println()
	}
}
