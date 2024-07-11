package main

import (
	"fmt"
	"os"

	"github.com/extrame/xls"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
)

func main() {
	//var sourceName, distName string
	//flag.StringVar(&sourceName, "sourceName", "", "sourceName")
	//flag.StringVar(&distName, "distName", "", "distName")
	//flag.Parse()
	sourceName := "./data/"
	distName := "./output.xls"

	merge(sourceName, distName)
}

func merge(sourceName, distName string) error {
	//读取文件夹下文件
	//将文件内容输出到同一个文件中
	// fileName := make(chan string, 20)
	// content := make(chan []string, 1)
	fileNameArr, err := readDir(sourceName)
	if err != nil {
		return errors.WithMessage(err, "主文件读取错误")
	}
	if len(fileNameArr) == 0 {
		return errors.WithMessage(err, "文件夹为空")
	}
	// var wg sync.WaitGroup
	// wg.Add(len(fileNameArr))
	data := make(chan [][]string)
	for _, v := range fileNameArr {
		// readExcel(wg, v, data)
		//通过扩展名判断文件类型
		if v[len(v)-4:] != ".xls" {
			readExcel(v, data)
		} else {
			go readXls(v, data)
		}
	}
	// wg.Wait()

	for {
		select {
		case content, ok := <-data:
			if ok {
				writeExcel(distName, content)
				//go writeExcel(distName, content)
			}
		}
	}
	//fmt.Println("res:", sum)
	return nil
}

func readDir(dirName string) (fileNameArr []string, err error) {
	files, errs := os.ReadDir(dirName)
	if errs != nil {
		defer func() {
			err = errors.WithMessage(errs, "读取文件夹失败")
		}()
		return
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if file.Name() == "." || file.Name() == ".." {
			continue
		}
		fileNameArr = append(fileNameArr, dirName+file.Name())
	}
	return fileNameArr, nil
}

// func readExcel(wg sync.WaitGroup, fileName string, data chan []string) {
func readExcel(fileName string, data chan [][]string) {
	// defer wg.Done()
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		err = errors.WithMessage(err, "读取文件失败")
		return
	}
	defer func() {
		if err = f.Close(); err != nil {
			err = errors.WithMessage(err, "获取文件关闭失败")
			fmt.Println(err)
		}
	}()
	rows, err := f.GetRows("Sheet1")
	fmt.Println("rows:", rows)
	if err != nil {
		err = errors.WithMessage(err, "获取文件内容失败")
		return
	}

	for _, val := range rows {
		fmt.Println("val:", val)
		// data <- val
		fmt.Println("data:", data)
	}
}

func readXls(fileName string, data chan [][]string) {
	if xlFile, err := xls.Open(fileName, "utf-8"); err == nil {
		//第一个sheet
		sheet := xlFile.GetSheet(0)
		if sheet.MaxRow != 0 {
			temp := make([][]string, sheet.MaxRow)
			for i := 0; i < int(sheet.MaxRow); i++ {
				row := sheet.Row(i)
				datas := make([]string, 0)
				if row.LastCol() > 0 {
					for j := 0; j < row.LastCol(); j++ {
						col := row.Col(j)
						datas = append(datas, col)
					}
					temp[i] = datas
				}
			}
			//res = append(res, temp...)
			data <- temp
		}
	} else {
		fmt.Println("open_err:", err)
	}
	return
}

func writeExcel(fileName string, data [][]string) {
	f := excelize.NewFile()
	// f.Path = "../excel_files/TMP_07.xlsx"
	// 创建一个工作表
	//index, _ := f.NewSheet("Sheet1")
	// 设置单元格的值
	//f.SetCellValue("Sheet1", "B2", 100)
	// 设置工作簿的默认工作表
	//f.SetActiveSheet(index)
	// 根据指定路径保存文件
	//if err := f.SaveAs("Book1.xlsx"); err != nil {
	//	fmt.Println(err)
	//}

	sheetName := f.GetSheetName(f.GetActiveSheetIndex())
	// 生成流写入对象
	streamSheet, err := f.NewStreamWriter(sheetName)
	if err != nil {
		fmt.Println(err)
	}
	// 设置一整行的值   只有值
	// 先构造数据
	// 姓名  年龄  性别  工资
	// 使用faker 模块构造测试数据
	//写个表头
	//结算单号 发票类型 发票号码 发票代码 销方编号(企业) 销方纳税人识别号 销方名称 销方地址电话 销方银行名称账号 购方编号(企业) 购方纳税人识别号 购方名称 购方地址电话 购方银行名称账号 税率 不含税金额 税额 含税金额 发票开票日期 发票来源 备注 发票状态 原发票号码 原发票代码 红字信息表编号 扩展字段1 扩展字段2 扩展字段3 数电发票号码 创建时间 购方银行账号 购方银行名称 购方电话 购方地址 红冲时间 红冲状态 机器编码 开票点名称 开票点代码 开机票号 税控终端码 退回时间 退回备注 退回状态 特殊发票标识 销方银行名称 销方电话 销方地址 销方银行账号 校验码 作废红冲原因 作废时间 密文 预作废状态 发票开具后自动发送邮件的邮箱 收款人姓名 系统来源 复核人姓名 开票人姓名

	titleCont := 0
	// var row = make([]string, len(v))
	for k, v := range data { // 行
		if k == 0 && titleCont > 0 {
			fmt.Println("v:", v)
			fmt.Println("k:", k)
			titleCont++
			//写入第一行标题
		}
		// row[0] = faker.Name()
		// row[1] = rand.Intn(100)
		// row[2] = faker.Gender()
		// row[3] = rand.Intn(10000) / 100

		// streamSheet.SetRow(fmt.Sprintf("A%d", i), row)
	}

	// 执行了 flush 才算是写进去了
	streamSheet.Flush()

	f.Save()
}
