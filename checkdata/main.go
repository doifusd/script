package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/extrame/xls"
	"github.com/xuri/excelize/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type queue struct {
	sync.WaitGroup
	count chan int
}

type InvoiceData struct {
	InvoiceCode string `json:"invoice_code"`
	InvoiceNo   string `json:"invoice_no"`
}

var db *gorm.DB

func init() {
	var err error
	dsn := "zhangzj28:Y0puLj7qUT3gI6FW@tcp(rm-2ze50gp4067jmfh1r.mysql.rds.aliyuncs.com:3306)/lecoo_service?charset=utf8mb4&parseTime=True&loc=Local"
	//dsn_dev := "p_lecoo_service:BkWsDQ72Fr@tcp(rm-2ze6qp1n47r2vrx0v.mysql.rds.aliyuncs.com:3306)/lecoo_service?charset=utf8mb4&parseTime=True&loc=Local"
	//dsn := "root:root@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Silent),
		QueryFields: true,
	})

	//db.Logger = logger.Default.LogMode(logger.Silent)
	if err != nil {
		fmt.Println("err:", err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(64)
}

func main() {
	readLog()
	//var totalCnt int64
	////result := db.Debug().Table("server_invoice").Where("invoice_code=?", "3700164320").Where("invoice_no=?", "47795657").Count(&totalCnt).Error
	//result := db.Table("server_invoice").Where("invoice_code=?", "3700164320").Where("invoice_no=?", "47795657").Count(&totalCnt).Error
	//if result != nil {
	//	fmt.Println("err:", result)
	//}
	//fmt.Println("totalCnt:", totalCnt)
}

//var notFoundSeller = map[string]bool{
//	"91420102MA49C2KK43": true,
//	"91310114MA1GUXG096": true,
//	"91510100MA6CDYPG2C": true,
//	"91440300MA5FLW0U7X": true,
//	"91370613MA94GECF1B": true,
//	"91310110MA1G9CTXXT": true,
//	"91320192MA22GLRY71": true,
//	"91370103MA3URP195N": true,
//	"91370602MA3U2M4E08": true,
//	"91610102MAB0REEYX8": true,
//	"91440101MA9Y0CGM9X": true,
//	"91430124MA4T50W1X8": true,
//	"91370602MA3WAHW9XQ": true,
//	"91350582MA8U99LH5J": true,
//}

// 读取excel文件
func readLog() {
	//读取目录下文件
	//判断获取的文件日期
	startTime := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println("run start: ", startTime)
	//读取文件夹下excel
	//将数据库中不存在的数据存储到日志中
	pathName := "/Users/sky/Documents/script/checkdata/data/data22/"
	files, err := os.ReadDir(pathName)
	if err != nil {
		fmt.Println("file-err:", err)
		return
	}
	//
	////ReadCsv(fileName)
	for _, file := range files {
		fileName := pathName + "" + file.Name()
		ReadTxt(fileName)
		//ReadXls(fileName)cd ..
		//ReadXlsx(fileName)
		//ReadXls(fileName)
	}
	//ReadTxt("/Users/sky/Documents/script/checkdata/data/data22/30.txt")

	////将数据库中不存在的数据存储到日志中，将原始日志中不存在单独存储
	//reader := bufio.NewReader(file)
	//i := 0
	//// var wg sync.WaitGroup
	//
	//waitNum := make(chan int, 10)
	//for {
	//	conent, errs := reader.ReadString('\n')
	//	if errs != nil {
	//		fmt.Println("err:", errs)
	//		//log.Fatal("read err:", errs)
	//		//等于eof 退出
	//		//errs != io.EOF
	//		break
	//	}
	//	res := strings.ReplaceAll(conent, "\\\\u", "\\u")
	//	str := strings.ReplaceAll(res, "\\\"", "\"")
	//
	//	str2, errs2 := zh2Unicode([]byte(str))
	//	if errs2 != nil {
	//		fmt.Println("to chinses err:", errs2)
	//	}
	//	isJSON := json.Valid([]byte(str2))
	//	if isJSON {
	//		waitNum <- 1
	//		//todo go routine
	//		// res := httpClient(str2)
	//		// wg.Add(1)
	//		go httpClient(waitNum, str2, i)
	//		// wg.Wait()
	//		time.Sleep(time.Microsecond * 5000)
	//	}
	//	i++
	//}
	//stopTime := time.Now().Format("2006-01-02 15:04:05")
	//fmt.Println("run complete: ", stopTime)
}

func ReadTxt(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("open file failed: %s \n", err.Error())
		return
	}
	defer file.Close()
	line := bufio.NewReader(file)
	waitNum := make(chan int, 10)
	i := 0
	for {
		content, _, err := line.ReadLine()
		if err == io.EOF {
			break
		}
		contentStr := string(content)
		//空格分割
		contentArr := strings.Fields(contentStr)
		//if notFoundSeller[contentArr[2]] == true {
		//	continue
		//}
		//fmt.Println("contentArr:", contentArr)
		//go checkData(contentArr[1], contentArr[0], contentArr[2], waitNum)
		go checkData(contentArr[1], contentArr[0], contentArr[2], waitNum)
		i++
		//fmt.Println("read line:", i)
	}
}

func ReadCsv(file_path string) (res [][]string) {
	file, err := os.Open(file_path)
	if err != nil {
		fmt.Errorf("open_err:", err)
		return
	}
	defer file.Close()
	// 初始化csv-reader
	reader := csv.NewReader(file)
	// 设置返回记录中每行数据期望的字段数，-1 表示返回所有字段
	reader.FieldsPerRecord = -1
	// 允许懒引号（忘记遇到哪个问题才加的这行）
	reader.LazyQuotes = true
	// 返回csv中的所有内容
	record, read_err := reader.ReadAll()
	if read_err != nil {
		fmt.Errorf("read_err:", read_err)
		return
	}
	waitNum := make(chan int, 10)
	for _, value := range record {
		//strArr := strings.Split(value[0], " ")
		waitNum <- 1
		//go checkData(value[2], value[1], waitNum)
		go updateData(value[2], value[1], value[3], waitNum)
		//go checkData(strArr[0], strArr[1], waitNum)
	}
	return record
}

func ReadXls(file_path string) (res [][]string) {
	if xlFile, err := xls.Open(file_path, "utf-8"); err == nil {
		//第一个sheet
		sheet := xlFile.GetSheet(0)
		if sheet.MaxRow != 0 {
			//temp := make([][]string, sheet.MaxRow)
			for i := 0; i < int(sheet.MaxRow); i++ {
				row := sheet.Row(i)
				//data := make([]string, 0)
				if row.LastCol() > 0 {
					//for j := 0; j < row.LastCol(); j++ {
					//	col := row.Col(j)
					//	data = append(data, col)
					//}
					//temp[i] = data
					fmt.Println("data:", row.Col(2), row.Col(3))
				}
			}
			//res = append(res, temp...)
		}
	} else {
		fmt.Println("open_err:", err)
	}
	return res
}

func ReadXlsx(file_path string) (res [][]string) {
	xlFile, err := excelize.OpenFile(file_path)
	rows, err := xlFile.GetRows("Sheet1")
	if err != nil {
		return
	}
	fmt.Println("rows:", rows)
	//if xlFile, err := excelize.OpenFile(file_path); err == nil {
	//	for index, sheet := range xlFile.Sheets {
	//		//第一个sheet
	//		if index == 0 {
	//			temp := make([][]string, len(sheet.Rows))
	//			for k, row := range sheet.Rows {
	//				var data []string
	//				for _, cell := range row.Cells {
	//					data = append(data, cell.Value)
	//				}
	//				temp[k] = data
	//			}
	//			res = append(res, temp...)
	//		}
	//	}
	//} else {
	//	fmt.Println("open_err:", err)
	//}
	return res
}

func zh2Unicode(raw []byte) (str string, err error) {
	str, err = strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	if err != nil {
		return
	}
	return
}

func checkData(invoiceCode, invoiceNo, sellerTaxNo string, waitNum chan int) {
	defer func(waitNum chan int) {
		<-waitNum
	}(waitNum)
	//var invoiceRes InvoiceData
	//result := db.Debug().Table("server_invoice").Where("invoice_code=?", invoiceCode).Where("invoice_no=?", invoiceNo).Order("id desc").First(&invoiceRes)
	//fmt.Println("result.RowsAffected:", result.RowsAffected)
	//if result.Error != nil {
	//	fmt.Println("err:", result)
	//}
	//if result.RowsAffected != 0 {
	//	fmt.Printf("%s %s %s\n", invoiceCode, invoiceNo, sellerTaxNo)
	//}
	var totalCnt int64
	result := db.Table("server_invoice").Where("invoice_code=?", invoiceCode).Where("invoice_no=?", invoiceNo).Count(&totalCnt).Error
	if result != nil {
		fmt.Println("err:", result)
	}
	//fmt.Println("totalCnt:", totalCnt)
	//if totalCnt != 0 {
	if totalCnt == 0 {
		fmt.Printf("%s %s %s\n", invoiceCode, invoiceNo, sellerTaxNo)
	}
}

func updateData(invoiceCode, invoiceNo, SellerTaxNo string, waitNum chan int) {
	defer func(waitNum chan int) {
		<-waitNum
	}(waitNum)
	var cnt int64
	//result := db.Table("server_invoice").Where("invoice_code=?", invoiceCode).Where("invoice_no=?", invoiceNo).Order("id desc").First(&invoiceRes)
	result := db.Table("invoice_result").Where("invoice_code=?", invoiceCode).Where("invoice_no=?", invoiceNo).Count(&cnt)
	if result.Error != nil {
		fmt.Println("err:", result.Error)
	}
	if cnt > 0 {
		res := db.Table("invoice_result").Where("invoice_code=?", invoiceCode).Where("invoice_no=?", invoiceNo).Update("seller_tax_no", SellerTaxNo)
		if res.Error != nil {
			fmt.Println("update err:", res.Error)
		}
	}
}
